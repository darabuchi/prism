# Parser - 订阅解析器

## 概述

`pkg/parser` 提供通用的代理订阅解析功能，支持多种订阅格式的自动识别和解析。

## 特性

- **多格式支持**：自动识别并解析 Clash YAML、V2Ray Base64 链接、逐行 JSON 格式
- **智能解析**：按优先级尝试不同格式，直到成功解析
- **协议兼容**：支持 18+ 种代理协议（SS、SSR、VMess、VLess、Trojan、Hysteria 等）
- **Mihomo 集成**：解析结果可直接用于 mihomo adapter.ParseProxy

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
	proxies, err := parser.Parse(subscriptionData)
	if err != nil {
		log.Errorf("解析失败: %v", err)
		return
	}

	log.Infof("解析到 %d 个节点", len(proxies))
	for _, proxy := range proxies {
		log.Infof("节点: %v", proxy["name"])
	}
}
```

### 与 Mihomo Adapter 集成

```go
import (
	"github.com/darabuchi/prism/pkg/parser"
	"github.com/metacubex/mihomo/adapter"
)

// 解析订阅并创建 mihomo 代理
func ParseAndCreateProxies(data []byte) error {
	// 解析订阅
	proxies, err := parser.Parse(data)
	if err != nil {
		return err
	}

	// 转换为 mihomo Proxy
	for _, proxyMap := range proxies {
		proxy, err := adapter.ParseProxy(proxyMap)
		if err != nil {
			log.Warnf("代理解析失败: %v", err)
			continue
		}

		// 使用 proxy...
		log.Infof("创建代理: %s", proxy.Name())
	}

	return nil
}
```

### 验证代理类型

```go
import "github.com/darabuchi/prism/pkg/parser"

func main() {
	// 验证代理类型是否支持
	if parser.ValidateProxyType("vmess") {
		log.Info("VMess 协议已支持")
	}

	if !parser.ValidateProxyType("unknown") {
		log.Warn("不支持的协议类型")
	}
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
返回 []map[string]any
```

## API 文档

### Parse

```go
func Parse(data []byte) ([]map[string]any, error)
```

解析订阅数据，自动识别格式。

**参数**：
- `data`: 订阅内容字节数组

**返回**：
- `[]map[string]any`: 代理配置列表，可直接用于 mihomo adapter.ParseProxy
- `error`: 解析错误

### ValidateProxyType

```go
func ValidateProxyType(proxyType string) bool
```

验证代理类型是否支持。

**参数**：
- `proxyType`: 代理类型字符串（如 "ss", "vmess"）

**返回**：
- `bool`: 是否支持该代理类型

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
proxies, err := parser.Parse(data)
if err != nil {
	log.Errorf("解析失败: %v", err)
	return err
}

if len(proxies) == 0 {
	log.Warn("未解析到任何节点")
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
