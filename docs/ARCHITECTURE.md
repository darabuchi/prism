# Prism 技术架构设计

## 整体架构

Prism 采用分层架构设计，由代理核心、API 服务层、Web 界面和客户端组成。

```
┌─────────────────────────────────────────────────────────────┐
│                    客户端层                                    │
├─────────────┬─────────────┬─────────────┬─────────────────────┤
│  Web 界面   │  桌面客户端  │  移动端应用  │  第三方客户端        │
│  (React)    │  (Tauri)    │ (Flutter)   │  (API 调用)         │
└─────────────┴─────────────┴─────────────┴─────────────────────┘
                              │
                    ┌─────────┴─────────┐
                    │    API 网关层      │
                    │   (HTTP/WebSocket) │
                    └─────────┬─────────┘
                              │
┌─────────────────────────────┴─────────────────────────────────┐
│                      API 服务层                               │
├─────────────┬─────────────┬─────────────┬─────────────────────┤
│  配置管理    │  状态监控    │  规则管理    │  节点管理            │
│  Service    │  Service    │  Service    │  Service            │
└─────────────┴─────────────┴─────────────┴─────────────────────┘
                              │
┌─────────────────────────────┴─────────────────────────────────┐
│                     代理核心层                                 │
├─────────────┬─────────────┬─────────────┬─────────────────────┤
│ Mihomo/Clash│  配置解析    │  流量转发    │  规则匹配           │
│    核心      │    模块      │    引擎      │    引擎             │
└─────────────┴─────────────┴─────────────┴─────────────────────┘
```

## 核心组件

### 1. 代理核心层 (Core Layer)

基于 mihomo/clash 构建，负责实际的代理功能。

**主要职责:**
- 网络流量代理转发
- 支持多种协议 (HTTP, HTTPS, SOCKS5, VMess, VLESS, Trojan 等)
- 规则匹配和路由
- 负载均衡和故障转移
- 流量统计和监控

**技术选型:**
- **语言**: Go 1.21+
- **核心库**: mihomo/clash core
- **网络库**: 基于 Go 标准库 + 优化
- **配置格式**: YAML

**核心模块:**
```go
// 核心接口设计
type ProxyCore interface {
    Start() error
    Stop() error
    Reload(config *Config) error
    GetConnections() []Connection
    GetProxies() []Proxy
    TestDelay(proxy string) (int64, error)
}
```

### 2. API 服务层 (Service Layer)

提供 RESTful API 和 WebSocket 服务，连接核心层和客户端。

**主要职责:**
- 配置管理 API
- 实时状态监控
- 节点测速和健康检查
- 规则和订阅管理
- 用户认证和权限控制

**技术选型:**
- **Web框架**: GoFiber v2 (高性能、Express风格、零内存分配路由)
- **数据库**: SQLite/MySQL/PostgreSQL/GaussDB (多数据库支持)
- **ORM**: GORM (与GoFiber完美集成)
- **缓存**: BoltDB/LevelDB (嵌入式键值存储，根据场景选择)
- **认证**: JWT Token (使用fiber/jwt中间件)
- **实时通信**: GoFiber WebSocket (github.com/gofiber/websocket/v2)
- **中间件**: GoFiber官方中间件生态 (cors, logger, recover, compress等)
- **配置存储**: YAML/JSON 文件
- **日志**: github.com/lazygophers/log
- **工具包**: github.com/lazygophers/utils (json, stringx, xtime, bufiox, randx, anyx, candy)
- **原子操作**: go.uber.org/atomic

**GoFiber API 设计:**
```go
// 使用 GoFiber 路由设计
app := fiber.New(fiber.Config{
    Prefork:       false,
    CaseSensitive: true,
    StrictRouting: true,
    ServerHeader:  "Prism",
})

// API 路由
api := app.Group("/api/v1")
api.Get("/status", handlers.GetSystemStatus)              // 获取系统状态
api.Get("/proxies", handlers.GetProxyList)               // 获取代理列表
api.Post("/proxies/:name/test", handlers.TestProxyDelay) // 测试代理延迟
api.Get("/connections", handlers.GetConnections)         // 获取连接信息
api.Post("/config/reload", handlers.ReloadConfig)        // 重载配置
api.Get("/rules", handlers.GetRuleList)                  // 获取规则列表
api.Post("/subscriptions", handlers.ManageSubscription)   // 管理订阅

// WebSocket 路由
app.Get("/ws", websocket.New(handlers.WebSocketHandler))  // 实时数据推送
```

