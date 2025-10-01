# GeoIP Package

GeoIP 包提供了 IP 地理位置信息的读取、合并和处理功能。

## 功能特性

- 🔍 **多数据源支持**: 支持从多个数据源读取 GeoIP 信息
- 🔀 **智能合并**: 自动合并多个数据源的结果，填充缺失字段
- 📚 **MaxMind DB 集成**: 内置 MaxMind GeoIP2 数据库支持
- 🛠️ **实用工具**: 提供 Clone、Append 等实用函数

## 安装

```bash
go get github.com/darabuchi/prism/pkg/geoip
```

## 核心接口

### Reader 接口

```go
type Reader interface {
    Lookup(ip net.IP) (*prism.GeoIP, error)
    LookupString(ipStr string) (*prism.GeoIP, error)
    Close() error
}
```

## 使用示例

### 使用 MaxMind 数据库

```go
package main

import (
    "fmt"
    "github.com/darabuchi/prism/pkg/geoip"
)

func main() {
    // 打开 MaxMind 数据库
    reader, err := geoip.NewMaxMindReader("GeoLite2-City.mmdb")
    if err != nil {
        panic(err)
    }
    defer reader.Close()

    // 查询 IP
    result, err := reader.LookupString("8.8.8.8")
    if err != nil {
        panic(err)
    }

    fmt.Println(result.String())
    // 输出: 8.8.8.8 | Mountain View, California, United States | Google LLC | AS15169 Google LLC
}
```

### 合并多个数据源

```go
// 创建多个读取器
reader1, _ := geoip.NewMaxMindReader("GeoLite2-City.mmdb")
reader2, _ := geoip.NewMaxMindReader("GeoLite2-ASN.mmdb")
defer reader1.Close()
defer reader2.Close()

// 自动合并结果
result, err := geoip.LookupString("8.8.8.8", reader1, reader2)
if err != nil {
    panic(err)
}

// result 包含两个数据库的合并信息
fmt.Printf("Country: %s\n", result.Country)
fmt.Printf("ASN: %d\n", result.ASN)
```

### 手动合并 GeoIP 数据

```go
geoip1 := &prism.GeoIP{
    IP:      "8.8.8.8",
    Country: "United States",
    City:    "Mountain View",
}

geoip2 := &prism.GeoIP{
    CountryCode: "US",
    Region:      "California",
    ASN:         15169,
}

// 合并多个 GeoIP 对象
merged := geoip.Merge(geoip1, geoip2)

fmt.Println(merged.String())
// 输出: 8.8.8.8 | Mountain View, California, United States | AS15169
```

### 追加缺失字段

```go
base := &prism.GeoIP{
    IP:      "8.8.8.8",
    Country: "United States",
}

additional := &prism.GeoIP{
    City: "Mountain View",
    ASN:  15169,
}

// 只追加缺失的字段，不覆盖已有字段
geoip.Append(base, additional)

fmt.Printf("City: %s\n", base.City)    // Mountain View
fmt.Printf("Country: %s\n", base.Country) // United States
```

### 克隆 GeoIP 对象

```go
original := &prism.GeoIP{
    IP:      "8.8.8.8",
    Country: "United States",
    ASN:     15169,
}

// 克隆对象
cloned := geoip.Clone(original)

// 修改克隆对象不影响原对象
cloned.City = "Los Angeles"
```

## 合并规则

`Merge` 和 `Append` 函数使用以下规则合并 GeoIP 数据：

1. **字符串字段**: 使用第一个非空值
2. **数字字段**: 使用第一个非零值
3. **布尔字段**: 使用 OR 逻辑（任一为 true 则为 true）
4. **派生字段**: 自动调用 `FillDerivedFields()` 填充 AS 和 IPVersion

### 示例

```go
geoip1 := &prism.GeoIP{Country: "United States", City: "Mountain View"}
geoip2 := &prism.GeoIP{Country: "US", Region: "California"} // Country 会被忽略

merged := geoip.Merge(geoip1, geoip2)
// 结果: Country = "United States", City = "Mountain View", Region = "California"
```

## 自定义 Reader

实现 `Reader` 接口即可创建自定义数据源：

```go
type CustomReader struct {
    // your fields
}

func (r *CustomReader) Lookup(ip net.IP) (*prism.GeoIP, error) {
    // 实现查询逻辑
    return &prism.GeoIP{
        IP:      ip.String(),
        Country: "Custom Country",
    }, nil
}

func (r *CustomReader) LookupString(ipStr string) (*prism.GeoIP, error) {
    ip := net.ParseIP(ipStr)
    if ip == nil {
        return nil, fmt.Errorf("invalid IP")
    }
    return r.Lookup(ip)
}

func (r *CustomReader) Close() error {
    return nil
}
```

## API 文档

### 函数

- `LookupIP(ip net.IP, readers ...Reader) (*prism.GeoIP, error)` - 从多个 Reader 查询并合并
- `LookupString(ipStr string, readers ...Reader) (*prism.GeoIP, error)` - 查询 IP 字符串
- `Merge(geoips ...*prism.GeoIP) *prism.GeoIP` - 合并多个 GeoIP 对象
- `Append(dest *prism.GeoIP, source *prism.GeoIP)` - 追加缺失字段到 dest
- `Clone(g *prism.GeoIP) *prism.GeoIP` - 克隆 GeoIP 对象

### MaxMindReader

- `NewMaxMindReader(dbPath string) (*MaxMindReader, error)` - 创建 MaxMind 读取器
- `Lookup(ip net.IP) (*prism.GeoIP, error)` - 查询 IP
- `LookupString(ipStr string) (*prism.GeoIP, error)` - 查询 IP 字符串
- `Close() error` - 关闭数据库
- `Metadata() maxminddb.Metadata` - 获取数据库元数据

## 最佳实践

1. **使用 defer Close()**: 确保数据库文件正确关闭
2. **错误处理**: 始终检查错误返回值
3. **多数据源**: 结合多个数据库获得更完整的信息
4. **缓存结果**: 对于频繁查询的 IP，考虑缓存结果

## 许可证

本包是 Prism 项目的一部分。
