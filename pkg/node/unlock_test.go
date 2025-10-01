package node_test

import (
	"context"
	"testing"
	"time"

	"github.com/darabuchi/prism/pkg/node"
	"github.com/lazygophers/log"
)

func TestUnlockNetflix(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := n.UnlockNetflix(ctx)
	if err != nil {
		t.Fatalf("UnlockNetflix failed: %v", err)
	}

	if result.Platform != "Netflix" {
		t.Errorf("expected platform 'Netflix', got '%s'", result.Platform)
	}

	log.Infof("Netflix unlock test:")
	log.Infof("  Platform: %s", result.Platform)
	log.Infof("  Status: %s", result.Status)
	log.Infof("  Region: %s", result.Region)
	log.Infof("  Message: %s", result.Message)
}

func TestUnlockDisneyPlus(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := n.UnlockDisneyPlus(ctx)
	if err != nil {
		t.Fatalf("UnlockDisneyPlus failed: %v", err)
	}

	if result.Platform != "Disney+" {
		t.Errorf("expected platform 'Disney+', got '%s'", result.Platform)
	}

	log.Infof("Disney+ unlock test:")
	log.Infof("  Platform: %s", result.Platform)
	log.Infof("  Status: %s", result.Status)
	log.Infof("  Region: %s", result.Region)
	log.Infof("  Message: %s", result.Message)
}

func TestUnlockYouTubePremium(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := n.UnlockYouTubePremium(ctx)
	if err != nil {
		t.Fatalf("UnlockYouTubePremium failed: %v", err)
	}

	if result.Platform != "YouTube Premium" {
		t.Errorf("expected platform 'YouTube Premium', got '%s'", result.Platform)
	}

	log.Infof("YouTube Premium unlock test:")
	log.Infof("  Platform: %s", result.Platform)
	log.Infof("  Status: %s", result.Status)
	log.Infof("  Region: %s", result.Region)
	log.Infof("  Message: %s", result.Message)
}

func TestUnlockHulu(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := n.UnlockHulu(ctx)
	if err != nil {
		t.Fatalf("UnlockHulu failed: %v", err)
	}

	if result.Platform != "Hulu" {
		t.Errorf("expected platform 'Hulu', got '%s'", result.Platform)
	}

	log.Infof("Hulu unlock test:")
	log.Infof("  Platform: %s", result.Platform)
	log.Infof("  Status: %s", result.Status)
	log.Infof("  Region: %s", result.Region)
	log.Infof("  Message: %s", result.Message)
}

func TestUnlockHBOMax(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := n.UnlockHBOMax(ctx)
	if err != nil {
		t.Fatalf("UnlockHBOMax failed: %v", err)
	}

	if result.Platform != "HBO Max" {
		t.Errorf("expected platform 'HBO Max', got '%s'", result.Platform)
	}

	log.Infof("HBO Max unlock test:")
	log.Infof("  Platform: %s", result.Platform)
	log.Infof("  Status: %s", result.Status)
	log.Infof("  Region: %s", result.Region)
	log.Infof("  Message: %s", result.Message)
}

func TestUnlockPrimeVideo(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := n.UnlockPrimeVideo(ctx)
	if err != nil {
		t.Fatalf("UnlockPrimeVideo failed: %v", err)
	}

	if result.Platform != "Prime Video" {
		t.Errorf("expected platform 'Prime Video', got '%s'", result.Platform)
	}

	log.Infof("Prime Video unlock test:")
	log.Infof("  Platform: %s", result.Platform)
	log.Infof("  Status: %s", result.Status)
	log.Infof("  Region: %s", result.Region)
	log.Infof("  Message: %s", result.Message)
}

func TestUnlockBilibili(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := n.UnlockBilibili(ctx)
	if err != nil {
		t.Fatalf("UnlockBilibili failed: %v", err)
	}

	if result.Platform != "Bilibili" {
		t.Errorf("expected platform 'Bilibili', got '%s'", result.Platform)
	}

	log.Infof("Bilibili unlock test:")
	log.Infof("  Platform: %s", result.Platform)
	log.Infof("  Status: %s", result.Status)
	log.Infof("  Region: %s", result.Region)
	log.Infof("  Message: %s", result.Message)
}

