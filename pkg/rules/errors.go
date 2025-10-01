package rules

import "errors"

var (
	// ErrInvalidRule 无效的规则格式
	ErrInvalidRule = errors.New("invalid rule format")

	// ErrInvalidRuleType 无效的规则类型
	ErrInvalidRuleType = errors.New("invalid rule type")

	// ErrInvalidAction 无效的动作类型
	ErrInvalidAction = errors.New("invalid action type")

	// ErrInvalidPayload 无效的载荷内容
	ErrInvalidPayload = errors.New("invalid payload")

	// ErrInvalidIPCIDR 无效的 IP CIDR
	ErrInvalidIPCIDR = errors.New("invalid IP CIDR")

	// ErrInvalidRegex 无效的正则表达式
	ErrInvalidRegex = errors.New("invalid regex pattern")

	// ErrInvalidPortRange 无效的端口范围
	ErrInvalidPortRange = errors.New("invalid port range")
)
