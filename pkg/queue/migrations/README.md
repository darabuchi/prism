# 队列配置数据库迁移

## 数据库表结构

### queue_configs 表

队列配置表用于存储所有队列的配置信息。

| 字段 | 类型 | 说明 |
|-----|------|------|
| id | BIGINT | 主键ID |
| name | VARCHAR(100) | 队列名称（唯一标识） |
| type | VARCHAR(20) | 队列类型：memory\|redis\|nsq\|kafka\|rabbitmq\|zeromq |
| address | VARCHAR(500) | 队列地址（memory类型可为空） |
| timeout | INT | 超时时间（秒），默认30 |
| max_depth | INT | 最大队列深度，默认10000，0表示无限制 |
| enabled | BOOLEAN | 是否启用，默认TRUE |
| retry_policy | TEXT | 重试策略配置（JSON格式） |
| enable_delayed_queue | BOOLEAN | 是否启用延时队列，默认FALSE |
| delay_check_interval | INT | 延时检查间隔（毫秒），默认1000 |
| options | TEXT | 扩展配置（JSON格式） |
| description | VARCHAR(500) | 队列描述 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

### 索引

- `name`: 唯一索引
- `type`: 普通索引
- `enabled`: 普通索引
- `created_at`: 普通索引

## retry_policy JSON 格式

```json
{
  "strategy": "exponential",
  "max_retries": 3,
  "initial_delay": 1000,
  "max_delay": 60000,
  "multiplier": 2.0,
  "increment": 1000
}
```

### 重试策略类型

- `none`: 不重试
- `fixed`: 固定间隔重试
- `linear`: 线性递增重试
- `exponential`: 指数退避重试（推荐）

## options JSON 格式

扩展配置可以存储任意键值对，用于特定队列类型的额外配置：

```json
{
  "pool_size": 10,
  "tls_enabled": true,
  "custom_setting": "value"
}
```

## 使用示例

### 1. 创建配置

```sql
INSERT INTO queue_configs (
    name, type, address, timeout, max_depth, enabled,
    retry_policy, enable_delayed_queue, description
) VALUES (
    'task_queue',
    'redis',
    'redis://localhost:6379/0',
    60,
    50000,
    TRUE,
    '{"strategy":"exponential","max_retries":5,"initial_delay":2000,"max_delay":120000,"multiplier":2.0}',
    TRUE,
    '任务队列'
);
```

### 2. 查询配置

```sql
-- 获取所有启用的配置
SELECT * FROM queue_configs WHERE enabled = TRUE;

-- 获取指定类型的配置
SELECT * FROM queue_configs WHERE type = 'redis' AND enabled = TRUE;

-- 获取指定名称的配置
SELECT * FROM queue_configs WHERE name = 'task_queue';
```

### 3. 更新配置

```sql
-- 更新重试策略
UPDATE queue_configs
SET retry_policy = '{"strategy":"fixed","max_retries":3,"initial_delay":5000}',
    updated_at = CURRENT_TIMESTAMP
WHERE name = 'task_queue';

-- 禁用队列
UPDATE queue_configs
SET enabled = FALSE,
    updated_at = CURRENT_TIMESTAMP
WHERE name = 'task_queue';
```

### 4. 删除配置

```sql
DELETE FROM queue_configs WHERE name = 'task_queue';
```

## 在代码中使用

```go
import (
    "context"
    "github.com/darabuchi/prism/pkg/queue"
)

// 创建配置加载器
loader := queue.NewDBConfigLoader(repo)

// 加载单个队列配置
cfg, err := loader.LoadConfig(context.Background(), "task_queue")
if err != nil {
    // 处理错误
}

// 创建队列实例
q, err := queue.NewQueue[MyTask](cfg)
if err != nil {
    // 处理错误
}

// 加载所有配置
allConfigs, err := loader.LoadAllConfigs(context.Background())
for name, cfg := range allConfigs {
    // 为每个配置创建队列实例
}
```

## 迁移说明

1. 执行 `001_create_queue_configs.sql` 创建表结构
2. 根据需要插入初始配置数据
3. 在应用程序中使用 `ConfigLoader` 加载配置
