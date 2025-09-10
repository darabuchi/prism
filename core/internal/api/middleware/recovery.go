package middleware

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Recovery 恢复中间件
func Recovery(logger *logrus.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// 获取堆栈信息
		stack := make([]byte, 4<<10) // 4KB
		length := runtime.Stack(stack, false)

		// 记录panic信息
		logger.WithFields(logrus.Fields{
			"method":    c.Request.Method,
			"path":      c.Request.URL.Path,
			"client_ip": c.ClientIP(),
			"error":     recovered,
			"stack":     string(stack[:length]),
		}).Error("Panic recovered")

		// 返回统一的错误响应
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":      5000,
			"message":   "服务器内部错误",
			"data":      map[string]string{"detail": fmt.Sprintf("%v", recovered)},
			"timestamp": getCurrentTimestamp(),
		})

		c.Abort()
	})
}

func getCurrentTimestamp() int64 {
	return time.Now().Unix()
}
