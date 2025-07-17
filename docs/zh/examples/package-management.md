# 包管理示例

本页面提供了详细的包管理示例，展示如何使用 Go Pip SDK 进行各种包操作。

## 基本包操作

### 安装单个包

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func installSinglePackage() {
    manager := pip.NewManager(nil)
    
    // 基本安装
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
    
    packages := []*pip.PackageSpec{
        // 精确版本
        {Name: "django", Version: "4.2.0"},
        
        // 版本范围
        {Name: "numpy", Version: ">=1.20.0,<2.0.0"},
        
        // 最小版本
        {Name: "requests", Version: ">=2.25.0"},
        
        // 兼容版本
        {Name: "click", Version: "~=7.1.0"},
    }
    
    for _, pkg := range packages {
        fmt.Printf("安装 %s %s...\n", pkg.Name, pkg.Version)
        if err := manager.InstallPackage(pkg); err != nil {
            fmt.Printf("❌ 安装 %s 失败: %v\n", pkg.Name, err)
            continue
        }
        fmt.Printf("✅ %s 安装成功\n", pkg.Name)
    }
}
```

### 安装带额外依赖的包

```go
func installWithExtras() {
    manager := pip.NewManager(nil)
    
    // 安装 FastAPI 及其开发依赖
    pkg := &pip.PackageSpec{
        Name:   "fastapi",
        Extras: []string{"dev", "test"},
    }
    
    if err := manager.InstallPackage(pkg); err != nil {
        log.Fatalf("安装 FastAPI 失败: %v", err)
    }
    
    fmt.Println("✅ FastAPI 及开发依赖安装成功")
}
```

## 高级安装选项

### 可编辑安装

```go
func editableInstall() {
    manager := pip.NewManager(nil)
    
    // 可编辑安装本地包
    pkg := &pip.PackageSpec{
        Name:     "./my-local-package",
        Editable: true,
    }
    
    if err := manager.InstallPackage(pkg); err != nil {
        log.Fatalf("可编辑安装失败: %v", err)
    }
    
    fmt.Println("✅ 本地包可编辑安装成功")
}
```

### 从 Git 仓库安装

```go
func installFromGit() {
    manager := pip.NewManager(nil)
    
    packages := []*pip.PackageSpec{
        // 从 GitHub 安装
        {Name: "git+https://github.com/user/repo.git"},
        
        // 指定分支
        {Name: "git+https://github.com/user/repo.git@develop"},
        
        // 指定标签
        {Name: "git+https://github.com/user/repo.git@v1.0.0"},
        
        // 指定提交
        {Name: "git+https://github.com/user/repo.git@abc123"},
    }
    
    for _, pkg := range packages {
        fmt.Printf("从 Git 安装 %s...\n", pkg.Name)
        if err := manager.InstallPackage(pkg); err != nil {
            fmt.Printf("❌ 安装失败: %v\n", err)
            continue
        }
        fmt.Printf("✅ 安装成功\n")
    }
}
```

### 用户级安装

```go
func userInstall() {
    manager := pip.NewManager(nil)
    
    // 用户级安装（避免权限问题）
    pkg := &pip.PackageSpec{
        Name:        "jupyter",
        UserInstall: true,
    }
    
    if err := manager.InstallPackage(pkg); err != nil {
        log.Fatalf("用户级安装失败: %v", err)
    }
    
    fmt.Println("✅ Jupyter 用户级安装成功")
}
```

## 批量包操作

### 批量安装

```go
func batchInstall() {
    manager := pip.NewManager(nil)
    
    packages := []*pip.PackageSpec{
        {Name: "requests", Version: ">=2.25.0"},
        {Name: "click", Version: ">=7.0"},
        {Name: "pydantic", Version: ">=1.8.0"},
        {Name: "fastapi", Version: ">=0.68.0"},
        {Name: "uvicorn", Version: ">=0.15.0"},
    }
    
    fmt.Printf("开始批量安装 %d 个包...\n", len(packages))
    
    var successful []string
    var failed []string
    
    for _, pkg := range packages {
        fmt.Printf("正在安装 %s...\n", pkg.Name)
        if err := manager.InstallPackage(pkg); err != nil {
            failed = append(failed, pkg.Name)
            fmt.Printf("❌ %s 安装失败: %v\n", pkg.Name, err)
        } else {
            successful = append(successful, pkg.Name)
            fmt.Printf("✅ %s 安装成功\n", pkg.Name)
        }
    }
    
    fmt.Printf("\n安装完成: %d 成功, %d 失败\n", len(successful), len(failed))
    if len(failed) > 0 {
        fmt.Printf("失败的包: %v\n", failed)
    }
}
```

### 并发安装

```go
func concurrentInstall() {
    manager := pip.NewManager(nil)
    
    packages := []*pip.PackageSpec{
        {Name: "requests"},
        {Name: "click"},
        {Name: "pydantic"},
        {Name: "fastapi"},
        {Name: "uvicorn"},
    }
    
    var wg sync.WaitGroup
    errChan := make(chan error, len(packages))
    semaphore := make(chan struct{}, 3) // 限制并发数为 3
    
    for _, pkg := range packages {
        wg.Add(1)
        go func(p *pip.PackageSpec) {
            defer wg.Done()
            
            // 获取信号量
            semaphore <- struct{}{}
            defer func() { <-semaphore }()
            
            fmt.Printf("正在安装 %s...\n", p.Name)
            if err := manager.InstallPackage(p); err != nil {
                errChan <- fmt.Errorf("安装 %s 失败: %w", p.Name, err)
            } else {
                fmt.Printf("✅ %s 安装成功\n", p.Name)
            }
        }(pkg)
    }
    
    wg.Wait()
    close(errChan)
    
    var errors []error
    for err := range errChan {
        errors = append(errors, err)
    }
    
    if len(errors) > 0 {
        fmt.Printf("并发安装中有 %d 个包失败:\n", len(errors))
        for _, err := range errors {
            fmt.Printf("- %v\n", err)
        }
    } else {
        fmt.Println("所有包并发安装成功！")
    }
}
```

## 包查询和信息

### 列出已安装的包

```go
func listInstalledPackages() {
    manager := pip.NewManager(nil)
    
    packages, err := manager.ListPackages()
    if err != nil {
        log.Fatalf("列出包失败: %v", err)
    }
    
    fmt.Printf("已安装的包 (%d 个):\n", len(packages))
    fmt.Println("名称\t\t版本\t\t位置")
    fmt.Println("----\t\t----\t\t----")
    
    for _, pkg := range packages {
        fmt.Printf("%-20s %-15s %s\n", pkg.Name, pkg.Version, pkg.Location)
    }
}
```

### 搜索包

```go
func searchPackages() {
    manager := pip.NewManager(nil)
    
    queries := []string{
        "web framework",
        "machine learning",
        "data analysis",
        "testing",
    }
    
    for _, query := range queries {
        fmt.Printf("\n搜索: %s\n", query)
        fmt.Println(strings.Repeat("=", 50))
        
        results, err := manager.SearchPackages(query)
        if err != nil {
            fmt.Printf("搜索失败: %v\n", err)
            continue
        }
        
        // 显示前5个结果
        for i, result := range results {
            if i >= 5 {
                break
            }
            fmt.Printf("%d. %s (%s)\n", i+1, result.Name, result.Version)
            fmt.Printf("   %s\n", result.Summary)
        }
        
        if len(results) > 5 {
            fmt.Printf("   ... 还有 %d 个结果\n", len(results)-5)
        }
    }
}
```

### 显示包详细信息

```go
func showPackageDetails() {
    manager := pip.NewManager(nil)
    
    packageNames := []string{"requests", "django", "numpy"}
    
    for _, name := range packageNames {
        fmt.Printf("\n%s 包信息:\n", name)
        fmt.Println(strings.Repeat("=", 50))
        
        info, err := manager.ShowPackage(name)
        if err != nil {
            fmt.Printf("获取 %s 信息失败: %v\n", name, err)
            continue
        }
        
        fmt.Printf("名称: %s\n", info.Name)
        fmt.Printf("版本: %s\n", info.Version)
        fmt.Printf("摘要: %s\n", info.Summary)
        fmt.Printf("作者: %s\n", info.Author)
        fmt.Printf("许可证: %s\n", info.License)
        fmt.Printf("主页: %s\n", info.Homepage)
        
        if len(info.Requires) > 0 {
            fmt.Printf("依赖: %s\n", strings.Join(info.Requires, ", "))
        }
        
        if len(info.RequiredBy) > 0 {
            fmt.Printf("被依赖: %s\n", strings.Join(info.RequiredBy, ", "))
        }
    }
}
```

## 包更新和维护

### 检查过时的包

```go
func checkOutdatedPackages() {
    manager := pip.NewManager(nil)
    
    fmt.Println("检查过时的包...")
    outdated, err := manager.CheckOutdated()
    if err != nil {
        log.Fatalf("检查过时包失败: %v", err)
    }
    
    if len(outdated) == 0 {
        fmt.Println("✅ 所有包都是最新的")
        return
    }
    
    fmt.Printf("发现 %d 个过时的包:\n", len(outdated))
    fmt.Println("包名\t\t当前版本\t最新版本")
    fmt.Println("----\t\t--------\t--------")
    
    for _, pkg := range outdated {
        fmt.Printf("%-20s %-15s %s\n", pkg.Name, pkg.CurrentVersion, pkg.LatestVersion)
    }
}
```

### 更新包

```go
func updatePackages() {
    manager := pip.NewManager(nil)
    
    // 更新单个包
    fmt.Println("更新 requests 包...")
    if err := manager.UpdatePackage("requests"); err != nil {
        fmt.Printf("更新 requests 失败: %v\n", err)
    } else {
        fmt.Println("✅ requests 更新成功")
    }
    
    // 更新所有包
    fmt.Println("更新所有包...")
    if err := manager.UpdateAllPackages(); err != nil {
        fmt.Printf("更新所有包失败: %v\n", err)
    } else {
        fmt.Println("✅ 所有包更新成功")
    }
}
```

### 卸载包

```go
func uninstallPackages() {
    manager := pip.NewManager(nil)
    
    packagesToRemove := []string{
        "old-package",
        "unused-dependency",
        "test-package",
    }
    
    for _, name := range packagesToRemove {
        fmt.Printf("卸载 %s...\n", name)
        if err := manager.UninstallPackage(name); err != nil {
            fmt.Printf("❌ 卸载 %s 失败: %v\n", name, err)
        } else {
            fmt.Printf("✅ %s 卸载成功\n", name)
        }
    }
}
```

## 需求文件管理

### 从需求文件安装

```go
func installFromRequirements() {
    manager := pip.NewManager(nil)
    
    requirementFiles := []string{
        "requirements.txt",
        "dev-requirements.txt",
        "test-requirements.txt",
    }
    
    for _, file := range requirementFiles {
        if _, err := os.Stat(file); os.IsNotExist(err) {
            fmt.Printf("⚠️  文件 %s 不存在，跳过\n", file)
            continue
        }
        
        fmt.Printf("从 %s 安装依赖...\n", file)
        if err := manager.InstallRequirements(file); err != nil {
            fmt.Printf("❌ 从 %s 安装失败: %v\n", file, err)
        } else {
            fmt.Printf("✅ 从 %s 安装成功\n", file)
        }
    }
}
```

### 生成需求文件

```go
func generateRequirements() {
    manager := pip.NewManager(nil)
    
    // 生成基本需求文件
    fmt.Println("生成 requirements.txt...")
    if err := manager.GenerateRequirements("requirements.txt"); err != nil {
        log.Fatalf("生成需求文件失败: %v", err)
    }
    
    // 生成冻结的需求文件（精确版本）
    fmt.Println("生成 requirements-freeze.txt...")
    packages, err := manager.FreezePackages()
    if err != nil {
        log.Fatalf("冻结包列表失败: %v", err)
    }
    
    file, err := os.Create("requirements-freeze.txt")
    if err != nil {
        log.Fatalf("创建文件失败: %v", err)
    }
    defer file.Close()
    
    for _, pkg := range packages {
        fmt.Fprintf(file, "%s==%s\n", pkg.Name, pkg.Version)
    }
    
    fmt.Println("✅ 需求文件生成完成")
}
```

## 依赖分析

### 分析包依赖

```go
func analyzeDependencies() {
    manager := pip.NewManager(nil)
    
    packageName := "django"
    fmt.Printf("分析 %s 的依赖关系...\n", packageName)
    
    deps, err := manager.GetDependencies(packageName)
    if err != nil {
        log.Fatalf("获取依赖失败: %v", err)
    }
    
    printDependencyTree(deps, 0)
}

