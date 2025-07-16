# Go Pip SDK Makefile

# Variables
BINARY_NAME=pip-cli
PACKAGE_PATH=./cmd/pip-cli
PKG_PATH=./pkg/...
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Build flags
BUILD_FLAGS=-ldflags="-s -w"
TEST_FLAGS=-v -race

.PHONY: all build clean test test-coverage test-integration lint fmt vet deps help

# Default target
all: clean fmt vet test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(BUILD_FLAGS) -o bin/$(BINARY_NAME) $(PACKAGE_PATH)
	@echo "Binary built: bin/$(BINARY_NAME)"

# Build for multiple platforms
build-all: build-linux build-windows build-darwin

build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o bin/$(BINARY_NAME)-linux-amd64 $(PACKAGE_PATH)

build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe $(PACKAGE_PATH)

build-darwin:
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 $(PACKAGE_PATH)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf bin/
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) $(TEST_FLAGS) $(PKG_PATH)

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) $(TEST_FLAGS) -coverprofile=$(COVERAGE_FILE) $(PKG_PATH)
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

# Run integration tests (requires pip installation)
test-integration:
	@echo "Running integration tests..."
	$(GOTEST) $(TEST_FLAGS) -run Integration $(PKG_PATH)

# Run short tests (skip integration tests)
test-short:
	@echo "Running short tests..."
	$(GOTEST) $(TEST_FLAGS) -short $(PKG_PATH)

# Lint the code
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Format the code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Vet the code
vet:
	@echo "Vetting code..."
	$(GOVET) $(PKG_PATH)

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	cp bin/$(BINARY_NAME) $(GOPATH)/bin/

# Run the CLI tool
run: build
	@echo "Running $(BINARY_NAME)..."
	./bin/$(BINARY_NAME) $(ARGS)

# Run examples
run-basic-example:
	@echo "Running basic example..."
	$(GOCMD) run examples/basic/main.go

run-venv-example:
	@echo "Running virtual environment example..."
	$(GOCMD) run examples/venv/main.go

# Development helpers
dev-setup:
	@echo "Setting up development environment..."
	$(GOMOD) download
	@echo "Installing development tools..."
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Check if pip is available for testing
check-pip:
	@echo "Checking pip availability..."
	@if command -v pip >/dev/null 2>&1; then \
		echo "✓ pip is available"; \
		pip --version; \
	else \
		echo "✗ pip is not available"; \
		echo "Some tests may fail without pip installed"; \
	fi

# Benchmark tests
benchmark:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem $(PKG_PATH)

# Generate documentation
docs:
	@echo "Generating documentation..."
	$(GOCMD) doc -all $(PKG_PATH)

# Security scan
security:
	@echo "Running security scan..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not installed. Install it with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Release preparation
release-check: clean fmt vet test-short lint
	@echo "Release check completed successfully"

# Help
help:
	@echo "Available targets:"
	@echo "  all              - Clean, format, vet, test, and build"
	@echo "  build            - Build the binary"
	@echo "  build-all        - Build for all platforms"
	@echo "  clean            - Clean build artifacts"
	@echo "  test             - Run all tests"
	@echo "  test-coverage    - Run tests with coverage report"
	@echo "  test-integration - Run integration tests only"
	@echo "  test-short       - Run tests excluding integration tests"
	@echo "  lint             - Run linter"
	@echo "  fmt              - Format code"
	@echo "  vet              - Vet code"
	@echo "  deps             - Download and tidy dependencies"
	@echo "  install          - Install binary to GOPATH/bin"
	@echo "  run              - Build and run CLI (use ARGS=... for arguments)"
	@echo "  run-basic-example - Run basic usage example"
	@echo "  run-venv-example  - Run virtual environment example"
	@echo "  dev-setup        - Set up development environment"
	@echo "  check-pip        - Check if pip is available"
	@echo "  benchmark        - Run benchmark tests"
	@echo "  docs             - Generate documentation"
	@echo "  security         - Run security scan"
	@echo "  release-check    - Run all checks for release"
	@echo "  help             - Show this help message"
