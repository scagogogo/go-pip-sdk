# 安装

本指南将引导您完成 Go Pip SDK 的安装过程。

## 前提条件

在安装 Go Pip SDK 之前，请确保您的系统满足以下要求：

### Go 环境

- **Go 1.19 或更高版本**：SDK 需要 Go 1.19+ 才能正常工作
- 正确配置的 `GOPATH` 和 `GOROOT`
- 启用 Go 模块支持

检查您的 Go 版本：

```bash
go version
```

### Python 环境

- **Python 3.7 或更高版本**：SDK 管理 Python 包，因此需要 Python
- pip 包管理器（通常与 Python 一起安装）

检查您的 Python 版本：

```bash
python --version
# 或
python3 --version
```

检查 pip 是否可用：

```bash
pip --version
# 或
pip3 --version
```

## 安装方法

### 方法 1：使用 go get（推荐）

这是安装 Go Pip SDK 最简单的方法：

```bash
go get github.com/scagogogo/go-pip-sdk
```

### 方法 2：使用 Go 模块

在您的项目中，创建或更新 `go.mod` 文件：

```bash
go mod init your-project-name
go get github.com/scagogogo/go-pip-sdk
```

### 方法 3：从源码安装

如果您想从源码安装或贡献代码：

```bash
git clone https://github.com/scagogogo/go-pip-sdk.git
cd go-pip-sdk
go mod download
go install ./...
```

## 验证安装

创建一个简单的测试文件来验证安装：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    // 创建新的 pip 管理器
    manager := pip.NewManager(nil)
    
    // 检查 pip 是否已安装
    installed, err := manager.IsInstalled()
    if err != nil {
        log.Fatalf("检查 pip 安装失败: %v", err)
    }
    
    if installed {
        fmt.Println("✅ Go Pip SDK 安装成功！")
        
        // 获取 pip 版本
        version, err := manager.GetVersion()
        if err == nil {
            fmt.Printf("Pip 版本: %s\n", version)
        }
    } else {
        fmt.Println("⚠️  Pip 未安装，但 SDK 可以帮助您安装它")
    }
}
```

运行测试：

```bash
go run main.go
```

## 常见问题

### 问题 1：找不到 Python

如果您遇到 "python not found" 错误：

**解决方案：**
- 确保 Python 已安装并在 PATH 中
- 在 Windows 上，您可能需要使用 `py` 命令而不是 `python`
- 考虑使用自定义配置指定 Python 路径

```go
config := &pip.Config{
    PythonPath: "/usr/bin/python3", // 指定 Python 路径
}
manager := pip.NewManager(config)
```

### 问题 2：权限错误

如果您遇到权限相关的错误：

**解决方案：**
- 在 Unix 系统上使用 `sudo`（不推荐用于虚拟环境）
- 使用虚拟环境（推荐）
- 使用用户级安装

### 问题 3：网络连接问题

如果下载失败：

**解决方案：**
- 检查网络连接
- 配置代理设置
- 使用镜像源

```go
config := &pip.Config{
    DefaultIndex: "https://pypi.tuna.tsinghua.edu.cn/simple/", // 使用清华镜像
}
manager := pip.NewManager(config)
```

## 更新

要更新到最新版本：

```bash
go get -u github.com/scagogogo/go-pip-sdk
```

## 卸载

要从项目中移除 SDK：

```bash
go mod edit -droprequire github.com/scagogogo/go-pip-sdk
go mod tidy
```

## 下一步

安装完成后，您可以：

- 阅读[快速开始](./getting-started.md)指南
- 查看[配置](./configuration.md)选项
- 探索[示例](/zh/examples/)
- 查看[API 参考](/zh/api/)

## 获取帮助

如果您在安装过程中遇到问题：

- 查看[常见问题](./error-handling.md)
- 在 [GitHub Issues](https://github.com/scagogogo/go-pip-sdk/issues) 中搜索或报告问题
- 参与 [GitHub Discussions](https://github.com/scagogogo/go-pip-sdk/discussions)
