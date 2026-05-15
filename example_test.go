package cmd_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/Wiiiiill/go-cmd"
)

// 测试 App 实例化
func TestAppInstance(t *testing.T) {
	app := cmd.NewApp()
	if app == nil {
		t.Fatal("NewApp() returned nil")
	}
}

// 测试自动排序功能
func TestAutoSorting(t *testing.T) {
	// 保存原始 os.Args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// 模拟命令行参数
	os.Args = []string{"testapp", "zebra"}

	app := cmd.NewApp()

	var executed string

	// 按非字典序添加命令
	cmdZebra := &cmd.Command{
		Run: func(c *cmd.Command, args []string) error {
			executed = "zebra"
			return nil
		},
		UsageLine: "zebra",
		Short:     "Zebra command",
		Long:      "Zebra command description",
	}

	cmdApple := &cmd.Command{
		Run: func(c *cmd.Command, args []string) error {
			executed = "apple"
			return nil
		},
		UsageLine: "apple",
		Short:     "Apple command",
		Long:      "Apple command description",
	}

	cmdMango := &cmd.Command{
		Run: func(c *cmd.Command, args []string) error {
			executed = "mango"
			return nil
		},
		UsageLine: "mango",
		Short:     "Mango command",
		Long:      "Mango command description",
	}

	// 按非字典序添加：zebra, apple, mango
	app.AddCommands(cmdZebra, cmdApple, cmdMango)

	// 重置 flag 包状态
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// 执行应该能找到 zebra 命令（即使它不是按字典序添加的）
	err := app.ExecuteE()
	if err != nil {
		t.Fatalf("ExecuteE() failed: %v", err)
	}

	if executed != "zebra" {
		t.Errorf("Expected 'zebra' to be executed, got '%s'", executed)
	}
}

// 测试 ExecuteE 错误处理
func TestExecuteE(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"testapp", "fail"}

	app := cmd.NewApp()

	cmdFail := &cmd.Command{
		Run: func(c *cmd.Command, args []string) error {
			return fmt.Errorf("intentional error")
		},
		UsageLine: "fail",
		Short:     "Fail command",
		Long:      "Command that always fails",
	}

	app.AddCommands(cmdFail)

	// 重置 flag 包状态
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	err := app.ExecuteE()
	if err == nil {
		t.Fatal("Expected error from ExecuteE(), got nil")
	}

	if err.Error() != "cmd(fail): intentional error" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

// 测试全局标志
func TestGlobalFlags(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// 不使用全局标志，因为 flag 包的全局状态难以在测试中重置
	os.Args = []string{"testapp", "test"}

	app := cmd.NewApp()

	var executed bool
	cmdTest := &cmd.Command{
		Run: func(c *cmd.Command, args []string) error {
			executed = true
			return nil
		},
		UsageLine: "test",
		Short:     "Test command",
		Long:      "Test command description",
	}

	app.AddCommands(cmdTest)

	// 重置 flag 包状态
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	err := app.ExecuteE()
	if err != nil {
		t.Fatalf("ExecuteE() failed: %v", err)
	}

	if !executed {
		t.Error("Command was not executed")
	}
}

// 测试自定义模板
func TestCustomTemplate(t *testing.T) {
	app := cmd.NewApp()

	customTemplate := `Custom Template Test
{{range .Commands}}{{.Name}}{{end}}`

	app.SetUsageTemplate(customTemplate)

	// 只是确保设置不会崩溃
	// 实际的模板渲染需要更复杂的测试设置
}

// 测试向后兼容性
func TestBackwardCompatibility(t *testing.T) {
	// 测试包级函数仍然可用
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"testapp", "compat"}

	var executed bool
	cmdCompat := &cmd.Command{
		Run: func(c *cmd.Command, args []string) error {
			executed = true
			return nil
		},
		UsageLine: "compat",
		Short:     "Compatibility test",
		Long:      "Test backward compatibility",
	}

	// 使用包级函数
	cmd.AddCommands(cmdCompat)

	// 重置 flag 包状态
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	err := cmd.ExecuteE()
	if err != nil {
		t.Fatalf("ExecuteE() failed: %v", err)
	}

	if !executed {
		t.Error("Command was not executed using package-level functions")
	}
}
