package prism

// RuleType 规则类型（数字枚举）
type RuleType int

const (
	// RuleTypeUnknown 未知类型
	RuleTypeUnknown RuleType = 0

	// RuleTypeDomain 域名精确匹配
	// 示例: DOMAIN,google.com,PROXY
	RuleTypeDomain RuleType = 1

	// RuleTypeDomainSuffix 域名后缀匹配
	// 示例: DOMAIN-SUFFIX,google.com,PROXY (匹配 *.google.com)
	RuleTypeDomainSuffix RuleType = 2

	// RuleTypeDomainKeyword 域名关键字匹配
	// 示例: DOMAIN-KEYWORD,google,PROXY
	RuleTypeDomainKeyword RuleType = 3

	// RuleTypeDomainRegex 域名正则表达式匹配
	// 示例: DOMAIN-REGEX,^.*\.google\.com$,PROXY
	RuleTypeDomainRegex RuleType = 4

	// RuleTypeGEOIP GeoIP 国家代码匹配
	// 示例: GEOIP,CN,DIRECT
	RuleTypeGEOIP RuleType = 5

	// RuleTypeIPCIDR IP CIDR 匹配
	// 示例: IP-CIDR,192.168.0.0/16,DIRECT
	RuleTypeIPCIDR RuleType = 6

	// RuleTypeIPCIDR6 IPv6 CIDR 匹配
	// 示例: IP-CIDR6,2001:db8::/32,PROXY
	RuleTypeIPCIDR6 RuleType = 7

	// RuleTypeSrcIP 源 IP 匹配
	// 示例: SRC-IP,192.168.1.1,DIRECT
	RuleTypeSrcIP RuleType = 8

	// RuleTypeSrcIPCIDR 源 IP CIDR 匹配
	// 示例: SRC-IP-CIDR,192.168.0.0/16,DIRECT
	RuleTypeSrcIPCIDR RuleType = 9

	// RuleTypeDstPort 目标端口匹配
	// 示例: DST-PORT,80,DIRECT
	RuleTypeDstPort RuleType = 10

	// RuleTypeSrcPort 源端口匹配
	// 示例: SRC-PORT,7890,DIRECT
	RuleTypeSrcPort RuleType = 11

	// RuleTypeProcess 进程名称匹配
	// 示例: PROCESS-NAME,chrome,PROXY
	RuleTypeProcess RuleType = 12

	// RuleTypeProcessPath 进程路径匹配
	// 示例: PROCESS-PATH,/usr/bin/wget,PROXY
	RuleTypeProcessPath RuleType = 13

	// RuleTypeIPASN IP ASN 匹配
	// 示例: IP-ASN,13335,PROXY
	RuleTypeIPASN RuleType = 14

	// RuleTypeRuleSet 规则集匹配
	// 示例: RULE-SET,reject,REJECT
	RuleTypeRuleSet RuleType = 15

	// RuleTypeInType 入站类型匹配
	// 示例: IN-TYPE,HTTP,PROXY
	RuleTypeInType RuleType = 16

	// RuleTypeInName 入站名称匹配
	// 示例: IN-NAME,socks-in,DIRECT
	RuleTypeInName RuleType = 17

	// RuleTypeInUser 入站用户匹配
	// 示例: IN-USER,admin,PROXY
	RuleTypeInUser RuleType = 18

	// RuleTypeNetwork 网络类型匹配
	// 示例: NETWORK,TCP,DIRECT
	RuleTypeNetwork RuleType = 19

	// RuleTypeUID 用户 ID 匹配
	// 示例: UID,1000,DIRECT
	RuleTypeUID RuleType = 20

	// RuleTypeDSCP DSCP 值匹配
	// 示例: DSCP,46,PROXY
	RuleTypeDSCP RuleType = 21

	// RuleTypeProcessNameRegex 进程名正则匹配
	// 示例: PROCESS-NAME-REGEX,^chrome.*$,PROXY
	RuleTypeProcessNameRegex RuleType = 22

	// RuleTypeProcessPathRegex 进程路径正则匹配
	// 示例: PROCESS-PATH-REGEX,^/usr/bin/.*$,DIRECT
	RuleTypeProcessPathRegex RuleType = 23

	// RuleTypeIPSuffix IP 后缀匹配
	// 示例: IP-SUFFIX,8.8.8.8/24,1,DIRECT
	RuleTypeIPSuffix RuleType = 24

	// RuleTypeGeoSite GeoSite 域名地理位置匹配
	// 示例: GEOSITE,cn,DIRECT
	RuleTypeGeoSite RuleType = 25

	// RuleTypeMatch 匹配所有流量（通常作为最后一条规则）
	// 示例: MATCH,PROXY
	RuleTypeMatch RuleType = 26
)

