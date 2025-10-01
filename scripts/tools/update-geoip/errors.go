package main

import "github.com/darabuchi/prism"

// GeoIP 相关错误码定义
// 错误码范围：3300-3399 (resource.geoip)
// 所有错误码已在 error_code.json 中注册
// 这些常量是从 error_code.gen.go 中引用的别名，方便本地使用
const (
	// ErrInvalidIP 无效的 IP 地址
	ErrInvalidIP = prism.ErrGeoipInvalidIp

	// ErrDatabaseOpenFailed 打开 GeoIP 数据库失败
	ErrDatabaseOpenFailed = prism.ErrGeoipDatabaseOpenFailed

	// ErrDatabaseWriteFailed 写入 GeoIP 数据库失败
	ErrDatabaseWriteFailed = prism.ErrGeoipDatabaseWriteFailed

	// ErrConversionFailed GeoIP 数据转换失败
	ErrConversionFailed = prism.ErrGeoipConversionFailed

	// ErrFileOperationFailed GeoIP 文件操作失败
	ErrFileOperationFailed = prism.ErrGeoipFileOperationFailed

	// ErrWriterCreateFailed 创建 GeoIP 写入器失败
	ErrWriterCreateFailed = prism.ErrGeoipWriterCreateFailed

	// ErrInsertFailed 插入 GeoIP 数据失败
	ErrInsertFailed = prism.ErrGeoipInsertFailed
)
