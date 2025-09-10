package service

import (
	"testing"

	"github.com/prism/core/internal/storage"
	"github.com/prism/core/internal/testutil"
)

func TestSubscriptionService_Create(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	service := NewSubscriptionService(db)

	req := &CreateSubscriptionRequest{
		Name:           "Test Subscription",
		URL:            "https://example.com/subscription",
		UserAgent:      "Test-Agent/1.0",
		AutoUpdate:     true,
		UpdateInterval: 3600,
		NodePoolIDs:    []uint{},
	}

	subscription, err := service.CreateSubscription(req)
	if err != nil {
		t.Fatalf("Failed to create subscription: %v", err)
	}

	if subscription.ID == 0 {
		t.Error("Expected subscription ID to be set")
	}

	if subscription.Name != req.Name {
		t.Errorf("Expected name %q, got %q", req.Name, subscription.Name)
	}

	if subscription.Status != "active" {
		t.Errorf("Expected status 'active', got %q", subscription.Status)
	}
}

func TestSubscriptionService_List(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	testutil.SeedTestData(t, db)

	service := NewSubscriptionService(db)

	req := &ListSubscriptionsRequest{
		Page: 1,
		Size: 10,
	}

	response, err := service.ListSubscriptions(req)
	if err != nil {
		t.Fatalf("Failed to list subscriptions: %v", err)
	}

	if response.Total == 0 {
		t.Error("Expected at least one subscription")
	}

	if len(response.Subscriptions) == 0 {
		t.Error("Expected subscriptions in response")
	}

	// 测试分页
	req.Size = 1
	response, err = service.ListSubscriptions(req)
	if err != nil {
		t.Fatalf("Failed to list subscriptions with pagination: %v", err)
	}

	if len(response.Subscriptions) > 1 {
		t.Errorf("Expected at most 1 subscription, got %d", len(response.Subscriptions))
	}
}

func TestSubscriptionService_GetByID(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	testutil.SeedTestData(t, db)

	service := NewSubscriptionService(db)

	// 获取已存在的订阅
	var subscription storage.Subscription
	db.DB.First(&subscription)

	result, err := service.GetSubscription(subscription.ID)
	if err != nil {
		t.Fatalf("Failed to get subscription by ID: %v", err)
	}

	if result.ID != subscription.ID {
		t.Errorf("Expected ID %d, got %d", subscription.ID, result.ID)
	}

	// 测试不存在的订阅
	_, err = service.GetSubscription(9999)
	if err == nil {
		t.Error("Expected error when getting non-existent subscription")
	}
}

func TestSubscriptionService_Update(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	testutil.SeedTestData(t, db)

	service := NewSubscriptionService(db)

	// 获取已存在的订阅
	var subscription storage.Subscription
	db.DB.First(&subscription)

	newName := "Updated Subscription Name"
	newStatus := "inactive"
	req := &UpdateSubscriptionRequest{
		Name:   &newName,
		Status: &newStatus,
	}

	result, err := service.UpdateSubscription(subscription.ID, req)
	if err != nil {
		t.Fatalf("Failed to update subscription: %v", err)
	}

	if result.Name != newName {
		t.Errorf("Expected name %q, got %q", newName, result.Name)
	}

	if result.Status != newStatus {
		t.Errorf("Expected status %q, got %q", newStatus, result.Status)
	}
}

func TestSubscriptionService_Delete(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	testutil.SeedTestData(t, db)

	service := NewSubscriptionService(db)

	// 获取已存在的订阅
	var subscription storage.Subscription
	db.DB.First(&subscription)

	err := service.DeleteSubscription(subscription.ID)
	if err != nil {
		t.Fatalf("Failed to delete subscription: %v", err)
	}

	// 验证删除
	_, err = service.GetSubscription(subscription.ID)
	if err == nil {
		t.Error("Expected error when getting deleted subscription")
	}
}

// TODO: Fix hash constraint issue with subscription content parsing
/*
func TestSubscriptionService_UpdateSubscription(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	// Note: Not seeding test data to avoid hash conflicts

	service := NewSubscriptionService(db)

	// 创建测试服务器
	mockServer := testutil.MockHTTPServer(map[string]string{
		"/subscription": testutil.MockSubscriptionData,
	})
	defer mockServer.Close()

	// 创建订阅
	req := &CreateSubscriptionRequest{
		Name:           "Test Subscription",
		URL:            mockServer.URL + "/subscription",
		UserAgent:      "Test-Agent/1.0",
		AutoUpdate:     true,
		UpdateInterval: 3600,
	}

	subscription, err := service.CreateSubscription(req)
	if err != nil {
		t.Fatalf("Failed to create subscription: %v", err)
	}

	// 更新订阅
	result, err := service.UpdateSubscriptionContent(subscription.ID)
	if err != nil {
		t.Fatalf("Failed to update subscription: %v", err)
	}

	if result.TotalFetched == 0 {
		t.Error("Expected to fetch some nodes")
	}

	// 验证节点已创建
	var nodeCount int64
	db.DB.Model(&storage.Node{}).Where("subscription_id = ?", subscription.ID).Count(&nodeCount)
	if nodeCount == 0 {
		t.Error("Expected nodes to be created")
	}
}
*/

// TODO: Implement GetSubscriptionStats method in SubscriptionService
/*
func TestSubscriptionService_GetSubscriptionStats(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	testutil.SeedTestData(t, db)

	service := NewSubscriptionService(db)

	// 获取已存在的订阅
	var subscription storage.Subscription
	db.DB.First(&subscription)

	stats, err := service.GetSubscriptionStats(subscription.ID)
	if err != nil {
		t.Fatalf("Failed to get subscription stats: %v", err)
	}

	if stats.SubscriptionID != subscription.ID {
		t.Errorf("Expected subscription ID %d, got %d", subscription.ID, stats.SubscriptionID)
	}

	if stats.TotalNodes != subscription.TotalNodes {
		t.Errorf("Expected total nodes %d, got %d", subscription.TotalNodes, stats.TotalNodes)
	}

	if stats.SurvivalRate == 0 {
		t.Error("Expected survival rate to be calculated")
	}

	if len(stats.ProtocolDistribution) == 0 {
		t.Error("Expected protocol distribution data")
	}
}
*/

// TODO: Implement GetLogs method in SubscriptionService
/*
func TestSubscriptionService_GetLogs(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	service := NewSubscriptionService(db)

	// 创建订阅
	subscription := &storage.Subscription{
		Name:   "Test Subscription",
		URL:    "https://example.com/subscription",
		Status: "active",
	}
	db.DB.Create(subscription)

	// 创建日志
	log := &storage.SubscriptionLog{
		SubscriptionID: subscription.ID,
		UpdateType:     "manual",
		Success:        true,
		Message:        "Test update",
		NodesAdded:     5,
		Duration:       1500,
		CreatedAt:      time.Now(),
	}
	db.DB.Create(log)

	// 获取日志
	req := &LogsRequest{
		Page: 1,
		Size: 10,
	}

	response, err := service.GetLogs(subscription.ID, req)
	if err != nil {
		t.Fatalf("Failed to get logs: %v", err)
	}

	if response.Total == 0 {
		t.Error("Expected at least one log entry")
	}

	if len(response.Logs) == 0 {
		t.Error("Expected logs in response")
	}

	if response.Logs[0].Message != log.Message {
		t.Errorf("Expected message %q, got %q", log.Message, response.Logs[0].Message)
	}
}
*/
