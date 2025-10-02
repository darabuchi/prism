package rules

import (
	"strconv"
	"strings"

	"github.com/darabuchi/prism"
)

// Network 网络类型匹配规则
type Network struct {
	base
	network string
}

// Match 匹配网络类型
func (n *Network) Match(metadata *Metadata) bool {
	if metadata.Network == "" {
		return false
	}
	return strings.EqualFold(metadata.Network, n.network)
}

// NewNetwork 创建网络类型匹配规则
func NewNetwork(payload string, action prism.Payload) *Network {
	network := strings.ToUpper(strings.TrimSpace(payload))
	return &Network{
		base:    newBase(TypeNetwork, payload, action),
		network: network,
	}
}

// UID 用户 ID 匹配规则
type UID struct {
	base
	uid uint32
}

// Match 匹配用户 ID
func (u *UID) Match(metadata *Metadata) bool {
	return metadata.UID == u.uid
}

// NewUID 创建用户 ID 匹配规则
func NewUID(payload string, action prism.Payload) (*UID, error) {
	uid, err := strconv.ParseUint(payload, 10, 32)
	if err != nil {
		return nil, err
	}

	return &UID{
		base: newBase(TypeUID, payload, action),
		uid:  uint32(uid),
	}, nil
}

// DSCP DSCP 值匹配规则
type DSCP struct {
	base
	dscp uint8
}

// Match 匹配 DSCP 值
func (d *DSCP) Match(metadata *Metadata) bool {
	return metadata.DSCP == d.dscp
}

// NewDSCP 创建 DSCP 值匹配规则
func NewDSCP(payload string, action prism.Payload) (*DSCP, error) {
	dscp, err := strconv.ParseUint(payload, 10, 8)
	if err != nil {
		return nil, err
	}

	return &DSCP{
		base: newBase(TypeDSCP, payload, action),
		dscp: uint8(dscp),
	}, nil
}
