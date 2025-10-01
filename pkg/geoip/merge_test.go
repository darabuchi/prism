package geoip

import (
	"testing"

	"github.com/darabuchi/prism"
)

func TestMerge(t *testing.T) {
	tests := []struct {
		name     string
		geoips   []*prism.GeoIP
		expected *prism.GeoIP
	}{
		{
			name:   "空输入",
			geoips: []*prism.GeoIP{},
			expected: &prism.GeoIP{},
		},
		{
			name: "单个输入",
			geoips: []*prism.GeoIP{
				{
					IP:          "8.8.8.8",
					Country:     "United States",
					CountryCode: "US",
					ASN:         15169,
					ASName:      "Google LLC",
				},
			},
			expected: &prism.GeoIP{
				IP:          "8.8.8.8",
				IPVersion:   4,
				Country:     "United States",
				CountryCode: "US",
				ASN:         15169,
				ASName:      "Google LLC",
				AS:          "AS15169 Google LLC",
			},
		},
		{
			name: "合并互补信息",
			geoips: []*prism.GeoIP{
				{
					IP:      "8.8.8.8",
					Country: "United States",
					City:    "Mountain View",
				},
				{
					CountryCode: "US",
					Region:      "California",
					ASN:         15169,
				},
				{
					ASName: "Google LLC",
					ISP:    "Google LLC",
				},
			},
			expected: &prism.GeoIP{
				IP:          "8.8.8.8",
				IPVersion:   4,
				Country:     "United States",
				CountryCode: "US",
				City:        "Mountain View",
				Region:      "California",
				ASN:         15169,
				ASName:      "Google LLC",
				AS:          "AS15169 Google LLC",
				ISP:         "Google LLC",
			},
		},
		{
			name: "优先使用第一个非空值",
			geoips: []*prism.GeoIP{
				{
					IP:      "8.8.8.8",
					Country: "United States",
					City:    "Mountain View",
				},
				{
					Country: "US", // 应该被忽略，因为第一个已经有了
					City:    "Los Angeles", // 应该被忽略
					Region:  "California", // 第一个没有，应该被使用
				},
			},
			expected: &prism.GeoIP{
				IP:        "8.8.8.8",
				IPVersion: 4,
				Country:   "United States",
				City:      "Mountain View",
				Region:    "California",
			},
		},
		{
			name: "布尔值 OR 逻辑",
			geoips: []*prism.GeoIP{
				{
					IP:    "8.8.8.8",
					Proxy: false,
					Hosting: true,
				},
				{
					Proxy: true,
					Hosting: false,
				},
			},
			expected: &prism.GeoIP{
				IP:        "8.8.8.8",
				IPVersion: 4,
				Proxy:     true,
				Hosting:   true,
			},
		},
		{
			name: "忽略 nil 输入",
			geoips: []*prism.GeoIP{
				nil,
				{
					IP:      "8.8.8.8",
					Country: "United States",
				},
				nil,
				{
					City: "Mountain View",
				},
				nil,
			},
			expected: &prism.GeoIP{
				IP:        "8.8.8.8",
				IPVersion: 4,
				Country:   "United States",
				City:      "Mountain View",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Merge(tt.geoips...)

			if result.IP != tt.expected.IP {
				t.Errorf("IP = %v, want %v", result.IP, tt.expected.IP)
			}
			if result.IPVersion != tt.expected.IPVersion {
				t.Errorf("IPVersion = %v, want %v", result.IPVersion, tt.expected.IPVersion)
			}
			if result.Country != tt.expected.Country {
				t.Errorf("Country = %v, want %v", result.Country, tt.expected.Country)
			}
			if result.CountryCode != tt.expected.CountryCode {
				t.Errorf("CountryCode = %v, want %v", result.CountryCode, tt.expected.CountryCode)
			}
			if result.City != tt.expected.City {
				t.Errorf("City = %v, want %v", result.City, tt.expected.City)
			}
			if result.Region != tt.expected.Region {
				t.Errorf("Region = %v, want %v", result.Region, tt.expected.Region)
			}
			if result.ASN != tt.expected.ASN {
				t.Errorf("ASN = %v, want %v", result.ASN, tt.expected.ASN)
			}
			if result.ASName != tt.expected.ASName {
				t.Errorf("ASName = %v, want %v", result.ASName, tt.expected.ASName)
			}
			if result.AS != tt.expected.AS {
				t.Errorf("AS = %v, want %v", result.AS, tt.expected.AS)
			}
			if result.ISP != tt.expected.ISP {
				t.Errorf("ISP = %v, want %v", result.ISP, tt.expected.ISP)
			}
			if result.Proxy != tt.expected.Proxy {
				t.Errorf("Proxy = %v, want %v", result.Proxy, tt.expected.Proxy)
			}
			if result.Hosting != tt.expected.Hosting {
				t.Errorf("Hosting = %v, want %v", result.Hosting, tt.expected.Hosting)
			}
		})
	}
}

