package dns

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"time"

	mdns "github.com/miekg/dns"
)

// Dialer 网络连接创建接口
type Dialer interface {
	// Dial 创建网络连接
	Dial(network, address string) (net.Conn, error)

	// DialContext 使用上下文创建网络连接
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

// DefaultDialer 默认拨号器
type DefaultDialer struct{}

// Dial 实现 Dialer 接口
func (d *DefaultDialer) Dial(network, address string) (net.Conn, error) {
	return net.Dial(network, address)
}

// DialContext 实现 Dialer 接口
func (d *DefaultDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	dialer := &net.Dialer{}
	return dialer.DialContext(ctx, network, address)
}

// Client DNS 客户端（用于选项应用）
type Client struct {
	dialer  Dialer
	timeout time.Duration
}

// Option 客户端选项
type Option func(*Client)

// WithDialer 设置自定义拨号器
func WithDialer(dialer Dialer) Option {
	return func(c *Client) {
		c.dialer = dialer
	}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// parseQueryType 解析查询类型
func parseQueryType(queryType string) (uint16, error) {
	switch queryType {
	case "A":
		return mdns.TypeA, nil
	case "AAAA":
		return mdns.TypeAAAA, nil
	case "CNAME":
		return mdns.TypeCNAME, nil
	case "MX":
		return mdns.TypeMX, nil
	case "TXT":
		return mdns.TypeTXT, nil
	case "NS":
		return mdns.TypeNS, nil
	case "SOA":
		return mdns.TypeSOA, nil
	case "PTR":
		return mdns.TypePTR, nil
	case "SRV":
		return mdns.TypeSRV, nil
	case "CAA":
		return mdns.TypeCAA, nil
	default:
		return 0, fmt.Errorf("unsupported query type: %s", queryType)
	}
}

// createDNSMessage 创建 DNS 查询消息
func createDNSMessage(domain, queryType string) (*mdns.Msg, error) {
	qtype, err := parseQueryType(queryType)
	if err != nil {
		return nil, err
	}

	msg := new(mdns.Msg)
	msg.SetQuestion(mdns.Fqdn(domain), qtype)
	msg.RecursionDesired = true

	return msg, nil
}

// parseAnswer 解析 DNS 应答
func parseAnswer(msg *mdns.Msg) []string {
	var answers []string

	for _, answer := range msg.Answer {
		switch rr := answer.(type) {
		case *mdns.A:
			answers = append(answers, rr.A.String())
		case *mdns.AAAA:
			answers = append(answers, rr.AAAA.String())
		case *mdns.CNAME:
			answers = append(answers, rr.Target)
		case *mdns.MX:
			answers = append(answers, fmt.Sprintf("%d %s", rr.Preference, rr.Mx))
		case *mdns.TXT:
			for _, txt := range rr.Txt {
				answers = append(answers, txt)
			}
		case *mdns.NS:
			answers = append(answers, rr.Ns)
		case *mdns.SOA:
			answers = append(answers, fmt.Sprintf("%s %s %d %d %d %d %d",
				rr.Ns, rr.Mbox, rr.Serial, rr.Refresh, rr.Retry, rr.Expire, rr.Minttl))
		case *mdns.PTR:
			answers = append(answers, rr.Ptr)
		case *mdns.SRV:
			answers = append(answers, fmt.Sprintf("%d %d %d %s",
				rr.Priority, rr.Weight, rr.Port, rr.Target))
		case *mdns.CAA:
			answers = append(answers, fmt.Sprintf("%d %s %s",
				rr.Flag, rr.Tag, rr.Value))
		}
	}

	return answers
}

// BaseResolver provides default implementations for mihomo Resolver interface
type BaseResolver struct{}

// LookupIP implements resolver.Resolver
func (b *BaseResolver) LookupIP(ctx context.Context, host string) ([]netip.Addr, error) {
	return nil, fmt.Errorf("not implemented")
}

// LookupIPv4 implements resolver.Resolver
func (b *BaseResolver) LookupIPv4(ctx context.Context, host string) ([]netip.Addr, error) {
	return nil, fmt.Errorf("not implemented")
}

// LookupIPv6 implements resolver.Resolver
func (b *BaseResolver) LookupIPv6(ctx context.Context, host string) ([]netip.Addr, error) {
	return nil, fmt.Errorf("not implemented")
}

// ResolveECH implements resolver.Resolver
func (b *BaseResolver) ResolveECH(ctx context.Context, host string) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

// ExchangeContext implements resolver.Resolver
func (b *BaseResolver) ExchangeContext(ctx context.Context, m *mdns.Msg) (*mdns.Msg, error) {
	return nil, fmt.Errorf("not implemented")
}

// Invalid implements resolver.Resolver
func (b *BaseResolver) Invalid() bool {
	return false
}

// ClearCache implements resolver.Resolver
func (b *BaseResolver) ClearCache() {
	// No cache to clear in base implementation
}

// ResetConnection implements resolver.Resolver
func (b *BaseResolver) ResetConnection() {
	// No connection to reset in base implementation
}
