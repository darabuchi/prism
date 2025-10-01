# Message Queue

消息队列抽象层，支持多种消息队列后端。使用泛型设计，类型安全，支持延时执行和重试策略。

## 支持的队列类型

- **Memory**: 内存队列（开发/测试用）✅ 已实现
- **Redis**: 基于 Redis Streams 的消息队列 🚧 待实现
- **NSQ**: 分布式实时消息平台 🚧 待实现
- **Kafka**: 高吞吐量分布式消息系统 🚧 待实现
- **RabbitMQ**: 可靠的消息代理 🚧 待实现
- **ZeroMQ**: 高性能异步消息库 🚧 待实现

## 核心特性

✨ **泛型支持** - 类型安全的消息队列
🔄 **智能重试** - 支持 4 种重试策略（无重试、固定间隔、线性递增、指数退避）
⏱️ **延时执行** - 支持延时消息和延时队列
📊 **监控指标** - 内置指标收集（已处理、失败、重试、丢弃、活跃任务）
🔒 **并发安全** - 所有实现都是并发安全的
⚡ **高性能** - 优化的 O(n) 复杂度延时队列处理器
🛡️ **优雅关闭** - 使用 channel 关闭机制确保 worker 正常退出

## 核心接口

```go
type Queue[T any] interface {
    // 推送消息
    Push(payload T) error
    PushWithDelay(payload T, delay time.Duration) error
    PushWithRetryPolicy(payload T, policy *RetryPolicy) error
    BatchPush(payloads []T) error

    // 消费消息
    Pop() (*Message[T], error)
    Process(concurrency int, handler func(*Message[T]) *HandlerResult) error

    // 队列信息
    Depth() int
    Close() error
}
```

## 快速开始

### 创建队列

```go
import "github.com/darabuchi/prism/pkg/queue"

// 创建内存队列
cfg := &queue.Config{
    Type:     queue.TypeMemory,
    MaxDepth: 10000,
    RetryPolicy: queue.ExponentialRetryPolicy(3, 1000, 2.0),
}

q, err := queue.New[string](cfg)
if err != nil {
    log.Fatal(err)
}
defer q.Close()
```

### 推送消息

```go
// 推送单个消息
err := q.Push("Hello World")

// 推送延时消息（需启用延时队列）
err := q.PushWithDelay("Delayed message", 5*time.Second)

// 推送带自定义重试策略的消息
policy := queue.FixedRetryPolicy(5, 3000) // 重试5次，每次延时3秒
err := q.PushWithRetryPolicy("Important task", policy)

// 批量推送
messages := []string{"msg1", "msg2", "msg3"}
err := q.BatchPush(messages)
```

### 消费消息

```go
// 使用 Process 启动消费者（推荐）
err := q.Process(5, func(msg *queue.Message[string]) *queue.HandlerResult {
    fmt.Printf("Processing: %s\n", msg.Payload)

    // 处理成功
    if processSuccess() {
        return queue.Success()
    }

    // 处理失败，但不重试
    if shouldNotRetry() {
        return queue.Fail(errors.New("permanent failure"))
    }

    // 处理失败，需要重试
    if needRetry() {
        return queue.RetryWithError(errors.New("temporary failure"))
    }

    // 处理失败，使用自定义延迟重试
    return queue.RetryWithDelay(10*time.Second, errors.New("retry after 10s"))
})
```

### 查询指标

```go
metrics := q.(*queue.MemoryQueue[string]).Metrics()
fmt.Printf("Processed: %d\n", metrics["processed"])
fmt.Printf("Failed: %d\n", metrics["failed"])
fmt.Printf("Retried: %d\n", metrics["retried"])
fmt.Printf("Dropped: %d\n", metrics["dropped"])
fmt.Printf("Active tasks: %d\n", metrics["active_tasks"])
fmt.Printf("Queue depth: %d\n", metrics["queue_depth"])
```

## 配置说明

### Config 结构

```go
type Config struct {
    Type               Type                   // 队列类型（必填）
    Address            string                 // 连接地址（TypeMemory 可选）
    Timeout            int                    // 操作超时（秒，默认：30）
    MaxDepth           int                    // 队列最大深度（默认：10000，0=无限制）
    RetryPolicy        *RetryPolicy           // 默认重试策略
    EnableDelayedQueue bool                   // 是否启用延时队列（默认：false）
    DelayCheckInterval int                    // 延时检查间隔（毫秒，默认：1000）
    Options            map[string]interface{} // 类型特定配置
}
```

