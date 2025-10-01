# Node - 节点管理包

## 概述

`pkg/node` 提供节点管理的基础功能，包括节点结构体、ID 生成算法等。
是 Prism 项目中节点管理的核心工具包，可被外部项目引用。

## 设计原则

- **高度独立**：不依赖 internal 包，可被外部项目引用
- **规范遵循**：严格遵循后端开发规范，使用指定的依赖包
- **与 Fire 保持一致**：ID 生成算法与 Fire 项目的 CalculateClashHash 逻辑保持一致

## 功能特性

- **Node 结构体**：封装 Mihomo 代理适配器的节点管理结构
- **ID 生成**：基于 URL 格式和 SHA256 的节点唯一标识生成
- **配置管理**：支持 Clash 格式配置和 Base64 编码的原始配置

## 类型定义

代理协议类型（`prism.ProxyType`）和熔断器健康状态（`prism.HealthState`）定义在项目根目录的 `types.go` 文件中。

支持 18+ 种代理协议类型：
- SS/SSR/VMess/VLess/Trojan
- Hysteria/Hysteria2/TUIC
- SOCKS5/HTTP/Snell
- WireGuard/SSH/Mieru/AnyTLS
- Direct/Reject/DNS

## 使用方式

### 创建节点

```go
import "github.com/darabuchi/prism/pkg/node"

// 从 Clash 配置创建节点
config := map[string]any{
    "name":   "香港节点",
    "type":   "vmess",
    "server": "hk.example.com",
    "port":   443,
    "uuid":   "xxx-xxx-xxx",
}

n, err := node.NewNode(config)
if err != nil {
    log.Errorf("err:%v", err)
    return
}

log.Infof("节点 ID: %s", n.UniqueId())
log.Infof("节点名称: %s", n.Name())
log.Infof("节点类型: %s", n.ProxyType())
log.Infof("服务器: %s:%d", n.Server(), n.Port())
```

### 使用节点进行代理连接

```go
// 创建连接元数据
metadata := &constant.Metadata{
    Host:    "example.com",
    DstPort: 443,
    NetWork: constant.TCP,
}

// 使用节点建立代理连接
conn, err := n.DialContext(context.Background(), metadata)
if err != nil {
    log.Errorf("err:%v", err)
    return
}
defer conn.Close()

// 使用连接进行数据传输
// ...
```

### 生成节点 ID

```go
import "github.com/darabuchi/prism/pkg/node"

config := map[string]any{
    "type":   "vmess",
    "server": "hk.example.com",
    "port":   443,
    "uuid":   "xxx-xxx-xxx",
}

// 生成唯一 ID（SHA256 哈希）
// ID 基于 URL 格式：vmess://hk.example.com:443?uuid=xxx-xxx-xxx
id := node.GenerateId(config)
log.Infof("节点 ID: %s", id) // 输出：64位十六进制字符串
```

### 使用代理类型常量

```go
import "github.com/darabuchi/prism"

// 使用 prism 包中的类型常量
proxyType := prism.TypeVMess
log.Infof("协议类型: %s", proxyType) // 输出：vmess

// 在配置中使用
config := map[string]any{
    "type": string(prism.TypeVMess),
    // ...
}
```

## API 文档

### Node 结构体

```go
type Node struct {
    // 私有字段
}
```

节点结构体，封装 Mihomo 代理适配器和节点配置。

**构造函数**：

- `NewNode(config map[string]any) (*Node, error)` - 从 Clash 配置创建节点（自动创建 Mihomo 适配器）

**基础信息方法**：

- `UniqueId() string` - 获取节点唯一标识（SHA256 哈希）
- `Name() string` - 获取节点名称（来自 adapter）
- `Type() constant.AdapterType` - 获取适配器类型（实现 ProxyAdapter 接口）
- `ProxyType() prism.ProxyType` - 获取代理协议类型
- `Addr() string` - 获取代理地址（host:port 格式）
- `Server() string` - 获取服务器地址（从 Addr 解析）
- `Port() int` - 获取服务器端口（从 Addr 解析）
- `Config() map[string]any` - 获取原始配置

**网络连接方法**：

- `DialContext(ctx, metadata) (Conn, error)` - 建立 TCP 代理连接
- `ListenPacketContext(ctx, metadata) (PacketConn, error)` - 建立 UDP 代理连接
- `DialContextWithDialer(ctx, dialer, metadata) (Conn, error)` - 使用自定义拨号器建立 TCP 连接
- `ListenPacketWithDialer(ctx, dialer, metadata) (PacketConn, error)` - 使用自定义拨号器建立 UDP 连接
- `StreamConnContext(ctx, c, metadata) (net.Conn, error)` - 在现有连接上包装协议

**功能检查方法**：

- `SupportUDP() bool` - 检查是否支持 UDP
- `SupportUOT() bool` - 检查是否支持 UDP over TCP
- `SupportWithDialer() NetWork` - 获取支持的网络类型
- `IsL3Protocol(metadata) bool` - 检查是否为 L3 协议

**其他方法**：

- `ProxyInfo() ProxyInfo` - 获取代理信息（XUDP、TFO、MPTCP 等）
- `MarshalJSON() ([]byte, error)` - 序列化为 JSON
- `Unwrap(metadata, touch) Proxy` - 解包代理
- `Close() error` - 关闭节点连接

