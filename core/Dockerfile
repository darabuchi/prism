# 多阶段构建
# 阶段 1: 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git ca-certificates tzdata

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X main.version=${VERSION:-dev} -X main.buildTime=$(date -u '+%Y-%m-%d_%H:%M:%S') -X main.gitCommit=${GIT_COMMIT:-unknown} -w -s" \
    -o prism-core \
    cmd/prism-core/main.go

# 阶段 2: 运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 创建非 root 用户
RUN addgroup -g 1001 -S prism && \
    adduser -u 1001 -S prism -G prism

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/prism-core .

# 复制配置文件
COPY --from=builder /app/config ./config

# 创建数据目录
RUN mkdir -p data logs && \
    chown -R prism:prism /app

# 切换到非 root 用户
USER prism

# 暴露端口
EXPOSE 9090 7890

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:9090/api/v1/health || exit 1

# 设置环境变量
ENV PRISM_CONFIG_PATH=/app/config/config.yaml
ENV PRISM_DATA_PATH=/app/data
ENV PRISM_LOG_LEVEL=info

# 启动应用
CMD ["./prism-core"]