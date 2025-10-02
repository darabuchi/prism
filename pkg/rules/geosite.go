package rules

import (
	"strings"
	"sync"

	"github.com/darabuchi/prism"
)

// GeositeProvider GeoSite 数据提供者
//
// 用于查询域名所属的 GeoSite 类别
type GeositeProvider func(domain string) []string

var (
	// 全局 GeoSite 提供者
	geositeProvider GeositeProvider
	geositeMu       sync.RWMutex
)

// SetGeositeProvider 设置全局 GeoSite 提供者
//
// 必须在使用 GEOSITE 规则前调用此函数
// provider 应返回域名所属的所有类别列表，如 ["cn", "google", "ads"]
func SetGeositeProvider(provider GeositeProvider) {
	geositeMu.Lock()
	defer geositeMu.Unlock()
	geositeProvider = provider
}

// GetGeositeProvider 获取全局 GeoSite 提供者
func GetGeositeProvider() GeositeProvider {
	geositeMu.RLock()
	defer geositeMu.RUnlock()
	return geositeProvider
}

// queryGeosite 查询域名的 GeoSite 类别
func queryGeosite(domain string) []string {
	provider := GetGeositeProvider()
	if provider == nil {
		return nil
	}
	return provider(domain)
}

// GeoSite GeoSite 域名地理位置匹配规则
//
// 基于 V2Ray 的 geosite.dat 数据库
type GeoSite struct {
	base
	category string
}

// Match 匹配 GeoSite 类别
func (g *GeoSite) Match(metadata *Metadata) bool {
	// 获取域名
	domain := metadata.Domain
	if domain == "" {
		if metadata.Host != "" {
			domain = strings.ToLower(metadata.Host)
		} else {
			return false
		}
	}

	// 查询 GeoSite 类别
	categories := queryGeosite(domain)
	if categories == nil {
		return false
	}

	// 检查是否包含目标类别
	for _, cat := range categories {
		if strings.EqualFold(cat, g.category) {
			return true
		}
	}

	return false
}

// NewGeoSite 创建 GeoSite 匹配规则
func NewGeoSite(payload string, action prism.Payload) *GeoSite {
	category := strings.ToLower(strings.TrimSpace(payload))
	return &GeoSite{
		base:     newBase(TypeGeoSite, payload, action),
		category: category,
	}
}

// DomainMatcher 域名匹配器
//
// 用于高效匹配域名到 GeoSite 类别
// 这是一个简单的实现，可以根据需要优化
type DomainMatcher struct {
	// 精确域名匹配
	domains map[string][]string

	// 域名后缀匹配
	suffixes map[string][]string

	// 域名关键字匹配
	keywords map[string][]string

	mu sync.RWMutex
}

// NewDomainMatcher 创建域名匹配器
func NewDomainMatcher() *DomainMatcher {
	return &DomainMatcher{
		domains:  make(map[string][]string),
		suffixes: make(map[string][]string),
		keywords: make(map[string][]string),
	}
}

// AddDomain 添加精确域名匹配
func (m *DomainMatcher) AddDomain(domain string, categories ...string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	domain = strings.ToLower(domain)
	m.domains[domain] = append(m.domains[domain], categories...)
}

// AddSuffix 添加域名后缀匹配
func (m *DomainMatcher) AddSuffix(suffix string, categories ...string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	suffix = strings.ToLower(suffix)
	suffix = strings.TrimPrefix(suffix, ".")
	m.suffixes[suffix] = append(m.suffixes[suffix], categories...)
}

// AddKeyword 添加域名关键字匹配
func (m *DomainMatcher) AddKeyword(keyword string, categories ...string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	keyword = strings.ToLower(keyword)
	m.keywords[keyword] = append(m.keywords[keyword], categories...)
}

// Match 匹配域名，返回所有匹配的类别
func (m *DomainMatcher) Match(domain string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	domain = strings.ToLower(domain)
	var categories []string

	// 精确匹配
	if cats, ok := m.domains[domain]; ok {
		categories = append(categories, cats...)
	}

	// 后缀匹配
	for suffix, cats := range m.suffixes {
		if domain == suffix || strings.HasSuffix(domain, "."+suffix) {
			categories = append(categories, cats...)
		}
	}

	// 关键字匹配
	for keyword, cats := range m.keywords {
		if strings.Contains(domain, keyword) {
			categories = append(categories, cats...)
		}
	}

	// 去重
	return uniqueStrings(categories)
}

// uniqueStrings 去重字符串切片
func uniqueStrings(slice []string) []string {
	if len(slice) == 0 {
		return slice
	}

	seen := make(map[string]struct{}, len(slice))
	result := make([]string, 0, len(slice))

	for _, s := range slice {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			result = append(result, s)
		}
	}

	return result
}
