# 高级用法示例

本页面提供了 Go Pip SDK 的高级使用示例，包括性能优化、错误处理、自定义配置和企业级应用场景。

## 高级配置

### 企业级配置

```go
package main

import (
    "time"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

func createEnterpriseConfig() *pip.Config {
    return &pip.Config{
        // Python 配置
        PythonPath: "/opt/python/bin/python3",
        PipPath:    "/opt/python/bin/pip",
        
        // 网络配置
        DefaultIndex: "https://pypi.company.com/simple/",
        ExtraIndexes: []string{
            "https://pypi.org/simple/",
            "https://pypi.tuna.tsinghua.edu.cn/simple/",
        },
        TrustedHosts: []string{
            "pypi.company.com",
            "pypi.org",
            "pypi.tuna.tsinghua.edu.cn",
        },
        Timeout: 300 * time.Second, // 5 分钟超时
        Retries: 5,
        
        // 缓存配置
        CacheDir: "/var/cache/pip",
        NoCache:  false,
        
        // 日志配置
        LogLevel: "INFO",
        LogFile:  "/var/log/pip-sdk.log",
        
        // 环境变量
        Environment: map[string]string{
            "PIP_NO_CACHE_DIR":              "0",
            "PIP_DISABLE_PIP_VERSION_CHECK": "1",
            "PIP_TIMEOUT":                   "300",
            "PIP_RETRIES":                   "5",
        },
        
        // 额外选项
        ExtraOptions: map[string]string{
            "cert":         "/etc/ssl/certs/company-ca.pem",
            "client-cert":  "/etc/ssl/certs/client.pem",
            "trusted-host": "pypi.company.com",
        },
    }
}

func enterpriseSetup() {
    config := createEnterpriseConfig()
    manager := pip.NewManager(config)
    
    // 验证配置
    if err := validateEnterpriseSetup(manager); err != nil {
        log.Fatalf("企业配置验证失败: %v", err)
    }
    
    fmt.Println("✅ 企业级配置设置成功")
}

func validateEnterpriseSetup(manager *pip.Manager) error {
    // 检查 Python 环境
    pythonPath, err := manager.GetPythonPath()
    if err != nil {
        return fmt.Errorf("Python 路径验证失败: %w", err)
    }
    fmt.Printf("Python 路径: %s\n", pythonPath)
    
    // 检查 pip 版本
    version, err := manager.GetVersion()
    if err != nil {
        return fmt.Errorf("Pip 版本检查失败: %w", err)
    }
    fmt.Printf("Pip 版本: %s\n", version)
    
    // 检查网络连接
    if err := testNetworkConnectivity(manager); err != nil {
        return fmt.Errorf("网络连接测试失败: %w", err)
    }
    
    return nil
}
```

### 多环境配置管理

```go
func multiEnvironmentConfig() {
    environments := map[string]*pip.Config{
        "development": {
            LogLevel:     "DEBUG",
            DefaultIndex: "https://pypi.org/simple/",
            Timeout:      60 * time.Second,
            Retries:      3,
        },
        "staging": {
            LogLevel:     "INFO",
            DefaultIndex: "https://pypi.company.com/simple/",
            Timeout:      120 * time.Second,
            Retries:      5,
        },
        "production": {
            LogLevel:     "WARN",
            DefaultIndex: "https://pypi.company.com/simple/",
            Timeout:      300 * time.Second,
            Retries:      10,
            Environment: map[string]string{
                "PIP_NO_CACHE_DIR": "1",
            },
        },
    }
    
    env := os.Getenv("ENVIRONMENT")
    if env == "" {
        env = "development"
    }
    
    config, exists := environments[env]
    if !exists {
        log.Fatalf("未知环境: %s", env)
    }
    
    manager := pip.NewManager(config)
    fmt.Printf("使用 %s 环境配置\n", env)
}
```

## 性能优化

### 并发包管理

