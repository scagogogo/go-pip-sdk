# API Reference

The Go Pip SDK provides a comprehensive set of APIs for managing Python packages, virtual environments, and projects. This section documents all public interfaces, types, and functions.

## Core Components

### [Manager](/api/manager)
The main interface for pip operations. Provides methods for package management, virtual environment operations, and project initialization.

### [Package Operations](/api/package-operations)
Functions for installing, uninstalling, listing, and managing Python packages.

### [Virtual Environments](/api/virtual-environments)
APIs for creating, activating, and managing Python virtual environments.

### [Project Management](/api/project-management)
Tools for initializing Python projects with standard structure and configuration files.

### [Types](/api/types)
Core data structures and interfaces used throughout the SDK.

### [Errors](/api/errors)
Error types and handling mechanisms for robust error management.

### [Logger](/api/logger)
Logging system with configurable levels and output formats.

### [Installer](/api/installer)
Cross-platform pip installation functionality.

## Quick Reference

### Main Interface

```go
type PipManager interface {
    // System operations
    IsInstalled() (bool, error)
    Install() error
    GetVersion() (string, error)

    // Package operations
    InstallPackage(pkg *PackageSpec) error
    UninstallPackage(name string) error
    ListPackages() ([]*Package, error)
    ShowPackage(name string) (*PackageInfo, error)
    SearchPackages(query string) ([]*SearchResult, error)
    FreezePackages() ([]*Package, error)

    // Virtual environment operations
    CreateVenv(path string) error
    ActivateVenv(path string) error
    DeactivateVenv() error
    RemoveVenv(path string) error

    // Project operations
    InitProject(path string, opts *ProjectOptions) error
    InstallRequirements(path string) error
    GenerateRequirements(path string) error
}
```

### Key Types

```go
// Manager configuration
type Config struct {
    PythonPath   string
    PipPath      string
    DefaultIndex string
    TrustedHosts []string
    Timeout      time.Duration
    Retries      int
    LogLevel     string
    CacheDir     string
    ExtraOptions map[string]string
    Environment  map[string]string
}

// Package specification
type PackageSpec struct {
    Name           string
    Version        string
    Extras         []string
    Index          string
    Options        map[string]string
    Editable       bool
    Upgrade        bool
    ForceReinstall bool
}

// Installed package information
type Package struct {
    Name      string
    Version   string
    Location  string
    Editable  bool
    Installer string
}
```

## Error Handling

All API functions return errors that implement the standard Go error interface. The SDK provides rich error types with additional context:

```go
type PipErrorDetails struct {
    Type        ErrorType
    Message     string
    Command     string
    Output      string
    ExitCode    int
    Suggestions []string
    Context     map[string]string
    Cause       error
}
```

## Usage Patterns

### Basic Usage

```go
// Create manager with default configuration
manager := pip.NewManager(nil)

// Check if pip is installed
installed, err := manager.IsInstalled()
if err != nil {
    return err
}

if !installed {
    // Install pip if not available
    if err := manager.Install(); err != nil {
        return err
    }
}
```

### With Custom Configuration

```go
config := &pip.Config{
    Timeout:      60 * time.Second,
    Retries:      5,
    DefaultIndex: "https://pypi.org/simple/",
    LogLevel:     "DEBUG",
}

manager := pip.NewManager(config)
```

### With Context

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

manager := pip.NewManagerWithContext(ctx, nil)
```

## Thread Safety

The Manager type is safe for concurrent use. Multiple goroutines can call methods on the same Manager instance simultaneously.

## Platform Support

All APIs work consistently across supported platforms:
- Windows (x86, x64, ARM64)
- macOS (Intel, Apple Silicon)
- Linux (x86, x64, ARM, ARM64)

Platform-specific behavior is handled automatically by the SDK.
