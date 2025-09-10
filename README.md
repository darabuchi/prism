# Prism 代理核心

基于 [mihomo](https://github.com/MetaCubeX/mihomo) 和 [clash](https://github.com/Dreamacro/clash) 二次开发的代理核心项目，专为节点池用户提供便捷的代理服务。

## 项目概述

Prism 是一个现代化的代理核心解决方案，旨在为节点池用户提供稳定、高效、易用的代理服务。项目采用渐进式开发策略，优先完成 Web 版本，随后扩展到多个平台。

### 核心特性

- 🔥 基于成熟的 mihomo/clash 核心
- 🌐 简洁直观的 Web 管理界面
- 📱 计划支持多平台客户端 (macOS, Windows, Linux, Android)
- 🎯 专为节点池用户优化，无复杂用户管理
- ⚡ 高性能代理转发和智能负载均衡
- 🔄 订阅自动抓取和节点多对多关系管理
- 📄 规则文件化存储，灵活配置代理规则
- 📊 详细的流量统计和节点性能监控

### 设计理念

- **性能优先**: 作为代理客户端，性能是第一优先级
- **最小化数据库**: 仅3个核心表，移除影响性能的日志和统计表
- **简化优先**: 去除用户管理和代理组，专注核心功能
- **订阅中心**: 以订阅为核心，自动管理节点池
- **多对多关系**: 节点可属于多个订阅，订阅可包含多个节点
- **文件存储**: 规则和日志文件化存储，缓存使用BoltDB/LevelDB
- **自动化**: 订阅自动更新，节点自动测速

### 目标平台

- **第一阶段**: Web 管理界面 (简化架构)
- **第二阶段**: 桌面客户端 (macOS, Windows, Linux)
- **第三阶段**: 移动端应用 (Android)

## 技术栈

### 核心引擎
- **代理核心**: mihomo/clash (Go)
- **API 服务**: Go + GoFiber v2 (高性能、零内存分配)
- **数据库**: SQLite/MySQL/PostgreSQL/GaussDB (多数据库支持)
- **缓存**: BoltDB/LevelDB (嵌入式键值存储)

### Web 界面
- **前端框架**: React 18 + TypeScript
- **状态管理**: Zustand (轻量级)
- **UI 组件**: Ant Design 5.x
- **构建工具**: Vite
- **样式**: Tailwind CSS + CSS Modules

### 跨平台客户端
- **桌面端**: Tauri + Rust + Web 技术
- **移动端**: Flutter (计划)

## 项目结构

```
prism/
├── cmd/            # 应用入口点
├── internal/       # 私有应用代码
│   ├── api/        # GoFiber API 处理器
│   ├── core/       # 代理核心集成
│   ├── models/     # 数据模型
│   ├── service/    # 业务逻辑层
│   └── database/   # 数据库操作
├── web/            # React + TypeScript Web界面
├── data/           # 数据文件
│   ├── rules/      # 规则文件存储
│   ├── configs/    # 配置文件
│   └── logs/       # 日志文件
├── docs/           # 项目文档
├── scripts/        # 构建和部署脚本
└── migrations/     # 数据库迁移文件
```

## 开发状态

🚧 项目处于早期规划阶段

- [x] 项目初始化和架构设计
- [x] 数据库设计 (简化版)
- [x] GoFiber API 框架设计
- [x] React + TypeScript 前端架构
- [ ] 订阅管理系统开发
- [ ] 节点管理和测速系统
- [ ] 规则文件管理系统
- [ ] Web 版本完整开发
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