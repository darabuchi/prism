package queue

import (
	"errors"
	"time"
)

// RabbitMQQueue RabbitMQ��
type RabbitMQQueue[T any] struct {
	cfg *Config
}

// NewRabbitMQQueue � RabbitMQ 
func NewRabbitMQQueue[T any](cfg *Config) (*RabbitMQQueue[T], error) {
	if cfg.Address == "" {
		return nil, ErrInvalidConfig
	}
	return &RabbitMQQueue[T]{cfg: cfg}, errors.New("RabbitMQ queue not yet implemented")
}

func (q *RabbitMQQueue[T]) Push(payload T) error {
	return errors.New("not implemented")
}

func (q *RabbitMQQueue[T]) PushWithDelay(payload T, delay time.Duration) error {
	return errors.New("not implemented")
}

func (q *RabbitMQQueue[T]) PushWithRetryPolicy(payload T, policy *RetryPolicy) error {
	return errors.New("not implemented")
}

func (q *RabbitMQQueue[T]) BatchPush(payloads []T) error {
	return errors.New("not implemented")
}

func (q *RabbitMQQueue[T]) Pop() (*Message[T], error) {
	return nil, errors.New("not implemented")
}

func (q *RabbitMQQueue[T]) Process(concurrency int, handler func(*Message[T]) *HandlerResult) error {
	return errors.New("not implemented")
}

func (q *RabbitMQQueue[T]) Depth() int {
	return 0
}

func (q *RabbitMQQueue[T]) Close() error {
	return nil
}