```go
func concurrentPackageManagement() {
    manager := pip.NewManager(nil)
    
    packages := []*pip.PackageSpec{
        {Name: "requests", Version: ">=2.25.0"},
        {Name: "click", Version: ">=7.0"},
        {Name: "pydantic", Version: ">=1.8.0"},
        {Name: "fastapi", Version: ">=0.68.0"},
        {Name: "uvicorn", Version: ">=0.15.0"},
        {Name: "sqlalchemy", Version: ">=1.4.0"},
        {Name: "alembic", Version: ">=1.7.0"},
        {Name: "pytest", Version: ">=6.0"},
    }
    
    // 使用工作池进行并发安装
    const maxWorkers = 4
    semaphore := make(chan struct{}, maxWorkers)
    var wg sync.WaitGroup
    results := make(chan result, len(packages))
    
    type result struct {
        pkg *pip.PackageSpec
        err error
    }
    
    for _, pkg := range packages {
        wg.Add(1)
        go func(p *pip.PackageSpec) {
            defer wg.Done()
            
            // 获取信号量
            semaphore <- struct{}{}
            defer func() { <-semaphore }()
            
            // 为每个 goroutine 创建独立的管理器
            workerManager := pip.NewManager(manager.GetConfig())
            err := workerManager.InstallPackage(p)
            results <- result{pkg: p, err: err}
        }(pkg)
    }
    
    // 等待所有任务完成
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // 收集结果
    var successful, failed int
    for res := range results {
        if res.err != nil {
            fmt.Printf("❌ %s 安装失败: %v\n", res.pkg.Name, res.err)
            failed++
        } else {
            fmt.Printf("✅ %s 安装成功\n", res.pkg.Name)
            successful++
        }
    }
    
    fmt.Printf("\n并发安装完成: %d 成功, %d 失败\n", successful, failed)
}
```

### 缓存优化

```go
func optimizedCaching() {
    // 配置高效缓存
    config := &pip.Config{
        CacheDir: "/tmp/pip-cache-optimized",
        NoCache:  false,
        ExtraOptions: map[string]string{
            "cache-dir": "/tmp/pip-cache-optimized",
            "no-deps":   "", // 某些情况下跳过依赖检查
        },
    }
    
    manager := pip.NewManager(config)
    
    // 预热缓存
    fmt.Println("预热包缓存...")
    commonPackages := []string{
        "requests", "click", "pydantic", "fastapi",
        "numpy", "pandas", "matplotlib", "pytest",
    }
    
    for _, pkgName := range commonPackages {
        pkg := &pip.PackageSpec{Name: pkgName}
        if err := manager.DownloadPackage(pkg, config.CacheDir); err != nil {
            fmt.Printf("下载 %s 到缓存失败: %v\n", pkgName, err)
        } else {
            fmt.Printf("✅ %s 已缓存\n", pkgName)
        }
    }
    
    fmt.Println("缓存预热完成")
}
```

## 高级错误处理

### 智能重试机制

```go
func intelligentRetry() {
    manager := pip.NewManager(nil)
    
    pkg := &pip.PackageSpec{Name: "some-package"}
    
    if err := installWithIntelligentRetry(manager, pkg); err != nil {
        log.Fatalf("最终安装失败: %v", err)
    }
}

func installWithIntelligentRetry(manager *pip.Manager, pkg *pip.PackageSpec) error {
    maxRetries := 5
    baseDelay := time.Second
    
    for attempt := 0; attempt < maxRetries; attempt++ {
        err := manager.InstallPackage(pkg)
        if err == nil {
            return nil // 成功
        }
        
        // 根据错误类型决定重试策略
        switch pip.GetErrorType(err) {
        case pip.ErrorTypeNetworkError, pip.ErrorTypeNetworkTimeout:
            // 网络错误：指数退避重试
            delay := time.Duration(1<<uint(attempt)) * baseDelay
            fmt.Printf("网络错误，%v 后重试 (%d/%d)...\n", delay, attempt+1, maxRetries)
            time.Sleep(delay)
            continue
            
        case pip.ErrorTypePackageNotFound:
            // 包未找到：尝试搜索相似包
            if attempt == 0 {
                if err := suggestAlternatives(manager, pkg.Name); err == nil {
                    continue // 重试一次
                }
            }
            return err // 不再重试
            
        case pip.ErrorTypeVersionConflict:
            // 版本冲突：尝试放宽版本约束
            if attempt == 0 && pkg.Version != "" {
                fmt.Println("版本冲突，尝试放宽版本约束...")
                pkg.Version = "" // 移除版本约束
                continue
            }
            return err
            
        case pip.ErrorTypePermissionDenied:
            // 权限错误：尝试用户级安装
            if attempt == 0 && !pkg.UserInstall {
                fmt.Println("权限被拒绝，尝试用户级安装...")
                pkg.UserInstall = true
                continue
            }
            return err
            
        default:
            return err // 其他错误不重试
        }
    }
    
    return fmt.Errorf("重试 %d 次后仍然失败", maxRetries)
}

func suggestAlternatives(manager *pip.Manager, packageName string) error {
    results, err := manager.SearchPackages(packageName)
    if err != nil {
        return err
    }
    
    if len(results) > 0 {
        fmt.Printf("包 '%s' 未找到，相似包:\n", packageName)
        for i, result := range results[:3] { // 显示前3个
            fmt.Printf("%d. %s - %s\n", i+1, result.Name, result.Summary)
        }
    }
    
    return fmt.Errorf("未找到包 %s", packageName)
}
```

