package queue

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

// MemoryQueue 内存队列实现
type MemoryQueue[T any] struct {
	queue        chan *Message[T] // 消息队列（立即执行）
	delayedQueue chan *Message[T] // 延时队列
	cfg          *Config          // 配置
	closed       atomic.Bool      // 关闭标志
	workerCh     chan struct{}    // worker 停止信号（通过关闭通知）
	delayedCh    chan struct{}    // 延时处理器停止信号（通过关闭通知）
	processDone  sync.WaitGroup   // Process goroutine 计数
	delayDone    sync.WaitGroup   // 延时队列 goroutine 计数

	// 监控指标
	metrics struct {
		processed   atomic.Int64 // 已处理消息数
		failed      atomic.Int64 // 失败消息数
		retried     atomic.Int64 // 重试消息数
		dropped     atomic.Int64 // 丢弃消息数
		activeTasks atomic.Int64 // 活跃任务数
	}
}

// NewMemoryQueue 创建内存队列
func NewMemoryQueue[T any](cfg *Config) (*MemoryQueue[T], error) {
	q := &MemoryQueue[T]{
		queue:     make(chan *Message[T], cfg.MaxDepth),
		workerCh:  make(chan struct{}),
		delayedCh: make(chan struct{}),
		cfg:       cfg,
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

// BatchPush 批量推送消息到队列（优化版）
func (q *MemoryQueue[T]) BatchPush(payloads []T) error {
	if q.closed.Load() {
		return ErrQueueClosed
	}

	if len(payloads) == 0 {
		return nil
	}

	// 预先生成所有消息
	messages := make([]*Message[T], 0, len(payloads))
	now := time.Now()
	for _, payload := range payloads {
		msg := &Message[T]{
			ID:        uuid.New().String(),
			Payload:   payload,
			Timestamp: now,
		}
		messages = append(messages, msg)
	}

	// 批量推送
	for _, msg := range messages {
		select {
		case q.queue <- msg:
		default:
			return ErrQueueFull
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
func (q *MemoryQueue[T]) Process(concurrency int, handler func(*Message[T]) (*RetryInfo, error)) error {
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
func (q *MemoryQueue[T]) worker(handler func(*Message[T]) (*RetryInfo, error)) {
	defer q.processDone.Done()

	for {
		select {
		case <-q.workerCh:
			// channel 已关闭，退出
			return

		case msg, ok := <-q.queue:
			if !ok {
				// channel 已关闭
				return
			}

			if msg == nil {
				continue
			}

			// 检查消息是否准备好
			if !msg.IsReady() {
				// 未准备好，重新入队
				q.requeueMessage(msg)
				continue
			}

			// 增加活跃任务计数
			q.metrics.activeTasks.Add(1)

			// 处理消息
			retryInfo, err := handler(msg)

			// 减少活跃任务计数
			q.metrics.activeTasks.Add(-1)

			// 更新指标
			q.metrics.processed.Add(1)

			// 记录错误信息（如果有）
			if err != nil {
				msg.LastError = err
				q.metrics.failed.Add(1)
			}

			// 检查是否需要重试
			if retryInfo != nil && retryInfo.ShouldRetry {
				msg.RetryCount++
				q.metrics.retried.Add(1)

				// 检查是否超过最大重试次数
				if !msg.ShouldRetry(q.cfg.RetryPolicy) {
					// 超过最大重试次数，丢弃消息
					log.Printf("[Queue] Message %s exceeded max retries, dropped. LastError: %v", msg.ID, msg.LastError)
					q.metrics.dropped.Add(1)
					continue
				}

				// 计算下次重试时间
				if retryInfo.Delay > 0 {
					// 使用自定义延迟
					msg.NextRetryAt = time.Now().Add(retryInfo.Delay)
				} else {
					// 使用策略计算延迟
					msg.CalculateNextRetry(q.cfg.RetryPolicy)
				}

				// 重新入队（根据是否启用延时队列决定）
				if q.cfg.EnableDelayedQueue && !msg.NextRetryAt.IsZero() {
					if err := q.pushDelayedMessage(msg); err != nil {
						log.Printf("[Queue] Failed to push message %s to delayed queue: %v", msg.ID, err)
						q.metrics.dropped.Add(1)
					}
				} else {
					q.requeueMessage(msg)
				}
			}
			// 处理成功或不重试，消息完成
		}
	}
}

// delayedProcessor 延时队列处理器（优化版 - O(n)）
func (q *MemoryQueue[T]) delayedProcessor() {
	defer q.delayDone.Done()

	ticker := time.NewTicker(time.Duration(q.cfg.DelayCheckInterval) * time.Millisecond)
	defer ticker.Stop()

	// 临时存储未到期的消息
	pending := make([]*Message[T], 0)

	for {
		select {
		case <-q.delayedCh:
			// channel 已关闭，退出
			return

		case <-ticker.C:
			// 使用 map 优化查找（O(n) 而非 O(n²)）
			now := time.Now()
			readyMessages := make([]*Message[T], 0)
			newPending := make([]*Message[T], 0, len(pending))

			// 检查所有待处理的消息
			for _, msg := range pending {
				if msg.IsReady() {
					readyMessages = append(readyMessages, msg)
				} else {
					newPending = append(newPending, msg)
				}
			}

			pending = newPending

			// 将准备好的消息推送到主队列
			for _, msg := range readyMessages {
				select {
				case q.queue <- msg:
				case <-q.delayedCh:
					return
				default:
					// 队列满，保留在待处理列表
					log.Printf("[Queue] Main queue full, message %s delayed", msg.ID)
					pending = append(pending, msg)
				}
			}

			// 从延时队列读取新消息
		drainLoop:
			for {
				select {
				case msg, ok := <-q.delayedQueue:
					if !ok {
						// channel 已关闭
						return
					}

					if msg.IsReady() {
						// 立即推送到主队列
						select {
						case q.queue <- msg:
						case <-q.delayedCh:
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

// requeueMessage 重新入队消息（非阻塞，队列满时记录警告）
func (q *MemoryQueue[T]) requeueMessage(msg *Message[T]) {
	select {
	case q.queue <- msg:
	case <-q.workerCh:
		// worker 已停止
	default:
		// 队列满，记录警告并丢弃消息
		log.Printf("[Queue] WARNING: Queue full, message %s (retry %d) dropped", msg.ID, msg.RetryCount)
		q.metrics.dropped.Add(1)
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

// Metrics 获取监控指标
func (q *MemoryQueue[T]) Metrics() map[string]int64 {
	return map[string]int64{
		"processed":    q.metrics.processed.Load(),
		"failed":       q.metrics.failed.Load(),
		"retried":      q.metrics.retried.Load(),
		"dropped":      q.metrics.dropped.Load(),
		"active_tasks": q.metrics.activeTasks.Load(),
		"queue_depth":  int64(q.Depth()),
	}
}

// Close 关闭队列
func (q *MemoryQueue[T]) Close() error {
	if q.closed.Swap(true) {
		// 已经关闭
		return nil
	}

	// 关闭 worker 停止信号 channel（通知所有 worker 退出）
	close(q.workerCh)

	// 关闭延时处理器停止信号 channel
	if q.delayedQueue != nil {
		close(q.delayedCh)
	}

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
