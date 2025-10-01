package queue

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

// MemoryQueue 内存队列实现
type MemoryQueue[T any] struct {
	queue      chan *Message[T] // 消息队列
	cfg        *Config          // 配置
	closed     atomic.Bool      // 关闭标志
	stopCh     chan struct{}    // 停止信号
	processDone sync.WaitGroup  // Process goroutine 计数
}

// NewMemoryQueue 创建内存队列
func NewMemoryQueue[T any](cfg *Config) (*MemoryQueue[T], error) {
	q := &MemoryQueue[T]{
		queue:  make(chan *Message[T], cfg.MaxDepth),
		cfg:    cfg,
		stopCh: make(chan struct{}),
	}
	q.closed.Store(false)
	return q, nil
}

// Push 推送单个消息到队列
func (q *MemoryQueue[T]) Push(payload T) error {
	if q.closed.Load() {
		return ErrQueueClosed
	}

	msg := &Message[T]{
		ID:        uuid.New().String(),
		Payload:   payload,
		Timestamp: time.Now(),
		RetryCount: 0,
	}

	select {
	case q.queue <- msg:
		return nil
	default:
		return ErrQueueFull
	}
}

// BatchPush 批量推送消息到队列
func (q *MemoryQueue[T]) BatchPush(payloads []T) error {
	if q.closed.Load() {
		return ErrQueueClosed
	}

	for _, payload := range payloads {
		if err := q.Push(payload); err != nil {
			return err
		}
	}

	return nil
}

// Pop 从队列中弹出一个消息
func (q *MemoryQueue[T]) Pop() (*Message[T], error) {
	if q.closed.Load() {
		return nil, ErrQueueClosed
	}

	select {
	case msg := <-q.queue:
		return msg, nil
	default:
		return nil, ErrQueueEmpty
	}
}

// Process 启动消费者处理队列消息
func (q *MemoryQueue[T]) Process(concurrency int, handler func(*Message[T]) error) error {
	if concurrency <= 0 {
		return ErrInvalidConcurrency
	}

	if q.closed.Load() {
		return ErrQueueClosed
	}

	// 启动多个 worker goroutine
	for i := 0; i < concurrency; i++ {
		q.processDone.Add(1)
		go q.worker(handler)
	}

	return nil
}

// worker 消息处理工作协程
func (q *MemoryQueue[T]) worker(handler func(*Message[T]) error) {
	defer q.processDone.Done()

	for {
		select {
		case <-q.stopCh:
			return
		case msg := <-q.queue:
			if msg == nil {
				continue
			}

			// 处理消息
			err := handler(msg)
			if err != nil {
				// 处理失败，检查是否需要重试
				if msg.RetryCount < q.cfg.MaxRetry {
					msg.RetryCount++
					// 重新入队
					select {
					case q.queue <- msg:
					case <-q.stopCh:
						return
					default:
						// 队列满，丢弃消息
					}
				}
			}
		}
	}
}

// Depth 获取当前队列深度
func (q *MemoryQueue[T]) Depth() int {
	return len(q.queue)
}

// Close 关闭队列
func (q *MemoryQueue[T]) Close() error {
	if q.closed.Swap(true) {
		// 已经关闭
		return nil
	}

	// 发送停止信号
	close(q.stopCh)

	// 等待所有 worker 退出
	q.processDone.Wait()

	// 关闭队列 channel
	close(q.queue)

	return nil
}
