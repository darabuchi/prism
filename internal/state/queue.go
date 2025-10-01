package state

import (
	"github.com/darabuchi/prism/pkg/queue"
	"github.com/lazygophers/log"
)

var (
	// TODO: 定义具体的队列实例
	// QueueSubscriptionUpdate *queue.MemoryQueue[*TaskSubscriptionUpdate]
	// QueueNodeTest *queue.MemoryQueue[*TaskNodeTest]
)

// InitQueue 初始化队列系统
func InitQueue() error {
	log.Info("try init queue")

	// TODO: 初始化具体的队列
	// QueueSubscriptionUpdate = queue.NewMemoryQueue[*TaskSubscriptionUpdate](&queue.Config{
	// 	MaxSize:            500,
	// 	ConcurrentWorkers:  5,
	// 	MaxRetries:         3,
	// 	RetryDelayBase:     time.Second,
	// 	MaxPendingMessages: 1000,
	// })

	log.Info("queue initialized successfully")

	return nil
}
