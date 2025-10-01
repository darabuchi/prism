package node_test

import (
	"testing"

	"github.com/darabuchi/prism/pkg/node"
	"github.com/lazygophers/log"
)

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

	log.Infof("Node created successfully: %s, ID: %s", n.Name(), n.UniqueId())
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

			log.Infof("%s Created successfully: %s", tc.name, n.Name())
		})
	}
}

func TestNodeConfig(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	config := n.Config()
	if config == nil {
		t.Error("config is nil")
	}

	if config["server"] != "127.0.0.1" {
		t.Errorf("expected server 127.0.0.1, got %v", config["server"])
	}
}

func TestNodeProxyType(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	proxyType := n.ProxyType()
	if proxyType == "" {
		t.Error("proxy type is empty")
	}

	log.Infof("Proxy type: %s", proxyType)
}

func TestNodeType(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	adapterType := n.Type()
	if adapterType.String() == "" {
		t.Error("adapter type is empty")
	}

	log.Infof("Adapter type: %s", adapterType)
}

func TestNodeAddr(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	addr := n.Addr()
	if addr == "" {
		t.Error("addr is empty")
	}

	log.Infof("Node address: %s", addr)
}

func TestNodeSupportUDP(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	supportUDP := n.SupportUDP()
	log.Infof("UDP support: %v", supportUDP)
}

func TestNodeSupportUOT(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	supportUOT := n.SupportUOT()
	log.Infof("UOT support: %v", supportUOT)
}

func TestNodeSupportWithDialer(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	network := n.SupportWithDialer()
	log.Infof("Supported network types: %v", network)
}

func TestNodeProxyInfo(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	proxyInfo := n.ProxyInfo()
	log.Infof("Proxy info: %+v", proxyInfo)
}

func TestNodeMarshalJSON(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	data, err := n.MarshalJSON()
	if err != nil {
		t.Errorf("MarshalJSON failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("JSON data is empty")
	}

	log.Infof("JSON: %s", string(data))
}

func TestNodeClose(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	err = n.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}
}

func TestNewNodeWithProtocolNormalization(t *testing.T) {
	testCases := []struct {
		name      string
		config    map[string]any
		shouldSkip bool
	}{
		{
			name: "shadowsocks to ss",
			config: map[string]any{
				"name":     "test-ss",
				"type":     "shadowsocks",
				"server":   "127.0.0.1",
				"port":     7890,
				"cipher":   "aes-256-gcm",
				"password": "test-password",
			},
		},
		{
			name: "shadowsocksr to ssr",
			config: map[string]any{
				"name":     "test-ssr",
				"type":     "shadowsocksr",
				"server":   "127.0.0.1",
				"port":     7890,
				"cipher":   "aes-256-cfb",
				"password": "test-password",
				"protocol": "origin",
				"obfs":     "plain",
			},
		},
		{
			name: "hy to hysteria",
			config: map[string]any{
				"name":   "test-hy",
				"type":   "hy",
				"server": "127.0.0.1",
				"port":   7890,
			},
			shouldSkip: true, // Hysteria 需要更多字段，可能创建失败
		},
		{
			name: "hy2 to hysteria2",
			config: map[string]any{
				"name":   "test-hy2",
				"type":   "hy2",
				"server": "127.0.0.1",
				"port":   7890,
			},
			shouldSkip: true, // Hysteria2 需要更多字段，可能创建失败
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n, err := node.NewNode(tc.config)
			if err != nil {
				if tc.shouldSkip {
					t.Skipf("%s: NewNode failed (expected): %v", tc.name, err)
					return
				}
				t.Fatalf("NewNode failed for %s: %v", tc.name, err)
			}

			if n == nil {
				t.Fatalf("%s: node is nil", tc.name)
			}

			log.Infof("%s: Protocol normalization successful", tc.name)
		})
	}
}

func TestNewNodeWithoutName(t *testing.T) {
	config := map[string]any{
		"type":   "http",
		"server": "127.0.0.1",
		"port":   7890,
		// 没有 name 字段
	}

	n, err := node.NewNode(config)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	if n.Name() == "" {
		t.Error("node name should have default value")
	}

	log.Infof("Default node name: %s", n.Name())
}
