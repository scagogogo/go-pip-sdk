# Virtual Environment Examples

This page provides practical examples of working with virtual environments using the Go Pip SDK.

## Basic Virtual Environment Workflow

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
    
    venvPath := "./example-env"
    
    // Create virtual environment
    fmt.Println("Creating virtual environment...")
    err := venvManager.CreateVenv(venvPath)
    if err != nil {
        log.Fatal(err)
    }
    
    // Activate virtual environment
    fmt.Println("Activating virtual environment...")
    err = venvManager.ActivateVenv(venvPath)
    if err != nil {
        log.Fatal(err)
    }
    
    // Install packages
    packages := []string{"requests", "flask", "pytest"}
    for _, pkgName := range packages {
        fmt.Printf("Installing %s...\n", pkgName)
        pkg := &pip.PackageSpec{Name: pkgName}
        err = manager.InstallPackage(pkg)
        if err != nil {
            log.Printf("Failed to install %s: %v", pkgName, err)
            continue
        }
    }
    
    // List installed packages
    fmt.Println("\nInstalled packages:")
    installedPkgs, err := manager.ListPackages()
    if err != nil {
        log.Fatal(err)
    }
    
    for _, pkg := range installedPkgs {
        fmt.Printf("  %s %s\n", pkg.Name, pkg.Version)
    }
    
    // Deactivate virtual environment
    fmt.Println("\nDeactivating virtual environment...")
    err = venvManager.DeactivateVenv()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Virtual environment workflow completed!")
}
```

## Multiple Virtual Environments

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func createProjectEnvironment(name string, packages []string) error {
    manager := pip.NewManager(nil)
    venvManager := pip.NewVenvManager(manager)
    
    venvPath := fmt.Sprintf("./%s-env", name)
    
    // Create and activate environment
    fmt.Printf("Setting up environment for %s...\n", name)
    err := venvManager.CreateVenv(venvPath)
    if err != nil {
        return err
    }
    
    err = venvManager.ActivateVenv(venvPath)
    if err != nil {
        return err
    }
    
    // Install packages
    for _, pkgName := range packages {
        pkg := &pip.PackageSpec{Name: pkgName}
        err = manager.InstallPackage(pkg)
        if err != nil {
            log.Printf("Warning: Failed to install %s: %v", pkgName, err)
        }
    }
    
    // Deactivate
    err = venvManager.DeactivateVenv()
    if err != nil {
        return err
    }
    
    fmt.Printf("Environment %s created successfully\n", name)
    return nil
}

func main() {
    // Create different environments for different projects
    projects := map[string][]string{
        "webapp": {"flask", "gunicorn", "psycopg2-binary"},
        "datascience": {"numpy", "pandas", "matplotlib", "jupyter"},
        "testing": {"pytest", "coverage", "mock"},
    }
    
    for project, packages := range projects {
        err := createProjectEnvironment(project, packages)
        if err != nil {
            log.Printf("Failed to create environment for %s: %v", project, err)
        }
    }
}
```

## Virtual Environment with Requirements

```go
package main

import (
    "fmt"
    "log"
    "os"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func setupProjectWithRequirements() error {
    manager := pip.NewManager(nil)
    venvManager := pip.NewVenvManager(manager)
    projectManager := pip.NewProjectManager(manager)
    
    venvPath := "./project-env"
    
    // Create virtual environment
    err := venvManager.CreateVenv(venvPath)
    if err != nil {
        return err
    }
    
    // Activate virtual environment
    err = venvManager.ActivateVenv(venvPath)
    if err != nil {
        return err
    }
    
    // Create requirements.txt
    requirements := `flask>=2.0.0
requests>=2.25.0
python-dotenv>=0.19.0
pytest>=6.0.0
black>=21.0.0`
    
    err = os.WriteFile("requirements.txt", []byte(requirements), 0644)
    if err != nil {
        return err
    }
    
    // Install from requirements
    fmt.Println("Installing from requirements.txt...")
    err = projectManager.InstallRequirements("requirements.txt")
    if err != nil {
        return err
    }
    
    // Generate frozen requirements
    fmt.Println("Generating frozen requirements...")
    err = projectManager.GenerateRequirements("requirements-frozen.txt")
    if err != nil {
        return err
    }
    
    fmt.Println("Project setup completed!")
    return nil
}

func main() {
    err := setupProjectWithRequirements()
    if err != nil {
        log.Fatal(err)
    }
}
```

## Virtual Environment Information

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func inspectVirtualEnvironment(venvPath string) {
    manager := pip.NewManager(nil)
    venvManager := pip.NewVenvManager(manager)
    
    // Get virtual environment info
    info, err := venvManager.GetVenvInfo(venvPath)
    if err != nil {
        log.Printf("Failed to get venv info: %v", err)
        return
    }
    
    fmt.Printf("Virtual Environment Information:\n")
    fmt.Printf("  Path: %s\n", info.Path)
    fmt.Printf("  Active: %t\n", info.IsActive)
    fmt.Printf("  Python Path: %s\n", info.PythonPath)
    fmt.Printf("  Pip Path: %s\n", info.PipPath)
    
    // Check if valid
    isValid := venvManager.IsVenvValid(venvPath)
    fmt.Printf("  Valid: %t\n", isValid)
    
    if isValid {
        // Activate and get package list
        err = venvManager.ActivateVenv(venvPath)
        if err != nil {
            log.Printf("Failed to activate: %v", err)
            return
        }
        
        packages, err := manager.ListPackages()
        if err != nil {
            log.Printf("Failed to list packages: %v", err)
            return
        }
        
        fmt.Printf("  Installed Packages: %d\n", len(packages))
        for _, pkg := range packages {
            fmt.Printf("    %s %s\n", pkg.Name, pkg.Version)
        }
        
        // Deactivate
        venvManager.DeactivateVenv()
    }
}

