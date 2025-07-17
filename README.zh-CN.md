# Go Pip SDK

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.19-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/go-pip-sdk)](https://goreportcard.com/report/github.com/scagogogo/go-pip-sdk)
[![Documentation](https://img.shields.io/badge/docs-online-blue.svg)](https://scagogogo.github.io/go-pip-sdk/zh/)
[![Build Status](https://img.shields.io/github/actions/workflow/status/scagogogo/go-pip-sdk/ci.yml?branch=main)](https://github.com/scagogogo/go-pip-sdk/actions)
[![Coverage Status](https://img.shields.io/codecov/c/github/scagogogo/go-pip-sdk)](https://codecov.io/gh/scagogogo/go-pip-sdk)

一个全面的、生产就绪的 Go SDK，用于管理 Python pip 操作、虚拟环境和 Python 项目。该库为所有常见的 pip 操作提供了清晰、类型安全的接口，具备企业级特性和跨平台支持。

[English](README.md) | **简体中文**

## ✨ 特性

### 🚀 核心功能
- **跨平台支持** - 在 Windows、macOS 和 Linux 上无缝工作
- **完整的 pip 操作** - 安装、卸载、列表、显示、冻结、搜索包
- **虚拟环境管理** - 创建、激活、停用、删除、克隆虚拟环境
- **项目初始化** - 使用可定制模板引导 Python 项目
- **自动 pip 安装** - 检测并使用多种安装方法安装 pip

### 🏢 企业级特性
- **生产就绪** - 在企业环境中经过实战测试
- **全面的日志记录** - 支持多种输出格式（JSON、文本）的结构化日志
- **高级错误处理** - 丰富的错误类型，提供可操作的建议和重试机制
- **配置管理** - 灵活的配置，支持环境变量
- **安全特性** - 证书验证、受信任主机和安全包安装

### 🛠️ 开发者体验
- **类型安全 API** - 完整的 Go 类型安全和全面的接口
- **广泛测试** - 95%+ 测试覆盖率，包含单元和集成测试
- **丰富文档** - 完整的 API 文档和示例
- **命令行界面** - 功能完整的 CLI 工具，可直接使用
- **Docker 支持** - 官方 Docker 镜像和容器化部署选项

## 📦 安装

### 使用 Go Modules（推荐）

```bash
go get github.com/scagogogo/go-pip-sdk
```

### 使用特定版本

```bash
go get github.com/scagogogo/go-pip-sdk@v1.0.0
```

### 系统要求

- **Go**: 1.19 或更高版本
- **Python**: 3.7 或更高版本（用于 pip 操作）
- **操作系统**: Windows 10+、macOS 10.15+ 或 Linux（任何现代发行版）

## 🚀 快速开始

### 基本用法

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

    // 检查 pip 是否已安装，如果缺失则安装
    if installed, err := manager.IsInstalled(); err != nil {
        log.Fatal(err)
    } else if !installed {
        fmt.Println("正在安装 pip...")
        if err := manager.Install(); err != nil {
            log.Fatal(err)
        }
        fmt.Println("✅ Pip 安装成功！")
    }

    // 使用版本约束安装包
    pkg := &pip.PackageSpec{
        Name:    "requests",
        Version: ">=2.25.0,<3.0.0",
        Extras:  []string{"security"}, // 安装额外依赖
    }

    fmt.Printf("正在安装 %s...\n", pkg.Name)
    if err := manager.InstallPackage(pkg); err != nil {
        log.Fatal(err)
    }

    fmt.Println("✅ 包安装成功！")

    // 列出已安装的包
    packages, err := manager.ListPackages()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("找到 %d 个已安装的包\n", len(packages))
}
```

### 使用自定义配置

```go
package main

import (
    "time"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    // 创建自定义配置
    config := &pip.Config{
        PythonPath:   "/usr/bin/python3",
        Timeout:      120 * time.Second,
        Retries:      5,
        LogLevel:     "INFO",
        DefaultIndex: "https://pypi.tuna.tsinghua.edu.cn/simple/",
        TrustedHosts: []string{"pypi.tuna.tsinghua.edu.cn"},
        Environment: map[string]string{
            "PIP_CACHE_DIR": "/tmp/pip-cache",
        },
    }

    manager := pip.NewManager(config)

    // 您的 pip 操作代码...
}
```

## 📚 核心功能

### 包管理

```go
// 使用各种选项安装包
pkg := &pip.PackageSpec{
    Name:           "fastapi",
    Version:        ">=0.68.0,<1.0.0",
    Extras:         []string{"all"},
    Upgrade:        true,
    ForceReinstall: false,
    UserInstall:    false,
}
err := manager.InstallPackage(pkg)

// 从需求文件安装
err = manager.InstallRequirements("requirements.txt")

// 从 Git 仓库安装
gitPkg := &pip.PackageSpec{
    Name: "git+https://github.com/user/repo.git@v1.0.0",
}
err = manager.InstallPackage(gitPkg)

// 卸载包
err = manager.UninstallPackage("requests")

// 列出已安装的包及详细信息
packages, err := manager.ListPackages()
for _, pkg := range packages {
    fmt.Printf("%s==%s (%s)\n", pkg.Name, pkg.Version, pkg.Location)
}

// 显示详细包信息
info, err := manager.ShowPackage("requests")
fmt.Printf("名称: %s\n版本: %s\n摘要: %s\n",
    info.Name, info.Version, info.Summary)

// 搜索包
results, err := manager.SearchPackages("web framework")

// 检查过时的包
outdated, err := manager.CheckOutdated()

// 冻结包（类似 pip freeze）
packages, err := manager.FreezePackages()
```

### 虚拟环境管理

```go
// 使用选项创建虚拟环境
opts := &pip.VenvOptions{
    PythonVersion:      "3.9",
    SystemSitePackages: false,
    Prompt:             "my-project",
    UpgradePip:         true,
}
err := manager.CreateVenvWithOptions("/path/to/venv", opts)

// 激活虚拟环境
err = manager.ActivateVenv("/path/to/venv")

// 检查虚拟环境是否激活
isActive, venvPath := manager.IsVenvActive()
if isActive {
    fmt.Printf("当前激活的虚拟环境: %s\n", venvPath)
}

// 列出所有虚拟环境
venvs, err := manager.ListVenvs("/path/to/envs")

// 获取详细的虚拟环境信息
info, err := manager.GetVenvInfo("/path/to/venv")
fmt.Printf("Python 版本: %s\n包数量: %d\n",
    info.PythonVersion, info.PackageCount)

// 克隆虚拟环境
err = manager.CloneVenv("/path/to/source", "/path/to/target")

// 删除虚拟环境
err = manager.RemoveVenv("/path/to/venv")
```

### 项目初始化

```go
// 初始化一个全面的 Python 项目
opts := &pip.ProjectOptions{
    Name:            "my-awesome-project",
    Version:         "0.1.0",
    Description:     "一个全面的 Python 项目",
    Author:          "您的姓名",
    AuthorEmail:     "your.email@example.com",
    License:         "MIT",
    Homepage:        "https://github.com/user/my-awesome-project",
    Repository:      "https://github.com/user/my-awesome-project.git",

    // 依赖
    Dependencies: []string{
        "fastapi>=0.68.0",
        "uvicorn[standard]>=0.15.0",
        "pydantic>=1.8.0",
    },
    DevDependencies: []string{
        "pytest>=6.0.0",
        "black>=21.0.0",
        "flake8>=3.8.0",
        "mypy>=0.812",
    },

    // 项目结构
    CreateVenv:          true,
    CreateSrc:           true,
    CreateTests:         true,
    CreateDocs:          true,
    CreateGithubActions: true,
    CreateDockerfile:    true,

    // 配置文件
    CreateSetupPy:       true,
    CreatePyprojectToml: true,
    CreateGitignore:     true,
    CreateReadme:        true,
}

err := manager.InitProject("/path/to/project", opts)

// 读取项目配置
config, err := manager.ReadProjectConfig("/path/to/project")

// 更新项目版本
err = manager.UpdateProjectVersion("/path/to/project", "1.0.0")

// 构建项目
buildOpts := &pip.BuildOptions{
    OutputDir: "./dist",
    Format:    "wheel",
    Clean:     true,
}
err = manager.BuildProject("/path/to/project", buildOpts)
```

## ⚙️ 配置

### 基本配置

```go
config := &pip.Config{
    // Python 设置
    PythonPath: "/usr/bin/python3",
    PipPath:    "/usr/bin/pip3",

    // 网络设置
    Timeout:      120 * time.Second,
    Retries:      5,
    DefaultIndex: "https://pypi.tuna.tsinghua.edu.cn/simple/",
    ExtraIndexes: []string{
        "https://pypi.org/simple/",
        "https://mirrors.aliyun.com/pypi/simple/",
    },
    TrustedHosts: []string{
        "pypi.tuna.tsinghua.edu.cn",
        "mirrors.aliyun.com",
    },

    // 缓存设置
    CacheDir: "/tmp/pip-cache",
    NoCache:  false,

    // 日志
    LogLevel: "INFO",
    LogFile:  "/var/log/pip-sdk.log",

    // 环境变量
    Environment: map[string]string{
        "PIP_CACHE_DIR":              "/tmp/pip-cache",
        "PIP_DISABLE_PIP_VERSION_CHECK": "1",
        "PIP_TIMEOUT":                "120",
    },
}

manager := pip.NewManager(config)
```

### 企业级配置

```go
// 具有安全特性的企业级配置
config := &pip.Config{
    PythonPath:   "/opt/python/bin/python3",
    DefaultIndex: "https://pypi.company.com/simple/",
    ExtraIndexes: []string{
        "https://pypi.tuna.tsinghua.edu.cn/simple/",
    },
    TrustedHosts: []string{
        "pypi.company.com",
        "pypi.tuna.tsinghua.edu.cn",
    },
    Timeout: 300 * time.Second,
    Retries: 10,

    // 安全设置
    ExtraOptions: map[string]string{
        "cert":         "/etc/ssl/certs/company-ca.pem",
        "client-cert":  "/etc/ssl/certs/client.pem",
        "trusted-host": "pypi.company.com",
    },

    // 审计日志
    LogLevel: "INFO",
    LogFile:  "/var/log/pip-operations.log",
}
```

## 🚀 高级用法

### 自定义日志

```go
// 创建具有多个输出的结构化日志记录器
loggerConfig := &pip.LoggerConfig{
    Level:      pip.LogLevelDebug,
    Format:     pip.LogFormatJSON,
    EnableFile: true,
    LogFile:    "/var/log/pip-sdk.log",
    MaxSize:    100, // 100MB
    MaxBackups: 5,
    MaxAge:     30, // 30天
    Compress:   true,

    // 为所有日志条目添加自定义字段
    Fields: map[string]interface{}{
        "service":     "pip-manager",
        "version":     "1.0.0",
        "environment": os.Getenv("ENVIRONMENT"),
    },
}

logger, err := pip.NewLogger(loggerConfig)
if err != nil {
    log.Fatal(err)
}
defer logger.Close()

// 设置自定义日志记录器
manager.SetCustomLogger(logger)
```

### 高级错误处理

```go
// 具有重试逻辑的全面错误处理
func installWithRetry(manager *pip.Manager, pkg *pip.PackageSpec, maxRetries int) error {
    var lastErr error

    for attempt := 0; attempt < maxRetries; attempt++ {
        err := manager.InstallPackage(pkg)
        if err == nil {
            return nil // 成功
        }

        lastErr = err

        // 处理不同的错误类型
        switch pip.GetErrorType(err) {
        case pip.ErrorTypeNetworkError, pip.ErrorTypeNetworkTimeout:
            // 使用指数退避重试网络错误
            delay := time.Duration(1<<uint(attempt)) * time.Second
            fmt.Printf("网络错误，%v 后重试... (%d/%d)\n", delay, attempt+1, maxRetries)
            time.Sleep(delay)
            continue

        case pip.ErrorTypePermissionDenied:
            // 对权限错误尝试用户级安装
            if attempt == 0 {
                pkg.UserInstall = true
                continue
            }
            return err

        case pip.ErrorTypePackageNotFound:
            // 为缺失的包建议替代方案
            if results, searchErr := manager.SearchPackages(pkg.Name); searchErr == nil && len(results) > 0 {
                fmt.Printf("包 '%s' 未找到。相似包:\n", pkg.Name)
                for i, result := range results[:3] {
                    fmt.Printf("%d. %s - %s\n", i+1, result.Name, result.Summary)
                }
            }
            return err

        default:
            return err // 不重试其他错误
        }
    }

    return fmt.Errorf("重试 %d 次后失败: %w", maxRetries, lastErr)
}
```

### 上下文支持和取消

```go
// 使用上下文进行超时和取消
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
defer cancel()

manager := pip.NewManagerWithContext(ctx, config)

// 操作将遵循上下文取消
err := manager.InstallPackage(pkg)
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        fmt.Println("操作超时")
    } else if errors.Is(err, context.Canceled) {
        fmt.Println("操作被取消")
    }
}
```

## 📖 文档

- 📖 **[在线文档](https://scagogogo.github.io/go-pip-sdk/zh/)** - 完整的 API 文档和指南
- 🚀 **[快速开始](https://scagogogo.github.io/go-pip-sdk/zh/guide/getting-started)** - 入门指南和安装说明
- 📚 **[API 参考](https://scagogogo.github.io/go-pip-sdk/zh/api/)** - 详细的 API 文档和示例
- 💡 **[示例](https://scagogogo.github.io/go-pip-sdk/zh/examples/)** - 全面的代码示例和用例
- 🔧 **[配置指南](https://scagogogo.github.io/go-pip-sdk/zh/guide/configuration)** - 配置选项和最佳实践
- 🐛 **[故障排除](https://scagogogo.github.io/go-pip-sdk/zh/guide/troubleshooting)** - 常见问题和解决方案

## 🔍 API 概览

### 核心接口

```go
// 主要的 pip 管理器接口
type Manager interface {
    // 包操作
    InstallPackage(pkg *PackageSpec) error
    UninstallPackage(name string) error
    ListPackages() ([]*Package, error)
    ShowPackage(name string) (*PackageInfo, error)
    SearchPackages(query string) ([]*SearchResult, error)

    // 虚拟环境操作
    CreateVenv(path string) error
    ActivateVenv(path string) error
    DeactivateVenv() error

    // 项目操作
    InitProject(path string, opts *ProjectOptions) error
    BuildProject(path string, opts *BuildOptions) error

    // 系统操作
    IsInstalled() (bool, error)
    GetVersion() (string, error)
}
```

### 关键数据类型

```go
// 包安装规范
type PackageSpec struct {
    Name           string   // 包名
    Version        string   // 版本约束
    Extras         []string // 额外依赖
    Upgrade        bool     // 如果已安装则升级
    ForceReinstall bool     // 强制重新安装
    UserInstall    bool     // 安装到用户目录
    Editable       bool     // 可编辑安装
}

// 项目初始化选项
type ProjectOptions struct {
    Name            string   // 项目名称
    Version         string   // 初始版本
    Description     string   // 项目描述
    Author          string   // 作者姓名
    AuthorEmail     string   // 作者邮箱
    License         string   // 许可证类型
    Dependencies    []string // 运行时依赖
    DevDependencies []string // 开发依赖
    CreateVenv      bool     // 创建虚拟环境
    CreateSrc       bool     // 创建 src/ 目录
    CreateTests     bool     // 创建 tests/ 目录
}

// 配置选项
type Config struct {
    PythonPath   string        // Python 可执行文件路径
    Timeout      time.Duration // 操作超时时间
    Retries      int           // 重试次数
    DefaultIndex string        // 默认包索引
    TrustedHosts []string      // 受信任的主机
    LogLevel     string        // 日志级别
    Environment  map[string]string // 环境变量
}
```

### 错误处理

```go
// 不同失败场景的错误类型
const (
    ErrorTypePipNotInstalled    ErrorType = "pip_not_installed"
    ErrorTypePythonNotFound     ErrorType = "python_not_found"
    ErrorTypePackageNotFound    ErrorType = "package_not_found"
    ErrorTypePermissionDenied   ErrorType = "permission_denied"
    ErrorTypeNetworkError       ErrorType = "network_error"
    ErrorTypeVersionConflict    ErrorType = "version_conflict"
    ErrorTypeCommandFailed      ErrorType = "command_failed"
)

// 检查错误类型
if pip.IsErrorType(err, pip.ErrorTypeNetworkError) {
    // 专门处理网络错误
}
```

## 🧪 测试

运行全面的测试套件：

```bash
# 运行所有测试
go test ./...

# 运行带覆盖率报告的测试
go test -cover ./...

# 生成详细的覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 仅运行单元测试（跳过集成测试）
go test -short ./...

# 运行集成测试（需要 Python 和 pip）
go test -run Integration ./...

# 运行基准测试
go test -bench=. ./...

# 运行带竞态检测的测试
go test -race ./...
```

### 测试分类

- **单元测试**：不需要外部依赖的快速测试
- **集成测试**：需要 Python 和 pip 安装的测试
- **基准测试**：关键操作的性能测试
- **示例测试**：确保文档示例正常工作

## 🤝 贡献

我们欢迎贡献！请查看我们的[贡献指南](CONTRIBUTING.md)了解详细信息。

### 贡献者快速开始

1. **Fork 仓库**
   ```bash
   git clone https://github.com/your-username/go-pip-sdk.git
   cd go-pip-sdk
   ```

2. **设置开发环境**
   ```bash
   # 安装依赖
   go mod download

   # 安装开发工具
   make install-tools

   # 运行测试确保一切正常
   make test
   ```

3. **创建功能分支**
   ```bash
   git checkout -b feature/amazing-feature
   ```

4. **进行更改并测试**
   ```bash
   # 运行测试
   make test

   # 运行代码检查
   make lint

   # 格式化代码
   make fmt
   ```

5. **提交并推送**
   ```bash
   git commit -m 'feat: 添加惊人的功能'
   git push origin feature/amazing-feature
   ```

6. **创建 Pull Request**

### 开发指南

- 遵循 [Go 代码审查注释](https://github.com/golang/go/wiki/CodeReviewComments)
- 为新功能编写测试
- 为 API 更改更新文档
- 使用约定式提交消息
- 确保所有 CI 检查通过

## 📋 系统要求

### 运行时要求
- **Go**: 1.19 或更高版本
- **Python**: 3.7 或更高版本（用于 pip 操作）
- **操作系统**: Windows 10+、macOS 10.15+ 或 Linux

### 开发要求
- **Go**: 1.19 或更高版本
- **Make**: 用于构建自动化
- **Git**: 用于版本控制
- **Python**: 3.7+ 带 pip（用于集成测试）

## 📄 许可证

该项目基于 MIT 许可证开源 - 查看 [LICENSE](LICENSE) 文件了解详情。

### 第三方许可证

该项目使用了几个第三方库。查看 [THIRD_PARTY_LICENSES.md](THIRD_PARTY_LICENSES.md) 了解详情。

## 🙏 致谢

- **Python pip 团队** - 创建了优秀的 pip 包管理器，启发了这个项目
- **Go 团队** - 提供了出色的编程语言和标准库
- **贡献者** - 感谢所有为这个项目做出贡献的开发者
- **社区** - 特别感谢提供反馈和报告问题的用户

## 📞 支持

### 获取帮助

- 📖 **[文档](https://scagogogo.github.io/go-pip-sdk/zh/)** - 全面的指南和 API 参考
- 🐛 **[问题跟踪](https://github.com/scagogogo/go-pip-sdk/issues)** - 报告错误或请求功能
- 💬 **[讨论](https://github.com/scagogogo/go-pip-sdk/discussions)** - 提问和分享想法
- 📧 **[邮件](mailto:support@scagogogo.com)** - 企业用户直接支持

### 企业支持

对于企业用户，我们提供：
- 优先支持和错误修复
- 自定义功能开发
- 培训和咨询服务
- SLA 支持协议

联系 [enterprise@scagogogo.com](mailto:enterprise@scagogogo.com) 了解更多信息。

## 🔗 相关项目

- [pip](https://pip.pypa.io/) - Python 包安装器
- [virtualenv](https://virtualenv.pypa.io/) - Python 虚拟环境工具
- [pipenv](https://pipenv.pypa.io/) - Python 开发工作流工具
- [poetry](https://python-poetry.org/) - Python 依赖管理和打包工具

---

<div align="center">

**[⬆ 回到顶部](#go-pip-sdk)**

由 Go Pip SDK 团队用 ❤️ 制作

</div>
