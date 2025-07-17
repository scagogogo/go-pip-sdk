# 日志记录 API

日志记录 API 提供了完整的日志系统功能，包括日志记录器创建、配置和管理。

## 核心接口

### Logger

日志记录器接口。

```go
type Logger interface {
    // 基本日志方法
    Trace(msg string, fields ...map[string]interface{})
    Debug(msg string, fields ...map[string]interface{})
    Info(msg string, fields ...map[string]interface{})
    Warn(msg string, fields ...map[string]interface{})
    Error(msg string, fields ...map[string]interface{})
    
    // 带字段的日志方法
    TraceWithFields(msg string, fields map[string]interface{})
    DebugWithFields(msg string, fields map[string]interface{})
    InfoWithFields(msg string, fields map[string]interface{})
    WarnWithFields(msg string, fields map[string]interface{})
    ErrorWithFields(msg string, fields map[string]interface{})
    
    // 上下文方法
    WithFields(fields map[string]interface{}) Logger
    WithField(key string, value interface{}) Logger
    
    // 控制方法
    SetLevel(level LogLevel)
    GetLevel() LogLevel
    Close() error
}
```

## 创建日志记录器

### NewLogger

创建新的日志记录器。

```go
func NewLogger(config *LoggerConfig) (Logger, error)
```

**参数：**
- `config` - 日志记录器配置

**返回值：**
- `Logger` - 日志记录器实例
- `error` - 错误信息

**示例：**

```go
config := &pip.LoggerConfig{
    Level:  pip.LogLevelInfo,
    Output: os.Stdout,
    Format: pip.LogFormatText,
    Color:  true,
}

logger, err := pip.NewLogger(config)
if err != nil {
    log.Fatalf("创建日志记录器失败: %v", err)
}
defer logger.Close()
```

### NewFileLogger

创建文件日志记录器。

```go
func NewFileLogger(filename string, config *LoggerConfig) (Logger, error)
```

**参数：**
- `filename` - 日志文件名
- `config` - 日志记录器配置

**返回值：**
- `Logger` - 日志记录器实例
- `error` - 错误信息

**示例：**

```go
logger, err := pip.NewFileLogger("app.log", &pip.LoggerConfig{
    Level:      pip.LogLevelDebug,
    Format:     pip.LogFormatJSON,
    MaxSize:    10, // 10MB
    MaxBackups: 5,
    MaxAge:     30, // 30天
    Compress:   true,
})
```

### NewMultiLogger

创建多输出日志记录器。

```go
func NewMultiLogger(loggers ...Logger) Logger
```

**参数：**
- `loggers` - 日志记录器列表

**返回值：**
- `Logger` - 多输出日志记录器

**示例：**

```go
consoleLogger, _ := pip.NewLogger(&pip.LoggerConfig{
    Level:  pip.LogLevelInfo,
    Output: os.Stdout,
    Color:  true,
})

fileLogger, _ := pip.NewFileLogger("app.log", &pip.LoggerConfig{
    Level:  pip.LogLevelDebug,
    Format: pip.LogFormatJSON,
})

multiLogger := pip.NewMultiLogger(consoleLogger, fileLogger)
```

## 日志级别

### LogLevel

日志级别枚举。

```go
type LogLevel int

const (
    LogLevelTrace LogLevel = iota
    LogLevelDebug
    LogLevelInfo
    LogLevelWarn
    LogLevelError
)
```

### 级别转换函数

```go
func ParseLogLevel(level string) (LogLevel, error)
func (l LogLevel) String() string
```

**示例：**

```go
// 从字符串解析日志级别
level, err := pip.ParseLogLevel("DEBUG")
if err != nil {
    log.Fatalf("无效的日志级别: %v", err)
}

// 转换为字符串
levelStr := level.String() // "DEBUG"
```

## 日志格式

### LogFormat

日志格式枚举。

```go
type LogFormat int

const (
    LogFormatText LogFormat = iota
    LogFormatJSON
    LogFormatCustom
)
```

