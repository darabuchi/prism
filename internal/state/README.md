# State 状态管理模块

## 概述

`internal/state` 模块负责全局状态管理，包括配置加载、数据库连接、缓存连接、多语言资源和队列系统的初始化。

## 架构设计

```
┌─────────────────────────────────────┐
│           State (全局状态)           │
├─────────────────────────────────────┤
│  - Config  (配置)                   │
│  - I18n    (多语言)                  │
└─────────────────────────────────────┘
         ↓
┌─────────────────────────────────────┐
│     Load() 初始化流程                │
├─────────────────────────────────────┤
│  1. LoadConfig()      加载配置      │
│  2. LoadI18n()        加载多语言     │
│  3. ConnectDatabase() 连接数据库     │
│  4. ConnectCache()    连接缓存       │
│  5. InitQueue()       初始化队列     │
└─────────────────────────────────────┘
```

## 文件说明

| 文件 | 说明 |
|------|------|
| `state.go` | 全局状态定义和 Load() 初始化函数 |
| `config.go` | 配置结构和加载逻辑 |
| `i18n.go` | 多语言资源加载（支持 7 种语言） |
| `table.go` | 数据库表管理和 Model 定义 |
| `cache.go` | 缓存连接管理 |
| `queue.go` | 队列系统初始化 |

## 快速开始

### 初始化

在应用启动时调用 `Load()` 函数初始化所有状态：

```go
package main

import (
    "github.com/darabuchi/prism/internal/state"
    "github.com/lazygophers/log"
)

func main() {
    // 初始化全局状态
    if err := state.Load(); err != nil {
        log.Fatalf("failed to load state: %v", err)
    }

    // 现在可以使用全局状态
    // ...
}
```

### 访问配置

```go
// 访问数据库配置
dbConfig := state.State.Config.Db

// 访问缓存配置
cacheConfig := state.State.Config.Cache
```

### 访问数据库

```go
// 获取数据库客户端
db := state.Db()

// 创建新的操作作用域
scoop := state.NewScoop()

// 使用事务
err := state.CommitOrRollback(func(tx *db.Scoop) error {
    // 事务逻辑
    return nil
})
```

### 访问缓存

```go
// 获取缓存客户端
cache := state.Cache()

// 设置值
err := cache.Set(ctx, "key", "value", time.Hour)

// 获取值
val, err := cache.Get(ctx, "key")
```

### 使用多语言

```go
// 使用当前上下文语言
msg := state.Localize("common.success")

// 使用指定语言
msg := state.LocalizeWithLang("zh-CN", "common.internal_error")
```

## 配置文件

配置文件支持 YAML、JSON、TOML 格式，默认查找 `config.yaml`：

```yaml
# config.yaml
db:
  type: sqlite
  address: resource
  name: prism
  username: prism
  password: prism2025
  debug: false

cache:
  type: memory
  address: ""
  port: 0
  db: 0
  data_dir: resource/cache/
```

### 配置字段说明

#### 数据库配置 (db)

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| type | string | 数据库类型 (sqlite/mysql/postgres) | sqlite |
| address | string | 数据库地址或文件路径 | resource |
| port | int | 数据库端口 | 0 |
| name | string | 数据库名称 | prism |
| username | string | 用户名 | prism |
| password | string | 密码 | prism2025 |
| debug | bool | 是否开启调试模式 | false |

#### 缓存配置 (cache)

| 字段 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| type | string | 缓存类型 (memory/redis) | memory |
| address | string | 缓存服务器地址 | "" |
| port | int | 缓存服务器端口 | 0 |
| password | string | 密码 | "" |
| db | int | Redis 数据库编号 | 0 |
| data_dir | string | 数据目录 | resource/cache/ |

## 多语言支持

内置 **7 种语言**支持：

- `en` - English (默认)
- `zh-CN` - 简体中文
- `zh` - 中文（回退）
- `zh-TW` - 繁體中文
- `fr` - Français (法语)
- `ru` - Русский (俄语)
- `es` - Español (西班牙语)
- `ar` - العربية (阿拉伯语)

多语言资源文件位于 `resource/localize/` 目录，使用 `embed.FS` 嵌入到二进制中。

## 数据库 Model

TODO: 添加 Model 定义后更新此部分

```go
// 示例：定义订阅 Model
type ModelSubscription struct {
    ID    uint64 `gorm:"primaryKey"`
    Title string
    URL   string
    // ...
}

// 在 table.go 中初始化
Subscription = db.NewModel[ModelSubscription](Db()).
    SetNotFound(xerror.NewError(3000)).
    SetDuplicatedKeyError(xerror.NewError(3001))
```

## 队列系统

TODO: 添加队列定义后更新此部分

```go
// 示例：定义订阅更新任务队列
QueueSubscriptionUpdate = queue.NewMemoryQueue[*TaskSubscriptionUpdate](&queue.Config{
    MaxSize:            500,
    ConcurrentWorkers:  5,
    MaxRetries:         3,
    RetryDelayBase:     time.Second,
})
```

## 错误处理

所有初始化函数都会返回错误，建议在应用启动时处理：

```go
if err := state.Load(); err != nil {
    log.Fatalf("initialization failed: %v", err)
}
```

初始化过程中会使用 `pterm` 显示友好的进度提示：

```
✓ 配置加载成功
✓ 多语言资源加载成功
✓ 数据库连接成功
✓ 缓存连接成功
✓ 队列启动成功
```

## 最佳实践

1. **单次初始化** - 在 `main()` 函数开始时调用 `state.Load()`，全局只初始化一次
2. **避免直接访问** - 使用提供的 `Db()`, `Cache()` 等函数访问资源
3. **事务管理** - 使用 `CommitOrRollback()` 自动处理事务提交和回滚
4. **配置验证** - 配置加载会自动验证必填字段，确保配置完整性
5. **资源清理** - 使用 `atexit` 自动注册资源清理逻辑

## 相关文档

- [错误码说明](../../ERROR_CODES.md)
- [多语言资源](../../resource/localize/README.md)
- [队列系统](../../pkg/queue/README.md)
