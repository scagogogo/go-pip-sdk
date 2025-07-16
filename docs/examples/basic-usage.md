# Basic Usage

This guide covers the fundamental operations of the Go Pip SDK with practical examples.

## Installation and Setup

First, install the SDK:

```bash
go get github.com/scagogogo/go-pip-sdk
```

Import the package in your Go code:

```go
import "github.com/scagogogo/go-pip-sdk/pkg/pip"
```

## Creating a Manager

### Default Manager

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    // Create manager with default configuration
    manager := pip.NewManager(nil)
    
    // Check if pip is available
    installed, err := manager.IsInstalled()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Pip installed: %t\n", installed)
}
```

### Custom Configuration

```go
func main() {
    // Create custom configuration
    config := &pip.Config{
        Timeout:  60 * time.Second,
        Retries:  5,
        LogLevel: "DEBUG",
    }
    
    manager := pip.NewManager(config)
    
    // Use the manager...
}
```

## Basic Package Operations

### Installing Packages

```go
func installPackage() error {
    manager := pip.NewManager(nil)
    
    // Simple package installation
    pkg := &pip.PackageSpec{
        Name: "requests",
    }
    
    if err := manager.InstallPackage(pkg); err != nil {
        return fmt.Errorf("failed to install package: %w", err)
    }
    
    fmt.Println("Package installed successfully!")
    return nil
}
```

### Installing with Version Constraints

```go
func installWithVersion() error {
    manager := pip.NewManager(nil)
    
    pkg := &pip.PackageSpec{
        Name:    "django",
        Version: ">=4.0,<5.0",
    }
    
    return manager.InstallPackage(pkg)
}
```

### Installing with Extras

```go
func installWithExtras() error {
    manager := pip.NewManager(nil)
    
    pkg := &pip.PackageSpec{
        Name:   "fastapi",
        Extras: []string{"dev", "test"},
    }
    
    return manager.InstallPackage(pkg)
}
```

### Uninstalling Packages

```go
func uninstallPackage() error {
    manager := pip.NewManager(nil)
    
    if err := manager.UninstallPackage("requests"); err != nil {
        return fmt.Errorf("failed to uninstall package: %w", err)
    }
    
    fmt.Println("Package uninstalled successfully!")
    return nil
}
```

## Listing and Information

### List All Packages

```go
func listPackages() error {
    manager := pip.NewManager(nil)
    
    packages, err := manager.ListPackages()
    if err != nil {
        return err
    }
    
    fmt.Printf("Found %d installed packages:\n", len(packages))
    for _, pkg := range packages {
        fmt.Printf("  %s %s\n", pkg.Name, pkg.Version)
    }
    
    return nil
}
```

### Get Package Information

```go
func getPackageInfo() error {
    manager := pip.NewManager(nil)
    
    info, err := manager.ShowPackage("requests")
    if err != nil {
        return err
    }
    
    fmt.Printf("Package: %s\n", info.Name)
    fmt.Printf("Version: %s\n", info.Version)
    fmt.Printf("Summary: %s\n", info.Summary)
    fmt.Printf("Author: %s\n", info.Author)
    fmt.Printf("License: %s\n", info.License)
    fmt.Printf("Dependencies: %v\n", info.Requires)
    
    return nil
}
```

### Freeze Packages

```go
func freezePackages() error {
    manager := pip.NewManager(nil)
    
    packages, err := manager.FreezePackages()
    if err != nil {
        return err
    }
    
    fmt.Println("Frozen packages:")
    for _, pkg := range packages {
        fmt.Printf("%s==%s\n", pkg.Name, pkg.Version)
    }
    
    return nil
}
```

## Error Handling

### Basic Error Handling

```go
func handleErrors() {
    manager := pip.NewManager(nil)
    
    pkg := &pip.PackageSpec{Name: "nonexistent-package"}
    
    if err := manager.InstallPackage(pkg); err != nil {
        fmt.Printf("Installation failed: %v\n", err)
        
        // Check specific error type
        if pip.IsErrorType(err, pip.ErrorTypePackageNotFound) {
            fmt.Println("Package not found - check the name")
        }
    }
}
```

### Comprehensive Error Handling

```go
func comprehensiveErrorHandling() {
    manager := pip.NewManager(nil)
    
    pkg := &pip.PackageSpec{Name: "some-package"}
    
    if err := manager.InstallPackage(pkg); err != nil {
        switch pip.GetErrorType(err) {
        case pip.ErrorTypePackageNotFound:
            fmt.Println("Package not found")
        case pip.ErrorTypePermissionDenied:
            fmt.Println("Permission denied - try running with elevated privileges")
        case pip.ErrorTypeNetworkError:
            fmt.Println("Network error - check your internet connection")
        case pip.ErrorTypeCommandFailed:
            fmt.Println("Pip command failed")
        default:
            fmt.Printf("Unknown error: %v\n", err)
        }
    }
}
```

## Working with Requirements Files

### Installing from Requirements

```go
func installFromRequirements() error {
    manager := pip.NewManager(nil)
    
    if err := manager.InstallRequirements("requirements.txt"); err != nil {
        return fmt.Errorf("failed to install requirements: %w", err)
    }
    
    fmt.Println("Requirements installed successfully!")
    return nil
}
```

### Generating Requirements

```go
func generateRequirements() error {
    manager := pip.NewManager(nil)
    
    if err := manager.GenerateRequirements("requirements.txt"); err != nil {
        return fmt.Errorf("failed to generate requirements: %w", err)
    }
    
    fmt.Println("Requirements file generated!")
    return nil
}
```

## Complete Example

Here's a complete example that demonstrates multiple operations:

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    // Create manager
    manager := pip.NewManager(nil)
    
    // Check pip installation
    if installed, err := manager.IsInstalled(); err != nil || !installed {
        fmt.Println("Installing pip...")
        if err := manager.Install(); err != nil {
            log.Fatalf("Failed to install pip: %v", err)
        }
    }
    
    // Get pip version
    version, err := manager.GetVersion()
    if err != nil {
        log.Fatalf("Failed to get pip version: %v", err)
    }
    fmt.Printf("Using pip version: %s\n", version)
    
    // Install packages
    packages := []*pip.PackageSpec{
        {Name: "requests", Version: ">=2.25.0"},
        {Name: "click", Version: ">=7.0"},
    }
    
    for _, pkg := range packages {
        fmt.Printf("Installing %s...\n", pkg.Name)
        if err := manager.InstallPackage(pkg); err != nil {
            fmt.Printf("Failed to install %s: %v\n", pkg.Name, err)
            continue
        }
        fmt.Printf("✓ %s installed\n", pkg.Name)
    }
    
    // List installed packages
    installed, err := manager.ListPackages()
    if err != nil {
        log.Fatalf("Failed to list packages: %v", err)
    }
    
    fmt.Printf("\nInstalled packages (%d):\n", len(installed))
    for _, pkg := range installed {
        fmt.Printf("  %s %s\n", pkg.Name, pkg.Version)
    }
    
    // Generate requirements file
    if err := manager.GenerateRequirements("requirements.txt"); err != nil {
        log.Printf("Failed to generate requirements: %v", err)
    } else {
        fmt.Println("\n✓ Requirements file generated")
    }
}
```

## Best Practices

1. **Always check for errors**:
   ```go
   if err := manager.InstallPackage(pkg); err != nil {
       return fmt.Errorf("installation failed: %w", err)
   }
   ```

2. **Use specific error handling**:
   ```go
   if pip.IsErrorType(err, pip.ErrorTypePackageNotFound) {
       // Handle package not found specifically
   }
   ```

3. **Configure appropriate timeouts**:
   ```go
   config := &pip.Config{
       Timeout: 120 * time.Second,
   }
   ```

4. **Validate package specifications**:
   ```go
   if pkg.Name == "" {
       return errors.New("package name cannot be empty")
   }
   ```

## Next Steps

- [Package Management](/examples/package-management) - Advanced package operations
- [Virtual Environments](/examples/virtual-environments) - Working with virtual environments
- [Project Initialization](/examples/project-initialization) - Setting up Python projects
