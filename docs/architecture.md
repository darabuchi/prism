# Prism - 项目架构设计文档

## 项目概述

Prism 是基于 mihomo/clash 内核二次开发的代理软件，专门为节点池用户提供统一的多平台代理解决方案。

### 核心目标
- 支持 macOS、Windows、Linux、Android 四大平台
- 提供统一的节点池管理体验
- 简化用户配置流程
- 支持高级代理功能和规则管理

## 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                    Prism 生态系统                             │
├─────────────────────────────────────────────────────────────┤
│  桌面客户端         │  移动客户端         │  命令行工具          │
│  ┌─────────────┐   │  ┌─────────────┐   │  ┌─────────────┐    │
│  │ GUI (Tauri) │   │  │Android App  │   │  │    CLI      │    │
│  │             │   │  │             │   │  │             │    │
│  │ macOS       │   │  │ Java/Kotlin │   │  │    Go       │    │
│  │ Windows     │   │  │             │   │  │             │    │
│  │ Linux       │   │  └─────────────┘   │  └─────────────┘    │
│  └─────────────┘   │                   │                    │
├─────────────────────┴───────────────────┴────────────────────┤
│                      API 层 (RESTful)                        │
├─────────────────────────────────────────────────────────────┤
│                    Prism Core (Go)                           │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │            基于 mihomo/clash 内核扩展                    │ │
│  ├─────────────────────────────────────────────────────────┤ │
│  │ • 节点池管理     • 订阅管理     • 规则引擎              │ │
│  │ • 流量统计       • 延迟测试     • 配置同步              │ │
│  │ • 自动切换       • 负载均衡     • 日志记录              │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## 技术栈选择

### Core (prism-core)
- **语言**: Go 1.21+
- **基础**: mihomo/clash 内核
- **网络**: 基于 clash 的代理协议支持
- **存储**: SQLite + JSON 配置
- **API**: RESTful API + WebSocket

### 桌面客户端 (prism-desktop)
- **框架**: Tauri 2.0
- **前端**: React 18 + TypeScript
- **UI库**: Ant Design 或 Mantine
- **状态管理**: Zustand
- **打包**: 原生应用打包

### Android 客户端 (prism-android)
- **语言**: Kotlin
- **架构**: MVVM + Jetpack Compose
- **网络**: Retrofit + OkHttp
- **数据库**: Room
- **依赖注入**: Hilt

### 命令行工具 (prism-cli)
- **语言**: Go
- **CLI框架**: Cobra
- **配置**: Viper
- **输出**: 彩色终端输出

## 目录结构

```
prism/
├── README.md
├── docs/                    # 文档目录
│   ├── architecture.md
│   ├── api-specification.md
│   ├── development-guide.md
│   └── deployment-guide.md
├── core/                    # 核心代理引擎
│   ├── cmd/
│   ├── internal/
│   ├── pkg/
│   ├── config/
│   ├── go.mod
│   └── Dockerfile
├── desktop/                 # 桌面客户端 (Tauri)
│   ├── src-tauri/
│   ├── src/
│   ├── public/
│   ├── package.json
│   └── tauri.conf.json
├── android/                 # Android 客户端
│   ├── app/
│   ├── build.gradle
│   └── gradle/
├── cli/                     # 命令行工具
│   ├── cmd/
│   ├── internal/
│   ├── go.mod
│   └── Makefile
├── scripts/                 # 构建和部署脚本
│   ├── build.sh
│   ├── release.sh
│   └── docker/
└── .github/                 # CI/CD 配置
    └── workflows/
```

## 模块设计

### 1. Core 模块 (prism-core)

#### 核心组件
- **ProxyCore**: 代理核心，基于 mihomo 扩展
- **NodePoolManager**: 节点池管理器
- **SubscriptionManager**: 订阅管理器  
- **RuleEngine**: 规则引擎
- **ConfigManager**: 配置管理器
- **APIServer**: RESTful API 服务器

#### 关键特性
- 支持多种代理协议 (SS, SSR, V2Ray, Trojan, Hysteria)
- 智能节点选择和负载均衡
- 实时流量监控和统计
- 灵活的规则配置系统
- 配置热重载

### 2. 桌面客户端 (prism-desktop)

#### 功能模块
- **Dashboard**: 仪表盘，显示连接状态和流量
- **NodePool**: 节点池管理界面
- **Subscription**: 订阅管理
- **Rules**: 规则配置
- **Settings**: 应用设置
- **Logs**: 日志查看

#### 特色功能
- 系统托盘集成
- 开机自启动
- 代理模式快速切换
- 实时流量图表
- 一键导入配置

### 3. Android 客户端 (prism-android)

#### 核心功能
- VPN 模式代理 (VpnService)
- 节点池管理
- 订阅自动更新
- 流量统计
- 规则管理

#### Android 特有功能
- 按应用分流
- 省电模式优化
- 快捷方式支持
- 通知栏控制

### 4. 命令行工具 (prism-cli)

#### 主要命令
- `prism start`: 启动代理服务
- `prism stop`: 停止代理服务  
- `prism status`: 查看运行状态
- `prism config`: 配置管理
- `prism subscribe`: 订阅管理
- `prism test`: 节点测试

## 数据流设计

### 配置同步流程
```
用户操作 → 客户端 → API 调用 → Core 处理 → 配置更新 → 通知其他客户端
```

### 节点测试流程  
```
客户端请求 → Core 接收 → 批量测试 → 结果存储 → 实时反馈 → 界面更新
```

### 流量统计流程
```
代理流量 → Core 统计 → 数据库存储 → API 查询 → 客户端展示
```

## 安全考虑

### 数据安全
- 敏感配置加密存储
- API 访问鉴权 (JWT Token)
- 本地数据库加密
- 网络传输 TLS 加密

### 应用安全
- 代码签名 (桌面应用)
- 应用包完整性校验
- 自动更新安全验证
- 权限最小化原则

## 性能优化

### Core 优化
- 连接池复用
- 内存缓存优化
- 异步 I/O 处理
- 智能节点选择算法

### 客户端优化
- 懒加载和虚拟滚动
- 状态管理优化
- 网络请求缓存
- 界面响应优化

## 部署架构

### 开发环境
- 本地开发服务器
- 热重载支持
- 调试工具集成
- 模拟数据服务

### 生产环境
- Docker 容器化部署
- 多平台交叉编译
- 自动化测试流水线
- 版本管理和发布

这个架构设计为后续的开发提供了清晰的指导方向，确保各个组件能够协调工作，同时保持代码的可维护性和扩展性。