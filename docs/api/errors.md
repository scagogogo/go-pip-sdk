# Errors

The Go Pip SDK provides a comprehensive error handling system with structured error types, helpful suggestions, and context-aware error messages.

## Error Types

### ErrorType

```go
type ErrorType string
```

Enumeration of all possible error types in the SDK.

**Constants:**
```go
const (
    ErrorTypePipNotInstalled    ErrorType = "pip_not_installed"
    ErrorTypePythonNotFound     ErrorType = "python_not_found"
    ErrorTypePackageNotFound    ErrorType = "package_not_found"
    ErrorTypeVenvNotFound       ErrorType = "venv_not_found"
    ErrorTypeVenvAlreadyExists  ErrorType = "venv_already_exists"
    ErrorTypeInvalidPackageSpec ErrorType = "invalid_package_spec"
    ErrorTypeCommandFailed      ErrorType = "command_failed"
    ErrorTypePermissionDenied   ErrorType = "permission_denied"
    ErrorTypeNetworkError       ErrorType = "network_error"
    ErrorTypeUnsupportedOS      ErrorType = "unsupported_os"
    ErrorTypeFileNotFound       ErrorType = "file_not_found"
    ErrorTypeInvalidPath        ErrorType = "invalid_path"
    ErrorTypeFeatureDisabled    ErrorType = "feature_disabled"
    ErrorTypeTimeout            ErrorType = "timeout"
    ErrorTypeInvalidConfig      ErrorType = "invalid_config"
)
```

## Core Error Types

### PipErrorDetails

```go
type PipErrorDetails struct {
    Type        ErrorType         // Error type classification
    Message     string            // Human-readable error message
    Command     string            // Command that failed (if applicable)
    Output      string            // Command output (if applicable)
    ExitCode    int               // Command exit code (if applicable)
    Suggestions []string          // Helpful suggestions for resolution
    Context     map[string]string // Additional context information
    Cause       error             // Underlying cause (if any)
}
```

The main error type that implements the standard Go `error` interface with additional context.

**Methods:**

#### Error

```go
func (e *PipErrorDetails) Error() string
```

Returns a formatted error message including type, message, command, and suggestions.

#### Unwrap

```go
func (e *PipErrorDetails) Unwrap() error
```

Returns the underlying cause error for error chain unwrapping.

#### Is

```go
func (e *PipErrorDetails) Is(target error) bool
```

Checks if the error matches a target error type.

#### WithSuggestion

```go
func (e *PipErrorDetails) WithSuggestion(suggestion string) *PipErrorDetails
```

Adds a suggestion to the error and returns the error for chaining.

**Example:**
```go
err := pip.NewPipError(pip.ErrorTypePackageNotFound, "package 'nonexistent' not found").
    WithSuggestion("Check the package name spelling").
    WithSuggestion("Search for the package on PyPI")

fmt.Println(err.Error())
// Output: [package_not_found] package 'nonexistent' not found | Suggestions: Check the package name spelling, Search for the package on PyPI
```

## Constructor Functions

### NewPipError

```go
func NewPipError(errorType ErrorType, message string) *PipErrorDetails
```

Creates a new pip error with the specified type and message.

**Parameters:**
- `errorType` (ErrorType): The error type classification.
- `message` (string): Human-readable error message.

**Returns:**
- `*PipErrorDetails`: New error instance.

### NewCommandError

```go
func NewCommandError(command string, output string, exitCode int, cause error) *PipErrorDetails
```

Creates a new command execution error with automatic suggestion generation.

**Parameters:**
- `command` (string): The command that failed.
- `output` (string): Command output.
- `exitCode` (int): Command exit code.
- `cause` (error): Underlying cause error.

**Returns:**
- `*PipErrorDetails`: New command error instance.

### WrapError

```go
func WrapError(err error, errorType ErrorType, message string) *PipErrorDetails
```

Wraps a generic error into a PipErrorDetails with additional context.

**Parameters:**
- `err` (error): Original error to wrap.
- `errorType` (ErrorType): Error type classification.
- `message` (string): Additional context message.

**Returns:**
- `*PipErrorDetails`: Wrapped error.

## Predefined Errors

