# Getting Started

Welcome to the Go Pip SDK! This guide will help you get up and running with the SDK in just a few minutes.

## Prerequisites

Before you begin, ensure you have:

- **Go 1.19 or later** installed on your system
- **Python 3.7 or later** (for pip operations)
- Basic familiarity with Go programming

## Installation

Install the SDK using Go modules:

```bash
go get github.com/scagogogo/go-pip-sdk
```

## Your First Program

Let's create a simple program that checks if pip is installed and installs a package:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    // Create a new pip manager with default configuration
    manager := pip.NewManager(nil)
    
    // Check if pip is installed
    installed, err := manager.IsInstalled()
    if err != nil {
        log.Fatalf("Failed to check pip installation: %v", err)
    }
    
    if !installed {
        fmt.Println("Pip is not installed. Installing...")
        if err := manager.Install(); err != nil {
            log.Fatalf("Failed to install pip: %v", err)
        }
        fmt.Println("Pip installed successfully!")
    } else {
        fmt.Println("Pip is already installed.")
    }
    
    // Get pip version
    version, err := manager.GetVersion()
    if err != nil {
        log.Fatalf("Failed to get pip version: %v", err)
    }
    fmt.Printf("Pip version: %s\n", version)
    
    // Install a package
    pkg := &pip.PackageSpec{
        Name:    "requests",
        Version: ">=2.25.0",
    }
    
    fmt.Printf("Installing package: %s\n", pkg.Name)
    if err := manager.InstallPackage(pkg); err != nil {
        log.Fatalf("Failed to install package: %v", err)
    }
    
    fmt.Println("Package installed successfully!")
    
    // List installed packages
    packages, err := manager.ListPackages()
    if err != nil {
        log.Fatalf("Failed to list packages: %v", err)
    }
    
    fmt.Printf("Found %d installed packages:\n", len(packages))
    for _, pkg := range packages {
        fmt.Printf("  %s %s\n", pkg.Name, pkg.Version)
    }
}
```

Save this as `main.go` and run it:

```bash
go run main.go
```

## Basic Concepts

### Manager

The `Manager` is the central component that provides all pip functionality. It implements the `PipManager` interface and handles:

- System operations (checking pip installation, getting version)
- Package operations (install, uninstall, list, show)
- Virtual environment operations
- Project management

### Configuration

You can customize the manager's behavior using a `Config` struct:

```go
config := &pip.Config{
    PythonPath:   "/usr/bin/python3",
    PipPath:      "/usr/bin/pip3",
    Timeout:      60 * time.Second,
    Retries:      3,
    LogLevel:     "DEBUG",
    DefaultIndex: "https://pypi.org/simple/",
}

manager := pip.NewManager(config)
```

### Package Specifications

When installing packages, you use `PackageSpec` to specify requirements:

```go
// Basic package
pkg := &pip.PackageSpec{
    Name: "requests",
}

// Package with version constraint
pkg := &pip.PackageSpec{
    Name:    "django",
    Version: ">=4.0,<5.0",
}

// Package with extras
pkg := &pip.PackageSpec{
    Name:   "fastapi",
    Extras: []string{"dev", "test"},
}

// Package with custom options
pkg := &pip.PackageSpec{
    Name:    "numpy",
    Upgrade: true,
    Options: map[string]string{
        "no-cache-dir": "",
        "timeout":      "120",
    },
}
```

## Common Operations

### Installing Packages

```go
// Install a single package
pkg := &pip.PackageSpec{Name: "requests"}
if err := manager.InstallPackage(pkg); err != nil {
    return err
}

// Install multiple packages
packages := []*pip.PackageSpec{
    {Name: "requests"},
    {Name: "click"},
    {Name: "pydantic"},
}

for _, pkg := range packages {
    if err := manager.InstallPackage(pkg); err != nil {
        fmt.Printf("Failed to install %s: %v\n", pkg.Name, err)
    }
}
```

### Working with Virtual Environments

```go
// Create a virtual environment
venvPath := "/path/to/my-venv"
if err := manager.CreateVenv(venvPath); err != nil {
    return err
}

// Activate the virtual environment
if err := manager.ActivateVenv(venvPath); err != nil {
    return err
}

// Install packages in the virtual environment
pkg := &pip.PackageSpec{Name: "requests"}
if err := manager.InstallPackage(pkg); err != nil {
    return err
}

// Deactivate when done
if err := manager.DeactivateVenv(); err != nil {
    return err
}
```

### Getting Package Information

```go
// List all installed packages
packages, err := manager.ListPackages()
if err != nil {
    return err
}

for _, pkg := range packages {
    fmt.Printf("%s %s\n", pkg.Name, pkg.Version)
}

// Get detailed information about a package
info, err := manager.ShowPackage("requests")
if err != nil {
    return err
}

fmt.Printf("Name: %s\n", info.Name)
fmt.Printf("Version: %s\n", info.Version)
fmt.Printf("Summary: %s\n", info.Summary)
fmt.Printf("Dependencies: %v\n", info.Requires)
```

## Error Handling

The SDK provides structured error handling with specific error types:

```go
if err := manager.InstallPackage(pkg); err != nil {
    switch pip.GetErrorType(err) {
    case pip.ErrorTypePackageNotFound:
        fmt.Printf("Package %s not found\n", pkg.Name)
    case pip.ErrorTypePermissionDenied:
        fmt.Println("Permission denied - try running with elevated privileges")
    case pip.ErrorTypeNetworkError:
        fmt.Println("Network error - check your internet connection")
    default:
        fmt.Printf("Installation failed: %v\n", err)
    }
}
```

## Logging

Enable logging to see what the SDK is doing:

```go
// Create a logger
logger, err := pip.NewLogger(&pip.LoggerConfig{
    Level:  pip.LogLevelInfo,
    Output: os.Stdout,
    Prefix: "[my-app]",
})
if err != nil {
    return err
}
defer logger.Close()

// Set the logger on the manager
manager := pip.NewManager(nil)
manager.SetCustomLogger(logger)

// Now all operations will be logged
manager.InstallPackage(&pip.PackageSpec{Name: "requests"})
```

## Next Steps

Now that you have the basics down, explore these topics:

- [Configuration](/guide/configuration) - Learn about all configuration options
- [Package Management](/guide/package-management) - Advanced package operations
- [Virtual Environments](/guide/virtual-environments) - Working with virtual environments
- [Project Management](/guide/project-management) - Initializing Python projects
- [API Reference](/api/) - Complete API documentation
- [Examples](/examples/) - More code examples

## Getting Help

If you run into issues:

1. Check the [API documentation](/api/) for detailed information
2. Look at the [examples](/examples/) for common use cases
3. Search the [issue tracker](https://github.com/scagogogo/go-pip-sdk/issues)
4. Ask questions in [discussions](https://github.com/scagogogo/go-pip-sdk/discussions)
