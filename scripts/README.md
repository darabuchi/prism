# 构建脚本说明

## GoReleaser 本地构建

### 安装 GoReleaser

```bash
# macOS/Linux (使用 Homebrew)
brew install goreleaser

# 或使用 go install
go install github.com/goreleaser/goreleaser/v2@latest

# 验证安装
goreleaser --version
```

### 本地快照构建

生成本地构建，不推送到 GitHub：

```bash
# 清理之前的构建
goreleaser release --clean --snapshot

# 构建产物位于 dist/ 目录
ls -la dist/
```

### 测试发布流程

模拟完整的发布流程（不实际发布）：

```bash
# 跳过验证和发布
goreleaser release --skip=validate,publish --clean

# 仅构建不打包
goreleaser build --snapshot --clean
```

### 构建特定平台

```bash
# 仅构建当前平台
goreleaser build --single-target --snapshot --clean

# 构建 Linux
GOOS=linux GOARCH=amd64 goreleaser build --single-target --snapshot --clean

# 构建 Windows
GOOS=windows GOARCH=amd64 goreleaser build --single-target --snapshot --clean

# 构建 macOS
GOOS=darwin GOARCH=arm64 goreleaser build --single-target --snapshot --clean
```

### 检查配置

验证 `.goreleaser.yaml` 配置是否正确：

```bash
goreleaser check
```

### 发布新版本

创建 Git 标签并发布：

```bash
# 创建标签
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# 发布到 GitHub（需要 GITHUB_TOKEN 环境变量）
export GITHUB_TOKEN="your_github_token"
goreleaser release --clean
```

## Makefile 命令

项目 Makefile 提供了便捷的构建命令：

```bash
# 本地构建当前平台
make release-local

# 构建所有平台（快照）
make release-snapshot

# 检查 goreleaser 配置
make release-check

# 发布版本（需要 Git 标签）
make release
```

## 构建产物

GoReleaser 会生成以下内容：

```
dist/
├── prism_linux_amd64_v1/         # Linux AMD64 构建
│   └── prism
├── prism_darwin_amd64_v1/        # macOS Intel 构建
│   └── prism
├── prism_darwin_arm64/           # macOS Apple Silicon 构建
│   └── prism
├── prism_windows_amd64_v1/       # Windows AMD64 构建
│   └── prism.exe
├── prism_*.tar.gz                # 压缩包
├── prism_*.zip                   # Windows 压缩包
├── checksums.txt                 # SHA256 校验和
└── metadata.json                 # 构建元数据
```

## 环境变量

### GITHUB_TOKEN

用于发布到 GitHub Release：

```bash
# 创建 GitHub Personal Access Token
# https://github.com/settings/tokens

# 设置环境变量
export GITHUB_TOKEN="ghp_xxxxxxxxxxxx"

# 或使用 .envrc (direnv)
echo 'export GITHUB_TOKEN="ghp_xxxxxxxxxxxx"' > .envrc
direnv allow
```

### 其他环境变量

```bash
# Docker Registry 认证
export DOCKER_USERNAME="your-username"
export DOCKER_PASSWORD="your-password"

# Homebrew Tap Token
export HOMEBREW_TAP_GITHUB_TOKEN="ghp_xxxxxxxxxxxx"
```

## 版本号管理

遵循 [语义化版本](https://semver.org/lang/zh-CN/)：

```
主版本号.次版本号.修订号

例如: 1.2.3
- 1: 主版本号（不兼容的 API 修改）
- 2: 次版本号（向下兼容的功能性新增）
- 3: 修订号（向下兼容的问题修正）
```

### 创建标签

```bash
# 发布版本
git tag -a v1.0.0 -m "Release version 1.0.0"

# 预发布版本
git tag -a v1.0.0-beta.1 -m "Beta release"
git tag -a v1.0.0-rc.1 -m "Release candidate"

# 推送标签
git push origin v1.0.0
```

## 常见问题

### 1. 构建失败：go.sum 不一致

```bash
# 更新依赖
go mod tidy
go mod download

# 重新构建
goreleaser release --snapshot --clean
```

### 2. 权限错误

```bash
# 检查 GITHUB_TOKEN 权限
# 需要 repo 和 write:packages 权限
```

### 3. Docker 构建失败

```bash
# 确保 Docker 正在运行
docker info

# 登录 GitHub Container Registry
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin
```

## 参考文档

- [GoReleaser 官方文档](https://goreleaser.com/)
- [GoReleaser v2 迁移指南](https://goreleaser.com/blog/goreleaser-v2/)
- [GitHub Actions 集成](https://goreleaser.com/ci/actions/)
