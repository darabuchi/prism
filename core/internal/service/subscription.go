package service

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
	"gorm.io/gorm"

	"github.com/prism/core/internal/storage"
)

// SubscriptionService 订阅服务
type SubscriptionService struct {
	db *storage.Database
}

// NewSubscriptionService 创建订阅服务
func NewSubscriptionService(db *storage.Database) *SubscriptionService {
	return &SubscriptionService{db: db}
}

// CreateSubscription 创建订阅
func (s *SubscriptionService) CreateSubscription(req *CreateSubscriptionRequest) (*storage.Subscription, error) {
	subscription := &storage.Subscription{
		Name:           req.Name,
		URL:            req.URL,
		UserAgent:      req.UserAgent,
		AutoUpdate:     req.AutoUpdate,
		UpdateInterval: req.UpdateInterval,
	}

	// 验证订阅链接有效性
	if err := s.validateSubscriptionURL(subscription.URL); err != nil {
		return nil, fmt.Errorf("invalid subscription URL: %w", err)
	}

	// 保存到数据库
	if err := s.db.Create(subscription).Error; err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	// 关联到节点池
	if len(req.NodePoolIDs) > 0 {
		if err := s.associateNodePools(subscription.ID, req.NodePoolIDs); err != nil {
			return nil, fmt.Errorf("failed to associate node pools: %w", err)
		}
	}

	return subscription, nil
}

// UpdateSubscription 更新订阅
func (s *SubscriptionService) UpdateSubscription(id uint, req *UpdateSubscriptionRequest) (*storage.Subscription, error) {
	var subscription storage.Subscription
	if err := s.db.First(&subscription, id).Error; err != nil {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}

	// 更新字段
	if req.Name != nil {
		subscription.Name = *req.Name
	}
	if req.UserAgent != nil {
		subscription.UserAgent = *req.UserAgent
	}
	if req.AutoUpdate != nil {
		subscription.AutoUpdate = *req.AutoUpdate
	}
	if req.UpdateInterval != nil {
		subscription.UpdateInterval = *req.UpdateInterval
	}
	if req.Status != nil {
		subscription.Status = *req.Status
	}

	if err := s.db.Save(&subscription).Error; err != nil {
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	return &subscription, nil
}

// GetSubscription 获取单个订阅
func (s *SubscriptionService) GetSubscription(id uint) (*storage.Subscription, error) {
	var subscription storage.Subscription
	if err := s.db.Preload("NodePools").Preload("Nodes").First(&subscription, id).Error; err != nil {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}
	return &subscription, nil
}

// ListSubscriptions 获取订阅列表
func (s *SubscriptionService) ListSubscriptions(req *ListSubscriptionsRequest) (*ListSubscriptionsResponse, error) {
	var subscriptions []storage.Subscription
	var total int64

	query := s.db.Model(&storage.Subscription{})

	// 应用过滤条件
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.AutoUpdate != nil {
		query = query.Where("auto_update = ?", *req.AutoUpdate)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count subscriptions: %w", err)
	}

	// 分页查询
	offset := (req.Page - 1) * req.Size
	if err := query.Offset(offset).Limit(req.Size).Order("created_at DESC").Find(&subscriptions).Error; err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}

	return &ListSubscriptionsResponse{
		Total:         int(total),
		Page:          req.Page,
		Size:          req.Size,
		Subscriptions: subscriptions,
	}, nil
}

