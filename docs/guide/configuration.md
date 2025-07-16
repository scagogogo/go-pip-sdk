# Configuration

The Go Pip SDK provides extensive configuration options to customize its behavior for different environments and use cases.

## Basic Configuration

### Default Configuration

```go
// Create manager with default configuration
manager := pip.NewManager(nil)
```

The default configuration includes:
- Automatic Python/pip detection
- 30-second timeout for operations
- 3 retries for failed operations
- INFO log level
- Standard PyPI index

### Custom Configuration

```go
config := &pip.Config{
    PythonPath:   "/usr/bin/python3",
    PipPath:      "/usr/bin/pip3",
    Timeout:      60 * time.Second,
    Retries:      5,
    LogLevel:     "DEBUG",
    DefaultIndex: "https://pypi.org/simple/",
}

manager := pip.NewManager(config)
```

## Configuration Options

### Core Settings

#### PythonPath
Path to the Python executable.

```go
config.PythonPath = "/usr/local/bin/python3.9"
```

**Default**: Auto-detected from PATH

#### PipPath
Path to the pip executable.

```go
config.PipPath = "/usr/local/bin/pip3"
```

**Default**: Auto-detected from PATH

#### Timeout
Timeout duration for pip operations.

```go
config.Timeout = 120 * time.Second
```

**Default**: 30 seconds

#### Retries
Number of retries for failed operations.

```go
config.Retries = 5
```

**Default**: 3

### Package Index Settings

#### DefaultIndex
Default package index URL.

```go
config.DefaultIndex = "https://pypi.org/simple/"
```

**Default**: PyPI (https://pypi.org/simple/)

#### TrustedHosts
List of trusted hosts for package downloads.

```go
config.TrustedHosts = []string{
    "pypi.org",
    "pypi.python.org",
    "files.pythonhosted.org",
}
```

### Logging Configuration

#### LogLevel
Logging level for operations.

```go
config.LogLevel = "DEBUG"  // DEBUG, INFO, WARN, ERROR
```

**Default**: "INFO"

### Advanced Settings

#### CacheDir
Custom cache directory for pip.

```go
config.CacheDir = "/tmp/pip-cache"
```

**Default**: System default

#### ExtraOptions
Additional pip command-line options.

```go
config.ExtraOptions = map[string]string{
    "no-cache-dir":              "",
    "disable-pip-version-check": "",
    "timeout":                   "60",
}
```

#### Environment
Environment variables for pip operations.

```go
config.Environment = map[string]string{
    "PIP_DISABLE_PIP_VERSION_CHECK": "1",
    "PIP_NO_COLOR":                  "1",
    "PYTHONUNBUFFERED":              "1",
}
```

## Environment-Specific Configurations

### Development Environment

```go
func developmentConfig() *pip.Config {
    return &pip.Config{
        LogLevel:     "DEBUG",
        Timeout:      120 * time.Second,
        Retries:      5,
        DefaultIndex: "https://pypi.org/simple/",
        ExtraOptions: map[string]string{
            "verbose": "",
        },
        Environment: map[string]string{
            "PIP_DISABLE_PIP_VERSION_CHECK": "1",
        },
    }
}
```

### Production Environment

```go
func productionConfig() *pip.Config {
    return &pip.Config{
        LogLevel:     "WARN",
        Timeout:      300 * time.Second,
        Retries:      10,
        DefaultIndex: "https://pypi.org/simple/",
        TrustedHosts: []string{
            "pypi.org",
            "pypi.python.org",
            "files.pythonhosted.org",
        },
        ExtraOptions: map[string]string{
            "no-cache-dir": "",
            "quiet":        "",
        },
    }
}
```

### Corporate Environment

```go
func corporateConfig() *pip.Config {
    return &pip.Config{
        DefaultIndex: "https://pypi.company.com/simple/",
        TrustedHosts: []string{
            "pypi.company.com",
            "pypi.org",
        },
        Timeout: 180 * time.Second,
        Retries: 3,
        ExtraOptions: map[string]string{
            "cert":         "/etc/ssl/certs/company.pem",
            "trusted-host": "pypi.company.com",
        },
        Environment: map[string]string{
            "HTTPS_PROXY": "http://proxy.company.com:8080",
            "HTTP_PROXY":  "http://proxy.company.com:8080",
        },
    }
}
```

## Dynamic Configuration

### Configuration from Environment Variables

```go
func configFromEnv() *pip.Config {
    config := &pip.Config{}
    
    if pythonPath := os.Getenv("PYTHON_PATH"); pythonPath != "" {
        config.PythonPath = pythonPath
    }
    
    if pipPath := os.Getenv("PIP_PATH"); pipPath != "" {
        config.PipPath = pipPath
    }
    
    if indexURL := os.Getenv("PIP_INDEX_URL"); indexURL != "" {
        config.DefaultIndex = indexURL
    }
    
    if logLevel := os.Getenv("PIP_LOG_LEVEL"); logLevel != "" {
        config.LogLevel = logLevel
    }
    
    return config
}
```

### Configuration from File

```go
func configFromFile(path string) (*pip.Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    
    var config pip.Config
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

## Validation

### Configuration Validation

```go
func validateConfig(config *pip.Config) error {
    if config.Timeout <= 0 {
        return errors.New("timeout must be positive")
    }
    
    if config.Retries < 0 {
        return errors.New("retries cannot be negative")
    }
    
    if config.PythonPath != "" {
        if _, err := os.Stat(config.PythonPath); err != nil {
            return fmt.Errorf("python path not found: %w", err)
        }
    }
    
    return nil
}
```

## Best Practices

1. **Use environment-specific configurations**:
   ```go
   var config *pip.Config
   if os.Getenv("ENV") == "production" {
       config = productionConfig()
   } else {
       config = developmentConfig()
   }
   ```

2. **Validate configuration before use**:
   ```go
   if err := validateConfig(config); err != nil {
       return fmt.Errorf("invalid configuration: %w", err)
   }
   ```

3. **Use reasonable timeouts**:
   ```go
   config.Timeout = 300 * time.Second  // 5 minutes for large packages
   ```

4. **Configure appropriate retry counts**:
   ```go
   config.Retries = 5  // More retries for unreliable networks
   ```

5. **Set up proper logging**:
   ```go
   config.LogLevel = "INFO"  // Balance between verbosity and usefulness
   ```
