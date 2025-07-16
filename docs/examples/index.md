# Examples

This section provides practical examples of using the Go Pip SDK for various tasks. Each example includes complete, runnable code with explanations.

## Quick Examples

### Basic Package Installation

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    manager := pip.NewManager(nil)
    
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

### Virtual Environment Setup

```go
package main

import (
    "fmt"
    "log"
    "path/filepath"
    
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    manager := pip.NewManager(nil)
    
    // Create virtual environment
    venvPath := filepath.Join(".", "my-venv")
    if err := manager.CreateVenv(venvPath); err != nil {
        log.Fatal(err)
    }
    
    // Activate virtual environment
    if err := manager.ActivateVenv(venvPath); err != nil {
        log.Fatal(err)
    }
    
    // Install packages in virtual environment
    packages := []*pip.PackageSpec{
        {Name: "requests"},
        {Name: "click"},
        {Name: "pydantic"},
    }
    
    for _, pkg := range packages {
        if err := manager.InstallPackage(pkg); err != nil {
            fmt.Printf("Failed to install %s: %v\n", pkg.Name, err)
        } else {
            fmt.Printf("Installed %s\n", pkg.Name)
        }
    }
    
    fmt.Println("Virtual environment setup complete!")
}
```

### Project Initialization

```go
package main

import (
    "log"
    
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    manager := pip.NewManager(nil)
    
    opts := &pip.ProjectOptions{
        Name:            "my-awesome-project",
        Version:         "0.1.0",
        Description:     "An awesome Python project",
        Author:          "Your Name",
        AuthorEmail:     "your.email@example.com",
        License:         "MIT",
        Dependencies:    []string{"requests>=2.25.0", "click>=7.0"},
        DevDependencies: []string{"pytest>=6.0", "black>=21.0"},
        CreateVenv:      true,
        VenvPath:        "./venv",
    }
    
    if err := manager.InitProject("./my-project", opts); err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Project initialized successfully!")
}
```

## Detailed Examples

### [Basic Usage](/examples/basic-usage)
Learn the fundamentals of the SDK with simple package management operations.

### [Package Management](/examples/package-management)
Advanced package installation, version management, and dependency handling.

### [Virtual Environments](/examples/virtual-environments)
Complete guide to creating and managing Python virtual environments.

### [Project Initialization](/examples/project-initialization)
Bootstrap new Python projects with proper structure and configuration.

### [Advanced Usage](/examples/advanced-usage)
Complex scenarios including error handling, logging, and custom configurations.

## Common Patterns

### Error Handling Pattern

```go
func installPackageWithErrorHandling(manager *pip.Manager, pkg *pip.PackageSpec) error {
    if err := manager.InstallPackage(pkg); err != nil {
        switch pip.GetErrorType(err) {
        case pip.ErrorTypePackageNotFound:
            return fmt.Errorf("package %s not found - check the name", pkg.Name)
        case pip.ErrorTypePermissionDenied:
            return fmt.Errorf("permission denied - try running with elevated privileges")
        case pip.ErrorTypeNetworkError:
            return fmt.Errorf("network error - check your internet connection")
        default:
            return fmt.Errorf("installation failed: %w", err)
        }
    }
    return nil
}
```

### Retry Pattern

```go
func installWithRetry(manager *pip.Manager, pkg *pip.PackageSpec, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        err := manager.InstallPackage(pkg)
        if err == nil {
            return nil
        }
        
        if pip.IsErrorType(err, pip.ErrorTypeNetworkError) {
            fmt.Printf("Network error, retrying... (%d/%d)\n", i+1, maxRetries)
            time.Sleep(time.Duration(i+1) * time.Second)
            continue
        }
        
        return err // Non-recoverable error
    }
    return fmt.Errorf("failed after %d retries", maxRetries)
}
```

### Logging Pattern

```go
func setupLogging() (*pip.Logger, error) {
    return pip.NewLogger(&pip.LoggerConfig{
        Level:      pip.LogLevelInfo,
        Output:     os.Stdout,
        Prefix:     "[my-app]",
        EnableFile: true,
        LogFile:    "pip-operations.log",
    })
}

func main() {
    logger, err := setupLogging()
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()
    
    manager := pip.NewManager(nil)
    manager.SetCustomLogger(logger)
    
    // All operations will now be logged
    // ...
}
```

## Testing Examples

### Unit Testing with Virtual Environments

```go
func TestWithCleanEnvironment(t *testing.T) {
    manager := pip.NewManager(nil)
    
    // Create temporary virtual environment
    tempDir, err := os.MkdirTemp("", "test-venv-*")
    require.NoError(t, err)
    defer os.RemoveAll(tempDir)
    
    venvPath := filepath.Join(tempDir, "venv")
    require.NoError(t, manager.CreateVenv(venvPath))
    require.NoError(t, manager.ActivateVenv(venvPath))
    
    // Test package installation in clean environment
    pkg := &pip.PackageSpec{Name: "requests"}
    require.NoError(t, manager.InstallPackage(pkg))
    
    // Verify installation
    packages, err := manager.ListPackages()
    require.NoError(t, err)
    
    found := false
    for _, p := range packages {
        if p.Name == "requests" {
            found = true
            break
        }
    }
    require.True(t, found, "requests package should be installed")
}
```

### Integration Testing

```go
func TestFullWorkflow(t *testing.T) {
    manager := pip.NewManager(nil)
    
    // Check pip installation
    installed, err := manager.IsInstalled()
    require.NoError(t, err)
    
    if !installed {
        require.NoError(t, manager.Install())
    }
    
    // Create project
    tempDir, err := os.MkdirTemp("", "test-project-*")
    require.NoError(t, err)
    defer os.RemoveAll(tempDir)
    
    opts := &pip.ProjectOptions{
        Name:         "test-project",
        Version:      "0.1.0",
        Author:       "Test Author",
        AuthorEmail:  "test@example.com",
        Dependencies: []string{"requests"},
        CreateVenv:   true,
    }
    
    require.NoError(t, manager.InitProject(tempDir, opts))
    
    // Verify project structure
    files := []string{
        "setup.py",
        "pyproject.toml",
        "requirements.txt",
        "README.md",
        ".gitignore",
    }
    
    for _, file := range files {
        path := filepath.Join(tempDir, file)
        _, err := os.Stat(path)
        require.NoError(t, err, "file %s should exist", file)
    }
}
```

## Performance Examples

### Concurrent Package Installation

```go
func installPackagesConcurrently(packages []*pip.PackageSpec) error {
    var wg sync.WaitGroup
    errChan := make(chan error, len(packages))
    
    for _, pkg := range packages {
        wg.Add(1)
        go func(p *pip.PackageSpec) {
            defer wg.Done()
            
            manager := pip.NewManager(nil)
            if err := manager.InstallPackage(p); err != nil {
                errChan <- fmt.Errorf("failed to install %s: %w", p.Name, err)
            }
        }(pkg)
    }
    
    wg.Wait()
    close(errChan)
    
    var errors []error
    for err := range errChan {
        errors = append(errors, err)
    }
    
    if len(errors) > 0 {
        return fmt.Errorf("installation errors: %v", errors)
    }
    
    return nil
}
```

## Next Steps

- Explore the [API Reference](/api/) for detailed documentation
- Check out the [Guide](/guide/) for comprehensive tutorials
- Browse the source code in the [GitHub repository](https://github.com/scagogogo/go-pip-sdk)
