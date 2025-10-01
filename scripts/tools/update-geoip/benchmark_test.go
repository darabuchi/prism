package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
)

// BenchmarkConvertCSV 测试 CSV 转换性能
func BenchmarkConvertCSV(b *testing.B) {
	// 使用实际的测试文件（需要先下载）
	cacheDir := "/tmp/prism_geoip_bench"
	os.MkdirAll(cacheDir, 0755)
	defer os.RemoveAll(cacheDir)

	// 创建小型测试数据
	testFile := filepath.Join(cacheDir, "test.csv")
	createTestCSV(b, testFile)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		converter, _ := NewConverter()
		if err := converter.ConvertCSV(testFile, ConvertIPInfoCSV); err != nil {
			b.Fatalf("转换失败: %v", err)
		}
	}
}

// BenchmarkConvertText 测试文本转换性能
func BenchmarkConvertText(b *testing.B) {
	cacheDir := "/tmp/prism_geoip_bench"
	os.MkdirAll(cacheDir, 0755)
	defer os.RemoveAll(cacheDir)

	testFile := filepath.Join(cacheDir, "test.txt")
	createTestText(b, testFile)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		converter, _ := NewConverter()
		ds := DownloadedSource{
			Name:   "test",
			Path:   testFile,
			Source: Source{Type: SourceTypeText, Format: "cidr-cn"},
		}
		if err := convertTextSource(converter, ds); err != nil {
			b.Fatalf("转换失败: %v", err)
		}
	}
}

// BenchmarkMemoryUsage 内存使用基准测试
func BenchmarkMemoryUsage(b *testing.B) {
	b.Run("SmallDataset", func(b *testing.B) {
		benchmarkMemoryUsage(b, 1000)
	})

	b.Run("MediumDataset", func(b *testing.B) {
		benchmarkMemoryUsage(b, 10000)
	})

	b.Run("LargeDataset", func(b *testing.B) {
		benchmarkMemoryUsage(b, 100000)
	})
}

func benchmarkMemoryUsage(b *testing.B, recordCount int) {
	cacheDir := "/tmp/prism_geoip_bench"
	os.MkdirAll(cacheDir, 0755)
	defer os.RemoveAll(cacheDir)

	testFile := filepath.Join(cacheDir, "test.csv")
	createTestCSVWithSize(b, testFile, recordCount)

	var before, after runtime.MemStats

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		runtime.GC()
		runtime.ReadMemStats(&before)

		converter, _ := NewConverter()
		converter.ConvertCSV(testFile, ConvertIPInfoCSV)

		runtime.ReadMemStats(&after)

		allocMB := float64(after.Alloc-before.Alloc) / 1024 / 1024
		b.ReportMetric(allocMB, "MB/op")
	}
}

// 辅助函数：创建测试 CSV 文件
func createTestCSV(b *testing.B, path string) {
	createTestCSVWithSize(b, path, 1000)
}

func createTestCSVWithSize(b *testing.B, path string, recordCount int) {
	f, err := os.Create(path)
	if err != nil {
		b.Fatalf("创建测试文件失败: %v", err)
	}
	defer f.Close()

	// 写入表头
	f.WriteString("start_ip,end_ip,country,country_name,continent,continent_name,asn,as_name,as_domain\n")

	// 写入测试数据
	for i := 0; i < recordCount; i++ {
		baseIP := 1000000 + i*256
		startIP := formatIP(baseIP)
		endIP := formatIP(baseIP + 255)
		f.WriteString(startIP + "," + endIP + ",US,United States,NA,North America,AS15169,Google LLC,google.com\n")
	}
}

// 辅助函数：创建测试文本文件
func createTestText(b *testing.B, path string) {
	f, err := os.Create(path)
	if err != nil {
		b.Fatalf("创建测试文件失败: %v", err)
	}
	defer f.Close()

	// 写入测试 CIDR
	for i := 0; i < 1000; i++ {
		f.WriteString("10." + formatByte(i/256) + "." + formatByte(i%256) + ".0/24\n")
	}
}

func formatIP(n int) string {
	return formatByte(n>>24) + "." + formatByte((n>>16)&0xFF) + "." + formatByte((n>>8)&0xFF) + "." + formatByte(n&0xFF)
}

func formatByte(n int) string {
	return strconv.Itoa(n)
}
