# 快速开始

欢迎使用 Go Pip SDK！本指南将帮助您在几分钟内启动并运行 SDK。

## 前提条件

开始之前，请确保您有：

- **Go 1.19 或更高版本** 安装在您的系统上
- **Python 3.7 或更高版本**（用于 pip 操作）
- 对 Go 编程的基本了解

## 安装

使用 Go 模块安装 SDK：

```bash
go get github.com/scagogogo/go-pip-sdk
```

## 您的第一个程序

让我们创建一个简单的程序来检查 pip 是否已安装并安装一个包：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    // 使用默认配置创建新的 pip 管理器
    manager := pip.NewManager(nil)
    
    // 检查 pip 是否已安装
    installed, err := manager.IsInstalled()
    if err != nil {
        log.Fatalf("检查 pip 安装失败: %v", err)
    }
    
    if !installed {
        fmt.Println("Pip 未安装。正在安装...")
        if err := manager.Install(); err != nil {
            log.Fatalf("安装 pip 失败: %v", err)
        }
        fmt.Println("Pip 安装成功！")
    } else {
        fmt.Println("Pip 已经安装。")
    }
    
    // 获取 pip 版本
    version, err := manager.GetVersion()
    if err != nil {
        log.Fatalf("获取 pip 版本失败: %v", err)
    }
    fmt.Printf("Pip 版本: %s\n", version)
    
    // 安装包
    pkg := &pip.PackageSpec{
        Name:    "requests",
        Version: ">=2.25.0",
    }
    
    fmt.Printf("正在安装包: %s\n", pkg.Name)
    if err := manager.InstallPackage(pkg); err != nil {
        log.Fatalf("安装包失败: %v", err)
    }
    
    fmt.Println("包安装成功！")
    
    // 列出已安装的包
    packages, err := manager.ListPackages()
    if err != nil {
        log.Fatalf("列出包失败: %v", err)
    }
    
    fmt.Printf("找到 %d 个已安装的包:\n", len(packages))
    for _, pkg := range packages {
        fmt.Printf("  %s %s\n", pkg.Name, pkg.Version)
    }
}
```

将此保存为 `main.go` 并运行：

```bash
go run main.go
```

## 基本概念

### 管理器

`Manager` 是提供所有 pip 功能的中心组件。它实现了 `PipManager` 接口并处理：

- 系统操作（检查 pip 安装、获取版本）
- 包操作（安装、卸载、列出、显示）
- 虚拟环境操作
- 项目管理

### 配置

您可以使用 `Config` 结构自定义管理器的行为：

```go
config := &pip.Config{
    PythonPath:   "/usr/bin/python3",
    PipPath:      "/usr/bin/pip3",
    Timeout:      60 * time.Second,
    Retries:      3,
    LogLevel:     "DEBUG",
    DefaultIndex: "https://pypi.org/simple/",
}

manager := pip.NewManager(config)
```

### 包规范

安装包时，您使用 `PackageSpec` 来指定要求：

```go
// 基本包
pkg := &pip.PackageSpec{
    Name: "requests",
}

// 带版本约束的包
pkg := &pip.PackageSpec{
    Name:    "django",
    Version: ">=4.0,<5.0",
}

// 带额外依赖的包
pkg := &pip.PackageSpec{
    Name:   "fastapi",
    Extras: []string{"dev", "test"},
}

// 带自定义选项的包
pkg := &pip.PackageSpec{
    Name:    "numpy",
    Upgrade: true,
    Options: map[string]string{
        "no-cache-dir": "",
        "timeout":      "120",
    },
}
```

## 常见操作

### 安装包

```go
// 安装单个包
pkg := &pip.PackageSpec{Name: "requests"}
if err := manager.InstallPackage(pkg); err != nil {
    return err
}

// 安装多个包
packages := []*pip.PackageSpec{
    {Name: "requests"},
    {Name: "click"},
    {Name: "pydantic"},
}

for _, pkg := range packages {
    if err := manager.InstallPackage(pkg); err != nil {
        fmt.Printf("安装 %s 失败: %v\n", pkg.Name, err)
    }
}
```

### 使用虚拟环境

```go
// 创建虚拟环境
venvPath := "/path/to/my-venv"
if err := manager.CreateVenv(venvPath); err != nil {
    return err
}

// 激活虚拟环境
if err := manager.ActivateVenv(venvPath); err != nil {
    return err
}

// 在虚拟环境中安装包
pkg := &pip.PackageSpec{Name: "requests"}
if err := manager.InstallPackage(pkg); err != nil {
    return err
}

// 完成后停用
if err := manager.DeactivateVenv(); err != nil {
    return err
}
```

### 获取包信息

```go
// 列出所有已安装的包
packages, err := manager.ListPackages()
if err != nil {
    return err
}

for _, pkg := range packages {
    fmt.Printf("%s %s\n", pkg.Name, pkg.Version)
}

// 获取包的详细信息
info, err := manager.ShowPackage("requests")
if err != nil {
    return err
}

fmt.Printf("名称: %s\n", info.Name)
fmt.Printf("版本: %s\n", info.Version)
fmt.Printf("摘要: %s\n", info.Summary)
fmt.Printf("依赖: %v\n", info.Requires)
```

## 错误处理

SDK 提供了具有特定错误类型的结构化错误处理：

```go
if err := manager.InstallPackage(pkg); err != nil {
    switch pip.GetErrorType(err) {
    case pip.ErrorTypePackageNotFound:
        fmt.Printf("包 %s 未找到\n", pkg.Name)
    case pip.ErrorTypePermissionDenied:
        fmt.Println("权限被拒绝 - 尝试以提升的权限运行")
    case pip.ErrorTypeNetworkError:
        fmt.Println("网络错误 - 检查您的互联网连接")
    default:
        fmt.Printf("安装失败: %v\n", err)
    }
}
```

## 日志记录

启用日志记录以查看 SDK 正在做什么：

```go
// 创建日志记录器
logger, err := pip.NewLogger(&pip.LoggerConfig{
    Level:  pip.LogLevelInfo,
    Output: os.Stdout,
    Prefix: "[my-app]",
})
if err != nil {
    return err
}
defer logger.Close()

// 在管理器上设置日志记录器
manager := pip.NewManager(nil)
manager.SetCustomLogger(logger)

// 现在所有操作都将被记录
manager.InstallPackage(&pip.PackageSpec{Name: "requests"})
```

## 下一步

现在您已经掌握了基础知识，探索这些主题：

- [配置](/zh/guide/configuration) - 了解所有配置选项
- [包管理](/zh/guide/package-management) - 高级包操作
- [虚拟环境](/zh/guide/virtual-environments) - 使用虚拟环境
- [项目管理](/zh/guide/project-management) - 初始化 Python 项目
- [API 参考](/zh/api/) - 完整的 API 文档
- [示例](/zh/examples/) - 更多代码示例

## 获取帮助

如果遇到问题：

1. 查看 [API 文档](/zh/api/) 获取详细信息
2. 查看 [示例](/zh/examples/) 了解常见用例
3. 搜索 [问题跟踪器](https://github.com/scagogogo/go-pip-sdk/issues)
4. 在 [讨论](https://github.com/scagogogo/go-pip-sdk/discussions) 中提问
