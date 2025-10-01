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
		log.Infof("ICMP 延迟测试符合预期（不支持）: %v", err)
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
