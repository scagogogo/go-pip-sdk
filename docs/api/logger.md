# Logger

The Go Pip SDK provides a comprehensive logging system with configurable levels, structured output, and specialized logging methods for pip operations.

## Log Levels

### LogLevel

```go
type LogLevel int

const (
    LogLevelDebug LogLevel = iota
    LogLevelInfo
    LogLevelWarn
    LogLevelError
    LogLevelFatal
)
```

Enumeration of available log levels in order of severity.

**Methods:**

#### String

```go
func (l LogLevel) String() string
```

Returns the string representation of the log level.

**Example:**
```go
level := pip.LogLevelInfo
fmt.Println(level.String()) // Output: "INFO"
```

## Core Types

### Logger

```go
type Logger struct {
    level      LogLevel
    output     io.Writer
    prefix     string
    flags      int
    enableFile bool
    logFile    *os.File
}
```

The main logger implementation with structured logging capabilities.

### LoggerConfig

```go
type LoggerConfig struct {
    Level      LogLevel  // Minimum log level to output
    Output     io.Writer // Output destination (default: os.Stdout)
    Prefix     string    // Log prefix (default: "[pip-sdk]")
    EnableFile bool      // Enable file logging
    LogFile    string    // Log file path (if EnableFile is true)
    Flags      int       // Log flags (default: log.LstdFlags)
}
```

Configuration for logger creation.

## Constructor Functions

### NewLogger

```go
func NewLogger(config *LoggerConfig) (*Logger, error)
```

Creates a new logger instance with the specified configuration.

**Parameters:**
- `config` (*LoggerConfig): Logger configuration. If `nil`, default configuration is used.

**Returns:**
- `*Logger`: New logger instance.
- `error`: Error if logger creation fails.

**Example:**
```go
// Create with default configuration
logger, err := pip.NewLogger(nil)
if err != nil {
    return err
}
defer logger.Close()

// Create with custom configuration
config := &pip.LoggerConfig{
    Level:      pip.LogLevelDebug,
    Output:     os.Stdout,
    Prefix:     "[my-app]",
    EnableFile: true,
    LogFile:    "/var/log/pip-sdk.log",
}

logger, err := pip.NewLogger(config)
if err != nil {
    return err
}
defer logger.Close()
```

### DefaultLoggerConfig

```go
func DefaultLoggerConfig() *LoggerConfig
```

Returns a default logger configuration with sensible defaults.

**Returns:**
- `*LoggerConfig`: Default configuration.

**Default Values:**
```go
&LoggerConfig{
    Level:      LogLevelInfo,
    Output:     os.Stdout,
    Prefix:     "[pip-sdk]",
    EnableFile: false,
    Flags:      log.LstdFlags,
}
```

## Core Logging Methods

### Debug

```go
func (l *Logger) Debug(format string, args ...interface{})
```

Logs a debug message. Only output if log level is Debug or lower.

### Info

```go
func (l *Logger) Info(format string, args ...interface{})
```

Logs an informational message.

### Warn

```go
func (l *Logger) Warn(format string, args ...interface{})
```

Logs a warning message.

### Error

```go
func (l *Logger) Error(format string, args ...interface{})
```

Logs an error message.

### Fatal

```go
func (l *Logger) Fatal(format string, args ...interface{})
```

Logs a fatal message and exits the program with status code 1.

**Example:**
```go
logger.Debug("Debugging package installation: %s", pkg.Name)
logger.Info("Installing package: %s", pkg.Name)
logger.Warn("Package %s is deprecated", pkg.Name)
logger.Error("Failed to install package: %v", err)
// logger.Fatal("Critical error: %v", err) // This will exit the program
```

## Specialized Logging Methods

### LogCommand

```go
func (l *Logger) LogCommand(command string, args []string, duration time.Duration)
```

Logs command execution with timing information.

**Parameters:**
- `command` (string): Command that was executed.
- `args` ([]string): Command arguments.
- `duration` (time.Duration): Execution duration.

