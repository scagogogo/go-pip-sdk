# 日志记录

Go Pip SDK 提供了强大而灵活的日志记录系统，帮助您监控操作、调试问题和跟踪性能。

## 日志系统概述

### 日志级别

SDK 支持以下日志级别（按严重程度排序）：

- **TRACE**：最详细的日志，包含所有执行细节
- **DEBUG**：调试信息，用于开发和故障排除
- **INFO**：一般信息，记录重要操作
- **WARN**：警告信息，表示潜在问题
- **ERROR**：错误信息，记录失败操作

### 默认日志配置

```go
// 使用默认日志配置
manager := pip.NewManager(nil)

// 默认配置等同于：
config := &pip.Config{
    LogLevel: "INFO",
}
manager = pip.NewManager(config)
```

## 配置日志记录

### 基本日志配置

```go
package main

import (
    "os"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    // 配置日志级别
    config := &pip.Config{
        LogLevel: "DEBUG", // 设置为调试级别
    }
    
    manager := pip.NewManager(config)
    
    // 现在所有操作都会产生详细的日志
    pkg := &pip.PackageSpec{Name: "requests"}
    if err := manager.InstallPackage(pkg); err != nil {
        log.Fatalf("安装失败: %v", err)
    }
}
```

### 自定义日志记录器

```go
// 创建自定义日志记录器
logger, err := pip.NewLogger(&pip.LoggerConfig{
    Level:      pip.LogLevelDebug,
    Output:     os.Stdout,
    Prefix:     "[MyApp]",
    EnableFile: true,
    LogFile:    "pip-operations.log",
    MaxSize:    10, // MB
    MaxBackups: 5,
    MaxAge:     30, // 天
})
if err != nil {
    log.Fatalf("创建日志记录器失败: %v", err)
}
defer logger.Close()

// 将自定义日志记录器设置到管理器
manager := pip.NewManager(nil)
manager.SetCustomLogger(logger)
```

## 日志输出格式

### 控制台输出

```go
// 配置控制台日志输出
logger, err := pip.NewLogger(&pip.LoggerConfig{
    Level:  pip.LogLevelInfo,
    Output: os.Stdout,
    Format: pip.LogFormatText, // 或 pip.LogFormatJSON
    Color:  true,              // 启用彩色输出
})
```

### 文件输出

```go
// 配置文件日志输出
logger, err := pip.NewLogger(&pip.LoggerConfig{
    Level:      pip.LogLevelDebug,
    EnableFile: true,
    LogFile:    "/var/log/pip-sdk.log",
    MaxSize:    50,  // 50MB
    MaxBackups: 10,  // 保留10个备份文件
    MaxAge:     30,  // 保留30天
    Compress:   true, // 压缩旧日志文件
})
```

### 同时输出到控制台和文件

```go
logger, err := pip.NewLogger(&pip.LoggerConfig{
    Level:      pip.LogLevelInfo,
    Output:     os.Stdout,     // 控制台输出
    EnableFile: true,
    LogFile:    "app.log",     // 文件输出
    Format:     pip.LogFormatJSON,
})
```

## 结构化日志记录

### JSON 格式日志

```go
// 启用 JSON 格式日志
logger, err := pip.NewLogger(&pip.LoggerConfig{
    Level:  pip.LogLevelInfo,
    Output: os.Stdout,
    Format: pip.LogFormatJSON,
    Fields: map[string]interface{}{
        "service":     "pip-manager",
        "version":     "1.0.0",
        "environment": "production",
    },
})

// JSON 日志输出示例：
// {"level":"info","msg":"Installing package","package":"requests","timestamp":"2024-01-15T10:30:00Z","service":"pip-manager"}
```

### 添加上下文字段

```go
// 为特定操作添加上下文
contextLogger := logger.WithFields(map[string]interface{}{
    "operation": "package_install",
    "user_id":   "12345",
    "session":   "abc-def-ghi",
})

manager.SetCustomLogger(contextLogger)
```

## 日志级别控制

### 动态调整日志级别

```go
// 运行时调整日志级别
manager.SetLogLevel(pip.LogLevelTrace) // 最详细
manager.SetLogLevel(pip.LogLevelDebug) // 调试
manager.SetLogLevel(pip.LogLevelInfo)  // 信息
manager.SetLogLevel(pip.LogLevelWarn)  // 警告
manager.SetLogLevel(pip.LogLevelError) // 仅错误
```

