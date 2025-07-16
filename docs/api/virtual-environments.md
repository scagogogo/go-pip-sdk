# Virtual Environments

The virtual environment API provides comprehensive functionality for creating, managing, and working with Python virtual environments across different platforms.

## Core Operations

### CreateVenv

```go
func (m *Manager) CreateVenv(path string) error
```

Creates a new virtual environment at the specified path.

**Parameters:**
- `path` (string): Path where the virtual environment should be created.

**Returns:**
- `error`: Error if creation fails.

**Creation Methods:**
The SDK tries multiple methods in order:
1. `python -m venv` (Python 3.3+)
2. `python -m virtualenv` (if virtualenv is installed)
3. `virtualenv` command (standalone installation)

**Example:**
```go
// Create virtual environment
venvPath := "/path/to/my-venv"
if err := manager.CreateVenv(venvPath); err != nil {
    return fmt.Errorf("failed to create virtual environment: %w", err)
}
fmt.Printf("Virtual environment created at: %s\n", venvPath)
```

### ActivateVenv

```go
func (m *Manager) ActivateVenv(path string) error
```

Activates a virtual environment by updating the manager's configuration to use the virtual environment's Python and pip executables.

**Parameters:**
- `path` (string): Path to the virtual environment to activate.

**Returns:**
- `error`: Error if activation fails.

**Effects:**
- Updates `PythonPath` and `PipPath` in manager configuration
- Sets `VIRTUAL_ENV` environment variable
- Modifies `PATH` to prioritize virtual environment binaries

**Example:**
```go
venvPath := "/path/to/my-venv"
if err := manager.ActivateVenv(venvPath); err != nil {
    return fmt.Errorf("failed to activate virtual environment: %w", err)
}

// Now all pip operations will use the virtual environment
pkg := &pip.PackageSpec{Name: "requests"}
if err := manager.InstallPackage(pkg); err != nil {
    return err
}
```

### DeactivateVenv

```go
func (m *Manager) DeactivateVenv() error
```

Deactivates the currently active virtual environment by resetting the manager configuration to use system Python and pip.

**Returns:**
- `error`: Error if deactivation fails.

**Example:**
```go
if err := manager.DeactivateVenv(); err != nil {
    return fmt.Errorf("failed to deactivate virtual environment: %w", err)
}
fmt.Println("Virtual environment deactivated")
```

### RemoveVenv

```go
func (m *Manager) RemoveVenv(path string) error
```

Removes a virtual environment by deleting its directory and contents.

**Parameters:**
- `path` (string): Path to the virtual environment to remove.

**Returns:**
- `error`: Error if removal fails.

**Example:**
```go
venvPath := "/path/to/my-venv"
if err := manager.RemoveVenv(venvPath); err != nil {
    return fmt.Errorf("failed to remove virtual environment: %w", err)
}
fmt.Printf("Virtual environment removed: %s\n", venvPath)
```

## Information and Validation

### GetVenvInfo

```go
func (m *Manager) GetVenvInfo(path string) (*VenvInfo, error)
```

Retrieves information about a virtual environment.

**Parameters:**
- `path` (string): Path to the virtual environment.

**Returns:**
- `*VenvInfo`: Information about the virtual environment.
- `error`: Error if information retrieval fails.

**Example:**
```go
info, err := manager.GetVenvInfo("/path/to/my-venv")
if err != nil {
    return fmt.Errorf("failed to get venv info: %w", err)
}

fmt.Printf("Virtual Environment Info:\n")
fmt.Printf("  Path: %s\n", info.Path)
fmt.Printf("  Active: %t\n", info.IsActive)
fmt.Printf("  Python Path: %s\n", info.PythonPath)
fmt.Printf("  Python Version: %s\n", info.PythonVersion)
fmt.Printf("  Created: %s\n", info.CreatedAt.Format("2006-01-02 15:04:05"))
```

## Data Types

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

### VenvManager

```go
type VenvManager struct {
    manager *Manager // Reference to the main manager
}
```

The `VenvManager` is an internal helper that provides platform-specific virtual environment operations.

## Platform-Specific Behavior

### Windows

- Virtual environment structure: `venv/Scripts/`
- Python executable: `python.exe`
- Pip executable: `pip.exe`
- Activation script: `Scripts/activate.bat`

### Unix-like (macOS, Linux)

- Virtual environment structure: `venv/bin/`
- Python executable: `python`
- Pip executable: `pip`
- Activation script: `bin/activate`

## Advanced Usage

### Creating Virtual Environment with Specific Python Version

