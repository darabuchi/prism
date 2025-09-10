package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prism/core/internal/config"
	"github.com/prism/core/internal/core"
)

// Server API 服务器
type Server struct {
	config    *config.Config
	proxyCore *core.ProxyCore
	engine    *gin.Engine
}

// NewServer 创建 API 服务器
func NewServer(cfg *config.Config, proxyCore *core.ProxyCore) *Server {
	// 设置 Gin 模式
	if cfg.Log.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	// 中间件
	engine.Use(gin.Recovery())
	engine.Use(LoggerMiddleware())
	engine.Use(CORSMiddleware(cfg.API.AllowOrigins))

	server := &Server{
		config:    cfg,
		proxyCore: proxyCore,
		engine:    engine,
	}

	server.setupRoutes()
	return server
}

// Engine 返回 Gin 引擎
func (s *Server) Engine() *gin.Engine {
	return s.engine
}

// setupRoutes 设置路由
func (s *Server) setupRoutes() {
	// API v1 路由组
	v1 := s.engine.Group("/api/v1")
	{
		// 系统相关
		system := v1.Group("/system")
		{
			system.GET("/status", s.getSystemStatus)
			system.GET("/config", s.getSystemConfig)
			system.PUT("/config", s.updateSystemConfig)
		}

		// 代理相关
		proxy := v1.Group("/proxy")
		{
			proxy.GET("/status", s.getProxyStatus)
			proxy.PUT("/mode", s.setProxyMode)
			proxy.GET("/traffic", s.getTrafficStats)
		}

		// 节点池相关
		nodepools := v1.Group("/nodepools")
		{
			nodepools.GET("", s.getNodePools)
			nodepools.POST("", s.createNodePool)
			nodepools.GET("/:id", s.getNodePool)
			nodepools.PUT("/:id", s.updateNodePool)
			nodepools.DELETE("/:id", s.deleteNodePool)
		}

		// 健康检查
		v1.GET("/health", s.healthCheck)
	}

	// WebSocket 路由
	s.engine.GET("/api/v1/ws", s.handleWebSocket)
}

// Response 标准响应格式
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func (s *Server) Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Error 错误响应
func (s *Server) Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// getSystemStatus 获取系统状态
func (s *Server) getSystemStatus(c *gin.Context) {
	status := map[string]interface{}{
		"version":        "1.0.0", // TODO: 从构建信息获取
		"uptime":         0,       // TODO: 计算运行时间
		"memory_usage":   0,       // TODO: 获取内存使用
		"cpu_usage":      0,       // TODO: 获取 CPU 使用
		"proxy_status":   "running",
		"total_upload":   0,
		"total_download": 0,
	}

	s.Success(c, status)
}

// getSystemConfig 获取系统配置
func (s *Server) getSystemConfig(c *gin.Context) {
	// 返回安全的配置信息（不包含敏感信息）
	config := map[string]interface{}{
		"api_port":   s.config.API.Port,
		"proxy_port": s.config.Proxy.Port,
		"allow_lan":  s.config.Proxy.AllowLAN,
		"proxy_mode": s.config.Proxy.Mode,
		"log_level":  s.config.Log.Level,
	}

	s.Success(c, config)
}

// updateSystemConfig 更新系统配置
func (s *Server) updateSystemConfig(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		s.Error(c, 1001, "Invalid request body")
		return
	}

	// TODO: 实现配置更新逻辑

	s.Success(c, "Configuration updated successfully")
}

// getProxyStatus 获取代理状态
func (s *Server) getProxyStatus(c *gin.Context) {
	status := s.proxyCore.GetStatus()
	s.Success(c, status)
}

// setProxyMode 设置代理模式
func (s *Server) setProxyMode(c *gin.Context) {
	var req struct {
		Mode string `json:"mode" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		s.Error(c, 1001, "Invalid request body")
		return
	}

	if err := s.proxyCore.SetMode(req.Mode); err != nil {
		s.Error(c, 2001, err.Error())
		return
	}

	s.Success(c, "Proxy mode updated successfully")
}

// getTrafficStats 获取流量统计
func (s *Server) getTrafficStats(c *gin.Context) {
	stats := s.proxyCore.GetTrafficStats()
	s.Success(c, stats)
}

// getNodePools 获取节点池列表
func (s *Server) getNodePools(c *gin.Context) {
	// TODO: 实现节点池管理
	s.Success(c, []interface{}{})
}

// createNodePool 创建节点池
func (s *Server) createNodePool(c *gin.Context) {
	// TODO: 实现节点池创建
	s.Success(c, "Node pool created successfully")
}

// getNodePool 获取节点池详情
func (s *Server) getNodePool(c *gin.Context) {
	// TODO: 实现节点池详情获取
	s.Success(c, map[string]interface{}{})
}

// updateNodePool 更新节点池
func (s *Server) updateNodePool(c *gin.Context) {
	// TODO: 实现节点池更新
	s.Success(c, "Node pool updated successfully")
}

// deleteNodePool 删除节点池
func (s *Server) deleteNodePool(c *gin.Context) {
	// TODO: 实现节点池删除
	s.Success(c, "Node pool deleted successfully")
}

// healthCheck 健康检查
func (s *Server) healthCheck(c *gin.Context) {
	status := "healthy"
	if !s.proxyCore.IsRunning() {
		status = "unhealthy"
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status":  status,
		"version": "1.0.0",
		"uptime":  0,
	})
}

// handleWebSocket 处理 WebSocket 连接
func (s *Server) handleWebSocket(c *gin.Context) {
	// TODO: 实现 WebSocket 处理
	c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "WebSocket not implemented yet",
	})
}
