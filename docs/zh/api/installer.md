# 安装器 API

安装器 API 提供了 pip 自动安装和管理功能，确保 pip 在系统中正确安装和配置。

## 核心接口

### Installer

pip 安装器接口。

```go
type Installer interface {
    // 检查和安装
    IsInstalled() (bool, error)
    Install() error
    Uninstall() error
    
    // 版本管理
    GetVersion() (string, error)
    GetLatestVersion() (string, error)
    Upgrade() error
    
    // 配置管理
    Configure(config *InstallerConfig) error
    GetConfig() *InstallerConfig
    
    // 状态查询
    GetStatus() (*InstallerStatus, error)
    Validate() error
}
```

## 创建安装器

### NewInstaller

创建新的 pip 安装器。

```go
func NewInstaller() Installer
```

**返回值：**
- `Installer` - 安装器实例

**示例：**

```go
installer := pip.NewInstaller()

// 检查 pip 是否已安装
installed, err := installer.IsInstalled()
if err != nil {
    log.Fatalf("检查 pip 安装状态失败: %v", err)
}

if !installed {
    fmt.Println("Pip 未安装，正在安装...")
    if err := installer.Install(); err != nil {
        log.Fatalf("安装 pip 失败: %v", err)
    }
    fmt.Println("✅ Pip 安装成功")
}
```

### NewInstallerWithConfig

使用配置创建安装器。

```go
func NewInstallerWithConfig(config *InstallerConfig) Installer
```

**参数：**
- `config` - 安装器配置

**返回值：**
- `Installer` - 安装器实例

**示例：**

```go
config := &pip.InstallerConfig{
    PythonPath:    "/usr/bin/python3",
    InstallMethod: pip.InstallMethodGetPip,
    ForceReinstall: false,
    UserInstall:   false,
}

installer := pip.NewInstallerWithConfig(config)
```

## 安装方法

### Install

安装 pip。

```go
func (i *Installer) Install() error
```

**返回值：**
- `error` - 错误信息

**示例：**

```go
if err := installer.Install(); err != nil {
    log.Fatalf("安装 pip 失败: %v", err)
}
fmt.Println("Pip 安装成功")
```

### InstallWithMethod

使用指定方法安装 pip。

```go
func (i *Installer) InstallWithMethod(method InstallMethod) error
```

**参数：**
- `method` - 安装方法

**返回值：**
- `error` - 错误信息

**示例：**

```go
// 使用 get-pip.py 安装
if err := installer.InstallWithMethod(pip.InstallMethodGetPip); err != nil {
    log.Fatalf("使用 get-pip.py 安装失败: %v", err)
}

// 使用包管理器安装
if err := installer.InstallWithMethod(pip.InstallMethodPackageManager); err != nil {
    log.Fatalf("使用包管理器安装失败: %v", err)
}
```

### Uninstall

卸载 pip。

```go
func (i *Installer) Uninstall() error
```

**返回值：**
- `error` - 错误信息

**示例：**

```go
if err := installer.Uninstall(); err != nil {
    log.Fatalf("卸载 pip 失败: %v", err)
}
fmt.Println("Pip 卸载成功")
```

## 版本管理

### GetVersion

获取当前安装的 pip 版本。

```go
func (i *Installer) GetVersion() (string, error)
```

**返回值：**
- `string` - pip 版本号
- `error` - 错误信息

**示例：**

```go
version, err := installer.GetVersion()
if err != nil {
    log.Fatalf("获取 pip 版本失败: %v", err)
}
fmt.Printf("当前 pip 版本: %s\n", version)
```

### GetLatestVersion

获取最新的 pip 版本。

```go
func (i *Installer) GetLatestVersion() (string, error)
```

**返回值：**
- `string` - 最新版本号
- `error` - 错误信息

**示例：**

```go
latest, err := installer.GetLatestVersion()
if err != nil {
    log.Fatalf("获取最新版本失败: %v", err)
}
fmt.Printf("最新 pip 版本: %s\n", latest)
```

### Upgrade

升级 pip 到最新版本。

```go
func (i *Installer) Upgrade() error
```

**返回值：**
- `error` - 错误信息

**示例：**

```go
current, _ := installer.GetVersion()
latest, _ := installer.GetLatestVersion()

if current != latest {
    fmt.Printf("升级 pip 从 %s 到 %s...\n", current, latest)
    if err := installer.Upgrade(); err != nil {
        log.Fatalf("升级 pip 失败: %v", err)
    }
    fmt.Println("✅ Pip 升级成功")
}
```

