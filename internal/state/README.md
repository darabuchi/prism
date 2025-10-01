# 状态管理模块

## 说明

此模块负责全局状态管理，包括配置、数据库连接、日志初始化等。

## 文件

- `state.go`: 全局状态定义
- `config.go`: 配置加载
- `table.go`: 数据库表管理
- `logger.go`: 日志初始化

## 使用示例

```go
// 访问配置
port := state.State.Config.Server.HttpPort

// 访问数据库
db := state.Db()

// 访问 Model
state.Subscription.Create(sub)
state.Node.First("id = ?", nodeID)
```
