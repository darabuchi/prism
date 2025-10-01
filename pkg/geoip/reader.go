package geoip

import (
	"fmt"
	"net"

	"github.com/darabuchi/prism"
	"github.com/oschwald/maxminddb-golang"
)

// MaxMindReader MaxMind GeoIP2 数据库读取器
type MaxMindReader struct {
	db *maxminddb.Reader
}

// NewMaxMindReader 创建 MaxMind 数据库读取器
//
// dbPath 为 MaxMind GeoIP2 数据库文件路径
// 支持 GeoLite2-City.mmdb, GeoIP2-City.mmdb 等
func NewMaxMindReader(dbPath string) (*MaxMindReader, error) {
	db, err := maxminddb.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open MaxMind database: %w", err)
	}

	return &MaxMindReader{
		db: db,
	}, nil
}

// Lookup 查询 IP 的 GeoIP 信息
func (r *MaxMindReader) Lookup(ip net.IP) (*prism.GeoIP, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var geoip prism.GeoIP

	err := r.db.Lookup(ip, &geoip)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup IP: %w", err)
	}

	// 设置 IP 地址
	geoip.IP = ip.String()

	// 填充派生字段
	geoip.FillDerivedFields()

	return &geoip, nil
}

// LookupString 查询 IP 字符串的 GeoIP 信息
func (r *MaxMindReader) LookupString(ipStr string) (*prism.GeoIP, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ipStr)
	}

	return r.Lookup(ip)
}

// Close 关闭数据库
func (r *MaxMindReader) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

// Metadata 返回数据库元数据
func (r *MaxMindReader) Metadata() maxminddb.Metadata {
	if r.db != nil {
		return r.db.Metadata
	}
	return maxminddb.Metadata{}
}
