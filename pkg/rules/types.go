package rules

import "net/netip"

// RuleType 规则类型
type RuleType string

const (
	// TypeDomain 域名精确匹配
	// 示例: DOMAIN,google.com,PROXY
	TypeDomain RuleType = "DOMAIN"

	// TypeDomainSuffix 域名后缀匹配
	// 示例: DOMAIN-SUFFIX,google.com,PROXY (匹配 *.google.com)
	TypeDomainSuffix RuleType = "DOMAIN-SUFFIX"

	// TypeDomainKeyword 域名关键字匹配
	// 示例: DOMAIN-KEYWORD,google,PROXY
	TypeDomainKeyword RuleType = "DOMAIN-KEYWORD"

	// TypeDomainRegex 域名正则表达式匹配
	// 示例: DOMAIN-REGEX,^.*\.google\.com$,PROXY
	TypeDomainRegex RuleType = "DOMAIN-REGEX"

	// TypeGEOIP GeoIP 国家代码匹配
	// 示例: GEOIP,CN,DIRECT
	TypeGEOIP RuleType = "GEOIP"

	// TypeIPCIDR IP CIDR 匹配
	// 示例: IP-CIDR,192.168.0.0/16,DIRECT
	TypeIPCIDR RuleType = "IP-CIDR"

	// TypeIPCIDR6 IPv6 CIDR 匹配
	// 示例: IP-CIDR6,2001:db8::/32,PROXY
	TypeIPCIDR6 RuleType = "IP-CIDR6"

	// TypeSrcIP 源 IP 匹配
	// 示例: SRC-IP,192.168.1.1,DIRECT
	TypeSrcIP RuleType = "SRC-IP"

	// TypeSrcIPCIDR 源 IP CIDR 匹配
	// 示例: SRC-IP-CIDR,192.168.0.0/16,DIRECT
	TypeSrcIPCIDR RuleType = "SRC-IP-CIDR"

	// TypeDstPort 目标端口匹配
	// 示例: DST-PORT,80,DIRECT
	TypeDstPort RuleType = "DST-PORT"

	// TypeSrcPort 源端口匹配
	// 示例: SRC-PORT,7890,DIRECT
	TypeSrcPort RuleType = "SRC-PORT"

	// TypeProcess 进程名称匹配
	// 示例: PROCESS-NAME,chrome,PROXY
	TypeProcess RuleType = "PROCESS-NAME"

	// TypeProcessPath 进程路径匹配
	// 示例: PROCESS-PATH,/usr/bin/wget,PROXY
	TypeProcessPath RuleType = "PROCESS-PATH"

	// TypeIPASN IP ASN 匹配
	// 示例: IP-ASN,13335,PROXY
	TypeIPASN RuleType = "IP-ASN"

	// TypeRuleSet 规则集匹配
	// 示例: RULE-SET,reject,REJECT
	TypeRuleSet RuleType = "RULE-SET"

	// TypeInType 入站类型匹配
	// 示例: IN-TYPE,HTTP,PROXY
	TypeInType RuleType = "IN-TYPE"

	// TypeInName 入站名称匹配
	// 示例: IN-NAME,socks-in,DIRECT
	TypeInName RuleType = "IN-NAME"

	// TypeInUser 入站用户匹配
	// 示例: IN-USER,admin,PROXY
	TypeInUser RuleType = "IN-USER"

	// TypeNetwork 网络类型匹配
	// 示例: NETWORK,TCP,DIRECT
	TypeNetwork RuleType = "NETWORK"

	// TypeUID 用户 ID 匹配
	// 示例: UID,1000,DIRECT
	TypeUID RuleType = "UID"

	// TypeDSCP DSCP 值匹配
	// 示例: DSCP,46,PROXY
	TypeDSCP RuleType = "DSCP"

	// TypeProcessNameRegex 进程名正则匹配
	// 示例: PROCESS-NAME-REGEX,^chrome.*$,PROXY
	TypeProcessNameRegex RuleType = "PROCESS-NAME-REGEX"

	// TypeProcessPathRegex 进程路径正则匹配
	// 示例: PROCESS-PATH-REGEX,^/usr/bin/.*$,DIRECT
	TypeProcessPathRegex RuleType = "PROCESS-PATH-REGEX"

	// TypeIPSuffix IP 后缀匹配
	// 示例: IP-SUFFIX,8.8.8.8/24,1,DIRECT
	TypeIPSuffix RuleType = "IP-SUFFIX"

	// TypeGeoSite GeoSite 域名地理位置匹配
	// 示例: GEOSITE,cn,DIRECT
	TypeGeoSite RuleType = "GEOSITE"

	// TypeMatch 匹配所有流量（通常作为最后一条规则）
	// 示例: MATCH,PROXY
	TypeMatch RuleType = "MATCH"
)

