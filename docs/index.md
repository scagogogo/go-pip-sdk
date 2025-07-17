---
layout: home

hero:
  name: "Go Pip SDK"
  text: "Python Package Management in Go"
  tagline: A comprehensive Go SDK for managing Python pip operations, virtual environments, and Python projects with CLI support
  image:
    src: /logo.svg
    alt: Go Pip SDK
  actions:
    - theme: brand
      text: Get Started
      link: /guide/getting-started
    - theme: alt
      text: View on GitHub
      link: https://github.com/scagogogo/go-pip-sdk

features:
  - icon: 🚀
    title: Cross-Platform Support
    details: Works seamlessly on Windows, macOS, and Linux with automatic platform detection and adaptation.
  
  - icon: 📦
    title: Complete Package Management
    details: Install, uninstall, list, show, and freeze Python packages with full pip compatibility.
  
  - icon: 🐍
    title: Virtual Environment Management
    details: Create, activate, deactivate, and remove virtual environments with ease.
  
  - icon: 🏗️
    title: Project Initialization
    details: Bootstrap Python projects with standard structure, setup.py, pyproject.toml, and more.
  
  - icon: 🔧
    title: Automatic Pip Installation
    details: Detects and installs pip automatically if missing, supporting multiple installation methods.
  
  - icon: 📝
    title: Comprehensive Logging
    details: Detailed operation logs with multiple levels and customizable output formats.
  
  - icon: ⚡
    title: Rich Error Handling
    details: Structured error types with helpful suggestions and context-aware error messages.
  
  - icon: 🧪
    title: Well Tested
    details: Extensive unit and integration tests with 82.3% code coverage for reliability.
  
  - icon: 🔒
    title: Type Safe
    details: Fully typed Go interfaces with comprehensive documentation and examples.
---

## Quick Start

Install the SDK and start managing Python packages in your Go applications:

```bash
go get github.com/scagogogo/go-pip-sdk
```

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
    
    // Install a package
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

## Why Go Pip SDK?

- **Native Go Integration**: No need for external Python scripts or subprocess calls
- **Production Ready**: Battle-tested with comprehensive error handling and logging
- **Developer Friendly**: Clean APIs with extensive documentation and examples
- **Flexible Configuration**: Customizable settings for different environments and use cases

## Community

- 📖 [Documentation](https://scagogogo.github.io/go-pip-sdk/)
- 🐛 [Issue Tracker](https://github.com/scagogogo/go-pip-sdk/issues)
- 💬 [Discussions](https://github.com/scagogogo/go-pip-sdk/discussions)
- 📧 [Contributing Guide](/guide/contributing)

## License

Released under the [MIT License](https://github.com/scagogogo/go-pip-sdk/blob/main/LICENSE).
