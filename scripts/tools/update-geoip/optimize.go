package main

import (
	"runtime"
	"runtime/debug"
	"time"

	"github.com/lazygophers/log"
)

// OptimizeMemory 优化内存设置
func OptimizeMemory() {
	// 设置 GC 百分比，降低内存使用（默认 100）
	// 值越小，GC 越频繁，内存占用越低，但 CPU 使用会增加
	debug.SetGCPercent(50)

	// 设置内存限制（可选）
	// debug.SetMemoryLimit(1024 * 1024 * 1024) // 1GB

	log.Infof("内存优化: GC 百分比设置为 50%%")
}

// OptimizeGoroutines 优化 goroutine 设置
func OptimizeGoroutines() {
	// 限制最大 CPU 使用数（默认为所有核心）
	// 可以设置为较少的核心以降低 CPU 占用
	numCPU := runtime.NumCPU()

	// 使用 50% 的 CPU 核心
	maxProcs := numCPU / 2
	if maxProcs < 1 {
		maxProcs = 1
	}

	runtime.GOMAXPROCS(maxProcs)
	log.Infof("CPU 优化: 使用 %d/%d 个 CPU 核心", maxProcs, numCPU)
}

// PeriodicGC 定期触发 GC
func PeriodicGC(interval time.Duration, done <-chan struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			before := m.Alloc / 1024 / 1024

			runtime.GC()

			runtime.ReadMemStats(&m)
			after := m.Alloc / 1024 / 1024

			log.Debugf("定期 GC: %d MB -> %d MB (释放 %d MB)", before, after, before-after)

		case <-done:
			return
		}
	}
}

// CompactMemory 压缩内存
func CompactMemory() {
	// 手动触发 GC
	runtime.GC()

	// 返回内存给操作系统
	debug.FreeOSMemory()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Debugf("内存压缩完成: 当前使用 %d MB", m.Alloc/1024/1024)
}
