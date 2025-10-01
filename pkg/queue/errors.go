package queue

import "errors"

var (
	// ErrUnsupportedQueueType 不支持的队列类型
	ErrUnsupportedQueueType = errors.New("unsupported queue type")

	// ErrQueueClosed 队列已关闭
	ErrQueueClosed = errors.New("queue is closed")

	// ErrQueueEmpty 队列为空
	ErrQueueEmpty = errors.New("queue is empty")

	// ErrQueueFull 队列已满
	ErrQueueFull = errors.New("queue is full")

	// ErrTimeout 操作超时
	ErrTimeout = errors.New("operation timeout")

	// ErrMaxRetriesExceeded 超过最大重试次数
	ErrMaxRetriesExceeded = errors.New("maximum retries exceeded")

	// ErrInvalidConfig 无效的配置
	ErrInvalidConfig = errors.New("invalid configuration")

	// ErrConnectionFailed 连接失败
	ErrConnectionFailed = errors.New("connection failed")

	// ErrInvalidConcurrency 无效的并发数
	ErrInvalidConcurrency = errors.New("invalid concurrency")

	// ErrQueueDisabled 队列已禁用
	ErrQueueDisabled = errors.New("queue is disabled")

	// ErrConfigNotFound 配置未找到
	ErrConfigNotFound = errors.New("configuration not found")
)
