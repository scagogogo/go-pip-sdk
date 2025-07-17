# 错误处理

Go Pip SDK 提供了全面的错误处理机制，帮助您识别、诊断和处理各种操作中可能出现的错误。

## 错误类型

### 基本错误接口

所有 SDK 错误都实现了标准的 Go `error` 接口：

```go
if err := manager.InstallPackage(pkg); err != nil {
    fmt.Printf("操作失败: %v\n", err)
}
```

### 结构化错误

SDK 提供了丰富的错误信息：

```go
if err := manager.InstallPackage(pkg); err != nil {
    if pipErr, ok := err.(*pip.PipErrorDetails); ok {
        fmt.Printf("错误类型: %s\n", pipErr.Type)
        fmt.Printf("错误消息: %s\n", pipErr.Message)
        fmt.Printf("失败命令: %s\n", pipErr.Command)
        fmt.Printf("退出代码: %d\n", pipErr.ExitCode)
        fmt.Printf("输出: %s\n", pipErr.Output)
        
        if len(pipErr.Suggestions) > 0 {
            fmt.Println("建议:")
            for _, suggestion := range pipErr.Suggestions {
                fmt.Printf("- %s\n", suggestion)
            }
        }
    }
}
```

## 错误分类

### 系统错误

#### Python 未找到

```go
_, err := manager.IsInstalled()
if err != nil {
    if pip.IsErrorType(err, pip.ErrorTypePythonNotFound) {
        fmt.Println("Python 未安装或未在 PATH 中")
        
        // 建议解决方案
        fmt.Println("解决方案:")
        fmt.Println("1. 安装 Python")
        fmt.Println("2. 将 Python 添加到 PATH")
        fmt.Println("3. 在配置中指定 Python 路径")
        
        // 使用自定义 Python 路径
        config := &pip.Config{
            PythonPath: "/usr/bin/python3",
        }
        manager = pip.NewManager(config)
    }
}
```

#### Pip 未安装

```go
installed, err := manager.IsInstalled()
if err != nil {
    if pip.IsErrorType(err, pip.ErrorTypePipNotFound) {
        fmt.Println("Pip 未安装")
        
        // 自动安装 pip
        if err := manager.Install(); err != nil {
            fmt.Printf("自动安装 pip 失败: %v\n", err)
        } else {
            fmt.Println("Pip 安装成功")
        }
    }
}
```

### 网络错误

#### 连接超时

```go
err := manager.InstallPackage(pkg)
if err != nil {
    if pip.IsErrorType(err, pip.ErrorTypeNetworkTimeout) {
        fmt.Println("网络连接超时")
        
        // 增加超时时间重试
        config := &pip.Config{
            Timeout: 300 * time.Second, // 5 分钟
        }
        retryManager := pip.NewManager(config)
        
        if err := retryManager.InstallPackage(pkg); err != nil {
            fmt.Printf("重试失败: %v\n", err)
        }
    }
}
```

#### 网络连接错误

```go
err := manager.InstallPackage(pkg)
if err != nil {
    if pip.IsErrorType(err, pip.ErrorTypeNetworkError) {
        fmt.Println("网络连接失败")
        
        // 尝试使用镜像源
        config := &pip.Config{
            DefaultIndex: "https://pypi.tuna.tsinghua.edu.cn/simple/",
            TrustedHosts: []string{"pypi.tuna.tsinghua.edu.cn"},
        }
        mirrorManager := pip.NewManager(config)
        
        if err := mirrorManager.InstallPackage(pkg); err != nil {
            fmt.Printf("使用镜像源仍然失败: %v\n", err)
        }
    }
}
```

### 包相关错误

#### 包未找到

```go
err := manager.InstallPackage(pkg)
if err != nil {
    if pip.IsErrorType(err, pip.ErrorTypePackageNotFound) {
        fmt.Printf("包 '%s' 未找到\n", pkg.Name)
        
        // 搜索相似包名
        results, searchErr := manager.SearchPackages(pkg.Name)
        if searchErr == nil && len(results) > 0 {
            fmt.Println("您是否想要安装以下包之一:")
            for i, result := range results[:5] { // 显示前5个结果
                fmt.Printf("%d. %s - %s\n", i+1, result.Name, result.Summary)
            }
        }
    }
}
```