// DeleteSubscription 删除订阅
func (s *SubscriptionService) DeleteSubscription(id uint) error {
	result := s.db.Delete(&storage.Subscription{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete subscription: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}
	return nil
}

// UpdateSubscriptionContent 手动更新订阅内容
func (s *SubscriptionService) UpdateSubscriptionContent(id uint) (*UpdateResult, error) {
	var subscription storage.Subscription
	if err := s.db.First(&subscription, id).Error; err != nil {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}

	startTime := time.Now()
	result := &UpdateResult{
		SubscriptionID: id,
		StartTime:      startTime,
	}

	// 记录更新日志
	log := &storage.SubscriptionLog{
		SubscriptionID: id,
		UpdateType:     "manual",
	}

	// 获取订阅内容
	content, httpStatus, responseTime, err := s.fetchSubscriptionContent(subscription.URL, subscription.UserAgent)
	log.HTTPStatus = &httpStatus
	log.ResponseTime = &responseTime

	if err != nil {
		log.Success = false
		log.ErrorMessage = err.Error()
		s.db.Create(log)

		// 更新订阅错误信息
		subscription.ErrorCount++
		subscription.ErrorMessage = err.Error()
		s.db.Save(&subscription)

		return nil, fmt.Errorf("failed to fetch subscription: %w", err)
	}

	// 解析节点配置
	nodes, err := s.parseSubscriptionContent(content)
	if err != nil {
		log.Success = false
		log.ErrorMessage = err.Error()
		s.db.Create(log)
		return nil, fmt.Errorf("failed to parse subscription content: %w", err)
	}

	log.TotalFetched = len(nodes)

	// 处理节点数据
	updateResult, err := s.processNodes(&subscription, nodes)
	if err != nil {
		log.Success = false
		log.ErrorMessage = err.Error()
		s.db.Create(log)
		return nil, fmt.Errorf("failed to process nodes: %w", err)
	}

	// 更新日志统计
	log.Success = true
	log.ValidNodes = updateResult.ValidNodes
	log.NewNodes = updateResult.NewNodes
	log.GlobalNewNodes = updateResult.GlobalNewNodes
	log.UpdatedNodes = updateResult.UpdatedNodes
	log.RemovedNodes = updateResult.RemovedNodes
	s.db.Create(log)

	// 更新订阅信息
	now := time.Now()
	subscription.LastUpdate = &now
	subscription.LastSuccess = &now
	subscription.ErrorCount = 0
	subscription.ErrorMessage = ""
	subscription.TotalNodes = updateResult.TotalNodes
	subscription.ActiveNodes = updateResult.ActiveNodes
	subscription.UniqueNewNodes = updateResult.GlobalNewNodes
	s.db.Save(&subscription)

	result.EndTime = time.Now()
	result.Duration = int(result.EndTime.Sub(startTime).Milliseconds())
	result.TotalFetched = updateResult.TotalFetched
	result.ValidNodes = updateResult.ValidNodes
	result.NewNodes = updateResult.NewNodes
	result.GlobalNewNodes = updateResult.GlobalNewNodes
	result.UpdatedNodes = updateResult.UpdatedNodes
	result.RemovedNodes = updateResult.RemovedNodes

	return result, nil
}

// GetSubscriptionStats 获取订阅统计信息
func (s *SubscriptionService) GetSubscriptionStats(id uint) (*SubscriptionStats, error) {
	var subscription storage.Subscription
	if err := s.db.First(&subscription, id).Error; err != nil {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}

	stats := &SubscriptionStats{
		SubscriptionID: id,
		TotalNodes:     subscription.TotalNodes,
		ActiveNodes:    subscription.ActiveNodes,
		SurvivalRate:   0,
	}

	if subscription.TotalNodes > 0 {
		stats.SurvivalRate = float64(subscription.ActiveNodes) / float64(subscription.TotalNodes) * 100
	}

	// 获取协议分布
	var protocolStats []struct {
		Protocol string `json:"protocol"`
		Count    int    `json:"count"`
	}
	s.db.Model(&storage.Node{}).
		Where("subscription_id = ?", id).
		Select("protocol, COUNT(*) as count").
		Group("protocol").
		Scan(&protocolStats)
	// 转换 protocolStats 为 []interface{}
	protocolInterfaces := make([]interface{}, len(protocolStats))
	for i, v := range protocolStats {
		protocolInterfaces[i] = v
	}
	stats.ProtocolDistribution = protocolInterfaces

	// 获取地区分布
	var countryStats []struct {
		Country string `json:"country"`
		Count   int    `json:"count"`
	}
	s.db.Model(&storage.Node{}).
		Where("subscription_id = ? AND country IS NOT NULL", id).
		Select("country, COUNT(*) as count").
		Group("country").
		Scan(&countryStats)
	// 转换 countryStats 为 []interface{}
	countryInterfaces := make([]interface{}, len(countryStats))
	for i, v := range countryStats {
		countryInterfaces[i] = v
	}
	stats.CountryDistribution = countryInterfaces

	// 获取最近更新记录
	s.db.Where("subscription_id = ?", id).
		Order("created_at DESC").
		Limit(10).
		Find(&stats.RecentLogs)

	return stats, nil
}

// GetSubscriptionLogs 获取订阅日志
func (s *SubscriptionService) GetSubscriptionLogs(id uint, req *LogsRequest) (*LogsResponse, error) {
	var logs []storage.SubscriptionLog
	var total int64

	query := s.db.Model(&storage.SubscriptionLog{}).Where("subscription_id = ?", id)

	// 应用过滤条件
	if req.Success != nil {
		query = query.Where("success = ?", *req.Success)
	}
	if req.UpdateType != "" {
		query = query.Where("update_type = ?", req.UpdateType)
	}
	if !req.StartTime.IsZero() {
		query = query.Where("created_at >= ?", req.StartTime)
	}
	if !req.EndTime.IsZero() {
		query = query.Where("created_at <= ?", req.EndTime)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count logs: %w", err)
	}

	// 分页查询
	offset := (req.Page - 1) * req.Size
	if err := query.Offset(offset).Limit(req.Size).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("failed to list logs: %w", err)
	}

	return &LogsResponse{
		Total: int(total),
		Page:  req.Page,
		Size:  req.Size,
		Logs:  logs,
	}, nil
}

