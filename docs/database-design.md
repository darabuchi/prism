# Prism - 数据库设计文档

## 设计概述

Prism 数据库设计专注于订阅管理和节点管理两个核心功能，支持自动更新、去重、性能监控等特性。设计遵循高性能、高效率原则，支持多种数据库后端。

## 支持的数据库

1. **SQLite** (默认) - 使用 modernc.org/sqlite (纯 Go 实现)
2. **MySQL** - 生产环境推荐
3. **PostgreSQL** - 高并发场景推荐
4. **GaussDB** - 企业级部署

## 数据库架构图

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Subscriptions  │    │   NodePools     │    │      Nodes      │
│                 │    │                 │    │                 │
│ • 订阅信息      │ 1:N│ • 池管理        │ 1:N│ • 节点详情      │
│ • 更新统计      │────│ • 存活统计      │────│ • 性能指标      │
│ • 自动更新      │    │ • 去重逻辑      │    │ • 解锁信息      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │              ┌─────────────────┐
         │              ┌─────────────────┐     │   NodeTests     │
         │              │ SubscriptionLog │     │                 │
         │              │                 │     │ • 测试记录      │
         └──────────────│ • 更新历史      │     │ • 延迟数据      │
                        │ • 节点变化      │     │ • 速度测试      │
                        │ • 错误记录      │     │ • 解锁检测      │
                        └─────────────────┘     └─────────────────┘
