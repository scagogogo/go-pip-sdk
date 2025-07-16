# Types

This section documents all the core data types and interfaces used throughout the Go Pip SDK.

## Interfaces

### PipManager

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

The main interface that defines all pip operations. Implemented by the `Manager` struct.

## Configuration Types

### Config

```go
type Config struct {
    PythonPath   string            // Path to Python executable
    PipPath      string            // Path to pip executable
    DefaultIndex string            // Default package index URL
    TrustedHosts []string          // Trusted hosts for package downloads
    Timeout      time.Duration     // Timeout for operations
    Retries      int               // Number of retries for failed operations
    LogLevel     string            // Logging level (DEBUG, INFO, WARN, ERROR)
    CacheDir     string            // Cache directory for pip
    ExtraOptions map[string]string // Additional pip options
    Environment  map[string]string // Environment variables
}
```

Configuration for the pip manager.

**Example:**
```go
config := &pip.Config{
    PythonPath:   "/usr/bin/python3",
    PipPath:      "/usr/bin/pip3",
    DefaultIndex: "https://pypi.org/simple/",
    TrustedHosts: []string{"pypi.org", "pypi.python.org"},
    Timeout:      60 * time.Second,
    Retries:      3,
    LogLevel:     "INFO",
    CacheDir:     "/tmp/pip-cache",
    ExtraOptions: map[string]string{
        "no-warn-script-location": "",
    },
    Environment: map[string]string{
        "PIP_DISABLE_PIP_VERSION_CHECK": "1",
    },
}
```

## Package Types

### PackageSpec

```go
type PackageSpec struct {
    Name           string            // Package name (required)
    Version        string            // Version constraint (e.g., ">=1.0.0")
    Extras         []string          // Extra dependencies
    Index          string            // Custom index URL
    Options        map[string]string // Additional pip options
    Editable       bool              // Install in editable mode
    Upgrade        bool              // Upgrade if already installed
    ForceReinstall bool              // Force reinstallation
}
```

Specification for package installation.

**Version Constraint Examples:**
- `">=1.0.0"` - At least version 1.0.0
- `"==2.1.0"` - Exactly version 2.1.0
- `">=1.0,<2.0"` - Version 1.x
- `"~=1.4.2"` - Compatible release (>=1.4.2, ==1.4.*)

**Example:**
```go
pkg := &pip.PackageSpec{
    Name:    "django",
    Version: ">=4.0,<5.0",
    Extras:  []string{"postgres", "redis"},
    Options: map[string]string{
        "no-cache-dir": "",
        "timeout":      "120",
    },
    Upgrade: true,
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

Represents an installed package.

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

Detailed information about a package.

### SearchResult

```go
type SearchResult struct {
    Name        string // Package name
    Version     string // Latest version
    Summary     string // Short description
    Score       float64 // Search relevance score
}
```

Result from package search (note: pip search is disabled in newer versions).

## Project Types

### ProjectOptions

```go
type ProjectOptions struct {
    Name            string   // Project name
    Version         string   // Project version
    Description     string   // Project description
    Author          string   // Author name
    AuthorEmail     string   // Author email
    License         string   // License
    PythonVersion   string   // Required Python version
    Dependencies    []string // Runtime dependencies
    DevDependencies []string // Development dependencies
    Template        string   // Project template
    CreateVenv      bool     // Create virtual environment
    VenvPath        string   // Virtual environment path
}
```

Options for project initialization.

**Example:**
```go
opts := &pip.ProjectOptions{
    Name:            "my-project",
    Version:         "0.1.0",
    Description:     "A sample Python project",
    Author:          "John Doe",
    AuthorEmail:     "john@example.com",
    License:         "MIT",
    PythonVersion:   ">=3.8",
    Dependencies:    []string{"requests>=2.25.0", "click>=7.0"},
    DevDependencies: []string{"pytest>=6.0", "black>=21.0"},
    CreateVenv:      true,
    VenvPath:        "./venv",
}
```

## Virtual Environment Types

### VenvInfo

```go
type VenvInfo struct {
    Path          string    // Virtual environment path
    IsActive      bool      // Whether currently active
    PythonPath    string    // Path to Python executable
    PythonVersion string    // Python version in the venv
    CreatedAt     time.Time // Creation timestamp
}
```

Information about a virtual environment.

### VenvManager

```go
type VenvManager struct {
    manager *Manager // Reference to the main manager
}
```

Internal helper for virtual environment operations.

## System Types

### OSType

```go
type OSType string