### 错误聚合和报告

```go
func errorAggregationExample() {
    manager := pip.NewManager(nil)
    collector := pip.NewErrorCollector()
    
    packages := []string{
        "requests", "nonexistent-package", "click",
        "another-bad-package", "pydantic",
    }
    
    // 收集所有错误
    for _, pkgName := range packages {
        pkg := &pip.PackageSpec{Name: pkgName}
        if err := manager.InstallPackage(pkg); err != nil {
            collector.Add(fmt.Errorf("安装 %s 失败: %w", pkgName, err))
        }
    }
    
    // 生成错误报告
    if collector.HasErrors() {
        report := generateDetailedErrorReport(collector.Errors())
        fmt.Println("错误报告:")
        fmt.Println(report)
        
        // 保存错误报告
        if err := saveErrorReport(report, "error-report.json"); err != nil {
            log.Printf("保存错误报告失败: %v", err)
        }
    }
}

func generateDetailedErrorReport(errors []error) string {
    report := map[string]interface{}{
        "timestamp": time.Now().UTC(),
        "total_errors": len(errors),
        "errors": make([]map[string]interface{}, 0, len(errors)),
        "summary": make(map[string]int),
    }
    
    for _, err := range errors {
        errorInfo := map[string]interface{}{
            "message": err.Error(),
            "type":    "unknown",
        }
        
        if pipErr, ok := err.(*pip.PipErrorDetails); ok {
            errorInfo["type"] = string(pipErr.Type)
            errorInfo["command"] = pipErr.Command
            errorInfo["exit_code"] = pipErr.ExitCode
            errorInfo["suggestions"] = pipErr.Suggestions
            
            // 统计错误类型
            report["summary"].(map[string]int)[string(pipErr.Type)]++
        }
        
        report["errors"] = append(report["errors"].([]map[string]interface{}), errorInfo)
    }
    
    jsonData, _ := json.MarshalIndent(report, "", "  ")
    return string(jsonData)
}
```

## 自定义扩展

### 自定义日志记录器

```go
func customLoggingExample() {
    // 创建自定义日志记录器
    logger, err := pip.NewLogger(&pip.LoggerConfig{
        Level:  pip.LogLevelInfo,
        Output: os.Stdout,
        Format: pip.LogFormatJSON,
        Fields: map[string]interface{}{
            "service":     "pip-manager",
            "version":     "1.0.0",
            "environment": os.Getenv("ENVIRONMENT"),
        },
        Filter: func(entry *pip.LogEntry) bool {
            // 过滤敏感信息
            if strings.Contains(entry.Message, "password") {
                return false
            }
            return true
        },
    })
    if err != nil {
        log.Fatalf("创建日志记录器失败: %v", err)
    }
    defer logger.Close()
    
    // 设置到管理器
    manager := pip.NewManager(nil)
    manager.SetCustomLogger(logger)
    
    // 使用带自定义日志的管理器
    pkg := &pip.PackageSpec{Name: "requests"}
    if err := manager.InstallPackage(pkg); err != nil {
        logger.ErrorWithFields("包安装失败", map[string]interface{}{
            "package": pkg.Name,
            "error":   err.Error(),
        })
    } else {
        logger.InfoWithFields("包安装成功", map[string]interface{}{
            "package": pkg.Name,
        })
    }
}
```

### 自定义包索引

