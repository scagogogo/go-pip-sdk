# Package Management

Advanced examples for managing Python packages with the Go Pip SDK.

## Advanced Installation

### Installing with Custom Options

```go
func installWithOptions() error {
    manager := pip.NewManager(nil)
    
    pkg := &pip.PackageSpec{
        Name:    "numpy",
        Version: ">=1.20.0",
        Options: map[string]string{
            "no-cache-dir":     "",
            "timeout":          "300",
            "trusted-host":     "pypi.org",
            "extra-index-url":  "https://pypi.org/simple/",
        },
    }
    
    return manager.InstallPackage(pkg)
}
```

### Installing in Editable Mode

```go
func installEditable() error {
    manager := pip.NewManager(nil)
    
    pkg := &pip.PackageSpec{
        Name:     "/path/to/local/package",
        Editable: true,
    }
    
    return manager.InstallPackage(pkg)
}
```

### Force Reinstallation

```go
func forceReinstall() error {
    manager := pip.NewManager(nil)
    
    pkg := &pip.PackageSpec{
        Name:           "requests",
        ForceReinstall: true,
    }
    
    return manager.InstallPackage(pkg)
}
```

### Upgrade Packages

```go
func upgradePackage() error {
    manager := pip.NewManager(nil)
    
    pkg := &pip.PackageSpec{
        Name:    "requests",
        Upgrade: true,
    }
    
    return manager.InstallPackage(pkg)
}
```

## Batch Operations

### Installing Multiple Packages

```go
func installMultiplePackages() error {
    manager := pip.NewManager(nil)
    
    packages := []*pip.PackageSpec{
        {Name: "requests", Version: ">=2.25.0"},
        {Name: "click", Version: ">=7.0"},
        {Name: "pydantic", Version: ">=1.8.0"},
        {Name: "fastapi", Extras: []string{"dev"}},
    }
    
    var errors []error
    
    for _, pkg := range packages {
        fmt.Printf("Installing %s...\n", pkg.Name)
        if err := manager.InstallPackage(pkg); err != nil {
            errors = append(errors, fmt.Errorf("failed to install %s: %w", pkg.Name, err))
            continue
        }
        fmt.Printf("✓ %s installed successfully\n", pkg.Name)
    }
    
    if len(errors) > 0 {
        return fmt.Errorf("some installations failed: %v", errors)
    }
    
    return nil
}
```

### Concurrent Installation

