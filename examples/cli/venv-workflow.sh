#!/bin/bash

# Virtual Environment Workflow Example
# This script demonstrates virtual environment operations with pip-cli

set -e

echo "=== Go Pip SDK CLI - Virtual Environment Workflow ==="
echo

# Build the CLI tool first
echo "Building pip-cli..."
make build
echo

# Create a temporary directory for our example
TEMP_DIR=$(mktemp -d)
VENV_PATH="$TEMP_DIR/example-venv"

echo "Working in temporary directory: $TEMP_DIR"
echo "Virtual environment path: $VENV_PATH"
echo

# Create virtual environment
echo "1. Creating virtual environment:"
./bin/pip-cli venv create "$VENV_PATH"
echo

# Check virtual environment info
echo "2. Checking virtual environment info:"
./bin/pip-cli venv info "$VENV_PATH"
echo

# Activate virtual environment
echo "3. Activating virtual environment:"
./bin/pip-cli venv activate "$VENV_PATH"
echo

# Install packages in virtual environment
echo "4. Installing packages in virtual environment:"
./bin/pip-cli install requests click
echo

# List packages in virtual environment
echo "5. Listing packages in virtual environment:"
./bin/pip-cli list
echo

# Show package info
echo "6. Showing click package info:"
./bin/pip-cli show click
echo

# Freeze packages
echo "7. Freezing packages in virtual environment:"
./bin/pip-cli freeze > "$TEMP_DIR/requirements.txt"
cat "$TEMP_DIR/requirements.txt"
echo

# Deactivate virtual environment
echo "8. Deactivating virtual environment:"
./bin/pip-cli venv deactivate
echo

# List packages after deactivation (should show system packages)
echo "9. Listing packages after deactivation:"
./bin/pip-cli list
echo

# Remove virtual environment
echo "10. Removing virtual environment:"
./bin/pip-cli venv remove "$VENV_PATH"
echo

# Clean up
echo "Cleaning up temporary directory..."
rm -rf "$TEMP_DIR"

echo "=== Virtual environment workflow example completed ==="
