# 类型定义

本页面包含 Go Pip SDK 中使用的所有数据结构和类型定义。

## 核心类型

### Config

SDK 配置结构。

```go
type Config struct {
    // Python 配置
    PythonPath string        // Python 可执行文件路径
    PipPath    string        // Pip 可执行文件路径
    
    // 网络配置
    DefaultIndex string      // 默认包索引 URL
    ExtraIndexes []string    // 额外包索引 URL
    TrustedHosts []string    // 受信任的主机
    Timeout      time.Duration // 网络超时时间
    Retries      int         // 重试次数
    
    // 缓存配置
    CacheDir     string      // 缓存目录
    NoCache      bool        // 禁用缓存
    
    // 日志配置
    LogLevel     string      // 日志级别
    LogFile      string      // 日志文件路径
    
    // 环境变量
    Environment  map[string]string // 环境变量
    
    // 额外选项
    ExtraOptions map[string]string // 额外的 pip 选项
}
```

### PackageSpec

包规范，定义要安装的包及其选项。

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
    Target         string            // 安装目标目录
    Prefix         string            // 安装前缀
}
```

### Package

已安装包的信息。

```go
type Package struct {
    Name      string    // 包名
    Version   string    // 版本
    Location  string    // 安装位置
    Editable  bool      // 是否为可编辑安装
    Installer string    // 安装器
    Metadata  *Metadata // 包元数据
}
```

### PackageInfo

包的详细信息。

```go
type PackageInfo struct {
    Name         string            // 包名
    Version      string            // 版本
    Summary      string            // 摘要
    Description  string            // 详细描述
    Author       string            // 作者
    AuthorEmail  string            // 作者邮箱
    Maintainer   string            // 维护者
    MaintainerEmail string         // 维护者邮箱
    License      string            // 许可证
    Homepage     string            // 主页
    Repository   string            // 仓库地址
    Documentation string           // 文档地址
    Keywords     []string          // 关键词
    Classifiers  []string          // 分类器
    Requires     []string          // 依赖包
    RequiredBy   []string          // 被依赖的包
    Location     string            // 安装位置
    Files        []string          // 包含的文件
    Metadata     map[string]string // 元数据
}
```

## 搜索和查询类型

### SearchResult

搜索结果。

```go
type SearchResult struct {
    Name        string  // 包名
    Version     string  // 最新版本
    Summary     string  // 摘要
    Description string  // 描述
    Score       float64 // 相关性评分
    Author      string  // 作者
    Homepage    string  // 主页
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
    Homepage       string // 主页
}
```

### DependencyTree

依赖树结构。

```go
type DependencyTree struct {
    Package      *Package           // 包信息
    Dependencies []*DependencyTree  // 依赖列表
    Depth        int                // 依赖深度
}
```

### DependencyConflict

依赖冲突信息。

```go
type DependencyConflict struct {
    Package     string   // 冲突的包
    Required    []string // 需要的版本
    Installed   string   // 已安装版本
    Conflicts   []string // 冲突的包
    Description string   // 冲突描述
}
```

## 虚拟环境类型

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
    Config        *VenvConfig // 环境配置
}
```

### VenvConfig

虚拟环境配置。

```go
type VenvConfig struct {
    Path          string
    PythonVersion string
    Options       *VenvOptions
}
```

## 项目管理类型

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
    BuildSystem  *BuildSystem
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

### BuildSystem

构建系统配置。

```go
type BuildSystem struct {
    Requires     []string // 构建依赖
    BuildBackend string   // 构建后端
}
```

## 命令执行类型

### CommandResult

命令执行结果。

```go
type CommandResult struct {
    Command   string        // 执行的命令
    Args      []string      // 命令参数
    Stdout    string        // 标准输出
    Stderr    string        // 错误输出
    ExitCode  int           // 退出代码
    Duration  time.Duration // 执行时间
    Success   bool          // 是否成功
}
```

### CommandOptions

命令执行选项。

```go
type CommandOptions struct {
    WorkingDir  string            // 工作目录
    Environment map[string]string // 环境变量
    Timeout     time.Duration     // 超时时间
    Input       string            // 标准输入
}
```

## 状态和监控类型

### ManagerStatus

管理器状态信息。

```go
type ManagerStatus struct {
    PythonVersion string    // Python 版本
    PipVersion    string    // Pip 版本
    VirtualEnv    string    // 当前虚拟环境
    PackageCount  int       // 已安装包数量
    CacheSize     int64     // 缓存大小
    LastUpdate    time.Time // 最后更新时间
    IsHealthy     bool      // 是否健康
}
```

