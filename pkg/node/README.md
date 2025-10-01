# Node - 节点工具包

## 概述

`pkg/node` 提供节点相关的基础工具函数和常量定义，是 Prism 项目中节点管理的基础工具包。

## 设计原则

根据 Prism 项目的编码规范：
- **高度独立**：不依赖 internal 包，可被外部项目引用
- **职责单一**：只提供基础工具函数和常量定义，不包含业务逻辑
- **规范遵循**：严格遵循后端开发规范，使用指定的依赖包

## 功能特性

- **常量定义**：代理协议类型、熔断器健康状态
- **ID 生成**：基于 SHA256 的节点唯一标识生成
- **唯一键生成**：格式化的节点唯一键（type://server:port）
- **配置编解码**：Base64 JSON 格式的配置编码和解码

## 支持的代理协议

| 协议 | 类型常量 | 说明 |
|------|---------|------|
| Shadowsocks | `TypeShadowsocks` | Shadowsocks 代理 |
| ShadowsocksR | `TypeShadowsocksR` | ShadowsocksR 代理 |
| VMess | `TypeVMess` | V2Ray VMess 协议 |
| VLess | `TypeVLess` | V2Ray VLess 协议 |
| Trojan | `TypeTrojan` | Trojan 代理协议 |
| Hysteria | `TypeHysteria` | Hysteria 协议 |
| Hysteria2 | `TypeHysteria2` | Hysteria2 协议 |
| SOCKS5 | `TypeSocks5` | SOCKS5 代理 |
| HTTP | `TypeHTTP` | HTTP(S) 代理 |
| Snell | `TypeSnell` | Snell 协议 |
| WireGuard | `TypeWireGuard` | WireGuard VPN |
| TUIC | `TypeTuic` | TUIC 协议 |
| SSH | `TypeSSH` | SSH 隧道 |
| Mieru | `TypeMieru` | Mieru 协议 |
| AnyTLS | `TypeAnyTLS` | AnyTLS 协议 |
| Direct | `TypeDirect` | 直连 |
| Reject | `TypeReject` | 拒绝连接 |
| DNS | `TypeDNS` | DNS 查询 |

## 熔断器健康状态

```go
const (
    HealthStateOpen     HealthState = 0 // 开路（熔断器开启，拒绝请求）
    HealthStateHalfOpen HealthState = 1 // 半开（熔断器尝试恢复）
    HealthStateClosed   HealthState = 2 // 闭路（熔断器关闭，正常工作）
)
```

### 状态机

```
      失败次数超过阈值
Closed ──────────────→ Open
  ↑                      │
  │                      │ 超时后进入半开
  │                      ↓
  └──────────────── HalfOpen
      请求成功
```

## 使用方式

### 生成节点 ID

```go
import "github.com/darabuchi/prism/pkg/node"

config := map[string]any{
    "name":   "香港节点",
    "type":   "vmess",
    "server": "hk.example.com",
    "port":   443,
    "uuid":   "xxx-xxx-xxx",
}

// 生成唯一 ID（SHA256 哈希）
id, err := node.GenerateID(config)
if err != nil {
    log.Errorf("err:%v", err)
    return
}
log.Infof("节点 ID: %s", id) // 输出：64位十六进制字符串
```

### 生成唯一键

```go
// 生成格式化的唯一键
uniqueKey, err := node.GenerateUniqueKey(config)
if err != nil {
    log.Errorf("err:%v", err)
    return
}
log.Infof("唯一键: %s", uniqueKey) // 输出：vmess://hk.example.com:443
```

### 配置编解码

```go
// 编码配置为 Base64 JSON
encoded, err := node.EncodeConfig(config)
if err != nil {
    log.Errorf("err:%v", err)
    return
}
log.Infof("编码配置: %s", encoded)

// 解码配置
decoded, err := node.DecodeConfig(encoded)
if err != nil {
    log.Errorf("err:%v", err)
    return
}
log.Infof("节点类型: %s", decoded["type"])
```

### 使用代理类型常量

```go
import "github.com/darabuchi/prism/pkg/node"

// 检查协议类型
proxyType := node.TypeVMess
log.Infof("协议类型: %s", proxyType) // 输出：vmess

// 在配置中使用
config := map[string]any{
    "type": string(node.TypeVMess),
    // ...
}
```

### 使用健康状态

```go
import "github.com/darabuchi/prism/pkg/node"

state := node.HealthStateClosed
log.Infof("状态: %s", state.String()) // 输出：closed

switch state {
case node.HealthStateOpen:
    log.Info("熔断器开启，拒绝请求")
case node.HealthStateHalfOpen:
    log.Info("熔断器半开，尝试恢复")
case node.HealthStateClosed:
    log.Info("熔断器关闭，正常工作")
}
```

## API 文档

### 常量和类型

#### ProxyType

