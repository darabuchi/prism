package geoip

import (
	"fmt"
	"io"
	"net"

	"github.com/darabuchi/prism"
	"github.com/lazygophers/log"
	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/inserter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
)

// Writer MMDB 数据库写入器
// 用于创建和写入自定义的 Prism GeoIP 数据库
type Writer struct {
	tree *mmdbwriter.Tree
}

// NewWriter 创建新的 MMDB 写入器
func NewWriter() (*Writer, error) {
	tree, err := mmdbwriter.New(
		mmdbwriter.Options{
			DatabaseType:            "Prism-GeoIP",
			Description:             map[string]string{"en": "Prism GeoIP Database"},
			DisableIPv4Aliasing:     true,
			IncludeReservedNetworks: true,
			RecordSize:              32,
			DisableMetadataPointers: true,
			Inserter:                createInserter(),
		},
	)
	if err != nil {
		log.Errorf("创建 MMDB 写入器失败: %v", err)
		return nil, fmt.Errorf("create mmdb writer failed: %w", err)
	}

	return &Writer{tree: tree}, nil
}

// InsertGeoIP 插入单个 IP 的 GeoIP 数据
func (w *Writer) InsertGeoIP(ip string, geo *prism.GeoIP) error {
	if geo == nil {
		return nil
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		log.Warnf("无效的 IP 地址: %s", ip)
		return fmt.Errorf("invalid IP address: %s", ip)
	}

	// 创建 IPNet
	var ipnet *net.IPNet
	if parsedIP.To4() != nil {
		_, ipnet, _ = net.ParseCIDR(ip + "/32")
	} else {
		_, ipnet, _ = net.ParseCIDR(ip + "/128")
	}

	// 转换为 MMDB 类型
	data := geoIPToMMDBType(geo)

	// 插入数据
	if err := w.tree.Insert(ipnet, data); err != nil {
		log.Errorf("插入 IP 数据失败，IP: %s, Country: %s, 错误: %v",
			ip, geo.CountryCode, err)
		return fmt.Errorf("insert IP %s failed: %w", ip, err)
	}

	return nil
}

// InsertGeoIPRange 插入 IP 范围的 GeoIP 数据
func (w *Writer) InsertGeoIPRange(startIP, endIP string, geo *prism.GeoIP) error {
	if geo == nil {
		return nil
	}

	start := net.ParseIP(startIP)
	end := net.ParseIP(endIP)

	if start == nil || end == nil {
		log.Warnf("无效的 IP 范围: %s - %s", startIP, endIP)
		return fmt.Errorf("invalid IP range: %s - %s", startIP, endIP)
	}

	// 转换为 MMDB 类型
	data := geoIPToMMDBType(geo)

	// 插入范围
	if err := w.tree.InsertRange(start, end, data); err != nil {
		log.Errorf("插入 IP 范围失败，范围: %s-%s, Country: %s, ASN: %d, 错误: %v",
			startIP, endIP, geo.CountryCode, geo.ASN, err)
		return fmt.Errorf("insert IP range %s-%s failed: %w", startIP, endIP, err)
	}

	return nil
}

// InsertGeoIPNet 插入 IP 网络的 GeoIP 数据
func (w *Writer) InsertGeoIPNet(ipnet *net.IPNet, geo *prism.GeoIP) error {
	if geo == nil {
		return nil
	}

	// 转换为 MMDB 类型
	data := geoIPToMMDBType(geo)

	// 插入数据
	if err := w.tree.Insert(ipnet, data); err != nil {
		log.Errorf("插入 IP 网络失败，网络: %s, Country: %s, 错误: %v",
			ipnet.String(), geo.CountryCode, err)
		return fmt.Errorf("insert IP network %s failed: %w", ipnet.String(), err)
	}

	return nil
}

// WriteTo 将数据库写入到 io.Writer
func (w *Writer) WriteTo(writer io.Writer) (int64, error) {
	written, err := w.tree.WriteTo(writer)
	if err != nil {
		log.Errorf("写入 MMDB 数据库失败: %v", err)
		return 0, fmt.Errorf("write mmdb database failed: %w", err)
	}

	log.Infof("成功写入 MMDB 数据库，大小: %.2f MB", float64(written)/(1024*1024))
	return written, nil
}

// createInserter 创建自定义的插入器，用于合并重叠 IP 的数据
func createInserter() inserter.FuncGenerator {
	return func(old mmdbtype.DataType) inserter.Func {
		return func(new mmdbtype.DataType) (mmdbtype.DataType, error) {
			// 如果旧数据为空，直接返回新数据
			if old == nil {
				return new, nil
			}

			// 将 MMDB 类型转换为 GeoIP
			oldGeo := mmdbTypeToGeoIP(old)
			newGeo := mmdbTypeToGeoIP(new)

			// 使用 pkg/geoip 的 Merge 函数合并
			merged := Merge(oldGeo, newGeo)

			// 转换回 MMDB 类型
			return geoIPToMMDBType(merged), nil
		}
	}
}