### 自定义格式化器

```go
type Formatter interface {
    Format(entry *LogEntry) ([]byte, error)
}

type CustomFormatter struct {
    Template string
}

func (f *CustomFormatter) Format(entry *LogEntry) ([]byte, error)
```

**示例：**

```go
formatter := &pip.CustomFormatter{
    Template: "[{{.Level}}] {{.Timestamp}} - {{.Message}} {{.Fields}}",
}

config := &pip.LoggerConfig{
    Level:     pip.LogLevelInfo,
    Output:    os.Stdout,
    Format:    pip.LogFormatCustom,
    Formatter: formatter,
}
```

## 配置选项

### LoggerConfig

日志记录器配置结构。

```go
type LoggerConfig struct {
    // 基本配置
    Level  LogLevel  // 日志级别
    Output io.Writer // 输出目标
    Format LogFormat // 日志格式
    
    // 文件配置
    EnableFile bool   // 启用文件输出
    LogFile    string // 日志文件路径
    MaxSize    int    // 最大文件大小 (MB)
    MaxBackups int    // 最大备份文件数
    MaxAge     int    // 最大保留天数
    Compress   bool   // 是否压缩旧文件
    
    // 显示配置
    Color        bool   // 启用彩色输出
    PrettyPrint  bool   // 美化输出
    EnableCaller bool   // 显示调用者信息
    EnableStack  bool   // 显示堆栈跟踪
    
    // 字段配置
    Fields    map[string]interface{} // 默认字段
    Prefix    string                 // 日志前缀
    Formatter Formatter              // 自定义格式化器
    
    // 过滤配置
    Filter     func(*LogEntry) bool // 日志过滤器
    SampleRate float64              // 采样率 (0.0-1.0)
    
    // 外部输出
    ExternalSinks []LogSink // 外部日志接收器
}
```

## 日志方法

### 基本日志记录

```go
// 基本用法
logger.Info("应用程序启动")
logger.Error("操作失败")

// 带字段的日志
logger.InfoWithFields("包安装成功", map[string]interface{}{
    "package": "requests",
    "version": "2.25.1",
    "duration": "2.3s",
})

// 使用上下文
contextLogger := logger.WithFields(map[string]interface{}{
    "operation": "package_install",
    "user_id":   "12345",
})
contextLogger.Info("开始安装包")
contextLogger.Error("安装失败")
```

### 条件日志记录

```go
// 只在调试模式下记录
if logger.GetLevel() <= pip.LogLevelDebug {
    logger.Debug("详细调试信息")
}

// 使用延迟计算避免性能开销
logger.DebugWithFields("复杂计算结果", func() map[string]interface{} {
    return map[string]interface{}{
        "result": expensiveCalculation(),
    }
}())
```

## 高级功能

### 日志过滤

```go
config := &pip.LoggerConfig{
    Level:  pip.LogLevelDebug,
    Output: os.Stdout,
    Filter: func(entry *pip.LogEntry) bool {
        // 过滤掉包含敏感信息的日志
        if strings.Contains(entry.Message, "password") {
            return false
        }
        return true
    },
}
```

### 日志采样

```go
config := &pip.LoggerConfig{
    Level:      pip.LogLevelInfo,
    Output:     os.Stdout,
    SampleRate: 0.1, // 只记录10%的日志
}
```

### 日志轮转

```go
config := &pip.LoggerConfig{
    Level:      pip.LogLevelInfo,
    EnableFile: true,
    LogFile:    "app.log",
    MaxSize:    100,  // 100MB
    MaxBackups: 10,   // 10个备份文件
    MaxAge:     30,   // 30天
    Compress:   true, // 压缩旧文件
}
```

## 外部日志系统

### LogSink

外部日志接收器接口。

```go
type LogSink interface {
    Write(entry *LogEntry) error
    Close() error
}
```

### 内置接收器

#### ElasticsearchSink

