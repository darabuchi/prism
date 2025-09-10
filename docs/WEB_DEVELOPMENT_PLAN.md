# Prism Web 版本开发计划

## 开发概述

Web 版本作为 Prism 项目的第一阶段，将提供完整的代理核心管理功能，包括节点管理、规则配置、流量监控等核心特性。本计划详细规划了 Web 版本的开发流程、技术实现和交付时间。

## 项目目标

### 核心目标
- 🎯 提供直观易用的 Web 管理界面
- ⚡ 实现高性能的代理核心服务
- 📊 支持实时流量监控和统计
- 🔧 完善的节点池管理功能
- 🌐 响应式设计，支持移动端访问

### 用户价值
- **节点池用户**: 方便管理大量代理节点
- **个人用户**: 简化代理配置和使用
- **高级用户**: 提供详细的控制和监控选项

## 技术栈确定

### 后端技术栈
```go
// 核心技术选择
Framework: GoFiber v2 (高性能、Express风格、零内存分配)
Database: SQLite(默认) / MySQL / PostgreSQL / GaussDB
Cache: BoltDB / LevelDB (嵌入式键值存储)
ORM: GORM (与GoFiber完美集成)
Auth: JWT Token (fiber/jwt中间件)
WebSocket: GoFiber WebSocket (github.com/gofiber/websocket/v2)
Middleware: GoFiber官方中间件生态 (cors, logger, recover, compress)
Config: Viper (配置管理)
Logging: github.com/lazygophers/log
Utils: github.com/lazygophers/utils (json, stringx, xtime, bufiox, randx, anyx, candy)
Atomic: go.uber.org/atomic
Testing: Testify
```

### 前端技术栈
```typescript
// 核心技术选择
Framework: React 18 + TypeScript (强类型、现代化React)
StateManagement: Zustand (轻量级状态管理)
UILibrary: Ant Design 5.x (企业级UI组件)
BuildTool: Vite (极速构建工具)
Styling: 
  - Tailwind CSS (原子化CSS)
  - CSS Modules (模块化样式)
  - Ant Design主题定制
Routing: React Router v6 (最新路由系统)
HTTPClient: Axios (HTTP请求)
WebSocket: native WebSocket API (实时通信)
Charts: ECharts (数据可视化)
Forms: React Hook Form + Yup (表单管理和验证)
DateHandling: dayjs (日期处理)
Icons: Ant Design Icons + Lucide React
Testing: Jest + React Testing Library (单元测试)
DevTools: React DevTools + Redux DevTools
```

## 开发阶段规划

### 🔨 阶段 1: 项目基础搭建 (Week 1-2)

#### 后端基础架构
**时间**: 5-7 天
**负责人**: 后端开发

**任务清单**:
- [ ] Go 项目结构搭建 (标准 Go 布局)
- [ ] GoFiber v2 框架集成和中间件配置
- [ ] 多数据库支持设计 (SQLite/MySQL/PostgreSQL/GaussDB) 和 GORM 集成
- [ ] JWT 认证中间件实现 (使用 fiber/jwt)
- [ ] GoFiber 路由系统设计和 API 端点定义
- [ ] lazygophers/log 日志系统集成
- [ ] lazygophers/utils 工具包集成
- [ ] GoFiber 中间件配置 (CORS, Logger, Recover, Compress)
- [ ] Docker 化配置

**交付物**:
```
prism/
├── cmd/
│   └── server/
│       └── main.go           # 程序入口
├── internal/
│   ├── api/                  # API 路由和处理器
│   ├── config/               # 配置管理
│   ├── core/                 # 代理核心集成
│   ├── database/             # 数据库模型
│   ├── middleware/           # 中间件
│   └── service/              # 业务逻辑
├── pkg/                      # 公共包
├── configs/                  # 配置文件
├── scripts/                  # 构建脚本
└── Dockerfile
```

#### 前端基础架构
**时间**: 3-5 天
**负责人**: 前端开发

**任务清单**:
- [ ] React 18 + TypeScript 项目初始化 (使用 create-vite)
- [ ] Vite 构建配置优化 (环境变量、代理、分包)
- [ ] Ant Design 5.x 集成和主题定制
- [ ] React Router v6 路由系统设计
- [ ] Zustand 状态管理架构设计
- [ ] Axios HTTP 客户端封装 (拦截器、错误处理)
- [ ] TypeScript 类型定义建立
- [ ] 基础组件库和 Hooks 建立
- [ ] Tailwind CSS + CSS Modules 样式系统
- [ ] 响应式布局框架 (移动端适配)

