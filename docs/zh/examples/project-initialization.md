# 项目初始化示例

本页面提供了详细的项目初始化示例，展示如何使用 Go Pip SDK 创建和管理 Python 项目。

## 基本项目初始化

### 简单项目创建

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func createBasicProject() {
    manager := pip.NewManager(nil)
    
    // 基本项目选项
    opts := &pip.ProjectOptions{
        Name:        "my-awesome-project",
        Version:     "0.1.0",
        Description: "一个很棒的 Python 项目",
        Author:      "您的姓名",
        AuthorEmail: "your.email@example.com",
        License:     "MIT",
    }
    
    projectPath := "./my-project"
    fmt.Printf("初始化项目: %s\n", opts.Name)
    
    if err := manager.InitProject(projectPath, opts); err != nil {
        log.Fatalf("项目初始化失败: %v", err)
    }
    
    fmt.Println("✅ 项目初始化成功！")
    
    // 读取项目配置验证
    config, err := manager.ReadProjectConfig(projectPath)
    if err != nil {
        log.Printf("读取项目配置失败: %v", err)
        return
    }
    
    fmt.Printf("项目名称: %s\n", config.Name)
    fmt.Printf("项目版本: %s\n", config.Version)
    fmt.Printf("项目描述: %s\n", config.Description)
}
```

### 完整项目配置

```go
func createAdvancedProject() {
    manager := pip.NewManager(nil)
    
    // 完整的项目配置
    opts := &pip.ProjectOptions{
        // 基本信息
        Name:        "advanced-python-project",
        Version:     "1.0.0",
        Description: "一个高级的 Python 项目示例",
        Author:      "开发团队",
        AuthorEmail: "team@example.com",
        License:     "MIT",
        
        // 项目 URL
        Homepage:    "https://github.com/username/project",
        Repository:  "https://github.com/username/project.git",
        
        // Python 版本要求
        PythonRequires: ">=3.8",
        
        // 依赖管理
        Dependencies: []string{
            "requests>=2.25.0",
            "click>=7.0",
            "pydantic>=1.8.0",
        },
        DevDependencies: []string{
            "pytest>=6.0",
            "black>=21.0",
            "flake8>=3.8",
            "mypy>=0.812",
        },
        
        // 虚拟环境
        CreateVenv: true,
        VenvPath:   "./venv",
        
        // 项目结构
        CreateSrc:   true,
        CreateTests: true,
        CreateDocs:  true,
        
        // 配置文件
        CreateSetupPy:      true,
        CreatePyprojectToml: true,
        CreateManifestIn:   true,
        CreateGitignore:    true,
        CreateReadme:       true,
        
        // CI/CD
        CreateGithubActions: true,
        CreateDockerfile:    true,
    }
    
    projectPath := "./advanced-project"
    fmt.Println("创建高级项目...")
    
    if err := manager.InitProject(projectPath, opts); err != nil {
        log.Fatalf("初始化高级项目失败: %v", err)
    }
    
    fmt.Println("✅ 高级项目创建成功")
    
    // 显示项目结构
    fmt.Println("\n项目结构:")
    showProjectStructure(projectPath)
}