```go
type ProxyType string
```

代理协议类型，支持 18+ 种协议。

#### HealthState

```go
type HealthState int32
```

熔断器健康状态，支持 3 种状态。

**方法**：
- `String() string` - 返回状态的字符串表示

### 函数

#### GenerateID

```go
func GenerateID(config map[string]any) (string, error)
```

生成节点唯一 ID（基于配置的 SHA256 哈希）。

**参数**：
- `config`: 节点配置（必须是有效的 JSON 对象）

**返回**：
- `string`: 64位十六进制 SHA256 哈希字符串
- `error`: 错误信息

**可能的错误**：
- 配置序列化失败

#### GenerateUniqueKey

```go
func GenerateUniqueKey(config map[string]any) (string, error)
```

生成节点唯一键（格式：type://server:port）。

**参数**：
- `config`: 节点配置（必须包含 type、server、port 字段）

**返回**：
- `string`: 格式化的唯一键
- `error`: 错误信息

**可能的错误**：
- 缺少 type 字段
- 缺少 server 字段
- port 字段类型错误或值无效

#### EncodeConfig

```go
func EncodeConfig(config map[string]any) (string, error)
```

将配置编码为 Base64 JSON 字符串。

**参数**：
- `config`: 节点配置

**返回**：
- `string`: Base64 编码的 JSON 字符串
- `error`: 错误信息

**可能的错误**：
- 配置序列化失败

#### DecodeConfig

```go
func DecodeConfig(encoded string) (map[string]any, error)
```

从 Base64 JSON 字符串解码配置。

**参数**：
- `encoded`: Base64 编码的 JSON 字符串

**返回**：
- `map[string]any`: 解码后的配置
- `error`: 错误信息

**可能的错误**：
- Base64 解码失败
- JSON 反序列化失败

## 依赖包

本包严格遵循项目的后端开发规范，使用以下依赖：

```go
import (
    "github.com/lazygophers/log"                    // 日志处理
    "github.com/lazygophers/lrpc/middleware/xerror" // 错误处理
    "github.com/lazygophers/utils/cryptox"          // 加密工具（SHA256/Base64）
    "github.com/lazygophers/utils/json"             // JSON 处理
)
```

**禁止使用的包**：
- ❌ `encoding/json` - 使用 `github.com/lazygophers/utils/json`
- ❌ `encoding/base64` - 使用 `github.com/lazygophers/utils/cryptox`
- ❌ `crypto/sha256` - 使用 `github.com/lazygophers/utils/cryptox`
- ❌ `fmt.Errorf` - 使用 `github.com/lazygophers/lrpc/middleware/xerror`

## 错误处理

所有函数使用 `xerror` 包处理错误：

```go
import "github.com/lazygophers/lrpc/middleware/xerror"

id, err := node.GenerateID(config)
if err != nil {
    // xerror 会自动包装错误信息
    log.Errorf("err:%v", err)
    return err
}
```

## 日志记录

使用 `lazygophers/log` 记录日志：

```go
import "github.com/lazygophers/log"

// 错误日志
log.Errorf("err:%v", err)

// 信息日志
log.Infof("节点 ID: %s", id)

// 调试日志
log.Debugf("配置: %+v", config)
```

## 设计说明

### 为什么不包含 Node 结构体？

根据项目的架构设计和编码规范：
- `pkg/node` 只提供基础工具函数和常量
- 复杂的 Node 结构体、状态管理、业务逻辑应放在 `internal/` 中实现
- 这样可以保持 pkg 包的高度独立性和可重用性

### 为什么使用 map[string]any？

- 节点配置格式多样（18+ 种协议，每种协议参数不同）
- 使用 `map[string]any` 提供最大的灵活性
- 具体的配置结构应在 `internal/model` 中定义

### 如何使用节点数据？

节点的完整数据结构（包括状态、测试结果、评分等）应该：
1. 在 `internal/model` 中定义 Model 结构
2. 使用 GORM 标签，符合数据库设计规范
3. 在 `internal/service` 中实现业务逻辑
4. 使用 `pkg/node` 提供的工具函数处理基础操作

## 测试

TODO: 添加单元测试

```go
func TestGenerateID(t *testing.T) {
    config := map[string]any{
        "type":   "vmess",
        "server": "example.com",
        "port":   443,
    }

    id, err := node.GenerateID(config)
    if err != nil {
        t.Fatalf("GenerateID failed: %v", err)
    }

    if len(id) != 64 { // SHA256 produces 64 hex chars
        t.Errorf("expected 64 chars, got %d", len(id))
    }
}
```

## 相关文档

- [后端开发规范](../../docs/编码规范/后端开发规范.md)
- [数据库设计规范](../../docs/编码规范/数据库设计规范.md)
- [Parser 订阅解析器](../parser/README.md)
- [系统设计文档](../../docs/系统设计文档.md)