**交付物**:
```
web/
├── src/
│   ├── components/           # 通用组件
│   ├── pages/               # 页面组件
│   ├── stores/              # 状态管理
│   ├── services/            # API 服务
│   ├── hooks/               # 自定义 Hooks
│   ├── utils/               # 工具函数
│   ├── types/               # TypeScript 类型
│   └── assets/              # 静态资源
├── public/
├── index.html
├── vite.config.ts
└── package.json
```

### ⚙️ 阶段 2: 核心功能开发 (Week 3-8)

#### 2.1 代理核心集成 (Week 3-4)
**主要目标**: 集成 mihomo/clash 核心，实现基础代理功能

**后端任务**:
- [ ] mihomo/clash 核心库集成
- [ ] 配置文件解析和验证
- [ ] 代理服务启动/停止控制
- [ ] 核心状态监控接口
- [ ] 配置热重载机制

```go
// 核心服务接口设计 - 使用推荐的包
import (
    "context"
    "github.com/gofiber/fiber/v2"
    "github.com/lazygophers/log"
    "github.com/lazygophers/utils/json"
    "github.com/lazygophers/utils/stringx"
    "github.com/lazygophers/utils/xtime"
    "go.uber.org/atomic"
)

type CoreService interface {
    Start(config *Config) error
    Stop() error
    Reload(config *Config) error
    GetStatus() *Status
    GetProxies() []*Proxy
    GetConnections() []*Connection
}

// 使用 GoFiber v2 的完整API示例
func SetupAPI(app *fiber.App) {
    // 全局中间件
    app.Use(cors.New(cors.Config{
        AllowOrigins: "*",
        AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
        AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    }))
    
    app.Use(compress.New())
    app.Use(recover.New())
    
    // 自定义日志中间件
    app.Use(func(c *fiber.Ctx) error {
        start := xtime.Now()
        
        err := c.Next()
        
        duration := xtime.Since(start)
        log.Info("HTTP请求", 
            log.String("method", c.Method()),
            log.String("path", c.Path()),
            log.Int("status", c.Response().StatusCode()),
            log.Duration("duration", duration),
            log.String("ip", c.IP()),
            log.String("userAgent", c.Get("User-Agent")),
        )
        
        return err
    })
    
    // API 路由组
    api := app.Group("/api/v1")
    
    // JWT 中间件保护的路由
    protected := api.Use(jwtware.New(jwtware.Config{
        SigningKey: []byte("your-secret-key"),
        ContextKey: "jwt",
        ErrorHandler: func(c *fiber.Ctx, err error) error {
            return c.Status(401).JSON(fiber.Map{
                "error": "未授权访问",
                "code": "UNAUTHORIZED",
            })
        },
    }))
    
    // 节点管理路由
    nodes := protected.Group("/nodes")
    nodes.Get("/", handlers.GetNodeList)           // 获取节点列表
    nodes.Post("/", handlers.CreateNode)           // 创建节点
    nodes.Get("/:id", handlers.GetNode)            // 获取单个节点
    nodes.Put("/:id", handlers.UpdateNode)         // 更新节点
    nodes.Delete("/:id", handlers.DeleteNode)     // 删除节点
    nodes.Post("/:id/test", handlers.TestNode)    // 测试节点
    nodes.Post("/batch", handlers.BatchOperation) // 批量操作
    
    // 订阅管理路由
    subs := protected.Group("/subscriptions")
    subs.Get("/", handlers.GetSubscriptionList)
    subs.Post("/", handlers.CreateSubscription)
    subs.Put("/:id", handlers.UpdateSubscription)
    subs.Delete("/:id", handlers.DeleteSubscription)
    subs.Post("/:id/update", handlers.UpdateSubscriptionNodes)
    
    // WebSocket 路由 (实时数据推送)
    app.Get("/ws", websocket.New(func(c *websocket.Conn) {
        handlers.HandleWebSocket(c)
    }))
}

// 节点列表处理器示例
func GetNodeList(c *fiber.Ctx) error {
    // 分页参数
    page := c.QueryInt("page", 1)
    size := c.QueryInt("size", 20)
    
    // 过滤参数
    filter := &models.NodeFilter{
        Type:        c.Query("type"),
        CountryCode: c.Query("country"),
        Status:      c.Query("status"),
        Enabled:     c.QueryBool("enabled", true),
    }
    
    // 排序参数
    sort := c.Query("sort", "created_at")
    order := c.Query("order", "desc")
    
    // 调用服务层
    result, err := services.NodeService.GetNodes(c.Context(), page, size, filter, sort, order)
    if err != nil {
        log.Error("获取节点列表失败", log.Error(err))
        return c.Status(500).JSON(fiber.Map{
            "error": "内部服务器错误",
            "code":  "INTERNAL_ERROR",
        })
    }
    
    return c.JSON(fiber.Map{
        "code":    200,
        "message": "获取成功",
        "data":    result.Data,
        "pagination": fiber.Map{
            "page":       page,
            "size":       size,
            "total":      result.Total,
            "totalPages": (result.Total + int64(size) - 1) / int64(size),
        },
        "timestamp": xtime.Now().Unix(),
    })
}

// 数据模型示例 - 使用 GORM 支持多数据库
type Node struct {
    ID       string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
    Name     string    `gorm:"type:varchar(255);not null" json:"name"`
    Type     string    `gorm:"type:varchar(50);not null" json:"type"`
    Server   string    `gorm:"type:varchar(255);not null" json:"server"`
    Port     int       `gorm:"not null" json:"port"`
    Config   string    `gorm:"type:text" json:"config"` // JSON字符串
    Delay    int64     `gorm:"default:-1" json:"delay"`
    Status   string    `gorm:"type:varchar(20);default:inactive" json:"status"`
    CreatedAt xtime.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt xtime.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
```

