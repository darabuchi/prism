package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter 速率限制器
type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int           // 每分钟允许的请求数
	window   time.Duration // 时间窗口
}

type visitor struct {
	count    int
	lastSeen time.Time
}

// NewRateLimiter 创建速率限制器
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		window:   window,
	}

	// 启动清理协程
	go rl.cleanupVisitors()

	return rl
}

// RateLimit 速率限制中间件
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		rl.mu.Lock()
		v, exists := rl.visitors[ip]
		if !exists {
			rl.visitors[ip] = &visitor{
				count:    1,
				lastSeen: time.Now(),
			}
			rl.mu.Unlock()
			c.Next()
			return
		}

		// 检查时间窗口
		if time.Since(v.lastSeen) > rl.window {
			v.count = 1
			v.lastSeen = time.Now()
		} else {
			v.count++
		}

		if v.count > rl.rate {
			rl.mu.Unlock()
			c.JSON(http.StatusTooManyRequests, map[string]interface{}{
				"code":      4291,
				"message":   "请求过于频繁",
				"data":      map[string]string{"detail": "请稍后再试"},
				"timestamp": time.Now().Unix(),
			})
			c.Abort()
			return
		}

		rl.mu.Unlock()
		c.Next()
	}
}

// cleanupVisitors 清理过期的访问者记录
func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			for ip, v := range rl.visitors {
				if time.Since(v.lastSeen) > rl.window*2 {
					delete(rl.visitors, ip)
				}
			}
			rl.mu.Unlock()
		}
	}
}

// DefaultRateLimit 默认速率限制中间件（100次/分钟）
func DefaultRateLimit() gin.HandlerFunc {
	return NewRateLimiter(100, time.Minute).RateLimit()
}
