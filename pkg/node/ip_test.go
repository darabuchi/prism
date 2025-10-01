package node_test

import (
	"context"
	"testing"
	"time"

	"github.com/darabuchi/prism/pkg/node"
	"github.com/lazygophers/log"
)

func TestGetIP(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test default service (ipify)
	ip, metadata, err := n.GetIP(ctx, "")
	if err != nil {
		t.Logf("GetIP (default) failed: %v", err)
		t.Skip("skipping IP lookup test")
		return
	}

	if ip == "" {
		t.Error("expected non-empty IP address")
	}

	log.Infof("IP address (default): %s", ip)
	if metadata != nil {
		log.Infof("Metadata: %+v", metadata)
	}
}

func TestGetIPWithSpecificService(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	services := []string{
		"ipify",
		"icanhazip",
		"ifconfig.me",
		"ident.me",
		"ip-api.com",
		"ipinfo.io",
	}

	for _, serviceName := range services {
		ip, metadata, err := n.GetIP(ctx, serviceName)
		if err != nil {
			t.Logf("GetIP (%s) failed: %v", serviceName, err)
			continue
		}

		if ip == "" {
			t.Errorf("expected non-empty IP address for service %s", serviceName)
		}

		log.Infof("IP address (%s): %s", serviceName, ip)
		if metadata != nil {
			log.Infof("  Country: %s (%s)", metadata.Country, metadata.CountryCode)
			log.Infof("  Region: %s", metadata.Region)
			log.Infof("  City: %s", metadata.City)
			log.Infof("  ISP: %s", metadata.ISP)
			log.Infof("  ASN: %d (%s)", metadata.ASN, metadata.AS)
			log.Infof("  Timezone: %s", metadata.Timezone)
		}
	}
}

func TestGetIPConcurrent(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	results, statistics, metadataStats, err := n.GetIPConcurrent(ctx)
	if err != nil {
		t.Fatalf("GetIPConcurrent failed: %v", err)
	}

	log.Infof("IP lookup results:")
	successCount := 0
	failCount := 0

	for _, result := range results {
		if result.Error != nil {
			log.Warnf("  [%s] failed: %v (delay: %v)", result.Source, result.Error, result.Delay)
			failCount++
		} else {
			log.Infof("  [%s] %s (delay: %v)", result.Source, result.IP, result.Delay)
			successCount++
		}
	}

	log.Infof("Success: %d, Failed: %d", successCount, failCount)

	if len(statistics) == 0 {
		t.Error("expected at least one IP in statistics")
		return
	}

	log.Infof("\nIP Statistics:")
	totalPercentage := 0.0
	for i, stat := range statistics {
		log.Infof("  %d. %s", i+1, stat.IP)
		log.Infof("     Count: %d", stat.Count)
		log.Infof("     Percentage: %.2f%%", stat.Percentage)
		log.Infof("     Sources: %v", stat.Sources)
		totalPercentage += stat.Percentage
	}

	// Verify percentage calculation
	if totalPercentage < 99.9 || totalPercentage > 100.1 {
		t.Errorf("total percentage should be ~100%%, got %.2f%%", totalPercentage)
	}

	// Verify most common IP
	if statistics[0].Count < 1 {
		t.Error("expected at least 1 count for most common IP")
	}

	// Verify statistics are sorted by count
	for i := 1; i < len(statistics); i++ {
		if statistics[i].Count > statistics[i-1].Count {
			t.Error("statistics should be sorted by count in descending order")
		}
	}

	// Display metadata statistics
	if metadataStats != nil {
		log.Infof("\nMetadata Statistics:")

		if len(metadataStats.Country) > 0 {
			log.Infof("  Country:")
			for country, stat := range metadataStats.Country {
				log.Infof("    %s: %d (%.2f%%) %v", country, stat.Count, stat.Percentage, stat.Sources)
			}
		}

		if len(metadataStats.Region) > 0 {
			log.Infof("  Region:")
			for region, stat := range metadataStats.Region {
				log.Infof("    %s: %d (%.2f%%) %v", region, stat.Count, stat.Percentage, stat.Sources)
			}
		}

		if len(metadataStats.City) > 0 {
			log.Infof("  City:")
			for city, stat := range metadataStats.City {
				log.Infof("    %s: %d (%.2f%%) %v", city, stat.Count, stat.Percentage, stat.Sources)
			}
		}

		if len(metadataStats.ISP) > 0 {
			log.Infof("  ISP:")
			for isp, stat := range metadataStats.ISP {
				log.Infof("    %s: %d (%.2f%%) %v", isp, stat.Count, stat.Percentage, stat.Sources)
			}
		}

		if len(metadataStats.ASN) > 0 {
			log.Infof("  ASN:")
			for asn, stat := range metadataStats.ASN {
				log.Infof("    %d: %d (%.2f%%) %v", asn, stat.Count, stat.Percentage, stat.Sources)
			}
		}

		if len(metadataStats.Timezone) > 0 {
			log.Infof("  Timezone:")
			for tz, stat := range metadataStats.Timezone {
				log.Infof("    %s: %d (%.2f%%) %v", tz, stat.Count, stat.Percentage, stat.Sources)
			}
		}
	}
}

