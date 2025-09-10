# GoFiber v2 开发指南

## 概述

本指南详细介绍如何在 Prism 项目中使用 GoFiber v2 框架构建高性能的 Web API。GoFiber 是一个基于 Fasthttp 的 Express 风格的 Go Web 框架，具有零内存分配和极高的性能。

## GoFiber 核心特性

### 🚀 性能优势
- **零内存分配路由器**: 比标准库快 10-15 倍
- **基于 Fasthttp**: 比 net/http 快 10 倍
- **低内存占用**: 最小的内存分配
- **快速启动**: 应用启动时间极短

### 🎯 Express 风格 API
- 熟悉的中间件模式
- 链式路由处理
- 灵活的路径匹配
- 内置的静态文件服务

## 项目结构

```
internal/
├── api/
│   ├── handlers/          # HTTP 处理器
│   │   ├── auth.go
│   │   ├── nodes.go
│   │   ├── subscriptions.go
│   │   └── websocket.go
│   ├── middleware/        # 中间件
│   │   ├── auth.go
│   │   ├── cors.go
│   │   ├── logger.go
│   │   └── rate_limit.go
│   ├── routes/            # 路由定义
│   │   ├── api.go
│   │   └── routes.go
│   └── dto/              # 数据传输对象
│       ├── request.go
│       └── response.go
├── service/              # 业务逻辑层
└── models/               # 数据模型
```

## 基础配置

### main.go 应用入口
```go
package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/compress"
    "github.com/gofiber/fiber/v2/middleware/recover"
    
    "prism/internal/api/routes"
    "prism/internal/database"
    logger "github.com/lazygophers/log"
)

func main() {
    // 创建 Fiber 应用
    app := fiber.New(fiber.Config{
        // 应用配置
        AppName:      "Prism API v1.0.0",
        ServerHeader: "Prism",
        
        // 性能配置
        Prefork:       false,  // 不使用 prefork 模式
        CaseSensitive: true,   // 路径大小写敏感
        StrictRouting: true,   // 严格路由模式
        
        // 请求配置
        BodyLimit:    4 * 1024 * 1024,  // 4MB 请求体限制
        ReadTimeout:  time.Second * 10,  // 读取超时
        WriteTimeout: time.Second * 10,  // 写入超时
        IdleTimeout:  time.Second * 120, // 空闲超时
        
        // 错误处理
        ErrorHandler: func(c *fiber.Ctx, err error) error {
            code := fiber.StatusInternalServerError
            message := "Internal Server Error"
            
            if e, ok := err.(*fiber.Error); ok {
                code = e.Code
                message = e.Message
            }
            
            logger.Error("HTTP错误", 
                logger.String("path", c.Path()),
                logger.String("method", c.Method()),
                logger.Int("status", code),
                logger.String("error", message),
            )
            
            return c.Status(code).JSON(fiber.Map{
                "code":    code,
                "message": message,
                "path":    c.Path(),
                "method":  c.Method(),
            })
        },
    })
    
    // 全局中间件
    app.Use(recover.New())
    app.Use(compress.New())
    app.Use(cors.New(cors.Config{
        AllowOrigins:     "http://localhost:3000,https://localhost:3000",
        AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
        AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
        AllowCredentials: true,
    }))
    
    // 连接数据库
    db, err := database.Connect()
    if err != nil {
        log.Fatal("数据库连接失败:", err)
    }
    
    // 设置路由
    routes.SetupRoutes(app, db)
    
    // 优雅关闭
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    
    go func() {
        <-c
        logger.Info("正在关闭服务器...")
        app.Shutdown()
    }()
    
    // 启动服务器
    port := os.Getenv("PORT")
    if port == "" {
        port = "9090"
    }
    
    logger.Info("服务器启动", logger.String("port", port))
    log.Fatal(app.Listen(":" + port))
}
```

