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

	log.Infof("节点创建成功: %s, ID: %s", n.Name(), n.UniqueId())
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

	log.Infof("代理类型: %s", proxyType)
}
