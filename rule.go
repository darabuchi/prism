package prism

import "net/netip"

// Rule 规则接口
//
// 所有规则类型都需要实现此接口
type Rule interface {
	// Type 返回规则类型
	Type() RuleType

	// Match 判断元数据是否匹配规则
	Match(metadata *Metadata) bool

	// Action 返回规则动作（目标：DIRECT、PROXY 等）
	Action() Payload

	// Content 返回规则内容（匹配模式：域名、IP、端口等）
	Content() string

	// String 返回规则的字符串表示
	String() string
}

// Metadata 连接元数据
//
// 包含连接的各种信息，用于规则匹配
type Metadata struct {
	// 网络信息
	Network string     // 网络类型: tcp, udp
	Type    string     // 连接类型: http, https, socks5
	SrcIP   netip.Addr // 源 IP
	DstIP   netip.Addr // 目标 IP
	SrcPort uint16     // 源端口
	DstPort uint16     // 目标端口

	// 域名信息
	Host   string // 主机名/域名
	Domain string // 规范化的域名（小写）

	// GeoIP 信息
	DstGeoIP *GeoIPInfo // 目标 IP 的 GeoIP 信息
	SrcGeoIP *GeoIPInfo // 源 IP 的 GeoIP 信息

	// 进程信息
	ProcessName string // 进程名称
	ProcessPath string // 进程路径

	// 入站信息
	InboundType string // 入站类型: HTTP, SOCKS5, etc.
	InboundName string // 入站名称
	InboundUser string // 入站用户（认证用户名）

	// 系统信息
	UID  uint32 // 用户 ID (Linux/Android)
	DSCP uint8  // DSCP 值 (Differentiated Services Code Point)
}

// GeoIPInfo GeoIP 信息
type GeoIPInfo struct {
	CountryCode   string // 国家代码
	Country       string // 国家名称
	Continent     string // 大洲
	ContinentCode string // 大洲代码
	ASN           int    // 自治系统号
	ASName        string // 自治系统名称
}
