package node_test

import (
	"context"
	"testing"
	"time"

	"github.com/darabuchi/prism/pkg/node"
	"github.com/lazygophers/log"
)

// 测试用的代理配置
var testProxyConfig = map[string]any{
	"name":   "测试节点",
	"type":   "http",
	"server": "127.0.0.1",
	"port":   7890,
}

func TestNewNode(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	if n.Name() == "" {
		t.Error("node name is empty")
	}

	if n.Server() != "127.0.0.1" {
		t.Errorf("expected server 127.0.0.1, got %s", n.Server())
	}

	if n.Port() != 7890 {
		t.Errorf("expected port 7890, got %d", n.Port())
	}

	if n.UniqueId() == "" {
		t.Error("unique id is empty")
	}

	log.Infof("节点创建成功: %s, ID: %s", n.Name(), n.UniqueId())
}

func TestDelayHTTP(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 测试 HTTP 延迟
	delay, err := n.DelayHTTP(ctx, "http://www.gstatic.com/generate_204")
	if err != nil {
		t.Logf("HTTP 延迟测试失败（可能是代理未运行）: %v", err)
		t.Skip("skipping HTTP delay test")
		return
	}

	if delay <= 0 {
		t.Errorf("expected positive delay, got %v", delay)
	}

	log.Infof("HTTP 延迟: %v", delay)
}

func TestDelayTCP(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 测试 TCP 延迟
	delay, err := n.DelayTCP(ctx, "www.google.com:80")
	if err != nil {
		t.Logf("TCP 延迟测试失败（可能是代理未运行）: %v", err)
		t.Skip("skipping TCP delay test")
		return
	}

	if delay <= 0 {
		t.Errorf("expected positive delay, got %v", delay)
	}

	log.Infof("TCP 延迟: %v", delay)
}

func TestDelayICMP(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 测试 ICMP 延迟（预期失败，因为 HTTP 代理不支持 ICMP）
	_, err = n.DelayICMP(ctx, "8.8.8.8")
	if err == nil {
		t.Error("expected ICMP to be unsupported for HTTP proxy")
	} else {
		log.Infof("ICMP 延迟测试符合预期（不支持）: %v", err)
	}
}

func TestGenerateId(t *testing.T) {
	config1 := map[string]any{
		"type":   "vmess",
		"server": "example.com",
		"port":   443,
		"uuid":   "test-uuid",
	}

	config2 := map[string]any{
		"type":   "vmess",
		"server": "example.com",
		"port":   443,
		"uuid":   "test-uuid",
		"name":   "不同的名字", // name 不应该影响 ID
	}

	id1 := node.GenerateId(config1)
	id2 := node.GenerateId(config2)

	if id1 != id2 {
		t.Error("相同配置（除 name 外）应该生成相同的 ID")
	}

	if len(id1) != 64 {
		t.Errorf("expected 64 character ID, got %d", len(id1))
	}

	log.Infof("生成的 ID: %s", id1)
}

func TestNodeWithDifferentProtocols(t *testing.T) {
	testCases := []struct {
		name   string
		config map[string]any
	}{
		{
			name: "HTTP Proxy",
			config: map[string]any{
				"name":   "HTTP 节点",
				"type":   "http",
				"server": "127.0.0.1",
				"port":   7890,
			},
		},
		{
			name: "SOCKS5 Proxy",
			config: map[string]any{
				"name":   "SOCKS5 节点",
				"type":   "socks5",
				"server": "127.0.0.1",
				"port":   7891,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n, err := node.NewNode(tc.config)
			if err != nil {
				t.Fatalf("NewNode failed for %s: %v", tc.name, err)
			}

			if n.Name() == "" {
				t.Errorf("%s: node name is empty", tc.name)
			}

			if n.UniqueId() == "" {
				t.Errorf("%s: unique id is empty", tc.name)
			}

			log.Infof("%s 创建成功: %s", tc.name, n.Name())
		})
	}
}
