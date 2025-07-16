package pip

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// ProjectManager handles Python project operations
type ProjectManager struct {
	manager *Manager
}

// NewProjectManager creates a new project manager
func NewProjectManager(manager *Manager) *ProjectManager {
	return &ProjectManager{
		manager: manager,
	}
}

// InitProject initializes a new Python project
func (m *Manager) InitProject(path string, opts *ProjectOptions) error {
	if path == "" {
		return &PipError{
			Type:    "invalid_path",
			Message: "project path cannot be empty",
		}
	}

	if opts == nil {
		opts = &ProjectOptions{}
	}

	m.logInfo("Initializing project: %s", path)

	// Create project directory if it doesn't exist
	if err := os.MkdirAll(path, 0755); err != nil {
		return &PipError{
			Type:    "project_creation_failed",
			Message: fmt.Sprintf("failed to create project directory: %v", err),
		}
	}

	projectManager := NewProjectManager(m)

	// Set default values
	if opts.Name == "" {
		opts.Name = filepath.Base(path)
	}
	if opts.Version == "" {
		opts.Version = "0.1.0"
	}
	if opts.License == "" {
		opts.License = "MIT"
	}
	if opts.PythonVersion == "" {
		opts.PythonVersion = ">=3.7"
	}

	// Create virtual environment if requested
	if opts.CreateVenv {
		venvPath := opts.VenvPath
		if venvPath == "" {
			venvPath = filepath.Join(path, "venv")
		}

		if err := m.CreateVenv(venvPath); err != nil {
			m.logError("Failed to create virtual environment: %v", err)
		} else {
			m.logInfo("Created virtual environment: %s", venvPath)
		}
	}

	// Create project files
	if err := projectManager.createProjectFiles(path, opts); err != nil {
		return err
	}

	m.logInfo("Project initialized successfully: %s", path)
	return nil
}

// createProjectFiles creates the standard Python project files
func (pm *ProjectManager) createProjectFiles(projectPath string, opts *ProjectOptions) error {
	// Create requirements.txt
	if err := pm.createRequirementsFile(projectPath, opts); err != nil {
		return err
	}

	// Create setup.py
	if err := pm.createSetupFile(projectPath, opts); err != nil {
		return err
	}

	// Create pyproject.toml
	if err := pm.createPyprojectFile(projectPath, opts); err != nil {
		return err
	}

	// Create README.md
	if err := pm.createReadmeFile(projectPath, opts); err != nil {
		return err
	}

	// Create .gitignore
	if err := pm.createGitignoreFile(projectPath); err != nil {
		return err
	}

	// Create main package directory
	if err := pm.createPackageDirectory(projectPath, opts); err != nil {
		return err
	}

	// Create extra files
	if err := pm.createExtraFiles(projectPath, opts); err != nil {
		return err
	}

	return nil
}

// createRequirementsFile creates requirements.txt
func (pm *ProjectManager) createRequirementsFile(projectPath string, opts *ProjectOptions) error {
	filePath := filepath.Join(projectPath, "requirements.txt")

	var content strings.Builder
	content.WriteString("# Production dependencies\n")

	for _, dep := range opts.Dependencies {
		content.WriteString(dep + "\n")
	}

	if len(opts.DevDependencies) > 0 {
		content.WriteString("\n# Development dependencies\n")
		content.WriteString("# Install with: pip install -r requirements-dev.txt\n")

		// Create requirements-dev.txt
		devFilePath := filepath.Join(projectPath, "requirements-dev.txt")
		var devContent strings.Builder
		devContent.WriteString("-r requirements.txt\n\n")
		devContent.WriteString("# Development dependencies\n")

		for _, dep := range opts.DevDependencies {
			devContent.WriteString(dep + "\n")
		}

		if err := os.WriteFile(devFilePath, []byte(devContent.String()), 0644); err != nil {
			return err
		}
	}

	return os.WriteFile(filePath, []byte(content.String()), 0644)
}

