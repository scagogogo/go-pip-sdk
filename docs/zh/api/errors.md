# 错误处理 API

错误处理 API 提供了完整的错误类型定义、错误检查和错误处理功能。

## 错误类型

### ErrorType

错误类型枚举。

```go
type ErrorType string

const (
    // 系统错误
    ErrorTypePythonNotFound    ErrorType = "python_not_found"
    ErrorTypePipNotFound       ErrorType = "pip_not_found"
    ErrorTypePermissionDenied  ErrorType = "permission_denied"
    ErrorTypeInsufficientSpace ErrorType = "insufficient_space"
    
    // 网络错误
    ErrorTypeNetworkError      ErrorType = "network_error"
    ErrorTypeNetworkTimeout    ErrorType = "network_timeout"
    ErrorTypeConnectionRefused ErrorType = "connection_refused"
    
    // 包相关错误
    ErrorTypePackageNotFound   ErrorType = "package_not_found"
    ErrorTypeVersionConflict   ErrorType = "version_conflict"
    ErrorTypeCorruptedPackage  ErrorType = "corrupted_package"
    ErrorTypeDependencyError   ErrorType = "dependency_error"
    
    // 虚拟环境错误
    ErrorTypeVenvExists        ErrorType = "venv_exists"
    ErrorTypeVenvNotFound      ErrorType = "venv_not_found"
    ErrorTypeVenvCorrupted     ErrorType = "venv_corrupted"
    
    // 项目错误
    ErrorTypeProjectExists     ErrorType = "project_exists"
    ErrorTypeProjectNotFound   ErrorType = "project_not_found"
    ErrorTypeInvalidProjectName ErrorType = "invalid_project_name"
    ErrorTypeConfigError       ErrorType = "config_error"
    
    // 通用错误
    ErrorTypeTimeout           ErrorType = "timeout"
    ErrorTypeInvalidArgument   ErrorType = "invalid_argument"
    ErrorTypeInternalError     ErrorType = "internal_error"
)
```

## 错误结构

### PipErrorDetails

详细错误信息结构。

```go
type PipErrorDetails struct {
    Type        ErrorType              // 错误类型
    Message     string                 // 错误消息
    Command     string                 // 失败的命令
    ExitCode    int                    // 退出代码
    Output      string                 // 命令输出
    Stderr      string                 // 错误输出
    Suggestions []string               // 解决建议
    Context     map[string]interface{} // 错误上下文
    Cause       error                  // 根本原因
    Timestamp   time.Time              // 错误时间
}
```

**方法：**

```go
func (e *PipErrorDetails) Error() string
func (e *PipErrorDetails) Unwrap() error
func (e *PipErrorDetails) Is(target error) bool
```

**示例：**

```go
if err := manager.InstallPackage(pkg); err != nil {
    if pipErr, ok := err.(*pip.PipErrorDetails); ok {
        fmt.Printf("错误类型: %s\n", pipErr.Type)
        fmt.Printf("错误消息: %s\n", pipErr.Message)
        fmt.Printf("失败命令: %s\n", pipErr.Command)
        fmt.Printf("退出代码: %d\n", pipErr.ExitCode)
        
        if len(pipErr.Suggestions) > 0 {
            fmt.Println("建议:")
            for _, suggestion := range pipErr.Suggestions {
                fmt.Printf("- %s\n", suggestion)
            }
        }
    }
}
```

## 错误检查函数

### IsErrorType

检查错误是否为特定类型。

```go
func IsErrorType(err error, errorType ErrorType) bool
```

**参数：**
- `err` - 要检查的错误
- `errorType` - 错误类型

**返回值：**
- `bool` - 是否匹配

**示例：**

```go
if err := manager.InstallPackage(pkg); err != nil {
    if pip.IsErrorType(err, pip.ErrorTypePackageNotFound) {
        fmt.Println("包未找到")
    } else if pip.IsErrorType(err, pip.ErrorTypeNetworkError) {
        fmt.Println("网络错误")
    }
}
```

### GetErrorType

获取错误类型。

```go
func GetErrorType(err error) ErrorType
```

**参数：**
- `err` - 错误对象

**返回值：**
- `ErrorType` - 错误类型

**示例：**

```go
if err := manager.InstallPackage(pkg); err != nil {
    errorType := pip.GetErrorType(err)
    fmt.Printf("错误类型: %s\n", errorType)
}
```

### IsRetryableError

检查错误是否可重试。

```go
func IsRetryableError(err error) bool
```

**参数：**
- `err` - 错误对象

**返回值：**
- `bool` - 是否可重试

**示例：**

```go
if err := manager.InstallPackage(pkg); err != nil {
    if pip.IsRetryableError(err) {
        fmt.Println("这是一个可重试的错误")
        // 执行重试逻辑
    }
}
```

## 错误创建函数

### NewPipError

创建新的 pip 错误。

```go
func NewPipError(errorType ErrorType, message string) *PipErrorDetails
```

**示例：**

```go
err := pip.NewPipError(
    pip.ErrorTypePackageNotFound,
    "Package 'nonexistent-package' not found",
)
```

### NewPipErrorWithContext

创建带上下文的 pip 错误。

```go
func NewPipErrorWithContext(errorType ErrorType, message string, context map[string]interface{}) *PipErrorDetails
```

**示例：**

```go
context := map[string]interface{}{
    "package_name": "requests",
    "version":      "2.25.1",
    "index_url":    "https://pypi.org/simple/",
}

err := pip.NewPipErrorWithContext(
    pip.ErrorTypeNetworkError,
    "Failed to download package",
    context,
)
```

### WrapError

包装现有错误。

```go
func WrapError(err error, errorType ErrorType, message string) *PipErrorDetails
```