## 配置管理

### InstallerConfig

安装器配置结构。

```go
type InstallerConfig struct {
    // Python 配置
    PythonPath string // Python 可执行文件路径
    
    // 安装配置
    InstallMethod  InstallMethod // 安装方法
    ForceReinstall bool          // 强制重新安装
    UserInstall    bool          // 用户级安装
    
    // 网络配置
    IndexURL     string        // 包索引 URL
    TrustedHosts []string      // 受信任的主机
    Timeout      time.Duration // 网络超时
    Retries      int           // 重试次数
    
    // 代理配置
    ProxyURL string // 代理 URL
    
    // 临时目录
    TempDir string // 临时目录路径
    
    // 日志配置
    Verbose bool   // 详细输出
    Quiet   bool   // 静默模式
    LogFile string // 日志文件
}
```

### InstallMethod

安装方法枚举。

```go
type InstallMethod int

const (
    InstallMethodAuto           InstallMethod = iota // 自动选择
    InstallMethodGetPip                              // 使用 get-pip.py
    InstallMethodPackageManager                      // 使用系统包管理器
    InstallMethodEnsurepip                           // 使用 ensurepip 模块
    InstallMethodBootstrap                           // 使用 bootstrap 脚本
)
```

### Configure

配置安装器。

```go
func (i *Installer) Configure(config *InstallerConfig) error
```

**参数：**
- `config` - 安装器配置

**返回值：**
- `error` - 错误信息

**示例：**

```go
config := &pip.InstallerConfig{
    PythonPath:     "/usr/bin/python3",
    InstallMethod:  pip.InstallMethodGetPip,
    UserInstall:    true,
    IndexURL:       "https://pypi.tuna.tsinghua.edu.cn/simple/",
    TrustedHosts:   []string{"pypi.tuna.tsinghua.edu.cn"},
    Timeout:        120 * time.Second,
    Retries:        5,
}

if err := installer.Configure(config); err != nil {
    log.Fatalf("配置安装器失败: %v", err)
}
```

## 状态查询

### InstallerStatus

安装器状态信息。

```go
type InstallerStatus struct {
    IsInstalled    bool      // 是否已安装
    Version        string    // 当前版本
    LatestVersion  string    // 最新版本
    InstallPath    string    // 安装路径
    PythonVersion  string    // Python 版本
    InstallMethod  string    // 安装方法
    LastChecked    time.Time // 最后检查时间
    NeedsUpgrade   bool      // 是否需要升级
}
```

### GetStatus

获取安装器状态。

```go
func (i *Installer) GetStatus() (*InstallerStatus, error)
```

**返回值：**
- `*InstallerStatus` - 状态信息
- `error` - 错误信息

**示例：**

```go
status, err := installer.GetStatus()
if err != nil {
    log.Fatalf("获取状态失败: %v", err)
}

fmt.Printf("Pip 状态:\n")
fmt.Printf("  已安装: %t\n", status.IsInstalled)
fmt.Printf("  当前版本: %s\n", status.Version)
fmt.Printf("  最新版本: %s\n", status.LatestVersion)
fmt.Printf("  安装路径: %s\n", status.InstallPath)
fmt.Printf("  需要升级: %t\n", status.NeedsUpgrade)
```

### Validate

验证 pip 安装。

```go
func (i *Installer) Validate() error
```

**返回值：**
- `error` - 验证错误

**示例：**

```go
if err := installer.Validate(); err != nil {
    fmt.Printf("Pip 安装验证失败: %v\n", err)
    
    // 尝试修复
    if err := installer.Install(); err != nil {
        log.Fatalf("修复安装失败: %v", err)
    }
} else {
    fmt.Println("✅ Pip 安装验证通过")
}
```

## 高级功能

### DownloadGetPip

下载 get-pip.py 脚本。

```go
func (i *Installer) DownloadGetPip(destPath string) error
```

**参数：**
- `destPath` - 下载目标路径

**返回值：**
- `error` - 错误信息

**示例：**

```go
if err := installer.DownloadGetPip("./get-pip.py"); err != nil {
    log.Fatalf("下载 get-pip.py 失败: %v", err)
}
```

### InstallFromScript