**API 端点**:
```
POST /api/core/start        # 启动代理核心
POST /api/core/stop         # 停止代理核心
POST /api/core/reload       # 重载配置
GET  /api/core/status       # 获取运行状态
```

#### 2.2 订阅管理系统 (Week 4-5)
**主要目标**: 实现订阅的自动抓取和节点管理

**后端任务**:
- [ ] 订阅数据模型设计
- [ ] 订阅 CRUD API 实现
- [ ] 订阅链接解析引擎 (支持Clash、V2Ray、SS等格式)
- [ ] 自动抓取和更新机制
- [ ] 节点去重和合并逻辑
- [ ] 多对多关系管理

```go
type Subscription struct {
    ID                  string    `json:"id"`
    Name                string    `json:"name"`
    URL                 string    `json:"url"`
    Type                string    `json:"type"`     // clash, v2ray, ss
    NodeCount           int       `json:"node_count"`
    LastUpdateAt        time.Time `json:"last_update_at"`
    UpdateIntervalHours int       `json:"update_interval_hours"`
    AutoUpdate          bool      `json:"auto_update"`
    Status              string    `json:"status"`   // active, error, updating
    Nodes               []Node    `json:"nodes"`    // 多对多关系
}
```

**前端任务**:
- [ ] 订阅列表页面开发
- [ ] 订阅添加/编辑表单
- [ ] 订阅更新状态显示
- [ ] 关联节点管理界面
- [ ] 自动更新配置

**API 端点**:
```
GET    /api/subscriptions           # 获取订阅列表
POST   /api/subscriptions           # 添加订阅
GET    /api/subscriptions/{id}      # 获取单个订阅
PUT    /api/subscriptions/{id}      # 更新订阅
DELETE /api/subscriptions/{id}      # 删除订阅
POST   /api/subscriptions/{id}/update # 手动更新订阅
GET    /api/subscriptions/{id}/nodes  # 获取订阅关联的节点
```

#### 2.3 节点管理系统 (Week 5-6)
**主要目标**: 实现节点的测速、监控和多对多关系管理

**后端任务**:
- [ ] 节点测速引擎开发
- [ ] 节点状态监控
- [ ] 多订阅节点去重逻辑
- [ ] 节点地理信息识别
- [ ] 性能统计和成功率计算

```go
type ProxyNode struct {
    ID            string        `json:"id"`
    Name          string        `json:"name"`
    Type          string        `json:"type"`     // vmess, vless, trojan, ss
    Server        string        `json:"server"`
    Port          int           `json:"port"`
    DelayMS       int           `json:"delay_ms"`
    Status        string        `json:"status"`   // online, offline, testing
    SuccessRate   float64       `json:"success_rate"`
    Subscriptions []Subscription `json:"subscriptions"` // 多对多关系
}
```

