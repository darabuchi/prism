package queue

import "time"

// Type 队列类型枚举
type Type string

const (
	TypeMemory   Type = "memory"   // 内存队列
	TypeRedis    Type = "redis"    // Redis队列
	TypeNSQ      Type = "nsq"      // NSQ队列
	TypeKafka    Type = "kafka"    // Kafka队列
	TypeRabbitMQ Type = "rabbitmq" // RabbitMQ队列
	TypeZeroMQ   Type = "zeromq"   // ZeroMQ队列
)

// IsValid 检查队列类型是否有效
func (t Type) IsValid() bool {
	switch t {
	case TypeMemory, TypeRedis, TypeNSQ, TypeKafka, TypeRabbitMQ, TypeZeroMQ:
		return true
	default:
		return false
	}
}

// String 返回队列类型字符串
func (t Type) String() string {
	return string(t)
}

// Message 队列消息（泛型）
type Message[T any] struct {
	ID        string    // 消息唯一标识
	Payload   T         // 消息内容（泛型）
	Timestamp time.Time // 消息时间戳
	RetryCount int      // 重试次数
}

// Config 队列配置
type Config struct {
	// 队列类型
	Type Type `json:"type" yaml:"type" validate:"required,oneof=memory redis nsq kafka rabbitmq zeromq" default:"memory"`

	// 连接地址
	Address string `json:"address" yaml:"address" validate:"required_unless=Type memory" default:""`

	// 最大重试次数
	MaxRetry int `json:"max_retry" yaml:"max_retry" validate:"min=0,max=100" default:"3"`

	// 操作超时（秒）
	Timeout int `json:"timeout" yaml:"timeout" validate:"min=1,max=300" default:"30"`

	// 队列最大深度（0 表示无限制）
	MaxDepth int `json:"max_depth" yaml:"max_depth" validate:"min=0" default:"10000"`

	// 类型特定配置
	Options map[string]interface{} `json:"options,omitempty" yaml:"options,omitempty"`
}

// Validate 验证配置
func (c *Config) Validate() error {
	if !c.Type.IsValid() {
		return ErrInvalidConfig
	}

	if c.Type != TypeMemory && c.Address == "" {
		return ErrInvalidConfig
	}

	if c.MaxRetry < 0 {
		return ErrInvalidConfig
	}

	if c.Timeout < 1 {
		return ErrInvalidConfig
	}

	return nil
}

// ApplyDefaults 应用默认值
func (c *Config) ApplyDefaults() {
	if c.Type == "" {
		c.Type = TypeMemory
	}
	if c.MaxRetry == 0 {
		c.MaxRetry = 3
	}
	if c.Timeout == 0 {
		c.Timeout = 30
	}
	if c.MaxDepth == 0 {
		c.MaxDepth = 10000
	}
	if c.Options == nil {
		c.Options = make(map[string]interface{})
	}
}
