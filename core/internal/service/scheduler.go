package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/prism/core/internal/storage"
)

// SchedulerService 调度服务
type SchedulerService struct {
	db              *storage.Database
	subscriptionSvc *SubscriptionService
	nodeSvc         *NodeService
	ctx             context.Context
	cancel          context.CancelFunc
	updateTicker    *time.Ticker
	testTicker      *time.Ticker
	cleanupTicker   *time.Ticker
	wg              sync.WaitGroup
	running         bool
	mu              sync.RWMutex
}

// NewSchedulerService 创建调度服务
func NewSchedulerService(db *storage.Database, subscriptionSvc *SubscriptionService, nodeSvc *NodeService) *SchedulerService {
	ctx, cancel := context.WithCancel(context.Background())
	return &SchedulerService{
		db:              db,
		subscriptionSvc: subscriptionSvc,
		nodeSvc:         nodeSvc,
		ctx:             ctx,
		cancel:          cancel,
		running:         false,
	}
}

// Start 启动调度服务
func (s *SchedulerService) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("scheduler already running")
	}

	log.Println("Starting scheduler service...")

	// 创建定时器
	s.updateTicker = time.NewTicker(5 * time.Minute) // 每5分钟检查订阅更新
	s.testTicker = time.NewTicker(30 * time.Minute)  // 每30分钟进行节点测试
	s.cleanupTicker = time.NewTicker(6 * time.Hour)  // 每6小时清理过期数据

	// 启动协程处理定时任务
	s.wg.Add(3)
	go s.handleSubscriptionUpdates()
	go s.handleNodeTests()
	go s.handleDataCleanup()

	s.running = true
	log.Println("Scheduler service started")

	return nil
}

// Stop 停止调度服务
func (s *SchedulerService) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return fmt.Errorf("scheduler not running")
	}

	log.Println("Stopping scheduler service...")

	// 停止定时器
	if s.updateTicker != nil {
		s.updateTicker.Stop()
	}
	if s.testTicker != nil {
		s.testTicker.Stop()
	}
	if s.cleanupTicker != nil {
		s.cleanupTicker.Stop()
	}

	// 取消上下文
	s.cancel()

	// 等待所有协程结束
	s.wg.Wait()

	s.running = false
	log.Println("Scheduler service stopped")

	return nil
}

// IsRunning 检查是否运行中
func (s *SchedulerService) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// TriggerSubscriptionUpdate 手动触发订阅更新
func (s *SchedulerService) TriggerSubscriptionUpdate() {
	go s.updateSubscriptions()
}

// TriggerNodeTest 手动触发节点测试
func (s *SchedulerService) TriggerNodeTest() {
	go s.testNodes()
}

// handleSubscriptionUpdates 处理订阅自动更新
func (s *SchedulerService) handleSubscriptionUpdates() {
	defer s.wg.Done()

	// 启动时立即执行一次
	s.updateSubscriptions()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.updateTicker.C:
			s.updateSubscriptions()
		}
	}
}

// handleNodeTests 处理节点定时测试
func (s *SchedulerService) handleNodeTests() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.testTicker.C:
			s.testNodes()
		}
	}
}

// handleDataCleanup 处理数据清理
func (s *SchedulerService) handleDataCleanup() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.cleanupTicker.C:
			s.cleanupOldData()
		}
	}
}

// updateSubscriptions 更新订阅
func (s *SchedulerService) updateSubscriptions() {
	log.Println("Starting scheduled subscription updates...")

	var subscriptions []storage.Subscription
	if err := s.db.Where("auto_update = true AND status = 'active'").Find(&subscriptions).Error; err != nil {
		log.Printf("Failed to get subscriptions for update: %v", err)
		return
	}

	updateCount := 0
	errorCount := 0

	for _, subscription := range subscriptions {
		// 检查是否需要更新
		if !s.shouldUpdateSubscription(&subscription) {
			continue
		}

		log.Printf("Updating subscription: %s (ID: %d)", subscription.Name, subscription.ID)

		result, err := s.subscriptionSvc.UpdateSubscriptionContent(subscription.ID)
		if err != nil {
			log.Printf("Failed to update subscription %d: %v", subscription.ID, err)
			errorCount++

			// 更新错误计数
			subscription.ErrorCount++
			if subscription.ErrorCount >= 5 {
				subscription.Status = "error"
				log.Printf("Subscription %d marked as error after 5 failures", subscription.ID)
			}
			subscription.ErrorMessage = err.Error()
			s.db.Save(&subscription)
			continue
		}

		updateCount++
		log.Printf("Successfully updated subscription %d: %+v", subscription.ID, result)

		// 短暂延迟避免请求过于频繁
		time.Sleep(2 * time.Second)
	}

	log.Printf("Subscription update completed: %d updated, %d errors", updateCount, errorCount)
}

// shouldUpdateSubscription 检查是否应该更新订阅
func (s *SchedulerService) shouldUpdateSubscription(subscription *storage.Subscription) bool {
	if subscription.LastUpdate == nil {
		return true // 从未更新过
	}

	// 检查更新间隔
	updateInterval := time.Duration(subscription.UpdateInterval) * time.Second
	return time.Since(*subscription.LastUpdate) >= updateInterval
}