**前端任务**:
- [ ] 节点列表页面（支持多订阅显示）
- [ ] 节点测速和状态显示
- [ ] 节点详情页面（显示所属订阅）
- [ ] 批量测速功能
- [ ] 地理位置和性能图表

**API 端点**:
```
GET    /api/nodes                # 获取节点列表
GET    /api/nodes/{id}          # 获取单个节点
POST   /api/nodes/{id}/test     # 测试节点延迟
POST   /api/nodes/batch-test    # 批量测试节点
GET    /api/nodes/{id}/subscriptions # 获取节点所属订阅
```

#### 2.4 规则文件系统 (Week 6-7)
**主要目标**: 实现基于文件的规则存储和管理

**后端任务**:
- [ ] 规则文件读写管理
- [ ] 规则编译和生成
- [ ] 内置规则预设
- [ ] 远程规则同步
- [ ] 规则文件热重载

```go
type RuleService struct {
    rulesDir string
}

type Rules struct {
    Direct []string `json:"direct"`
    Proxy  []string `json:"proxy"`
    Reject []string `json:"reject"`
}
```

**前端任务**:
- [ ] 规则文件管理界面
- [ ] 规则编辑器（支持语法高亮）
- [ ] 规则模板选择
- [ ] 规则导入导出
- [ ] 规则预览和测试

**文件结构**:
```
data/rules/
├── builtin/     # 内置规则
├── custom/      # 自定义规则
├── remote/      # 远程规则缓存
└── compiled/    # 编译后的规则
```

#### 2.5 实时监控系统 (Week 7-8)
**主要目标**: 实现流量监控和连接状态展示

**后端任务**:
- [ ] WebSocket 实时数据推送
- [ ] 流量统计收集
- [ ] 连接信息监控
- [ ] 历史数据存储
- [ ] 性能指标计算

**前端任务**:
- [ ] 实时监控面板
- [ ] 流量图表展示（ECharts）
- [ ] 连接列表页面
- [ ] 性能指标仪表盘
- [ ] WebSocket 连接管理

### 🎨 阶段 3: 用户界面完善 (Week 9-11)

#### 3.1 界面设计和交互优化 (Week 9-10)
**主要目标**: 提升用户体验和界面美观度

**任务清单**:
- [ ] UI/UX 设计评审和优化
- [ ] 深色/浅色主题支持
- [ ] 响应式设计完善
- [ ] 加载状态和错误处理
- [ ] 操作反馈和提示优化
- [ ] 快捷键支持
- [ ] 可访问性改进

#### 3.2 系统设置和配置 (Week 10-11)
**主要目标**: 实现系统配置和参数管理

**后端任务**:
- [ ] 系统设置数据模型
- [ ] 配置文件管理
- [ ] 参数验证和类型转换
- [ ] 配置热重载
- [ ] 默认配置预设

**前端任务**:
- [ ] 系统设置页面
- [ ] 配置表单和验证
- [ ] 参数分类管理
- [ ] 配置导入导出
- [ ] 重置默认设置

### 🔧 阶段 4: 系统完善和优化 (Week 12-14)

#### 4.1 性能优化 (Week 12)
- [ ] 后端性能调优
- [ ] 数据库查询优化
- [ ] 前端打包优化
- [ ] 资源懒加载
- [ ] 缓存策略优化

#### 4.2 测试和质量保证 (Week 13)
- [ ] 单元测试编写
- [ ] 集成测试
- [ ] 端到端测试
- [ ] 性能测试
- [ ] 安全测试

#### 4.3 部署和发布准备 (Week 14)
- [ ] 生产环境配置
- [ ] Docker 镜像构建
- [ ] CI/CD 流程完善
- [ ] 文档完善
- [ ] 版本发布准备

## API 设计规范

### RESTful API 设计
```yaml
# OpenAPI 3.0 规范示例
openapi: 3.0.0
info:
  title: Prism API
  version: 1.0.0

paths:
  /api/nodes:
    get:
      summary: 获取节点列表
      parameters:
        - name: page
          in: query
          schema:
            type: integer
        - name: size
          in: query
          schema:
            type: integer
      responses:
        200:
          description: 成功返回节点列表
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/Node'
                  pagination:
                    $ref: '#/components/schemas/Pagination'
```

