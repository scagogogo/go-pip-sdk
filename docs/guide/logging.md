# Logging

This guide covers how to configure and use logging with the Go Pip SDK.

## Overview

The Go Pip SDK includes a comprehensive logging system that helps you debug issues, monitor operations, and track the execution of pip commands.

## Basic Logging Setup

### Default Logger

```go
package main

import (
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    // Create manager with default logging
    manager := pip.NewManager(nil)
    
    // Operations will be logged with default settings
    pkg := &pip.PackageSpec{Name: "requests"}
    err := manager.InstallPackage(pkg)
    if err != nil {
        // Error will be automatically logged
        return
    }
}
```

### Custom Logger Configuration

```go
config := &pip.Config{
    Logger: &pip.LoggerConfig{
        Level:      pip.LogLevelDebug,
        Output:     os.Stdout,
        Prefix:     "[pip-sdk]",
        TimeFormat: "2006-01-02 15:04:05",
        Colors:     true,
    },
}

manager := pip.NewManager(config)
```

## Log Levels

The SDK supports different log levels:

```go
const (
    LogLevelDebug = "DEBUG"
    LogLevelInfo  = "INFO"
    LogLevelWarn  = "WARN"
    LogLevelError = "ERROR"
    LogLevelFatal = "FATAL"
)
```

### Setting Log Level

```go
// Only show warnings and errors
config := &pip.Config{
    Logger: &pip.LoggerConfig{
        Level: pip.LogLevelWarn,
    },
}

// Show all debug information
config := &pip.Config{
    Logger: &pip.LoggerConfig{
        Level: pip.LogLevelDebug,
    },
}
```

## Logging to Files

### Single Log File

```go
logFile, err := os.OpenFile("pip-sdk.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
if err != nil {
    log.Fatal(err)
}
defer logFile.Close()

config := &pip.Config{
    Logger: &pip.LoggerConfig{
        Level:  pip.LogLevelInfo,
        Output: logFile,
    },
}

manager := pip.NewManager(config)
```

### Multiple Outputs

```go
import (
    "io"
    "os"
)

// Log to both file and console
logFile, err := os.OpenFile("pip-sdk.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
if err != nil {
    log.Fatal(err)
}
defer logFile.Close()

multiWriter := io.MultiWriter(os.Stdout, logFile)

config := &pip.Config{
    Logger: &pip.LoggerConfig{
        Level:  pip.LogLevelInfo,
        Output: multiWriter,
        Colors: false, // Disable colors for file output
    },
}
```

## Structured Logging

### JSON Logging

```go
import (
    "encoding/json"
    "log/slog"
    "os"
)

// Create a structured logger
jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
})
logger := slog.New(jsonHandler)

// Use with custom logging wrapper
type StructuredLogger struct {
    logger *slog.Logger
}

func (sl *StructuredLogger) Debug(msg string, args ...interface{}) {
    sl.logger.Debug(msg, args...)
}

func (sl *StructuredLogger) Info(msg string, args ...interface{}) {
    sl.logger.Info(msg, args...)
}

func (sl *StructuredLogger) Warn(msg string, args ...interface{}) {
    sl.logger.Warn(msg, args...)
}

func (sl *StructuredLogger) Error(msg string, args ...interface{}) {
    sl.logger.Error(msg, args...)
}

// Set custom logger
manager.SetLogger(&StructuredLogger{logger: logger})
```

## Logging Operations

### Package Operations

```go
// Enable debug logging to see detailed pip commands
config := &pip.Config{
    Logger: &pip.LoggerConfig{
        Level: pip.LogLevelDebug,
    },
}

manager := pip.NewManager(config)

// This will log:
// [DEBUG] Installing package: requests
// [DEBUG] Executing command: pip install requests
// [DEBUG] Command output: Successfully installed requests-2.28.1
pkg := &pip.PackageSpec{Name: "requests"}
err := manager.InstallPackage(pkg)
```

### Virtual Environment Operations

```go
venvManager := pip.NewVenvManager(manager)

// This will log virtual environment operations
err := venvManager.CreateVenv("./myenv")
// [INFO] Creating virtual environment: ./myenv
// [DEBUG] Executing command: python -m venv ./myenv
// [INFO] Successfully created virtual environment using venv module

err = venvManager.ActivateVenv("./myenv")
// [INFO] Activating virtual environment: ./myenv
// [INFO] Virtual environment activated: ./myenv
```

## Custom Logging Integration

### Logrus Integration

```go
import (
    "github.com/sirupsen/logrus"
)

type LogrusAdapter struct {
    logger *logrus.Logger
}

func NewLogrusAdapter() *LogrusAdapter {
    logger := logrus.New()
    logger.SetLevel(logrus.DebugLevel)
    logger.SetFormatter(&logrus.JSONFormatter{})
    
    return &LogrusAdapter{logger: logger}
}

func (la *LogrusAdapter) Debug(msg string, args ...interface{}) {
    la.logger.Debug(fmt.Sprintf(msg, args...))
}

func (la *LogrusAdapter) Info(msg string, args ...interface{}) {
    la.logger.Info(fmt.Sprintf(msg, args...))
}

func (la *LogrusAdapter) Warn(msg string, args ...interface{}) {
    la.logger.Warn(fmt.Sprintf(msg, args...))
}

func (la *LogrusAdapter) Error(msg string, args ...interface{}) {
    la.logger.Error(fmt.Sprintf(msg, args...))
}

// Use with manager
manager := pip.NewManager(nil)
manager.SetLogger(NewLogrusAdapter())
```

### Zap Integration

