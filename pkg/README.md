# 公共包

## 说明

可被外部项目引用的公共包，高度独立，无内部依赖。

## 子包

### queue/ - 消息队列 ✅
高性能泛型消息队列抽象层，支持多种后端。

**已实现：**
- Memory: 内存队列（完整实现）

**待实现：**
- Redis: 基于 Redis Streams
- NSQ: 分布式实时消息平台
- Kafka: 高吞吐量分布式消息系统
- RabbitMQ: 可靠的消息代理
- ZeroMQ: 高性能异步消息库

**核心特性：**
- 泛型支持（类型安全）
- 智能重试（4种策略）
- 延时执行
- 监控指标
- 并发安全
- 优雅关闭

详见：[queue/README.md](queue/README.md)

### mihomo/ - Mihomo 代理核心
第三方代理核心模块（git submodule）。

### parser/ - 订阅解析器 ✅
通用代理订阅解析器，支持多种订阅格式。

**已实现：**
- Clash YAML 格式解析
- V2Ray Base64 链接解析
- 逐行 JSON 格式解析

**核心特性：**
- 自动格式识别
- 18+ 种协议支持（SS、SSR、VMess、VLess、Trojan、Hysteria 等）
- Mihomo 集成
- 智能回退机制

详见：[parser/README.md](parser/README.md)

## 设计原则

1. **独立性**: 不依赖 internal 包
2. **文档性**: 完善的 GoDoc 文档
3. **测试性**: 完整的单元测试
4. **稳定性**: API 保持向后兼容

## 使用示例

### Queue 队列

```go
import "github.com/darabuchi/prism/pkg/queue"

// 创建内存队列
q, err := queue.New[string](&queue.Config{
    Type:     queue.TypeMemory,
    MaxDepth: 10000,
})

// 推送消息
q.Push("task data")

// 处理消息
q.Process(5, func(msg *queue.Message[string]) (*queue.RetryInfo, error) {
    // 处理逻辑
    return nil, nil
})
```

### Parser 订阅解析器

```go
import "github.com/darabuchi/prism/pkg/parser"

// 解析订阅内容
data := []byte(`
proxies:
  - name: "香港节点"
    type: ss
    server: hk.example.com
    port: 443
    cipher: aes-256-gcm
    password: "password123"
`)

proxies, err := parser.Parse(data)
if err != nil {
    log.Errorf("解析失败: %v", err)
    return
}

log.Infof("解析到 %d 个节点", len(proxies))
```

## 图标说明

- ✅ 已完整实现
- 🚧 待实现/规划中