### 路由配置
```go
// internal/api/routes/routes.go
package routes

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/websocket/v2"
    "gorm.io/gorm"
    
    "prism/internal/api/handlers"
    "prism/internal/api/middleware"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
    // 初始化处理器
    nodeHandler := handlers.NewNodeHandler(db)
    authHandler := handlers.NewAuthHandler(db)
    subHandler := handlers.NewSubscriptionHandler(db)
    wsHandler := handlers.NewWebSocketHandler(db)
    
    // 健康检查
    app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status":    "ok",
            "service":   "prism-api",
            "version":   "1.0.0",
            "timestamp": time.Now().Unix(),
        })
    })
    
    // 静态文件服务 (前端构建文件)
    app.Static("/", "./web/dist", fiber.Static{
        Compress:      true,
        ByteRange:     true,
        Browse:        false,
        CacheDuration: 24 * time.Hour,
    })
    
    // API 路由组
    api := app.Group("/api/v1")
    
    // 公开路由 (不需要认证)
    public := api.Group("")
    public.Post("/auth/login", authHandler.Login)
    public.Post("/auth/register", authHandler.Register)
    public.Get("/system/info", handlers.GetSystemInfo)
    
    // 受保护的路由 (需要JWT认证)
    protected := api.Group("", middleware.JWTAuth())
    
    // 用户管理
    users := protected.Group("/users")
    users.Get("/profile", authHandler.GetProfile)
    users.Put("/profile", authHandler.UpdateProfile)
    users.Post("/logout", authHandler.Logout)
    
    // 节点管理
    nodes := protected.Group("/nodes")
    nodes.Get("/", nodeHandler.GetList)
    nodes.Post("/", nodeHandler.Create)
    nodes.Get("/:id", nodeHandler.GetByID)
    nodes.Put("/:id", nodeHandler.Update)
    nodes.Delete("/:id", nodeHandler.Delete)
    nodes.Post("/:id/test", nodeHandler.Test)
    nodes.Post("/batch", nodeHandler.BatchOperation)
    
    // 订阅管理
    subs := protected.Group("/subscriptions")
    subs.Get("/", subHandler.GetList)
    subs.Post("/", subHandler.Create)
    subs.Get("/:id", subHandler.GetByID)
    subs.Put("/:id", subHandler.Update)
    subs.Delete("/:id", subHandler.Delete)
    subs.Post("/:id/update", subHandler.UpdateNodes)
    
    // WebSocket 路由 (实时数据)
    app.Get("/ws", websocket.New(wsHandler.HandleConnection))
    
    // 管理员路由
    admin := protected.Group("/admin", middleware.AdminAuth())
    admin.Get("/users", authHandler.GetUsers)
    admin.Get("/system/stats", handlers.GetSystemStats)
    admin.Get("/logs", handlers.GetLogs)
    
    // 404 处理
    app.Use(func(c *fiber.Ctx) error {
        return c.Status(404).JSON(fiber.Map{
            "code":    404,
            "message": "路由未找到",
            "path":    c.Path(),
        })
    })
}
```

## 中间件开发

### JWT 认证中间件
```go
// internal/api/middleware/auth.go
package middleware

import (
    "strings"
    
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v4"
    "github.com/lazygophers/log"
    "github.com/lazygophers/utils/stringx"
)

type JWTClaims struct {
    UserID string `json:"user_id"`
    Role   string `json:"role"`
    jwt.StandardClaims
}

func JWTAuth() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // 从 Header 获取 Token
        authHeader := c.Get("Authorization")
        if stringx.IsEmpty(authHeader) {
            return c.Status(401).JSON(fiber.Map{
                "code":    401,
                "message": "缺少认证令牌",
            })
        }
        
        // 解析 Bearer Token
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            return c.Status(401).JSON(fiber.Map{
                "code":    401,
                "message": "无效的认证令牌格式",
            })
        }
        
        tokenString := parts[1]
        
        // 验证 JWT Token
        token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
            return []byte("your-secret-key"), nil
        })
        
        if err != nil || !token.Valid {
            log.Warn("JWT认证失败", 
                log.String("error", err.Error()),
                log.String("ip", c.IP()),
                log.String("userAgent", c.Get("User-Agent")),
            )
            return c.Status(401).JSON(fiber.Map{
                "code":    401,
                "message": "认证令牌无效或已过期",
            })
        }
        
        // 获取用户信息
        claims, ok := token.Claims.(*JWTClaims)
        if !ok {
            return c.Status(401).JSON(fiber.Map{
                "code":    401,
                "message": "无效的令牌声明",
            })
        }
        
        // 将用户信息存储到上下文
        c.Locals("user_id", claims.UserID)
        c.Locals("user_role", claims.Role)
        
        return c.Next()
    }
}

func AdminAuth() fiber.Handler {
    return func(c *fiber.Ctx) error {
        role := c.Locals("user_role")
        if role != "admin" {
            return c.Status(403).JSON(fiber.Map{
                "code":    403,
                "message": "需要管理员权限",
            })
        }
        return c.Next()
    }
}
```