func TestAppend(t *testing.T) {
	tests := []struct {
		name     string
		dest     *prism.GeoIP
		source   *prism.GeoIP
		expected *prism.GeoIP
	}{
		{
			name:     "nil dest",
			dest:     nil,
			source:   &prism.GeoIP{IP: "8.8.8.8"},
			expected: nil,
		},
		{
			name:     "nil source",
			dest:     &prism.GeoIP{IP: "8.8.8.8"},
			source:   nil,
			expected: &prism.GeoIP{IP: "8.8.8.8", IPVersion: 4},
		},
		{
			name: "追加缺失字段",
			dest: &prism.GeoIP{
				IP:      "8.8.8.8",
				Country: "United States",
			},
			source: &prism.GeoIP{
				City: "Mountain View",
				ASN:  15169,
			},
			expected: &prism.GeoIP{
				IP:        "8.8.8.8",
				IPVersion: 4,
				Country:   "United States",
				City:      "Mountain View",
				ASN:       15169,
			},
		},
		{
			name: "不覆盖已有字段",
			dest: &prism.GeoIP{
				IP:      "8.8.8.8",
				Country: "United States",
				City:    "Mountain View",
			},
			source: &prism.GeoIP{
				Country: "US",
				City:    "Los Angeles",
				Region:  "California",
			},
			expected: &prism.GeoIP{
				IP:        "8.8.8.8",
				IPVersion: 4,
				Country:   "United States",
				City:      "Mountain View",
				Region:    "California",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Append(tt.dest, tt.source)

			if tt.expected == nil {
				if tt.dest != nil {
					t.Error("expected dest to remain nil")
				}
				return
			}

			if tt.dest.IP != tt.expected.IP {
				t.Errorf("IP = %v, want %v", tt.dest.IP, tt.expected.IP)
			}
			if tt.dest.Country != tt.expected.Country {
				t.Errorf("Country = %v, want %v", tt.dest.Country, tt.expected.Country)
			}
			if tt.dest.City != tt.expected.City {
				t.Errorf("City = %v, want %v", tt.dest.City, tt.expected.City)
			}
			if tt.dest.Region != tt.expected.Region {
				t.Errorf("Region = %v, want %v", tt.dest.Region, tt.expected.Region)
			}
			if tt.dest.ASN != tt.expected.ASN {
				t.Errorf("ASN = %v, want %v", tt.dest.ASN, tt.expected.ASN)
			}
		})
	}
}

func TestClone(t *testing.T) {
	original := &prism.GeoIP{
		IP:          "8.8.8.8",
		IPVersion:   4,
		Country:     "United States",
		CountryCode: "US",
		City:        "Mountain View",
		ASN:         15169,
		ASName:      "Google LLC",
		Proxy:       true,
	}

	cloned := Clone(original)

	// 验证值相同
	if cloned.IP != original.IP {
		t.Errorf("IP = %v, want %v", cloned.IP, original.IP)
	}
	if cloned.Country != original.Country {
		t.Errorf("Country = %v, want %v", cloned.Country, original.Country)
	}
	if cloned.ASN != original.ASN {
		t.Errorf("ASN = %v, want %v", cloned.ASN, original.ASN)
	}
	if cloned.Proxy != original.Proxy {
		t.Errorf("Proxy = %v, want %v", cloned.Proxy, original.Proxy)
	}

	// 验证是不同的对象
	if cloned == original {
		t.Error("cloned should be a different object")
	}

	// 修改克隆对象不应影响原对象
	cloned.City = "Los Angeles"
	if original.City == "Los Angeles" {
		t.Error("modifying clone should not affect original")
	}
}

func TestCloneNil(t *testing.T) {
	cloned := Clone(nil)
	if cloned == nil {
		t.Error("Clone(nil) should return empty GeoIP, not nil")
	}
	if cloned.IP != "" {
		t.Errorf("Clone(nil) should return empty GeoIP, got IP = %v", cloned.IP)
	}
}
