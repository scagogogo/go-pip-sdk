# 包操作 API

包操作 API 提供了完整的 Python 包管理功能，包括安装、卸载、搜索、列出和更新包。

## 核心方法

### InstallPackage

安装指定的 Python 包。

```go
func (m *Manager) InstallPackage(pkg *PackageSpec) error
```

**参数：**
- `pkg` - 包规范，定义要安装的包及其选项

**返回值：**
- `error` - 错误信息

**示例：**

```go
// 基本安装
pkg := &pip.PackageSpec{
    Name: "requests",
}
if err := manager.InstallPackage(pkg); err != nil {
    log.Fatalf("安装失败: %v", err)
}

// 指定版本安装
pkg := &pip.PackageSpec{
    Name:    "django",
    Version: "4.2.0",
}

// 使用版本约束
pkg := &pip.PackageSpec{
    Name:    "numpy",
    Version: ">=1.20.0,<2.0.0",
}

// 可编辑安装
pkg := &pip.PackageSpec{
    Name:     "./my-local-package",
    Editable: true,
}

// 强制重新安装
pkg := &pip.PackageSpec{
    Name:           "problematic-package",
    ForceReinstall: true,
}
```

### UninstallPackage

卸载指定的 Python 包。

```go
func (m *Manager) UninstallPackage(name string) error
```

**参数：**
- `name` - 要卸载的包名

**返回值：**
- `error` - 错误信息

**示例：**

```go
if err := manager.UninstallPackage("requests"); err != nil {
    log.Fatalf("卸载失败: %v", err)
}
```

### ListPackages

列出所有已安装的包。

```go
func (m *Manager) ListPackages() ([]*Package, error)
```

**返回值：**
- `[]*Package` - 已安装包的列表
- `error` - 错误信息

**示例：**

```go
packages, err := manager.ListPackages()
if err != nil {
    log.Fatalf("列出包失败: %v", err)
}

fmt.Printf("找到 %d 个已安装的包:\n", len(packages))
for _, pkg := range packages {
    fmt.Printf("- %s %s\n", pkg.Name, pkg.Version)
}
```

### ShowPackage

显示指定包的详细信息。

```go
func (m *Manager) ShowPackage(name string) (*PackageInfo, error)
```

**参数：**
- `name` - 包名

**返回值：**
- `*PackageInfo` - 包的详细信息
- `error` - 错误信息

**示例：**

```go
info, err := manager.ShowPackage("requests")
if err != nil {
    log.Fatalf("获取包信息失败: %v", err)
}

fmt.Printf("包名: %s\n", info.Name)
fmt.Printf("版本: %s\n", info.Version)
fmt.Printf("摘要: %s\n", info.Summary)
fmt.Printf("作者: %s\n", info.Author)
fmt.Printf("许可证: %s\n", info.License)
fmt.Printf("依赖: %v\n", info.Requires)
```

### SearchPackages

搜索 PyPI 上的包。

```go
func (m *Manager) SearchPackages(query string) ([]*SearchResult, error)
```

**参数：**
- `query` - 搜索查询字符串

**返回值：**
- `[]*SearchResult` - 搜索结果列表
- `error` - 错误信息

**示例：**

```go
results, err := manager.SearchPackages("web framework")
if err != nil {
    log.Fatalf("搜索失败: %v", err)
}

fmt.Printf("找到 %d 个搜索结果:\n", len(results))
for _, result := range results {
    fmt.Printf("- %s: %s\n", result.Name, result.Summary)
}
```

### FreezePackages

获取当前环境的包列表（类似 `pip freeze`）。

```go
func (m *Manager) FreezePackages() ([]*Package, error)
```

**返回值：**
- `[]*Package` - 包列表，包含精确版本
- `error` - 错误信息

**示例：**

```go
packages, err := manager.FreezePackages()
if err != nil {
    log.Fatalf("冻结包列表失败: %v", err)
}

fmt.Println("当前环境包列表:")
for _, pkg := range packages {
    fmt.Printf("%s==%s\n", pkg.Name, pkg.Version)
}
```

## 需求文件操作

### InstallRequirements

从需求文件安装包。

