package handler

import (
	"time"
)

// APIResponse API 统一响应格式
type APIResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(data interface{}) *APIResponse {
	return &APIResponse{
		Code:      0,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(code int, message, detail string) *APIResponse {
	data := map[string]interface{}{
		"detail": detail,
	}
	return &APIResponse{
		Code:      code,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

// getCurrentTime 获取当前时间
func getCurrentTime() time.Time {
	return time.Now()
}

// PaginationResponse 分页响应
type PaginationResponse struct {
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
	Data  interface{} `json:"data"`
}

// ErrorCodes 错误码定义
const (
	// 成功
	CodeSuccess = 0

	// 参数错误 1000-1999
	CodeInvalidParams    = 1001
	CodeAuthFailed       = 1002
	CodePermissionDenied = 1003

	// 资源不存在 2000-2999
	CodeNodePoolNotFound     = 2001
	CodeNodeNotFound         = 2002
	CodeSubscriptionNotFound = 2003

	// 配置错误 3000-3999
	CodeConfigError      = 3001
	CodeProxyStartFailed = 3002

	// 服务器内部错误 5000-5999
	CodeInternalError = 5000
	CodeDatabaseError = 5001
	CodeNetworkError  = 5002
)

// GetErrorMessage 根据错误码获取错误信息
func GetErrorMessage(code int) string {
	messages := map[int]string{
		CodeSuccess:              "成功",
		CodeInvalidParams:        "参数错误",
		CodeAuthFailed:           "认证失败",
		CodePermissionDenied:     "权限不足",
		CodeNodePoolNotFound:     "节点池不存在",
		CodeNodeNotFound:         "节点不存在",
		CodeSubscriptionNotFound: "订阅不存在",
		CodeConfigError:          "配置错误",
		CodeProxyStartFailed:     "代理启动失败",
		CodeInternalError:        "服务器内部错误",
		CodeDatabaseError:        "数据库错误",
		CodeNetworkError:         "网络错误",
	}

	if msg, ok := messages[code]; ok {
		return msg
	}
	return "未知错误"
}
