# Advanced Usage Examples

This page provides advanced examples and patterns for using the Go Pip SDK in complex scenarios.

## Custom Package Managers

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

type CustomPackageManager struct {
    pip.PipManager
    retryCount int
    timeout    time.Duration
}

func NewCustomPackageManager(config *pip.Config) *CustomPackageManager {
    if config == nil {
        config = &pip.Config{}
    }
    
    if config.Timeout == 0 {
        config.Timeout = 60 * time.Second
    }
    
    return &CustomPackageManager{
        PipManager: pip.NewManager(config),
        retryCount: 3,
        timeout:    config.Timeout,
    }
}

func (cpm *CustomPackageManager) InstallWithRetry(pkg *pip.PackageSpec) error {
    var lastErr error
    
    for attempt := 0; attempt < cpm.retryCount; attempt++ {
        ctx, cancel := context.WithTimeout(context.Background(), cpm.timeout)
        defer cancel()
        
        // Create a new manager with context for this attempt
        manager := pip.NewManagerWithContext(ctx, nil)
        
        err := manager.InstallPackage(pkg)
        if err == nil {
            if attempt > 0 {
                fmt.Printf("Package %s installed successfully on attempt %d\n", 
                    pkg.Name, attempt+1)
            }
            return nil
        }
        
        lastErr = err
        
        // Check if it's a retryable error
        if pip.IsErrorType(err, pip.ErrorTypeNetworkError) || 
           pip.IsErrorType(err, pip.ErrorTypeTimeout) {
            
            if attempt < cpm.retryCount-1 {
                waitTime := time.Duration(attempt+1) * 5 * time.Second
                fmt.Printf("Attempt %d failed, retrying in %v...\n", 
                    attempt+1, waitTime)
                time.Sleep(waitTime)
                continue
            }
        }
        
        // Non-retryable error
        break
    }
    
    return fmt.Errorf("failed after %d attempts: %w", cpm.retryCount, lastErr)
}

func main() {
    config := &pip.Config{
        Timeout: 30 * time.Second,
        Logger: &pip.LoggerConfig{
            Level: pip.LogLevelDebug,
        },
    }
    
    manager := NewCustomPackageManager(config)
    
    pkg := &pip.PackageSpec{
        Name: "requests",
    }
    
    err := manager.InstallWithRetry(pkg)
    if err != nil {
        log.Fatal(err)
    }
}
```

## Parallel Package Installation

```go
package main

