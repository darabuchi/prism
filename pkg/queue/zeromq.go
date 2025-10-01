package queue

import (
	"errors"
	"time"
)

// ZeroMQQueue ZeroMQ��
type ZeroMQQueue[T any] struct {
	cfg *Config
}

// NewZeroMQQueue � ZeroMQ 
func NewZeroMQQueue[T any](cfg *Config) (*ZeroMQQueue[T], error) {
	if cfg.Address == "" {
		return nil, ErrInvalidConfig
	}
	return &ZeroMQQueue[T]{cfg: cfg}, errors.New("ZeroMQ queue not yet implemented")
}

func (q *ZeroMQQueue[T]) Push(payload T) error {
	return errors.New("not implemented")
}

func (q *ZeroMQQueue[T]) PushWithDelay(payload T, delay time.Duration) error {
	return errors.New("not implemented")
}

func (q *ZeroMQQueue[T]) PushWithRetryPolicy(payload T, policy *RetryPolicy) error {
	return errors.New("not implemented")
}

func (q *ZeroMQQueue[T]) BatchPush(payloads []T) error {
	return errors.New("not implemented")
}

func (q *ZeroMQQueue[T]) Pop() (*Message[T], error) {
	return nil, errors.New("not implemented")
}

func (q *ZeroMQQueue[T]) Process(concurrency int, handler func(*Message[T]) (*RetryInfo, error)) error {
	return errors.New("not implemented")
}

func (q *ZeroMQQueue[T]) Depth() int {
	return 0
}

func (q *ZeroMQQueue[T]) Close() error {
	return nil
}
