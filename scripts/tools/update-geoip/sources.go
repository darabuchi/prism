package main

import "fmt"

// SourceType 数据源类型
type SourceType string

const (
	SourceTypeMMDB SourceType = "mmdb" // MaxMind MMDB 格式
	SourceTypeCSV  SourceType = "csv"  // CSV 格式
	SourceTypeText SourceType = "txt"  // 文本格式（CIDR 列表）
)

// Source 数据源定义
type Source struct {
	URL         string     // 数据源 URL
	Type        SourceType // 数据源类型
	Format      string     // 数据格式（ipinfo/dbip/cidr）
	Description string     // 描述
	Priority    int        // 优先级（越小越优先）
	Enabled     bool       // 是否启用
	CacheDays   int        // 缓存有效期（天），0 表示使用全局配置
}

// GetSources 获取所有数据源
// 参考 fire 项目的数据源配置
func GetSources(cfg *Config) map[string]Source {
	sources := make(map[string]Source)

	// === IPInfo 数据源（免费，高质量） ===
	// IPInfo 免费数据库每月更新一次
	sources["ipinfo-country-asn"] = Source{
		URL:         fmt.Sprintf("https://ipinfo.io/data/free/country_asn.csv.gz?token=%s", cfg.IPInfoToken),
		Type:        SourceTypeCSV,
		Format:      "ipinfo",
		Description: "IPInfo Country + ASN Database (Free)",
		Priority:    1,
		Enabled:     true,
		CacheDays:   30, // 每月更新
	}

	// === Sapics IP Location DB (基于 DBIP) ===
	// DBIP 数据库每月更新（GitHub 仓库每月同步）
	// DBIP ASN
	sources["sapics-dbip-asn-ipv4"] = Source{
		URL:         "https://raw.githubusercontent.com/sapics/ip-location-db/refs/heads/main/dbip-asn/dbip-asn-ipv4.csv",
		Type:        SourceTypeCSV,
		Format:      "dbip-asn",
		Description: "DBIP ASN Database (IPv4)",
		Priority:    2,
		Enabled:     true,
		CacheDays:   30, // 每月更新
	}
	sources["sapics-dbip-asn-ipv6"] = Source{
		URL:         "https://raw.githubusercontent.com/sapics/ip-location-db/refs/heads/main/dbip-asn/dbip-asn-ipv6.csv",
		Type:        SourceTypeCSV,
		Format:      "dbip-asn",
		Description: "DBIP ASN Database (IPv6)",
		Priority:    2,
		Enabled:     true,
		CacheDays:   30, // 每月更新
	}

	// DBIP City
	sources["sapics-dbip-city-ipv4"] = Source{
		URL:         "https://raw.githubusercontent.com/sapics/ip-location-db/refs/heads/main/dbip-city/dbip-city-ipv4.csv.gz",
		Type:        SourceTypeCSV,
		Format:      "dbip-city",
		Description: "DBIP City Database (IPv4)",
		Priority:    3,
		Enabled:     true,
		CacheDays:   30, // 每月更新
	}
	sources["sapics-dbip-city-ipv6"] = Source{
		URL:         "https://raw.githubusercontent.com/sapics/ip-location-db/refs/heads/main/dbip-city/dbip-city-ipv6.csv.gz",
		Type:        SourceTypeCSV,
		Format:      "dbip-city",
		Description: "DBIP City Database (IPv6)",
		Priority:    3,
		Enabled:     true,
		CacheDays:   30, // 每月更新
	}

	// DBIP Country
	sources["sapics-dbip-country-ipv4"] = Source{
		URL:         "https://raw.githubusercontent.com/sapics/ip-location-db/refs/heads/main/dbip-country/dbip-country-ipv4.csv",
		Type:        SourceTypeCSV,
		Format:      "dbip-country",
		Description: "DBIP Country Database (IPv4)",
		Priority:    4,
		Enabled:     true,
		CacheDays:   30, // 每月更新
	}
	sources["sapics-dbip-country-ipv6"] = Source{
		URL:         "https://raw.githubusercontent.com/sapics/ip-location-db/refs/heads/main/dbip-country/dbip-country-ipv6.csv",
		Type:        SourceTypeCSV,
		Format:      "dbip-country",
		Description: "DBIP Country Database (IPv6)",
		Priority:    4,
		Enabled:     true,
		CacheDays:   30, // 每月更新
	}

	// GeoLite2 ASN - MaxMind 每周二更新 GeoLite2
	sources["sapics-geolite2-asn-ipv4"] = Source{
		URL:         "https://raw.githubusercontent.com/sapics/ip-location-db/refs/heads/main/geolite2-asn/geolite2-asn-ipv4.csv",
		Type:        SourceTypeCSV,
		Format:      "dbip-asn",
		Description: "GeoLite2 ASN Database (IPv4)",
		Priority:    5,
		Enabled:     true,
		CacheDays:   7, // 每周更新
	}
	sources["sapics-geolite2-asn-ipv6"] = Source{
		URL:         "https://raw.githubusercontent.com/sapics/ip-location-db/refs/heads/main/geolite2-asn/geolite2-asn-ipv6.csv",
		Type:        SourceTypeCSV,
		Format:      "dbip-asn",
		Description: "GeoLite2 ASN Database (IPv6)",
		Priority:    5,
		Enabled:     true,
		CacheDays:   7, // 每周更新
	}

	// === 中国 IP 列表（高精度中国 IP 数据） ===
	// 基础中国 IP 列表 - GitHub 仓库不定期更新
	sources["china-ip-chnroute"] = Source{
		URL:         "https://github.com/mayaxcn/china-ip-list/raw/master/chnroute.txt",
		Type:        SourceTypeText,
		Format:      "cidr-cn",
		Description: "China IP List - CHNRoute (IPv4)",
		Priority:    10,
		Enabled:     true,
		CacheDays:   14, // 每两周检查更新
	}
	sources["china-ip-chnroute-v6"] = Source{
		URL:         "https://github.com/mayaxcn/china-ip-list/raw/master/chnroute_v6.txt",
		Type:        SourceTypeText,
		Format:      "cidr-cn",
		Description: "China IP List - CHNRoute (IPv6)",
		Priority:    10,
		Enabled:     true,
		CacheDays:   14, // 每两周检查更新
	}

	// 中国运营商 IP（高质量数据源）- GitHub Pages 每日自动更新
	sources["china-operator-cernet"] = Source{
		URL:         "https://gaoyifan.github.io/china-operator-ip/cernet.txt",
		Type:        SourceTypeText,
		Format:      "cidr-cernet",
		Description: "China Education Network (CERNET IPv4)",
		Priority:    11,
		Enabled:     true,
		CacheDays:   1, // 每日更新
	}
	sources["china-operator-cernet6"] = Source{
		URL:         "https://gaoyifan.github.io/china-operator-ip/cernet6.txt",
		Type:        SourceTypeText,
		Format:      "cidr-cernet",
		Description: "China Education Network (CERNET IPv6)",
		Priority:    11,
		Enabled:     true,
		CacheDays:   1, // 每日更新
	}
	sources["china-operator-chinanet"] = Source{
		URL:         "https://gaoyifan.github.io/china-operator-ip/chinanet.txt",
		Type:        SourceTypeText,
		Format:      "cidr-chinanet",
		Description: "China Telecom (ChinaNet IPv4)",
		Priority:    11,
		Enabled:     true,
		CacheDays:   1, // 每日更新
	}
	sources["china-operator-chinanet6"] = Source{
		URL:         "https://gaoyifan.github.io/china-operator-ip/chinanet6.txt",
		Type:        SourceTypeText,
		Format:      "cidr-chinanet",
		Description: "China Telecom (ChinaNet IPv6)",
		Priority:    11,
		Enabled:     true,
		CacheDays:   1, // 每日更新
	}
	sources["china-operator-cmcc"] = Source{
		URL:         "https://gaoyifan.github.io/china-operator-ip/cmcc.txt",
		Type:        SourceTypeText,
		Format:      "cidr-cmcc",
		Description: "China Mobile (CMCC IPv4)",
		Priority:    11,
		Enabled:     true,
		CacheDays:   1, // 每日更新
	}
	sources["china-operator-cmcc6"] = Source{
		URL:         "https://gaoyifan.github.io/china-operator-ip/cmcc6.txt",
		Type:        SourceTypeText,
		Format:      "cidr-cmcc",
		Description: "China Mobile (CMCC IPv6)",
		Priority:    11,
		Enabled:     true,
		CacheDays:   1, // 每日更新
	}
	sources["china-operator-unicom"] = Source{
		URL:         "https://gaoyifan.github.io/china-operator-ip/unicom.txt",
		Type:        SourceTypeText,
		Format:      "cidr-unicom",
		Description: "China Unicom (IPv4)",
		Priority:    11,
		Enabled:     true,
		CacheDays:   1, // 每日更新
	}
	sources["china-operator-unicom6"] = Source{
		URL:         "https://gaoyifan.github.io/china-operator-ip/unicom6.txt",
		Type:        SourceTypeText,
		Format:      "cidr-unicom",
		Description: "China Unicom (IPv6)",
		Priority:    11,
		Enabled:     true,
		CacheDays:   1, // 每日更新
	}
	sources["china-operator-drpeng"] = Source{
		URL:         "https://gaoyifan.github.io/china-operator-ip/drpeng.txt",
		Type:        SourceTypeText,
		Format:      "cidr-drpeng",
		Description: "Dr.Peng Network (IPv4)",
		Priority:    11,
		Enabled:     true,
		CacheDays:   1, // 每日更新
	}
	sources["china-operator-drpeng6"] = Source{
		URL:         "https://gaoyifan.github.io/china-operator-ip/drpeng6.txt",
		Type:        SourceTypeText,
		Format:      "cidr-drpeng",
		Description: "Dr.Peng Network (IPv6)",
		Priority:    11,
		Enabled:     true,
		CacheDays:   1, // 每日更新
	}

	// 其他中国 IP 列表 - 不同仓库更新频率不同
	sources["china-ip-17mon"] = Source{
		URL:         "https://github.com/17mon/china_ip_list/raw/refs/heads/master/china_ip_list.txt",
		Type:        SourceTypeText,
		Format:      "cidr-cn",
		Description: "17mon China IP List",
		Priority:    12,
		Enabled:     true,
		CacheDays:   30, // 更新较慢，每月检查
	}
	sources["china-ip-chnroutes2"] = Source{
		URL:         "https://github.com/misakaio/chnroutes2/raw/refs/heads/master/chnroutes.txt",
		Type:        SourceTypeText,
		Format:      "cidr-cn",
		Description: "CHNRoutes2 China IP List",
		Priority:    12,
		Enabled:     true,
		CacheDays:   7, // 每周更新
	}
	sources["china-ip-hackl0us"] = Source{
		URL:         "https://github.com/Hackl0us/GeoIP2-CN/raw/release/CN-ip-cidr.txt",
		Type:        SourceTypeText,
		Format:      "cidr-cn",
		Description: "GeoIP2-CN China IP List",
		Priority:    12,
		Enabled:     true,
		CacheDays:   7, // 每周更新
	}

	// 过滤启用的数据源
	enabledSources := make(map[string]Source)
	for name, source := range sources {
		if source.Enabled {
			enabledSources[name] = source
		}
	}

	return enabledSources
}
