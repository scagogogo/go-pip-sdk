# Project Initialization Examples

This page provides practical examples of initializing Python projects using the Go Pip SDK.

## Basic Project Initialization

```go
package main

import (
    "log"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    manager := pip.NewManager(nil)
    projectManager := pip.NewProjectManager(manager)
    
    // Initialize a basic project
    err := projectManager.InitProject("./myproject", nil)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println("Project initialized successfully!")
}
```

## Web Application Project

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func createWebAppProject() error {
    manager := pip.NewManager(nil)
    projectManager := pip.NewProjectManager(manager)
    
    opts := &pip.ProjectOptions{
        Name:        "webapp",
        Version:     "1.0.0",
        Description: "A Flask web application",
        Author:      "Developer",
        Email:       "dev@example.com",
        License:     "MIT",
        CreateVenv:  true,
        VenvPath:    "./webapp/venv",
        Dependencies: []string{
            "flask>=2.0.0",
            "flask-sqlalchemy>=2.5.0",
            "flask-migrate>=3.1.0",
            "python-dotenv>=0.19.0",
            "gunicorn>=20.1.0",
        },
        DevDependencies: []string{
            "pytest>=6.0.0",
            "pytest-flask>=1.2.0",
            "black>=21.0.0",
            "flake8>=3.9.0",
            "mypy>=0.910",
        },
    }
    
    fmt.Println("Creating web application project...")
    err := projectManager.InitProject("./webapp", opts)
    if err != nil {
        return err
    }
    
    fmt.Println("Web application project created successfully!")
    return nil
}

func main() {
    err := createWebAppProject()
    if err != nil {
        log.Fatal(err)
    }
}
```

## Data Science Project

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func createDataScienceProject() error {
    manager := pip.NewManager(nil)
    projectManager := pip.NewProjectManager(manager)
    
    opts := &pip.ProjectOptions{
        Name:        "datascience-project",
        Version:     "0.1.0",
        Description: "Data science and machine learning project",
        Author:      "Data Scientist",
        Email:       "datascientist@example.com",
        License:     "MIT",
        CreateVenv:  true,
        VenvPath:    "./datascience-project/venv",
        Dependencies: []string{
            "numpy>=1.21.0",
            "pandas>=1.3.0",
            "matplotlib>=3.4.0",
            "seaborn>=0.11.0",
            "scikit-learn>=1.0.0",
            "jupyter>=1.0.0",
            "ipykernel>=6.0.0",
            "plotly>=5.0.0",
        },
        DevDependencies: []string{
            "pytest>=6.0.0",
            "black>=21.0.0",
            "flake8>=3.9.0",
            "mypy>=0.910",
            "pre-commit>=2.15.0",
        },
    }
    
    fmt.Println("Creating data science project...")
    err := projectManager.InitProject("./datascience-project", opts)
    if err != nil {
        return err
    }
    
    // Create additional data science specific directories
    additionalDirs := []string{
        "./datascience-project/data/raw",
        "./datascience-project/data/processed",
        "./datascience-project/notebooks",
        "./datascience-project/models",
        "./datascience-project/reports",
    }
    
    for _, dir := range additionalDirs {
        err := os.MkdirAll(dir, 0755)
        if err != nil {
            return err
        }
    }
    
    fmt.Println("Data science project created successfully!")
    return nil
}

func main() {
    err := createDataScienceProject()
    if err != nil {
        log.Fatal(err)
    }
}
```

## API Project with FastAPI

