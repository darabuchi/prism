# Prism - 智能代理管理系统

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Prism 是一个智能代理管理系统，支持订阅管理、节点测试、智能路由等功能。

## ✨ 特性

- 🚀 **订阅管理**: 支持多种订阅格式（Clash、V2Ray、Surge）
- 🔍 **智能测试**: 自动测试节点延迟、速度、解锁能力
- 🎯 **智能路由**: 基于规则的智能路由选择
- 📊 **监控告警**: Prometheus 监控 + Grafana 可视化
- 🔐 **安全可靠**: 支持 MITM、证书管理、熔断保护
- 🌍 **多语言支持**: 内置 7 种语言（英语、简体中文、繁体中文、法语、俄语、西班牙语、阿拉伯语）
- 🐳 **易于部署**: Docker、Systemd、Kubernetes 多种部署方式

## 📋 目录结构

```
prism/
├── cmd/                    # 命令行入口
├── prism/                  # 核心包
├── internal/               # 内部实现
│   ├── state/             # 状态管理
│   └── api/               # HTTP API
├── pkg/                    # 公共包
├── config/                 # 配置文件
├── docs/                   # 文档
└── deployments/            # 部署配置
```

详细说明见 [项目结构文档](docs/项目结构.md)

## 🚀 快速开始

### 前置要求

- Go 1.21+
- SQLite3（或 MySQL/PostgreSQL）

### 安装

```bash
# 克隆仓库
git clone https://github.com/ice-cream-heaven/prism.git
cd prism

# 下载依赖
go mod download

# 编译
make build

# 运行
make run
```

### Docker 部署

```bash
cd deployments/docker
docker-compose up -d
```

## 📖 文档

- [系统设计文档](docs/系统设计文档.md)
- [项目结构](docs/项目结构.md)
- [API 文档](docs/API文档.md)
- [部署指南](docs/详细设计/部署与运维.md)
- [错误码注册表](error_code.json) - 标准化错误码定义

### 核心包文档

- [i18n 国际化](pkg/i18n/README.md) - 多语言支持（7种语言）
- [Queue 队列系统](pkg/queue/README.md) - 高性能消息队列

### 详细设计

- [数据库访问层](docs/详细设计/数据库访问层.md)
- [缓存系统](docs/详细设计/缓存系统.md)
- [错误处理](docs/详细设计/错误处理.md)
- [并发控制](docs/详细设计/并发控制.md)
- [测试策略](docs/详细设计/测试策略.md)

## 🛠️ 开发

### 运行测试

```bash
# 所有测试
make test

# 单元测试
make test-unit

# 集成测试
make test-integration

# 覆盖率报告
make test-coverage
```

### 代码检查

```bash
# 格式化
make fmt

# Lint
make lint
```

## 🏗️ 架构

```
┌─────────────────────────────────────┐
│          HTTP API / SOCKS5          │
├─────────────────────────────────────┤
│        Business Logic Layer         │
│  ┌──────────┬──────────┬─────────┐ │
│  │Subscription│  Node   │  Route  │ │
│  └──────────┴──────────┴─────────┘ │
├─────────────────────────────────────┤
│         Data Access Layer           │
│  ┌──────────┬──────────┬─────────┐ │
│  │ Database │  Cache   │  Queue  │ │
│  └──────────┴──────────┴─────────┘ │
└─────────────────────────────────────┘
```

## 📊 监控

访问 `http://localhost:9090/metrics` 查看 Prometheus 指标。

## 🤝 贡献

欢迎贡献！请查看 [贡献指南](CONTRIBUTING.md)。

## 📄 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 🙏 致谢

- [LRPC](https://github.com/lazygophers/lrpc) - 框架支持
- [Gin](https://github.com/gin-gonic/gin) - HTTP 框架
- [GORM](https://gorm.io/) - ORM 库

## 📮 联系方式

- Issue: [GitHub Issues](https://github.com/ice-cream-heaven/prism/issues)
- Email: support@ice-cream-heaven.com
