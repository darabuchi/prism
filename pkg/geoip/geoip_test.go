package geoip

import (
	"net"
	"testing"

	"github.com/darabuchi/prism"
)

// mockReader 用于测试的模拟读取器
type mockReader struct {
	lookupFunc func(ip net.IP) (*prism.GeoIP, error)
}

func (m *mockReader) Lookup(ip net.IP) (*prism.GeoIP, error) {
	if m.lookupFunc != nil {
		return m.lookupFunc(ip)
	}
	return nil, nil
}

func (m *mockReader) LookupString(ipStr string) (*prism.GeoIP, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, nil
	}
	return m.Lookup(ip)
}

func (m *mockReader) Close() error {
	return nil
}

func TestLookupIP(t *testing.T) {
	ip := net.ParseIP("8.8.8.8")

	tests := []struct {
		name    string
		readers []Reader
		wantErr bool
	}{
		{
			name:    "无读取器",
			readers: []Reader{},
			wantErr: true,
		},
		{
			name: "单个读取器成功",
			readers: []Reader{
				&mockReader{
					lookupFunc: func(ip net.IP) (*prism.GeoIP, error) {
						return &prism.GeoIP{
							IP:      ip.String(),
							Country: "United States",
						}, nil
					},
				},
			},
			wantErr: false,
		},
		{
			name: "多个读取器合并",
			readers: []Reader{
				&mockReader{
					lookupFunc: func(ip net.IP) (*prism.GeoIP, error) {
						return &prism.GeoIP{
							IP:      ip.String(),
							Country: "United States",
						}, nil
					},
				},
				&mockReader{
					lookupFunc: func(ip net.IP) (*prism.GeoIP, error) {
						return &prism.GeoIP{
							City: "Mountain View",
							ASN:  15169,
						}, nil
					},
				},
			},
			wantErr: false,
		},
		{
			name: "忽略失败的读取器",
			readers: []Reader{
				&mockReader{
					lookupFunc: func(ip net.IP) (*prism.GeoIP, error) {
						return nil, nil // 失败
					},
				},
				&mockReader{
					lookupFunc: func(ip net.IP) (*prism.GeoIP, error) {
						return &prism.GeoIP{
							IP:      ip.String(),
							Country: "United States",
						}, nil
					},
				},
			},
			wantErr: false,
		},
		{
			name: "所有读取器失败",
			readers: []Reader{
				&mockReader{
					lookupFunc: func(ip net.IP) (*prism.GeoIP, error) {
						return nil, nil
					},
				},
				&mockReader{
					lookupFunc: func(ip net.IP) (*prism.GeoIP, error) {
						return nil, nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "忽略 nil 读取器",
			readers: []Reader{
				nil,
				&mockReader{
					lookupFunc: func(ip net.IP) (*prism.GeoIP, error) {
						return &prism.GeoIP{
							IP:      ip.String(),
							Country: "United States",
						}, nil
					},
				},
				nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := LookupIP(ip, tt.readers...)

			if (err != nil) != tt.wantErr {
				t.Errorf("LookupIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result == nil {
					t.Error("expected non-nil result")
					return
				}
				if result.IP != ip.String() {
					t.Errorf("IP = %v, want %v", result.IP, ip.String())
				}
			}
		})
	}
}

func TestLookupString(t *testing.T) {
	tests := []struct {
		name    string
		ipStr   string
		readers []Reader
		wantErr bool
	}{
		{
			name:  "有效 IPv4",
			ipStr: "8.8.8.8",
			readers: []Reader{
				&mockReader{
					lookupFunc: func(ip net.IP) (*prism.GeoIP, error) {
						return &prism.GeoIP{
							IP:      ip.String(),
							Country: "United States",
						}, nil
					},
				},
			},
			wantErr: false,
		},
		{
			name:  "有效 IPv6",
			ipStr: "2001:4860:4860::8888",
			readers: []Reader{
				&mockReader{
					lookupFunc: func(ip net.IP) (*prism.GeoIP, error) {
						return &prism.GeoIP{
							IP:      ip.String(),
							Country: "United States",
						}, nil
					},
				},
			},
			wantErr: false,
		},
		{
			name:  "无效 IP",
			ipStr: "invalid-ip",
			readers: []Reader{
				&mockReader{
					lookupFunc: func(ip net.IP) (*prism.GeoIP, error) {
						return &prism.GeoIP{}, nil
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := LookupString(tt.ipStr, tt.readers...)

			if (err != nil) != tt.wantErr {
				t.Errorf("LookupString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result == nil {
				t.Error("expected non-nil result")
			}
		})
	}
}

func TestLookupIPMerge(t *testing.T) {
	ip := net.ParseIP("8.8.8.8")

	reader1 := &mockReader{
		lookupFunc: func(ip net.IP) (*prism.GeoIP, error) {
			return &prism.GeoIP{
				IP:      ip.String(),
				Country: "United States",
				City:    "Mountain View",
			}, nil
		},
	}

	reader2 := &mockReader{
		lookupFunc: func(ip net.IP) (*prism.GeoIP, error) {
			return &prism.GeoIP{
				CountryCode: "US",
				Region:      "California",
				ASN:         15169,
			}, nil
		},
	}

	reader3 := &mockReader{
		lookupFunc: func(ip net.IP) (*prism.GeoIP, error) {
			return &prism.GeoIP{
				ASName: "Google LLC",
				ISP:    "Google LLC",
			}, nil
		},
	}

	result, err := LookupIP(ip, reader1, reader2, reader3)
	if err != nil {
		t.Fatalf("LookupIP() error = %v", err)
	}

	// 验证合并结果
	if result.IP != "8.8.8.8" {
		t.Errorf("IP = %v, want %v", result.IP, "8.8.8.8")
	}
	if result.Country != "United States" {
		t.Errorf("Country = %v, want %v", result.Country, "United States")
	}
	if result.City != "Mountain View" {
		t.Errorf("City = %v, want %v", result.City, "Mountain View")
	}
	if result.CountryCode != "US" {
		t.Errorf("CountryCode = %v, want %v", result.CountryCode, "US")
	}
	if result.Region != "California" {
		t.Errorf("Region = %v, want %v", result.Region, "California")
	}
	if result.ASN != 15169 {
		t.Errorf("ASN = %v, want %v", result.ASN, 15169)
	}
	if result.ASName != "Google LLC" {
		t.Errorf("ASName = %v, want %v", result.ASName, "Google LLC")
	}
	if result.ISP != "Google LLC" {
		t.Errorf("ISP = %v, want %v", result.ISP, "Google LLC")
	}
	if result.AS != "AS15169 Google LLC" {
		t.Errorf("AS = %v, want %v", result.AS, "AS15169 Google LLC")
	}
}