func showProjectStructure(projectPath string) {
    // 这里可以实现显示项目目录结构的逻辑
    fmt.Printf("📁 %s/\n", projectPath)
    fmt.Println("├── src/")
    fmt.Println("│   └── advanced_project/")
    fmt.Println("│       ├── __init__.py")
    fmt.Println("│       └── main.py")
    fmt.Println("├── tests/")
    fmt.Println("│   ├── __init__.py")
    fmt.Println("│   └── test_main.py")
    fmt.Println("├── docs/")
    fmt.Println("├── venv/")
    fmt.Println("├── .github/workflows/")
    fmt.Println("├── setup.py")
    fmt.Println("├── pyproject.toml")
    fmt.Println("├── requirements.txt")
    fmt.Println("├── dev-requirements.txt")
    fmt.Println("├── .gitignore")
    fmt.Println("├── README.md")
    fmt.Println("└── Dockerfile")
}
```

## 项目模板

### 使用内置模板

```go
func createProjectFromTemplate() {
    manager := pip.NewManager(nil)
    
    // 列出可用模板
    templates, err := manager.ListProjectTemplates()
    if err != nil {
        log.Fatalf("列出模板失败: %v", err)
    }
    
    fmt.Println("可用模板:")
    for _, template := range templates {
        fmt.Printf("- %s: %s\n", template.Name, template.Description)
    }
    
    // 使用特定模板
    opts := &pip.ProjectOptions{
        TemplateName: "fastapi-web-api",
        Name:         "my-api-project",
        Version:      "0.1.0",
        Author:       "API 开发者",
        AuthorEmail:  "api@example.com",
    }
    
    if err := manager.InitProject("./my-api", opts); err != nil {
        log.Fatalf("从模板初始化失败: %v", err)
    }
    
    fmt.Println("✅ 从模板创建项目成功")
}
```

### 创建自定义模板

```go
func createCustomTemplate() {
    manager := pip.NewManager(nil)
    
    // 定义自定义模板
    template := &pip.ProjectTemplate{
        Name:        "my-custom-template",
        Description: "我的自定义项目模板",
        Version:     "1.0.0",
        
        // 模板文件内容
        Files: map[string]string{
            "{{.Name}}/__init__.py": `"""
{{.Description}}
"""
__version__ = "{{.Version}}"
__author__ = "{{.Author}}"
`,
            "{{.Name}}/main.py": `#!/usr/bin/env python3
"""
{{.Name}} 主模块
"""

import click

@click.command()
@click.option('--name', default='World', help='要问候的名字')
def hello(name):
    """简单的问候程序"""
    click.echo(f'Hello {name}!')

if __name__ == '__main__':
    hello()
`,
            "tests/test_{{.Name}}.py": `import pytest
from {{.Name}}.main import hello

def test_hello():
    """测试 hello 函数"""
    # 这里添加测试代码
    pass
`,
            "README.md": `# {{.Name}}

{{.Description}}

## 安装

\`\`\`bash
pip install -e .
\`\`\`

## 使用

\`\`\`bash
python -m {{.Name}}
\`\`\`

## 开发

\`\`\`bash
pip install -e ".[dev]"
pytest
\`\`\`
`,
        },
        
        // 默认依赖
        Dependencies: []string{
            "click>=7.0",
            "rich>=10.0",
        },
        DevDependencies: []string{
            "pytest>=6.0",
            "black>=21.0",
            "flake8>=3.8",
        },
        
        // 模板变量
        Variables: map[string]interface{}{
            "DefaultPort": 8000,
            "UseAsync":    true,
        },
    }
    
    // 注册模板
    if err := manager.RegisterTemplate(template); err != nil {
        log.Fatalf("注册模板失败: %v", err)
    }
    
    fmt.Println("✅ 自定义模板注册成功")
    
    // 使用自定义模板创建项目
    opts := &pip.ProjectOptions{
        TemplateName: "my-custom-template",
        Name:         "test-project",
        Version:      "0.1.0",
        Description:  "使用自定义模板的测试项目",
        Author:       "模板用户",
        AuthorEmail:  "user@example.com",
    }
    
    if err := manager.InitProject("./test-project", opts); err != nil {
        log.Fatalf("使用自定义模板创建项目失败: %v", err)
    }
    
    fmt.Println("✅ 使用自定义模板创建项目成功")
}
```

## 特定类型的项目

### Web API 项目

```go
func createWebAPIProject() {
    manager := pip.NewManager(nil)
    
    opts := &pip.ProjectOptions{
        Name:        "my-web-api",
        Version:     "0.1.0",
        Description: "FastAPI Web API 项目",
        Author:      "API 开发者",
        AuthorEmail: "api@example.com",
        License:     "MIT",
        
        // Web API 特定依赖
        Dependencies: []string{
            "fastapi>=0.68.0",
            "uvicorn[standard]>=0.15.0",
            "pydantic>=1.8.0",
            "sqlalchemy>=1.4.0",
            "alembic>=1.7.0",
            "python-multipart>=0.0.5",
            "python-jose[cryptography]>=3.3.0",
            "passlib[bcrypt]>=1.7.4",
        },
        DevDependencies: []string{
            "pytest>=6.0",
            "pytest-asyncio>=0.15.0",
            "httpx>=0.24.0",
            "black>=21.0",
            "flake8>=3.8",
            "mypy>=0.812",
        },
        
        CreateVenv:          true,
        CreateSrc:           true,
        CreateTests:         true,
        CreateDocs:          true,
        CreateGithubActions: true,
        CreateDockerfile:    true,
    }
    
    if err := manager.InitProject("./my-web-api", opts); err != nil {
        log.Fatalf("创建 Web API 项目失败: %v", err)
    }
    
    fmt.Println("✅ Web API 项目创建成功")
}
```

### 数据科学项目

```go
func createDataScienceProject() {
    manager := pip.NewManager(nil)
    
    opts := &pip.ProjectOptions{
        Name:        "data-analysis-project",
        Version:     "0.1.0",
        Description: "数据分析项目",
        Author:      "数据科学家",
        AuthorEmail: "datascientist@example.com",
        License:     "MIT",
        
        // 数据科学依赖
        Dependencies: []string{
            "numpy>=1.20.0",
            "pandas>=1.3.0",
            "matplotlib>=3.4.0",
            "seaborn>=0.11.0",
            "scikit-learn>=1.0.0",
            "jupyter>=1.0.0",
            "ipython>=7.0.0",
            "plotly>=5.0.0",
            "scipy>=1.7.0",
        },
        DevDependencies: []string{
            "pytest>=6.0",
            "black>=21.0",
            "flake8>=3.8",
            "mypy>=0.812",
            "nbqa>=1.1.0",
        },
        
        CreateVenv:  true,
        CreateSrc:   true,
        CreateTests: true,
        CreateDocs:  true,
    }
    
    if err := manager.InitProject("./data-analysis", opts); err != nil {
        log.Fatalf("创建数据科学项目失败: %v", err)
    }
    
    fmt.Println("✅ 数据科学项目创建成功")
}
```

### CLI 工具项目

```go
func createCLIProject() {
    manager := pip.NewManager(nil)
    
    opts := &pip.ProjectOptions{
        Name:        "my-cli-tool",
        Version:     "0.1.0",
        Description: "命令行工具项目",
        Author:      "CLI 开发者",
        AuthorEmail: "cli@example.com",
        License:     "MIT",
        
        // CLI 工具依赖
        Dependencies: []string{
            "click>=7.0",
            "rich>=10.0",
            "typer>=0.4.0",
            "pydantic>=1.8.0",
            "toml>=0.10.0",
        },
        DevDependencies: []string{
            "pytest>=6.0",
            "black>=21.0",
            "flake8>=3.8",
            "mypy>=0.812",
            "click-testing>=0.1.0",
        },
        
        CreateVenv:          true,
        CreateSrc:           true,
        CreateTests:         true,
        CreateDocs:          true,
        CreateGithubActions: true,
    }
    
    if err := manager.InitProject("./my-cli-tool", opts); err != nil {
        log.Fatalf("创建 CLI 项目失败: %v", err)
    }
    
    fmt.Println("✅ CLI 工具项目创建成功")
}
```

## 项目配置管理

### 读取和更新项目配置

```go
func manageProjectConfig() {
    manager := pip.NewManager(nil)
    projectPath := "./my-project"
    
    // 读取项目配置
    config, err := manager.ReadProjectConfig(projectPath)
    if err != nil {
        log.Fatalf("读取项目配置失败: %v", err)
    }
    
    fmt.Printf("当前项目配置:\n")
    fmt.Printf("  名称: %s\n", config.Name)
    fmt.Printf("  版本: %s\n", config.Version)
    fmt.Printf("  描述: %s\n", config.Description)
    fmt.Printf("  依赖数量: %d\n", len(config.Dependencies))
    
    // 更新项目版本
    newVersion := "1.1.0"
    fmt.Printf("更新项目版本到 %s...\n", newVersion)
    if err := manager.UpdateProjectVersion(projectPath, newVersion); err != nil {
        log.Fatalf("更新项目版本失败: %v", err)
    }
    fmt.Println("✅ 项目版本更新成功")
    
    // 添加新依赖
    fmt.Println("添加新依赖...")
    if err := manager.AddDependency(projectPath, "requests", ">=2.25.0"); err != nil {
        log.Printf("添加依赖失败: %v", err)
    } else {
        fmt.Println("✅ 依赖添加成功")
    }
    
    // 移除依赖
    fmt.Println("移除旧依赖...")
    if err := manager.RemoveDependency(projectPath, "old-package"); err != nil {
        log.Printf("移除依赖失败: %v", err)
    } else {
        fmt.Println("✅ 依赖移除成功")
    }
}
```

### 项目验证

```go
func validateProject() {
    manager := pip.NewManager(nil)
    projectPath := "./my-project"
    
    fmt.Println("验证项目结构...")
    issues, err := manager.ValidateProject(projectPath)
    if err != nil {
        log.Fatalf("验证项目失败: %v", err)
    }
    
    if len(issues) == 0 {
        fmt.Println("✅ 项目结构验证通过")
    } else {
        fmt.Printf("发现 %d 个问题:\n", len(issues))
        for _, issue := range issues {
            fmt.Printf("- %s: %s\n", issue.Type, issue.Message)
            if issue.Suggestion != "" {
                fmt.Printf("  建议: %s\n", issue.Suggestion)
            }
        }
    }
}
```

## 项目构建和发布

### 构建项目

```go
func buildProject() {
    manager := pip.NewManager(nil)
    projectPath := "./my-project"
    
    // 构建配置
    buildOpts := &pip.BuildOptions{
        OutputDir: "./dist",
        Format:    "wheel", // 或 "sdist"
        Clean:     true,
        Verbose:   true,
    }
    
    fmt.Println("构建项目...")
    if err := manager.BuildProject(projectPath, buildOpts); err != nil {
        log.Fatalf("构建项目失败: %v", err)
    }
    
    fmt.Println("✅ 项目构建成功")
    fmt.Println("构建文件位于: ./dist/")
}
```

### 发布项目

```go
func publishProject() {
    manager := pip.NewManager(nil)
    projectPath := "./my-project"
    
    // 发布到测试 PyPI
    publishOpts := &pip.PublishOptions{
        Repository: "testpypi",
        Username:   "your-username",
        Password:   "your-password", // 建议使用环境变量
    }
    
    fmt.Println("发布到测试 PyPI...")
    if err := manager.PublishProject(projectPath, publishOpts); err != nil {
        log.Fatalf("发布项目失败: %v", err)
    }
    
    fmt.Println("✅ 项目发布成功")
}
```

## 项目迁移和升级

### 迁移现有项目

```go
func migrateExistingProject() {
    manager := pip.NewManager(nil)
    
    // 从 requirements.txt 迁移到现代项目结构
    migrationOpts := &pip.MigrationOptions{
        SourceFormat:    "requirements",
        TargetFormat:    "pyproject",
        CreateVenv:      true,
        UpdateGitignore: true,
        BackupOriginal:  true,
    }
    
    fmt.Println("迁移现有项目...")
    if err := manager.MigrateProject("./legacy-project", migrationOpts); err != nil {
        log.Fatalf("项目迁移失败: %v", err)
    }
    
    fmt.Println("✅ 项目迁移成功")
}
```

### 升级项目

```go
func upgradeProject() {
    manager := pip.NewManager(nil)
    
    upgradeOpts := &pip.UpgradeOptions{
        UpdateDependencies: true,
        UpdateTemplates:    true,
        BackupOriginal:     true,
        Force:              false,
    }
    
    fmt.Println("升级项目...")
    if err := manager.UpgradeProject("./my-project", upgradeOpts); err != nil {
        log.Fatalf("项目升级失败: %v", err)
    }
    
    fmt.Println("✅ 项目升级成功")
}
```

## 批量项目操作

### 批量创建项目

```go
func batchCreateProjects() {
    manager := pip.NewManager(nil)
    
    projects := []struct {
        name        string
        projectType string
        description string
    }{
        {"web-api-service", "fastapi", "Web API 服务"},
        {"data-processor", "data-science", "数据处理工具"},
        {"cli-utility", "cli", "命令行工具"},
    }
    
    for _, proj := range projects {
        fmt.Printf("创建项目: %s (%s)\n", proj.name, proj.description)
        
        opts := &pip.ProjectOptions{
            Name:        proj.name,
            Version:     "0.1.0",
            Description: proj.description,
            Author:      "批量创建",
            AuthorEmail: "batch@example.com",
            License:     "MIT",
            CreateVenv:  true,
            CreateSrc:   true,
            CreateTests: true,
        }
        
        // 根据项目类型设置特定依赖
        switch proj.projectType {
        case "fastapi":
            opts.Dependencies = []string{"fastapi>=0.68.0", "uvicorn>=0.15.0"}
        case "data-science":
            opts.Dependencies = []string{"numpy>=1.20.0", "pandas>=1.3.0"}
        case "cli":
            opts.Dependencies = []string{"click>=7.0", "rich>=10.0"}
        }
        
        if err := manager.InitProject("./"+proj.name, opts); err != nil {
            fmt.Printf("❌ 创建 %s 失败: %v\n", proj.name, err)
        } else {
            fmt.Printf("✅ %s 创建成功\n", proj.name)
        }
    }
}
```

## 完整示例：项目管理工具

```go
package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"
    
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    manager := pip.NewManager(nil)
    
    for {
        fmt.Println("\n=== Python 项目管理工具 ===")
        fmt.Println("1. 创建新项目")
        fmt.Println("2. 从模板创建项目")
        fmt.Println("3. 验证项目")
        fmt.Println("4. 构建项目")
        fmt.Println("5. 迁移项目")
        fmt.Println("6. 升级项目")
        fmt.Println("7. 退出")
        fmt.Print("请选择操作 (1-7): ")
        
        reader := bufio.NewReader(os.Stdin)
        choice, _ := reader.ReadString('\n')
        choice = strings.TrimSpace(choice)
        
        switch choice {
        case "1":
            createProjectInteractive(manager, reader)
        case "2":
            createFromTemplateInteractive(manager, reader)
        case "3":
            validateProjectInteractive(manager, reader)
        case "4":
            buildProjectInteractive(manager, reader)
        case "5":
            migrateProjectInteractive(manager, reader)
        case "6":
            upgradeProjectInteractive(manager, reader)
        case "7":
            fmt.Println("再见！")
            return
        default:
            fmt.Println("无效选择，请重试")
        }
    }
}

