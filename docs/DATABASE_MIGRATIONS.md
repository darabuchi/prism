# 数据库迁移脚本

## 概述

本文档包含 Prism 项目的数据库迁移脚本，支持多数据库引擎（SQLite、MySQL、PostgreSQL、GaussDB）。使用 GORM 的 AutoMigrate 功能和手动迁移脚本相结合的方式。

## 迁移策略

### 1. 自动迁移 (GORM AutoMigrate)
开发环境和简单部署使用 GORM 自动迁移
```go
db.AutoMigrate(
    &User{},
    &ProxyNode{},
    &ProxyGroup{},
    &Subscription{},
    &ProxyRule{},
    // ... 其他模型
)
```

### 2. 手动迁移脚本
生产环境使用版本化迁移脚本确保数据安全

## 迁移文件结构

```
migrations/
├── 001_create_users_table.up.sql
├── 001_create_users_table.down.sql
├── 002_create_proxy_nodes_table.up.sql
├── 002_create_proxy_nodes_table.down.sql
├── 003_create_proxy_groups_table.up.sql
├── 003_create_proxy_groups_table.down.sql
├── 004_create_subscriptions_table.up.sql
├── 004_create_subscriptions_table.down.sql
├── 005_create_proxy_rules_table.up.sql
├── 005_create_proxy_rules_table.down.sql
├── 006_create_configurations_table.up.sql
├── 006_create_configurations_table.down.sql
├── 007_create_statistics_tables.up.sql
├── 007_create_statistics_tables.down.sql
├── 008_create_system_tables.up.sql
├── 008_create_system_tables.down.sql
├── 009_add_indexes.up.sql
├── 009_add_indexes.down.sql
└── 010_insert_default_data.up.sql
```

## 迁移脚本内容

