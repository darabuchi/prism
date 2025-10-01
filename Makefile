.PHONY: help build build-server build-cli build-all run run-server test test-unit test-integration test-coverage clean lint fmt web-dev web-build desktop-dev desktop-build release release-local release-snapshot release-check update-geoip clean-geoip-cache gen-error-code validate-error-code

# 默认目标
.DEFAULT_GOAL := help

# 项目名称
PROJECT_NAME := prism
# 构建输出目录
BUILD_DIR := bin
# 主程序路径
MAIN_PATH := ./cmd/main.go

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
	@echo ""
	@echo "Backend:"
	@echo "  build             - 编译项目"
	@echo "  run               - 运行服务器"
	@echo "  test              - 运行所有测试"
	@echo "  test-unit         - 运行单元测试"
	@echo "  test-integration  - 运行集成测试"
	@echo "  test-coverage     - 生成测试覆盖率报告"
	@echo "  lint              - 运行代码检查"
	@echo "  fmt               - 格式化代码"
	@echo ""
	@echo "Release:"
	@echo "  release-local     - 本地构建当前平台（goreleaser）"
	@echo "  release-snapshot  - 构建所有平台快照版本（goreleaser）"
	@echo "  release-check     - 检查 goreleaser 配置"
	@echo "  release           - 发布版本（需要 Git 标签）"
	@echo ""
	@echo "Frontend:"
	@echo "  web-dev           - 启动 Web 前端开发服务器"
	@echo "  web-build         - 构建 Web 前端生产版本"
	@echo "  desktop-dev       - 启动桌面应用开发模式"
	@echo "  desktop-build     - 构建桌面应用"
	@echo ""
	@echo "Utils:"
	@echo "  clean             - 清理构建产物"
	@echo "  deps              - 下载依赖"
	@echo "  tidy              - 整理依赖"
	@echo ""
	@echo "GeoIP:"
	@echo "  update-geoip      - 更新 GeoIP 数据库"
	@echo "  clean-geoip-cache - 清理 GeoIP 缓存（强制重新下载）"
	@echo ""
	@echo "Code Generation:"
	@echo "  gen-error-code      - 根据 error_code.json 生成错误码文件"
	@echo "  validate-error-code - 验证 error_code.json 格式"

## build: 编译项目
build:
	@echo "Building $(PROJECT_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(PROJECT_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(PROJECT_NAME)"

## run: 运行服务器
run: build
	@echo "Running $(PROJECT_NAME) server..."
	@$(BUILD_DIR)/$(PROJECT_NAME) server

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

## web-dev: 启动 Web 前端开发服务器
web-dev:
	@echo "Starting web frontend development server..."
	@if [ -d "web/frontend" ]; then \
		cd web/frontend && pnpm dev; \
	else \
		echo "web/frontend directory not found"; \
	fi

## web-build: 构建 Web 前端生产版本
web-build:
	@echo "Building web frontend..."
	@if [ -d "web/frontend" ]; then \
		cd web/frontend && pnpm build; \
	fi
	@if [ -d "web/admin" ]; then \
		cd web/admin && pnpm build; \
	fi
	@echo "Web frontend build complete"

## desktop-dev: 启动桌面应用开发模式（Tauri）
desktop-dev:
	@echo "Starting desktop app development (Tauri)..."
	@if [ -d "desktop" ]; then \
		cd desktop && pnpm tauri dev; \
	else \
		echo "desktop directory not found"; \
	fi

## desktop-build: 构建桌面应用（Tauri）
desktop-build:
	@echo "Building desktop app (Tauri)..."
	@if [ -d "desktop" ]; then \
		cd desktop && pnpm tauri build; \
	else \
		echo "desktop directory not found"; \
	fi
	@echo "Desktop app build complete"

## release-local: 本地构建当前平台（goreleaser）
release-local:
	@echo "Building local release..."
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser build --single-target --snapshot --clean; \
	else \
		echo "goreleaser not found. Install with: brew install goreleaser"; \
		exit 1; \
	fi
	@echo "Local release build complete: dist/"

## release-snapshot: 构建所有平台快照版本（goreleaser）
release-snapshot:
	@echo "Building snapshot release for all platforms..."
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser release --snapshot --clean; \
	else \
		echo "goreleaser not found. Install with: brew install goreleaser"; \
		exit 1; \
	fi
	@echo "Snapshot release build complete: dist/"

## release-check: 检查 goreleaser 配置
release-check:
	@echo "Checking goreleaser configuration..."
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser check; \
	else \
		echo "goreleaser not found. Install with: brew install goreleaser"; \
		exit 1; \
	fi

## release: 发布版本（需要 Git 标签和 GITHUB_TOKEN）
release:
	@echo "Releasing version..."
	@if [ -z "$(GITHUB_TOKEN)" ]; then \
		echo "GITHUB_TOKEN is not set. Please export GITHUB_TOKEN"; \
		exit 1; \
	fi
	@if command -v goreleaser >/dev/null 2>&1; then \
		goreleaser release --clean; \
	else \
		echo "goreleaser not found. Install with: brew install goreleaser"; \
		exit 1; \
	fi
	@echo "Release complete"

## update-geoip: 更新 GeoIP 数据库
update-geoip:
	@echo "Updating GeoIP databases..."
	@go run ./scripts/tools/update-geoip

## clean-geoip-cache: 清理 GeoIP 缓存（强制重新下载）
clean-geoip-cache:
	@echo "Cleaning GeoIP cache..."
	@rm -rf /tmp/prism_geoip
	@echo "GeoIP cache cleaned"

## gen-error-code: 根据 error_code.json 生成错误码文件
gen-error-code:
	@echo "Generating error code file..."
	@cd scripts/tools/gen-error-code && $(MAKE) run
	@echo "Error code file generated: error_code.gen.go"

## validate-error-code: 验证 error_code.json 格式
validate-error-code:
	@echo "Validating error_code.json..."
	@cd scripts/tools/gen-error-code && $(MAKE) validate
	@echo "Error code validation passed"
