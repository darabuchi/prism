package rules

import (
	"net/netip"
	"testing"
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
