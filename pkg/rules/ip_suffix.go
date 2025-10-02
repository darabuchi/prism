package rules

import (
	"fmt"
	"net/netip"
	"strconv"
	"strings"
)

// IPSuffix IP 后缀匹配规则
//
// 匹配 IP 地址的特定后缀位
// 格式: IP-SUFFIX,8.8.8.8/24,1,DIRECT
// 表示匹配 8.8.8.0/24 段中后缀为 .1 的 IP (即 8.8.8.1)
type IPSuffix struct {
	base
	prefix     netip.Prefix
	suffix     uint32
	suffixBits int
	noResolve  bool
}

// Match 匹配 IP 后缀
func (i *IPSuffix) Match(metadata *Metadata) bool {
	// 如果设置了 no-resolve，且目标是域名，则不匹配
	if i.noResolve && metadata.Host != "" && !metadata.DstIP.IsValid() {
		return false
	}

	if !metadata.DstIP.IsValid() {
		return false
	}

	// 检查 IP 是否在前缀范围内
	if !i.prefix.Contains(metadata.DstIP) {
		return false
	}

	// 提取 IP 的后缀部分并匹配
	ip := metadata.DstIP
	if ip.Is4() {
		ipBytes := ip.As4()
		// 计算实际的 IP 值
		ipVal := uint32(ipBytes[0])<<24 | uint32(ipBytes[1])<<16 | uint32(ipBytes[2])<<8 | uint32(ipBytes[3])

		// 提取后缀部分
		mask := uint32((1 << i.suffixBits) - 1)
		actualSuffix := ipVal & mask

		return actualSuffix == i.suffix
	}

	// IPv6 暂不支持
	return false
}

// NoResolve 设置是否跳过域名解析
func (i *IPSuffix) NoResolve(noResolve bool) *IPSuffix {
	i.noResolve = noResolve
	return i
}

// NewIPSuffix 创建 IP 后缀匹配规则
//
// payload 格式: "IP/PREFIX,SUFFIX"
// 示例: "8.8.8.8/24,1" 匹配 8.8.8.1
func NewIPSuffix(payload string, action ActionType) (*IPSuffix, error) {
	parts := strings.Split(payload, ",")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid IP-SUFFIX format, expected IP/PREFIX,SUFFIX")
	}

	// 解析前缀
	prefix, err := netip.ParsePrefix(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, err
	}

	// 解析后缀值
	suffix, err := strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 32)
	if err != nil {
		return nil, err
	}

	// 计算后缀位数
	prefixBits := prefix.Bits()
	suffixBits := 32 - prefixBits // IPv4

	// 验证后缀值不超过范围
	if suffix >= (1 << suffixBits) {
		return nil, fmt.Errorf("suffix %d exceeds maximum for /%d prefix", suffix, prefixBits)
	}

	return &IPSuffix{
		base:       newBase(TypeIPSuffix, payload, action),
		prefix:     prefix,
		suffix:     uint32(suffix),
		suffixBits: suffixBits,
	}, nil
}