### 条件日志记录

```go
// 根据环境变量设置日志级别
logLevel := os.Getenv("LOG_LEVEL")
if logLevel == "" {
    logLevel = "INFO" // 默认级别
}

config := &pip.Config{
    LogLevel: logLevel,
}
manager := pip.NewManager(config)
```

## 日志过滤和采样

### 日志过滤

```go
// 创建带过滤器的日志记录器
logger, err := pip.NewLogger(&pip.LoggerConfig{
    Level:  pip.LogLevelDebug,
    Output: os.Stdout,
    Filter: func(entry *pip.LogEntry) bool {
        // 过滤掉某些包的日志
        if packageName, ok := entry.Fields["package"]; ok {
            if packageName == "noisy-package" {
                return false
            }
        }
        return true
    },
})
```

### 日志采样

```go
// 对高频日志进行采样
logger, err := pip.NewLogger(&pip.LoggerConfig{
    Level:      pip.LogLevelDebug,
    Output:     os.Stdout,
    SampleRate: 0.1, // 只记录10%的日志
})
```

## 性能监控日志

### 操作耗时记录

```go
// 启用性能监控
logger, err := pip.NewLogger(&pip.LoggerConfig{
    Level:           pip.LogLevelInfo,
    Output:          os.Stdout,
    EnableTiming:    true,
    EnableMetrics:   true,
})

manager := pip.NewManager(nil)
manager.SetCustomLogger(logger)

// 安装包时会自动记录耗时
pkg := &pip.PackageSpec{Name: "requests"}
if err := manager.InstallPackage(pkg); err != nil {
    log.Fatalf("安装失败: %v", err)
}

// 日志输出示例：
// INFO: Package installation completed package=requests duration=2.5s
```

### 自定义性能指标

```go
// 记录自定义指标
func installPackageWithMetrics(manager *pip.Manager, pkg *pip.PackageSpec) error {
    start := time.Now()
    
    err := manager.InstallPackage(pkg)
    
    duration := time.Since(start)
    
    // 记录性能指标
    manager.LogMetric("package_install_duration", duration.Seconds(), map[string]interface{}{
        "package": pkg.Name,
        "success": err == nil,
    })
    
    return err
}
```

## 日志轮转和归档

### 自动日志轮转

```go
logger, err := pip.NewLogger(&pip.LoggerConfig{
    Level:      pip.LogLevelInfo,
    EnableFile: true,
    LogFile:    "pip-sdk.log",
    
    // 轮转配置
    MaxSize:    100,  // 100MB 后轮转
    MaxBackups: 30,   // 保留30个备份文件
    MaxAge:     90,   // 保留90天
    Compress:   true, // 压缩旧文件
    
    // 轮转策略
    RotateOnStartup: true, // 启动时轮转
})
```

### 手动日志轮转

```go
// 手动触发日志轮转
if err := logger.Rotate(); err != nil {
    log.Printf("日志轮转失败: %v", err)
}
```

## 日志聚合和监控

### 发送日志到外部系统

```go
// 配置日志发送到 ELK Stack
logger, err := pip.NewLogger(&pip.LoggerConfig{
    Level:  pip.LogLevelInfo,
    Output: os.Stdout,
    Format: pip.LogFormatJSON,
    
    // 外部日志系统配置
    ExternalSinks: []pip.LogSink{
        &pip.ElasticsearchSink{
            URL:   "http://elasticsearch:9200",
            Index: "pip-sdk-logs",
        },
        &pip.SyslogSink{
            Network: "udp",
            Address: "syslog-server:514",
        },
    },
})
```

### 日志告警

```go
// 配置错误日志告警
logger, err := pip.NewLogger(&pip.LoggerConfig{
    Level:  pip.LogLevelInfo,
    Output: os.Stdout,
    
    // 告警配置
    AlertConfig: &pip.AlertConfig{
        ErrorThreshold: 10,        // 10个错误后告警
        TimeWindow:     time.Hour, // 1小时时间窗口
        WebhookURL:     "https://hooks.slack.com/...",
    },
})
```

## 调试和故障排除

### 启用详细调试日志

