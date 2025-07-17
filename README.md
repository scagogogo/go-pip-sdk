# Go Pip SDK

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.19-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/go-pip-sdk)](https://goreportcard.com/report/github.com/scagogogo/go-pip-sdk)
[![Documentation](https://img.shields.io/badge/docs-online-blue.svg)](https://scagogogo.github.io/go-pip-sdk/)
[![Build Status](https://img.shields.io/github/actions/workflow/status/scagogogo/go-pip-sdk/ci.yml?branch=main)](https://github.com/scagogogo/go-pip-sdk/actions)
[![Coverage Status](https://img.shields.io/codecov/c/github/scagogogo/go-pip-sdk)](https://codecov.io/gh/scagogogo/go-pip-sdk)

A comprehensive, production-ready Go SDK for managing Python pip operations, virtual environments, and Python projects. This library provides a clean, type-safe interface for all common pip operations with enterprise-grade features and cross-platform support.

**English** | [简体中文](README.zh-CN.md)

## ✨ Features

### 🚀 Core Capabilities
- **Cross-platform support** - Works seamlessly on Windows, macOS, and Linux
- **Complete pip operations** - Install, uninstall, list, show, freeze, search packages
- **Virtual environment management** - Create, activate, deactivate, remove, clone virtual environments
- **Project initialization** - Bootstrap Python projects with customizable templates
- **Automatic pip installation** - Detects and installs pip if missing with multiple installation methods

### 🏢 Enterprise Features
- **Production-ready** - Battle-tested in enterprise environments
- **Comprehensive logging** - Structured logging with multiple output formats (JSON, text)
- **Advanced error handling** - Rich error types with actionable suggestions and retry mechanisms
- **Configuration management** - Flexible configuration with environment variable support
- **Security features** - Certificate validation, trusted hosts, and secure package installation

### 🛠️ Developer Experience
- **Type-safe API** - Full Go type safety with comprehensive interfaces
- **Extensive testing** - 95%+ test coverage with unit and integration tests
- **Rich documentation** - Complete API documentation with examples
- **Command-line interface** - Full-featured CLI tool for direct usage
- **Docker support** - Official Docker images and containerized deployment options

## 📦 Installation

### Using Go Modules (Recommended)

```bash
go get github.com/scagogogo/go-pip-sdk
```

### Using Specific Version

```bash
go get github.com/scagogogo/go-pip-sdk@v1.0.0
```

### Requirements

- **Go**: 1.19 or later
- **Python**: 3.7 or later (for pip operations)
- **Operating System**: Windows 10+, macOS 10.15+, or Linux (any modern distribution)

## 🚀 Quick Start

### Basic Usage

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

    // Check if pip is installed and install if missing
    if installed, err := manager.IsInstalled(); err != nil {
        log.Fatal(err)
    } else if !installed {
        fmt.Println("Installing pip...")
        if err := manager.Install(); err != nil {
            log.Fatal(err)
        }
        fmt.Println("✅ Pip installed successfully!")
    }

    // Install a package with version constraints
    pkg := &pip.PackageSpec{
        Name:    "requests",
        Version: ">=2.25.0,<3.0.0",
        Extras:  []string{"security"}, // Install with extras
    }

    fmt.Printf("Installing %s...\n", pkg.Name)
    if err := manager.InstallPackage(pkg); err != nil {
        log.Fatal(err)
    }

    fmt.Println("✅ Package installed successfully!")

    // List installed packages
    packages, err := manager.ListPackages()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d installed packages\n", len(packages))
}
```

### With Custom Configuration

```go
package main

import (
    "time"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    // Create custom configuration
    config := &pip.Config{
        PythonPath:   "/usr/bin/python3",
        Timeout:      120 * time.Second,
        Retries:      5,
        LogLevel:     "INFO",
        DefaultIndex: "https://pypi.org/simple/",
        TrustedHosts: []string{"pypi.org", "pypi.python.org"},
        Environment: map[string]string{
            "PIP_CACHE_DIR": "/tmp/pip-cache",
        },
    }

    manager := pip.NewManager(config)

    // Your pip operations here...
}
```

## 📚 Core Operations

### Package Management

```go
// Install packages with various options
pkg := &pip.PackageSpec{
    Name:           "fastapi",
    Version:        ">=0.68.0,<1.0.0",
    Extras:         []string{"all"},
    Upgrade:        true,
    ForceReinstall: false,
    UserInstall:    false,
}
err := manager.InstallPackage(pkg)

// Install from requirements file
err = manager.InstallRequirements("requirements.txt")

// Install from Git repository
gitPkg := &pip.PackageSpec{
    Name: "git+https://github.com/user/repo.git@v1.0.0",
}
err = manager.InstallPackage(gitPkg)

// Uninstall packages
err = manager.UninstallPackage("requests")

// List installed packages with details
packages, err := manager.ListPackages()
for _, pkg := range packages {
    fmt.Printf("%s==%s (%s)\n", pkg.Name, pkg.Version, pkg.Location)
}

// Show detailed package information
info, err := manager.ShowPackage("requests")
fmt.Printf("Name: %s\nVersion: %s\nSummary: %s\n",
    info.Name, info.Version, info.Summary)

// Search for packages
results, err := manager.SearchPackages("web framework")

// Check for outdated packages
outdated, err := manager.CheckOutdated()

// Freeze packages (like pip freeze)
packages, err := manager.FreezePackages()
```

### Virtual Environment Management

```go
// Create virtual environment with options
opts := &pip.VenvOptions{
    PythonVersion:      "3.9",
    SystemSitePackages: false,
    Prompt:             "my-project",
    UpgradePip:         true,
}
err := manager.CreateVenvWithOptions("/path/to/venv", opts)

// Activate virtual environment
err = manager.ActivateVenv("/path/to/venv")

// Check if virtual environment is active
isActive, venvPath := manager.IsVenvActive()
if isActive {
    fmt.Printf("Active virtual environment: %s\n", venvPath)
}

// List all virtual environments
venvs, err := manager.ListVenvs("/path/to/envs")

// Get detailed virtual environment information
info, err := manager.GetVenvInfo("/path/to/venv")
fmt.Printf("Python version: %s\nPackages: %d\n",
    info.PythonVersion, info.PackageCount)

// Clone virtual environment
err = manager.CloneVenv("/path/to/source", "/path/to/target")

// Remove virtual environment
err = manager.RemoveVenv("/path/to/venv")
```

### Project Initialization

```go
// Initialize a comprehensive Python project
opts := &pip.ProjectOptions{
    Name:            "my-awesome-project",
    Version:         "0.1.0",
    Description:     "A comprehensive Python project",
    Author:          "Your Name",
    AuthorEmail:     "your.email@example.com",
    License:         "MIT",
    Homepage:        "https://github.com/user/my-awesome-project",
    Repository:      "https://github.com/user/my-awesome-project.git",

    // Dependencies
    Dependencies: []string{
        "fastapi>=0.68.0",
        "uvicorn[standard]>=0.15.0",
        "pydantic>=1.8.0",
    },
    DevDependencies: []string{
        "pytest>=6.0.0",
        "black>=21.0.0",
        "flake8>=3.8.0",
        "mypy>=0.812",
    },

    // Project structure
    CreateVenv:          true,
    CreateSrc:           true,
    CreateTests:         true,
    CreateDocs:          true,
    CreateGithubActions: true,
    CreateDockerfile:    true,

    // Configuration files
    CreateSetupPy:       true,
    CreatePyprojectToml: true,
    CreateGitignore:     true,
    CreateReadme:        true,
}

err := manager.InitProject("/path/to/project", opts)

// Read project configuration
config, err := manager.ReadProjectConfig("/path/to/project")

// Update project version
err = manager.UpdateProjectVersion("/path/to/project", "1.0.0")

// Build project
buildOpts := &pip.BuildOptions{
    OutputDir: "./dist",
    Format:    "wheel",
    Clean:     true,
}
err = manager.BuildProject("/path/to/project", buildOpts)
```

## ⚙️ Configuration

### Basic Configuration

```go
config := &pip.Config{
    // Python settings
    PythonPath: "/usr/bin/python3",
    PipPath:    "/usr/bin/pip3",

    // Network settings
    Timeout:      120 * time.Second,
    Retries:      5,
    DefaultIndex: "https://pypi.org/simple/",
    ExtraIndexes: []string{
        "https://pypi.python.org/simple/",
        "https://test.pypi.org/simple/",
    },
    TrustedHosts: []string{"pypi.org", "pypi.python.org"},

    // Cache settings
    CacheDir: "/tmp/pip-cache",
    NoCache:  false,

    // Logging
    LogLevel: "INFO",
    LogFile:  "/var/log/pip-sdk.log",

    // Environment variables
    Environment: map[string]string{
        "PIP_CACHE_DIR":              "/tmp/pip-cache",
        "PIP_DISABLE_PIP_VERSION_CHECK": "1",
        "PIP_TIMEOUT":                "120",
    },
}

manager := pip.NewManager(config)
```

### Enterprise Configuration

```go
// Enterprise-grade configuration with security features
config := &pip.Config{
    PythonPath:   "/opt/python/bin/python3",
    DefaultIndex: "https://pypi.company.com/simple/",
    ExtraIndexes: []string{
        "https://pypi.org/simple/",
    },
    TrustedHosts: []string{
        "pypi.company.com",
        "pypi.org",
    },
    Timeout: 300 * time.Second,
    Retries: 10,

    // Security settings
    ExtraOptions: map[string]string{
        "cert":         "/etc/ssl/certs/company-ca.pem",
        "client-cert":  "/etc/ssl/certs/client.pem",
        "trusted-host": "pypi.company.com",
    },

    // Logging for audit trails
    LogLevel: "INFO",
    LogFile:  "/var/log/pip-operations.log",
}
```

## 🚀 Advanced Usage

### Custom Logging

```go
// Create structured logger with multiple outputs
loggerConfig := &pip.LoggerConfig{
    Level:      pip.LogLevelDebug,
    Format:     pip.LogFormatJSON,
    EnableFile: true,
    LogFile:    "/var/log/pip-sdk.log",
    MaxSize:    100, // 100MB
    MaxBackups: 5,
    MaxAge:     30, // 30 days
    Compress:   true,

    // Add custom fields to all log entries
    Fields: map[string]interface{}{
        "service":     "pip-manager",
        "version":     "1.0.0",
        "environment": os.Getenv("ENVIRONMENT"),
    },
}

logger, err := pip.NewLogger(loggerConfig)
if err != nil {
    log.Fatal(err)
}
defer logger.Close()

// Set custom logger
manager.SetCustomLogger(logger)

// Use contextual logging
contextLogger := logger.WithFields(map[string]interface{}{
    "operation": "package_install",
    "user_id":   "12345",
})
```

### Advanced Error Handling

```go
// Comprehensive error handling with retry logic
func installWithRetry(manager *pip.Manager, pkg *pip.PackageSpec, maxRetries int) error {
    var lastErr error

    for attempt := 0; attempt < maxRetries; attempt++ {
        err := manager.InstallPackage(pkg)
        if err == nil {
            return nil // Success
        }

        lastErr = err

        // Handle different error types
        switch pip.GetErrorType(err) {
        case pip.ErrorTypeNetworkError, pip.ErrorTypeNetworkTimeout:
            // Retry network errors with exponential backoff
            delay := time.Duration(1<<uint(attempt)) * time.Second
            fmt.Printf("Network error, retrying in %v... (%d/%d)\n", delay, attempt+1, maxRetries)
            time.Sleep(delay)
            continue

        case pip.ErrorTypePermissionDenied:
            // Try user installation for permission errors
            if attempt == 0 {
                pkg.UserInstall = true
                continue
            }
            return err

        case pip.ErrorTypePackageNotFound:
            // Suggest alternatives for missing packages
            if results, searchErr := manager.SearchPackages(pkg.Name); searchErr == nil && len(results) > 0 {
                fmt.Printf("Package '%s' not found. Similar packages:\n", pkg.Name)
                for i, result := range results[:3] {
                    fmt.Printf("%d. %s - %s\n", i+1, result.Name, result.Summary)
                }
            }
            return err

        default:
            return err // Don't retry other errors
        }
    }

    return fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}
```

### Context Support and Cancellation

```go
// Use context for timeout and cancellation
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
defer cancel()

manager := pip.NewManagerWithContext(ctx, config)

// Operations will respect context cancellation
err := manager.InstallPackage(pkg)
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        fmt.Println("Operation timed out")
    } else if errors.Is(err, context.Canceled) {
        fmt.Println("Operation was canceled")
    }
}
```

### Performance Monitoring

```go
// Monitor operation performance
func monitoredInstall(manager *pip.Manager, pkg *pip.PackageSpec) error {
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        fmt.Printf("Package installation took %v\n", duration)
    }()

    return manager.InstallPackage(pkg)
}
```

## Command-Line Interface

The SDK includes a full-featured CLI tool for direct usage:

### Installation

```bash
# Using Go install
go install github.com/scagogogo/go-pip-sdk/cmd/pip-cli@latest

# Or build from source
git clone https://github.com/scagogogo/go-pip-sdk.git
cd go-pip-sdk
make build
```

### Usage

```bash
# Install packages
pip-cli install requests
pip-cli install "django>=4.0,<5.0"

# Manage virtual environments
pip-cli venv create ./myenv
pip-cli venv activate ./myenv
pip-cli venv deactivate

# Initialize projects
pip-cli project init ./myproject

# Get help
pip-cli help
```

### Docker Usage

```bash
# Run with Docker
docker run --rm scagogogo/pip-cli:latest help

# Use Docker Compose for development
cd cmd/pip-cli
docker-compose up pip-cli-dev
```

## 📖 Documentation

- 📖 **[Online Documentation](https://scagogogo.github.io/go-pip-sdk/)** - Complete API documentation and guides
- 🚀 **[Getting Started](https://scagogogo.github.io/go-pip-sdk/guide/getting-started)** - Quick start guide and installation
- 📚 **[API Reference](https://scagogogo.github.io/go-pip-sdk/api/)** - Detailed API documentation with examples
- 💡 **[Examples](https://scagogogo.github.io/go-pip-sdk/examples/)** - Comprehensive code examples and use cases
- 🔧 **[Configuration](https://scagogogo.github.io/go-pip-sdk/guide/configuration)** - Configuration options and best practices
- 🐛 **[Troubleshooting](https://scagogogo.github.io/go-pip-sdk/guide/troubleshooting)** - Common issues and solutions

## Examples

See the [examples](examples/) directory for more comprehensive examples:

- [Basic Usage](examples/basic/main.go)
- [Virtual Environment Management](examples/venv/main.go)
- [Project Initialization](examples/project/main.go)
- [Advanced Configuration](examples/advanced/main.go)

## 🔍 API Reference

### Core Interfaces

```go
// Main pip manager interface
type Manager interface {
    // Package operations
    InstallPackage(pkg *PackageSpec) error
    UninstallPackage(name string) error
    ListPackages() ([]*Package, error)
    ShowPackage(name string) (*PackageInfo, error)
    SearchPackages(query string) ([]*SearchResult, error)

    // Virtual environment operations
    CreateVenv(path string) error
    ActivateVenv(path string) error
    DeactivateVenv() error

    // Project operations
    InitProject(path string, opts *ProjectOptions) error
    BuildProject(path string, opts *BuildOptions) error

    // System operations
    IsInstalled() (bool, error)
    GetVersion() (string, error)
}
```

### Key Data Types

```go
// Package specification for installation
type PackageSpec struct {
    Name           string   // Package name
    Version        string   // Version constraint
    Extras         []string // Extra dependencies
    Upgrade        bool     // Upgrade if already installed
    ForceReinstall bool     // Force reinstallation
    UserInstall    bool     // Install to user directory
    Editable       bool     // Editable installation
}

// Project initialization options
type ProjectOptions struct {
    Name            string   // Project name
    Version         string   // Initial version
    Description     string   // Project description
    Author          string   // Author name
    AuthorEmail     string   // Author email
    License         string   // License type
    Dependencies    []string // Runtime dependencies
    DevDependencies []string // Development dependencies
    CreateVenv      bool     // Create virtual environment
    CreateSrc       bool     // Create src/ directory
    CreateTests     bool     // Create tests/ directory
}

// Configuration options
type Config struct {
    PythonPath   string        // Python executable path
    Timeout      time.Duration // Operation timeout
    Retries      int           // Number of retries
    DefaultIndex string        // Default package index
    TrustedHosts []string      // Trusted hosts
    LogLevel     string        // Logging level
    Environment  map[string]string // Environment variables
}
```

### Error Handling

```go
// Error types for different failure scenarios
const (
    ErrorTypePipNotInstalled    ErrorType = "pip_not_installed"
    ErrorTypePythonNotFound     ErrorType = "python_not_found"
    ErrorTypePackageNotFound    ErrorType = "package_not_found"
    ErrorTypePermissionDenied   ErrorType = "permission_denied"
    ErrorTypeNetworkError       ErrorType = "network_error"
    ErrorTypeVersionConflict    ErrorType = "version_conflict"
    ErrorTypeCommandFailed      ErrorType = "command_failed"
)

// Check error types
if pip.IsErrorType(err, pip.ErrorTypeNetworkError) {
    // Handle network errors specifically
}
```

## 🧪 Testing

Run the comprehensive test suite:

```bash
# Run all tests
go test ./...

# Run tests with coverage report
go test -cover ./...

# Generate detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run only unit tests (skip integration tests)
go test -short ./...

# Run integration tests (requires Python and pip)
go test -run Integration ./...

# Run benchmarks
go test -bench=. ./...

# Run tests with race detection
go test -race ./...
```

### Test Categories

- **Unit Tests**: Fast tests that don't require external dependencies
- **Integration Tests**: Tests that require Python and pip installation
- **Benchmark Tests**: Performance tests for critical operations
- **Example Tests**: Ensure documentation examples work correctly

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Quick Start for Contributors

1. **Fork the repository**
   ```bash
   git clone https://github.com/your-username/go-pip-sdk.git
   cd go-pip-sdk
   ```

2. **Set up development environment**
   ```bash
   # Install dependencies
   go mod download

   # Install development tools
   make install-tools

   # Run tests to ensure everything works
   make test
   ```

3. **Create your feature branch**
   ```bash
   git checkout -b feature/amazing-feature
   ```

4. **Make your changes and test**
   ```bash
   # Run tests
   make test

   # Run linting
   make lint

   # Format code
   make fmt
   ```

5. **Commit and push**
   ```bash
   git commit -m 'feat: add amazing feature'
   git push origin feature/amazing-feature
   ```

6. **Open a Pull Request**

### Development Guidelines

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Write tests for new functionality
- Update documentation for API changes
- Use conventional commit messages
- Ensure all CI checks pass

## 📋 Requirements

### Runtime Requirements
- **Go**: 1.19 or later
- **Python**: 3.7 or later (for pip operations)
- **Operating System**: Windows 10+, macOS 10.15+, or Linux

### Development Requirements
- **Go**: 1.19 or later
- **Make**: For build automation
- **Git**: For version control
- **Python**: 3.7+ with pip (for integration tests)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

### Third-Party Licenses

This project uses several third-party libraries. See [THIRD_PARTY_LICENSES.md](THIRD_PARTY_LICENSES.md) for details.

## 🙏 Acknowledgments

- **Python pip team** - For creating the excellent pip package manager that inspired this project
- **Go team** - For providing an amazing programming language and standard library
- **Contributors** - Thanks to all the developers who have contributed to this project
- **Community** - Special thanks to users who provide feedback and report issues

## 📞 Support

### Getting Help

- 📖 **[Documentation](https://scagogogo.github.io/go-pip-sdk/)** - Comprehensive guides and API reference
- 🐛 **[Issue Tracker](https://github.com/scagogogo/go-pip-sdk/issues)** - Report bugs or request features
- 💬 **[Discussions](https://github.com/scagogogo/go-pip-sdk/discussions)** - Ask questions and share ideas
- 📧 **[Email](mailto:support@scagogogo.com)** - Direct support for enterprise users

### Enterprise Support

For enterprise users, we offer:
- Priority support and bug fixes
- Custom feature development
- Training and consulting services
- SLA-backed support agreements

Contact us at [enterprise@scagogogo.com](mailto:enterprise@scagogogo.com) for more information.

---

<div align="center">

**[⬆ Back to Top](#go-pip-sdk)**

Made with ❤️ by the Go Pip SDK team

</div>
