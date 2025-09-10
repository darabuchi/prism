package storage

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	// 使用纯 Go 实现的 SQLite 驱动
	_ "modernc.org/sqlite"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver string `mapstructure:"driver" yaml:"driver"`
	DSN    string `mapstructure:"dsn" yaml:"dsn"`
	Debug  bool   `mapstructure:"debug" yaml:"debug"`

	// 连接池配置
	MaxOpenConns    int           `mapstructure:"max_open_conns" yaml:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime"`
}

// Database 数据库管理器
type Database struct {
	*gorm.DB
	config *DatabaseConfig
}

// NewDatabase 创建数据库连接
func NewDatabase(config *DatabaseConfig) (*Database, error) {
	var dialector gorm.Dialector

	// 设置默认值
	if config.MaxOpenConns == 0 {
		config.MaxOpenConns = 25
	}
	if config.MaxIdleConns == 0 {
		config.MaxIdleConns = 5
	}
	if config.ConnMaxLifetime == 0 {
		config.ConnMaxLifetime = 5 * time.Minute
	}

	// 根据驱动类型选择适当的方言
	switch config.Driver {
	case "sqlite":
		dialector = sqlite.Open(config.DSN)
	case "mysql":
		dialector = mysql.Open(config.DSN)
	case "postgres", "postgresql":
		dialector = postgres.Open(config.DSN)
	default:
		// 默认使用 SQLite
		if config.DSN == "" {
			config.DSN = "data/prism.db"
		}
		dialector = sqlite.Open(config.DSN)
		config.Driver = "sqlite"
	}

	// 配置 GORM
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false, // 使用复数表名
		},
		DisableForeignKeyConstraintWhenMigrating: false,
	}

	// 设置日志级别
	if config.Debug {
		gormConfig.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: false,
				Colorful:                  true,
			},
		)
	} else {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	// 创建数据库连接
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 获取底层 sql.DB 对象
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	database := &Database{
		DB:     db,
		config: config,
	}

	return database, nil
}

// AutoMigrate 自动迁移数据库表
func (db *Database) AutoMigrate() error {
	models := []interface{}{
		&Subscription{},
		&NodePool{},
		&Node{},
		&NodeTest{},
		&SubscriptionLog{},
		&NodePoolSubscription{},
	}

	for _, model := range models {
		if err := db.DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %w", model, err)
		}
	}

	// 创建自定义索引（如果需要）
	if err := db.createCustomIndexes(); err != nil {
		return fmt.Errorf("failed to create custom indexes: %w", err)
	}

	return nil
}

// createCustomIndexes 创建自定义索引
func (db *Database) createCustomIndexes() error {
	// 复合索引
	indexes := []string{
		// 节点表的复合索引
		"CREATE INDEX IF NOT EXISTS idx_nodes_server_port ON nodes(server, port)",
		"CREATE INDEX IF NOT EXISTS idx_nodes_subscription_status ON nodes(subscription_id, status)",
		"CREATE INDEX IF NOT EXISTS idx_nodes_country_protocol ON nodes(country, protocol)",
		"CREATE INDEX IF NOT EXISTS idx_nodes_delay_status ON nodes(delay, status) WHERE delay IS NOT NULL",

		// 测试记录表的复合索引
		"CREATE INDEX IF NOT EXISTS idx_node_tests_node_type_time ON node_tests(node_id, test_type, tested_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_node_tests_success_time ON node_tests(success, tested_at DESC)",

		// 订阅日志表的复合索引
		"CREATE INDEX IF NOT EXISTS idx_subscription_logs_sub_success ON subscription_logs(subscription_id, success, created_at DESC)",
	}

	// 只有 SQLite 需要手动创建索引
	if db.config != nil && db.config.Driver == "sqlite" || 
	   db.config == nil && db.Dialector.Name() == "sqlite" {
		for _, indexSQL := range indexes {
			if err := db.Exec(indexSQL).Error; err != nil {
				// 索引可能已存在，记录但不返回错误
				fmt.Printf("Warning: failed to create index: %v\n", err)
			}
		}
	}

	return nil
}

// Close 关闭数据库连接
func (db *Database) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Ping 测试数据库连接
func (db *Database) Ping() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// GetStats 获取数据库统计信息
func (db *Database) GetStats() map[string]interface{} {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration,
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_idle_time_closed": stats.MaxIdleTimeClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}
}

// Transaction 执行事务
func (db *Database) Transaction(fn func(*gorm.DB) error) error {
	return db.DB.Transaction(fn)
}

// HealthCheck 健康检查
func (db *Database) HealthCheck() error {
	// 检查连接
	if err := db.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	// 执行简单查询
	var count int64
	if err := db.Model(&Subscription{}).Count(&count).Error; err != nil {
		return fmt.Errorf("database query failed: %w", err)
	}

	return nil
}

// Cleanup 清理过期数据
func (db *Database) Cleanup(olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)

	// 清理过期的测试记录
	result := db.Where("tested_at < ?", cutoff).Delete(&NodeTest{})
	if result.Error != nil {
		return fmt.Errorf("failed to cleanup node tests: %w", result.Error)
	}

	// 清理过期的订阅日志
	result = db.Where("created_at < ?", cutoff).Delete(&SubscriptionLog{})
	if result.Error != nil {
		return fmt.Errorf("failed to cleanup subscription logs: %w", result.Error)
	}

	return nil
}

// OptimizeDatabase 优化数据库
func (db *Database) OptimizeDatabase() error {
	switch db.config.Driver {
	case "sqlite":
		// SQLite 优化
		if err := db.Exec("VACUUM").Error; err != nil {
			return fmt.Errorf("failed to vacuum sqlite database: %w", err)
		}
		if err := db.Exec("ANALYZE").Error; err != nil {
			return fmt.Errorf("failed to analyze sqlite database: %w", err)
		}
	case "mysql":
		// MySQL 优化
		if err := db.Exec("OPTIMIZE TABLE subscriptions, node_pools, nodes, node_tests, subscription_logs").Error; err != nil {
			return fmt.Errorf("failed to optimize mysql tables: %w", err)
		}
	case "postgres", "postgresql":
		// PostgreSQL 优化
		if err := db.Exec("VACUUM ANALYZE").Error; err != nil {
			return fmt.Errorf("failed to vacuum analyze postgresql database: %w", err)
		}
	}

	return nil
}
