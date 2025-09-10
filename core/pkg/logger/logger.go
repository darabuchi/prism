package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var globalLogger *logrus.Logger

// Init 初始化日志系统
func Init(level, format string) {
	globalLogger = logrus.New()

	// 设置日志级别
	switch level {
	case "debug":
		globalLogger.SetLevel(logrus.DebugLevel)
	case "info":
		globalLogger.SetLevel(logrus.InfoLevel)
	case "warn":
		globalLogger.SetLevel(logrus.WarnLevel)
	case "error":
		globalLogger.SetLevel(logrus.ErrorLevel)
	default:
		globalLogger.SetLevel(logrus.InfoLevel)
	}

	// 设置日志格式
	switch format {
	case "json":
		globalLogger.SetFormatter(&logrus.JSONFormatter{})
	default:
		globalLogger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	// 设置输出
	globalLogger.SetOutput(os.Stdout)
}

// GetLogger 获取全局日志实例
func GetLogger() *logrus.Logger {
	if globalLogger == nil {
		Init("info", "text")
	}
	return globalLogger
}

// Debug 调试日志
func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

// Info 信息日志
func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

// Warn 警告日志
func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

// Error 错误日志
func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

// Fatal 致命错误日志
func Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}

// WithFields 带字段的日志
func WithFields(fields logrus.Fields) *logrus.Entry {
	return GetLogger().WithFields(fields)
}