```

## 数据表设计

### 1. 订阅表 (subscriptions)

存储订阅链接和基础配置信息。

```sql
CREATE TABLE subscriptions (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    name                VARCHAR(255) NOT NULL COMMENT '订阅名称',
    url                 TEXT NOT NULL COMMENT '订阅链接',
    user_agent          VARCHAR(255) DEFAULT 'clash' COMMENT '请求 User-Agent',
    
    -- 更新配置
    auto_update         BOOLEAN DEFAULT TRUE COMMENT '自动更新开关',
    update_interval     INTEGER DEFAULT 3600 COMMENT '更新间隔(秒)',
    
    -- 统计信息
    total_nodes         INTEGER DEFAULT 0 COMMENT '当前订阅总节点数',
    active_nodes        INTEGER DEFAULT 0 COMMENT '存活节点数',
    unique_new_nodes    INTEGER DEFAULT 0 COMMENT '全局去重后新增节点数',
    
    -- 状态信息
    status              VARCHAR(50) DEFAULT 'active' COMMENT '订阅状态: active/inactive/error',
    last_update         TIMESTAMP NULL COMMENT '最后更新时间',
    last_success        TIMESTAMP NULL COMMENT '最后成功时间',
    error_message       TEXT NULL COMMENT '最后错误信息',
    error_count         INTEGER DEFAULT 0 COMMENT '连续错误次数',
    
    -- 时间戳
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- 索引
    INDEX idx_status (status),
    INDEX idx_last_update (last_update),
    INDEX idx_auto_update (auto_update),
    UNIQUE KEY uk_url (url)
);
```

### 2. 节点池表 (node_pools)

逻辑分组管理，支持多订阅聚合。

```sql
CREATE TABLE node_pools (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    name                VARCHAR(255) NOT NULL COMMENT '节点池名称',
    description         TEXT COMMENT '节点池描述',
    
    -- 统计信息
    total_subscriptions INTEGER DEFAULT 0 COMMENT '关联订阅数',
    total_nodes         INTEGER DEFAULT 0 COMMENT '节点总数',
    active_nodes        INTEGER DEFAULT 0 COMMENT '存活节点数',
    survival_rate       DECIMAL(5,2) DEFAULT 0.00 COMMENT '存活率(%)',
    
    -- 配置信息
    enabled             BOOLEAN DEFAULT TRUE COMMENT '启用状态',
    priority            INTEGER DEFAULT 0 COMMENT '优先级',
    
    -- 时间戳
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- 索引
    INDEX idx_enabled (enabled),
    INDEX idx_priority (priority),
    UNIQUE KEY uk_name (name)
);
```

### 3. 节点表 (nodes)

存储节点的详细信息和性能数据。

```sql
CREATE TABLE nodes (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    subscription_id     INTEGER NOT NULL COMMENT '所属订阅ID',
    node_pool_id        INTEGER NULL COMMENT '所属节点池ID',
    
    -- 节点基础信息
    name                VARCHAR(255) NOT NULL COMMENT '节点名称',
    hash                VARCHAR(64) NOT NULL COMMENT '节点哈希(用于去重)',
    server              VARCHAR(255) NOT NULL COMMENT '服务器地址',
    port                INTEGER NOT NULL COMMENT '端口',
    protocol            VARCHAR(50) NOT NULL COMMENT '协议类型: ss/ssr/vmess/trojan/hysteria',
    
    -- clash 配置 (JSON 存储)
    clash_config        JSON NOT NULL COMMENT 'Clash 完整配置',
    
    -- 地理信息
    country             VARCHAR(10) NULL COMMENT '国家代码',
    country_name        VARCHAR(100) NULL COMMENT '国家名称',
    city                VARCHAR(100) NULL COMMENT '城市',
    isp                 VARCHAR(100) NULL COMMENT '运营商',
    
    -- 性能指标 (最新测试结果)
    delay               INTEGER NULL COMMENT '延迟(ms)',
    upload_speed        BIGINT NULL COMMENT '上传速度(bps)',
    download_speed      BIGINT NULL COMMENT '下载速度(bps)',
    loss_rate           DECIMAL(5,2) NULL COMMENT '丢包率(%)',
    
    -- 连通性状态
    status              VARCHAR(50) DEFAULT 'unknown' COMMENT '状态: online/offline/testing/unknown',
    last_test           TIMESTAMP NULL COMMENT '最后测试时间',
    last_online         TIMESTAMP NULL COMMENT '最后在线时间',
    continuous_failures INTEGER DEFAULT 0 COMMENT '连续失败次数',
    
    -- 流媒体解锁信息 (JSON 存储)
    streaming_unlock    JSON NULL COMMENT '流媒体解锁情况',
    
    -- 时间戳
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- 索引优化
    INDEX idx_subscription_id (subscription_id),
    INDEX idx_node_pool_id (node_pool_id),
    INDEX idx_hash (hash),
    INDEX idx_status (status),
    INDEX idx_delay (delay),
    INDEX idx_last_test (last_test),
    INDEX idx_protocol (protocol),
    INDEX idx_country (country),
    INDEX idx_server_port (server, port),
    
    -- 外键约束
    FOREIGN KEY (subscription_id) REFERENCES subscriptions(id) ON DELETE CASCADE,
    FOREIGN KEY (node_pool_id) REFERENCES node_pools(id) ON DELETE SET NULL,
    
    -- 唯一约束 (防止重复节点)
    UNIQUE KEY uk_hash (hash)
);
```

### 4. 节点测试记录表 (node_tests)

历史测试数据，用于趋势分析和智能选择。

```sql
CREATE TABLE node_tests (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    node_id             INTEGER NOT NULL COMMENT '节点ID',
    
    -- 测试类型和配置
    test_type           VARCHAR(50) NOT NULL COMMENT '测试类型: delay/speed/streaming',
    test_config         JSON NULL COMMENT '测试配置参数',
    
    -- 测试结果
    delay               INTEGER NULL COMMENT '延迟(ms)',
    upload_speed        BIGINT NULL COMMENT '上传速度(bps)',
    download_speed      BIGINT NULL COMMENT '下载速度(bps)',
    loss_rate           DECIMAL(5,2) NULL COMMENT '丢包率(%)',
    
    -- 流媒体解锁结果
    streaming_results   JSON NULL COMMENT '流媒体测试结果',
    
    -- 测试状态
    success             BOOLEAN NOT NULL DEFAULT FALSE COMMENT '测试是否成功',
    error_message       TEXT NULL COMMENT '错误信息',
    
    -- 时间戳
    tested_at           TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- 索引
    INDEX idx_node_id (node_id),
    INDEX idx_test_type (test_type),
    INDEX idx_tested_at (tested_at),
    INDEX idx_success (success),
    
    -- 外键约束
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE CASCADE
);
```

### 5. 订阅更新日志表 (subscription_logs)

记录订阅更新历史，便于问题追踪和统计分析。

```sql
CREATE TABLE subscription_logs (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    subscription_id     INTEGER NOT NULL COMMENT '订阅ID',
    
    -- 更新信息
    update_type         VARCHAR(50) NOT NULL COMMENT '更新类型: auto/manual/retry',
    
    -- 更新结果
    success             BOOLEAN NOT NULL DEFAULT FALSE COMMENT '更新是否成功',
    
    -- 节点统计
    total_fetched       INTEGER DEFAULT 0 COMMENT '获取到的节点总数',
    valid_nodes         INTEGER DEFAULT 0 COMMENT '有效节点数',
    new_nodes           INTEGER DEFAULT 0 COMMENT '新增节点数(订阅内去重)',
    global_new_nodes    INTEGER DEFAULT 0 COMMENT '全局新增节点数',
    updated_nodes       INTEGER DEFAULT 0 COMMENT '更新的节点数',
    removed_nodes       INTEGER DEFAULT 0 COMMENT '移除的节点数',
    
    -- 错误信息
    error_message       TEXT NULL COMMENT '错误信息',
    http_status         INTEGER NULL COMMENT 'HTTP 状态码',
    response_time       INTEGER NULL COMMENT '响应时间(ms)',
    
    -- 时间戳
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- 索引
    INDEX idx_subscription_id (subscription_id),
    INDEX idx_success (success),
    INDEX idx_update_type (update_type),
    INDEX idx_created_at (created_at),
    
    -- 外键约束
    FOREIGN KEY (subscription_id) REFERENCES subscriptions(id) ON DELETE CASCADE
);
```

### 6. 节点池关联表 (node_pool_subscriptions)

多对多关系表，支持一个订阅属于多个节点池。

```sql
CREATE TABLE node_pool_subscriptions (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    node_pool_id        INTEGER NOT NULL COMMENT '节点池ID',
    subscription_id     INTEGER NOT NULL COMMENT '订阅ID',
    
    -- 配置信息
    enabled             BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    priority            INTEGER DEFAULT 0 COMMENT '优先级',
    
    -- 时间戳
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- 索引
    INDEX idx_node_pool_id (node_pool_id),
    INDEX idx_subscription_id (subscription_id),
    INDEX idx_enabled (enabled),
    
    -- 外键约束
    FOREIGN KEY (node_pool_id) REFERENCES node_pools(id) ON DELETE CASCADE,
    FOREIGN KEY (subscription_id) REFERENCES subscriptions(id) ON DELETE CASCADE,
    
    -- 唯一约束
    UNIQUE KEY uk_pool_subscription (node_pool_id, subscription_id)
);
```

## 数据流媒体解锁字段设计

```json
{
  "netflix": {
    "available": true,
    "region": "US",
    "tested_at": "2024-01-15T10:00:00Z"
  },
  "youtube": {
    "available": true,
    "region": "US",
    "tested_at": "2024-01-15T10:00:00Z"
  },
  "disney_plus": {
    "available": false,
    "region": null,
    "tested_at": "2024-01-15T10:00:00Z"
  },
  "hulu": {
    "available": true,
    "region": "US",
    "tested_at": "2024-01-15T10:00:00Z"
  },
  "amazon_prime": {
    "available": true,
    "region": "US",
    "tested_at": "2024-01-15T10:00:00Z"
  },
  "chatgpt": {
    "available": true,
    "tested_at": "2024-01-15T10:00:00Z"
  }
}
```

## Clash 配置字段设计

```json
{
  "name": "🇺🇸 US-01",
  "type": "vmess",
  "server": "us1.example.com",
  "port": 443,
  "uuid": "12345678-1234-1234-1234-123456789abc",
  "alterId": 0,
  "cipher": "auto",
  "tls": true,
  "network": "ws",
  "ws-opts": {
    "path": "/path",
    "headers": {
      "Host": "us1.example.com"
    }
  }
}
```

## 性能优化策略

### 1. 索引策略
- **主键索引**: 所有表使用自增主键
- **外键索引**: 所有外键列都有对应索引
- **查询索引**: 基于常用查询条件建立复合索引
- **时间索引**: 时间字段建立索引支持范围查询

### 2. 数据分区 (MySQL/PostgreSQL)
```sql
-- 按月分区节点测试记录表
PARTITION BY RANGE (YEAR(tested_at)*100 + MONTH(tested_at)) (
    PARTITION p202401 VALUES LESS THAN (202402),
    PARTITION p202402 VALUES LESS THAN (202403),
    -- ...
    PARTITION p999999 VALUES LESS THAN MAXVALUE
);
```

### 3. 缓存策略
- **节点状态缓存**: Redis 缓存最新节点状态
- **统计数据缓存**: 缓存存活率等计算密集的统计数据
- **配置缓存**: 缓存订阅配置和节点池配置

### 4. 读写分离
- **主从复制**: MySQL/PostgreSQL 支持读写分离
- **查询路由**: 统计查询路由到从库
- **写入队列**: 异步写入测试数据

### 5. 数据归档
- **历史数据**: 定期归档老旧测试记录
- **日志清理**: 自动清理过期日志数据
- **备份策略**: 定期全量和增量备份

## 事务处理策略

### 1. 订阅更新事务
```sql
BEGIN TRANSACTION;
-- 更新订阅信息
-- 批量插入/更新节点
-- 更新统计数据
-- 记录更新日志
COMMIT;
```

### 2. 节点测试事务
```sql
BEGIN TRANSACTION;
-- 插入测试记录
-- 更新节点状态
-- 更新统计数据
COMMIT;
```

### 3. 数据一致性
- **外键约束**: 确保关联数据完整性
- **唯一约束**: 防止重复数据
- **检查约束**: 确保数据有效性
- **触发器**: 自动更新统计数据

## 扩展性考虑

### 1. 水平扩展
- **分库分表**: 按订阅或地区分片
- **负载均衡**: 数据库连接池和负载均衡
- **缓存集群**: Redis 集群支持

### 2. 垂直扩展  
- **硬件升级**: CPU、内存、存储升级
- **配置优化**: 数据库参数调优
- **索引优化**: 定期分析和优化索引

### 3. 监控告警
- **性能监控**: 查询响应时间、连接数
- **容量监控**: 磁盘使用率、表大小
- **异常告警**: 错误率、超时告警

这个数据库设计支持高并发读写、复杂统计查询，并为后续的智能节点选择算法提供了丰富的数据基础。