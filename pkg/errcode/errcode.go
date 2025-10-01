package errcode

import (
	"github.com/darabuchi/prism/pkg/i18n"
	"github.com/lazygophers/lrpc/middleware/xerror"
)

// 错误码定义
const (
	// 成功
	Success int32 = 0

	// 通用错误 (1000-1999)
	ErrInternal            int32 = 1000
	ErrInvalidRequest      int32 = 1001
	ErrNotFound            int32 = 1002
	ErrAlreadyExists       int32 = 1003
	ErrTimeout             int32 = 1004
	ErrTooManyRequests     int32 = 1005
	ErrServiceUnavailable  int32 = 1006
	ErrMethodNotAllowed    int32 = 1007
	ErrUnprocessableEntity int32 = 1008

	// 认证授权 (2000-2999)
	ErrUnauthorized       int32 = 2000
	ErrForbidden          int32 = 2001
	ErrInvalidToken       int32 = 2002
	ErrTokenExpired       int32 = 2003
	ErrInvalidCredentials int32 = 2004

	// 订阅相关 (3000-3099)
	ErrSubscriptionNotFound int32 = 3000
	ErrSubscriptionExists   int32 = 3001
	ErrSubscriptionInvalid  int32 = 3002
	ErrSubscriptionUpdateFailed int32 = 3003

	// 节点相关 (3100-3199)
	ErrNodeNotFound   int32 = 3100
	ErrNodeTestFailed int32 = 3101
	ErrNodeUnavailable int32 = 3102

	// 路由相关 (3200-3299)
	ErrRouteNotFound int32 = 3200
	ErrRouteInvalid  int32 = 3201

	// 验证错误 (4000-4999)
	ErrValidationFailed int32 = 4000
	ErrURLRequired      int32 = 4001
	ErrURLInvalid       int32 = 4002
	ErrTitleRequired    int32 = 4003
	ErrPortInvalid      int32 = 4004
	ErrEmailInvalid     int32 = 4005
	ErrFieldRequired    int32 = 4006
	ErrFieldTooLong     int32 = 4007
	ErrFieldTooShort    int32 = 4008

	// 服务错误 (5000-5999)
	ErrDatabaseError    int32 = 5000
	ErrCacheError       int32 = 5001
	ErrConnectionFailed int32 = 5002

	// 队列错误 (6000-6999)
	ErrQueueFull                int32 = 6000
	ErrQueueEmpty               int32 = 6001
	ErrQueueClosed              int32 = 6002
	ErrQueueTimeout             int32 = 6003
	ErrQueueInvalidConcurrency  int32 = 6004
	ErrQueueUnsupportedType     int32 = 6005
	ErrQueueMaxRetriesExceeded  int32 = 6006

	// 配置错误 (7000-7999)
	ErrConfigInvalid     int32 = 7000
	ErrConfigNotFound    int32 = 7001
	ErrConfigParseFailed int32 = 7002

	// 代理错误 (8000-8099)
	ErrProxyConnectionFailed int32 = 8000
	ErrProxyTimeout          int32 = 8001
	ErrProxyAuthFailed       int32 = 8002

	// 网络错误 (8100-8199)
	ErrNetworkUnreachable   int32 = 8100
	ErrDNSResolveFailed     int32 = 8101
)

// Init 初始化错误码系统
// dirPath: i18n 资源文件目录路径
func Init(dirPath string) error {
	// 加载多语言资源
	if err := i18n.LoadFromDir(dirPath); err != nil {
		return err
	}

	// 设置 xerror 的 i18n 实现
	xerror.SetI18n(i18n.NewLocalizer())

	return nil
}

// New 创建带本地化消息的错误
func New(code int32, langs ...string) *xerror.Error {
	return xerror.NewError(code, langs...)
}

// NewWithMsg 创建带自定义消息的错误
func NewWithMsg(code int32, msg string) *xerror.Error {
	return xerror.NewErrorWithMsg(code, msg)
}
