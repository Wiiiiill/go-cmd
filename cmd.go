package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"text/template"
	"unicode"
	"unicode/utf8"
)

// App represents a CLI application with its commands and configuration
type App struct {
	commands       Commands
	usageTemplate  string
	setFlags       func(f *flag.FlagSet)
	exitMu         sync.Mutex
	exitStatus     int
	commandsSorted bool
}

// NewApp creates a new CLI application with default settings
func NewApp() *App {
	return &App{
		commands:      Commands{},
		usageTemplate: defaultUsageTemplate(),
		exitStatus:    0,
	}
}

func defaultUsageTemplate() string {
	return `{{.AppName}} is a command line tool
Usage:
	{{.AppName}} command [arguments]

The commands are:
{{range .Commands}}{{if .Runnable}}
	{{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use "{{.AppName}} help [command]" for more information about a command.
`
}

var (
	_defaultApp = NewApp()

	// Deprecated: use App instance instead
	_usageTemplate = `[webgo] is a web service base on web.go
Usage:
	[webgo] command [arguments]

The commands are:
{{range .}}{{if .Runnable}}
	{{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use "[webgo] help [command]" for more information about a command.
`

	_commands   = Commands{}
	_exitMu     sync.Mutex
	_exitStatus = 0
	_setFlags   func(f *flag.FlagSet)
)

// SetUsageTemplate sets a custom usage template for the app
func (a *App) SetUsageTemplate(usageTemplate string) {
	a.usageTemplate = usageTemplate
}

// SetFlags sets flags that will be added to all commands
func (a *App) SetFlags(f func(f *flag.FlagSet)) {
	a.setFlags = f
}

// AddCommands adds one or more commands to the app
func (a *App) AddCommands(cmds ...*Command) {
	a.commands = append(a.commands, cmds...)
	a.commandsSorted = false
}

// sortCommands sorts commands by name if not already sorted
func (a *App) sortCommands() {
	if !a.commandsSorted && len(a.commands) > 0 {
		sort.Slice(a.commands, func(i, j int) bool {
			return a.commands[i].Name() < a.commands[j].Name()
		})
		a.commandsSorted = true
	}
}

// getCommand gets a command by name
func (a *App) getCommand(name string) (*Command, error) {
	if len(a.commands) == 0 {
		return nil, fmt.Errorf("no commands")
	}

	a.sortCommands()
	cmd := a.commands.Search(name)

	if cmd == nil {
		return nil, fmt.Errorf("unknown sub command %q", name)
	}

	return cmd, nil
}

// ExecuteE executes the app and returns any error instead of exiting
func (a *App) ExecuteE() error {
	flag.Usage = func() { a.usage() }
	flag.Parse()
	log.SetFlags(0)

	args := flag.Args()

	if len(args) < 1 {
		a.printUsage(os.Stderr)
		return fmt.Errorf("no command specified")
	}

	if args[0] == "help" {
		return a.help(args[1:])
	}

	name := args[0]
	cmd, err := a.getCommand(name)

	if err != nil {
		return fmt.Errorf("cmd(%s): %w", name, err)
	}

	a.addFlags(&cmd.Flag)
	cmd.Flag.Usage = func() { cmd.UsageWithApp(a) }
	if err := cmd.Flag.Parse(args[1:]); err != nil {
		return fmt.Errorf("cmd(%s): %w", name, err)
	}

	if err := cmd.Run(cmd, cmd.Flag.Args()); err != nil {
		return fmt.Errorf("cmd(%s): %w", name, err)
	}

	return nil
}

// Execute executes the app and exits on error (for backward compatibility)
func (a *App) Execute() {
	if err := a.ExecuteE(); err != nil {
		log.Printf("%v\n", err)
		a.setExitStatus(1)
		a.exit()
	}
	a.exit()
}

func (a *App) addFlags(f *flag.FlagSet) {
	if a.setFlags != nil {
		a.setFlags(f)
	}
}

func (a *App) setExitStatus(n int) {
	a.exitMu.Lock()
	if a.exitStatus < n {
		a.exitStatus = n
	}
	a.exitMu.Unlock()
}

func (a *App) exit() {
	os.Exit(a.exitStatus)
}

// Package-level functions for backward compatibility

// SetUsageTemplate set value to usageTemplate
func SetUsageTemplate(usageTemplate string) {
	_usageTemplate = usageTemplate
	_defaultApp.SetUsageTemplate(usageTemplate)
}

// SetFlags set flags to all commands
func SetFlags(f func(f *flag.FlagSet)) {
	_setFlags = f
	_defaultApp.SetFlags(f)
}

// AddCommands Add Command.
func AddCommands(cmds ...*Command) {
	_commands = append(_commands, cmds...)
	_defaultApp.AddCommands(cmds...)
}

// getCommand get Command by name.
func getCommand(name string) (*Command, error) {
	if len(_commands) == 0 {
		return nil, fmt.Errorf("no commands")
	}

	cmd := _commands.Search(name)

	if cmd == nil {
		return nil, fmt.Errorf("unknown sub command %q", name)
	}

	return cmd, nil
}

