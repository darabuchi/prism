package rules

import (
	"net/netip"
	"os"
	"testing"

	prism "github.com/darabuchi/prism"
)

func TestDomain(t *testing.T) {
	rule := NewDomain("google.com", ActionProxy)

	tests := []struct {
		name     string
		metadata *Metadata
		want     bool
	}{
		{
			name:     "exact match",
			metadata: &Metadata{Host: "google.com"},
			want:     true,
		},
		{
			name:     "case insensitive",
			metadata: &Metadata{Host: "Google.COM"},
			want:     true,
		},
		{
			name:     "subdomain no match",
			metadata: &Metadata{Host: "www.google.com"},
			want:     false,
		},
		{
			name:     "no match",
			metadata: &Metadata{Host: "example.com"},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Match(tt.metadata); got != tt.want {
				t.Errorf("Domain.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainSuffix(t *testing.T) {
	rule := NewDomainSuffix("google.com", ActionProxy)

	tests := []struct {
		name     string
		metadata *Metadata
		want     bool
	}{
		{
			name:     "exact match",
			metadata: &Metadata{Host: "google.com"},
			want:     true,
		},
		{
			name:     "subdomain match",
			metadata: &Metadata{Host: "www.google.com"},
			want:     true,
		},
		{
			name:     "deep subdomain match",
			metadata: &Metadata{Host: "mail.dev.google.com"},
			want:     true,
		},
		{
			name:     "no match",
			metadata: &Metadata{Host: "example.com"},
			want:     false,
		},
		{
			name:     "partial no match",
			metadata: &Metadata{Host: "notgoogle.com"},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Match(tt.metadata); got != tt.want {
				t.Errorf("DomainSuffix.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainKeyword(t *testing.T) {
	rule := NewDomainKeyword("google", ActionProxy)

	tests := []struct {
		name     string
		metadata *Metadata
		want     bool
	}{
		{
			name:     "contains keyword",
			metadata: &Metadata{Host: "google.com"},
			want:     true,
		},
		{
			name:     "contains in subdomain",
			metadata: &Metadata{Host: "www.google.com"},
			want:     true,
		},
		{
			name:     "contains in middle",
			metadata: &Metadata{Host: "mygooglemail.com"},
			want:     true,
		},
		{
			name:     "no match",
			metadata: &Metadata{Host: "example.com"},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Match(tt.metadata); got != tt.want {
				t.Errorf("DomainKeyword.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPCIDR(t *testing.T) {
	rule, err := NewIPCIDR("192.168.0.0/16", ActionDirect, false)
	if err != nil {
		t.Fatalf("NewIPCIDR failed: %v", err)
	}

	tests := []struct {
		name     string
		metadata *Metadata
		want     bool
	}{
		{
			name: "in range",
			metadata: &Metadata{
				DstIP: netip.MustParseAddr("192.168.1.1"),
			},
			want: true,
		},
		{
			name: "edge case - start",
			metadata: &Metadata{
				DstIP: netip.MustParseAddr("192.168.0.0"),
			},
			want: true,
		},
		{
			name: "edge case - end",
			metadata: &Metadata{
				DstIP: netip.MustParseAddr("192.168.255.255"),
			},
			want: true,
		},
		{
			name: "out of range",
			metadata: &Metadata{
				DstIP: netip.MustParseAddr("192.169.1.1"),
			},
			want: false,
		},
		{
			name: "different network",
			metadata: &Metadata{
				DstIP: netip.MustParseAddr("10.0.0.1"),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Match(tt.metadata); got != tt.want {
				t.Errorf("IPCIDR.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPort(t *testing.T) {
	rule, err := NewPort("80,443,8000-9000", ActionDirect, true)
	if err != nil {
		t.Fatalf("NewPort failed: %v", err)
	}

	tests := []struct {
		name     string
		metadata *Metadata
		want     bool
	}{
		{
			name:     "single port 80",
			metadata: &Metadata{DstPort: 80},
			want:     true,
		},
		{
			name:     "single port 443",
			metadata: &Metadata{DstPort: 443},
			want:     true,
		},
		{
			name:     "in range - start",
			metadata: &Metadata{DstPort: 8000},
			want:     true,
		},
		{
			name:     "in range - middle",
			metadata: &Metadata{DstPort: 8500},
			want:     true,
		},
		{
			name:     "in range - end",
			metadata: &Metadata{DstPort: 9000},
			want:     true,
		},
		{
			name:     "not in list",
			metadata: &Metadata{DstPort: 22},
			want:     false,
		},
		{
			name:     "out of range",
			metadata: &Metadata{DstPort: 9001},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Match(tt.metadata); got != tt.want {
				t.Errorf("Port.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcess(t *testing.T) {
	rule := NewProcess("chrome", ActionProxy, false)

	tests := []struct {
		name     string
		metadata *Metadata
		want     bool
	}{
		{
			name:     "exact match",
			metadata: &Metadata{ProcessName: "chrome"},
			want:     true,
		},
		{
			name:     "case insensitive",
			metadata: &Metadata{ProcessName: "Chrome"},
			want:     true,
		},
		{
			name:     "no match",
			metadata: &Metadata{ProcessName: "firefox"},
			want:     false,
		},
		{
			name:     "empty process",
			metadata: &Metadata{ProcessName: ""},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Match(tt.metadata); got != tt.want {
				t.Errorf("Process.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatch(t *testing.T) {
	rule := NewMatch(ActionProxy)

	tests := []struct {
		name     string
		metadata *Metadata
	}{
		{
			name:     "always match - empty",
			metadata: &Metadata{},
		},
		{
			name: "always match - with data",
			metadata: &Metadata{
				Host:    "google.com",
				DstPort: 443,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Match(tt.metadata); !got {
				t.Errorf("Match.Match() = %v, want true", got)
			}
		})
	}
}

func TestParseRule(t *testing.T) {
	tests := []struct {
		name    string
		ruleStr string
		wantErr bool
		check   func(Rule) bool
	}{
		{
			name:    "domain rule",
			ruleStr: "DOMAIN,google.com,PROXY",
			wantErr: false,
			check: func(r Rule) bool {
				return r.Type() == TypeDomain && r.Action() == ActionProxy && r.Payload() == "google.com"
			},
		},
		{
			name:    "domain suffix rule",
			ruleStr: "DOMAIN-SUFFIX,google.com,DIRECT",
			wantErr: false,
			check: func(r Rule) bool {
				return r.Type() == TypeDomainSuffix && r.Action() == ActionDirect
			},
		},
		{
			name:    "ip cidr rule",
			ruleStr: "IP-CIDR,192.168.0.0/16,DIRECT",
			wantErr: false,
			check: func(r Rule) bool {
				return r.Type() == TypeIPCIDR && r.Action() == ActionDirect
			},
		},
		{
			name:    "match rule",
			ruleStr: "MATCH,PROXY",
			wantErr: false,
			check: func(r Rule) bool {
				return r.Type() == TypeMatch && r.Action() == ActionProxy
			},
		},
		{
			name:    "comment",
			ruleStr: "# This is a comment",
			wantErr: false,
			check:   func(r Rule) bool { return r == nil },
		},
		{
			name:    "empty line",
			ruleStr: "",
			wantErr: false,
			check:   func(r Rule) bool { return r == nil },
		},
		{
			name:    "invalid format",
			ruleStr: "DOMAIN",
			wantErr: true,
		},
		{
			name:    "invalid rule type",
			ruleStr: "UNKNOWN,test,PROXY",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule, err := ParseRule(tt.ruleStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil && !tt.check(rule) {
				t.Errorf("ParseRule() returned unexpected rule: %v", rule)
			}
		})
	}
}

func TestMatchFirst(t *testing.T) {
	rules := []Rule{
		NewDomain("google.com", ActionProxy),
		NewDomainSuffix("example.com", ActionDirect),
		NewMatch(ActionProxy),
	}

	tests := []struct {
		name     string
		metadata *Metadata
		wantType RuleType
		wantOk   bool
	}{
		{
			name:     "match domain",
			metadata: &Metadata{Host: "google.com"},
			wantType: TypeDomain,
			wantOk:   true,
		},
		{
			name:     "match suffix",
			metadata: &Metadata{Host: "www.example.com"},
			wantType: TypeDomainSuffix,
			wantOk:   true,
		},
		{
			name:     "match all",
			metadata: &Metadata{Host: "other.com"},
			wantType: TypeMatch,
			wantOk:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule, ok := MatchFirst(rules, tt.metadata)
			if ok != tt.wantOk {
				t.Errorf("MatchFirst() ok = %v, want %v", ok, tt.wantOk)
			}
			if ok && rule.Type() != tt.wantType {
				t.Errorf("MatchFirst() type = %v, want %v", rule.Type(), tt.wantType)
			}
		})
	}
}

// ===== 入站规则测试 =====

func TestInType(t *testing.T) {
	rule := NewInType("HTTP", ActionProxy)

	tests := []struct {
		name     string
		metadata *Metadata
		want     bool
	}{
		{
			name:     "exact match",
			metadata: &Metadata{InboundType: "HTTP"},
			want:     true,
		},
		{
			name:     "case insensitive",
			metadata: &Metadata{InboundType: "http"},
			want:     true,
		},
		{
			name:     "no match",
			metadata: &Metadata{InboundType: "SOCKS5"},
			want:     false,
		},
		{
			name:     "empty inbound type",
			metadata: &Metadata{InboundType: ""},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Match(tt.metadata); got != tt.want {
				t.Errorf("InType.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInName(t *testing.T) {
	rule := NewInName("proxy-in", ActionDirect)

	tests := []struct {
		name     string
		metadata *Metadata
		want     bool
	}{
		{
			name:     "exact match",
			metadata: &Metadata{InboundName: "proxy-in"},
			want:     true,
		},
		{
			name:     "case insensitive",
			metadata: &Metadata{InboundName: "Proxy-In"},
			want:     true,
		},
		{
			name:     "no match",
			metadata: &Metadata{InboundName: "other-in"},
			want:     false,
		},
		{
			name:     "empty inbound name",
			metadata: &Metadata{InboundName: ""},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Match(tt.metadata); got != tt.want {
				t.Errorf("InName.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInUser(t *testing.T) {
	rule := NewInUser("alice", ActionProxy)

	tests := []struct {
		name     string
		metadata *Metadata
		want     bool
	}{
		{
			name:     "exact match",
			metadata: &Metadata{InboundUser: "alice"},
			want:     true,
		},
		{
			name:     "no match - different user",
			metadata: &Metadata{InboundUser: "bob"},
			want:     false,
		},
		{
			name:     "no match - case sensitive",
			metadata: &Metadata{InboundUser: "Alice"},
			want:     false,
		},
		{
			name:     "empty user",
			metadata: &Metadata{InboundUser: ""},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Match(tt.metadata); got != tt.want {
				t.Errorf("InUser.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ===== 网络层规则测试 =====

func TestNetwork(t *testing.T) {
	rule := NewNetwork("TCP", ActionDirect)

	tests := []struct {
		name     string
		metadata *Metadata
		want     bool
	}{
		{
			name:     "tcp match",
			metadata: &Metadata{Network: "TCP"},
			want:     true,
		},
		{
			name:     "case insensitive",
			metadata: &Metadata{Network: "tcp"},
			want:     true,
		},
		{
			name:     "no match - udp",
			metadata: &Metadata{Network: "UDP"},
			want:     false,
		},
		{
			name:     "empty network",
			metadata: &Metadata{Network: ""},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Match(tt.metadata); got != tt.want {
				t.Errorf("Network.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUID(t *testing.T) {
	rule, err := NewUID("1000", ActionDirect)
	if err != nil {
		t.Fatalf("NewUID failed: %v", err)
	}

	tests := []struct {
		name     string
		metadata *Metadata
		want     bool
	}{
		{
			name:     "exact match",
			metadata: &Metadata{UID: 1000},
			want:     true,
		},
		{
			name:     "no match",
			metadata: &Metadata{UID: 1001},
			want:     false,
		},
		{
			name:     "zero uid no match",
			metadata: &Metadata{UID: 0},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Match(tt.metadata); got != tt.want {
				t.Errorf("UID.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUIDInvalidPayload(t *testing.T) {
	tests := []string{
		"invalid",
		"-1",
		"99999999999",
	}

	for _, payload := range tests {
		t.Run(payload, func(t *testing.T) {
			_, err := NewUID(payload, ActionDirect)
			if err == nil {
				t.Errorf("NewUID(%s) should return error", payload)
			}
		})
	}
}

func TestDSCP(t *testing.T) {
	rule, err := NewDSCP("46", ActionDirect)
	if err != nil {
		t.Fatalf("NewDSCP failed: %v", err)
	}

	tests := []struct {
		name     string
		metadata *Metadata
		want     bool
	}{
		{
			name:     "exact match",
			metadata: &Metadata{DSCP: 46},
			want:     true,
		},
		{
			name:     "no match",
			metadata: &Metadata{DSCP: 0},
			want:     false,
		},
		{
			name:     "different value",
			metadata: &Metadata{DSCP: 10},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Match(tt.metadata); got != tt.want {
				t.Errorf("DSCP.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDSCPInvalidPayload(t *testing.T) {
	tests := []string{
		"invalid",
		"-1",
		"256",
		"1000",
	}

	for _, payload := range tests {
		t.Run(payload, func(t *testing.T) {
			_, err := NewDSCP(payload, ActionDirect)
			if err == nil {
				t.Errorf("NewDSCP(%s) should return error", payload)
			}
		})
	}
}

// ===== 正则表达式扩展测试 =====

func TestProcessRegex(t *testing.T) {
	// 进程名称正则
	nameRule, err := NewProcessRegex("^chrome.*", ActionProxy, false)
	if err != nil {
		t.Fatalf("NewProcessRegex failed: %v", err)
	}

	tests := []struct {
		name     string
		metadata *Metadata
		want     bool
	}{
		{
			name:     "match - chrome",
			metadata: &Metadata{ProcessName: "chrome"},
			want:     true,
		},
		{
			name:     "match - chrome.exe",
			metadata: &Metadata{ProcessName: "chrome.exe"},
			want:     true,
		},
		{
			name:     "no match - firefox",
			metadata: &Metadata{ProcessName: "firefox"},
			want:     false,
		},
		{
			name:     "empty process name",
			metadata: &Metadata{ProcessName: ""},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nameRule.Match(tt.metadata); got != tt.want {
				t.Errorf("ProcessRegex.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessPathRegex(t *testing.T) {
	// 进程路径正则
	pathRule, err := NewProcessRegex("/usr/bin/.*", ActionDirect, true)
	if err != nil {
		t.Fatalf("NewProcessRegex failed: %v", err)
	}

	tests := []struct {
		name     string
		metadata *Metadata
		want     bool
	}{
		{
			name:     "match - /usr/bin/curl",
			metadata: &Metadata{ProcessPath: "/usr/bin/curl"},
			want:     true,
		},
		{
			name:     "match - /usr/bin/wget",
			metadata: &Metadata{ProcessPath: "/usr/bin/wget"},
			want:     true,
		},
		{
			name:     "no match - /usr/local/bin/app",
			metadata: &Metadata{ProcessPath: "/usr/local/bin/app"},
			want:     false,
		},
		{
			name:     "empty process path",
			metadata: &Metadata{ProcessPath: ""},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pathRule.Match(tt.metadata); got != tt.want {
				t.Errorf("ProcessRegex.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessRegexInvalid(t *testing.T) {
	_, err := NewProcessRegex("[invalid", ActionProxy, false)
	if err == nil {
		t.Error("NewProcessRegex with invalid regex should return error")
	}
}

func TestIPSuffix(t *testing.T) {
	// 匹配 8.8.8.0/24 段中后缀为 1 的 IP (即 8.8.8.1)
	rule, err := NewIPSuffix("8.8.8.0/24,1", ActionDirect)
	if err != nil {
		t.Fatalf("NewIPSuffix failed: %v", err)
	}

	tests := []struct {
		name     string
		metadata *Metadata
		want     bool
	}{
		{
			name:     "exact match - 8.8.8.1",
			metadata: &Metadata{DstIP: netip.MustParseAddr("8.8.8.1")},
			want:     true,
		},
		{
			name:     "no match - 8.8.8.2",
			metadata: &Metadata{DstIP: netip.MustParseAddr("8.8.8.2")},
			want:     false,
		},
		{
			name:     "no match - different subnet",
			metadata: &Metadata{DstIP: netip.MustParseAddr("8.8.9.1")},
			want:     false,
		},
		{
			name:     "invalid IP",
			metadata: &Metadata{DstIP: netip.Addr{}},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Match(tt.metadata); got != tt.want {
				t.Errorf("IPSuffix.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPSuffixLargerRange(t *testing.T) {
	// 匹配 192.168.0.0/16 段中后缀为 256 的 IP (即 192.168.1.0)
	rule, err := NewIPSuffix("192.168.0.0/16,256", ActionDirect)
	if err != nil {
		t.Fatalf("NewIPSuffix failed: %v", err)
	}

	tests := []struct {
		name     string
		metadata *Metadata
		want     bool
	}{
		{
			name:     "match - 192.168.1.0",
			metadata: &Metadata{DstIP: netip.MustParseAddr("192.168.1.0")},
			want:     true,
		},
		{
			name:     "no match - 192.168.0.0",
			metadata: &Metadata{DstIP: netip.MustParseAddr("192.168.0.0")},
			want:     false,
		},
		{
			name:     "no match - 192.168.1.1",
			metadata: &Metadata{DstIP: netip.MustParseAddr("192.168.1.1")},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Match(tt.metadata); got != tt.want {
				t.Errorf("IPSuffix.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPSuffixNoResolve(t *testing.T) {
	rule, err := NewIPSuffix("8.8.8.0/24,1", ActionDirect)
	if err != nil {
		t.Fatalf("NewIPSuffix failed: %v", err)
	}
	rule.NoResolve(true)

	// 有域名但没有 IP，应该不匹配
	metadata := &Metadata{
		Host:  "example.com",
		DstIP: netip.Addr{},
	}

	if got := rule.Match(metadata); got {
		t.Errorf("IPSuffix.Match() with no-resolve should return false for domain without IP")
	}
}

func TestIPSuffixInvalidPayload(t *testing.T) {
	tests := []struct {
		name    string
		payload string
	}{
		{"missing suffix", "8.8.8.0/24"},
		{"invalid IP", "invalid/24,1"},
		{"invalid suffix", "8.8.8.0/24,invalid"},
		{"suffix too large", "8.8.8.0/24,999999"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewIPSuffix(tt.payload, ActionDirect)
			if err == nil {
				t.Errorf("NewIPSuffix(%s) should return error", tt.payload)
			}
		})
	}
}

// ===== GEOSITE 规则测试 =====

func TestGeoSite(t *testing.T) {
	// 设置测试用的 GeoSite 提供者
	oldProvider := GetGeositeProvider()
	defer SetGeositeProvider(oldProvider)

	SetGeositeProvider(func(domain string) []string {
		// 简单的测试映射
		switch domain {
		case "google.com", "www.google.com":
			return []string{"google", "cn"}
		case "facebook.com":
			return []string{"facebook", "social"}
		default:
			return nil
		}
	})

	rule := NewGeoSite("google", ActionProxy)

	tests := []struct {
		name     string
		metadata *Metadata
		want     bool
	}{
		{
			name:     "match - google.com",
			metadata: &Metadata{Domain: "google.com"},
			want:     true,
		},
		{
			name:     "match - www.google.com",
			metadata: &Metadata{Domain: "www.google.com"},
			want:     true,
		},
		{
			name:     "match - case insensitive",
			metadata: &Metadata{Host: "Google.COM"},
			want:     true,
		},
		{
			name:     "no match - facebook.com",
			metadata: &Metadata{Domain: "facebook.com"},
			want:     false,
		},
		{
			name:     "no match - unknown domain",
			metadata: &Metadata{Domain: "example.com"},
			want:     false,
		},
		{
			name:     "no match - empty domain",
			metadata: &Metadata{},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Match(tt.metadata); got != tt.want {
				t.Errorf("GeoSite.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainMatcher(t *testing.T) {
	matcher := NewDomainMatcher()

	// 添加测试数据
	matcher.AddDomain("google.com", "search", "tech")
	matcher.AddSuffix("github.com", "dev", "tech")
	matcher.AddKeyword("ads", "advertising")

	tests := []struct {
		name   string
		domain string
		want   []string
	}{
		{
			name:   "exact domain match",
			domain: "google.com",
			want:   []string{"search", "tech"},
		},
		{
			name:   "suffix match",
			domain: "www.github.com",
			want:   []string{"dev", "tech"},
		},
		{
			name:   "keyword match",
			domain: "doubleclick.ads.google.com",
			want:   []string{"advertising"},
		},
		{
			name:   "no match",
			domain: "example.com",
			want:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matcher.Match(tt.domain)
			if len(got) != len(tt.want) {
				t.Errorf("DomainMatcher.Match() returned %d categories, want %d", len(got), len(tt.want))
				return
			}
			// 检查是否包含所有预期的类别
			gotMap := make(map[string]bool)
			for _, cat := range got {
				gotMap[cat] = true
			}
			for _, cat := range tt.want {
				if !gotMap[cat] {
					t.Errorf("DomainMatcher.Match() missing category %s", cat)
				}
			}
		})
	}
}

// ===== 规则引擎测试 =====

func TestEngine(t *testing.T) {
	engine := NewEngine()

	// 添加各种规则
	engine.AddRule(NewDomain("google.com", ActionProxy))
	engine.AddRule(NewDomainSuffix("github.com", ActionDirect))
	cidr, _ := NewIPCIDR("192.168.0.0/16", ActionDirect, false)
	engine.AddRule(cidr)
	engine.AddRule(NewMatch(ActionProxy))

	tests := []struct {
		name       string
		metadata   *Metadata
		wantAction ActionType
		wantMatch  bool
	}{
		{
			name:       "domain exact match",
			metadata:   &Metadata{Domain: "google.com"},
			wantAction: ActionProxy,
			wantMatch:  true,
		},
		{
			name:       "domain suffix match",
			metadata:   &Metadata{Domain: "www.github.com"},
			wantAction: ActionDirect,
			wantMatch:  true,
		},
		{
			name:       "ip cidr match",
			metadata:   &Metadata{DstIP: netip.MustParseAddr("192.168.1.1")},
			wantAction: ActionDirect,
			wantMatch:  true,
		},
		{
			name:       "match all",
			metadata:   &Metadata{Host: "example.com"},
			wantAction: ActionProxy,
			wantMatch:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule, ok := engine.Match(tt.metadata)
			if ok != tt.wantMatch {
				t.Errorf("Engine.Match() ok = %v, want %v", ok, tt.wantMatch)
			}
			if ok && rule.Action() != tt.wantAction {
				t.Errorf("Engine.Match() action = %v, want %v", rule.Action(), tt.wantAction)
			}
		})
	}
}

func TestSuffixTree(t *testing.T) {
	tree := NewSuffixTree()

	rule1 := NewDomainSuffix("google.com", ActionProxy)
	rule2 := NewDomainSuffix("github.com", ActionDirect)

	tree.Add("google.com", rule1)
	tree.Add("github.com", rule2)

	tests := []struct {
		name       string
		domain     string
		wantRules  int
		wantAction ActionType
	}{
		{
			name:       "match google.com",
			domain:     "google.com",
			wantRules:  1,
			wantAction: ActionProxy,
		},
		{
			name:       "match www.google.com",
			domain:     "www.google.com",
			wantRules:  1,
			wantAction: ActionProxy,
		},
		{
			name:       "match github.com",
			domain:     "github.com",
			wantRules:  1,
			wantAction: ActionDirect,
		},
		{
			name:      "no match example.com",
			domain:    "example.com",
			wantRules: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules := tree.Match(tt.domain)
			if len(rules) != tt.wantRules {
				t.Errorf("SuffixTree.Match() returned %d rules, want %d", len(rules), tt.wantRules)
			}
			if len(rules) > 0 && rules[0].Action() != tt.wantAction {
				t.Errorf("SuffixTree.Match() action = %v, want %v", rules[0].Action(), tt.wantAction)
			}
		})
	}
}

func TestIPTrie(t *testing.T) {
	trie := NewIPTrie()

	rule1, _ := NewIPCIDR("192.168.0.0/16", ActionDirect, false)
	rule2, _ := NewIPCIDR("10.0.0.0/8", ActionProxy, false)

	trie.Add(netip.MustParsePrefix("192.168.0.0/16"), rule1)
	trie.Add(netip.MustParsePrefix("10.0.0.0/8"), rule2)

	tests := []struct {
		name       string
		ip         string
		wantRules  int
		wantAction ActionType
	}{
		{
			name:       "match 192.168.1.1",
			ip:         "192.168.1.1",
			wantRules:  1,
			wantAction: ActionDirect,
		},
		{
			name:       "match 10.0.0.1",
			ip:         "10.0.0.1",
			wantRules:  1,
			wantAction: ActionProxy,
		},
		{
			name:      "no match 8.8.8.8",
			ip:        "8.8.8.8",
			wantRules: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules := trie.Match(netip.MustParseAddr(tt.ip))
			if len(rules) != tt.wantRules {
				t.Errorf("IPTrie.Match() returned %d rules, want %d", len(rules), tt.wantRules)
			}
			if len(rules) > 0 && rules[0].Action() != tt.wantAction {
				t.Errorf("IPTrie.Match() action = %v, want %v", rules[0].Action(), tt.wantAction)
			}
		})
	}
}

func TestEngineClear(t *testing.T) {
	engine := NewEngine()

	// 添加规则
	engine.AddRule(NewDomain("google.com", ActionProxy))
	engine.AddRule(NewDomainSuffix("github.com", ActionDirect))

	// 验证规则存在
	if len(engine.Rules()) != 2 {
		t.Errorf("Engine.Rules() returned %d rules, want 2", len(engine.Rules()))
	}

	// 清空规则
	engine.Clear()

	// 验证规则已清空
	if len(engine.Rules()) != 0 {
		t.Errorf("Engine.Rules() after Clear() returned %d rules, want 0", len(engine.Rules()))
	}

	// 验证匹配失败
	metadata := &Metadata{Domain: "google.com"}
	if _, ok := engine.Match(metadata); ok {
		t.Error("Engine.Match() after Clear() should return false")
	}
}

// ===== 解析器扩展测试 =====

// ===== GeoIP 规则测试 =====

func TestGEOIP(t *testing.T) {
	// 设置测试用的 GeoIP 提供者
	oldProvider := GetGeoIPProvider()
	defer SetGeoIPProvider(oldProvider)

	SetGeoIPProvider(func(ip string) *prism.GeoIP {
		// 简单的测试映射
		if ip == "1.1.1.1" {
			return &prism.GeoIP{
				CountryCode: "US",
				Country:     "United States",
			}
		}
		if ip == "114.114.114.114" {
			return &prism.GeoIP{
				CountryCode: "CN",
				Country:     "China",
			}
		}
		return nil
	})

	rule := NewGEOIP("CN", ActionDirect)

	tests := []struct {
		name     string
		metadata *Metadata
		want     bool
	}{
		{
			name: "match CN IP",
			metadata: &Metadata{
				DstIP: netip.MustParseAddr("114.114.114.114"),
			},
			want: true,
		},
		{
			name: "no match US IP",
			metadata: &Metadata{
				DstIP: netip.MustParseAddr("1.1.1.1"),
			},
			want: false,
		},
		{
			name: "match with cached GeoIP",
			metadata: &Metadata{
				DstIP: netip.MustParseAddr("8.8.8.8"),
				DstGeoIP: &GeoIPInfo{
					CountryCode: "CN",
				},
			},
			want: true,
		},
		{
			name:     "no match - invalid IP",
			metadata: &Metadata{},
			want:     false,
		},
		{
			name: "no match - unknown IP",
			metadata: &Metadata{
				DstIP: netip.MustParseAddr("192.168.1.1"),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Match(tt.metadata); got != tt.want {
				t.Errorf("GEOIP.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGEOIPNoResolve(t *testing.T) {
	oldProvider := GetGeoIPProvider()
	defer SetGeoIPProvider(oldProvider)

	SetGeoIPProvider(func(ip string) *prism.GeoIP {
		return &prism.GeoIP{CountryCode: "CN"}
	})

	rule := NewGEOIP("CN", ActionDirect)
	rule.NoResolve(true)

	// 有域名但没有 IP，应该不匹配
	metadata := &Metadata{
		Host:  "example.com",
		DstIP: netip.Addr{},
	}

	if got := rule.Match(metadata); got {
		t.Errorf("GEOIP.Match() with no-resolve should return false for domain without IP")
	}

	// 有 IP 应该匹配
	metadata.DstIP = netip.MustParseAddr("1.1.1.1")
	if got := rule.Match(metadata); !got {
		t.Errorf("GEOIP.Match() with no-resolve should return true for IP")
	}
}

func TestGeoIPProvider(t *testing.T) {
	// 保存旧的提供者
	oldProvider := GetGeoIPProvider()
	defer SetGeoIPProvider(oldProvider)

	// 设置新的提供者
	testProvider := func(ip string) *prism.GeoIP {
		return &prism.GeoIP{CountryCode: "TEST"}
	}

	SetGeoIPProvider(testProvider)

	// 验证提供者已设置
	provider := GetGeoIPProvider()
	if provider == nil {
		t.Fatal("GetGeoIPProvider() returned nil")
	}

	// 测试查询
	result := provider("1.1.1.1")
	if result == nil || result.CountryCode != "TEST" {
		t.Errorf("Provider returned unexpected result: %v", result)
	}
}

// ===== IP-ASN 规则测试 =====

func TestIPASN(t *testing.T) {
	// 设置测试用的 GeoIP 提供者
	oldProvider := GetGeoIPProvider()
	defer SetGeoIPProvider(oldProvider)

	SetGeoIPProvider(func(ip string) *prism.GeoIP {
		if ip == "1.1.1.1" {
			return &prism.GeoIP{
				ASN:    13335,
				ASName: "Cloudflare",
			}
		}
		if ip == "8.8.8.8" {
			return &prism.GeoIP{
				ASN:    15169,
				ASName: "Google",
			}
		}
		return nil
	})

	rule, err := NewIPASN("13335", ActionProxy)
	if err != nil {
		t.Fatalf("NewIPASN failed: %v", err)
	}

	tests := []struct {
		name     string
		metadata *Metadata
		want     bool
	}{
		{
			name: "match ASN 13335",
			metadata: &Metadata{
				DstIP: netip.MustParseAddr("1.1.1.1"),
			},
			want: true,
		},
		{
			name: "no match ASN 15169",
			metadata: &Metadata{
				DstIP: netip.MustParseAddr("8.8.8.8"),
			},
			want: false,
		},
		{
			name: "match with cached GeoIP",
			metadata: &Metadata{
				DstIP: netip.MustParseAddr("2.2.2.2"),
				DstGeoIP: &GeoIPInfo{
					ASN: 13335,
				},
			},
			want: true,
		},
		{
			name:     "no match - invalid IP",
			metadata: &Metadata{},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Match(tt.metadata); got != tt.want {
				t.Errorf("IPASN.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPASNNoResolve(t *testing.T) {
	oldProvider := GetGeoIPProvider()
	defer SetGeoIPProvider(oldProvider)

	SetGeoIPProvider(func(ip string) *prism.GeoIP {
		return &prism.GeoIP{ASN: 13335}
	})

	rule, _ := NewIPASN("13335", ActionProxy)
	rule.NoResolve(true)

	// 有域名但没有 IP，应该不匹配
	metadata := &Metadata{
		Host:  "example.com",
		DstIP: netip.Addr{},
	}

	if got := rule.Match(metadata); got {
		t.Errorf("IPASN.Match() with no-resolve should return false for domain without IP")
	}
}

func TestIPASNInvalid(t *testing.T) {
	tests := []string{
		"invalid",
		"not-a-number",
	}

	for _, payload := range tests {
		t.Run(payload, func(t *testing.T) {
			_, err := NewIPASN(payload, ActionProxy)
			if err == nil {
				t.Errorf("NewIPASN(%s) should return error", payload)
			}
		})
	}
}

func TestIPASNNegative(t *testing.T) {
	// -1 虽然可以解析，但逻辑上是无效的 ASN
	rule, err := NewIPASN("-1", ActionProxy)
	if err != nil {
		t.Fatalf("NewIPASN failed: %v", err)
	}

	// 但不应该匹配任何有效的 ASN
	metadata := &Metadata{
		DstGeoIP: &GeoIPInfo{ASN: -1},
	}

	if !rule.Match(metadata) {
		t.Error("IPASN should match -1")
	}
}

// ===== DomainRegex 规则测试 =====

func TestDomainRegex(t *testing.T) {
	rule, err := NewDomainRegex("^.*\\.cn$", ActionDirect)
	if err != nil {
		t.Fatalf("NewDomainRegex failed: %v", err)
	}

	tests := []struct {
		name     string
		metadata *Metadata
		want     bool
	}{
		{
			name:     "match .cn domain",
			metadata: &Metadata{Host: "example.cn"},
			want:     true,
		},
		{
			name:     "match subdomain .cn",
			metadata: &Metadata{Host: "www.example.cn"},
			want:     true,
		},
		{
			name:     "no match .com domain",
			metadata: &Metadata{Host: "example.com"},
			want:     false,
		},
		{
			name:     "no match empty",
			metadata: &Metadata{},
			want:     false,
		},
		{
			name:     "use Domain field",
			metadata: &Metadata{Domain: "test.cn"},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Match(tt.metadata); got != tt.want {
				t.Errorf("DomainRegex.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainRegexInvalid(t *testing.T) {
	_, err := NewDomainRegex("[invalid", ActionDirect)
	if err == nil {
		t.Error("NewDomainRegex with invalid regex should return error")
	}
}

// ===== Parser 功能测试 =====

func TestParseRules(t *testing.T) {
	ruleStrs := []string{
		"DOMAIN,google.com,PROXY",
		"# comment",
		"",
		"DOMAIN-SUFFIX,example.com,DIRECT",
		"IP-CIDR,192.168.0.0/16,DIRECT",
	}

	rules, err := ParseRules(ruleStrs)
	if err != nil {
		t.Fatalf("ParseRules failed: %v", err)
	}

	// 应该有 3 条有效规则（跳过注释和空行）
	if len(rules) != 3 {
		t.Errorf("ParseRules() returned %d rules, want 3", len(rules))
	}

	// 验证规则类型
	if rules[0].Type() != TypeDomain {
		t.Errorf("First rule type = %v, want %v", rules[0].Type(), TypeDomain)
	}
	if rules[1].Type() != TypeDomainSuffix {
		t.Errorf("Second rule type = %v, want %v", rules[1].Type(), TypeDomainSuffix)
	}
	if rules[2].Type() != TypeIPCIDR {
		t.Errorf("Third rule type = %v, want %v", rules[2].Type(), TypeIPCIDR)
	}
}

func TestParseRulesError(t *testing.T) {
	ruleStrs := []string{
		"DOMAIN,google.com,PROXY",
		"INVALID,test,PROXY", // 无效规则
	}

	_, err := ParseRules(ruleStrs)
	if err == nil {
		t.Error("ParseRules should return error for invalid rule")
	}
}

func TestLoadRulesFromFile(t *testing.T) {
	// 创建临时文件
	tmpfile := "/tmp/test_rules.txt"
	content := `# Test rules file
DOMAIN,google.com,PROXY
DOMAIN-SUFFIX,example.com,DIRECT

# More rules
IP-CIDR,192.168.0.0/16,DIRECT
MATCH,PROXY
`
	err := os.WriteFile(tmpfile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(tmpfile)

	rules, err := LoadRulesFromFile(tmpfile)
	if err != nil {
		t.Fatalf("LoadRulesFromFile failed: %v", err)
	}

	// 应该有 4 条规则
	if len(rules) != 4 {
		t.Errorf("LoadRulesFromFile() returned %d rules, want 4", len(rules))
	}
}

func TestLoadRulesFromFileNotExist(t *testing.T) {
	_, err := LoadRulesFromFile("/tmp/nonexistent_file_12345.txt")
	if err == nil {
		t.Error("LoadRulesFromFile should return error for non-existent file")
	}
}

func TestLoadRulesFromFileInvalidRule(t *testing.T) {
	tmpfile := "/tmp/test_invalid_rules.txt"
	content := `DOMAIN,google.com,PROXY
INVALID,test,PROXY
`
	err := os.WriteFile(tmpfile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(tmpfile)

	_, err = LoadRulesFromFile(tmpfile)
	if err == nil {
		t.Error("LoadRulesFromFile should return error for invalid rule")
	}
}

func TestParseRuleExtended(t *testing.T) {
	tests := []struct {
		name    string
		ruleStr string
		wantErr bool
		check   func(Rule) bool
	}{
		{
			name:    "IN-TYPE rule",
			ruleStr: "IN-TYPE,HTTP,PROXY",
			wantErr: false,
			check: func(r Rule) bool {
				return r.Type() == TypeInType && r.Action() == ActionProxy
			},
		},
		{
			name:    "IN-NAME rule",
			ruleStr: "IN-NAME,proxy-in,DIRECT",
			wantErr: false,
			check: func(r Rule) bool {
				return r.Type() == TypeInName && r.Action() == ActionDirect
			},
		},
		{
			name:    "IN-USER rule",
			ruleStr: "IN-USER,alice,PROXY",
			wantErr: false,
			check: func(r Rule) bool {
				return r.Type() == TypeInUser && r.Action() == ActionProxy
			},
		},
		{
			name:    "NETWORK rule",
			ruleStr: "NETWORK,TCP,DIRECT",
			wantErr: false,
			check: func(r Rule) bool {
				return r.Type() == TypeNetwork && r.Action() == ActionDirect
			},
		},
		{
			name:    "UID rule",
			ruleStr: "UID,1000,DIRECT",
			wantErr: false,
			check: func(r Rule) bool {
				return r.Type() == TypeUID && r.Action() == ActionDirect
			},
		},
		{
			name:    "DSCP rule",
			ruleStr: "DSCP,46,DIRECT",
			wantErr: false,
			check: func(r Rule) bool {
				return r.Type() == TypeDSCP && r.Action() == ActionDirect
			},
		},
		{
			name:    "PROCESS-NAME-REGEX rule",
			ruleStr: "PROCESS-NAME-REGEX,^chrome.*,PROXY",
			wantErr: false,
			check: func(r Rule) bool {
				return r.Type() == TypeProcessNameRegex && r.Action() == ActionProxy
			},
		},
		{
			name:    "PROCESS-PATH-REGEX rule",
			ruleStr: "PROCESS-PATH-REGEX,/usr/bin/.*,DIRECT",
			wantErr: false,
			check: func(r Rule) bool {
				return r.Type() == TypeProcessPathRegex && r.Action() == ActionDirect
			},
		},
		{
			name:    "IP-SUFFIX rule",
			ruleStr: "IP-SUFFIX,8.8.8.0/24,1,DIRECT",
			wantErr: false,
			check: func(r Rule) bool {
				return r.Type() == TypeIPSuffix && r.Action() == ActionDirect
			},
		},
		{
			name:    "GEOSITE rule",
			ruleStr: "GEOSITE,google,PROXY",
			wantErr: false,
			check: func(r Rule) bool {
				return r.Type() == TypeGeoSite && r.Action() == ActionProxy
			},
		},
		{
			name:    "invalid regex",
			ruleStr: "PROCESS-NAME-REGEX,[invalid,PROXY",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule, err := ParseRule(tt.ruleStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil && !tt.check(rule) {
				t.Errorf("ParseRule() returned unexpected rule: %v", rule)
			}
		})
	}
}

// ===== IPCIDR IPv6 和 no-resolve 测试 =====

func TestIPCIDR6(t *testing.T) {
	rule, err := NewIPCIDR("2001:db8::/32", ActionDirect, true)
	if err != nil {
		t.Fatalf("NewIPCIDR failed: %v", err)
	}

	tests := []struct {
		name     string
		metadata *Metadata
		want     bool
	}{
		{
			name: "in range IPv6",
			metadata: &Metadata{
				DstIP: netip.MustParseAddr("2001:db8::1"),
			},
			want: true,
		},
		{
			name: "out of range IPv6",
			metadata: &Metadata{
				DstIP: netip.MustParseAddr("2001:db9::1"),
			},
			want: false,
		},
		{
			name: "IPv4 not match",
			metadata: &Metadata{
				DstIP: netip.MustParseAddr("192.168.1.1"),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Match(tt.metadata); got != tt.want {
				t.Errorf("IPCIDR6.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPCIDRNoResolve(t *testing.T) {
	rule, _ := NewIPCIDR("192.168.0.0/16", ActionDirect, false)
	rule.NoResolve(true)

	// 有域名但没有 IP，应该不匹配
	metadata := &Metadata{
		Host:  "example.com",
		DstIP: netip.Addr{},
	}

	if got := rule.Match(metadata); got {
		t.Errorf("IPCIDR.Match() with no-resolve should return false for domain without IP")
	}

	// 有 IP 应该匹配
	metadata.DstIP = netip.MustParseAddr("192.168.1.1")
	if got := rule.Match(metadata); !got {
		t.Errorf("IPCIDR.Match() with no-resolve should return true for IP")
	}
}

func TestIPCIDRInvalid(t *testing.T) {
	tests := []string{
		"invalid",
		"256.256.256.256/24",
		"192.168.0.0/33",
	}

	for _, payload := range tests {
		t.Run(payload, func(t *testing.T) {
			_, err := NewIPCIDR(payload, ActionDirect, false)
			if err == nil {
				t.Errorf("NewIPCIDR(%s) should return error", payload)
			}
		})
	}
}

// ===== Process 路径匹配测试 =====

func TestProcessPath(t *testing.T) {
	rule := NewProcess("/usr/bin/curl", ActionDirect, true)

	tests := []struct {
		name     string
		metadata *Metadata
		want     bool
	}{
		{
			name:     "exact match path",
			metadata: &Metadata{ProcessPath: "/usr/bin/curl"},
			want:     true,
		},
		{
			name:     "contains path",
			metadata: &Metadata{ProcessPath: "/usr/bin/curl --help"},
			want:     true,
		},
		{
			name:     "no match different path",
			metadata: &Metadata{ProcessPath: "/usr/bin/wget"},
			want:     false,
		},
		{
			name:     "empty path",
			metadata: &Metadata{ProcessPath: ""},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Match(tt.metadata); got != tt.want {
				t.Errorf("Process.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ===== Port 边界情况测试 =====

func TestPortSinglePort(t *testing.T) {
	rule, err := NewPort("443", ActionDirect, true)
	if err != nil {
		t.Fatalf("NewPort failed: %v", err)
	}

	if !rule.Match(&Metadata{DstPort: 443}) {
		t.Error("Port should match single port 443")
	}
	if rule.Match(&Metadata{DstPort: 80}) {
		t.Error("Port should not match port 80")
	}
}

func TestPortSourcePort(t *testing.T) {
	rule, err := NewPort("7890", ActionReject, false)
	if err != nil {
		t.Fatalf("NewPort failed: %v", err)
	}

	tests := []struct {
		name     string
		metadata *Metadata
		want     bool
	}{
		{
			name:     "match source port",
			metadata: &Metadata{SrcPort: 7890},
			want:     true,
		},
		{
			name:     "no match different source port",
			metadata: &Metadata{SrcPort: 8080},
			want:     false,
		},
		{
			name:     "dst port not matched",
			metadata: &Metadata{DstPort: 7890},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rule.Match(tt.metadata); got != tt.want {
				t.Errorf("Port.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPortInvalidPayload(t *testing.T) {
	tests := []string{
		"invalid",
		"70000",
		"-1",
		"100-50", // 倒序范围
	}

	for _, payload := range tests {
		t.Run(payload, func(t *testing.T) {
			_, err := NewPort(payload, ActionDirect, true)
			if err == nil {
				t.Errorf("NewPort(%s) should return error", payload)
			}
		})
	}
}

// ===== base.String() 测试 =====

func TestBaseString(t *testing.T) {
	rule := NewDomain("google.com", ActionProxy)

	str := rule.String()
	if str == "" {
		t.Error("Rule.String() should not return empty string")
	}

	// 验证字符串包含关键信息
	if !contains(str, "DOMAIN") {
		t.Error("Rule.String() should contain rule type")
	}
	if !contains(str, "google.com") {
		t.Error("Rule.String() should contain payload")
	}
	if !contains(str, "PROXY") {
		t.Error("Rule.String() should contain action")
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ===== Engine.AddRules 测试 =====

func TestEngineAddRules(t *testing.T) {
	engine := NewEngine()

	cidr, _ := NewIPCIDR("10.0.0.0/8", ActionDirect, false)
	rules := []Rule{
		NewDomain("google.com", ActionProxy),
		NewDomainSuffix("example.com", ActionDirect),
		cidr,
	}

	engine.AddRules(rules)

	allRules := engine.Rules()
	if len(allRules) != 3 {
		t.Errorf("Engine.Rules() returned %d rules, want 3", len(allRules))
	}
}

// ===== Domain/DomainKeyword Domain 字段测试 =====

func TestDomainWithDomainField(t *testing.T) {
	rule := NewDomain("google.com", ActionProxy)

	// 测试使用 Domain 字段
	metadata := &Metadata{
		Domain: "google.com",
	}

	if !rule.Match(metadata) {
		t.Error("Domain should match using Domain field")
	}
}

func TestDomainKeywordWithDomainField(t *testing.T) {
	rule := NewDomainKeyword("google", ActionProxy)

	// 测试使用 Domain 字段
	metadata := &Metadata{
		Domain: "www.google.com",
	}

	if !rule.Match(metadata) {
		t.Error("DomainKeyword should match using Domain field")
	}
}

// ===== MatchFirst 完整测试 =====

func TestMatchFirstNoMatch(t *testing.T) {
	rules := []Rule{
		NewDomain("google.com", ActionProxy),
		NewDomainSuffix("example.com", ActionDirect),
	}

	metadata := &Metadata{Host: "other.com"}

	if _, ok := MatchFirst(rules, metadata); ok {
		t.Error("MatchFirst should return false when no rule matches")
	}
}
