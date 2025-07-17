# 基本用法

本指南展示了 Go Pip SDK 的基本使用方法，适合初学者快速上手。

## 安装和设置

### 导入 SDK

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)
```

### 创建管理器

```go
func main() {
    // 使用默认配置创建管理器
    manager := pip.NewManager(nil)
    
    // 或者使用自定义配置
    config := &pip.Config{
        Timeout:  60 * time.Second,
        LogLevel: "INFO",
    }
    manager = pip.NewManager(config)
}
```

## 检查系统状态

### 检查 Python 和 Pip

```go
func checkSystem(manager *pip.Manager) error {
    // 检查 Python 是否可用
    pythonPath, err := manager.GetPythonPath()
    if err != nil {
        return fmt.Errorf("Python 未找到: %w", err)
    }
    fmt.Printf("Python 路径: %s\n", pythonPath)
    
    // 检查 pip 是否已安装
    installed, err := manager.IsInstalled()
    if err != nil {
        return fmt.Errorf("检查 pip 失败: %w", err)
    }
    
    if !installed {
        fmt.Println("Pip 未安装，正在安装...")
        if err := manager.Install(); err != nil {
            return fmt.Errorf("安装 pip 失败: %w", err)
        }
        fmt.Println("✅ Pip 安装成功")
    } else {
        fmt.Println("✅ Pip 已安装")
    }
    
    // 获取 pip 版本
    version, err := manager.GetVersion()
    if err != nil {
        return fmt.Errorf("获取 pip 版本失败: %w", err)
    }
    fmt.Printf("Pip 版本: %s\n", version)
    
    return nil
}
```

## 基本包操作

### 安装单个包

```go
func installSinglePackage() {
    manager := pip.NewManager(nil)
    
    // 创建包规范
    pkg := &pip.PackageSpec{
        Name: "requests",
    }
    
    fmt.Printf("正在安装 %s...\n", pkg.Name)
    if err := manager.InstallPackage(pkg); err != nil {
        log.Fatalf("安装失败: %v", err)
    }
    
    fmt.Println("✅ 包安装成功！")
}
```

### 安装指定版本的包

```go
func installSpecificVersion() {
    manager := pip.NewManager(nil)
    
    // 安装特定版本
    pkg := &pip.PackageSpec{
        Name:    "django",
        Version: "4.2.0",
    }
    
    if err := manager.InstallPackage(pkg); err != nil {
        log.Fatalf("安装 Django 4.2.0 失败: %v", err)
    }
    
    fmt.Println("✅ Django 4.2.0 安装成功")
}
```

### 安装多个包

```go
func installMultiplePackages() {
    manager := pip.NewManager(nil)
    
    packages := []*pip.PackageSpec{
        {Name: "requests"},
        {Name: "click"},
        {Name: "pydantic"},
    }
    
    for _, pkg := range packages {
        fmt.Printf("正在安装 %s...\n", pkg.Name)
        if err := manager.InstallPackage(pkg); err != nil {
            fmt.Printf("❌ 安装 %s 失败: %v\n", pkg.Name, err)
            continue
        }
        fmt.Printf("✅ %s 安装成功\n", pkg.Name)
    }
}
```

### 列出已安装的包

```go
func listInstalledPackages() {
    manager := pip.NewManager(nil)
    
    packages, err := manager.ListPackages()
    if err != nil {
        log.Fatalf("列出包失败: %v", err)
    }
    
    fmt.Printf("找到 %d 个已安装的包:\n", len(packages))
    for _, pkg := range packages {
        fmt.Printf("- %s %s\n", pkg.Name, pkg.Version)
    }
}
```

### 查看包详细信息

```go
func showPackageInfo() {
    manager := pip.NewManager(nil)
    
    info, err := manager.ShowPackage("requests")
    if err != nil {
        log.Fatalf("获取包信息失败: %v", err)
    }
    
    fmt.Printf("包名: %s\n", info.Name)
    fmt.Printf("版本: %s\n", info.Version)
    fmt.Printf("摘要: %s\n", info.Summary)
    fmt.Printf("作者: %s\n", info.Author)
    fmt.Printf("许可证: %s\n", info.License)
    
    if len(info.Requires) > 0 {
        fmt.Printf("依赖: %s\n", strings.Join(info.Requires, ", "))
    }
}
```

### 搜索包

```go
func searchPackages() {
    manager := pip.NewManager(nil)
    
    results, err := manager.SearchPackages("web framework")
    if err != nil {
        log.Fatalf("搜索失败: %v", err)
    }
    
    fmt.Printf("找到 %d 个搜索结果:\n", len(results))
    for i, result := range results {
        if i >= 5 { // 只显示前5个结果
            break
        }
        fmt.Printf("%d. %s - %s\n", i+1, result.Name, result.Summary)
    }
}
```

### 卸载包

```go
func uninstallPackage() {
    manager := pip.NewManager(nil)
    
    packageName := "old-package"
    fmt.Printf("正在卸载 %s...\n", packageName)
    
    if err := manager.UninstallPackage(packageName); err != nil {
        log.Fatalf("卸载失败: %v", err)
    }
    
    fmt.Printf("✅ %s 卸载成功\n", packageName)
}
```

## 需求文件操作

### 从需求文件安装

```go
func installFromRequirements() {
    manager := pip.NewManager(nil)
    
    fmt.Println("从 requirements.txt 安装包...")
    if err := manager.InstallRequirements("requirements.txt"); err != nil {
        log.Fatalf("从需求文件安装失败: %v", err)
    }
    
    fmt.Println("✅ 所有依赖安装完成")
}
```

### 生成需求文件

```go
func generateRequirements() {
    manager := pip.NewManager(nil)
    
    fmt.Println("生成 requirements.txt...")
    if err := manager.GenerateRequirements("requirements.txt"); err != nil {
        log.Fatalf("生成需求文件失败: %v", err)
    }
    
    fmt.Println("✅ requirements.txt 生成完成")
}
```

## 基本错误处理

### 简单错误处理

```go
func basicErrorHandling() {
    manager := pip.NewManager(nil)
    
    pkg := &pip.PackageSpec{
        Name: "nonexistent-package-12345",
    }
    
    if err := manager.InstallPackage(pkg); err != nil {
        // 检查错误类型
        if pip.IsErrorType(err, pip.ErrorTypePackageNotFound) {
            fmt.Printf("包 '%s' 未找到\n", pkg.Name)
        } else if pip.IsErrorType(err, pip.ErrorTypeNetworkError) {
            fmt.Println("网络连接错误")
        } else {
            fmt.Printf("安装失败: %v\n", err)
        }
        return
    }
    
    fmt.Println("包安装成功")
}
```

### 重试机制

```go
func installWithRetry(manager *pip.Manager, pkg *pip.PackageSpec) error {
    maxRetries := 3
    
    for i := 0; i < maxRetries; i++ {
        err := manager.InstallPackage(pkg)
        if err == nil {
            return nil // 成功
        }
        
        // 只对网络错误重试
        if pip.IsErrorType(err, pip.ErrorTypeNetworkError) {
            fmt.Printf("网络错误，重试 %d/%d...\n", i+1, maxRetries)
            time.Sleep(time.Second * time.Duration(i+1))
            continue
        }
        
        // 其他错误不重试
        return err
    }
    
    return fmt.Errorf("重试 %d 次后仍然失败", maxRetries)
}
```

## 完整示例

### 基本包管理工具

```go
package main