func TestGetIPConcurrentMultipleTimes(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// Run multiple times to check consistency
	const runs = 3
	allIPs := make(map[string]int)

	for i := 0; i < runs; i++ {
		log.Infof("\n=== Run %d/%d ===", i+1, runs)

		_, statistics, _, err := n.GetIPConcurrent(ctx)
		if err != nil {
			t.Logf("GetIPConcurrent run %d failed: %v", i+1, err)
			continue
		}

		if len(statistics) > 0 {
			topIP := statistics[0].IP
			allIPs[topIP]++
			log.Infof("Top IP: %s (%.2f%%)", topIP, statistics[0].Percentage)
		}

		// Add delay between runs to avoid rate limiting
		if i < runs-1 {
			time.Sleep(2 * time.Second)
		}
	}

	log.Infof("\n=== Summary ===")
	log.Infof("Different IPs seen across %d runs:", runs)
	for ip, count := range allIPs {
		log.Infof("  %s: %d times", ip, count)
	}

	if len(allIPs) == 0 {
		t.Error("expected at least one IP across all runs")
	}
}

func TestIPLookupWithCanceledContext(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, _, err = n.GetIP(ctx, "")
	if err == nil {
		t.Error("expected error when context is canceled")
	}

	log.Infof("Canceled context test passed: got expected error")
}

func TestIPLookupWithTimeout(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	// Very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	time.Sleep(10 * time.Millisecond) // Wait for timeout

	_, _, err = n.GetIP(ctx, "")
	if err == nil {
		t.Log("Note: Expected timeout error, but request might have completed very quickly")
	} else {
		log.Infof("Timeout test passed: got expected error")
	}
}

func TestIPLookupIndividualServices(t *testing.T) {
	n, err := node.NewNode(testProxyConfig)
	if err != nil {
		t.Fatalf("NewNode failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	services := []string{
		"ipify",
		"icanhazip",
		"ifconfig.me",
		"ident.me",
		"api.ipify",
		"checkip.amazonaws",
		"ipecho.net",
		"myip.com",
		"ip-api.com",
		"ipinfo.io",
		"cloudflare",
		"google-dns",
	}

	successCount := 0
	for _, serviceName := range services {
		ip, metadata, err := n.GetIP(ctx, serviceName)
		if err != nil {
			log.Warnf("Service %s failed: %v", serviceName, err)
			continue
		}

		if ip == "" {
			t.Errorf("Service %s returned empty IP", serviceName)
			continue
		}

		log.Infof("Service %s: %s", serviceName, ip)
		if metadata != nil && (metadata.Country != "" || metadata.City != "" || metadata.ASN > 0) {
			log.Infof("  Metadata: Country=%s, City=%s, ASN=%d", metadata.Country, metadata.City, metadata.ASN)
		}
		successCount++
	}

	if successCount == 0 {
		t.Error("expected at least one service to succeed")
	}

	log.Infof("Successfully queried %d/%d services", successCount, len(services))
}
