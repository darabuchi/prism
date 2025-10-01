# Rules 包 - 规则分流引擎

`rules` 包提供了强大而灵活的规则分流功能，用于根据连接的各种属性（域名、IP、端口、进程等）进行流量分类和路由决策。

## 特性

- 🎯 **多种规则类型**：支持域名、IP、端口、进程、GeoIP、ASN 等多种匹配方式
- ⚡ **高性能匹配**：优化的匹配算法，最小化字符串操作和内存分配
- 📝 **简单易用**：直观的规则语法，易于编写和维护
- 🔧 **灵活扩展**：清晰的接口设计，易于添加自定义规则类型
- ✅ **完整测试**：全面的单元测试覆盖

## 支持的规则类型

### 域名规则

| 规则类型 | 说明 | 示例 | 匹配示例 |
|---------|------|------|---------|
| `DOMAIN` | 域名精确匹配 | `DOMAIN,google.com,PROXY` | `google.com` ✓ `www.google.com` ✗ |
| `DOMAIN-SUFFIX` | 域名后缀匹配 | `DOMAIN-SUFFIX,google.com,PROXY` | `google.com` ✓ `www.google.com` ✓ |
| `DOMAIN-KEYWORD` | 域名关键字匹配 | `DOMAIN-KEYWORD,google,PROXY` | `google.com` ✓ `www.google.co.uk` ✓ |
| `DOMAIN-REGEX` | 域名正则表达式 | `DOMAIN-REGEX,^.*\.cn$,DIRECT` | `example.cn` ✓ `www.test.cn` ✓ |

### IP 规则

| 规则类型 | 说明 | 示例 | 匹配示例 |
|---------|------|------|---------|
| `IP-CIDR` | IPv4 CIDR 匹配 | `IP-CIDR,192.168.0.0/16,DIRECT` | `192.168.1.1` ✓ `10.0.0.1` ✗ |
| `IP-CIDR6` | IPv6 CIDR 匹配 | `IP-CIDR6,2001:db8::/32,PROXY` | `2001:db8::1` ✓ |
| `GEOIP` | GeoIP 国家代码 | `GEOIP,CN,DIRECT` | 中国 IP ✓ |
| `IP-ASN` | ASN 号码匹配 | `IP-ASN,13335,PROXY` | Cloudflare IP ✓ |

### 端口规则

| 规则类型 | 说明 | 示例 |
|---------|------|------|
| `DST-PORT` | 目标端口匹配 | `DST-PORT,80,DIRECT` |
| `DST-PORT` | 目标端口范围 | `DST-PORT,8000-9000,PROXY` |
| `SRC-PORT` | 源端口匹配 | `SRC-PORT,7890,REJECT` |

### 进程规则

| 规则类型 | 说明 | 示例 |
|---------|------|------|
| `PROCESS-NAME` | 进程名称匹配 | `PROCESS-NAME,chrome,PROXY` |
| `PROCESS-PATH` | 进程路径匹配 | `PROCESS-PATH,/usr/bin/wget,DIRECT` |

### 特殊规则

| 规则类型 | 说明 | 示例 |
|---------|------|------|
| `MATCH` | 匹配所有流量 | `MATCH,PROXY` |

## 快速开始

### 基本用法

```go
package main

import (
	"fmt"
	"net/netip"

	"github.com/darabuchi/prism/pkg/rules"
)

func main() {
	// 解析单条规则
	rule, err := rules.ParseRule("DOMAIN,google.com,PROXY")
	if err != nil {
		panic(err)
	}

	// 创建连接元数据
	metadata := &rules.Metadata{
		Host: "google.com",
		DstPort: 443,
	}

	// 匹配规则
	if rule.Match(metadata) {
		fmt.Printf("匹配成功，动作: %s\n", rule.Action())
	}
}
```

### 从文件加载规则

```go
// 从文件加载规则列表
rules, err := rules.LoadRulesFromFile("rules.txt")
if err != nil {
	panic(err)
}

// 匹配第一个符合的规则
metadata := &rules.Metadata{
	Host:    "www.google.com",
	DstPort: 443,
}

if matchedRule, ok := rules.MatchFirst(rules, metadata); ok {
	fmt.Printf("匹配规则: %s\n", matchedRule)
	fmt.Printf("动作: %s\n", matchedRule.Action())
}
```

### 规则文件格式

```text
# 注释行以 # 或 // 开头

# 域名规则
DOMAIN,google.com,PROXY
DOMAIN-SUFFIX,google.com,PROXY
DOMAIN-KEYWORD,google,PROXY

# IP 规则
IP-CIDR,192.168.0.0/16,DIRECT
IP-CIDR,10.0.0.0/8,DIRECT
GEOIP,CN,DIRECT

# 端口规则
DST-PORT,80,DIRECT
DST-PORT,443,PROXY
DST-PORT,8000-9000,REJECT

# 进程规则
PROCESS-NAME,chrome,PROXY
PROCESS-PATH,/usr/bin/wget,DIRECT

# 匹配所有（通常放在最后）
MATCH,PROXY
```

