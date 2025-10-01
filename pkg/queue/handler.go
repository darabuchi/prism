package queue

import "time"

// RetryInfo 重试信息
type RetryInfo struct {
	// ShouldRetry 是否需要重试
	ShouldRetry bool

	// Delay 自定义重试延迟（可选，为 0 使用策略计算）
	Delay time.Duration
}

// NoRetry 返回不重试的信息
func NoRetry() *RetryInfo {
	return &RetryInfo{
		ShouldRetry: false,
	}
}

// Retry 返回需要重试的信息（使用策略计算延迟）
func Retry() *RetryInfo {
	return &RetryInfo{
		ShouldRetry: true,
	}
}

// RetryAfter 返回需要重试的信息（自定义延迟）
func RetryAfter(delay time.Duration) *RetryInfo {
	return &RetryInfo{
		ShouldRetry: true,
		Delay:       delay,
	}
}