import (
    "fmt"
    "log"
    "sync"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

type PackageInstaller struct {
    manager     pip.PipManager
    concurrency int
}

type InstallResult struct {
    Package *pip.PackageSpec
    Error   error
}

func NewPackageInstaller(manager pip.PipManager, concurrency int) *PackageInstaller {
    return &PackageInstaller{
        manager:     manager,
        concurrency: concurrency,
    }
}

func (pi *PackageInstaller) InstallPackages(packages []*pip.PackageSpec) []InstallResult {
    jobs := make(chan *pip.PackageSpec, len(packages))
    results := make(chan InstallResult, len(packages))
    
    // Start workers
    var wg sync.WaitGroup
    for i := 0; i < pi.concurrency; i++ {
        wg.Add(1)
        go pi.worker(jobs, results, &wg)
    }
    
    // Send jobs
    for _, pkg := range packages {
        jobs <- pkg
    }
    close(jobs)
    
    // Wait for completion
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Collect results
    var installResults []InstallResult
    for result := range results {
        installResults = append(installResults, result)
    }
    
    return installResults
}

func (pi *PackageInstaller) worker(jobs <-chan *pip.PackageSpec, results chan<- InstallResult, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for pkg := range jobs {
        fmt.Printf("Installing %s...\n", pkg.Name)
        err := pi.manager.InstallPackage(pkg)
        
        results <- InstallResult{
            Package: pkg,
            Error:   err,
        }
        
        if err != nil {
            fmt.Printf("Failed to install %s: %v\n", pkg.Name, err)
        } else {
            fmt.Printf("Successfully installed %s\n", pkg.Name)
        }
    }
}

func main() {
    manager := pip.NewManager(nil)
    installer := NewPackageInstaller(manager, 3) // 3 concurrent installations
    
    packages := []*pip.PackageSpec{
        {Name: "requests"},
        {Name: "flask"},
        {Name: "django"},
        {Name: "numpy"},
        {Name: "pandas"},
        {Name: "matplotlib"},
    }
    
    results := installer.InstallPackages(packages)
    
    // Summary
    successful := 0
    failed := 0
    
    for _, result := range results {
        if result.Error == nil {
            successful++
        } else {
            failed++
        }
    }
    
    fmt.Printf("\nInstallation Summary:\n")
    fmt.Printf("Successful: %d\n", successful)
    fmt.Printf("Failed: %d\n", failed)
}
```

## Environment Manager

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

type EnvironmentManager struct {
    manager     pip.PipManager
    venvManager pip.VenvManager
    baseDir     string
}

type Environment struct {
    Name     string
    Path     string
    Packages []string
    Active   bool
}

func NewEnvironmentManager(baseDir string) *EnvironmentManager {
    manager := pip.NewManager(nil)
    venvManager := pip.NewVenvManager(manager)
    
    return &EnvironmentManager{
        manager:     manager,
        venvManager: venvManager,
        baseDir:     baseDir,
    }
}

func (em *EnvironmentManager) CreateEnvironment(env *Environment) error {
    envPath := filepath.Join(em.baseDir, env.Name)
    
    // Create virtual environment
    err := em.venvManager.CreateVenv(envPath)
    if err != nil {
        return fmt.Errorf("failed to create environment: %w", err)
    }
    
    // Activate environment
    err = em.venvManager.ActivateVenv(envPath)
    if err != nil {
        return fmt.Errorf("failed to activate environment: %w", err)
    }
    
    // Install packages
    for _, pkgName := range env.Packages {
        pkg := &pip.PackageSpec{Name: pkgName}
        err = em.manager.InstallPackage(pkg)
        if err != nil {
            log.Printf("Warning: Failed to install %s: %v", pkgName, err)
        }
    }
    
    // Deactivate
    err = em.venvManager.DeactivateVenv()
    if err != nil {
        return fmt.Errorf("failed to deactivate environment: %w", err)
    }
    
    env.Path = envPath
    return nil
}

func (em *EnvironmentManager) ListEnvironments() ([]*Environment, error) {
    entries, err := os.ReadDir(em.baseDir)
    if err != nil {
        return nil, err
    }
    
    var environments []*Environment
    
    for _, entry := range entries {
        if entry.IsDir() {
            envPath := filepath.Join(em.baseDir, entry.Name())
            
            if em.venvManager.IsVenvValid(envPath) {
                env := &Environment{
                    Name: entry.Name(),
                    Path: envPath,
                }
                
                // Get package list
                err = em.venvManager.ActivateVenv(envPath)
                if err == nil {
                    packages, err := em.manager.ListPackages()
                    if err == nil {
                        for _, pkg := range packages {
                            env.Packages = append(env.Packages, pkg.Name)
                        }
                    }
                    em.venvManager.DeactivateVenv()
                }
                
                environments = append(environments, env)
            }
        }
    }
    
    return environments, nil
}

func (em *EnvironmentManager) CloneEnvironment(sourceName, targetName string) error {
    sourcePath := filepath.Join(em.baseDir, sourceName)
    
    if !em.venvManager.IsVenvValid(sourcePath) {
        return fmt.Errorf("source environment %s is not valid", sourceName)
    }
    
    // Activate source environment
    err := em.venvManager.ActivateVenv(sourcePath)
    if err != nil {
        return err
    }
    
    // Get package list
    packages, err := em.manager.FreezePackages()
    if err != nil {
        return err
    }
    
    // Deactivate source
    em.venvManager.DeactivateVenv()
    
    // Create target environment
    var packageNames []string
    for _, pkg := range packages {
        packageNames = append(packageNames, pkg.Name+"=="+pkg.Version)
    }
    
    targetEnv := &Environment{
        Name:     targetName,
        Packages: packageNames,
    }
    
    return em.CreateEnvironment(targetEnv)
}

func main() {
    envManager := NewEnvironmentManager("./environments")
    
    // Create environments directory
    os.MkdirAll("./environments", 0755)
    
    // Create development environment
    devEnv := &Environment{
        Name: "development",
        Packages: []string{
            "flask",
            "requests",
            "pytest",
            "black",
            "flake8",
        },
    }
    
    fmt.Println("Creating development environment...")
    err := envManager.CreateEnvironment(devEnv)
    if err != nil {
        log.Fatal(err)
    }
    
    // Create production environment
    prodEnv := &Environment{
        Name: "production",
        Packages: []string{
            "flask",
            "requests",
            "gunicorn",
        },
    }
    
    fmt.Println("Creating production environment...")
    err = envManager.CreateEnvironment(prodEnv)
    if err != nil {
        log.Fatal(err)
    }
    
    // List all environments
    fmt.Println("\nListing environments...")
    environments, err := envManager.ListEnvironments()
    if err != nil {
        log.Fatal(err)
    }
    
    for _, env := range environments {
        fmt.Printf("Environment: %s\n", env.Name)
        fmt.Printf("  Path: %s\n", env.Path)
        fmt.Printf("  Packages: %d\n", len(env.Packages))
        for _, pkg := range env.Packages {
            fmt.Printf("    - %s\n", pkg)
        }
        fmt.Println()
    }
    
    // Clone development to testing
    fmt.Println("Cloning development to testing...")
    err = envManager.CloneEnvironment("development", "testing")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Environment management completed!")
}
```

## Package Dependency Analyzer

```go
package main

import (
    "fmt"
    "log"
    "strings"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

type DependencyAnalyzer struct {
    manager pip.PipManager
}

type PackageDependency struct {
    Name         string
    Version      string
    Dependencies []string
    Dependents   []string
}

func NewDependencyAnalyzer(manager pip.PipManager) *DependencyAnalyzer {
    return &DependencyAnalyzer{manager: manager}
}

func (da *DependencyAnalyzer) AnalyzeDependencies() (map[string]*PackageDependency, error) {
    packages, err := da.manager.ListPackages()
    if err != nil {
        return nil, err
    }
    
    dependencies := make(map[string]*PackageDependency)
    
    // Initialize package dependencies
    for _, pkg := range packages {
        dependencies[pkg.Name] = &PackageDependency{
            Name:         pkg.Name,
            Version:      pkg.Version,
            Dependencies: []string{},
            Dependents:   []string{},
        }
    }
    
    // Get detailed information for each package
    for _, pkg := range packages {
        info, err := da.manager.ShowPackage(pkg.Name)
        if err != nil {
            log.Printf("Warning: Could not get info for %s: %v", pkg.Name, err)
            continue
        }
        
        // Parse requires field
        if info.Requires != "" {
            requires := strings.Split(info.Requires, ", ")
            for _, req := range requires {
                req = strings.TrimSpace(req)
                if req != "" {
                    dependencies[pkg.Name].Dependencies = append(
                        dependencies[pkg.Name].Dependencies, req)
                    
                    // Add to dependents
                    if dep, exists := dependencies[req]; exists {
                        dep.Dependents = append(dep.Dependents, pkg.Name)
                    }
                }
            }
        }
    }
    
    return dependencies, nil
}

func (da *DependencyAnalyzer) FindOrphanPackages(dependencies map[string]*PackageDependency) []string {
    var orphans []string
    
    for name, dep := range dependencies {
        if len(dep.Dependents) == 0 {
            orphans = append(orphans, name)
        }
    }
    
    return orphans
}

func (da *DependencyAnalyzer) FindCircularDependencies(dependencies map[string]*PackageDependency) [][]string {
    var cycles [][]string
    visited := make(map[string]bool)
    recStack := make(map[string]bool)
    
    var dfs func(string, []string) bool
    dfs = func(node string, path []string) bool {
        visited[node] = true
        recStack[node] = true
        path = append(path, node)
        
        if dep, exists := dependencies[node]; exists {
            for _, neighbor := range dep.Dependencies {
                if !visited[neighbor] {
                    if dfs(neighbor, path) {
                        return true
                    }
                } else if recStack[neighbor] {
                    // Found cycle
                    cycleStart := -1
                    for i, p := range path {
                        if p == neighbor {
                            cycleStart = i
                            break
                        }
                    }
                    if cycleStart != -1 {
                        cycle := append(path[cycleStart:], neighbor)
                        cycles = append(cycles, cycle)
                    }
                    return true
                }
            }
        }
        
        recStack[node] = false
        return false
    }
    
    for name := range dependencies {
        if !visited[name] {
            dfs(name, []string{})
        }
    }
    
    return cycles
}

func (da *DependencyAnalyzer) GenerateDependencyTree(packageName string, dependencies map[string]*PackageDependency) {
    fmt.Printf("Dependency tree for %s:\n", packageName)
    da.printDependencyTree(packageName, dependencies, 0, make(map[string]bool))
}

func (da *DependencyAnalyzer) printDependencyTree(packageName string, dependencies map[string]*PackageDependency, level int, visited map[string]bool) {
    indent := strings.Repeat("  ", level)
    
    if visited[packageName] {
        fmt.Printf("%s%s (circular)\n", indent, packageName)
        return
    }
    
    visited[packageName] = true
    fmt.Printf("%s%s\n", indent, packageName)
    
    if dep, exists := dependencies[packageName]; exists {
        for _, depName := range dep.Dependencies {
            da.printDependencyTree(depName, dependencies, level+1, visited)
        }
    }
    
    visited[packageName] = false
}

func main() {
    manager := pip.NewManager(nil)
    analyzer := NewDependencyAnalyzer(manager)
    
    fmt.Println("Analyzing package dependencies...")
    dependencies, err := analyzer.AnalyzeDependencies()
    if err != nil {
        log.Fatal(err)
    }
    
    // Print summary
    fmt.Printf("Found %d packages\n\n", len(dependencies))
    
    // Find orphan packages
    orphans := analyzer.FindOrphanPackages(dependencies)
    if len(orphans) > 0 {
        fmt.Printf("Orphan packages (no dependents):\n")
        for _, orphan := range orphans {
            fmt.Printf("  - %s\n", orphan)
        }
        fmt.Println()
    }
    
    // Find circular dependencies
    cycles := analyzer.FindCircularDependencies(dependencies)
    if len(cycles) > 0 {
        fmt.Printf("Circular dependencies found:\n")
        for i, cycle := range cycles {
            fmt.Printf("  Cycle %d: %s\n", i+1, strings.Join(cycle, " -> "))
        }
        fmt.Println()
    }
    
    // Generate dependency tree for a specific package
    if len(dependencies) > 0 {
        var firstPackage string
        for name := range dependencies {
            firstPackage = name
            break
        }
        analyzer.GenerateDependencyTree(firstPackage, dependencies)
    }
}
```

## Configuration Management

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "github.com/scagogogo/go-pip-sdk/pkg/pip"
)

