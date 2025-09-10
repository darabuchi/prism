package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prism/core/pkg/logger"
	"github.com/sirupsen/logrus"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		log := logger.GetLogger()
		log.WithFields(logrus.Fields{
			"client_ip":   param.ClientIP,
			"method":      param.Method,
			"path":        param.Path,
			"status_code": param.StatusCode,
			"latency":     param.Latency,
			"user_agent":  param.Request.UserAgent(),
		}).Info("API Request")
		return ""
	})
}

// CORSMiddleware CORS 中间件
func CORSMiddleware(allowOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// 检查是否允许该来源
		allowed := false
		for _, allowedOrigin := range allowOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}
		
		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// AuthMiddleware 认证中间件
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现 JWT 认证
		// 这里暂时跳过认证，后续可以添加 JWT 验证
		c.Next()
	}
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现限流逻辑
		// 可以使用 go-rate 或其他限流库
		c.Next()
	}
}