// ActionType 规则动作类型
type ActionType string

const (
	// ActionProxy 使用代理
	ActionProxy ActionType = "PROXY"

	// ActionDirect 直连
	ActionDirect ActionType = "DIRECT"

	// ActionReject 拒绝连接
	ActionReject ActionType = "REJECT"

	// ActionRejectDrop 拒绝连接并丢弃数据包
	ActionRejectDrop ActionType = "REJECT-DROP"
)

// IsDirect 判断是否为直连动作
func (a ActionType) IsDirect() bool {
	return a == ActionDirect
}

// IsProxy 判断是否为代理动作
// 包括 PROXY 和所有自定义服务名称（OPENAI, NETFLIX, YOUTUBE 等）
func (a ActionType) IsProxy() bool {
	if a == ActionProxy {
		return true
	}
	// 不是直连也不是拒绝，则认为是代理类动作
	return !a.IsDirect() && !a.IsReject()
}

// IsReject 判断是否为拒绝动作
func (a ActionType) IsReject() bool {
	return a == ActionReject || a == ActionRejectDrop
}

// Name 获取动作的本地化名称
// 支持多语言，默认返回英文名称
func (a ActionType) Name(lang ...string) string {
	// 获取语言代码，默认为英文
	langCode := "en"
	if len(lang) > 0 {
		langCode = lang[0]
	}

	// 本地化映射
	names := map[string]map[ActionType]string{
		"en": {
			ActionProxy:      "Proxy",
			ActionDirect:     "Direct",
			ActionReject:     "Reject",
			ActionRejectDrop: "Reject Drop",
		},
		"zh": {
			ActionProxy:      "代理",
			ActionDirect:     "直连",
			ActionReject:     "拒绝",
			ActionRejectDrop: "拒绝丢弃",
		},
		"zh-CN": {
			ActionProxy:      "代理",
			ActionDirect:     "直连",
			ActionReject:     "拒绝",
			ActionRejectDrop: "拒绝丢弃",
		},
		"zh-TW": {
			ActionProxy:      "代理",
			ActionDirect:     "直連",
			ActionReject:     "拒絕",
			ActionRejectDrop: "拒絕丟棄",
		},
		"ja": {
			ActionProxy:      "プロキシ",
			ActionDirect:     "直接接続",
			ActionReject:     "拒否",
			ActionRejectDrop: "拒否してドロップ",
		},
	}

	// 查找对应语言的名称
	if langMap, ok := names[langCode]; ok {
		if name, ok := langMap[a]; ok {
			return name
		}
	}

	// 如果是自定义服务名称（如 OPENAI, NETFLIX 等），返回格式化的名称
	actionStr := string(a)

	// 对于标准动作，返回默认英文名称
	if a == ActionProxy || a == ActionDirect || a == ActionReject || a == ActionRejectDrop {
		if names["en"][a] != "" {
			return names["en"][a]
		}
	}

	// 对于自定义服务，返回首字母大写的格式
	if len(actionStr) > 0 {
		// 转换为标题格式：OPENAI -> OpenAI, YOUTUBE -> YouTube
		return formatServiceName(actionStr)
	}

	return actionStr
}