```go
func customPackageIndex() {
    // 配置自定义包索引
    config := &pip.Config{
        DefaultIndex: "https://custom-pypi.company.com/simple/",
        ExtraIndexes: []string{
            "https://pypi.org/simple/",
            "https://backup-pypi.company.com/simple/",
        },
        TrustedHosts: []string{
            "custom-pypi.company.com",
            "backup-pypi.company.com",
        },
        ExtraOptions: map[string]string{
            "cert":        "/etc/ssl/certs/company.pem",
            "client-cert": "/etc/ssl/certs/client.pem",
        },
    }
    
    manager := pip.NewManager(config)
    
    // 测试自定义索引连接
    if err := testCustomIndex(manager); err != nil {
        log.Fatalf("自定义索引测试失败: %v", err)
    }
    
    fmt.Println("✅ 自定义包索引配置成功")
}

func testCustomIndex(manager *pip.Manager) error {
    // 尝试搜索包来测试索引连接
    results, err := manager.SearchPackages("requests")
    if err != nil {
        return fmt.Errorf("搜索测试失败: %w", err)
    }
    
    if len(results) == 0 {
        return fmt.Errorf("自定义索引中未找到测试包")
    }
    
    fmt.Printf("自定义索引测试成功，找到 %d 个结果\n", len(results))
    return nil
}
```

## 监控和指标

### 性能监控

```go
func performanceMonitoring() {
    manager := pip.NewManager(nil)
    
    // 创建性能监控器
    monitor := &PerformanceMonitor{
        metrics: make(map[string]*Metric),
    }
    
    packages := []string{"requests", "click", "pydantic"}
    
    for _, pkgName := range packages {
        // 监控包安装性能
        start := time.Now()
        pkg := &pip.PackageSpec{Name: pkgName}
        
        err := manager.InstallPackage(pkg)
        duration := time.Since(start)
        
        // 记录指标
        monitor.RecordMetric("package_install_duration", duration.Seconds(), map[string]string{
            "package": pkgName,
            "success": fmt.Sprintf("%t", err == nil),
        })
        
        if err != nil {
            monitor.RecordMetric("package_install_errors", 1, map[string]string{
                "package": pkgName,
                "error_type": string(pip.GetErrorType(err)),
            })
        }
    }
    
    // 生成性能报告
    report := monitor.GenerateReport()
    fmt.Println("性能报告:")
    fmt.Println(report)
}

type PerformanceMonitor struct {
    metrics map[string]*Metric
    mutex   sync.RWMutex
}

type Metric struct {
    Name  string
    Value float64
    Tags  map[string]string
    Time  time.Time
}

func (m *PerformanceMonitor) RecordMetric(name string, value float64, tags map[string]string) {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    key := fmt.Sprintf("%s_%d", name, time.Now().UnixNano())
    m.metrics[key] = &Metric{
        Name:  name,
        Value: value,
        Tags:  tags,
        Time:  time.Now(),
    }
}

func (m *PerformanceMonitor) GenerateReport() string {
    m.mutex.RLock()
    defer m.mutex.RUnlock()
    
    report := "性能指标报告\n"
    report += "================\n"
    
    for _, metric := range m.metrics {
        report += fmt.Sprintf("%s: %.2f (时间: %s)\n", 
            metric.Name, metric.Value, metric.Time.Format("15:04:05"))
        for k, v := range metric.Tags {
            report += fmt.Sprintf("  %s: %s\n", k, v)
        }
        report += "\n"
    }
    
    return report
}
```

### 健康检查

