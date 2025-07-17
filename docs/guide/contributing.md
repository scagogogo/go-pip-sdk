# Contributing to Go Pip SDK

Thank you for your interest in contributing to the Go Pip SDK! This guide will help you get started with contributing to the project.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Code Style](#code-style)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Reporting Issues](#reporting-issues)
- [Documentation](#documentation)
- [Community Guidelines](#community-guidelines)

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

- Go 1.19 or later
- Git
- Python 3.7+ (for testing pip functionality)
- Make (optional, for using Makefile commands)

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:

```bash
git clone https://github.com/YOUR_USERNAME/go-pip-sdk.git
cd go-pip-sdk
```

3. Add the upstream repository:

```bash
git remote add upstream https://github.com/scagogogo/go-pip-sdk.git
```

## Development Setup

### Install Dependencies

```bash
# Install Go dependencies
go mod download

# Install development tools
make install-tools
```

### Build the Project

```bash
# Build the project
make build

# Or manually
go build ./...
```

### Run Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific tests
go test ./pkg/pip/...
```

## Code Style

### Go Code Standards

We follow standard Go conventions:

- Use `gofmt` for formatting
- Use `golint` for linting
- Follow effective Go practices
- Write clear, self-documenting code
- Use meaningful variable and function names

### Code Formatting

Before submitting code, ensure it's properly formatted:

```bash
# Format code
make fmt

# Or manually
gofmt -w .
```

### Linting

Run linters to check code quality:

```bash
# Run all linters
make lint

# Run specific linters
golangci-lint run
```

## Testing

### Writing Tests

- Write unit tests for all new functionality
- Use table-driven tests where appropriate
- Mock external dependencies
- Aim for high test coverage (>80%)

### Test Structure

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected ExpectedType
        wantErr  bool
    }{
        {
            name:     "valid case",
            input:    validInput,
            expected: expectedOutput,
            wantErr:  false,
        },
        // Add more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := FunctionName(tt.input)
            
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### Integration Tests

For integration tests that require Python/pip:

```go
func TestIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    // Integration test code
}
```

Run integration tests:

```bash
# Run all tests including integration
make test-integration

# Skip integration tests
go test -short ./...
```

## Submitting Changes

### Commit Guidelines

Follow conventional commit format:

```
type(scope): description

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Examples:
```
feat(manager): add support for custom pip indexes
fix(installer): handle network timeout errors properly
docs(api): update manager documentation
```

### Pull Request Process

1. Create a feature branch:
```bash
git checkout -b feature/your-feature-name
```

2. Make your changes and commit:
```bash
git add .
git commit -m "feat: add new feature"
```

3. Push to your fork:
```bash
git push origin feature/your-feature-name
```

4. Create a Pull Request on GitHub

### Pull Request Requirements

- [ ] Code follows project style guidelines
- [ ] Tests pass locally
- [ ] New functionality includes tests
- [ ] Documentation is updated
- [ ] Commit messages follow conventional format
- [ ] No merge conflicts with main branch

## Reporting Issues

### Bug Reports

When reporting bugs, include:

- Go version
- Operating system
- Python/pip version
- Steps to reproduce
- Expected behavior
- Actual behavior
- Error messages/logs

### Feature Requests

For feature requests, provide:

- Clear description of the feature
- Use case and motivation
- Proposed API (if applicable)
- Examples of usage

## Documentation

### Code Documentation

- Document all public functions and types
- Use clear, concise comments
- Include examples in documentation
- Follow Go documentation conventions

Example:
```go
// InstallPackage installs a Python package using pip.
// It returns an error if the installation fails.
//
// Example:
//   pkg := &PackageSpec{Name: "requests", Version: ">=2.25.0"}
//   err := manager.InstallPackage(pkg)
//   if err != nil {
//       log.Fatal(err)
//   }
func (m *Manager) InstallPackage(pkg *PackageSpec) error {
    // Implementation
}
```

### Documentation Site

The documentation site is built with VitePress:

```bash
# Install dependencies
cd docs
npm install

# Start development server
npm run dev

# Build documentation
npm run build
```

## Community Guidelines

### Code of Conduct

- Be respectful and inclusive
- Welcome newcomers
- Focus on constructive feedback
- Help others learn and grow

### Communication

- Use GitHub issues for bug reports and feature requests
- Use GitHub discussions for questions and general discussion
- Be clear and concise in communication
- Provide context and examples

### Getting Help

- Check existing issues and documentation first
- Ask questions in GitHub discussions
- Provide minimal reproducible examples
- Be patient and respectful

## Release Process

Releases are handled by maintainers:

1. Version bump in appropriate files
2. Update CHANGELOG.md
3. Create release tag
4. Publish release notes

## License

By contributing, you agree that your contributions will be licensed under the same license as the project (MIT License).

## Thank You

Thank you for contributing to Go Pip SDK! Your contributions help make this project better for everyone.
