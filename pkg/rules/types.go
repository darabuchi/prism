package rules

import "net/netip"

// RuleType 规则类型
type RuleType string

const (
	// TypeDomain 域名精确匹配
	// 示例: DOMAIN,google.com,PROXY
	TypeDomain RuleType = "DOMAIN"

	// TypeDomainSuffix 域名后缀匹配
	// 示例: DOMAIN-SUFFIX,google.com,PROXY (匹配 *.google.com)
	TypeDomainSuffix RuleType = "DOMAIN-SUFFIX"

	// TypeDomainKeyword 域名关键字匹配
	// 示例: DOMAIN-KEYWORD,google,PROXY
	TypeDomainKeyword RuleType = "DOMAIN-KEYWORD"

	// TypeDomainRegex 域名正则表达式匹配
	// 示例: DOMAIN-REGEX,^.*\.google\.com$,PROXY
	TypeDomainRegex RuleType = "DOMAIN-REGEX"

	// TypeGEOIP GeoIP 国家代码匹配
	// 示例: GEOIP,CN,DIRECT
	TypeGEOIP RuleType = "GEOIP"

	// TypeIPCIDR IP CIDR 匹配
	// 示例: IP-CIDR,192.168.0.0/16,DIRECT
	TypeIPCIDR RuleType = "IP-CIDR"

	// TypeIPCIDR6 IPv6 CIDR 匹配
	// 示例: IP-CIDR6,2001:db8::/32,PROXY
	TypeIPCIDR6 RuleType = "IP-CIDR6"

	// TypeSrcIP 源 IP 匹配
	// 示例: SRC-IP,192.168.1.1,DIRECT
	TypeSrcIP RuleType = "SRC-IP"

	// TypeSrcIPCIDR 源 IP CIDR 匹配
	// 示例: SRC-IP-CIDR,192.168.0.0/16,DIRECT
	TypeSrcIPCIDR RuleType = "SRC-IP-CIDR"

	// TypeDstPort 目标端口匹配
	// 示例: DST-PORT,80,DIRECT
	TypeDstPort RuleType = "DST-PORT"

	// TypeSrcPort 源端口匹配
	// 示例: SRC-PORT,7890,DIRECT
	TypeSrcPort RuleType = "SRC-PORT"

	// TypeProcess 进程名称匹配
	// 示例: PROCESS-NAME,chrome,PROXY
	TypeProcess RuleType = "PROCESS-NAME"

	// TypeProcessPath 进程路径匹配
	// 示例: PROCESS-PATH,/usr/bin/wget,PROXY
	TypeProcessPath RuleType = "PROCESS-PATH"

	// TypeIPASN IP ASN 匹配
	// 示例: IP-ASN,13335,PROXY
	TypeIPASN RuleType = "IP-ASN"

	// TypeRuleSet 规则集匹配
	// 示例: RULE-SET,reject,REJECT
	TypeRuleSet RuleType = "RULE-SET"

	// TypeMatch 匹配所有流量（通常作为最后一条规则）
	// 示例: MATCH,PROXY
	TypeMatch RuleType = "MATCH"
)

// ActionType 规则动作类型
type ActionType string

const (
	// ActionProxy 使用代理
	ActionProxy ActionType = "PROXY"

	// ActionDirect 直连
	ActionDirect ActionType = "DIRECT"

	// ActionReject 拒绝连接
	ActionReject ActionType = "REJECT"

	// ActionRejectDrop 拒绝连接并丢弃数据包
	ActionRejectDrop ActionType = "REJECT-DROP"
)

// Metadata 连接元数据
//
// 包含连接的各种信息，用于规则匹配
type Metadata struct {
	// 网络信息
	Network     string      // 网络类型: tcp, udp
	Type        string      // 连接类型: http, https, socks5
	SrcIP       netip.Addr  // 源 IP
	DstIP       netip.Addr  // 目标 IP
	SrcPort     uint16      // 源端口
	DstPort     uint16      // 目标端口

	// 域名信息
	Host        string      // 主机名/域名
	Domain      string      // 规范化的域名（小写）

	// GeoIP 信息
	DstGeoIP    *GeoIPInfo  // 目标 IP 的 GeoIP 信息
	SrcGeoIP    *GeoIPInfo  // 源 IP 的 GeoIP 信息

	// 进程信息
	ProcessName string      // 进程名称
	ProcessPath string      // 进程路径

	// 其他信息
	UID         uint32      // 用户 ID
	InboundName string      // 入站名称
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

// Rule 规则接口
//
// 所有规则类型都需要实现此接口
type Rule interface {
	// Type 返回规则类型
	Type() RuleType

	// Match 判断元数据是否匹配规则
	Match(metadata *Metadata) bool

	// Action 返回规则动作
	Action() ActionType

	// Payload 返回规则载荷（匹配内容）
	Payload() string

	// String 返回规则的字符串表示
	String() string
}
