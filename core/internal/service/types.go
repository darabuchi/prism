package service

import (
	"time"

	"github.com/prism/core/internal/storage"
)

// CreateSubscriptionRequest 创建订阅请求
type CreateSubscriptionRequest struct {
	Name           string `json:"name" binding:"required"`
	URL            string `json:"url" binding:"required"`
	UserAgent      string `json:"user_agent"`
	AutoUpdate     bool   `json:"auto_update"`
	UpdateInterval int    `json:"update_interval"`
	NodePoolIDs    []uint `json:"node_pool_ids"`
}

// UpdateSubscriptionRequest 更新订阅请求
type UpdateSubscriptionRequest struct {
	Name           *string `json:"name"`
	UserAgent      *string `json:"user_agent"`
	AutoUpdate     *bool   `json:"auto_update"`
	UpdateInterval *int    `json:"update_interval"`
	Status         *string `json:"status"`
}

// ListSubscriptionsRequest 订阅列表请求
type ListSubscriptionsRequest struct {
	Page       int    `form:"page,default=1" binding:"min=1"`
	Size       int    `form:"size,default=20" binding:"min=1,max=100"`
	Status     string `form:"status"`
	AutoUpdate *bool  `form:"auto_update"`
}

// ListSubscriptionsResponse 订阅列表响应
type ListSubscriptionsResponse struct {
	Total         int                    `json:"total"`
	Page          int                    `json:"page"`
	Size          int                    `json:"size"`
	Subscriptions []storage.Subscription `json:"subscriptions"`
}

