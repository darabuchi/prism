package queue

import "errors"

// NSQQueue NSQ队列实现
type NSQQueue[T any] struct {
	cfg *Config
	// TODO: 添加 NSQ producer 和 consumer
	// producer *nsq.Producer
	// consumers []*nsq.Consumer
}

// NewNSQQueue 创建 NSQ 队列
func NewNSQQueue[T any](cfg *Config) (*NSQQueue[T], error) {
	if cfg.Address == "" {
		return nil, ErrInvalidConfig
	}

	// TODO: 初始化 NSQ producer
	// 依赖: github.com/nsqio/go-nsq

	return &NSQQueue[T]{
		cfg: cfg,
	}, errors.New("NSQ queue not yet implemented")
}

// Push 推送单个消息到队列
func (q *NSQQueue[T]) Push(payload T) error {
	// TODO: 使用 Producer.Publish() 推送消息
	// body := serialize(payload)
	// return q.producer.Publish(topic, body)
	return errors.New("not implemented")
}

// BatchPush 批量推送消息到队列
func (q *NSQQueue[T]) BatchPush(payloads []T) error {
	// TODO: 使用 Producer.MultiPublish() 批量推送
	for _, payload := range payloads {
		if err := q.Push(payload); err != nil {
			return err
		}
	}
	return nil
}

// Pop 从队列中弹出一个消息
func (q *NSQQueue[T]) Pop() (*Message[T], error) {
	// NSQ 不支持 Pop 操作，必须使用 Process 消费
	return nil, errors.New("NSQ does not support Pop operation, use Process instead")
}

// Process 启动消费者处理队列消息
func (q *NSQQueue[T]) Process(concurrency int, handler func(*Message[T]) error) error {
	// TODO: 创建 Consumer 并设置 concurrency
	// consumer.AddConcurrentHandlers(nsq.HandlerFunc(func(message *nsq.Message) error {
	//     msg := deserialize[T](message.Body)
	//     return handler(msg)
	// }), concurrency)
	return errors.New("not implemented")
}

// Depth 获取当前队列深度
func (q *NSQQueue[T]) Depth() int {
	// TODO: 通过 nsqd HTTP API 获取队列深度
	return 0
}

// Close 关闭队列连接
func (q *NSQQueue[T]) Close() error {
	// TODO: 停止 producer 和所有 consumers
	return nil
}
