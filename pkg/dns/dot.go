package dns

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	mdns "github.com/miekg/dns"
)

// DoTResolver DNS over TLS 解析器
type DoTResolver struct {
	BaseResolver
	dialer  Dialer
	timeout time.Duration
}

// NewDoTResolver 创建 DoT DNS 解析器
func NewDoTResolver(opts ...Option) *DoTResolver {
	c := &Client{
		dialer:  &DefaultDialer{},
		timeout: 10 * time.Second,
	}

	for _, opt := range opts {
		opt(c)
	}

	return &DoTResolver{
		dialer:  c.dialer,
		timeout: c.timeout,
	}
}

// Type 返回解析器类型
func (r *DoTResolver) Type() string {
	return "dot"
}

// Query 执行 DoT DNS 查询
func (r *DoTResolver) Query(ctx context.Context, domain, queryType, server string) (*QueryResult, error) {
	start := time.Now()

	msg, err := createDNSMessage(domain, queryType)
	if err != nil {
		return nil, err
	}

	// 确保服务器地址包含端口
	if _, _, err := net.SplitHostPort(server); err != nil {
		server = net.JoinHostPort(server, "853")
	}

	// 创建 TCP 连接
	conn, err := r.dialer.DialContext(ctx, "tcp", server)
	if err != nil {
		return nil, fmt.Errorf("failed to dial TCP for DoT: %w", err)
	}
	defer conn.Close()

	// 提取主机名用于 TLS
	host, _, _ := net.SplitHostPort(server)

	// 升级为 TLS 连接
	tlsConn := tls.Client(conn, &tls.Config{
		ServerName: host,
		MinVersion: tls.VersionTLS12,
	})

	// TLS 握手
	if err := tlsConn.HandshakeContext(ctx); err != nil {
		return nil, fmt.Errorf("TLS handshake failed: %w", err)
	}
	defer tlsConn.Close()

	// 设置超时
	if err := tlsConn.SetDeadline(time.Now().Add(r.timeout)); err != nil {
		return nil, fmt.Errorf("failed to set deadline: %w", err)
	}

	// 创建 DNS 客户端
	dnsConn := &mdns.Conn{Conn: tlsConn}
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
		Protocol:   "dot",
		StatusCode: response.Rcode,
	}, nil
}
