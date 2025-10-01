package queue

import (
	"errors"
	"time"
)

// NSQQueue NSQž°
type NSQQueue[T any] struct {
	cfg *Config
}

// NewNSQQueue ú NSQ 
func NewNSQQueue[T any](cfg *Config) (*NSQQueue[T], error) {
	if cfg.Address == "" {
		return nil, ErrInvalidConfig
	}
	return &NSQQueue[T]{cfg: cfg}, errors.New("NSQ queue not yet implemented")
}

func (q *NSQQueue[T]) Push(payload T) error {
	return errors.New("not implemented")
}

func (q *NSQQueue[T]) PushWithDelay(payload T, delay time.Duration) error {
	return errors.New("not implemented")
}

func (q *NSQQueue[T]) PushWithRetryPolicy(payload T, policy *RetryPolicy) error {
	return errors.New("not implemented")
}

func (q *NSQQueue[T]) BatchPush(payloads []T) error {
	return errors.New("not implemented")
}

func (q *NSQQueue[T]) Pop() (*Message[T], error) {
	return nil, errors.New("not implemented")
}

func (q *NSQQueue[T]) Process(concurrency int, handler func(*Message[T]) error) error {
	return errors.New("not implemented")
}

func (q *NSQQueue[T]) Depth() int {
	return 0
}

func (q *NSQQueue[T]) Close() error {
	return nil
}
