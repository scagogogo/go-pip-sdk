# Go Pip SDK

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.19-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/go-pip-sdk)](https://goreportcard.com/report/github.com/scagogogo/go-pip-sdk)
[![Documentation](https://img.shields.io/badge/docs-online-blue.svg)](https://scagogogo.github.io/go-pip-sdk/zh/)

用于管理 Python pip 操作、虚拟环境和 Python 项目的综合 Go SDK。该库为所有常见的 pip 操作提供了清晰、类型安全的接口，并支持跨平台使用。

[English](README.md) | 简体中文

## 特性

- 🚀 **跨平台支持** - 在 Windows、macOS 和 Linux 上工作
- 📦 **完整的 pip 操作** - 安装、卸载、列表、显示、冻结包
- 🐍 **虚拟环境管理** - 创建、激活、停用、删除虚拟环境
- 🏗️ **项目初始化** - 使用标准结构引导 Python 项目
- 🔧 **自动 pip 安装** - 检测并在缺失时安装 pip
- 📝 **全面的日志记录** - 多级别的详细操作日志
- ⚡ **错误处理** - 丰富的错误类型和有用的建议
- 🧪 **充分测试** - 广泛的单元和集成测试

## 安装

```bash
go get github.com/scagogogo/go-pip-sdk
```

## 快速开始

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    // 创建新的 pip 管理器
    manager := pip.NewManager(nil) // 使用默认配置
    
    // 检查 pip 是否已安装
    installed, err := manager.IsInstalled()
    if err != nil {
        log.Fatal(err)
    }
    
    if !installed {
        fmt.Println("正在安装 pip...")
        if err := manager.Install(); err != nil {
            log.Fatal(err)
        }
    }
    
    // 安装包
    pkg := &pip.PackageSpec{
        Name:    "requests",
        Version: ">=2.25.0",
    }
    
    if err := manager.InstallPackage(pkg); err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("包安装成功！")
}
```

## 主要功能

### 包管理

```go
// 安装包
pkg := &pip.PackageSpec{
    Name:    "django",
    Version: ">=4.0,<5.0",
    Extras:  []string{"postgres", "redis"},
}
err := manager.InstallPackage(pkg)

// 卸载包
err = manager.UninstallPackage("requests")

// 列出已安装的包
packages, err := manager.ListPackages()

// 显示包信息
info, err := manager.ShowPackage("requests")

// 冻结包（类似 pip freeze）
packages, err := manager.FreezePackages()
```

### 虚拟环境管理

```go
// 创建虚拟环境
err := manager.CreateVenv("/path/to/venv")

// 激活虚拟环境
err = manager.ActivateVenv("/path/to/venv")

// 停用当前虚拟环境
err = manager.DeactivateVenv()

// 删除虚拟环境
err = manager.RemoveVenv("/path/to/venv")
```

### 项目初始化

```go
// 初始化新的 Python 项目
opts := &pip.ProjectOptions{
    Name:            "my-project",
    Version:         "0.1.0",
    Author:          "Your Name",
    AuthorEmail:     "your.email@example.com",
    Dependencies:    []string{"requests>=2.25.0"},
    DevDependencies: []string{"pytest>=6.0"},
    CreateVenv:      true,
}

err := manager.InitProject("/path/to/project", opts)
```

### 高级配置

```go
// 自定义配置
config := &pip.Config{
    PythonPath:   "/usr/bin/python3",
    PipPath:      "/usr/bin/pip3",
    Timeout:      60 * time.Second,
    Retries:      3,
    DefaultIndex: "https://pypi.org/simple/",
    LogLevel:     "DEBUG",
}

manager := pip.NewManager(config)

// 使用上下文进行超时控制
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

manager := pip.NewManagerWithContext(ctx, config)
```

## 文档

- 📖 **[在线文档](https://scagogogo.github.io/go-pip-sdk/zh/)** - 完整的 API 文档和指南
- 🚀 **[快速开始](https://scagogogo.github.io/go-pip-sdk/zh/guide/getting-started)** - 入门指南
- 📚 **[API 参考](https://scagogogo.github.io/go-pip-sdk/zh/api/)** - 详细的 API 文档
- 💡 **[示例](https://scagogogo.github.io/go-pip-sdk/zh/examples/)** - 代码示例和用例

## API 概览

### 核心类型

- `Manager` - pip 操作的主要接口
- `PackageSpec` - 包安装规范
- `Package` - 表示已安装的包
- `PackageInfo` - 详细的包信息
- `ProjectOptions` - 项目初始化选项
- `Config` - 管理器配置
- `Logger` - 自定义日志接口

### 错误类型

- `ErrorTypePipNotInstalled` - Pip 未安装
- `ErrorTypePythonNotFound` - 未找到 Python 解释器
- `ErrorTypePackageNotFound` - 包未找到
- `ErrorTypePermissionDenied` - 权限被拒绝
- `ErrorTypeNetworkError` - 网络连接问题
- `ErrorTypeCommandFailed` - 命令执行失败

## 测试

运行测试套件：

```bash
# 运行所有测试
go test ./...

# 运行带覆盖率的测试
go test -cover ./...

# 运行基准测试
go test -bench=. ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 贡献

我们欢迎贡献！请查看我们的[贡献指南](CONTRIBUTING.md)了解如何开始。

### 开发设置

1. 克隆仓库：
   ```bash
   git clone https://github.com/scagogogo/go-pip-sdk.git
   cd go-pip-sdk
   ```

2. 安装依赖：
   ```bash
   go mod download
   ```

3. 运行测试：
   ```bash
   make test
   ```

4. 运行 linter：
   ```bash
   make lint
   ```

## 许可证

该项目基于 [MIT 许可证](LICENSE) 开源。

## 致谢

- 受 Python pip 包管理器启发
- 使用 Go 优秀的标准库构建
- 感谢所有贡献者和用户

## 支持

- 📖 [文档](https://scagogogo.github.io/go-pip-sdk/zh/)
- 🐛 [问题跟踪](https://github.com/scagogogo/go-pip-sdk/issues)
- 💬 [讨论](https://github.com/scagogogo/go-pip-sdk/discussions)

## 相关项目

- [pip](https://pip.pypa.io/) - Python 包安装器
- [virtualenv](https://virtualenv.pypa.io/) - Python 虚拟环境工具
- [pipenv](https://pipenv.pypa.io/) - Python 开发工作流工具
