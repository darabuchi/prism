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
Framework: Fiber (高性能、Express-like API)
Database: SQLite(默认) / MySQL / PostgreSQL / GaussDB
Cache: BoltDB / LevelDB (嵌入式键值存储)
ORM: GORM (支持多数据库)
Auth: JWT Token
WebSocket: Fiber WebSocket
Config: Viper (配置管理)
Logging: github.com/lazygophers/log
Utils: github.com/lazygophers/utils (json, stringx, xtime, bufiox, randx, anyx, candy)
Atomic: go.uber.org/atomic
Testing: Testify
```

### 前端技术栈
```json
{
  "framework": "React 18 + TypeScript",
  "stateManagement": "Zustand",
  "uiLibrary": "Ant Design",
  "buildTool": "Vite",
  "styling": "Tailwind CSS",
  "charts": "ECharts",
  "http": "Axios",
  "websocket": "native WebSocket API",
  "testing": "Jest + React Testing Library"
}
```

## 开发阶段规划

### 🔨 阶段 1: 项目基础搭建 (Week 1-2)

#### 后端基础架构
**时间**: 5-7 天
**负责人**: 后端开发

**任务清单**:
- [ ] Go 项目结构搭建
- [ ] Fiber 框架集成和中间件配置
- [ ] 多数据库支持设计 (SQLite/MySQL/PostgreSQL/GaussDB) 和 GORM 集成
- [ ] JWT 认证中间件实现
- [ ] 基础 API 路由定义
- [ ] lazygophers/log 日志系统集成
- [ ] lazygophers/utils 工具包集成
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
- [ ] React + TypeScript 项目初始化
- [ ] Vite 构建配置优化
- [ ] Ant Design 主题定制
- [ ] 路由系统设计 (React Router)
- [ ] Zustand 状态管理配置
- [ ] Axios HTTP 客户端封装
- [ ] 基础组件库建立
- [ ] 响应式布局框架

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

// 使用 Fiber 框架的API处理器示例
func setupNodeAPI(app *fiber.App) {
    api := app.Group("/api")
    
    // 使用原子计数器
    var requestCounter atomic.Int64
    
    // 中间件 - 请求计数和日志
    api.Use(func(c *fiber.Ctx) error {
        count := requestCounter.Inc()
        
        // 使用 lazygophers/log 记录请求
        log.Info("API请求", 
            log.String("method", c.Method()),
            log.String("path", c.Path()),
            log.Int64("count", count),
            log.String("ip", c.IP()),
        )
        
        return c.Next()
    })
    
    // 节点管理端点
    api.Get("/nodes", func(c *fiber.Ctx) error {
        // 使用 stringx 进行参数处理
        pageStr := c.Query("page", "1")
        if !stringx.IsNumeric(pageStr) {
            return c.Status(400).JSON(fiber.Map{
                "error": "页码参数无效",
            })
        }
        
        nodes, err := nodeService.GetNodes(c.Context())
        if err != nil {
            log.Error("获取节点失败", log.Error(err))
            return c.Status(500).JSON(fiber.Map{
                "error": "内部服务器错误",
            })
        }
        
        // 使用 lazygophers/utils/json 进行JSON操作
        response := fiber.Map{
            "data": nodes,
            "timestamp": xtime.Now().Unix(),
            "count": len(nodes),
        }
        
        return c.JSON(response)
    })
    
    // 节点测速端点
    api.Post("/nodes/:id/test", func(c *fiber.Ctx) error {
        nodeID := c.Params("id")
        if stringx.IsEmpty(nodeID) {
            return c.Status(400).JSON(fiber.Map{
                "error": "节点ID不能为空",
            })
        }
        
        // 异步测速
        go func() {
            start := xtime.Now()
            delay, err := testNodeDelay(nodeID)
            duration := xtime.Since(start)
            
            if err != nil {
                log.Error("节点测速失败",
                    log.String("nodeId", nodeID),
                    log.Error(err),
                    log.Duration("duration", duration),
                )
                return
            }
            
            log.Info("节点测速完成",
                log.String("nodeId", nodeID), 
                log.Int64("delay", delay),
                log.Duration("testDuration", duration),
            )
        }()
        
        return c.JSON(fiber.Map{
            "message": "测速任务已启动",
            "nodeId": nodeID,
        })
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

#### 2.2 节点管理系统 (Week 4-5)
**主要目标**: 实现节点的增删改查和批量管理

**后端任务**:
- [ ] 节点数据模型设计
- [ ] 节点 CRUD API 实现
- [ ] 节点测速功能
- [ ] 节点健康检查
- [ ] 订阅链接解析

```go
type Node struct {
    ID       string    `json:"id"`
    Name     string    `json:"name"`
    Type     string    `json:"type"`     // vmess, vless, trojan, ss
    Server   string    `json:"server"`
    Port     int       `json:"port"`
    Config   NodeConfig `json:"config"`
    Delay    int64     `json:"delay"`
    Status   string    `json:"status"`   // active, error, testing
    UpdateAt time.Time `json:"update_at"`
}
```

**前端任务**:
- [ ] 节点列表页面开发
- [ ] 节点添加/编辑表单
- [ ] 节点测速和状态显示
- [ ] 批量操作功能
- [ ] 搜索和过滤功能

**API 端点**:
```
GET    /api/nodes              # 获取节点列表
POST   /api/nodes              # 添加节点
GET    /api/nodes/{id}         # 获取单个节点
PUT    /api/nodes/{id}         # 更新节点
DELETE /api/nodes/{id}         # 删除节点
POST   /api/nodes/{id}/test    # 测试节点延迟
POST   /api/nodes/batch        # 批量操作
```

#### 2.3 订阅管理系统 (Week 5-6)
**主要目标**: 支持订阅链接管理和自动更新

**后端任务**:
- [ ] 订阅数据模型设计
- [ ] 订阅链接解析引擎
- [ ] 自动更新调度器
- [ ] 订阅内容缓存
- [ ] 错误处理和重试机制

```go
type Subscription struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    URL         string    `json:"url"`
    UpdatedAt   time.Time `json:"updated_at"`
    NodeCount   int       `json:"node_count"`
    Status      string    `json:"status"`
    AutoUpdate  bool      `json:"auto_update"`
    UpdateInterval int    `json:"update_interval"` // 小时
}
```

**前端任务**:
- [ ] 订阅管理页面
- [ ] 订阅添加和配置
- [ ] 订阅更新状态监控
- [ ] 订阅节点预览
- [ ] 自动更新设置

#### 2.4 规则配置系统 (Week 6-7)
**主要目标**: 实现灵活的路由规则配置

**后端任务**:
- [ ] 规则数据模型和存储
- [ ] 规则引擎集成
- [ ] 预设规则模板
- [ ] 规则验证和测试
- [ ] 规则优先级管理

```go
type Rule struct {
    ID       string `json:"id"`
    Type     string `json:"type"`      // DOMAIN, DOMAIN-SUFFIX, IP-CIDR
    Payload  string `json:"payload"`   // 规则内容
    Proxy    string `json:"proxy"`     // 代理策略
    Priority int    `json:"priority"`
}
```

**前端任务**:
- [ ] 规则管理页面
- [ ] 规则编辑器（支持语法高亮）
- [ ] 规则模板选择
- [ ] 规则测试工具
- [ ] 拖拽排序功能

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

#### 3.2 用户认证和权限 (Week 10-11)
**主要目标**: 实现安全的用户认证系统

**后端任务**:
- [ ] 用户认证中间件
- [ ] 权限控制系统
- [ ] 密码安全处理
- [ ] 会话管理
- [ ] API 访问控制

**前端任务**:
- [ ] 登录/注册页面
- [ ] 用户设置页面
- [ ] 权限状态管理
- [ ] 自动登录和记住密码
- [ ] 退出登录处理

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

### 表结构设计
```sql
-- 节点表
CREATE TABLE nodes (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    server VARCHAR(255) NOT NULL,
    port INTEGER NOT NULL,
    config TEXT NOT NULL,
    delay INTEGER DEFAULT -1,
    status VARCHAR(20) DEFAULT 'inactive',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 订阅表
CREATE TABLE subscriptions (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url TEXT NOT NULL,
    updated_at DATETIME,
    node_count INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'pending',
    auto_update BOOLEAN DEFAULT true,
    update_interval INTEGER DEFAULT 24,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 规则表
CREATE TABLE rules (
    id VARCHAR(36) PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    payload TEXT NOT NULL,
    proxy VARCHAR(255) NOT NULL,
    priority INTEGER DEFAULT 0,
    enabled BOOLEAN DEFAULT true,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 用户表
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'user',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_login DATETIME
);
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