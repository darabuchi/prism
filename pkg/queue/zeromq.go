package queue

import "errors"

// ZeroMQQueue ZeroMQ队列实现
type ZeroMQQueue[T any] struct {
	cfg *Config
	// TODO: 添加 ZeroMQ sockets
	// pushSocket *zmq4.Socket
	// pullSocket *zmq4.Socket
}

// NewZeroMQQueue 创建 ZeroMQ 队列
func NewZeroMQQueue[T any](cfg *Config) (*ZeroMQQueue[T], error) {
	if cfg.Address == "" {
		return nil, ErrInvalidConfig
	}

	// TODO: 初始化 ZeroMQ sockets
	// 依赖: github.com/pebbe/zmq4

	return &ZeroMQQueue[T]{
		cfg: cfg,
	}, errors.New("ZeroMQ queue not yet implemented")
}

// Push 推送单个消息到队列
func (q *ZeroMQQueue[T]) Push(payload T) error {
	// TODO: 使用 PUSH socket 发送消息
	// _, err := q.pushSocket.SendBytes(serialize(payload), 0)
	return errors.New("not implemented")
}

// BatchPush 批量推送消息到队列
func (q *ZeroMQQueue[T]) BatchPush(payloads []T) error {
	// TODO: 批量推送消息
	for _, payload := range payloads {
		if err := q.Push(payload); err != nil {
			return err
		}
	}
	return nil
}

// Pop 从队列中弹出一个消息
func (q *ZeroMQQueue[T]) Pop() (*Message[T], error) {
	// TODO: 使用 PULL socket 非阻塞接收
	// bytes, err := q.pullSocket.RecvBytes(zmq4.DONTWAIT)
	return nil, errors.New("not implemented")
}

// Process 启动消费者处理队列消息
func (q *ZeroMQQueue[T]) Process(concurrency int, handler func(*Message[T]) error) error {
	// TODO: 启动 concurrency 个 goroutine，每个使用 PULL socket 接收消息
	return errors.New("not implemented")
}

// Depth 获取当前队列深度
func (q *ZeroMQQueue[T]) Depth() int {
	// ZeroMQ 不支持查询队列深度
	return 0
}

// Close 关闭队列连接
func (q *ZeroMQQueue[T]) Close() error {
	// TODO: 关闭所有 sockets
	return nil
}