// formatServiceName 格式化服务名称
func formatServiceName(name string) string {
	// 特殊服务名称映射
	specialNames := map[string]string{
		"OPENAI":         "OpenAI",
		"CLAUDE":         "Claude",
		"GEMINI":         "Gemini",
		"YOUTUBE":        "YouTube",
		"NETFLIX":        "Netflix",
		"DISNEY":         "Disney+",
		"SPOTIFY":        "Spotify",
		"TIKTOK":         "TikTok",
		"TWITCH":         "Twitch",
		"HBO":            "HBO",
		"HULU":           "Hulu",
		"PRIME-VIDEO":    "Prime Video",
		"PANDORA":        "Pandora",
		"SOUNDCLOUD":     "SoundCloud",
		"DAZN":           "DAZN",
		"VIMEO":          "Vimeo",
		"BILIBILI":       "哔哩哔哩",
		"BILIBILI-HK":    "哔哩哔哩港澳台",
		"IQIYI":          "爱奇艺",
		"IQIYI-HK":       "爱奇艺港澳台",
		"TENCENT-VIDEO":  "腾讯视频",
		"YOUKU":          "优酷",
		"NETEASE-MUSIC":  "网易云音乐",
		"CCTV":           "CCTV",
		"DOUYU":          "斗鱼",
		"HIMALAYA":       "喜马拉雅",
		"APP-STORE":      "App Store",
		"ICLOUD":         "iCloud",
		"APPLE-TV":       "Apple TV",
		"APPLE-MUSIC":    "Apple Music",
		"TESTFLIGHT":     "TestFlight",
		"APPLE":          "Apple",
		"ONEDRIVE":       "OneDrive",
		"GDRIVE":         "Google Drive",
		"DROPBOX":        "Dropbox",
		"GOOGLE":         "Google",
		"MICROSOFT":      "Microsoft",
		"AMAZON":         "Amazon",
		"FACEBOOK":       "Facebook",
		"ADOBE":          "Adobe",
		"GITHUB":         "GitHub",
		"GITLAB":         "GitLab",
		"DOCKER":         "Docker",
		"HEROKU":         "Heroku",
		"DIGITALOCEAN":   "DigitalOcean",
		"VERCEL":         "Vercel",
		"CLOUDFLARE":     "Cloudflare",
		"BINANCE":        "Binance",
		"OKX":            "OKX",
		"CRYPTO":         "Crypto.com",
		"CRYPTOCURRENCY": "Cryptocurrency",
		"PAYPAL":         "PayPal",
		"TELEGRAM":       "Telegram",
		"TWITTER":        "Twitter",
		"INSTAGRAM":      "Instagram",
		"WHATSAPP":       "WhatsApp",
		"DISCORD":        "Discord",
		"LINE":           "Line",
		"THREADS":        "Threads",
		"REDDIT":         "Reddit",
		"LINKEDIN":       "LinkedIn",
		"WIKIPEDIA":      "Wikipedia",
		"STEAM":          "Steam",
		"EPIC":           "Epic Games",
		"PLAYSTATION":    "PlayStation",
		"EBAY":           "eBay",
		"SHOPIFY":        "Shopify",
		"BBC":            "BBC",
		"CNN":            "CNN",
		"BLOOMBERG":      "Bloomberg",
		"NYTIMES":        "New York Times",
		"SCHOLAR":        "Scholar",
	}

	if special, ok := specialNames[name]; ok {
		return special
	}

	// 默认：首字母大写，其余小写
	if len(name) == 0 {
		return name
	}

	// 如果全是大写，转换为首字母大写
	allUpper := true
	for _, c := range name {
		if c >= 'a' && c <= 'z' {
			allUpper = false
			break
		}
	}

	if allUpper && len(name) > 0 {
		return string(name[0]) + string(name[1:])
	}

	return name
}

// Metadata 连接元数据
//
// 包含连接的各种信息，用于规则匹配
type Metadata struct {
	// 网络信息
	Network     string      // 网络类型: tcp, udp
	Type        string      // 连接类型: http, https, socks5
	SrcIP       netip.Addr  // 源 IP
	DstIP       netip.Addr  // 目标 IP
	SrcPort     uint16      // 源端口
	DstPort     uint16      // 目标端口

	// 域名信息
	Host        string      // 主机名/域名
	Domain      string      // 规范化的域名（小写）

	// GeoIP 信息
	DstGeoIP    *GeoIPInfo  // 目标 IP 的 GeoIP 信息
	SrcGeoIP    *GeoIPInfo  // 源 IP 的 GeoIP 信息

	// 进程信息
	ProcessName string      // 进程名称
	ProcessPath string      // 进程路径

	// 入站信息
	InboundType string      // 入站类型: HTTP, SOCKS5, etc.
	InboundName string      // 入站名称
	InboundUser string      // 入站用户（认证用户名）

	// 系统信息
	UID         uint32      // 用户 ID (Linux/Android)
	DSCP        uint8       // DSCP 值 (Differentiated Services Code Point)
}

// GeoIPInfo GeoIP 信息
type GeoIPInfo struct {
	CountryCode   string // 国家代码
	Country       string // 国家名称
	Continent     string // 大洲
	ContinentCode string // 大洲代码
	ASN           int    // 自治系统号
	ASName        string // 自治系统名称
}

// Rule 规则接口
//
// 所有规则类型都需要实现此接口
type Rule interface {
	// Type 返回规则类型
	Type() RuleType

	// Match 判断元数据是否匹配规则
	Match(metadata *Metadata) bool

	// Action 返回规则动作
	Action() ActionType

	// Payload 返回规则载荷（匹配内容）
	Payload() string

	// String 返回规则的字符串表示
	String() string
}
