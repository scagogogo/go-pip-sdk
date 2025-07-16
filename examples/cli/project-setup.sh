#!/bin/bash

# Project Setup Example
# This script demonstrates project initialization and setup with pip-cli

set -e

echo "=== Go Pip SDK CLI - Project Setup Example ==="
echo

# Build the CLI tool first
echo "Building pip-cli..."
make build
echo

# Create a temporary directory for our example
TEMP_DIR=$(mktemp -d)
PROJECT_PATH="$TEMP_DIR/my-web-app"

echo "Working in temporary directory: $TEMP_DIR"
echo "Project path: $PROJECT_PATH"
echo

# Initialize new project
echo "1. Initializing new Python project:"
./bin/pip-cli project init "$PROJECT_PATH"
echo

# Show project structure
echo "2. Project structure created:"
find "$PROJECT_PATH" -type f | head -20
echo

# Navigate to project and activate virtual environment
echo "3. Activating project virtual environment:"
./bin/pip-cli venv activate "$PROJECT_PATH/venv"
echo

# Install web development packages
echo "4. Installing web development packages:"
./bin/pip-cli install fastapi uvicorn jinja2
echo

# Install development dependencies
echo "5. Installing development dependencies:"
./bin/pip-cli install pytest black flake8
echo

# List all installed packages
echo "6. Listing all installed packages:"
./bin/pip-cli list
echo

# Show specific package info
echo "7. Showing FastAPI package info:"
./bin/pip-cli show fastapi
echo

# Generate requirements file
echo "8. Generating requirements.txt:"
./bin/pip-cli freeze > "$PROJECT_PATH/requirements.txt"
echo "Requirements file content:"
cat "$PROJECT_PATH/requirements.txt"
echo

# Create a simple Python application
echo "9. Creating a simple FastAPI application:"
cat > "$PROJECT_PATH/app.py" << 'EOF'
from fastapi import FastAPI

app = FastAPI()

@app.get("/")
def read_root():
    return {"Hello": "World"}

@app.get("/items/{item_id}")
def read_item(item_id: int, q: str = None):
    return {"item_id": item_id, "q": q}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
EOF

echo "FastAPI application created at $PROJECT_PATH/app.py"
echo

# Show final project structure
echo "10. Final project structure:"
find "$PROJECT_PATH" -type f | sort
echo

# Deactivate virtual environment
echo "11. Deactivating virtual environment:"
./bin/pip-cli venv deactivate
echo

# Clean up
echo "Cleaning up temporary directory..."
rm -rf "$TEMP_DIR"

echo "=== Project setup example completed ==="
echo "You can use this workflow to set up real Python projects!"
