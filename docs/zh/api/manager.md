# 管理器 API

`Manager` 是 Go Pip SDK 的核心组件，提供了所有包管理、虚拟环境和项目操作的统一接口。

## 接口定义

```go
type PipManager interface {
    // 系统操作
    IsInstalled() (bool, error)
    Install() error
    GetVersion() (string, error)
    SetTimeout(timeout time.Duration)
    SetRetries(retries int)
    
    // 包操作
    InstallPackage(pkg *PackageSpec) error
    UninstallPackage(name string) error
    ListPackages() ([]*Package, error)
    ShowPackage(name string) (*PackageInfo, error)
    SearchPackages(query string) ([]*SearchResult, error)
    FreezePackages() ([]*Package, error)
    InstallRequirements(path string) error
    GenerateRequirements(path string) error
    
    // 虚拟环境操作
    CreateVenv(path string) error
    CreateVenvWithOptions(path string, opts *VenvOptions) error
    ActivateVenv(path string) error
    DeactivateVenv() error
    IsVenvActive() (bool, string)
    RemoveVenv(path string) error
    ListVenvs(baseDir string) ([]*VenvInfo, error)
    GetVenvInfo(path string) (*VenvInfo, error)
    
    // 项目操作
    InitProject(path string, opts *ProjectOptions) error
    ReadProjectConfig(path string) (*ProjectConfig, error)
    UpdateProjectVersion(path string, version string) error
    BuildProject(path string, opts *BuildOptions) error
    
    // 配置和日志
    SetCustomLogger(logger Logger)
    SetLogLevel(level LogLevel)
    GetConfig() *Config
    UpdateConfig(config *Config) error
}
```

## 构造函数

### NewManager

创建一个新的管理器实例。

```go
func NewManager(config *Config) PipManager
```

**参数：**
- `config` - 管理器配置，如果为 `nil` 则使用默认配置

**返回值：**
- `PipManager` - 管理器实例

**示例：**

```go
// 使用默认配置
manager := pip.NewManager(nil)

// 使用自定义配置
config := &pip.Config{
    PythonPath: "/usr/bin/python3",
    Timeout:    60 * time.Second,
    LogLevel:   "DEBUG",
}
manager := pip.NewManager(config)
```

### NewManagerWithContext

创建一个带上下文的管理器实例。

```go
func NewManagerWithContext(ctx context.Context, config *Config) PipManager
```

**参数：**
- `ctx` - 上下文，用于取消操作
- `config` - 管理器配置

**返回值：**
- `PipManager` - 管理器实例

**示例：**

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

manager := pip.NewManagerWithContext(ctx, nil)
```

## 系统操作

### IsInstalled

检查 pip 是否已安装。

```go
func (m *Manager) IsInstalled() (bool, error)
```

**返回值：**
- `bool` - 如果 pip 已安装则为 true
- `error` - 错误信息

**示例：**

```go
installed, err := manager.IsInstalled()
if err != nil {
    log.Fatalf("检查 pip 安装状态失败: %v", err)
}

if !installed {
    fmt.Println("Pip 未安装")
}
```

### Install

安装 pip。

```go
func (m *Manager) Install() error
```

**返回值：**
- `error` - 错误信息

**示例：**

```go
if err := manager.Install(); err != nil {
    log.Fatalf("安装 pip 失败: %v", err)
}
fmt.Println("Pip 安装成功")
```

### GetVersion

获取 pip 版本。

```go
func (m *Manager) GetVersion() (string, error)
```

**返回值：**
- `string` - pip 版本号
- `error` - 错误信息

**示例：**

```go
version, err := manager.GetVersion()
if err != nil {
    log.Fatalf("获取 pip 版本失败: %v", err)
}
fmt.Printf("Pip 版本: %s\n", version)
```

## 配置管理

### SetTimeout

设置操作超时时间。

```go
func (m *Manager) SetTimeout(timeout time.Duration)
```

**参数：**
- `timeout` - 超时时间

**示例：**

```go
// 设置 2 分钟超时
manager.SetTimeout(2 * time.Minute)
```

### SetRetries

设置重试次数。

```go
func (m *Manager) SetRetries(retries int)
```

**参数：**
- `retries` - 重试次数

**示例：**

```go
// 设置最多重试 5 次
manager.SetRetries(5)
```

### GetConfig

获取当前配置。

```go
func (m *Manager) GetConfig() *Config
```

**返回值：**
- `*Config` - 当前配置

**示例：**

```go
config := manager.GetConfig()
fmt.Printf("Python 路径: %s\n", config.PythonPath)
fmt.Printf("超时时间: %v\n", config.Timeout)
```

### UpdateConfig

更新配置。

```go
func (m *Manager) UpdateConfig(config *Config) error
```

**参数：**
- `config` - 新的配置

**返回值：**
- `error` - 错误信息

**示例：**

```go
newConfig := &pip.Config{
    Timeout:      120 * time.Second,
    DefaultIndex: "https://pypi.tuna.tsinghua.edu.cn/simple/",
}

if err := manager.UpdateConfig(newConfig); err != nil {
    log.Fatalf("更新配置失败: %v", err)
}
```

## 日志管理

### SetCustomLogger

设置自定义日志记录器。

```go
func (m *Manager) SetCustomLogger(logger Logger)
```

**参数：**
- `logger` - 自定义日志记录器

**示例：**

```go
logger, err := pip.NewLogger(&pip.LoggerConfig{
    Level:  pip.LogLevelDebug,
    Output: os.Stdout,
    Format: pip.LogFormatJSON,
})
if err != nil {
    log.Fatal(err)
}