```go
func installConcurrently() error {
    packages := []*pip.PackageSpec{
        {Name: "requests"},
        {Name: "click"},
        {Name: "pydantic"},
    }
    
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

## Dependency Management

### Installing with Dependency Resolution

```go
func installWithDependencies() error {
    manager := pip.NewManager(nil)
    
    // Install a package that has many dependencies
    pkg := &pip.PackageSpec{
        Name:    "django",
        Version: ">=4.0,<5.0",
        Options: map[string]string{
            "upgrade-strategy": "eager",
        },
    }
    
    if err := manager.InstallPackage(pkg); err != nil {
        return err
    }
    
    // Check what was installed
    packages, err := manager.ListPackages()
    if err != nil {
        return err
    }
    
    fmt.Printf("Installed %d packages (including dependencies)\n", len(packages))
    return nil
}
```

### Checking Package Dependencies

```go
func checkDependencies() error {
    manager := pip.NewManager(nil)
    
    info, err := manager.ShowPackage("django")
    if err != nil {
        return err
    }
    
    fmt.Printf("Package: %s\n", info.Name)
    fmt.Printf("Dependencies: %v\n", info.Requires)
    fmt.Printf("Required by: %v\n", info.RequiredBy)
    
    return nil
}
```

## Version Management

### Installing Specific Versions

```go
func installSpecificVersions() error {
    manager := pip.NewManager(nil)
    
    packages := []*pip.PackageSpec{
        {Name: "django", Version: "==4.2.0"},      // Exact version
        {Name: "requests", Version: ">=2.25.0"},   // Minimum version
        {Name: "click", Version: ">=7.0,<8.0"},    // Version range
        {Name: "pydantic", Version: "~=1.8.0"},    // Compatible release
    }
    
    for _, pkg := range packages {
        if err := manager.InstallPackage(pkg); err != nil {
            return fmt.Errorf("failed to install %s %s: %w", pkg.Name, pkg.Version, err)
        }
    }
    
    return nil
}
```

### Upgrading All Packages

```go
func upgradeAllPackages() error {
    manager := pip.NewManager(nil)
    
    // Get list of installed packages
    packages, err := manager.ListPackages()
    if err != nil {
        return err
    }
    
    // Upgrade each package
    for _, pkg := range packages {
        upgradeSpec := &pip.PackageSpec{
            Name:    pkg.Name,
            Upgrade: true,
        }
        
        fmt.Printf("Upgrading %s...\n", pkg.Name)
        if err := manager.InstallPackage(upgradeSpec); err != nil {
            fmt.Printf("Failed to upgrade %s: %v\n", pkg.Name, err)
            continue
        }
        fmt.Printf("✓ %s upgraded\n", pkg.Name)
    }
    
    return nil
}
```

## Custom Package Indexes

### Using Private Package Index

```go
func usePrivateIndex() error {
    config := &pip.Config{
        DefaultIndex: "https://pypi.company.com/simple/",
        TrustedHosts: []string{"pypi.company.com"},
    }
    
    manager := pip.NewManager(config)
    
    pkg := &pip.PackageSpec{
        Name:  "company-private-package",
        Index: "https://pypi.company.com/simple/",
    }
    
    return manager.InstallPackage(pkg)
}
```

### Multiple Package Indexes

```go
func useMultipleIndexes() error {
    manager := pip.NewManager(nil)
    
    packages := []*pip.PackageSpec{
        {
            Name:  "public-package",
            Index: "https://pypi.org/simple/",
        },
        {
            Name:  "private-package",
            Index: "https://pypi.company.com/simple/",
            Options: map[string]string{
                "trusted-host": "pypi.company.com",
            },
        },
    }
    
    for _, pkg := range packages {
        if err := manager.InstallPackage(pkg); err != nil {
            return err
        }
    }
    
    return nil
}
```

## Package Validation

### Verifying Installation

```go
func verifyInstallation() error {
    manager := pip.NewManager(nil)
    
    packageName := "requests"
    
    // Install package
    pkg := &pip.PackageSpec{Name: packageName}
    if err := manager.InstallPackage(pkg); err != nil {
        return err
    }
    
    // Verify installation
    packages, err := manager.ListPackages()
    if err != nil {
        return err
    }
    
    found := false
    for _, p := range packages {
        if p.Name == packageName {
            found = true
            fmt.Printf("✓ %s %s installed successfully\n", p.Name, p.Version)
            break
        }
    }
    
    if !found {
        return fmt.Errorf("package %s not found after installation", packageName)
    }
    
    return nil
}
```

### Checking Package Integrity

```go
func checkPackageIntegrity() error {
    manager := pip.NewManager(nil)
    
    info, err := manager.ShowPackage("requests")
    if err != nil {
        return err
    }
    
    fmt.Printf("Package: %s\n", info.Name)
    fmt.Printf("Version: %s\n", info.Version)
    fmt.Printf("Location: %s\n", info.Location)
    fmt.Printf("Files: %d\n", len(info.Files))
    
    // Check if package location exists
    if _, err := os.Stat(info.Location); err != nil {
        return fmt.Errorf("package location not found: %w", err)
    }
    
    fmt.Println("✓ Package integrity check passed")
    return nil
}
```

## Error Recovery

### Retry Failed Installations

```go
func retryInstallation() error {
    manager := pip.NewManager(nil)
    
    pkg := &pip.PackageSpec{Name: "some-package"}
    maxRetries := 3
    
    for i := 0; i < maxRetries; i++ {
        err := manager.InstallPackage(pkg)
        if err == nil {
            fmt.Printf("✓ Package installed on attempt %d\n", i+1)
            return nil
        }
        
        if pip.IsErrorType(err, pip.ErrorTypeNetworkError) {
            fmt.Printf("Network error on attempt %d, retrying...\n", i+1)
            time.Sleep(time.Duration(i+1) * time.Second)
            continue
        }
        
        // Non-recoverable error
        return err
    }
    
    return fmt.Errorf("failed after %d attempts", maxRetries)
}
```

### Rollback on Failure

```go
func installWithRollback() error {
    manager := pip.NewManager(nil)
    
    // Get current package list
    beforeInstall, err := manager.ListPackages()
    if err != nil {
        return err
    }
    
    // Try to install new packages
    packages := []*pip.PackageSpec{
        {Name: "package1"},
        {Name: "package2"},
        {Name: "package3"},
    }
    
    var installedPackages []string
    
    for _, pkg := range packages {
        if err := manager.InstallPackage(pkg); err != nil {
            fmt.Printf("Installation failed at %s, rolling back...\n", pkg.Name)
            
            // Rollback: uninstall previously installed packages
            for _, installed := range installedPackages {
                if uninstallErr := manager.UninstallPackage(installed); uninstallErr != nil {
                    fmt.Printf("Warning: failed to uninstall %s during rollback: %v\n", installed, uninstallErr)
                }
            }
            
            return fmt.Errorf("installation failed and rolled back: %w", err)
        }
        
        installedPackages = append(installedPackages, pkg.Name)
    }
    
    fmt.Println("✓ All packages installed successfully")
    return nil
}
```

## Performance Optimization

### Parallel Installation with Worker Pool

```go
func parallelInstallation() error {
    packages := []*pip.PackageSpec{
        {Name: "requests"},
        {Name: "click"},
        {Name: "pydantic"},
        {Name: "fastapi"},
        {Name: "uvicorn"},
    }
    
    const numWorkers = 3
    packageChan := make(chan *pip.PackageSpec, len(packages))
    resultChan := make(chan error, len(packages))
    
    // Start workers
    for i := 0; i < numWorkers; i++ {
        go func() {
            manager := pip.NewManager(nil)
            for pkg := range packageChan {
                err := manager.InstallPackage(pkg)
                if err != nil {
                    resultChan <- fmt.Errorf("failed to install %s: %w", pkg.Name, err)
                } else {
                    resultChan <- nil
                }
            }
        }()
    }
    
    // Send packages to workers
    for _, pkg := range packages {
        packageChan <- pkg
    }
    close(packageChan)
    
    // Collect results
    var errors []error
    for i := 0; i < len(packages); i++ {
        if err := <-resultChan; err != nil {
            errors = append(errors, err)
        }
    }
    
    if len(errors) > 0 {
        return fmt.Errorf("installation errors: %v", errors)
    }
    
    return nil
}
```

## Next Steps

- [Virtual Environments](/examples/virtual-environments) - Managing Python virtual environments
- [Project Initialization](/examples/project-initialization) - Setting up Python projects
- [Advanced Usage](/examples/advanced-usage) - Complex scenarios and configurations
