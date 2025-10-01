package errcode

import "net/http"

// 错误码范围定义：
// 1000-1999: 通用错误
// 2000-2999: 认证和授权错误
// 3000-3999: 资源错误（订阅、节点、路由等）
// 4000-4999: 验证错误
// 5000-5999: 服务错误
// 6000-6999: 队列错误
// 7000-7999: 配置错误
// 8000-8999: 网络和代理错误

var (
	// ==================== 通用错误 (1000-1999) ====================

	// Success 成功
	Success = &ErrCode{Code: 0, Key: "common.success", HTTPCode: http.StatusOK}

	// ErrInternal 内部错误
	ErrInternal = &ErrCode{Code: 1000, Key: "common.internal_error", HTTPCode: http.StatusInternalServerError}

	// ErrInvalidRequest 无效的请求
	ErrInvalidRequest = &ErrCode{Code: 1001, Key: "common.invalid_request", HTTPCode: http.StatusBadRequest}

	// ErrNotFound 资源未找到
	ErrNotFound = &ErrCode{Code: 1002, Key: "common.not_found", HTTPCode: http.StatusNotFound}

	// ErrAlreadyExists 资源已存在
	ErrAlreadyExists = &ErrCode{Code: 1003, Key: "common.already_exists", HTTPCode: http.StatusConflict}

	// ErrTimeout 操作超时
	ErrTimeout = &ErrCode{Code: 1004, Key: "common.timeout", HTTPCode: http.StatusRequestTimeout}

	// ErrTooManyRequests 请求过多
	ErrTooManyRequests = &ErrCode{Code: 1005, Key: "common.too_many_requests", HTTPCode: http.StatusTooManyRequests}

	// ErrServiceUnavailable 服务不可用
	ErrServiceUnavailable = &ErrCode{Code: 1006, Key: "common.service_unavailable", HTTPCode: http.StatusServiceUnavailable}

	// ErrMethodNotAllowed 方法不允许
	ErrMethodNotAllowed = &ErrCode{Code: 1007, Key: "common.method_not_allowed", HTTPCode: http.StatusMethodNotAllowed}

	// ErrUnprocessableEntity 无法处理的实体
	ErrUnprocessableEntity = &ErrCode{Code: 1008, Key: "common.unprocessable_entity", HTTPCode: http.StatusUnprocessableEntity}

	// ==================== 认证和授权错误 (2000-2999) ====================

	// ErrUnauthorized 未授权
	ErrUnauthorized = &ErrCode{Code: 2000, Key: "auth.unauthorized", HTTPCode: http.StatusUnauthorized}

	// ErrForbidden 禁止访问
	ErrForbidden = &ErrCode{Code: 2001, Key: "auth.forbidden", HTTPCode: http.StatusForbidden}

	// ErrInvalidToken 无效的令牌
	ErrInvalidToken = &ErrCode{Code: 2002, Key: "auth.invalid_token", HTTPCode: http.StatusUnauthorized}

	// ErrTokenExpired 令牌已过期
	ErrTokenExpired = &ErrCode{Code: 2003, Key: "auth.token_expired", HTTPCode: http.StatusUnauthorized}

	// ErrInvalidCredentials 无效的凭证
	ErrInvalidCredentials = &ErrCode{Code: 2004, Key: "auth.invalid_credentials", HTTPCode: http.StatusUnauthorized}

	// ==================== 资源错误 (3000-3999) ====================

	// ErrSubscriptionNotFound 订阅不存在
	ErrSubscriptionNotFound = &ErrCode{Code: 3000, Key: "subscription.not_found", HTTPCode: http.StatusNotFound}

	// ErrSubscriptionExists 订阅已存在
	ErrSubscriptionExists = &ErrCode{Code: 3001, Key: "subscription.already_exists", HTTPCode: http.StatusConflict}

	// ErrSubscriptionInvalid 订阅无效
	ErrSubscriptionInvalid = &ErrCode{Code: 3002, Key: "subscription.invalid", HTTPCode: http.StatusBadRequest}

	// ErrSubscriptionUpdateFailed 订阅更新失败
	ErrSubscriptionUpdateFailed = &ErrCode{Code: 3003, Key: "subscription.update_failed", HTTPCode: http.StatusInternalServerError}

	// ErrNodeNotFound 节点不存在
	ErrNodeNotFound = &ErrCode{Code: 3100, Key: "node.not_found", HTTPCode: http.StatusNotFound}

	// ErrNodeTestFailed 节点测试失败
	ErrNodeTestFailed = &ErrCode{Code: 3101, Key: "node.test_failed", HTTPCode: http.StatusInternalServerError}

	// ErrNodeUnavailable 节点不可用
	ErrNodeUnavailable = &ErrCode{Code: 3102, Key: "node.unavailable", HTTPCode: http.StatusServiceUnavailable}

	// ErrRouteNotFound 路由规则不存在
	ErrRouteNotFound = &ErrCode{Code: 3200, Key: "route.not_found", HTTPCode: http.StatusNotFound}

	// ErrRouteInvalid 路由规则无效
	ErrRouteInvalid = &ErrCode{Code: 3201, Key: "route.invalid", HTTPCode: http.StatusBadRequest}

	// ==================== 验证错误 (4000-4999) ====================

	// ErrValidationFailed 验证失败
	ErrValidationFailed = &ErrCode{Code: 4000, Key: "validation.failed", HTTPCode: http.StatusBadRequest}

	// ErrURLRequired URL 不能为空
	ErrURLRequired = &ErrCode{Code: 4001, Key: "validation.url_required", HTTPCode: http.StatusBadRequest}

	// ErrURLInvalid URL 格式无效
	ErrURLInvalid = &ErrCode{Code: 4002, Key: "validation.url_invalid", HTTPCode: http.StatusBadRequest}

	// ErrTitleRequired 标题不能为空
	ErrTitleRequired = &ErrCode{Code: 4003, Key: "validation.title_required", HTTPCode: http.StatusBadRequest}

	// ErrPortInvalid 端口号无效
	ErrPortInvalid = &ErrCode{Code: 4004, Key: "validation.port_invalid", HTTPCode: http.StatusBadRequest}

	// ErrEmailInvalid 邮箱格式无效
	ErrEmailInvalid = &ErrCode{Code: 4005, Key: "validation.email_invalid", HTTPCode: http.StatusBadRequest}

	// ErrFieldRequired 字段不能为空
	ErrFieldRequired = &ErrCode{Code: 4006, Key: "validation.field_required", HTTPCode: http.StatusBadRequest}

	// ErrFieldTooLong 字段过长
	ErrFieldTooLong = &ErrCode{Code: 4007, Key: "validation.field_too_long", HTTPCode: http.StatusBadRequest}

	// ErrFieldTooShort 字段过短
	ErrFieldTooShort = &ErrCode{Code: 4008, Key: "validation.field_too_short", HTTPCode: http.StatusBadRequest}

	// ==================== 服务错误 (5000-5999) ====================

	// ErrDatabaseError 数据库错误
	ErrDatabaseError = &ErrCode{Code: 5000, Key: "service.database_error", HTTPCode: http.StatusInternalServerError}

	// ErrCacheError 缓存错误
	ErrCacheError = &ErrCode{Code: 5001, Key: "service.cache_error", HTTPCode: http.StatusInternalServerError}

	// ErrConnectionFailed 连接失败
	ErrConnectionFailed = &ErrCode{Code: 5002, Key: "service.connection_failed", HTTPCode: http.StatusServiceUnavailable}

	// ==================== 队列错误 (6000-6999) ====================

	// ErrQueueFull 队列已满
	ErrQueueFull = &ErrCode{Code: 6000, Key: "queue.full", HTTPCode: http.StatusServiceUnavailable}

	// ErrQueueEmpty 队列为空
	ErrQueueEmpty = &ErrCode{Code: 6001, Key: "queue.empty", HTTPCode: http.StatusNotFound}

	// ErrQueueClosed 队列已关闭
	ErrQueueClosed = &ErrCode{Code: 6002, Key: "queue.closed", HTTPCode: http.StatusServiceUnavailable}

	// ErrQueueTimeout 队列操作超时
	ErrQueueTimeout = &ErrCode{Code: 6003, Key: "queue.timeout", HTTPCode: http.StatusRequestTimeout}

	// ErrQueueInvalidConcurrency 无效的并发数
	ErrQueueInvalidConcurrency = &ErrCode{Code: 6004, Key: "queue.invalid_concurrency", HTTPCode: http.StatusBadRequest}

	// ErrQueueUnsupportedType 不支持的队列类型
	ErrQueueUnsupportedType = &ErrCode{Code: 6005, Key: "queue.unsupported_type", HTTPCode: http.StatusBadRequest}

	// ErrQueueMaxRetriesExceeded 超过最大重试次数
	ErrQueueMaxRetriesExceeded = &ErrCode{Code: 6006, Key: "queue.max_retries_exceeded", HTTPCode: http.StatusInternalServerError}

	// ==================== 配置错误 (7000-7999) ====================

	// ErrConfigInvalid 配置无效
	ErrConfigInvalid = &ErrCode{Code: 7000, Key: "config.invalid", HTTPCode: http.StatusBadRequest}

	// ErrConfigNotFound 配置未找到
	ErrConfigNotFound = &ErrCode{Code: 7001, Key: "config.not_found", HTTPCode: http.StatusNotFound}

	// ErrConfigParseFailed 配置解析失败
	ErrConfigParseFailed = &ErrCode{Code: 7002, Key: "config.parse_failed", HTTPCode: http.StatusInternalServerError}

	// ==================== 网络和代理错误 (8000-8999) ====================

	// ErrProxyConnectionFailed 代理连接失败
	ErrProxyConnectionFailed = &ErrCode{Code: 8000, Key: "proxy.connection_failed", HTTPCode: http.StatusServiceUnavailable}

	// ErrProxyTimeout 代理超时
	ErrProxyTimeout = &ErrCode{Code: 8001, Key: "proxy.timeout", HTTPCode: http.StatusGatewayTimeout}

	// ErrProxyAuthFailed 代理认证失败
	ErrProxyAuthFailed = &ErrCode{Code: 8002, Key: "proxy.auth_failed", HTTPCode: http.StatusProxyAuthRequired}

	// ErrNetworkUnreachable 网络不可达
	ErrNetworkUnreachable = &ErrCode{Code: 8003, Key: "network.unreachable", HTTPCode: http.StatusServiceUnavailable}

	// ErrDNSResolveFailed DNS 解析失败
	ErrDNSResolveFailed = &ErrCode{Code: 8004, Key: "network.dns_resolve_failed", HTTPCode: http.StatusServiceUnavailable}
)
