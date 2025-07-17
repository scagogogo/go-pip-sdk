# 包管理

Go Pip SDK 提供了全面的 Python 包管理功能，包括安装、卸载、列出、搜索和更新包。

## 基本包操作

### 安装包

#### 安装单个包

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    manager := pip.NewManager(nil)
    
    // 基本包安装
    pkg := &pip.PackageSpec{
        Name: "requests",
    }
    
    if err := manager.InstallPackage(pkg); err != nil {
        log.Fatalf("安装包失败: %v", err)
    }
    
    fmt.Println("包安装成功！")
}
```

#### 指定版本安装

```go
// 安装特定版本
pkg := &pip.PackageSpec{
    Name:    "django",
    Version: "4.2.0",
}

// 使用版本约束
pkg := &pip.PackageSpec{
    Name:    "numpy",
    Version: ">=1.20.0,<2.0.0",
}

// 安装最新版本
pkg := &pip.PackageSpec{
    Name:    "flask",
    Version: "latest",
}
```

#### 安装多个包

```go
packages := []*pip.PackageSpec{
    {Name: "requests"},
    {Name: "click", Version: ">=7.0"},
    {Name: "pydantic", Version: "^1.8.0"},
}

for _, pkg := range packages {
    if err := manager.InstallPackage(pkg); err != nil {
        fmt.Printf("安装 %s 失败: %v\n", pkg.Name, err)
        continue
    }
    fmt.Printf("成功安装 %s\n", pkg.Name)
}
```

### 高级安装选项

#### 可编辑安装

```go
pkg := &pip.PackageSpec{
    Name:     "my-local-package",
    Editable: true,
    // 对于本地包，Name 字段应该是路径
}
```

#### 从特定索引安装

```go
pkg := &pip.PackageSpec{
    Name:  "private-package",
    Index: "https://private.pypi.server.com/simple/",
}
```

#### 安装额外依赖

```go
pkg := &pip.PackageSpec{
    Name:   "fastapi",
    Extras: []string{"dev", "test"},
}
```

#### 强制重新安装

```go
pkg := &pip.PackageSpec{
    Name:           "problematic-package",
    ForceReinstall: true,
}
```

#### 升级包

```go
pkg := &pip.PackageSpec{
    Name:    "outdated-package",
    Upgrade: true,
}
```

### 卸载包

```go
// 卸载单个包
if err := manager.UninstallPackage("requests"); err != nil {
    log.Fatalf("卸载包失败: %v", err)
}

// 卸载多个包
packages := []string{"requests", "click", "pydantic"}
for _, name := range packages {
    if err := manager.UninstallPackage(name); err != nil {
        fmt.Printf("卸载 %s 失败: %v\n", name, err)
        continue
    }
    fmt.Printf("成功卸载 %s\n", name)
}
```

## 包信息查询

### 列出已安装的包

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

### 显示包详细信息

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

### 搜索包

```go
results, err := manager.SearchPackages("web framework")
if err != nil {
    log.Fatalf("搜索包失败: %v", err)
}

fmt.Printf("找到 %d 个搜索结果:\n", len(results))
for _, result := range results {
    fmt.Printf("- %s: %s\n", result.Name, result.Summary)
}
```

### 冻结包列表

```go
// 获取当前环境的包列表（类似 pip freeze）
packages, err := manager.FreezePackages()
if err != nil {
    log.Fatalf("冻结包列表失败: %v", err)
}

fmt.Println("当前环境包列表:")
for _, pkg := range packages {
    fmt.Printf("%s==%s\n", pkg.Name, pkg.Version)
}
```

## 需求文件管理

### 从需求文件安装

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

### 生成需求文件

```go
// 生成 requirements.txt
if err := manager.GenerateRequirements("requirements.txt"); err != nil {
    log.Fatalf("生成需求文件失败: %v", err)
}

// 生成到自定义文件
if err := manager.GenerateRequirements("current-env.txt"); err != nil {
    log.Fatalf("生成环境文件失败: %v", err)
}
```

## 包管理最佳实践

### 1. 使用虚拟环境

```go
// 创建并激活虚拟环境
venvPath := "./my-project-env"
if err := manager.CreateVenv(venvPath); err != nil {
    log.Fatalf("创建虚拟环境失败: %v", err)
}

if err := manager.ActivateVenv(venvPath); err != nil {
    log.Fatalf("激活虚拟环境失败: %v", err)
}

// 在虚拟环境中安装包
pkg := &pip.PackageSpec{Name: "requests"}
if err := manager.InstallPackage(pkg); err != nil {
    log.Fatalf("安装包失败: %v", err)
}
```

### 2. 版本锁定

```go
// 开发依赖 - 使用灵活版本
devPackages := []*pip.PackageSpec{
    {Name: "pytest", Version: ">=6.0"},
    {Name: "black", Version: ">=21.0"},
    {Name: "flake8", Version: ">=3.8"},
}

// 生产依赖 - 使用精确版本
prodPackages := []*pip.PackageSpec{
    {Name: "django", Version: "4.2.7"},
    {Name: "psycopg2", Version: "2.9.5"},
    {Name: "redis", Version: "4.5.1"},
}
```

### 3. 错误处理

```go
func installPackageWithRetry(manager *pip.Manager, pkg *pip.PackageSpec, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        err := manager.InstallPackage(pkg)
        if err == nil {
            return nil
        }
        
        // 检查错误类型
        if pip.IsErrorType(err, pip.ErrorTypeNetworkError) {
            fmt.Printf("网络错误，重试 %d/%d...\n", i+1, maxRetries)
            time.Sleep(time.Duration(i+1) * time.Second)
            continue
        }
        
        // 非网络错误，直接返回
        return err
    }
    
    return fmt.Errorf("重试 %d 次后仍然失败", maxRetries)
}
```

### 4. 批量操作

```go
func installPackagesBatch(manager *pip.Manager, packages []*pip.PackageSpec) error {
    var errors []error
    
    for _, pkg := range packages {
        if err := manager.InstallPackage(pkg); err != nil {
            errors = append(errors, fmt.Errorf("安装 %s 失败: %w", pkg.Name, err))
            continue
        }
        fmt.Printf("✅ 成功安装 %s\n", pkg.Name)
    }
    
    if len(errors) > 0 {
        return fmt.Errorf("批量安装中有 %d 个包失败", len(errors))
    }
    
    return nil
}
```

## 性能优化

### 1. 使用缓存

```go
config := &pip.Config{
    CacheDir: "/tmp/pip-cache",
    ExtraOptions: map[string]string{
        "cache-dir": "/tmp/pip-cache",
    },
}
manager := pip.NewManager(config)
```

### 2. 并行安装

```go
func installPackagesConcurrently(packages []*pip.PackageSpec) error {
    var wg sync.WaitGroup
    errChan := make(chan error, len(packages))
    
    for _, pkg := range packages {
        wg.Add(1)
        go func(p *pip.PackageSpec) {
            defer wg.Done()
            
            manager := pip.NewManager(nil)
            if err := manager.InstallPackage(p); err != nil {
                errChan <- fmt.Errorf("安装 %s 失败: %w", p.Name, err)
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
        return fmt.Errorf("并发安装中有 %d 个包失败", len(errors))
    }
    
    return nil
}
```

## 下一步

- 学习[虚拟环境管理](./virtual-environments.md)
- 了解[项目管理](./project-management.md)
- 查看[错误处理](./error-handling.md)
- 探索[包管理示例](/zh/examples/package-management.md)
