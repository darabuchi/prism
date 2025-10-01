package queue

import (
	"fmt"
	"time"
)

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

	// 取 DelayUntil 和 NextRetryAt 的最大值
	maxTime := m.DelayUntil
	if m.NextRetryAt.After(maxTime) {
		maxTime = m.NextRetryAt
	}

	// 如果最大时间为零值或已过期，则准备好
	return maxTime.IsZero() || !now.Before(maxTime)
}

// Age 返回消息的年龄（从创建到现在的时长）
func (m *Message[T]) Age() time.Duration {
	return time.Since(m.Timestamp)
}

// HasFailed 判断消息是否有失败记录
func (m *Message[T]) HasFailed() bool {
	return m.LastError != nil
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
	Type Type `json:"type" yaml:"type" validate:"required,oneof=memory redis nsq kafka rabbitmq zeromq"`

	// 连接地址
	Address string `json:"address" yaml:"address"`

	// 操作超时（秒）
	Timeout int `json:"timeout" yaml:"timeout" validate:"min=1,max=300"`

	// 队列最大深度（0 表示无限制）
	MaxDepth int `json:"max_depth" yaml:"max_depth" validate:"min=0"`

	// 默认重试策略
	RetryPolicy *RetryPolicy `json:"retry_policy,omitempty" yaml:"retry_policy,omitempty"`

	// 是否启用延时队列
	EnableDelayedQueue bool `json:"enable_delayed_queue" yaml:"enable_delayed_queue"`

	// 延时队列检查间隔（毫秒）
	DelayCheckInterval int `json:"delay_check_interval" yaml:"delay_check_interval" validate:"min=100"`

	// 类型特定配置
	Options map[string]interface{} `json:"options,omitempty" yaml:"options,omitempty"`
}

// Validate 验证配置
func (c *Config) Validate() error {
	// 使用 validator 进行结构化验证
	if err := ValidateStruct(c); err != nil {
		return err
	}

	// 额外的业务逻辑验证
	if !c.Type.IsValid() {
		return fmt.Errorf("invalid queue type: %s", c.Type)
	}

	if c.Type != TypeMemory && c.Address == "" {
		return fmt.Errorf("address is required for queue type: %s", c.Type)
	}

	// 验证重试策略
	if c.RetryPolicy != nil {
		if err := c.RetryPolicy.Validate(); err != nil {
			return fmt.Errorf("invalid retry policy: %w", err)
		}
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