// String 返回规则类型的字符串表示
func (r RuleType) String() string {
	switch r {
	case RuleTypeDomain:
		return "DOMAIN"
	case RuleTypeDomainSuffix:
		return "DOMAIN-SUFFIX"
	case RuleTypeDomainKeyword:
		return "DOMAIN-KEYWORD"
	case RuleTypeDomainRegex:
		return "DOMAIN-REGEX"
	case RuleTypeGEOIP:
		return "GEOIP"
	case RuleTypeIPCIDR:
		return "IP-CIDR"
	case RuleTypeIPCIDR6:
		return "IP-CIDR6"
	case RuleTypeSrcIP:
		return "SRC-IP"
	case RuleTypeSrcIPCIDR:
		return "SRC-IP-CIDR"
	case RuleTypeDstPort:
		return "DST-PORT"
	case RuleTypeSrcPort:
		return "SRC-PORT"
	case RuleTypeProcess:
		return "PROCESS-NAME"
	case RuleTypeProcessPath:
		return "PROCESS-PATH"
	case RuleTypeIPASN:
		return "IP-ASN"
	case RuleTypeRuleSet:
		return "RULE-SET"
	case RuleTypeInType:
		return "IN-TYPE"
	case RuleTypeInName:
		return "IN-NAME"
	case RuleTypeInUser:
		return "IN-USER"
	case RuleTypeNetwork:
		return "NETWORK"
	case RuleTypeUID:
		return "UID"
	case RuleTypeDSCP:
		return "DSCP"
	case RuleTypeProcessNameRegex:
		return "PROCESS-NAME-REGEX"
	case RuleTypeProcessPathRegex:
		return "PROCESS-PATH-REGEX"
	case RuleTypeIPSuffix:
		return "IP-SUFFIX"
	case RuleTypeGeoSite:
		return "GEOSITE"
	case RuleTypeMatch:
		return "MATCH"
	default:
		return "UNKNOWN"
	}
}

// ParseRuleType 从字符串解析规则类型
func ParseRuleType(s string) RuleType {
	switch s {
	case "DOMAIN":
		return RuleTypeDomain
	case "DOMAIN-SUFFIX":
		return RuleTypeDomainSuffix
	case "DOMAIN-KEYWORD":
		return RuleTypeDomainKeyword
	case "DOMAIN-REGEX":
		return RuleTypeDomainRegex
	case "GEOIP":
		return RuleTypeGEOIP
	case "IP-CIDR":
		return RuleTypeIPCIDR
	case "IP-CIDR6":
		return RuleTypeIPCIDR6
	case "SRC-IP":
		return RuleTypeSrcIP
	case "SRC-IP-CIDR":
		return RuleTypeSrcIPCIDR
	case "DST-PORT":
		return RuleTypeDstPort
	case "SRC-PORT":
		return RuleTypeSrcPort
	case "PROCESS-NAME":
		return RuleTypeProcess
	case "PROCESS-PATH":
		return RuleTypeProcessPath
	case "IP-ASN":
		return RuleTypeIPASN
	case "RULE-SET":
		return RuleTypeRuleSet
	case "IN-TYPE":
		return RuleTypeInType
	case "IN-NAME":
		return RuleTypeInName
	case "IN-USER":
		return RuleTypeInUser
	case "NETWORK":
		return RuleTypeNetwork
	case "UID":
		return RuleTypeUID
	case "DSCP":
		return RuleTypeDSCP
	case "PROCESS-NAME-REGEX":
		return RuleTypeProcessNameRegex
	case "PROCESS-PATH-REGEX":
		return RuleTypeProcessPathRegex
	case "IP-SUFFIX":
		return RuleTypeIPSuffix
	case "GEOSITE":
		return RuleTypeGeoSite
	case "MATCH":
		return RuleTypeMatch
	default:
		return RuleTypeUnknown
	}
}