// geoIPToMMDBType 将 prism.GeoIP 转换为 MMDB 类型
func geoIPToMMDBType(geo *prism.GeoIP) mmdbtype.DataType {
	if geo == nil {
		return mmdbtype.Map{}
	}

	data := mmdbtype.Map{}

	// IP 信息
	if geo.IP != "" {
		data[mmdbtype.String("ip")] = mmdbtype.String(geo.IP)
	}
	if geo.IPVersion != 0 {
		data[mmdbtype.String("ip_version")] = mmdbtype.Uint64(uint64(geo.IPVersion))
	}

	// 地理位置信息
	if geo.Country != "" {
		data[mmdbtype.String("country")] = mmdbtype.String(geo.Country)
	}
	if geo.CountryCode != "" {
		data[mmdbtype.String("cc")] = mmdbtype.String(geo.CountryCode)
	}
	if geo.Region != "" {
		data[mmdbtype.String("region")] = mmdbtype.String(geo.Region)
	}
	if geo.RegionCode != "" {
		data[mmdbtype.String("rc")] = mmdbtype.String(geo.RegionCode)
	}
	if geo.City != "" {
		data[mmdbtype.String("city")] = mmdbtype.String(geo.City)
	}
	if geo.Latitude != 0 {
		data[mmdbtype.String("lat")] = mmdbtype.Float64(geo.Latitude)
	}
	if geo.Longitude != 0 {
		data[mmdbtype.String("lon")] = mmdbtype.Float64(geo.Longitude)
	}
	if geo.Postal != "" {
		data[mmdbtype.String("postal")] = mmdbtype.String(geo.Postal)
	}
	if geo.Timezone != "" {
		data[mmdbtype.String("tz")] = mmdbtype.String(geo.Timezone)
	}

	// ASN 信息
	if geo.ASN != 0 {
		data[mmdbtype.String("asn")] = mmdbtype.Uint64(uint64(geo.ASN))
	}
	if geo.ASName != "" {
		data[mmdbtype.String("as_name")] = mmdbtype.String(geo.ASName)
	}
	if geo.AS != "" {
		data[mmdbtype.String("as")] = mmdbtype.String(geo.AS)
	}

	// ISP 信息
	if geo.ISP != "" {
		data[mmdbtype.String("isp")] = mmdbtype.String(geo.ISP)
	}
	if geo.Org != "" {
		data[mmdbtype.String("org")] = mmdbtype.String(geo.Org)
	}

	// 其他信息
	if geo.Continent != "" {
		data[mmdbtype.String("continent")] = mmdbtype.String(geo.Continent)
	}
	if geo.ContinentCode != "" {
		data[mmdbtype.String("ccode")] = mmdbtype.String(geo.ContinentCode)
	}

	// 布尔值
	data[mmdbtype.String("proxy")] = mmdbtype.Bool(geo.Proxy)
	data[mmdbtype.String("hosting")] = mmdbtype.Bool(geo.Hosting)

	return data
}

// mmdbTypeToGeoIP 将 MMDB 类型转换为 prism.GeoIP
func mmdbTypeToGeoIP(data mmdbtype.DataType) *prism.GeoIP {
	if data == nil {
		return &prism.GeoIP{}
	}

	m, ok := data.(mmdbtype.Map)
	if !ok {
		return &prism.GeoIP{}
	}

	geo := &prism.GeoIP{}

	// 辅助函数：安全获取字符串
	getString := func(key string) string {
		if val, ok := m[mmdbtype.String(key)]; ok {
			if str, ok := val.(mmdbtype.String); ok {
				return string(str)
			}
		}
		return ""
	}

	// 辅助函数：安全获取 uint64
	getUint64 := func(key string) uint64 {
		if val, ok := m[mmdbtype.String(key)]; ok {
			if u, ok := val.(mmdbtype.Uint64); ok {
				return uint64(u)
			}
		}
		return 0
	}

	// 辅助函数：安全获取 float64
	getFloat64 := func(key string) float64 {
		if val, ok := m[mmdbtype.String(key)]; ok {
			if f, ok := val.(mmdbtype.Float64); ok {
				return float64(f)
			}
		}
		return 0
	}

	// 辅助函数：安全获取 bool
	getBool := func(key string) bool {
		if val, ok := m[mmdbtype.String(key)]; ok {
			if b, ok := val.(mmdbtype.Bool); ok {
				return bool(b)
			}
		}
		return false
	}

	// 填充 GeoIP 结构
	geo.IP = getString("ip")
	geo.IPVersion = int(getUint64("ip_version"))
	geo.Country = getString("country")
	geo.CountryCode = getString("cc")
	geo.Region = getString("region")
	geo.RegionCode = getString("rc")
	geo.City = getString("city")
	geo.Latitude = getFloat64("lat")
	geo.Longitude = getFloat64("lon")
	geo.Postal = getString("postal")
	geo.Timezone = getString("tz")
	geo.ASN = int(getUint64("asn"))
	geo.ASName = getString("as_name")
	geo.AS = getString("as")
	geo.ISP = getString("isp")
	geo.Org = getString("org")
	geo.Continent = getString("continent")
	geo.ContinentCode = getString("ccode")
	geo.Proxy = getBool("proxy")
	geo.Hosting = getBool("hosting")

	return geo
}