### WebSocket 事件设计
```typescript
// WebSocket 消息类型定义
interface WSMessage {
  type: 'traffic' | 'connections' | 'status' | 'logs';
  data: any;
  timestamp: number;
}

// 流量数据
interface TrafficData {
  up: number;
  down: number;
  total: {
    up: number;
    down: number;
  };
}

// 连接信息
interface ConnectionInfo {
  id: string;
  metadata: {
    network: string;
    type: string;
    sourceIP: string;
    destinationIP: string;
    host: string;
  };
  upload: number;
  download: number;
  start: string;
  chains: string[];
  rule: string;
}
```

## 数据库设计

**详细设计请参考**: [数据库设计文档 (最小化版)](./DATABASE_DESIGN_MINIMAL.md)

### 最小化架构概述

基于代理客户端性能优先的考虑，数据库采用**最小化设计**，只保留3个核心表，移除所有可能影响代理性能的日志、统计和监控表。

#### 🎯 核心表结构 (仅3个表)
- **subscriptions**: 订阅管理，自动抓取和更新节点
- **proxy_nodes**: 代理节点，支持多协议和基本性能信息
- **subscription_nodes**: 订阅节点关联表（多对多关系）

#### ❌ 移除的表 (性能考虑)
- ~~subscription_logs~~ → 改用文件日志
- ~~node_tests~~ → 改用内存/BoltDB缓存
- ~~traffic_stats~~ → 改用LevelDB缓存
- ~~connection_logs~~ → 改用文件日志
- ~~system_settings~~ → 改用YAML配置文件
- ~~operation_logs~~ → 改用文件日志

#### 📁 替代存储方案
- **文件日志**: app.log, subscription.log, node_test.log, proxy.log
- **内存缓存**: BoltDB (节点测试) + LevelDB (流量统计)
- **配置文件**: settings.yaml 系统配置管理
- **规则文件**: builtin、custom、remote、compiled 分类管理

#### ⚡ 性能优势
- **启动时间**: 减少 80%+ (3个表 vs 15个表)
- **内存占用**: 减少 60%+ (无大量日志数据)
- **数据库操作**: 减少 90%+ (最小化写入操作)
- **代理延迟**: 几乎无影响 (无实时统计写入)

### GORM 模型示例 (最小化版)
```go
// 使用推荐的包
import (
    "gorm.io/gorm"
    "github.com/lazygophers/utils/xtime"
    "github.com/google/uuid"
)

// BaseModel 基础模型
type BaseModel struct {
    ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
    CreatedAt xtime.Time     `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt xtime.Time     `gorm:"autoUpdateTime" json:"updated_at"`  
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// Subscription 订阅模型 (简化版)
type Subscription struct {
    BaseModel
    Name                string     `gorm:"size:255;not null" json:"name"`
    URL                 string     `gorm:"type:text;not null" json:"url"`
    Type                string     `gorm:"size:50;not null;default:clash" json:"type"`
    UpdateIntervalHours int        `gorm:"not null;default:24" json:"update_interval_hours"`
    AutoUpdate          bool       `gorm:"not null;default:true" json:"auto_update"`
    Enabled             bool       `gorm:"not null;default:true" json:"enabled"`
    LastUpdateAt        *xtime.Time `json:"last_update_at"`
    NextUpdateAt        *xtime.Time `json:"next_update_at"`
    Status              string     `gorm:"size:20;not null;default:pending" json:"status"`
    NodeCount           int        `gorm:"not null;default:0" json:"node_count"`
    
    // 多对多关系
    Nodes []ProxyNode `gorm:"many2many:subscription_nodes" json:"nodes,omitempty"`
}

// ProxyNode 代理节点模型 (简化版)
type ProxyNode struct {
    BaseModel
    Name        string     `gorm:"size:255;not null" json:"name"`
    Type        string     `gorm:"size:50;not null" json:"type"`
    Server      string     `gorm:"size:255;not null" json:"server"`
    Port        int        `gorm:"not null" json:"port"`
    Config      string     `gorm:"type:text;not null" json:"config"`
    CountryCode string     `gorm:"size:10" json:"country_code"`
    DelayMS     int        `gorm:"default:-1" json:"delay_ms"`
    LastTestAt  *xtime.Time `json:"last_test_at"`
    Status      string     `gorm:"size:20;not null;default:unknown" json:"status"`
    Enabled     bool       `gorm:"not null;default:true" json:"enabled"`
    SortOrder   int        `gorm:"not null;default:0" json:"sort_order"`
    
    // 多对多关系
    Subscriptions []Subscription `gorm:"many2many:subscription_nodes" json:"subscriptions,omitempty"`
}

