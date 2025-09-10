package testutil

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

// MockHTTPServer 创建模拟 HTTP 服务器
func MockHTTPServer(responses map[string]string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		if response, ok := responses[url]; ok {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(response))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
		}
	}))
}

// MockSubscriptionData 模拟订阅数据
const MockSubscriptionData = `vmess://eyJhZGQiOiJ0ZXN0MS5leGFtcGxlLmNvbSIsImFpZCI6IjAiLCJob3N0IjoiIiwiaWQiOiJ1dWlkLTEiLCJuZXQiOiJ0Y3AiLCJwYXRoIjoiIiwicG9ydCI6IjQ0MyIsInBzIjoidGVzdC1ub2RlLTEiLCJzY3kiOiJhdXRvIiwic25pIjoiIiwidGxzIjoiIiwidHlwZSI6Im5vbmUiLCJ2IjoiMiJ9
vmess://eyJhZGQiOiJ0ZXN0Mi5leGFtcGxlLmNvbSIsImFpZCI6IjAiLCJob3N0IjoiIiwiaWQiOiJ1dWlkLTIiLCJuZXQiOiJ0Y3AiLCJwYXRoIjoiIiwicG9ydCI6IjQ0MyIsInBzIjoidGVzdC1ub2RlLTIiLCJzY3kiOiJhdXRvIiwic25pIjoiIiwidGxzIjoiIiwidHlwZSI6Im5vbmUiLCJ2IjoiMiJ9`

// MockYAMLSubscriptionData 模拟 YAML 格式订阅数据
const MockYAMLSubscriptionData = `proxies:
  - name: "test-node-1"
    type: vmess
    server: test1.example.com
    port: 443
    uuid: uuid-1
    alterId: 0
    cipher: auto
    network: tcp
    
  - name: "test-node-2"
    type: vmess
    server: test2.example.com
    port: 443
    uuid: uuid-2
    alterId: 0
    cipher: auto
    network: tcp`

// ValidateJSONResponse 验证 JSON 响应格式
func ValidateJSONResponse(body string, expectedKeys ...string) bool {
	for _, key := range expectedKeys {
		if !strings.Contains(body, `"`+key+`"`) {
			return false
		}
	}
	return true
}

// AssertContains 断言字符串包含子字符串
func AssertContains(t interface {
	Errorf(format string, args ...interface{})
}, haystack, needle string) {
	if !strings.Contains(haystack, needle) {
		t.Errorf("Expected %q to contain %q", haystack, needle)
	}
}

// AssertNotContains 断言字符串不包含子字符串
func AssertNotContains(t interface {
	Errorf(format string, args ...interface{})
}, haystack, needle string) {
	if strings.Contains(haystack, needle) {
		t.Errorf("Expected %q to not contain %q", haystack, needle)
	}
}
