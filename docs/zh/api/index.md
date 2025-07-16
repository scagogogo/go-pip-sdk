# API 参考

Go Pip SDK 提供了一套全面的 API，用于管理 Python 包、虚拟环境和项目。本节记录了所有公开接口、类型和函数。

## 核心组件

### [管理器](/zh/api/manager)
pip 操作的主要接口。提供包管理、虚拟环境操作和项目初始化的方法。

### [包操作](/zh/api/package-operations)
用于安装、卸载、列出和管理 Python 包的函数。

### [虚拟环境](/zh/api/virtual-environments)
用于创建、激活和管理 Python 虚拟环境的 API。

### [项目管理](/zh/api/project-management)
用于使用标准结构和配置文件初始化 Python 项目的工具。

### [类型](/zh/api/types)
整个 SDK 中使用的核心数据结构和接口。

### [错误](/zh/api/errors)
用于强大错误管理的错误类型和处理机制。

### [日志](/zh/api/logger)
具有可配置级别和输出格式的日志系统。

### [安装器](/zh/api/installer)
跨平台 pip 安装功能。

## 快速参考

### 主要接口

```go
type PipManager interface {
    // 系统操作
    IsInstalled() (bool, error)
    Install() error
    GetVersion() (string, error)

    // 包操作
    InstallPackage(pkg *PackageSpec) error
    UninstallPackage(name string) error
    ListPackages() ([]*Package, error)
    ShowPackage(name string) (*PackageInfo, error)
    SearchPackages(query string) ([]*SearchResult, error)
    FreezePackages() ([]*Package, error)

    // 虚拟环境操作
    CreateVenv(path string) error
    ActivateVenv(path string) error
    DeactivateVenv() error
    RemoveVenv(path string) error

    // 项目操作
    InitProject(path string, opts *ProjectOptions) error
    InstallRequirements(path string) error
    GenerateRequirements(path string) error
}
```

### 关键类型

```go
// 管理器配置
type Config struct {
    PythonPath   string
    PipPath      string
    DefaultIndex string
    TrustedHosts []string
    Timeout      time.Duration
    Retries      int
    LogLevel     string
    CacheDir     string
    ExtraOptions map[string]string
    Environment  map[string]string
}

// 包规范
type PackageSpec struct {
    Name           string
    Version        string
    Extras         []string
    Index          string
    Options        map[string]string
    Editable       bool
    Upgrade        bool
    ForceReinstall bool
}

// 已安装包信息
type Package struct {
    Name      string
    Version   string
    Location  string
    Editable  bool
    Installer string
}
```

## 错误处理

所有 API 函数都返回实现标准 Go error 接口的错误。SDK 提供了具有附加上下文的丰富错误类型：

```go
type PipErrorDetails struct {
    Type        ErrorType
    Message     string
    Command     string
    Output      string
    ExitCode    int
    Suggestions []string
    Context     map[string]string
    Cause       error
}
```

## 使用模式

### 基本用法

```go
// 使用默认配置创建管理器
manager := pip.NewManager(nil)

// 检查 pip 是否已安装
installed, err := manager.IsInstalled()
if err != nil {
    return err
}

if !installed {
    // 如果不可用则安装 pip
    if err := manager.Install(); err != nil {
        return err
    }
}
```

### 使用自定义配置

```go
config := &pip.Config{
    Timeout:      60 * time.Second,
    Retries:      5,
    DefaultIndex: "https://pypi.org/simple/",
    LogLevel:     "DEBUG",
}

manager := pip.NewManager(config)
```

### 使用上下文

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

manager := pip.NewManagerWithContext(ctx, nil)
```

## 线程安全

Manager 类型对于并发使用是安全的。多个 goroutine 可以同时调用同一个 Manager 实例的方法。

## 平台支持

所有 API 在支持的平台上都能一致工作：
- Windows (x86, x64, ARM64)
- macOS (Intel, Apple Silicon)
- Linux (x86, x64, ARM, ARM64)

平台特定的行为由 SDK 自动处理。