// ExecuteE executes the default app and returns any error
func ExecuteE() error {
	return _defaultApp.ExecuteE()
}

// Execute func
func Execute() {
	flag.Usage = usage
	flag.Parse() // catch -h argument
	log.SetFlags(0)

	args := flag.Args()

	if len(args) < 1 {
		usage()
	}

	if args[0] == "help" {
		help(args[1:])
		return
	}

	name := args[0]
	cmd, err := getCommand(name)

	if err != nil {
		fatalf("cmd(%s): %v \n", name, err)
	}

	addFlags(&cmd.Flag)
	cmd.Flag.Usage = func() { cmd.Usage() }
	cmd.Flag.Parse(args[1:])

	if err := cmd.Run(cmd, cmd.Flag.Args()); err != nil {
		logf("cmd(%s): %v\n", name, err)
	}

	exit()
}

// Command struct
type Command struct {
	Run       func(cmd *Command, args []string) error
	Flag      flag.FlagSet
	UsageLine string
	Short     string
	Long      string
}

// Name string
func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.IndexRune(name, ' ')
	if i >= 0 {
		name = name[:i]
	}
	return name
}

// Usage u
func (c *Command) Usage() {
	help([]string{c.Name()})
	os.Exit(2)
}

// UsageWithApp prints usage for a command using the app instance
func (c *Command) UsageWithApp(a *App) {
	a.help([]string{c.Name()})
	os.Exit(2)
}

// Runnable bool
func (c *Command) Runnable() bool {
	return c.Run != nil
}

type Commands []*Command

// Search use binary search to find and return the smallest index *Command
func (c *Commands) Search(name string) *Command {

	i := sort.Search(len(*c), func(i int) bool { return (*c)[i].Name() >= name })

	if i < len(*c) && (*c)[i].Name() == name {
		return (*c)[i]
	}

	return nil
}

func (a *App) usage() {
	a.printUsage(os.Stderr)
	os.Exit(2)
}

func (a *App) printUsage(w io.Writer) {
	a.sortCommands()
	bw := bufio.NewWriter(w)

	// Prepare template data with app name
	data := struct {
		AppName  string
		Commands Commands
	}{
		AppName:  getAppName(),
		Commands: a.commands,
	}

	runTemplate(bw, a.usageTemplate, data)
	bw.Flush()
}

func (a *App) help(args []string) error {
	if len(args) == 0 {
		a.printUsage(os.Stdout)
		return nil
	}
	if len(args) != 1 {
		return fmt.Errorf("usage: help command\n\nToo many arguments given")
	}

	name := args[0]

	cmd, err := a.getCommand(name)

	if err != nil {
		return fmt.Errorf("help(%s): %w", name, err)
	}

	if cmd.Runnable() {
		fmt.Fprintf(os.Stdout, "usage: %s\n", cmd.UsageLine)
	}

	runTemplate(os.Stdout, cmd.Long, nil)
	return nil
}

func getAppName() string {
	if len(os.Args) > 0 {
		return os.Args[0]
	}
	return "app"
}

func usage() {
	printUsage(os.Stderr)
	os.Exit(2)
}

func printUsage(w io.Writer) {
	bw := bufio.NewWriter(w)
	runTemplate(bw, _usageTemplate, _commands)
	bw.Flush()
}

type errWriter struct {
	w   io.Writer
	err error
}

func (w *errWriter) Write(b []byte) (int, error) {
	n, err := w.w.Write(b)
	if err != nil {
		w.err = err
	}
	return n, err
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToTitle(r)) + s[n:]
}

func runTemplate(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{
		"trim":       strings.TrimSpace,
		"capitalize": capitalize,
	})
	template.Must(t.Parse(text))
	ew := &errWriter{w: w}
	err := t.Execute(ew, data)
	if ew.err != nil {
		if strings.Contains(ew.err.Error(), "pipe") {
			os.Exit(1)
		}
		fatalf("writing output: %v", ew.err)
	}
	if err != nil {
		panic(err)
	}
}

func help(args []string) {
	if len(args) == 0 {
		printUsage(os.Stdout)
		return
	}
	if len(args) != 1 {
		fatalf("usage: help command\n\nToo many arguments given.\n")
	}

	name := args[0]

	cmd, err := getCommand(name)

	if err != nil {
		fatalf("help(%s): %v \n", name, err)
	}

	if cmd.Runnable() {
		fmt.Fprintf(os.Stdout, "usage: %s\n", cmd.UsageLine)
	}

	runTemplate(os.Stdout, cmd.Long, nil)
}

func addFlags(f *flag.FlagSet) {
	if _setFlags != nil {
		_setFlags(f)
	}
}

func logf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func errorf(format string, args ...interface{}) {
	logf(format, args...)
	setExitStatus(1)
}

func fatalf(format string, args ...interface{}) {
	errorf(format, args...)
	exit()
}

func setExitStatus(n int) {
	_exitMu.Lock()
	if _exitStatus < n {
		_exitStatus = n
	}
	_exitMu.Unlock()
}

func exit() {
	os.Exit(_exitStatus)
}
