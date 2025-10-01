package node_test

import (
	"context"
	"testing"
	"time"

	"github.com/darabuchi/prism/pkg/node"
	"github.com/lazygophers/log"
)

func TestSpeedTestDownloadHTTP(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test HTTP download with default URL
	result, err := n.SpeedTestDownload(ctx, "", 10*time.Second)
	if err != nil {
		t.Logf("Speed test download failed: %v", err)
		t.Skip("skipping HTTP download test")
		return
	}

	if result.TestType != "http" {
		t.Errorf("expected test type 'http', got '%s'", result.TestType)
	}

	if result.SpeedMbps <= 0 {
		t.Error("expected speed > 0 Mbps")
	}

	log.Infof("HTTP download speed test:")
	log.Infof("  Test type: %s", result.TestType)
	log.Infof("  Speed: %.2f Mbps", result.SpeedMbps)
	log.Infof("  Bytes: %d", result.BytesTransfer)
	log.Infof("  Duration: %v", result.Duration)
}

func TestSpeedTestDownloadWithCustomHTTPURL(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with custom HTTP URL
	result, err := n.SpeedTestDownload(ctx, "http://cachefly.cachefly.net/1mb.test", 5*time.Second)
	if err != nil {
		t.Logf("Speed test download failed: %v", err)
		t.Skip("skipping custom HTTP download test")
		return
	}

	if result.TestType != "http" {
		t.Errorf("expected test type 'http', got '%s'", result.TestType)
	}

	log.Infof("Custom HTTP download: %.2f Mbps (%d bytes in %v)",
		result.SpeedMbps, result.BytesTransfer, result.Duration)
}

func TestSpeedTestDownloadOokla(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Test Ookla download with automatic server selection
	result, err := n.SpeedTestDownload(ctx, "ookla://", 0)
	if err != nil {
		t.Logf("Ookla download test failed: %v", err)
		t.Skip("skipping Ookla download test")
		return
	}

	if result.TestType != "ookla" {
		t.Errorf("expected test type 'ookla', got '%s'", result.TestType)
	}

	log.Infof("Ookla download speed test:")
	log.Infof("  Speed: %.2f Mbps", result.SpeedMbps)
	log.Infof("  Latency: %v", result.Latency)
}

func TestSpeedTestUploadHTTP(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test HTTP upload
	result, err := n.SpeedTestUpload(ctx, "", 1024*1024, 10*time.Second)
	if err != nil {
		t.Logf("Speed test upload failed: %v", err)
		t.Skip("skipping HTTP upload test")
		return
	}

	if result.TestType != "http" {
		t.Errorf("expected test type 'http', got '%s'", result.TestType)
	}

	if result.SpeedMbps <= 0 {
		t.Error("expected speed > 0 Mbps")
	}

	log.Infof("HTTP upload speed test:")
	log.Infof("  Test type: %s", result.TestType)
	log.Infof("  Speed: %.2f Mbps", result.SpeedMbps)
	log.Infof("  Bytes: %d", result.BytesTransfer)
	log.Infof("  Duration: %v", result.Duration)
}

func TestSpeedTestUploadOokla(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Test Ookla upload
	result, err := n.SpeedTestUpload(ctx, "ookla://", 0, 0)
	if err != nil {
		t.Logf("Ookla upload test failed: %v", err)
		t.Skip("skipping Ookla upload test")
		return
	}

	if result.TestType != "ookla" {
		t.Errorf("expected test type 'ookla', got '%s'", result.TestType)
	}

	log.Infof("Ookla upload speed test:")
	log.Infof("  Speed: %.2f Mbps", result.SpeedMbps)
	log.Infof("  Latency: %v", result.Latency)
}

func TestSpeedTest(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	// Test full Ookla speedtest
	download, upload, err := n.SpeedTest(ctx, "ookla://")
	if err != nil {
		t.Logf("Ookla speedtest failed: %v", err)
		t.Skip("skipping Ookla full test")
		return
	}

	if download.TestType != "ookla" {
		t.Errorf("expected download test type 'ookla', got '%s'", download.TestType)
	}

	if upload.TestType != "ookla" {
		t.Errorf("expected upload test type 'ookla', got '%s'", upload.TestType)
	}

	log.Infof("Ookla full speedtest:")
	log.Infof("  Server: %s (%s) [%s]", download.ServerName, download.ServerCountry, download.ServerID)
	log.Infof("  Latency: %v", download.Latency)
	log.Infof("  Download: %.2f Mbps", download.SpeedMbps)
	log.Infof("  Upload: %.2f Mbps", upload.SpeedMbps)
}

func TestSpeedTestWithSpeedtestURL(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Test using speedtest:// prefix
	result, err := n.SpeedTestDownload(ctx, "speedtest://", 0)
	if err != nil {
		t.Logf("Speedtest download failed: %v", err)
		t.Skip("skipping speedtest:// download test")
		return
	}

	if result.TestType != "ookla" {
		t.Errorf("expected test type 'ookla', got '%s'", result.TestType)
	}

	log.Infof("Speedtest download (speedtest:// prefix): %.2f Mbps", result.SpeedMbps)
}

func TestSpeedTestUnsupportedURL(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test with unsupported URL (SpeedTest only supports ookla://)
	_, _, err = n.SpeedTest(ctx, "http://example.com")
	if err != node.ErrUnsupportedSpeedTestURL {
		t.Errorf("expected ErrUnsupportedSpeedTestURL, got %v", err)
	}

	log.Infof("Unsupported URL test passed: got expected error")
}
