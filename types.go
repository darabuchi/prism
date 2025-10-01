package prism

import (
	"fmt"
	"net"
	"strings"
)

// ProxyType 代理协议类型
//
// 定义所有支持的代理协议类型常量
// 这些常量用于标识节点的协议类型，与 Mihomo/Meta 核心的协议类型保持一致
type ProxyType string

const (
	// TypeShadowsocks Shadowsocks 代理协议
	TypeShadowsocks ProxyType = "ss"

	// TypeShadowsocksR ShadowsocksR 代理协议
	TypeShadowsocksR ProxyType = "ssr"

	// TypeVMess V2Ray VMess 协议
	TypeVMess ProxyType = "vmess"

	// TypeVLess V2Ray VLess 协议
	TypeVLess ProxyType = "vless"

	// TypeTrojan Trojan 代理协议
	TypeTrojan ProxyType = "trojan"

	// TypeHysteria Hysteria 协议
	TypeHysteria ProxyType = "hysteria"

	// TypeHysteria2 Hysteria2 协议
	TypeHysteria2 ProxyType = "hysteria2"

	// TypeSocks5 SOCKS5 代理协议
	TypeSocks5 ProxyType = "socks5"

	// TypeHTTP HTTP(S) 代理协议
	TypeHTTP ProxyType = "http"

	// TypeSnell Snell 协议
	TypeSnell ProxyType = "snell"

	// TypeWireGuard WireGuard VPN 协议
	TypeWireGuard ProxyType = "wireguard"

	// TypeTuic TUIC 协议
	TypeTuic ProxyType = "tuic"

	// TypeSSH SSH 隧道协议
	TypeSSH ProxyType = "ssh"

	// TypeMieru Mieru 协议
	TypeMieru ProxyType = "mieru"

	// TypeAnyTLS AnyTLS 协议
	TypeAnyTLS ProxyType = "anytls"

	// TypeDirect 直连（不使用代理）
	TypeDirect ProxyType = "direct"

	// TypeReject 拒绝连接
	TypeReject ProxyType = "reject"

	// TypeDNS DNS 查询
	TypeDNS ProxyType = "dns"
)

// HealthState 熔断器健康状态
//
// 用于表示节点的健康状态，基于熔断器模式
// 熔断器会根据节点的成功/失败次数自动切换状态
type HealthState int32

const (
	// HealthStateOpen 开路状态（熔断器开启）
	// 节点连续失败次数超过阈值，熔断器开启，拒绝所有请求
	HealthStateOpen HealthState = 0

	// HealthStateHalfOpen 半开状态（熔断器尝试恢复）
	// 熔断器开启一段时间后，进入半开状态，允许少量请求通过以测试节点是否恢复
	HealthStateHalfOpen HealthState = 1

	// HealthStateClosed 闭路状态（熔断器关闭）
	// 节点工作正常，熔断器关闭，允许所有请求通过
	HealthStateClosed HealthState = 2
)

// String 返回健康状态的字符串表示
func (h HealthState) String() string {
	switch h {
	case HealthStateOpen:
		return "open"
	case HealthStateHalfOpen:
		return "half_open"
	case HealthStateClosed:
		return "closed"
	default:
		return "unknown"
	}
}

// GeoIP 地理位置信息
//
// 包含 IP 地址的地理位置、ASN、ISP 等信息
// 通常从 GeoIP 数据库或在线服务获取
type GeoIP struct {
	// IP 信息
	IP        string `json:"ip,omitempty" maxminddb:"-"`        // IP 地址
	IPVersion int    `json:"ip_version,omitempty" maxminddb:"-"` // IP 版本 (4 或 6)

	// 地理位置信息
	Country     string  `json:"country,omitempty" maxminddb:"country>names>en"`           // 国家名称
	CountryCode string  `json:"country_code,omitempty" maxminddb:"country>iso_code"`      // 国家代码 (ISO 3166-1 alpha-2)
	Region      string  `json:"region,omitempty" maxminddb:"subdivisions>0>names>en"`     // 地区/省份名称
	RegionCode  string  `json:"region_code,omitempty" maxminddb:"subdivisions>0>iso_code"` // 地区/省份代码
	City        string  `json:"city,omitempty" maxminddb:"city>names>en"`                 // 城市名称
	Latitude    float64 `json:"latitude,omitempty" maxminddb:"location>latitude"`         // 纬度
	Longitude   float64 `json:"longitude,omitempty" maxminddb:"location>longitude"`       // 经度
	Postal      string  `json:"postal,omitempty" maxminddb:"postal>code"`                 // 邮政编码
	Timezone    string  `json:"timezone,omitempty" maxminddb:"location>time_zone"`        // 时区

	// ASN 信息
	ASN    int    `json:"asn,omitempty" maxminddb:"traits>autonomous_system_number"`         // 自治系统号
	ASName string `json:"as_name,omitempty" maxminddb:"traits>autonomous_system_organization"` // 自治系统名称
	AS     string `json:"as,omitempty" maxminddb:"-"`                                        // AS 字符串表示 (如 "AS15169")

	// ISP 信息
	ISP string `json:"isp,omitempty" maxminddb:"traits>isp"`          // 互联网服务提供商
	Org string `json:"org,omitempty" maxminddb:"traits>organization"` // 组织名称

	// 其他信息
	Continent     string `json:"continent,omitempty" maxminddb:"continent>names>en"`       // 大洲
	ContinentCode string `json:"continent_code,omitempty" maxminddb:"continent>code"`      // 大洲代码
	Proxy         bool   `json:"proxy,omitempty" maxminddb:"traits>is_anonymous_proxy"`    // 是否为代理/VPN
	Hosting       bool   `json:"hosting,omitempty" maxminddb:"traits>is_hosting_provider"` // 是否为托管服务器
}

