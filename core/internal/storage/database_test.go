package storage

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *Database {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	database := &Database{DB: db}
	if err := database.AutoMigrate(); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return database
}

func TestDatabase_AutoMigrate(t *testing.T) {
	db := setupTestDB(t)

	// 验证所有表都已创建
	tables := []string{"subscriptions", "nodes", "node_pools", "node_pool_subscriptions", "subscription_logs", "node_tests"}

	for _, table := range tables {
		if !db.DB.Migrator().HasTable(table) {
			t.Errorf("Table %s was not created", table)
		}
	}
}

func TestDatabase_CRUD_Subscription(t *testing.T) {
	db := setupTestDB(t)

	// Create
	subscription := &Subscription{
		Name:           "Test Subscription",
		URL:            "https://example.com/subscription",
		UserAgent:      "Test-Agent/1.0",
		Status:         "active",
		AutoUpdate:     true,
		UpdateInterval: 3600,
		TotalNodes:     5,
		ActiveNodes:    3,
	}

	err := db.DB.Create(subscription).Error
	if err != nil {
		t.Fatalf("Failed to create subscription: %v", err)
	}

	if subscription.ID == 0 {
		t.Error("Expected subscription ID to be set after creation")
	}

	// Read
	var retrieved Subscription
	err = db.DB.First(&retrieved, subscription.ID).Error
	if err != nil {
		t.Fatalf("Failed to retrieve subscription: %v", err)
	}

	if retrieved.Name != subscription.Name {
		t.Errorf("Expected name %q, got %q", subscription.Name, retrieved.Name)
	}

	// Update
	retrieved.Status = "inactive"
	err = db.DB.Save(&retrieved).Error
	if err != nil {
		t.Fatalf("Failed to update subscription: %v", err)
	}

	var updated Subscription
	db.DB.First(&updated, subscription.ID)
	if updated.Status != "inactive" {
		t.Errorf("Expected status 'inactive', got %q", updated.Status)
	}

	// Delete
	err = db.DB.Delete(&updated, updated.ID).Error
	if err != nil {
		t.Fatalf("Failed to delete subscription: %v", err)
	}

	var count int64
	db.DB.Model(&Subscription{}).Where("id = ?", updated.ID).Count(&count)
	if count != 0 {
		t.Error("Expected subscription to be deleted")
	}
}

func TestDatabase_CRUD_Node(t *testing.T) {
	db := setupTestDB(t)

	// 先创建订阅
	subscription := &Subscription{
		Name:   "Test Subscription",
		URL:    "https://example.com/subscription",
		Status: "active",
	}
	db.DB.Create(subscription)

	// Create Node
	delay := 100
	uploadSpeed := int64(1000)
	downloadSpeed := int64(5000)
	node := &Node{
		SubscriptionID: subscription.ID,
		Name:           "Test Node",
		Hash:           "test-hash",
		Server:         "test.example.com",
		Port:           443,
		Protocol:       "vmess",
		Country:        "US",
		CountryName:    "United States",
		City:           "Los Angeles",
		Status:         "active",
		Delay:          &delay,
		UploadSpeed:    &uploadSpeed,
		DownloadSpeed:  &downloadSpeed,
		ClashConfig: JSON{
			"name":   "Test Node",
			"type":   "vmess",
			"server": "test.example.com",
			"port":   443,
		},
	}

	err := db.DB.Create(node).Error
	if err != nil {
		t.Fatalf("Failed to create node: %v", err)
	}

	// Read
	var retrieved Node
	err = db.DB.First(&retrieved, node.ID).Error
	if err != nil {
		t.Fatalf("Failed to retrieve node: %v", err)
	}

	if retrieved.Server != node.Server {
		t.Errorf("Expected server %q, got %q", node.Server, retrieved.Server)
	}

	// Update
	retrieved.Status = "inactive"
	err = db.DB.Save(&retrieved).Error
	if err != nil {
		t.Fatalf("Failed to update node: %v", err)
	}

	// Delete
	err = db.DB.Delete(&retrieved, retrieved.ID).Error
	if err != nil {
		t.Fatalf("Failed to delete node: %v", err)
	}
}

func TestDatabase_NodePoolAssociations(t *testing.T) {
	db := setupTestDB(t)

	// 创建订阅
	subscription := &Subscription{
		Name:   "Test Subscription",
		URL:    "https://example.com/subscription",
		Status: "active",
	}
	db.DB.Create(subscription)

	// 创建节点池
	nodePool := &NodePool{
		Name:        "Test Pool",
		Description: "Test node pool",
		Enabled:     true,
		Priority:    1,
	}
	db.DB.Create(nodePool)

	// 创建关联
	association := &NodePoolSubscription{
		NodePoolID:     nodePool.ID,
		SubscriptionID: subscription.ID,
		Enabled:        true,
		Priority:       1,
	}
	err := db.DB.Create(association).Error
	if err != nil {
		t.Fatalf("Failed to create association: %v", err)
	}

	// 验证关联
	var count int64
	db.DB.Model(&NodePoolSubscription{}).
		Where("node_pool_id = ? AND subscription_id = ?", nodePool.ID, subscription.ID).
		Count(&count)

	if count != 1 {
		t.Errorf("Expected 1 association, got %d", count)
	}
}

func TestDatabase_SubscriptionLog(t *testing.T) {
	db := setupTestDB(t)

	// 创建订阅
	subscription := &Subscription{
		Name:   "Test Subscription",
		URL:    "https://example.com/subscription",
		Status: "active",
	}
	db.DB.Create(subscription)

	// 创建日志
	responseTime := 1500
	log := &SubscriptionLog{
		SubscriptionID: subscription.ID,
		UpdateType:     "manual",
		Success:        true,
		ErrorMessage:   "Update successful",
		NewNodes:       5,
		UpdatedNodes:   2,
		RemovedNodes:   1,
		ResponseTime:   &responseTime,
		CreatedAt:      time.Now(),
	}

	err := db.DB.Create(log).Error
	if err != nil {
		t.Fatalf("Failed to create subscription log: %v", err)
	}

	// 验证日志
	var retrieved SubscriptionLog
	err = db.DB.First(&retrieved, log.ID).Error
	if err != nil {
		t.Fatalf("Failed to retrieve subscription log: %v", err)
	}

	if retrieved.ErrorMessage != log.ErrorMessage {
		t.Errorf("Expected error message %q, got %q", log.ErrorMessage, retrieved.ErrorMessage)
	}
}
