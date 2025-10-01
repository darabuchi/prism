.PHONY: help build run test test-unit test-integration test-coverage clean lint fmt

# 默认目标
.DEFAULT_GOAL := help

# 项目名称
PROJECT_NAME := prism
# 构建输出目录
BUILD_DIR := bin
# 主程序路径
MAIN_PATH := ./cmd/prism

# Go 相关变量
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOMOD := $(GOCMD) mod
GOFMT := $(GOCMD) fmt
GOVET := $(GOCMD) vet

# 构建标志
BUILD_FLAGS := -v
LDFLAGS := -s -w

## help: 显示帮助信息
help:
	@echo "Available targets:"
	@echo "  build             - 编译项目"
	@echo "  run               - 运行项目"
	@echo "  test              - 运行所有测试"
	@echo "  test-unit         - 运行单元测试"
	@echo "  test-integration  - 运行集成测试"
	@echo "  test-coverage     - 生成测试覆盖率报告"
	@echo "  clean             - 清理构建产物"
	@echo "  lint              - 运行代码检查"
	@echo "  fmt               - 格式化代码"
	@echo "  deps              - 下载依赖"
	@echo "  tidy              - 整理依赖"

## build: 编译项目
build:
	@echo "Building $(PROJECT_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(PROJECT_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(PROJECT_NAME)"

## run: 运行项目
run: build
	@echo "Running $(PROJECT_NAME)..."
	@$(BUILD_DIR)/$(PROJECT_NAME)

## test: 运行所有测试
test:
	@echo "Running all tests..."
	$(GOTEST) -v -race ./...

## test-unit: 运行单元测试
test-unit:
	@echo "Running unit tests..."
	$(GOTEST) -v -race -short ./...

## test-integration: 运行集成测试
test-integration:
	@echo "Running integration tests..."
	$(GOTEST) -v -race ./tests/integration/...

## test-coverage: 生成测试覆盖率报告
test-coverage:
	@echo "Generating coverage report..."
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## clean: 清理构建产物
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

## lint: 运行代码检查
lint:
	@echo "Running linters..."
	$(GOVET) ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, skipping..."; \
	fi

## fmt: 格式化代码
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

## deps: 下载依赖
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download

## tidy: 整理依赖
tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

## docker-build: 构建 Docker 镜像
docker-build:
	@echo "Building Docker image..."
	docker build -t $(PROJECT_NAME):latest -f deployments/docker/Dockerfile .

## docker-run: 运行 Docker 容器
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 -p 1080:1080 $(PROJECT_NAME):latest
