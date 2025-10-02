.PHONY: run build clean help

# 默认目标
.DEFAULT_GOAL := help

# 运行规则收集器
run:
	@echo "Running rule collector..."
	go run *.go

# 构建二进制文件
build:
	@echo "Building rule-collector..."
	go build -o bin/rule-collector *.go
	@echo "Binary built: bin/rule-collector"

# 清理缓存和输出文件
clean:
	@echo "Cleaning cache and output files..."
	rm -rf tmp/cache/*
	rm -rf output/*
	rm -f bin/rule-collector
	@echo "Clean completed"

# 清理缓存
clean-cache:
	@echo "Cleaning cache..."
	rm -rf tmp/cache/*
	@echo "Cache cleaned"

# 清理输出
clean-output:
	@echo "Cleaning output..."
	rm -rf output/*
	@echo "Output cleaned"

# 强制刷新（清理缓存后重新运行）
refresh: clean-cache run

# 显示帮助信息
help:
	@echo "Prism Rule Collector - Makefile Commands"
	@echo ""
	@echo "Usage:"
	@echo "  make run           - Run the rule collector"
	@echo "  make build         - Build binary executable"
	@echo "  make clean         - Clean cache and output files"
	@echo "  make clean-cache   - Clean cache only"
	@echo "  make clean-output  - Clean output only"
	@echo "  make refresh       - Clean cache and re-run"
	@echo "  make help          - Show this help message"
