---
layout: home

hero:
  name: "Go Pip SDK"
  text: "Go 语言中的 Python 包管理"
  tagline: 用于管理 Python pip 操作、虚拟环境和 Python 项目的综合 Go SDK
  image:
    src: /logo.svg
    alt: Go Pip SDK
  actions:
    - theme: brand
      text: 快速开始
      link: /zh/guide/getting-started
    - theme: alt
      text: 查看 GitHub
      link: https://github.com/scagogogo/go-pip-sdk

features:
  - icon: 🚀
    title: 跨平台支持
    details: 在 Windows、macOS 和 Linux 上无缝工作，具有自动平台检测和适配功能。
  
  - icon: 📦
    title: 完整的包管理
    details: 安装、卸载、列表、显示和冻结 Python 包，完全兼容 pip。
  
  - icon: 🐍
    title: 虚拟环境管理
    details: 轻松创建、激活、停用和删除虚拟环境。
  
  - icon: 🏗️
    title: 项目初始化
    details: 使用标准结构、setup.py、pyproject.toml 等引导 Python 项目。
  
  - icon: 🔧
    title: 自动 Pip 安装
    details: 如果缺失，自动检测并安装 pip，支持多种安装方法。
  
  - icon: 📝
    title: 全面的日志记录
    details: 详细的操作日志，具有多个级别和可自定义的输出格式。
  
  - icon: ⚡
    title: 丰富的错误处理
    details: 结构化错误类型，具有有用的建议和上下文感知的错误消息。
  
  - icon: 🧪
    title: 充分测试
    details: 广泛的单元和集成测试，82.3% 的代码覆盖率确保可靠性。
  
  - icon: 🔒
    title: 类型安全
    details: 完全类型化的 Go 接口，具有全面的文档和示例。
---

## 快速开始

安装 SDK 并开始在您的 Go 应用程序中管理 Python 包：

```bash
go get github.com/scagogogo/go-pip-sdk
```

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

## 为什么选择 Go Pip SDK？

- **原生 Go 集成**: 无需外部 Python 脚本或子进程调用
- **生产就绪**: 经过实战测试，具有全面的错误处理和日志记录
- **开发者友好**: 清晰的 API，具有广泛的文档和示例
- **灵活配置**: 针对不同环境和用例的可自定义设置

## 社区

- 📖 [文档](https://scagogogo.github.io/go-pip-sdk/zh/)
- 🐛 [问题跟踪](https://github.com/scagogogo/go-pip-sdk/issues)
- 💬 [讨论](https://github.com/scagogogo/go-pip-sdk/discussions)
- 📧 [贡献指南](/zh/guide/contributing)

## 许可证

基于 [MIT 许可证](https://github.com/scagogogo/go-pip-sdk/blob/main/LICENSE) 发布。