// validateSubscriptionURL 验证订阅链接
func (s *SubscriptionService) validateSubscriptionURL(rawURL string) error {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("unsupported URL scheme: %s", parsedURL.Scheme)
	}

	return nil
}

// fetchSubscriptionContent 获取订阅内容
func (s *SubscriptionService) fetchSubscriptionContent(url, userAgent string) (string, int, int, error) {
	startTime := time.Now()

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", 0, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "*/*")

	resp, err := client.Do(req)
	if err != nil {
		return "", 0, int(time.Since(startTime).Milliseconds()), fmt.Errorf("failed to fetch subscription: %w", err)
	}
	defer resp.Body.Close()

	responseTime := int(time.Since(startTime).Milliseconds())

	if resp.StatusCode != http.StatusOK {
		return "", resp.StatusCode, responseTime, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode, responseTime, fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), resp.StatusCode, responseTime, nil
}

// parseSubscriptionContent 解析订阅内容
func (s *SubscriptionService) parseSubscriptionContent(content string) ([]map[string]interface{}, error) {
	// 尝试 base64 解码
	if decoded, err := base64.StdEncoding.DecodeString(content); err == nil {
		content = string(decoded)
	}

	var config struct {
		Proxies []map[string]interface{} `yaml:"proxies"`
	}

	// 尝试解析 YAML
	if err := yaml.Unmarshal([]byte(content), &config); err == nil && len(config.Proxies) > 0 {
		return config.Proxies, nil
	}

	// 尝试解析单个代理链接
	lines := strings.Split(strings.TrimSpace(content), "\n")
	var proxies []map[string]interface{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		proxy, err := s.parseProxyURL(line)
		if err != nil {
			continue // 跳过无效的代理链接
		}

		proxies = append(proxies, proxy)
	}

	if len(proxies) == 0 {
		return nil, fmt.Errorf("no valid proxy configurations found")
	}

	return proxies, nil
}

// parseProxyURL 解析单个代理链接
func (s *SubscriptionService) parseProxyURL(proxyURL string) (map[string]interface{}, error) {
	// 这里应该实现各种代理协议的解析逻辑
	// 为了简化，这里只返回一个示例配置
	return map[string]interface{}{
		"name":   "Parsed Node",
		"type":   "ss",
		"server": "example.com",
		"port":   443,
	}, nil
}

// associateNodePools 关联节点池
func (s *SubscriptionService) associateNodePools(subscriptionID uint, nodePoolIDs []uint) error {
	for _, poolID := range nodePoolIDs {
		association := &storage.NodePoolSubscription{
			NodePoolID:     poolID,
			SubscriptionID: subscriptionID,
			Enabled:        true,
			Priority:       0,
		}
		if err := s.db.Create(association).Error; err != nil {
			return fmt.Errorf("failed to associate with pool %d: %w", poolID, err)
		}
	}
	return nil
}

// processNodes 处理节点数据
func (s *SubscriptionService) processNodes(subscription *storage.Subscription, nodeConfigs []map[string]interface{}) (*ProcessNodesResult, error) {
	result := &ProcessNodesResult{}

	// 在事务中处理所有节点操作
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 获取现有节点哈希
		var existingHashes []string
		tx.Model(&storage.Node{}).
			Where("subscription_id = ?", subscription.ID).
			Pluck("hash", &existingHashes)
		existingHashMap := make(map[string]bool)
		for _, hash := range existingHashes {
			existingHashMap[hash] = true
		}

		// 处理每个节点配置
		newHashes := make(map[string]bool)
		for _, config := range nodeConfigs {
			result.TotalFetched++

			// 验证节点配置
			node, err := s.validateAndCreateNode(subscription.ID, config)
			if err != nil {
				continue // 跳过无效节点
			}

			result.ValidNodes++
			newHashes[node.Hash] = true

			// 检查是否已存在（订阅内去重）
			if existingHashMap[node.Hash] {
				result.UpdatedNodes++
				// 更新现有节点
				tx.Model(&storage.Node{}).
					Where("hash = ? AND subscription_id = ?", node.Hash, subscription.ID).
					Updates(map[string]interface{}{
						"name":         node.Name,
						"clash_config": node.ClashConfig,
						"updated_at":   time.Now(),
					})
			} else {
				// 检查全局去重
				var globalCount int64
				tx.Model(&storage.Node{}).Where("hash = ?", node.Hash).Count(&globalCount)

				if globalCount == 0 {
					result.GlobalNewNodes++
				}

				result.NewNodes++
				// 插入新节点
				if err := tx.Create(node).Error; err != nil {
					return fmt.Errorf("failed to create node: %w", err)
				}
			}
		}

		// 删除不再存在的节点
		var removedCount int64
		tx.Model(&storage.Node{}).
			Where("subscription_id = ? AND hash NOT IN ?", subscription.ID, getKeys(newHashes)).
			Count(&removedCount)
		result.RemovedNodes = int(removedCount)

		tx.Where("subscription_id = ? AND hash NOT IN ?", subscription.ID, getKeys(newHashes)).
			Delete(&storage.Node{})

		// 统计节点数量
		var totalNodes, activeNodes int64
		tx.Model(&storage.Node{}).Where("subscription_id = ?", subscription.ID).Count(&totalNodes)
		tx.Model(&storage.Node{}).Where("subscription_id = ? AND status = 'online'", subscription.ID).Count(&activeNodes)

		result.TotalNodes = int(totalNodes)
		result.ActiveNodes = int(activeNodes)

		return nil
	})

	return result, err
}

// validateAndCreateNode 验证并创建节点
func (s *SubscriptionService) validateAndCreateNode(subscriptionID uint, config map[string]interface{}) (*storage.Node, error) {
	// 验证必需字段
	name, ok := config["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("invalid node name")
	}

	server, ok := config["server"].(string)
	if !ok || server == "" {
		return nil, fmt.Errorf("invalid server")
	}

	port, ok := config["port"].(int)
	if !ok || port <= 0 || port > 65535 {
		return nil, fmt.Errorf("invalid port")
	}

	protocol, ok := config["type"].(string)
	if !ok || protocol == "" {
		return nil, fmt.Errorf("invalid protocol")
	}

	// 计算节点哈希
	hash := s.calculateNodeHash(config)

	node := &storage.Node{
		SubscriptionID: subscriptionID,
		Name:           name,
		Hash:           hash,
		Server:         server,
		Port:           port,
		Protocol:       protocol,
		ClashConfig:    storage.JSON(config),
		Status:         "unknown",
	}

	// 解析地理信息（如果节点名称包含地区信息）
	s.parseNodeGeography(node)

	return node, nil
}

// calculateNodeHash 计算节点哈希用于去重
func (s *SubscriptionService) calculateNodeHash(config map[string]interface{}) string {
	// 简化的哈希计算，实际应该使用更复杂的算法
	server := config["server"].(string)
	port := config["port"].(int)
	protocol := config["type"].(string)
	return fmt.Sprintf("%s:%d:%s", server, port, protocol)
}

// parseNodeGeography 解析节点地理信息
func (s *SubscriptionService) parseNodeGeography(node *storage.Node) {
	name := strings.ToUpper(node.Name)

	// 简单的地区解析逻辑
	countryMap := map[string][]string{
		"HK": {"香港", "HK", "HONG KONG"},
		"US": {"美国", "US", "USA", "UNITED STATES"},
		"JP": {"日本", "JP", "JAPAN"},
		"SG": {"新加坡", "SG", "SINGAPORE"},
		"UK": {"英国", "UK", "UNITED KINGDOM"},
		"DE": {"德国", "DE", "GERMANY"},
	}

	for code, keywords := range countryMap {
		for _, keyword := range keywords {
			if strings.Contains(name, keyword) {
				node.Country = code
				break
			}
		}
		if node.Country != "" {
			break
		}
	}
}

// getKeys 获取 map 的 keys
func getKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
