package queue

import (
	"container/heap"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

// messageHeap 实现 heap.Interface，按到期时间排序的消息最小堆
type messageHeap[T any] []*Message[T]

func (h messageHeap[T]) Len() int { return len(h) }

func (h messageHeap[T]) Less(i, j int) bool {
	// 获取两个消息的到期时间（DelayUntil 和 NextRetryAt 的最大值）
	timeI := h[i].DelayUntil
	if h[i].NextRetryAt.After(timeI) {
		timeI = h[i].NextRetryAt
	}

	timeJ := h[j].DelayUntil
	if h[j].NextRetryAt.After(timeJ) {
		timeJ = h[j].NextRetryAt
	}

	// 零值时间视为最早（立即执行）
	if timeI.IsZero() {
		return true
	}
	if timeJ.IsZero() {
		return false
	}

	return timeI.Before(timeJ)
}

func (h messageHeap[T]) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *messageHeap[T]) Push(x interface{}) {
	*h = append(*h, x.(*Message[T]))
}

func (h *messageHeap[T]) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

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
	messagePool  sync.Pool        // 消息对象池

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

	// 初始化消息对象池
	q.messagePool = sync.Pool{
		New: func() interface{} {
			return &Message[T]{}
		},
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

// BatchPush 批量推送消息到队列（优化版 - 使用 sync.Pool）
func (q *MemoryQueue[T]) BatchPush(payloads []T) error {
	if q.closed.Load() {
		return ErrQueueClosed
	}

	if len(payloads) == 0 {
		return nil
	}

	// 预先生成所有消息（使用对象池）
	messages := make([]*Message[T], 0, len(payloads))
	now := time.Now()
	for _, payload := range payloads {
		msg := q.messagePool.Get().(*Message[T])
		// 重置消息字段
		msg.ID = uuid.New().String()
		msg.Payload = payload
		msg.Timestamp = now
		msg.DelayUntil = time.Time{}
		msg.RetryCount = 0
		msg.RetryPolicy = nil
		msg.LastError = nil
		msg.NextRetryAt = time.Time{}

		messages = append(messages, msg)
	}

	// 批量推送
	for i, msg := range messages {
		select {
		case q.queue <- msg:
		default:
			// 队列满，将未推送的消息归还对象池
			for j := i; j < len(messages); j++ {
				q.messagePool.Put(messages[j])
			}
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

// delayedProcessor 延时队列处理器（heap 优化版 - O(log n) 插入和删除）
func (q *MemoryQueue[T]) delayedProcessor() {
	defer q.delayDone.Done()

	// 初始化消息最小堆
	h := &messageHeap[T]{}
	heap.Init(h)

	// 动态 ticker，初始使用配置的间隔
	checkInterval := time.Duration(q.cfg.DelayCheckInterval) * time.Millisecond
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-q.delayedCh:
			// channel 已关闭，退出
			return

		case <-ticker.C:
	
			// 处理所有已到期的消息（堆顶）
			for h.Len() > 0 {
				msg := (*h)[0] // 查看堆顶元素

				if !msg.IsReady() {
					// 堆顶元素未到期，后续元素也都未到期
					break
				}

				// 从堆中移除已到期的消息
				heap.Pop(h)

				// 推送到主队列
				select {
				case q.queue <- msg:
				case <-q.delayedCh:
					return
				default:
					// 队列满，重新加入堆
					if q.cfg.MaxPendingMessages > 0 && h.Len() >= q.cfg.MaxPendingMessages {
						// 待处理列表已满，丢弃当前消息
						log.Printf("[Queue] WARNING: Pending heap full (%d), dropping ready message %s",
							q.cfg.MaxPendingMessages, msg.ID)
						q.metrics.dropped.Add(1)
					} else {
						log.Printf("[Queue] Main queue full, message %s delayed", msg.ID)
						heap.Push(h, msg)
					}
				}
			}

			// 动态调整下次检查的间隔
			if h.Len() > 0 {
				nextMsg := (*h)[0]
				nextTime := nextMsg.DelayUntil
				if nextMsg.NextRetryAt.After(nextTime) {
					nextTime = nextMsg.NextRetryAt
				}

				if !nextTime.IsZero() {
					untilNext := time.Until(nextTime)
					if untilNext > 0 && untilNext < checkInterval {
						// 下次消息到期时间更近，调整 ticker
						ticker.Reset(untilNext)
					} else {
						// 使用默认间隔
						ticker.Reset(checkInterval)
					}
				}
			} else {
				// 堆为空，使用默认间隔
				ticker.Reset(checkInterval)
			}

			// 从延时队列读取新消息并加入堆
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
							// 队列满，加入堆
							if q.cfg.MaxPendingMessages > 0 && h.Len() >= q.cfg.MaxPendingMessages {
								// 待处理堆已满，丢弃堆顶（最早到期）的消息
								oldMsg := heap.Pop(h).(*Message[T])
								log.Printf("[Queue] WARNING: Pending heap full (%d), dropping oldest message %s",
									q.cfg.MaxPendingMessages, oldMsg.ID)
								q.metrics.dropped.Add(1)
							}
							heap.Push(h, msg)
						}
					} else {
						// 未到期，加入堆
						if q.cfg.MaxPendingMessages > 0 && h.Len() >= q.cfg.MaxPendingMessages {
							// 待处理堆已满，丢弃堆顶（最早到期）的消息
							oldMsg := heap.Pop(h).(*Message[T])
							log.Printf("[Queue] WARNING: Pending heap full (%d), dropping oldest message %s",
								q.cfg.MaxPendingMessages, oldMsg.ID)
							q.metrics.dropped.Add(1)
						}
						heap.Push(h, msg)
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
func (q *MemoryQueue[T]) Metrics() QueueMetrics {
	return QueueMetrics{
		Processed:   q.metrics.processed.Load(),
		Failed:      q.metrics.failed.Load(),
		Retried:     q.metrics.retried.Load(),
		Dropped:     q.metrics.dropped.Load(),
		ActiveTasks: q.metrics.activeTasks.Load(),
		QueueDepth:  int64(q.Depth()),
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
