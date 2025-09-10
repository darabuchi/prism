# Prism 代理核心

基于 [mihomo](https://github.com/MetaCubeX/mihomo) 和 [clash](https://github.com/Dreamacro/clash) 二次开发的代理核心项目，专为节点池用户提供便捷的代理服务。

## 项目概述

Prism 是一个现代化的代理核心解决方案，旨在为节点池用户提供稳定、高效、易用的代理服务。项目采用渐进式开发策略，优先完成 Web 版本，随后扩展到多个平台。

### 核心特性

- 🔥 基于成熟的 mihomo/clash 核心
- 🌐 优先支持 Web 界面管理
- 📱 计划支持多平台客户端 (macOS, Windows, Linux, Android)
- 🎯 专为节点池用户优化
- ⚡ 高性能代理转发
- 🔧 灵活的配置管理
- 📊 详细的连接统计

### 目标平台

- **第一阶段**: Web 管理界面
- **第二阶段**: 桌面客户端 (macOS, Windows, Linux)
- **第三阶段**: 移动端应用 (Android)

## 技术栈

### 核心引擎
- **代理核心**: mihomo/clash (Go)
- **API 服务**: Go + Gin/Echo
- **配置管理**: YAML + JSON

### Web 界面
- **前端框架**: React/Vue.js (待定)
- **状态管理**: Redux/Vuex
- **UI 组件**: Material-UI/Ant Design
- **构建工具**: Vite/Webpack

### 跨平台客户端
- **桌面端**: Tauri + Rust + Web 技术
- **移动端**: Flutter/React Native (待定)

## 项目结构

```
prism/
├── core/           # 代理核心 (基于 mihomo/clash)
├── api/            # RESTful API 服务
├── web/            # Web 管理界面
├── desktop/        # 桌面客户端 (未来)
├── mobile/         # 移动端应用 (未来)
├── docs/           # 项目文档
└── scripts/        # 构建和部署脚本
```

## 开发状态

🚧 项目处于早期规划阶段

- [x] 项目初始化
- [ ] 核心架构设计
- [ ] Web 版本开发
- [ ] 桌面客户端开发
- [ ] 移动端应用开发

## 快速开始

> 注意：项目尚在开发中，以下是预期的使用方式

### Web 版本
```bash
# 克隆项目
git clone https://github.com/ice-cream-heaven/prism.git
cd prism

# 启动核心服务
cd core && go run .

# 启动 Web 界面
cd web && npm install && npm run dev
```

### API 接口
```bash
# 获取代理状态
curl http://localhost:9090/api/status

# 更新配置
curl -X POST http://localhost:9090/api/config
```

## 贡献指南

我们欢迎社区贡献！请查看 [CONTRIBUTING.md](./CONTRIBUTING.md) 了解详细信息。

### 开发环境要求

- Go 1.21+
- Node.js 18+
- Git

### 提交规范

使用 [Conventional Commits](https://conventionalcommits.org/) 规范：

- `feat:` 新功能
- `fix:` 修复
- `docs:` 文档更新
- `style:` 代码格式
- `refactor:` 重构
- `test:` 测试相关
- `chore:` 构建过程或辅助工具的变动

## 许可证

本项目基于 [MIT License](./LICENSE) 开源。

## 致谢

- [mihomo](https://github.com/MetaCubeX/mihomo) - 提供强大的代理核心
- [clash](https://github.com/Dreamacro/clash) - 原始代理核心实现

## 联系我们

- 项目主页: https://github.com/ice-cream-heaven/prism
- 问题反馈: https://github.com/ice-cream-heaven/prism/issues
- 讨论区: https://github.com/ice-cream-heaven/prism/discussions

---

**注意**: 本项目仅用于学习和合法的网络代理用途，请遵守当地法律法规。