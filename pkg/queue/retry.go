package queue

import (
	"math"
	"time"
)

// RetryStrategy 重试策略类型
type RetryStrategy string

const (
	// RetryStrategyNone 不重试
	RetryStrategyNone RetryStrategy = "none"

	// RetryStrategyFixed 固定间隔重试
	RetryStrategyFixed RetryStrategy = "fixed"

	// RetryStrategyLinear 线性递增重试
	RetryStrategyLinear RetryStrategy = "linear"

	// RetryStrategyExponential 指数退避重试
	RetryStrategyExponential RetryStrategy = "exponential"
)

// IsValid 检查重试策略是否有效
func (s RetryStrategy) IsValid() bool {
	switch s {
	case RetryStrategyNone, RetryStrategyFixed, RetryStrategyLinear, RetryStrategyExponential:
		return true
	default:
		return false
	}
}

// String 返回重试策略字符串
func (s RetryStrategy) String() string {
	return string(s)
}

// RetryPolicy 重试策略配置
type RetryPolicy struct {
	// 重试策略类型
	Strategy RetryStrategy `json:"strategy" yaml:"strategy" validate:"required,oneof=none fixed linear exponential" default:"exponential"`

	// 最大重试次数（0 表示不限制）
	MaxRetries int `json:"max_retries" yaml:"max_retries" validate:"min=0" default:"3"`

	// 初始延迟（毫秒）
	InitialDelay int `json:"initial_delay" yaml:"initial_delay" validate:"min=0" default:"1000"`

	// 最大延迟（毫秒，0 表示不限制）
	MaxDelay int `json:"max_delay" yaml:"max_delay" validate:"min=0" default:"60000"`

	// 延迟倍数（指数退避时使用）
	Multiplier float64 `json:"multiplier" yaml:"multiplier" validate:"min=1" default:"2.0"`

	// 增量（线性递增时使用，毫秒）
	Increment int `json:"increment" yaml:"increment" validate:"min=0" default:"1000"`
}

// ApplyDefaults 应用默认值
func (p *RetryPolicy) ApplyDefaults() {
	if p.Strategy == "" {
		p.Strategy = RetryStrategyExponential
	}
	if p.MaxRetries == 0 {
		p.MaxRetries = 3
	}
	if p.InitialDelay == 0 {
		p.InitialDelay = 1000
	}
	if p.MaxDelay == 0 {
		p.MaxDelay = 60000
	}
	if p.Multiplier == 0 {
		p.Multiplier = 2.0
	}
	if p.Increment == 0 {
		p.Increment = 1000
	}
}

// Validate 验证重试策略配置
func (p *RetryPolicy) Validate() error {
	if !p.Strategy.IsValid() {
		return ErrInvalidConfig
	}

	if p.MaxRetries < 0 {
		return ErrInvalidConfig
	}

	if p.InitialDelay < 0 {
		return ErrInvalidConfig
	}

	if p.Multiplier < 1 {
		return ErrInvalidConfig
	}

	return nil
}

// CalculateDelay 计算下次重试延迟时间
func (p *RetryPolicy) CalculateDelay(retryCount int) time.Duration {
	if retryCount < 0 {
		retryCount = 0
	}

	var delayMs int

	switch p.Strategy {
	case RetryStrategyNone:
		return 0

	case RetryStrategyFixed:
		// 固定延迟
		delayMs = p.InitialDelay

	case RetryStrategyLinear:
		// 线性递增：initialDelay + increment * retryCount
		delayMs = p.InitialDelay + p.Increment*retryCount

	case RetryStrategyExponential:
		// 指数退避：initialDelay * (multiplier ^ retryCount)
		delay := float64(p.InitialDelay) * math.Pow(p.Multiplier, float64(retryCount))
		delayMs = int(delay)

	default:
		delayMs = p.InitialDelay
	}

	// 应用最大延迟限制
	if p.MaxDelay > 0 && delayMs > p.MaxDelay {
		delayMs = p.MaxDelay
	}

	return time.Duration(delayMs) * time.Millisecond
}

// ShouldRetry 判断是否应该重试
func (p *RetryPolicy) ShouldRetry(retryCount int) bool {
	if p.Strategy == RetryStrategyNone {
		return false
	}

	if p.MaxRetries == 0 {
		// 0 表示不限制重试次数
		return true
	}

	return retryCount < p.MaxRetries
}

// DefaultRetryPolicy 默认重试策略
func DefaultRetryPolicy() *RetryPolicy {
	policy := &RetryPolicy{}
	policy.ApplyDefaults()
	return policy
}

// NoRetryPolicy 不重试策略
func NoRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		Strategy:   RetryStrategyNone,
		MaxRetries: 0,
	}
}

// FixedRetryPolicy 固定间隔重试策略
func FixedRetryPolicy(maxRetries int, delayMs int) *RetryPolicy {
	return &RetryPolicy{
		Strategy:     RetryStrategyFixed,
		MaxRetries:   maxRetries,
		InitialDelay: delayMs,
		MaxDelay:     delayMs,
	}
}

// LinearRetryPolicy 线性递增重试策略
func LinearRetryPolicy(maxRetries int, initialDelayMs int, incrementMs int) *RetryPolicy {
	return &RetryPolicy{
		Strategy:     RetryStrategyLinear,
		MaxRetries:   maxRetries,
		InitialDelay: initialDelayMs,
		Increment:    incrementMs,
		MaxDelay:     60000, // 默认最大 60 秒
	}
}

// ExponentialRetryPolicy 指数退避重试策略
func ExponentialRetryPolicy(maxRetries int, initialDelayMs int, multiplier float64) *RetryPolicy {
	return &RetryPolicy{
		Strategy:     RetryStrategyExponential,
		MaxRetries:   maxRetries,
		InitialDelay: initialDelayMs,
		Multiplier:   multiplier,
		MaxDelay:     60000, // 默认最大 60 秒
	}
}
