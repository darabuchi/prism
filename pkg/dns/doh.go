package dns

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"

	mdns "github.com/miekg/dns"
)

// DoHResolver DNS over HTTPS 解析器
type DoHResolver struct {
	BaseResolver
	dialer  Dialer
	timeout time.Duration
}

// NewDoHResolver 创建 DoH DNS 解析器
func NewDoHResolver(opts ...Option) *DoHResolver {
	c := &Client{
		dialer:  &DefaultDialer{},
		timeout: 10 * time.Second,
	}

	for _, opt := range opts {
		opt(c)
	}

	return &DoHResolver{
		dialer:  c.dialer,
		timeout: c.timeout,
	}
}

// Type 返回解析器类型
func (r *DoHResolver) Type() string {
	return "doh"
}

// Query 执行 DoH DNS 查询
func (r *DoHResolver) Query(ctx context.Context, domain, queryType, server string) (*QueryResult, error) {
	start := time.Now()

	msg, err := createDNSMessage(domain, queryType)
	if err != nil {
		return nil, err
	}

	// 打包 DNS 消息
	packed, err := msg.Pack()
	if err != nil {
		return nil, fmt.Errorf("failed to pack DNS message: %w", err)
	}

	// 确保服务器地址是完整的 URL
	if server[0] != 'h' {
		server = "https://" + server + "/dns-query"
	}

	// 创建 HTTP 客户端
	transport := &http.Transport{
		DialContext: r.dialer.DialContext,
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   r.timeout,
	}

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "POST", server, bytes.NewReader(packed))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/dns-message")
	req.Header.Set("Accept", "application/dns-message")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send DoH request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DoH request failed with status: %d", resp.StatusCode)
	}

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read DoH response: %w", err)
	}

	// 解析响应
	response := new(mdns.Msg)
	if err := response.Unpack(body); err != nil {
		return nil, fmt.Errorf("failed to unpack DNS response: %w", err)
	}

	return &QueryResult{
		Question:   domain,
		Type:       queryType,
		Answers:    parseAnswer(response),
		Duration:   time.Since(start),
		Server:     server,
		Protocol:   "doh",
		StatusCode: response.Rcode,
	}, nil
}
