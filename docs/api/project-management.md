# Project Management

The project management API provides functionality for initializing Python projects with proper structure, configuration files, and virtual environments.

## Core Operations

### InitProject

```go
func (m *Manager) InitProject(path string, opts *ProjectOptions) error
```

Initializes a new Python project at the specified path with the given options.

**Parameters:**
- `path` (string): Path where the project should be created.
- `opts` (*ProjectOptions): Project configuration options.

**Returns:**
- `error`: Error if initialization fails.

**Example:**
```go
opts := &pip.ProjectOptions{
    Name:            "my-project",
    Version:         "0.1.0",
    Description:     "A sample Python project",
    Author:          "Your Name",
    AuthorEmail:     "your.email@example.com",
    License:         "MIT",
    Dependencies:    []string{"requests>=2.25.0"},
    DevDependencies: []string{"pytest>=6.0"},
    CreateVenv:      true,
    VenvPath:        "./venv",
}

if err := manager.InitProject("/path/to/project", opts); err != nil {
    return err
}
```

## Project Structure

When initializing a project, the following structure is created:

```
project/
├── setup.py              # Package setup script
├── pyproject.toml         # Modern Python project configuration
├── requirements.txt       # Runtime dependencies
├── requirements-dev.txt   # Development dependencies
├── README.md             # Project documentation
├── .gitignore            # Git ignore file
├── src/                  # Source code directory
│   └── project_name/     # Main package
│       └── __init__.py
├── tests/                # Test directory
│   └── __init__.py
└── venv/                 # Virtual environment (if CreateVenv is true)
```

## Configuration Files

### setup.py

Generated setup.py includes:
- Package metadata (name, version, author, etc.)
- Dependencies and development dependencies
- Entry points and scripts
- Package discovery

### pyproject.toml

Modern Python project configuration with:
- Build system requirements
- Project metadata
- Tool configurations (pytest, black, etc.)
- Dependency specifications

### requirements.txt

Lists runtime dependencies in pip-compatible format.

### requirements-dev.txt

Lists development dependencies for testing, linting, and documentation.

## Data Types

### ProjectOptions

```go
type ProjectOptions struct {
    Name            string   // Project name (required)
    Version         string   // Project version (default: "0.1.0")
    Description     string   // Project description
    Author          string   // Author name
    AuthorEmail     string   // Author email
    License         string   // License (default: "MIT")
    PythonVersion   string   // Required Python version (default: ">=3.7")
    Dependencies    []string // Runtime dependencies
    DevDependencies []string // Development dependencies
    Template        string   // Project template ("basic", "library", "application")
    CreateVenv      bool     // Create virtual environment
    VenvPath        string   // Virtual environment path (default: "./venv")
}
```

## Examples

### Basic Project

```go
opts := &pip.ProjectOptions{
    Name:        "hello-world",
    Version:     "1.0.0",
    Author:      "John Doe",
    AuthorEmail: "john@example.com",
}

manager.InitProject("./hello-world", opts)
```

### Library Project

```go
opts := &pip.ProjectOptions{
    Name:            "my-library",
    Version:         "0.1.0",
    Description:     "A useful Python library",
    Author:          "Jane Smith",
    AuthorEmail:     "jane@example.com",
    License:         "Apache-2.0",
    Dependencies:    []string{"numpy>=1.20.0", "pandas>=1.3.0"},
    DevDependencies: []string{"pytest>=6.0", "sphinx>=4.0"},
    Template:        "library",
    CreateVenv:      true,
}

manager.InitProject("./my-library", opts)
```

### Application Project

```go
opts := &pip.ProjectOptions{
    Name:            "web-app",
    Version:         "1.0.0",
    Description:     "A web application",
    Author:          "Team Lead",
    AuthorEmail:     "team@company.com",
    Dependencies:    []string{"fastapi>=0.68.0", "uvicorn>=0.15.0"},
    DevDependencies: []string{"pytest>=6.0", "black>=21.0", "flake8>=3.9.0"},
    Template:        "application",
    CreateVenv:      true,
    VenvPath:        "./env",
}

manager.InitProject("./web-app", opts)
```

## Best Practices

1. **Always specify a virtual environment**:
   ```go
   opts.CreateVenv = true
   opts.VenvPath = "./venv"
   ```

2. **Use semantic versioning**:
   ```go
   opts.Version = "1.0.0"  // MAJOR.MINOR.PATCH
   ```

3. **Pin dependency versions for reproducibility**:
   ```go
   opts.Dependencies = []string{
       "requests>=2.25.0,<3.0.0",
       "click>=7.0,<8.0",
   }
   ```

4. **Include development dependencies**:
   ```go
   opts.DevDependencies = []string{
       "pytest>=6.0",
       "black>=21.0",
       "flake8>=3.9.0",
       "mypy>=0.910",
   }
   ```