```go
package main

import (
    "fmt"
    "log"
    "os"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func createAPIProject() error {
    manager := pip.NewManager(nil)
    projectManager := pip.NewProjectManager(manager)
    
    opts := &pip.ProjectOptions{
        Name:        "api-project",
        Version:     "1.0.0",
        Description: "FastAPI REST API project",
        Author:      "API Developer",
        Email:       "api@example.com",
        License:     "MIT",
        CreateVenv:  true,
        VenvPath:    "./api-project/venv",
        Dependencies: []string{
            "fastapi>=0.70.0",
            "uvicorn[standard]>=0.15.0",
            "pydantic>=1.8.0",
            "sqlalchemy>=1.4.0",
            "alembic>=1.7.0",
            "python-multipart>=0.0.5",
            "python-jose[cryptography]>=3.3.0",
            "passlib[bcrypt]>=1.7.4",
        },
        DevDependencies: []string{
            "pytest>=6.0.0",
            "pytest-asyncio>=0.16.0",
            "httpx>=0.24.0",
            "black>=21.0.0",
            "flake8>=3.9.0",
            "mypy>=0.910",
        },
    }
    
    fmt.Println("Creating API project...")
    err := projectManager.InitProject("./api-project", opts)
    if err != nil {
        return err
    }
    
    // Create API-specific structure
    apiDirs := []string{
        "./api-project/app/api/v1",
        "./api-project/app/core",
        "./api-project/app/models",
        "./api-project/app/schemas",
        "./api-project/app/crud",
        "./api-project/app/db",
        "./api-project/alembic/versions",
    }
    
    for _, dir := range apiDirs {
        err := os.MkdirAll(dir, 0755)
        if err != nil {
            return err
        }
    }
    
    // Create main.py
    mainPy := `from fastapi import FastAPI

app = FastAPI(title="API Project", version="1.0.0")

@app.get("/")
async def root():
    return {"message": "Hello World"}

@app.get("/health")
async def health_check():
    return {"status": "healthy"}
`
    
    err = os.WriteFile("./api-project/app/main.py", []byte(mainPy), 0644)
    if err != nil {
        return err
    }
    
    fmt.Println("API project created successfully!")
    return nil
}

func main() {
    err := createAPIProject()
    if err != nil {
        log.Fatal(err)
    }
}
```

## CLI Tool Project

```go
package main

import (
    "fmt"
    "log"
    "os"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func createCLIProject() error {
    manager := pip.NewManager(nil)
    projectManager := pip.NewProjectManager(manager)
    
    opts := &pip.ProjectOptions{
        Name:        "cli-tool",
        Version:     "1.0.0",
        Description: "Command-line interface tool",
        Author:      "CLI Developer",
        Email:       "cli@example.com",
        License:     "MIT",
        CreateVenv:  true,
        VenvPath:    "./cli-tool/venv",
        Dependencies: []string{
            "click>=8.0.0",
            "rich>=10.0.0",
            "typer>=0.4.0",
            "pydantic>=1.8.0",
            "pyyaml>=5.4.0",
        },
        DevDependencies: []string{
            "pytest>=6.0.0",
            "pytest-click>=1.1.0",
            "black>=21.0.0",
            "flake8>=3.9.0",
            "mypy>=0.910",
        },
    }
    
    fmt.Println("Creating CLI tool project...")
    err := projectManager.InitProject("./cli-tool", opts)
    if err != nil {
        return err
    }
    
    // Create CLI-specific structure
    cliDirs := []string{
        "./cli-tool/src/cli_tool/commands",
        "./cli-tool/src/cli_tool/utils",
        "./cli-tool/src/cli_tool/config",
    }
    
    for _, dir := range cliDirs {
        err := os.MkdirAll(dir, 0755)
        if err != nil {
            return err
        }
    }
    
    // Create main CLI file
    cliPy := `import click
from rich.console import Console

console = Console()

@click.group()
@click.version_option(version="1.0.0")
def cli():
    """CLI Tool - A command-line interface tool"""
    pass

@cli.command()
@click.argument('name')
def hello(name):
    """Say hello to someone"""
    console.print(f"Hello, {name}!", style="bold green")

@cli.command()
def status():
    """Show status information"""
    console.print("CLI Tool is running!", style="bold blue")

if __name__ == "__main__":
    cli()
`
    
    err = os.WriteFile("./cli-tool/src/cli_tool/main.py", []byte(cliPy), 0644)
    if err != nil {
        return err
    }
    
    fmt.Println("CLI tool project created successfully!")
    return nil
}

func main() {
    err := createCLIProject()
    if err != nil {
        log.Fatal(err)
    }
}
```

## Multiple Projects Setup

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

type ProjectTemplate struct {
    Name         string
    Description  string
    Dependencies []string
    DevDeps      []string
}

