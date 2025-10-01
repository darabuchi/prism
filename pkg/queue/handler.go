package queue

import "time"

// HandlerResult 处理器返回结果
type HandlerResult struct {
	// 是否需要重试
	Retry bool

	// 错误信息（可选）
	Error error

	// 自定义重试延迟（可选，为 0 使用策略计算）
	RetryDelay time.Duration
}

// Success 返回成功结果
func Success() *HandlerResult {
	return &HandlerResult{
		Retry: false,
		Error: nil,
	}
}

// Fail 返回失败但不重试的结果
func Fail(err error) *HandlerResult {
	return &HandlerResult{
		Retry: false,
		Error: err,
	}
}

// RetryWithError 返回需要重试的结果（带错误信息）
func RetryWithError(err error) *HandlerResult {
	return &HandlerResult{
		Retry: true,
		Error: err,
	}
}

// RetryWithDelay 返回需要重试的结果（自定义延迟）
func RetryWithDelay(delay time.Duration, err error) *HandlerResult {
	return &HandlerResult{
		Retry:      true,
		Error:      err,
		RetryDelay: delay,
	}
}

// ShouldRetry 判断是否需要重试
func (r *HandlerResult) ShouldRetry() bool {
	return r != nil && r.Retry
}

// IsSuccess 判断是否成功（无错误且不重试）
func (r *HandlerResult) IsSuccess() bool {
	return r != nil && !r.Retry && r.Error == nil
}