### 重试策略

```go
type RetryPolicy struct {
    Strategy     RetryStrategy // 重试策略类型
    MaxRetries   int           // 最大重试次数（0=无限制）
    InitialDelay int           // 初始延迟（毫秒）
    MaxDelay     int           // 最大延迟（毫秒，0=无限制）
    Multiplier   float64       // 延迟倍数（指数退避）
    Increment    int           // 增量（线性递增，毫秒）
}
```

#### 重试策略类型

1. **None** - 不重试
```go
policy := queue.NoRetryPolicy()
```

2. **Fixed** - 固定间隔重试
```go
// 重试3次，每次间隔5秒
policy := queue.FixedRetryPolicy(3, 5000)
```

3. **Linear** - 线性递增重试
```go
// 重试5次，初始延迟1秒，每次增加2秒
policy := queue.LinearRetryPolicy(5, 1000, 2000)
// 重试间隔: 1s, 3s, 5s, 7s, 9s
```

4. **Exponential** - 指数退避重试（推荐）
```go
// 重试3次，初始延迟1秒，倍数2.0
policy := queue.ExponentialRetryPolicy(3, 1000, 2.0)
// 重试间隔: 1s, 2s, 4s
```

## 延时队列

延时队列允许消息在未来某个时间点执行：

```go
cfg := &queue.Config{
    Type:               queue.TypeMemory,
    EnableDelayedQueue: true,
    DelayCheckInterval: 500, // 每500ms检查一次
}

q, _ := queue.New[Task](cfg)

// 5秒后执行
q.PushWithDelay(Task{ID: "1"}, 5*time.Second)

// 自定义延迟重试
q.Process(3, func(msg *queue.Message[Task]) *queue.HandlerResult {
    if needRetryAfter10s() {
        return queue.RetryWithDelay(10*time.Second, err)
    }
    return queue.Success()
})
```

## Handler 返回结果

Handler 必须返回 `*HandlerResult` 明确指示处理结果：

```go
type HandlerResult struct {
    Retry      bool          // 是否需要重试
    Error      error         // 错误信息（可选）
    RetryDelay time.Duration // 自定义重试延迟（可选）
}
```

辅助函数：

```go
// 成功
queue.Success()

// 失败但不重试
queue.Fail(err)

// 失败需要重试（使用策略计算延迟）
queue.RetryWithError(err)

// 失败需要重试（自定义延迟）
queue.RetryWithDelay(5*time.Second, err)
```

## 数据库配置

队列配置可以存储在数据库中：

```go
// 创建配置
qc := &queue.QueueConfig{
    Name:               "task_queue",
    Type:               queue.TypeRedis,
    Address:            "redis://localhost:6379",
    Timeout:            60,
    MaxDepth:           50000,
    Enabled:            true,
    EnableDelayedQueue: true,
}

qc.SetRetryPolicy(queue.ExponentialRetryPolicy(5, 2000, 2.0))

// 转换为运行时配置
cfg, err := qc.ToConfig()

// 创建队列
q, err := queue.New[Task](cfg)
```

## 完整示例

### 示例1：任务队列

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/darabuchi/prism/pkg/queue"
)

type Task struct {
    ID      string
    Payload string
}

