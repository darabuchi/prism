package parser

// ProxyType 代理类型
type ProxyType string

const (
	TypeShadowsocks   ProxyType = "ss"
	TypeShadowsocksR  ProxyType = "ssr"
	TypeVMess         ProxyType = "vmess"
	TypeVLess         ProxyType = "vless"
	TypeTrojan        ProxyType = "trojan"
	TypeHysteria      ProxyType = "hysteria"
	TypeHysteria2     ProxyType = "hysteria2"
	TypeSocks5        ProxyType = "socks5"
	TypeHTTP          ProxyType = "http"
	TypeSnell         ProxyType = "snell"
	TypeWireGuard     ProxyType = "wireguard"
	TypeTuic          ProxyType = "tuic"
	TypeSSH           ProxyType = "ssh"
	TypeMieru         ProxyType = "mieru"
	TypeAnyTLS        ProxyType = "anytls"
	TypeDirect        ProxyType = "direct"
	TypeReject        ProxyType = "reject"
	TypeDNS           ProxyType = "dns"
)

// ValidateProxyType 验证代理类型是否支持
func ValidateProxyType(proxyType string) bool {
	switch ProxyType(proxyType) {
	case TypeShadowsocks, TypeShadowsocksR,
		TypeVMess, TypeVLess,
		TypeTrojan, TypeHysteria, TypeHysteria2,
		TypeSocks5, TypeHTTP,
		TypeSnell, TypeWireGuard, TypeTuic,
		TypeSSH, TypeMieru, TypeAnyTLS,
		TypeDirect, TypeReject, TypeDNS:
		return true
	default:
		return false
	}
}
