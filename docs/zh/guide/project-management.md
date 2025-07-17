# 项目管理

Go Pip SDK 提供了强大的项目管理功能，帮助您快速初始化 Python 项目，设置标准的项目结构，并管理项目依赖。

## 项目初始化

### 基本项目初始化

```go
package main

import (
    "log"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
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
    
    // 初始化项目
    projectPath := "./my-project"
    if err := manager.InitProject(projectPath, opts); err != nil {
        log.Fatalf("初始化项目失败: %v", err)
    }
    
    fmt.Println("项目初始化成功！")
}
```

### 完整项目配置

```go
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

if err := manager.InitProject("./advanced-project", opts); err != nil {
    log.Fatalf("初始化高级项目失败: %v", err)
}
```

## 项目结构

### 标准项目结构

初始化后的项目将具有以下结构：

```
my-project/
├── src/
│   └── my_project/
│       ├── __init__.py
│       └── main.py
├── tests/
│   ├── __init__.py
│   └── test_main.py
├── docs/
│   ├── index.md
│   └── api.md
├── venv/                    # 虚拟环境（如果启用）
├── .github/
│   └── workflows/
│       └── ci.yml          # GitHub Actions（如果启用）
├── setup.py                # 包安装脚本
├── pyproject.toml          # 现代 Python 项目配置
├── requirements.txt        # 生产依赖
├── dev-requirements.txt    # 开发依赖
├── MANIFEST.in            # 包含文件清单
├── .gitignore             # Git 忽略文件
├── README.md              # 项目说明
├── LICENSE                # 许可证文件
└── Dockerfile             # Docker 配置（如果启用）
```

### 自定义项目模板

```go
// 创建自定义项目模板
template := &pip.ProjectTemplate{
    Name: "web-api-template",
    Files: map[string]string{
        "app/__init__.py":     "",
        "app/main.py":         webApiMainTemplate,
        "app/models.py":       modelsTemplate,
        "app/routes.py":       routesTemplate,
        "config.py":           configTemplate,
        "docker-compose.yml":  dockerComposeTemplate,
    },
    Dependencies: []string{
        "fastapi>=0.68.0",
        "uvicorn>=0.15.0",
        "sqlalchemy>=1.4.0",
        "alembic>=1.7.0",
    },
}

// 使用自定义模板初始化项目
opts := &pip.ProjectOptions{
    Template: template,
    Name:     "my-web-api",
    Version:  "0.1.0",
}

if err := manager.InitProjectFromTemplate("./my-api", opts); err != nil {
    log.Fatalf("从模板初始化项目失败: %v", err)
}
```

## 依赖管理

### 添加依赖

```go
// 添加生产依赖
if err := manager.AddDependency("./my-project", "numpy", ">=1.20.0"); err != nil {
    log.Fatalf("添加依赖失败: %v", err)
}

// 添加开发依赖
if err := manager.AddDevDependency("./my-project", "pytest-cov", ">=2.12.0"); err != nil {
    log.Fatalf("添加开发依赖失败: %v", err)
}
```

### 移除依赖

```go
// 移除依赖
if err := manager.RemoveDependency("./my-project", "old-package"); err != nil {
    log.Fatalf("移除依赖失败: %v", err)
}
```

### 更新依赖

```go
// 更新所有依赖到最新版本
if err := manager.UpdateDependencies("./my-project"); err != nil {
    log.Fatalf("更新依赖失败: %v", err)
}

// 更新特定依赖
if err := manager.UpdateDependency("./my-project", "requests"); err != nil {
    log.Fatalf("更新 requests 失败: %v", err)
}
```

### 锁定依赖版本

```go
// 生成锁定文件（类似 poetry.lock 或 Pipfile.lock）
if err := manager.LockDependencies("./my-project"); err != nil {
    log.Fatalf("锁定依赖失败: %v", err)
}

// 从锁定文件安装精确版本
if err := manager.InstallFromLock("./my-project"); err != nil {
    log.Fatalf("从锁定文件安装失败: %v", err)
}
```

## 项目配置管理

### 读取项目配置

```go
// 读取项目配置
config, err := manager.ReadProjectConfig("./my-project")
if err != nil {
    log.Fatalf("读取项目配置失败: %v", err)
}

fmt.Printf("项目名称: %s\n", config.Name)
fmt.Printf("版本: %s\n", config.Version)
fmt.Printf("依赖数量: %d\n", len(config.Dependencies))
```

### 更新项目配置

```go
// 更新项目版本
if err := manager.UpdateProjectVersion("./my-project", "1.1.0"); err != nil {
    log.Fatalf("更新项目版本失败: %v", err)
}

// 更新项目描述
if err := manager.UpdateProjectDescription("./my-project", "更新的项目描述"); err != nil {
    log.Fatalf("更新项目描述失败: %v", err)
}
```

## 项目构建和打包

### 构建项目

```go
// 构建项目包
buildOpts := &pip.BuildOptions{
    OutputDir: "./dist",
    Format:    "wheel", // 或 "sdist"
    Clean:     true,
}

if err := manager.BuildProject("./my-project", buildOpts); err != nil {
    log.Fatalf("构建项目失败: %v", err)
}
```

