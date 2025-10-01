package dns

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	mdns "github.com/miekg/dns"
)

// mockDialer 用于测试的模拟拨号器
type mockDialer struct {
	dialFunc        func(network, address string) (net.Conn, error)
	dialContextFunc func(ctx context.Context, network, address string) (net.Conn, error)
}

func (m *mockDialer) Dial(network, address string) (net.Conn, error) {
	if m.dialFunc != nil {
		return m.dialFunc(network, address)
	}
	return net.Dial(network, address)
}

func (m *mockDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	if m.dialContextFunc != nil {
		return m.dialContextFunc(ctx, network, address)
	}
	dialer := &net.Dialer{}
	return dialer.DialContext(ctx, network, address)
}

// mockConn 用于测试的模拟连接
type mockConn struct {
	readFunc        func(b []byte) (int, error)
	writeFunc       func(b []byte) (int, error)
	closeFunc       func() error
	setDeadlineFunc func(t time.Time) error
}

func (m *mockConn) Read(b []byte) (int, error) {
	if m.readFunc != nil {
		return m.readFunc(b)
	}
	return 0, nil
}

func (m *mockConn) Write(b []byte) (int, error) {
	if m.writeFunc != nil {
		return m.writeFunc(b)
	}
	return len(b), nil
}

func (m *mockConn) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return nil
}

func (m *mockConn) LocalAddr() net.Addr {
	return &net.TCPAddr{}
}

func (m *mockConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{}
}

func (m *mockConn) SetDeadline(t time.Time) error {
	if m.setDeadlineFunc != nil {
		return m.setDeadlineFunc(t)
	}
	return nil
}

func (m *mockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

// 测试用的 DNS 服务器
const (
	aliDNS     = "223.5.5.5"
	aliDNS2    = "223.6.6.6"
	tencentDNS = "119.29.29.29"
	aliDoH     = "https://dns.alidns.com/dns-query"
	tencentDoH = "https://doh.pub/dns-query"
	aliDoT     = "dns.alidns.com"
	tencentDoT = "dot.pub"
	testDomain = "www.baidu.com"
)

func TestNewResolver(t *testing.T) {
	tests := []struct {
		protocol string
		expected string
	}{
		{"udp", "udp"},
		{"tcp", "tcp"},
		{"doh", "doh"},
		{"https", "doh"},
		{"dot", "dot"},
		{"tls", "dot"},
		{"doh3", "doh3"},
		{"h3", "doh3"},
		{"unknown", "udp"}, // fallback to udp
	}

	for _, tt := range tests {
		t.Run(tt.protocol, func(t *testing.T) {
			resolver := NewResolver(tt.protocol)
			if resolver.Type() != tt.expected {
				t.Errorf("expected type %s, got %s", tt.expected, resolver.Type())
			}
		})
	}
}

func TestDefaultDialer(t *testing.T) {
	dialer := &DefaultDialer{}

	// Test Dial
	conn, err := dialer.Dial("udp", aliDNS+":53")
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	conn.Close()

	// Test DialContext
	ctx := context.Background()
	conn, err = dialer.DialContext(ctx, "udp", aliDNS+":53")
	if err != nil {
		t.Fatalf("DialContext failed: %v", err)
	}
	conn.Close()
}

func TestParseQueryType(t *testing.T) {
	tests := []struct {
		queryType string
		wantErr   bool
	}{
		{"A", false},
		{"AAAA", false},
		{"CNAME", false},
		{"MX", false},
		{"TXT", false},
		{"NS", false},
		{"SOA", false},
		{"PTR", false},
		{"SRV", false},
		{"CAA", false},
		{"INVALID", true},
	}

	for _, tt := range tests {
		t.Run(tt.queryType, func(t *testing.T) {
			_, err := parseQueryType(tt.queryType)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseQueryType(%s) error = %v, wantErr %v", tt.queryType, err, tt.wantErr)
			}
		})
	}
}

func TestCreateDNSMessage(t *testing.T) {
	msg, err := createDNSMessage(testDomain, "A")
	if err != nil {
		t.Fatalf("createDNSMessage failed: %v", err)
	}

	if len(msg.Question) == 0 {
		t.Error("expected at least one question")
	}

	if !msg.RecursionDesired {
		t.Error("expected RecursionDesired to be true")
	}
}

