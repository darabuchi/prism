package node

import (
	"context"
	"errors"
	"net"

	"github.com/darabuchi/prism"
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/xerror"
	"github.com/lazygophers/utils/candy"
	"github.com/metacubex/mihomo/adapter"
	"github.com/metacubex/mihomo/component/ca"
	"github.com/metacubex/mihomo/constant"
)

// Error definitions
var (
	ErrUnsupportedSpeedTestURL = errors.New("unsupported speed test URL format")
)

func init() {
	// 重置证书，避免证书问题
	// 参考 fire 项目的实现
	ca.ResetCertificate()
}

// 确保 Node 实现了 constant.ProxyAdapter 接口
var _ constant.ProxyAdapter = (*Node)(nil)

// Node 节点结构体
//
// 封装了 Mihomo 的 ProxyAdapter，提供节点管理的基础功能
// 包含节点配置、唯一标识、原始数据等信息
type Node struct {
	// adapter Mihomo 代理适配器，提供实际的代理功能
	adapter constant.ProxyAdapter

	// config 节点配置（Clash 格式）
	config map[string]any

	// uniqueId 节点唯一标识（基于配置的 SHA256 哈希）
	uniqueId string
}

// NewNode 创建节点实例
//
// 从 Clash 格式的配置创建节点实例，自动创建 Mihomo 代理适配器
// 配置必须包含 type、server、port 等基本字段
//
// 参数:
//   - config: Clash 格式的节点配置
//
// 返回:
//   - *Node: 节点实例
//   - error: 创建失败时返回错误
//
// 示例:
//
//	config := map[string]any{
//	    "type":   "vmess",
//	    "server": "example.com",
//	    "port":   443,
//	    "uuid":   "xxx-xxx-xxx",
//	}
//	node, err := NewNode(config)
func NewNode(config map[string]any) (*Node, error) {
	// 确保配置有 name 字段
	if _, ok := config["name"]; !ok {
		config["name"] = "prism-node"
	}

	// 标准化协议类型名称
	// 将长格式转换为 Mihomo 识别的短格式
	if typ, ok := config["type"]; ok {
		if t, ok := typ.(string); ok {
			switch t {
			case "shadowsocks":
				config["type"] = "ss"
			case "shadowsocksr":
				config["type"] = "ssr"
			case "hy":
				config["type"] = "hysteria"
			case "hy2":
				config["type"] = "hysteria2"
			}
		}
	}

	// 使用 Mihomo 的 ParseProxy 创建代理适配器
	// 这是 Mihomo 提供的标准方法，用于解析 Clash 格式配置
	proxyAdapter, err := adapter.ParseProxy(config)
	if err != nil {
		log.Errorf("parse proxy config failed: %v, config: %+v", err, config)
		return nil, xerror.WrapError(err, xerror.ErrSystemError, "parse proxy config failed")
	}

	// 生成节点唯一 ID
	uniqueId := GenerateId(config)

	return &Node{
		adapter:  proxyAdapter,
		config:   config,
		uniqueId: uniqueId,
	}, nil
}

// UniqueId 获取节点唯一标识
func (n *Node) UniqueId() string {
	return n.uniqueId
}

// Config 获取节点配置
func (n *Node) Config() map[string]any {
	return n.config
}

// Name 获取节点名称
//
// 使用 Mihomo adapter 的 Name() 方法
func (n *Node) Name() string {
	return n.adapter.Name()
}

// Type 获取节点类型
//
// 实现 constant.ProxyAdapter 接口
// 返回 Mihomo 的 AdapterType
func (n *Node) Type() constant.AdapterType {
	return n.adapter.Type()
}

// ProxyType 获取节点协议类型
//
// 返回 Prism 定义的 ProxyType
func (n *Node) ProxyType() prism.ProxyType {
	return prism.ProxyType(n.adapter.Type().String())
}

// Server 获取服务器地址
//
// 从 Mihomo adapter 的 Addr() 中解析出 host 部分
func (n *Node) Server() string {
	host, _, _ := net.SplitHostPort(n.adapter.Addr())
	return host
}

// Port 获取服务器端口
//
// 从 Mihomo adapter 的 Addr() 中解析出 port 部分
// 使用 candy.ToInt 进行类型转换
func (n *Node) Port() int {
	_, portStr, _ := net.SplitHostPort(n.adapter.Addr())
	return candy.ToInt(portStr)
}

