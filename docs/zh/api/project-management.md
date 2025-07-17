# 项目管理 API

项目管理 API 提供了完整的 Python 项目初始化、配置管理和构建功能。

## 核心方法

### InitProject

初始化一个新的 Python 项目。

```go
func (m *Manager) InitProject(path string, opts *ProjectOptions) error
```

**参数：**
- `path` - 项目路径
- `opts` - 项目选项配置

**返回值：**
- `error` - 错误信息

**示例：**

```go
opts := &pip.ProjectOptions{
    Name:        "my-awesome-project",
    Version:     "0.1.0",
    Description: "一个很棒的 Python 项目",
    Author:      "您的姓名",
    AuthorEmail: "your.email@example.com",
    License:     "MIT",
    CreateVenv:  true,
    CreateSrc:   true,
    CreateTests: true,
    Dependencies: []string{
        "requests>=2.25.0",
        "click>=7.0",
    },
}

if err := manager.InitProject("./my-project", opts); err != nil {
    log.Fatalf("项目初始化失败: %v", err)
}
```

### ReadProjectConfig

读取项目配置信息。

```go
func (m *Manager) ReadProjectConfig(path string) (*ProjectConfig, error)
```

**参数：**
- `path` - 项目路径

**返回值：**
- `*ProjectConfig` - 项目配置信息
- `error` - 错误信息

**示例：**

```go
config, err := manager.ReadProjectConfig("./my-project")
if err != nil {
    log.Fatalf("读取项目配置失败: %v", err)
}

fmt.Printf("项目名称: %s\n", config.Name)
fmt.Printf("版本: %s\n", config.Version)
fmt.Printf("依赖数量: %d\n", len(config.Dependencies))
```

### UpdateProjectVersion

更新项目版本。

```go
func (m *Manager) UpdateProjectVersion(path string, version string) error
```

**参数：**
- `path` - 项目路径
- `version` - 新版本号

**返回值：**
- `error` - 错误信息

**示例：**

```go
if err := manager.UpdateProjectVersion("./my-project", "1.1.0"); err != nil {
    log.Fatalf("更新项目版本失败: %v", err)
}
```

### BuildProject

构建项目包。

```go
func (m *Manager) BuildProject(path string, opts *BuildOptions) error
```

**参数：**
- `path` - 项目路径
- `opts` - 构建选项

**返回值：**
- `error` - 错误信息

**示例：**

```go
buildOpts := &pip.BuildOptions{
    OutputDir: "./dist",
    Format:    "wheel",
    Clean:     true,
}

if err := manager.BuildProject("./my-project", buildOpts); err != nil {
    log.Fatalf("构建项目失败: %v", err)
}
```

## 数据类型

### ProjectOptions

项目初始化选项。

```go
type ProjectOptions struct {
    // 基本信息
    Name        string
    Version     string
    Description string
    Author      string
    AuthorEmail string
    License     string
    
    // 项目 URL
    Homepage   string
    Repository string
    
    // Python 版本要求
    PythonRequires string
    
    // 依赖管理
    Dependencies    []string
    DevDependencies []string
    
    // 虚拟环境
    CreateVenv bool
    VenvPath   string
    
    // 项目结构
    CreateSrc   bool
    CreateTests bool
    CreateDocs  bool
    
    // 配置文件
    CreateSetupPy      bool
    CreatePyprojectToml bool
    CreateManifestIn   bool
    CreateGitignore    bool
    CreateReadme       bool
    
    // CI/CD
    CreateGithubActions bool
    CreateDockerfile    bool
    
    // 模板
    Template     *ProjectTemplate
    TemplateName string
}
```

### ProjectConfig

项目配置信息。

```go
type ProjectConfig struct {
    Name         string
    Version      string
    Description  string
    Author       string
    AuthorEmail  string
    License      string
    Homepage     string
    Repository   string
    Dependencies []string
    DevDependencies []string
    PythonRequires  string
    Scripts      map[string]string
    EntryPoints  map[string]map[string]string
}
```

### BuildOptions

构建选项。

```go
type BuildOptions struct {
    OutputDir string   // 输出目录
    Format    string   // 构建格式 (wheel, sdist)
    Clean     bool     // 构建前清理
    Verbose   bool     // 详细输出
    Parallel  bool     // 并行构建
    Optimize  bool     // 优化构建
}
```

### ProjectTemplate

项目模板。