**Example:**
```go
start := time.Now()
// Execute command...
duration := time.Since(start)

logger.LogCommand("pip", []string{"install", "requests"}, duration)
// Output: [INFO] Command executed: pip install requests (duration: 2.5s)
```

### LogCommandError

```go
func (l *Logger) LogCommandError(command string, args []string, err error, output string)
```

Logs command execution errors with detailed information.

**Parameters:**
- `command` (string): Command that failed.
- `args` ([]string): Command arguments.
- `err` (error): Error that occurred.
- `output` (string): Command output.

### LogPackageOperation

```go
func (l *Logger) LogPackageOperation(operation, packageName string, success bool, duration time.Duration)
```

Logs package operations (install, uninstall, etc.) with timing.

**Parameters:**
- `operation` (string): Operation type (e.g., "install", "uninstall").
- `packageName` (string): Name of the package.
- `success` (bool): Whether the operation succeeded.
- `duration` (time.Duration): Operation duration.

**Example:**
```go
start := time.Now()
err := manager.InstallPackage(pkg)
duration := time.Since(start)

logger.LogPackageOperation("install", pkg.Name, err == nil, duration)
```

### LogVenvOperation

```go
func (l *Logger) LogVenvOperation(operation, path string, success bool)
```

Logs virtual environment operations.

**Parameters:**
- `operation` (string): Operation type (e.g., "create", "activate", "remove").
- `path` (string): Virtual environment path.
- `success` (bool): Whether the operation succeeded.

### LogProjectOperation

```go
func (l *Logger) LogProjectOperation(operation, path string, success bool)
```

Logs project operations.

**Parameters:**
- `operation` (string): Operation type (e.g., "init", "build").
- `path` (string): Project path.
- `success` (bool): Whether the operation succeeded.

## Configuration Methods

### SetLevel

```go
func (l *Logger) SetLevel(level LogLevel)
```

Sets the minimum log level for output.

### GetLevel

```go
func (l *Logger) GetLevel() LogLevel
```

Returns the current log level.

### IsEnabled

```go
func (l *Logger) IsEnabled(level LogLevel) bool
```

Checks if a specific log level is enabled.

**Example:**
```go
if logger.IsEnabled(pip.LogLevelDebug) {
    // Expensive debug operation
    debugInfo := generateDebugInfo()
    logger.Debug("Debug info: %s", debugInfo)
}
```

### Close

```go
func (l *Logger) Close() error
```

Closes the logger and any open log files.

**Returns:**
- `error`: Error if closing fails.

## Utility Functions

### ParseLogLevel

```go
func ParseLogLevel(level string) (LogLevel, error)
```

Parses a log level string into a LogLevel constant.

**Parameters:**
- `level` (string): Log level string (case-insensitive).

**Returns:**
- `LogLevel`: Parsed log level.
- `error`: Error if parsing fails.

**Supported Strings:**
- "DEBUG"
- "INFO"
- "WARN", "WARNING"
- "ERROR"
- "FATAL"

**Example:**
```go
level, err := pip.ParseLogLevel("DEBUG")
if err != nil {
    return err
}

logger.SetLevel(level)
```

## Integration with Manager

The logger can be integrated with the pip manager for automatic operation logging.

### SetCustomLogger

```go
func (m *Manager) SetCustomLogger(logger *Logger)
```

Sets a custom logger for the manager to use.

**Example:**
```go
// Create custom logger
logger, err := pip.NewLogger(&pip.LoggerConfig{
    Level:      pip.LogLevelDebug,
    EnableFile: true,
    LogFile:    "pip-operations.log",
})
if err != nil {
    return err
}

// Set logger on manager
manager := pip.NewManager(nil)
manager.SetCustomLogger(logger)

// All operations will now be logged
manager.InstallPackage(&pip.PackageSpec{Name: "requests"})
```

## Log Output Format

