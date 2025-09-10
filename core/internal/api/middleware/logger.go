package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// responseBodyWriter 包装响应写入器以捕获响应体
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 自定义日志格式
		return ""
	})
}

// RequestLogger 详细请求日志中间件
func RequestLogger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// 读取请求体
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			requestBody = string(bodyBytes)
			c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		// 包装响应写入器
		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		// 处理请求
		c.Next()

		// 计算延迟
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		// 日志字段
		fields := logrus.Fields{
			"method":     method,
			"path":       path,
			"status":     statusCode,
			"latency":    latency.String(),
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}

		// 添加请求体（仅对非GET请求且小于1KB的情况）
		if method != "GET" && len(requestBody) > 0 && len(requestBody) < 1024 {
			fields["request_body"] = requestBody
		}

		// 添加响应体（仅对错误响应）
		if statusCode >= 400 && w.body.Len() < 1024 {
			fields["response_body"] = w.body.String()
		}

		// 根据状态码确定日志级别
		if statusCode >= 500 {
			logger.WithFields(fields).Error("Server Error")
		} else if statusCode >= 400 {
			logger.WithFields(fields).Warn("Client Error")
		} else {
			logger.WithFields(fields).Info("Request Processed")
		}
	}
}
