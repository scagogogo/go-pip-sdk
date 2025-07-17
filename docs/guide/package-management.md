# Package Management

This guide covers how to manage Python packages using the Go Pip SDK.

## Installing Packages

### Basic Installation

```go
package main

import (
    "log"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    manager := pip.NewManager(nil)
    
    // Install a single package
    pkg := &pip.PackageSpec{
        Name: "requests",
    }
    
    err := manager.InstallPackage(pkg)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println("Package installed successfully")
}
```

### Installing with Version Constraints

```go
// Install specific version
pkg := &pip.PackageSpec{
    Name:    "requests",
    Version: "==2.28.0",
}

// Install with minimum version
pkg := &pip.PackageSpec{
    Name:    "requests",
    Version: ">=2.25.0",
}

// Install with version range
pkg := &pip.PackageSpec{
    Name:    "requests",
    Version: ">=2.25.0,<3.0.0",
}
```

### Installing with Extras

```go
pkg := &pip.PackageSpec{
    Name:   "requests",
    Extras: []string{"security", "socks"},
}
```

### Installing with Options

```go
pkg := &pip.PackageSpec{
    Name:           "requests",
    Upgrade:        true,
    ForceReinstall: true,
    Index:          "https://pypi.org/simple/",
    Options: map[string]string{
        "no-deps":     "",
        "target":      "/custom/path",
        "timeout":     "300",
    },
}
```

## Uninstalling Packages

```go
err := manager.UninstallPackage("requests")
if err != nil {
    log.Fatal(err)
}
```

## Listing Installed Packages

```go
packages, err := manager.ListPackages()
if err != nil {
    log.Fatal(err)
}

for _, pkg := range packages {
    fmt.Printf("%s %s\n", pkg.Name, pkg.Version)
}
```

## Getting Package Information

```go
info, err := manager.ShowPackage("requests")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Name: %s\n", info.Name)
fmt.Printf("Version: %s\n", info.Version)
fmt.Printf("Summary: %s\n", info.Summary)
fmt.Printf("Location: %s\n", info.Location)
```

## Freezing Dependencies

```go
packages, err := manager.FreezePackages()
if err != nil {
    log.Fatal(err)
}

// Write to requirements.txt
file, err := os.Create("requirements.txt")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

for _, pkg := range packages {
    fmt.Fprintf(file, "%s==%s\n", pkg.Name, pkg.Version)
}
```

## Installing from Requirements File

```go
// Using the project manager
projectManager := pip.NewProjectManager(manager)

err := projectManager.InstallRequirements("requirements.txt")
if err != nil {
    log.Fatal(err)
}
```

## Batch Operations

### Installing Multiple Packages

```go
packages := []*pip.PackageSpec{
    {Name: "requests"},
    {Name: "flask"},
    {Name: "django", Version: ">=3.0"},
}

for _, pkg := range packages {
    err := manager.InstallPackage(pkg)
    if err != nil {
        log.Printf("Failed to install %s: %v", pkg.Name, err)
        continue
    }
    log.Printf("Installed %s successfully", pkg.Name)
}
```

## Error Handling

The SDK provides detailed error information:

```go
err := manager.InstallPackage(&pip.PackageSpec{Name: "nonexistent-package"})
if err != nil {
    if pipErr, ok := err.(*pip.PipErrorDetails); ok {
        fmt.Printf("Error Type: %s\n", pipErr.Type)
        fmt.Printf("Message: %s\n", pipErr.Message)
        fmt.Printf("Command: %s\n", pipErr.Command)
        fmt.Printf("Exit Code: %d\n", pipErr.ExitCode)
        
        if len(pipErr.Suggestions) > 0 {
            fmt.Println("Suggestions:")
            for _, suggestion := range pipErr.Suggestions {
                fmt.Printf("  - %s\n", suggestion)
            }
        }
    }
}
```

## Using CLI Tool

The SDK includes a CLI tool for package management:

```bash
# Install packages
pip-cli install requests flask

# Install with version
pip-cli install "requests>=2.25.0"

# List packages
pip-cli list

# Show package info
pip-cli show requests

# Uninstall package
pip-cli uninstall requests

# Freeze dependencies
pip-cli freeze > requirements.txt
```

## Best Practices

1. **Use Virtual Environments**: Always work within virtual environments
2. **Pin Dependencies**: Use specific versions in production
3. **Handle Errors**: Always check for and handle errors appropriately
4. **Use Requirements Files**: Manage dependencies with requirements.txt
5. **Test Installations**: Verify packages work after installation

## Next Steps

- Learn about [Virtual Environments](./virtual-environments.md)
- Explore [Project Management](./project-management.md)
- Check [Error Handling](./error-handling.md) guide
