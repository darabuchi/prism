package prism

// ProxyType 代理协议类型
//
// 定义所有支持的代理协议类型常量
// 这些常量用于标识节点的协议类型，与 Mihomo/Meta 核心的协议类型保持一致
type ProxyType string

const (
	// TypeShadowsocks Shadowsocks 代理协议
	TypeShadowsocks ProxyType = "ss"

	// TypeShadowsocksR ShadowsocksR 代理协议
	TypeShadowsocksR ProxyType = "ssr"

	// TypeVMess V2Ray VMess 协议
	TypeVMess ProxyType = "vmess"

	// TypeVLess V2Ray VLess 协议
	TypeVLess ProxyType = "vless"

	// TypeTrojan Trojan 代理协议
	TypeTrojan ProxyType = "trojan"

	// TypeHysteria Hysteria 协议
	TypeHysteria ProxyType = "hysteria"

	// TypeHysteria2 Hysteria2 协议
	TypeHysteria2 ProxyType = "hysteria2"

	// TypeSocks5 SOCKS5 代理协议
	TypeSocks5 ProxyType = "socks5"

	// TypeHTTP HTTP(S) 代理协议
	TypeHTTP ProxyType = "http"

	// TypeSnell Snell 协议
	TypeSnell ProxyType = "snell"

	// TypeWireGuard WireGuard VPN 协议
	TypeWireGuard ProxyType = "wireguard"

	// TypeTuic TUIC 协议
	TypeTuic ProxyType = "tuic"

	// TypeSSH SSH 隧道协议
	TypeSSH ProxyType = "ssh"

	// TypeMieru Mieru 协议
	TypeMieru ProxyType = "mieru"

	// TypeAnyTLS AnyTLS 协议
	TypeAnyTLS ProxyType = "anytls"

	// TypeDirect 直连（不使用代理）
	TypeDirect ProxyType = "direct"

	// TypeReject 拒绝连接
	TypeReject ProxyType = "reject"

	// TypeDNS DNS 查询
	TypeDNS ProxyType = "dns"
)

// HealthState 熔断器健康状态
//
// 用于表示节点的健康状态，基于熔断器模式
// 熔断器会根据节点的成功/失败次数自动切换状态
type HealthState int32

const (
	// HealthStateOpen 开路状态（熔断器开启）
	// 节点连续失败次数超过阈值，熔断器开启，拒绝所有请求
	HealthStateOpen HealthState = 0

	// HealthStateHalfOpen 半开状态（熔断器尝试恢复）
	// 熔断器开启一段时间后，进入半开状态，允许少量请求通过以测试节点是否恢复
	HealthStateHalfOpen HealthState = 1

	// HealthStateClosed 闭路状态（熔断器关闭）
	// 节点工作正常，熔断器关闭，允许所有请求通过
	HealthStateClosed HealthState = 2
)

// String 返回健康状态的字符串表示
func (h HealthState) String() string {
	switch h {
	case HealthStateOpen:
		return "open"
	case HealthStateHalfOpen:
		return "half_open"
	case HealthStateClosed:
		return "closed"
	default:
		return "unknown"
	}
}