### 请求日志中间件
```go
// internal/api/middleware/logger.go
package middleware

import (
    "time"
    
    "github.com/gofiber/fiber/v2"
    "github.com/lazygophers/log"
    "github.com/lazygophers/utils/xtime"
)

func Logger() fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := xtime.Now()
        
        // 执行下一个中间件
        err := c.Next()
        
        // 记录请求信息
        duration := xtime.Since(start)
        status := c.Response().StatusCode()
        
        // 构建日志字段
        fields := []log.Field{
            log.String("method", c.Method()),
            log.String("path", c.Path()),
            log.String("query", c.Request().URI().QueryString()),
            log.Int("status", status),
            log.Duration("duration", duration),
            log.String("ip", c.IP()),
            log.String("userAgent", c.Get("User-Agent")),
            log.Int("bodySize", len(c.Response().Body())),
        }
        
        // 添加用户信息 (如果已认证)
        if userID := c.Locals("user_id"); userID != nil {
            fields = append(fields, log.String("userId", userID.(string)))
        }
        
        // 根据状态码选择日志级别
        switch {
        case status >= 500:
            log.Error("HTTP请求", fields...)
        case status >= 400:
            log.Warn("HTTP请求", fields...)
        default:
            log.Info("HTTP请求", fields...)
        }
        
        return err
    }
}
```

### 限流中间件
```go
// internal/api/middleware/rate_limit.go
package middleware

import (
    "time"
    
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/limiter"
    "github.com/gofiber/storage/memory"
)

func RateLimit() fiber.Handler {
    return limiter.New(limiter.Config{
        Max:        100,              // 最大请求数
        Expiration: 1 * time.Minute, // 时间窗口
        KeyGenerator: func(c *fiber.Ctx) string {
            return c.IP() // 基于 IP 限流
        },
        LimitReached: func(c *fiber.Ctx) error {
            return c.Status(429).JSON(fiber.Map{
                "code":    429,
                "message": "请求过于频繁，请稍后再试",
            })
        },
        Storage: memory.New(), // 使用内存存储
    })
}

func APIRateLimit() fiber.Handler {
    return limiter.New(limiter.Config{
        Max:        1000,             // API 更高的限制
        Expiration: 1 * time.Hour,   // 更长的时间窗口
        KeyGenerator: func(c *fiber.Ctx) string {
            if userID := c.Locals("user_id"); userID != nil {
                return userID.(string) // 基于用户限流
            }
            return c.IP() // 未认证用户基于 IP
        },
        LimitReached: func(c *fiber.Ctx) error {
            return c.Status(429).JSON(fiber.Map{
                "code":    429,
                "message": "API 调用频率超出限制",
            })
        },
    })
}
```

## 处理器开发

