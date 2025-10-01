package node_test

import (
	"context"
	"testing"
	"time"

	"github.com/darabuchi/prism/pkg/node"
	"github.com/lazygophers/log"
)

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

	log.Infof("TCP delay: %v", delay)
}

func TestDelayTCPWithDifferentPorts(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 测试不同的端口
	testAddrs := []string{
		"www.google.com:80",
		"www.google.com:443",
	}

	for _, addr := range testAddrs {
		delay, err := n.DelayTCP(ctx, addr)
		if err != nil {
			t.Logf("TCP 延迟测试失败 (%s): %v", addr, err)
			continue
		}

		if delay <= 0 {
			t.Errorf("expected positive delay for %s, got %v", addr, delay)
		}

		log.Infof("TCP delay (%s): %v", addr, delay)
	}
}

func TestDelayTCPWithInvalidAddr(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 测试无效地址
	_, err = n.DelayTCP(ctx, "invalid-address")
	if err == nil {
		t.Error("expected error for invalid address")
	}

	// 测试空地址
	_, err = n.DelayTCP(ctx, "")
	if err == nil {
		t.Error("expected error for empty address")
	}

	// 测试无效端口（port为0）
	_, err = n.DelayTCP(ctx, "example.com:0")
	if err == nil {
		t.Error("expected error for invalid port 0")
	}
}

func TestDelayTCPWithUnreachableHost(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 测试不可达的主机
	// 注意：代理可能会接受连接请求，所以这个测试可能不会失败
	_, err = n.DelayTCP(ctx, "192.0.2.1:80") // 192.0.2.0/24 is TEST-NET-1, non-routable
	if err != nil {
		log.Infof("Unreachable host test passed (connection failed): %v", err)
	} else {
		log.Infof("Unreachable host test (proxy accepted connection request)")
	}
}
