package service

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/prism/core/internal/storage"
)

// NodeService 节点服务
type NodeService struct {
	db        *storage.Database
	testTasks map[string]*storage.TestTask
	taskMutex sync.RWMutex
}

// NewNodeService 创建节点服务
func NewNodeService(db *storage.Database) *NodeService {
	return &NodeService{
		db:        db,
		testTasks: make(map[string]*storage.TestTask),
	}
}

// GetNode 获取单个节点
func (s *NodeService) GetNode(id uint) (*storage.Node, error) {
	var node storage.Node
	if err := s.db.Preload("Subscription").Preload("NodePool").First(&node, id).Error; err != nil {
		return nil, fmt.Errorf("node not found: %w", err)
	}
	return &node, nil
}

// ListNodes 获取节点列表
func (s *NodeService) ListNodes(req *ListNodesRequest) (*ListNodesResponse, error) {
	var nodes []storage.Node
	var total int64

	query := s.db.Model(&storage.Node{})

	// 应用过滤条件
	if req.SubscriptionID != nil {
		query = query.Where("subscription_id = ?", *req.SubscriptionID)
	}
	if req.NodePoolID != nil {
		query = query.Where("node_pool_id = ?", *req.NodePoolID)
	}
	if req.Country != "" {
		query = query.Where("country = ?", req.Country)
	}
	if req.Protocol != "" {
		query = query.Where("protocol = ?", req.Protocol)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.MinDelay != nil {
		query = query.Where("delay >= ?", *req.MinDelay)
	}
	if req.MaxDelay != nil {
		query = query.Where("delay <= ?", *req.MaxDelay)
	}
	if req.StreamingUnlock != "" {
		query = query.Where("JSON_EXTRACT(streaming_unlock, ?) = true", fmt.Sprintf("$.%s.available", req.StreamingUnlock))
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count nodes: %w", err)
	}

	// 应用排序
	orderClause := "created_at DESC"
	if req.Sort != "" && req.Order != "" {
		validSorts := []string{"delay", "upload_speed", "download_speed", "last_test"}
		validOrders := []string{"asc", "desc"}

		sortValid := false
		for _, validSort := range validSorts {
			if req.Sort == validSort {
				sortValid = true
				break
			}
		}

		orderValid := false
		for _, validOrder := range validOrders {
			if req.Order == validOrder {
				orderValid = true
				break
			}
		}

		if sortValid && orderValid {
			orderClause = fmt.Sprintf("%s %s", req.Sort, strings.ToUpper(req.Order))
		}
	}

	// 分页查询
	offset := (req.Page - 1) * req.Size
	if err := query.Offset(offset).Limit(req.Size).Order(orderClause).Find(&nodes).Error; err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	return &ListNodesResponse{
		Total: int(total),
		Page:  req.Page,
		Size:  req.Size,
		Nodes: nodes,
	}, nil
}

// TestNode 测试单个节点
func (s *NodeService) TestNode(id uint, req *TestNodeRequest) (*storage.NodeTest, error) {
	// 获取节点
	var node storage.Node
	if err := s.db.First(&node, id).Error; err != nil {
		return nil, fmt.Errorf("node not found: %w", err)
	}

	// 执行测试
	testResult := s.performNodeTest(&node, req.TestTypes, req.TestConfig)

	// 保存测试记录
	nodeTest := &storage.NodeTest{
		NodeID:           id,
		TestType:         strings.Join(req.TestTypes, ","),
		TestConfig:       storage.JSON(req.TestConfig),
		Delay:            testResult.Delay,
		UploadSpeed:      testResult.UploadSpeed,
		DownloadSpeed:    testResult.DownloadSpeed,
		LossRate:         testResult.LossRate,
		StreamingResults: testResult.StreamingResults,
		Success:          testResult.Success,
		ErrorMessage:     testResult.ErrorMessage,
	}

	if err := s.db.Create(nodeTest).Error; err != nil {
		return nil, fmt.Errorf("failed to save test result: %w", err)
	}

	// 更新节点状态
	if testResult.Success {
		node.UpdateStatus("online", testResult.Delay)
		if testResult.UploadSpeed != nil {
			node.UploadSpeed = testResult.UploadSpeed
		}
		if testResult.DownloadSpeed != nil {
			node.DownloadSpeed = testResult.DownloadSpeed
		}
		if testResult.LossRate != nil {
			node.LossRate = testResult.LossRate
		}
		if testResult.StreamingResults != nil {
			node.StreamingUnlock = testResult.StreamingResults
		}
	} else {
		node.UpdateStatus("offline", nil)
	}

	if err := s.db.Save(&node).Error; err != nil {
		return nil, fmt.Errorf("failed to update node status: %w", err)
	}

	return nodeTest, nil
}

// BatchTestNodes 批量测试节点
func (s *NodeService) BatchTestNodes(req *BatchTestRequest) (*BatchTestResponse, error) {
	taskID := uuid.New().String()

	// 创建测试任务
	task := &storage.TestTask{
		ID:        taskID,
		NodeIDs:   req.NodeIDs,
		TestTypes: req.TestTypes,
		Status:    "running",
		Progress:  0,
		Total:     len(req.NodeIDs),
		Results:   make([]storage.TestResult, 0),
		StartedAt: time.Now(),
	}

	s.taskMutex.Lock()
	s.testTasks[taskID] = task
	s.taskMutex.Unlock()

	// 异步执行批量测试
	go s.executeBatchTest(task, req.TestConfig)

	return &BatchTestResponse{
		TaskID:     taskID,
		TotalNodes: len(req.NodeIDs),
		Status:     "running",
	}, nil
}

// GetTestTaskStatus 获取测试任务状态
func (s *NodeService) GetTestTaskStatus(taskID string) (*TestTaskStatus, error) {
	s.taskMutex.RLock()
	task, exists := s.testTasks[taskID]
	s.taskMutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("test task not found")
	}

	return &TestTaskStatus{
		TaskID:      task.ID,
		Status:      task.Status,
		Progress:    task.Progress,
		Total:       task.Total,
		Results:     task.Results,
		StartedAt:   task.StartedAt,
		CompletedAt: task.CompletedAt,
		Error:       task.Error,
	}, nil
}

