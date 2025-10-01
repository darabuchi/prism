package i18n

// keyToCodeMap 将 i18n 键映射到错误码
// 这个映射表用于将 YAML 文件中的键（如 "common.success"）映射到对应的错误码
var keyToCodeMap = map[string]int32{
	// 成功
	"common.success": 0,

	// 通用错误 (1000-1999)
	"common.internal_error":       1000,
	"common.invalid_request":      1001,
	"common.not_found":            1002,
	"common.already_exists":       1003,
	"common.timeout":              1004,
	"common.too_many_requests":    1005,
	"common.service_unavailable":  1006,
	"common.method_not_allowed":   1007,
	"common.unprocessable_entity": 1008,

	// 认证授权 (2000-2999)
	"auth.unauthorized":        2000,
	"auth.forbidden":           2001,
	"auth.invalid_token":       2002,
	"auth.token_expired":       2003,
	"auth.invalid_credentials": 2004,

	// 订阅相关 (3000-3099)
	"subscription.not_found":     3000,
	"subscription.already_exists": 3001,
	"subscription.invalid":       3002,
	"subscription.update_failed": 3003,

	// 节点相关 (3100-3199)
	"node.not_found":   3100,
	"node.test_failed": 3101,
	"node.unavailable": 3102,

	// 路由相关 (3200-3299)
	"route.not_found": 3200,
	"route.invalid":   3201,

	// 验证错误 (4000-4999)
	"validation.failed":          4000,
	"validation.url_required":    4001,
	"validation.url_invalid":     4002,
	"validation.title_required":  4003,
	"validation.port_invalid":    4004,
	"validation.email_invalid":   4005,
	"validation.field_required":  4006,
	"validation.field_too_long":  4007,
	"validation.field_too_short": 4008,

	// 服务错误 (5000-5999)
	"service.database_error":    5000,
	"service.cache_error":       5001,
	"service.connection_failed": 5002,

	// 队列错误 (6000-6999)
	"queue.full":                   6000,
	"queue.empty":                  6001,
	"queue.closed":                 6002,
	"queue.timeout":                6003,
	"queue.invalid_concurrency":    6004,
	"queue.unsupported_type":       6005,
	"queue.max_retries_exceeded":   6006,

	// 配置错误 (7000-7999)
	"config.invalid":      7000,
	"config.not_found":    7001,
	"config.parse_failed": 7002,

	// 代理错误 (8000-8099)
	"proxy.connection_failed": 8000,
	"proxy.timeout":           8001,
	"proxy.auth_failed":       8002,

	// 网络错误 (8100-8199)
	"network.unreachable":        8100,
	"network.dns_resolve_failed": 8101,
}
