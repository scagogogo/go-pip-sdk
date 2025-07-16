# 示例

本节提供使用 Go Pip SDK 执行各种任务的实际示例。每个示例都包含完整的、可运行的代码和解释。

## 快速示例

### 基本包安装

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    manager := pip.NewManager(nil)
    
    pkg := &pip.PackageSpec{
        Name:    "requests",
        Version: ">=2.25.0",
    }
    
    if err := manager.InstallPackage(pkg); err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("包安装成功！")
}
```

### 虚拟环境设置

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
        log.Fatal(err)
    }
    
    // 激活虚拟环境
    if err := manager.ActivateVenv(venvPath); err != nil {
        log.Fatal(err)
    }
    
    // 在虚拟环境中安装包
    packages := []*pip.PackageSpec{
        {Name: "requests"},
        {Name: "click"},
        {Name: "pydantic"},
    }
    
    for _, pkg := range packages {
        if err := manager.InstallPackage(pkg); err != nil {
            fmt.Printf("安装 %s 失败: %v\n", pkg.Name, err)
        } else {
            fmt.Printf("已安装 %s\n", pkg.Name)
        }
    }
    
    fmt.Println("虚拟环境设置完成！")
}
```

### 项目初始化

```go
package main

import (
    "log"
    
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func main() {
    manager := pip.NewManager(nil)
    
    opts := &pip.ProjectOptions{
        Name:            "my-awesome-project",
        Version:         "0.1.0",
        Description:     "一个很棒的 Python 项目",
        Author:          "您的姓名",
        AuthorEmail:     "your.email@example.com",
        License:         "MIT",
        Dependencies:    []string{"requests>=2.25.0", "click>=7.0"},
        DevDependencies: []string{"pytest>=6.0", "black>=21.0"},
        CreateVenv:      true,
        VenvPath:        "./venv",
    }
    
    if err := manager.InitProject("./my-project", opts); err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("项目初始化成功！")
}
```

## 详细示例

### [基本用法](/zh/examples/basic-usage)
通过简单的包管理操作学习 SDK 的基础知识。

### [包管理](/zh/examples/package-management)
高级包安装、版本管理和依赖处理。

### [虚拟环境](/zh/examples/virtual-environments)
创建和管理 Python 虚拟环境的完整指南。

### [项目初始化](/zh/examples/project-initialization)
使用适当的结构和配置引导新的 Python 项目。

### [高级用法](/zh/examples/advanced-usage)
复杂场景，包括错误处理、日志记录和自定义配置。

## 常见模式

### 错误处理模式

```go
func handlePipError(err error) {
    if err == nil {
        return
    }

    pipErr, ok := err.(*pip.PipErrorDetails)
    if !ok {
        fmt.Printf("意外错误: %v\n", err)
        return
    }

    fmt.Printf("错误类型: %s\n", pipErr.Type)
    fmt.Printf("消息: %s\n", pipErr.Message)

    if pipErr.Command != "" {
        fmt.Printf("失败的命令: %s\n", pipErr.Command)
    }

    if pipErr.ExitCode != 0 {
        fmt.Printf("退出代码: %d\n", pipErr.ExitCode)
    }

    if len(pipErr.Suggestions) > 0 {
        fmt.Println("建议:")
        for _, suggestion := range pipErr.Suggestions {
            fmt.Printf("  - %s\n", suggestion)
        }
    }

    if pipErr.Cause != nil {
        fmt.Printf("根本原因: %v\n", pipErr.Cause)
    }
}
```

### 重试模式

```go
func installWithRetry(manager *pip.Manager, pkg *pip.PackageSpec, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        err := manager.InstallPackage(pkg)
        if err == nil {
            return nil // 成功
        }
        
        lastErr := err
        
        // 检查错误是否可恢复
        if pip.IsErrorType(err, pip.ErrorTypeNetworkError) {
            fmt.Printf("网络错误，正在重试... (%d/%d)\n", i+1, maxRetries)
            time.Sleep(time.Duration(i+1) * time.Second)
            continue
        }
        
        return lastErr // 不可恢复的错误
    }
    return fmt.Errorf("重试 %d 次后失败", maxRetries)
}
```

### 日志记录模式

```go
func setupLogging() (*pip.Logger, error) {
    return pip.NewLogger(&pip.LoggerConfig{
        Level:      pip.LogLevelInfo,
        Output:     os.Stdout,
        Prefix:     "[my-app]",
        EnableFile: true,
        LogFile:    "pip-operations.log",
    })
}

func main() {
    logger, err := setupLogging()
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Close()
    
    manager := pip.NewManager(nil)
    manager.SetCustomLogger(logger)
    
    // 现在所有操作都将被记录
    // ...
}
```

## 测试示例

### 使用虚拟环境进行单元测试

```go
func TestWithCleanEnvironment(t *testing.T) {
    manager := pip.NewManager(nil)
    
    // 创建临时虚拟环境
    tempDir, err := os.MkdirTemp("", "test-venv-*")
    require.NoError(t, err)
    defer os.RemoveAll(tempDir)
    
    venvPath := filepath.Join(tempDir, "venv")
    require.NoError(t, manager.CreateVenv(venvPath))
    require.NoError(t, manager.ActivateVenv(venvPath))
    
    // 在干净环境中测试包安装
    pkg := &pip.PackageSpec{Name: "requests"}
    require.NoError(t, manager.InstallPackage(pkg))
    
    // 验证安装
    packages, err := manager.ListPackages()
    require.NoError(t, err)
    
    found := false
    for _, p := range packages {
        if p.Name == "requests" {
            found = true
            break
        }
    }
    require.True(t, found, "应该安装 requests 包")
}
```

### 集成测试

```go
func TestFullWorkflow(t *testing.T) {
    manager := pip.NewManager(nil)
    
    // 检查 pip 安装
    installed, err := manager.IsInstalled()
    require.NoError(t, err)
    
    if !installed {
        require.NoError(t, manager.Install())
    }
    
    // 创建项目
    tempDir, err := os.MkdirTemp("", "test-project-*")
    require.NoError(t, err)
    defer os.RemoveAll(tempDir)
    
    opts := &pip.ProjectOptions{
        Name:         "test-project",
        Version:      "0.1.0",
        Author:       "测试作者",
        AuthorEmail:  "test@example.com",
        Dependencies: []string{"requests"},
        CreateVenv:   true,
    }
    
    require.NoError(t, manager.InitProject(tempDir, opts))
    
    // 验证项目结构
    files := []string{
        "setup.py",
        "pyproject.toml",
        "requirements.txt",
        "README.md",
        ".gitignore",
    }
    
    for _, file := range files {
        path := filepath.Join(tempDir, file)
        _, err := os.Stat(path)
        require.NoError(t, err, "文件 %s 应该存在", file)
    }
}
```

## 性能示例

### 并发包安装

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
        return fmt.Errorf("安装错误: %v", errors)
    }
    
    return nil
}
```

## 下一步

- 探索 [API 参考](/zh/api/) 获取详细文档
- 查看 [指南](/zh/guide/) 获取全面教程
- 浏览 [GitHub 仓库](https://github.com/scagogogo/go-pip-sdk) 中的源代码