import (
    "fmt"
    "log"
    "os"
    "strings"
    "time"
    
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("用法: go run main.go <command> [args...]")
        fmt.Println("命令:")
        fmt.Println("  install <package>  - 安装包")
        fmt.Println("  uninstall <package> - 卸载包")
        fmt.Println("  list              - 列出已安装的包")
        fmt.Println("  search <query>    - 搜索包")
        fmt.Println("  info <package>    - 显示包信息")
        return
    }
    
    manager := pip.NewManager(&pip.Config{
        LogLevel: "INFO",
        Timeout:  60 * time.Second,
    })
    
    command := os.Args[1]
    
    switch command {
    case "install":
        if len(os.Args) < 3 {
            log.Fatal("请指定要安装的包名")
        }
        installPackageCommand(manager, os.Args[2])
        
    case "uninstall":
        if len(os.Args) < 3 {
            log.Fatal("请指定要卸载的包名")
        }
        uninstallPackageCommand(manager, os.Args[2])
        
    case "list":
        listPackagesCommand(manager)
        
    case "search":
        if len(os.Args) < 3 {
            log.Fatal("请指定搜索查询")
        }
        searchPackagesCommand(manager, strings.Join(os.Args[2:], " "))
        
    case "info":
        if len(os.Args) < 3 {
            log.Fatal("请指定包名")
        }
        showPackageInfoCommand(manager, os.Args[2])
        
    default:
        log.Fatalf("未知命令: %s", command)
    }
}

