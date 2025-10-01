# GeoIP 更新工具性能优化指南

## 性能目标

1. **内存使用**（最高优先级）：保持在合理范围内，避免 OOM
2. **CPU 使用**：降低 CPU 占用率
3. **处理速度**：在前两者约束下尽可能快

## 性能测试

### 基准测试

```bash
# 运行所有基准测试
go test -bench=. -benchmem -benchtime=10s

# 运行特定测试
go test -bench=BenchmarkMemoryUsage -benchmem

# 生成 CPU profile
go test -bench=BenchmarkConvertCSV -cpuprofile=cpu.prof -memprofile=mem.prof

# 查看 profile
go tool pprof -http=:8080 cpu.prof
go tool pprof -http=:8080 mem.prof
```

### 实际数据测试

```bash
# 启用性能分析运行
./update-geoip --profile --profile-dir=./profiles

# 查看生成的 profile 文件
go tool pprof -http=:8080 ./profiles/cpu.prof
go tool pprof -http=:8080 ./profiles/mem.prof
```

## 性能优化策略

### 1. 内存优化（最高优先级）

#### 当前实现
- ✅ 流式读取：CSV/文本文件逐行处理
- ✅ 批量插入：5000 条记录批量插入，减少单次插入开销
- ✅ 对象池：sync.Pool 复用 GeoIP 和 StringBuilder 对象
- ✅ 激进 GC：每处理 3 个数据源触发一次压缩 GC
- ✅ 内存压缩：使用 debug.FreeOSMemory 归还内存给 OS
- ✅ GC 优化：GC 百分比设置为 50%（更频繁但内存更低）
- ✅ 资源释放：defer 确保文件句柄及时关闭
- ✅ 避免缓冲：使用 io.Copy 流式写入

#### 待优化项
- [ ] 增量写入：分阶段写入 MMDB 树（如果可能）
- [ ] 字符串内联：减少字符串分配
- [ ] 并发处理：goroutine 池（需平衡内存）

### 2. CPU 优化

#### 当前实现
- ✅ CPU 核心限制：默认使用 50% CPU 核心
- ✅ 避免正则表达式：使用字符串操作
- ✅ 减少类型转换：最小化接口转换
- ✅ 高效解析：优化 IP 和 CIDR 解析
- ✅ 批量操作：减少函数调用开销

#### 待优化项
- [ ] 并发处理：goroutine 池并发转换
- [ ] SIMD 优化：使用 SIMD 指令加速
- [ ] 缓存热点：缓存常用数据
- [ ] 减少日志：降低日志级别

### 3. 速度优化

#### 当前实现
- ✅ 优先级排序：按优先级处理数据源
- ✅ 跳过无效数据：快速过滤无效记录
- ✅ 进度显示：每 10000 条记录显示一次

#### 待优化项
- [ ] 并发下载：并行下载多个数据源
- [ ] 增量更新：只更新变化的数据
- [ ] 预处理：提前验证和规范化数据

## 性能基线

### 内存使用（目标）
- 峰值内存：< 500 MB（小数据集）
- 峰值内存：< 2 GB（完整数据集）
- GC 压力：< 10 次/秒

### CPU 使用（目标）
- 平均 CPU：< 50%（单核）
- 峰值 CPU：< 80%（单核）

### 处理速度（参考）
- CSV 处理：> 50,000 行/秒
- CIDR 处理：> 100,000 行/秒
- 总处理时间：< 5 分钟（完整数据集）

## 监控指标

### 关键指标
1. **内存分配（Alloc）**：当前分配的内存
2. **总分配（TotalAlloc）**：累计分配的内存
3. **系统内存（Sys）**：从系统申请的内存
4. **GC 次数（NumGC）**：垃圾回收次数
5. **Goroutine 数量**：并发度

### 查看方式
```bash
# 运行时监控
./update-geoip --profile

# 实时内存监控
watch -n 1 'ps aux | grep update-geoip'

# pprof 交互式分析
go tool pprof -http=:8080 ./profiles/mem.prof
```

## 优化流程

1. **建立基线**
   ```bash
   go test -bench=BenchmarkMemoryUsage/LargeDataset -benchmem > baseline.txt
   ```

2. **实施优化**
   - 修改代码
   - 确保测试通过

3. **对比测试**
   ```bash
   go test -bench=BenchmarkMemoryUsage/LargeDataset -benchmem > optimized.txt
   benchcmp baseline.txt optimized.txt
   ```

4. **Profile 分析**
   ```bash
   go tool pprof -http=:8080 mem.prof
   # 查看 Top、Graph、Flame Graph
   ```

5. **验证改进**
   - 内存降低 > 20%
   - 速度保持或提升
   - 功能正确性不变

## 常见问题

### Q: 内存持续增长怎么办？
A:
1. 检查是否有内存泄漏（goroutine 泄漏、引用循环）
2. 增加 GC 频率
3. 使用对象池减少分配
4. 检查 mmdbwriter 是否正确释放

### Q: CPU 占用过高怎么办？
A:
1. 减少字符串拼接和转换
2. 优化热点函数（pprof 识别）
3. 降低日志级别
4. 使用更高效的数据结构

### Q: 处理速度慢怎么办？
A:
1. 启用并发处理（注意内存）
2. 优化 I/O（批量读写）
3. 减少数据验证
4. 缓存重复计算

## 参考资源

- [Go Performance Tools](https://go.dev/doc/diagnostics)
- [pprof User Guide](https://github.com/google/pprof/blob/main/doc/README.md)
- [Memory Optimization](https://go.dev/blog/pprof)
- [Benchmarking Go](https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go)