### 001_create_users_table.up.sql
```sql
-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    salt VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    avatar_url VARCHAR(500),
    last_login_at TIMESTAMP,
    last_login_ip VARCHAR(45),
    login_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 创建用户会话表
CREATE TABLE IF NOT EXISTS user_sessions (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    token_hash VARCHAR(255) NOT NULL,
    device_type VARCHAR(50),
    device_info TEXT,
    ip_address VARCHAR(45),
    user_agent TEXT,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

### 001_create_users_table.down.sql
```sql
-- 删除用户相关表
DROP TABLE IF EXISTS user_sessions;
DROP TABLE IF EXISTS users;
```

### 002_create_proxy_nodes_table.up.sql
```sql
-- 创建代理节点表
CREATE TABLE IF NOT EXISTS proxy_nodes (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    server VARCHAR(255) NOT NULL,
    port INTEGER NOT NULL,
    config TEXT NOT NULL,
    country_code VARCHAR(10),
    country_name VARCHAR(100),
    region VARCHAR(100),
    city VARCHAR(100),
    isp VARCHAR(100),
    tags TEXT,
    delay_ms INTEGER DEFAULT -1,
    last_test_at TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'unknown',
    error_count INTEGER NOT NULL DEFAULT 0,
    success_count INTEGER NOT NULL DEFAULT 0,
    total_upload BIGINT NOT NULL DEFAULT 0,
    total_download BIGINT NOT NULL DEFAULT 0,
    subscription_id VARCHAR(36),
    sort_order INTEGER NOT NULL DEFAULT 0,
    enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

### 002_create_proxy_nodes_table.down.sql
```sql
-- 删除代理节点表
DROP TABLE IF EXISTS proxy_nodes;
```

### 003_create_proxy_groups_table.up.sql
```sql
-- 创建代理组表
CREATE TABLE IF NOT EXISTS proxy_groups (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    strategy VARCHAR(50),
    url VARCHAR(500),
    interval_seconds INTEGER DEFAULT 300,
    tolerance_ms INTEGER DEFAULT 150,
    config TEXT,
    sort_order INTEGER NOT NULL DEFAULT 0,
    enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 创建代理组节点关联表
CREATE TABLE IF NOT EXISTS proxy_group_nodes (
    id VARCHAR(36) PRIMARY KEY,
    group_id VARCHAR(36) NOT NULL,
    node_id VARCHAR(36) NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(group_id, node_id)
);
```

### 003_create_proxy_groups_table.down.sql
```sql
-- 删除代理组相关表
DROP TABLE IF EXISTS proxy_group_nodes;
DROP TABLE IF EXISTS proxy_groups;
```

### 004_create_subscriptions_table.up.sql
```sql
-- 创建订阅表
CREATE TABLE IF NOT EXISTS subscriptions (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url TEXT NOT NULL,
    type VARCHAR(50) NOT NULL DEFAULT 'clash',
    user_agent VARCHAR(500),
    headers TEXT,
    node_count INTEGER NOT NULL DEFAULT 0,
    last_update_at TIMESTAMP,
    next_update_at TIMESTAMP,
    update_interval_hours INTEGER NOT NULL DEFAULT 24,
    auto_update BOOLEAN NOT NULL DEFAULT true,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    error_message TEXT,
    error_count INTEGER NOT NULL DEFAULT 0,
    success_count INTEGER NOT NULL DEFAULT 0,
    total_traffic BIGINT NOT NULL DEFAULT 0,
    tags TEXT,
    sort_order INTEGER NOT NULL DEFAULT 0,
    enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 创建订阅日志表
CREATE TABLE IF NOT EXISTS subscription_logs (
    id VARCHAR(36) PRIMARY KEY,
    subscription_id VARCHAR(36) NOT NULL,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    message TEXT,
    nodes_added INTEGER NOT NULL DEFAULT 0,
    nodes_updated INTEGER NOT NULL DEFAULT 0,
    nodes_removed INTEGER NOT NULL DEFAULT 0,
    duration_ms INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (subscription_id) REFERENCES subscriptions(id) ON DELETE CASCADE
);

-- 添加外键约束
ALTER TABLE proxy_nodes ADD CONSTRAINT fk_proxy_nodes_subscription 
    FOREIGN KEY (subscription_id) REFERENCES subscriptions(id) ON DELETE SET NULL;
```

### 004_create_subscriptions_table.down.sql
```sql
-- 删除外键约束
ALTER TABLE proxy_nodes DROP CONSTRAINT IF EXISTS fk_proxy_nodes_subscription;

-- 删除订阅相关表
DROP TABLE IF EXISTS subscription_logs;
DROP TABLE IF EXISTS subscriptions;
```

### 005_create_proxy_rules_table.up.sql
```sql
-- 创建代理规则表
CREATE TABLE IF NOT EXISTS proxy_rules (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    payload TEXT NOT NULL,
    proxy_target VARCHAR(255) NOT NULL,
    priority INTEGER NOT NULL DEFAULT 0,
    description TEXT,
    source VARCHAR(100),
    tags TEXT,
    hit_count BIGINT NOT NULL DEFAULT 0,
    last_hit_at TIMESTAMP,
    enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 创建规则集表
CREATE TABLE IF NOT EXISTS rule_sets (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    url TEXT,
    type VARCHAR(50) NOT NULL DEFAULT 'manual',
    format VARCHAR(50) NOT NULL DEFAULT 'clash',
    rule_count INTEGER NOT NULL DEFAULT 0,
    last_update_at TIMESTAMP,
    update_interval_hours INTEGER DEFAULT 24,
    auto_update BOOLEAN NOT NULL DEFAULT false,
    tags TEXT,
    enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 创建规则集规则关联表
CREATE TABLE IF NOT EXISTS rule_set_rules (
    id VARCHAR(36) PRIMARY KEY,
    rule_set_id VARCHAR(36) NOT NULL,
    rule_id VARCHAR(36) NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (rule_set_id) REFERENCES rule_sets(id) ON DELETE CASCADE,
    FOREIGN KEY (rule_id) REFERENCES proxy_rules(id) ON DELETE CASCADE,
    UNIQUE(rule_set_id, rule_id)
);
```

### 005_create_proxy_rules_table.down.sql
```sql
-- 删除规则相关表
DROP TABLE IF EXISTS rule_set_rules;
DROP TABLE IF EXISTS rule_sets;
DROP TABLE IF EXISTS proxy_rules;
```

### 006_create_configurations_table.up.sql
```sql
-- 创建配置表
CREATE TABLE IF NOT EXISTS configurations (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    config_data TEXT NOT NULL,
    format VARCHAR(20) NOT NULL DEFAULT 'yaml',
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT false,
    version VARCHAR(50),
    tags TEXT,
    created_by VARCHAR(36),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 创建配置历史表
CREATE TABLE IF NOT EXISTS configuration_history (
    id VARCHAR(36) PRIMARY KEY,
    config_id VARCHAR(36) NOT NULL,
    version VARCHAR(50) NOT NULL,
    config_data TEXT NOT NULL,
    change_summary TEXT,
    changed_by VARCHAR(36),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (config_id) REFERENCES configurations(id) ON DELETE CASCADE,
    FOREIGN KEY (changed_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 创建系统设置表
CREATE TABLE IF NOT EXISTS system_settings (
    id VARCHAR(36) PRIMARY KEY,
    category VARCHAR(100) NOT NULL,
    key_name VARCHAR(255) NOT NULL,
    value_data TEXT,
    data_type VARCHAR(50) NOT NULL,
    description TEXT,
    is_public BOOLEAN NOT NULL DEFAULT false,
    is_readonly BOOLEAN NOT NULL DEFAULT false,
    validation_rule TEXT,
    default_value TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(category, key_name)
);
```

### 006_create_configurations_table.down.sql
```sql
-- 删除配置相关表
DROP TABLE IF EXISTS system_settings;
DROP TABLE IF EXISTS configuration_history;
DROP TABLE IF EXISTS configurations;
```

### 007_create_statistics_tables.up.sql
```sql
-- 创建流量统计表
CREATE TABLE IF NOT EXISTS traffic_stats (
    id VARCHAR(36) PRIMARY KEY,
    stat_type VARCHAR(50) NOT NULL,
    target_id VARCHAR(36),
    date_key VARCHAR(10) NOT NULL,
    hour_key VARCHAR(13),
    upload_bytes BIGINT NOT NULL DEFAULT 0,
    download_bytes BIGINT NOT NULL DEFAULT 0,
    connection_count INTEGER NOT NULL DEFAULT 0,
    request_count INTEGER NOT NULL DEFAULT 0,
    error_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 创建连接日志表
CREATE TABLE IF NOT EXISTS connection_logs (
    id VARCHAR(36) PRIMARY KEY,
    session_id VARCHAR(100),
    user_id VARCHAR(36),
    node_id VARCHAR(36),
    source_ip VARCHAR(45),
    destination_host VARCHAR(255),
    destination_ip VARCHAR(45),
    destination_port INTEGER,
    protocol VARCHAR(20),
    rule_matched VARCHAR(255),
    upload_bytes BIGINT NOT NULL DEFAULT 0,
    download_bytes BIGINT NOT NULL DEFAULT 0,
    duration_ms INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL,
    error_message TEXT,
    started_at TIMESTAMP NOT NULL,
    ended_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (node_id) REFERENCES proxy_nodes(id) ON DELETE SET NULL
);

-- 创建系统指标表
CREATE TABLE IF NOT EXISTS system_metrics (
    id VARCHAR(36) PRIMARY KEY,
    metric_name VARCHAR(100) NOT NULL,
    metric_value DECIMAL(10,2) NOT NULL,
    metric_unit VARCHAR(20),
    tags TEXT,
    recorded_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### 007_create_statistics_tables.down.sql
```sql
-- 删除统计相关表
DROP TABLE IF EXISTS system_metrics;
DROP TABLE IF EXISTS connection_logs;
DROP TABLE IF EXISTS traffic_stats;
```

### 008_create_system_tables.up.sql
```sql
-- 创建审计日志表
CREATE TABLE IF NOT EXISTS audit_logs (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36),
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(100),
    resource_id VARCHAR(36),
    resource_name VARCHAR(255),
    old_values TEXT,
    new_values TEXT,
    ip_address VARCHAR(45),
    user_agent TEXT,
    request_id VARCHAR(100),
    status VARCHAR(20) NOT NULL,
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- 创建通知表
CREATE TABLE IF NOT EXISTS notifications (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36),
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    action_url VARCHAR(500),
    action_text VARCHAR(100),
    is_read BOOLEAN NOT NULL DEFAULT false,
    read_at TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

### 008_create_system_tables.down.sql
```sql
-- 删除系统管理表
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS audit_logs;
```

### 009_add_indexes.up.sql
```sql
-- 用户表索引
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);

-- 用户会话表索引
CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_token ON user_sessions(token_hash);
CREATE INDEX IF NOT EXISTS idx_user_sessions_expires ON user_sessions(expires_at);

-- 代理节点表索引
CREATE INDEX IF NOT EXISTS idx_proxy_nodes_type ON proxy_nodes(type);
CREATE INDEX IF NOT EXISTS idx_proxy_nodes_country ON proxy_nodes(country_code);
CREATE INDEX IF NOT EXISTS idx_proxy_nodes_status ON proxy_nodes(status);
CREATE INDEX IF NOT EXISTS idx_proxy_nodes_delay ON proxy_nodes(delay_ms);
CREATE INDEX IF NOT EXISTS idx_proxy_nodes_subscription ON proxy_nodes(subscription_id);
CREATE INDEX IF NOT EXISTS idx_proxy_nodes_enabled ON proxy_nodes(enabled);

-- 代理组表索引
CREATE INDEX IF NOT EXISTS idx_proxy_groups_type ON proxy_groups(type);
CREATE INDEX IF NOT EXISTS idx_proxy_groups_enabled ON proxy_groups(enabled);

-- 代理组节点关联表索引
CREATE INDEX IF NOT EXISTS idx_proxy_group_nodes_group ON proxy_group_nodes(group_id);
CREATE INDEX IF NOT EXISTS idx_proxy_group_nodes_node ON proxy_group_nodes(node_id);

-- 订阅表索引
CREATE INDEX IF NOT EXISTS idx_subscriptions_status ON subscriptions(status);
CREATE INDEX IF NOT EXISTS idx_subscriptions_next_update ON subscriptions(next_update_at);
CREATE INDEX IF NOT EXISTS idx_subscriptions_enabled ON subscriptions(enabled);

-- 订阅日志表索引
CREATE INDEX IF NOT EXISTS idx_subscription_logs_subscription ON subscription_logs(subscription_id);
CREATE INDEX IF NOT EXISTS idx_subscription_logs_created ON subscription_logs(created_at);

-- 代理规则表索引
CREATE INDEX IF NOT EXISTS idx_proxy_rules_type ON proxy_rules(type);
CREATE INDEX IF NOT EXISTS idx_proxy_rules_priority ON proxy_rules(priority DESC);
CREATE INDEX IF NOT EXISTS idx_proxy_rules_enabled ON proxy_rules(enabled);
CREATE INDEX IF NOT EXISTS idx_proxy_rules_source ON proxy_rules(source);

-- 规则集表索引
CREATE INDEX IF NOT EXISTS idx_rule_sets_type ON rule_sets(type);
CREATE INDEX IF NOT EXISTS idx_rule_sets_enabled ON rule_sets(enabled);

-- 规则集规则关联表索引
CREATE INDEX IF NOT EXISTS idx_rule_set_rules_set ON rule_set_rules(rule_set_id);
CREATE INDEX IF NOT EXISTS idx_rule_set_rules_rule ON rule_set_rules(rule_id);

-- 配置表索引
CREATE INDEX IF NOT EXISTS idx_configurations_type ON configurations(type);
CREATE INDEX IF NOT EXISTS idx_configurations_active ON configurations(is_active);
CREATE INDEX IF NOT EXISTS idx_configurations_creator ON configurations(created_by);

-- 配置历史表索引
CREATE INDEX IF NOT EXISTS idx_configuration_history_config ON configuration_history(config_id);
CREATE INDEX IF NOT EXISTS idx_configuration_history_created ON configuration_history(created_at);

-- 系统设置表索引
CREATE INDEX IF NOT EXISTS idx_system_settings_category ON system_settings(category);
CREATE INDEX IF NOT EXISTS idx_system_settings_public ON system_settings(is_public);

-- 流量统计表索引
CREATE INDEX IF NOT EXISTS idx_traffic_stats_type_target ON traffic_stats(stat_type, target_id);
CREATE INDEX IF NOT EXISTS idx_traffic_stats_date ON traffic_stats(date_key);
CREATE INDEX IF NOT EXISTS idx_traffic_stats_hour ON traffic_stats(hour_key);
CREATE UNIQUE INDEX IF NOT EXISTS idx_traffic_stats_unique ON traffic_stats(stat_type, target_id, date_key, hour_key);

-- 连接日志表索引
CREATE INDEX IF NOT EXISTS idx_connection_logs_user ON connection_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_connection_logs_node ON connection_logs(node_id);
CREATE INDEX IF NOT EXISTS idx_connection_logs_started ON connection_logs(started_at);
CREATE INDEX IF NOT EXISTS idx_connection_logs_destination ON connection_logs(destination_host);

-- 系统指标表索引
CREATE INDEX IF NOT EXISTS idx_system_metrics_name ON system_metrics(metric_name);
CREATE INDEX IF NOT EXISTS idx_system_metrics_recorded ON system_metrics(recorded_at);

-- 审计日志表索引
CREATE INDEX IF NOT EXISTS idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created ON audit_logs(created_at);

-- 通知表索引
CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_unread ON notifications(user_id, is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_created ON notifications(created_at);
```

### 009_add_indexes.down.sql
```sql
-- 删除所有索引
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_status;
-- ... (其他所有索引的删除语句)
```

### 010_insert_default_data.up.sql
```sql
-- 插入默认系统设置
INSERT OR IGNORE INTO system_settings (id, category, key_name, value_data, data_type, description, is_public, is_readonly) VALUES
('sys-001', 'general', 'app_name', '"Prism"', 'string', '应用名称', true, false),
('sys-002', 'general', 'app_version', '"1.0.0"', 'string', '应用版本', true, true),
('sys-003', 'proxy', 'http_port', '7890', 'number', 'HTTP代理端口', false, false),
('sys-004', 'proxy', 'socks_port', '7891', 'number', 'SOCKS5代理端口', false, false),
('sys-005', 'proxy', 'mixed_port', '7892', 'number', '混合代理端口', false, false),
('sys-006', 'proxy', 'api_port', '9090', 'number', 'API端口', false, false),
('sys-007', 'security', 'jwt_secret', '', 'string', 'JWT密钥', false, false),
('sys-008', 'security', 'jwt_expire_hours', '168', 'number', 'JWT过期时间(小时)', false, false),
('sys-009', 'ui', 'theme', '"light"', 'string', '界面主题', true, false),
('sys-010', 'ui', 'language', '"zh-CN"', 'string', '界面语言', true, false);

-- 创建默认管理员用户 (用户名: admin, 密码: admin123)
-- 注意: 生产环境请立即修改默认密码
INSERT OR IGNORE INTO users (id, username, password_hash, salt, role, status) VALUES
('admin-001', 'admin', 'hash_placeholder', 'salt_placeholder', 'admin', 'active');

-- 插入默认代理规则
INSERT OR IGNORE INTO proxy_rules (id, name, type, payload, proxy_target, priority, description, source, enabled) VALUES
('rule-001', 'Local Network', 'IP-CIDR', '127.0.0.0/8', 'DIRECT', 1000, '本地回环地址', 'builtin', true),
('rule-002', 'Local Network', 'IP-CIDR', '192.168.0.0/16', 'DIRECT', 1000, '局域网地址', 'builtin', true),
('rule-003', 'Local Network', 'IP-CIDR', '10.0.0.0/8', 'DIRECT', 1000, '局域网地址', 'builtin', true),
('rule-004', 'Local Network', 'IP-CIDR', '172.16.0.0/12', 'DIRECT', 1000, '局域网地址', 'builtin', true),
('rule-005', 'China IP', 'GEOIP', 'CN', 'DIRECT', 500, '中国大陆IP', 'builtin', true);

-- 创建默认代理组
INSERT OR IGNORE INTO proxy_groups (id, name, type, strategy, url, interval_seconds, tolerance_ms, enabled) VALUES
('group-001', 'Auto Select', 'url-test', '', 'http://www.gstatic.com/generate_204', 300, 150, true),
('group-002', 'Manual Select', 'select', '', '', 0, 0, true),
('group-003', 'Fallback', 'fallback', '', 'http://www.gstatic.com/generate_204', 300, 150, true);
```

### 010_insert_default_data.down.sql
```sql
-- 删除默认数据
DELETE FROM proxy_groups WHERE id LIKE 'group-%';
DELETE FROM proxy_rules WHERE id LIKE 'rule-%';
DELETE FROM users WHERE id = 'admin-001';
DELETE FROM system_settings WHERE id LIKE 'sys-%';
```

## 多数据库兼容性处理

### 数据类型映射

#### SQLite
```sql
-- SQLite 特定类型
TEXT        -- 对应 VARCHAR/TEXT
INTEGER     -- 对应 INT/BIGINT
REAL        -- 对应 DECIMAL
BOOLEAN     -- 对应 BOOLEAN (存储为 INTEGER)
TIMESTAMP   -- 对应 DATETIME (存储为 TEXT)
```

#### MySQL
```sql
-- MySQL 特定类型
VARCHAR(n)          -- 字符串
TEXT                -- 长文本
INT, BIGINT         -- 整数
DECIMAL(10,2)       -- 小数
BOOLEAN             -- 布尔值
TIMESTAMP           -- 时间戳
DATETIME            -- 日期时间
```

#### PostgreSQL
```sql
-- PostgreSQL 特定类型
VARCHAR(n)          -- 字符串
TEXT                -- 长文本
INTEGER, BIGINT     -- 整数
DECIMAL(10,2)       -- 小数
BOOLEAN             -- 布尔值
TIMESTAMP           -- 时间戳
TIMESTAMPTZ         -- 带时区时间戳
```

#### GaussDB
```sql
-- GaussDB 特定类型
VARCHAR2(n)         -- 字符串
CLOB                -- 长文本
INTEGER, NUMBER     -- 整数
NUMBER(10,2)        -- 小数
BOOLEAN             -- 布尔值
TIMESTAMP           -- 时间戳
```

## 迁移管理工具

### Go 迁移管理器
```go
package migration

import (
    "fmt"
    "gorm.io/gorm"
    "github.com/lazygophers/log"
)

type Migration struct {
    Version string
    Name    string
    UpSQL   string
    DownSQL string
}

type Migrator struct {
    db *gorm.DB
    migrations []Migration
}

func NewMigrator(db *gorm.DB) *Migrator {
    return &Migrator{
        db: db,
        migrations: []Migration{
            {Version: "001", Name: "create_users_table", UpSQL: users001Up, DownSQL: users001Down},
            {Version: "002", Name: "create_proxy_nodes_table", UpSQL: nodes002Up, DownSQL: nodes002Down},
            // ... 其他迁移
        },
    }
}

func (m *Migrator) Up() error {
    // 创建迁移记录表
    if err := m.createMigrationTable(); err != nil {
        return err
    }
    
    for _, migration := range m.migrations {
        if err := m.runMigration(migration); err != nil {
            log.Error("迁移失败", 
                log.String("version", migration.Version),
                log.String("name", migration.Name),
                log.Error(err),
            )
            return err
        }
    }
    
    return nil
}

func (m *Migrator) Down(targetVersion string) error {
    // 回滚到指定版本
    for i := len(m.migrations) - 1; i >= 0; i-- {
        migration := m.migrations[i]
        if migration.Version == targetVersion {
            break
        }
        
        if err := m.rollbackMigration(migration); err != nil {
            return err
        }
    }
    
    return nil
}

func (m *Migrator) createMigrationTable() error {
    return m.db.Exec(`
        CREATE TABLE IF NOT EXISTS schema_migrations (
            version VARCHAR(50) PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            executed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        )
    `).Error
}

func (m *Migrator) runMigration(migration Migration) error {
    // 检查是否已执行
    var count int64
    m.db.Raw("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", migration.Version).Scan(&count)
    if count > 0 {
        log.Info("迁移已执行，跳过", log.String("version", migration.Version))
        return nil
    }
    
    // 执行迁移
    if err := m.db.Exec(migration.UpSQL).Error; err != nil {
        return err
    }
    
    // 记录迁移
    return m.db.Exec(
        "INSERT INTO schema_migrations (version, name) VALUES (?, ?)",
        migration.Version, migration.Name,
    ).Error
}

func (m *Migrator) rollbackMigration(migration Migration) error {
    // 执行回滚
    if err := m.db.Exec(migration.DownSQL).Error; err != nil {
        return err
    }
    
    // 删除迁移记录
    return m.db.Exec(
        "DELETE FROM schema_migrations WHERE version = ?",
        migration.Version,
    ).Error
}
```

### 命令行工具
```go
// cmd/migrate/main.go
package main

import (
    "flag"
    "log"
    "prism/internal/database"
    "prism/internal/migration"
)

func main() {
    var (
        action = flag.String("action", "up", "Migration action: up, down, status")
        target = flag.String("target", "", "Target version for rollback")
        config = flag.String("config", "config.yaml", "Database config file")
    )
    flag.Parse()
    
    // 连接数据库
    db, err := database.Connect(*config)
    if err != nil {
        log.Fatal("Database connection failed:", err)
    }
    
    migrator := migration.NewMigrator(db)
    
    switch *action {
    case "up":
        if err := migrator.Up(); err != nil {
            log.Fatal("Migration failed:", err)
        }
        log.Println("Migration completed successfully")
    case "down":
        if *target == "" {
            log.Fatal("Target version required for rollback")
        }
        if err := migrator.Down(*target); err != nil {
            log.Fatal("Rollback failed:", err)
        }
        log.Println("Rollback completed successfully")
    case "status":
        migrator.ShowStatus()
    default:
        log.Fatal("Invalid action:", *action)
    }
}
```

### 使用示例
```bash
# 执行所有迁移
go run cmd/migrate/main.go -action=up -config=config.yaml

# 回滚到指定版本
go run cmd/migrate/main.go -action=down -target=005 -config=config.yaml

# 查看迁移状态
go run cmd/migrate/main.go -action=status -config=config.yaml
```

## 数据备份和恢复

### 备份脚本
```bash
#!/bin/bash
# backup.sh

DATABASE_TYPE="${DATABASE_TYPE:-sqlite}"
BACKUP_DIR="${BACKUP_DIR:-./backups}"
DATE=$(date +%Y%m%d_%H%M%S)

case $DATABASE_TYPE in
    "sqlite")
        cp prism.db "$BACKUP_DIR/prism_$DATE.db"
        ;;
    "mysql")
        mysqldump -h $DB_HOST -u $DB_USER -p$DB_PASSWORD $DB_NAME > "$BACKUP_DIR/prism_$DATE.sql"
        ;;
    "postgres")
        pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME > "$BACKUP_DIR/prism_$DATE.sql"
        ;;
esac

echo "Backup completed: $BACKUP_DIR/prism_$DATE.*"
```

### 恢复脚本
```bash
#!/bin/bash
# restore.sh

BACKUP_FILE="$1"
DATABASE_TYPE="${DATABASE_TYPE:-sqlite}"

case $DATABASE_TYPE in
    "sqlite")
        cp "$BACKUP_FILE" prism.db
        ;;
    "mysql")
        mysql -h $DB_HOST -u $DB_USER -p$DB_PASSWORD $DB_NAME < "$BACKUP_FILE"
        ;;
    "postgres")
        psql -h $DB_HOST -U $DB_USER -d $DB_NAME < "$BACKUP_FILE"
        ;;
esac

echo "Restore completed from: $BACKUP_FILE"
```

## 测试数据

### 开发环境测试数据
```sql
-- 测试用户
INSERT INTO users (id, username, email, password_hash, salt, role) VALUES
('test-user-001', 'testuser', 'test@example.com', 'test_hash', 'test_salt', 'user');

-- 测试订阅
INSERT INTO subscriptions (id, name, url, type, enabled) VALUES
('test-sub-001', 'Test Subscription', 'https://example.com/sub', 'clash', true);

-- 测试节点
INSERT INTO proxy_nodes (id, name, type, server, port, config, subscription_id, enabled) VALUES
('test-node-001', 'Test Node 1', 'vmess', 'example1.com', 443, '{"uuid":"test"}', 'test-sub-001', true),
('test-node-002', 'Test Node 2', 'trojan', 'example2.com', 443, '{"password":"test"}', 'test-sub-001', true);

-- 测试规则
INSERT INTO proxy_rules (id, name, type, payload, proxy_target, priority, enabled) VALUES
('test-rule-001', 'Test Domain Rule', 'DOMAIN', 'example.com', 'Test Node 1', 100, true),
('test-rule-002', 'Test IP Rule', 'IP-CIDR', '1.1.1.0/24', 'DIRECT', 200, true);
```

### 清理测试数据
```sql
DELETE FROM proxy_rules WHERE id LIKE 'test-rule-%';
DELETE FROM proxy_nodes WHERE id LIKE 'test-node-%';
DELETE FROM subscriptions WHERE id LIKE 'test-sub-%';
DELETE FROM users WHERE id LIKE 'test-user-%';
```