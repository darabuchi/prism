# Node - 代理节点管理

## 概述

`pkg/node` 提供代理节点的核心数据结构和操作方法，是 Prism 项目中节点管理的基础包。

## 特性

- **节点抽象**：统一的节点数据结构，支持所有代理协议
- **唯一标识**：基于 SHA256 的节点 ID 生成
- **配置编解码**：Base64 JSON 配置的编码和解码
- **状态管理**：节点存活、启用、熔断器状态
- **测试结果**：延迟、速度、解锁、地理信息等测试结果
- **评分系统**：节点综合评分计算
- **独立性**：无内部依赖，可被外部项目引用

## 核心概念

### 节点唯一标识

每个节点有两个标识符：

1. **ID**：基于完整配置的 SHA256 哈希（64位十六进制字符串）
   - 唯一标识节点配置
   - 用于数据库主键和跨订阅去重

2. **UniqueKey**：格式化的唯一键（如 `vmess://example.com:443`）
   - 便于人类阅读
   - 用于快速查找和显示

### 节点状态

节点维护多个状态：

- **Alive**: 是否存活（通过延迟测试）
- **Enabled**: 是否启用（用户控制）
- **HealthState**: 熔断器状态（Open/HalfOpen/Closed）

### 测试结果

节点支持多维度测试：

- **延迟测试**: 平均/最小/最大延迟（ms）
- **速度测试**: 下载/上传速度（bytes/s）
- **地理信息**: 入口/出口国家、IP、ASN、ISP
- **解锁检测**: OpenAI、Netflix、YouTube、Disney+ 解锁状态
- **线路质量**: 是否为优质线路

## 数据结构

### Node

节点核心结构：

```go
type Node struct {
    id        string         // 唯一标识（SHA256）
    uniqueKey string         // 唯一键
    config    map[string]any // 完整配置
    info      *Info          // 节点信息
}
```

### Config

节点配置信息：

```go
type Config struct {
    Name    string         // 节点名称
    Type    string         // 协议类型
    Server  string         // 服务器地址
    Port    int            // 端口
    Config  map[string]any // 完整配置
    RawJSON string         // Base64 编码的原始 JSON 配置
}
```

### Status

节点运行时状态：

```go
type Status struct {
    Alive       bool        // 是否存活
    Enabled     bool        // 是否启用
    HealthState HealthState // 熔断器状态

    // 熔断器统计
    SuccessCount uint64 // 成功请求数
    FailureCount uint64 // 失败请求数

    // 流量统计
    UploadTotal   int64 // 总上传流量（字节）
    DownloadTotal int64 // 总下载流量（字节）
}
```

### TestResult

测试结果：

```go
type TestResult struct {
    // 延迟测试
    DelayAvg       int   // 平均延迟(ms)
    DelayMin       int   // 最小延迟(ms)
    DelayMax       int   // 最大延迟(ms)
    DelayTestedAt  int64 // 延迟测试时间

    // 速度测试
    DownloadSpeed      int64 // 下载速度(bytes/s)
    UploadSpeed        int64 // 上载速度(bytes/s)

    // 地理信息
    InboundCountry  string // 入口国家代码
    OutboundCountry string // 出口国家代码

    // 解锁检测 (0:未测试 1:失败 2:部分解锁 3:完全解锁)
    UnlockOpenAI   int
    UnlockNetflix  int
    UnlockYouTube  int
    UnlockDisney   int

    // 线路质量
    HighQualityRoutes bool     // 是否为优质线路
    RouterList        []string // 路由列表
}
```

### Info

节点完整信息：

```go
type Info struct {
    ID        string     // 唯一标识
    UniqueKey string     // 唯一键
    Config    Config     // 配置信息
    Status    Status     // 运行状态
    Test      TestResult // 测试结果
    Score     float64    // 综合评分

    // 时间戳
    FirstAliveAt int64 // 首次存活时间
    LastAliveAt  int64 // 最后存活时间
    DeathCount   int64 // 死亡计数
    CreatedAt    int64 // 创建时间
    UpdatedAt    int64 // 更新时间
}
```

