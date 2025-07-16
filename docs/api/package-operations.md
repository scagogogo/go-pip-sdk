# Package Operations

The package operations API provides comprehensive functionality for managing Python packages, including installation, uninstallation, listing, and information retrieval.

## Installation

### InstallPackage

```go
func (m *Manager) InstallPackage(pkg *PackageSpec) error
```

Installs a Python package with the specified configuration.

**Parameters:**
- `pkg` (*PackageSpec): Package specification with installation options.

**Returns:**
- `error`: Error if installation fails.

**Example:**
```go
// Basic installation
pkg := &pip.PackageSpec{
    Name:    "requests",
    Version: ">=2.25.0",
}
if err := manager.InstallPackage(pkg); err != nil {
    return err
}

// Advanced installation with options
pkg := &pip.PackageSpec{
    Name:           "django",
    Version:        ">=4.0,<5.0",
    Extras:         []string{"postgres", "redis"},
    Upgrade:        true,
    ForceReinstall: false,
    Index:          "https://pypi.org/simple/",
    Options: map[string]string{
        "no-cache-dir": "",
        "timeout":      "60",
    },
}
if err := manager.InstallPackage(pkg); err != nil {
    return err
}
```

### InstallRequirements

```go
func (m *Manager) InstallRequirements(path string) error
```

Installs packages from a requirements.txt file.

**Parameters:**
- `path` (string): Path to the requirements.txt file.

**Returns:**
- `error`: Error if installation fails.

**Example:**
```go
if err := manager.InstallRequirements("requirements.txt"); err != nil {
    return fmt.Errorf("failed to install requirements: %w", err)
}
```

## Uninstallation

### UninstallPackage

```go
func (m *Manager) UninstallPackage(name string) error
```

Uninstalls a Python package.

**Parameters:**
- `name` (string): Name of the package to uninstall.

**Returns:**
- `error`: Error if uninstallation fails.

**Example:**
```go
if err := manager.UninstallPackage("requests"); err != nil {
    return fmt.Errorf("failed to uninstall package: %w", err)
}
```

## Package Information

### ListPackages

```go
func (m *Manager) ListPackages() ([]*Package, error)
```

Lists all installed packages.

**Returns:**
- `[]*Package`: Slice of installed packages.
- `error`: Error if listing fails.

**Example:**
```go
packages, err := manager.ListPackages()
if err != nil {
    return fmt.Errorf("failed to list packages: %w", err)
}

for _, pkg := range packages {
    fmt.Printf("%s %s\n", pkg.Name, pkg.Version)
}
```

### ShowPackage

```go
func (m *Manager) ShowPackage(name string) (*PackageInfo, error)
```

Shows detailed information about a specific package.

**Parameters:**
- `name` (string): Name of the package to show.

**Returns:**
- `*PackageInfo`: Detailed package information.
- `error`: Error if package information retrieval fails.

**Example:**
```go
info, err := manager.ShowPackage("requests")
if err != nil {
    return fmt.Errorf("failed to get package info: %w", err)
}

fmt.Printf("Name: %s\n", info.Name)
fmt.Printf("Version: %s\n", info.Version)
fmt.Printf("Summary: %s\n", info.Summary)
fmt.Printf("Author: %s\n", info.Author)
fmt.Printf("License: %s\n", info.License)
fmt.Printf("Dependencies: %v\n", info.Requires)
```

### FreezePackages

```go
func (m *Manager) FreezePackages() ([]*Package, error)
```

Returns a list of installed packages with exact versions (equivalent to `pip freeze`).

**Returns:**
- `[]*Package`: Slice of packages with exact versions.
- `error`: Error if freeze operation fails.

**Example:**
```go
packages, err := manager.FreezePackages()
if err != nil {
    return fmt.Errorf("failed to freeze packages: %w", err)
}

for _, pkg := range packages {
    fmt.Printf("%s==%s\n", pkg.Name, pkg.Version)
}
```

### SearchPackages

