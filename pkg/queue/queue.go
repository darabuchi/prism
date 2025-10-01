package queue

import (
	"context"
	"time"
)

// Message represents a queue message
type Message struct {
	ID        string                 // Message ID
	Topic     string                 // Topic/Queue name
	Payload   []byte                 // Message payload
	Metadata  map[string]string      // Optional metadata
	Timestamp time.Time              // Message timestamp
	RetryCount int                   // Retry count
}

// Queue defines the message queue interface
type Queue interface {
	// Publish sends a message to the queue
	Publish(ctx context.Context, topic string, payload []byte) error

	// PublishWithMetadata sends a message with metadata
	PublishWithMetadata(ctx context.Context, topic string, payload []byte, metadata map[string]string) error

	// Subscribe subscribes to a topic and returns a channel for receiving messages
	Subscribe(ctx context.Context, topic string) (<-chan *Message, error)

	// Consume consumes messages from a topic with a handler function
	Consume(ctx context.Context, topic string, handler func(*Message) error) error

	// Ack acknowledges a message
	Ack(ctx context.Context, msg *Message) error

	// Nack negatively acknowledges a message (for redelivery)
	Nack(ctx context.Context, msg *Message) error

	// Close closes the queue connection
	Close() error
}

// Config represents queue configuration
type Config struct {
	Type     string                 // Queue type: memory, redis, nsq, kafka, rabbitmq, zeromq
	Address  string                 // Connection address(es)
	Options  map[string]interface{} // Type-specific options
	MaxRetry int                    // Maximum retry count
	Timeout  time.Duration          // Operation timeout
}

// Option is a functional option for queue configuration
type Option func(*Config)

// WithAddress sets the queue address
func WithAddress(addr string) Option {
	return func(c *Config) {
		c.Address = addr
	}
}

// WithMaxRetry sets the maximum retry count
func WithMaxRetry(maxRetry int) Option {
	return func(c *Config) {
		c.MaxRetry = maxRetry
	}
}

// WithTimeout sets the operation timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithOptions sets type-specific options
func WithOptions(opts map[string]interface{}) Option {
	return func(c *Config) {
		c.Options = opts
	}
}

// New creates a new queue instance based on the configuration
func New(queueType string, opts ...Option) (Queue, error) {
	cfg := &Config{
		Type:     queueType,
		MaxRetry: 3,
		Timeout:  30 * time.Second,
		Options:  make(map[string]interface{}),
	}

	for _, opt := range opts {
		opt(cfg)
	}

	switch queueType {
	case "memory":
		return NewMemoryQueue(cfg)
	case "redis":
		return NewRedisQueue(cfg)
	case "nsq":
		return NewNSQQueue(cfg)
	case "kafka":
		return NewKafkaQueue(cfg)
	case "rabbitmq":
		return NewRabbitMQQueue(cfg)
	case "zeromq":
		return NewZeroMQQueue(cfg)
	default:
		return nil, ErrUnsupportedQueueType
	}
}
