# Prism 数据库设计 (最小化客户端版)

## 设计理念

作为代理客户端，性能是第一优先级。数据库设计遵循**最小必要原则**，只保留核心功能必需的表，移除所有可能影响代理性能的统计、日志和监控表。

### 核心原则
- **性能优先**: 最小化数据库操作，避免影响代理转发性能
- **功能精简**: 只保留订阅管理和节点管理核心功能
- **内存优化**: 减少数据库连接和查询频率
- **启动快速**: 减少启动时的数据库初始化时间

## 最小化表结构 (仅3个核心表)

### 1. subscriptions (订阅表)
```sql
CREATE TABLE subscriptions (
    id VARCHAR(36) PRIMARY KEY,                   -- UUID
    name VARCHAR(255) NOT NULL,                   -- 订阅名称
    url TEXT NOT NULL,                            -- 订阅链接
    type VARCHAR(50) NOT NULL DEFAULT 'clash',    -- 订阅类型: clash, v2ray, ss
    
    -- 基本配置
    update_interval_hours INT NOT NULL DEFAULT 24, -- 更新间隔(小时)
    auto_update BOOLEAN NOT NULL DEFAULT true,    -- 是否自动更新
    enabled BOOLEAN NOT NULL DEFAULT true,        -- 是否启用
    
    -- 状态信息 (最小化)
    last_update_at DATETIME,                      -- 最后更新时间
    next_update_at DATETIME,                      -- 下次更新时间
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- 状态: active, error, disabled
    node_count INT NOT NULL DEFAULT 0,            -- 节点数量
    
    -- 系统字段
    created_at DATETIME NOT NULL,                 -- 创建时间
    updated_at DATETIME NOT NULL,                 -- 更新时间
    deleted_at DATETIME                           -- 软删除时间
);

-- 最小化索引
CREATE INDEX idx_subscriptions_enabled ON subscriptions(enabled);
CREATE INDEX idx_subscriptions_auto_update ON subscriptions(auto_update, next_update_at);
```

### 2. proxy_nodes (代理节点表)
```sql
CREATE TABLE proxy_nodes (
    id VARCHAR(36) PRIMARY KEY,                   -- UUID
    name VARCHAR(255) NOT NULL,                   -- 节点名称
    type VARCHAR(50) NOT NULL,                    -- 协议类型: vmess, vless, trojan, ss, ssr
    server VARCHAR(255) NOT NULL,                 -- 服务器地址
    port INT NOT NULL,                            -- 端口号
    config TEXT NOT NULL,                         -- 节点配置 JSON格式
    
    -- 地理信息 (简化)
    country_code VARCHAR(10),                     -- 国家代码: US, JP, HK
    
    -- 基本性能信息
    delay_ms INT DEFAULT -1,                      -- 延迟(毫秒) -1表示未测试
    last_test_at DATETIME,                        -- 最后测速时间
    status VARCHAR(20) NOT NULL DEFAULT 'unknown', -- 状态: online, offline, unknown
    
    -- 基本元数据
    enabled BOOLEAN NOT NULL DEFAULT true,        -- 是否启用
    sort_order INT NOT NULL DEFAULT 0,            -- 排序权重
    
    -- 系统字段
    created_at DATETIME NOT NULL,                 -- 创建时间
    updated_at DATETIME NOT NULL,                 -- 更新时间
    deleted_at DATETIME                           -- 软删除时间
);

-- 最小化索引
CREATE INDEX idx_proxy_nodes_enabled ON proxy_nodes(enabled);
CREATE INDEX idx_proxy_nodes_type ON proxy_nodes(type);
CREATE INDEX idx_proxy_nodes_server_port ON proxy_nodes(server, port);
CREATE INDEX idx_proxy_nodes_delay ON proxy_nodes(delay_ms);
```

### 3. subscription_nodes (订阅节点关联表)
```sql
CREATE TABLE subscription_nodes (
    id VARCHAR(36) PRIMARY KEY,                   -- UUID
    subscription_id VARCHAR(36) NOT NULL,         -- 订阅ID
    node_id VARCHAR(36) NOT NULL,                 -- 节点ID
    
    -- 关联信息 (最小化)
    node_index INT NOT NULL,                      -- 节点在订阅中的索引位置
    is_primary BOOLEAN NOT NULL DEFAULT false,    -- 是否为主要来源(用于去重)
    
    -- 系统字段
    created_at DATETIME NOT NULL,                 -- 关联创建时间
    
    FOREIGN KEY (subscription_id) REFERENCES subscriptions(id) ON DELETE CASCADE,
    FOREIGN KEY (node_id) REFERENCES proxy_nodes(id) ON DELETE CASCADE,
    UNIQUE(subscription_id, node_id)
);

-- 最小化索引
CREATE INDEX idx_subscription_nodes_subscription ON subscription_nodes(subscription_id);
CREATE INDEX idx_subscription_nodes_node ON subscription_nodes(node_id);
```