### 节点管理处理器
```go
// internal/api/handlers/nodes.go
package handlers

import (
    "strconv"
    
    "github.com/gofiber/fiber/v2"
    "github.com/lazygophers/log"
    "github.com/lazygophers/utils/stringx"
    "github.com/lazygophers/utils/xtime"
    "gorm.io/gorm"
    
    "prism/internal/models"
    "prism/internal/service"
)

type NodeHandler struct {
    nodeService *service.NodeService
}

func NewNodeHandler(db *gorm.DB) *NodeHandler {
    return &NodeHandler{
        nodeService: service.NewNodeService(db),
    }
}

// GetList 获取节点列表
func (h *NodeHandler) GetList(c *fiber.Ctx) error {
    // 分页参数
    page, _ := strconv.Atoi(c.Query("page", "1"))
    size, _ := strconv.Atoi(c.Query("size", "20"))
    
    if page < 1 {
        page = 1
    }
    if size < 1 || size > 100 {
        size = 20
    }
    
    // 过滤参数
    filter := &models.NodeFilter{
        Type:        c.Query("type"),
        CountryCode: c.Query("country"),
        Status:      c.Query("status"),
        Search:      c.Query("search"),
    }
    
    // 启用状态过滤
    if enabledStr := c.Query("enabled"); enabledStr != "" {
        enabled := enabledStr == "true"
        filter.Enabled = &enabled
    }
    
    // 排序参数
    sort := c.Query("sort", "created_at")
    order := c.Query("order", "desc")
    
    // 调用服务层
    result, err := h.nodeService.GetNodes(c.Context(), page, size, filter, sort, order)
    if err != nil {
        log.Error("获取节点列表失败", log.Error(err))
        return c.Status(500).JSON(fiber.Map{
            "code":    500,
            "message": "获取节点列表失败",
            "error":   err.Error(),
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

// Create 创建节点
func (h *NodeHandler) Create(c *fiber.Ctx) error {
    var req models.CreateNodeRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "code":    400,
            "message": "请求参数格式错误",
            "error":   err.Error(),
        })
    }
    
    // 验证请求参数
    if err := req.Validate(); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "code":    400,
            "message": "请求参数验证失败",
            "error":   err.Error(),
        })
    }
    
    // 获取当前用户
    userID := c.Locals("user_id").(string)
    
    // 调用服务层创建节点
    node, err := h.nodeService.CreateNode(c.Context(), &req, userID)
    if err != nil {
        log.Error("创建节点失败", 
            log.Error(err),
            log.String("userId", userID),
            log.String("nodeName", req.Name),
        )
        return c.Status(500).JSON(fiber.Map{
            "code":    500,
            "message": "创建节点失败",
            "error":   err.Error(),
        })
    }
    
    log.Info("节点创建成功", 
        log.String("nodeId", node.ID),
        log.String("nodeName", node.Name),
        log.String("userId", userID),
    )
    
    return c.Status(201).JSON(fiber.Map{
        "code":    201,
        "message": "节点创建成功",
        "data":    node,
    })
}

// Test 测试节点
func (h *NodeHandler) Test(c *fiber.Ctx) error {
    nodeID := c.Params("id")
    if stringx.IsEmpty(nodeID) {
        return c.Status(400).JSON(fiber.Map{
            "code":    400,
            "message": "节点ID不能为空",
        })
    }
    
    // 获取测试配置
    timeout := c.QueryInt("timeout", 5000) // 默认5秒超时
    
    // 异步测试节点
    go func() {
        start := xtime.Now()
        result, err := h.nodeService.TestNode(nodeID, timeout)
        duration := xtime.Since(start)
        
        if err != nil {
            log.Error("节点测试失败",
                log.String("nodeId", nodeID),
                log.Error(err),
                log.Duration("duration", duration),
            )
            return
        }
        
        log.Info("节点测试完成",
            log.String("nodeId", nodeID),
            log.Int64("delay", result.Delay),
            log.String("status", result.Status),
            log.Duration("testDuration", duration),
        )
    }()
    
    return c.JSON(fiber.Map{
        "code":    200,
        "message": "节点测试已启动",
        "data": fiber.Map{
            "nodeId": nodeID,
            "status": "testing",
        },
    })
}

// BatchOperation 批量操作
func (h *NodeHandler) BatchOperation(c *fiber.Ctx) error {
    var req models.BatchOperationRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "code":    400,
            "message": "请求参数格式错误",
        })
    }
    
    if len(req.NodeIDs) == 0 {
        return c.Status(400).JSON(fiber.Map{
            "code":    400,
            "message": "节点ID列表不能为空",
        })
    }
    
    userID := c.Locals("user_id").(string)
    
    result, err := h.nodeService.BatchOperation(c.Context(), &req, userID)
    if err != nil {
        log.Error("批量操作失败",
            log.Error(err),
            log.String("operation", req.Operation),
            log.Int("nodeCount", len(req.NodeIDs)),
            log.String("userId", userID),
        )
        return c.Status(500).JSON(fiber.Map{
            "code":    500,
            "message": "批量操作失败",
            "error":   err.Error(),
        })
    }
    
    return c.JSON(fiber.Map{
        "code":    200,
        "message": "批量操作完成",
        "data":    result,
    })
}
```