```go
// Configure manager to use specific Python version
config := &pip.Config{
    PythonPath: "/usr/bin/python3.9",
}
manager := pip.NewManager(config)

// Create virtual environment with this Python version
if err := manager.CreateVenv("/path/to/venv"); err != nil {
    return err
}
```

### Working with Multiple Virtual Environments

```go
// Create multiple virtual environments
venvs := []string{
    "/path/to/dev-venv",
    "/path/to/test-venv",
    "/path/to/prod-venv",
}

for _, venvPath := range venvs {
    if err := manager.CreateVenv(venvPath); err != nil {
        fmt.Printf("Failed to create %s: %v\n", venvPath, err)
        continue
    }
    
    // Activate and install packages
    if err := manager.ActivateVenv(venvPath); err != nil {
        continue
    }
    
    // Install environment-specific packages
    pkg := &pip.PackageSpec{Name: "requests"}
    manager.InstallPackage(pkg)
    
    // Deactivate before moving to next
    manager.DeactivateVenv()
}
```

### Installing Requirements in Virtual Environment

```go
// Create and activate virtual environment
venvPath := "/path/to/project-venv"
if err := manager.CreateVenv(venvPath); err != nil {
    return err
}

if err := manager.ActivateVenv(venvPath); err != nil {
    return err
}

// Install requirements
if err := manager.InstallRequirements("requirements.txt"); err != nil {
    return err
}

fmt.Println("Virtual environment set up with requirements")
```

## Error Handling

Virtual environment operations can fail for various reasons:

```go
if err := manager.CreateVenv(venvPath); err != nil {
    switch pip.GetErrorType(err) {
    case pip.ErrorTypeVenvAlreadyExists:
        fmt.Printf("Virtual environment already exists at %s\n", venvPath)
    case pip.ErrorTypePythonNotFound:
        fmt.Println("Python interpreter not found")
    case pip.ErrorTypePermissionDenied:
        fmt.Println("Permission denied - check directory permissions")
    case pip.ErrorTypeInvalidPath:
        fmt.Println("Invalid path specified")
    default:
        fmt.Printf("Failed to create virtual environment: %v\n", err)
    }
}
```

## Best Practices

1. **Always check if virtual environment exists before creating**:
   ```go
   if info, err := manager.GetVenvInfo(venvPath); err == nil {
       fmt.Printf("Virtual environment already exists: %s\n", info.Path)
       return nil
   }
   ```

2. **Use absolute paths for virtual environments**:
   ```go
   venvPath, err := filepath.Abs("./my-venv")
   if err != nil {
       return err
   }
   manager.CreateVenv(venvPath)
   ```

3. **Clean up virtual environments when done**:
   ```go
   defer func() {
       if err := manager.RemoveVenv(venvPath); err != nil {
           fmt.Printf("Warning: failed to clean up venv: %v\n", err)
       }
   }()
   ```

4. **Validate virtual environment before use**:
   ```go
   info, err := manager.GetVenvInfo(venvPath)
   if err != nil {
       return fmt.Errorf("invalid virtual environment: %w", err)
   }

   if !info.IsActive {
       if err := manager.ActivateVenv(venvPath); err != nil {
           return err
       }
   }
   ```

## Common Patterns

### Project Setup with Virtual Environment

```go
func SetupProject(projectPath string) error {
    manager := pip.NewManager(nil)

    // Create virtual environment in project directory
    venvPath := filepath.Join(projectPath, "venv")
    if err := manager.CreateVenv(venvPath); err != nil {
        return err
    }

    // Activate virtual environment
    if err := manager.ActivateVenv(venvPath); err != nil {
        return err
    }

    // Install project dependencies
    reqPath := filepath.Join(projectPath, "requirements.txt")
    if _, err := os.Stat(reqPath); err == nil {
        if err := manager.InstallRequirements(reqPath); err != nil {
            return err
        }
    }

    return nil
}
```

### Testing with Isolated Environments

```go
func TestWithCleanEnvironment(t *testing.T) {
    manager := pip.NewManager(nil)

    // Create temporary virtual environment
    tempDir, err := os.MkdirTemp("", "test-venv-*")
    if err != nil {
        t.Fatal(err)
    }
    defer os.RemoveAll(tempDir)

    venvPath := filepath.Join(tempDir, "venv")
    if err := manager.CreateVenv(venvPath); err != nil {
        t.Fatal(err)
    }

    if err := manager.ActivateVenv(venvPath); err != nil {
        t.Fatal(err)
    }

    // Run tests with clean environment
    // ...
}
```