### 发布项目

```go
// 发布到 PyPI
publishOpts := &pip.PublishOptions{
    Repository: "pypi", // 或 "testpypi"
    Username:   "your-username",
    Password:   "your-password", // 建议使用环境变量
}

if err := manager.PublishProject("./my-project", publishOpts); err != nil {
    log.Fatalf("发布项目失败: %v", err)
}
```

## 项目验证和测试

### 验证项目结构

```go
// 验证项目结构是否正确
issues, err := manager.ValidateProject("./my-project")
if err != nil {
    log.Fatalf("验证项目失败: %v", err)
}

if len(issues) > 0 {
    fmt.Println("发现以下问题:")
    for _, issue := range issues {
        fmt.Printf("- %s: %s\n", issue.Type, issue.Message)
    }
} else {
    fmt.Println("项目结构验证通过！")
}
```

### 运行项目测试

```go
// 运行项目测试
testOpts := &pip.TestOptions{
    Coverage:    true,
    Verbose:     true,
    FailFast:    false,
    TestPattern: "test_*.py",
}

result, err := manager.RunTests("./my-project", testOpts)
if err != nil {
    log.Fatalf("运行测试失败: %v", err)
}

fmt.Printf("测试结果: %d 通过, %d 失败\n", result.Passed, result.Failed)
fmt.Printf("代码覆盖率: %.2f%%\n", result.Coverage)
```

## 项目模板

### 内置模板

SDK 提供了多种内置项目模板：

```go
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
}

if err := manager.InitProject("./my-api", opts); err != nil {
    log.Fatalf("从模板初始化失败: %v", err)
}
```

### 创建自定义模板

```go
// 创建自定义模板
template := &pip.ProjectTemplate{
    Name:        "my-custom-template",
    Description: "我的自定义项目模板",
    Version:     "1.0.0",
    
    // 模板文件
    Files: map[string]string{
        "{{.Name}}/__init__.py": initTemplate,
        "{{.Name}}/main.py":     mainTemplate,
        "tests/test_{{.Name}}.py": testTemplate,
        "README.md":             readmeTemplate,
    },
    
    // 默认依赖
    Dependencies: []string{
        "click>=7.0",
        "rich>=10.0",
    },
    DevDependencies: []string{
        "pytest>=6.0",
        "black>=21.0",
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
```

## 项目迁移

### 从现有项目迁移

```go
// 从 requirements.txt 迁移到现代项目结构
migrationOpts := &pip.MigrationOptions{
    SourceFormat: "requirements",
    TargetFormat: "pyproject",
    CreateVenv:   true,
    UpdateGitignore: true,
}

if err := manager.MigrateProject("./legacy-project", migrationOpts); err != nil {
    log.Fatalf("项目迁移失败: %v", err)
}
```

### 项目升级

```go
// 升级项目到新的 SDK 版本
upgradeOpts := &pip.UpgradeOptions{
    UpdateDependencies: true,
    UpdateTemplates:    true,
    BackupOriginal:     true,
}

if err := manager.UpgradeProject("./my-project", upgradeOpts); err != nil {
    log.Fatalf("项目升级失败: %v", err)
}
```

## 最佳实践

### 1. 项目命名

```go
// 使用有意义的项目名称
opts := &pip.ProjectOptions{
    Name: "user-authentication-service", // 好
    // Name: "project1",                 // 不好
}
```

### 2. 版本管理

```go
// 遵循语义化版本
opts := &pip.ProjectOptions{
    Version: "1.2.3", // 主版本.次版本.修订版本
}
```

### 3. 依赖管理

```go
// 指定版本范围而不是精确版本
opts := &pip.ProjectOptions{
    Dependencies: []string{
        "requests>=2.25.0,<3.0.0", // 好
        // "requests==2.25.1",      // 过于严格
    },
}
```

### 4. 项目文档

```go
// 始终包含完整的项目文档
opts := &pip.ProjectOptions{
    CreateReadme:    true,
    CreateDocs:      true,
    CreateChangelog: true,
}
```

## 故障排除

### 常见问题

1. **项目初始化失败**
   - 检查目录权限
   - 确保目标目录不存在或为空
   - 验证项目名称有效性

2. **依赖安装失败**
   - 检查网络连接
   - 验证包名称和版本
   - 检查 Python 版本兼容性

3. **构建失败**
   - 确保所有必需文件存在
   - 检查 setup.py 或 pyproject.toml 配置
   - 验证依赖完整性

### 调试技巧

```go
// 启用详细日志
config := &pip.Config{
    LogLevel: "DEBUG",
}
manager := pip.NewManager(config)

// 验证项目状态
status, err := manager.GetProjectStatus("./my-project")
if err != nil {
    log.Printf("获取项目状态失败: %v", err)
} else {
    fmt.Printf("项目状态: %+v\n", status)
}
```

## 下一步

- 学习[错误处理](./error-handling.md)
- 了解[日志记录](./logging.md)
- 查看[项目管理示例](/zh/examples/project-initialization.md)
- 探索[API 参考](/zh/api/project-management.md)