## 使用方式

### 创建节点

```go
package main

import (
    "github.com/darabuchi/prism/pkg/node"
    "github.com/lazygophers/log"
)

func main() {
    // 从配置创建节点
    config := map[string]any{
        "name":   "香港节点",
        "type":   "vmess",
        "server": "hk.example.com",
        "port":   443,
        "uuid":   "xxx-xxx-xxx",
        "alterId": 0,
    }

    n, err := node.New(config)
    if err != nil {
        log.Errorf("创建节点失败: %v", err)
        return
    }

    log.Infof("节点 ID: %s", n.ID())
    log.Infof("节点名称: %s", n.Name())
    log.Infof("节点地址: %s", n.Address())
}
```

### 从编码配置创建

```go
// 从 Base64 编码的配置创建节点
encoded := "eyJuYW1lIjoi6aaZ6Kef6IqC54K5IiwidHlwZSI6InZtZXNzIi...=="

n, err := node.NewFromEncoded(encoded)
if err != nil {
    log.Errorf("解码失败: %v", err)
    return
}

log.Infof("节点: %s", n.Name())
```

### 节点状态管理

```go
// 设置节点存活状态
n.SetAlive(true)

// 设置节点启用状态
n.SetEnabled(true)

// 设置熔断器状态
n.SetHealthState(node.HealthStateClosed)

// 检查节点状态
if n.IsAlive() && n.IsEnabled() {
    log.Info("节点可用")
}

// 更新完整状态
status := node.Status{
    Alive:        true,
    Enabled:      true,
    HealthState:  node.HealthStateClosed,
    SuccessCount: 100,
    FailureCount: 5,
}
n.UpdateStatus(status)
```

### 测试结果更新

```go
// 更新延迟测试结果
test := node.TestResult{
    DelayAvg:      50,
    DelayMin:      45,
    DelayMax:      60,
    DelayTestedAt: time.Now().Unix(),
}
n.UpdateTest(test)

// 更新速度测试结果
test.DownloadSpeed = 10 * 1024 * 1024 // 10 MB/s
test.UploadSpeed = 5 * 1024 * 1024    // 5 MB/s
test.SpeedTestedAt = time.Now().Unix()
n.UpdateTest(test)

// 更新解锁检测结果
test.UnlockOpenAI = 3  // 完全解锁
test.UnlockNetflix = 3 // 完全解锁
test.UnlockYouTube = 2 // 部分解锁
test.UnlockDisney = 1  // 失败
n.UpdateTest(test)
```

### 节点评分

```go
// 设置综合评分
n.SetScore(85.5)

// 获取节点信息（包含评分）
info := n.Info()
log.Infof("节点评分: %.2f", info.Score)
```

### 节点克隆

```go
// 克隆节点（用于测试或备份）
cloned := n.Clone()
log.Infof("克隆节点: %s", cloned.Name())
```

### ID 生成和配置编码

```go
import "github.com/darabuchi/prism/pkg/node"

// 生成节点 ID
id, err := node.GenerateID(config)
if err != nil {
    log.Errorf("生成 ID 失败: %v", err)
}
log.Infof("节点 ID: %s", id)

// 生成唯一键
uniqueKey := node.GenerateUniqueKey(config)
log.Infof("唯一键: %s", uniqueKey)

// 编码配置
encoded, err := node.EncodeConfig(config)
if err != nil {
    log.Errorf("编码失败: %v", err)
}
log.Infof("编码配置: %s", encoded)

// 解码配置
decoded, err := node.DecodeConfig(encoded)
if err != nil {
    log.Errorf("解码失败: %v", err)
}
```

## 健康状态说明

### HealthState 状态机

```
      失败次数超过阈值
Closed ──────────────→ Open
  ↑                      │
  │                      │ 超时后进入半开
  │                      ↓
  └──────────────── HalfOpen
      请求成功
```

