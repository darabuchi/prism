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
type HealthState int

const (
	HealthStateOpen     HealthState = 0 // 开路（熔断器开启，拒绝请求）
	HealthStateHalfOpen HealthState = 1 // 半开（熔断器尝试恢复）
	HealthStateClosed   HealthState = 2 // 闭路（熔断器关闭，正常工作）
)

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

// Config 节点配置
type Config struct {
	Name    string         `json:"name" yaml:"name"`       // 节点名称
	Type    string         `json:"type" yaml:"type"`       // 协议类型
	Server  string         `json:"server" yaml:"server"`   // 服务器地址
	Port    int            `json:"port" yaml:"port"`       // 端口
	Config  map[string]any `json:"-" yaml:"-"`             // 完整配置（不序列化）
	RawJSON string         `json:"raw_json" yaml:"raw_json"` // Base64 编码的原始 JSON 配置
}

// Status 节点运行时状态
type Status struct {
	Alive       bool        `json:"alive"`        // 是否存活
	Enabled     bool        `json:"enabled"`      // 是否启用
	HealthState HealthState `json:"health_state"` // 熔断器状态

	// 熔断器统计
	SuccessCount uint64 `json:"success_count"` // 成功请求数
	FailureCount uint64 `json:"failure_count"` // 失败请求数

	// 流量统计
	UploadTotal   int64 `json:"upload_total"`   // 总上传流量（字节）
	DownloadTotal int64 `json:"download_total"` // 总下载流量（字节）
}

// TestResult 节点测试结果
type TestResult struct {
	// 延迟测试 (ms)
	DelayAvg       int   `json:"delay_avg"`        // 平均延迟
	DelayMin       int   `json:"delay_min"`        // 最小延迟
	DelayMax       int   `json:"delay_max"`        // 最大延迟
	DelayTestedAt  int64 `json:"delay_tested_at"`  // 延迟测试时间
	DelayNextTestAt int64 `json:"delay_next_test_at"` // 下次延迟测试时间

	// 速度测试 (bytes/s)
	DownloadSpeed      int64 `json:"download_speed"`       // 下载速度
	UploadSpeed        int64 `json:"upload_speed"`         // 上传速度
	SpeedTestedAt      int64 `json:"speed_tested_at"`      // 速度测试时间
	SpeedNextTestAt    int64 `json:"speed_next_test_at"`   // 下次速度测试时间

	// 地理信息
	InboundCountry  string `json:"inbound_country"`  // 入口国家代码
	InboundIP       string `json:"inbound_ip"`       // 入口 IP
	InboundASN      string `json:"inbound_asn"`      // 入口 ASN
	InboundISP      string `json:"inbound_isp"`      // 入口 ISP
	OutboundCountry string `json:"outbound_country"` // 出口国家代码
	OutboundIP      string `json:"outbound_ip"`      // 出口 IP
	OutboundASN     string `json:"outbound_asn"`     // 出口 ASN
	OutboundISP     string `json:"outbound_isp"`     // 出口 ISP
	RegionTestedAt  int64  `json:"region_tested_at"` // 地理信息测试时间
	RegionNextTestAt int64 `json:"region_next_test_at"` // 下次地理信息测试时间

	// 解锁检测 (0:未测试 1:失败 2:部分解锁 3:完全解锁)
	UnlockOpenAI   int   `json:"unlock_openai"`    // OpenAI 解锁状态
	UnlockNetflix  int   `json:"unlock_netflix"`   // Netflix 解锁状态
	UnlockYouTube  int   `json:"unlock_youtube"`   // YouTube 解锁状态
	UnlockDisney   int   `json:"unlock_disney"`    // Disney+ 解锁状态
	UnlockTestedAt int64 `json:"unlock_tested_at"` // 解锁测试时间
	UnlockNextTestAt int64 `json:"unlock_next_test_at"` // 下次解锁测试时间

	// 线路质量
	HighQualityRoutes bool     `json:"high_quality_routes"` // 是否为优质线路
	RouterList        []string `json:"router_list"`         // 路由列表
}

// Info 节点完整信息
type Info struct {
	ID        string      `json:"id"`         // 唯一标识（SHA256）
	UniqueKey string      `json:"unique_key"` // 唯一键
	Config    Config      `json:"config"`     // 配置信息
	Status    Status      `json:"status"`     // 运行状态
	Test      TestResult  `json:"test"`       // 测试结果
	Score     float64     `json:"score"`      // 综合评分

	// 时间戳
	FirstAliveAt int64 `json:"first_alive_at"` // 首次存活时间
	LastAliveAt  int64 `json:"last_alive_at"`  // 最后存活时间
	DeathCount   int64 `json:"death_count"`    // 死亡计数
	CreatedAt    int64 `json:"created_at"`     // 创建时间
	UpdatedAt    int64 `json:"updated_at"`     // 更新时间
}
