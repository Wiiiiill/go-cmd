package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Wiiiiill/go-cmd"
)

var (
	// 全局标志
	verbose bool
	config  string
)

// 版本命令
var cmdVersion = &cmd.Command{
	Run:       runVersion,
	UsageLine: "version",
	Short:     "显示版本信息",
	Long: `version 命令显示应用程序的版本号和构建信息。

示例：
	demo version
`,
}

func runVersion(c *cmd.Command, args []string) error {
	fmt.Println("Demo App v2.0.0")
	fmt.Println("使用新的 App 实例 API")
	if verbose {
		fmt.Println("详细模式: 启用")
		fmt.Printf("配置文件: %s\n", config)
	}
	return nil
}

// 构建命令
var cmdBuild = &cmd.Command{
	Run:       runBuild,
	UsageLine: "build [选项]",
	Short:     "构建项目",
	Long: `build 命令用于构建项目。

选项：
	-output string  输出文件名

示例：
	demo build
	demo build -output myapp
`,
}

var buildOutput string

func runBuild(c *cmd.Command, args []string) error {
	// 命令特定标志必须在这里定义，但会在 Execute 中自动解析
	if verbose {
		fmt.Println("开始构建...")
		fmt.Printf("输出文件: %s\n", buildOutput)
	}

	fmt.Println("✓ 构建成功！")
	return nil
}

func init() {
	// 在 init 中为命令添加标志
	cmdBuild.Flag.StringVar(&buildOutput, "output", "app", "输出文件名")
}

// 测试命令
var cmdTest = &cmd.Command{
	Run:       runTest,
	UsageLine: "test [包路径...]",
	Short:     "运行测试",
	Long: `test 命令运行项目的测试用例。

示例：
	demo test
	demo test ./...
`,
}

func runTest(c *cmd.Command, args []string) error {
	targets := []string{"./..."}
	if len(args) > 0 {
		targets = args
	}

	if verbose {
		fmt.Printf("运行测试: %v\n", targets)
	}

	fmt.Println("✓ 所有测试通过！")
	return nil
}

// 部署命令（演示错误处理）
var cmdDeploy = &cmd.Command{
	Run:       runDeploy,
	UsageLine: "deploy <环境>",
	Short:     "部署应用",
	Long: `deploy 命令将应用部署到指定环境。

参数：
	环境  部署目标环境 (dev, staging, prod)

示例：
	demo deploy dev
	demo deploy prod
`,
}

func runDeploy(c *cmd.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("缺少环境参数，请指定 dev, staging 或 prod")
	}

	env := args[0]
	validEnvs := map[string]bool{"dev": true, "staging": true, "prod": true}

	if !validEnvs[env] {
		return fmt.Errorf("无效的环境 '%s'，有效值: dev, staging, prod", env)
	}

	if verbose {
		fmt.Printf("部署到环境: %s\n", env)
	}

	fmt.Printf("✓ 成功部署到 %s 环境！\n", env)
	return nil
}

func main() {
	// 方式 1: 使用新的 App 实例 API（推荐）
	app := cmd.NewApp()

	// 设置全局标志
	app.SetFlags(func(f *flag.FlagSet) {
		f.BoolVar(&verbose, "verbose", false, "显示详细信息")
		f.BoolVar(&verbose, "v", false, "显示详细信息（简写）")
		f.StringVar(&config, "config", "config.json", "配置文件路径")
	})

	// 自定义帮助模板
	app.SetUsageTemplate(`Demo App - 演示 go-cmd v2.0 新特性
用法:
	{{.AppName}} <命令> [选项]

可用命令：
{{range .Commands}}{{if .Runnable}}
	{{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

全局选项:
	-v, -verbose     显示详细信息
	-config string   配置文件路径 (默认 "config.json")

使用 "{{.AppName}} help <命令>" 查看命令的详细信息。

新特性:
	✓ 命令自动排序 - 无需手动按字典序添加
	✓ ExecuteE 方法 - 返回错误而不是直接退出
	✓ App 实例化 - 更好的可测试性和并发安全
	✓ 动态帮助模板 - 自动使用程序名
`)

	// 添加命令（注意：按非字典序添加，演示自动排序功能）
	// 添加顺序: version, build, test, deploy
	// 字典序应该是: build, deploy, test, version
	app.AddCommands(cmdVersion, cmdBuild, cmdTest, cmdDeploy)

	// 使用 ExecuteE 进行自定义错误处理
	if err := app.ExecuteE(); err != nil {
		log.Printf("错误: %v\n", err)
		// 可以在这里添加自定义错误处理逻辑
		// 例如：发送错误通知、记录日志等
	}

	// 如果想要传统的行为（错误时自动退出），可以使用:
	// app.Execute()
}
