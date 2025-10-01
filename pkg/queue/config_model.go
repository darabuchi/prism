package queue

import (
	"encoding/json"
	"time"
)

// QueueConfig 队列配置数据库模型
type QueueConfig struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;uniqueIndex;type:varchar(100);not null" json:"name" validate:"required,min=1,max=100"`
	Type      Type      `gorm:"column:type;type:varchar(20);not null" json:"type" validate:"required,oneof=memory redis nsq kafka rabbitmq zeromq"`
	Address   string    `gorm:"column:address;type:varchar(500)" json:"address" validate:"required_unless=Type memory"`
	Timeout   int       `gorm:"column:timeout;default:30" json:"timeout" validate:"min=1,max=300"`
	MaxDepth  int       `gorm:"column:max_depth;default:10000" json:"max_depth" validate:"min=0"`
	Enabled   bool      `gorm:"column:enabled;default:true" json:"enabled"`

	// 重试策略（JSON 存储）
	RetryPolicyJSON string `gorm:"column:retry_policy;type:text" json:"-"`

	// 延时队列配置
	EnableDelayedQueue bool `gorm:"column:enable_delayed_queue;default:false" json:"enable_delayed_queue"`
	DelayCheckInterval int  `gorm:"column:delay_check_interval;default:1000" json:"delay_check_interval" validate:"min=100"`

	// 扩展配置（JSON 存储）
	OptionsJSON string `gorm:"column:options;type:text" json:"-"`

	// 描述和元信息
	Description string    `gorm:"column:description;type:varchar(500)" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// 运行时字段（不存储到数据库）
	retryPolicy *RetryPolicy           `gorm:"-" json:"retry_policy,omitempty"`
	options     map[string]interface{} `gorm:"-" json:"options,omitempty"`
}

// TableName 指定表名
func (QueueConfig) TableName() string {
	return "queue_configs"
}

// ToConfig 转换为运行时配置
func (qc *QueueConfig) ToConfig() (*Config, error) {
	cfg := &Config{
		Type:               qc.Type,
		Address:            qc.Address,
		Timeout:            qc.Timeout,
		MaxDepth:           qc.MaxDepth,
		EnableDelayedQueue: qc.EnableDelayedQueue,
		DelayCheckInterval: qc.DelayCheckInterval,
	}

	// 解析重试策略
	if qc.RetryPolicyJSON != "" {
		var policy RetryPolicy
		if err := json.Unmarshal([]byte(qc.RetryPolicyJSON), &policy); err != nil {
			return nil, err
		}
		cfg.RetryPolicy = &policy
	}

	// 解析扩展配置
	if qc.OptionsJSON != "" {
		var options map[string]interface{}
		if err := json.Unmarshal([]byte(qc.OptionsJSON), &options); err != nil {
			return nil, err
		}
		cfg.Options = options
	}

	return cfg, nil
}

// FromConfig 从运行时配置创建数据库模型
func FromConfig(name string, cfg *Config) (*QueueConfig, error) {
	qc := &QueueConfig{
		Name:               name,
		Type:               cfg.Type,
		Address:            cfg.Address,
		Timeout:            cfg.Timeout,
		MaxDepth:           cfg.MaxDepth,
		EnableDelayedQueue: cfg.EnableDelayedQueue,
		DelayCheckInterval: cfg.DelayCheckInterval,
		Enabled:            true,
	}

	// 序列化重试策略
	if cfg.RetryPolicy != nil {
		policyJSON, err := json.Marshal(cfg.RetryPolicy)
		if err != nil {
			return nil, err
		}
		qc.RetryPolicyJSON = string(policyJSON)
	}

	// 序列化扩展配置
	if cfg.Options != nil {
		optionsJSON, err := json.Marshal(cfg.Options)
		if err != nil {
			return nil, err
		}
		qc.OptionsJSON = string(optionsJSON)
	}

	return qc, nil
}

// GetRetryPolicy 获取重试策略
func (qc *QueueConfig) GetRetryPolicy() (*RetryPolicy, error) {
	if qc.retryPolicy != nil {
		return qc.retryPolicy, nil
	}

	if qc.RetryPolicyJSON == "" {
		return nil, nil
	}

	var policy RetryPolicy
	if err := json.Unmarshal([]byte(qc.RetryPolicyJSON), &policy); err != nil {
		return nil, err
	}

	qc.retryPolicy = &policy
	return qc.retryPolicy, nil
}

// SetRetryPolicy 设置重试策略
func (qc *QueueConfig) SetRetryPolicy(policy *RetryPolicy) error {
	qc.retryPolicy = policy

	if policy == nil {
		qc.RetryPolicyJSON = ""
		return nil
	}

	policyJSON, err := json.Marshal(policy)
	if err != nil {
		return err
	}

	qc.RetryPolicyJSON = string(policyJSON)
	return nil
}

// GetOptions 获取扩展配置
func (qc *QueueConfig) GetOptions() (map[string]interface{}, error) {
	if qc.options != nil {
		return qc.options, nil
	}

	if qc.OptionsJSON == "" {
		return nil, nil
	}

	var options map[string]interface{}
	if err := json.Unmarshal([]byte(qc.OptionsJSON), &options); err != nil {
		return nil, err
	}

	qc.options = options
	return qc.options, nil
}

// SetOptions 设置扩展配置
func (qc *QueueConfig) SetOptions(options map[string]interface{}) error {
	qc.options = options

	if options == nil {
		qc.OptionsJSON = ""
		return nil
	}

	optionsJSON, err := json.Marshal(options)
	if err != nil {
		return err
	}

	qc.OptionsJSON = string(optionsJSON)
	return nil
}
