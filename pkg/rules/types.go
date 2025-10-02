package rules

import (
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

// Rule 规则接口（引用 prism.Rule）
type Rule = prism.Rule

// Metadata 连接元数据（引用 prism.Metadata）
type Metadata = prism.Metadata

// GeoIPInfo GeoIP 信息（引用 prism.GeoIPInfo）
type GeoIPInfo = prism.GeoIPInfo
