package node

import (
	"maps"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/candy"
	"github.com/lazygophers/utils/cryptox"
)

// GenerateId 生成节点唯一ID
//
// 该函数基于节点的配置生成唯一标识符，与 fire 项目的 CalculateClashHash 逻辑保持一致。
// 生成过程：
//  1. 如果配置中有 unique_id 字段，直接返回该字段值
//  2. 排除不影响节点唯一性的字段（name、哈希值、额外信息等）
//  3. 将配置转换为 URL 格式：scheme=type, host=server:port, query=其他字段
//  4. 对 URL 字符串计算 SHA256 哈希
//
// 参数:
//   - config: 节点配置，必须包含 type、server 等基本字段
//
// 返回:
//   - string: 64 位十六进制 SHA256 哈希字符串，或 unique_id 字段值
//
// 示例:
//
//	config := map[string]any{
//	    "type":   "vmess",
//	    "server": "example.com",
//	    "port":   443,
//	    "uuid":   "xxx-xxx-xxx",
//	}
//	id := GenerateId(config)
func GenerateId(config map[string]any) string {
	// 克隆配置，避免修改原始数据
	m := maps.Clone(config)

	// 如果配置中已经有 unique_id，直接使用
	if v, ok := m["unique_id"]; ok {
		return candy.ToString(v)
	}

	// 删除不影响节点唯一性的字段
	// 这些字段只是元数据，不应该影响节点的身份识别
	delete(m, "name")        // 节点名称（用户自定义，可变）
	delete(m, "md5")         // MD5 哈希值（已弃用）
	delete(m, "sha256")      // SHA256 哈希值（避免循环）
	delete(m, "sha384")      // SHA384 哈希值（避免循环）
	delete(m, "sha512")      // SHA512 哈希值（避免循环）
	delete(m, "unique_key")  // 旧版唯一键（已弃用）
	delete(m, "unique_id")   // 唯一ID（已处理）
	delete(m, "extra_info")  // 额外信息（元数据）

	// deepField 递归处理各种类型的字段值，转换为字符串
	var deepField func(v any) string
	deepField = func(v any) string {
		if v == nil {
			return ""
		}

		switch x := v.(type) {
		case string:
			return x

		case []byte:
			return string(x)

		case int:
			return strconv.Itoa(x)

		case int64:
			return strconv.FormatInt(x, 10)

		case []string:
			// 过滤空字符串，排序去重后用分号连接
			x = candy.Filter(x, func(s string) bool {
				return s != ""
			})
			return strings.Join(candy.Sort(candy.Unique(x)), ";")

		case float64:
			return strconv.FormatFloat(x, 'f', -1, 64)

		case bool:
			if x {
				return "true"
			}
			return "false"

		case map[string]string:
			// 将 map 转换为 URL 编码格式
			sub := &url.Values{}
			for kk, vv := range x {
				if vv == "" {
					continue
				}
				sub.Set(kk, vv)
			}
			return sub.Encode()

		case []interface{}:
			// 递归处理数组中的每个元素
			var ss []string
			for _, x := range x {
				ss = append(ss, deepField(x))
			}
			ss = candy.Filter(ss, func(s string) bool {
				return s != ""
			})
			return strings.Join(candy.Sort(candy.Unique(ss)), ";")

		case map[string]any:
			// 递归处理嵌套的 map
			sub := &url.Values{}
			for k, v := range x {
				sub.Set(k, deepField(v))
			}
			return sub.Encode()

		default:
			// 不支持的类型，记录日志并返回空字符串
			log.Warnf("unsupported type %T for field value", x)
			return ""
		}
	}

	// 构建 URL 对象
	u := &url.URL{}

	// 设置 scheme 为节点类型
	if typ, ok := m["type"]; ok {
		u.Scheme = candy.ToString(typ)
		delete(m, "type")
	}

	// 设置 host 为 server:port
	if host, ok := m["server"]; ok {
		u.Host = candy.ToString(host)
		delete(m, "server")
	}

	if port, ok := m["port"]; ok {
		u.Host = net.JoinHostPort(u.Host, candy.ToString(port))
		delete(m, "port")
	} else {
		// 如果没有 port 字段，使用默认端口 8388（Shadowsocks 默认端口）
		u.Host = net.JoinHostPort(u.Host, "8388")
	}

	// 将剩余的字段转换为 URL query 参数
	query := &url.Values{}
	for k, v := range m {
		query.Set(k, deepField(v))
	}

	u.RawQuery = query.Encode()

	// 对完整的 URL 字符串计算 SHA256 哈希
	return cryptox.Sha256(u.String())
}
