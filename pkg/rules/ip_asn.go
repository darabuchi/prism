package rules

import (
	"strconv"

	"github.com/darabuchi/prism"
)

// IPASN IP ASN 匹配规则
type IPASN struct {
	base
	asn       int
	noResolve bool // 是否跳过域名解析
}

// Match 匹配 IP ASN
func (i *IPASN) Match(metadata *Metadata) bool {
	// 如果设置了 no-resolve，且目标是域名，则不匹配
	if i.noResolve && metadata.Host != "" && !metadata.DstIP.IsValid() {
		return false
	}

	// 优先使用已有的 GeoIP 信息
	if metadata.DstGeoIP != nil {
		return metadata.DstGeoIP.ASN == i.asn
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
			return info.ASN == i.asn
		}
	}

	return false
}

// NoResolve 设置是否跳过域名解析
func (i *IPASN) NoResolve(noResolve bool) *IPASN {
	i.noResolve = noResolve
	return i
}

// NewIPASN 创建 IP ASN 匹配规则
func NewIPASN(payload string, action prism.Payload) (*IPASN, error) {
	asn, err := strconv.Atoi(payload)
	if err != nil {
		return nil, err
	}

	return &IPASN{
		base: newBase(TypeIPASN, payload, action),
		asn:  asn,
	}, nil
}
