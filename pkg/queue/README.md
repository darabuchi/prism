# Message Queue

消息队列抽象层，支持多种消息队列后端。

## 支持的队列类型

- **Memory**: 内存队列（开发/测试用）
- **Redis**: 基于 Redis Streams 的消息队列
- **NSQ**: 分布式实时消息平台
- **Kafka**: 高吞吐量分布式消息系统
- **RabbitMQ**: 可靠的消息代理
- **ZeroMQ**: 高性能异步消息库

## 使用示例

### 创建队列

```go
import "github.com/ice-cream-heaven/prism/pkg/queue"

// 内存队列
q, err := queue.New("memory")

// Redis 队列
q, err := queue.New("redis",
    queue.WithAddress("localhost:6379"),
    queue.WithMaxRetry(3),
)

// Kafka 队列
q, err := queue.New("kafka",
    queue.WithAddress("localhost:9092"),
    queue.WithOptions(map[string]interface{}{
        "group_id": "my-consumer-group",
    }),
)
```

### 发布消息

```go
// 简单发布
err := q.Publish(ctx, "my-topic", []byte("Hello World"))

// 带元数据发布
err := q.PublishWithMetadata(ctx, "my-topic",
    []byte("Hello World"),
    map[string]string{
        "priority": "high",
        "source": "api",
    },
)
```

### 订阅消息

```go
// 方式一：通过 Channel
msgChan, err := q.Subscribe(ctx, "my-topic")
if err != nil {
    log.Fatal(err)
}

for msg := range msgChan {
    log.Printf("Received: %s", string(msg.Payload))
    q.Ack(ctx, msg)
}

// 方式二：使用 Handler
err := q.Consume(ctx, "my-topic", func(msg *queue.Message) error {
    log.Printf("Received: %s", string(msg.Payload))
    // 返回 error 会触发重试
    return nil
})
```

### 消息确认

```go
// 成功处理，确认消息
q.Ack(ctx, msg)

// 处理失败，重新入队
q.Nack(ctx, msg)
```

## 配置选项

### Memory Queue

无需额外配置。

### Redis Queue

```go
queue.New("redis",
    queue.WithAddress("localhost:6379"),
    queue.WithOptions(map[string]interface{}{
        "password": "secret",
        "db":       0,
    }),
)
```

### NSQ Queue

```go
queue.New("nsq",
    queue.WithAddress("localhost:4150"),  // nsqd address
    queue.WithOptions(map[string]interface{}{
        "lookupd": []string{"localhost:4161"},  // nsqlookupd addresses
        "channel": "my-channel",
    }),
)
```

### Kafka Queue

```go
queue.New("kafka",
    queue.WithAddress("localhost:9092"),
    queue.WithOptions(map[string]interface{}{
        "group_id":     "my-group",
        "sasl_enabled": false,
    }),
)
```

### RabbitMQ Queue

```go
queue.New("rabbitmq",
    queue.WithAddress("amqp://guest:guest@localhost:5672/"),
    queue.WithOptions(map[string]interface{}{
        "exchange":      "my-exchange",
        "exchange_type": "topic",
    }),
)
```

### ZeroMQ Queue

```go
queue.New("zeromq",
    queue.WithAddress("tcp://localhost:5555"),
    queue.WithOptions(map[string]interface{}{
        "socket_type": "pub-sub",  // or "push-pull"
    }),
)
```

## 实现状态

| 队列类型 | 状态 | 说明 |
|---------|------|------|
| Memory  | ✅ 已实现 | 适用于开发和测试 |
| Redis   | 🚧 待实现 | 需要 github.com/redis/go-redis/v9 |
| NSQ     | 🚧 待实现 | 需要 github.com/nsqio/go-nsq |
| Kafka   | 🚧 待实现 | 需要 github.com/IBM/sarama |
| RabbitMQ | 🚧 待实现 | 需要 github.com/rabbitmq/amqp091-go |
| ZeroMQ  | 🚧 待实现 | 需要 github.com/pebbe/zmq4 |

## 架构设计

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │
       ├── Publish()
       ├── Subscribe()
       └── Consume()
       │
       v
┌─────────────┐
│   Queue     │ (Interface)
│  Interface  │
└──────┬──────┘
       │
       ├─────────────┬─────────────┬─────────────┐
       v             v             v             v
  ┌────────┐   ┌────────┐   ┌────────┐   ┌────────┐
  │ Memory │   │ Redis  │   │  NSQ   │   │ Kafka  │
  └────────┘   └────────┘   └────────┘   └────────┘
```

## 注意事项

1. **内存队列**：仅适用于单机场景，进程重启后消息会丢失
2. **重试机制**：所有队列都支持消息重试，最大重试次数可配置
3. **消息顺序**：除 Kafka 外，其他队列不保证严格的消息顺序
4. **持久化**：Memory 队列不支持持久化，其他队列需要根据配置决定
5. **性能考虑**：生产环境推荐使用 Kafka、NSQ 或 RabbitMQ

## 使用场景

- **任务调度**: 异步任务分发和处理
- **订阅更新**: 订阅内容变更通知
- **节点测试**: 节点延迟测试任务队列
- **事件驱动**: 系统事件通知
- **日志收集**: 异步日志处理

## 依赖

实现不同的队列后端需要添加相应的依赖：

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
