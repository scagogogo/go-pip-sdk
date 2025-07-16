# pip-cli

A command-line interface for the Go Pip SDK that provides a convenient way to manage Python packages, virtual environments, and projects from the command line.

## Installation

### Build from Source

```bash
# Clone the repository
git clone https://github.com/scagogogo/go-pip-sdk.git
cd go-pip-sdk

# Build the CLI tool
make build

# Install globally (optional)
make install
```

### Using Go Install

```bash
go install github.com/scagogogo/go-pip-sdk/cmd/pip-cli@latest
```

## Usage

### Global Options

- `-timeout duration`: Timeout for operations (default: 30s)
- `-verbose`: Enable verbose logging
- `-python string`: Path to Python executable
- `-pip string`: Path to pip executable

### Commands

#### Package Management

**Install a package:**
```bash
pip-cli install requests
pip-cli install "requests>=2.25.0"
pip-cli install django ">=4.0,<5.0"
```

**Uninstall a package:**
```bash
pip-cli uninstall requests
```

**List installed packages:**
```bash
pip-cli list
```

**Show package information:**
```bash
pip-cli show requests
```

**Freeze packages (requirements format):**
```bash
pip-cli freeze
pip-cli freeze > requirements.txt
```

#### Virtual Environment Management

**Create a virtual environment:**
```bash
pip-cli venv create ./myenv
pip-cli venv create /path/to/venv
```

**Activate a virtual environment:**
```bash
pip-cli venv activate ./myenv
```

**Deactivate current virtual environment:**
```bash
pip-cli venv deactivate
```

**Remove a virtual environment:**
```bash
pip-cli venv remove ./myenv
```

**Get virtual environment info:**
```bash
pip-cli venv info ./myenv
```

#### Project Management

**Initialize a new Python project:**
```bash
pip-cli project init ./myproject
```

This creates a new project with:
- `setup.py` and `pyproject.toml`
- `requirements.txt`
- Basic project structure
- Virtual environment (in `./myproject/venv`)

#### Utility Commands

**Show version:**
```bash
pip-cli version
```

**Get help:**
```bash
pip-cli help
pip-cli help install
pip-cli help venv
```

## Examples

### Basic Workflow

```bash
# Create a new project
pip-cli project init ./my-web-app

# Navigate to project
cd my-web-app

# Activate the virtual environment
pip-cli venv activate ./venv

# Install dependencies
pip-cli install fastapi uvicorn

# List installed packages
pip-cli list

# Freeze dependencies
pip-cli freeze > requirements.txt
```

### Working with Existing Projects

```bash
# Activate existing virtual environment
pip-cli venv activate ./venv

# Install from requirements
pip-cli install -r requirements.txt

# Show package information
pip-cli show fastapi

# Upgrade a package
pip-cli install --upgrade fastapi
```

### Verbose Mode

Enable verbose logging to see detailed operations:

```bash
pip-cli -verbose install requests
pip-cli -verbose venv create ./myenv
```

### Custom Python/Pip Paths

Specify custom Python or pip executables:

```bash
pip-cli -python /usr/local/bin/python3.9 -pip /usr/local/bin/pip3.9 install requests
```

### Timeout Configuration

Set custom timeout for long-running operations:

```bash
pip-cli -timeout 5m install tensorflow
```

## Error Handling

The CLI tool provides detailed error messages and suggestions:

```bash
$ pip-cli install nonexistent-package
Installation failed: [package_not_found] package 'nonexistent-package' not found | Suggestions: Check the package name spelling, Search for the package on PyPI
```

## Integration with Go Pip SDK

The CLI tool is built on top of the Go Pip SDK and demonstrates how to use the SDK in real applications. You can use it as a reference for building your own tools.

## Development

### Building

```bash
make build
```

### Testing

```bash
make test
make test-coverage
```

### Linting

```bash
make lint
make fmt
make vet
```

### Cross-Platform Builds

```bash
make build-all
```

This creates binaries for:
- Linux (amd64)
- Windows (amd64)
- macOS (amd64)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Run `make release-check`
6. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.
