package queue

import "errors"

// RabbitMQQueue RabbitMQ队列实现
type RabbitMQQueue[T any] struct {
	cfg *Config
	// TODO: 添加 RabbitMQ connection 和 channel
	// conn *amqp.Connection
	// channel *amqp.Channel
}

// NewRabbitMQQueue 创建 RabbitMQ 队列
func NewRabbitMQQueue[T any](cfg *Config) (*RabbitMQQueue[T], error) {
	if cfg.Address == "" {
		return nil, ErrInvalidConfig
	}

	// TODO: 初始化 RabbitMQ 连接
	// 依赖: github.com/rabbitmq/amqp091-go

	return &RabbitMQQueue[T]{
		cfg: cfg,
	}, errors.New("RabbitMQ queue not yet implemented")
}

// Push 推送单个消息到队列
func (q *RabbitMQQueue[T]) Push(payload T) error {
	// TODO: 使用 Channel.Publish() 推送消息
	// err := q.channel.Publish(
	//     "",        // exchange
	//     queueName, // routing key
	//     false,     // mandatory
	//     false,     // immediate
	//     amqp.Publishing{
	//         ContentType: "application/json",
	//         Body:        serialize(payload),
	//     },
	// )
	return errors.New("not implemented")
}

// BatchPush 批量推送消息到队列
func (q *RabbitMQQueue[T]) BatchPush(payloads []T) error {
	// TODO: 批量推送消息
	for _, payload := range payloads {
		if err := q.Push(payload); err != nil {
			return err
		}
	}
	return nil
}

// Pop 从队列中弹出一个消息
func (q *RabbitMQQueue[T]) Pop() (*Message[T], error) {
	// TODO: 使用 Channel.Get() 获取单个消息
	// msg, ok, err := q.channel.Get(queueName, true)
	return nil, errors.New("not implemented")
}

// Process 启动消费者处理队列消息
func (q *RabbitMQQueue[T]) Process(concurrency int, handler func(*Message[T]) error) error {
	// TODO: 启动 concurrency 个 goroutine 消费消息
	// deliveries, err := q.channel.Consume(queueName, ...)
	// 每个 goroutine 从 deliveries channel 读取并处理
	return errors.New("not implemented")
}

// Depth 获取当前队列深度
func (q *RabbitMQQueue[T]) Depth() int {
	// TODO: 使用 Channel.QueueInspect() 获取队列深度
	return 0
}

// Close 关闭队列连接
func (q *RabbitMQQueue[T]) Close() error {
	// TODO: 关闭 channel 和 connection
	return nil
}