## 移除的表 (性能考虑)

### ❌ 完全移除的表
```sql
-- 性能影响较大的表，完全移除：

-- 1. subscription_logs - 订阅日志
-- 原因：频繁写入影响性能，改用文件日志

-- 2. node_tests - 节点测试记录
-- 原因：大量测试数据影响数据库性能，改用内存缓存

-- 3. traffic_stats - 流量统计
-- 原因：实时统计会严重影响代理性能


-- 6. system_settings - 系统设置
-- 原因：改用配置文件管理
```

## 替代方案

### 📄 文件日志系统
```
data/
├── logs/
│   ├── app.log              # 应用日志
│   ├── subscription.log     # 订阅更新日志
│   └── node_test.log        # 节点测试日志
├── cache/
│   ├── node_tests.json      # 节点测试结果缓存
│   └── traffic_stats.json   # 流量统计缓存
└── config/
    └── settings.yaml        # 系统配置文件
```

### ⚡ 内存缓存策略
- **BoltDB**: 缓存节点测试结果和性能数据
- **LevelDB**: 缓存流量统计数据
- **内存**: 缓存活跃节点信息和配置

## Go 模型定义 (简化版)

### 基础模型
```go
package models

import (
    "time"
    "gorm.io/gorm"
    "github.com/lazygophers/utils/xtime"
    "github.com/google/uuid"
)

// BaseModel 基础模型
type BaseModel struct {
    ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
    CreatedAt xtime.Time     `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt xtime.Time     `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
    if m.ID == "" {
        m.ID = uuid.New().String()
    }
    return nil
}
```

### 订阅模型
```go
// Subscription 订阅模型 (简化版)
type Subscription struct {
    BaseModel
    Name                string     `gorm:"size:255;not null" json:"name"`
    URL                 string     `gorm:"type:text;not null" json:"url"`
    Type                string     `gorm:"size:50;not null;default:clash" json:"type"`
    UpdateIntervalHours int        `gorm:"not null;default:24" json:"update_interval_hours"`
    AutoUpdate          bool       `gorm:"not null;default:true" json:"auto_update"`
    Enabled             bool       `gorm:"not null;default:true" json:"enabled"`
    LastUpdateAt        *xtime.Time `json:"last_update_at"`
    NextUpdateAt        *xtime.Time `json:"next_update_at"`
    Status              string     `gorm:"size:20;not null;default:pending" json:"status"`
    NodeCount           int        `gorm:"not null;default:0" json:"node_count"`
    
    // 关联关系
    Nodes []ProxyNode `gorm:"many2many:subscription_nodes" json:"nodes,omitempty"`
}

// TableName 指定表名
func (Subscription) TableName() string {
    return "subscriptions"
}
```

### 节点模型
```go
// ProxyNode 代理节点模型 (简化版)
type ProxyNode struct {
    BaseModel
    Name        string     `gorm:"size:255;not null" json:"name"`
    Type        string     `gorm:"size:50;not null" json:"type"`
    Server      string     `gorm:"size:255;not null" json:"server"`
    Port        int        `gorm:"not null" json:"port"`
    Config      string     `gorm:"type:text;not null" json:"config"`
    CountryCode string     `gorm:"size:10" json:"country_code"`
    DelayMS     int        `gorm:"default:-1" json:"delay_ms"`
    LastTestAt  *xtime.Time `json:"last_test_at"`
    Status      string     `gorm:"size:20;not null;default:unknown" json:"status"`
    Enabled     bool       `gorm:"not null;default:true" json:"enabled"`
    SortOrder   int        `gorm:"not null;default:0" json:"sort_order"`
    
    // 关联关系
    Subscriptions []Subscription `gorm:"many2many:subscription_nodes" json:"subscriptions,omitempty"`
}

// TableName 指定表名
func (ProxyNode) TableName() string {
    return "proxy_nodes"
}
```

### 关联模型
```go
// SubscriptionNode 订阅节点关联模型 (简化版)
type SubscriptionNode struct {
    BaseModel
    SubscriptionID string `gorm:"size:36;not null" json:"subscription_id"`
    NodeID         string `gorm:"size:36;not null" json:"node_id"`
    NodeIndex      int    `gorm:"not null" json:"node_index"`
    IsPrimary      bool   `gorm:"not null;default:false" json:"is_primary"`
    
    // 关联
    Subscription *Subscription `gorm:"foreignKey:SubscriptionID" json:"subscription,omitempty"`
    Node         *ProxyNode    `gorm:"foreignKey:NodeID" json:"node,omitempty"`
}

// TableName 指定表名
func (SubscriptionNode) TableName() string {
    return "subscription_nodes"
}
```