type ProjectConfig struct {
    Name            string            `json:"name"`
    Version         string            `json:"version"`
    Description     string            `json:"description"`
    Dependencies    map[string]string `json:"dependencies"`
    DevDependencies map[string]string `json:"dev_dependencies"`
    PythonVersion   string            `json:"python_version"`
    VenvPath        string            `json:"venv_path"`
    IndexURL        string            `json:"index_url"`
    ExtraIndexURLs  []string          `json:"extra_index_urls"`
}

type ConfigManager struct {
    manager     pip.PipManager
    venvManager pip.VenvManager
    configPath  string
}

func NewConfigManager(configPath string) *ConfigManager {
    manager := pip.NewManager(nil)
    venvManager := pip.NewVenvManager(manager)
    
    return &ConfigManager{
        manager:     manager,
        venvManager: venvManager,
        configPath:  configPath,
    }
}

func (cm *ConfigManager) LoadConfig() (*ProjectConfig, error) {
    data, err := os.ReadFile(cm.configPath)
    if err != nil {
        return nil, err
    }
    
    var config ProjectConfig
    err = json.Unmarshal(data, &config)
    if err != nil {
        return nil, err
    }
    
    return &config, nil
}

func (cm *ConfigManager) SaveConfig(config *ProjectConfig) error {
    data, err := json.MarshalIndent(config, "", "  ")
    if err != nil {
        return err
    }
    
    return os.WriteFile(cm.configPath, data, 0644)
}

