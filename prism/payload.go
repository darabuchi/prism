package prism

// Payload 规则分流目标类型
// 表示流量的处理方式或目标服务
type Payload string

const (
	// PayloadProxy 使用代理
	PayloadProxy Payload = "PROXY"

	// PayloadDirect 直连
	PayloadDirect Payload = "DIRECT"

	// PayloadReject 拒绝连接
	PayloadReject Payload = "REJECT"

	// PayloadRejectDrop 拒绝连接并丢弃数据包
	PayloadRejectDrop Payload = "REJECT-DROP"
)

// IsDirect 判断是否为直连动作
func (p Payload) IsDirect() bool {
	return p == PayloadDirect
}

// IsProxy 判断是否为代理动作
// 包括 PROXY 和所有自定义服务名称（OPENAI, NETFLIX, YOUTUBE 等）
func (p Payload) IsProxy() bool {
	if p == PayloadProxy {
		return true
	}
	// 不是直连也不是拒绝，则认为是代理类动作
	return !p.IsDirect() && !p.IsReject()
}

// IsReject 判断是否为拒绝动作
func (p Payload) IsReject() bool {
	return p == PayloadReject || p == PayloadRejectDrop
}

// Name 获取分流目标的本地化名称
// 支持多语言，默认返回英文名称
//
// 参数:
//   - lang: 可选的语言代码，如 "en", "zh", "zh-CN", "zh-TW", "ja"
//
// 返回:
//   - 本地化后的名称字符串
//
// 示例:
//
//	PayloadProxy.Name()           // "Proxy"
//	PayloadProxy.Name("zh")       // "代理"
//	Payload("OPENAI").Name()      // "OpenAI"
//	Payload("OPENAI").Name("zh")  // "OpenAI"
func (p Payload) Name(lang ...string) string {
	// 获取语言代码，默认为英文
	langCode := "en"
	if len(lang) > 0 {
		langCode = lang[0]
	}

	// 本地化映射
	names := map[string]map[Payload]string{
		"en": {
			PayloadProxy:      "Proxy",
			PayloadDirect:     "Direct",
			PayloadReject:     "Reject",
			PayloadRejectDrop: "Reject Drop",
		},
		"zh": {
			PayloadProxy:      "代理",
			PayloadDirect:     "直连",
			PayloadReject:     "拒绝",
			PayloadRejectDrop: "拒绝丢弃",
		},
		"zh-CN": {
			PayloadProxy:      "代理",
			PayloadDirect:     "直连",
			PayloadReject:     "拒绝",
			PayloadRejectDrop: "拒绝丢弃",
		},
		"zh-TW": {
			PayloadProxy:      "代理",
			PayloadDirect:     "直連",
			PayloadReject:     "拒絕",
			PayloadRejectDrop: "拒絕丟棄",
		},
		"ja": {
			PayloadProxy:      "プロキシ",
			PayloadDirect:     "直接接続",
			PayloadReject:     "拒否",
			PayloadRejectDrop: "拒否してドロップ",
		},
	}

	// 查找对应语言的名称
	if langMap, ok := names[langCode]; ok {
		if name, ok := langMap[p]; ok {
			return name
		}
	}

	// 如果是自定义服务名称（如 OPENAI, NETFLIX 等），返回格式化的名称
	payloadStr := string(p)

	// 对于标准动作，返回默认英文名称
	if p == PayloadProxy || p == PayloadDirect || p == PayloadReject || p == PayloadRejectDrop {
		if names["en"][p] != "" {
			return names["en"][p]
		}
	}

	// 对于自定义服务，返回首字母大写的格式
	if len(payloadStr) > 0 {
		// 转换为标题格式：OPENAI -> OpenAI, YOUTUBE -> YouTube
		return formatServiceName(payloadStr)
	}

	return payloadStr
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