### 函数

#### GenerateId

```go
func GenerateId(config map[string]any) string
```

生成节点唯一 ID（基于配置的 SHA256 哈希）。

与 Fire 项目的 CalculateClashHash 逻辑保持一致：
1. 如果配置中有 unique_id 字段，直接返回
2. 排除不影响节点唯一性的字段（name、哈希值等）
3. 将配置转换为 URL 格式并计算 SHA256

**参数**：
- `config`: 节点配置（必须包含 type、server、port 等字段）

**返回**：
- `string`: 64位十六进制 SHA256 哈希字符串

**示例**：
```go
config := map[string]any{
    "type":   "vmess",
    "server": "example.com",
    "port":   443,
}
id := GenerateId(config)
// 输出类似: "a1b2c3d4..."
```

## 依赖包

本包严格遵循项目的后端开发规范，使用以下依赖：

```go
import (
    "github.com/darabuchi/prism"                    // 项目类型定义
    "github.com/lazygophers/log"                    // 日志处理
    "github.com/lazygophers/lrpc/middleware/xerror" // 错误处理
    "github.com/lazygophers/utils/candy"            // 类型转换工具
    "github.com/lazygophers/utils/cryptox"          // 加密工具（SHA256）
    "github.com/lazygophers/utils/json"             // JSON 处理
    "github.com/metacubex/mihomo/constant"          // Mihomo 核心接口
)
```

## 设计说明

### Node 结构体

`pkg/node` 中的 Node 结构体是基础节点管理结构，封装了：
- Mihomo 代理适配器（constant.ProxyAdapter）
- 节点配置（Clash 格式）
- 唯一标识（SHA256 哈希）

**接口实现**：
- **实现 constant.ProxyAdapter**：Node 完整实现了 Mihomo 的 ProxyAdapter 接口
- **透明代理**：所有 ProxyAdapter 方法都转发给内部的 adapter
- **额外功能**：在 ProxyAdapter 基础上增加了 UniqueId、Server、Port、ProxyType 等便捷方法

**核心特性**：
- **自动创建适配器**：NewNode 时自动调用 mihomo 的 adapter.ParseProxy 创建代理适配器
- **即开即用**：创建节点后即可直接使用 DialContext 等方法进行代理连接
- **协议标准化**：自动将 shadowsocks、hy 等长格式转换为 ss、hysteria 等短格式
- **证书初始化**：包初始化时自动调用 ca.ResetCertificate 避免证书问题
- **方法从 adapter 读取**：Name、Type、Addr、Server、Port 等方法都从 adapter 读取，确保数据一致性

复杂的业务逻辑（如状态管理、测试结果、评分、数据库操作等）应放在 `internal/` 中实现。

### ID 生成算法

GenerateId 函数与 Fire 项目的 CalculateClashHash 保持完全一致：

1. **检查 unique_id**：如果配置中有 unique_id 字段，直接返回
2. **过滤字段**：删除 name、md5、sha256 等不影响唯一性的字段
3. **构建 URL**：
   - scheme = type（如 vmess）
   - host = server:port（如 example.com:443）
   - query = 其他所有字段的 URL 编码
4. **计算哈希**：对 URL 字符串计算 SHA256

这确保了相同配置的节点在 Prism 和 Fire 项目中具有相同的 ID。

### 配置格式

使用 `map[string]any` 的原因：
- 支持 18+ 种代理协议，每种协议参数不同
- 提供最大的灵活性和扩展性
- 与 Mihomo/Clash 配置格式保持一致

### 与 internal 的关系

- `pkg/node`：基础节点管理（可被外部引用）
- `internal/model`：数据库模型（GORM）
- `internal/service`：业务逻辑（测试、评分、熔断等）
- `internal/proxy`：代理管理（结合 pkg/node 和业务模型）

## 测试

TODO: 添加单元测试

```go
func TestNewNode(t *testing.T) {
    config := map[string]any{
        "type":   "vmess",
        "server": "example.com",
        "port":   443,
    }

    node, err := NewNode(config)
    if err != nil {
        t.Fatalf("NewNode failed: %v", err)
    }

    if node.ProxyType() != prism.TypeVMess {
        t.Errorf("expected vmess, got %s", node.ProxyType())
    }

    if len(node.UniqueId()) != 64 { // SHA256 produces 64 hex chars
        t.Errorf("expected 64 chars, got %d", len(node.UniqueId()))
    }
}

func TestGenerateId(t *testing.T) {
    config := map[string]any{
        "type":   "vmess",
        "server": "example.com",
        "port":   443,
    }

    id := GenerateId(config)
    if len(id) != 64 {
        t.Errorf("expected 64 chars, got %d", len(id))
    }

    // 相同配置应生成相同的 ID
    id2 := GenerateId(config)
    if id != id2 {
        t.Errorf("expected same ID for same config")
    }
}
```

## 相关文档

- [项目类型定义](../../types.go)
- [后端开发规范](../../docs/编码规范/后端开发规范.md)
- [Parser 订阅解析器](../parser/README.md)
- [系统设计文档](../../docs/系统设计文档.md)
