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

	log.Infof("TCP 延迟: %v", delay)
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

		log.Infof("TCP 延迟 (%s): %v", addr, delay)
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
}
