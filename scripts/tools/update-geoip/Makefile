.PHONY: help build run test bench profile clean

help:
	@echo "GeoIP 更新工具 - Makefile"
	@echo ""
	@echo "可用命令:"
	@echo "  make build     - 编译工具"
	@echo "  make run       - 运行工具"
	@echo "  make test      - 运行测试"
	@echo "  make bench     - 运行基准测试"
	@echo "  make profile   - 运行性能分析"
	@echo "  make analyze   - 分析性能数据"
	@echo "  make clean     - 清理生成文件"

build:
	@echo "编译 update-geoip..."
	@go build -o update-geoip

run:
	@echo "运行 update-geoip..."
	@./update-geoip

test:
	@echo "运行测试..."
	@go test -v

bench:
	@echo "运行基准测试..."
	@go test -bench=. -benchmem -benchtime=10s

bench-mem:
	@echo "运行内存基准测试..."
	@go test -bench=BenchmarkMemoryUsage -benchmem -benchtime=5s

bench-csv:
	@echo "运行 CSV 基准测试..."
	@go test -bench=BenchmarkConvertCSV -benchmem -benchtime=10s

bench-text:
	@echo "运行文本基准测试..."
	@go test -bench=BenchmarkConvertText -benchmem -benchtime=10s

profile:
	@echo "运行性能分析..."
	@mkdir -p profiles
	@./update-geoip --profile --profile-dir=./profiles

profile-cpu:
	@echo "分析 CPU profile..."
	@go tool pprof -http=:8080 ./profiles/cpu.prof

profile-mem:
	@echo "分析内存 profile..."
	@go tool pprof -http=:8080 ./profiles/mem.prof

profile-goroutine:
	@echo "分析 goroutine profile..."
	@go tool pprof -http=:8080 ./profiles/mem.prof.goroutine

analyze: profile-cpu

clean:
	@echo "清理生成文件..."
	@rm -f update-geoip
	@rm -rf profiles/
	@rm -f *.prof
	@rm -f *.test
	@echo "清理完成"