func createProjectInteractive(manager *pip.Manager, reader *bufio.Reader) {
    fmt.Print("项目名称: ")
    name, _ := reader.ReadString('\n')
    name = strings.TrimSpace(name)
    
    fmt.Print("项目描述: ")
    description, _ := reader.ReadString('\n')
    description = strings.TrimSpace(description)
    
    fmt.Print("作者姓名: ")
    author, _ := reader.ReadString('\n')
    author = strings.TrimSpace(author)
    
    fmt.Print("作者邮箱: ")
    email, _ := reader.ReadString('\n')
    email = strings.TrimSpace(email)
    
    opts := &pip.ProjectOptions{
        Name:        name,
        Version:     "0.1.0",
        Description: description,
        Author:      author,
        AuthorEmail: email,
        License:     "MIT",
        CreateVenv:  true,
        CreateSrc:   true,
        CreateTests: true,
        CreateDocs:  true,
    }
    
    projectPath := "./" + name
    if err := manager.InitProject(projectPath, opts); err != nil {
        fmt.Printf("❌ 项目创建失败: %v\n", err)
    } else {
        fmt.Printf("✅ 项目 '%s' 创建成功\n", name)
    }
}

// ... 其他交互函数的实现
```

## 下一步

- 学习[高级用法示例](./advanced-usage.md)
- 查看[包管理示例](./package-management.md)
- 探索[虚拟环境示例](./virtual-environments.md)
- 了解[基本用法](./basic-usage.md)
