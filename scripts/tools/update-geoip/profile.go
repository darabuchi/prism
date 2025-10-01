package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/lazygophers/log"
)

// Profiler 性能分析器
type Profiler struct {
	cpuFile    *os.File
	memFile    string
	startTime  time.Time
	startMem   runtime.MemStats
	enabled    bool
	sampleRate int // 采样率（每N个操作采样一次）
}

// NewProfiler 创建性能分析器
func NewProfiler(enabled bool, outputDir string) (*Profiler, error) {
	if !enabled {
		return &Profiler{enabled: false}, nil
	}

	os.MkdirAll(outputDir, 0755)

	// CPU profiling
	cpuFile, err := os.Create(fmt.Sprintf("%s/cpu.prof", outputDir))
	if err != nil {
		return nil, fmt.Errorf("创建 CPU profile 文件失败: %w", err)
	}

	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		cpuFile.Close()
		return nil, fmt.Errorf("启动 CPU profiling 失败: %w", err)
	}

	p := &Profiler{
		cpuFile:    cpuFile,
		memFile:    fmt.Sprintf("%s/mem.prof", outputDir),
		startTime:  time.Now(),
		enabled:    true,
		sampleRate: 10000, // 每10000条记录采样一次
	}

	runtime.ReadMemStats(&p.startMem)

	log.Infof("性能分析已启动，输出目录: %s", outputDir)
	return p, nil
}

// Stop 停止性能分析
func (p *Profiler) Stop() error {
	if !p.enabled {
		return nil
	}

	// 停止 CPU profiling
	pprof.StopCPUProfile()
	p.cpuFile.Close()

	// 写入内存 profile
	memFile, err := os.Create(p.memFile)
	if err != nil {
		return fmt.Errorf("创建内存 profile 文件失败: %w", err)
	}
	defer memFile.Close()

	runtime.GC()
	if err := pprof.WriteHeapProfile(memFile); err != nil {
		return fmt.Errorf("写入内存 profile 失败: %w", err)
	}

	// 写入 goroutine profile
	goroutineFile, err := os.Create(p.memFile + ".goroutine")
	if err != nil {
		return fmt.Errorf("创建 goroutine profile 文件失败: %w", err)
	}
	defer goroutineFile.Close()

	if err := pprof.Lookup("goroutine").WriteTo(goroutineFile, 0); err != nil {
		return fmt.Errorf("写入 goroutine profile 失败: %w", err)
	}

	// 统计信息
	var endMem runtime.MemStats
	runtime.ReadMemStats(&endMem)

	duration := time.Since(p.startTime)
	allocMB := float64(endMem.TotalAlloc-p.startMem.TotalAlloc) / 1024 / 1024
	currentMB := float64(endMem.Alloc) / 1024 / 1024

	log.Infof("")
	log.Infof("性能分析完成:")
	log.Infof("  执行时间: %v", duration)
	log.Infof("  总分配内存: %.2f MB", allocMB)
	log.Infof("  当前内存使用: %.2f MB", currentMB)
	log.Infof("  GC 次数: %d", endMem.NumGC-p.startMem.NumGC)
	log.Infof("")
	log.Infof("分析文件:")
	log.Infof("  CPU: %s", p.cpuFile.Name())
	log.Infof("  内存: %s", p.memFile)
	log.Infof("  Goroutine: %s", p.memFile+".goroutine")
	log.Infof("")
	log.Infof("使用以下命令查看:")
	log.Infof("  go tool pprof -http=:8080 %s", p.cpuFile.Name())
	log.Infof("  go tool pprof -http=:8080 %s", p.memFile)

	return nil
}

// Checkpoint 记录检查点
func (p *Profiler) Checkpoint(name string) {
	if !p.enabled {
		return
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	log.Debugf("性能检查点 [%s]: 内存 %.2f MB, Goroutines %d",
		name, float64(m.Alloc)/1024/1024, runtime.NumGoroutine())
}

// ShouldSample 是否应该采样（用于降低采样开销）
func (p *Profiler) ShouldSample(counter int) bool {
	if !p.enabled {
		return false
	}
	return counter%p.sampleRate == 0
}

// MemoryStats 内存统计信息
type MemoryStats struct {
	AllocMB      float64
	TotalAllocMB float64
	SysMB        float64
	NumGC        uint32
	Goroutines   int
}

// GetMemoryStats 获取当前内存统计
func GetMemoryStats() MemoryStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return MemoryStats{
		AllocMB:      float64(m.Alloc) / 1024 / 1024,
		TotalAllocMB: float64(m.TotalAlloc) / 1024 / 1024,
		SysMB:        float64(m.Sys) / 1024 / 1024,
		NumGC:        m.NumGC,
		Goroutines:   runtime.NumGoroutine(),
	}
}

// LogMemoryStats 记录内存统计
func LogMemoryStats(prefix string) {
	stats := GetMemoryStats()
	log.Infof("%s: Alloc=%.2f MB, TotalAlloc=%.2f MB, Sys=%.2f MB, NumGC=%d, Goroutines=%d",
		prefix, stats.AllocMB, stats.TotalAllocMB, stats.SysMB, stats.NumGC, stats.Goroutines)
}
