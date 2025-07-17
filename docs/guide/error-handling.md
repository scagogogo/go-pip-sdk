# Error Handling

This guide covers how to handle errors effectively when using the Go Pip SDK.

## Overview

The Go Pip SDK provides comprehensive error handling with detailed error information, suggestions for resolution, and different error types for different scenarios.

## Error Types

The SDK defines several error types to help you handle different scenarios:

```go
const (
    ErrorTypePipNotInstalled    = "pip_not_installed"
    ErrorTypePythonNotFound     = "python_not_found"
    ErrorTypePackageNotFound    = "package_not_found"
    ErrorTypePermissionDenied   = "permission_denied"
    ErrorTypeNetworkError       = "network_error"
    ErrorTypeTimeout            = "timeout"
    ErrorTypeCommandFailed      = "command_failed"
    ErrorTypeInvalidPackageSpec = "invalid_package_spec"
    ErrorTypeVenvNotFound       = "venv_not_found"
    ErrorTypeVenvAlreadyExists  = "venv_already_exists"
)
```

## Basic Error Handling

### Simple Error Checking

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    manager := pip.NewManager(nil)
    
    pkg := &pip.PackageSpec{Name: "requests"}
    err := manager.InstallPackage(pkg)
    if err != nil {
        log.Printf("Failed to install package: %v", err)
        return
    }
    
    fmt.Println("Package installed successfully")
}
```

### Detailed Error Information

```go
err := manager.InstallPackage(&pip.PackageSpec{Name: "nonexistent-package"})
if err != nil {
    if pipErr, ok := err.(*pip.PipErrorDetails); ok {
        fmt.Printf("Error Type: %s\n", pipErr.Type)
        fmt.Printf("Message: %s\n", pipErr.Message)
        fmt.Printf("Command: %s\n", pipErr.Command)
        fmt.Printf("Exit Code: %d\n", pipErr.ExitCode)
        fmt.Printf("Output: %s\n", pipErr.Output)
        
        if len(pipErr.Suggestions) > 0 {
            fmt.Println("Suggestions:")
            for _, suggestion := range pipErr.Suggestions {
                fmt.Printf("  - %s\n", suggestion)
            }
        }
    }
}
```

## Handling Specific Error Types

### Package Not Found

```go
err := manager.InstallPackage(&pip.PackageSpec{Name: "nonexistent-package"})
if err != nil {
    if pip.IsErrorType(err, pip.ErrorTypePackageNotFound) {
        fmt.Println("Package not found. Please check the package name.")
        fmt.Println("You can search for packages at https://pypi.org")
        return
    }
}
```

### Permission Denied

```go
err := manager.InstallPackage(&pip.PackageSpec{Name: "requests"})
if err != nil {
    if pip.IsErrorType(err, pip.ErrorTypePermissionDenied) {
        fmt.Println("Permission denied. Try one of the following:")
        fmt.Println("1. Use a virtual environment")
        fmt.Println("2. Install with --user flag")
        fmt.Println("3. Run with sudo (not recommended)")
        return
    }
}
```

### Network Errors

```go
err := manager.InstallPackage(&pip.PackageSpec{Name: "requests"})
if err != nil {
    if pip.IsErrorType(err, pip.ErrorTypeNetworkError) {
        fmt.Println("Network error occurred. Please check:")
        fmt.Println("1. Internet connection")
        fmt.Println("2. Proxy settings")
        fmt.Println("3. PyPI server status")
        return
    }
}
```

### Pip Not Installed

```go
installed, err := manager.IsInstalled()
if err != nil {
    if pip.IsErrorType(err, pip.ErrorTypePipNotInstalled) {
        fmt.Println("Pip is not installed. Please install pip first:")
        fmt.Println("- macOS: brew install python")
        fmt.Println("- Ubuntu: sudo apt install python3-pip")
        fmt.Println("- Windows: Download from python.org")
        return
    }
}
```

## Error Recovery Strategies

### Retry with Backoff

```go
import (
    "time"
    "math"
)

func installWithRetry(manager pip.PipManager, pkg *pip.PackageSpec, maxRetries int) error {
    for attempt := 0; attempt < maxRetries; attempt++ {
        err := manager.InstallPackage(pkg)
        if err == nil {
            return nil
        }
        
        // Check if it's a retryable error
        if pip.IsErrorType(err, pip.ErrorTypeNetworkError) || 
           pip.IsErrorType(err, pip.ErrorTypeTimeout) {
            
            if attempt < maxRetries-1 {
                // Exponential backoff
                delay := time.Duration(math.Pow(2, float64(attempt))) * time.Second
                fmt.Printf("Attempt %d failed, retrying in %v...\n", attempt+1, delay)
                time.Sleep(delay)
                continue
            }
        }
        
        // Non-retryable error or max retries reached
        return err
    }
    
    return fmt.Errorf("failed after %d attempts", maxRetries)
}
```

### Fallback Strategies

```go
func installWithFallback(manager pip.PipManager, pkg *pip.PackageSpec) error {
    // Try normal installation first
    err := manager.InstallPackage(pkg)
    if err == nil {
        return nil
    }
    
    // If permission denied, try with --user flag
    if pip.IsErrorType(err, pip.ErrorTypePermissionDenied) {
        fmt.Println("Permission denied, trying user installation...")
        pkg.Options = map[string]string{"user": ""}
        err = manager.InstallPackage(pkg)
        if err == nil {
            return nil
        }
    }
    
    // If network error, try with different index
    if pip.IsErrorType(err, pip.ErrorTypeNetworkError) {
        fmt.Println("Network error, trying alternative index...")
        pkg.Index = "https://mirrors.aliyun.com/pypi/simple/"
        err = manager.InstallPackage(pkg)
        if err == nil {
            return nil
        }
    }
    
    return err
}
```

## Validation and Prevention

### Package Spec Validation

```go
func validatePackageSpec(pkg *pip.PackageSpec) error {
    if pkg == nil {
        return fmt.Errorf("package specification cannot be nil")
    }
    
    if pkg.Name == "" {
        return fmt.Errorf("package name cannot be empty")
    }
    
    // Validate package name format
    if !isValidPackageName(pkg.Name) {
        return fmt.Errorf("invalid package name: %s", pkg.Name)
    }
    
    // Validate version specification
    if pkg.Version != "" && !isValidVersionSpec(pkg.Version) {
        return fmt.Errorf("invalid version specification: %s", pkg.Version)
    }
    
    return nil
}

