package storage

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// JSON 自定义类型，支持 JSON 字段存储
type JSON map[string]interface{}

func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("cannot scan non-[]byte value into JSON field")
	}

	return json.Unmarshal(bytes, j)
}

// Subscription 订阅表模型
type Subscription struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string `gorm:"size:255;not null;comment:订阅名称" json:"name"`
	URL       string `gorm:"type:text;not null;comment:订阅链接" json:"url"`
	UserAgent string `gorm:"size:255;default:'clash';comment:请求User-Agent" json:"user_agent"`

	// 更新配置
	AutoUpdate     bool `gorm:"default:true;comment:自动更新开关" json:"auto_update"`
	UpdateInterval int  `gorm:"default:3600;comment:更新间隔(秒)" json:"update_interval"`

	// 统计信息
	TotalNodes     int `gorm:"default:0;comment:当前订阅总节点数" json:"total_nodes"`
	ActiveNodes    int `gorm:"default:0;comment:存活节点数" json:"active_nodes"`
	UniqueNewNodes int `gorm:"default:0;comment:全局去重后新增节点数" json:"unique_new_nodes"`

	// 状态信息
	Status       string     `gorm:"size:50;default:'active';comment:订阅状态" json:"status"`
	LastUpdate   *time.Time `gorm:"comment:最后更新时间" json:"last_update"`
	LastSuccess  *time.Time `gorm:"comment:最后成功时间" json:"last_success"`
	ErrorMessage string     `gorm:"type:text;comment:最后错误信息" json:"error_message"`
	ErrorCount   int        `gorm:"default:0;comment:连续错误次数" json:"error_count"`

	// 时间戳
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联
	Nodes     []Node            `gorm:"foreignKey:SubscriptionID;constraint:OnDelete:CASCADE" json:"nodes,omitempty"`
	NodePools []NodePool        `gorm:"many2many:node_pool_subscriptions" json:"node_pools,omitempty"`
	Logs      []SubscriptionLog `gorm:"foreignKey:SubscriptionID;constraint:OnDelete:CASCADE" json:"logs,omitempty"`
}

// TableName 指定表名
func (Subscription) TableName() string {
	return "subscriptions"
}

// BeforeCreate GORM 钩子
func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
	// 设置默认值
	if s.UserAgent == "" {
		s.UserAgent = "clash"
	}
	if s.UpdateInterval == 0 {
		s.UpdateInterval = 3600
	}
	if s.Status == "" {
		s.Status = "active"
	}
	return nil
}

