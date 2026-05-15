# Changelog

## [2.0.0] - 2026-05-15

### 🎉 重大更新

这是一个重大版本更新，引入了多项改进，同时保持向后兼容性。

### ✨ 新增功能

#### 1. App 实例化 API
- 新增 `NewApp()` 函数创建独立的 CLI 应用实例
- 每个 App 实例拥有独立的命令列表和配置
- 提高了可测试性和并发安全性

```go
// 新的推荐方式
app := cmd.NewApp()
app.AddCommands(cmdVersion, cmdBuild)
app.Execute()
```

#### 2. ExecuteE 错误处理
- 新增 `ExecuteE() error` 方法，返回错误而不是直接退出
- 允许用户自定义错误处理逻辑
- 适用于需要精细控制错误处理的场景

```go
if err := app.ExecuteE(); err != nil {
    log.Printf("错误: %v", err)
    // 自定义错误处理
}
```

#### 3. 自动命令排序
- 命令在执行时自动按字典序排序
- **无需手动按字典序注册命令**
- 消除了之前版本中容易出错的约束

```go
// 可以按任意顺序添加命令
app.AddCommands(cmdVersion, cmdBuild, cmdTest, cmdClean)
// 内部会自动排序为: build, clean, test, version
```

#### 4. 动态帮助模板
- 帮助模板自动使用当前程序名（从 `os.Args[0]` 获取）
- 不再需要硬编码应用名称
- 模板数据结构更新为包含 `AppName` 和 `Commands` 字段

```go
// 新的模板格式
{{.AppName}} is a command line tool
Usage:
    {{.AppName}} command [arguments]

The commands are:
{{range .Commands}}{{if .Runnable}}
    {{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}
```

### 🔄 改进

- **性能优化**: 命令排序采用延迟执行，仅在首次查找时排序一次
- **代码结构**: 全局状态封装到 App 结构体中，代码更清晰
- **错误信息**: 使用 `%w` 进行错误包装，提供更好的错误追踪

### 🔧 内部变化

- 引入 `App` 结构体封装应用状态
- 添加 `sortCommands()` 方法实现自动排序
- 添加 `commandsSorted` 标志避免重复排序
- 更新 `Command.UsageWithApp()` 方法支持 App 实例

### ⚠️ 向后兼容性

**完全向后兼容！** 所有现有代码无需修改即可继续工作：

- 包级函数（`AddCommands`, `Execute`, `SetFlags`, `SetUsageTemplate`）仍然可用
- 现有的命令定义和使用方式保持不变
- 旧的模板格式仍然支持（通过包级函数使用时）

### 📝 迁移指南

#### 从 v1.x 迁移到 v2.0

**选项 1: 不做任何改动（推荐用于现有项目）**
```go
// v1.x 代码继续工作
cmd.AddCommands(cmdBuild, cmdClean, cmdTest)
cmd.Execute()
```

**选项 2: 迁移到新 API（推荐用于新项目）**
```go
// v2.0 新 API
app := cmd.NewApp()
app.AddCommands(cmdBuild, cmdClean, cmdTest)  // 无需按字典序
app.Execute()  // 或使用 app.ExecuteE() 进行错误处理
```

**选项 3: 使用 ExecuteE 进行错误处理**
```go
// 使用包级函数 + ExecuteE
cmd.AddCommands(cmdBuild, cmdClean, cmdTest)
if err := cmd.ExecuteE(); err != nil {
    log.Fatalf("Error: %v", err)
}
```

### 🐛 修复

- 修复了命令必须按字典序注册的限制
- 改进了错误消息的格式和可读性

### 📚 文档更新

- 更新 README.md，添加新 API 的使用示例
- 移除了"命令注册顺序（重要）"章节中的字典序要求说明
- 添加了 App 实例化和 ExecuteE 的文档
- 更新了所有示例代码

### 🧪 测试

- 添加了完整的单元测试套件
- 测试覆盖：App 实例化、自动排序、ExecuteE、向后兼容性
- 添加了示例程序演示新特性

### 🙏 致谢

感谢所有使用和反馈的用户！

---

## [1.0.0] - 之前版本

初始版本，提供基础的命令行工具构建功能。
