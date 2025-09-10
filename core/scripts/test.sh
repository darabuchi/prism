#!/bin/bash

# Prism Core 测试脚本
# 用于自动化运行各种测试场景

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 获取脚本目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 显示帮助信息
show_help() {
    echo "Prism Core 测试脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help       显示帮助信息"
    echo "  -u, --unit       运行单元测试"
    echo "  -i, --integration 运行集成测试"
    echo "  -c, --coverage   生成覆盖率报告"
    echo "  -l, --lint       运行代码检查"
    echo "  -f, --fmt        格式化代码"
    echo "  -a, --all        运行所有检查（默认）"
    echo "  -v, --verbose    详细输出"
    echo "  -q, --quick      快速测试（跳过长时间运行的测试）"
    echo "  --ci             CI 模式"
    echo "  --clean          清理测试文件"
    echo ""
    echo "示例:"
    echo "  $0                   # 运行所有测试"
    echo "  $0 -u -v           # 详细模式运行单元测试"
    echo "  $0 -c              # 生成覆盖率报告"
    echo "  $0 --clean         # 清理测试文件"
}

# 检查 Go 环境
check_go() {
    log_info "检查 Go 环境..."
    if ! command -v go &> /dev/null; then
        log_error "Go 未安装或不在 PATH 中"
        exit 1
    fi
    
    GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
    log_success "Go 版本: $GO_VERSION"
}

# 检查项目依赖
check_dependencies() {
    log_info "检查项目依赖..."
    cd "$PROJECT_DIR"
    if [ ! -f "go.mod" ]; then
        log_error "未找到 go.mod 文件"
        exit 1
    fi
    
    log_info "下载依赖..."
    go mod download
    log_success "依赖检查完成"
}

# 运行单元测试
run_unit_tests() {
    log_info "运行单元测试..."
    cd "$PROJECT_DIR"
    
    if [ "$VERBOSE" = true ]; then
        go test -v -race -short ./internal/service/... ./internal/storage/... ./internal/testutil/...
    else
        go test -race -short ./internal/service/... ./internal/storage/... ./internal/testutil/...
    fi
    
    log_success "单元测试完成"
}

# 运行集成测试
run_integration_tests() {
    log_info "运行集成测试..."
    cd "$PROJECT_DIR"
    
    if [ "$VERBOSE" = true ]; then
        go test -v -race ./internal/api/...
    else
        go test -race ./internal/api/...
    fi
    
    log_success "集成测试完成"
}

# 生成覆盖率报告
generate_coverage() {
    log_info "生成覆盖率报告..."
    cd "$PROJECT_DIR"
    
    mkdir -p coverage
    COVERAGE_FILE="coverage/coverage.out"
    COVERAGE_HTML="coverage/coverage.html"
    
    go test -race -coverprofile="$COVERAGE_FILE" ./...
    
    if [ -f "$COVERAGE_FILE" ]; then
        # 生成 HTML 报告
        go tool cover -html="$COVERAGE_FILE" -o "$COVERAGE_HTML"
        
        # 显示总覆盖率
        TOTAL_COVERAGE=$(go tool cover -func="$COVERAGE_FILE" | grep total | awk '{print $3}')
        log_success "总覆盖率: $TOTAL_COVERAGE"
        log_success "HTML 报告: $COVERAGE_HTML"
        
        # 如果不是 CI 环境，尝试打开浏览器
        if [ "$CI_MODE" != true ] && command -v open &> /dev/null; then
            log_info "尝试在浏览器中打开覆盖率报告..."
            open "$COVERAGE_HTML" || true
        fi
    else
        log_error "覆盖率文件生成失败"
        exit 1
    fi
}

# 运行代码检查
run_lint() {
    log_info "运行代码检查..."
    cd "$PROJECT_DIR"
    
    # 检查是否安装了 golangci-lint
    if command -v golangci-lint &> /dev/null; then
        golangci-lint run ./...
    else
        log_warning "golangci-lint 未安装，使用 go vet 替代"
        go vet ./...
    fi
    
    log_success "代码检查完成"
}

# 格式化代码
format_code() {
    log_info "格式化代码..."
    cd "$PROJECT_DIR"
    
    go fmt ./...
    
    # 检查是否有未格式化的文件
    UNFORMATTED=$(go fmt ./...)
    if [ -n "$UNFORMATTED" ]; then
        log_warning "以下文件已被格式化:"
        echo "$UNFORMATTED"
    fi
    
    log_success "代码格式化完成"
}

# 清理测试文件
clean_test_files() {
    log_info "清理测试文件..."
    cd "$PROJECT_DIR"
    
    rm -rf coverage/
    find . -name "*.test" -delete
    
    log_success "测试文件清理完成"
}

# 运行所有测试
run_all_tests() {
    log_info "运行所有测试和检查..."
    
    format_code
    run_lint
    run_unit_tests
    run_integration_tests
    
    if [ "$QUICK_MODE" != true ]; then
        generate_coverage
    fi
    
    log_success "所有测试和检查完成！"
}

# 解析命令行参数
UNIT_TESTS=false
INTEGRATION_TESTS=false
COVERAGE=false
LINT=false
FORMAT=false
ALL_TESTS=false
VERBOSE=false
QUICK_MODE=false
CI_MODE=false
CLEAN=false

# 如果没有参数，默认运行所有测试
if [ $# -eq 0 ]; then
    ALL_TESTS=true
fi

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -u|--unit)
            UNIT_TESTS=true
            shift
            ;;
        -i|--integration)
            INTEGRATION_TESTS=true
            shift
            ;;
        -c|--coverage)
            COVERAGE=true
            shift
            ;;
        -l|--lint)
            LINT=true
            shift
            ;;
        -f|--fmt)
            FORMAT=true
            shift
            ;;
        -a|--all)
            ALL_TESTS=true
            shift
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -q|--quick)
            QUICK_MODE=true
            shift
            ;;
        --ci)
            CI_MODE=true
            VERBOSE=true
            shift
            ;;
        --clean)
            CLEAN=true
            shift
            ;;
        *)
            log_error "未知选项: $1"
            show_help
            exit 1
            ;;
    esac
done

# 主逻辑
main() {
    log_info "Prism Core 测试脚本开始运行..."
    
    # 基本环境检查
    check_go
    check_dependencies
    
    # 根据参数执行相应操作
    if [ "$CLEAN" = true ]; then
        clean_test_files
        exit 0
    fi
    
    if [ "$ALL_TESTS" = true ]; then
        run_all_tests
    else
        if [ "$FORMAT" = true ]; then
            format_code
        fi
        
        if [ "$LINT" = true ]; then
            run_lint
        fi
        
        if [ "$UNIT_TESTS" = true ]; then
            run_unit_tests
        fi
        
        if [ "$INTEGRATION_TESTS" = true ]; then
            run_integration_tests
        fi
        
        if [ "$COVERAGE" = true ]; then
            generate_coverage
        fi
    fi
    
    log_success "测试脚本执行完成！"
}

# 错误处理
trap 'log_error "脚本执行失败，退出码: $?"' ERR

# 运行主函数
main