func printDependencyTree(tree *pip.DependencyTree, indent int) {
    prefix := strings.Repeat("  ", indent)
    fmt.Printf("%s- %s %s\n", prefix, tree.Package.Name, tree.Package.Version)
    
    for _, dep := range tree.Dependencies {
        printDependencyTree(dep, indent+1)
    }
}
```

### 检查依赖冲突

```go
func checkDependencyConflicts() {
    manager := pip.NewManager(nil)
    
    fmt.Println("检查依赖冲突...")
    conflicts, err := manager.CheckDependencies()
    if err != nil {
        log.Fatalf("检查依赖失败: %v", err)
    }
    
    if len(conflicts) == 0 {
        fmt.Println("✅ 没有发现依赖冲突")
        return
    }
    
    fmt.Printf("发现 %d 个依赖冲突:\n", len(conflicts))
    for _, conflict := range conflicts {
        fmt.Printf("- %s: %s\n", conflict.Package, conflict.Description)
        fmt.Printf("  需要: %v\n", conflict.Required)
        fmt.Printf("  已安装: %s\n", conflict.Installed)
        fmt.Printf("  冲突: %v\n", conflict.Conflicts)
    }
}
```

## 错误处理和重试

### 健壮的包安装

```go
func robustPackageInstall() {
    manager := pip.NewManager(nil)
    
    packages := []*pip.PackageSpec{
        {Name: "requests"},
        {Name: "nonexistent-package-12345"}, // 故意的错误包
        {Name: "click"},
    }
    
    for _, pkg := range packages {
        if err := installWithRetry(manager, pkg, 3); err != nil {
            fmt.Printf("❌ %s 最终安装失败: %v\n", pkg.Name, err)
        } else {
            fmt.Printf("✅ %s 安装成功\n", pkg.Name)
        }
    }
}

