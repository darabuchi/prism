# Prism - 构建和开发工具

.PHONY: doctor help build test clean install dev deps

# 默认目标
.DEFAULT_GOAL := help

# 版本信息
VERSION ?= v1.0.0
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD)

# 构建标志
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)"

# 平台配置
PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

# 颜色输出
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
BLUE := \033[34m
RESET := \033[0m

##@ 环境检查

doctor: ## 检查开发环境是否满足要求
	@echo "$(BLUE)🔍 检查 Prism 开发环境...$(RESET)"
	@echo ""
	@$(MAKE) -s check-system
	@$(MAKE) -s check-go
	@$(MAKE) -s check-node
	@$(MAKE) -s check-rust
	@$(MAKE) -s check-android
	@$(MAKE) -s check-docker
	@$(MAKE) -s check-git
	@echo ""
	@echo "$(GREEN)✅ 环境检查完成！$(RESET)"

check-system: ## 检查系统信息
	@echo "$(YELLOW)📋 系统信息:$(RESET)"
	@echo "  操作系统: $$(uname -s)"
	@echo "  架构: $$(uname -m)"
	@echo "  内核版本: $$(uname -r)"
	@echo ""

check-go: ## 检查 Go 环境
	@echo "$(YELLOW)🐹 Go 环境:$(RESET)"
	@if command -v go >/dev/null 2>&1; then \
		echo "  $(GREEN)✓$(RESET) Go 已安装: $$(go version)"; \
		echo "  GOPATH: $$(go env GOPATH)"; \
		echo "  GOROOT: $$(go env GOROOT)"; \
		GO_VERSION=$$(go version | grep -o 'go[0-9]\+\.[0-9]\+' | sed 's/go//'); \
		if [ "$$(printf '%s\n' "1.21" "$$GO_VERSION" | sort -V | head -n1)" = "1.21" ]; then \
			echo "  $(GREEN)✓$(RESET) Go 版本满足要求 (>= 1.21)"; \
		else \
			echo "  $(RED)✗$(RESET) Go 版本过低，需要 >= 1.21"; \
		fi; \
		if go mod download >/dev/null 2>&1; then \
			echo "  $(GREEN)✓$(RESET) Go 模块下载正常"; \
		else \
			echo "  $(YELLOW)⚠$(RESET) Go 模块下载可能有问题"; \
		fi; \
	else \
		echo "  $(RED)✗$(RESET) Go 未安装"; \
		echo "  安装方法:"; \
		echo "    macOS: brew install go"; \
		echo "    Ubuntu: sudo apt install golang-go"; \
		echo "    Windows: winget install GoLang.Go"; \
	fi
	@echo ""

check-node: ## 检查 Node.js 环境
	@echo "$(YELLOW)📦 Node.js 环境:$(RESET)"
	@if command -v node >/dev/null 2>&1; then \
		echo "  $(GREEN)✓$(RESET) Node.js 已安装: $$(node --version)"; \
		NODE_VERSION=$$(node --version | sed 's/v//'); \
		if [ "$$(printf '%s\n' "18.0.0" "$$NODE_VERSION" | sort -V | head -n1)" = "18.0.0" ]; then \
			echo "  $(GREEN)✓$(RESET) Node.js 版本满足要求 (>= 18.0.0)"; \
		else \
			echo "  $(RED)✗$(RESET) Node.js 版本过低，需要 >= 18.0.0"; \
		fi; \
	else \
		echo "  $(RED)✗$(RESET) Node.js 未安装"; \
		echo "  安装方法:"; \
		echo "    macOS: brew install node"; \
		echo "    Ubuntu: sudo apt install nodejs npm"; \
		echo "    Windows: winget install OpenJS.NodeJS"; \
	fi
	@if command -v npm >/dev/null 2>&1; then \
		echo "  $(GREEN)✓$(RESET) npm 已安装: $$(npm --version)"; \
	else \
		echo "  $(RED)✗$(RESET) npm 未安装"; \
	fi
	@echo ""

check-rust: ## 检查 Rust 环境
	@echo "$(YELLOW)🦀 Rust 环境:$(RESET)"
	@if command -v rustc >/dev/null 2>&1; then \
		echo "  $(GREEN)✓$(RESET) Rust 已安装: $$(rustc --version)"; \
		if command -v cargo >/dev/null 2>&1; then \
			echo "  $(GREEN)✓$(RESET) Cargo 已安装: $$(cargo --version)"; \
		fi; \
		if command -v tauri >/dev/null 2>&1; then \
			echo "  $(GREEN)✓$(RESET) Tauri CLI 已安装: $$(tauri --version)"; \
		else \
			echo "  $(YELLOW)⚠$(RESET) Tauri CLI 未安装"; \
			echo "  安装方法: cargo install tauri-cli"; \
		fi; \
	else \
		echo "  $(RED)✗$(RESET) Rust 未安装"; \
		echo "  安装方法: curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh"; \
	fi
	@echo ""

check-android: ## 检查 Android 开发环境
	@echo "$(YELLOW)🤖 Android 环境:$(RESET)"
	@if [ -n "$$ANDROID_HOME" ]; then \
		echo "  $(GREEN)✓$(RESET) ANDROID_HOME 已设置: $$ANDROID_HOME"; \
		if [ -d "$$ANDROID_HOME" ]; then \
			echo "  $(GREEN)✓$(RESET) Android SDK 目录存在"; \
		else \
			echo "  $(RED)✗$(RESET) Android SDK 目录不存在"; \
		fi; \
	else \
		echo "  $(YELLOW)⚠$(RESET) ANDROID_HOME 未设置"; \
		echo "  请安装 Android Studio 并设置环境变量"; \
	fi
	@if command -v adb >/dev/null 2>&1; then \
		echo "  $(GREEN)✓$(RESET) ADB 已安装: $$(adb --version | head -n1)"; \
	else \
		echo "  $(YELLOW)⚠$(RESET) ADB 未安装或不在 PATH 中"; \
	fi
	@if command -v java >/dev/null 2>&1; then \
		echo "  $(GREEN)✓$(RESET) Java 已安装: $$(java -version 2>&1 | head -n1)"; \
	else \
		echo "  $(RED)✗$(RESET) Java 未安装"; \
	fi
	@echo ""

check-docker: ## 检查 Docker 环境
	@echo "$(YELLOW)🐳 Docker 环境:$(RESET)"
	@if command -v docker >/dev/null 2>&1; then \
		echo "  $(GREEN)✓$(RESET) Docker 已安装: $$(docker --version)"; \
		if docker info >/dev/null 2>&1; then \
			echo "  $(GREEN)✓$(RESET) Docker 服务运行正常"; \
		else \
			echo "  $(YELLOW)⚠$(RESET) Docker 服务未运行"; \
		fi; \
	else \
		echo "  $(YELLOW)⚠$(RESET) Docker 未安装"; \
		echo "  安装方法:"; \
		echo "    macOS: brew install --cask docker"; \
		echo "    Ubuntu: sudo apt install docker.io"; \
		echo "    Windows: winget install Docker.DockerDesktop"; \
	fi
	@if command -v docker-compose >/dev/null 2>&1; then \
		echo "  $(GREEN)✓$(RESET) Docker Compose 已安装: $$(docker-compose --version)"; \
	else \
		echo "  $(YELLOW)⚠$(RESET) Docker Compose 未安装"; \
	fi
	@echo ""

check-git: ## 检查 Git 环境
	@echo "$(YELLOW)📝 Git 环境:$(RESET)"
	@if command -v git >/dev/null 2>&1; then \
		echo "  $(GREEN)✓$(RESET) Git 已安装: $$(git --version)"; \
		if git config user.name >/dev/null 2>&1; then \
			echo "  $(GREEN)✓$(RESET) Git 用户名已配置: $$(git config user.name)"; \
		else \
			echo "  $(YELLOW)⚠$(RESET) Git 用户名未配置"; \
		fi; \
		if git config user.email >/dev/null 2>&1; then \
			echo "  $(GREEN)✓$(RESET) Git 邮箱已配置: $$(git config user.email)"; \
		else \
			echo "  $(YELLOW)⚠$(RESET) Git 邮箱未配置"; \
		fi; \
	else \
		echo "  $(RED)✗$(RESET) Git 未安装"; \
	fi

##@ 依赖管理

deps: deps-core deps-desktop deps-android deps-cli ## 安装所有项目依赖

deps-core: ## 安装 Core 依赖
	@echo "$(BLUE)📦 安装 Core 依赖...$(RESET)"
	@cd core && go mod tidy && go mod download

deps-desktop: ## 安装桌面客户端依赖
	@echo "$(BLUE)📦 安装桌面客户端依赖...$(RESET)"
	@cd desktop && npm install

deps-android: ## 安装 Android 依赖
	@echo "$(BLUE)📦 安装 Android 依赖...$(RESET)"
	@cd android && ./gradlew build --refresh-dependencies

deps-cli: ## 安装 CLI 依赖
	@echo "$(BLUE)📦 安装 CLI 依赖...$(RESET)"
	@cd cli && go mod tidy && go mod download

##@ 开发

dev: ## 启动开发环境
	@echo "$(BLUE)🚀 启动开发环境...$(RESET)"
	@$(MAKE) dev-core &
	@sleep 3
	@$(MAKE) dev-desktop

dev-core: ## 启动 Core 开发服务器
	@echo "$(BLUE)🎯 启动 Core 开发服务器...$(RESET)"
	@cd core && go run cmd/prism-core/main.go

dev-desktop: ## 启动桌面客户端开发服务器
	@echo "$(BLUE)🖥️ 启动桌面客户端开发服务器...$(RESET)"
	@cd desktop && npm run tauri dev

dev-android: ## 构建 Android 开发版本
	@echo "$(BLUE)📱 构建 Android 开发版本...$(RESET)"
	@cd android && ./gradlew assembleDebug

##@ 构建

build: build-core build-desktop build-android build-cli ## 构建所有组件

build-core: ## 构建 Core
	@echo "$(BLUE)🔨 构建 Core...$(RESET)"
	@mkdir -p dist
	@cd core && go build $(LDFLAGS) -o ../dist/prism-core cmd/prism-core/main.go

build-desktop: ## 构建桌面客户端
	@echo "$(BLUE)🔨 构建桌面客户端...$(RESET)"
	@cd desktop && npm run tauri build

build-android: ## 构建 Android 应用
	@echo "$(BLUE)🔨 构建 Android 应用...$(RESET)"
	@cd android && ./gradlew assembleRelease

build-cli: ## 构建 CLI 工具
	@echo "$(BLUE)🔨 构建 CLI 工具...$(RESET)"
	@mkdir -p dist
	@cd cli && go build $(LDFLAGS) -o ../dist/prism-cli main.go

build-all: ## 构建所有平台版本
	@echo "$(BLUE)🔨 构建所有平台版本...$(RESET)"
	@mkdir -p dist
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*}; \
		GOARCH=$${platform#*/}; \
		echo "构建 $$GOOS/$$GOARCH..."; \
		cd core && env GOOS=$$GOOS GOARCH=$$GOARCH go build $(LDFLAGS) -o ../dist/prism-core-$$GOOS-$$GOARCH cmd/prism-core/main.go; \
		cd cli && env GOOS=$$GOOS GOARCH=$$GOARCH go build $(LDFLAGS) -o ../dist/prism-cli-$$GOOS-$$GOARCH main.go; \
		if [ "$$GOOS" = "windows" ]; then \
			mv ../dist/prism-core-$$GOOS-$$GOARCH ../dist/prism-core-$$GOOS-$$GOARCH.exe; \
			mv ../dist/prism-cli-$$GOOS-$$GOARCH ../dist/prism-cli-$$GOOS-$$GOARCH.exe; \
		fi; \
		cd ..; \
	done

##@ 测试

test: test-core test-desktop test-android test-cli ## 运行所有测试

test-core: ## 运行 Core 测试
	@echo "$(BLUE)🧪 运行 Core 测试...$(RESET)"
	@cd core && go test -v -race -coverprofile=coverage.out ./...

test-desktop: ## 运行桌面客户端测试
	@echo "$(BLUE)🧪 运行桌面客户端测试...$(RESET)"
	@cd desktop && npm test

test-android: ## 运行 Android 测试
	@echo "$(BLUE)🧪 运行 Android 测试...$(RESET)"
	@cd android && ./gradlew test

test-cli: ## 运行 CLI 测试
	@echo "$(BLUE)🧪 运行 CLI 测试...$(RESET)"
	@cd cli && go test -v ./...

##@ 代码质量

lint: lint-core lint-desktop lint-android lint-cli ## 运行所有代码检查

lint-core: ## 检查 Core 代码质量
	@echo "$(BLUE)🔍 检查 Core 代码质量...$(RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		cd core && golangci-lint run; \
	else \
		echo "$(YELLOW)⚠$(RESET) golangci-lint 未安装，跳过检查"; \
		echo "安装方法: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

lint-desktop: ## 检查桌面客户端代码质量
	@echo "$(BLUE)🔍 检查桌面客户端代码质量...$(RESET)"
	@cd desktop && npm run lint

lint-android: ## 检查 Android 代码质量
	@echo "$(BLUE)🔍 检查 Android 代码质量...$(RESET)"
	@cd android && ./gradlew ktlintCheck

lint-cli: ## 检查 CLI 代码质量
	@echo "$(BLUE)🔍 检查 CLI 代码质量...$(RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		cd cli && golangci-lint run; \
	else \
		echo "$(YELLOW)⚠$(RESET) golangci-lint 未安装，跳过检查"; \
	fi

fmt: ## 格式化所有代码
	@echo "$(BLUE)✨ 格式化代码...$(RESET)"
	@cd core && go fmt ./...
	@cd cli && go fmt ./...
	@cd desktop && npm run format
	@cd android && ./gradlew ktlintFormat

##@ Docker

docker-build: ## 构建 Docker 镜像
	@echo "$(BLUE)🐳 构建 Docker 镜像...$(RESET)"
	@docker build -t prism/core:$(VERSION) -f core/Dockerfile .

docker-run: ## 运行 Docker 容器
	@echo "$(BLUE)🐳 运行 Docker 容器...$(RESET)"
	@docker run -d -p 9090:9090 -p 7890:7890 --name prism-core prism/core:$(VERSION)

docker-compose-up: ## 启动 Docker Compose
	@echo "$(BLUE)🐳 启动 Docker Compose...$(RESET)"
	@docker-compose up -d

docker-compose-down: ## 停止 Docker Compose
	@echo "$(BLUE)🐳 停止 Docker Compose...$(RESET)"
	@docker-compose down

##@ 维护

clean: ## 清理构建文件
	@echo "$(BLUE)🧹 清理构建文件...$(RESET)"
	@rm -rf dist/
	@cd core && go clean
	@cd cli && go clean
	@cd desktop && rm -rf dist/ target/
	@cd android && ./gradlew clean

install: ## 安装到系统
	@echo "$(BLUE)📦 安装到系统...$(RESET)"
	@sudo cp dist/prism-core /usr/local/bin/
	@sudo cp dist/prism-cli /usr/local/bin/
	@echo "$(GREEN)✅ 安装完成！$(RESET)"

uninstall: ## 从系统卸载
	@echo "$(BLUE)🗑️ 从系统卸载...$(RESET)"
	@sudo rm -f /usr/local/bin/prism-core
	@sudo rm -f /usr/local/bin/prism-cli
	@echo "$(GREEN)✅ 卸载完成！$(RESET)"

##@ 工具

setup-dev: ## 设置开发环境
	@echo "$(BLUE)⚙️ 设置开发环境...$(RESET)"
	@$(MAKE) doctor
	@$(MAKE) deps
	@echo "$(GREEN)✅ 开发环境设置完成！$(RESET)"

release: ## 创建发布版本
	@echo "$(BLUE)🚀 创建发布版本 $(VERSION)...$(RESET)"
	@$(MAKE) clean
	@$(MAKE) test
	@$(MAKE) build-all
	@echo "$(GREEN)✅ 发布版本 $(VERSION) 构建完成！$(RESET)"

##@ 帮助

help: ## 显示帮助信息
	@awk 'BEGIN {FS = ":.*##"; printf "\n$(BLUE)Prism 构建工具$(RESET)\n\n使用方法:\n  make $(YELLOW)<target>$(RESET)\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  $(YELLOW)%-15s$(RESET) %s\n", $$1, $$2 } /^##@/ { printf "\n$(BLUE)%s$(RESET)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)