## WebSocket 处理

### WebSocket 处理器
```go
// internal/api/handlers/websocket.go
package handlers

import (
    "context"
    "encoding/json"
    "time"
    
    "github.com/gofiber/websocket/v2"
    "github.com/lazygophers/log"
    "go.uber.org/atomic"
    "gorm.io/gorm"
    
    "prism/internal/service"
)

type WebSocketHandler struct {
    nodeService *service.NodeService
    clients     map[string]*websocket.Conn
    clientCount atomic.Int64
}

type WSMessage struct {
    Type      string      `json:"type"`
    Data      interface{} `json:"data"`
    Timestamp int64       `json:"timestamp"`
}

type WSResponse struct {
    Type      string      `json:"type"`
    Data      interface{} `json:"data"`
    Success   bool        `json:"success"`
    Message   string      `json:"message"`
    Timestamp int64       `json:"timestamp"`
}

func NewWebSocketHandler(db *gorm.DB) *WebSocketHandler {
    return &WebSocketHandler{
        nodeService: service.NewNodeService(db),
        clients:     make(map[string]*websocket.Conn),
    }
}

func (h *WebSocketHandler) HandleConnection(c *websocket.Conn) {
    defer func() {
        c.Close()
        h.clientCount.Dec()
    }()
    
    clientID := c.Locals("user_id").(string)
    h.clients[clientID] = c
    count := h.clientCount.Inc()
    
    log.Info("WebSocket连接建立", 
        log.String("clientId", clientID),
        log.Int64("totalClients", count),
    )
    
    // 发送欢迎消息
    h.sendMessage(c, &WSResponse{
        Type:      "welcome",
        Success:   true,
        Message:   "WebSocket 连接建立成功",
        Timestamp: time.Now().Unix(),
    })
    
    // 启动数据推送 goroutine
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    go h.startDataPush(ctx, c, clientID)
    
    // 处理客户端消息
    for {
        var msg WSMessage
        if err := c.ReadJSON(&msg); err != nil {
            log.Warn("WebSocket消息读取失败", 
                log.String("clientId", clientID),
                log.Error(err),
            )
            break
        }
        
        h.handleMessage(c, clientID, &msg)
    }
    
    delete(h.clients, clientID)
    log.Info("WebSocket连接关闭", 
        log.String("clientId", clientID),
        log.Int64("totalClients", h.clientCount.Load()),
    )
}

func (h *WebSocketHandler) handleMessage(c *websocket.Conn, clientID string, msg *WSMessage) {
    log.Debug("收到WebSocket消息",
        log.String("clientId", clientID),
        log.String("type", msg.Type),
    )
    
    switch msg.Type {
    case "ping":
        h.sendMessage(c, &WSResponse{
            Type:      "pong",
            Success:   true,
            Data:      msg.Data,
            Timestamp: time.Now().Unix(),
        })
        
    case "subscribe_stats":
        // 订阅实时统计
        h.sendMessage(c, &WSResponse{
            Type:      "stats_subscribed",
            Success:   true,
            Message:   "已订阅实时统计数据",
            Timestamp: time.Now().Unix(),
        })
        
    case "test_node":
        // 实时节点测试
        if nodeID, ok := msg.Data.(string); ok {
            go h.handleNodeTest(c, clientID, nodeID)
        }
        
    default:
        h.sendMessage(c, &WSResponse{
            Type:      "error",
            Success:   false,
            Message:   "未知的消息类型",
            Timestamp: time.Now().Unix(),
        })
    }
}

func (h *WebSocketHandler) handleNodeTest(c *websocket.Conn, clientID, nodeID string) {
    // 发送开始测试消息
    h.sendMessage(c, &WSResponse{
        Type:    "node_test_start",
        Success: true,
        Data: map[string]interface{}{
            "nodeId": nodeID,
            "status": "testing",
        },
        Timestamp: time.Now().Unix(),
    })
    
    // 执行测试
    result, err := h.nodeService.TestNode(nodeID, 5000)
    
    if err != nil {
        h.sendMessage(c, &WSResponse{
            Type:    "node_test_error",
            Success: false,
            Data: map[string]interface{}{
                "nodeId": nodeID,
                "error":  err.Error(),
            },
            Message:   "节点测试失败",
            Timestamp: time.Now().Unix(),
        })
        return
    }
    
    // 发送测试结果
    h.sendMessage(c, &WSResponse{
        Type:    "node_test_result",
        Success: true,
        Data: map[string]interface{}{
            "nodeId": nodeID,
            "delay":  result.Delay,
            "status": result.Status,
        },
        Message:   "节点测试完成",
        Timestamp: time.Now().Unix(),
    })
}

func (h *WebSocketHandler) startDataPush(ctx context.Context, c *websocket.Conn, clientID string) {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            // 推送实时统计数据
            stats, err := h.nodeService.GetRealtimeStats()
            if err != nil {
                log.Error("获取实时统计失败", 
                    log.String("clientId", clientID),
                    log.Error(err),
                )
                continue
            }
            
            h.sendMessage(c, &WSResponse{
                Type:      "realtime_stats",
                Success:   true,
                Data:      stats,
                Timestamp: time.Now().Unix(),
            })
        }
    }
}

func (h *WebSocketHandler) sendMessage(c *websocket.Conn, response *WSResponse) {
    if err := c.WriteJSON(response); err != nil {
        log.Error("WebSocket消息发送失败", log.Error(err))
    }
}

// Broadcast 广播消息给所有客户端
func (h *WebSocketHandler) Broadcast(msgType string, data interface{}) {
    message := &WSResponse{
        Type:      msgType,
        Success:   true,
        Data:      data,
        Timestamp: time.Now().Unix(),
    }
    
    for clientID, conn := range h.clients {
        if err := conn.WriteJSON(message); err != nil {
            log.Warn("广播消息失败",
                log.String("clientId", clientID),
                log.Error(err),
            )
            // 移除失效连接
            delete(h.clients, clientID)
        }
    }
}
```