// NodePool 节点池表模型
type NodePool struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"size:255;not null;uniqueIndex;comment:节点池名称" json:"name"`
	Description string `gorm:"type:text;comment:节点池描述" json:"description"`

	// 统计信息
	TotalSubscriptions int     `gorm:"default:0;comment:关联订阅数" json:"total_subscriptions"`
	TotalNodes         int     `gorm:"default:0;comment:节点总数" json:"total_nodes"`
	ActiveNodes        int     `gorm:"default:0;comment:存活节点数" json:"active_nodes"`
	SurvivalRate       float64 `gorm:"type:decimal(5,2);default:0.00;comment:存活率(%)" json:"survival_rate"`

	// 配置信息
	Enabled  bool `gorm:"default:true;comment:启用状态" json:"enabled"`
	Priority int  `gorm:"default:0;comment:优先级" json:"priority"`

	// 时间戳
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联
	Nodes         []Node         `gorm:"foreignKey:NodePoolID;constraint:OnDelete:SET NULL" json:"nodes,omitempty"`
	Subscriptions []Subscription `gorm:"many2many:node_pool_subscriptions" json:"subscriptions,omitempty"`
}

// TableName 指定表名
func (NodePool) TableName() string {
	return "node_pools"
}

// UpdateStats 更新统计信息
func (np *NodePool) UpdateStats(db *gorm.DB) error {
	var stats struct {
		TotalNodes  int64
		ActiveNodes int64
	}

	err := db.Model(&Node{}).
		Where("node_pool_id = ?", np.ID).
		Select("COUNT(*) as total_nodes, COUNT(CASE WHEN status = 'online' THEN 1 END) as active_nodes").
		Scan(&stats).Error

	if err != nil {
		return err
	}

	np.TotalNodes = int(stats.TotalNodes)
	np.ActiveNodes = int(stats.ActiveNodes)

	if np.TotalNodes > 0 {
		np.SurvivalRate = float64(np.ActiveNodes) / float64(np.TotalNodes) * 100
	} else {
		np.SurvivalRate = 0
	}

	return db.Save(np).Error
}

// Node 节点表模型
type Node struct {
	ID             uint  `gorm:"primaryKey;autoIncrement" json:"id"`
	SubscriptionID uint  `gorm:"not null;index;comment:所属订阅ID" json:"subscription_id"`
	NodePoolID     *uint `gorm:"index;comment:所属节点池ID" json:"node_pool_id"`

	// 节点基础信息
	Name     string `gorm:"size:255;not null;comment:节点名称" json:"name"`
	Hash     string `gorm:"size:64;not null;uniqueIndex;comment:节点哈希(用于去重)" json:"hash"`
	Server   string `gorm:"size:255;not null;index;comment:服务器地址" json:"server"`
	Port     int    `gorm:"not null;comment:端口" json:"port"`
	Protocol string `gorm:"size:50;not null;index;comment:协议类型" json:"protocol"`

	// clash 配置 (JSON 存储)
	ClashConfig JSON `gorm:"type:json;not null;comment:Clash完整配置" json:"clash_config"`

	// 地理信息
	Country     string `gorm:"size:10;index;comment:国家代码" json:"country"`
	CountryName string `gorm:"size:100;comment:国家名称" json:"country_name"`
	City        string `gorm:"size:100;comment:城市" json:"city"`
	ISP         string `gorm:"size:100;comment:运营商" json:"isp"`

	// 性能指标 (最新测试结果)
	Delay         *int     `gorm:"index;comment:延迟(ms)" json:"delay"`
	UploadSpeed   *int64   `gorm:"comment:上传速度(bps)" json:"upload_speed"`
	DownloadSpeed *int64   `gorm:"comment:下载速度(bps)" json:"download_speed"`
	LossRate      *float64 `gorm:"type:decimal(5,2);comment:丢包率(%)" json:"loss_rate"`

	// 连通性状态
	Status             string     `gorm:"size:50;default:'unknown';index;comment:状态" json:"status"`
	LastTest           *time.Time `gorm:"index;comment:最后测试时间" json:"last_test"`
	LastOnline         *time.Time `gorm:"comment:最后在线时间" json:"last_online"`
	ContinuousFailures int        `gorm:"default:0;comment:连续失败次数" json:"continuous_failures"`

	// 流媒体解锁信息 (JSON 存储)
	StreamingUnlock JSON `gorm:"type:json;comment:流媒体解锁情况" json:"streaming_unlock"`

	// 时间戳
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联
	Subscription *Subscription `gorm:"foreignKey:SubscriptionID;constraint:OnDelete:CASCADE" json:"subscription,omitempty"`
	NodePool     *NodePool     `gorm:"foreignKey:NodePoolID;constraint:OnDelete:SET NULL" json:"node_pool,omitempty"`
	Tests        []NodeTest    `gorm:"foreignKey:NodeID;constraint:OnDelete:CASCADE" json:"tests,omitempty"`
}

// TableName 指定表名
func (Node) TableName() string {
	return "nodes"
}

// BeforeCreate GORM 钩子
func (n *Node) BeforeCreate(tx *gorm.DB) error {
	if n.Status == "" {
		n.Status = "unknown"
	}
	return nil
}

// IsOnline 检查节点是否在线
func (n *Node) IsOnline() bool {
	return n.Status == "online"
}

// UpdateStatus 更新节点状态
func (n *Node) UpdateStatus(status string, delay *int) {
	n.Status = status
	now := time.Now()
	n.LastTest = &now

	if status == "online" {
		n.LastOnline = &now
		n.ContinuousFailures = 0
		if delay != nil {
			n.Delay = delay
		}
	} else {
		n.ContinuousFailures++
	}
}

// NodeTest 节点测试记录表模型
type NodeTest struct {
	ID     uint `gorm:"primaryKey;autoIncrement" json:"id"`
	NodeID uint `gorm:"not null;index;comment:节点ID" json:"node_id"`

	// 测试类型和配置
	TestType   string `gorm:"size:50;not null;index;comment:测试类型" json:"test_type"`
	TestConfig JSON   `gorm:"type:json;comment:测试配置参数" json:"test_config"`

	// 测试结果
	Delay         *int     `gorm:"comment:延迟(ms)" json:"delay"`
	UploadSpeed   *int64   `gorm:"comment:上传速度(bps)" json:"upload_speed"`
	DownloadSpeed *int64   `gorm:"comment:下载速度(bps)" json:"download_speed"`
	LossRate      *float64 `gorm:"type:decimal(5,2);comment:丢包率(%)" json:"loss_rate"`

	// 流媒体解锁结果
	StreamingResults JSON `gorm:"type:json;comment:流媒体测试结果" json:"streaming_results"`

	// 测试状态
	Success      bool   `gorm:"not null;default:false;index;comment:测试是否成功" json:"success"`
	ErrorMessage string `gorm:"type:text;comment:错误信息" json:"error_message"`

	// 时间戳
	TestedAt time.Time `gorm:"autoCreateTime;index;comment:测试时间" json:"tested_at"`

	// 关联
	Node *Node `gorm:"foreignKey:NodeID;constraint:OnDelete:CASCADE" json:"node,omitempty"`
}

// TableName 指定表名
func (NodeTest) TableName() string {
	return "node_tests"
}

// SubscriptionLog 订阅更新日志表模型
type SubscriptionLog struct {
	ID             uint `gorm:"primaryKey;autoIncrement" json:"id"`
	SubscriptionID uint `gorm:"not null;index;comment:订阅ID" json:"subscription_id"`

	// 更新信息
	UpdateType string `gorm:"size:50;not null;index;comment:更新类型" json:"update_type"`

	// 更新结果
	Success bool `gorm:"not null;default:false;index;comment:更新是否成功" json:"success"`

	// 节点统计
	TotalFetched   int `gorm:"default:0;comment:获取到的节点总数" json:"total_fetched"`
	ValidNodes     int `gorm:"default:0;comment:有效节点数" json:"valid_nodes"`
	NewNodes       int `gorm:"default:0;comment:新增节点数(订阅内去重)" json:"new_nodes"`
	GlobalNewNodes int `gorm:"default:0;comment:全局新增节点数" json:"global_new_nodes"`
	UpdatedNodes   int `gorm:"default:0;comment:更新的节点数" json:"updated_nodes"`
	RemovedNodes   int `gorm:"default:0;comment:移除的节点数" json:"removed_nodes"`

	// 错误信息
	ErrorMessage string `gorm:"type:text;comment:错误信息" json:"error_message"`
	HTTPStatus   *int   `gorm:"comment:HTTP状态码" json:"http_status"`
	ResponseTime *int   `gorm:"comment:响应时间(ms)" json:"response_time"`

	// 时间戳
	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"created_at"`

	// 关联
	Subscription *Subscription `gorm:"foreignKey:SubscriptionID;constraint:OnDelete:CASCADE" json:"subscription,omitempty"`
}

