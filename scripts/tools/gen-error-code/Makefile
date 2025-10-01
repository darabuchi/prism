.PHONY: all build run clean install help

# 项目根目录（相对于 Makefile 的位置）
PROJECT_ROOT := ../../..

# 输入输出文件路径
INPUT_FILE := $(PROJECT_ROOT)/error_code.json
OUTPUT_FILE := $(PROJECT_ROOT)/error_code.gen.go

# 构建目标
BINARY := gen-error-code
BUILD_DIR := ./bin

all: build run

# 构建工具
build:
	@echo "构建错误码生成工具..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY) .
	@echo "✅ 构建完成: $(BUILD_DIR)/$(BINARY)"

# 运行工具生成错误码文件
run: build
	@echo "生成错误码文件..."
	@$(BUILD_DIR)/$(BINARY) \
		--input $(INPUT_FILE) \
		--output $(OUTPUT_FILE)
	@echo "✅ 错误码文件已生成: $(OUTPUT_FILE)"

# 安装工具到 GOPATH/bin
install:
	@echo "安装错误码生成工具..."
	@go install .
	@echo "✅ 工具已安装到 GOPATH/bin"

# 清理构建产物
clean:
	@echo "清理构建产物..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(OUTPUT_FILE)
	@echo "✅ 清理完成"

# 验证错误码 JSON 格式
validate: build
	@echo "验证错误码 JSON 格式..."
	@$(BUILD_DIR)/$(BINARY) \
		--input $(INPUT_FILE) \
		--output /tmp/error_code_test.gen.go
	@rm -f /tmp/error_code_test.gen.go
	@echo "✅ 错误码 JSON 格式验证通过"

# 显示帮助信息
help:
	@echo "错误码生成工具 Makefile"
	@echo ""
	@echo "使用方法:"
	@echo "  make         - 构建并运行，生成错误码文件"
	@echo "  make build   - 仅构建工具"
	@echo "  make run     - 构建并运行工具"
	@echo "  make install - 安装工具到 GOPATH/bin"
	@echo "  make validate - 验证错误码 JSON 格式"
	@echo "  make clean   - 清理构建产物"
	@echo "  make help    - 显示此帮助信息"
	@echo ""
	@echo "配置:"
	@echo "  输入文件: $(INPUT_FILE)"
	@echo "  输出文件: $(OUTPUT_FILE)"