## 数据验证和序列化

### 请求/响应模型
```go
// internal/api/dto/request.go
package dto

import (
    "errors"
    "github.com/lazygophers/utils/stringx"
)

type CreateNodeRequest struct {
    Name     string            `json:"name" validate:"required,min=1,max=255"`
    Type     string            `json:"type" validate:"required,oneof=vmess vless trojan ss ssr"`
    Server   string            `json:"server" validate:"required,hostname"`
    Port     int               `json:"port" validate:"required,min=1,max=65535"`
    Config   map[string]interface{} `json:"config" validate:"required"`
    Tags     []string          `json:"tags"`
    Enabled  *bool             `json:"enabled"`
}

func (r *CreateNodeRequest) Validate() error {
    if stringx.IsEmpty(r.Name) {
        return errors.New("节点名称不能为空")
    }
    
    if stringx.IsEmpty(r.Server) {
        return errors.New("服务器地址不能为空")
    }
    
    if r.Port < 1 || r.Port > 65535 {
        return errors.New("端口范围必须在 1-65535 之间")
    }
    
    validTypes := []string{"vmess", "vless", "trojan", "ss", "ssr"}
    if !stringx.Contains(validTypes, r.Type) {
        return errors.New("不支持的节点类型")
    }
    
    return nil
}

type UpdateNodeRequest struct {
    Name    *string           `json:"name,omitempty"`
    Config  map[string]interface{} `json:"config,omitempty"`
    Tags    []string          `json:"tags,omitempty"`
    Enabled *bool             `json:"enabled,omitempty"`
}

type BatchOperationRequest struct {
    NodeIDs   []string `json:"node_ids" validate:"required,min=1"`
    Operation string   `json:"operation" validate:"required,oneof=delete enable disable test"`
}

// internal/api/dto/response.go
package dto

type NodeResponse struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Type        string            `json:"type"`
    Server      string            `json:"server"`
    Port        int               `json:"port"`
    CountryCode string            `json:"country_code"`
    CountryName string            `json:"country_name"`
    DelayMS     int               `json:"delay_ms"`
    Status      string            `json:"status"`
    Tags        []string          `json:"tags"`
    Enabled     bool              `json:"enabled"`
    CreatedAt   int64             `json:"created_at"`
    UpdatedAt   int64             `json:"updated_at"`
}

type PaginatedResponse struct {
    Data       interface{} `json:"data"`
    Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
    Page       int   `json:"page"`
    Size       int   `json:"size"`
    Total      int64 `json:"total"`
    TotalPages int64 `json:"total_pages"`
}

type APIResponse struct {
    Code      int         `json:"code"`
    Message   string      `json:"message"`
    Data      interface{} `json:"data,omitempty"`
    Error     string      `json:"error,omitempty"`
    Timestamp int64       `json:"timestamp"`
}
```