### ValidationIssue

验证问题。

```go
type ValidationIssue struct {
    Type        string // 问题类型
    Severity    string // 严重程度
    Message     string // 问题描述
    Suggestion  string // 解决建议
    Location    string // 问题位置
}
```

### VenvIssue

虚拟环境问题。

```go
type VenvIssue struct {
    Type        string // 问题类型
    Severity    string // 严重程度
    Message     string // 问题描述
    Suggestion  string // 解决建议
    Path        string // 问题路径
}
```

## 日志类型

### LogLevel

日志级别。

```go
type LogLevel int

const (
    LogLevelTrace LogLevel = iota
    LogLevelDebug
    LogLevelInfo
    LogLevelWarn
    LogLevelError
)
```

### LogEntry

日志条目。

```go
type LogEntry struct {
    Level     LogLevel               // 日志级别
    Message   string                 // 日志消息
    Fields    map[string]interface{} // 字段
    Timestamp time.Time              // 时间戳
    Caller    string                 // 调用者
}
```

### LoggerConfig

日志记录器配置。

```go
type LoggerConfig struct {
    Level      LogLevel
    Output     io.Writer
    Format     LogFormat
    EnableFile bool
    LogFile    string
    MaxSize    int
    MaxBackups int
    MaxAge     int
    Compress   bool
    Color      bool
    Fields     map[string]interface{}
}
```

## 发布和迁移类型

### PublishOptions

发布选项。

```go
type PublishOptions struct {
    Repository string // 仓库 (pypi, testpypi)
    Username   string // 用户名
    Password   string // 密码
    Token      string // API 令牌
    SignKey    string // 签名密钥
    Comment    string // 发布注释
}
```

### MigrationOptions

迁移选项。

```go
type MigrationOptions struct {
    SourceFormat    string // 源格式
    TargetFormat    string // 目标格式
    CreateVenv      bool   // 创建虚拟环境
    UpdateGitignore bool   // 更新 .gitignore
    BackupOriginal  bool   // 备份原文件
}
```

### UpgradeOptions

升级选项。

```go
type UpgradeOptions struct {
    UpdateDependencies bool // 更新依赖
    UpdateTemplates    bool // 更新模板
    BackupOriginal     bool // 备份原文件
    Force              bool // 强制升级
}
```

## 接口类型

### Logger

日志记录器接口。

```go
type Logger interface {
    Trace(msg string, fields ...map[string]interface{})
    Debug(msg string, fields ...map[string]interface{})
    Info(msg string, fields ...map[string]interface{})
    Warn(msg string, fields ...map[string]interface{})
    Error(msg string, fields ...map[string]interface{})
    WithFields(fields map[string]interface{}) Logger
    Close() error
}
```

### PipManager

Pip 管理器接口。

```go
type PipManager interface {
    // 系统操作
    IsInstalled() (bool, error)
    Install() error
    GetVersion() (string, error)
    
    // 包操作
    InstallPackage(pkg *PackageSpec) error
    UninstallPackage(name string) error
    ListPackages() ([]*Package, error)
    ShowPackage(name string) (*PackageInfo, error)
    
    // 虚拟环境操作
    CreateVenv(path string) error
    ActivateVenv(path string) error
    DeactivateVenv() error
    
    // 项目操作
    InitProject(path string, opts *ProjectOptions) error
    ReadProjectConfig(path string) (*ProjectConfig, error)
    
    // 配置和日志
    SetCustomLogger(logger Logger)
    GetConfig() *Config
}
```

## 常量定义

### 默认值

```go
const (
    DefaultTimeout      = 60 * time.Second
    DefaultRetries      = 3
    DefaultCacheDir     = "~/.pip/cache"
    DefaultLogLevel     = "INFO"
    DefaultPythonPath   = "python"
    DefaultPipPath      = "pip"
    DefaultIndexURL     = "https://pypi.org/simple/"
)
```

### 文件名

```go
const (
    RequirementsFile    = "requirements.txt"
    DevRequirementsFile = "dev-requirements.txt"
    SetupPyFile         = "setup.py"
    PyprojectTomlFile   = "pyproject.toml"
    ManifestInFile      = "MANIFEST.in"
    GitignoreFile       = ".gitignore"
    ReadmeFile          = "README.md"
)
```

## 下一步

- 查看[管理器 API](./manager.md)
- 了解[包操作 API](./package-operations.md)
- 探索[虚拟环境 API](./virtual-environments.md)
- 学习[错误处理 API](./errors.md)
