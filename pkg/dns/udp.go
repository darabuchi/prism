package dns

import (
	"context"
	"fmt"
	"net"
	"time"

	mdns "github.com/miekg/dns"
)

// UDPResolver UDP DNS 解析器
type UDPResolver struct {
	BaseResolver
	dialer  Dialer
	timeout time.Duration
}

// NewUDPResolver 创建 UDP DNS 解析器
func NewUDPResolver(opts ...Option) *UDPResolver {
	c := &Client{
		dialer:  &DefaultDialer{},
		timeout: 5 * time.Second,
	}

	for _, opt := range opts {
		opt(c)
	}

	return &UDPResolver{
		dialer:  c.dialer,
		timeout: c.timeout,
	}
}

// Type 返回解析器类型
func (r *UDPResolver) Type() string {
	return "udp"
}

// Query 执行 UDP DNS 查询
func (r *UDPResolver) Query(ctx context.Context, domain, queryType, server string) (*QueryResult, error) {
	start := time.Now()

	msg, err := createDNSMessage(domain, queryType)
	if err != nil {
		return nil, err
	}

	// 确保服务器地址包含端口
	if _, _, err := net.SplitHostPort(server); err != nil {
		server = net.JoinHostPort(server, "53")
	}

	// 创建 UDP 连接
	conn, err := r.dialer.DialContext(ctx, "udp", server)
	if err != nil {
		return nil, fmt.Errorf("failed to dial UDP: %w", err)
	}
	defer conn.Close()

	// 设置超时
	if err := conn.SetDeadline(time.Now().Add(r.timeout)); err != nil {
		return nil, fmt.Errorf("failed to set deadline: %w", err)
	}

	// 创建 DNS 客户端
	dnsConn := &mdns.Conn{Conn: conn}
	defer dnsConn.Close()

	// 发送查询
	if err := dnsConn.WriteMsg(msg); err != nil {
		return nil, fmt.Errorf("failed to write DNS message: %w", err)
	}

	// 接收响应
	response, err := dnsConn.ReadMsg()
	if err != nil {
		return nil, fmt.Errorf("failed to read DNS response: %w", err)
	}

	return &QueryResult{
		Question:   domain,
		Type:       queryType,
		Answers:    parseAnswer(response),
		Duration:   time.Since(start),
		Server:     server,
		Protocol:   "udp",
		StatusCode: response.Rcode,
	}, nil
}