从脚本安装 pip。

```go
func (i *Installer) InstallFromScript(scriptPath string) error
```

**参数：**
- `scriptPath` - 脚本路径

**返回值：**
- `error` - 错误信息

**示例：**

```go
if err := installer.InstallFromScript("./get-pip.py"); err != nil {
    log.Fatalf("从脚本安装失败: %v", err)
}
```

### CheckDependencies

检查 pip 依赖。

```go
func (i *Installer) CheckDependencies() ([]*Dependency, error)
```

**返回值：**
- `[]*Dependency` - 依赖列表
- `error` - 错误信息

**示例：**

```go
deps, err := installer.CheckDependencies()
if err != nil {
    log.Fatalf("检查依赖失败: %v", err)
}

fmt.Println("Pip 依赖:")
for _, dep := range deps {
    fmt.Printf("- %s %s (%s)\n", dep.Name, dep.Version, dep.Status)
}
```

## 平台特定功能

### 检测操作系统

```go
func DetectOS() OSType
func GetPackageManager() PackageManager
```

**示例：**

```go
osType := pip.DetectOS()
switch osType {
case pip.OSTypeLinux:
    fmt.Println("Linux 系统")
case pip.OSTypeMacOS:
    fmt.Println("macOS 系统")
case pip.OSTypeWindows:
    fmt.Println("Windows 系统")
}

pkgMgr := pip.GetPackageManager()
fmt.Printf("包管理器: %s\n", pkgMgr)
```

### 系统包管理器安装

```go
func (i *Installer) InstallWithPackageManager() error
```

**示例：**

```go
// 在 Ubuntu/Debian 上使用 apt 安装
// 在 CentOS/RHEL 上使用 yum/dnf 安装
// 在 macOS 上使用 brew 安装
if err := installer.InstallWithPackageManager(); err != nil {
    log.Fatalf("使用包管理器安装失败: %v", err)
}
```

## 错误处理

### 安装器特定错误

```go
const (
    ErrorTypePipAlreadyInstalled ErrorType = "pip_already_installed"
    ErrorTypePipNotInstalled     ErrorType = "pip_not_installed"
    ErrorTypeUnsupportedPlatform ErrorType = "unsupported_platform"
    ErrorTypeDownloadFailed      ErrorType = "download_failed"
    ErrorTypeInstallationFailed  ErrorType = "installation_failed"
)
```

**示例：**

```go
if err := installer.Install(); err != nil {
    switch pip.GetErrorType(err) {
    case pip.ErrorTypePipAlreadyInstalled:
        fmt.Println("Pip 已经安装")
    case pip.ErrorTypeUnsupportedPlatform:
        fmt.Println("不支持的平台")
    case pip.ErrorTypeDownloadFailed:
        fmt.Println("下载失败，请检查网络连接")
    case pip.ErrorTypeInstallationFailed:
        fmt.Println("安装失败，请检查权限")
    default:
        fmt.Printf("未知错误: %v\n", err)
    }
}
```

## 实用工具

### 环境检查

```go
func CheckPythonEnvironment() (*PythonEnvironment, error)
```

**示例：**

```go
env, err := pip.CheckPythonEnvironment()
if err != nil {
    log.Fatalf("检查 Python 环境失败: %v", err)
}

fmt.Printf("Python 环境:\n")
fmt.Printf("  版本: %s\n", env.Version)
fmt.Printf("  路径: %s\n", env.Path)
fmt.Printf("  架构: %s\n", env.Architecture)
fmt.Printf("  支持 pip: %t\n", env.SupportsPip)
```

### 清理工具

```go
func (i *Installer) Cleanup() error
```

**示例：**

```go
// 清理临时文件和缓存
if err := installer.Cleanup(); err != nil {
    log.Printf("清理失败: %v", err)
}
```

## 最佳实践

1. **检查后安装**：总是先检查 pip 是否已安装
2. **选择合适的安装方法**：根据环境选择最适合的安装方法
3. **处理权限问题**：考虑使用用户级安装避免权限问题
4. **网络配置**：在网络受限环境中配置代理和镜像源
5. **验证安装**：安装后验证 pip 是否正常工作

## 下一步

- 查看[管理器 API](./manager.md)
- 了解[错误处理 API](./errors.md)
- 探索[日志记录 API](./logger.md)
- 学习[类型定义](./types.md)