// testNodes 测试节点
func (s *SchedulerService) testNodes() {
	log.Println("Starting scheduled node tests...")

	// 获取需要测试的节点
	var nodes []storage.Node
	testInterval := 2 * time.Hour // 每2小时测试一次
	cutoffTime := time.Now().Add(-testInterval)

	query := s.db.Where("status != 'offline' AND (last_test IS NULL OR last_test < ?)", cutoffTime)

	// 优先测试存活节点和未知状态节点
	if err := query.Order("CASE WHEN status = 'online' THEN 1 WHEN status = 'unknown' THEN 2 ELSE 3 END, last_test ASC").
		Limit(50). // 限制一次测试的节点数量
		Find(&nodes).Error; err != nil {
		log.Printf("Failed to get nodes for testing: %v", err)
		return
	}

	if len(nodes) == 0 {
		log.Println("No nodes need testing")
		return
	}

	log.Printf("Testing %d nodes...", len(nodes))

	// 批量测试节点
	nodeIDs := make([]uint, len(nodes))
	for i, node := range nodes {
		nodeIDs[i] = node.ID
	}

	testTypes := []string{"delay", "speed"}
	testConfig := map[string]interface{}{
		"delay_url":  "http://www.gstatic.com/generate_204",
		"timeout":    5000,
		"concurrent": 10,
	}

	batchReq := &BatchTestRequest{
		NodeIDs:    nodeIDs,
		TestTypes:  testTypes,
		TestConfig: testConfig,
	}

	result, err := s.nodeSvc.BatchTestNodes(batchReq)
	if err != nil {
		log.Printf("Failed to start batch test: %v", err)
		return
	}

	log.Printf("Started batch test task: %s", result.TaskID)

	// 异步等待测试完成并更新节点池统计
	go s.waitForTestCompletion(result.TaskID)
}

// waitForTestCompletion 等待测试完成
func (s *SchedulerService) waitForTestCompletion(taskID string) {
	timeout := time.After(30 * time.Minute) // 30分钟超时
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			log.Printf("Test task %s timed out", taskID)
			return
		case <-ticker.C:
			status, err := s.nodeSvc.GetTestTaskStatus(taskID)
			if err != nil {
				log.Printf("Failed to get test task status: %v", err)
				return
			}

			if status.Status == "completed" {
				log.Printf("Test task %s completed: %d/%d successful",
					taskID,
					s.countSuccessfulResults(status.Results),
					len(status.Results))

				// 更新所有相关节点池的统计信息
				s.updateAllNodePoolStats()
				return
			} else if status.Status == "failed" {
				log.Printf("Test task %s failed: %s", taskID, status.Error)
				return
			}
		}
	}
}

// countSuccessfulResults 统计成功的测试结果
func (s *SchedulerService) countSuccessfulResults(results []storage.TestResult) int {
	count := 0
	for _, result := range results {
		if result.Success {
			count++
		}
	}
	return count
}

// updateAllNodePoolStats 更新所有节点池统计
func (s *SchedulerService) updateAllNodePoolStats() {
	var nodePools []storage.NodePool
	if err := s.db.Find(&nodePools).Error; err != nil {
		log.Printf("Failed to get node pools for stats update: %v", err)
		return
	}

	for _, nodePool := range nodePools {
		if err := nodePool.UpdateStats(s.db.DB); err != nil {
			log.Printf("Failed to update stats for node pool %d: %v", nodePool.ID, err)
		}
	}

	log.Printf("Updated stats for %d node pools", len(nodePools))
}

// cleanupOldData 清理过期数据
func (s *SchedulerService) cleanupOldData() {
	log.Println("Starting data cleanup...")

	// 清理7天前的测试记录
	testCleanupDuration := 7 * 24 * time.Hour
	if err := s.db.Cleanup(testCleanupDuration); err != nil {
		log.Printf("Failed to cleanup test data: %v", err)
	} else {
		log.Println("Test data cleanup completed")
	}

	// 清理30天前的订阅日志
	logCleanupTime := time.Now().Add(-30 * 24 * time.Hour)
	result := s.db.Where("created_at < ?", logCleanupTime).Delete(&storage.SubscriptionLog{})
	if result.Error != nil {
		log.Printf("Failed to cleanup subscription logs: %v", result.Error)
	} else {
		log.Printf("Cleaned up %d old subscription logs", result.RowsAffected)
	}

	// 数据库优化
	if err := s.db.OptimizeDatabase(); err != nil {
		log.Printf("Failed to optimize database: %v", err)
	} else {
		log.Println("Database optimization completed")
	}

	log.Println("Data cleanup completed")
}

// GetSchedulerStatus 获取调度器状态
func (s *SchedulerService) GetSchedulerStatus() *SchedulerStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	status := &SchedulerStatus{
		Running: s.running,
	}

	if s.running {
		// 获取下次执行时间
		now := time.Now()
		if s.updateTicker != nil {
			// 估算下次订阅更新时间（简化计算）
			status.NextSubscriptionUpdate = now.Add(5 * time.Minute).Truncate(5 * time.Minute)
		}
		if s.testTicker != nil {
			// 估算下次节点测试时间
			status.NextNodeTest = now.Add(30 * time.Minute).Truncate(30 * time.Minute)
		}
		if s.cleanupTicker != nil {
			// 估算下次数据清理时间
			status.NextDataCleanup = now.Add(6 * time.Hour).Truncate(6 * time.Hour)
		}
	}

	return status
}

// SchedulerStatus 调度器状态
type SchedulerStatus struct {
	Running                bool      `json:"running"`
	NextSubscriptionUpdate time.Time `json:"next_subscription_update,omitempty"`
	NextNodeTest           time.Time `json:"next_node_test,omitempty"`
	NextDataCleanup        time.Time `json:"next_data_cleanup,omitempty"`
}