The logger produces structured log entries with the following format:

```
YYYY-MM-DD HH:MM:SS.mmm [LEVEL] [prefix] filename:line - message
```

**Example Output:**
```
2024-01-15 14:30:25.123 [INFO] [pip-sdk] manager.go:45 - Installing package: requests
2024-01-15 14:30:27.456 [INFO] [pip-sdk] operations.go:123 - Command executed: pip install requests (duration: 2.3s)
2024-01-15 14:30:27.457 [INFO] [pip-sdk] operations.go:89 - Package install successful: requests (duration: 2.3s)
```

## Advanced Usage

### Multiple Output Destinations

```go
// Log to both file and stdout
file, err := os.OpenFile("pip.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
if err != nil {
    return err
}

multiWriter := io.MultiWriter(os.Stdout, file)

logger, err := pip.NewLogger(&pip.LoggerConfig{
    Level:  pip.LogLevelInfo,
    Output: multiWriter,
})
```

### Structured Logging with JSON

```go
type JSONLogger struct {
    *pip.Logger
    encoder *json.Encoder
}

func (jl *JSONLogger) Info(format string, args ...interface{}) {
    entry := map[string]interface{}{
        "timestamp": time.Now().Format(time.RFC3339),
        "level":     "INFO",
        "message":   fmt.Sprintf(format, args...),
    }
    jl.encoder.Encode(entry)
}
```

### Conditional Logging

```go
func logIfVerbose(logger *pip.Logger, verbose bool, format string, args ...interface{}) {
    if verbose {
        logger.Info(format, args...)
    }
}
```

### Log Rotation

```go
import "gopkg.in/natefinch/lumberjack.v2"

func createRotatingLogger() (*pip.Logger, error) {
    rotator := &lumberjack.Logger{
        Filename:   "pip-sdk.log",
        MaxSize:    10, // MB
        MaxBackups: 3,
        MaxAge:     28, // days
        Compress:   true,
    }

    return pip.NewLogger(&pip.LoggerConfig{
        Level:  pip.LogLevelInfo,
        Output: rotator,
    })
}
```

## Performance Considerations

1. **Check log level before expensive operations**:
   ```go
   if logger.IsEnabled(pip.LogLevelDebug) {
       expensiveDebugData := generateDebugData()
       logger.Debug("Debug data: %v", expensiveDebugData)
   }
   ```

2. **Use appropriate log levels**:
   - `Debug`: Detailed debugging information
   - `Info`: General operational information
   - `Warn`: Warning conditions that don't prevent operation
   - `Error`: Error conditions that prevent operation
   - `Fatal`: Critical errors that require program termination

3. **Avoid logging in tight loops**:
   ```go
   // Bad
   for _, pkg := range packages {
       logger.Debug("Processing package: %s", pkg.Name)
       // ... process package
   }

   // Better
   logger.Debug("Processing %d packages", len(packages))
   for _, pkg := range packages {
       // ... process package
   }
   ```

## Thread Safety

The Logger is safe for concurrent use. Multiple goroutines can call logging methods simultaneously without additional synchronization.

## Best Practices

1. **Always close loggers that use files**:
   ```go
   logger, err := pip.NewLogger(config)
   if err != nil {
       return err
   }
   defer logger.Close()
   ```

2. **Use appropriate log levels for different environments**:
   ```go
   var logLevel pip.LogLevel
   if os.Getenv("DEBUG") == "true" {
       logLevel = pip.LogLevelDebug
   } else {
       logLevel = pip.LogLevelInfo
   }
   ```

3. **Include context in log messages**:
   ```go
   logger.Error("Failed to install package %s in environment %s: %v",
       pkg.Name, venvPath, err)
   ```

4. **Use structured logging for machine-readable logs**:
   ```go
   logger.Info("operation=install package=%s version=%s duration=%v success=%t",
       pkg.Name, pkg.Version, duration, success)
   ```
```
