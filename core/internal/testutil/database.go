package testutil

import (
	"database/sql"
	"fmt"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/prism/core/internal/storage"
	
	// Force pure Go SQLite driver
	_ "modernc.org/sqlite"
)

// SetupTestDB 设置测试数据库
func SetupTestDB(t *testing.T) *storage.Database {
	// 直接使用纯 Go SQLite 驱动
	sqlDB, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open sqlite database: %v", err)
	}
	
	// 使用现有连接创建 GORM DB
	db, err := gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 静默模式，避免测试输出干扰
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// 创建 Database 实例
	database := &storage.Database{
		DB: db,
	}

	// 自动迁移所有模型
	err = database.AutoMigrate()
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return database
}

// CleanupTestDB 清理测试数据库
func CleanupTestDB(t *testing.T, db *storage.Database) {
	// 删除所有表的数据
	db.DB.Exec("DELETE FROM subscription_logs")
	db.DB.Exec("DELETE FROM node_tests")
	db.DB.Exec("DELETE FROM test_results")
	db.DB.Exec("DELETE FROM nodes")
	db.DB.Exec("DELETE FROM node_pool_subscriptions")
	db.DB.Exec("DELETE FROM node_pools")
	db.DB.Exec("DELETE FROM subscriptions")
}

// SeedTestData 创建测试数据
func SeedTestData(t *testing.T, db *storage.Database) {
	// 创建测试订阅
	subscription := &storage.Subscription{
		Name:           "Test Subscription",
		URL:            "https://example.com/subscription",
		UserAgent:      "Test-Agent/1.0",
		Status:         "active",
		AutoUpdate:     true,
		UpdateInterval: 3600,
		TotalNodes:     10,
		ActiveNodes:    8,
	}
	if err := db.DB.Create(subscription).Error; err != nil {
		t.Fatalf("Failed to create test subscription: %v", err)
	}

	// 创建测试节点池
	nodePool := &storage.NodePool{
		Name:        "Test Pool",
		Description: "Test node pool",
		Enabled:     true,
		Priority:    1,
	}
	if err := db.DB.Create(nodePool).Error; err != nil {
		t.Fatalf("Failed to create test node pool: %v", err)
	}

	// 创建测试节点
	for i := 0; i < 10; i++ {
		delay := i*50 + 100 // 100-550ms 延迟
		node := &storage.Node{
			SubscriptionID: subscription.ID,
			Name:           fmt.Sprintf("Test Node %d", i+1),
			Hash:           fmt.Sprintf("hash-%d", i+1),
			Server:         fmt.Sprintf("test%d.example.com", i+1),
			Port:           443,
			Protocol:       "vmess",
			Country:        "US",
			CountryName:    "United States",
			City:           "Los Angeles",
			ISP:            "Test ISP",
			Delay:          &delay,
			Status:         "online",
			ClashConfig: storage.JSON{
				"name":   fmt.Sprintf("Test Node %d", i+1),
				"type":   "vmess",
				"server": fmt.Sprintf("test%d.example.com", i+1),
				"port":   443,
			},
		}
		if i >= 8 {
			node.Status = "offline" // 最后两个节点设为不活跃
		}
		if err := db.DB.Create(node).Error; err != nil {
			t.Fatalf("Failed to create test node %d: %v", i+1, err)
		}
	}

	// 关联节点池和订阅
	// TODO: Fix NodePoolSubscription migration issue
	// association := &storage.NodePoolSubscription{
	// 	NodePoolID:     nodePool.ID,
	// 	SubscriptionID: subscription.ID,
	// 	Enabled:        true,
	// 	Priority:       1,
	// }
	// if err := db.DB.Create(association).Error; err != nil {
	// 	t.Fatalf("Failed to create node pool association: %v", err)
	// }
}