func isValidPackageName(name string) bool {
    // Package names should contain only letters, numbers, hyphens, and underscores
    for _, char := range name {
        if !((char >= 'a' && char <= 'z') || 
             (char >= 'A' && char <= 'Z') || 
             (char >= '0' && char <= '9') || 
             char == '-' || char == '_' || char == '.') {
            return false
        }
    }
    return true
}
```

### Environment Validation

```go
func validateEnvironment(manager pip.PipManager) error {
    // Check if pip is installed
    installed, err := manager.IsInstalled()
    if err != nil {
        return fmt.Errorf("failed to check pip installation: %w", err)
    }
    
    if !installed {
        return fmt.Errorf("pip is not installed")
    }
    
    // Check pip version
    version, err := manager.GetVersion()
    if err != nil {
        return fmt.Errorf("failed to get pip version: %w", err)
    }
    
    fmt.Printf("Using pip version: %s\n", version)
    
    // Check if we're in a virtual environment (recommended)
    venvManager := pip.NewVenvManager(manager)
    if !venvManager.IsVenvActive() {
        fmt.Println("Warning: Not in a virtual environment")
        fmt.Println("Consider using a virtual environment for better isolation")
    }
    
    return nil
}
```

## Error Logging

### Structured Logging

```go
import (
    "log/slog"
    "os"
)

func setupErrorLogging() *slog.Logger {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    }))
    
    return logger
}

func logError(logger *slog.Logger, operation string, err error) {
    if pipErr, ok := err.(*pip.PipErrorDetails); ok {
        logger.Error("pip operation failed",
            "operation", operation,
            "error_type", pipErr.Type,
            "message", pipErr.Message,
            "command", pipErr.Command,
            "exit_code", pipErr.ExitCode,
            "suggestions", pipErr.Suggestions,
        )
    } else {
        logger.Error("operation failed",
            "operation", operation,
            "error", err.Error(),
        )
    }
}
```

## Error Handler Utility

```go
type ErrorHandler struct {
    logger *slog.Logger
    retryConfig RetryConfig
}

type RetryConfig struct {
    MaxRetries int
    BaseDelay  time.Duration
    MaxDelay   time.Duration
}

func NewErrorHandler(logger *slog.Logger) *ErrorHandler {
    return &ErrorHandler{
        logger: logger,
        retryConfig: RetryConfig{
            MaxRetries: 3,
            BaseDelay:  time.Second,
            MaxDelay:   30 * time.Second,
        },
    }
}

func (eh *ErrorHandler) HandleWithRetry(operation string, fn func() error) error {
    var lastErr error
    
    for attempt := 0; attempt < eh.retryConfig.MaxRetries; attempt++ {
        err := fn()
        if err == nil {
            return nil
        }
        
        lastErr = err
        
        // Log the error
        eh.logError(operation, err, attempt)
        
        // Check if retryable
        if !eh.isRetryable(err) {
            return err
        }
        
        // Calculate delay
        delay := eh.calculateDelay(attempt)
        if attempt < eh.retryConfig.MaxRetries-1 {
            eh.logger.Info("retrying operation",
                "operation", operation,
                "attempt", attempt+1,
                "delay", delay,
            )
            time.Sleep(delay)
        }
    }
    
    return fmt.Errorf("operation failed after %d attempts: %w", 
        eh.retryConfig.MaxRetries, lastErr)
}

func (eh *ErrorHandler) isRetryable(err error) bool {
    return pip.IsErrorType(err, pip.ErrorTypeNetworkError) ||
           pip.IsErrorType(err, pip.ErrorTypeTimeout)
}

func (eh *ErrorHandler) calculateDelay(attempt int) time.Duration {
    delay := time.Duration(math.Pow(2, float64(attempt))) * eh.retryConfig.BaseDelay
    if delay > eh.retryConfig.MaxDelay {
        delay = eh.retryConfig.MaxDelay
    }
    return delay
}
```

## Best Practices

1. **Always Check Errors**: Never ignore errors from SDK operations
2. **Use Specific Error Types**: Handle different error types appropriately
3. **Provide User-Friendly Messages**: Convert technical errors to user-friendly messages
4. **Implement Retry Logic**: For transient errors like network issues
5. **Log Errors Properly**: Use structured logging for better debugging
6. **Validate Input**: Validate package specifications before operations
7. **Graceful Degradation**: Provide fallback options when possible

## Next Steps

- Learn about [Logging](./logging.md) configuration
- Explore [Examples](../examples/) for error handling patterns
- Check the [API Reference](../api/) for detailed error information
