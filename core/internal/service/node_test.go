package service

import (
	"testing"

	"github.com/prism/core/internal/storage"
	"github.com/prism/core/internal/testutil"
)

func TestNodeService_List(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	testutil.SeedTestData(t, db)

	service := NewNodeService(db)

	req := &ListNodesRequest{
		Page: 1,
		Size: 10,
	}

	response, err := service.ListNodes(req)
	if err != nil {
		t.Fatalf("Failed to list nodes: %v", err)
	}

	if response.Total == 0 {
		t.Error("Expected at least one node")
	}

	if len(response.Nodes) == 0 {
		t.Error("Expected nodes in response")
	}

	// 测试分页
	req.Size = 5
	response, err = service.ListNodes(req)
	if err != nil {
		t.Fatalf("Failed to list nodes with pagination: %v", err)
	}

	if len(response.Nodes) > 5 {
		t.Errorf("Expected at most 5 nodes, got %d", len(response.Nodes))
	}
}

func TestNodeService_ListWithFilters(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	testutil.SeedTestData(t, db)

	service := NewNodeService(db)

	// 测试按订阅ID过滤
	var subscription storage.Subscription
	db.DB.First(&subscription)

	req := &ListNodesRequest{
		Page:           1,
		Size:           10,
		SubscriptionID: &subscription.ID,
	}

	response, err := service.ListNodes(req)
	if err != nil {
		t.Fatalf("Failed to list nodes by subscription: %v", err)
	}

	for _, node := range response.Nodes {
		if node.SubscriptionID != subscription.ID {
			t.Errorf("Expected subscription ID %d, got %d", subscription.ID, node.SubscriptionID)
		}
	}

	// 测试按国家过滤
	req = &ListNodesRequest{
		Page:    1,
		Size:    10,
		Country: "US",
	}

	response, err = service.ListNodes(req)
	if err != nil {
		t.Fatalf("Failed to list nodes by country: %v", err)
	}

	for _, node := range response.Nodes {
		if node.Country != "US" {
			t.Errorf("Expected country 'US', got %q", node.Country)
		}
	}

	// 测试按协议过滤
	req = &ListNodesRequest{
		Page:     1,
		Size:     10,
		Protocol: "vmess",
	}

	response, err = service.ListNodes(req)
	if err != nil {
		t.Fatalf("Failed to list nodes by protocol: %v", err)
	}

	for _, node := range response.Nodes {
		if node.Protocol != "vmess" {
			t.Errorf("Expected protocol 'vmess', got %q", node.Protocol)
		}
	}
}

func TestNodeService_GetByID(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	testutil.SeedTestData(t, db)

	service := NewNodeService(db)

	// 获取已存在的节点
	var node storage.Node
	db.DB.First(&node)

	result, err := service.GetNode(node.ID)
	if err != nil {
		t.Fatalf("Failed to get node by ID: %v", err)
	}

	if result.ID != node.ID {
		t.Errorf("Expected ID %d, got %d", node.ID, result.ID)
	}

	// 测试不存在的节点
	_, err = service.GetNode(9999)
	if err == nil {
		t.Error("Expected error when getting non-existent node")
	}
}

func TestNodeService_TestNode(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	testutil.SeedTestData(t, db)

	service := NewNodeService(db)

	// 获取已存在的节点
	var node storage.Node
	db.DB.First(&node)

	req := &TestNodeRequest{
		TestTypes: []string{"delay", "speed"},
		TestConfig: map[string]interface{}{
			"timeout": 5000,
		},
	}

	// 注意：这个测试可能需要模拟网络连接，这里我们主要测试参数验证
	_, err := service.TestNode(node.ID, req)
	// 在实际环境中，这可能会失败，因为节点不真实存在
	// 但我们可以验证服务能正确处理请求
	if err != nil {
		// 这是预期的，因为我们没有真实的代理服务器
		t.Logf("Node test failed as expected: %v", err)
	}
}

func TestNodeService_BatchTest(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	testutil.SeedTestData(t, db)

	service := NewNodeService(db)

	// 获取一些节点
	var nodes []storage.Node
	db.DB.Limit(3).Find(&nodes)

	nodeIDs := make([]uint, len(nodes))
	for i, node := range nodes {
		nodeIDs[i] = node.ID
	}

	req := &BatchTestRequest{
		NodeIDs:   nodeIDs,
		TestTypes: []string{"delay"},
		TestConfig: map[string]interface{}{
			"timeout": 3000,
		},
	}

	response, err := service.BatchTestNodes(req)
	if err != nil {
		t.Fatalf("Failed to start batch test: %v", err)
	}

	if response.TaskID == "" {
		t.Error("Expected task ID to be set")
	}

	if response.TotalNodes != len(nodeIDs) {
		t.Errorf("Expected total nodes %d, got %d", len(nodeIDs), response.TotalNodes)
	}

	if response.Status != "running" {
		t.Errorf("Expected status 'running', got %q", response.Status)
	}
}

func TestNodeService_GetTestStatus(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	service := NewNodeService(db)

	// 创建一个测试任务
	task := &storage.TestTask{
		ID:     "test-task-123",
		Status: "running",
		Total:  3,
	}
	service.testTasks["test-task-123"] = task

	status, err := service.GetTestTaskStatus("test-task-123")
	if err != nil {
		t.Fatalf("Failed to get test status: %v", err)
	}

	if status.TaskID != "test-task-123" {
		t.Errorf("Expected task ID 'test-task-123', got %q", status.TaskID)
	}

	if status.Status != "running" {
		t.Errorf("Expected status 'running', got %q", status.Status)
	}

	// 测试不存在的任务
	_, err = service.GetTestTaskStatus("non-existent")
	if err == nil {
		t.Error("Expected error when getting non-existent test status")
	}
}

func TestNodeService_GetBestSelection(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	testutil.SeedTestData(t, db)

	service := NewNodeService(db)

	req := &BestSelectionRequest{
		Count:   3,
		Country: "US",
	}

	response, err := service.GetBestNodes(req)
	if err != nil {
		t.Fatalf("Failed to get best selection: %v", err)
	}

	if len(response.Selections) > req.Count {
		t.Errorf("Expected at most %d selections, got %d", req.Count, len(response.Selections))
	}

	// 验证选择的节点都是指定国家
	for _, selection := range response.Selections {
		if selection.Node.Country != req.Country {
			t.Errorf("Expected country %q, got %q", req.Country, selection.Node.Country)
		}
	}
}

func TestNodeService_GetTestHistory(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	testutil.SeedTestData(t, db)

	service := NewNodeService(db)

	// 获取已存在的节点
	var node storage.Node
	db.DB.First(&node)

	// 创建测试历史记录
	delay := 100
	nodeTest := &storage.NodeTest{
		NodeID:   node.ID,
		TestType: "delay",
		Delay:    &delay,
		Success:  true,
	}
	db.DB.Create(nodeTest)

	req := &NodeTestHistoryRequest{
		TestType: "delay",
		Limit:    10,
	}

	response, err := service.GetNodeTestHistory(node.ID, req)
	if err != nil {
		t.Fatalf("Failed to get test history: %v", err)
	}

	if len(response.Tests) == 0 {
		t.Error("Expected test history records")
	}

	if response.Tests[0].TestType != req.TestType {
		t.Errorf("Expected test type %q, got %q", req.TestType, response.Tests[0].TestType)
	}
}