```go
func healthCheckExample() {
    manager := pip.NewManager(nil)
    
    // 执行全面健康检查
    health := performHealthCheck(manager)
    
    fmt.Printf("健康检查结果: %s\n", health.Status)
    fmt.Printf("检查项目: %d 个\n", len(health.Checks))
    
    for _, check := range health.Checks {
        status := "✅"
        if !check.Passed {
            status = "❌"
        }
        fmt.Printf("%s %s: %s\n", status, check.Name, check.Message)
    }
    
    if !health.IsHealthy() {
        fmt.Println("系统存在问题，需要修复")
    }
}

type HealthCheck struct {
    Status string
    Checks []CheckResult
}

type CheckResult struct {
    Name    string
    Passed  bool
    Message string
}

func (h *HealthCheck) IsHealthy() bool {
    for _, check := range h.Checks {
        if !check.Passed {
            return false
        }
    }
    return true
}

func performHealthCheck(manager *pip.Manager) *HealthCheck {
    health := &HealthCheck{
        Checks: make([]CheckResult, 0),
    }
    
    // 检查 Python
    if pythonPath, err := manager.GetPythonPath(); err != nil {
        health.Checks = append(health.Checks, CheckResult{
            Name:    "Python 可用性",
            Passed:  false,
            Message: fmt.Sprintf("Python 未找到: %v", err),
        })
    } else {
        health.Checks = append(health.Checks, CheckResult{
            Name:    "Python 可用性",
            Passed:  true,
            Message: fmt.Sprintf("Python 路径: %s", pythonPath),
        })
    }
    
    // 检查 Pip
    if installed, err := manager.IsInstalled(); err != nil || !installed {
        health.Checks = append(health.Checks, CheckResult{
            Name:    "Pip 可用性",
            Passed:  false,
            Message: "Pip 未安装或不可用",
        })
    } else {
        if version, err := manager.GetVersion(); err != nil {
            health.Checks = append(health.Checks, CheckResult{
                Name:    "Pip 版本",
                Passed:  false,
                Message: fmt.Sprintf("获取版本失败: %v", err),
            })
        } else {
            health.Checks = append(health.Checks, CheckResult{
                Name:    "Pip 版本",
                Passed:  true,
                Message: fmt.Sprintf("Pip 版本: %s", version),
            })
        }
    }
    
    // 检查网络连接
    if err := testNetworkConnectivity(manager); err != nil {
        health.Checks = append(health.Checks, CheckResult{
            Name:    "网络连接",
            Passed:  false,
            Message: fmt.Sprintf("网络连接失败: %v", err),
        })
    } else {
        health.Checks = append(health.Checks, CheckResult{
            Name:    "网络连接",
            Passed:  true,
            Message: "网络连接正常",
        })
    }
    
    // 设置总体状态
    if health.IsHealthy() {
        health.Status = "健康"
    } else {
        health.Status = "不健康"
    }
    
    return health
}

func testNetworkConnectivity(manager *pip.Manager) error {
    // 尝试搜索一个常见包来测试网络连接
    _, err := manager.SearchPackages("requests")
    return err
}
```

## 企业级集成

### CI/CD 集成

```go
func cicdIntegration() {
    // CI/CD 环境配置
    config := &pip.Config{
        LogLevel: "INFO",
        Environment: map[string]string{
            "PIP_NO_CACHE_DIR":              "1",
            "PIP_DISABLE_PIP_VERSION_CHECK": "1",
            "PIP_QUIET":                     "1",
        },
        Timeout: 600 * time.Second, // 10 分钟超时
        Retries: 3,
    }
    
    manager := pip.NewManager(config)
    
    // CI/CD 流水线步骤
    steps := []struct {
        name string
        fn   func() error
    }{
        {"环境验证", func() error { return validateCIEnvironment(manager) }},
        {"依赖安装", func() error { return installCIDependencies(manager) }},
        {"测试执行", func() error { return runTests(manager) }},
        {"构建验证", func() error { return validateBuild(manager) }},
    }
    
    for _, step := range steps {
        fmt.Printf("执行步骤: %s...\n", step.name)
        if err := step.fn(); err != nil {
            fmt.Printf("❌ %s 失败: %v\n", step.name, err)
            os.Exit(1)
        }
        fmt.Printf("✅ %s 完成\n", step.name)
    }
    
    fmt.Println("🎉 CI/CD 流水线执行成功")
}

func validateCIEnvironment(manager *pip.Manager) error {
    // 验证 CI 环境
    health := performHealthCheck(manager)
    if !health.IsHealthy() {
        return fmt.Errorf("环境健康检查失败")
    }
    return nil
}

func installCIDependencies(manager *pip.Manager) error {
    // 安装 CI 依赖
    if err := manager.InstallRequirements("requirements.txt"); err != nil {
        return fmt.Errorf("安装生产依赖失败: %w", err)
    }
    
    if err := manager.InstallRequirements("dev-requirements.txt"); err != nil {
        return fmt.Errorf("安装开发依赖失败: %w", err)
    }
    
    return nil
}

func runTests(manager *pip.Manager) error {
    // 运行测试（这里只是示例）
    fmt.Println("运行测试套件...")
    // 实际实现会调用 pytest 或其他测试框架
    return nil
}

func validateBuild(manager *pip.Manager) error {
    // 验证构建（这里只是示例）
    fmt.Println("验证构建产物...")
    // 实际实现会检查构建输出
    return nil
}
```

## 下一步

- 查看[基本用法](./basic-usage.md)了解基础功能
- 学习[包管理示例](./package-management.md)掌握包操作
- 探索[虚拟环境示例](./virtual-environments.md)了解环境管理
- 阅读[项目初始化示例](./project-initialization.md)学习项目管理