// GetNodeTestHistory 获取节点测试历史
func (s *NodeService) GetNodeTestHistory(nodeID uint, req *NodeTestHistoryRequest) (*NodeTestHistoryResponse, error) {
	var tests []storage.NodeTest

	query := s.db.Model(&storage.NodeTest{}).Where("node_id = ?", nodeID)

	// 应用过滤条件
	if req.TestType != "" {
		query = query.Where("test_type LIKE ?", "%"+req.TestType+"%")
	}
	if !req.StartTime.IsZero() {
		query = query.Where("tested_at >= ?", req.StartTime)
	}
	if !req.EndTime.IsZero() {
		query = query.Where("tested_at <= ?", req.EndTime)
	}

	// 查询测试记录
	if err := query.Order("tested_at DESC").Limit(req.Limit).Find(&tests).Error; err != nil {
		return nil, fmt.Errorf("failed to get test history: %w", err)
	}

	return &NodeTestHistoryResponse{
		Tests: tests,
	}, nil
}

// GetBestNodes 智能节点选择
func (s *NodeService) GetBestNodes(req *BestSelectionRequest) (*BestSelectionResponse, error) {
	query := s.db.Model(&storage.Node{}).Where("status = 'online'")

	// 应用过滤条件
	if req.NodePoolID != nil {
		query = query.Where("node_pool_id = ?", *req.NodePoolID)
	}
	if req.Country != "" {
		query = query.Where("country = ?", req.Country)
	}
	if req.Protocol != "" {
		query = query.Where("protocol = ?", req.Protocol)
	}

	// 流媒体解锁条件
	for _, service := range req.StreamingUnlock {
		query = query.Where("JSON_EXTRACT(streaming_unlock, ?) = true", fmt.Sprintf("$.%s.available", service))
	}

	var nodes []storage.Node
	if err := query.Find(&nodes).Error; err != nil {
		return nil, fmt.Errorf("failed to query nodes: %w", err)
	}

	if len(nodes) == 0 {
		return &BestSelectionResponse{Selections: []NodeSelection{}}, nil
	}

	// 计算节点评分并排序
	scoredNodes := s.calculateNodeScores(nodes, req.StreamingUnlock)

	// 选择最佳节点
	selections := make([]NodeSelection, 0, req.Count)
	for i := 0; i < req.Count && i < len(scoredNodes); i++ {
		selections = append(selections, scoredNodes[i])
	}

	return &BestSelectionResponse{
		Selections: selections,
	}, nil
}

