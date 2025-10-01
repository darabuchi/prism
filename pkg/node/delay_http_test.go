package node_test

import (
	"context"
	"testing"
	"time"

	"github.com/darabuchi/prism/pkg/node"
	"github.com/lazygophers/log"
)

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

func TestDelayHTTPWithCustomURL(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 测试自定义 URL
	testURLs := []string{
		"http://www.gstatic.com/generate_204",
		"http://cp.cloudflare.com/generate_204",
	}

	for _, url := range testURLs {
		delay, err := n.DelayHTTP(ctx, url)
		if err != nil {
			t.Logf("HTTP 延迟测试失败 (%s): %v", url, err)
			continue
		}

		if delay <= 0 {
			t.Errorf("expected positive delay for %s, got %v", url, delay)
		}

		log.Infof("HTTP 延迟 (%s): %v", url, delay)
	}
}

func TestDelayHTTPWithEmptyURL(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 测试空 URL（应使用默认值）
	delay, err := n.DelayHTTP(ctx, "")
	if err != nil {
		t.Logf("HTTP 延迟测试失败（可能是代理未运行）: %v", err)
		t.Skip("skipping HTTP delay test")
		return
	}

	if delay <= 0 {
		t.Errorf("expected positive delay, got %v", delay)
	}

	log.Infof("HTTP 延迟（默认 URL）: %v", delay)
}