func installPackageCommand(manager *pip.Manager, packageName string) {
    pkg := &pip.PackageSpec{Name: packageName}
    
    fmt.Printf("正在安装 %s...\n", packageName)
    if err := manager.InstallPackage(pkg); err != nil {
        log.Fatalf("安装失败: %v", err)
    }
    
    fmt.Printf("✅ %s 安装成功\n", packageName)
}

func uninstallPackageCommand(manager *pip.Manager, packageName string) {
    fmt.Printf("正在卸载 %s...\n", packageName)
    if err := manager.UninstallPackage(packageName); err != nil {
        log.Fatalf("卸载失败: %v", err)
    }
    
    fmt.Printf("✅ %s 卸载成功\n", packageName)
}

func listPackagesCommand(manager *pip.Manager) {
    packages, err := manager.ListPackages()
    if err != nil {
        log.Fatalf("列出包失败: %v", err)
    }
    
    fmt.Printf("已安装的包 (%d 个):\n", len(packages))
    for _, pkg := range packages {
        fmt.Printf("  %s %s\n", pkg.Name, pkg.Version)
    }
}

func searchPackagesCommand(manager *pip.Manager, query string) {
    results, err := manager.SearchPackages(query)
    if err != nil {
        log.Fatalf("搜索失败: %v", err)
    }
    
    fmt.Printf("搜索 '%s' 的结果 (%d 个):\n", query, len(results))
    for i, result := range results {
        if i >= 10 { // 只显示前10个结果
            break
        }
        fmt.Printf("  %s - %s\n", result.Name, result.Summary)
    }
}

func showPackageInfoCommand(manager *pip.Manager, packageName string) {
    info, err := manager.ShowPackage(packageName)
    if err != nil {
        log.Fatalf("获取包信息失败: %v", err)
    }
    
    fmt.Printf("包信息:\n")
    fmt.Printf("  名称: %s\n", info.Name)
    fmt.Printf("  版本: %s\n", info.Version)
    fmt.Printf("  摘要: %s\n", info.Summary)
    fmt.Printf("  作者: %s\n", info.Author)
    fmt.Printf("  许可证: %s\n", info.License)
    
    if len(info.Requires) > 0 {
        fmt.Printf("  依赖: %s\n", strings.Join(info.Requires, ", "))
    }
    
    if info.Homepage != "" {
        fmt.Printf("  主页: %s\n", info.Homepage)
    }
}
```

## 运行示例

保存上面的代码为 `main.go`，然后可以这样使用：

```bash
# 安装包
go run main.go install requests

# 列出已安装的包
go run main.go list

# 搜索包
go run main.go search "web framework"

# 查看包信息
go run main.go info requests

# 卸载包
go run main.go uninstall requests
```

## 常见问题

### 1. Python 未找到

```go
// 解决方案：指定 Python 路径
config := &pip.Config{
    PythonPath: "/usr/bin/python3", // Linux/macOS
    // PythonPath: "C:\\Python39\\python.exe", // Windows
}
manager := pip.NewManager(config)
```

### 2. 网络连接问题

```go
// 解决方案：使用镜像源和增加超时
config := &pip.Config{
    DefaultIndex: "https://pypi.tuna.tsinghua.edu.cn/simple/",
    TrustedHosts: []string{"pypi.tuna.tsinghua.edu.cn"},
    Timeout:      120 * time.Second,
    Retries:      5,
}
manager := pip.NewManager(config)
```

### 3. 权限问题

```go
// 解决方案：使用用户级安装
pkg := &pip.PackageSpec{
    Name:        "package-name",
    UserInstall: true,
}
```

## 下一步

- 学习[包管理](./package-management.md)的高级功能
- 了解[虚拟环境](./virtual-environments.md)管理
- 探索[项目初始化](./project-initialization.md)
- 查看[高级用法](./advanced-usage.md)示例