func installWithRetry(manager *pip.Manager, pkg *pip.PackageSpec, maxRetries int) error {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        err := manager.InstallPackage(pkg)
        if err == nil {
            return nil // 成功
        }
        
        lastErr = err
        
        // 检查是否为可重试的错误
        if pip.IsErrorType(err, pip.ErrorTypeNetworkError) ||
           pip.IsErrorType(err, pip.ErrorTypeNetworkTimeout) {
            
            fmt.Printf("网络错误，重试 %d/%d...\n", i+1, maxRetries)
            time.Sleep(time.Duration(i+1) * time.Second) // 指数退避
            continue
        }
        
        // 非网络错误，不重试
        return err
    }
    
    return fmt.Errorf("重试 %d 次后仍然失败: %w", maxRetries, lastErr)
}
```

## 完整示例：包管理工具

```go
package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"
    "time"
    
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    manager := pip.NewManager(&pip.Config{
        LogLevel: "INFO",
        Timeout:  60 * time.Second,
    })
    
    for {
        fmt.Println("\n=== Python 包管理工具 ===")
        fmt.Println("1. 安装包")
        fmt.Println("2. 卸载包")
        fmt.Println("3. 列出已安装的包")
        fmt.Println("4. 搜索包")
        fmt.Println("5. 显示包信息")
        fmt.Println("6. 检查过时的包")
        fmt.Println("7. 更新包")
        fmt.Println("8. 生成需求文件")
        fmt.Println("9. 退出")
        fmt.Print("请选择操作 (1-9): ")
        
        reader := bufio.NewReader(os.Stdin)
        choice, _ := reader.ReadString('\n')
        choice = strings.TrimSpace(choice)
        
        switch choice {
        case "1":
            installPackageInteractive(manager, reader)
        case "2":
            uninstallPackageInteractive(manager, reader)
        case "3":
            listPackagesInteractive(manager)
        case "4":
            searchPackagesInteractive(manager, reader)
        case "5":
            showPackageInfoInteractive(manager, reader)
        case "6":
            checkOutdatedInteractive(manager)
        case "7":
            updatePackageInteractive(manager, reader)
        case "8":
            generateRequirementsInteractive(manager)
        case "9":
            fmt.Println("再见！")
            return
        default:
            fmt.Println("无效选择，请重试")
        }
    }
}

func installPackageInteractive(manager *pip.Manager, reader *bufio.Reader) {
    fmt.Print("请输入包名: ")
    packageName, _ := reader.ReadString('\n')
    packageName = strings.TrimSpace(packageName)
    
    fmt.Print("请输入版本约束 (可选，直接回车跳过): ")
    version, _ := reader.ReadString('\n')
    version = strings.TrimSpace(version)
    
    pkg := &pip.PackageSpec{
        Name:    packageName,
        Version: version,
    }
    
    fmt.Printf("正在安装 %s...\n", packageName)
    if err := manager.InstallPackage(pkg); err != nil {
        fmt.Printf("❌ 安装失败: %v\n", err)
    } else {
        fmt.Printf("✅ %s 安装成功\n", packageName)
    }
}

// ... 其他交互函数的实现
```

## 下一步

- 学习[虚拟环境示例](./virtual-environments.md)
- 查看[项目初始化示例](./project-initialization.md)
- 探索[高级用法示例](./advanced-usage.md)
- 了解[基本用法](./basic-usage.md)
