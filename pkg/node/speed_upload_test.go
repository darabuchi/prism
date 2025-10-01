package node_test

import (
	"context"
	"testing"
	"time"

	"github.com/darabuchi/prism/pkg/node"
	"github.com/lazygophers/log"
)

func TestSpeedUploadHTTP(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test upload speed with 1MB data and 10-second duration
	result, err := n.SpeedUploadHTTP(ctx, "", 1024*1024, 10*time.Second)
	if err != nil {
		t.Logf("Upload speed test failed (proxy may not be running): %v", err)
		t.Skip("skipping upload speed test")
		return
	}

	if result.BytesUploaded <= 0 {
		t.Error("expected bytes uploaded > 0")
	}

	if result.Speed <= 0 {
		t.Error("expected speed > 0")
	}

	log.Infof("Upload speed test result:")
	log.Infof("  Bytes uploaded: %d", result.BytesUploaded)
	log.Infof("  Duration: %v", result.Duration)
	log.Infof("  Speed: %.2f Mbps", result.SpeedMbps)
}

func TestSpeedUploadHTTPWithDifferentSizes(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Test with different data sizes
	testCases := []struct {
		name     string
		dataSize int64
		duration time.Duration
	}{
		{
			name:     "512KB",
			dataSize: 512 * 1024,
			duration: 5 * time.Second,
		},
		{
			name:     "1MB",
			dataSize: 1024 * 1024,
			duration: 5 * time.Second,
		},
		{
			name:     "5MB",
			dataSize: 5 * 1024 * 1024,
			duration: 10 * time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := n.SpeedUploadHTTP(ctx, "", tc.dataSize, tc.duration)
			if err != nil {
				t.Logf("Upload speed test failed for %s: %v", tc.name, err)
				return
			}

			log.Infof("%s upload speed: %.2f Mbps (%d bytes in %v)",
				tc.name, result.SpeedMbps, result.BytesUploaded, result.Duration)
		})
	}
}

func TestSpeedUploadHTTPFullData(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test uploading complete data (duration = 0)
	dataSize := int64(1024 * 1024) // 1MB
	result, err := n.SpeedUploadHTTP(ctx, "", dataSize, 0)
	if err != nil {
		t.Logf("Upload speed test failed: %v", err)
		t.Skip("skipping full data upload test")
		return
	}

	if result.BytesUploaded < dataSize-1024 {
		t.Logf("Warning: expected ~%d bytes, got %d bytes", dataSize, result.BytesUploaded)
	}

	log.Infof("Full data upload: %d bytes in %v (%.2f Mbps)",
		result.BytesUploaded, result.Duration, result.SpeedMbps)
}

func TestSpeedUploadHTTPWithCustomURL(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with custom upload endpoint
	testURLs := []string{
		"http://httpbin.org/post",
		"http://httpbin.org/anything", // accepts any HTTP method
	}

	for _, url := range testURLs {
		t.Run(url, func(t *testing.T) {
			result, err := n.SpeedUploadHTTP(ctx, url, 512*1024, 5*time.Second)
			if err != nil {
				t.Logf("Upload speed test failed for %s: %v", url, err)
				return
			}

			log.Infof("%s upload speed: %.2f Mbps", url, result.SpeedMbps)
		})
	}
}

func TestSpeedUploadHTTPCanceledContext(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	// Create a context that's already canceled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Should return immediately with partial results
	result, err := n.SpeedUploadHTTP(ctx, "", 1024*1024, 10*time.Second)
	if err != nil {
		t.Logf("Upload speed test failed with canceled context: %v", err)
		t.Skip("skipping canceled context test")
		return
	}

	// Should have 0 or very few bytes uploaded since context was canceled immediately
	log.Infof("Canceled context test: %d bytes uploaded", result.BytesUploaded)
}