const (
    OSWindows OSType = "windows"
    OSMacOS   OSType = "darwin"
    OSLinux   OSType = "linux"
    OSUnknown OSType = "unknown"
)
```

Operating system type enumeration.

**Methods:**
```go
func (o OSType) String() string
```

Returns the string representation of the OS type.

## Installer Types

### Installer

```go
type Installer struct {
    manager *Manager // Reference to the main manager
}
```

Handles pip installation across different operating systems.

## Type Validation

### Package Name Validation

Package names must follow Python package naming conventions:
- Contain only letters, numbers, hyphens, underscores, and periods
- Start with a letter or number
- Be case-insensitive

```go
func isValidPackageName(name string) bool {
    if name == "" {
        return false
    }
    // Implementation would check against Python package naming rules
    return true
}
```

### Version Constraint Validation

Version constraints must follow PEP 440 specification:

```go
// Valid version constraints
validConstraints := []string{
    ">=1.0.0",
    "==2.1.0",
    ">=1.0,<2.0",
    "~=1.4.2",
    "!=1.5.0",
}
```

## Type Conversion Utilities

### String to OSType

```go
func ParseOSType(s string) OSType {
    switch strings.ToLower(s) {
    case "windows":
        return OSWindows
    case "darwin", "macos":
        return OSMacOS
    case "linux":
        return OSLinux
    default:
        return OSUnknown
    }
}
```

### Package to String

```go
func (p *Package) String() string {
    return fmt.Sprintf("%s==%s", p.Name, p.Version)
}
```

### PackageSpec to Command Arguments

```go
func (ps *PackageSpec) ToArgs() []string {
    var args []string

    // Build package specification
    spec := ps.Name
    if ps.Version != "" {
        spec += ps.Version
    }
    if len(ps.Extras) > 0 {
        spec += "[" + strings.Join(ps.Extras, ",") + "]"
    }
    args = append(args, spec)

    // Add options
    if ps.Upgrade {
        args = append(args, "--upgrade")
    }
    if ps.ForceReinstall {
        args = append(args, "--force-reinstall")
    }
    if ps.Editable {
        args = append(args, "--editable")
    }

    return args
}
```

## JSON Serialization

All major types support JSON serialization:

```go
// Serialize package info to JSON
info := &PackageInfo{
    Name:    "requests",
    Version: "2.28.1",
    Summary: "Python HTTP for Humans.",
}

data, err := json.Marshal(info)
if err != nil {
    return err
}

// Deserialize from JSON
var info2 PackageInfo
if err := json.Unmarshal(data, &info2); err != nil {
    return err
}
```

## Type Safety Best Practices

1. **Always validate input parameters**:
   ```go
   func (m *Manager) InstallPackage(pkg *PackageSpec) error {
       if pkg == nil {
           return errors.New("package specification cannot be nil")
       }
       if pkg.Name == "" {
           return errors.New("package name cannot be empty")
       }
       // ... rest of validation
   }
   ```

2. **Use type assertions safely**:
   ```go
   if pipErr, ok := err.(*PipErrorDetails); ok {
       // Handle pip-specific error
       fmt.Printf("Pip error: %s\n", pipErr.Type)
   }
   ```

3. **Provide sensible defaults**:
   ```go
   func NewPackageSpec(name string) *PackageSpec {
       return &PackageSpec{
           Name:           name,
           Options:        make(map[string]string),
           Upgrade:        false,
           ForceReinstall: false,
           Editable:       false,
       }
   }
   ```

## Compatibility

All types are designed to be:
- **Backward compatible**: New fields are added with `omitempty` JSON tags
- **Forward compatible**: Unknown fields are ignored during deserialization
- **Cross-platform**: No platform-specific fields in core types
