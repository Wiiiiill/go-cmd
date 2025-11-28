# go-cmd

一个简洁优雅的 Go 语言命令行工具库，帮助你快速构建功能强大的 CLI 应用程序。

## 📋 目录

- [功能特性](#功能特性)
- [安装](#安装)
- [快速开始](#快速开始)
- [核心概念](#核心概念)
- [API 文档](#api-文档)
- [完整示例](#完整示例)
- [高级用法](#高级用法)
- [最佳实践](#最佳实践)

## ✨ 功能特性

- **简洁易用** - 仅需几行代码即可创建强大的命令行工具
- **子命令支持** - 轻松管理多个子命令
- **标志参数** - 完整支持 Go 标准库 flag 包的所有功能
- **帮助系统** - 自动生成帮助文档和使用说明
- **自定义模板** - 支持自定义使用说明的显示模板
- **二分查找** - 高效的命令查找机制
- **零依赖** - 仅依赖 Go 标准库

## 📦 安装

```bash
go get github.com/Wiiiiill/go-cmd
```

在你的项目中引入：

```go
import "github.com/Wiiiiill/go-cmd"
```

## 🚀 快速开始

### 最简示例

创建一个带有单个命令的 CLI 工具：

```go
package main

import (
	"fmt"
	"github.com/Wiiiiill/go-cmd"
)

var cmdHello = &cmd.Command{
	Run:       runHello,
	UsageLine: "hello [name]",
	Short:     "打印问候语",
	Long:      "hello 命令用于打印友好的问候语。\n",
}

func runHello(c *cmd.Command, args []string) error {
	name := "World"
	if len(args) > 0 {
		name = args[0]
	}
	fmt.Printf("Hello, %s!\n", name)
	return nil
}

func main() {
	cmd.AddCommands(cmdHello)
	cmd.Execute()
}
```

运行示例：

```bash
# 显示帮助信息
$ ./myapp help

# 执行 hello 命令
$ ./myapp hello
Hello, World!

$ ./myapp hello Go
Hello, Go!

# 查看命令帮助
$ ./myapp help hello
```

## 🎯 核心概念

### Command 结构

`Command` 是库的核心结构，定义一个可执行的子命令：

```go
type Command struct {
	Run       func(cmd *Command, args []string) error  // 命令执行函数
	Flag      flag.FlagSet                              // 命令特定的标志参数
	UsageLine string                                    // 使用说明（格式：命令名 [参数]）
	Short     string                                    // 简短描述（显示在命令列表中）
	Long      string                                    // 详细说明（显示在 help 命令中）
}
```

### 主要方法

| 方法 | 说明 |
|------|------|
| `AddCommands(...*Command)` | 添加一个或多个命令 |
| `Execute()` | 执行命令行参数解析和命令调用 |
| `SetFlags(func(*flag.FlagSet))` | 设置全局标志参数 |
| `SetUsageTemplate(string)` | 自定义帮助信息模板 |

## 📚 API 文档

### AddCommands

添加命令到应用程序。

```go
func AddCommands(cmds ...*Command)
```

**示例：**

```go
cmd.AddCommands(cmdVersion, cmdBuild, cmdRun)
```

### Execute

解析命令行参数并执行相应的命令。这通常是 `main()` 函数中最后调用的函数。

```go
func Execute()
```

**功能：**
- 解析命令行参数
- 处理 `-h` 和 `--help` 标志
- 路由到相应的子命令
- 处理错误和退出状态

### SetFlags

设置所有命令共享的全局标志参数。

```go
func SetFlags(f func(f *flag.FlagSet))
```

**示例：**

```go
var (
	verbose bool
	config  string
)

func main() {
	cmd.SetFlags(func(f *flag.FlagSet) {
		f.BoolVar(&verbose, "verbose", false, "详细输出模式")
		f.StringVar(&config, "config", "config.json", "配置文件路径")
	})
	
	cmd.AddCommands(cmdRun)
	cmd.Execute()
}
```

### SetUsageTemplate

自定义帮助信息的显示模板。

```go
func SetUsageTemplate(usageTemplate string)
```

**默认模板：**

```
[webgo] is a web service base on web.go
Usage:
	[webgo] command [arguments]

The commands are:
{{range .}}{{if .Runnable}}
	{{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use "[webgo] help [command]" for more information about a command.
```

**自定义示例：**

```go
const customTemplate = `MyApp - 我的应用程序
用法：
	myapp <命令> [选项]

可用命令：
{{range .}}{{if .Runnable}}
	{{.Name | printf "%-15s"}} {{.Short}}{{end}}{{end}}

使用 "myapp help <命令>" 查看命令的详细信息。
更多信息请访问: https://github.com/username/myapp
`

func main() {
	cmd.SetUsageTemplate(customTemplate)
	// ...
}
```

## 💡 完整示例

### 示例 1：带版本命令的应用

```go
package main

import (
	"flag"
	"fmt"
	"github.com/Wiiiiill/go-cmd"
)

var (
	_version = "v1.0.0"
	_osarch  = "linux/amd64"  // 通常通过 ldflags 设置
	_force   bool
)

var cmdVersion = &cmd.Command{
	Run:       runVersion,
	UsageLine: "version",
	Short:     "显示版本信息",
	Long: `version 命令显示应用程序的版本号和构建信息。

使用示例：
	myapp version
	myapp version -force
`,
}

func runVersion(c *cmd.Command, args []string) error {
	fmt.Printf("版本: %s\n", _version)
	fmt.Printf("平台: %s\n", _osarch)
	if _force {
		fmt.Println("强制模式: 启用")
	}
	return nil
}

func main() {
	// 设置全局标志
	cmd.SetFlags(func(f *flag.FlagSet) {
		f.BoolVar(&_force, "force", false, "强制执行模式")
	})

	// 添加命令
	cmd.AddCommands(cmdVersion)
	
	// 执行
	cmd.Execute()
}
```

### 示例 2：多命令应用

```go
package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/Wiiiiill/go-cmd"
)

var (
	// 全局标志
	verbose bool
	
	// build 命令标志
	output string
	tags   string
)

var cmdBuild = &cmd.Command{
	Run:       runBuild,
	UsageLine: "build [选项] [包路径]",
	Short:     "编译项目",
	Long: `build 命令用于编译 Go 项目。

选项：
	-o string    输出文件名
	-tags string 构建标签

示例：
	myapp build
	myapp build -o myapp.exe
	myapp build -tags prod ./cmd/server
`,
}

var cmdTest = &cmd.Command{
	Run:       runTest,
	UsageLine: "test [包路径...]",
	Short:     "运行测试",
	Long: `test 命令运行项目的测试用例。

示例：
	myapp test
	myapp test ./...
	myapp test -verbose ./pkg/utils
`,
}

var cmdClean = &cmd.Command{
	Run:       runClean,
	UsageLine: "clean",
	Short:     "清理构建缓存",
	Long:      "clean 命令清理所有构建产物和缓存文件。\n",
}

func runBuild(c *cmd.Command, args []string) error {
	// 为 build 命令添加特定标志
	c.Flag.StringVar(&output, "o", "", "输出文件名")
	c.Flag.StringVar(&tags, "tags", "", "构建标签")
	
	target := "."
	if len(args) > 0 {
		target = args[0]
	}
	
	if verbose {
		fmt.Printf("正在构建: %s\n", target)
		fmt.Printf("输出: %s\n", output)
		fmt.Printf("标签: %s\n", tags)
	}
	
	fmt.Println("构建成功！")
	return nil
}

func runTest(c *cmd.Command, args []string) error {
	targets := []string{"./..."}
	if len(args) > 0 {
		targets = args
	}
	
	if verbose {
		fmt.Printf("运行测试: %v\n", targets)
	}
	
	fmt.Println("所有测试通过！")
	return nil
}

func runClean(c *cmd.Command, args []string) error {
	if verbose {
		fmt.Println("清理构建缓存...")
	}
	
	// 执行清理逻辑
	fmt.Println("清理完成！")
	return nil
}

func main() {
	// 设置全局标志
	cmd.SetFlags(func(f *flag.FlagSet) {
		f.BoolVar(&verbose, "verbose", false, "显示详细信息")
		f.BoolVar(&verbose, "v", false, "显示详细信息（简写）")
	})
	
	// 自定义帮助模板
	cmd.SetUsageTemplate(`MyApp - 项目构建工具
用法:
	myapp <命令> [选项]

可用命令：
{{range .}}{{if .Runnable}}
	{{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

全局选项:
	-v, -verbose  显示详细信息

使用 "myapp help <命令>" 查看命令的详细信息。
`)
	
	// 添加所有命令
	cmd.AddCommands(cmdBuild, cmdTest, cmdClean)
	
	// 执行
	cmd.Execute()
}
```

### 示例 3：实际项目应用（Web 服务）

```go
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"github.com/Wiiiiill/go-cmd"
)

var (
	// 全局配置
	configFile string
	
	// server 命令配置
	port int
	host string
	
	// migrate 命令配置
	migrateUp bool
)

var cmdServer = &cmd.Command{
	Run:       runServer,
	UsageLine: "server [选项]",
	Short:     "启动 Web 服务器",
	Long: `server 命令启动 HTTP/HTTPS Web 服务器。

选项：
	-host string  监听地址 (默认 "0.0.0.0")
	-port int     监听端口 (默认 8080)

示例：
	myapp server
	myapp server -host localhost -port 3000
	myapp server -config production.json
`,
}

var cmdMigrate = &cmd.Command{
	Run:       runMigrate,
	UsageLine: "migrate [选项]",
	Short:     "数据库迁移",
	Long: `migrate 命令执行数据库架构迁移。

选项：
	-up   执行向上迁移（默认）

示例：
	myapp migrate
	myapp migrate -up
`,
}

var cmdVersion = &cmd.Command{
	Run:       runVersion,
	UsageLine: "version",
	Short:     "显示版本信息",
	Long:      "显示应用程序版本和构建信息。\n",
}

func runServer(c *cmd.Command, args []string) error {
	// 添加命令特定标志
	c.Flag.StringVar(&host, "host", "0.0.0.0", "监听地址")
	c.Flag.IntVar(&port, "port", 8080, "监听端口")
	
	fmt.Printf("使用配置文件: %s\n", configFile)
	fmt.Printf("服务器启动在 %s:%d\n", host, port)
	
	// 处理优雅关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		<-sigChan
		fmt.Println("\n正在关闭服务器...")
		os.Exit(0)
	}()
	
	// 模拟服务器运行
	select {}
}

func runMigrate(c *cmd.Command, args []string) error {
	c.Flag.BoolVar(&migrateUp, "up", true, "执行向上迁移")
	
	fmt.Println("开始数据库迁移...")
	
	if migrateUp {
		fmt.Println("执行向上迁移")
		// 执行迁移逻辑
	}
	
	fmt.Println("迁移完成！")
	return nil
}

func runVersion(c *cmd.Command, args []string) error {
	fmt.Println("MyApp v2.0.0")
	fmt.Println("Build: 2024-01-01")
	return nil
}

func main() {
	// 设置全局标志
	cmd.SetFlags(func(f *flag.FlagSet) {
		f.StringVar(&configFile, "config", "config.json", "配置文件路径")
	})
	
	// 添加命令
	cmd.AddCommands(cmdServer, cmdMigrate, cmdVersion)
	
	// 执行
	cmd.Execute()
}
```

## 🔥 高级用法

### 错误处理

命令的 `Run` 函数返回 `error`，你可以返回错误来指示命令执行失败：

```go
func runDeploy(c *cmd.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("缺少部署目标参数")
	}
	
	target := args[0]
	if err := deploy(target); err != nil {
		return fmt.Errorf("部署失败: %w", err)
	}
	
	return nil
}
```

### 命令特定标志

每个命令可以有自己的标志参数，通过 `Command.Flag` 字段设置：

```go
var (
	recursive bool
	exclude   string
)

func runCopy(c *cmd.Command, args []string) error {
	// 在 Run 函数中添加命令特定标志
	c.Flag.BoolVar(&recursive, "r", false, "递归复制")
	c.Flag.StringVar(&exclude, "exclude", "", "排除模式")
	
	// 使用标志
	if recursive {
		fmt.Println("递归模式已启用")
	}
	
	return nil
}
```

### 动态命令注册

你可以根据条件动态注册命令：

```go
func main() {
	commands := []*cmd.Command{cmdVersion}
	
	// 仅在开发环境添加调试命令
	if os.Getenv("ENV") == "development" {
		commands = append(commands, cmdDebug)
	}
	
	cmd.AddCommands(commands...)
	cmd.Execute()
}
```

### 使用 ldflags 设置版本信息

在编译时通过 ldflags 注入版本信息：

```go
var (
	version = "dev"      // 默认值
	commit  = "unknown"
	date    = "unknown"
)

var cmdVersion = &cmd.Command{
	Run: func(c *cmd.Command, args []string) error {
		fmt.Printf("版本: %s\n", version)
		fmt.Printf("提交: %s\n", commit)
		fmt.Printf("日期: %s\n", date)
		return nil
	},
	UsageLine: "version",
	Short:     "显示版本信息",
	Long:      "显示详细的版本和构建信息。\n",
}
```

编译命令：

```bash
go build -ldflags "-X main.version=v1.0.0 -X main.commit=$(git rev-parse HEAD) -X main.date=$(date -u +%Y-%m-%d)"
```

## 📖 最佳实践

### 1. 组织命令文件

对于大型项目，建议将每个命令放在单独的文件中：

```
myapp/
├── main.go
├── cmd_version.go
├── cmd_build.go
├── cmd_test.go
└── cmd_deploy.go
```

### 2. 命令命名规范

- 使用清晰、简洁的命令名
- 使用动词形式（如 `build`, `run`, `deploy`）
- 保持命令名全小写
- 多个单词使用连字符（如 `cache-clear`）

### 3. 帮助文档编写

- `Short`：一行简短描述（不超过 50 字符）
- `Long`：详细说明，包括用法示例和选项说明
- `UsageLine`：格式为 "命令名 [必选参数] [可选参数]"

### 4. 标志参数设计

- 提供简短和完整两种形式（如 `-v` 和 `-verbose`）
- 为所有标志提供合理的默认值
- 在帮助文档中说明每个标志的作用

### 5. 错误处理

- 返回有意义的错误信息
- 使用 `fmt.Errorf` 和 `%w` 进行错误包装
- 避免在 `Run` 函数中调用 `os.Exit`

## 🔧 常见问题

### Q: 如何添加子命令的子命令？

A: 当前库设计为单层命令结构。如果需要多层子命令，建议在命令名中使用连字符，如 `cache-clear`、`cache-list`。

### Q: 如何实现命令别名？

A: 可以创建多个 `Command` 实例指向同一个 `Run` 函数：

```go
var cmdRun = &cmd.Command{
	Run:       doRun,
	UsageLine: "run [选项]",
	Short:     "运行应用",
	Long:      "...",
}

var cmdStart = &cmd.Command{
	Run:       doRun,  // 相同的 Run 函数
	UsageLine: "start [选项]",
	Short:     "运行应用（别名：run）",
	Long:      "...",
}

func main() {
	cmd.AddCommands(cmdRun, cmdStart)
	cmd.Execute()
}
```

### Q: 命令执行顺序如何控制？

A: 命令在帮助信息中按字母顺序显示。库内部使用二分查找以提高查找效率。

### Q: 如何禁用自动帮助命令？

A: 帮助命令是内置的。如果需要自定义帮助行为，可以通过 `SetUsageTemplate` 修改帮助信息的显示方式。

## 📄 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📞 联系方式

- GitHub: [https://github.com/Wiiiiill/go-cmd](https://github.com/Wiiiiill/go-cmd)
- Issues: [https://github.com/Wiiiiill/go-cmd/issues](https://github.com/Wiiiiill/go-cmd/issues)

---

**Happy Coding! 🎉**