manager.SetCustomLogger(logger)
```

### SetLogLevel

设置日志级别。

```go
func (m *Manager) SetLogLevel(level LogLevel)
```

**参数：**
- `level` - 日志级别

**示例：**

```go
// 设置为调试级别
manager.SetLogLevel(pip.LogLevelDebug)

// 设置为错误级别
manager.SetLogLevel(pip.LogLevelError)
```

## 实用方法

### AddTrustedHost

添加受信任的主机。

```go
func (m *Manager) AddTrustedHost(host string) error
```

**参数：**
- `host` - 主机名

**返回值：**
- `error` - 错误信息

**示例：**

```go
if err := manager.AddTrustedHost("pypi.tuna.tsinghua.edu.cn"); err != nil {
    log.Printf("添加受信任主机失败: %v", err)
}
```

### GetPythonPath

获取当前使用的 Python 路径。

```go
func (m *Manager) GetPythonPath() (string, error)
```

**返回值：**
- `string` - Python 可执行文件路径
- `error` - 错误信息

**示例：**

```go
pythonPath, err := manager.GetPythonPath()
if err != nil {
    log.Fatalf("获取 Python 路径失败: %v", err)
}
fmt.Printf("Python 路径: %s\n", pythonPath)
```

### GetPipPath

获取当前使用的 pip 路径。

```go
func (m *Manager) GetPipPath() (string, error)
```

**返回值：**
- `string` - pip 可执行文件路径
- `error` - 错误信息

**示例：**

```go
pipPath, err := manager.GetPipPath()
if err != nil {
    log.Fatalf("获取 pip 路径失败: %v", err)
}
fmt.Printf("Pip 路径: %s\n", pipPath)
```

## 状态查询

### GetStatus

获取管理器状态信息。

```go
func (m *Manager) GetStatus() (*ManagerStatus, error)
```

**返回值：**
- `*ManagerStatus` - 状态信息
- `error` - 错误信息

**示例：**

```go
status, err := manager.GetStatus()
if err != nil {
    log.Fatalf("获取状态失败: %v", err)
}

fmt.Printf("Python 版本: %s\n", status.PythonVersion)
fmt.Printf("Pip 版本: %s\n", status.PipVersion)
fmt.Printf("虚拟环境: %s\n", status.VirtualEnv)
fmt.Printf("已安装包数量: %d\n", status.PackageCount)
```

### IsHealthy

检查管理器是否处于健康状态。

```go
func (m *Manager) IsHealthy() (bool, []string, error)
```

**返回值：**
- `bool` - 是否健康
- `[]string` - 健康检查问题列表
- `error` - 错误信息

**示例：**

```go
healthy, issues, err := manager.IsHealthy()
if err != nil {
    log.Fatalf("健康检查失败: %v", err)
}

if !healthy {
    fmt.Println("发现以下问题:")
    for _, issue := range issues {
        fmt.Printf("- %s\n", issue)
    }
}
```

## 高级功能

### ExecuteCommand

执行自定义 pip 命令。

```go
func (m *Manager) ExecuteCommand(args []string) (*CommandResult, error)
```

**参数：**
- `args` - 命令参数

**返回值：**
- `*CommandResult` - 命令执行结果
- `error` - 错误信息

**示例：**

```go
// 执行 pip list --format=json
result, err := manager.ExecuteCommand([]string{"list", "--format=json"})
if err != nil {
    log.Fatalf("执行命令失败: %v", err)
}

fmt.Printf("输出: %s\n", result.Stdout)
fmt.Printf("错误: %s\n", result.Stderr)
fmt.Printf("退出代码: %d\n", result.ExitCode)
```

### Clone

克隆管理器实例。

```go
func (m *Manager) Clone() PipManager
```

**返回值：**
- `PipManager` - 克隆的管理器实例

**示例：**

```go
// 克隆管理器用于不同的配置
clonedManager := manager.Clone()
clonedManager.SetTimeout(30 * time.Second)
```

## 错误处理

所有管理器方法都返回详细的错误信息。使用 `pip.IsErrorType()` 函数检查特定错误类型：

```go
if err := manager.InstallPackage(pkg); err != nil {
    switch {
    case pip.IsErrorType(err, pip.ErrorTypePythonNotFound):
        fmt.Println("Python 未找到")
    case pip.IsErrorType(err, pip.ErrorTypePipNotFound):
        fmt.Println("Pip 未找到")
    case pip.IsErrorType(err, pip.ErrorTypeNetworkError):
        fmt.Println("网络错误")
    case pip.IsErrorType(err, pip.ErrorTypePackageNotFound):
        fmt.Println("包未找到")
    default:
        fmt.Printf("其他错误: %v\n", err)
    }
}
```

## 线程安全

`Manager` 类型是线程安全的，可以在多个 goroutine 中同时使用：

```go
var wg sync.WaitGroup
packages := []string{"requests", "click", "pydantic"}

for _, name := range packages {
    wg.Add(1)
    go func(packageName string) {
        defer wg.Done()
        
        pkg := &pip.PackageSpec{Name: packageName}
        if err := manager.InstallPackage(pkg); err != nil {
            fmt.Printf("安装 %s 失败: %v\n", packageName, err)
        }
    }(name)
}

wg.Wait()
```

## 最佳实践

1. **重用管理器实例**：创建一个管理器实例并在整个应用程序中重用它
2. **适当的超时设置**：根据网络条件设置合理的超时时间
3. **错误处理**：始终检查和处理错误
4. **日志记录**：在生产环境中启用适当的日志级别
5. **资源清理**：在适当的时候停用虚拟环境和关闭日志记录器

## 下一步

- 查看[包操作 API](./package-operations.md)
- 了解[虚拟环境 API](./virtual-environments.md)
- 探索[项目管理 API](./project-management.md)
- 学习[错误处理](./errors.md)