// SubscriptionNode 关联表模型 (最小化)
type SubscriptionNode struct {
    BaseModel
    SubscriptionID string `gorm:"size:36;not null" json:"subscription_id"`
    NodeID         string `gorm:"size:36;not null" json:"node_id"`
    NodeIndex      int    `gorm:"not null" json:"node_index"`
    IsPrimary      bool   `gorm:"not null;default:false" json:"is_primary"`
}
```

### 多数据库配置
```go
// 数据库配置支持
type DatabaseConfig struct {
    Type     string // sqlite, mysql, postgres, gaussdb
    Host     string
    Port     int
    Database string
    Username string
    Password string
    SSLMode  string
}

// 自动适配数据库方言
func NewDatabase(config *DatabaseConfig) (*gorm.DB, error) {
    switch config.Type {
    case "sqlite":
        return gorm.Open(sqlite.Open(config.Database), &gorm.Config{})
    case "mysql":
        dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
            config.Username, config.Password, config.Host, config.Port, config.Database)
        return gorm.Open(mysql.Open(dsn), &gorm.Config{})
    case "postgres":
        dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai",
            config.Host, config.Username, config.Password, config.Database, config.Port, config.SSLMode)
        return gorm.Open(postgres.Open(dsn), &gorm.Config{})
    }
}
```

## 前端组件设计

### 核心组件结构
```typescript
// 组件层次结构
src/
├── components/
│   ├── Layout/              # 布局组件
│   │   ├── Header/
│   │   ├── Sidebar/
│   │   └── Footer/
│   ├── Node/                # 节点相关组件
│   │   ├── NodeList/
│   │   ├── NodeCard/
│   │   ├── NodeForm/
│   │   └── NodeTest/
│   ├── Subscription/        # 订阅相关组件
│   ├── Rule/                # 规则相关组件
│   ├── Monitor/             # 监控相关组件
│   └── Common/              # 通用组件
│       ├── LoadingSpinner/
│       ├── ErrorBoundary/
│       ├── ConfirmDialog/
│       └── Toast/
```

### 状态管理设计
```typescript
// Zustand store 结构
interface AppState {
  // 用户状态
  user: {
    isAuthenticated: boolean;
    userInfo: UserInfo | null;
    login: (credentials: LoginCredentials) => Promise<void>;
    logout: () => void;
  };
  
  // 节点状态
  nodes: {
    list: Node[];
    loading: boolean;
    selectedNodes: string[];
    fetchNodes: () => Promise<void>;
    addNode: (node: CreateNodeRequest) => Promise<void>;
    updateNode: (id: string, node: UpdateNodeRequest) => Promise<void>;
    deleteNode: (id: string) => Promise<void>;
    testNode: (id: string) => Promise<number>;
  };
  
  // 订阅状态
  subscriptions: {
    list: Subscription[];
    loading: boolean;
    fetchSubscriptions: () => Promise<void>;
    addSubscription: (sub: CreateSubscriptionRequest) => Promise<void>;
    updateSubscription: (id: string) => Promise<void>;
  };
  
  // 监控状态
  monitor: {
    traffic: TrafficData;
    connections: ConnectionInfo[];
    isConnected: boolean;
    connect: () => void;
    disconnect: () => void;
  };
}
```

## 测试策略

### 后端测试
```go
// 使用推荐的包进行测试
import (
    "testing"
    "github.com/lazygophers/log"
    "github.com/lazygophers/utils/json"
    "github.com/lazygophers/utils/stringx"
    "github.com/stretchr/testify/assert"
    "go.uber.org/atomic"
)

// 单元测试示例
func TestNodeService_CreateNode(t *testing.T) {
    service := NewNodeService(mockDB)
    node := &Node{
        Name:   "Test Node",
        Type:   "vmess", 
        Server: "example.com",
        Port:   443,
    }
    
    // 使用 lazygophers/log 记录测试日志
    log.Info("开始创建节点测试", log.String("name", node.Name))
    
    createdNode, err := service.CreateNode(node)
    assert.NoError(t, err)
    assert.Equal(t, node.Name, createdNode.Name)
    
    // 使用 lazygophers/utils/json 进行JSON操作
    nodeJSON, _ := json.Marshal(createdNode)
    log.Debug("创建的节点", log.String("json", string(nodeJSON)))
}