// createSetupFile creates setup.py
func (pm *ProjectManager) createSetupFile(projectPath string, opts *ProjectOptions) error {
	filePath := filepath.Join(projectPath, "setup.py")

	tmpl := `#!/usr/bin/env python3
"""Setup script for {{.Name}}."""

from setuptools import setup, find_packages

with open("README.md", "r", encoding="utf-8") as fh:
    long_description = fh.read()

with open("requirements.txt", "r", encoding="utf-8") as fh:
    requirements = [line.strip() for line in fh if line.strip() and not line.startswith("#")]

setup(
    name="{{.Name}}",
    version="{{.Version}}",
    author="{{.Author}}",
    author_email="{{.AuthorEmail}}",
    description="{{.Description}}",
    long_description=long_description,
    long_description_content_type="text/markdown",
    packages=find_packages(),
    classifiers=[
        "Development Status :: 3 - Alpha",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: {{.License}} License",
        "Operating System :: OS Independent",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.7",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
    ],
    python_requires="{{.PythonVersion}}",
    install_requires=requirements,
    entry_points={
        "console_scripts": [
            "{{.Name}}={{.Name}}.main:main",
        ],
    },
)
`

	t, err := template.New("setup").Parse(tmpl)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return t.Execute(file, opts)
}

// createPyprojectFile creates pyproject.toml
func (pm *ProjectManager) createPyprojectFile(projectPath string, opts *ProjectOptions) error {
	filePath := filepath.Join(projectPath, "pyproject.toml")

	tmpl := `[build-system]
requires = ["setuptools>=45", "wheel", "setuptools_scm[toml]>=6.2"]
build-backend = "setuptools.build_meta"

[project]
name = "{{.Name}}"
version = "{{.Version}}"
description = "{{.Description}}"
readme = "README.md"
requires-python = "{{.PythonVersion}}"
license = {text = "{{.License}}"}
authors = [
    {name = "{{.Author}}", email = "{{.AuthorEmail}}"},
]
classifiers = [
    "Development Status :: 3 - Alpha",
    "Intended Audience :: Developers",
    "License :: OSI Approved :: {{.License}} License",
    "Operating System :: OS Independent",
    "Programming Language :: Python :: 3",
]
dependencies = [
{{range .Dependencies}}    "{{.}}",
{{end}}]

[project.optional-dependencies]
dev = [
{{range .DevDependencies}}    "{{.}}",
{{end}}]

[project.scripts]
{{.Name}} = "{{.Name}}.main:main"

[tool.setuptools.packages.find]
where = ["."]
include = ["{{.Name}}*"]

[tool.black]
line-length = 88
target-version = ['py37']

[tool.isort]
profile = "black"
line_length = 88

[tool.pytest.ini_options]
testpaths = ["tests"]
python_files = ["test_*.py"]
python_classes = ["Test*"]
python_functions = ["test_*"]
`

	t, err := template.New("pyproject").Parse(tmpl)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return t.Execute(file, opts)
}

// createReadmeFile creates README.md
func (pm *ProjectManager) createReadmeFile(projectPath string, opts *ProjectOptions) error {
	filePath := filepath.Join(projectPath, "README.md")

	tmpl := `# {{.Name}}

{{.Description}}

## Installation

` + "```bash" + `
pip install {{.Name}}
` + "```" + `

## Usage

` + "```python" + `
import {{.Name}}

# Your code here
` + "```" + `

## Development

1. Clone the repository:
` + "```bash" + `
git clone <repository-url>
cd {{.Name}}
` + "```" + `

2. Create a virtual environment:
` + "```bash" + `
python -m venv venv
source venv/bin/activate  # On Windows: venv\\Scripts\\activate
` + "```" + `

3. Install dependencies:
` + "```bash" + `
pip install -r requirements-dev.txt
` + "```" + `

## License

This project is licensed under the {{.License}} License.
`

	t, err := template.New("readme").Parse(tmpl)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return t.Execute(file, opts)
}