// UpdateResult 更新结果
type UpdateResult struct {
	SubscriptionID uint      `json:"subscription_id"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	Duration       int       `json:"duration"`         // 毫秒
	TotalFetched   int       `json:"total_fetched"`    // 获取到的节点总数
	ValidNodes     int       `json:"valid_nodes"`      // 有效节点数
	NewNodes       int       `json:"new_nodes"`        // 新增节点数(订阅内去重)
	GlobalNewNodes int       `json:"global_new_nodes"` // 全局新增节点数
	UpdatedNodes   int       `json:"updated_nodes"`    // 更新的节点数
	RemovedNodes   int       `json:"removed_nodes"`    // 移除的节点数
}

// ProcessNodesResult 处理节点结果
type ProcessNodesResult struct {
	TotalFetched   int `json:"total_fetched"`
	ValidNodes     int `json:"valid_nodes"`
	NewNodes       int `json:"new_nodes"`
	GlobalNewNodes int `json:"global_new_nodes"`
	UpdatedNodes   int `json:"updated_nodes"`
	RemovedNodes   int `json:"removed_nodes"`
	TotalNodes     int `json:"total_nodes"`
	ActiveNodes    int `json:"active_nodes"`
}

// SubscriptionStats 订阅统计信息
type SubscriptionStats struct {
	SubscriptionID       uint                      `json:"subscription_id"`
	TotalNodes           int                       `json:"total_nodes"`
	ActiveNodes          int                       `json:"active_nodes"`
	SurvivalRate         float64                   `json:"survival_rate"`
	ProtocolDistribution []interface{}             `json:"protocol_distribution"`
	CountryDistribution  []interface{}             `json:"country_distribution"`
	RecentLogs           []storage.SubscriptionLog `json:"recent_logs"`
}

// LogsRequest 日志查询请求
type LogsRequest struct {
	Page       int       `form:"page,default=1" binding:"min=1"`
	Size       int       `form:"size,default=20" binding:"min=1,max=100"`
	Success    *bool     `form:"success"`
	UpdateType string    `form:"update_type"`
	StartTime  time.Time `form:"start_time"`
	EndTime    time.Time `form:"end_time"`
}

// LogsResponse 日志响应
type LogsResponse struct {
	Total int                       `json:"total"`
	Page  int                       `json:"page"`
	Size  int                       `json:"size"`
	Logs  []storage.SubscriptionLog `json:"logs"`
}

// CreateNodePoolRequest 创建节点池请求
type CreateNodePoolRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	Priority    int    `json:"priority"`
}

// UpdateNodePoolRequest 更新节点池请求
type UpdateNodePoolRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Enabled     *bool   `json:"enabled"`
	Priority    *int    `json:"priority"`
}

// ListNodePoolsResponse 节点池列表响应
type ListNodePoolsResponse struct {
	NodePools []storage.NodePool `json:"node_pools"`
}

// AssociateSubscriptionsRequest 关联订阅请求
type AssociateSubscriptionsRequest struct {
	SubscriptionIDs []uint `json:"subscription_ids" binding:"required"`
	Enabled         bool   `json:"enabled"`
	Priority        int    `json:"priority"`
}

// ListNodesRequest 节点列表请求
type ListNodesRequest struct {
	Page            int    `form:"page,default=1" binding:"min=1"`
	Size            int    `form:"size,default=20" binding:"min=1,max=100"`
	SubscriptionID  *uint  `form:"subscription_id"`
	NodePoolID      *uint  `form:"node_pool_id"`
	Country         string `form:"country"`
	Protocol        string `form:"protocol"`
	Status          string `form:"status"`
	Sort            string `form:"sort"`  // delay/upload_speed/download_speed/last_test
	Order           string `form:"order"` // asc/desc
	MinDelay        *int   `form:"min_delay"`
	MaxDelay        *int   `form:"max_delay"`
	StreamingUnlock string `form:"streaming_unlock"` // netflix/youtube/chatgpt等
}

// ListNodesResponse 节点列表响应
type ListNodesResponse struct {
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
	Nodes []storage.Node `json:"nodes"`
}

// TestNodeRequest 测试节点请求
type TestNodeRequest struct {
	TestTypes  []string               `json:"test_types" binding:"required"`
	TestConfig map[string]interface{} `json:"test_config"`
}

// BatchTestRequest 批量测试请求
type BatchTestRequest struct {
	NodeIDs    []uint                 `json:"node_ids" binding:"required"`
	TestTypes  []string               `json:"test_types" binding:"required"`
	TestConfig map[string]interface{} `json:"test_config"`
}

// BatchTestResponse 批量测试响应
type BatchTestResponse struct {
	TaskID     string `json:"task_id"`
	TotalNodes int    `json:"total_nodes"`
	Status     string `json:"status"`
}

// TestTaskStatus 测试任务状态
type TestTaskStatus struct {
	TaskID      string               `json:"task_id"`
	Status      string               `json:"status"`
	Progress    int                  `json:"progress"`
	Total       int                  `json:"total"`
	Results     []storage.TestResult `json:"results"`
	StartedAt   time.Time            `json:"started_at"`
	CompletedAt *time.Time           `json:"completed_at,omitempty"`
	Error       string               `json:"error,omitempty"`
}

// NodeTestHistoryRequest 节点测试历史请求
type NodeTestHistoryRequest struct {
	TestType  string    `form:"test_type"`
	StartTime time.Time `form:"start_time"`
	EndTime   time.Time `form:"end_time"`
	Limit     int       `form:"limit,default=100" binding:"max=500"`
}

// NodeTestHistoryResponse 节点测试历史响应
type NodeTestHistoryResponse struct {
	Tests []storage.NodeTest `json:"tests"`
}

// BestSelectionRequest 最佳节点选择请求
type BestSelectionRequest struct {
	NodePoolID      *uint    `form:"node_pool_id"`
	Country         string   `form:"country"`
	Protocol        string   `form:"protocol"`
	StreamingUnlock []string `form:"streaming_unlock"`
	Count           int      `form:"count,default=1" binding:"min=1,max=10"`
}

// NodeSelection 节点选择结果
type NodeSelection struct {
	Node            storage.Node   `json:"node"`
	SelectionReason map[string]int `json:"selection_reason"`
}

// BestSelectionResponse 最佳节点选择响应
type BestSelectionResponse struct {
	Selections []NodeSelection `json:"selections"`
}

// OverviewStats 整体统计
type OverviewStats struct {
	TotalSubscriptions   int     `json:"total_subscriptions"`
	ActiveSubscriptions  int     `json:"active_subscriptions"`
	TotalNodePools       int     `json:"total_node_pools"`
	TotalNodes           int     `json:"total_nodes"`
	ActiveNodes          int     `json:"active_nodes"`
	OverallSurvivalRate  float64 `json:"overall_survival_rate"`
	TotalTestsToday      int     `json:"total_tests_today"`
	SuccessfulTestsToday int     `json:"successful_tests_today"`
}

// GeoDistribution 地区分布
type GeoDistribution struct {
	Country     string `json:"country"`
	CountryName string `json:"country_name"`
	NodeCount   int    `json:"node_count"`
	ActiveCount int    `json:"active_count"`
}

// ProtocolDistribution 协议分布
type ProtocolDistribution struct {
	Protocol    string `json:"protocol"`
	NodeCount   int    `json:"node_count"`
	ActiveCount int    `json:"active_count"`
}

// PerformanceTrendRequest 性能趋势请求
type PerformanceTrendRequest struct {
	Period     string `form:"period" binding:"required"` // hour/day/week/month
	NodePoolID *uint  `form:"node_pool_id"`
	Country    string `form:"country"`
}

// TrendDataPoint 趋势数据点
type TrendDataPoint struct {
	Timestamp    time.Time `json:"timestamp"`
	AvgDelay     float64   `json:"avg_delay"`
	AvgUpload    float64   `json:"avg_upload"`
	AvgDownload  float64   `json:"avg_download"`
	SurvivalRate float64   `json:"survival_rate"`
	TestCount    int       `json:"test_count"`
}

// PerformanceTrendResponse 性能趋势响应
type PerformanceTrendResponse struct {
	Period     string           `json:"period"`
	DataPoints []TrendDataPoint `json:"data_points"`
}
