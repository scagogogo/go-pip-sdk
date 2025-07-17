# Project Management

This guide covers how to manage Python projects using the Go Pip SDK's project management features.

## Overview

The Go Pip SDK provides comprehensive project management capabilities, including project initialization, dependency management, and project structure creation.

## Creating New Projects

### Basic Project Creation

```go
package main

import (
    "log"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    manager := pip.NewManager(nil)
    projectManager := pip.NewProjectManager(manager)
    
    // Create a new project with default options
    err := projectManager.InitProject("./myproject", nil)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println("Project created successfully")
}
```

### Project with Custom Options

```go
opts := &pip.ProjectOptions{
    Name:        "myproject",
    Version:     "1.0.0",
    Description: "My awesome Python project",
    Author:      "John Doe",
    Email:       "john@example.com",
    License:     "MIT",
    CreateVenv:  true,
    VenvPath:    "./venv",
    Dependencies: []string{
        "requests>=2.25.0",
        "flask>=2.0.0",
        "pytest>=6.0.0",
    },
    DevDependencies: []string{
        "black",
        "flake8",
        "mypy",
    },
}

err := projectManager.InitProject("./myproject", opts)
if err != nil {
    log.Fatal(err)
}
```

## Project Structure

When you create a project, the SDK generates the following structure:

```
myproject/
├── venv/                 # Virtual environment (if CreateVenv: true)
├── src/
│   └── myproject/
│       ├── __init__.py
│       └── main.py
├── tests/
│   ├── __init__.py
│   └── test_main.py
├── requirements.txt      # Production dependencies
├── requirements-dev.txt  # Development dependencies
├── setup.py             # Package setup file
├── pyproject.toml       # Modern Python project configuration
├── README.md            # Project documentation
├── .gitignore           # Git ignore file
└── LICENSE              # License file
```

## Managing Dependencies

### Installing Project Dependencies

```go
// Install production dependencies
err := projectManager.InstallRequirements("./myproject/requirements.txt")
if err != nil {
    log.Fatal(err)
}

// Install development dependencies
err = projectManager.InstallRequirements("./myproject/requirements-dev.txt")
if err != nil {
    log.Fatal(err)
}
```

### Generating Requirements Files

```go
// Generate requirements.txt from current environment
err := projectManager.GenerateRequirements("./myproject/requirements.txt")
if err != nil {
    log.Fatal(err)
}
```

### Adding Dependencies

```go
// Add a new dependency to the project
pkg := &pip.PackageSpec{
    Name:    "numpy",
    Version: ">=1.20.0",
}

err := manager.InstallPackage(pkg)
if err != nil {
    log.Fatal(err)
}

// Update requirements.txt
err = projectManager.GenerateRequirements("./myproject/requirements.txt")
if err != nil {
    log.Fatal(err)
}
```

## Working with Virtual Environments

### Project with Virtual Environment

```go
opts := &pip.ProjectOptions{
    Name:       "myproject",
    CreateVenv: true,
    VenvPath:   "./myproject/venv",
}

err := projectManager.InitProject("./myproject", opts)
if err != nil {
    log.Fatal(err)
}

// The virtual environment is automatically created and can be activated
venvManager := pip.NewVenvManager(manager)
err = venvManager.ActivateVenv("./myproject/venv")
if err != nil {
    log.Fatal(err)
}
```

