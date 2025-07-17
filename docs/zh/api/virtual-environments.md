# 虚拟环境 API

虚拟环境 API 提供了完整的 Python 虚拟环境管理功能，包括创建、激活、停用和删除虚拟环境。

## 核心方法

### CreateVenv

创建新的虚拟环境。

```go
func (m *Manager) CreateVenv(path string) error
```

**参数：**
- `path` - 虚拟环境路径

**返回值：**
- `error` - 错误信息

**示例：**

```go
venvPath := "./my-env"
if err := manager.CreateVenv(venvPath); err != nil {
    log.Fatalf("创建虚拟环境失败: %v", err)
}
fmt.Println("✅ 虚拟环境创建成功")
```

### CreateVenvWithOptions

使用自定义选项创建虚拟环境。

```go
func (m *Manager) CreateVenvWithOptions(path string, opts *VenvOptions) error
```

**参数：**
- `path` - 虚拟环境路径
- `opts` - 虚拟环境选项

**返回值：**
- `error` - 错误信息

**示例：**

```go
opts := &pip.VenvOptions{
    PythonVersion: "3.9",
    SystemSitePackages: false,
    Prompt: "my-project",
    UpgradePip: true,
    InstallSetuptools: true,
    InstallWheel: true,
}

if err := manager.CreateVenvWithOptions("./my-env", opts); err != nil {
    log.Fatalf("创建虚拟环境失败: %v", err)
}
```

### ActivateVenv

激活虚拟环境。

```go
func (m *Manager) ActivateVenv(path string) error
```

**参数：**
- `path` - 虚拟环境路径

**返回值：**
- `error` - 错误信息

**示例：**

```go
if err := manager.ActivateVenv("./my-env"); err != nil {
    log.Fatalf("激活虚拟环境失败: %v", err)
}
fmt.Println("✅ 虚拟环境已激活")
```

### DeactivateVenv

停用当前虚拟环境。

```go
func (m *Manager) DeactivateVenv() error
```

**返回值：**
- `error` - 错误信息

**示例：**

```go
if err := manager.DeactivateVenv(); err != nil {
    log.Fatalf("停用虚拟环境失败: %v", err)
}
fmt.Println("✅ 虚拟环境已停用")
```

### IsVenvActive

检查虚拟环境是否激活。

```go
func (m *Manager) IsVenvActive() (bool, string)
```

**返回值：**
- `bool` - 是否激活
- `string` - 虚拟环境路径

**示例：**

```go
isActive, venvPath := manager.IsVenvActive()
if isActive {
    fmt.Printf("当前虚拟环境: %s\n", venvPath)
} else {
    fmt.Println("未在虚拟环境中")
}
```

### RemoveVenv

删除虚拟环境。

```go
func (m *Manager) RemoveVenv(path string) error
```

**参数：**
- `path` - 虚拟环境路径

**返回值：**
- `error` - 错误信息

**示例：**

```go
if err := manager.RemoveVenv("./old-env"); err != nil {
    log.Fatalf("删除虚拟环境失败: %v", err)
}
fmt.Println("✅ 虚拟环境已删除")
```

## 查询和管理

### ListVenvs

列出指定目录下的所有虚拟环境。

```go
func (m *Manager) ListVenvs(baseDir string) ([]*VenvInfo, error)
```

**参数：**
- `baseDir` - 基础目录

**返回值：**
- `[]*VenvInfo` - 虚拟环境信息列表
- `error` - 错误信息

**示例：**

```go
venvs, err := manager.ListVenvs("./envs")
if err != nil {
    log.Fatalf("列出虚拟环境失败: %v", err)
}

fmt.Printf("找到 %d 个虚拟环境:\n", len(venvs))
for _, venv := range venvs {
    fmt.Printf("- %s (Python %s)\n", venv.Name, venv.PythonVersion)
}
```

### GetVenvInfo

获取虚拟环境详细信息。

```go
func (m *Manager) GetVenvInfo(path string) (*VenvInfo, error)
```

**参数：**
- `path` - 虚拟环境路径

**返回值：**
- `*VenvInfo` - 虚拟环境信息
- `error` - 错误信息

**示例：**

```go
info, err := manager.GetVenvInfo("./my-env")
if err != nil {
    log.Fatalf("获取虚拟环境信息失败: %v", err)
}

fmt.Printf("名称: %s\n", info.Name)
fmt.Printf("路径: %s\n", info.Path)
fmt.Printf("Python 版本: %s\n", info.PythonVersion)
fmt.Printf("创建时间: %s\n", info.CreatedAt.Format("2006-01-02 15:04:05"))
fmt.Printf("最后使用: %s\n", info.LastUsed.Format("2006-01-02 15:04:05"))
```

## 数据类型

### VenvOptions

虚拟环境创建选项。

```go
type VenvOptions struct {
    PythonVersion      string // Python 版本
    SystemSitePackages bool   // 是否访问系统包
    Prompt             string // 提示符名称
    UpgradePip         bool   // 是否升级 pip
    InstallSetuptools  bool   // 是否安装 setuptools
    InstallWheel       bool   // 是否安装 wheel
    Clear              bool   // 是否清理现有环境
    Symlinks           bool   // 是否使用符号链接
    Copies             bool   // 是否使用复制
    WithoutPip         bool   // 是否不安装 pip
}
```

### VenvInfo

虚拟环境信息。

