package service

import (
	"fmt"
	"time"

	"github.com/prism/core/internal/storage"
)

// StatsService 统计服务
type StatsService struct {
	db *storage.Database
}

// NewStatsService 创建统计服务
func NewStatsService(db *storage.Database) *StatsService {
	return &StatsService{db: db}
}

// GetOverviewStats 获取整体统计
func (s *StatsService) GetOverviewStats() (*OverviewStats, error) {
	stats := &OverviewStats{}

	// 订阅统计
	var totalSubscriptions, activeSubscriptions int64
	s.db.Model(&storage.Subscription{}).Count(&totalSubscriptions)
	s.db.Model(&storage.Subscription{}).Where("status = 'active'").Count(&activeSubscriptions)
	stats.TotalSubscriptions = int(totalSubscriptions)
	stats.ActiveSubscriptions = int(activeSubscriptions)

	// 节点池统计
	var totalNodePools int64
	s.db.Model(&storage.NodePool{}).Count(&totalNodePools)
	stats.TotalNodePools = int(totalNodePools)

	// 节点统计
	var totalNodes, activeNodes int64
	s.db.Model(&storage.Node{}).Count(&totalNodes)
	s.db.Model(&storage.Node{}).Where("status = 'online'").Count(&activeNodes)
	stats.TotalNodes = int(totalNodes)
	stats.ActiveNodes = int(activeNodes)

	// 计算整体存活率
	if totalNodes > 0 {
		stats.OverallSurvivalRate = float64(activeNodes) / float64(totalNodes) * 100
	}

	// 今日测试统计
	today := time.Now().Truncate(24 * time.Hour)
	var totalTestsToday, successfulTestsToday int64
	s.db.Model(&storage.NodeTest{}).Where("tested_at >= ?", today).Count(&totalTestsToday)
	s.db.Model(&storage.NodeTest{}).Where("tested_at >= ? AND success = true", today).Count(&successfulTestsToday)
	stats.TotalTestsToday = int(totalTestsToday)
	stats.SuccessfulTestsToday = int(successfulTestsToday)

	return stats, nil
}