func (cm *ConfigManager) SetupProject(config *ProjectConfig) error {
    // Create virtual environment if specified
    if config.VenvPath != "" {
        fmt.Printf("Creating virtual environment at %s...\n", config.VenvPath)
        err := cm.venvManager.CreateVenv(config.VenvPath)
        if err != nil {
            return err
        }
        
        err = cm.venvManager.ActivateVenv(config.VenvPath)
        if err != nil {
            return err
        }
    }
    
    // Install dependencies
    fmt.Println("Installing dependencies...")
    for name, version := range config.Dependencies {
        pkg := &pip.PackageSpec{
            Name:    name,
            Version: version,
            Index:   config.IndexURL,
        }
        
        err := cm.manager.InstallPackage(pkg)
        if err != nil {
            log.Printf("Warning: Failed to install %s: %v", name, err)
        }
    }
    
    // Install dev dependencies
    fmt.Println("Installing development dependencies...")
    for name, version := range config.DevDependencies {
        pkg := &pip.PackageSpec{
            Name:    name,
            Version: version,
            Index:   config.IndexURL,
        }
        
        err := cm.manager.InstallPackage(pkg)
        if err != nil {
            log.Printf("Warning: Failed to install dev dependency %s: %v", name, err)
        }
    }
    
    return nil
}