```go
type ProjectTemplate struct {
    Name        string
    Description string
    Version     string
    Files       map[string]string
    Dependencies    []string
    DevDependencies []string
    Variables   map[string]interface{}
    Hooks       map[string]func() error
}
```

## 高级功能

### AddDependency

添加项目依赖。

```go
func (m *Manager) AddDependency(projectPath, packageName, version string) error
```

**示例：**

```go
if err := manager.AddDependency("./my-project", "numpy", ">=1.20.0"); err != nil {
    log.Fatalf("添加依赖失败: %v", err)
}
```

### RemoveDependency

移除项目依赖。

```go
func (m *Manager) RemoveDependency(projectPath, packageName string) error
```

**示例：**

```go
if err := manager.RemoveDependency("./my-project", "old-package"); err != nil {
    log.Fatalf("移除依赖失败: %v", err)
}
```

### ValidateProject

验证项目结构。

```go
func (m *Manager) ValidateProject(path string) ([]*ValidationIssue, error)
```

**示例：**

```go
issues, err := manager.ValidateProject("./my-project")
if err != nil {
    log.Fatalf("验证项目失败: %v", err)
}

if len(issues) > 0 {
    fmt.Println("发现以下问题:")
    for _, issue := range issues {
        fmt.Printf("- %s: %s\n", issue.Type, issue.Message)
    }
}
```

### PublishProject

发布项目到 PyPI。

```go
func (m *Manager) PublishProject(path string, opts *PublishOptions) error
```

**示例：**

```go
publishOpts := &pip.PublishOptions{
    Repository: "pypi",
    Username:   "your-username",
    Password:   "your-password",
}

if err := manager.PublishProject("./my-project", publishOpts); err != nil {
    log.Fatalf("发布项目失败: %v", err)
}
```

## 模板管理

### ListProjectTemplates

列出可用的项目模板。

```go
func (m *Manager) ListProjectTemplates() ([]*ProjectTemplate, error)
```

**示例：**

```go
templates, err := manager.ListProjectTemplates()
if err != nil {
    log.Fatalf("列出模板失败: %v", err)
}

fmt.Println("可用模板:")
for _, template := range templates {
    fmt.Printf("- %s: %s\n", template.Name, template.Description)
}
```

### RegisterTemplate

注册自定义模板。

```go
func (m *Manager) RegisterTemplate(template *ProjectTemplate) error
```

**示例：**

```go
template := &pip.ProjectTemplate{
    Name:        "my-custom-template",
    Description: "我的自定义项目模板",
    Files: map[string]string{
        "{{.Name}}/__init__.py": initTemplate,
        "{{.Name}}/main.py":     mainTemplate,
    },
    Dependencies: []string{
        "click>=7.0",
        "rich>=10.0",
    },
}

if err := manager.RegisterTemplate(template); err != nil {
    log.Fatalf("注册模板失败: %v", err)
}
```

## 项目迁移

### MigrateProject

迁移项目到新格式。

```go
func (m *Manager) MigrateProject(path string, opts *MigrationOptions) error
```

**示例：**

```go
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

### UpgradeProject

升级项目到新版本。

```go
func (m *Manager) UpgradeProject(path string, opts *UpgradeOptions) error
```

**示例：**

```go
upgradeOpts := &pip.UpgradeOptions{
    UpdateDependencies: true,
    UpdateTemplates:    true,
    BackupOriginal:     true,
}

if err := manager.UpgradeProject("./my-project", upgradeOpts); err != nil {
    log.Fatalf("项目升级失败: %v", err)
}
```

## 错误处理

项目管理操作可能遇到的错误类型：

```go
if err := manager.InitProject(path, opts); err != nil {
    switch {
    case pip.IsErrorType(err, pip.ErrorTypeProjectExists):
        fmt.Println("项目已存在")
    case pip.IsErrorType(err, pip.ErrorTypeInvalidProjectName):
        fmt.Println("无效的项目名称")
    case pip.IsErrorType(err, pip.ErrorTypePermissionDenied):
        fmt.Println("权限被拒绝")
    default:
        fmt.Printf("其他错误: %v\n", err)
    }
}
```

## 最佳实践

1. **使用有意义的项目名称**
2. **遵循语义化版本控制**
3. **指定合理的依赖版本范围**
4. **始终创建虚拟环境**
5. **包含完整的项目文档**

## 下一步

- 查看[虚拟环境 API](./virtual-environments.md)
- 了解[错误处理](./errors.md)
- 探索[类型定义](./types.md)
- 学习[日志记录](./logger.md)
