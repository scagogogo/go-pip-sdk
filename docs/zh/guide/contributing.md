# 为 Go Pip SDK 贡献代码

感谢您对为 Go Pip SDK 项目贡献代码的兴趣！本指南将帮助您开始为项目做出贡献。

## 目录

- [开始](#开始)
- [开发环境设置](#开发环境设置)
- [代码风格](#代码风格)
- [测试](#测试)
- [提交更改](#提交更改)
- [报告问题](#报告问题)
- [文档](#文档)
- [社区准则](#社区准则)

## 开始

### 前置条件

在开始之前，请确保您已安装以下软件：

- Go 1.19 或更高版本
- Git
- Python 3.7+（用于测试 pip 功能）
- Make（可选，用于使用 Makefile 命令）

### Fork 和克隆

1. 在 GitHub 上 Fork 仓库
2. 在本地克隆您的 Fork：

```bash
git clone https://github.com/YOUR_USERNAME/go-pip-sdk.git
cd go-pip-sdk
```

3. 添加上游仓库：

```bash
git remote add upstream https://github.com/scagogogo/go-pip-sdk.git
```

## 开发环境设置

### 安装依赖

```bash
# 安装 Go 依赖
go mod download

# 安装开发工具
make install-tools
```

### 构建项目

```bash
# 构建项目
make build

# 或手动构建
go build ./...
```

### 运行测试

```bash
# 运行所有测试
make test

# 运行带覆盖率的测试
make test-coverage

# 运行特定测试
go test ./pkg/pip/...
```

## 代码风格

### Go 代码标准

我们遵循标准的 Go 约定：

- 使用 `gofmt` 进行格式化
- 使用 `golint` 进行代码检查
- 遵循有效的 Go 实践
- 编写清晰、自文档化的代码
- 使用有意义的变量和函数名

### 代码格式化

在提交代码之前，请确保代码格式正确：

```bash
# 格式化代码
make fmt

# 或手动格式化
gofmt -w .
```

### 代码检查

运行代码检查工具检查代码质量：

```bash
# 运行所有检查工具
make lint

# 运行特定检查工具
golangci-lint run
```

## 测试

### 编写测试

- 为所有新功能编写单元测试
- 在适当的地方使用表驱动测试
- 模拟外部依赖
- 追求高测试覆盖率（>80%）

### 测试结构

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected ExpectedType
        wantErr  bool
    }{
        {
            name:     "valid case",
            input:    validInput,
            expected: expectedOutput,
            wantErr:  false,
        },
        // 添加更多测试用例
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := FunctionName(tt.input)
            
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### 集成测试

对于需要 Python/pip 的集成测试：

```go
func TestIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("在短模式下跳过集成测试")
    }
    
    // 集成测试代码
}
```

运行集成测试：

```bash
# 运行包括集成测试在内的所有测试
make test-integration

# 跳过集成测试
go test -short ./...
```

## 提交更改

### 提交准则

遵循约定式提交格式：

```
type(scope): description

[optional body]

[optional footer]
```

类型：
- `feat`: 新功能
- `fix`: 错误修复
- `docs`: 文档更改
- `style`: 代码风格更改
- `refactor`: 代码重构
- `test`: 添加或更新测试
- `chore`: 维护任务

示例：
```
feat(manager): 添加对自定义 pip 索引的支持
fix(installer): 正确处理网络超时错误
docs(api): 更新管理器文档
```

### Pull Request 流程

1. 创建功能分支：
```bash
git checkout -b feature/your-feature-name
```

2. 进行更改并提交：
```bash
git add .
git commit -m "feat: 添加新功能"
```

3. 推送到您的 Fork：
```bash
git push origin feature/your-feature-name
```

4. 在 GitHub 上创建 Pull Request

### Pull Request 要求

- [ ] 代码遵循项目风格准则
- [ ] 测试在本地通过
- [ ] 新功能包含测试
- [ ] 文档已更新
- [ ] 提交消息遵循约定格式
- [ ] 与主分支没有合并冲突

## 报告问题

### 错误报告

报告错误时，请包含：

- Go 版本
- 操作系统
- Python/pip 版本
- 重现步骤
- 预期行为
- 实际行为
- 错误消息/日志

### 功能请求

对于功能请求，请提供：

- 功能的清晰描述
- 用例和动机
- 建议的 API（如适用）
- 使用示例

## 文档

### 代码文档

- 为所有公共函数和类型编写文档
- 使用清晰、简洁的注释
- 在文档中包含示例
- 遵循 Go 文档约定

示例：
```go
// InstallPackage 使用 pip 安装 Python 包。
// 如果安装失败，它会返回错误。
//
// 示例：
//   pkg := &PackageSpec{Name: "requests", Version: ">=2.25.0"}
//   err := manager.InstallPackage(pkg)
//   if err != nil {
//       log.Fatal(err)
//   }
func (m *Manager) InstallPackage(pkg *PackageSpec) error {
    // 实现
}
```

### 文档网站

文档网站使用 VitePress 构建：

```bash
# 安装依赖
cd docs
npm install

# 启动开发服务器
npm run dev

# 构建文档
npm run build
```

## 社区准则

### 行为准则

- 保持尊重和包容
- 欢迎新人
- 专注于建设性反馈
- 帮助他人学习和成长

### 沟通

- 使用 GitHub issues 报告错误和功能请求
- 使用 GitHub discussions 进行问题和一般讨论
- 沟通要清晰简洁
- 提供上下文和示例

### 获取帮助

- 首先检查现有问题和文档
- 在 GitHub discussions 中提问
- 提供最小可重现示例
- 保持耐心和尊重

## 发布流程

发布由维护者处理：

1. 在适当文件中更新版本
2. 更新 CHANGELOG.md
3. 创建发布标签
4. 发布发布说明

## 许可证

通过贡献，您同意您的贡献将在与项目相同的许可证（MIT 许可证）下获得许可。

## 感谢

感谢您为 Go Pip SDK 做出贡献！您的贡献有助于让这个项目对每个人都更好。
