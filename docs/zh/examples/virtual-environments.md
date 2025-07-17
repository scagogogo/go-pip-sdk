# 虚拟环境示例

本页面提供了详细的虚拟环境管理示例，展示如何使用 Go Pip SDK 创建、管理和使用 Python 虚拟环境。

## 基本虚拟环境操作

### 创建和激活虚拟环境

```go
package main

import (
    "fmt"
    "log"
    "path/filepath"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func basicVenvOperations() {
    manager := pip.NewManager(nil)
    venvPath := filepath.Join(".", "my-project-env")
    
    // 创建虚拟环境
    fmt.Println("创建虚拟环境...")
    if err := manager.CreateVenv(venvPath); err != nil {
        log.Fatalf("创建虚拟环境失败: %v", err)
    }
    fmt.Println("✅ 虚拟环境创建成功")
    
    // 激活虚拟环境
    fmt.Println("激活虚拟环境...")
    if err := manager.ActivateVenv(venvPath); err != nil {
        log.Fatalf("激活虚拟环境失败: %v", err)
    }
    fmt.Println("✅ 虚拟环境已激活")
    
    // 检查虚拟环境状态
    isActive, activePath := manager.IsVenvActive()
    if isActive {
        fmt.Printf("当前虚拟环境: %s\n", activePath)
    }
    
    // 在虚拟环境中安装包
    pkg := &pip.PackageSpec{Name: "requests"}
    if err := manager.InstallPackage(pkg); err != nil {
        log.Printf("安装包失败: %v", err)
    } else {
        fmt.Println("✅ 在虚拟环境中安装包成功")
    }
    
    // 停用虚拟环境
    fmt.Println("停用虚拟环境...")
    if err := manager.DeactivateVenv(); err != nil {
        log.Printf("停用虚拟环境失败: %v", err)
    } else {
        fmt.Println("✅ 虚拟环境已停用")
    }
}
```

### 使用自定义选项创建虚拟环境

```go
func createVenvWithOptions() {
    manager := pip.NewManager(nil)
    
    // 高级虚拟环境选项
    opts := &pip.VenvOptions{
        PythonVersion:      "3.9",
        SystemSitePackages: false,
        Prompt:             "my-project",
        UpgradePip:         true,
        InstallSetuptools:  true,
        InstallWheel:       true,
        Clear:              false,
    }
    
    venvPath := "./advanced-env"
    fmt.Println("使用自定义选项创建虚拟环境...")
    
    if err := manager.CreateVenvWithOptions(venvPath, opts); err != nil {
        log.Fatalf("创建虚拟环境失败: %v", err)
    }
    
    fmt.Println("✅ 高级虚拟环境创建成功")
    
    // 获取虚拟环境信息
    info, err := manager.GetVenvInfo(venvPath)
    if err != nil {
        log.Printf("获取虚拟环境信息失败: %v", err)
        return
    }
    
    fmt.Printf("环境信息:\n")
    fmt.Printf("  名称: %s\n", info.Name)
    fmt.Printf("  Python 版本: %s\n", info.PythonVersion)
    fmt.Printf("  Pip 版本: %s\n", info.PipVersion)
    fmt.Printf("  创建时间: %s\n", info.CreatedAt.Format("2006-01-02 15:04:05"))
}
```

## 虚拟环境管理

### 列出和管理多个虚拟环境

```go
func manageMultipleVenvs() {
    manager := pip.NewManager(nil)
    
    // 创建多个虚拟环境
    envConfigs := []struct {
        name    string
        python  string
        purpose string
    }{
        {"web-dev", "3.9", "Web 开发环境"},
        {"data-science", "3.8", "数据科学环境"},
        {"testing", "3.10", "测试环境"},
    }
    
    fmt.Println("创建多个虚拟环境...")
    for _, config := range envConfigs {
        venvPath := filepath.Join("./envs", config.name)
        
        opts := &pip.VenvOptions{
            PythonVersion: config.python,
            Prompt:        config.name,
            UpgradePip:    true,
        }
        
        fmt.Printf("创建 %s (%s)...\n", config.name, config.purpose)
        if err := manager.CreateVenvWithOptions(venvPath, opts); err != nil {
            fmt.Printf("❌ 创建 %s 失败: %v\n", config.name, err)
            continue
        }
        fmt.Printf("✅ %s 创建成功\n", config.name)
    }
    
    // 列出所有虚拟环境
    fmt.Println("\n列出所有虚拟环境:")
    venvs, err := manager.ListVenvs("./envs")
    if err != nil {
        log.Printf("列出虚拟环境失败: %v", err)
        return
    }
    
    fmt.Printf("找到 %d 个虚拟环境:\n", len(venvs))
    for _, venv := range venvs {
        fmt.Printf("- %s (Python %s, %d 个包)\n", 
            venv.Name, venv.PythonVersion, venv.PackageCount)
    }
}
```

