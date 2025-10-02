package rules

import (
	"strings"

	"github.com/darabuchi/prism"
)

// GEOIP GeoIP 国家代码匹配规则
type GEOIP struct {
	base
	countryCode string
	noResolve   bool // 是否跳过域名解析
}

// Match 匹配 GeoIP 国家代码
func (g *GEOIP) Match(metadata *Metadata) bool {
	// 如果设置了 no-resolve，且目标是域名，则不匹配
	if g.noResolve && metadata.Host != "" && !metadata.DstIP.IsValid() {
		return false
	}

	// 优先使用已有的 GeoIP 信息
	if metadata.DstGeoIP != nil {
		return strings.EqualFold(metadata.DstGeoIP.CountryCode, g.countryCode)
	}

	// 如果没有 GeoIP 信息，但有 IP 地址，查询 GeoIP 数据库
	if metadata.DstIP.IsValid() {
		info := queryGeoIP(metadata.DstIP.String())
		if info != nil {
			// 缓存 GeoIP 信息
			metadata.DstGeoIP = &GeoIPInfo{
				CountryCode:   info.CountryCode,
				Country:       info.Country,
				Continent:     info.Continent,
				ContinentCode: info.ContinentCode,
				ASN:           info.ASN,
				ASName:        info.ASName,
			}
			return strings.EqualFold(info.CountryCode, g.countryCode)
		}
	}

	return false
}

// NoResolve 设置是否跳过域名解析
func (g *GEOIP) NoResolve(noResolve bool) *GEOIP {
	g.noResolve = noResolve
	return g
}

// NewGEOIP 创建 GeoIP 国家代码匹配规则
func NewGEOIP(payload string, action prism.Payload) *GEOIP {
	countryCode := strings.ToUpper(strings.TrimSpace(payload))
	return &GEOIP{
		base:        newBase(TypeGEOIP, payload, action),
		countryCode: countryCode,
	}
}