The SDK provides predefined error instances for common scenarios:

```go
var (
    ErrPipNotInstalled = NewPipError(ErrorTypePipNotInstalled, "pip is not installed").
        WithSuggestion("Install pip using your system package manager or download get-pip.py")

    ErrPythonNotFound = NewPipError(ErrorTypePythonNotFound, "Python interpreter not found").
        WithSuggestion("Install Python from https://python.org or use your system package manager")

    ErrPackageNotFound = NewPipError(ErrorTypePackageNotFound, "package not found").
        WithSuggestion("Check the package name spelling and availability on PyPI")

    ErrVenvNotFound = NewPipError(ErrorTypeVenvNotFound, "virtual environment not found").
        WithSuggestion("Create a virtual environment first using CreateVenv()")

    ErrVenvAlreadyExists = NewPipError(ErrorTypeVenvAlreadyExists, "virtual environment already exists").
        WithSuggestion("Use a different path or remove the existing virtual environment")

    ErrInvalidPackageSpec = NewPipError(ErrorTypeInvalidPackageSpec, "invalid package specification").
        WithSuggestion("Ensure package name is not empty and version constraints are valid")

    ErrPermissionDenied = NewPipError(ErrorTypePermissionDenied, "permission denied").
        WithSuggestion("Run with elevated privileges or check file/directory permissions")

    ErrNetworkError = NewPipError(ErrorTypeNetworkError, "network error").
        WithSuggestion("Check your internet connection and proxy settings")

    ErrUnsupportedOS = NewPipError(ErrorTypeUnsupportedOS, "unsupported operating system").
        WithSuggestion("This operation is not supported on your operating system")

    ErrFeatureDisabled = NewPipError(ErrorTypeFeatureDisabled, "feature is disabled").
        WithSuggestion("This feature has been disabled or is not available")
)
```

## Utility Functions

### IsErrorType

```go
func IsErrorType(err error, errorType ErrorType) bool
```

Checks if an error is of a specific type.

**Parameters:**
- `err` (error): Error to check.
- `errorType` (ErrorType): Expected error type.

**Returns:**
- `bool`: `true` if the error matches the specified type.

**Example:**
```go
if err := manager.InstallPackage(pkg); err != nil {
    if pip.IsErrorType(err, pip.ErrorTypePackageNotFound) {
        fmt.Println("Package not found - check the name")
        return nil // Handle gracefully
    }
    return err // Propagate other errors
}
```

### GetErrorType

```go
func GetErrorType(err error) ErrorType
```

Returns the error type if it's a PipErrorDetails, otherwise returns empty string.

**Parameters:**
- `err` (error): Error to examine.

**Returns:**
- `ErrorType`: Error type or empty string.

**Example:**
```go
if err := manager.CreateVenv(path); err != nil {
    switch pip.GetErrorType(err) {
    case pip.ErrorTypeVenvAlreadyExists:
        fmt.Println("Virtual environment already exists")
    case pip.ErrorTypePermissionDenied:
        fmt.Println("Permission denied")
    default:
        fmt.Printf("Unexpected error: %v\n", err)
    }
}
```

## Error Handler

### ErrorHandler

```go
type ErrorHandler struct {
    logger *Logger
}
```

Provides centralized error handling and logging.

#### NewErrorHandler

```go
func NewErrorHandler(logger *Logger) *ErrorHandler
```

Creates a new error handler with the specified logger.

#### Handle

```go
func (eh *ErrorHandler) Handle(err error, context string) error
```

Processes an error and logs it appropriately based on the error type.

**Parameters:**
- `err` (error): Error to handle.
- `context` (string): Context description for logging.

