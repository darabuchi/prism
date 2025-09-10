# Prism

基于 mihomo/clash 内核二次开发的多平台代理客户端，专为节点池用户设计。

## 🚀 特性

- 🎯 **多平台支持**: macOS、Windows、Linux、Android
- 🔧 **节点池管理**: 统一管理多个订阅源和节点池
- ⚡ **高性能**: 基于 mihomo/clash 高性能代理内核
- 🎨 **现代 UI**: 使用 Tauri + React 构建的现代化桌面界面
- 📱 **移动优化**: 原生 Android 应用，支持 VPN 模式
- 🛠️ **CLI 工具**: 命令行工具支持服务器部署
- 🔄 **自动更新**: 支持订阅自动更新和节点测试
- 📊 **流量统计**: 实时流量监控和历史统计
- 🎛️ **规则管理**: 灵活的代理规则配置

## 📁 项目结构

```
prism/
├── core/           # Go 代理核心服务
├── desktop/        # Tauri 桌面客户端
├── android/        # Android 客户端
├── cli/            # 命令行工具
├── docs/           # 项目文档
└── scripts/        # 构建脚本
```

## 🏗️ 架构设计

Prism 采用模块化架构设计：

- **Core**: 基于 mihomo 的代理核心，提供 RESTful API
- **Desktop**: 使用 Tauri 构建的跨平台桌面应用
- **Android**: 使用 Kotlin + Jetpack Compose 的原生应用
- **CLI**: Go 编写的命令行管理工具

## 🛠️ 技术栈

### Core
- **Go 1.21+** - 核心语言
- **mihomo/clash** - 代理内核
- **Gin** - Web 框架
- **SQLite** - 数据存储

### Desktop
- **Tauri 2.0** - 应用框架
- **React 18** - 前端框架
- **TypeScript** - 类型安全
- **Ant Design** - UI 组件库

### Android
- **Kotlin** - 开发语言
- **Jetpack Compose** - UI 框架
- **MVVM** - 架构模式
- **Room** - 数据库

### CLI
- **Go** - 开发语言
- **Cobra** - CLI 框架

## 📚 文档

- [架构设计](docs/architecture.md) - 系统架构和技术选型
- [API 规格](docs/api-specification.md) - RESTful API 接口文档
- [开发指南](docs/development-guide.md) - 多平台开发详细指南
- [部署指南](docs/deployment-guide.md) - 部署和分发说明

## 🚀 快速开始

### 开发环境准备

1. **安装基础依赖**
   ```bash
   # macOS
   brew install go node rust
   
   # Ubuntu/Debian
   sudo apt install golang nodejs npm rustc
   
   # Windows (使用 winget)
   winget install GoLang.Go
   winget install OpenJS.NodeJS
   winget install Rustlang.Rustup
   ```

2. **克隆项目**
   ```bash
   git clone https://github.com/yourusername/prism.git
   cd prism
   ```

### 运行开发版本

1. **启动 Core 服务**
   ```bash
   cd core
   go mod tidy
   go run cmd/prism-core/main.go
   ```

2. **启动桌面客户端**
   ```bash
   cd desktop
   npm install
   npm run tauri dev
   ```

3. **构建 Android 应用**
   ```bash
   cd android
   ./gradlew assembleDebug
   ```

### 生产构建

```bash
# 使用构建脚本
./scripts/build.sh

# 或使用 Docker
docker-compose up --build
```

## 📖 使用说明

### 桌面客户端

1. 启动应用后，首先配置 Core 服务地址（默认：`http://localhost:9090`）
2. 添加节点池或订阅链接
3. 选择代理模式（直连/全局/规则）
4. 选择节点或让系统自动选择最优节点

### Android 应用

1. 安装 APK 文件
2. 授予 VPN 权限
3. 配置节点池和规则
4. 启用 VPN 服务

### CLI 工具

```bash
# 查看状态
prism status

# 启动代理
prism start

# 切换节点
prism node select <node-id>

# 更新订阅
prism subscribe update
```

## 🤝 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目基于 [MIT License](LICENSE) 开源协议。

## 🙏 致谢

- [mihomo](https://github.com/MetaCubeX/mihomo) - 高性能代理内核
- [clash](https://github.com/Dreamacro/clash) - 原始 clash 项目
- [Tauri](https://tauri.app/) - 现代化桌面应用框架

## 📞 联系我们

- 提交 Issue: [GitHub Issues](https://github.com/yourusername/prism/issues)
- 讨论: [GitHub Discussions](https://github.com/yourusername/prism/discussions)

---

⭐ 如果这个项目对你有帮助，请给我们一个 Star！