```go
func (m *Manager) SearchPackages(query string) ([]*SearchResult, error)
```

Searches for packages on PyPI. **Note**: This feature is disabled in newer pip versions.

**Parameters:**
- `query` (string): Search query.

**Returns:**
- `[]*SearchResult`: Search results (currently returns error due to pip limitation).
- `error`: Error indicating that search is disabled.

**Example:**
```go
results, err := manager.SearchPackages("http")
if err != nil {
    // Expected error: pip search is disabled
    fmt.Printf("Search disabled: %v\n", err)
    return err
}
```

## Requirements Management

### GenerateRequirements

```go
func (m *Manager) GenerateRequirements(path string) error
```

Generates a requirements.txt file from currently installed packages.

**Parameters:**
- `path` (string): Path where to save the requirements.txt file.

**Returns:**
- `error`: Error if generation fails.

**Example:**
```go
if err := manager.GenerateRequirements("requirements.txt"); err != nil {
    return fmt.Errorf("failed to generate requirements: %w", err)
}
```

## Data Types

### PackageSpec

```go
type PackageSpec struct {
    Name           string            // Package name (required)
    Version        string            // Version constraint (e.g., ">=1.0.0", "==2.1.0")
    Extras         []string          // Extra dependencies (e.g., ["dev", "test"])
    Index          string            // Custom index URL
    Options        map[string]string // Additional pip options
    Editable       bool              // Install in editable mode
    Upgrade        bool              // Upgrade if already installed
    ForceReinstall bool              // Force reinstallation
}
```

### Package

```go
type Package struct {
    Name      string // Package name
    Version   string // Installed version
    Location  string // Installation location
    Editable  bool   // Whether installed in editable mode
    Installer string // Installer used (pip, conda, etc.)
}
```

### PackageInfo

```go
type PackageInfo struct {
    Name        string            // Package name
    Version     string            // Version
    Summary     string            // Short description
    HomePage    string            // Homepage URL
    Author      string            // Author name
    AuthorEmail string            // Author email
    License     string            // License
    Location    string            // Installation location
    Requires    []string          // Dependencies
    RequiredBy  []string          // Packages that depend on this
    Files       []string          // Installed files
    Metadata    map[string]string // Additional metadata
}
```

## Common Usage Patterns

### Installing Multiple Packages

```go
packages := []*pip.PackageSpec{
    {Name: "requests", Version: ">=2.25.0"},
    {Name: "click", Version: ">=7.0"},
    {Name: "pydantic", Version: ">=1.8.0"},
}

for _, pkg := range packages {
    if err := manager.InstallPackage(pkg); err != nil {
        return fmt.Errorf("failed to install %s: %w", pkg.Name, err)
    }
}
```

### Installing with Development Dependencies

```go
// Install package with development extras
pkg := &pip.PackageSpec{
    Name:   "fastapi",
    Extras: []string{"dev", "test"},
}
if err := manager.InstallPackage(pkg); err != nil {
    return err
}
```

### Installing from Custom Index

```go
pkg := &pip.PackageSpec{
    Name:  "private-package",
    Index: "https://private-pypi.company.com/simple/",
    Options: map[string]string{
        "trusted-host": "private-pypi.company.com",
    },
}
if err := manager.InstallPackage(pkg); err != nil {
    return err
}
```

## Error Handling

Package operations can fail for various reasons. The SDK provides specific error types:

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

## Best Practices

1. **Always validate package specifications**:
   ```go
   if pkg.Name == "" {
       return errors.New("package name is required")
   }
   ```

2. **Use version constraints for reproducible builds**:
   ```go
   pkg := &pip.PackageSpec{
       Name:    "django",
       Version: ">=4.0,<5.0",  // Pin major version
   }
   ```

3. **Handle network timeouts gracefully**:
   ```go
   config := &pip.Config{
       Timeout: 120 * time.Second,
       Retries: 3,
   }
   manager := pip.NewManager(config)
   ```