```go
type ElasticsearchSink struct {
    URL   string
    Index string
    Auth  *ElasticsearchAuth
}

func NewElasticsearchSink(url, index string) *ElasticsearchSink
```

#### SyslogSink

```go
type SyslogSink struct {
    Network string
    Address string
    Tag     string
}

func NewSyslogSink(network, address, tag string) *SyslogSink
```

#### WebhookSink

```go
type WebhookSink struct {
    URL     string
    Headers map[string]string
    Timeout time.Duration
}

func NewWebhookSink(url string) *WebhookSink
```

**示例：**

```go
config := &pip.LoggerConfig{
    Level:  pip.LogLevelInfo,
    Output: os.Stdout,
    ExternalSinks: []pip.LogSink{
        pip.NewElasticsearchSink("http://localhost:9200", "pip-logs"),
        pip.NewSyslogSink("udp", "localhost:514", "pip-sdk"),
        pip.NewWebhookSink("https://hooks.slack.com/services/..."),
    },
}
```

## 性能监控

### 性能日志记录

```go
// 记录操作耗时
func logOperationDuration(logger pip.Logger, operation string, fn func() error) error {
    start := time.Now()
    err := fn()
    duration := time.Since(start)
    
    fields := map[string]interface{}{
        "operation": operation,
        "duration":  duration.String(),
        "success":   err == nil,
    }
    
    if err != nil {
        fields["error"] = err.Error()
        logger.ErrorWithFields("操作失败", fields)
    } else {
        logger.InfoWithFields("操作完成", fields)
    }
    
    return err
}
```

### 指标记录

```go
// 记录自定义指标
func (l *Logger) LogMetric(name string, value float64, tags map[string]interface{}) {
    l.InfoWithFields("metric", map[string]interface{}{
        "metric_name":  name,
        "metric_value": value,
        "tags":         tags,
        "timestamp":    time.Now().Unix(),
    })
}
```

## 日志分析

### 日志查询

```go
type LogQuery struct {
    Level     LogLevel
    StartTime time.Time
    EndTime   time.Time
    Fields    map[string]interface{}
    Message   string
}

func QueryLogs(logFile string, query *LogQuery) ([]*LogEntry, error)
```

### 日志统计

```go
type LogStats struct {
    TotalEntries int64
    ByLevel      map[LogLevel]int64
    ErrorRate    float64
    LastEntry    time.Time
}

func AnalyzeLogs(logFile string) (*LogStats, error)
```

## 最佳实践

### 1. 结构化日志

```go
// 好的做法：使用结构化字段
logger.InfoWithFields("用户登录", map[string]interface{}{
    "user_id":    "12345",
    "ip_address": "192.168.1.1",
    "user_agent": "Mozilla/5.0...",
})

// 不好的做法：在消息中嵌入变量
logger.Info(fmt.Sprintf("用户 %s 从 %s 登录", userID, ipAddress))
```

### 2. 适当的日志级别

```go
// 生产环境
logger.SetLevel(pip.LogLevelWarn)

// 开发环境
logger.SetLevel(pip.LogLevelDebug)

// 故障排除
logger.SetLevel(pip.LogLevelTrace)
```

### 3. 敏感信息处理

```go
// 过滤敏感信息
config := &pip.LoggerConfig{
    Filter: func(entry *pip.LogEntry) bool {
        // 移除密码字段
        delete(entry.Fields, "password")
        delete(entry.Fields, "token")
        return true
    },
}
```

### 4. 性能考虑

```go
// 避免昂贵的字符串操作
if logger.GetLevel() <= pip.LogLevelDebug {
    logger.Debug(expensiveStringOperation())
}

// 使用延迟计算
logger.DebugWithFields("调试信息", func() map[string]interface{} {
    return map[string]interface{}{
        "data": expensiveDataGeneration(),
    }
}())
```

## 下一步

- 查看[错误处理 API](./errors.md)
- 了解[类型定义](./types.md)
- 探索[安装器 API](./installer.md)
- 学习[管理器 API](./manager.md)
