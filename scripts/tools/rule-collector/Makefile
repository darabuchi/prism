.PHONY: help build run clean clean-cache clean-output refresh

# 默认目标
.DEFAULT_GOAL := help

help:
	@echo "Prism 规则收集器 - Makefile"
	@echo ""
	@echo "可用命令:"
	@echo "  make build         - 编译工具"
	@echo "  make run           - 运行工具"
	@echo "  make clean         - 清理所有生成文件"
	@echo "  make clean-cache   - 仅清理缓存"
	@echo "  make clean-output  - 仅清理输出"
	@echo "  make refresh       - 清理缓存并重新运行"

build:
	@echo "编译 rule-collector..."
	@go build -o rule-collector *.go
	@echo "编译完成: ./rule-collector"

run:
	@echo "运行 rule-collector..."
	@go run *.go

clean:
	@echo "清理生成文件..."
	@rm -f rule-collector
	@rm -rf output/
	@echo "清理完成"

clean-cache:
	@echo "清理缓存..."
	@rm -rf /tmp/prism_rule_collector/
	@echo "缓存清理完成"

clean-output:
	@echo "清理输出..."
	@rm -rf output/
	@echo "输出清理完成"

refresh: clean-cache run