```go
func (m *Manager) InstallRequirements(path string) error
```

**参数：**
- `path` - 需求文件路径

**返回值：**
- `error` - 错误信息

**示例：**

```go
// 从 requirements.txt 安装
if err := manager.InstallRequirements("requirements.txt"); err != nil {
    log.Fatalf("从需求文件安装失败: %v", err)
}

// 从自定义需求文件安装
if err := manager.InstallRequirements("dev-requirements.txt"); err != nil {
    log.Fatalf("安装开发依赖失败: %v", err)
}
```

### GenerateRequirements

生成需求文件。

```go
func (m *Manager) GenerateRequirements(path string) error
```

**参数：**
- `path` - 输出文件路径

**返回值：**
- `error` - 错误信息

**示例：**

```go
// 生成 requirements.txt
if err := manager.GenerateRequirements("requirements.txt"); err != nil {
    log.Fatalf("生成需求文件失败: %v", err)
}
```

## 高级包操作

### UpdatePackage

更新指定包到最新版本。

```go
func (m *Manager) UpdatePackage(name string) error
```

**参数：**
- `name` - 要更新的包名

**返回值：**
- `error` - 错误信息

**示例：**

```go
if err := manager.UpdatePackage("requests"); err != nil {
    log.Fatalf("更新包失败: %v", err)
}
```

### UpdateAllPackages

更新所有已安装的包。

```go
func (m *Manager) UpdateAllPackages() error
```

**返回值：**
- `error` - 错误信息

**示例：**

```go
if err := manager.UpdateAllPackages(); err != nil {
    log.Fatalf("更新所有包失败: %v", err)
}
```

### CheckOutdated

检查过时的包。

```go
func (m *Manager) CheckOutdated() ([]*OutdatedPackage, error)
```

**返回值：**
- `[]*OutdatedPackage` - 过时包的列表
- `error` - 错误信息

**示例：**

```go
outdated, err := manager.CheckOutdated()
if err != nil {
    log.Fatalf("检查过时包失败: %v", err)
}

if len(outdated) > 0 {
    fmt.Println("发现过时的包:")
    for _, pkg := range outdated {
        fmt.Printf("- %s: %s -> %s\n", pkg.Name, pkg.CurrentVersion, pkg.LatestVersion)
    }
}
```

### DownloadPackage

下载包但不安装。

```go
func (m *Manager) DownloadPackage(pkg *PackageSpec, destDir string) error
```

**参数：**
- `pkg` - 包规范
- `destDir` - 下载目标目录

**返回值：**
- `error` - 错误信息

**示例：**

```go
pkg := &pip.PackageSpec{
    Name:    "requests",
    Version: "2.25.1",
}

if err := manager.DownloadPackage(pkg, "./downloads"); err != nil {
    log.Fatalf("下载包失败: %v", err)
}
```

## 数据类型

### PackageSpec

包规范结构，用于定义要安装的包。

```go
type PackageSpec struct {
    Name           string            // 包名
    Version        string            // 版本约束
    Extras         []string          // 额外依赖
    Index          string            // 自定义索引 URL
    Options        map[string]string // 额外选项
    Editable       bool              // 可编辑安装
    Upgrade        bool              // 升级包
    ForceReinstall bool              // 强制重新安装
    NoDeps         bool              // 不安装依赖
    UserInstall    bool              // 用户级安装
}
```

**示例：**

```go
// 完整的包规范
pkg := &pip.PackageSpec{
    Name:    "fastapi",
    Version: ">=0.68.0",
    Extras:  []string{"dev", "test"},
    Index:   "https://pypi.org/simple/",
    Options: map[string]string{
        "timeout": "120",
        "retries": "5",
    },
    Upgrade: true,
}
```

### Package

已安装包的信息。

```go
type Package struct {
    Name      string // 包名
    Version   string // 版本
    Location  string // 安装位置
    Editable  bool   // 是否为可编辑安装
    Installer string // 安装器
}
```

### PackageInfo

包的详细信息。

