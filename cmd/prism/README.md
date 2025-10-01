# Prism 主程序

## 说明

此目录包含 Prism 应用程序的主入口点。

## 文件

- `main.go`: 程序入口，负责初始化和启动服务

## 职责

1. 加载配置文件
2. 初始化日志系统
3. 连接数据库
4. 启动 HTTP 服务器
5. 启动代理服务器
6. 处理信号和优雅关闭

## 示例

```go
package main

import (
	"github.com/ice-cream-heaven/prism/internal/state"
	"github.com/lazygophers/log"
)

func main() {
	// 加载配置
	if err := state.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化服务
	// ...
}
```