## 性能优化

### 连接池配置
```go
// 数据库连接池优化
func setupDatabase() *gorm.DB {
    db, err := gorm.Open(sqlite.Open("prism.db"), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    
    if err != nil {
        panic("数据库连接失败")
    }
    
    sqlDB, _ := db.DB()
    
    // 连接池配置
    sqlDB.SetMaxIdleConns(10)           // 最大空闲连接
    sqlDB.SetMaxOpenConns(100)          // 最大打开连接
    sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生命周期
    
    return db
}
```

### 缓存配置
```go
// 使用 GoFiber 内置缓存
import "github.com/gofiber/fiber/v2/middleware/cache"

app.Use(cache.New(cache.Config{
    Duration: 1 * time.Minute,
    KeyGenerator: func(c *fiber.Ctx) string {
        return c.Path() + ":" + c.Method()
    },
    CacheControl: true,
}))
```

## 测试

### 单元测试示例
```go
// internal/api/handlers/nodes_test.go
package handlers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "testing"
    
    "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestNodeHandler_Create(t *testing.T) {
    app := fiber.New()
    handler := NewNodeHandler(setupTestDB())
    
    app.Post("/nodes", handler.Create)
    
    payload := map[string]interface{}{
        "name":   "Test Node",
        "type":   "vmess",
        "server": "example.com",
        "port":   443,
        "config": map[string]interface{}{
            "uuid": "test-uuid",
        },
    }
    
    body, _ := json.Marshal(payload)
    req, _ := http.NewRequest("POST", "/nodes", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := app.Test(req)
    require.NoError(t, err)
    
    assert.Equal(t, 201, resp.StatusCode)
    
    // 验证响应
    var response map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&response)
    
    assert.Equal(t, float64(201), response["code"])
    assert.Equal(t, "节点创建成功", response["message"])
}
```

## 部署配置

### Docker 配置
```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/web/dist ./web/dist

EXPOSE 9090
CMD ["./main"]
```

### Docker Compose
```yaml
# docker-compose.yml
version: '3.8'
services:
  prism-api:
    build: .
    ports:
      - "9090:9090"
    environment:
      - ENV=production
      - DB_TYPE=sqlite
      - DB_PATH=/data/prism.db
      - JWT_SECRET=your-secret-key
    volumes:
      - ./data:/data
    restart: unless-stopped
```

这个指南提供了在 Prism 项目中使用 GoFiber v2 的完整开发方案，包括项目结构、中间件、处理器、WebSocket、测试和部署等各个方面。