## 缓存服务设计

### 节点测试结果缓存
```go
type NodeTestCache struct {
    cache *bolt.DB
}

func (c *NodeTestCache) SaveTestResult(nodeID string, result *TestResult) error {
    return c.cache.Update(func(tx *bolt.Tx) error {
        bucket := tx.Bucket([]byte("node_tests"))
        data, _ := json.Marshal(result)
        return bucket.Put([]byte(nodeID), data)
    })
}

func (c *NodeTestCache) GetTestResult(nodeID string) (*TestResult, error) {
    var result TestResult
    err := c.cache.View(func(tx *bolt.Tx) error {
        bucket := tx.Bucket([]byte("node_tests"))
        data := bucket.Get([]byte(nodeID))
        if data == nil {
            return errors.New("not found")
        }
        return json.Unmarshal(data, &result)
    })
    return &result, err
}
```

### 配置文件管理
```go
type ConfigManager struct {
    configPath string
    config     *Config
}

type Config struct {
    Proxy    ProxyConfig    `yaml:"proxy"`
    Update   UpdateConfig   `yaml:"update"`
    Logging  LoggingConfig  `yaml:"logging"`
    Performance PerformanceConfig `yaml:"performance"`
}

type ProxyConfig struct {
    HTTPPort  int  `yaml:"http_port"`
    SOCKSPort int  `yaml:"socks_port"`
    MixedPort int  `yaml:"mixed_port"`
}

type UpdateConfig struct {
    CheckInterval int  `yaml:"check_interval"`
    AutoUpdate    bool `yaml:"auto_update"`
    Timeout       int  `yaml:"timeout"`
}

type PerformanceConfig struct {
    MaxConnections int `yaml:"max_connections"`
    BufferSize     int `yaml:"buffer_size"`
    EnableCache    bool `yaml:"enable_cache"`
}

func (c *ConfigManager) Load() error {
    data, err := os.ReadFile(c.configPath)
    if err != nil {
        return err
    }
    return yaml.Unmarshal(data, &c.config)
}

func (c *ConfigManager) Save() error {
    data, err := yaml.Marshal(c.config)
    if err != nil {
        return err
    }
    return os.WriteFile(c.configPath, data, 0644)
}
```

## 性能优化策略

### 数据库连接优化
```go
func setupDatabase() *gorm.DB {
    db, err := gorm.Open(sqlite.Open("prism.db"), &gorm.Config{
        Logger:                 logger.Default.LogMode(logger.Silent), // 禁用SQL日志
        DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束提升性能
        PrepareStmt:            true, // 预编译SQL语句
        CreateBatchSize:        100,  // 批量创建大小
    })
    
    if err != nil {
        panic("数据库连接失败")
    }
    
    sqlDB, _ := db.DB()
    
    // 连接池优化
    sqlDB.SetMaxIdleConns(2)            // 最小空闲连接
    sqlDB.SetMaxOpenConns(10)           // 最大打开连接
    sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生命周期
    
    return db
}
```

### 查询优化
```go
// 只查询必要字段
func GetActiveNodes() []ProxyNode {
    var nodes []ProxyNode
    db.Select("id", "name", "type", "server", "port", "config", "delay_ms", "status").
       Where("enabled = ? AND deleted_at IS NULL", true).
       Order("sort_order ASC, delay_ms ASC").
       Find(&nodes)
    return nodes
}

// 批量更新节点状态
func UpdateNodesStatus(updates []NodeStatusUpdate) error {
    return db.Transaction(func(tx *gorm.DB) error {
        for _, update := range updates {
            if err := tx.Model(&ProxyNode{}).
                Where("id = ?", update.NodeID).
                Updates(map[string]interface{}{
                    "status": update.Status,
                    "delay_ms": update.DelayMS,
                    "last_test_at": xtime.Now(),
                }).Error; err != nil {
                return err
            }
        }
        return nil
    })
}
```

### 启动优化
```go
func QuickStart() {
    // 1. 最小化数据库初始化
    db := setupDatabase()
    
    // 2. 只加载必要数据
    subscriptions := loadActiveSubscriptions()
    nodes := loadActiveNodes()
    
    // 3. 异步加载非关键数据
    go func() {
        loadRulesFromFiles()
        initializeCache()
    }()
    
    // 4. 快速启动代理核心
    startProxyCore(nodes)
}
```

## 总结

这个最小化设计将数据库表从原来的15个减少到仅3个核心表，移除了所有可能影响代理性能的日志、统计和监控表。通过文件系统和缓存替代重型数据库操作，确保代理客户端的最佳性能表现。

### 性能提升预期
- **启动时间**: 减少80%以上
- **内存占用**: 减少60%以上  
- **数据库操作**: 减少90%以上
- **代理延迟**: 几乎无影响