**Returns:**
- `error`: Processed error (may be wrapped if it wasn't already a PipErrorDetails).

**Example:**
```go
logger, _ := pip.NewLogger(pip.DefaultLoggerConfig())
handler := pip.NewErrorHandler(logger)

if err := manager.InstallPackage(pkg); err != nil {
    return handler.Handle(err, "package installation")
}
```

## Error Handling Patterns

### Basic Error Checking

```go
if err := manager.InstallPackage(pkg); err != nil {
    if pip.IsErrorType(err, pip.ErrorTypePackageNotFound) {
        // Handle package not found specifically
        fmt.Printf("Package %s not found\n", pkg.Name)
        return nil
    }
    // Handle other errors
    return fmt.Errorf("installation failed: %w", err)
}
```

### Comprehensive Error Handling

```go
func handlePipError(err error) {
    if err == nil {
        return
    }

    pipErr, ok := err.(*pip.PipErrorDetails)
    if !ok {
        fmt.Printf("Unexpected error: %v\n", err)
        return
    }

    fmt.Printf("Error Type: %s\n", pipErr.Type)
    fmt.Printf("Message: %s\n", pipErr.Message)

    if pipErr.Command != "" {
        fmt.Printf("Failed Command: %s\n", pipErr.Command)
    }

    if pipErr.ExitCode != 0 {
        fmt.Printf("Exit Code: %d\n", pipErr.ExitCode)
    }

    if len(pipErr.Suggestions) > 0 {
        fmt.Println("Suggestions:")
        for _, suggestion := range pipErr.Suggestions {
            fmt.Printf("  - %s\n", suggestion)
        }
    }

    if pipErr.Cause != nil {
        fmt.Printf("Underlying cause: %v\n", pipErr.Cause)
    }
}
```

### Error Recovery

```go
func installWithRetry(manager *pip.Manager, pkg *pip.PackageSpec, maxRetries int) error {
    var lastErr error

    for i := 0; i < maxRetries; i++ {
        err := manager.InstallPackage(pkg)
        if err == nil {
            return nil // Success
        }

        lastErr = err

        // Check if error is recoverable
        if pip.IsErrorType(err, pip.ErrorTypeNetworkError) {
            fmt.Printf("Network error, retrying... (%d/%d)\n", i+1, maxRetries)
            time.Sleep(time.Duration(i+1) * time.Second)
            continue
        }

        // Non-recoverable error
        break
    }

    return lastErr
}
```

### Graceful Degradation

```go
func installOptionalPackage(manager *pip.Manager, pkg *pip.PackageSpec) {
    if err := manager.InstallPackage(pkg); err != nil {
        switch pip.GetErrorType(err) {
        case pip.ErrorTypePackageNotFound:
            fmt.Printf("Optional package %s not available, skipping\n", pkg.Name)
        case pip.ErrorTypePermissionDenied:
            fmt.Printf("Permission denied for %s, skipping\n", pkg.Name)
        default:
            fmt.Printf("Failed to install optional package %s: %v\n", pkg.Name, err)
        }
        // Continue execution without failing
    }
}
```

## Best Practices

1. **Always check error types for specific handling**:
   ```go
   if err != nil {
       if pip.IsErrorType(err, pip.ErrorTypePackageNotFound) {
           // Specific handling for package not found
       } else {
           // Generic error handling
       }
   }
   ```

2. **Use error wrapping for context**:
   ```go
   if err := manager.InstallPackage(pkg); err != nil {
       return fmt.Errorf("failed to install %s: %w", pkg.Name, err)
   }
   ```

3. **Log errors with appropriate levels**:
   ```go
   if err := manager.InstallPackage(pkg); err != nil {
       if pip.IsErrorType(err, pip.ErrorTypeNetworkError) {
           logger.Warn("Network issue during installation: %v", err)
       } else {
           logger.Error("Installation failed: %v", err)
       }
   }
   ```

4. **Provide helpful error messages to users**:
   ```go
   if err := manager.CreateVenv(path); err != nil {
       if pip.IsErrorType(err, pip.ErrorTypePermissionDenied) {
           return fmt.Errorf("cannot create virtual environment: permission denied. Try running with elevated privileges or choose a different location")
       }
       return err
   }
   ```

## Error Context

Errors can include additional context information:

```go
err := pip.NewPipError(pip.ErrorTypePackageNotFound, "package not found")
if pipErr, ok := err.(*pip.PipErrorDetails); ok {
    pipErr.Context["package_name"] = "nonexistent-package"
    pipErr.Context["index_url"] = "https://pypi.org/simple/"
    pipErr.Context["python_version"] = "3.9"
}
```

This context can be used for debugging, logging, or providing more detailed error information to users.
```