// GetGeoDistribution 获取地区分布统计
func (s *StatsService) GetGeoDistribution() ([]GeoDistribution, error) {
	var distributions []GeoDistribution

	rows, err := s.db.Raw(`
		SELECT 
			country,
			country_name,
			COUNT(*) as node_count,
			SUM(CASE WHEN status = 'online' THEN 1 ELSE 0 END) as active_count
		FROM nodes 
		WHERE country IS NOT NULL AND country != ''
		GROUP BY country, country_name
		ORDER BY node_count DESC
	`).Rows()

	if err != nil {
		return nil, fmt.Errorf("failed to get geo distribution: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var dist GeoDistribution
		if err := rows.Scan(&dist.Country, &dist.CountryName, &dist.NodeCount, &dist.ActiveCount); err != nil {
			return nil, fmt.Errorf("failed to scan geo distribution: %w", err)
		}
		distributions = append(distributions, dist)
	}

	return distributions, nil
}

// GetProtocolDistribution 获取协议分布统计
func (s *StatsService) GetProtocolDistribution() ([]ProtocolDistribution, error) {
	var distributions []ProtocolDistribution

	rows, err := s.db.Raw(`
		SELECT 
			protocol,
			COUNT(*) as node_count,
			SUM(CASE WHEN status = 'online' THEN 1 ELSE 0 END) as active_count
		FROM nodes 
		GROUP BY protocol
		ORDER BY node_count DESC
	`).Rows()

	if err != nil {
		return nil, fmt.Errorf("failed to get protocol distribution: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var dist ProtocolDistribution
		if err := rows.Scan(&dist.Protocol, &dist.NodeCount, &dist.ActiveCount); err != nil {
			return nil, fmt.Errorf("failed to scan protocol distribution: %w", err)
		}
		distributions = append(distributions, dist)
	}

	return distributions, nil
}

// GetPerformanceTrend 获取性能趋势
func (s *StatsService) GetPerformanceTrend(req *PerformanceTrendRequest) (*PerformanceTrendResponse, error) {
	var dataPoints []TrendDataPoint

	// 根据周期确定时间间隔和查询语句
	var groupBy string
	var startTime time.Time

	now := time.Now()
	switch req.Period {
	case "hour":
		groupBy = "DATE_FORMAT(tested_at, '%Y-%m-%d %H')"
		startTime = now.Add(-24 * time.Hour)
	case "day":
		groupBy = "DATE(tested_at)"
		startTime = now.Add(-30 * 24 * time.Hour)
	case "week":
		groupBy = "YEARWEEK(tested_at)"
		startTime = now.Add(-12 * 7 * 24 * time.Hour)
	case "month":
		groupBy = "DATE_FORMAT(tested_at, '%Y-%m')"
		startTime = now.Add(-12 * 30 * 24 * time.Hour)
	default:
		return nil, fmt.Errorf("invalid period: %s", req.Period)
	}

	// 构建查询条件
	query := `
		SELECT 
			%s as time_group,
			AVG(delay) as avg_delay,
			AVG(upload_speed) as avg_upload,
			AVG(download_speed) as avg_download,
			COUNT(*) as test_count,
			SUM(CASE WHEN success = true THEN 1 ELSE 0 END) * 100.0 / COUNT(*) as survival_rate
		FROM node_tests nt
		JOIN nodes n ON nt.node_id = n.id
		WHERE nt.tested_at >= ?
	`
	args := []interface{}{startTime}

	// 添加过滤条件
	if req.NodePoolID != nil {
		query += " AND n.node_pool_id = ?"
		args = append(args, *req.NodePoolID)
	}
	if req.Country != "" {
		query += " AND n.country = ?"
		args = append(args, req.Country)
	}

	query += fmt.Sprintf(" GROUP BY %s ORDER BY time_group", groupBy)
	query = fmt.Sprintf(query, groupBy)

	rows, err := s.db.Raw(query, args...).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to get performance trend: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var timeGroup string
		var point TrendDataPoint
		var avgDelay, avgUpload, avgDownload *float64

		if err := rows.Scan(&timeGroup, &avgDelay, &avgUpload, &avgDownload, &point.TestCount, &point.SurvivalRate); err != nil {
			return nil, fmt.Errorf("failed to scan trend data: %w", err)
		}

		// 解析时间戳
		timestamp, err := s.parseTimeGroup(timeGroup, req.Period)
		if err != nil {
			continue
		}
		point.Timestamp = timestamp

		// 设置平均值
		if avgDelay != nil {
			point.AvgDelay = *avgDelay
		}
		if avgUpload != nil {
			point.AvgUpload = *avgUpload / 1048576 // 转换为 Mbps
		}
		if avgDownload != nil {
			point.AvgDownload = *avgDownload / 1048576 // 转换为 Mbps
		}

		dataPoints = append(dataPoints, point)
	}

	return &PerformanceTrendResponse{
		Period:     req.Period,
		DataPoints: dataPoints,
	}, nil
}

// parseTimeGroup 解析时间分组
func (s *StatsService) parseTimeGroup(timeGroup, period string) (time.Time, error) {
	var layout string
	switch period {
	case "hour":
		layout = "2006-01-02 15"
		timeGroup += ":00:00"
	case "day":
		layout = "2006-01-02"
	case "week":
		// YEARWEEK 格式需要特殊处理
		var year, week int
		if _, err := fmt.Sscanf(timeGroup, "%d%d", &year, &week); err != nil {
			return time.Time{}, err
		}
		// 计算该周的第一天
		jan1 := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		daysToAdd := (week-1)*7 - int(jan1.Weekday())
		return jan1.AddDate(0, 0, daysToAdd), nil
	case "month":
		layout = "2006-01"
		timeGroup += "-01"
	}

	if period != "week" {
		return time.Parse(layout, timeGroup)
	}

	return time.Time{}, fmt.Errorf("unsupported period for time parsing: %s", period)
}

// GetSubscriptionPerformance 获取订阅性能统计
func (s *StatsService) GetSubscriptionPerformance(subscriptionID uint) (*SubscriptionStats, error) {
	var subscription storage.Subscription
	if err := s.db.First(&subscription, subscriptionID).Error; err != nil {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}

	stats := &SubscriptionStats{
		SubscriptionID: subscriptionID,
		TotalNodes:     subscription.TotalNodes,
		ActiveNodes:    subscription.ActiveNodes,
	}

	// 计算存活率
	if subscription.TotalNodes > 0 {
		stats.SurvivalRate = float64(subscription.ActiveNodes) / float64(subscription.TotalNodes) * 100
	}

	// 获取协议分布
	var protocolStats []struct {
		Protocol string `json:"protocol"`
		Count    int    `json:"count"`
	}
	s.db.Model(&storage.Node{}).
		Where("subscription_id = ?", subscriptionID).
		Select("protocol, COUNT(*) as count").
		Group("protocol").
		Scan(&protocolStats)
	stats.ProtocolDistribution = make([]interface{}, len(protocolStats))
	for i, ps := range protocolStats {
		stats.ProtocolDistribution[i] = ps
	}

	// 获取地区分布
	var countryStats []struct {
		Country string `json:"country"`
		Count   int    `json:"count"`
	}
	s.db.Model(&storage.Node{}).
		Where("subscription_id = ? AND country IS NOT NULL", subscriptionID).
		Select("country, COUNT(*) as count").
		Group("country").
		Scan(&countryStats)
	stats.CountryDistribution = make([]interface{}, len(countryStats))
	for i, cs := range countryStats {
		stats.CountryDistribution[i] = cs
	}

	// 获取最近更新记录
	s.db.Where("subscription_id = ?", subscriptionID).
		Order("created_at DESC").
		Limit(10).
		Find(&stats.RecentLogs)

	return stats, nil
}

// GetNodePoolPerformance 获取节点池性能统计
func (s *StatsService) GetNodePoolPerformance(nodePoolID uint) (*NodePoolStats, error) {
	var nodePool storage.NodePool
	if err := s.db.First(&nodePool, nodePoolID).Error; err != nil {
		return nil, fmt.Errorf("node pool not found: %w", err)
	}

	stats := &NodePoolStats{
		NodePoolID:   nodePoolID,
		TotalNodes:   nodePool.TotalNodes,
		ActiveNodes:  nodePool.ActiveNodes,
		SurvivalRate: nodePool.SurvivalRate,
	}

	// 获取最近测试统计
	var recentTestStats struct {
		TotalTests   int64
		SuccessTests int64
		AvgDelay     *float64
		AvgUpload    *float64
		AvgDownload  *float64
	}

	yesterday := time.Now().Add(-24 * time.Hour)
	s.db.Raw(`
		SELECT 
			COUNT(*) as total_tests,
			SUM(CASE WHEN success = true THEN 1 ELSE 0 END) as success_tests,
			AVG(delay) as avg_delay,
			AVG(upload_speed) as avg_upload,
			AVG(download_speed) as avg_download
		FROM node_tests nt
		JOIN nodes n ON nt.node_id = n.id
		WHERE n.node_pool_id = ? AND nt.tested_at >= ?
	`, nodePoolID, yesterday).Scan(&recentTestStats)

	stats.RecentTestCount = int(recentTestStats.TotalTests)
	if recentTestStats.TotalTests > 0 {
		stats.RecentSuccessRate = float64(recentTestStats.SuccessTests) / float64(recentTestStats.TotalTests) * 100
	}
	if recentTestStats.AvgDelay != nil {
		avgDelay := int(*recentTestStats.AvgDelay)
		stats.AvgDelay = &avgDelay
	}
	if recentTestStats.AvgUpload != nil {
		avgUpload := int64(*recentTestStats.AvgUpload)
		stats.AvgUploadSpeed = &avgUpload
	}
	if recentTestStats.AvgDownload != nil {
		avgDownload := int64(*recentTestStats.AvgDownload)
		stats.AvgDownloadSpeed = &avgDownload
	}

	return stats, nil
}

// NodePoolStats 节点池统计信息
type NodePoolStats struct {
	NodePoolID        uint    `json:"node_pool_id"`
	TotalNodes        int     `json:"total_nodes"`
	ActiveNodes       int     `json:"active_nodes"`
	SurvivalRate      float64 `json:"survival_rate"`
	RecentTestCount   int     `json:"recent_test_count"`
	RecentSuccessRate float64 `json:"recent_success_rate"`
	AvgDelay          *int    `json:"avg_delay"`
	AvgUploadSpeed    *int64  `json:"avg_upload_speed"`
	AvgDownloadSpeed  *int64  `json:"avg_download_speed"`
}