#### 版本冲突

```go
err := manager.InstallPackage(pkg)
if err != nil {
    if pip.IsErrorType(err, pip.ErrorTypeVersionConflict) {
        fmt.Printf("版本冲突: %s\n", err.Error())
        
        // 显示当前已安装的版本
        if info, showErr := manager.ShowPackage(pkg.Name); showErr == nil {
            fmt.Printf("当前安装版本: %s\n", info.Version)
            fmt.Printf("请求版本: %s\n", pkg.Version)
        }
        
        // 建议解决方案
        fmt.Println("解决方案:")
        fmt.Println("1. 更新版本约束")
        fmt.Println("2. 使用虚拟环境")
        fmt.Println("3. 强制重新安装")
    }
}
```

### 权限错误

```go
err := manager.InstallPackage(pkg)
if err != nil {
    if pip.IsErrorType(err, pip.ErrorTypePermissionDenied) {
        fmt.Println("权限被拒绝")
        
        // 建议解决方案
        fmt.Println("解决方案:")
        fmt.Println("1. 使用虚拟环境 (推荐)")
        fmt.Println("2. 使用 --user 标志")
        fmt.Println("3. 以管理员身份运行")
        
        // 尝试用户级安装
        userPkg := &pip.PackageSpec{
            Name: pkg.Name,
            Options: map[string]string{
                "user": "",
            },
        }
        
        if err := manager.InstallPackage(userPkg); err != nil {
            fmt.Printf("用户级安装也失败: %v\n", err)
        }
    }
}
```

## 错误处理模式

### 基本错误处理

```go
func handleBasicError(err error) {
    if err == nil {
        return
    }
    
    fmt.Printf("操作失败: %v\n", err)
    
    // 记录错误
    log.Printf("Error: %v", err)
}
```

### 详细错误处理

```go
func handleDetailedError(err error) {
    if err == nil {
        return
    }
    
    pipErr, ok := err.(*pip.PipErrorDetails)
    if !ok {
        fmt.Printf("未知错误: %v\n", err)
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
            fmt.Printf("- %s\n", suggestion)
        }
    }
    
    if pipErr.Cause != nil {
        fmt.Printf("根本原因: %v\n", pipErr.Cause)
    }
}
```

### 重试机制

```go
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

### 错误恢复

```go
func installPackageWithRecovery(manager *pip.Manager, pkg *pip.PackageSpec) error {
    err := manager.InstallPackage(pkg)
    if err == nil {
        return nil
    }
    
    switch pip.GetErrorType(err) {
    case pip.ErrorTypePackageNotFound:
        // 尝试搜索相似包名
        return suggestAlternativePackages(manager, pkg.Name)
        
    case pip.ErrorTypeVersionConflict:
        // 尝试安装兼容版本
        return installCompatibleVersion(manager, pkg)
        
    case pip.ErrorTypePermissionDenied:
        // 尝试用户级安装
        return installAsUser(manager, pkg)
        
    case pip.ErrorTypeNetworkError:
        // 尝试使用镜像源
        return installFromMirror(manager, pkg)
        
    default:
        return err
    }
}

func suggestAlternativePackages(manager *pip.Manager, packageName string) error {
    results, err := manager.SearchPackages(packageName)
    if err != nil {
        return fmt.Errorf("搜索替代包失败: %w", err)
    }
    
    if len(results) == 0 {
        return fmt.Errorf("未找到包 '%s' 的替代方案", packageName)
    }
    
    fmt.Printf("未找到包 '%s'，但找到以下相似包:\n", packageName)
    for i, result := range results[:3] {
        fmt.Printf("%d. %s - %s\n", i+1, result.Name, result.Summary)
    }
    
    return fmt.Errorf("请检查包名是否正确")
}
```

## 错误日志记录

### 基本日志记录

```go
func logError(err error, operation string) {
    if err == nil {
        return
    }
    
    log.Printf("[ERROR] %s failed: %v", operation, err)
    
    if pipErr, ok := err.(*pip.PipErrorDetails); ok {
        log.Printf("[ERROR] Type: %s, Command: %s, ExitCode: %d", 
            pipErr.Type, pipErr.Command, pipErr.ExitCode)
    }
}
```

### 结构化日志记录

```go
import (
    "encoding/json"
    "log"
)

