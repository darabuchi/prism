package rules

import (
	"net/netip"

	"github.com/darabuchi/prism"
)

// RuleType 规则类型（引用 prism.RuleType）
type RuleType = prism.RuleType

// 规则类型常量别名（兼容旧代码）
const (
	// TypeDomain 域名精确匹配
	TypeDomain = prism.RuleTypeDomain

	// TypeDomainSuffix 域名后缀匹配
	TypeDomainSuffix = prism.RuleTypeDomainSuffix

	// TypeDomainKeyword 域名关键字匹配
	TypeDomainKeyword = prism.RuleTypeDomainKeyword

	// TypeDomainRegex 域名正则表达式匹配
	TypeDomainRegex = prism.RuleTypeDomainRegex

	// TypeGEOIP GeoIP 国家代码匹配
	TypeGEOIP = prism.RuleTypeGEOIP

	// TypeIPCIDR IP CIDR 匹配
	TypeIPCIDR = prism.RuleTypeIPCIDR

	// TypeIPCIDR6 IPv6 CIDR 匹配
	TypeIPCIDR6 = prism.RuleTypeIPCIDR6

	// TypeSrcIP 源 IP 匹配
	TypeSrcIP = prism.RuleTypeSrcIP

	// TypeSrcIPCIDR 源 IP CIDR 匹配
	TypeSrcIPCIDR = prism.RuleTypeSrcIPCIDR

	// TypeDstPort 目标端口匹配
	TypeDstPort = prism.RuleTypeDstPort

	// TypeSrcPort 源端口匹配
	TypeSrcPort = prism.RuleTypeSrcPort

	// TypeProcess 进程名称匹配
	TypeProcess = prism.RuleTypeProcess

	// TypeProcessPath 进程路径匹配
	TypeProcessPath = prism.RuleTypeProcessPath

	// TypeIPASN IP ASN 匹配
	TypeIPASN = prism.RuleTypeIPASN

	// TypeRuleSet 规则集匹配
	TypeRuleSet = prism.RuleTypeRuleSet

	// TypeInType 入站类型匹配
	TypeInType = prism.RuleTypeInType

	// TypeInName 入站名称匹配
	TypeInName = prism.RuleTypeInName

	// TypeInUser 入站用户匹配
	TypeInUser = prism.RuleTypeInUser

	// TypeNetwork 网络类型匹配
	TypeNetwork = prism.RuleTypeNetwork

	// TypeUID 用户 ID 匹配
	TypeUID = prism.RuleTypeUID

	// TypeDSCP DSCP 值匹配
	TypeDSCP = prism.RuleTypeDSCP

	// TypeProcessNameRegex 进程名正则匹配
	TypeProcessNameRegex = prism.RuleTypeProcessNameRegex

	// TypeProcessPathRegex 进程路径正则匹配
	TypeProcessPathRegex = prism.RuleTypeProcessPathRegex

	// TypeIPSuffix IP 后缀匹配
	TypeIPSuffix = prism.RuleTypeIPSuffix

	// TypeGeoSite GeoSite 域名地理位置匹配
	TypeGeoSite = prism.RuleTypeGeoSite

	// TypeMatch 匹配所有流量（通常作为最后一条规则）
	TypeMatch = prism.RuleTypeMatch
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

	// 入站信息
	InboundType string      // 入站类型: HTTP, SOCKS5, etc.
	InboundName string      // 入站名称
	InboundUser string      // 入站用户（认证用户名）

	// 系统信息
	UID         uint32      // 用户 ID (Linux/Android)
	DSCP        uint8       // DSCP 值 (Differentiated Services Code Point)
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
	Action() prism.Payload

	// Payload 返回规则载荷（匹配内容）
	Payload() string

	// String 返回规则的字符串表示
	String() string
}