**示例：**

```go
originalErr := fmt.Errorf("connection refused")
wrappedErr := pip.WrapError(
    originalErr,
    pip.ErrorTypeNetworkError,
    "Failed to connect to PyPI",
)
```

## 错误处理模式

### 基本错误处理

```go
func handleBasicError(err error) {
    if err == nil {
        return
    }
    
    switch pip.GetErrorType(err) {
    case pip.ErrorTypePackageNotFound:
        fmt.Println("包未找到，请检查包名")
    case pip.ErrorTypeNetworkError:
        fmt.Println("网络连接失败，请检查网络")
    case pip.ErrorTypePermissionDenied:
        fmt.Println("权限被拒绝，请使用管理员权限或虚拟环境")
    default:
        fmt.Printf("未知错误: %v\n", err)
    }
}
```

### 重试错误处理

```go
func handleRetryableError(operation func() error, maxRetries int) error {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        err := operation()
        if err == nil {
            return nil
        }
        
        lastErr = err
        
        if !pip.IsRetryableError(err) {
            return err // 不可重试的错误
        }
        
        fmt.Printf("重试 %d/%d...\n", i+1, maxRetries)
        time.Sleep(time.Duration(i+1) * time.Second)
    }
    
    return fmt.Errorf("重试 %d 次后仍然失败: %w", maxRetries, lastErr)
}
```

### 错误恢复处理

```go
func handleErrorWithRecovery(err error, pkg *pip.PackageSpec, manager *pip.Manager) error {
    switch pip.GetErrorType(err) {
    case pip.ErrorTypePackageNotFound:
        return suggestAlternativePackages(manager, pkg.Name)
        
    case pip.ErrorTypeVersionConflict:
        return resolveVersionConflict(manager, pkg)
        
    case pip.ErrorTypePermissionDenied:
        return installAsUser(manager, pkg)
        
    case pip.ErrorTypeNetworkError:
        return retryWithMirror(manager, pkg)
        
    default:
        return err
    }
}
```

## 错误聚合

### ErrorCollector

错误收集器，用于收集多个错误。

```go
type ErrorCollector struct {
    errors []error
    mutex  sync.RWMutex
}

func NewErrorCollector() *ErrorCollector
func (c *ErrorCollector) Add(err error)
func (c *ErrorCollector) HasErrors() bool
func (c *ErrorCollector) Errors() []error
func (c *ErrorCollector) Error() string
```

**示例：**

```go
collector := pip.NewErrorCollector()

packages := []string{"pkg1", "pkg2", "pkg3"}
for _, pkgName := range packages {
    pkg := &pip.PackageSpec{Name: pkgName}
    if err := manager.InstallPackage(pkg); err != nil {
        collector.Add(err)
    }
}

if collector.HasErrors() {
    fmt.Printf("安装过程中发生 %d 个错误:\n", len(collector.Errors()))
    for i, err := range collector.Errors() {
        fmt.Printf("%d. %v\n", i+1, err)
    }
}
```

## 错误统计

### ErrorStats

错误统计信息。

```go
type ErrorStats struct {
    Total       int
    ByType      map[ErrorType]int
    ByOperation map[string]int
    LastError   time.Time
}

func (s *ErrorStats) RecordError(err error, operation string)
func (s *ErrorStats) GetStats() map[string]interface{}
func (s *ErrorStats) Reset()
```

**示例：**

```go
stats := &pip.ErrorStats{
    ByType:      make(map[pip.ErrorType]int),
    ByOperation: make(map[string]int),
}

// 记录错误
stats.RecordError(err, "install_package")

// 获取统计信息
data := stats.GetStats()
fmt.Printf("总错误数: %d\n", data["total"])
```

## 错误报告

### GenerateErrorReport

生成错误报告。

```go
func GenerateErrorReport(errors []error) string
```

**示例：**

```go
errors := []error{
    pip.NewPipError(pip.ErrorTypePackageNotFound, "Package A not found"),
    pip.NewPipError(pip.ErrorTypeNetworkError, "Network timeout"),
}

report := pip.GenerateErrorReport(errors)
fmt.Println(report)
```

### ExportErrorLog

导出错误日志。

```go
func ExportErrorLog(errors []error, filename string) error
```

**示例：**

```go
if err := pip.ExportErrorLog(errors, "error_log.json"); err != nil {
    log.Printf("导出错误日志失败: %v", err)
}
```

## 自定义错误处理

### ErrorHandler

错误处理器接口。

```go
type ErrorHandler interface {
    HandleError(err error, context map[string]interface{}) error
}
```

### 实现自定义错误处理器

```go
type CustomErrorHandler struct {
    logger Logger
}

func (h *CustomErrorHandler) HandleError(err error, context map[string]interface{}) error {
    // 记录错误
    h.logger.Error("Operation failed", map[string]interface{}{
        "error":   err.Error(),
        "context": context,
    })
    
    // 根据错误类型执行不同的处理逻辑
    switch pip.GetErrorType(err) {
    case pip.ErrorTypeNetworkError:
        return h.handleNetworkError(err, context)
    case pip.ErrorTypePackageNotFound:
        return h.handlePackageNotFound(err, context)
    default:
        return err
    }
}
```

## 最佳实践

1. **总是检查错误**：不要忽略任何错误
2. **使用类型检查**：使用 `IsErrorType` 检查特定错误类型
3. **提供有用的错误消息**：包含足够的上下文信息
4. **实现重试机制**：对可重试的错误进行重试
5. **记录错误**：记录错误以便调试和监控

## 下一步

- 查看[日志记录 API](./logger.md)
- 了解[类型定义](./types.md)
- 探索[安装器 API](./installer.md)
- 学习[管理器 API](./manager.md)