### 虚拟环境信息查询

```go
func queryVenvInfo() {
    manager := pip.NewManager(nil)
    
    venvPaths := []string{
        "./envs/web-dev",
        "./envs/data-science",
        "./envs/testing",
    }
    
    for _, path := range venvPaths {
        fmt.Printf("\n=== %s ===\n", filepath.Base(path))
        
        info, err := manager.GetVenvInfo(path)
        if err != nil {
            fmt.Printf("❌ 获取信息失败: %v\n", err)
            continue
        }
        
        fmt.Printf("名称: %s\n", info.Name)
        fmt.Printf("路径: %s\n", info.Path)
        fmt.Printf("Python 版本: %s\n", info.PythonVersion)
        fmt.Printf("Pip 版本: %s\n", info.PipVersion)
        fmt.Printf("包数量: %d\n", info.PackageCount)
        fmt.Printf("环境大小: %.2f MB\n", float64(info.Size)/(1024*1024))
        fmt.Printf("创建时间: %s\n", info.CreatedAt.Format("2006-01-02 15:04:05"))
        fmt.Printf("最后使用: %s\n", info.LastUsed.Format("2006-01-02 15:04:05"))
        fmt.Printf("是否激活: %t\n", info.IsActive)
    }
}
```

## 项目特定的虚拟环境

### Web 开发环境

```go
func setupWebDevEnvironment() {
    manager := pip.NewManager(nil)
    venvPath := "./web-project-env"
    
    // 创建 Web 开发虚拟环境
    opts := &pip.VenvOptions{
        PythonVersion: "3.9",
        Prompt:        "web-project",
        UpgradePip:    true,
    }
    
    fmt.Println("设置 Web 开发环境...")
    if err := manager.CreateVenvWithOptions(venvPath, opts); err != nil {
        log.Fatalf("创建环境失败: %v", err)
    }
    
    if err := manager.ActivateVenv(venvPath); err != nil {
        log.Fatalf("激活环境失败: %v", err)
    }
    
    // 安装 Web 开发相关包
    webPackages := []*pip.PackageSpec{
        {Name: "fastapi", Version: ">=0.68.0"},
        {Name: "uvicorn", Version: ">=0.15.0"},
        {Name: "pydantic", Version: ">=1.8.0"},
        {Name: "sqlalchemy", Version: ">=1.4.0"},
        {Name: "alembic", Version: ">=1.7.0"},
        {Name: "pytest", Version: ">=6.0"},
        {Name: "black", Version: ">=21.0"},
        {Name: "flake8", Version: ">=3.8"},
    }
    
    fmt.Println("安装 Web 开发包...")
    for _, pkg := range webPackages {
        fmt.Printf("安装 %s...\n", pkg.Name)
        if err := manager.InstallPackage(pkg); err != nil {
            fmt.Printf("❌ 安装 %s 失败: %v\n", pkg.Name, err)
        } else {
            fmt.Printf("✅ %s 安装成功\n", pkg.Name)
        }
    }
    
    // 生成需求文件
    if err := manager.GenerateRequirements("web-requirements.txt"); err != nil {
        log.Printf("生成需求文件失败: %v", err)
    } else {
        fmt.Println("✅ Web 开发需求文件生成成功")
    }
}
```

### 数据科学环境

```go
func setupDataScienceEnvironment() {
    manager := pip.NewManager(nil)
    venvPath := "./data-science-env"
    
    opts := &pip.VenvOptions{
        PythonVersion: "3.8",
        Prompt:        "data-science",
        UpgradePip:    true,
    }
    
    fmt.Println("设置数据科学环境...")
    if err := manager.CreateVenvWithOptions(venvPath, opts); err != nil {
        log.Fatalf("创建环境失败: %v", err)
    }
    
    if err := manager.ActivateVenv(venvPath); err != nil {
        log.Fatalf("激活环境失败: %v", err)
    }
    
    // 安装数据科学相关包
    dataPackages := []*pip.PackageSpec{
        {Name: "numpy", Version: ">=1.20.0"},
        {Name: "pandas", Version: ">=1.3.0"},
        {Name: "matplotlib", Version: ">=3.4.0"},
        {Name: "seaborn", Version: ">=0.11.0"},
        {Name: "scikit-learn", Version: ">=1.0.0"},
        {Name: "jupyter", Version: ">=1.0.0"},
        {Name: "ipython", Version: ">=7.0.0"},
        {Name: "plotly", Version: ">=5.0.0"},
    }
    
    fmt.Println("安装数据科学包...")
    for _, pkg := range dataPackages {
        fmt.Printf("安装 %s...\n", pkg.Name)
        if err := manager.InstallPackage(pkg); err != nil {
            fmt.Printf("❌ 安装 %s 失败: %v\n", pkg.Name, err)
        } else {
            fmt.Printf("✅ %s 安装成功\n", pkg.Name)
        }
    }
    
    fmt.Println("✅ 数据科学环境设置完成")
}
```

