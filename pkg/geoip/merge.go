package geoip

import (
	"github.com/darabuchi/prism"
)

// Merge 合并多个 GeoIP 数据
//
// 合并规则：
//   - 优先使用第一个非空值
//   - 对于数字字段，使用第一个非零值
//   - 对于布尔字段，使用 OR 逻辑（任一为 true 则为 true）
//   - IP 和 IPVersion 字段特殊处理，优先使用第一个
func Merge(geoips ...*prism.GeoIP) *prism.GeoIP {
	if len(geoips) == 0 {
		return &prism.GeoIP{}
	}

	result := &prism.GeoIP{}

	for _, g := range geoips {
		if g == nil {
			continue
		}

		// IP 信息 - 使用第一个非空值
		if result.IP == "" && g.IP != "" {
			result.IP = g.IP
		}
		if result.IPVersion == 0 && g.IPVersion != 0 {
			result.IPVersion = g.IPVersion
		}

		// 地理位置信息 - 使用第一个非空值
		if result.Country == "" && g.Country != "" {
			result.Country = g.Country
		}
		if result.CountryCode == "" && g.CountryCode != "" {
			result.CountryCode = g.CountryCode
		}
		if result.Region == "" && g.Region != "" {
			result.Region = g.Region
		}
		if result.RegionCode == "" && g.RegionCode != "" {
			result.RegionCode = g.RegionCode
		}
		if result.City == "" && g.City != "" {
			result.City = g.City
		}
		if result.Latitude == 0 && g.Latitude != 0 {
			result.Latitude = g.Latitude
		}
		if result.Longitude == 0 && g.Longitude != 0 {
			result.Longitude = g.Longitude
		}
		if result.Postal == "" && g.Postal != "" {
			result.Postal = g.Postal
		}
		if result.Timezone == "" && g.Timezone != "" {
			result.Timezone = g.Timezone
		}

		// ASN 信息 - 使用第一个非零值
		if result.ASN == 0 && g.ASN != 0 {
			result.ASN = g.ASN
		}
		if result.ASName == "" && g.ASName != "" {
			result.ASName = g.ASName
		}
		if result.AS == "" && g.AS != "" {
			result.AS = g.AS
		}

		// ISP 信息 - 使用第一个非空值
		if result.ISP == "" && g.ISP != "" {
			result.ISP = g.ISP
		}
		if result.Org == "" && g.Org != "" {
			result.Org = g.Org
		}

		// 其他信息
		if result.Continent == "" && g.Continent != "" {
			result.Continent = g.Continent
		}
		if result.ContinentCode == "" && g.ContinentCode != "" {
			result.ContinentCode = g.ContinentCode
		}

		// 布尔值使用 OR 逻辑
		result.Proxy = result.Proxy || g.Proxy
		result.Hosting = result.Hosting || g.Hosting
	}

	// 填充派生字段
	result.FillDerivedFields()

	return result
}

// Append 将 source 的非空字段追加到 dest
//
// 与 Merge 不同，Append 只修改已存在的 dest，不创建新对象
// 只有当 dest 中的字段为空时，才会从 source 复制
func Append(dest *prism.GeoIP, source *prism.GeoIP) {
	if dest == nil || source == nil {
		return
	}

	// IP 信息
	if dest.IP == "" && source.IP != "" {
		dest.IP = source.IP
	}
	if dest.IPVersion == 0 && source.IPVersion != 0 {
		dest.IPVersion = source.IPVersion
	}

	// 地理位置信息
	if dest.Country == "" && source.Country != "" {
		dest.Country = source.Country
	}
	if dest.CountryCode == "" && source.CountryCode != "" {
		dest.CountryCode = source.CountryCode
	}
	if dest.Region == "" && source.Region != "" {
		dest.Region = source.Region
	}
	if dest.RegionCode == "" && source.RegionCode != "" {
		dest.RegionCode = source.RegionCode
	}
	if dest.City == "" && source.City != "" {
		dest.City = source.City
	}
	if dest.Latitude == 0 && source.Latitude != 0 {
		dest.Latitude = source.Latitude
	}
	if dest.Longitude == 0 && source.Longitude != 0 {
		dest.Longitude = source.Longitude
	}
	if dest.Postal == "" && source.Postal != "" {
		dest.Postal = source.Postal
	}
	if dest.Timezone == "" && source.Timezone != "" {
		dest.Timezone = source.Timezone
	}

	// ASN 信息
	if dest.ASN == 0 && source.ASN != 0 {
		dest.ASN = source.ASN
	}
	if dest.ASName == "" && source.ASName != "" {
		dest.ASName = source.ASName
	}
	if dest.AS == "" && source.AS != "" {
		dest.AS = source.AS
	}

	// ISP 信息
	if dest.ISP == "" && source.ISP != "" {
		dest.ISP = source.ISP
	}
	if dest.Org == "" && source.Org != "" {
		dest.Org = source.Org
	}

	// 其他信息
	if dest.Continent == "" && source.Continent != "" {
		dest.Continent = source.Continent
	}
	if dest.ContinentCode == "" && source.ContinentCode != "" {
		dest.ContinentCode = source.ContinentCode
	}

	// 布尔值使用 OR 逻辑
	dest.Proxy = dest.Proxy || source.Proxy
	dest.Hosting = dest.Hosting || source.Hosting

	// 填充派生字段
	dest.FillDerivedFields()
}

// Clone 克隆一个 GeoIP 对象
func Clone(g *prism.GeoIP) *prism.GeoIP {
	if g == nil {
		return &prism.GeoIP{}
	}

	return &prism.GeoIP{
		IP:            g.IP,
		IPVersion:     g.IPVersion,
		Country:       g.Country,
		CountryCode:   g.CountryCode,
		Region:        g.Region,
		RegionCode:    g.RegionCode,
		City:          g.City,
		Latitude:      g.Latitude,
		Longitude:     g.Longitude,
		Postal:        g.Postal,
		Timezone:      g.Timezone,
		ASN:           g.ASN,
		ASName:        g.ASName,
		AS:            g.AS,
		ISP:           g.ISP,
		Org:           g.Org,
		Continent:     g.Continent,
		ContinentCode: g.ContinentCode,
		Proxy:         g.Proxy,
		Hosting:       g.Hosting,
	}
}
