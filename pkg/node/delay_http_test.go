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

	log.Infof("HTTP delay: %v", delay)
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

		log.Infof("HTTP delay (%s): %v", url, delay)
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

	log.Infof("HTTP delay (default URL): %v", delay)
}

func TestDelayHTTPWithInvalidURL(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 测试无效的 URL
	_, err = n.DelayHTTP(ctx, "http://invalid-domain-that-does-not-exist-12345.com")
	if err == nil {
		t.Error("expected error for invalid URL")
	} else {
		log.Infof("Invalid URL test passed: %v", err)
	}
}

func TestDelayHTTPWithCanceledContext(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	// 创建一个已取消的 context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// 应该立即失败
	_, err = n.DelayHTTP(ctx, "http://www.gstatic.com/generate_204")
	if err == nil {
		t.Error("expected error for canceled context")
	} else {
		log.Infof("Canceled context test passed: %v", err)
	}
}