```go
type PackageInfo struct {
    Name         string   // 包名
    Version      string   // 版本
    Summary      string   // 摘要
    Description  string   // 详细描述
    Author       string   // 作者
    AuthorEmail  string   // 作者邮箱
    License      string   // 许可证
    Homepage     string   // 主页
    Requires     []string // 依赖包
    RequiredBy   []string // 被依赖的包
    Location     string   // 安装位置
    Files        []string // 包含的文件
}
```

### SearchResult

搜索结果。

```go
type SearchResult struct {
    Name    string // 包名
    Version string // 最新版本
    Summary string // 摘要
    Score   float64 // 相关性评分
}
```

### OutdatedPackage

过时包信息。

```go
type OutdatedPackage struct {
    Name           string // 包名
    CurrentVersion string // 当前版本
    LatestVersion  string // 最新版本
    LatestType     string // 最新版本类型（stable, pre-release）
}
```

## 批量操作

### InstallPackages

批量安装多个包。

```go
func (m *Manager) InstallPackages(packages []*PackageSpec) error
```

**参数：**
- `packages` - 包规范列表

**返回值：**
- `error` - 错误信息

**示例：**

```go
packages := []*pip.PackageSpec{
    {Name: "requests"},
    {Name: "click", Version: ">=7.0"},
    {Name: "pydantic", Version: "^1.8.0"},
}

if err := manager.InstallPackages(packages); err != nil {
    log.Fatalf("批量安装失败: %v", err)
}
```

### UninstallPackages

批量卸载多个包。

```go
func (m *Manager) UninstallPackages(names []string) error
```

**参数：**
- `names` - 包名列表

**返回值：**
- `error` - 错误信息

**示例：**

```go
packages := []string{"requests", "click", "pydantic"}
if err := manager.UninstallPackages(packages); err != nil {
    log.Fatalf("批量卸载失败: %v", err)
}
```

## 依赖分析

### GetDependencies

获取包的依赖关系。

```go
func (m *Manager) GetDependencies(name string) (*DependencyTree, error)
```

**参数：**
- `name` - 包名

**返回值：**
- `*DependencyTree` - 依赖树
- `error` - 错误信息

**示例：**

```go
deps, err := manager.GetDependencies("django")
if err != nil {
    log.Fatalf("获取依赖失败: %v", err)
}

fmt.Printf("Django 的依赖:\n")
printDependencyTree(deps, 0)
```

### CheckDependencies

检查依赖冲突。

```go
func (m *Manager) CheckDependencies() ([]*DependencyConflict, error)
```

**返回值：**
- `[]*DependencyConflict` - 依赖冲突列表
- `error` - 错误信息

**示例：**

```go
conflicts, err := manager.CheckDependencies()
if err != nil {
    log.Fatalf("检查依赖失败: %v", err)
}

if len(conflicts) > 0 {
    fmt.Println("发现依赖冲突:")
    for _, conflict := range conflicts {
        fmt.Printf("- %s: %s\n", conflict.Package, conflict.Description)
    }
}
```

## 错误处理

包操作可能遇到的常见错误类型：

```go
if err := manager.InstallPackage(pkg); err != nil {
    switch {
    case pip.IsErrorType(err, pip.ErrorTypePackageNotFound):
        fmt.Printf("包 '%s' 未找到\n", pkg.Name)
    case pip.IsErrorType(err, pip.ErrorTypeVersionConflict):
        fmt.Printf("版本冲突: %v\n", err)
    case pip.IsErrorType(err, pip.ErrorTypeNetworkError):
        fmt.Printf("网络错误: %v\n", err)
    case pip.IsErrorType(err, pip.ErrorTypePermissionDenied):
        fmt.Printf("权限被拒绝: %v\n", err)
    default:
        fmt.Printf("其他错误: %v\n", err)
    }
}
```

## 最佳实践

1. **使用版本约束**：指定合理的版本范围而不是精确版本
2. **虚拟环境**：在虚拟环境中进行包操作
3. **错误处理**：适当处理不同类型的错误
4. **批量操作**：对多个包使用批量操作提高效率
5. **依赖检查**：定期检查依赖冲突

## 下一步

- 了解[虚拟环境 API](./virtual-environments.md)
- 查看[项目管理 API](./project-management.md)
- 学习[错误处理](./errors.md)
- 探索[类型定义](./types.md)
