package service

import (
	"fmt"

	"github.com/prism/core/internal/storage"
)

// NodePoolService 节点池服务
type NodePoolService struct {
	db *storage.Database
}

// NewNodePoolService 创建节点池服务
func NewNodePoolService(db *storage.Database) *NodePoolService {
	return &NodePoolService{db: db}
}

// CreateNodePool 创建节点池
func (s *NodePoolService) CreateNodePool(req *CreateNodePoolRequest) (*storage.NodePool, error) {
	nodePool := &storage.NodePool{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Priority:    req.Priority,
	}

	if err := s.db.Create(nodePool).Error; err != nil {
		return nil, fmt.Errorf("failed to create node pool: %w", err)
	}

	return nodePool, nil
}

// UpdateNodePool 更新节点池
func (s *NodePoolService) UpdateNodePool(id uint, req *UpdateNodePoolRequest) (*storage.NodePool, error) {
	var nodePool storage.NodePool
	if err := s.db.First(&nodePool, id).Error; err != nil {
		return nil, fmt.Errorf("node pool not found: %w", err)
	}

	// 更新字段
	if req.Name != nil {
		nodePool.Name = *req.Name
	}
	if req.Description != nil {
		nodePool.Description = *req.Description
	}
	if req.Enabled != nil {
		nodePool.Enabled = *req.Enabled
	}
	if req.Priority != nil {
		nodePool.Priority = *req.Priority
	}

	if err := s.db.Save(&nodePool).Error; err != nil {
		return nil, fmt.Errorf("failed to update node pool: %w", err)
	}

	return &nodePool, nil
}

// GetNodePool 获取单个节点池
func (s *NodePoolService) GetNodePool(id uint) (*storage.NodePool, error) {
	var nodePool storage.NodePool
	if err := s.db.Preload("Subscriptions").Preload("Nodes").First(&nodePool, id).Error; err != nil {
		return nil, fmt.Errorf("node pool not found: %w", err)
	}
	return &nodePool, nil
}

// ListNodePools 获取节点池列表
func (s *NodePoolService) ListNodePools() (*ListNodePoolsResponse, error) {
	var nodePools []storage.NodePool
	if err := s.db.Order("priority DESC, created_at DESC").Find(&nodePools).Error; err != nil {
		return nil, fmt.Errorf("failed to list node pools: %w", err)
	}

	return &ListNodePoolsResponse{
		NodePools: nodePools,
	}, nil
}

// DeleteNodePool 删除节点池
func (s *NodePoolService) DeleteNodePool(id uint) error {
	result := s.db.Delete(&storage.NodePool{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete node pool: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("node pool not found")
	}
	return nil
}

// AssociateSubscriptions 关联订阅到节点池
func (s *NodePoolService) AssociateSubscriptions(poolID uint, req *AssociateSubscriptionsRequest) error {
	// 验证节点池存在
	var nodePool storage.NodePool
	if err := s.db.First(&nodePool, poolID).Error; err != nil {
		return fmt.Errorf("node pool not found: %w", err)
	}

	// 删除现有关联
	if err := s.db.Where("node_pool_id = ?", poolID).Delete(&storage.NodePoolSubscription{}).Error; err != nil {
		return fmt.Errorf("failed to remove existing associations: %w", err)
	}

	// 创建新关联
	for _, subscriptionID := range req.SubscriptionIDs {
		association := &storage.NodePoolSubscription{
			NodePoolID:     poolID,
			SubscriptionID: subscriptionID,
			Enabled:        req.Enabled,
			Priority:       req.Priority,
		}
		if err := s.db.Create(association).Error; err != nil {
			return fmt.Errorf("failed to associate subscription %d: %w", subscriptionID, err)
		}
	}

	// 更新节点池统计
	if err := s.updateNodePoolStats(poolID); err != nil {
		return fmt.Errorf("failed to update node pool stats: %w", err)
	}

	return nil
}

// updateNodePoolStats 更新节点池统计信息
func (s *NodePoolService) updateNodePoolStats(poolID uint) error {
	var nodePool storage.NodePool
	if err := s.db.First(&nodePool, poolID).Error; err != nil {
		return fmt.Errorf("node pool not found: %w", err)
	}

	return nodePool.UpdateStats(s.db.DB)
}