- **Closed（闭路）**: 正常工作，接受所有请求
- **Open（开路）**: 熔断器开启，拒绝所有请求
- **HalfOpen（半开）**: 尝试恢复，允许少量请求测试

## 协议类型

支持的代理协议（18+种）：

| 协议 | 类型字符串 | 说明 |
|------|-----------|------|
| Shadowsocks | `ss` | Shadowsocks 代理 |
| ShadowsocksR | `ssr` | ShadowsocksR 代理 |
| VMess | `vmess` | V2Ray VMess 协议 |
| VLess | `vless` | V2Ray VLess 协议 |
| Trojan | `trojan` | Trojan 代理协议 |
| Hysteria | `hysteria` | Hysteria 协议 |
| Hysteria2 | `hysteria2` | Hysteria2 协议 |
| SOCKS5 | `socks5` | SOCKS5 代理 |
| HTTP | `http` | HTTP(S) 代理 |
| Snell | `snell` | Snell 协议 |
| WireGuard | `wireguard` | WireGuard VPN |
| TUIC | `tuic` | TUIC 协议 |
| SSH | `ssh` | SSH 隧道 |
| Mieru | `mieru` | Mieru 协议 |
| AnyTLS | `anytls` | AnyTLS 协议 |
| Direct | `direct` | 直连 |
| Reject | `reject` | 拒绝连接 |
| DNS | `dns` | DNS 查询 |

## 设计原则

1. **独立性**: 不依赖 internal 包，可被外部项目引用
2. **不可变性**: 核心字段（ID、UniqueKey）不可修改
3. **线程安全**: 提供克隆方法支持并发场景
4. **可扩展性**: 使用 map[string]any 支持任意配置
5. **可序列化**: 所有字段支持 JSON 序列化

## API 文档

### 节点创建

- `New(config map[string]any) (*Node, error)` - 从配置创建节点
- `NewFromEncoded(encoded string) (*Node, error)` - 从编码配置创建节点

### 节点信息

- `ID() string` - 获取节点唯一标识
- `UniqueKey() string` - 获取节点唯一键
- `Config() map[string]any` - 获取节点配置
- `Info() *Info` - 获取节点完整信息
- `Name() string` - 获取节点名称
- `Type() string` - 获取节点类型
- `Server() string` - 获取服务器地址
- `Port() int` - 获取端口
- `Address() string` - 获取完整地址

### 状态管理

- `IsAlive() bool` - 是否存活
- `IsEnabled() bool` - 是否启用
- `SetAlive(alive bool)` - 设置存活状态
- `SetEnabled(enabled bool)` - 设置启用状态
- `SetHealthState(state HealthState)` - 设置熔断器状态
- `UpdateStatus(status Status)` - 更新完整状态

### 测试结果

- `UpdateTest(test TestResult)` - 更新测试结果
- `SetScore(score float64)` - 设置评分

### 工具函数

- `GenerateID(config map[string]any) (string, error)` - 生成节点 ID
- `GenerateUniqueKey(config map[string]any) string` - 生成唯一键
- `EncodeConfig(config map[string]any) (string, error)` - 编码配置
- `DecodeConfig(encoded string) (map[string]any, error)` - 解码配置

### 其他

- `Clone() *Node` - 克隆节点

## 依赖

- `crypto/sha256` - SHA256 哈希计算
- `encoding/base64` - Base64 编解码
- `github.com/lazygophers/utils/json` - JSON 处理

## 测试

TODO: 添加单元测试

```go
func TestNode(t *testing.T) {
    // 测试节点创建
    // 测试 ID 生成
    // 测试配置编解码
    // 测试状态管理
}
```

## 相关文档

- [Parser 订阅解析器](../parser/README.md)
- [产品设计文档](../../docs/产品设计文档.md)
- [系统设计文档](../../docs/系统设计文档.md)
