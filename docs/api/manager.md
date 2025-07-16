# Manager

The `Manager` is the central component of the Go Pip SDK, implementing the `PipManager` interface and providing all core functionality for Python package management.

## Constructor Functions

### NewManager

```go
func NewManager(config *Config) *Manager
```

Creates a new pip manager instance with the specified configuration.

**Parameters:**
- `config` (*Config): Configuration for the manager. If `nil`, default configuration is used.

**Returns:**
- `*Manager`: A new manager instance.

**Example:**
```go
// Create with default configuration
manager := pip.NewManager(nil)

// Create with custom configuration
config := &pip.Config{
    Timeout:  60 * time.Second,
    Retries:  5,
    LogLevel: "DEBUG",
}
manager := pip.NewManager(config)
```

### NewManagerWithContext

```go
func NewManagerWithContext(ctx context.Context, config *Config) *Manager
```

Creates a new pip manager instance with a context for cancellation and timeouts.

**Parameters:**
- `ctx` (context.Context): Context for operation cancellation and timeouts.
- `config` (*Config): Configuration for the manager. If `nil`, default configuration is used.

**Returns:**
- `*Manager`: A new manager instance with the specified context.

**Example:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

manager := pip.NewManagerWithContext(ctx, nil)
```

## Configuration Methods

### SetConfig

```go
func (m *Manager) SetConfig(config *Config)
```

Updates the manager configuration.

**Parameters:**
- `config` (*Config): New configuration to apply.

### GetConfig

```go
func (m *Manager) GetConfig() *Config
```

Returns the current manager configuration.

**Returns:**
- `*Config`: Current configuration.

### SetLogger

```go
func (m *Manager) SetLogger(logger *log.Logger)
```

Sets a custom standard library logger.

**Parameters:**
- `logger` (*log.Logger): Standard library logger instance.

### SetCustomLogger

```go
func (m *Manager) SetCustomLogger(logger *Logger)
```

Sets a custom SDK logger with advanced features.

**Parameters:**
- `logger` (*Logger): SDK logger instance with structured logging.

### SetContext

```go
func (m *Manager) SetContext(ctx context.Context)
```

Updates the context for operations.

**Parameters:**
- `ctx` (context.Context): New context for operations.

## System Operations

### IsInstalled

```go
func (m *Manager) IsInstalled() (bool, error)
```

Checks if pip is installed and accessible.

**Returns:**
- `bool`: `true` if pip is installed and accessible, `false` otherwise.
- `error`: Error if the check fails.

**Example:**
```go
installed, err := manager.IsInstalled()
if err != nil {
    return fmt.Errorf("failed to check pip installation: %w", err)
}

if !installed {
    fmt.Println("Pip is not installed")
}
```

### Install

```go
func (m *Manager) Install() error
```

Installs pip on the system using the most appropriate method for the current platform.

**Returns:**
- `error`: Error if installation fails.

**Installation Methods:**
- **Windows**: Uses `ensurepip` module or downloads `get-pip.py`
- **macOS**: Uses `ensurepip`, Homebrew, or MacPorts
- **Linux**: Uses system package manager (apt, yum, dnf, etc.) or `ensurepip`

**Example:**
```go
if err := manager.Install(); err != nil {
    return fmt.Errorf("failed to install pip: %w", err)
}
fmt.Println("Pip installed successfully")
```

### GetVersion

```go
func (m *Manager) GetVersion() (string, error)
```

Gets the installed pip version.

**Returns:**
- `string`: Pip version string (e.g., "23.2.1").
- `error`: Error if version retrieval fails.

**Example:**
```go
version, err := manager.GetVersion()
if err != nil {
    return fmt.Errorf("failed to get pip version: %w", err)
}
fmt.Printf("Pip version: %s\n", version)
```

## Utility Functions

### DefaultConfig

```go
func DefaultConfig() *Config
```

Returns a default configuration with sensible defaults.

**Returns:**
- `*Config`: Default configuration.

**Default Values:**
```go
&Config{
    Timeout:      30 * time.Second,
    Retries:      3,
    LogLevel:     "INFO",
    CacheDir:     "",
    Environment:  make(map[string]string),
    ExtraOptions: make(map[string]string),
}
```

### GetOSType

```go
func GetOSType() OSType
```

Detects the current operating system type.

**Returns:**
- `OSType`: Operating system type (`OSWindows`, `OSMacOS`, `OSLinux`, or `OSUnknown`).

**Example:**
```go
osType := pip.GetOSType()
switch osType {
case pip.OSWindows:
    fmt.Println("Running on Windows")
case pip.OSMacOS:
    fmt.Println("Running on macOS")
case pip.OSLinux:
    fmt.Println("Running on Linux")
default:
    fmt.Println("Unknown operating system")
}
```

## Thread Safety

The `Manager` type is safe for concurrent use. All methods can be called simultaneously from multiple goroutines.

## Error Handling

All methods return structured errors that can be type-checked:

```go
if err := manager.Install(); err != nil {
    if pip.IsErrorType(err, pip.ErrorTypePipNotInstalled) {
        // Handle pip not installed error
    } else if pip.IsErrorType(err, pip.ErrorTypePermissionDenied) {
        // Handle permission error
    }
    return err
}
```

## Best Practices

1. **Always check if pip is installed** before performing operations:
   ```go
   if installed, err := manager.IsInstalled(); err != nil || !installed {
       if err := manager.Install(); err != nil {
           return err
       }
   }
   ```

2. **Use contexts for long-running operations**:
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
   defer cancel()
   manager := pip.NewManagerWithContext(ctx, nil)
   ```

3. **Configure appropriate timeouts and retries**:
   ```go
   config := &pip.Config{
       Timeout: 120 * time.Second,
       Retries: 5,
   }
   ```