// createGitignoreFile creates .gitignore
func (pm *ProjectManager) createGitignoreFile(projectPath string) error {
	filePath := filepath.Join(projectPath, ".gitignore")

	content := `# Byte-compiled / optimized / DLL files
__pycache__/
*.py[cod]
*$py.class

# C extensions
*.so

# Distribution / packaging
.Python
build/
develop-eggs/
dist/
downloads/
eggs/
.eggs/
lib/
lib64/
parts/
sdist/
var/
wheels/
*.egg-info/
.installed.cfg
*.egg
MANIFEST

# PyInstaller
*.manifest
*.spec

# Installer logs
pip-log.txt
pip-delete-this-directory.txt

# Unit test / coverage reports
htmlcov/
.tox/
.coverage
.coverage.*
.cache
nosetests.xml
coverage.xml
*.cover
.hypothesis/
.pytest_cache/

# Translations
*.mo
*.pot

# Django stuff:
*.log
local_settings.py
db.sqlite3

# Flask stuff:
instance/
.webassets-cache

# Scrapy stuff:
.scrapy

# Sphinx documentation
docs/_build/

# PyBuilder
target/

# Jupyter Notebook
.ipynb_checkpoints

# pyenv
.python-version

# celery beat schedule file
celerybeat-schedule

# SageMath parsed files
*.sage.py

# Environments
.env
.venv
env/
venv/
ENV/
env.bak/
venv.bak/

# Spyder project settings
.spyderproject
.spyproject

# Rope project settings
.ropeproject

# mkdocs documentation
/site

# mypy
.mypy_cache/
.dmypy.json
dmypy.json

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db
`

	return os.WriteFile(filePath, []byte(content), 0644)
}

// createPackageDirectory creates the main package directory
func (pm *ProjectManager) createPackageDirectory(projectPath string, opts *ProjectOptions) error {
	packagePath := filepath.Join(projectPath, opts.Name)

	// Create package directory
	if err := os.MkdirAll(packagePath, 0755); err != nil {
		return err
	}

	// Create __init__.py
	initFile := filepath.Join(packagePath, "__init__.py")
	initContent := fmt.Sprintf(`"""{{.Name}} package."""

__version__ = "{{.Version}}"
__author__ = "{{.Author}}"
__email__ = "{{.AuthorEmail}}"
`)

	t, err := template.New("init").Parse(initContent)
	if err != nil {
		return err
	}

	file, err := os.Create(initFile)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := t.Execute(file, opts); err != nil {
		return err
	}

	// Create main.py
	mainFile := filepath.Join(packagePath, "main.py")
	mainContent := `"""Main module."""

def main():
    """Main entry point."""
    print("Hello from {{.Name}}!")

if __name__ == "__main__":
    main()
`

	t, err = template.New("main").Parse(mainContent)
	if err != nil {
		return err
	}

	file, err = os.Create(mainFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return t.Execute(file, opts)
}

// createExtraFiles creates additional files specified in options
func (pm *ProjectManager) createExtraFiles(projectPath string, opts *ProjectOptions) error {
	for fileName, content := range opts.ExtraFiles {
		filePath := filepath.Join(projectPath, fileName)

		// Create directory if needed
		if dir := filepath.Dir(filePath); dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}

		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

// InstallRequirements installs packages from requirements.txt
func (m *Manager) InstallRequirements(path string) error {
	if path == "" {
		return &PipError{
			Type:    "invalid_path",
			Message: "requirements file path cannot be empty",
		}
	}

	m.logInfo("Installing requirements from: %s", path)

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &PipError{
			Type:    "file_not_found",
			Message: fmt.Sprintf("requirements file not found: %s", path),
		}
	}

	pipPath, err := m.findPipExecutable()
	if err != nil {
		return ErrPipNotInstalled
	}

	args := []string{"install", "-r", path}
	return m.executePipCommand(pipPath, args)
}

// GenerateRequirements generates requirements.txt file
func (m *Manager) GenerateRequirements(path string) error {
	if path == "" {
		return &PipError{
			Type:    "invalid_path",
			Message: "requirements file path cannot be empty",
		}
	}

	m.logInfo("Generating requirements file: %s", path)

	// Get frozen packages
	packages, err := m.FreezePackages()
	if err != nil {
		return err
	}

	// Create requirements content
	var content strings.Builder
	content.WriteString("# Generated requirements file\n")
	content.WriteString(fmt.Sprintf("# Generated on: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	for _, pkg := range packages {
		if pkg.Version != "" {
			content.WriteString(fmt.Sprintf("%s==%s\n", pkg.Name, pkg.Version))
		} else {
			content.WriteString(fmt.Sprintf("%s\n", pkg.Name))
		}
	}

	// Write to file
	return os.WriteFile(path, []byte(content.String()), 0644)
}
