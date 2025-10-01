package dns

import (
	"context"
	"time"
)

// DoH3Resolver DNS over HTTP/3 解析器
type DoH3Resolver struct {
	BaseResolver
	dialer  Dialer
	timeout time.Duration
	doh     *DoHResolver
}

// NewDoH3Resolver 创建 DoH3 DNS 解析器
func NewDoH3Resolver(opts ...Option) *DoH3Resolver {
	c := &Client{
		dialer:  &DefaultDialer{},
		timeout: 10 * time.Second,
	}

	for _, opt := range opts {
		opt(c)
	}

	return &DoH3Resolver{
		dialer:  c.dialer,
		timeout: c.timeout,
		doh:     NewDoHResolver(opts...),
	}
}

// Type 返回解析器类型
func (r *DoH3Resolver) Type() string {
	return "doh3"
}

// Query 执行 DoH3 DNS 查询
//
// 注意：当前实现使用 DoH 作为 fallback
// TODO: 实现完整的 HTTP/3 支持（需要 quic-go 库）
func (r *DoH3Resolver) Query(ctx context.Context, domain, queryType, server string) (*QueryResult, error) {
	// HTTP/3 support is complex and requires additional dependencies
	// For now, fallback to DoH over HTTP/2
	result, err := r.doh.Query(ctx, domain, queryType, server)
	if err != nil {
		return nil, err
	}

	// 修改协议标识为 doh3
	result.Protocol = "doh3"

	return result, nil
}
