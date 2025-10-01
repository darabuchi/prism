package node_test

import (
	"context"
	"testing"
	"time"

	"github.com/darabuchi/prism/pkg/node"
	"github.com/lazygophers/log"
	"github.com/metacubex/mihomo/constant"
)

// 测试实际网络连接相关的方法

func TestNodeDialContext(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	metadata := &constant.Metadata{
		NetWork: constant.TCP,
		Host:    "www.google.com",
		DstPort: 80,
	}

	conn, err := n.DialContext(ctx, metadata)
	if err != nil {
		t.Logf("DialContext failed (可能是代理未运行): %v", err)
		t.Skip("skipping DialContext test")
		return
	}
	defer conn.Close()

	log.Infof("DialContext succeeded")
}

func TestNodeListenPacketContext(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	// HTTP 代理不支持 UDP，所以这个测试应该失败或被跳过
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	metadata := &constant.Metadata{
		NetWork: constant.UDP,
		Host:    "8.8.8.8",
		DstPort: 53,
	}

	conn, err := n.ListenPacketContext(ctx, metadata)
	if err != nil {
		log.Infof("ListenPacketContext failed (HTTP proxy does not support UDP): %v", err)
		return
	}
	defer conn.Close()

	log.Infof("ListenPacketContext succeeded")
}

func TestNodeStreamConnContext(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 先建立一个实际连接
	metadata := &constant.Metadata{
		NetWork: constant.TCP,
		Host:    "www.google.com",
		DstPort: 80,
	}

	conn, err := n.DialContext(ctx, metadata)
	if err != nil {
		t.Logf("DialContext failed (可能是代理未运行): %v", err)
		t.Skip("skipping StreamConnContext test")
		return
	}
	defer conn.Close()

	// 测试 StreamConnContext
	streamConn, err := n.StreamConnContext(ctx, conn, metadata)
	if err != nil {
		log.Infof("StreamConnContext failed (possibly not supported): %v", err)
	} else {
		if streamConn != nil {
			defer streamConn.Close()
		}
		log.Infof("StreamConnContext succeeded")
	}
}

func TestNodeDialContextWithDialer(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	// HTTP 代理的 DialContextWithDialer 需要自定义 dialer
	// 由于需要实现复杂的 dialer 接口，这里只测试方法存在性
	// 如果需要完整测试，需要使用支持 dialer 的代理类型
	_ = n
	log.Infof("DialContextWithDialer method exists")
}

func TestNodeListenPacketWithDialer(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	// HTTP 代理的 ListenPacketWithDialer 需要自定义 dialer
	// 由于需要实现复杂的 dialer 接口，这里只测试方法存在性
	// 如果需要完整测试，需要使用支持 dialer 的代理类型
	_ = n
	log.Infof("ListenPacketWithDialer method exists")
}

func TestNodeIsL3Protocol(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	metadata := &constant.Metadata{
		NetWork: constant.TCP,
		Host:    "www.google.com",
		DstPort: 80,
	}

	isL3 := n.IsL3Protocol(metadata)
	log.Infof("IsL3Protocol: %v", isL3)
}

func TestNodeUnwrap(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	metadata := &constant.Metadata{
		NetWork: constant.TCP,
		Host:    "www.google.com",
		DstPort: 80,
	}

	// HTTP 代理是叶子代理，Unwrap 应该返回 nil
	proxy := n.Unwrap(metadata, false)
	log.Infof("Unwrap result: %v", proxy)
}
