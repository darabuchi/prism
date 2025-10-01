package node

import (
	"fmt"
)

// Node 代理节点
type Node struct {
	id        string         // 唯一标识（SHA256）
	uniqueKey string         // 唯一键
	config    map[string]any // 完整配置
	info      *Info          // 节点信息
}

// New 创建新节点
func New(config map[string]any) (*Node, error) {
	if config == nil || len(config) == 0 {
		return nil, fmt.Errorf("config is empty")
	}

	// 生成 ID
	id, err := GenerateID(config)
	if err != nil {
		return nil, fmt.Errorf("generate id failed: %w", err)
	}

	// 生成唯一键
	uniqueKey := GenerateUniqueKey(config)
	if uniqueKey == "" {
		return nil, fmt.Errorf("generate unique key failed: missing type/server/port")
	}

	// 编码配置
	rawJSON, err := EncodeConfig(config)
	if err != nil {
		return nil, fmt.Errorf("encode config failed: %w", err)
	}

	// 提取基本信息
	name, _ := config["name"].(string)
	proxyType, _ := config["type"].(string)
	server, _ := config["server"].(string)
	port, _ := config["port"].(float64)

	node := &Node{
		id:        id,
		uniqueKey: uniqueKey,
		config:    config,
		info: &Info{
			ID:        id,
			UniqueKey: uniqueKey,
			Config: Config{
				Name:    name,
				Type:    proxyType,
				Server:  server,
				Port:    int(port),
				Config:  config,
				RawJSON: rawJSON,
			},
			Status: Status{
				Alive:       false,
				Enabled:     true,
				HealthState: HealthStateOpen,
			},
			Score: 0,
		},
	}

	return node, nil
}

// NewFromEncoded 从编码的配置创建节点
func NewFromEncoded(encoded string) (*Node, error) {
	config, err := DecodeConfig(encoded)
	if err != nil {
		return nil, fmt.Errorf("decode config failed: %w", err)
	}

	return New(config)
}

// ID 返回节点唯一标识
func (n *Node) ID() string {
	return n.id
}

// UniqueKey 返回节点唯一键
func (n *Node) UniqueKey() string {
	return n.uniqueKey
}

// Config 返回节点配置
func (n *Node) Config() map[string]any {
	return n.config
}

// Info 返回节点完整信息
func (n *Node) Info() *Info {
	return n.info
}

// Name 返回节点名称
func (n *Node) Name() string {
	return n.info.Config.Name
}

// Type 返回节点类型
func (n *Node) Type() string {
	return n.info.Config.Type
}

// Server 返回服务器地址
func (n *Node) Server() string {
	return n.info.Config.Server
}

// Port 返回端口
func (n *Node) Port() int {
	return n.info.Config.Port
}

// Address 返回完整地址 (server:port)
func (n *Node) Address() string {
	return fmt.Sprintf("%s:%d", n.info.Config.Server, n.info.Config.Port)
}

// IsAlive 返回节点是否存活
func (n *Node) IsAlive() bool {
	return n.info.Status.Alive
}

// IsEnabled 返回节点是否启用
func (n *Node) IsEnabled() bool {
	return n.info.Status.Enabled
}

// SetAlive 设置节点存活状态
func (n *Node) SetAlive(alive bool) {
	n.info.Status.Alive = alive
}

// SetEnabled 设置节点启用状态
func (n *Node) SetEnabled(enabled bool) {
	n.info.Status.Enabled = enabled
}

// SetHealthState 设置熔断器健康状态
func (n *Node) SetHealthState(state HealthState) {
	n.info.Status.HealthState = state
}

// UpdateStatus 更新节点状态
func (n *Node) UpdateStatus(status Status) {
	n.info.Status = status
}

// UpdateTest 更新测试结果
func (n *Node) UpdateTest(test TestResult) {
	n.info.Test = test
}

// SetScore 设置节点评分
func (n *Node) SetScore(score float64) {
	n.info.Score = score
}

// Clone 克隆节点
func (n *Node) Clone() *Node {
	// 深拷贝配置
	config := make(map[string]any, len(n.config))
	for k, v := range n.config {
		config[k] = v
	}

	// 深拷贝 info
	infoCopy := *n.info
	infoCopy.Config.Config = config

	return &Node{
		id:        n.id,
		uniqueKey: n.uniqueKey,
		config:    config,
		info:      &infoCopy,
	}
}