// Fiber API 测试示例
func TestNodesAPI(t *testing.T) {
    app := fiber.New()
    setupRoutes(app)
    
    req := httptest.NewRequest("GET", "/api/nodes", nil)
    resp, err := app.Test(req)
    
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
    
    // 使用 stringx 进行字符串操作
    contentType := resp.Header.Get("Content-Type")
    assert.True(t, stringx.Contains(contentType, "application/json"))
}

// 原子操作测试
func TestAtomicCounters(t *testing.T) {
    var counter atomic.Int64
    counter.Store(0)
    
    // 并发测试
    for i := 0; i < 1000; i++ {
        go func() {
            counter.Inc()
        }()
    }
    
    // 等待所有 goroutine 完成
    time.Sleep(100 * time.Millisecond)
    assert.Equal(t, int64(1000), counter.Load())
}
```

### 前端测试
```typescript
// 组件测试示例
import { render, screen, fireEvent } from '@testing-library/react';
import { NodeList } from './NodeList';

test('renders node list correctly', () => {
  const mockNodes = [
    { id: '1', name: 'Node 1', type: 'vmess', delay: 100 },
    { id: '2', name: 'Node 2', type: 'trojan', delay: 200 },
  ];
  
  render(<NodeList nodes={mockNodes} />);
  
  expect(screen.getByText('Node 1')).toBeInTheDocument();
  expect(screen.getByText('Node 2')).toBeInTheDocument();
});

// 集成测试
test('node CRUD operations', async () => {
  render(<App />);
  
  // 添加节点
  fireEvent.click(screen.getByText('Add Node'));
  // ... 填充表单
  fireEvent.click(screen.getByText('Save'));
  
  // 验证节点添加成功
  await waitFor(() => {
    expect(screen.getByText('New Node')).toBeInTheDocument();
  });
});
```

## 部署策略

### 开发环境部署
```bash
# 后端开发服务
cd core && go run cmd/server/main.go

# 前端开发服务
cd web && npm run dev

# 数据库启动
sqlite3 prism.db
```

### Docker 部署
```dockerfile
# 多阶段构建
FROM node:18-alpine AS frontend-builder
WORKDIR /app/web
COPY web/package*.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

FROM golang:1.21-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend-builder /app/web/dist ./static
RUN go build -o prism cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=backend-builder /app/prism .
COPY --from=backend-builder /app/static ./static
EXPOSE 9090
CMD ["./prism"]
```

### 生产环境部署
```yaml
# docker-compose.yml
version: '3.8'
services:
  prism:
    image: prism:latest
    ports:
      - "9090:9090"
      - "7890:7890"
      - "7891:7891"
    volumes:
      - ./config:/app/config
      - ./data:/app/data
    environment:
      - PRISM_ENV=production
    restart: unless-stopped
```

## 风险评估和应对

### 技术风险
1. **mihomo/clash 兼容性风险**
   - 风险: 上游项目变更导致兼容性问题
   - 应对: 版本锁定 + 定期更新测试

2. **性能瓶颈风险**
   - 风险: 大量节点时性能下降
   - 应对: 性能测试 + 优化策略预案

3. **安全漏洞风险**
   - 风险: 代理服务安全问题
   - 应对: 安全审计 + 及时更新

### 项目风险
1. **开发进度风险**
   - 风险: 功能复杂度超出预期
   - 应对: 分阶段交付 + 功能优先级管理

2. **用户需求变更风险**
   - 风险: 需求频繁变更影响开发
   - 应对: 需求确认流程 + 变更管理

## 交付标准

### 功能完整性
- ✅ 节点管理功能 100% 完成
- ✅ 订阅管理功能 100% 完成
- ✅ 规则配置功能 100% 完成
- ✅ 实时监控功能 100% 完成
- ✅ 用户认证功能 100% 完成

### 质量标准
- ✅ 代码测试覆盖率 > 80%
- ✅ API 响应时间 < 200ms
- ✅ Web 界面加载时间 < 3s
- ✅ 支持 1000+ 节点管理
- ✅ 24/7 稳定运行

### 文档完整性
- ✅ 用户使用手册
- ✅ API 文档
- ✅ 部署指南
- ✅ 开发者文档
- ✅ 故障排查指南

---

**项目联系人**: 开发团队
**文档更新**: 2024年开发计划
**下次评审**: 开发启动后每周评审