## 虚拟环境克隆和备份

### 克隆虚拟环境

```go
func cloneVirtualEnvironment() {
    manager := pip.NewManager(nil)
    
    sourcePath := "./production-env"
    targetPath := "./staging-env"
    
    fmt.Printf("克隆虚拟环境从 %s 到 %s...\n", sourcePath, targetPath)
    
    if err := manager.CloneVenv(sourcePath, targetPath); err != nil {
        log.Fatalf("克隆虚拟环境失败: %v", err)
    }
    
    fmt.Println("✅ 虚拟环境克隆成功")
    
    // 验证克隆结果
    sourceInfo, _ := manager.GetVenvInfo(sourcePath)
    targetInfo, _ := manager.GetVenvInfo(targetPath)
    
    fmt.Printf("源环境包数量: %d\n", sourceInfo.PackageCount)
    fmt.Printf("目标环境包数量: %d\n", targetInfo.PackageCount)
}
```

### 导出和导入虚拟环境

```go
func exportImportVenv() {
    manager := pip.NewManager(nil)
    
    // 导出虚拟环境配置
    venvPath := "./my-project-env"
    configFile := "environment.yml"
    
    fmt.Printf("导出虚拟环境 %s 到 %s...\n", venvPath, configFile)
    if err := manager.ExportVenv(venvPath, configFile); err != nil {
        log.Fatalf("导出虚拟环境失败: %v", err)
    }
    fmt.Println("✅ 虚拟环境导出成功")
    
    // 从配置文件导入虚拟环境
    newVenvPath := "./imported-env"
    fmt.Printf("从 %s 导入虚拟环境到 %s...\n", configFile, newVenvPath)
    if err := manager.ImportVenv(configFile, newVenvPath); err != nil {
        log.Fatalf("导入虚拟环境失败: %v", err)
    }
    fmt.Println("✅ 虚拟环境导入成功")
}
```

## 虚拟环境维护

### 清理和更新虚拟环境

```go
func maintainVirtualEnvironments() {
    manager := pip.NewManager(nil)
    
    venvPaths := []string{
        "./web-project-env",
        "./data-science-env",
        "./testing-env",
    }
    
    for _, path := range venvPaths {
        fmt.Printf("\n维护虚拟环境: %s\n", path)
        
        // 验证虚拟环境
        issues, err := manager.ValidateVenv(path)
        if err != nil {
            fmt.Printf("❌ 验证失败: %v\n", err)
            continue
        }
        
        if len(issues) > 0 {
            fmt.Printf("发现 %d 个问题:\n", len(issues))
            for _, issue := range issues {
                fmt.Printf("- %s: %s\n", issue.Type, issue.Message)
            }
            
            // 尝试修复
            fmt.Println("尝试修复虚拟环境...")
            if err := manager.RepairVenv(path); err != nil {
                fmt.Printf("❌ 修复失败: %v\n", err)
            } else {
                fmt.Println("✅ 修复成功")
            }
        } else {
            fmt.Println("✅ 虚拟环境状态良好")
        }
        
        // 清理缓存
        fmt.Println("清理虚拟环境缓存...")
        if err := manager.CleanVenv(path); err != nil {
            fmt.Printf("❌ 清理失败: %v\n", err)
        } else {
            fmt.Println("✅ 缓存清理完成")
        }
        
        // 更新虚拟环境
        fmt.Println("更新虚拟环境...")
        if err := manager.UpdateVenv(path); err != nil {
            fmt.Printf("❌ 更新失败: %v\n", err)
        } else {
            fmt.Println("✅ 更新完成")
        }
    }
}
```

### 清理旧的虚拟环境

```go
func cleanupOldEnvironments() {
    manager := pip.NewManager(nil)
    baseDir := "./envs"
    
    fmt.Println("查找旧的虚拟环境...")
    venvs, err := manager.ListVenvs(baseDir)
    if err != nil {
        log.Fatalf("列出虚拟环境失败: %v", err)
    }
    
    cutoff := time.Now().AddDate(0, 0, -30) // 30 天前
    var toDelete []string
    
    for _, venv := range venvs {
        if venv.LastUsed.Before(cutoff) {
            toDelete = append(toDelete, venv.Path)
        }
    }
    
    if len(toDelete) == 0 {
        fmt.Println("没有找到需要清理的旧环境")
        return
    }
    
    fmt.Printf("找到 %d 个超过 30 天未使用的环境:\n", len(toDelete))
    for _, path := range toDelete {
        fmt.Printf("- %s\n", path)
    }
    
    fmt.Print("是否删除这些环境? (y/N): ")
    var response string
    fmt.Scanln(&response)
    
    if strings.ToLower(response) == "y" {
        for _, path := range toDelete {
            fmt.Printf("删除 %s...\n", path)
            if err := manager.RemoveVenv(path); err != nil {
                fmt.Printf("❌ 删除失败: %v\n", err)
            } else {
                fmt.Printf("✅ 删除成功\n")
            }
        }
    }
}
```

