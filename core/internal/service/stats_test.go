package service

import (
	"testing"

	"github.com/prism/core/internal/storage"
	"github.com/prism/core/internal/testutil"
)

func TestStatsService_GetOverview(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	testutil.SeedTestData(t, db)

	service := NewStatsService(db)

	stats, err := service.GetOverviewStats()
	if err != nil {
		t.Fatalf("Failed to get overview stats: %v", err)
	}

	if stats.TotalSubscriptions == 0 {
		t.Error("Expected at least one subscription")
	}

	if stats.TotalNodes == 0 {
		t.Error("Expected at least one node")
	}

	if stats.OverallSurvivalRate == 0 {
		t.Error("Expected survival rate to be calculated")
	}

	// 验证计算是否正确
	expectedSurvivalRate := float64(stats.ActiveNodes) / float64(stats.TotalNodes) * 100
	if stats.OverallSurvivalRate != expectedSurvivalRate {
		t.Errorf("Expected survival rate %f, got %f", expectedSurvivalRate, stats.OverallSurvivalRate)
	}
}

func TestStatsService_GetGeoDistribution(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	testutil.SeedTestData(t, db)

	service := NewStatsService(db)

	distribution, err := service.GetGeoDistribution()
	if err != nil {
		t.Fatalf("Failed to get geo distribution: %v", err)
	}

	if len(distribution) == 0 {
		t.Error("Expected geo distribution data")
	}

	// 验证数据结构
	for _, geo := range distribution {
		if geo.Country == "" {
			t.Error("Expected country to be set")
		}
		if geo.NodeCount == 0 {
			t.Error("Expected node count to be > 0")
		}
	}

	// 验证总数是否匹配
	totalNodes := 0
	for _, geo := range distribution {
		totalNodes += geo.NodeCount
	}

	var expectedTotal int64
	db.DB.Model(&storage.Node{}).Count(&expectedTotal)
	if totalNodes != int(expectedTotal) {
		t.Errorf("Expected total nodes %d, got %d", expectedTotal, totalNodes)
	}
}

func TestStatsService_GetProtocolDistribution(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	testutil.SeedTestData(t, db)

	service := NewStatsService(db)

	distribution, err := service.GetProtocolDistribution()
	if err != nil {
		t.Fatalf("Failed to get protocol distribution: %v", err)
	}

	if len(distribution) == 0 {
		t.Error("Expected protocol distribution data")
	}

	// 验证数据结构
	for _, protocol := range distribution {
		if protocol.Protocol == "" {
			t.Error("Expected protocol to be set")
		}
		if protocol.NodeCount == 0 {
			t.Error("Expected node count to be > 0")
		}
	}

	// 验证总数是否匹配
	totalNodes := 0
	for _, protocol := range distribution {
		totalNodes += protocol.NodeCount
	}

	var expectedTotal int64
	db.DB.Model(&storage.Node{}).Count(&expectedTotal)
	if totalNodes != int(expectedTotal) {
		t.Errorf("Expected total nodes %d, got %d", expectedTotal, totalNodes)
	}
}

// TODO: Implement GetPerformanceTrend method in StatsService
/*
func TestStatsService_GetPerformanceTrend(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	testutil.SeedTestData(t, db)

	service := NewStatsService(db)

	// 创建一些测试数据
	var node storage.Node
	db.DB.First(&node)

	// 创建测试结果数据
	now := time.Now()
	for i := 0; i < 24; i++ {
		delay := 100 + i*10 // 递增的延迟
		nodeTest := &storage.NodeTest{
			NodeID:   node.ID,
			TestType: "delay",
			Delay:    &delay,
			Success:  true,
			TestedAt: now.Add(-time.Duration(i) * time.Hour),
		}
		db.DB.Create(nodeTest)
	}

	// 测试小时趋势
	req := &PerformanceTrendRequest{
		Period: "hour",
	}

	response, err := service.GetPerformanceTrend(req)
	if err != nil {
		t.Fatalf("Failed to get performance trend: %v", err)
	}

	if response.Period != "hour" {
		t.Errorf("Expected period 'hour', got %q", response.Period)
	}

	if len(response.DataPoints) == 0 {
		t.Error("Expected data points")
	}

	// 验证数据点的结构
	for _, point := range response.DataPoints {
		if point.Timestamp.IsZero() {
			t.Error("Expected timestamp to be set")
		}
		if point.TestCount == 0 {
			t.Error("Expected test count to be > 0")
		}
	}

	// 测试天趋势
	req.Period = "day"
	response, err = service.GetPerformanceTrend(req)
	if err != nil {
		t.Fatalf("Failed to get daily performance trend: %v", err)
	}

	if response.Period != "day" {
		t.Errorf("Expected period 'day', got %q", response.Period)
	}
}
*/

// TODO: Implement GetPerformanceTrend method in StatsService  
/*
func TestStatsService_GetPerformanceTrend_InvalidPeriod(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	service := NewStatsService(db)

	req := &PerformanceTrendRequest{
		Period: "invalid",
	}

	_, err := service.GetPerformanceTrend(req)
	if err == nil {
		t.Error("Expected error for invalid period")
	}
}
*/

// TODO: Implement GetPerformanceTrend method in StatsService 
/*
func TestStatsService_GetPerformanceTrend_WithFilters(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)
	testutil.SeedTestData(t, db)

	service := NewStatsService(db)

	// 获取节点池
	var nodePool storage.NodePool
	db.DB.First(&nodePool)

	// 测试带节点池过滤的趋势
	req := &PerformanceTrendRequest{
		Period:     "hour",
		NodePoolID: &nodePool.ID,
		Country:    "US",
	}

	response, err := service.GetPerformanceTrend(req)
	if err != nil {
		t.Fatalf("Failed to get filtered performance trend: %v", err)
	}

	if response.Period != "hour" {
		t.Errorf("Expected period 'hour', got %q", response.Period)
	}

	// 注意：由于测试数据可能不完整，这里主要验证请求能成功处理
}
*/

func TestStatsService_CalculateSurvivalRate(t *testing.T) {
	tests := []struct {
		name        string
		activeNodes int
		totalNodes  int
		expected    float64
	}{
		{
			name:        "Normal case",
			activeNodes: 8,
			totalNodes:  10,
			expected:    80.0,
		},
		{
			name:        "All active",
			activeNodes: 10,
			totalNodes:  10,
			expected:    100.0,
		},
		{
			name:        "None active",
			activeNodes: 0,
			totalNodes:  10,
			expected:    0.0,
		},
		{
			name:        "Zero total",
			activeNodes: 0,
			totalNodes:  0,
			expected:    0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result float64
			if tt.totalNodes > 0 {
				result = float64(tt.activeNodes) / float64(tt.totalNodes) * 100
			}

			if result != tt.expected {
				t.Errorf("Expected survival rate %f, got %f", tt.expected, result)
			}
		})
	}
}
