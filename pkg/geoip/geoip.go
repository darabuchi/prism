package geoip

import (
	"fmt"
	"net"

	"github.com/darabuchi/prism"
)

// Reader GeoIP 数据读取器接口
type Reader interface {
	// Lookup 查询 IP 的 GeoIP 信息
	Lookup(ip net.IP) (*prism.GeoIP, error)

	// LookupString 查询 IP 字符串的 GeoIP 信息
	LookupString(ipStr string) (*prism.GeoIP, error)

	// Close 关闭读取器
	Close() error
}

// LookupIP 从多个 Reader 中查询 IP 的 GeoIP 信息并合并
//
// 按照 readers 的顺序依次查询，后面的结果会填充前面缺失的字段
func LookupIP(ip net.IP, readers ...Reader) (*prism.GeoIP, error) {
	if len(readers) == 0 {
		return nil, fmt.Errorf("no readers provided")
	}

	var results []*prism.GeoIP

	for _, reader := range readers {
		if reader == nil {
			continue
		}

		geoip, err := reader.Lookup(ip)
		if err != nil {
			continue
		}

		if geoip != nil {
			results = append(results, geoip)
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no GeoIP data found for IP: %s", ip.String())
	}

	// 合并所有结果
	return Merge(results...), nil
}

// LookupString 从多个 Reader 中查询 IP 字符串的 GeoIP 信息并合并
func LookupString(ipStr string, readers ...Reader) (*prism.GeoIP, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ipStr)
	}

	return LookupIP(ip, readers...)
}

// MustLookupIP 查询 IP 的 GeoIP 信息，失败时 panic
func MustLookupIP(ip net.IP, readers ...Reader) *prism.GeoIP {
	result, err := LookupIP(ip, readers...)
	if err != nil {
		panic(err)
	}
	return result
}

// MustLookupString 查询 IP 字符串的 GeoIP 信息，失败时 panic
func MustLookupString(ipStr string, readers ...Reader) *prism.GeoIP {
	result, err := LookupString(ipStr, readers...)
	if err != nil {
		panic(err)
	}
	return result
}
