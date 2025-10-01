# Message Queue

消息队列抽象层，支持多种消息队列后端。使用泛型设计，类型安全。

## 支持的队列类型

- **Memory**: 内存队列（开发/测试用）✅ 已实现
- **Redis**: 基于 Redis Streams 的消息队列 🚧 待实现
- **NSQ**: 分布式实时消息平台 🚧 待实现
- **Kafka**: 高吞吐量分布式消息系统 🚧 待实现
- **RabbitMQ**: 可靠的消息代理 🚧 待实现
- **ZeroMQ**: 高性能异步消息库 🚧 待实现

## 核心接口

```go
type Queue[T any] interface {
    Push(payload T) error
    BatchPush(payloads []T) error
    Pop() (*Message[T], error)
    Process(concurrency int, handler func(*Message[T]) error) error
    Depth() int
    Close() error
}
```

## 使用示例

### 创建队列

```go
import "github.com/ice-cream-heaven/prism/pkg/queue"

// 方式一：使用默认配置
cfg := &queue.Config{
    Type: queue.TypeMemory,
}

q, err := queue.New[string](cfg)
if err != nil {
    log.Fatal(err)
}
defer q.Close()

// 方式二：完整配置
cfg := &queue.Config{
    Type:     queue.TypeMemory,
    MaxRetry: 3,
    Timeout:  30,
    MaxDepth: 10000,
}
cfg.ApplyDefaults() // 应用默认值

q, err := queue.New[MyStruct](cfg)
```

### 推送消息

```go
// 推送单个消息
err := q.Push("Hello World")

// 批量推送
messages := []string{"msg1", "msg2", "msg3"}
err := q.BatchPush(messages)
```

### 消费消息

```go
// 方式一：使用 Pop（手动拉取）
for {
    msg, err := q.Pop()
    if err == queue.ErrQueueEmpty {
        time.Sleep(100 * time.Millisecond)
        continue
    }
    if err != nil {
        log.Printf("Error: %v", err)
        break
    }

    // 处理消息
    fmt.Printf("Received: %s\n", msg.Payload)
}

// 方式二：使用 Process（自动消费）
err := q.Process(5, func(msg *queue.Message[string]) error {
    fmt.Printf("Processing: %s\n", msg.Payload)
    // 返回 error 会触发重试
    return nil
})
```

### 查询队列深度

```go
depth := q.Depth()
fmt.Printf("Queue depth: %d\n", depth)
```

## 泛型使用

队列支持任意类型的消息：

```go
// 字符串队列
strQueue, _ := queue.New[string](cfg)
strQueue.Push("hello")

// 整数队列
intQueue, _ := queue.New[int](cfg)
intQueue.Push(42)

// 结构体队列
type Task struct {
    ID   string
    Name string
}

taskQueue, _ := queue.New[Task](cfg)
taskQueue.Push(Task{ID: "1", Name: "task1"})

// 指针队列
ptrQueue, _ := queue.New[*Task](cfg)
ptrQueue.Push(&Task{ID: "1", Name: "task1"})
```

## 配置说明

### Config 结构

```go
type Config struct {
    Type     Type                   // 队列类型（必填）
    Address  string                 // 连接地址
    MaxRetry int                    // 最大重试次数（默认：3）
    Timeout  int                    // 操作超时（秒，默认：30）
    MaxDepth int                    // 队列最大深度（默认：10000）
    Options  map[string]interface{} // 类型特定配置
}
```

### 配置示例（YAML）

```yaml
queue:
  type: memory
  max_retry: 3
  timeout: 30
  max_depth: 10000
```

```yaml
queue:
  type: redis
  address: localhost:6379
  max_retry: 5
  timeout: 60
  max_depth: 50000
  options:
    password: secret
    db: 0
```

### 配置示例（JSON）

```json
{
  "type": "memory",
  "max_retry": 3,
  "timeout": 30,
  "max_depth": 10000
}
```

## 队列类型常量

