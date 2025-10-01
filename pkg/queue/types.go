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
	ID           string        // 消息唯一标识
	Payload      T             // 消息内容（泛型）
	Timestamp    time.Time     // 消息创建时间戳
	DelayUntil   time.Time     // 延时执行时间（为零值表示立即执行）
	RetryCount   int           // 当前重试次数
	RetryPolicy  *RetryPolicy  // 重试策略（为 nil 使用队列默认策略）
	LastError    error         // 上次处理错误
	NextRetryAt  time.Time     // 下次重试时间
}

// IsReady 判断消息是否准备好执行
func (m *Message[T]) IsReady() bool {
	now := time.Now()

	// 检查延时执行时间
	if !m.DelayUntil.IsZero() && now.Before(m.DelayUntil) {
		return false
	}

	// 检查重试时间
	if !m.NextRetryAt.IsZero() && now.Before(m.NextRetryAt) {
		return false
	}

	return true
}

// ShouldRetry 判断是否应该重试
func (m *Message[T]) ShouldRetry(defaultPolicy *RetryPolicy) bool {
	policy := m.RetryPolicy
	if policy == nil {
		policy = defaultPolicy
	}

	if policy == nil {
		return false
	}

	return policy.ShouldRetry(m.RetryCount)
}

// CalculateNextRetry 计算下次重试时间
func (m *Message[T]) CalculateNextRetry(defaultPolicy *RetryPolicy) {
	policy := m.RetryPolicy
	if policy == nil {
		policy = defaultPolicy
	}

	if policy == nil {
		return
	}

	delay := policy.CalculateDelay(m.RetryCount)
	m.NextRetryAt = time.Now().Add(delay)
}

// Config 队列配置
type Config struct {
	// 队列类型
	Type Type `json:"type" yaml:"type" validate:"required,oneof=memory redis nsq kafka rabbitmq zeromq" default:"memory"`

	// 连接地址
	Address string `json:"address" yaml:"address" validate:"required_unless=Type memory" default:""`

	// 操作超时（秒）
	Timeout int `json:"timeout" yaml:"timeout" validate:"min=1,max=300" default:"30"`

	// 队列最大深度（0 表示无限制）
	MaxDepth int `json:"max_depth" yaml:"max_depth" validate:"min=0" default:"10000"`

	// 默认重试策略
	RetryPolicy *RetryPolicy `json:"retry_policy,omitempty" yaml:"retry_policy,omitempty"`

	// 是否启用延时队列
	EnableDelayedQueue bool `json:"enable_delayed_queue" yaml:"enable_delayed_queue" default:"false"`

	// 延时队列检查间隔（毫秒）
	DelayCheckInterval int `json:"delay_check_interval" yaml:"delay_check_interval" validate:"min=100" default:"1000"`

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
	if c.Timeout == 0 {
		c.Timeout = 30
	}
	if c.MaxDepth == 0 {
		c.MaxDepth = 10000
	}
	if c.RetryPolicy == nil {
		c.RetryPolicy = DefaultRetryPolicy()
	} else {
		c.RetryPolicy.ApplyDefaults()
	}
	if c.DelayCheckInterval == 0 {
		c.DelayCheckInterval = 1000
	}
	if c.Options == nil {
		c.Options = make(map[string]interface{})
	}
}