func createMultipleProjects() error {
    manager := pip.NewManager(nil)
    projectManager := pip.NewProjectManager(manager)
    
    templates := []ProjectTemplate{
        {
            Name:        "microservice-auth",
            Description: "Authentication microservice",
            Dependencies: []string{
                "fastapi>=0.70.0",
                "uvicorn[standard]>=0.15.0",
                "python-jose[cryptography]>=3.3.0",
                "passlib[bcrypt]>=1.7.4",
                "sqlalchemy>=1.4.0",
            },
            DevDeps: []string{"pytest>=6.0.0", "httpx>=0.24.0"},
        },
        {
            Name:        "microservice-users",
            Description: "User management microservice",
            Dependencies: []string{
                "fastapi>=0.70.0",
                "uvicorn[standard]>=0.15.0",
                "sqlalchemy>=1.4.0",
                "alembic>=1.7.0",
            },
            DevDeps: []string{"pytest>=6.0.0", "httpx>=0.24.0"},
        },
        {
            Name:        "shared-models",
            Description: "Shared data models library",
            Dependencies: []string{
                "pydantic>=1.8.0",
                "sqlalchemy>=1.4.0",
            },
            DevDeps: []string{"pytest>=6.0.0"},
        },
    }
    
    for _, template := range templates {
        fmt.Printf("Creating project: %s\n", template.Name)
        
        opts := &pip.ProjectOptions{
            Name:            template.Name,
            Version:         "1.0.0",
            Description:     template.Description,
            Author:          "Microservices Team",
            License:         "MIT",
            CreateVenv:      true,
            VenvPath:        fmt.Sprintf("./%s/venv", template.Name),
            Dependencies:    template.Dependencies,
            DevDependencies: template.DevDeps,
        }
        
        err := projectManager.InitProject(fmt.Sprintf("./%s", template.Name), opts)
        if err != nil {
            return fmt.Errorf("failed to create %s: %w", template.Name, err)
        }
        
        fmt.Printf("✓ Created %s\n", template.Name)
    }
    
    fmt.Println("All projects created successfully!")
    return nil
}

func main() {
    err := createMultipleProjects()
    if err != nil {
        log.Fatal(err)
    }
}
```

## Project with Custom Templates

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func createProjectFromTemplate(templateName, projectName string) error {
    manager := pip.NewManager(nil)
    projectManager := pip.NewProjectManager(manager)
    
    // Define templates
    templates := map[string]*pip.ProjectOptions{
        "django": {
            Name:        projectName,
            Description: "Django web application",
            Dependencies: []string{
                "django>=4.0.0",
                "djangorestframework>=3.14.0",
                "django-cors-headers>=3.13.0",
                "psycopg2-binary>=2.9.0",
                "celery>=5.2.0",
                "redis>=4.3.0",
            },
            DevDependencies: []string{
                "pytest-django>=4.5.0",
                "black>=21.0.0",
                "flake8>=3.9.0",
            },
        },
        "ml": {
            Name:        projectName,
            Description: "Machine learning project",
            Dependencies: []string{
                "scikit-learn>=1.0.0",
                "tensorflow>=2.8.0",
                "torch>=1.11.0",
                "numpy>=1.21.0",
                "pandas>=1.3.0",
                "matplotlib>=3.4.0",
                "seaborn>=0.11.0",
                "jupyter>=1.0.0",
            },
            DevDependencies: []string{
                "pytest>=6.0.0",
                "black>=21.0.0",
            },
        },
    }
    
    opts, exists := templates[templateName]
    if !exists {
        return fmt.Errorf("template %s not found", templateName)
    }
    
    opts.CreateVenv = true
    opts.VenvPath = filepath.Join(".", projectName, "venv")
    opts.Author = "Developer"
    opts.License = "MIT"
    opts.Version = "1.0.0"
    
    fmt.Printf("Creating %s project from %s template...\n", projectName, templateName)
    err := projectManager.InitProject(fmt.Sprintf("./%s", projectName), opts)
    if err != nil {
        return err
    }
    
    // Add template-specific files
    switch templateName {
    case "django":
        err = createDjangoFiles(projectName)
    case "ml":
        err = createMLFiles(projectName)
    }
    
    if err != nil {
        return err
    }
    
    fmt.Printf("Project %s created successfully from %s template!\n", projectName, templateName)
    return nil
}

func createDjangoFiles(projectName string) error {
    // Create Django-specific structure
    dirs := []string{
        filepath.Join(projectName, "apps"),
        filepath.Join(projectName, "config"),
        filepath.Join(projectName, "static"),
        filepath.Join(projectName, "templates"),
        filepath.Join(projectName, "media"),
    }
    
    for _, dir := range dirs {
        err := os.MkdirAll(dir, 0755)
        if err != nil {
            return err
        }
    }
    
    // Create manage.py placeholder
    managePy := `#!/usr/bin/env python