### 3. Web 界面层 (Frontend Layer)

现代化的 Web 管理界面，提供直观的用户体验。

**主要功能:**
- 代理节点管理和测速
- 实时流量监控和统计
- 规则配置和订阅管理
- 系统设置和主题切换
- 响应式设计，支持移动端

**技术选型:**
- **框架**: React 18 + TypeScript
- **状态管理**: Zustand (轻量) 或 Redux Toolkit
- **UI 组件**: Ant Design 或 Material-UI
- **图表库**: ECharts 或 Chart.js
- **构建工具**: Vite
- **样式**: Tailwind CSS + CSS Modules

**组件架构:**
```
src/
├── components/           # 通用组件
├── pages/               # 页面组件
├── stores/              # 状态管理
├── services/            # API 服务
├── hooks/               # 自定义 Hooks
├── utils/               # 工具函数
└── types/               # TypeScript 类型
```

## 数据流设计

### 配置数据流
```
用户配置 → Web界面 → API服务 → 配置验证 → 核心重载 → 状态更新
```

### 监控数据流
```
代理核心 → 统计数据 → API服务 → WebSocket → 前端实时更新
```

### 节点管理流
```
订阅链接 → 解析节点 → 存储配置 → 健康检查 → 状态展示
```

## 部署架构

### 单机部署 (推荐)
```
┌─────────────────────────────────┐
│         Prism 服务               │
├─────────────────────────────────┤
│  Web界面 (内嵌静态文件)          │
│  API服务 (端口: 9090)           │
│  代理核心 (端口: 7890/7891)      │
│  配置文件 (~/.prism/config/)    │
└─────────────────────────────────┘
```

### Docker 部署
```yaml
# docker-compose.yml
version: '3.8'
services:
  prism:
    image: prism:latest
    ports:
      - "9090:9090"    # API + Web
      - "7890:7890"    # HTTP 代理
      - "7891:7891"    # SOCKS5 代理
    volumes:
      - ./config:/app/config
      - ./data:/app/data
```

## 安全设计

### 认证授权
- Web 界面支持用户名/密码认证
- API 访问使用 JWT Token
- 支持白名单 IP 访问控制

### 数据安全
- 敏感配置文件加密存储
- HTTPS/TLS 传输加密
- 定期安全更新和漏洞修复

### 权限控制
```go
type Permission struct {
    Read   bool  // 查看配置和状态
    Write  bool  // 修改配置
    Admin  bool  // 系统管理
}
```

## 性能优化

### 核心优化
- 连接池复用
- 内存缓存热点数据
- 异步 I/O 操作
- 智能路由算法

### API 优化
- 响应数据压缩
- 合理的缓存策略
- 分页查询大数据集
- WebSocket 减少轮询

### 前端优化
- 组件懒加载
- 虚拟滚动长列表
- 图片和资源压缩
- Service Worker 缓存

## 扩展性设计

### 插件系统
```go
type Plugin interface {
    Name() string
    Init(ctx context.Context) error
    HandleTraffic(conn net.Conn) error
    Cleanup() error
}
```

### 配置扩展
- 支持多种配置格式 (YAML, JSON, TOML)
- 配置模板和预设
- 动态配置热重载

### API 扩展
- 版本化 API (v1, v2)
- 插件化中间件
- 自定义认证提供商

## 监控和日志

### 指标监控
- 连接数和流量统计
- 代理节点延迟和可用性
- 系统资源使用情况
- 错误率和异常统计

### 日志系统
```go
// 使用 github.com/lazygophers/log
import "github.com/lazygophers/log"

// 日志记录示例
log.Debug("调试信息", log.String("module", "core"))
log.Info("系统启动", log.Int("port", 9090))
log.Warn("配置警告", log.String("file", "config.yaml"))
log.Error("系统错误", log.Error(err))

// 结构化日志字段
log.Info("节点测速完成", 
    log.String("node", nodeName),
    log.Int64("delay", delay),
    log.String("status", "success"),
)
```

### 健康检查
- HTTP 健康检查端点
- 自动故障检测和恢复
- 服务依赖检查

## 开发和测试

### 开发环境
- 热重载开发服务器
- 模拟数据和 Mock 服务
- 开发工具集成

### 测试策略
- 单元测试覆盖率 > 80%
- 集成测试关键流程
- 端到端测试用户场景
- 性能测试和压力测试

### CI/CD 流程
```yaml
# GitHub Actions 示例
name: CI/CD
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run tests
        run: make test
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Build artifacts
        run: make build-all
```