# 命令行入口

## 说明

此目录包含应用程序的入口文件（只包含 Go 文件，无子目录）。

## 文件

### server.go
主服务器程序入口，启动以下服务：
- HTTP API 服务器
- SOCKS5 代理服务器
- 订阅管理服务
- 节点测试服务

### cli.go（可选）
命令行工具，用于：
- 管理订阅
- 测试节点
- 查看配置
- 数据库操作

## 构建

```bash
# 构建服务器
make build-server
# 或
go build -o bin/prism-server ./cmd/server.go

# 构建命令行工具
make build-cli
# 或
go build -o bin/prism-cli ./cmd/cli.go
```

## 运行

```bash
# 运行服务器
./bin/prism-server

# 使用命令行工具
./bin/prism-cli --help
```

## 设计原则

- 每个文件对应一个可执行程序
- 保持文件简洁，主要逻辑放在 internal 包中
- 只负责程序初始化和启动
- 通过命令行参数或配置文件进行配置