// Location 返回格式化的地理位置字符串
//
// 格式: "城市, 地区, 国家" 或根据可用信息调整
func (g *GeoIP) Location() string {
	var parts []string

	if g.City != "" {
		parts = append(parts, g.City)
	}
	if g.Region != "" {
		parts = append(parts, g.Region)
	}
	if g.Country != "" {
		parts = append(parts, g.Country)
	} else if g.CountryCode != "" {
		parts = append(parts, g.CountryCode)
	}

	if len(parts) == 0 {
		return "Unknown"
	}

	return strings.Join(parts, ", ")
}

// Coordinates 返回格式化的坐标字符串
func (g *GeoIP) Coordinates() string {
	if g.Latitude == 0 && g.Longitude == 0 {
		return ""
	}
	return fmt.Sprintf("%.4f, %.4f", g.Latitude, g.Longitude)
}

// ASString 返回完整的 AS 信息字符串
//
// 格式: "AS15169 Google LLC" 或 "AS15169" (如果没有名称)
func (g *GeoIP) ASString() string {
	if g.ASN == 0 {
		return ""
	}

	if g.ASName != "" {
		return fmt.Sprintf("AS%d %s", g.ASN, g.ASName)
	}

	return fmt.Sprintf("AS%d", g.ASN)
}

// IsIPv4 检查是否为 IPv4 地址
func (g *GeoIP) IsIPv4() bool {
	if g.IPVersion != 0 {
		return g.IPVersion == 4
	}

	// 如果未设置 IPVersion，从 IP 字符串判断
	if g.IP == "" {
		return false
	}

	ip := net.ParseIP(g.IP)
	if ip == nil {
		return false
	}

	return ip.To4() != nil
}

// IsIPv6 检查是否为 IPv6 地址
func (g *GeoIP) IsIPv6() bool {
	if g.IPVersion != 0 {
		return g.IPVersion == 6
	}

	// 如果未设置 IPVersion，从 IP 字符串判断
	if g.IP == "" {
		return false
	}

	ip := net.ParseIP(g.IP)
	if ip == nil {
		return false
	}

	return ip.To4() == nil && ip.To16() != nil
}

// HasLocation 检查是否有地理位置信息
func (g *GeoIP) HasLocation() bool {
	return g.Country != "" || g.CountryCode != "" || g.City != "" || g.Region != ""
}

// HasASN 检查是否有 ASN 信息
func (g *GeoIP) HasASN() bool {
	return g.ASN > 0
}

// String 返回 GeoIP 的字符串表示
func (g *GeoIP) String() string {
	var parts []string

	if g.IP != "" {
		parts = append(parts, g.IP)
	}

	location := g.Location()
	if location != "Unknown" {
		parts = append(parts, location)
	}

	if g.ISP != "" {
		parts = append(parts, g.ISP)
	}

	asStr := g.ASString()
	if asStr != "" {
		parts = append(parts, asStr)
	}

	if len(parts) == 0 {
		return "Empty GeoIP"
	}

	return strings.Join(parts, " | ")
}

// FillDerivedFields 填充派生字段
//
// 在从 MaxMind DB 或其他来源读取数据后调用此方法
// 用于填充那些不直接从数据库读取但可以派生的字段
func (g *GeoIP) FillDerivedFields() {
	// 从 ASN 生成 AS 字符串
	if g.ASN > 0 && g.AS == "" {
		if g.ASName != "" {
			g.AS = fmt.Sprintf("AS%d %s", g.ASN, g.ASName)
		} else {
			g.AS = fmt.Sprintf("AS%d", g.ASN)
		}
	}

	// 从 IP 字符串推断 IP 版本
	if g.IP != "" && g.IPVersion == 0 {
		ip := net.ParseIP(g.IP)
		if ip != nil {
			if ip.To4() != nil {
				g.IPVersion = 4
			} else if ip.To16() != nil {
				g.IPVersion = 6
			}
		}
	}
}
