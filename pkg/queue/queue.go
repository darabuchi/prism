package queue

import "time"

// Queue 队列接口（泛型）
type Queue[T any] interface {
	// Push 推送单个消息到队列
	Push(payload T) error

	// PushWithDelay 推送延时消息到队列
	// delay: 延时时长
	PushWithDelay(payload T, delay time.Duration) error

	// PushWithRetryPolicy 推送带自定义重试策略的消息
	PushWithRetryPolicy(payload T, policy *RetryPolicy) error

	// BatchPush 批量推送消息到队列
	BatchPush(payloads []T) error

	// Pop 从队列中弹出一个消息
	// 如果队列为空，返回 nil 和 ErrQueueEmpty
	Pop() (*Message[T], error)

	// Process 启动消费者处理队列消息
	// concurrency: 并发处理数量
	// handler: 消息处理函数，返回 HandlerResult 明确指示是否重试
	Process(concurrency int, handler func(*Message[T]) *HandlerResult) error

	// Depth 获取当前队列深度
	Depth() int

	// Close 关闭队列
	Close() error
}

// New 创建队列实例
func New[T any](cfg *Config) (Queue[T], error) {
	if cfg == nil {
		return nil, ErrInvalidConfig
	}

	// 应用默认值
	cfg.ApplyDefaults()

	// 验证配置
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// 根据类型创建队列
	switch cfg.Type {
	case TypeMemory:
		return NewMemoryQueue[T](cfg)
	case TypeRedis:
		return NewRedisQueue[T](cfg)
	case TypeNSQ:
		return NewNSQQueue[T](cfg)
	case TypeKafka:
		return NewKafkaQueue[T](cfg)
	case TypeRabbitMQ:
		return NewRabbitMQQueue[T](cfg)
	case TypeZeroMQ:
		return NewZeroMQQueue[T](cfg)
	default:
		return nil, ErrUnsupportedQueueType
	}
}