// DialContext 使用节点拨号连接
//
// 通过 Mihomo 代理适配器建立网络连接
// 这是使用节点进行代理连接的核心方法
//
// 参数:
//   - ctx: 上下文，用于控制连接超时和取消
//   - metadata: 连接元数据，包含目标地址、端口、协议等信息
//
// 返回:
//   - constant.Conn: Mihomo 连接对象
//   - error: 连接失败时返回错误
func (n *Node) DialContext(ctx context.Context, metadata *constant.Metadata) (constant.Conn, error) {
	return n.adapter.DialContext(ctx, metadata)
}

// ListenPacketContext 使用节点监听 UDP 数据包
//
// 通过 Mihomo 代理适配器创建 UDP 连接
// 用于 UDP 协议的代理连接
//
// 参数:
//   - ctx: 上下文，用于控制连接超时和取消
//   - metadata: 连接元数据，包含目标地址、端口、协议等信息
//
// 返回:
//   - constant.PacketConn: Mihomo UDP 连接对象
//   - error: 连接失败时返回错误
func (n *Node) ListenPacketContext(ctx context.Context, metadata *constant.Metadata) (constant.PacketConn, error) {
	return n.adapter.ListenPacketContext(ctx, metadata)
}

// Addr 获取代理地址
//
// 返回格式：host:port
func (n *Node) Addr() string {
	return n.adapter.Addr()
}

// SupportUDP 检查节点是否支持 UDP
//
// 实现 constant.ProxyAdapter 接口
func (n *Node) SupportUDP() bool {
	return n.adapter.SupportUDP()
}

// ProxyInfo 获取代理信息
//
// 实现 constant.ProxyAdapter 接口
// 返回额外的代理信息（XUDP、TFO、MPTCP 等）
func (n *Node) ProxyInfo() constant.ProxyInfo {
	return n.adapter.ProxyInfo()
}

// MarshalJSON 序列化为 JSON
//
// 实现 constant.ProxyAdapter 接口
func (n *Node) MarshalJSON() ([]byte, error) {
	return n.adapter.MarshalJSON()
}

// StreamConnContext 在现有连接上包装协议
//
// 实现 constant.ProxyAdapter 接口
// 在 net.Conn 上包装代理协议
func (n *Node) StreamConnContext(ctx context.Context, c net.Conn, metadata *constant.Metadata) (net.Conn, error) {
	return n.adapter.StreamConnContext(ctx, c, metadata)
}

// SupportUOT 检查是否支持 UDP over TCP
//
// 实现 constant.ProxyAdapter 接口
func (n *Node) SupportUOT() bool {
	return n.adapter.SupportUOT()
}

// SupportWithDialer 获取支持的网络类型
//
// 实现 constant.ProxyAdapter 接口
func (n *Node) SupportWithDialer() constant.NetWork {
	return n.adapter.SupportWithDialer()
}

// DialContextWithDialer 使用自定义拨号器建立连接
//
// 实现 constant.ProxyAdapter 接口
func (n *Node) DialContextWithDialer(ctx context.Context, dialer constant.Dialer, metadata *constant.Metadata) (constant.Conn, error) {
	return n.adapter.DialContextWithDialer(ctx, dialer, metadata)
}

// ListenPacketWithDialer 使用自定义拨号器监听 UDP
//
// 实现 constant.ProxyAdapter 接口
func (n *Node) ListenPacketWithDialer(ctx context.Context, dialer constant.Dialer, metadata *constant.Metadata) (constant.PacketConn, error) {
	return n.adapter.ListenPacketWithDialer(ctx, dialer, metadata)
}

// IsL3Protocol 检查是否为 L3 协议
//
// 实现 constant.ProxyAdapter 接口
func (n *Node) IsL3Protocol(metadata *constant.Metadata) bool {
	return n.adapter.IsL3Protocol(metadata)
}

// Unwrap 解包代理
//
// 实现 constant.ProxyAdapter 接口
func (n *Node) Unwrap(metadata *constant.Metadata, touch bool) constant.Proxy {
	return n.adapter.Unwrap(metadata, touch)
}

// Close 关闭节点连接
//
// 实现 constant.ProxyAdapter 接口
func (n *Node) Close() error {
	return nil
}
