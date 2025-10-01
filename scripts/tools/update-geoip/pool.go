package main

import (
	"strings"
	"sync"

	"github.com/darabuchi/prism"
)

// GeoIPPool GeoIP 对象池，减少内存分配
var geoIPPool = sync.Pool{
	New: func() interface{} {
		return &prism.GeoIP{}
	},
}

// GetGeoIP 从对象池获取 GeoIP 对象
func GetGeoIP() *prism.GeoIP {
	geo := geoIPPool.Get().(*prism.GeoIP)
	// 重置对象
	*geo = prism.GeoIP{}
	return geo
}

// PutGeoIP 归还 GeoIP 对象到对象池
func PutGeoIP(geo *prism.GeoIP) {
	if geo != nil {
		geoIPPool.Put(geo)
	}
}

// StringBuilderPool 字符串构建器池
var stringBuilderPool = sync.Pool{
	New: func() interface{} {
		return &strings.Builder{}
	},
}

// GetStringBuilder 从对象池获取 StringBuilder
func GetStringBuilder() *strings.Builder {
	sb := stringBuilderPool.Get().(*strings.Builder)
	sb.Reset()
	return sb
}

// PutStringBuilder 归还 StringBuilder 到对象池
func PutStringBuilder(sb *strings.Builder) {
	if sb != nil {
		stringBuilderPool.Put(sb)
	}
}