## Complete Project Setup Example

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    manager := pip.NewManager(nil)
    projectManager := pip.NewProjectManager(manager)
    venvManager := pip.NewVenvManager(manager)
    
    projectPath := "./mywebapp"
    
    // 1. Create project with virtual environment
    fmt.Println("Creating project...")
    opts := &pip.ProjectOptions{
        Name:        "mywebapp",
        Version:     "0.1.0",
        Description: "A web application built with Flask",
        Author:      "Developer",
        License:     "MIT",
        CreateVenv:  true,
        VenvPath:    projectPath + "/venv",
        Dependencies: []string{
            "flask>=2.0.0",
            "requests>=2.25.0",
            "python-dotenv>=0.19.0",
        },
        DevDependencies: []string{
            "pytest>=6.0.0",
            "black>=21.0.0",
            "flake8>=3.9.0",
        },
    }
    
    err := projectManager.InitProject(projectPath, opts)
    if err != nil {
        log.Fatal(err)
    }
    
    // 2. Activate virtual environment
    fmt.Println("Activating virtual environment...")
    err = venvManager.ActivateVenv(projectPath + "/venv")
    if err != nil {
        log.Fatal(err)
    }
    
    // 3. Install dependencies
    fmt.Println("Installing dependencies...")
    err = projectManager.InstallRequirements(projectPath + "/requirements.txt")
    if err != nil {
        log.Fatal(err)
    }
    
    err = projectManager.InstallRequirements(projectPath + "/requirements-dev.txt")
    if err != nil {
        log.Fatal(err)
    }
    
    // 4. Verify installation
    fmt.Println("Verifying installation...")
    packages, err := manager.ListPackages()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Installed %d packages:\n", len(packages))
    for _, pkg := range packages {
        fmt.Printf("  %s %s\n", pkg.Name, pkg.Version)
    }
    
    fmt.Println("Project setup completed successfully!")
}
```

## Using CLI Tool

The SDK includes CLI commands for project management:

```bash
# Initialize a new project
pip-cli project init ./myproject

# The CLI will create the project structure and optionally:
# - Create a virtual environment
# - Install dependencies
# - Set up development tools
```

## Project Templates

### Web Application Template

```go
opts := &pip.ProjectOptions{
    Name:        "webapp",
    Description: "Web application with Flask",
    Dependencies: []string{
        "flask>=2.0.0",
        "flask-sqlalchemy>=2.5.0",
        "flask-migrate>=3.1.0",
        "python-dotenv>=0.19.0",
    },
    DevDependencies: []string{
        "pytest>=6.0.0",
        "pytest-flask>=1.2.0",
        "black>=21.0.0",
        "flake8>=3.9.0",
    },
}
```

### Data Science Template

```go
opts := &pip.ProjectOptions{
    Name:        "datascience",
    Description: "Data science project",
    Dependencies: []string{
        "numpy>=1.20.0",
        "pandas>=1.3.0",
        "matplotlib>=3.4.0",
        "seaborn>=0.11.0",
        "scikit-learn>=1.0.0",
        "jupyter>=1.0.0",
    },
    DevDependencies: []string{
        "pytest>=6.0.0",
        "black>=21.0.0",
        "flake8>=3.9.0",
    },
}
```

## Configuration Files

### setup.py

The SDK generates a `setup.py` file for package distribution:

```python
from setuptools import setup, find_packages

setup(
    name="myproject",
    version="1.0.0",
    description="My awesome Python project",
    author="John Doe",
    author_email="john@example.com",
    packages=find_packages(where="src"),
    package_dir={"": "src"},
    install_requires=[
        "requests>=2.25.0",
        "flask>=2.0.0",
    ],
    extras_require={
        "dev": [
            "pytest>=6.0.0",
            "black>=21.0.0",
            "flake8>=3.9.0",
        ]
    },
    python_requires=">=3.7",
)
```

### pyproject.toml

Modern Python projects use `pyproject.toml`:

```toml
[build-system]
requires = ["setuptools>=45", "wheel"]
build-backend = "setuptools.build_meta"

[project]
name = "myproject"
version = "1.0.0"
description = "My awesome Python project"
authors = [{name = "John Doe", email = "john@example.com"}]
license = {text = "MIT"}
requires-python = ">=3.7"
dependencies = [
    "requests>=2.25.0",
    "flask>=2.0.0",
]

[project.optional-dependencies]
dev = [
    "pytest>=6.0.0",
    "black>=21.0.0",
    "flake8>=3.9.0",
]
```

## Best Practices

1. **Use Virtual Environments**: Always create projects with virtual environments
2. **Pin Dependencies**: Specify version constraints for dependencies
3. **Separate Dev Dependencies**: Keep development and production dependencies separate
4. **Documentation**: Include comprehensive README.md files
5. **Version Control**: Initialize git repositories for projects
6. **Testing**: Set up testing frameworks from the start

## Next Steps

- Learn about [Error Handling](./error-handling.md)
- Explore [Logging](./logging.md) configuration
- Check out [Examples](../examples/) for more use cases
