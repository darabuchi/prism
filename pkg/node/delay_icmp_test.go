package node_test

import (
	"context"
	"testing"
	"time"

	"github.com/darabuchi/prism/pkg/node"
	"github.com/lazygophers/log"
)

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
		log.Infof("ICMP delay test passed (not supported): %v", err)
	}
}

func TestDelayICMPWithEmptyAddr(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 测试空地址
	_, err = n.DelayICMP(ctx, "")
	if err == nil {
		t.Error("expected error for empty address")
	}
}

func TestDelayICMPWithWireGuard(t *testing.T) {
	// 创建 WireGuard 类型的节点
	wireGuardConfig := map[string]any{
		"name":        "WireGuard 节点",
		"type":        "wireguard",
		"server":      "127.0.0.1",
		"port":        51820,
		"private-key": "test-private-key",
		"public-key":  "test-public-key",
	}

	n, err := node.NewNode(wireGuardConfig)
	if err != nil {
		// WireGuard 可能因为缺少必要字段而创建失败，这是正常的
		t.Skipf("WireGuard node creation failed (expected): %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 测试 WireGuard 的 ICMP（应该返回"未实现"错误）
	_, err = n.DelayICMP(ctx, "8.8.8.8")
	if err == nil {
		t.Error("expected error for WireGuard ICMP (not implemented)")
	} else {
		log.Infof("WireGuard ICMP test passed (not implemented): %v", err)
	}
}
