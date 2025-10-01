package queue

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

// MemoryQueue 内存队列实现
type MemoryQueue[T any] struct {
	queue       chan *Message[T] // 消息队列（立即执行）
	delayedQueue chan *Message[T] // 延时队列
	cfg         *Config          // 配置
	closed      atomic.Bool      // 关闭标志
	stopCh      chan struct{}    // 停止信号
	processDone sync.WaitGroup   // Process goroutine 计数
	delayDone   sync.WaitGroup   // 延时队列 goroutine 计数
}

// NewMemoryQueue 创建内存队列
func NewMemoryQueue[T any](cfg *Config) (*MemoryQueue[T], error) {
	q := &MemoryQueue[T]{
		queue:  make(chan *Message[T], cfg.MaxDepth),
		stopCh: make(chan struct{}),
		cfg:    cfg,
	}

	// 如果启用延时队列，创建延时队列并启动处理器
	if cfg.EnableDelayedQueue {
		q.delayedQueue = make(chan *Message[T], cfg.MaxDepth)
		q.delayDone.Add(1)
		go q.delayedProcessor()
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
	}

	return q.pushMessage(msg)
}

// PushWithDelay 推送延时消息到队列
func (q *MemoryQueue[T]) PushWithDelay(payload T, delay time.Duration) error {
	if q.closed.Load() {
		return ErrQueueClosed
	}

	if !q.cfg.EnableDelayedQueue {
		return ErrInvalidConfig
	}

	msg := &Message[T]{
		ID:         uuid.New().String(),
		Payload:    payload,
		Timestamp:  time.Now(),
		DelayUntil: time.Now().Add(delay),
	}

	return q.pushDelayedMessage(msg)
}

// PushWithRetryPolicy 推送带自定义重试策略的消息
func (q *MemoryQueue[T]) PushWithRetryPolicy(payload T, policy *RetryPolicy) error {
	if q.closed.Load() {
		return ErrQueueClosed
	}

	msg := &Message[T]{
		ID:          uuid.New().String(),
		Payload:     payload,
		Timestamp:   time.Now(),
		RetryPolicy: policy,
	}

	return q.pushMessage(msg)
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

			// 检查消息是否准备好
			if !msg.IsReady() {
				// 未准备好，重新入队
				q.requeueMessage(msg)
				continue
			}

			// 处理消息
			err := handler(msg)
			if err != nil {
				// 处理失败，记录错误
				msg.LastError = err
				msg.RetryCount++

				// 判断是否需要重试
				if msg.ShouldRetry(q.cfg.RetryPolicy) {
					// 计算下次重试时间
					msg.CalculateNextRetry(q.cfg.RetryPolicy)

					// 重新入队（根据是否启用延时队列决定）
					if q.cfg.EnableDelayedQueue && !msg.NextRetryAt.IsZero() {
						q.pushDelayedMessage(msg)
					} else {
						q.requeueMessage(msg)
					}
				}
				// 超过重试次数，丢弃消息
			}
		}
	}
}

// delayedProcessor 延时队列处理器
func (q *MemoryQueue[T]) delayedProcessor() {
	defer q.delayDone.Done()

	ticker := time.NewTicker(time.Duration(q.cfg.DelayCheckInterval) * time.Millisecond)
	defer ticker.Stop()

	// 临时存储未到期的消息
	pending := make([]*Message[T], 0)

	for {
		select {
		case <-q.stopCh:
			return

		case <-ticker.C:
			// 检查延时队列中的消息
			readyMessages := make([]*Message[T], 0)

			// 检查所有待处理的消息
			for _, msg := range pending {
				if msg.IsReady() {
					readyMessages = append(readyMessages, msg)
				}
			}

			// 移除已准备好的消息
			newPending := make([]*Message[T], 0, len(pending))
			for _, msg := range pending {
				ready := false
				for _, rm := range readyMessages {
					if rm.ID == msg.ID {
						ready = true
						break
					}
				}
				if !ready {
					newPending = append(newPending, msg)
				}
			}
			pending = newPending

			// 将准备好的消息推送到主队列
			for _, msg := range readyMessages {
				select {
				case q.queue <- msg:
				case <-q.stopCh:
					return
				default:
					// 队列满，保留在待处理列表
					pending = append(pending, msg)
				}
			}

			// 从延时队列读取新消息
		drainLoop:
			for {
				select {
				case msg := <-q.delayedQueue:
					if msg.IsReady() {
						// 立即推送到主队列
						select {
						case q.queue <- msg:
						case <-q.stopCh:
							return
						default:
							// 队列满，添加到待处理列表
							pending = append(pending, msg)
						}
					} else {
						// 未到期，添加到待处理列表
						pending = append(pending, msg)
					}
				default:
					break drainLoop
				}
			}
		}
	}
}

// pushMessage 推送消息到主队列
func (q *MemoryQueue[T]) pushMessage(msg *Message[T]) error {
	select {
	case q.queue <- msg:
		return nil
	default:
		return ErrQueueFull
	}
}

// pushDelayedMessage 推送消息到延时队列
func (q *MemoryQueue[T]) pushDelayedMessage(msg *Message[T]) error {
	if q.delayedQueue == nil {
		return ErrInvalidConfig
	}

	select {
	case q.delayedQueue <- msg:
		return nil
	default:
		return ErrQueueFull
	}
}

// requeueMessage 重新入队消息（非阻塞）
func (q *MemoryQueue[T]) requeueMessage(msg *Message[T]) {
	select {
	case q.queue <- msg:
	case <-q.stopCh:
	default:
		// 队列满，丢弃消息
	}
}

// Depth 获取当前队列深度
func (q *MemoryQueue[T]) Depth() int {
	depth := len(q.queue)
	if q.delayedQueue != nil {
		depth += len(q.delayedQueue)
	}
	return depth
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

	// 等待延时队列处理器退出
	if q.delayedQueue != nil {
		q.delayDone.Wait()
		close(q.delayedQueue)
	}

	// 关闭队列 channel
	close(q.queue)

	return nil
}