```go
import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

type ZapAdapter struct {
    logger *zap.SugaredLogger
}

func NewZapAdapter() *ZapAdapter {
    config := zap.NewDevelopmentConfig()
    config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
    
    logger, _ := config.Build()
    sugar := logger.Sugar()
    
    return &ZapAdapter{logger: sugar}
}

func (za *ZapAdapter) Debug(msg string, args ...interface{}) {
    za.logger.Debugf(msg, args...)
}

func (za *ZapAdapter) Info(msg string, args ...interface{}) {
    za.logger.Infof(msg, args...)
}

func (za *ZapAdapter) Warn(msg string, args ...interface{}) {
    za.logger.Warnf(msg, args...)
}

func (za *ZapAdapter) Error(msg string, args ...interface{}) {
    za.logger.Errorf(msg, args...)
}
```

## Logging Best Practices

### Contextual Logging

```go
type ContextualLogger struct {
    logger pip.Logger
    context map[string]interface{}
}

func NewContextualLogger(logger pip.Logger) *ContextualLogger {
    return &ContextualLogger{
        logger:  logger,
        context: make(map[string]interface{}),
    }
}

func (cl *ContextualLogger) WithContext(key string, value interface{}) *ContextualLogger {
    newContext := make(map[string]interface{})
    for k, v := range cl.context {
        newContext[k] = v
    }
    newContext[key] = value
    
    return &ContextualLogger{
        logger:  cl.logger,
        context: newContext,
    }
}

func (cl *ContextualLogger) Info(msg string, args ...interface{}) {
    contextMsg := fmt.Sprintf("[%s] %s", cl.formatContext(), msg)
    cl.logger.Info(contextMsg, args...)
}

func (cl *ContextualLogger) formatContext() string {
    var parts []string
    for k, v := range cl.context {
        parts = append(parts, fmt.Sprintf("%s=%v", k, v))
    }
    return strings.Join(parts, " ")
}

// Usage
logger := NewContextualLogger(manager.GetLogger())
projectLogger := logger.WithContext("project", "myapp").WithContext("env", "development")
projectLogger.Info("Installing dependencies")
```

## Performance Logging

### Operation Timing

```go
type TimingLogger struct {
    logger pip.Logger
}

func (tl *TimingLogger) LogOperation(operation string, fn func() error) error {
    start := time.Now()
    tl.logger.Info("Starting operation: %s", operation)
    
    err := fn()
    duration := time.Since(start)
    
    if err != nil {
        tl.logger.Error("Operation failed: %s (took %v): %v", operation, duration, err)
    } else {
        tl.logger.Info("Operation completed: %s (took %v)", operation, duration)
    }
    
    return err
}

// Usage
timingLogger := &TimingLogger{logger: manager.GetLogger()}

err := timingLogger.LogOperation("install requests", func() error {
    return manager.InstallPackage(&pip.PackageSpec{Name: "requests"})
})
```

## CLI Logging

The CLI tool supports various logging options:

```bash
# Enable verbose logging
pip-cli -verbose install requests

# Set custom timeout with logging
pip-cli -timeout 60s -verbose install tensorflow

# Different log levels
pip-cli -log-level debug install requests
pip-cli -log-level error install requests
```

## Log Analysis

### Log Parsing

```go
import (
    "bufio"
    "regexp"
    "strings"
)

type LogEntry struct {
    Timestamp string
    Level     string
    Message   string
    Operation string
}

func ParseLogFile(filename string) ([]LogEntry, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    var entries []LogEntry
    scanner := bufio.NewScanner(file)
    
    // Pattern: [2023-01-01 12:00:00] [INFO] [pip-sdk] Installing package: requests
    pattern := regexp.MustCompile(`\[([^\]]+)\] \[([^\]]+)\] \[([^\]]+)\] (.+)`)
    
    for scanner.Scan() {
        line := scanner.Text()
        matches := pattern.FindStringSubmatch(line)
        
        if len(matches) == 5 {
            entry := LogEntry{
                Timestamp: matches[1],
                Level:     matches[2],
                Message:   matches[4],
            }
            
            // Extract operation from message
            if strings.Contains(entry.Message, "Installing package:") {
                entry.Operation = "install"
            } else if strings.Contains(entry.Message, "Uninstalling package:") {
                entry.Operation = "uninstall"
            }
            
            entries = append(entries, entry)
        }
    }
    
    return entries, scanner.Err()
}
```

## Monitoring and Alerting

### Error Rate Monitoring

```go
type ErrorRateMonitor struct {
    logger pip.Logger
    errors int
    total  int
    mutex  sync.Mutex
}

func (erm *ErrorRateMonitor) LogOperation(operation string, err error) {
    erm.mutex.Lock()
    defer erm.mutex.Unlock()
    
    erm.total++
    if err != nil {
        erm.errors++
        erm.logger.Error("Operation failed: %s: %v", operation, err)
    } else {
        erm.logger.Info("Operation succeeded: %s", operation)
    }
    
    // Alert if error rate is too high
    if erm.total >= 10 && float64(erm.errors)/float64(erm.total) > 0.5 {
        erm.logger.Error("High error rate detected: %d/%d operations failed", 
            erm.errors, erm.total)
    }
}
```

## Best Practices

1. **Use Appropriate Log Levels**: Debug for development, Info for production
2. **Include Context**: Add relevant context to log messages
3. **Structured Logging**: Use structured logging for better analysis
4. **Performance Monitoring**: Log operation timings
5. **Error Details**: Include full error details in logs
6. **Log Rotation**: Implement log rotation for long-running applications
7. **Security**: Don't log sensitive information like credentials

## Next Steps

- Explore [Examples](../examples/) for logging patterns
- Check the [API Reference](../api/logger.md) for detailed logging API
- Learn about [Error Handling](./error-handling.md) integration
