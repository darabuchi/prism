package rules

import (
	"net/netip"
)

// IPCIDR IP CIDR 匹配规则
type IPCIDR struct {
	base
	prefix  netip.Prefix
	isIPv6  bool
	noResolve bool // 是否跳过域名解析
}

// Match 匹配 IP CIDR
func (i *IPCIDR) Match(metadata *Metadata) bool {
	// 如果设置了 no-resolve，且目标是域名，则不匹配
	if i.noResolve && metadata.Host != "" && !metadata.DstIP.IsValid() {
		return false
	}

	if !metadata.DstIP.IsValid() {
		return false
	}

	// IPv4 规则不匹配 IPv6 地址，反之亦然
	if i.isIPv6 != metadata.DstIP.Is6() {
		return false
	}

	return i.prefix.Contains(metadata.DstIP)
}

// NoResolve 设置是否跳过域名解析
func (i *IPCIDR) NoResolve(noResolve bool) *IPCIDR {
	i.noResolve = noResolve
	return i
}

// NewIPCIDR 创建 IP CIDR 匹配规则
func NewIPCIDR(payload string, action ActionType, isIPv6 bool) (*IPCIDR, error) {
	prefix, err := netip.ParsePrefix(payload)
	if err != nil {
		return nil, err
	}

	ruleType := TypeIPCIDR
	if isIPv6 {
		ruleType = TypeIPCIDR6
	}

	return &IPCIDR{
		base:   newBase(ruleType, payload, action),
		prefix: prefix,
		isIPv6: isIPv6,
	}, nil
}
