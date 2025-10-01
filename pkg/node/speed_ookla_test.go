package node_test

import (
	"context"
	"testing"
	"time"

	"github.com/darabuchi/prism/pkg/node"
	"github.com/lazygophers/log"
)

func TestSpeedTestOokla(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	// Test full Ookla speedtest (may take 30-60 seconds)
	result, err := n.SpeedTestOokla(ctx, "")
	if err != nil {
		t.Logf("Ookla speedtest failed (proxy may not be running or speedtest.net unavailable): %v", err)
		t.Skip("skipping Ookla speedtest")
		return
	}

	if result.Download <= 0 {
		t.Error("expected download speed > 0")
	}

	if result.Upload <= 0 {
		t.Error("expected upload speed > 0")
	}

	log.Infof("Ookla speedtest result:")
	log.Infof("  Server: %s (%s) [%s]", result.ServerName, result.ServerCountry, result.ServerID)
	log.Infof("  Latency: %v", result.Latency)
	log.Infof("  Jitter: %v", result.Jitter)
	log.Infof("  Download: %.2f Mbps", result.Download)
	log.Infof("  Upload: %.2f Mbps", result.Upload)
	log.Infof("  Test Duration: %v", result.TestDuration)
}

func TestSpeedTestOoklaDownload(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Test only download speed (faster than full test)
	download, latency, err := n.SpeedTestOoklaDownload(ctx, "")
	if err != nil {
		t.Logf("Ookla download test failed: %v", err)
		t.Skip("skipping Ookla download test")
		return
	}

	if download <= 0 {
		t.Logf("Warning: download speed is 0, speedtest may have failed")
	}

	if latency <= 0 {
		t.Logf("Warning: latency is 0, speedtest may have failed")
	}

	log.Infof("Ookla download test result:")
	log.Infof("  Latency: %v", latency)
	log.Infof("  Download: %.2f Mbps", download)
}

func TestSpeedTestOoklaUpload(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Test only upload speed (faster than full test)
	upload, latency, err := n.SpeedTestOoklaUpload(ctx, "")
	if err != nil {
		t.Logf("Ookla upload test failed: %v", err)
		t.Skip("skipping Ookla upload test")
		return
	}

	if upload <= 0 {
		t.Error("expected upload speed > 0")
	}

	if latency <= 0 {
		t.Error("expected latency > 0")
	}

	log.Infof("Ookla upload test result:")
	log.Infof("  Latency: %v", latency)
	log.Infof("  Upload: %.2f Mbps", upload)
}

func TestSpeedTestOoklaWithSpecificServer(t *testing.T) {
	// This test is skipped by default because it requires knowing a valid server ID
	t.Skip("skipping specific server test - requires valid server ID")

	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	// Test with a specific server ID (you need to replace this with a valid ID)
	serverID := "12345"
	result, err := n.SpeedTestOokla(ctx, serverID)
	if err != nil {
		t.Logf("Ookla speedtest with specific server failed: %v", err)
		return
	}

	if result.ServerID != serverID {
		t.Errorf("expected server ID %s, got %s", serverID, result.ServerID)
	}

	log.Infof("Ookla speedtest with server %s: Download %.2f Mbps, Upload %.2f Mbps",
		serverID, result.Download, result.Upload)
}