// TableName 指定表名
func (SubscriptionLog) TableName() string {
	return "subscription_logs"
}

// NodePoolSubscription 节点池关联表模型
type NodePoolSubscription struct {
	ID             uint `gorm:"primaryKey;autoIncrement" json:"id"`
	NodePoolID     uint `gorm:"not null;index;comment:节点池ID" json:"node_pool_id"`
	SubscriptionID uint `gorm:"not null;index;comment:订阅ID" json:"subscription_id"`

	// 配置信息
	Enabled  bool `gorm:"default:true;comment:是否启用" json:"enabled"`
	Priority int  `gorm:"default:0;comment:优先级" json:"priority"`

	// 时间戳
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联
	NodePool     *NodePool     `gorm:"foreignKey:NodePoolID;constraint:OnDelete:CASCADE" json:"node_pool,omitempty"`
	Subscription *Subscription `gorm:"foreignKey:SubscriptionID;constraint:OnDelete:CASCADE" json:"subscription,omitempty"`
}

// TableName 指定表名
func (NodePoolSubscription) TableName() string {
	return "node_pool_subscriptions"
}

// BeforeCreate GORM 钩子，确保唯一约束
func (nps *NodePoolSubscription) BeforeCreate(tx *gorm.DB) error {
	var count int64
	err := tx.Model(&NodePoolSubscription{}).
		Where("node_pool_id = ? AND subscription_id = ?", nps.NodePoolID, nps.SubscriptionID).
		Count(&count).Error

	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("this subscription is already associated with the node pool")
	}

	return nil
}

// TestTask 测试任务模型（内存中使用，不持久化）
type TestTask struct {
	ID          string       `json:"id"`
	NodeIDs     []uint       `json:"node_ids"`
	TestTypes   []string     `json:"test_types"`
	Status      string       `json:"status"` // running/completed/failed
	Progress    int          `json:"progress"`
	Total       int          `json:"total"`
	Results     []TestResult `json:"results"`
	StartedAt   time.Time    `json:"started_at"`
	CompletedAt *time.Time   `json:"completed_at,omitempty"`
	Error       string       `json:"error,omitempty"`
}

// TestResult 单个节点测试结果
type TestResult struct {
	NodeID  uint                   `json:"node_id"`
	Success bool                   `json:"success"`
	Error   string                 `json:"error,omitempty"`
	Results map[string]interface{} `json:"results,omitempty"`
}
