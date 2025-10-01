package rules

import (
	"sync"

	prism "github.com/darabuchi/prism"
)

// GeoIPProvider GeoIP 查询提供者
//
// 用于查询 IP 地址的 GeoIP 信息
type GeoIPProvider func(ip string) *prism.GeoIP

var (
	// 全局 GeoIP 提供者
	geoipProvider GeoIPProvider
	providerMu    sync.RWMutex
)

// SetGeoIPProvider 设置全局 GeoIP 提供者
//
// 必须在使用 GEOIP 或 IP-ASN 规则前调用此函数
func SetGeoIPProvider(provider GeoIPProvider) {
	providerMu.Lock()
	defer providerMu.Unlock()
	geoipProvider = provider
}

// GetGeoIPProvider 获取全局 GeoIP 提供者
func GetGeoIPProvider() GeoIPProvider {
	providerMu.RLock()
	defer providerMu.RUnlock()
	return geoipProvider
}

// queryGeoIP 查询 GeoIP 信息
func queryGeoIP(ip string) *prism.GeoIP {
	provider := GetGeoIPProvider()
	if provider == nil {
		return nil
	}
	return provider(ip)
}