// performNodeTest 执行节点测试
func (s *NodeService) performNodeTest(node *storage.Node, testTypes []string, testConfig map[string]interface{}) *TestNodeResult {
	result := &TestNodeResult{
		Success: true,
	}

	// 模拟测试逻辑
	for _, testType := range testTypes {
		switch testType {
		case "delay":
			delay := s.testDelay(node, testConfig)
			result.Delay = &delay
			if delay > 1000 { // 延迟超过1秒视为失败
				result.Success = false
				result.ErrorMessage = "High delay"
			}

		case "speed":
			upload, download := s.testSpeed(node, testConfig)
			result.UploadSpeed = &upload
			result.DownloadSpeed = &download
			if upload == 0 || download == 0 {
				result.Success = false
				result.ErrorMessage = "Speed test failed"
			}

		case "streaming":
			streamingResults := s.testStreaming(node, testConfig)
			result.StreamingResults = streamingResults
		}
	}

	return result
}

// testDelay 测试延迟
func (s *NodeService) testDelay(node *storage.Node, config map[string]interface{}) int {
	// 模拟延迟测试，实际应该通过网络连接测试
	baseDelay := 50 + rand.Intn(200) // 50-250ms 基础延迟

	// 根据地区调整延迟
	switch node.Country {
	case "HK":
		baseDelay += rand.Intn(50)
	case "US":
		baseDelay += 100 + rand.Intn(100)
	case "JP":
		baseDelay += 30 + rand.Intn(70)
	default:
		baseDelay += rand.Intn(150)
	}

	return baseDelay
}

// testSpeed 测试速度
func (s *NodeService) testSpeed(node *storage.Node, config map[string]interface{}) (int64, int64) {
	// 模拟速度测试
	upload := int64(1048576 * (10 + rand.Intn(90)))    // 10-100 Mbps
	download := int64(1048576 * (20 + rand.Intn(180))) // 20-200 Mbps

	return upload, download
}

// testStreaming 测试流媒体解锁
func (s *NodeService) testStreaming(node *storage.Node, config map[string]interface{}) storage.JSON {
	results := make(storage.JSON)

	services := []string{"netflix", "youtube", "disney_plus", "hulu", "amazon_prime", "chatgpt"}

	for _, service := range services {
		available := rand.Float32() > 0.3 // 70% 解锁概率

		serviceResult := map[string]interface{}{
			"available": available,
			"tested_at": time.Now(),
		}

		if available && service != "chatgpt" {
			// 为流媒体服务添加区域信息
			serviceResult["region"] = node.Country
		}

		results[service] = serviceResult
	}

	return results
}

// executeBatchTest 执行批量测试
func (s *NodeService) executeBatchTest(task *storage.TestTask, testConfig map[string]interface{}) {
	defer func() {
		now := time.Now()
		task.CompletedAt = &now
		if task.Status == "running" {
			task.Status = "completed"
		}
	}()

	for i, nodeID := range task.NodeIDs {
		var node storage.Node
		if err := s.db.First(&node, nodeID).Error; err != nil {
			task.Results = append(task.Results, storage.TestResult{
				NodeID:  nodeID,
				Success: false,
				Error:   fmt.Sprintf("Node not found: %v", err),
			})
			continue
		}

		// 执行测试
		testResult := s.performNodeTest(&node, task.TestTypes, testConfig)

		// 保存测试记录
		nodeTest := &storage.NodeTest{
			NodeID:           nodeID,
			TestType:         strings.Join(task.TestTypes, ","),
			TestConfig:       storage.JSON(testConfig),
			Delay:            testResult.Delay,
			UploadSpeed:      testResult.UploadSpeed,
			DownloadSpeed:    testResult.DownloadSpeed,
			LossRate:         testResult.LossRate,
			StreamingResults: testResult.StreamingResults,
			Success:          testResult.Success,
			ErrorMessage:     testResult.ErrorMessage,
		}

		s.db.Create(nodeTest)

		// 更新节点状态
		if testResult.Success {
			node.UpdateStatus("online", testResult.Delay)
			if testResult.UploadSpeed != nil {
				node.UploadSpeed = testResult.UploadSpeed
			}
			if testResult.DownloadSpeed != nil {
				node.DownloadSpeed = testResult.DownloadSpeed
			}
			if testResult.StreamingResults != nil {
				node.StreamingUnlock = testResult.StreamingResults
			}
		} else {
			node.UpdateStatus("offline", nil)
		}
		s.db.Save(&node)

		// 添加测试结果到任务
		result := storage.TestResult{
			NodeID:  nodeID,
			Success: testResult.Success,
			Error:   testResult.ErrorMessage,
		}
		if testResult.Success {
			result.Results = map[string]interface{}{
				"delay":          testResult.Delay,
				"upload_speed":   testResult.UploadSpeed,
				"download_speed": testResult.DownloadSpeed,
				"streaming":      testResult.StreamingResults,
			}
		}
		task.Results = append(task.Results, result)

		// 更新进度
		task.Progress = i + 1

		// 模拟测试延迟
		time.Sleep(time.Duration(100+rand.Intn(900)) * time.Millisecond)
	}
}