"""Django's command-line utility for administrative tasks."""
import os
import sys

if __name__ == '__main__':
    os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'config.settings')
    try:
        from django.core.management import execute_from_command_line
    except ImportError as exc:
        raise ImportError(
            "Couldn't import Django. Are you sure it's installed and "
            "available on your PYTHONPATH environment variable? Did you "
            "forget to activate a virtual environment?"
        ) from exc
    execute_from_command_line(sys.argv)
`
    
    return os.WriteFile(filepath.Join(projectName, "manage.py"), []byte(managePy), 0755)
}

func createMLFiles(projectName string) error {
    // Create ML-specific structure
    dirs := []string{
        filepath.Join(projectName, "data", "raw"),
        filepath.Join(projectName, "data", "processed"),
        filepath.Join(projectName, "notebooks"),
        filepath.Join(projectName, "models"),
        filepath.Join(projectName, "src", "features"),
        filepath.Join(projectName, "src", "models"),
        filepath.Join(projectName, "src", "visualization"),
    }
    
    for _, dir := range dirs {
        err := os.MkdirAll(dir, 0755)
        if err != nil {
            return err
        }
    }
    
    // Create example notebook
    notebook := `{
 "cells": [
  {
   "cell_type": "markdown",
   "metadata": {},
   "source": [
    "# ` + projectName + ` - Exploratory Data Analysis\n",
    "\n",
    "This notebook contains the initial data exploration for the project."
   ]
  },
  {
   "cell_type": "code",
   "execution_count": null,
   "metadata": {},
   "outputs": [],
   "source": [
    "import pandas as pd\n",
    "import numpy as np\n",
    "import matplotlib.pyplot as plt\n",
    "import seaborn as sns\n",
    "\n",
    "# Set style\n",
    "plt.style.use('seaborn')\n",
    "sns.set_palette('husl')"
   ]
  }
 ],
 "metadata": {
  "kernelspec": {
   "display_name": "Python 3",
   "language": "python",
   "name": "python3"
  },
  "language_info": {
   "name": "python",
   "version": "3.9.0"
  }
 },
 "nbformat": 4,
 "nbformat_minor": 4
}
`
    
    return os.WriteFile(filepath.Join(projectName, "notebooks", "01-eda.ipynb"), []byte(notebook), 0644)
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: go run main.go <template> <project-name>")
        fmt.Println("Available templates: django, ml")
        os.Exit(1)
    }
    
    templateName := os.Args[1]
    projectName := os.Args[2]
    
    err := createProjectFromTemplate(templateName, projectName)
    if err != nil {
        log.Fatal(err)
    }
}
```

## CLI Usage Examples

```bash
# Initialize basic project
pip-cli project init ./myproject

# The CLI tool creates:
# - Project structure
# - Virtual environment
# - Basic configuration files
# - README.md template
```

## Best Practices

1. **Use descriptive project names**
2. **Always create virtual environments**
3. **Pin dependency versions** for reproducibility
4. **Include development dependencies**
5. **Create comprehensive README files**
6. **Set up testing from the start**
7. **Use version control** (git)

## Next Steps

- Learn about [Advanced Usage](./advanced-usage.md)
- Explore [Package Management](./package-management.md)
- Check the [Project Management API](../api/project-management.md)
