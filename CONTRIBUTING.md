# Contributing to Go Pip SDK

Thank you for your interest in contributing to the Go Pip SDK! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Code Style](#code-style)
- [Documentation](#documentation)

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally
3. Set up the development environment
4. Create a new branch for your changes
5. Make your changes
6. Test your changes
7. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.19 or later
- Python 3.7 or later (for testing pip operations)
- pip (Python package installer)
- make (for using the Makefile)

### Setup Steps

1. Clone the repository:
```bash
git clone https://github.com/your-username/go-pip-sdk.git
cd go-pip-sdk
```

2. Set up the development environment:
```bash
make dev-setup
```

3. Verify the setup:
```bash
make check-pip
make test-short
```

## Making Changes

### Branch Naming

Use descriptive branch names:
- `feature/add-new-functionality`
- `bugfix/fix-installation-issue`
- `docs/update-readme`
- `refactor/improve-error-handling`

### Commit Messages

Follow conventional commit format:
- `feat: add new package installation feature`
- `fix: resolve virtual environment activation issue`
- `docs: update API documentation`
- `test: add integration tests for pip operations`
- `refactor: improve error handling structure`

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run only unit tests (skip integration tests)
make test-short

# Run integration tests (requires pip installation)
make test-integration

# Run tests with coverage
make test-coverage

# Run benchmarks
make benchmark
```

### Test Categories

1. **Unit Tests**: Test individual functions and methods in isolation
2. **Integration Tests**: Test actual pip operations (require pip installation)
3. **Benchmark Tests**: Performance testing

### Writing Tests

- Place test files next to the code they test with `_test.go` suffix
- Use table-driven tests where appropriate
- Include both positive and negative test cases
- Mock external dependencies when possible
- Use descriptive test names

Example test structure:
```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr bool
    }{
        {
            name:    "valid input",
            input:   validInput,
            want:    expectedOutput,
            wantErr: false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := FunctionName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("FunctionName() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("FunctionName() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Code Style

### Go Style Guidelines

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Use `go vet` to check for common mistakes
- Follow the project's linting rules (see `.golangci.yml`)

### Code Formatting

```bash
# Format code
make fmt

# Vet code
make vet

# Run linter
make lint
```

### Naming Conventions

- Use descriptive names for variables, functions, and types
- Follow Go naming conventions (camelCase for private, PascalCase for public)
- Use meaningful package names
- Avoid abbreviations unless they're well-known

### Error Handling

- Use the custom error types defined in `pkg/pip/errors.go`
- Provide helpful error messages with suggestions
- Wrap errors with context when appropriate
- Use the error handler for consistent logging

Example:
```go
if err != nil {
    return NewPipError(ErrorTypeCommandFailed, "failed to install package").
        WithSuggestion("check package name spelling").
        WithContext("package", pkg.Name).
        WithCause(err)
}
```

## Documentation

### Code Documentation

- Add godoc comments for all public functions, types, and packages
- Include examples in documentation where helpful
- Keep documentation up to date with code changes

### README Updates

- Update the README.md if you add new features
- Include usage examples for new functionality
- Update the feature list if applicable

### API Documentation

- Document all public APIs
- Include parameter descriptions
- Provide usage examples
- Document error conditions

## Submitting Changes

### Pull Request Process

1. Ensure your code follows the style guidelines
2. Add or update tests for your changes
3. Update documentation as needed
4. Run the full test suite
5. Create a pull request with a clear description

### Pull Request Template

```markdown
## Description
Brief description of the changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Added new tests for changes

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No breaking changes (or clearly documented)
```

### Review Process

1. Automated checks must pass (tests, linting, etc.)
2. At least one maintainer review required
3. Address review feedback
4. Maintainer will merge when approved

## Release Process

Releases are handled by maintainers:

1. Version bump in relevant files
2. Update CHANGELOG.md
3. Create release tag
4. Build and publish binaries
5. Update documentation

## Getting Help

- Check existing issues and discussions
- Create a new issue for bugs or feature requests
- Join discussions for questions and ideas
- Contact maintainers for urgent matters

## Recognition

Contributors will be recognized in:
- CONTRIBUTORS.md file
- Release notes
- Project documentation

Thank you for contributing to Go Pip SDK!
