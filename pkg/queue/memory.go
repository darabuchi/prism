package queue

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

// MemoryQueue implements an in-memory message queue
type MemoryQueue struct {
	mu         sync.RWMutex
	topics     map[string]*topic
	maxRetry   int
	timeout    time.Duration
	closed     bool
}

type topic struct {
	subscribers []chan *Message
	mu          sync.RWMutex
}

// NewMemoryQueue creates a new in-memory queue
func NewMemoryQueue(cfg *Config) (*MemoryQueue, error) {
	return &MemoryQueue{
		topics:   make(map[string]*topic),
		maxRetry: cfg.MaxRetry,
		timeout:  cfg.Timeout,
	}, nil
}

// Publish sends a message to the queue
func (q *MemoryQueue) Publish(ctx context.Context, topicName string, payload []byte) error {
	return q.PublishWithMetadata(ctx, topicName, payload, nil)
}

// PublishWithMetadata sends a message with metadata
func (q *MemoryQueue) PublishWithMetadata(ctx context.Context, topicName string, payload []byte, metadata map[string]string) error {
	q.mu.RLock()
	if q.closed {
		q.mu.RUnlock()
		return ErrQueueClosed
	}
	q.mu.RUnlock()

	msg := &Message{
		ID:        uuid.New().String(),
		Topic:     topicName,
		Payload:   payload,
		Metadata:  metadata,
		Timestamp: time.Now(),
		RetryCount: 0,
	}

	q.mu.Lock()
	t, ok := q.topics[topicName]
	if !ok {
		t = &topic{
			subscribers: make([]chan *Message, 0),
		}
		q.topics[topicName] = t
	}
	q.mu.Unlock()

	// Send to all subscribers
	t.mu.RLock()
	defer t.mu.RUnlock()

	for _, sub := range t.subscribers {
		select {
		case sub <- msg:
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(q.timeout):
			return ErrTimeout
		}
	}

	return nil
}

// Subscribe subscribes to a topic and returns a channel for receiving messages
func (q *MemoryQueue) Subscribe(ctx context.Context, topicName string) (<-chan *Message, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return nil, ErrQueueClosed
	}

	t, ok := q.topics[topicName]
	if !ok {
		t = &topic{
			subscribers: make([]chan *Message, 0),
		}
		q.topics[topicName] = t
	}

	msgChan := make(chan *Message, 100)

	t.mu.Lock()
	t.subscribers = append(t.subscribers, msgChan)
	t.mu.Unlock()

	return msgChan, nil
}

// Consume consumes messages from a topic with a handler function
func (q *MemoryQueue) Consume(ctx context.Context, topicName string, handler func(*Message) error) error {
	msgChan, err := q.Subscribe(ctx, topicName)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case msg := <-msgChan:
				if err := handler(msg); err != nil {
					// Retry logic
					if msg.RetryCount < q.maxRetry {
						msg.RetryCount++
						// Re-publish for retry
						_ = q.PublishWithMetadata(context.Background(), topicName, msg.Payload, msg.Metadata)
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

// Ack acknowledges a message (no-op for memory queue)
func (q *MemoryQueue) Ack(ctx context.Context, msg *Message) error {
	return nil
}

// Nack negatively acknowledges a message (triggers retry)
func (q *MemoryQueue) Nack(ctx context.Context, msg *Message) error {
	if msg.RetryCount >= q.maxRetry {
		return ErrMaxRetriesExceeded
	}

	msg.RetryCount++
	return q.PublishWithMetadata(ctx, msg.Topic, msg.Payload, msg.Metadata)
}

// Close closes the queue
func (q *MemoryQueue) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return nil
	}

	q.closed = true

	// Close all subscriber channels
	for _, t := range q.topics {
		t.mu.Lock()
		for _, sub := range t.subscribers {
			close(sub)
		}
		t.subscribers = nil
		t.mu.Unlock()
	}

	return nil
}
