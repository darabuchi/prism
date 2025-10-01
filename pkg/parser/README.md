# Parser - 订阅解析器

## 概述

`pkg/parser` 提供通用的代理订阅解析功能，支持多种订阅格式的自动识别和解析。

## 特性

- **多格式支持**：自动识别并解析 Clash YAML、V2Ray Base64 链接、逐行 JSON 格式
- **智能解析**：按优先级尝试不同格式，直到成功解析
- **协议兼容**：支持 18+ 种代理协议（SS、SSR、VMess、VLess、Trojan、Hysteria 等）
- **节点创建**：自动创建 node.Node 实例，包含 Mihomo 适配器
- **容错处理**：跳过无效节点，继续解析其他节点

## 支持的格式

### 1. Clash YAML 格式

标准的 Clash 配置文件格式：

```yaml
proxies:
  - name: "节点1"
    type: ss
    server: example.com
    port: 443
    cipher: aes-256-gcm
    password: "password"
  - name: "节点2"
    type: vmess
    server: example.com
    port: 443
    uuid: "uuid"
    alterId: 0
```

### 2. V2Ray 格式

Base64 编码的代理链接，每行一个：

```
vmess://base64encodedconfig
ss://base64encodedconfig
trojan://password@host:port
```

### 3. 逐行 JSON 格式

每行一个 JSON 对象，可选 Clash YAML 的 `-` 前缀：

```json
{"name":"节点1","type":"ss","server":"example.com","port":443}
- {"name":"节点2","type":"vmess","server":"example.com","port":443}
```

## 支持的代理协议

| 协议 | 类型 | 说明 |
|------|------|------|
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

## 使用方式

### 基本用法

```go
package main

import (
	"github.com/darabuchi/prism/pkg/parser"
	"github.com/lazygophers/log"
)

func main() {
	// 从订阅 URL 获取的数据
	subscriptionData := []byte(`
proxies:
  - name: "香港节点"
    type: ss
    server: hk.example.com
    port: 443
    cipher: aes-256-gcm
    password: "password123"
`)

	// 解析订阅
	nodes, err := parser.Parse(subscriptionData)
	if err != nil {
		log.Errorf("解析失败: %v", err)
		return
	}

	log.Infof("解析到 %d 个节点", len(nodes))
	for _, node := range nodes {
		log.Infof("节点: %s", node.Name())
	}
}
```

### 与 Mihomo Adapter 集成

```go
import (
	"context"
	"github.com/darabuchi/prism/pkg/parser"
	"github.com/metacubex/mihomo/constant"
)

// 解析订阅并使用节点
func ParseAndUseNodes(data []byte) error {
	// 解析订阅（自动创建 mihomo adapter）
	nodes, err := parser.Parse(data)
	if err != nil {
		return err
	}

	// 直接使用节点进行代理连接
	for _, node := range nodes {
		log.Infof("节点: %s, 类型: %s, 地址: %s:%d",
			node.Name(), node.ProxyType(), node.Server(), node.Port())

		// 使用节点建立连接
		metadata := &constant.Metadata{
			Host:    "example.com",
			DstPort: 443,
			NetWork: constant.TCP,
		}
		conn, err := node.DialContext(context.Background(), metadata)
		if err != nil {
			log.Warnf("连接失败: %v", err)
			continue
		}
		conn.Close()
	}

	return nil
}
```

## 解析流程

```
输入数据
    ↓
尝试 Clash YAML 解析
    ↓ (失败)
尝试 V2Ray Base64 解析
    ↓ (失败)
尝试逐行 JSON 解析
    ↓
创建 node.Node 实例
    ↓
按唯一 ID 去重
    ↓
返回 []*node.Node
```

## API 文档

### Parse

```go
func Parse(data []byte) ([]*node.Node, error)
```

解析订阅数据，自动识别格式并创建节点实例。

**参数**：
- `data`: 订阅内容字节数组

**返回**：
- `[]*node.Node`: 节点列表（已自动创建 mihomo adapter，已去重）
- `error`: 解析错误（仅当所有节点创建失败时返回错误）

**特性**：
- 自动识别格式（Clash YAML、V2Ray Base64、逐行 JSON）
- 自动创建 mihomo 代理适配器
- 按唯一 ID 去重
- 跳过无效节点，记录警告但继续处理其他节点

## 设计原则

1. **简单性**：单一入口 `Parse()` 函数，自动识别格式
2. **兼容性**：基于 mihomo 的转换器，保证协议兼容性
3. **健壮性**：多格式回退机制，提高解析成功率
4. **独立性**：无内部依赖，可被外部项目引用

## 依赖

- `github.com/metacubex/mihomo/common/convert` - V2Ray 格式转换
- `gopkg.in/yaml.v3` - YAML 解析
- `github.com/lazygophers/utils/json` - JSON 解析

## 错误处理

所有解析函数在失败时会记录 Debug 日志并返回错误。调用方应处理错误：

```go
nodes, err := parser.Parse(data)
if err != nil {
	log.Errorf("解析失败: %v", err)
	return err
}

if len(nodes) == 0 {
	log.Warn("未解析到任何节点")
}

// 使用节点
for _, node := range nodes {
	log.Infof("节点: %s, ID: %s", node.Name(), node.UniqueId())
}
```

## 测试

TODO: 添加单元测试

```go
func TestParse(t *testing.T) {
	// 测试 Clash YAML 格式
	// 测试 V2Ray 格式
	// 测试逐行 JSON 格式
}
```

## 相关文档

- [Mihomo 文档](../mihomo/README.md)
- [产品设计文档](../../docs/产品设计文档.md)
