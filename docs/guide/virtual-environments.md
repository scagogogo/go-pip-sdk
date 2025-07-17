# Virtual Environments

This guide covers how to manage Python virtual environments using the Go Pip SDK.

## Overview

Virtual environments are isolated Python environments that allow you to install packages without affecting the system Python installation. The Go Pip SDK provides comprehensive virtual environment management capabilities.

## Creating Virtual Environments

### Basic Creation

```go
package main

import (
    "log"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    manager := pip.NewManager(nil)
    venvManager := pip.NewVenvManager(manager)
    
    // Create a new virtual environment
    err := venvManager.CreateVenv("/path/to/myenv")
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println("Virtual environment created successfully")
}
```

### Creating with Custom Python

```go
config := &pip.Config{
    PythonPath: "/usr/local/bin/python3.9",
}
manager := pip.NewManager(config)
venvManager := pip.NewVenvManager(manager)

err := venvManager.CreateVenv("/path/to/myenv")
```

## Activating Virtual Environments

```go
err := venvManager.ActivateVenv("/path/to/myenv")
if err != nil {
    log.Fatal(err)
}

log.Println("Virtual environment activated")
```

## Deactivating Virtual Environments

```go
err := venvManager.DeactivateVenv()
if err != nil {
    log.Fatal(err)
}

log.Println("Virtual environment deactivated")
```

## Removing Virtual Environments

```go
err := venvManager.RemoveVenv("/path/to/myenv")
if err != nil {
    log.Fatal(err)
}

log.Println("Virtual environment removed")
```

## Getting Virtual Environment Information

```go
info, err := venvManager.GetVenvInfo("/path/to/myenv")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Path: %s\n", info.Path)
fmt.Printf("Active: %t\n", info.IsActive)
fmt.Printf("Python Path: %s\n", info.PythonPath)
fmt.Printf("Pip Path: %s\n", info.PipPath)
```

## Working with Packages in Virtual Environments

Once a virtual environment is activated, all package operations will be performed within that environment:

```go
// Activate virtual environment
err := venvManager.ActivateVenv("/path/to/myenv")
if err != nil {
    log.Fatal(err)
}

// Install packages in the virtual environment
pkg := &pip.PackageSpec{Name: "requests"}
err = manager.InstallPackage(pkg)
if err != nil {
    log.Fatal(err)
}

// List packages in the virtual environment
packages, err := manager.ListPackages()
if err != nil {
    log.Fatal(err)
}

for _, pkg := range packages {
    fmt.Printf("%s %s\n", pkg.Name, pkg.Version)
}
```

## Complete Workflow Example

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    manager := pip.NewManager(nil)
    venvManager := pip.NewVenvManager(manager)
    
    venvPath := "./myproject-env"
    
    // 1. Create virtual environment
    fmt.Println("Creating virtual environment...")
    err := venvManager.CreateVenv(venvPath)
    if err != nil {
        log.Fatal(err)
    }
    
    // 2. Activate virtual environment
    fmt.Println("Activating virtual environment...")
    err = venvManager.ActivateVenv(venvPath)
    if err != nil {
        log.Fatal(err)
    }
    
    // 3. Install packages
    fmt.Println("Installing packages...")
    packages := []string{"requests", "flask", "pytest"}
    for _, pkgName := range packages {
        pkg := &pip.PackageSpec{Name: pkgName}
        err = manager.InstallPackage(pkg)
        if err != nil {
            log.Printf("Failed to install %s: %v", pkgName, err)
            continue
        }
        fmt.Printf("✓ Installed %s\n", pkgName)
    }
    
    // 4. List installed packages
    fmt.Println("\nInstalled packages:")
    installedPkgs, err := manager.ListPackages()
    if err != nil {
        log.Fatal(err)
    }
    
    for _, pkg := range installedPkgs {
        fmt.Printf("  %s %s\n", pkg.Name, pkg.Version)
    }
    
    // 5. Deactivate when done
    fmt.Println("\nDeactivating virtual environment...")
    err = venvManager.DeactivateVenv()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Workflow completed successfully!")
}
```

## Using CLI Tool

The SDK includes CLI commands for virtual environment management:

```bash
# Create virtual environment
pip-cli venv create ./myenv

# Get virtual environment info
pip-cli venv info ./myenv

# Activate virtual environment
pip-cli venv activate ./myenv

# Install packages (in activated environment)
pip-cli install requests flask

# List packages
pip-cli list

# Deactivate virtual environment
pip-cli venv deactivate

# Remove virtual environment
pip-cli venv remove ./myenv
```

## Virtual Environment Validation

```go
// Check if a directory is a valid virtual environment
isValid := venvManager.IsVenvValid("/path/to/myenv")
if isValid {
    fmt.Println("Valid virtual environment")
} else {
    fmt.Println("Not a valid virtual environment")
}

// Check if a virtual environment is currently active
isActive := venvManager.IsVenvActive()
if isActive {
    fmt.Println("A virtual environment is currently active")
}
```

## Error Handling

```go
err := venvManager.CreateVenv("/invalid/path/myenv")
if err != nil {
    if pipErr, ok := err.(*pip.PipErrorDetails); ok {
        switch pipErr.Type {
        case pip.ErrorTypePermissionDenied:
            fmt.Println("Permission denied - check directory permissions")
        case pip.ErrorTypePythonNotFound:
            fmt.Println("Python not found - check Python installation")
        default:
            fmt.Printf("Error: %s\n", pipErr.Message)
        }
    }
}
```

## Best Practices

1. **Use Project-Specific Environments**: Create a separate virtual environment for each project
2. **Descriptive Names**: Use descriptive names for virtual environments
3. **Requirements Files**: Generate and maintain requirements.txt files
4. **Clean Up**: Remove unused virtual environments to save disk space
5. **Activate Before Operations**: Always activate the environment before installing packages

## Integration with Project Management

```go
// Create a project with virtual environment
projectManager := pip.NewProjectManager(manager)

opts := &pip.ProjectOptions{
    Name:       "myproject",
    CreateVenv: true,
    VenvPath:   "./myproject-env",
}

err := projectManager.InitProject("./myproject", opts)
if err != nil {
    log.Fatal(err)
}
```

## Next Steps

- Learn about [Project Management](./project-management.md)
- Explore [Error Handling](./error-handling.md)
- Check out [Examples](../examples/) for more use cases