func main() {
    // 创建队列
    cfg := &queue.Config{
        Type:     queue.TypeMemory,
        MaxDepth: 1000,
        RetryPolicy: queue.ExponentialRetryPolicy(3, 1000, 2.0),
        EnableDelayedQueue: true,
    }

    q, err := queue.New[Task](cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer q.Close()

    // 推送任务
    q.Push(Task{ID: "1", Payload: "immediate"})
    q.PushWithDelay(Task{ID: "2", Payload: "delayed"}, 5*time.Second)

    // 启动消费者
    q.Process(5, func(msg *queue.Message[Task]) *queue.HandlerResult {
        fmt.Printf("Processing task %s: %s (retry: %d)\n",
            msg.Payload.ID, msg.Payload.Payload, msg.RetryCount)

        // 模拟处理
        if err := processTask(msg.Payload); err != nil {
            return queue.RetryWithError(err)
        }

        return queue.Success()
    })

    // 运行一段时间
    time.Sleep(30 * time.Second)

    // 查看指标
    if mq, ok := q.(*queue.MemoryQueue[Task]); ok {
        metrics := mq.Metrics()
        fmt.Printf("Metrics: %+v\n", metrics)
    }
}

func processTask(task Task) error {
    // 处理逻辑
    return nil
}
```

### 示例2：订阅更新队列

```go
type SubscriptionUpdate struct {
    SubscriptionID string
    Action         string
    Data           map[string]interface{}
}

func main() {
    cfg := &queue.Config{
        Type:     queue.TypeMemory,
        MaxRetry: 5,
    }

    q, _ := queue.New[SubscriptionUpdate](cfg)
    defer q.Close()

    // 批量推送更新
    updates := []SubscriptionUpdate{
        {SubscriptionID: "sub1", Action: "update"},
        {SubscriptionID: "sub2", Action: "update"},
    }
    q.BatchPush(updates)

    // 处理更新
    q.Process(10, func(msg *queue.Message[SubscriptionUpdate]) *queue.HandlerResult {
        update := msg.Payload
        fmt.Printf("Updating subscription %s: %s\n",
            update.SubscriptionID, update.Action)

        if err := updateSubscription(update); err != nil {
            return queue.RetryWithError(err)
        }

        return queue.Success()
    })

    select {}
}

func updateSubscription(update SubscriptionUpdate) error {
    // 更新逻辑
    return nil
}
```

## 性能考虑

### 内存队列

- ✅ 零延迟
- ✅ 零序列化开销
- ✅ O(n) 延时队列处理器
- ✅ 内置监控指标
- ❌ 不支持持久化
- ❌ 仅单机

### 其他队列（待实现）

根据实际需求选择：

- **轻量级分布式**: Redis
- **中等规模**: NSQ
- **大规模/高吞吐**: Kafka
- **企业级/可靠性**: RabbitMQ
- **高性能/低延迟**: ZeroMQ

## 最佳实践

1. **使用泛型** - 确保类型安全
2. **合理设置重试策略** - 避免无限重试
3. **启用延时队列** - 对于需要延迟执行的场景
4. **监控指标** - 定期检查队列健康状况
5. **优雅关闭** - 确保 `Close()` 被调用
6. **错误处理** - 区分可重试和不可重试错误
7. **并发控制** - 根据负载调整 `concurrency` 参数
8. **日志记录** - 消息丢弃和超过重试次数会被记录

## 注意事项

1. **Handler 必须返回 HandlerResult** - 明确指示是否重试
2. **消息丢弃会被记录** - 队列满时会记录警告日志
3. **延时队列需显式启用** - 设置 `EnableDelayedQueue: true`
4. **Worker 通过关闭 channel 停止** - 确保优雅退出
5. **监控指标实时更新** - 使用 `Metrics()` 获取当前状态

## 错误处理

```go
var (
    ErrUnsupportedQueueType  // 不支持的队列类型
    ErrQueueClosed           // 队列已关闭
    ErrQueueEmpty            // 队列为空
    ErrQueueFull             // 队列已满
    ErrTimeout               // 操作超时
    ErrMaxRetriesExceeded    // 超过最大重试次数
    ErrInvalidConfig         // 无效的配置
    ErrConnectionFailed      // 连接失败
    ErrInvalidConcurrency    // 无效的并发数
    ErrQueueDisabled         // 队列已禁用
    ErrConfigNotFound        // 配置未找到
)
```

## 依赖

各队列后端的依赖（按需添加）：

```bash
# Redis
go get github.com/redis/go-redis/v9

# NSQ
go get github.com/nsqio/go-nsq

# Kafka
go get github.com/IBM/sarama

# RabbitMQ
go get github.com/rabbitmq/amqp091-go

# ZeroMQ
go get github.com/pebbe/zmq4

# Validator
go get github.com/darabuchi/lazygophers/utils/validator
```

## 参考资料

- [Redis Streams](https://redis.io/docs/data-types/streams/)
- [NSQ Documentation](https://nsq.io/)
- [Apache Kafka](https://kafka.apache.org/documentation/)
- [RabbitMQ Tutorials](https://www.rabbitmq.com/getstarted.html)
- [ZeroMQ Guide](https://zguide.zeromq.org/)
