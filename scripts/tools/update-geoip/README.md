# GeoIP 数据库更新工具

从可信的 GitHub 数据源下载并更新 GeoIP 数据库，支持缓存和代理配置。

## 功能特性

- ✅ 从多个可信 GitHub 数据源下载 GeoIP 数据库
- ✅ HTTP 缓存机制（默认 7 天）
- ✅ 代理支持
- ✅ ISP/Org 名称标准化
- ✅ 自动生成元数据文件

## 使用方法

### 通过 Makefile

```bash
# 更新 GeoIP 数据库
make update-geoip

# 清理缓存并强制重新下载
make clean-geoip-cache
make update-geoip
```

### 直接运行

```bash
# 查看帮助
go run ./scripts/tools/update-geoip --help

# 更新数据库（详细输出）
go run ./scripts/tools/update-geoip -v

# 强制更新（忽略缓存）
go run ./scripts/tools/update-geoip -f
```

## 配置

工具优先从项目根目录的 `.env` 文件读取配置，如果文件不存在则使用默认配置。

### 默认配置

```bash
GEOIP_CACHE_DIR=/tmp/prism_geoip     # 缓存目录
GEOIP_DATA_DIR=./resource/geoip      # 数据目录
GEOIP_CACHE_DAYS=7                   # 缓存有效期（天）
GEOIP_PROXY=                         # HTTP 代理地址
```

### 自定义配置

在项目根目录创建或编辑 `.env` 文件：

```bash
# GeoIP 配置
GEOIP_CACHE_DIR=/tmp/my_geoip_cache
GEOIP_DATA_DIR=./data/geoip
GEOIP_CACHE_DAYS=14
GEOIP_PROXY=http://127.0.0.1:7890
```

## 数据源

| 数据库 | 大小 | 描述 |
|--------|------|------|
| GeoLite2-City | ~59 MB | 城市级别地理位置数据 |
| GeoLite2-Country | ~9.4 MB | 国家级别地理位置数据 |
| GeoLite2-ASN | ~10 MB | ASN 信息数据 |
| Country (Loyalsoldier) | ~10 MB | CN 优化版地理位置数据 |

## 输出文件

更新完成后，会在数据目录生成以下文件：

- `*.mmdb` - MaxMind 数据库文件
- `isp_mappings.txt` - ISP/Org 名称标准化映射表
- `metadata.json` - 元数据文件（包含更新时间、数据源等信息）

## ISP 名称标准化

工具会自动标准化常见的 ISP/Org 名称，例如：

- `Google LLC`, `Google Inc.`, `GOOGLE` → `Google`
- `Amazon.com, Inc.`, `Amazon Technologies Inc.` → `Amazon`
- `China Telecom`, `CHINANET` → `中国电信`
- `Alibaba Cloud LLC`, `ALICLOUD` → `Alibaba Cloud`

完整映射表见 `resource/geoip/isp_mappings.txt`。

## 命令行参数

```
Flags:
  -f, --force     强制更新，忽略缓存
  -h, --help      help for update-geoip
  -v, --verbose   详细输出
```

## 示例

### 基本使用

```bash
# 首次运行或缓存过期时会下载数据
$ make update-geoip
Updating GeoIP databases...
开始更新 GeoIP 数据库...
处理数据源: GeoLite2-City
从 https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-City.mmdb 下载...
下载完成: 59.07 MB
...
GeoIP 数据库更新完成!
```

### 强制更新

```bash
# 清理缓存
$ make clean-geoip-cache
Cleaning GeoIP cache...
GeoIP cache cleaned

# 重新下载
$ make update-geoip
```

### 使用代理

在 `.env` 文件中配置：

```bash
GEOIP_PROXY=http://127.0.0.1:7890
```

## 在代码中使用

更新完成后，可以在 Go 代码中使用下载的数据库：

```go
import "github.com/darabuchi/prism/pkg/geoip"

// 打开数据库
reader, err := geoip.NewMaxMindReader("resource/geoip/GeoLite2-City.mmdb")
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
```

## 注意事项

- 首次下载需要较长时间（约 1-2 分钟，取决于网络速度）
- 使用代理可以提高国内下载速度
- 缓存默认 7 天有效，过期后自动重新下载
- 可以使用 `-f` 参数强制更新，忽略缓存
