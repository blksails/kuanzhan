# 发布指南

本文档说明如何使用 GoReleaser 配置来构建和发布 Kuanzhan CLI 工具。

## 📋 前提条件

### 必需工具
- [Go 1.23+](https://golang.org/dl/)
- [GoReleaser](https://goreleaser.com/install/)
- [Docker](https://www.docker.com/get-started) (可选，用于构建 Docker 镜像)
- [Git](https://git-scm.com/)

### 安装 GoReleaser

```bash
# macOS
brew install goreleaser

# Linux
curl -sfL https://goreleaser.com/static/run | bash

# Windows
choco install goreleaser
```

## 🔧 配置说明

### 1. 修改配置文件

编辑 `.goreleaser.yaml` 文件中的以下占位符：

```yaml
release:
  github:
    owner: "your-github-username"  # 替换为您的 GitHub 用户名
    name: "kuanzhan"               # 替换为您的仓库名

brews:
  - repository:
      owner: "your-github-username"  # 替换为您的 GitHub 用户名
      name: "homebrew-tap"

# 其他需要替换的 URL 和标识符
```

### 2. 创建 GitHub Token

1. 访问 [GitHub Personal Access Tokens](https://github.com/settings/tokens)
2. 创建新的 token，选择以下权限：
   - `repo` (完整仓库访问)
   - `write:packages` (包写入权限)
   - `read:packages` (包读取权限)
3. 将 token 添加到仓库的 Secrets 中：
   - `Settings` → `Secrets and variables` → `Actions`
   - 添加 `GITHUB_TOKEN` secret

### 3. 配置 Docker Hub（可选）

如果需要发布 Docker 镜像：

1. 在 GitHub 仓库 Secrets 中添加：
   - `DOCKER_USERNAME` - Docker Hub 用户名
   - `DOCKER_PASSWORD` - Docker Hub 密码或 token

2. 修改 `.goreleaser.yaml` 中的 Docker 配置：
   ```yaml
   dockers:
     - image_templates:
         - "your-dockerhub-username/kuanzhan:{{ .Tag }}"
         - "your-dockerhub-username/kuanzhan:latest"
   ```

## 🚀 发布流程

### 本地测试构建

```bash
# 检查配置是否正确
make check-release

# 创建快照版本（本地测试）
make snapshot

# 或者使用 goreleaser 直接命令
goreleaser release --snapshot --clean
```

### 正式发布

1. **准备发布**：
   ```bash
   # 确保代码已提交
   git add .
   git commit -m "feat: prepare for release"
   
   # 推送到远程仓库
   git push origin main
   ```

2. **创建标签**：
   ```bash
   # 创建版本标签
   git tag -a v1.0.0 -m "Release v1.0.0"
   
   # 推送标签到远程
   git push origin v1.0.0
   ```

3. **自动发布**：
   - GitHub Actions 会自动触发发布流程
   - 查看 `Actions` 标签页监控发布状态

### 手动发布

如果不使用 GitHub Actions：

```bash
# 设置环境变量
export GITHUB_TOKEN=your_github_token

# 执行发布
make release

# 或者使用 goreleaser 直接命令
goreleaser release --clean
```

## 📦 构建产物

成功发布后，将生成以下产物：

### 二进制文件
- `kuanzhan_Linux_x86_64.tar.gz`
- `kuanzhan_Linux_arm64.tar.gz`
- `kuanzhan_Darwin_x86_64.tar.gz`
- `kuanzhan_Darwin_arm64.tar.gz`
- `kuanzhan_Windows_x86_64.zip`
- `kuanzhan_Windows_arm64.zip`

### 包管理器
- **Homebrew**: `brew install your-username/tap/kuanzhan`
- **Linux 包**: `.deb`, `.rpm`, `.apk` 文件
- **Docker**: `docker pull your-username/kuanzhan:latest`

### 校验文件
- `checksums.txt` - 所有文件的校验和

## 🔧 开发工作流

### 日常开发

```bash
# 格式化代码
make fmt

# 运行测试
make test

# 运行测试覆盖率
make test-coverage

# 构建开发版本
make dev

# 运行应用
make run
```

### 多平台构建

```bash
# 构建所有平台
make build-all

# 构建 Docker 镜像
make docker-build

# 运行 Docker 容器
make docker-run
```

## 📋 发布检查清单

在发布前，请确保：

- [ ] 所有测试通过
- [ ] 代码已格式化
- [ ] 更新了 `CHANGELOG.md`
- [ ] 更新了 `README.md`
- [ ] 版本号遵循 [Semantic Versioning](https://semver.org/)
- [ ] GitHub token 已配置
- [ ] Docker Hub 凭据已配置（如果需要）

## 🔄 版本管理

### 版本号规则
- `v1.0.0` - 主要版本（破坏性更改）
- `v1.1.0` - 次要版本（新功能）
- `v1.1.1` - 补丁版本（bug 修复）

### 预发布版本
- `v1.0.0-rc1` - 候选版本
- `v1.0.0-alpha1` - 内测版本
- `v1.0.0-beta1` - 公测版本

## 🐛 故障排除

### 常见问题

1. **GitHub token 权限不足**
   - 确保 token 有 `repo` 和 `write:packages` 权限

2. **Docker 构建失败**
   - 检查 Dockerfile 语法
   - 确保 Docker 守护进程正在运行

3. **Homebrew 发布失败**
   - 确保有 `homebrew-tap` 仓库
   - 检查 tap 仓库权限

4. **Linux 包构建失败**
   - 检查 `maintainer` 字段格式
   - 确保包名符合规范

### 调试命令

```bash
# 检查配置
goreleaser check

# 详细输出
goreleaser release --debug

# 跳过发布，只构建
goreleaser build --clean
```

## 📚 更多资源

- [GoReleaser 文档](https://goreleaser.com/)
- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [Docker 多阶段构建](https://docs.docker.com/develop/dev-best-practices/dockerfile_best-practices/)
- [Semantic Versioning](https://semver.org/)

---

如有问题，请查看项目的 [Issues](https://github.com/your-org/kuanzhan/issues) 或创建新的 issue。 