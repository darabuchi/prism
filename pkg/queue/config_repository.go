package queue

import "context"

// ConfigRepository 队列配置仓储接口
type ConfigRepository interface {
	// GetByName 根据名称获取配置
	GetByName(ctx context.Context, name string) (*QueueConfig, error)

	// GetEnabled 获取所有启用的配置
	GetEnabled(ctx context.Context) ([]*QueueConfig, error)

	// GetByType 根据类型获取配置列表
	GetByType(ctx context.Context, queueType Type) ([]*QueueConfig, error)

	// Create 创建配置
	Create(ctx context.Context, config *QueueConfig) error

	// Update 更新配置
	Update(ctx context.Context, config *QueueConfig) error

	// Delete 删除配置
	Delete(ctx context.Context, name string) error

	// Enable 启用配置
	Enable(ctx context.Context, name string) error

	// Disable 禁用配置
	Disable(ctx context.Context, name string) error
}

// ConfigLoader 配置加载器接口
type ConfigLoader interface {
	// LoadConfig 加载队列配置
	LoadConfig(ctx context.Context, name string) (*Config, error)

	// LoadAllConfigs 加载所有启用的队列配置
	LoadAllConfigs(ctx context.Context) (map[string]*Config, error)

	// ReloadConfig 重新加载配置
	ReloadConfig(ctx context.Context, name string) (*Config, error)
}

// DBConfigLoader 基于数据库的配置加载器
type DBConfigLoader struct {
	repo ConfigRepository
}

// NewDBConfigLoader 创建数据库配置加载器
func NewDBConfigLoader(repo ConfigRepository) *DBConfigLoader {
	return &DBConfigLoader{repo: repo}
}

// LoadConfig 加载指定名称的队列配置
func (l *DBConfigLoader) LoadConfig(ctx context.Context, name string) (*Config, error) {
	qc, err := l.repo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	if !qc.Enabled {
		return nil, ErrQueueDisabled
	}

	return qc.ToConfig()
}

// LoadAllConfigs 加载所有启用的队列配置
func (l *DBConfigLoader) LoadAllConfigs(ctx context.Context) (map[string]*Config, error) {
	configs, err := l.repo.GetEnabled(ctx)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*Config, len(configs))
	for _, qc := range configs {
		cfg, err := qc.ToConfig()
		if err != nil {
			return nil, err
		}
		result[qc.Name] = cfg
	}

	return result, nil
}

// ReloadConfig 重新加载配置
func (l *DBConfigLoader) ReloadConfig(ctx context.Context, name string) (*Config, error) {
	return l.LoadConfig(ctx, name)
}
