package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/prism/core/internal/api/handler"
	"github.com/prism/core/internal/api/middleware"
	"github.com/prism/core/internal/service"
)

// Router API 路由器
type Router struct {
	engine              *gin.Engine
	logger              *logrus.Logger
	subscriptionHandler *handler.SubscriptionHandler
}

// NewRouter 创建路由器
func NewRouter(
	logger *logrus.Logger,
	subscriptionSvc *service.SubscriptionService,
	nodePoolSvc *service.NodePoolService,
	nodeSvc *service.NodeService,
	statsSvc *service.StatsService,
	schedulerSvc *service.SchedulerService,
) *Router {
	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()

	// 创建处理器
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionSvc)

	router := &Router{
		engine:              engine,
		logger:              logger,
		subscriptionHandler: subscriptionHandler,
	}

	router.setupMiddleware()
	router.setupRoutes()

	return router
}

// setupMiddleware 设置中间件
func (r *Router) setupMiddleware() {
	// 基础中间件
	r.engine.Use(middleware.CORS())
	r.engine.Use(middleware.Recovery(r.logger))
	r.engine.Use(middleware.RequestLogger(r.logger))
	r.engine.Use(middleware.DefaultRateLimit())
}

// setupRoutes 设置路由
func (r *Router) setupRoutes() {
	// 健康检查
	r.engine.GET("/health", r.healthCheck)
	r.engine.GET("/ping", r.ping)

	// API v1 路由组
	v1 := r.engine.Group("/api/v1")
	{
		// 订阅管理
		subscriptions := v1.Group("/subscriptions")
		{
			subscriptions.POST("", r.subscriptionHandler.CreateSubscription)
			subscriptions.GET("", r.subscriptionHandler.ListSubscriptions)
			subscriptions.GET("/:subscription_id", r.subscriptionHandler.GetSubscription)
			subscriptions.PUT("/:subscription_id", r.subscriptionHandler.UpdateSubscription)
			subscriptions.DELETE("/:subscription_id", r.subscriptionHandler.DeleteSubscription)

			// 订阅操作
			subscriptions.POST("/:subscription_id/update", r.subscriptionHandler.UpdateSubscriptionContent)
			subscriptions.POST("/:subscription_id/enable", r.subscriptionHandler.EnableSubscription)
			subscriptions.POST("/:subscription_id/disable", r.subscriptionHandler.DisableSubscription)

			// 订阅统计和日志
			subscriptions.GET("/:subscription_id/stats", r.subscriptionHandler.GetSubscriptionStats)
			subscriptions.GET("/:subscription_id/logs", r.subscriptionHandler.GetSubscriptionLogs)

			// 批量操作
			subscriptions.POST("/import", r.subscriptionHandler.ImportSubscription)
			subscriptions.GET("/export", r.subscriptionHandler.ExportSubscriptions)
		}

		// 节点池管理（待实现）
		nodepools := v1.Group("/nodepools")
		{
			nodepools.GET("", r.placeholderHandler)
			nodepools.POST("", r.placeholderHandler)
			nodepools.GET("/:pool_id", r.placeholderHandler)
			nodepools.PUT("/:pool_id", r.placeholderHandler)
			nodepools.DELETE("/:pool_id", r.placeholderHandler)
			nodepools.POST("/:pool_id/subscriptions", r.placeholderHandler)
		}

		// 节点管理（待实现）
		nodes := v1.Group("/nodes")
		{
			nodes.GET("", r.placeholderHandler)
			nodes.GET("/:node_id", r.placeholderHandler)
			nodes.POST("/:node_id/test", r.placeholderHandler)
			nodes.POST("/batch-test", r.placeholderHandler)
			nodes.GET("/test-tasks/:task_id", r.placeholderHandler)
			nodes.GET("/:node_id/test-history", r.placeholderHandler)
			nodes.GET("/best-selection", r.placeholderHandler)
		}

		// 统计分析（待实现）
		stats := v1.Group("/stats")
		{
			stats.GET("/overview", r.placeholderHandler)
			stats.GET("/geo-distribution", r.placeholderHandler)
			stats.GET("/protocol-distribution", r.placeholderHandler)
			stats.GET("/performance-trend", r.placeholderHandler)
		}

		// 自动化任务（待实现）
		tasks := v1.Group("/tasks")
		{
			tasks.GET("/auto-update", r.placeholderHandler)
			tasks.POST("/auto-update/trigger", r.placeholderHandler)
			tasks.GET("/scheduled-test", r.placeholderHandler)
		}

		// 系统管理（待实现）
		system := v1.Group("/system")
		{
			system.GET("/status", r.placeholderHandler)
			system.GET("/config", r.placeholderHandler)
			system.PUT("/config", r.placeholderHandler)
		}
	}
}

// healthCheck 健康检查
func (r *Router) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"version":   "1.0.0",
	})
}

// ping 简单的 ping 接口
func (r *Router) ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":   "pong",
		"timestamp": time.Now().Unix(),
	})
}

// placeholderHandler 占位处理器
func (r *Router) placeholderHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"code":    5001,
		"message": "功能尚未实现",
		"data": gin.H{
			"detail": "该接口正在开发中",
		},
		"timestamp": time.Now().Unix(),
	})
}

// GetEngine 获取 Gin 引擎
func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}

// Run 启动服务器
func (r *Router) Run(addr string) error {
	r.logger.Infof("Starting API server on %s", addr)
	return r.engine.Run(addr)
}