// calculateNodeScores 计算节点评分
func (s *NodeService) calculateNodeScores(nodes []storage.Node, streamingServices []string) []NodeSelection {
	selections := make([]NodeSelection, 0, len(nodes))

	for _, node := range nodes {
		scores := make(map[string]int)
		totalScore := 0

		// 延迟评分 (40%)
		delayScore := 100
		if node.Delay != nil {
			delay := *node.Delay
			if delay > 500 {
				delayScore = 20
			} else if delay > 200 {
				delayScore = 60
			} else if delay > 100 {
				delayScore = 80
			}
		}
		scores["delay_score"] = delayScore
		totalScore += delayScore * 40 / 100

		// 速度评分 (30%)
		speedScore := 80
		if node.UploadSpeed != nil && node.DownloadSpeed != nil {
			upload := *node.UploadSpeed / 1048576 // 转换为 Mbps
			download := *node.DownloadSpeed / 1048576
			if upload > 50 && download > 100 {
				speedScore = 100
			} else if upload > 20 && download > 50 {
				speedScore = 90
			}
		}
		scores["speed_score"] = speedScore
		totalScore += speedScore * 30 / 100

		// 稳定性评分 (20%)
		stabilityScore := 100 - node.ContinuousFailures*10
		if stabilityScore < 0 {
			stabilityScore = 0
		}
		scores["stability_score"] = stabilityScore
		totalScore += stabilityScore * 20 / 100

		// 流媒体解锁评分 (10%)
		streamingScore := 50
		if len(streamingServices) > 0 && node.StreamingUnlock != nil {
			unlocked := 0
			for _, service := range streamingServices {
				if serviceData, ok := node.StreamingUnlock[service].(map[string]interface{}); ok {
					if available, ok := serviceData["available"].(bool); ok && available {
						unlocked++
					}
				}
			}
			streamingScore = unlocked * 100 / len(streamingServices)
		}
		scores["streaming_score"] = streamingScore
		totalScore += streamingScore * 10 / 100

		selections = append(selections, NodeSelection{
			Node:            node,
			SelectionReason: scores,
		})
	}

	// 按总分排序
	for i := 0; i < len(selections)-1; i++ {
		for j := i + 1; j < len(selections); j++ {
			score1 := s.calculateTotalScore(selections[i].SelectionReason)
			score2 := s.calculateTotalScore(selections[j].SelectionReason)
			if score2 > score1 {
				selections[i], selections[j] = selections[j], selections[i]
			}
		}
	}

	return selections
}

// calculateTotalScore 计算总分
func (s *NodeService) calculateTotalScore(scores map[string]int) int {
	total := 0
	weights := map[string]int{
		"delay_score":     40,
		"speed_score":     30,
		"stability_score": 20,
		"streaming_score": 10,
	}

	for key, score := range scores {
		if weight, ok := weights[key]; ok {
			total += score * weight / 100
		}
	}

	return total
}

// TestNodeResult 节点测试结果
type TestNodeResult struct {
	Success          bool         `json:"success"`
	Delay            *int         `json:"delay"`
	UploadSpeed      *int64       `json:"upload_speed"`
	DownloadSpeed    *int64       `json:"download_speed"`
	LossRate         *float64     `json:"loss_rate"`
	StreamingResults storage.JSON `json:"streaming_results"`
	ErrorMessage     string       `json:"error_message"`
}
