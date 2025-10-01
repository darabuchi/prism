package prism

import (
	"testing"
)

func TestGeoIP_Location(t *testing.T) {
	tests := []struct {
		name     string
		geoip    GeoIP
		expected string
	}{
		{
			name: "完整地理位置信息",
			geoip: GeoIP{
				City:    "Beijing",
				Region:  "Beijing",
				Country: "China",
			},
			expected: "Beijing, Beijing, China",
		},
		{
			name: "只有国家",
			geoip: GeoIP{
				Country: "United States",
			},
			expected: "United States",
		},
		{
			name: "只有国家代码",
			geoip: GeoIP{
				CountryCode: "US",
			},
			expected: "US",
		},
		{
			name: "城市和国家",
			geoip: GeoIP{
				City:    "Tokyo",
				Country: "Japan",
			},
			expected: "Tokyo, Japan",
		},
		{
			name:     "空信息",
			geoip:    GeoIP{},
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.geoip.Location()
			if result != tt.expected {
				t.Errorf("Location() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGeoIP_Coordinates(t *testing.T) {
	tests := []struct {
		name     string
		geoip    GeoIP
		expected string
	}{
		{
			name: "有坐标",
			geoip: GeoIP{
				Latitude:  39.9042,
				Longitude: 116.4074,
			},
			expected: "39.9042, 116.4074",
		},
		{
			name:     "零坐标",
			geoip:    GeoIP{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.geoip.Coordinates()
			if result != tt.expected {
				t.Errorf("Coordinates() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGeoIP_ASString(t *testing.T) {
	tests := []struct {
		name     string
		geoip    GeoIP
		expected string
	}{
		{
			name: "完整 ASN 信息",
			geoip: GeoIP{
				ASN:    15169,
				ASName: "Google LLC",
			},
			expected: "AS15169 Google LLC",
		},
		{
			name: "只有 ASN",
			geoip: GeoIP{
				ASN: 15169,
			},
			expected: "AS15169",
		},
		{
			name:     "没有 ASN",
			geoip:    GeoIP{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.geoip.ASString()
			if result != tt.expected {
				t.Errorf("ASString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGeoIP_IsIPv4(t *testing.T) {
	tests := []struct {
		name     string
		geoip    GeoIP
		expected bool
	}{
		{
			name: "IPv4 通过版本字段",
			geoip: GeoIP{
				IP:        "8.8.8.8",
				IPVersion: 4,
			},
			expected: true,
		},
		{
			name: "IPv4 通过 IP 解析",
			geoip: GeoIP{
				IP: "8.8.8.8",
			},
			expected: true,
		},
		{
			name: "IPv6 通过版本字段",
			geoip: GeoIP{
				IP:        "2001:4860:4860::8888",
				IPVersion: 6,
			},
			expected: false,
		},
		{
			name: "IPv6 通过 IP 解析",
			geoip: GeoIP{
				IP: "2001:4860:4860::8888",
			},
			expected: false,
		},
		{
			name:     "空 IP",
			geoip:    GeoIP{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.geoip.IsIPv4()
			if result != tt.expected {
				t.Errorf("IsIPv4() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGeoIP_IsIPv6(t *testing.T) {
	tests := []struct {
		name     string
		geoip    GeoIP
		expected bool
	}{
		{
			name: "IPv6 通过版本字段",
			geoip: GeoIP{
				IP:        "2001:4860:4860::8888",
				IPVersion: 6,
			},
			expected: true,
		},
		{
			name: "IPv6 通过 IP 解析",
			geoip: GeoIP{
				IP: "2001:4860:4860::8888",
			},
			expected: true,
		},
		{
			name: "IPv4 通过版本字段",
			geoip: GeoIP{
				IP:        "8.8.8.8",
				IPVersion: 4,
			},
			expected: false,
		},
		{
			name: "IPv4 通过 IP 解析",
			geoip: GeoIP{
				IP: "8.8.8.8",
			},
			expected: false,
		},
		{
			name:     "空 IP",
			geoip:    GeoIP{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.geoip.IsIPv6()
			if result != tt.expected {
				t.Errorf("IsIPv6() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGeoIP_HasLocation(t *testing.T) {
	tests := []struct {
		name     string
		geoip    GeoIP
		expected bool
	}{
		{
			name: "有国家",
			geoip: GeoIP{
				Country: "China",
			},
			expected: true,
		},
		{
			name: "有国家代码",
			geoip: GeoIP{
				CountryCode: "CN",
			},
			expected: true,
		},
		{
			name: "有城市",
			geoip: GeoIP{
				City: "Beijing",
			},
			expected: true,
		},
		{
			name: "有地区",
			geoip: GeoIP{
				Region: "Beijing",
			},
			expected: true,
		},
		{
			name:     "没有位置信息",
			geoip:    GeoIP{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.geoip.HasLocation()
			if result != tt.expected {
				t.Errorf("HasLocation() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGeoIP_HasASN(t *testing.T) {
	tests := []struct {
		name     string
		geoip    GeoIP
		expected bool
	}{
		{
			name: "有 ASN",
			geoip: GeoIP{
				ASN: 15169,
			},
			expected: true,
		},
		{
			name:     "没有 ASN",
			geoip:    GeoIP{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.geoip.HasASN()
			if result != tt.expected {
				t.Errorf("HasASN() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGeoIP_String(t *testing.T) {
	tests := []struct {
		name     string
		geoip    GeoIP
		expected string
	}{
		{
			name: "完整信息",
			geoip: GeoIP{
				IP:      "8.8.8.8",
				Country: "United States",
				City:    "Mountain View",
				ISP:     "Google LLC",
				ASN:     15169,
				ASName:  "Google LLC",
			},
			expected: "8.8.8.8 | Mountain View, United States | Google LLC | AS15169 Google LLC",
		},
		{
			name: "只有 IP",
			geoip: GeoIP{
				IP: "8.8.8.8",
			},
			expected: "8.8.8.8",
		},
		{
			name:     "空信息",
			geoip:    GeoIP{},
			expected: "Empty GeoIP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.geoip.String()
			if result != tt.expected {
				t.Errorf("String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHealthState_String(t *testing.T) {
	tests := []struct {
		name     string
		state    HealthState
		expected string
	}{
		{
			name:     "开路状态",
			state:    HealthStateOpen,
			expected: "open",
		},
		{
			name:     "半开状态",
			state:    HealthStateHalfOpen,
			expected: "half_open",
		},
		{
			name:     "闭路状态",
			state:    HealthStateClosed,
			expected: "closed",
		},
		{
			name:     "未知状态",
			state:    HealthState(999),
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.state.String()
			if result != tt.expected {
				t.Errorf("String() = %v, want %v", result, tt.expected)
			}
		})
	}
}
