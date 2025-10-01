package dns

import (
	"context"
	"time"

	"github.com/metacubex/mihomo/component/resolver"
)

// Resolver DNS 解析器接口
type Resolver interface {
	resolver.Resolver

	// Query 执行 DNS 查询
	//
	// 参数:
	//   - ctx: 上下文
	//   - domain: 查询的域名
	//   - queryType: 查询类型（A, AAAA, CNAME, MX, TXT 等）
	//   - server: DNS 服务器地址
	//
	// 返回:
	//   - *QueryResult: 查询结果
	//   - error: 错误信息
	Query(ctx context.Context, domain, queryType, server string) (*QueryResult, error)

	// Type 返回解析器类型
	Type() string
}

// QueryResult DNS 查询结果
type QueryResult struct {
	Question   string        // 查询的域名
	Type       string        // 查询类型
	Answers    []string      // 应答记录
	Duration   time.Duration // 查询耗时
	Server     string        // DNS 服务器地址
	Protocol   string        // 协议类型
	StatusCode int           // 状态码
}

// NewResolver 创建解析器
//
// 参数:
//   - protocol: 协议类型（udp, tcp, doh, dot, doh3）
//   - opts: 选项
//
// 返回:
//   - Resolver: 解析器实例
func NewResolver(protocol string, opts ...Option) Resolver {
	switch protocol {
	case "udp":
		return NewUDPResolver(opts...)
	case "tcp":
		return NewTCPResolver(opts...)
	case "doh", "https":
		return NewDoHResolver(opts...)
	case "dot", "tls":
		return NewDoTResolver(opts...)
	case "doh3", "h3":
		return NewDoH3Resolver(opts...)
	default:
		return NewUDPResolver(opts...)
	}
}