```go
type VenvInfo struct {
    Name          string    // 环境名称
    Path          string    // 环境路径
    PythonVersion string    // Python 版本
    PipVersion    string    // Pip 版本
    CreatedAt     time.Time // 创建时间
    LastUsed      time.Time // 最后使用时间
    Size          int64     // 环境大小（字节）
    PackageCount  int       // 已安装包数量
    IsActive      bool      // 是否当前激活
    HasRequirements bool    // 是否有 requirements.txt
}
```

## 高级功能

### CloneVenv

克隆虚拟环境。

```go
func (m *Manager) CloneVenv(sourcePath, targetPath string) error
```

**参数：**
- `sourcePath` - 源虚拟环境路径
- `targetPath` - 目标虚拟环境路径

**返回值：**
- `error` - 错误信息

**示例：**

```go
if err := manager.CloneVenv("./source-env", "./target-env"); err != nil {
    log.Fatalf("克隆虚拟环境失败: %v", err)
}
```

### ExportVenv

导出虚拟环境配置。

```go
func (m *Manager) ExportVenv(path, outputFile string) error
```

**参数：**
- `path` - 虚拟环境路径
- `outputFile` - 输出文件路径

**返回值：**
- `error` - 错误信息

**示例：**

```go
if err := manager.ExportVenv("./my-env", "environment.yml"); err != nil {
    log.Fatalf("导出虚拟环境失败: %v", err)
}
```

### ImportVenv

从配置文件导入虚拟环境。

```go
func (m *Manager) ImportVenv(configFile, targetPath string) error
```

**参数：**
- `configFile` - 配置文件路径
- `targetPath` - 目标虚拟环境路径

**返回值：**
- `error` - 错误信息

**示例：**

```go
if err := manager.ImportVenv("environment.yml", "./imported-env"); err != nil {
    log.Fatalf("导入虚拟环境失败: %v", err)
}
```

### UpdateVenv

更新虚拟环境。

```go
func (m *Manager) UpdateVenv(path string) error
```

**参数：**
- `path` - 虚拟环境路径

**返回值：**
- `error` - 错误信息

**示例：**

```go
if err := manager.UpdateVenv("./my-env"); err != nil {
    log.Fatalf("更新虚拟环境失败: %v", err)
}
```

### CleanVenv

清理虚拟环境缓存。

```go
func (m *Manager) CleanVenv(path string) error
```

**参数：**
- `path` - 虚拟环境路径

**返回值：**
- `error` - 错误信息

**示例：**

```go
if err := manager.CleanVenv("./my-env"); err != nil {
    log.Fatalf("清理虚拟环境失败: %v", err)
}
```

## 批量操作

### CreateMultipleVenvs

批量创建虚拟环境。

```go
func (m *Manager) CreateMultipleVenvs(configs []*VenvConfig) error
```

**示例：**

```go
configs := []*pip.VenvConfig{
    {Path: "./env1", PythonVersion: "3.8"},
    {Path: "./env2", PythonVersion: "3.9"},
    {Path: "./env3", PythonVersion: "3.10"},
}

if err := manager.CreateMultipleVenvs(configs); err != nil {
    log.Fatalf("批量创建虚拟环境失败: %v", err)
}
```

### RemoveMultipleVenvs

批量删除虚拟环境。

```go
func (m *Manager) RemoveMultipleVenvs(paths []string) error
```

**示例：**

```go
paths := []string{"./old-env1", "./old-env2", "./old-env3"}
if err := manager.RemoveMultipleVenvs(paths); err != nil {
    log.Fatalf("批量删除虚拟环境失败: %v", err)
}
```

## 环境检查

### ValidateVenv

验证虚拟环境完整性。

```go
func (m *Manager) ValidateVenv(path string) ([]*VenvIssue, error)
```

**示例：**

```go
issues, err := manager.ValidateVenv("./my-env")
if err != nil {
    log.Fatalf("验证虚拟环境失败: %v", err)
}

if len(issues) > 0 {
    fmt.Println("发现以下问题:")
    for _, issue := range issues {
        fmt.Printf("- %s: %s\n", issue.Type, issue.Message)
    }
}
```

### RepairVenv

修复虚拟环境。

```go
func (m *Manager) RepairVenv(path string) error
```

**示例：**

```go
if err := manager.RepairVenv("./broken-env"); err != nil {
    log.Fatalf("修复虚拟环境失败: %v", err)
}
```

## 错误处理

虚拟环境操作可能遇到的错误类型：

```go
if err := manager.CreateVenv(path); err != nil {
    switch {
    case pip.IsErrorType(err, pip.ErrorTypeVenvExists):
        fmt.Println("虚拟环境已存在")
    case pip.IsErrorType(err, pip.ErrorTypePythonNotFound):
        fmt.Println("Python 未找到")
    case pip.IsErrorType(err, pip.ErrorTypePermissionDenied):
        fmt.Println("权限被拒绝")
    case pip.IsErrorType(err, pip.ErrorTypeInsufficientSpace):
        fmt.Println("磁盘空间不足")
    default:
        fmt.Printf("其他错误: %v\n", err)
    }
}
```

## 最佳实践

1. **为每个项目创建独立的虚拟环境**
2. **使用有意义的环境名称**
3. **定期清理不用的虚拟环境**
4. **在虚拟环境中安装项目依赖**
5. **导出环境配置以便复现**

## 下一步

- 查看[包操作 API](./package-operations.md)
- 了解[项目管理 API](./project-management.md)
- 探索[错误处理](./errors.md)
- 学习[类型定义](./types.md)