func main() {
    // Inspect multiple virtual environments
    venvPaths := []string{
        "./webapp-env",
        "./datascience-env",
        "./testing-env",
    }
    
    for _, venvPath := range venvPaths {
        fmt.Printf("\n" + "="*50 + "\n")
        inspectVirtualEnvironment(venvPath)
    }
}
```

## Virtual Environment Cleanup

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func findVirtualEnvironments(rootDir string) ([]string, error) {
    var venvs []string
    
    err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if info.IsDir() && (strings.HasSuffix(path, "-env") || 
                           strings.HasSuffix(path, "venv") ||
                           strings.Contains(path, "virtualenv")) {
            
            manager := pip.NewManager(nil)
            venvManager := pip.NewVenvManager(manager)
            
            if venvManager.IsVenvValid(path) {
                venvs = append(venvs, path)
            }
        }
        
        return nil
    })
    
    return venvs, err
}

func cleanupVirtualEnvironments(dryRun bool) error {
    manager := pip.NewManager(nil)
    venvManager := pip.NewVenvManager(manager)
    
    // Find all virtual environments
    venvs, err := findVirtualEnvironments(".")
    if err != nil {
        return err
    }
    
    fmt.Printf("Found %d virtual environments:\n", len(venvs))
    
    for _, venvPath := range venvs {
        info, err := venvManager.GetVenvInfo(venvPath)
        if err != nil {
            log.Printf("Failed to get info for %s: %v", venvPath, err)
            continue
        }
        
        fmt.Printf("\n%s:\n", venvPath)
        fmt.Printf("  Active: %t\n", info.IsActive)
        
        // Get size
        size, err := getDirSize(venvPath)
        if err == nil {
            fmt.Printf("  Size: %.2f MB\n", float64(size)/(1024*1024))
        }
        
        // Get package count
        if !info.IsActive {
            err = venvManager.ActivateVenv(venvPath)
            if err == nil {
                packages, err := manager.ListPackages()
                if err == nil {
                    fmt.Printf("  Packages: %d\n", len(packages))
                }
                venvManager.DeactivateVenv()
            }
        }
        
        // Ask for confirmation
        if !dryRun {
            fmt.Printf("  Remove this environment? (y/N): ")
            var response string
            fmt.Scanln(&response)
            
            if strings.ToLower(response) == "y" {
                err = venvManager.RemoveVenv(venvPath)
                if err != nil {
                    log.Printf("Failed to remove %s: %v", venvPath, err)
                } else {
                    fmt.Printf("  ✓ Removed\n")
                }
            }
        }
    }
    
    return nil
}

func getDirSize(path string) (int64, error) {
    var size int64
    err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            size += info.Size()
        }
        return err
    })
    return size, err
}

func main() {
    // Dry run first
    fmt.Println("=== DRY RUN ===")
    err := cleanupVirtualEnvironments(true)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("\n=== ACTUAL CLEANUP ===")
    err = cleanupVirtualEnvironments(false)
    if err != nil {
        log.Fatal(err)
    }
}
```

## CLI Examples

### Basic CLI Usage

```bash
# Create virtual environment
pip-cli venv create ./myproject-env

# Get environment information
pip-cli venv info ./myproject-env

# Activate environment
pip-cli venv activate ./myproject-env

# Install packages
pip-cli install flask requests pytest

# List packages in environment
pip-cli list

# Freeze requirements
pip-cli freeze > requirements.txt

# Deactivate environment
pip-cli venv deactivate

# Remove environment
pip-cli venv remove ./myproject-env
```

### Complete Project Setup

```bash
#!/bin/bash

PROJECT_NAME="mywebapp"
VENV_PATH="./${PROJECT_NAME}-env"

echo "Setting up project: $PROJECT_NAME"

# Create virtual environment
echo "Creating virtual environment..."
pip-cli venv create "$VENV_PATH"

# Activate virtual environment
echo "Activating virtual environment..."
pip-cli venv activate "$VENV_PATH"

# Install dependencies
echo "Installing dependencies..."
pip-cli install flask gunicorn python-dotenv

# Install development dependencies
echo "Installing development dependencies..."
pip-cli install pytest black flake8 mypy

# Generate requirements
echo "Generating requirements files..."
pip-cli freeze > requirements.txt

# Create project structure
mkdir -p src tests docs
touch src/__init__.py
touch tests/__init__.py
touch README.md

echo "Project setup completed!"
echo "Virtual environment: $VENV_PATH"
echo "To activate: pip-cli venv activate $VENV_PATH"
```

## Best Practices

1. **Use descriptive names** for virtual environments
2. **Always activate** before installing packages
3. **Generate requirements.txt** for reproducibility
4. **Clean up unused** virtual environments regularly
5. **Use separate environments** for different projects
6. **Document environment setup** in project README

## Next Steps

- Learn about [Project Initialization](./project-initialization.md)
- Explore [Advanced Usage](./advanced-usage.md)
- Check the [Virtual Environments API](../api/virtual-environments.md)