func TestUnlockBahamut(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := n.UnlockBahamut(ctx)
	if err != nil {
		t.Fatalf("UnlockBahamut failed: %v", err)
	}

	if result.Platform != "Bahamut" {
		t.Errorf("expected platform 'Bahamut', got '%s'", result.Platform)
	}

	log.Infof("Bahamut unlock test:")
	log.Infof("  Platform: %s", result.Platform)
	log.Infof("  Status: %s", result.Status)
	log.Infof("  Region: %s", result.Region)
	log.Infof("  Message: %s", result.Message)
}

func TestUnlockAbemaTV(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := n.UnlockAbemaTV(ctx)
	if err != nil {
		t.Fatalf("UnlockAbemaTV failed: %v", err)
	}

	if result.Platform != "AbemaTV" {
		t.Errorf("expected platform 'AbemaTV', got '%s'", result.Platform)
	}

	log.Infof("AbemaTV unlock test:")
	log.Infof("  Platform: %s", result.Platform)
	log.Infof("  Status: %s", result.Status)
	log.Infof("  Region: %s", result.Region)
	log.Infof("  Message: %s", result.Message)
}

func TestUnlockDAZN(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := n.UnlockDAZN(ctx)
	if err != nil {
		t.Fatalf("UnlockDAZN failed: %v", err)
	}

	if result.Platform != "DAZN" {
		t.Errorf("expected platform 'DAZN', got '%s'", result.Platform)
	}

	log.Infof("DAZN unlock test:")
	log.Infof("  Platform: %s", result.Platform)
	log.Infof("  Status: %s", result.Status)
	log.Infof("  Region: %s", result.Region)
	log.Infof("  Message: %s", result.Message)
}

func TestUnlockTikTok(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := n.UnlockTikTok(ctx)
	if err != nil {
		t.Fatalf("UnlockTikTok failed: %v", err)
	}

	if result.Platform != "TikTok" {
		t.Errorf("expected platform 'TikTok', got '%s'", result.Platform)
	}

	log.Infof("TikTok unlock test:")
	log.Infof("  Platform: %s", result.Platform)
	log.Infof("  Status: %s", result.Status)
	log.Infof("  Region: %s", result.Region)
	log.Infof("  Message: %s", result.Message)
}

func TestUnlockTest(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with Netflix
	result, err := n.UnlockTest(ctx, "netflix")
	if err != nil {
		t.Fatalf("UnlockTest(netflix) failed: %v", err)
	}

	if result.Platform != "Netflix" {
		t.Errorf("expected platform 'Netflix', got '%s'", result.Platform)
	}

	log.Infof("UnlockTest (Netflix):")
	log.Infof("  Status: %s", result.Status)
	log.Infof("  Message: %s", result.Message)

	// Test with unsupported platform
	_, err = n.UnlockTest(ctx, "unsupported-platform")
	if err == nil {
		t.Error("expected error for unsupported platform, got nil")
	}
	log.Infof("Unsupported platform test passed: got expected error")
}

func TestUnlockTestAll(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	results, err := n.UnlockTestAll(ctx)
	if err != nil {
		t.Fatalf("UnlockTestAll failed: %v", err)
	}

	expectedPlatforms := []string{
		"netflix", "disney", "youtube", "hulu", "hbo",
		"primevideo", "bilibili", "bahamut", "abematv", "dazn", "tiktok",
	}

	for _, platform := range expectedPlatforms {
		result, ok := results[platform]
		if !ok {
			t.Errorf("missing result for platform '%s'", platform)
			continue
		}

		log.Infof("Platform: %s", result.Platform)
		log.Infof("  Status: %s", result.Status)
		log.Infof("  Region: %s", result.Region)
		log.Infof("  Message: %s", result.Message)
	}

	if len(results) != len(expectedPlatforms) {
		t.Errorf("expected %d results, got %d", len(expectedPlatforms), len(results))
	}
}

func TestUnlockStatusString(t *testing.T) {
	tests := []struct {
		status   node.UnlockStatus
		expected string
	}{
		{node.UnlockStatusUnknown, "Unknown"},
		{node.UnlockStatusYes, "Yes"},
		{node.UnlockStatusNo, "No"},
		{node.UnlockStatusBanned, "Banned"},
		{node.UnlockStatusFailed, "Failed"},
		{node.UnlockStatusNotAvailable, "Not Available"},
	}

	for _, tt := range tests {
		result := tt.status.String()
		if result != tt.expected {
			t.Errorf("status %d: expected '%s', got '%s'", tt.status, tt.expected, result)
		}
	}
}
