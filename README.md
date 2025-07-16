# Go Pip SDK

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.19-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/go-pip-sdk)](https://goreportcard.com/report/github.com/scagogogo/go-pip-sdk)
[![Documentation](https://img.shields.io/badge/docs-online-blue.svg)](https://scagogogo.github.io/go-pip-sdk/)

A comprehensive Go SDK for managing Python pip operations, virtual environments, and Python projects. This library provides a clean, type-safe interface for all common pip operations with cross-platform support.

English | [简体中文](README.zh-CN.md)

## Features

- 🚀 **Cross-platform support** - Works on Windows, macOS, and Linux
- 📦 **Complete pip operations** - Install, uninstall, list, show, freeze packages
- 🐍 **Virtual environment management** - Create, activate, deactivate, remove virtual environments
- 🏗️ **Project initialization** - Bootstrap Python projects with standard structure
- 🔧 **Automatic pip installation** - Detects and installs pip if missing
- 📝 **Comprehensive logging** - Detailed operation logs with multiple levels
- ⚡ **Error handling** - Rich error types with helpful suggestions
- 🧪 **Well tested** - Extensive unit and integration tests

## Installation

```bash
go get github.com/scagogogo/go-pip-sdk
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    // Create a new pip manager
    manager := pip.NewManager(nil) // Uses default configuration
    
    // Check if pip is installed
    installed, err := manager.IsInstalled()
    if err != nil {
        log.Fatal(err)
    }
    
    if !installed {
        fmt.Println("Installing pip...")
        if err := manager.Install(); err != nil {
            log.Fatal(err)
        }
    }
    
    // Install a package
    pkg := &pip.PackageSpec{
        Name:    "requests",
        Version: ">=2.25.0",
    }
    
    if err := manager.InstallPackage(pkg); err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Package installed successfully!")
}
```

## Core Operations

### Package Management

```go
// Install a package
pkg := &pip.PackageSpec{
    Name:    "requests",
    Version: ">=2.25.0",
    Extras:  []string{"security", "socks"},
    Upgrade: true,
}
err := manager.InstallPackage(pkg)

// Uninstall a package
err = manager.UninstallPackage("requests")

// List installed packages
packages, err := manager.ListPackages()

// Show package information
info, err := manager.ShowPackage("requests")

// Freeze packages (like pip freeze)
packages, err := manager.FreezePackages()
```

### Virtual Environment Management

```go
// Create a virtual environment
err := manager.CreateVenv("/path/to/venv")

// Activate a virtual environment
err = manager.ActivateVenv("/path/to/venv")

// Deactivate current virtual environment
err = manager.DeactivateVenv()

// Remove a virtual environment
err = manager.RemoveVenv("/path/to/venv")

// Get virtual environment information
info, err := manager.GetVenvInfo("/path/to/venv")
```

### Project Initialization

```go
// Initialize a new Python project
opts := &pip.ProjectOptions{
    Name:            "my-project",
    Version:         "0.1.0",
    Description:     "My awesome Python project",
    Author:          "Your Name",
    AuthorEmail:     "your.email@example.com",
    License:         "MIT",
    Dependencies:    []string{"requests>=2.25.0", "click>=8.0.0"},
    DevDependencies: []string{"pytest>=6.0.0", "black>=21.0.0"},
    CreateVenv:      true,
    VenvPath:        "./venv",
}

err := manager.InitProject("/path/to/project", opts)
```

## Configuration

```go
config := &pip.Config{
    PythonPath:   "/usr/bin/python3",
    PipPath:      "/usr/bin/pip3",
    Timeout:      60 * time.Second,
    Retries:      3,
    LogLevel:     "INFO",
    DefaultIndex: "https://pypi.org/simple/",
    TrustedHosts: []string{"pypi.org", "pypi.python.org"},
    Environment: map[string]string{
        "PIP_CACHE_DIR": "/tmp/pip-cache",
    },
}

manager := pip.NewManager(config)
```

## Advanced Usage

### Custom Logging

```go
// Create custom logger
loggerConfig := &pip.LoggerConfig{
    Level:      pip.LogLevelDebug,
    EnableFile: true,
    LogFile:    "/var/log/pip-sdk.log",
}

logger, err := pip.NewLogger(loggerConfig)
if err != nil {
    log.Fatal(err)
}
defer logger.Close()

// Set custom logger
manager.SetCustomLogger(logger)
```

### Error Handling

```go
err := manager.InstallPackage(pkg)
if err != nil {
    if pip.IsErrorType(err, pip.ErrorTypePermissionDenied) {
        fmt.Println("Permission denied. Try running with elevated privileges.")
    } else if pip.IsErrorType(err, pip.ErrorTypeNetworkError) {
        fmt.Println("Network error. Check your internet connection.")
    } else {
        fmt.Printf("Installation failed: %v\n", err)
    }
}
```

### Context Support

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

manager := pip.NewManagerWithContext(ctx, nil)
```

## Documentation

- 📖 **[Online Documentation](https://scagogogo.github.io/go-pip-sdk/)** - Complete API documentation and guides
- 🚀 **[Getting Started](https://scagogogo.github.io/go-pip-sdk/guide/getting-started)** - Quick start guide
- 📚 **[API Reference](https://scagogogo.github.io/go-pip-sdk/api/)** - Detailed API documentation
- 💡 **[Examples](https://scagogogo.github.io/go-pip-sdk/examples/)** - Code examples and use cases

## Examples

See the [examples](examples/) directory for more comprehensive examples:

- [Basic Usage](examples/basic/main.go)
- [Virtual Environment Management](examples/venv/main.go)
- [Project Initialization](examples/project/main.go)
- [Advanced Configuration](examples/advanced/main.go)

## API Reference

### Types

- `Manager` - Main interface for pip operations
- `PackageSpec` - Specification for package installation
- `Package` - Represents an installed package
- `PackageInfo` - Detailed package information
- `ProjectOptions` - Options for project initialization
- `Config` - Manager configuration
- `Logger` - Custom logging interface

### Error Types

- `ErrorTypePipNotInstalled` - Pip is not installed
- `ErrorTypePythonNotFound` - Python interpreter not found
- `ErrorTypePackageNotFound` - Package not found
- `ErrorTypePermissionDenied` - Permission denied
- `ErrorTypeNetworkError` - Network connectivity issues
- `ErrorTypeCommandFailed` - Command execution failed

## Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run only unit tests (skip integration tests)
go test -short ./...

# Run integration tests (requires pip installation)
go test -run Integration ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Requirements

- Go 1.19 or later
- Python 3.7 or later (for pip operations)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by the Python pip package manager
- Built with Go's excellent standard library
- Thanks to all contributors and users

## Support

- 📖 [Documentation](https://scagogogo.github.io/go-pip-sdk/)
- 🐛 [Issue Tracker](https://github.com/scagogogo/go-pip-sdk/issues)
- 💬 [Discussions](https://github.com/scagogogo/go-pip-sdk/discussions)