## 批量虚拟环境操作

### 批量创建虚拟环境

```go
func batchCreateVenvs() {
    manager := pip.NewManager(nil)
    
    configs := []*pip.VenvConfig{
        {Path: "./envs/python38", PythonVersion: "3.8"},
        {Path: "./envs/python39", PythonVersion: "3.9"},
        {Path: "./envs/python310", PythonVersion: "3.10"},
    }
    
    fmt.Printf("批量创建 %d 个虚拟环境...\n", len(configs))
    
    if err := manager.CreateMultipleVenvs(configs); err != nil {
        log.Fatalf("批量创建失败: %v", err)
    }
    
    fmt.Println("✅ 批量创建完成")
}
```

### 批量删除虚拟环境

```go
func batchRemoveVenvs() {
    manager := pip.NewManager(nil)
    
    paths := []string{
        "./old-env1",
        "./old-env2",
        "./old-env3",
    }
    
    fmt.Printf("批量删除 %d 个虚拟环境...\n", len(paths))
    
    if err := manager.RemoveMultipleVenvs(paths); err != nil {
        log.Printf("批量删除失败: %v", err)
    } else {
        fmt.Println("✅ 批量删除完成")
    }
}
```

## 完整示例：虚拟环境管理器

```go
package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"
    "time"
    
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    manager := pip.NewManager(nil)
    
    for {
        fmt.Println("\n=== 虚拟环境管理器 ===")
        fmt.Println("1. 创建虚拟环境")
        fmt.Println("2. 列出虚拟环境")
        fmt.Println("3. 激活虚拟环境")
        fmt.Println("4. 停用虚拟环境")
        fmt.Println("5. 删除虚拟环境")
        fmt.Println("6. 克隆虚拟环境")
        fmt.Println("7. 清理旧环境")
        fmt.Println("8. 退出")
        fmt.Print("请选择操作 (1-8): ")
        
        reader := bufio.NewReader(os.Stdin)
        choice, _ := reader.ReadString('\n')
        choice = strings.TrimSpace(choice)
        
        switch choice {
        case "1":
            createVenvInteractive(manager, reader)
        case "2":
            listVenvsInteractive(manager)
        case "3":
            activateVenvInteractive(manager, reader)
        case "4":
            deactivateVenvInteractive(manager)
        case "5":
            removeVenvInteractive(manager, reader)
        case "6":
            cloneVenvInteractive(manager, reader)
        case "7":
            cleanupOldVenvsInteractive(manager)
        case "8":
            fmt.Println("再见！")
            return
        default:
            fmt.Println("无效选择，请重试")
        }
    }
}

func createVenvInteractive(manager *pip.Manager, reader *bufio.Reader) {
    fmt.Print("请输入虚拟环境路径: ")
    path, _ := reader.ReadString('\n')
    path = strings.TrimSpace(path)
    
    fmt.Print("请输入 Python 版本 (可选): ")
    version, _ := reader.ReadString('\n')
    version = strings.TrimSpace(version)
    
    fmt.Print("请输入提示符名称 (可选): ")
    prompt, _ := reader.ReadString('\n')
    prompt = strings.TrimSpace(prompt)
    
    if version == "" && prompt == "" {
        // 简单创建
        if err := manager.CreateVenv(path); err != nil {
            fmt.Printf("❌ 创建失败: %v\n", err)
        } else {
            fmt.Println("✅ 虚拟环境创建成功")
        }
    } else {
        // 使用选项创建
        opts := &pip.VenvOptions{
            UpgradePip: true,
        }
        if version != "" {
            opts.PythonVersion = version
        }
        if prompt != "" {
            opts.Prompt = prompt
        }
        
        if err := manager.CreateVenvWithOptions(path, opts); err != nil {
            fmt.Printf("❌ 创建失败: %v\n", err)
        } else {
            fmt.Println("✅ 虚拟环境创建成功")
        }
    }
}

// ... 其他交互函数的实现
```

## 下一步

- 学习[项目初始化示例](./project-initialization.md)
- 查看[包管理示例](./package-management.md)
- 探索[高级用法示例](./advanced-usage.md)
- 了解[基本用法](./basic-usage.md)