func logStructuredError(err error, operation string, context map[string]interface{}) {
    if err == nil {
        return
    }
    
    logEntry := map[string]interface{}{
        "level":     "error",
        "operation": operation,
        "error":     err.Error(),
        "timestamp": time.Now().UTC(),
        "context":   context,
    }
    
    if pipErr, ok := err.(*pip.PipErrorDetails); ok {
        logEntry["error_type"] = pipErr.Type
        logEntry["command"] = pipErr.Command
        logEntry["exit_code"] = pipErr.ExitCode
        logEntry["suggestions"] = pipErr.Suggestions
    }
    
    jsonLog, _ := json.Marshal(logEntry)
    log.Println(string(jsonLog))
}
```

## 错误监控和报告

### 错误统计

```go
type ErrorStats struct {
    Total       int
    ByType      map[pip.ErrorType]int
    ByOperation map[string]int
    mutex       sync.RWMutex
}

func (s *ErrorStats) RecordError(err error, operation string) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    s.Total++
    s.ByOperation[operation]++
    
    if pipErr, ok := err.(*pip.PipErrorDetails); ok {
        s.ByType[pipErr.Type]++
    }
}

func (s *ErrorStats) GetStats() map[string]interface{} {
    s.mutex.RLock()
    defer s.mutex.RUnlock()
    
    return map[string]interface{}{
        "total":        s.Total,
        "by_type":      s.ByType,
        "by_operation": s.ByOperation,
    }
}
```

### 错误报告

```go
func generateErrorReport(stats *ErrorStats) string {
    data := stats.GetStats()
    
    report := fmt.Sprintf("错误报告 (总计: %d)\n", data["total"])
    report += "===================\n\n"
    
    report += "按类型分组:\n"
    for errorType, count := range data["by_type"].(map[pip.ErrorType]int) {
        report += fmt.Sprintf("- %s: %d\n", errorType, count)
    }
    
    report += "\n按操作分组:\n"
    for operation, count := range data["by_operation"].(map[string]int) {
        report += fmt.Sprintf("- %s: %d\n", operation, count)
    }
    
    return report
}
```

## 最佳实践

### 1. 总是检查错误

```go
// 好的做法
if err := manager.InstallPackage(pkg); err != nil {
    handleError(err)
    return
}

// 不好的做法
manager.InstallPackage(pkg) // 忽略错误
```

### 2. 提供有用的错误消息

```go
func installPackage(manager *pip.Manager, packageName string) error {
    pkg := &pip.PackageSpec{Name: packageName}
    
    if err := manager.InstallPackage(pkg); err != nil {
        return fmt.Errorf("安装包 '%s' 失败: %w", packageName, err)
    }
    
    return nil
}
```

### 3. 使用适当的错误处理策略

```go
func robustInstallPackage(manager *pip.Manager, pkg *pip.PackageSpec) error {
    // 1. 首先尝试正常安装
    if err := manager.InstallPackage(pkg); err == nil {
        return nil
    }
    
    // 2. 如果失败，尝试重试机制
    if err := installWithRetry(manager, pkg, 3); err == nil {
        return nil
    }
    
    // 3. 如果仍然失败，尝试错误恢复
    return installPackageWithRecovery(manager, pkg)
}
```

### 4. 记录错误上下文

```go
func installPackageWithContext(manager *pip.Manager, pkg *pip.PackageSpec) error {
    context := map[string]interface{}{
        "package_name": pkg.Name,
        "version":      pkg.Version,
        "timestamp":    time.Now(),
    }
    
    if err := manager.InstallPackage(pkg); err != nil {
        logStructuredError(err, "install_package", context)
        return err
    }
    
    return nil
}
```

## 下一步

- 学习[日志记录](./logging.md)
- 查看[错误处理示例](/zh/examples/advanced-usage.md)
- 探索[API 参考](/zh/api/errors.md)
- 了解[贡献指南](./contributing.md)