### 使用 GeoIP 功能

GeoIP 和 IP-ASN 规则需要配置 GeoIP 数据提供者：

```go
package main

import (
	"net"

	"github.com/darabuchi/prism"
	"github.com/darabuchi/prism/pkg/geoip"
	"github.com/darabuchi/prism/pkg/rules"
)

func main() {
	// 创建 GeoIP 读取器
	reader, err := geoip.NewReader("geoip.mmdb")
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	// 设置 GeoIP 提供者
	rules.SetGeoIPProvider(func(ip string) *prism.GeoIP {
		info, _ := reader.LookupString(ip)
		return info
	})

	// 现在可以使用 GEOIP 和 IP-ASN 规则了
	rule, _ := rules.ParseRule("GEOIP,CN,DIRECT")

	metadata := &rules.Metadata{
		DstIP: netip.MustParseAddr("1.1.1.1"),
	}

	if rule.Match(metadata) {
		fmt.Println("匹配中国 IP")
	}
}
```

## 高级用法

### 创建自定义规则

```go
// 直接创建规则对象
domainRule := rules.NewDomain("example.com", rules.ActionProxy)
suffixRule := rules.NewDomainSuffix("google.com", rules.ActionProxy)
cidrRule, _ := rules.NewIPCIDR("192.168.0.0/16", rules.ActionDirect, false)
```

### 规则选项

某些规则支持额外选项：

```go
// IP-CIDR 规则支持 no-resolve 选项
rule, _ := rules.ParseRule("IP-CIDR,192.168.0.0/16,DIRECT,no-resolve")

// 或者通过代码设置
cidrRule, _ := rules.NewIPCIDR("192.168.0.0/16", rules.ActionDirect, false)
cidrRule.NoResolve(true)
```

`no-resolve` 选项表示如果目标是域名（未解析为 IP），则不匹配此规则。

### 批量解析规则

```go
ruleStrs := []string{
	"DOMAIN,google.com,PROXY",
	"DOMAIN-SUFFIX,example.com,DIRECT",
	"IP-CIDR,192.168.0.0/16,DIRECT",
	"MATCH,PROXY",
}

rules, err := rules.ParseRules(ruleStrs)
if err != nil {
	panic(err)
}
```

## 性能优化

### 规则顺序

将常用规则放在前面可以提高匹配性能：

```text
# 推荐：常用规则在前
DOMAIN,frequently-used.com,PROXY
DOMAIN-SUFFIX,common-site.com,DIRECT
# ... 其他规则

# 不推荐：很少匹配的规则在前
DOMAIN-REGEX,^very-rare-pattern$,PROXY
```

### 规则类型选择

选择最精确的规则类型：

- ✅ 优先使用 `DOMAIN` 而不是 `DOMAIN-SUFFIX`
- ✅ 优先使用 `DOMAIN-SUFFIX` 而不是 `DOMAIN-KEYWORD`
- ✅ 避免过度使用 `DOMAIN-REGEX`（性能较低）

### 预编译规则

在程序启动时预编译所有规则：

```go
// 启动时加载
var globalRules []rules.Rule

func init() {
	var err error
	globalRules, err = rules.LoadRulesFromFile("rules.txt")
	if err != nil {
		panic(err)
	}
}

// 使用时直接匹配
func matchRequest(metadata *rules.Metadata) rules.ActionType {
	if rule, ok := rules.MatchFirst(globalRules, metadata); ok {
		return rule.Action()
	}
	return rules.ActionDirect
}
```

## 规则动作

支持以下动作类型：

| 动作 | 说明 |
|------|------|
| `PROXY` | 使用代理 |
| `DIRECT` | 直接连接 |
| `REJECT` | 拒绝连接 |
| `REJECT-DROP` | 拒绝连接并丢弃数据包 |

## 错误处理

```go
rule, err := rules.ParseRule("INVALID,rule,format")
if err != nil {
	switch {
	case errors.Is(err, rules.ErrInvalidRule):
		fmt.Println("规则格式错误")
	case errors.Is(err, rules.ErrInvalidRuleType):
		fmt.Println("规则类型不支持")
	case errors.Is(err, rules.ErrInvalidIPCIDR):
		fmt.Println("IP CIDR 格式错误")
	default:
		fmt.Printf("未知错误: %v\n", err)
	}
}
```

## 测试

运行测试：

```bash
go test -v ./pkg/rules
```

运行基准测试：

```bash
go test -bench=. ./pkg/rules
```

## 兼容性

本包设计参考了 Clash、Surge 等主流代理工具的规则语法，大部分规则格式可以直接兼容。

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！
