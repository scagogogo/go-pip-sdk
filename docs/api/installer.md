# Installer

The installer component provides cross-platform pip installation functionality, automatically detecting the best installation method for each operating system.

## Core Functions

### Install

```go
func (m *Manager) Install() error
```

Installs pip on the system using the most appropriate method for the current platform.

**Returns:**
- `error`: Error if installation fails.

**Installation Methods by Platform:**

#### Windows
1. **ensurepip module**: `python -m ensurepip --upgrade`
2. **get-pip.py**: Downloads and executes the official installer
3. **Chocolatey**: `choco install python` (if available)
4. **Scoop**: `scoop install python` (if available)

#### macOS
1. **ensurepip module**: `python3 -m ensurepip --upgrade`
2. **Homebrew**: `brew install python` (if available)
3. **MacPorts**: `port install python39 +universal` (if available)
4. **get-pip.py**: Downloads and executes the official installer

#### Linux
1. **System package manager**: 
   - Ubuntu/Debian: `apt-get install python3-pip`
   - CentOS/RHEL: `yum install python3-pip` or `dnf install python3-pip`
   - Arch: `pacman -S python-pip`
   - Alpine: `apk add py3-pip`
2. **ensurepip module**: `python3 -m ensurepip --upgrade`
3. **get-pip.py**: Downloads and executes the official installer

## Installation Detection

### IsInstalled

```go
func (m *Manager) IsInstalled() (bool, error)
```

Checks if pip is installed and accessible on the system.

**Returns:**
- `bool`: `true` if pip is installed and working, `false` otherwise.
- `error`: Error if the check fails.

**Detection Process:**
1. Check if pip executable exists in PATH
2. Verify pip can be executed (`pip --version`)
3. Validate pip functionality (`pip list`)

### GetVersion

```go
func (m *Manager) GetVersion() (string, error)
```

Gets the installed pip version.

**Returns:**
- `string`: Pip version string (e.g., "23.2.1").
- `error`: Error if version retrieval fails.

## Platform-Specific Behavior

### Windows Installation

```go
// Windows-specific installation logic
func (i *Installer) installWindows() error {
    // Try ensurepip first
    if err := i.tryEnsurepip(); err == nil {
        return nil
    }
    
    // Try get-pip.py
    if err := i.tryGetPip(); err == nil {
        return nil
    }
    
    // Try package managers
    if err := i.tryChocolatey(); err == nil {
        return nil
    }
    
    return i.tryScoop()
}
```

### macOS Installation

```go
// macOS-specific installation logic
func (i *Installer) installMacOS() error {
    // Try ensurepip first
    if err := i.tryEnsurepip(); err == nil {
        return nil
    }
    
    // Try Homebrew
    if err := i.tryHomebrew(); err == nil {
        return nil
    }
    
    // Try MacPorts
    if err := i.tryMacPorts(); err == nil {
        return nil
    }
    
    return i.tryGetPip()
}
```

### Linux Installation

```go
// Linux-specific installation logic
func (i *Installer) installLinux() error {
    // Try system package manager first
    if err := i.trySystemPackageManager(); err == nil {
        return nil
    }
    
    // Try ensurepip
    if err := i.tryEnsurepip(); err == nil {
        return nil
    }
    
    return i.tryGetPip()
}
```

## Error Handling

The installer provides specific error types for different failure scenarios:

```go
if err := manager.Install(); err != nil {
    switch pip.GetErrorType(err) {
    case pip.ErrorTypePythonNotFound:
        fmt.Println("Python is not installed or not in PATH")
    case pip.ErrorTypePermissionDenied:
        fmt.Println("Permission denied - try running with elevated privileges")
    case pip.ErrorTypeNetworkError:
        fmt.Println("Network error - check internet connection")
    case pip.ErrorTypeUnsupportedOS:
        fmt.Println("Unsupported operating system")
    default:
        fmt.Printf("Installation failed: %v\n", err)
    }
}
```

## Advanced Usage

### Custom Installation Options

```go
// Configure custom installation behavior
config := &pip.Config{
    PythonPath: "/usr/local/bin/python3.9",  // Specific Python version
    Timeout:    300 * time.Second,           // Extended timeout
    Retries:    5,                           // More retries
}

manager := pip.NewManager(config)
if err := manager.Install(); err != nil {
    return err
}
```

### Installation with Logging

```go
// Enable detailed logging during installation
logger, err := pip.NewLogger(&pip.LoggerConfig{
    Level:  pip.LogLevelDebug,
    Output: os.Stdout,
})
if err != nil {
    return err
}

manager := pip.NewManager(nil)
manager.SetCustomLogger(logger)

// Installation will be logged in detail
if err := manager.Install(); err != nil {
    return err
}
```

### Offline Installation

```go
// For environments without internet access
// Pre-download get-pip.py and use local file
config := &pip.Config{
    ExtraOptions: map[string]string{
        "get-pip-path": "/path/to/get-pip.py",
    },
}

manager := pip.NewManager(config)
if err := manager.Install(); err != nil {
    return err
}
```

## Validation

### Post-Installation Verification

```go
func verifyInstallation(manager *pip.Manager) error {
    // Check if pip is installed
    installed, err := manager.IsInstalled()
    if err != nil {
        return fmt.Errorf("failed to check installation: %w", err)
    }
    
    if !installed {
        return errors.New("pip installation verification failed")
    }
    
    // Check pip version
    version, err := manager.GetVersion()
    if err != nil {
        return fmt.Errorf("failed to get version: %w", err)
    }
    
    fmt.Printf("Pip %s installed successfully\n", version)
    
    // Test basic functionality
    packages, err := manager.ListPackages()
    if err != nil {
        return fmt.Errorf("pip functionality test failed: %w", err)
    }
    
    fmt.Printf("Pip is working correctly (%d packages found)\n", len(packages))
    return nil
}
```

## Best Practices

1. **Always check before installing**:
   ```go
   if installed, err := manager.IsInstalled(); err != nil || !installed {
       if err := manager.Install(); err != nil {
           return err
       }
   }
   ```

2. **Handle platform differences gracefully**:
   ```go
   osType := pip.GetOSType()
   switch osType {
   case pip.OSWindows:
       // Windows-specific handling
   case pip.OSMacOS:
       // macOS-specific handling
   case pip.OSLinux:
       // Linux-specific handling
   }
   ```

3. **Use appropriate timeouts**:
   ```go
   config := &pip.Config{
       Timeout: 300 * time.Second,  // 5 minutes for installation
   }
   ```

4. **Verify installation success**:
   ```go
   if err := manager.Install(); err != nil {
       return err
   }
   
   // Verify the installation worked
   if err := verifyInstallation(manager); err != nil {
       return err
   }
   ```
