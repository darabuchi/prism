package queue

import (
	"errors"
	"time"
)

// KafkaQueue Kafka��
type KafkaQueue[T any] struct {
	cfg *Config
}

// NewKafkaQueue � Kafka 
func NewKafkaQueue[T any](cfg *Config) (*KafkaQueue[T], error) {
	if cfg.Address == "" {
		return nil, ErrInvalidConfig
	}
	return &KafkaQueue[T]{cfg: cfg}, errors.New("Kafka queue not yet implemented")
}

func (q *KafkaQueue[T]) Push(payload T) error {
	return errors.New("not implemented")
}

func (q *KafkaQueue[T]) PushWithDelay(payload T, delay time.Duration) error {
	return errors.New("not implemented")
}

func (q *KafkaQueue[T]) PushWithRetryPolicy(payload T, policy *RetryPolicy) error {
	return errors.New("not implemented")
}

func (q *KafkaQueue[T]) BatchPush(payloads []T) error {
	return errors.New("not implemented")
}

func (q *KafkaQueue[T]) Pop() (*Message[T], error) {
	return nil, errors.New("not implemented")
}

func (q *KafkaQueue[T]) Process(concurrency int, handler func(*Message[T]) (*RetryInfo, error)) error {
	return errors.New("not implemented")
}

func (q *KafkaQueue[T]) Depth() int {
	return 0
}

func (q *KafkaQueue[T]) Close() error {
	return nil
}
