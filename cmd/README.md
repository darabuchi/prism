# 命令行入口

## 说明

此目录包含应用程序的入口文件 `main.go`，使用 cobra 实现子命令架构。

## 文件

### main.go
主程序入口，包含所有子命令：

**server 命令** - 启动服务器
```bash
prism server [flags]
  --config string   配置文件路径 (default "./config.yaml")
  --port int        HTTP 端口 (default 8080)
  --daemon          后台运行
```

**subscription 命令** - 订阅管理
```bash
prism subscription list                    # 列出所有订阅
prism subscription add --url URL           # 添加订阅
prism subscription update --id ID          # 更新订阅
prism subscription remove --id ID          # 删除订阅
```

**node 命令** - 节点管理
```bash
prism node list                           # 列出所有节点
prism node test --id ID                   # 测试节点
```

**config 命令** - 配置管理
```bash
prism config show                         # 显示当前配置
prism config validate                     # 验证配置文件
```

**version 命令** - 版本信息
```bash
prism version                             # 显示版本信息
```

## 构建

```bash
# 构建
make build
# 或
go build -o bin/prism ./cmd/main.go
```

## 运行

```bash
# 启动服务器
./bin/prism server

# 使用自定义配置
./bin/prism server --config /path/to/config.yaml

# 管理订阅
./bin/prism subscription list
./bin/prism subscription add --url https://example.com/sub

# 测试节点
./bin/prism node test --id 123

# 查看配置
./bin/prism config show
```

## 设计原则

- 单文件入口，使用 cobra 管理所有子命令
- 保持文件简洁，主要逻辑放在 internal 包中
- 只负责程序初始化和命令分发
- 提供友好的命令行交互体验
- 支持丰富的命令行参数和标志
