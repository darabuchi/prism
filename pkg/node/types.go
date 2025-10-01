package node

// ProxyType 代理类型
type ProxyType string

const (
	TypeShadowsocks  ProxyType = "ss"
	TypeShadowsocksR ProxyType = "ssr"
	TypeVMess        ProxyType = "vmess"
	TypeVLess        ProxyType = "vless"
	TypeTrojan       ProxyType = "trojan"
	TypeHysteria     ProxyType = "hysteria"
	TypeHysteria2    ProxyType = "hysteria2"
	TypeSocks5       ProxyType = "socks5"
	TypeHTTP         ProxyType = "http"
	TypeSnell        ProxyType = "snell"
	TypeWireGuard    ProxyType = "wireguard"
	TypeTuic         ProxyType = "tuic"
	TypeSSH          ProxyType = "ssh"
	TypeMieru        ProxyType = "mieru"
	TypeAnyTLS       ProxyType = "anytls"
	TypeDirect       ProxyType = "direct"
	TypeReject       ProxyType = "reject"
	TypeDNS          ProxyType = "dns"
)

// HealthState 熔断器健康状态
type HealthState int32

const (
	HealthStateOpen     HealthState = 0 // 开路（熔断器开启，拒绝请求）
	HealthStateHalfOpen HealthState = 1 // 半开（熔断器尝试恢复）
	HealthStateClosed   HealthState = 2 // 闭路（熔断器关闭，正常工作）
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