func TestCreateDNSMessageWithInvalidType(t *testing.T) {
	_, err := createDNSMessage(testDomain, "INVALID")
	if err == nil {
		t.Error("expected error for invalid query type")
	}
}

func TestUDPResolver(t *testing.T) {
	resolver := NewUDPResolver()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	servers := []string{aliDNS, aliDNS2, tencentDNS}

	for _, server := range servers {
		t.Run(server, func(t *testing.T) {
			result, err := resolver.Query(ctx, testDomain, "A", server)
			if err != nil {
				t.Logf("UDP query to %s failed: %v", server, err)
				t.Skip("skipping UDP test")
				return
			}

			if result.Question != testDomain {
				t.Errorf("expected question %s, got %s", testDomain, result.Question)
			}

			if result.Type != "A" {
				t.Errorf("expected type A, got %s", result.Type)
			}

			if len(result.Answers) == 0 {
				t.Error("expected at least one answer")
			}

			if result.Protocol != "udp" {
				t.Errorf("expected protocol udp, got %s", result.Protocol)
			}

			t.Logf("UDP query result: %v", result.Answers)
		})
	}
}

func TestTCPResolver(t *testing.T) {
	resolver := NewTCPResolver()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	servers := []string{aliDNS, tencentDNS}

	for _, server := range servers {
		t.Run(server, func(t *testing.T) {
			result, err := resolver.Query(ctx, testDomain, "A", server)
			if err != nil {
				t.Logf("TCP query to %s failed: %v", server, err)
				t.Skip("skipping TCP test")
				return
			}

			if result.Question != testDomain {
				t.Errorf("expected question %s, got %s", testDomain, result.Question)
			}

			if len(result.Answers) == 0 {
				t.Error("expected at least one answer")
			}

			if result.Protocol != "tcp" {
				t.Errorf("expected protocol tcp, got %s", result.Protocol)
			}

			t.Logf("TCP query result: %v", result.Answers)
		})
	}
}