```go
// 开发环境：启用最详细的日志
if os.Getenv("ENV") == "development" {
    config := &pip.Config{
        LogLevel: "TRACE",
    }
    
    logger, _ := pip.NewLogger(&pip.LoggerConfig{
        Level:         pip.LogLevelTrace,
        Output:        os.Stdout,
        EnableCaller:  true, // 显示调用位置
        EnableStack:   true, // 显示堆栈跟踪
        PrettyPrint:   true, // 美化输出
    })
    
    manager := pip.NewManager(config)
    manager.SetCustomLogger(logger)
}
```

### 条件调试日志

```go
// 为特定包启用调试日志
func installPackageWithDebug(manager *pip.Manager, pkg *pip.PackageSpec) error {
    // 为特定包启用调试
    if pkg.Name == "problematic-package" {
        manager.SetLogLevel(pip.LogLevelTrace)
        defer manager.SetLogLevel(pip.LogLevelInfo)
    }
    
    return manager.InstallPackage(pkg)
}
```

## 日志分析

### 日志统计

```go
// 收集日志统计信息
type LogStats struct {
    TotalLogs    int64
    ErrorCount   int64
    WarnCount    int64
    InfoCount    int64
    DebugCount   int64
    LastError    time.Time
    mutex        sync.RWMutex
}

func (s *LogStats) RecordLog(level pip.LogLevel) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    s.TotalLogs++
    
    switch level {
    case pip.LogLevelError:
        s.ErrorCount++
        s.LastError = time.Now()
    case pip.LogLevelWarn:
        s.WarnCount++
    case pip.LogLevelInfo:
        s.InfoCount++
    case pip.LogLevelDebug:
        s.DebugCount++
    }
}
```

### 日志查询

```go
// 查询特定时间范围的日志
func queryLogs(logFile string, startTime, endTime time.Time) ([]pip.LogEntry, error) {
    file, err := os.Open(logFile)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    var entries []pip.LogEntry
    scanner := bufio.NewScanner(file)
    
    for scanner.Scan() {
        var entry pip.LogEntry
        if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
            continue
        }
        
        if entry.Timestamp.After(startTime) && entry.Timestamp.Before(endTime) {
            entries = append(entries, entry)
        }
    }
    
    return entries, scanner.Err()
}
```

## 最佳实践

### 1. 选择合适的日志级别

```go
// 生产环境
config := &pip.Config{
    LogLevel: "WARN", // 只记录警告和错误
}

// 开发环境
config := &pip.Config{
    LogLevel: "DEBUG", // 记录详细信息
}

// 故障排除
config := &pip.Config{
    LogLevel: "TRACE", // 记录所有信息
}
```

### 2. 使用结构化日志

```go
// 好的做法：使用结构化字段
logger.InfoWithFields("Package installed successfully", map[string]interface{}{
    "package": "requests",
    "version": "2.25.1",
    "duration": "2.3s",
})

// 不好的做法：在消息中嵌入变量
logger.Info(fmt.Sprintf("Package %s version %s installed in %s", 
    "requests", "2.25.1", "2.3s"))
```

### 3. 避免敏感信息泄露

```go
// 过滤敏感信息
logger, err := pip.NewLogger(&pip.LoggerConfig{
    Level:  pip.LogLevelInfo,
    Output: os.Stdout,
    Sanitizer: func(entry *pip.LogEntry) {
        // 移除敏感字段
        delete(entry.Fields, "password")
        delete(entry.Fields, "token")
        delete(entry.Fields, "api_key")
    },
})
```

### 4. 合理配置日志轮转

```go
// 根据应用规模配置轮转
logger, err := pip.NewLogger(&pip.LoggerConfig{
    Level:      pip.LogLevelInfo,
    EnableFile: true,
    LogFile:    "pip-sdk.log",
    
    // 高流量应用
    MaxSize:    50,   // 50MB
    MaxBackups: 20,   // 20个备份
    MaxAge:     7,    // 7天
    
    // 低流量应用
    // MaxSize:    10,   // 10MB
    // MaxBackups: 5,    // 5个备份
    // MaxAge:     30,   // 30天
})
```

## 下一步

- 查看[贡献指南](./contributing.md)
- 探索[日志记录示例](/zh/examples/advanced-usage.md)
- 了解[API 参考](/zh/api/logger.md)
- 学习[错误处理](./error-handling.md)
