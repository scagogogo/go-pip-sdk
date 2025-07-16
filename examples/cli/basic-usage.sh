#!/bin/bash

# Basic CLI Usage Example
# This script demonstrates basic pip-cli operations

set -e

echo "=== Go Pip SDK CLI - Basic Usage Example ==="
echo

# Build the CLI tool first
echo "Building pip-cli..."
make build
echo

# Check version
echo "1. Checking CLI version:"
./bin/pip-cli version
echo

# Show help
echo "2. Showing help:"
./bin/pip-cli help
echo

# List currently installed packages
echo "3. Listing installed packages:"
./bin/pip-cli list
echo

# Install a simple package
echo "4. Installing requests package:"
./bin/pip-cli install requests
echo

# Show package information
echo "5. Showing requests package info:"
./bin/pip-cli show requests
echo

# List packages again to see the new installation
echo "6. Listing packages after installation:"
./bin/pip-cli list
echo

# Freeze packages
echo "7. Freezing packages to requirements format:"
./bin/pip-cli freeze
echo

# Uninstall the package
echo "8. Uninstalling requests package:"
./bin/pip-cli uninstall requests
echo

echo "=== Basic usage example completed ==="