func TestDoHResolver(t *testing.T) {
	resolver := NewDoHResolver(WithTimeout(15 * time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	servers := []string{aliDoH, tencentDoH}

	for _, server := range servers {
		t.Run(server, func(t *testing.T) {
			result, err := resolver.Query(ctx, testDomain, "A", server)
			if err != nil {
				t.Logf("DoH query to %s failed: %v", server, err)
				t.Skip("skipping DoH test")
				return
			}

			if len(result.Answers) == 0 {
				t.Error("expected at least one answer")
			}

			if result.Protocol != "doh" {
				t.Errorf("expected protocol doh, got %s", result.Protocol)
			}

			t.Logf("DoH query result: %v", result.Answers)
		})
	}
}

func TestDoTResolver(t *testing.T) {
	resolver := NewDoTResolver(WithTimeout(15 * time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	servers := []string{aliDoT, tencentDoT}

	for _, server := range servers {
		t.Run(server, func(t *testing.T) {
			result, err := resolver.Query(ctx, testDomain, "A", server)
			if err != nil {
				t.Logf("DoT query to %s failed: %v", server, err)
				t.Skip("skipping DoT test")
				return
			}

			if len(result.Answers) == 0 {
				t.Error("expected at least one answer")
			}

			if result.Protocol != "dot" {
				t.Errorf("expected protocol dot, got %s", result.Protocol)
			}

			t.Logf("DoT query result: %v", result.Answers)
		})
	}
}

func TestDoH3Resolver(t *testing.T) {
	resolver := NewDoH3Resolver(WithTimeout(15 * time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	result, err := resolver.Query(ctx, testDomain, "A", aliDoH)
	if err != nil {
		t.Logf("DoH3 query failed: %v", err)
		t.Skip("skipping DoH3 test")
		return
	}

	if len(result.Answers) == 0 {
		t.Error("expected at least one answer")
	}

	if result.Protocol != "doh3" {
		t.Errorf("expected protocol doh3, got %s", result.Protocol)
	}
}

func TestResolverWithOptions(t *testing.T) {
	customDialer := &DefaultDialer{}
	customTimeout := 10 * time.Second

	resolver := NewUDPResolver(
		WithDialer(customDialer),
		WithTimeout(customTimeout),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := resolver.Query(ctx, testDomain, "A", aliDNS)
	if err != nil {
		t.Fatalf("Query with custom options failed: %v", err)
	}

	if len(result.Answers) == 0 {
		t.Error("expected at least one answer")
	}
}

func TestResolverWithFailingDialer(t *testing.T) {
	dialer := &mockDialer{
		dialContextFunc: func(ctx context.Context, network, address string) (net.Conn, error) {
			return nil, errors.New("dial failed")
		},
	}

	resolver := NewUDPResolver(WithDialer(dialer))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := resolver.Query(ctx, testDomain, "A", aliDNS)
	if err == nil {
		t.Error("expected error with failing dialer")
	}
}

func TestResolverWithCanceledContext(t *testing.T) {
	resolver := NewUDPResolver()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := resolver.Query(ctx, testDomain, "A", aliDNS)
	if err == nil {
		t.Error("expected error with canceled context")
	}
}

func TestDifferentQueryTypes(t *testing.T) {
	resolver := NewUDPResolver()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tests := []struct {
		domain    string
		queryType string
	}{
		{testDomain, "A"},
		{testDomain, "AAAA"},
		{"baidu.com", "MX"},
		{"baidu.com", "TXT"},
		{"baidu.com", "NS"},
		{"baidu.com", "SOA"},
	}

	for _, tt := range tests {
		t.Run(tt.queryType, func(t *testing.T) {
			result, err := resolver.Query(ctx, tt.domain, tt.queryType, aliDNS)
			if err != nil {
				t.Logf("Query type %s failed: %v", tt.queryType, err)
				t.Skip("skipping query type test")
				return
			}

			if result.Type != tt.queryType {
				t.Errorf("expected type %s, got %s", tt.queryType, result.Type)
			}

			t.Logf("Query type %s result: %v", tt.queryType, result.Answers)
		})
	}
}

func TestTCPResolverWithMockDialer(t *testing.T) {
	// Create a valid DNS response
	msg := new(mdns.Msg)
	msg.SetReply(&mdns.Msg{Question: []mdns.Question{{Name: testDomain + ".", Qtype: mdns.TypeA, Qclass: mdns.ClassINET}}})
	msg.Answer = append(msg.Answer, &mdns.A{
		Hdr: mdns.RR_Header{Name: testDomain + ".", Rrtype: mdns.TypeA, Class: mdns.ClassINET, Ttl: 300},
		A:   net.ParseIP("1.2.3.4"),
	})
	packed, _ := msg.Pack()

	// TCP DNS response with 2-byte length prefix
	tcpResponse := make([]byte, 2+len(packed))
	tcpResponse[0] = byte(len(packed) >> 8)
	tcpResponse[1] = byte(len(packed))
	copy(tcpResponse[2:], packed)

	readOffset := 0
	mockConn := &mockConn{
		readFunc: func(b []byte) (int, error) {
			if readOffset >= len(tcpResponse) {
				return 0, net.ErrClosed
			}
			n := copy(b, tcpResponse[readOffset:])
			readOffset += n
			return n, nil
		},
		writeFunc: func(b []byte) (int, error) {
			return len(b), nil
		},
		closeFunc: func() error {
			return nil
		},
		setDeadlineFunc: func(t time.Time) error {
			return nil
		},
	}

	dialer := &mockDialer{
		dialContextFunc: func(ctx context.Context, network, address string) (net.Conn, error) {
			readOffset = 0 // Reset for each dial
			return mockConn, nil
		},
	}

	resolver := NewTCPResolver(WithDialer(dialer))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := resolver.Query(ctx, testDomain, "A", aliDNS)
	if err != nil {
		t.Fatalf("TCP query with mock dialer failed: %v", err)
	}

	if len(result.Answers) == 0 {
		t.Error("expected at least one answer")
	}
}

func TestDoHResolverErrors(t *testing.T) {
	resolver := NewDoHResolver()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test with invalid query type
	_, err := resolver.Query(ctx, testDomain, "INVALID", aliDoH)
	if err == nil {
		t.Error("expected error for invalid query type")
	}
}

func TestDoTResolverErrors(t *testing.T) {
	resolver := NewDoTResolver()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test with invalid query type
	_, err := resolver.Query(ctx, testDomain, "INVALID", aliDoT)
	if err == nil {
		t.Error("expected error for invalid query type")
	}
}

func TestUDPResolverErrors(t *testing.T) {
	resolver := NewUDPResolver()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test with invalid query type
	_, err := resolver.Query(ctx, testDomain, "INVALID", aliDNS)
	if err == nil {
		t.Error("expected error for invalid query type")
	}
}

func TestDoH3ResolverErrors(t *testing.T) {
	resolver := NewDoH3Resolver()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test with invalid query type
	_, err := resolver.Query(ctx, testDomain, "INVALID", aliDoH)
	if err == nil {
		t.Error("expected error for invalid query type")
	}
}

func TestParseDifferentRecordTypes(t *testing.T) {
	resolver := NewUDPResolver()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tests := []struct {
		domain    string
		queryType string
	}{
		{"baidu.com", "PTR"},
		{"_xmpp-server._tcp.google.com", "SRV"},
		{"google.com", "CAA"},
	}

	for _, tt := range tests {
		t.Run(tt.queryType, func(t *testing.T) {
			_, err := resolver.Query(ctx, tt.domain, tt.queryType, aliDNS)
			// These might fail or return empty results, but should not crash
			if err != nil {
				t.Logf("Query type %s: %v", tt.queryType, err)
			}
		})
	}
}

func TestParseAnswerWithAllRecordTypes(t *testing.T) {
	// Test all DNS record types directly
	tests := []struct {
		name   string
		answer mdns.RR
		want   string
	}{
		{
			name:   "PTR",
			answer: &mdns.PTR{Ptr: "example.com."},
			want:   "example.com.",
		},
		{
			name:   "SRV",
			answer: &mdns.SRV{Priority: 10, Weight: 20, Port: 80, Target: "server.example.com."},
			want:   "10 20 80 server.example.com.",
		},
		{
			name:   "CAA",
			answer: &mdns.CAA{Flag: 0, Tag: "issue", Value: "letsencrypt.org"},
			want:   "0 issue letsencrypt.org",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &mdns.Msg{
				Answer: []mdns.RR{tt.answer},
			}
			answers := parseAnswer(msg)
			if len(answers) != 1 {
				t.Fatalf("expected 1 answer, got %d", len(answers))
			}
			if answers[0] != tt.want {
				t.Errorf("expected %q, got %q", tt.want, answers[0])
			}
		})
	}
}

func TestUDPResolverWithDeadlineError(t *testing.T) {
	mockConn := &mockConn{
		setDeadlineFunc: func(t time.Time) error {
			return errors.New("deadline error")
		},
	}

	dialer := &mockDialer{
		dialContextFunc: func(ctx context.Context, network, address string) (net.Conn, error) {
			return mockConn, nil
		},
	}

	resolver := NewUDPResolver(WithDialer(dialer))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := resolver.Query(ctx, testDomain, "A", aliDNS)
	if err == nil {
		t.Error("expected error for deadline setting failure")
	}
}

func TestTCPResolverWithDeadlineError(t *testing.T) {
	mockConn := &mockConn{
		setDeadlineFunc: func(t time.Time) error {
			return errors.New("deadline error")
		},
	}

	dialer := &mockDialer{
		dialContextFunc: func(ctx context.Context, network, address string) (net.Conn, error) {
			return mockConn, nil
		},
	}

	resolver := NewTCPResolver(WithDialer(dialer))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := resolver.Query(ctx, testDomain, "A", aliDNS)
	if err == nil {
		t.Error("expected error for deadline setting failure")
	}
}

func TestDoTResolverWithTLSError(t *testing.T) {
	dialer := &mockDialer{
		dialContextFunc: func(ctx context.Context, network, address string) (net.Conn, error) {
			// Return a connection that will fail TLS handshake
			return &mockConn{}, nil
		},
	}

	resolver := NewDoTResolver(WithDialer(dialer))
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := resolver.Query(ctx, testDomain, "A", "invalid.server:853")
	if err == nil {
		t.Error("expected error for TLS handshake failure")
	}
}

func TestDoHResolverWithShortURL(t *testing.T) {
	resolver := NewDoHResolver()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test with short URL (without https prefix)
	result, err := resolver.Query(ctx, testDomain, "A", "dns.alidns.com")
	if err != nil {
		t.Logf("DoH with short URL failed: %v", err)
		// This might fail due to network, but should handle URL formatting
	} else if len(result.Answers) > 0 {
		t.Logf("DoH with short URL succeeded: %v", result.Answers)
	}
}

func TestUDPResolverWithServerWithoutPort(t *testing.T) {
	resolver := NewUDPResolver()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test with server address without port
	result, err := resolver.Query(ctx, testDomain, "A", aliDNS)
	if err != nil {
		t.Fatalf("UDP query without port failed: %v", err)
	}

	if len(result.Answers) == 0 {
		t.Error("expected at least one answer")
	}
}

func TestTCPResolverWithServerWithoutPort(t *testing.T) {
	// Create a successful mock TCP connection
	msg := new(mdns.Msg)
	msg.SetReply(&mdns.Msg{Question: []mdns.Question{{Name: testDomain + ".", Qtype: mdns.TypeA, Qclass: mdns.ClassINET}}})
	msg.Answer = append(msg.Answer, &mdns.A{
		Hdr: mdns.RR_Header{Name: testDomain + ".", Rrtype: mdns.TypeA, Class: mdns.ClassINET, Ttl: 300},
		A:   net.ParseIP("1.2.3.4"),
	})
	packed, _ := msg.Pack()

	tcpResponse := make([]byte, 2+len(packed))
	tcpResponse[0] = byte(len(packed) >> 8)
	tcpResponse[1] = byte(len(packed))
	copy(tcpResponse[2:], packed)

	readOffset := 0
	mockConn := &mockConn{
		readFunc: func(b []byte) (int, error) {
			if readOffset >= len(tcpResponse) {
				return 0, net.ErrClosed
			}
			n := copy(b, tcpResponse[readOffset:])
			readOffset += n
			return n, nil
		},
		writeFunc: func(b []byte) (int, error) {
			return len(b), nil
		},
		closeFunc: func() error {
			return nil
		},
		setDeadlineFunc: func(t time.Time) error {
			return nil
		},
	}

	dialer := &mockDialer{
		dialContextFunc: func(ctx context.Context, network, address string) (net.Conn, error) {
			readOffset = 0
			return mockConn, nil
		},
	}

	resolver := NewTCPResolver(WithDialer(dialer))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test with server address without port (should add :53)
	result, err := resolver.Query(ctx, testDomain, "A", aliDNS)
	if err != nil {
		t.Fatalf("TCP query without port failed: %v", err)
	}

	if len(result.Answers) == 0 {
		t.Error("expected at least one answer")
	}
}

func TestDoTResolverWithServerWithoutPort(t *testing.T) {
	resolver := NewDoTResolver()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test with server address without port (should add :853)
	result, err := resolver.Query(ctx, testDomain, "A", aliDoT)
	if err != nil {
		t.Logf("DoT query without port failed: %v", err)
		t.Skip("skipping DoT test")
		return
	}

	if len(result.Answers) == 0 {
		t.Error("expected at least one answer")
	}
}

func TestMultipleResolvers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	servers := []struct {
		name     string
		resolver Resolver
		server   string
	}{
		{"UDP Ali", NewUDPResolver(), aliDNS},
		{"TCP Ali", NewTCPResolver(), aliDNS},
		{"UDP Tencent", NewUDPResolver(), tencentDNS},
		{"TCP Tencent", NewTCPResolver(), tencentDNS},
		{"DoH Ali", NewDoHResolver(WithTimeout(15 * time.Second)), aliDoH},
		{"DoH Tencent", NewDoHResolver(WithTimeout(15 * time.Second)), tencentDoH},
		{"DoT Ali", NewDoTResolver(WithTimeout(15 * time.Second)), aliDoT},
		{"DoT Tencent", NewDoTResolver(WithTimeout(15 * time.Second)), tencentDoT},
	}

	successCount := 0
	for _, s := range servers {
		t.Run(s.name, func(t *testing.T) {
			result, err := s.resolver.Query(ctx, testDomain, "A", s.server)
			if err != nil {
				t.Logf("%s query failed: %v", s.name, err)
				return
			}

			if len(result.Answers) > 0 {
				successCount++
				t.Logf("%s query succeeded: %v (duration: %v)", s.name, result.Answers, result.Duration)
			}
		})
	}

	if successCount == 0 {
		t.Error("expected at least one resolver to succeed")
	}

	t.Logf("Successfully queried %d/%d resolvers", successCount, len(servers))
}
