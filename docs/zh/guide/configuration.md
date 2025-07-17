# 配置

Go Pip SDK 提供了灵活的配置选项，允许您自定义其行为以满足您的特定需求。

## 配置概述

SDK 使用 `Config` 结构来管理所有配置选项。您可以在创建管理器时提供配置，或者使用默认设置。

```go
// 使用默认配置
manager := pip.NewManager(nil)

// 使用自定义配置
config := &pip.Config{
    PythonPath: "/usr/bin/python3",
    Timeout:    60 * time.Second,
}
manager := pip.NewManager(config)
```

## 配置选项

### 基本路径配置

#### PythonPath
指定 Python 可执行文件的路径。

```go
config := &pip.Config{
    PythonPath: "/usr/bin/python3", // Linux/macOS
    // PythonPath: "C:\\Python39\\python.exe", // Windows
}
```

**默认值：** 自动检测（`python3`、`python`）

#### PipPath
指定 pip 可执行文件的路径。

```go
config := &pip.Config{
    PipPath: "/usr/local/bin/pip3",
}
```

**默认值：** 自动检测（`pip3`、`pip`）

### 网络配置

#### DefaultIndex
设置默认的 PyPI 索引 URL。

```go
config := &pip.Config{
    DefaultIndex: "https://pypi.org/simple/",
    // 或使用镜像源
    // DefaultIndex: "https://pypi.tuna.tsinghua.edu.cn/simple/",
}
```

**默认值：** `https://pypi.org/simple/`

#### TrustedHosts
指定受信任的主机列表，用于 HTTPS 验证。

```go
config := &pip.Config{
    TrustedHosts: []string{
        "pypi.org",
        "pypi.tuna.tsinghua.edu.cn",
    },
}
```

#### Timeout
设置操作超时时间。

```go
config := &pip.Config{
    Timeout: 120 * time.Second, // 2 分钟超时
}
```

**默认值：** 60 秒

#### Retries
设置失败操作的重试次数。

```go
config := &pip.Config{
    Retries: 5, // 最多重试 5 次
}
```

**默认值：** 3 次

### 缓存配置

#### CacheDir
指定 pip 缓存目录。

```go
config := &pip.Config{
    CacheDir: "/tmp/pip-cache",
}
```

**默认值：** 系统默认缓存目录

### 日志配置

#### LogLevel
设置日志级别。

```go
config := &pip.Config{
    LogLevel: "DEBUG", // TRACE, DEBUG, INFO, WARN, ERROR
}
```

**默认值：** `INFO`

### 环境变量

#### Environment
设置额外的环境变量。

```go
config := &pip.Config{
    Environment: map[string]string{
        "PIP_NO_CACHE_DIR": "1",
        "PIP_DISABLE_PIP_VERSION_CHECK": "1",
        "PYTHONPATH": "/custom/python/path",
    },
}
```

### 额外选项

#### ExtraOptions
为 pip 命令添加额外的全局选项。

```go
config := &pip.Config{
    ExtraOptions: map[string]string{
        "no-cache-dir": "",
        "timeout": "120",
        "trusted-host": "pypi.org",
    },
}
```

## 完整配置示例

```go
package main

import (
    "time"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    config := &pip.Config{
        // 路径配置
        PythonPath: "/usr/bin/python3",
        PipPath:    "/usr/local/bin/pip3",
        
        // 网络配置
        DefaultIndex: "https://pypi.tuna.tsinghua.edu.cn/simple/",
        TrustedHosts: []string{
            "pypi.tuna.tsinghua.edu.cn",
            "pypi.org",
        },
        Timeout: 120 * time.Second,
        Retries: 5,
        
        // 缓存配置
        CacheDir: "/tmp/my-pip-cache",
        
        // 日志配置
        LogLevel: "DEBUG",
        
        // 环境变量
        Environment: map[string]string{
            "PIP_NO_CACHE_DIR": "0",
            "PIP_DISABLE_PIP_VERSION_CHECK": "1",
        },
        
        // 额外选项
        ExtraOptions: map[string]string{
            "timeout": "120",
            "retries": "5",
        },
    }
    
    manager := pip.NewManager(config)
    
    // 使用配置的管理器
    // ...
}
```

## 环境特定配置

### 开发环境

```go
func createDevConfig() *pip.Config {
    return &pip.Config{
        LogLevel: "DEBUG",
        Timeout:  30 * time.Second,
        Retries:  1,
        Environment: map[string]string{
            "PIP_DISABLE_PIP_VERSION_CHECK": "1",
        },
    }
}
```

### 生产环境

```go
func createProdConfig() *pip.Config {
    return &pip.Config{
        LogLevel: "WARN",
        Timeout:  300 * time.Second,
        Retries:  5,
        DefaultIndex: "https://pypi.org/simple/",
        CacheDir: "/var/cache/pip",
    }
}
```

### CI/CD 环境

```go
func createCIConfig() *pip.Config {
    return &pip.Config{
        LogLevel: "INFO",
        Timeout:  180 * time.Second,
        Retries:  3,
        Environment: map[string]string{
            "PIP_NO_CACHE_DIR": "1",
            "PIP_DISABLE_PIP_VERSION_CHECK": "1",
        },
        ExtraOptions: map[string]string{
            "no-cache-dir": "",
            "quiet": "",
        },
    }
}
```

## 配置验证

SDK 会自动验证配置并提供有用的错误消息：

```go
config := &pip.Config{
    PythonPath: "/invalid/path/python",
}

manager := pip.NewManager(config)
_, err := manager.IsInstalled()
if err != nil {
    // 错误消息会指出 Python 路径无效
    fmt.Printf("配置错误: %v\n", err)
}
```

## 动态配置更新

您可以在运行时更新某些配置选项：

```go
manager := pip.NewManager(nil)

// 更新超时设置
manager.SetTimeout(120 * time.Second)

// 更新重试次数
manager.SetRetries(5)

// 添加受信任的主机
manager.AddTrustedHost("custom.pypi.server.com")
```

## 配置最佳实践

1. **使用环境变量**：对于敏感信息，使用环境变量而不是硬编码
2. **分环境配置**：为不同环境创建不同的配置
3. **合理的超时**：根据网络条件设置适当的超时时间
4. **启用缓存**：在开发环境中启用缓存以提高性能
5. **日志级别**：在生产环境中使用较低的日志级别

## 下一步

- 了解[包管理](./package-management.md)
- 探索[虚拟环境](./virtual-environments.md)
- 查看[错误处理](./error-handling.md)
