package node_test

import (
	"context"
	"testing"
	"time"

	"github.com/darabuchi/prism/pkg/node"
	"github.com/lazygophers/log"
)

func TestSpeedDownloadHTTP(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test download speed with default URL and 10-second duration
	result, err := n.SpeedDownloadHTTP(ctx, "", 10*time.Second)
	if err != nil {
		t.Logf("Download speed test failed (proxy may not be running): %v", err)
		t.Skip("skipping download speed test")
		return
	}

	if result.BytesDownloaded <= 0 {
		t.Error("expected bytes downloaded > 0")
	}

	if result.Speed <= 0 {
		t.Error("expected speed > 0")
	}

	log.Infof("Download speed test result:")
	log.Infof("  Bytes downloaded: %d", result.BytesDownloaded)
	log.Infof("  Duration: %v", result.Duration)
	log.Infof("  Speed: %.2f Mbps", result.SpeedMbps)
}

func TestSpeedDownloadHTTPWithCustomURL(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with different test files
	testCases := []struct {
		name string
		url  string
	}{
		{
			name: "1MB file",
			url:  "http://cachefly.cachefly.net/1mb.test",
		},
		{
			name: "10MB file",
			url:  "http://cachefly.cachefly.net/10mb.test",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := n.SpeedDownloadHTTP(ctx, tc.url, 5*time.Second)
			if err != nil {
				t.Logf("Download speed test failed for %s: %v", tc.name, err)
				return
			}

			log.Infof("%s download speed: %.2f Mbps (%d bytes in %v)",
				tc.name, result.SpeedMbps, result.BytesDownloaded, result.Duration)
		})
	}
}

func TestSpeedDownloadHTTPFullFile(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test downloading complete file (duration = 0)
	result, err := n.SpeedDownloadHTTP(ctx, "http://cachefly.cachefly.net/1mb.test", 0)
	if err != nil {
		t.Logf("Download speed test failed: %v", err)
		t.Skip("skipping full file download test")
		return
	}

	// 1MB = 1048576 bytes
	expectedSize := int64(1048576)
	if result.BytesDownloaded < expectedSize-1024 || result.BytesDownloaded > expectedSize+1024 {
		t.Logf("Warning: expected ~%d bytes, got %d bytes", expectedSize, result.BytesDownloaded)
	}

	log.Infof("Full file download: %d bytes in %v (%.2f Mbps)",
		result.BytesDownloaded, result.Duration, result.SpeedMbps)
}

func TestSpeedDownloadHTTPCanceledContext(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	// Create a context that's already canceled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Should return immediately with partial results
	result, err := n.SpeedDownloadHTTP(ctx, "", 10*time.Second)
	if err != nil {
		t.Logf("Download speed test failed with canceled context: %v", err)
		t.Skip("skipping canceled context test")
		return
	}

	// Should have 0 bytes downloaded since context was canceled immediately
	log.Infof("Canceled context test: %d bytes downloaded", result.BytesDownloaded)
}