func (cm *ConfigManager) UpdateConfig() error {
    config, err := cm.LoadConfig()
    if err != nil {
        return err
    }
    
    // Get current packages
    packages, err := cm.manager.FreezePackages()
    if err != nil {
        return err
    }
    
    // Update dependencies
    config.Dependencies = make(map[string]string)
    for _, pkg := range packages {
        // Skip development packages (this is a simple heuristic)
        if !cm.isDevPackage(pkg.Name) {
            config.Dependencies[pkg.Name] = "==" + pkg.Version
        }
    }
    
    return cm.SaveConfig(config)
}

func (cm *ConfigManager) isDevPackage(name string) bool {
    devPackages := []string{
        "pytest", "black", "flake8", "mypy", "coverage",
        "pre-commit", "sphinx", "twine", "wheel",
    }
    
    for _, devPkg := range devPackages {
        if name == devPkg {
            return true
        }
    }
    
    return false
}

func main() {
    configPath := "./project.json"
    configManager := NewConfigManager(configPath)
    
    // Create example config
    config := &ProjectConfig{
        Name:        "example-project",
        Version:     "1.0.0",
        Description: "An example project",
        Dependencies: map[string]string{
            "requests": ">=2.25.0",
            "flask":    ">=2.0.0",
        },
        DevDependencies: map[string]string{
            "pytest": ">=6.0.0",
            "black":  ">=21.0.0",
        },
        PythonVersion: "3.9",
        VenvPath:      "./venv",
        IndexURL:      "https://pypi.org/simple/",
    }
    
    // Save config
    err := configManager.SaveConfig(config)
    if err != nil {
        log.Fatal(err)
    }
    
    // Setup project
    err = configManager.SetupProject(config)
    if err != nil {
        log.Fatal(err)
    }
    
    // Update config with current state
    err = configManager.UpdateConfig()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Project setup completed!")
}
```

## Best Practices

1. **Error Handling**: Always implement comprehensive error handling
2. **Concurrency**: Use goroutines for parallel operations when appropriate
3. **Configuration**: Use configuration files for complex setups
4. **Logging**: Implement detailed logging for debugging
5. **Testing**: Write tests for custom functionality
6. **Resource Management**: Properly manage virtual environments and resources
7. **Performance**: Monitor and optimize performance for large operations

## Next Steps

- Explore the [API Reference](../api/) for detailed documentation
- Check out [Error Handling](../guide/error-handling.md) patterns
- Learn about [Logging](../guide/logging.md) best practices
