# 虚拟环境

虚拟环境是 Python 开发的最佳实践，它允许您为不同的项目创建隔离的 Python 环境。Go Pip SDK 提供了完整的虚拟环境管理功能。

## 虚拟环境基础

### 什么是虚拟环境？

虚拟环境是一个独立的 Python 环境，具有：
- 独立的 Python 解释器
- 独立的包安装目录
- 独立的环境变量
- 与系统 Python 环境隔离

### 为什么使用虚拟环境？

1. **依赖隔离**：不同项目可以使用不同版本的包
2. **避免冲突**：防止包版本冲突
3. **清洁环境**：保持系统 Python 环境干净
4. **可重现性**：确保项目在不同环境中的一致性

## 创建虚拟环境

### 基本创建

```go
package main

import (
    "fmt"
    "log"
    "path/filepath"
    
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    manager := pip.NewManager(nil)
    
    // 创建虚拟环境
    venvPath := filepath.Join(".", "my-venv")
    if err := manager.CreateVenv(venvPath); err != nil {
        log.Fatalf("创建虚拟环境失败: %v", err)
    }
    
    fmt.Printf("虚拟环境已创建: %s\n", venvPath)
}
```

### 指定 Python 版本

```go
// 使用特定 Python 版本创建虚拟环境
config := &pip.Config{
    PythonPath: "/usr/bin/python3.9",
}
manager := pip.NewManager(config)

venvPath := "./python39-env"
if err := manager.CreateVenv(venvPath); err != nil {
    log.Fatalf("创建 Python 3.9 虚拟环境失败: %v", err)
}
```

### 创建带系统包的虚拟环境

```go
// 创建可访问系统包的虚拟环境
venvOptions := &pip.VenvOptions{
    SystemSitePackages: true,
}

if err := manager.CreateVenvWithOptions(venvPath, venvOptions); err != nil {
    log.Fatalf("创建虚拟环境失败: %v", err)
}
```

## 激活和停用虚拟环境

### 激活虚拟环境

```go
venvPath := "./my-venv"

// 激活虚拟环境
if err := manager.ActivateVenv(venvPath); err != nil {
    log.Fatalf("激活虚拟环境失败: %v", err)
}

fmt.Println("虚拟环境已激活")

// 现在所有的 pip 操作都会在这个虚拟环境中执行
pkg := &pip.PackageSpec{Name: "requests"}
if err := manager.InstallPackage(pkg); err != nil {
    log.Fatalf("在虚拟环境中安装包失败: %v", err)
}
```

### 停用虚拟环境

```go
// 停用当前虚拟环境
if err := manager.DeactivateVenv(); err != nil {
    log.Fatalf("停用虚拟环境失败: %v", err)
}

fmt.Println("虚拟环境已停用")
```

### 检查虚拟环境状态

```go
// 检查是否在虚拟环境中
isActive, venvPath := manager.IsVenvActive()
if isActive {
    fmt.Printf("当前在虚拟环境中: %s\n", venvPath)
} else {
    fmt.Println("当前不在虚拟环境中")
}
```

## 虚拟环境管理

### 列出虚拟环境

```go
// 列出指定目录下的所有虚拟环境
venvs, err := manager.ListVenvs("./envs")
if err != nil {
    log.Fatalf("列出虚拟环境失败: %v", err)
}

fmt.Printf("找到 %d 个虚拟环境:\n", len(venvs))
for _, venv := range venvs {
    fmt.Printf("- %s (Python %s)\n", venv.Path, venv.PythonVersion)
}
```

### 获取虚拟环境信息

```go
info, err := manager.GetVenvInfo("./my-venv")
if err != nil {
    log.Fatalf("获取虚拟环境信息失败: %v", err)
}

fmt.Printf("虚拟环境信息:\n")
fmt.Printf("路径: %s\n", info.Path)
fmt.Printf("Python 版本: %s\n", info.PythonVersion)
fmt.Printf("Pip 版本: %s\n", info.PipVersion)
fmt.Printf("已安装包数量: %d\n", len(info.InstalledPackages))
```

### 删除虚拟环境

```go
// 删除虚拟环境
venvPath := "./old-venv"
if err := manager.RemoveVenv(venvPath); err != nil {
    log.Fatalf("删除虚拟环境失败: %v", err)
}

fmt.Printf("虚拟环境已删除: %s\n", venvPath)
```

## 完整的虚拟环境工作流

```go
package main

import (
    "fmt"
    "log"
    "path/filepath"
    
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    manager := pip.NewManager(nil)
    venvPath := filepath.Join(".", "project-env")
    
    // 1. 创建虚拟环境
    fmt.Println("创建虚拟环境...")
    if err := manager.CreateVenv(venvPath); err != nil {
        log.Fatalf("创建虚拟环境失败: %v", err)
    }
    
    // 2. 激活虚拟环境
    fmt.Println("激活虚拟环境...")
    if err := manager.ActivateVenv(venvPath); err != nil {
        log.Fatalf("激活虚拟环境失败: %v", err)
    }
    
    // 3. 安装包
    fmt.Println("安装项目依赖...")
    packages := []*pip.PackageSpec{
        {Name: "requests"},
        {Name: "click"},
        {Name: "pydantic"},
    }
    
    for _, pkg := range packages {
        if err := manager.InstallPackage(pkg); err != nil {
            fmt.Printf("安装 %s 失败: %v\n", pkg.Name, err)
            continue
        }
        fmt.Printf("✅ 安装 %s 成功\n", pkg.Name)
    }
    
    // 4. 生成需求文件
    fmt.Println("生成需求文件...")
    if err := manager.GenerateRequirements("requirements.txt"); err != nil {
        log.Printf("生成需求文件失败: %v", err)
    }
    
    // 5. 列出已安装的包
    fmt.Println("已安装的包:")
    packages, err := manager.ListPackages()
    if err != nil {
        log.Printf("列出包失败: %v", err)
    } else {
        for _, pkg := range packages {
            fmt.Printf("- %s %s\n", pkg.Name, pkg.Version)
        }
    }
    
    // 6. 停用虚拟环境（可选）
    fmt.Println("停用虚拟环境...")
    if err := manager.DeactivateVenv(); err != nil {
        log.Printf("停用虚拟环境失败: %v", err)
    }
    
    fmt.Println("虚拟环境工作流完成！")
}
```

