# Installation

This guide covers how to install and set up the Go Pip SDK in your Go project.

## Prerequisites

- Go 1.19 or later
- Python 3.7 or later (for pip operations)
- pip installed on your system

## Installing the SDK

### Using go get

```bash
go get github.com/scagogogo/go-pip-sdk
```

### Using go mod

Add the dependency to your `go.mod` file:

```go
module your-project

go 1.19

require (
    github.com/scagogogo/go-pip-sdk v1.0.0
)
```

Then run:

```bash
go mod tidy
```

## Verifying Installation

Create a simple test file to verify the installation:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    // Create a new pip manager
    manager := pip.NewManager(nil)
    
    // Check if pip is installed
    installed, err := manager.IsInstalled()
    if err != nil {
        log.Fatal(err)
    }
    
    if installed {
        fmt.Println("✓ Pip is installed and accessible")
        
        // Get pip version
        version, err := manager.GetVersion()
        if err == nil {
            fmt.Printf("✓ Pip version: %s\n", version)
        }
    } else {
        fmt.Println("✗ Pip is not installed or not accessible")
    }
}
```

Run the test:

```bash
go run main.go
```

Expected output:
```
✓ Pip is installed and accessible
✓ Pip version: 23.3.1
```

## Installing the CLI Tool

The SDK also includes a command-line interface tool:

```bash
# Clone the repository
git clone https://github.com/scagogogo/go-pip-sdk.git
cd go-pip-sdk

# Build the CLI tool
make build

# The binary will be available at bin/pip-cli
./bin/pip-cli --help
```

### Installing CLI Globally

To install the CLI tool globally:

```bash
# Build and install
make install

# Or manually copy to your PATH
cp bin/pip-cli /usr/local/bin/
```

## System Requirements

### Python and Pip

The SDK requires Python and pip to be installed on your system:

**macOS:**
```bash
# Using Homebrew
brew install python

# Pip is included with Python 3.4+
```

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install python3 python3-pip
```

**Windows:**
- Download Python from [python.org](https://www.python.org/downloads/)
- Pip is included with Python 3.4+

### Verifying Python/Pip Installation

```bash
python3 --version
pip3 --version
```

## Configuration

The SDK can be configured with custom Python and pip paths:

```go
config := &pip.Config{
    PythonPath: "/usr/local/bin/python3",
    PipPath:    "/usr/local/bin/pip3",
    Timeout:    30 * time.Second,
}

manager := pip.NewManager(config)
```

## Next Steps

- Read the [Getting Started](./getting-started.md) guide
- Explore [Configuration](./configuration.md) options
- Check out [Examples](../examples/) for common use cases
