-- 队列配置表
CREATE TABLE IF NOT EXISTS queue_configs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    name VARCHAR(100) NOT NULL UNIQUE COMMENT '队列名称（唯一标识）',
    type VARCHAR(20) NOT NULL COMMENT '队列类型：memory|redis|nsq|kafka|rabbitmq|zeromq',
    address VARCHAR(500) DEFAULT '' COMMENT '队列地址（memory类型可为空）',
    timeout INT DEFAULT 30 COMMENT '超时时间（秒）',
    max_depth INT DEFAULT 10000 COMMENT '最大队列深度（0表示无限制）',
    enabled BOOLEAN DEFAULT TRUE COMMENT '是否启用',

    -- 重试策略配置（JSON格式）
    retry_policy TEXT COMMENT '重试策略配置（JSON）：{"strategy":"exponential","max_retries":3,"initial_delay":1000,"max_delay":60000,"multiplier":2.0}',

    -- 延时队列配置
    enable_delayed_queue BOOLEAN DEFAULT FALSE COMMENT '是否启用延时队列',
    delay_check_interval INT DEFAULT 1000 COMMENT '延时检查间隔（毫秒）',

    -- 扩展配置（JSON格式）
    options TEXT COMMENT '扩展配置（JSON）',

    -- 元信息
    description VARCHAR(500) DEFAULT '' COMMENT '队列描述',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    INDEX idx_type (type),
    INDEX idx_enabled (enabled),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='队列配置表';

-- 示例数据
INSERT INTO queue_configs (name, type, address, timeout, max_depth, enabled, retry_policy, enable_delayed_queue, description) VALUES
('default_memory', 'memory', '', 30, 10000, TRUE, '{"strategy":"exponential","max_retries":3,"initial_delay":1000,"max_delay":60000,"multiplier":2.0}', FALSE, '默认内存队列'),
('task_queue', 'redis', 'redis://localhost:6379/0', 60, 50000, TRUE, '{"strategy":"exponential","max_retries":5,"initial_delay":2000,"max_delay":120000,"multiplier":2.0}', TRUE, '任务队列（Redis Streams）'),
('event_queue', 'nsq', 'nsqlookupd://localhost:4161', 30, 0, TRUE, '{"strategy":"fixed","max_retries":3,"initial_delay":5000}', FALSE, '事件队列（NSQ）');