```go
const (
    TypeMemory   Type = "memory"
    TypeRedis    Type = "redis"
    TypeNSQ      Type = "nsq"
    TypeKafka    Type = "kafka"
    TypeRabbitMQ Type = "rabbitmq"
    TypeZeroMQ   Type = "zeromq"
)
```

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
)
```

## 完整示例

### 示例1：订阅更新任务队列

```go
package main

import (
    "fmt"
    "log"

    "github.com/ice-cream-heaven/prism/pkg/queue"
)

type SubscriptionTask struct {
    SubscriptionID string
    Action         string
}

func main() {
    // 创建队列
    cfg := &queue.Config{
        Type:     queue.TypeMemory,
        MaxDepth: 1000,
    }
    q, err := queue.New[SubscriptionTask](cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer q.Close()

    // 推送任务
    tasks := []SubscriptionTask{
        {SubscriptionID: "sub1", Action: "update"},
        {SubscriptionID: "sub2", Action: "update"},
    }
    if err := q.BatchPush(tasks); err != nil {
        log.Fatal(err)
    }

    // 启动消费者（5个并发）
    err = q.Process(5, func(msg *queue.Message[SubscriptionTask]) error {
        fmt.Printf("Processing subscription %s: %s\n",
            msg.Payload.SubscriptionID,
            msg.Payload.Action)
        // 执行更新逻辑
        return nil
    })

    if err != nil {
        log.Fatal(err)
    }

    // 等待处理完成
    select {}
}
```

### 示例2：节点测试任务队列

```go
type NodeTestTask struct {
    NodeID   string
    Priority int
    Timeout  int
}

func main() {
    cfg := &queue.Config{
        Type:     queue.TypeMemory,
        MaxRetry: 3,
    }
    q, err := queue.New[NodeTestTask](cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer q.Close()

    // 推送节点测试任务
    q.Push(NodeTestTask{
        NodeID:   "node1",
        Priority: 1,
        Timeout:  30,
    })

    // 处理任务
    q.Process(10, func(msg *queue.Message[NodeTestTask]) error {
        task := msg.Payload
        fmt.Printf("Testing node %s (retry: %d)\n",
            task.NodeID, msg.RetryCount)

        // 执行测试逻辑
        // 如果失败，返回 error 会触发重试
        return testNode(task.NodeID)
    })

    select {}
}
```

## 性能考虑

### 内存队列

- ✅ 零延迟
- ✅ 零序列化开销
- ❌ 不支持持久化
- ❌ 仅单机

### 其他队列（待实现）

根据实际需求选择：

- **轻量级分布式**: Redis
- **中等规模**: NSQ
- **大规模/高吞吐**: Kafka
- **企业级/可靠性**: RabbitMQ
- **高性能/低延迟**: ZeroMQ

## 注意事项

1. **类型安全**: 使用泛型确保编译时类型检查
2. **并发安全**: 所有实现都是并发安全的
3. **优雅关闭**: 使用 `Close()` 确保资源正确释放
4. **重试机制**: 失败的消息会自动重试，直到达到 `MaxRetry`
5. **队列满**: Push 操作在队列满时会返回 `ErrQueueFull`
6. **队列空**: Pop 操作在队列空时会返回 `ErrQueueEmpty`

## 开发指南

### 添加新的队列类型

1. 在 `types.go` 中添加新的 Type 常量
2. 创建新的实现文件（如 `pulsar.go`）
3. 实现 `Queue[T]` 接口
4. 在 `queue.go` 的 `New()` 函数中添加 case
5. 更新文档

### 运行测试

```bash
go test -v ./pkg/queue/...
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
```

## 参考资料

- [Redis Streams](https://redis.io/docs/data-types/streams/)
- [NSQ Documentation](https://nsq.io/)
- [Apache Kafka](https://kafka.apache.org/documentation/)
- [RabbitMQ Tutorials](https://www.rabbitmq.com/getstarted.html)
- [ZeroMQ Guide](https://zguide.zeromq.org/)
