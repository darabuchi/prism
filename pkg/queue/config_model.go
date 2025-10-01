package queue

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// QueueConfig 队列配置数据库模型（无 ORM 依赖）
type QueueConfig struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" validate:"required,min=1,max=100"`
	Type      Type      `json:"type" db:"type" validate:"required,oneof=memory redis nsq kafka rabbitmq zeromq"`
	Address   string    `json:"address" db:"address"`
	Timeout   int       `json:"timeout" db:"timeout" validate:"min=1,max=300"`
	MaxDepth  int       `json:"max_depth" db:"max_depth" validate:"min=0"`
	Enabled   bool      `json:"enabled" db:"enabled"`

	// 重试策略（JSON 存储）
	RetryPolicyJSON string `json:"-" db:"retry_policy"`

	// 延时队列配置
	EnableDelayedQueue bool `json:"enable_delayed_queue" db:"enable_delayed_queue"`
	DelayCheckInterval int  `json:"delay_check_interval" db:"delay_check_interval" validate:"min=100"`

	// 扩展配置（JSON 存储）
	OptionsJSON string `json:"-" db:"options"`

	// 描述和元信息
	Description string    `json:"description" db:"description" validate:"max=500"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// 运行时字段（不存储到数据库）- 并发安全
	mu          sync.RWMutex
	retryPolicy *RetryPolicy
	options     map[string]interface{}
}

// TableName 指定表名
func (QueueConfig) TableName() string {
	return "queue_configs"
}

// ToConfig 转换为运行时配置
func (qc *QueueConfig) ToConfig() (*Config, error) {
	qc.mu.RLock()
	defer qc.mu.RUnlock()

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

// GetRetryPolicy 获取重试策略（并发安全）
func (qc *QueueConfig) GetRetryPolicy() (*RetryPolicy, error) {
	qc.mu.RLock()
	if qc.retryPolicy != nil {
		policy := qc.retryPolicy
		qc.mu.RUnlock()
		return policy, nil
	}
	qc.mu.RUnlock()

	// 升级为写锁
	qc.mu.Lock()
	defer qc.mu.Unlock()

	// 双重检查
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

// SetRetryPolicy 设置重试策略（并发安全）
func (qc *QueueConfig) SetRetryPolicy(policy *RetryPolicy) error {
	qc.mu.Lock()
	defer qc.mu.Unlock()

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

// GetOptions 获取扩展配置（并发安全）
func (qc *QueueConfig) GetOptions() (map[string]interface{}, error) {
	qc.mu.RLock()
	if qc.options != nil {
		// 返回副本以避免外部修改
		optionsCopy := make(map[string]interface{}, len(qc.options))
		for k, v := range qc.options {
			optionsCopy[k] = v
		}
		qc.mu.RUnlock()
		return optionsCopy, nil
	}
	qc.mu.RUnlock()

	// 升级为写锁
	qc.mu.Lock()
	defer qc.mu.Unlock()

	// 双重检查
	if qc.options != nil {
		optionsCopy := make(map[string]interface{}, len(qc.options))
		for k, v := range qc.options {
			optionsCopy[k] = v
		}
		return optionsCopy, nil
	}

	if qc.OptionsJSON == "" {
		return nil, nil
	}

	var options map[string]interface{}
	if err := json.Unmarshal([]byte(qc.OptionsJSON), &options); err != nil {
		return nil, err
	}

	qc.options = options

	// 返回副本
	optionsCopy := make(map[string]interface{}, len(options))
	for k, v := range options {
		optionsCopy[k] = v
	}
	return optionsCopy, nil
}

// SetOptions 设置扩展配置（并发安全）
func (qc *QueueConfig) SetOptions(options map[string]interface{}) error {
	qc.mu.Lock()
	defer qc.mu.Unlock()

	// 存储副本
	if options != nil {
		qc.options = make(map[string]interface{}, len(options))
		for k, v := range options {
			qc.options[k] = v
		}
	} else {
		qc.options = nil
	}

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

// AfterLoad 加载后钩子 - 预解析 JSON 字段
func (qc *QueueConfig) AfterLoad() error {
	qc.mu.Lock()
	defer qc.mu.Unlock()

	// 预解析重试策略
	if qc.RetryPolicyJSON != "" && qc.retryPolicy == nil {
		var policy RetryPolicy
		if err := json.Unmarshal([]byte(qc.RetryPolicyJSON), &policy); err != nil {
			return err
		}
		qc.retryPolicy = &policy
	}

	// 预解析扩展配置
	if qc.OptionsJSON != "" && qc.options == nil {
		var options map[string]interface{}
		if err := json.Unmarshal([]byte(qc.OptionsJSON), &options); err != nil {
			return err
		}
		qc.options = options
	}

	return nil
}

// BeforeSave 保存前钩子 - 序列化缓存字段
func (qc *QueueConfig) BeforeSave() error {
	qc.mu.RLock()
	defer qc.mu.RUnlock()

	// 序列化重试策略
	if qc.retryPolicy != nil {
		policyJSON, err := json.Marshal(qc.retryPolicy)
		if err != nil {
			return err
		}
		qc.RetryPolicyJSON = string(policyJSON)
	}

	// 序列化扩展配置
	if qc.options != nil {
		optionsJSON, err := json.Marshal(qc.options)
		if err != nil {
			return err
		}
		qc.OptionsJSON = string(optionsJSON)
	}

	return nil
}

// Validate 验证配置
func (qc *QueueConfig) Validate() error {
	// 使用 validator 进行结构化验证
	if err := ValidateStruct(qc); err != nil {
		return err
	}

	// 额外的业务逻辑验证
	if qc.Name == "" {
		return fmt.Errorf("queue name is required")
	}

	if !qc.Type.IsValid() {
		return fmt.Errorf("invalid queue type: %s", qc.Type)
	}

	if qc.Type != TypeMemory && qc.Address == "" {
		return fmt.Errorf("address is required for queue type: %s", qc.Type)
	}

	return nil
}