## 高级虚拟环境操作

### 克隆虚拟环境

```go
// 克隆现有虚拟环境
sourceVenv := "./source-env"
targetVenv := "./target-env"

if err := manager.CloneVenv(sourceVenv, targetVenv); err != nil {
    log.Fatalf("克隆虚拟环境失败: %v", err)
}

fmt.Printf("虚拟环境已克隆: %s -> %s\n", sourceVenv, targetVenv)
```

### 比较虚拟环境

```go
// 比较两个虚拟环境的包差异
venv1 := "./env1"
venv2 := "./env2"

diff, err := manager.CompareVenvs(venv1, venv2)
if err != nil {
    log.Fatalf("比较虚拟环境失败: %v", err)
}

fmt.Printf("环境差异:\n")
fmt.Printf("仅在 %s 中: %v\n", venv1, diff.OnlyInFirst)
fmt.Printf("仅在 %s 中: %v\n", venv2, diff.OnlyInSecond)
fmt.Printf("版本不同: %v\n", diff.VersionDifferences)
```

### 导出和导入虚拟环境

```go
// 导出虚拟环境配置
venvPath := "./my-env"
exportPath := "./env-export.json"

if err := manager.ExportVenv(venvPath, exportPath); err != nil {
    log.Fatalf("导出虚拟环境失败: %v", err)
}

// 从配置文件重建虚拟环境
newVenvPath := "./restored-env"
if err := manager.ImportVenv(exportPath, newVenvPath); err != nil {
    log.Fatalf("导入虚拟环境失败: %v", err)
}
```

## 虚拟环境最佳实践

### 1. 命名约定

```go
// 使用描述性名称
projectVenv := "./myproject-env"
testVenv := "./myproject-test-env"
prodVenv := "./myproject-prod-env"
```

### 2. 目录结构

```
project/
├── src/
├── tests/
├── docs/
├── venv/          # 虚拟环境
├── requirements.txt
├── dev-requirements.txt
└── README.md
```

### 3. 自动化脚本

```go
func setupDevelopmentEnvironment(projectPath string) error {
    manager := pip.NewManager(nil)
    venvPath := filepath.Join(projectPath, "venv")
    
    // 创建虚拟环境
    if err := manager.CreateVenv(venvPath); err != nil {
        return fmt.Errorf("创建虚拟环境失败: %w", err)
    }
    
    // 激活虚拟环境
    if err := manager.ActivateVenv(venvPath); err != nil {
        return fmt.Errorf("激活虚拟环境失败: %w", err)
    }
    
    // 安装开发依赖
    reqFile := filepath.Join(projectPath, "dev-requirements.txt")
    if err := manager.InstallRequirements(reqFile); err != nil {
        return fmt.Errorf("安装开发依赖失败: %w", err)
    }
    
    return nil
}
```

### 4. 清理脚本

```go
func cleanupEnvironments(baseDir string) error {
    manager := pip.NewManager(nil)
    
    venvs, err := manager.ListVenvs(baseDir)
    if err != nil {
        return err
    }
    
    for _, venv := range venvs {
        // 删除超过 30 天未使用的虚拟环境
        if time.Since(venv.LastUsed) > 30*24*time.Hour {
            if err := manager.RemoveVenv(venv.Path); err != nil {
                fmt.Printf("删除 %s 失败: %v\n", venv.Path, err)
                continue
            }
            fmt.Printf("已删除旧虚拟环境: %s\n", venv.Path)
        }
    }
    
    return nil
}
```

## 故障排除

### 常见问题

1. **虚拟环境创建失败**
   - 检查 Python 是否正确安装
   - 确保有足够的磁盘空间
   - 检查目录权限

2. **激活失败**
   - 确保虚拟环境路径正确
   - 检查虚拟环境是否完整

3. **包安装失败**
   - 确保虚拟环境已激活
   - 检查网络连接
   - 验证包名称和版本

### 调试技巧

```go
// 启用详细日志
config := &pip.Config{
    LogLevel: "DEBUG",
}
manager := pip.NewManager(config)

// 检查虚拟环境状态
isActive, path := manager.IsVenvActive()
fmt.Printf("虚拟环境状态: 激活=%v, 路径=%s\n", isActive, path)

// 验证 Python 路径
pythonPath, err := manager.GetPythonPath()
if err != nil {
    fmt.Printf("获取 Python 路径失败: %v\n", err)
} else {
    fmt.Printf("当前 Python 路径: %s\n", pythonPath)
}
```

## 下一步

- 了解[项目管理](./project-management.md)
- 学习[错误处理](./error-handling.md)
- 查看[虚拟环境示例](/zh/examples/virtual-environments.md)
- 探索[API 参考](/zh/api/virtual-environments.md)
