package prism

// Payload 规则分流目标类型
// 表示流量的处理方式或目标服务
type Payload int

const (
	// PayloadUnknown 未知类型
	PayloadUnknown Payload = 0

	// PayloadDirect 直连
	PayloadDirect Payload = 1

	// PayloadReject 拒绝连接
	PayloadReject Payload = 2

	// PayloadRejectDrop 拒绝连接并丢弃数据包
	PayloadRejectDrop Payload = 3

	// PayloadProxy 使用代理
	PayloadProxy Payload = 4

	// ========== AI 服务 ==========
	// PayloadOpenAI OpenAI 服务
	PayloadOpenAI Payload = 100

	// PayloadClaude Claude 服务
	PayloadClaude Payload = 101

	// PayloadGemini Google Gemini 服务
	PayloadGemini Payload = 102

	// PayloadCopilot GitHub Copilot 服务
	PayloadCopilot Payload = 103

	// PayloadBing Bing AI / Copilot 服务
	PayloadBing Payload = 104

	// PayloadPerplexity Perplexity AI 服务
	PayloadPerplexity Payload = 105

	// PayloadCharacterAI Character.AI 服务
	PayloadCharacterAI Payload = 106

	// PayloadMidjourney Midjourney AI 绘图
	PayloadMidjourney Payload = 107

	// PayloadStableDiffusion Stable Diffusion AI 绘图
	PayloadStableDiffusion Payload = 108

	// PayloadHuggingFace Hugging Face AI 平台
	PayloadHuggingFace Payload = 109

	// PayloadCohere Cohere AI 服务
	PayloadCohere Payload = 110

	// PayloadMistralAI Mistral AI 服务
	PayloadMistralAI Payload = 111

	// PayloadPoe Poe AI 聚合平台
	PayloadPoe Payload = 112

	// PayloadNotionAI Notion AI 服务
	PayloadNotionAI Payload = 113

	// PayloadJasper Jasper AI 写作助手
	PayloadJasper Payload = 114

	// PayloadChatGPT ChatGPT 服务（独立于 OpenAI）
	PayloadChatGPT Payload = 115

	// PayloadBard Google Bard 服务（已合并到 Gemini）
	PayloadBard Payload = 116

	// PayloadLlama Meta Llama AI 服务
	PayloadLlama Payload = 117

	// PayloadReplicate Replicate AI 平台
	PayloadReplicate Payload = 118

	// PayloadRunwayML Runway ML AI 视频生成
	PayloadRunwayML Payload = 119

	// ========== 流媒体 - 国际 ==========
	// PayloadYouTube YouTube 视频服务
	PayloadYouTube Payload = 200

	// PayloadNetflix Netflix 流媒体
	PayloadNetflix Payload = 201

	// PayloadDisney Disney+ 流媒体
	PayloadDisney Payload = 202

	// PayloadSpotify Spotify 音乐服务
	PayloadSpotify Payload = 203

	// PayloadTikTok TikTok 短视频
	PayloadTikTok Payload = 204

	// PayloadTwitch Twitch 直播平台
	PayloadTwitch Payload = 205

	// PayloadHBO HBO 流媒体
	PayloadHBO Payload = 206

	// PayloadHulu Hulu 流媒体
	PayloadHulu Payload = 207

	// PayloadPrimeVideo Amazon Prime Video
	PayloadPrimeVideo Payload = 208

	// PayloadPandora Pandora 音乐
	PayloadPandora Payload = 209

	// PayloadSoundCloud SoundCloud 音乐
	PayloadSoundCloud Payload = 210

	// PayloadDAZN DAZN 体育流媒体
	PayloadDAZN Payload = 211

	// PayloadVimeo Vimeo 视频平台
	PayloadVimeo Payload = 212

	// ========== 流媒体 - 国内 ==========
	// PayloadBilibili 哔哩哔哩
	PayloadBilibili Payload = 300

	// PayloadBilibiliHK 哔哩哔哩港澳台
	PayloadBilibiliHK Payload = 301

	// PayloadIQIYI 爱奇艺
	PayloadIQIYI Payload = 302

	// PayloadIQIYIHK 爱奇艺港澳台
	PayloadIQIYIHK Payload = 303

	// PayloadTencentVideo 腾讯视频
	PayloadTencentVideo Payload = 304

	// PayloadYouku 优酷
	PayloadYouku Payload = 305

	// PayloadNeteaseMusic 网易云音乐
	PayloadNeteaseMusic Payload = 306

	// PayloadCCTV CCTV
	PayloadCCTV Payload = 307

	// PayloadDouyu 斗鱼
	PayloadDouyu Payload = 308

	// PayloadHimalaya 喜马拉雅
	PayloadHimalaya Payload = 309

	// ========== 苹果服务 ==========
	// PayloadAppStore App Store
	PayloadAppStore Payload = 400

	// PayloadICloud iCloud
	PayloadICloud Payload = 401

	// PayloadAppleTV Apple TV
	PayloadAppleTV Payload = 402

	// PayloadAppleMusic Apple Music
	PayloadAppleMusic Payload = 403

	// PayloadTestFlight TestFlight
	PayloadTestFlight Payload = 404

	// PayloadApple Apple 其他服务
	PayloadApple Payload = 405

	// ========== 云存储 ==========
	// PayloadOneDrive OneDrive
	PayloadOneDrive Payload = 500

	// PayloadGoogleDrive Google Drive
	PayloadGoogleDrive Payload = 501

	// PayloadDropbox Dropbox
	PayloadDropbox Payload = 502

	// ========== 科技公司 ==========
	// PayloadGoogle Google
	PayloadGoogle Payload = 600

	// PayloadMicrosoft Microsoft
	PayloadMicrosoft Payload = 601

	// PayloadAmazon Amazon
	PayloadAmazon Payload = 602

	// PayloadFacebook Facebook
	PayloadFacebook Payload = 603

	// PayloadAdobe Adobe
	PayloadAdobe Payload = 604

	// ========== 开发工具/VPS ==========
	// PayloadGitHub GitHub
	PayloadGitHub Payload = 700

	// PayloadGitLab GitLab
	PayloadGitLab Payload = 701

	// PayloadDocker Docker Hub
	PayloadDocker Payload = 702

	// PayloadHeroku Heroku
	PayloadHeroku Payload = 703

	// PayloadDigitalOcean DigitalOcean
	PayloadDigitalOcean Payload = 704

	// PayloadVercel Vercel
	PayloadVercel Payload = 705

	// PayloadCloudflare Cloudflare
	PayloadCloudflare Payload = 706

	// ========== 交易所 ==========
	// PayloadBinance Binance 币安
	PayloadBinance Payload = 800

	// PayloadOKX OKX 交易所
	PayloadOKX Payload = 801

	// PayloadCrypto Crypto.com
	PayloadCrypto Payload = 802

	// PayloadCryptocurrency 加密货币综合
	PayloadCryptocurrency Payload = 803

	// ========== 支付 ==========
	// PayloadPayPal PayPal
	PayloadPayPal Payload = 900

	// ========== 社交平台 ==========
	// PayloadTelegram Telegram
	PayloadTelegram Payload = 1000

	// PayloadTwitter Twitter
	PayloadTwitter Payload = 1001

	// PayloadInstagram Instagram
	PayloadInstagram Payload = 1002

	// PayloadWhatsApp WhatsApp
	PayloadWhatsApp Payload = 1003

	// PayloadDiscord Discord
	PayloadDiscord Payload = 1004

	// PayloadLine Line
	PayloadLine Payload = 1005

	// PayloadThreads Threads
	PayloadThreads Payload = 1006

	// PayloadReddit Reddit
	PayloadReddit Payload = 1007

	// PayloadLinkedIn LinkedIn
	PayloadLinkedIn Payload = 1008

	// ========== 其他服务 ==========
	// PayloadWikipedia Wikipedia
	PayloadWikipedia Payload = 1100

	// PayloadSteam Steam
	PayloadSteam Payload = 1101

	// PayloadEpic Epic Games
	PayloadEpic Payload = 1102

	// PayloadPlayStation PlayStation
	PayloadPlayStation Payload = 1103

	// PayloadEbay eBay
	PayloadEbay Payload = 1104

	// PayloadShopify Shopify
	PayloadShopify Payload = 1105

	// PayloadBBC BBC
	PayloadBBC Payload = 1106

	// PayloadCNN CNN
	PayloadCNN Payload = 1107

	// PayloadBloomberg Bloomberg
	PayloadBloomberg Payload = 1108

	// PayloadNYTimes New York Times
	PayloadNYTimes Payload = 1109

	// PayloadScholar 学术/政府资源
	PayloadScholar Payload = 1110
)

// ParsePayload 从字符串解析 Payload
// 支持从规则配置文件中的字符串转换为 Payload 枚举
func ParsePayload(s string) Payload {
	payloadMap := map[string]Payload{
		// 基础动作
		"UNKNOWN":      PayloadUnknown,
		"DIRECT":       PayloadDirect,
		"REJECT":       PayloadReject,
		"REJECT-DROP":  PayloadRejectDrop,
		"PROXY":        PayloadProxy,

		// AI 服务
		"OPENAI":           PayloadOpenAI,
		"CLAUDE":           PayloadClaude,
		"GEMINI":           PayloadGemini,
		"COPILOT":          PayloadCopilot,
		"BING":             PayloadBing,
		"PERPLEXITY":       PayloadPerplexity,
		"CHARACTER-AI":     PayloadCharacterAI,
		"CHARACTERAI":      PayloadCharacterAI,
		"MIDJOURNEY":       PayloadMidjourney,
		"STABLE-DIFFUSION": PayloadStableDiffusion,
		"HUGGINGFACE":      PayloadHuggingFace,
		"COHERE":           PayloadCohere,
		"MISTRAL":          PayloadMistralAI,
		"MISTRAL-AI":       PayloadMistralAI,
		"POE":              PayloadPoe,
		"NOTION-AI":        PayloadNotionAI,
		"NOTIONAI":         PayloadNotionAI,
		"JASPER":           PayloadJasper,
		"CHATGPT":          PayloadChatGPT,
		"BARD":             PayloadBard,
		"LLAMA":            PayloadLlama,
		"REPLICATE":        PayloadReplicate,
		"RUNWAY":           PayloadRunwayML,
		"RUNWAYML":         PayloadRunwayML,

		// 国际流媒体
		"YOUTUBE":      PayloadYouTube,
		"NETFLIX":      PayloadNetflix,
		"DISNEY":       PayloadDisney,
		"SPOTIFY":      PayloadSpotify,
		"TIKTOK":       PayloadTikTok,
		"TWITCH":       PayloadTwitch,
		"HBO":          PayloadHBO,
		"HULU":         PayloadHulu,
		"PRIME-VIDEO":  PayloadPrimeVideo,
		"PANDORA":      PayloadPandora,
		"SOUNDCLOUD":   PayloadSoundCloud,
		"DAZN":         PayloadDAZN,
		"VIMEO":        PayloadVimeo,

		// 国内流媒体
		"BILIBILI":     PayloadBilibili,
		"BILIBILI-HK":  PayloadBilibiliHK,
		"IQIYI":        PayloadIQIYI,
		"IQIYI-HK":     PayloadIQIYIHK,
		"TENCENT-VIDEO": PayloadTencentVideo,
		"YOUKU":        PayloadYouku,
		"NETEASE-MUSIC": PayloadNeteaseMusic,
		"CCTV":         PayloadCCTV,
		"DOUYU":        PayloadDouyu,
		"HIMALAYA":     PayloadHimalaya,

		// 苹果服务
		"APP-STORE":    PayloadAppStore,
		"ICLOUD":       PayloadICloud,
		"APPLE-TV":     PayloadAppleTV,
		"APPLE-MUSIC":  PayloadAppleMusic,
		"TESTFLIGHT":   PayloadTestFlight,
		"APPLE":        PayloadApple,

		// 云存储
		"ONEDRIVE":     PayloadOneDrive,
		"GDRIVE":       PayloadGoogleDrive,
		"DROPBOX":      PayloadDropbox,

		// 科技公司
		"GOOGLE":       PayloadGoogle,
		"MICROSOFT":    PayloadMicrosoft,
		"AMAZON":       PayloadAmazon,
		"FACEBOOK":     PayloadFacebook,
		"ADOBE":        PayloadAdobe,

		// 开发工具/VPS
		"GITHUB":       PayloadGitHub,
		"GITLAB":       PayloadGitLab,
		"DOCKER":       PayloadDocker,
		"HEROKU":       PayloadHeroku,
		"DIGITALOCEAN": PayloadDigitalOcean,
		"VERCEL":       PayloadVercel,
		"CLOUDFLARE":   PayloadCloudflare,

		// 交易所
		"BINANCE":      PayloadBinance,
		"OKX":          PayloadOKX,
		"CRYPTO":       PayloadCrypto,
		"CRYPTOCURRENCY": PayloadCryptocurrency,

		// 支付
		"PAYPAL":       PayloadPayPal,

		// 社交平台
		"TELEGRAM":     PayloadTelegram,
		"TWITTER":      PayloadTwitter,
		"INSTAGRAM":    PayloadInstagram,
		"WHATSAPP":     PayloadWhatsApp,
		"DISCORD":      PayloadDiscord,
		"LINE":         PayloadLine,
		"THREADS":      PayloadThreads,
		"REDDIT":       PayloadReddit,
		"LINKEDIN":     PayloadLinkedIn,

		// 其他服务
		"WIKIPEDIA":    PayloadWikipedia,
		"STEAM":        PayloadSteam,
		"EPIC":         PayloadEpic,
		"PLAYSTATION":  PayloadPlayStation,
		"EBAY":         PayloadEbay,
		"SHOPIFY":      PayloadShopify,
		"BBC":          PayloadBBC,
		"CNN":          PayloadCNN,
		"BLOOMBERG":    PayloadBloomberg,
		"NYTIMES":      PayloadNYTimes,
		"SCHOLAR":      PayloadScholar,
	}

	if payload, ok := payloadMap[s]; ok {
		return payload
	}

	return PayloadUnknown
}

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

// String 返回 Payload 的字符串表示
func (p Payload) String() string {
	return p.Name()
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
//	PayloadProxy.Name()        // "Proxy"
//	PayloadProxy.Name("zh")    // "代理"
//	PayloadOpenAI.Name()       // "OpenAI"
//	PayloadOpenAI.Name("zh")   // "OpenAI"
func (p Payload) Name(lang ...string) string {
	// 获取语言代码，默认为英文
	langCode := "en"
	if len(lang) > 0 {
		langCode = lang[0]
	}

	// 完整的本地化名称映射
	names := getPayloadNames()

	// 查找对应语言的名称
	if langMap, ok := names[langCode]; ok {
		if name, ok := langMap[p]; ok {
			return name
		}
	}

	// 如果没有找到，返回英文名称
	if name, ok := names["en"][p]; ok {
		return name
	}

	// 兜底返回
	return "Unknown"
}

// getPayloadNames 获取所有 Payload 的本地化名称映射
func getPayloadNames() map[string]map[Payload]string {
	return map[string]map[Payload]string{
		"en": {
			PayloadUnknown:         "Unknown",
			PayloadDirect:          "Direct",
			PayloadReject:          "Reject",
			PayloadRejectDrop:      "Reject Drop",
			PayloadProxy:           "Proxy",
			PayloadOpenAI:          "OpenAI",
			PayloadClaude:          "Claude",
			PayloadGemini:          "Gemini",
			PayloadCopilot:         "Copilot",
			PayloadBing:            "Bing AI",
			PayloadPerplexity:      "Perplexity",
			PayloadCharacterAI:     "Character.AI",
			PayloadMidjourney:      "Midjourney",
			PayloadStableDiffusion: "Stable Diffusion",
			PayloadHuggingFace:     "Hugging Face",
			PayloadCohere:          "Cohere",
			PayloadMistralAI:       "Mistral AI",
			PayloadPoe:             "Poe",
			PayloadNotionAI:        "Notion AI",
			PayloadJasper:          "Jasper",
			PayloadChatGPT:         "ChatGPT",
			PayloadBard:            "Bard",
			PayloadLlama:           "Llama",
			PayloadReplicate:       "Replicate",
			PayloadRunwayML:        "Runway ML",
			PayloadYouTube:         "YouTube",
			PayloadNetflix:         "Netflix",
			PayloadDisney:          "Disney+",
			PayloadSpotify:         "Spotify",
			PayloadTikTok:          "TikTok",
			PayloadTwitch:          "Twitch",
			PayloadHBO:             "HBO",
			PayloadHulu:            "Hulu",
			PayloadPrimeVideo:      "Prime Video",
			PayloadPandora:         "Pandora",
			PayloadSoundCloud:      "SoundCloud",
			PayloadDAZN:            "DAZN",
			PayloadVimeo:           "Vimeo",
			PayloadBilibili:        "Bilibili",
			PayloadBilibiliHK:      "Bilibili HK",
			PayloadIQIYI:           "iQIYI",
			PayloadIQIYIHK:         "iQIYI HK",
			PayloadTencentVideo:    "Tencent Video",
			PayloadYouku:           "Youku",
			PayloadNeteaseMusic:    "Netease Music",
			PayloadCCTV:            "CCTV",
			PayloadDouyu:           "Douyu",
			PayloadHimalaya:        "Himalaya",
			PayloadAppStore:        "App Store",
			PayloadICloud:          "iCloud",
			PayloadAppleTV:         "Apple TV",
			PayloadAppleMusic:      "Apple Music",
			PayloadTestFlight:      "TestFlight",
			PayloadApple:           "Apple",
			PayloadOneDrive:        "OneDrive",
			PayloadGoogleDrive:     "Google Drive",
			PayloadDropbox:         "Dropbox",
			PayloadGoogle:          "Google",
			PayloadMicrosoft:       "Microsoft",
			PayloadAmazon:          "Amazon",
			PayloadFacebook:        "Facebook",
			PayloadAdobe:           "Adobe",
			PayloadGitHub:          "GitHub",
			PayloadGitLab:          "GitLab",
			PayloadDocker:          "Docker",
			PayloadHeroku:          "Heroku",
			PayloadDigitalOcean:    "DigitalOcean",
			PayloadVercel:          "Vercel",
			PayloadCloudflare:      "Cloudflare",
			PayloadBinance:         "Binance",
			PayloadOKX:             "OKX",
			PayloadCrypto:          "Crypto.com",
			PayloadCryptocurrency:  "Cryptocurrency",
			PayloadPayPal:          "PayPal",
			PayloadTelegram:        "Telegram",
			PayloadTwitter:         "Twitter",
			PayloadInstagram:       "Instagram",
			PayloadWhatsApp:        "WhatsApp",
			PayloadDiscord:         "Discord",
			PayloadLine:            "Line",
			PayloadThreads:         "Threads",
			PayloadReddit:          "Reddit",
			PayloadLinkedIn:        "LinkedIn",
			PayloadWikipedia:       "Wikipedia",
			PayloadSteam:           "Steam",
			PayloadEpic:            "Epic Games",
			PayloadPlayStation:     "PlayStation",
			PayloadEbay:            "eBay",
			PayloadShopify:         "Shopify",
			PayloadBBC:             "BBC",
			PayloadCNN:             "CNN",
			PayloadBloomberg:       "Bloomberg",
			PayloadNYTimes:         "New York Times",
			PayloadScholar:         "Scholar",
		},
		"zh": {
			PayloadUnknown:         "未知",
			PayloadDirect:          "直连",
			PayloadReject:          "拒绝",
			PayloadRejectDrop:      "拒绝丢弃",
			PayloadProxy:           "代理",
			PayloadOpenAI:          "OpenAI",
			PayloadClaude:          "Claude",
			PayloadGemini:          "Gemini",
			PayloadCopilot:         "Copilot",
			PayloadBing:            "必应AI",
			PayloadPerplexity:      "Perplexity",
			PayloadCharacterAI:     "Character.AI",
			PayloadMidjourney:      "Midjourney",
			PayloadStableDiffusion: "Stable Diffusion",
			PayloadHuggingFace:     "Hugging Face",
			PayloadCohere:          "Cohere",
			PayloadMistralAI:       "Mistral AI",
			PayloadPoe:             "Poe",
			PayloadNotionAI:        "Notion AI",
			PayloadJasper:          "Jasper",
			PayloadChatGPT:         "ChatGPT",
			PayloadBard:            "Bard",
			PayloadLlama:           "Llama",
			PayloadReplicate:       "Replicate",
			PayloadRunwayML:        "Runway ML",
			PayloadYouTube:         "YouTube",
			PayloadNetflix:         "Netflix",
			PayloadDisney:          "Disney+",
			PayloadSpotify:         "Spotify",
			PayloadTikTok:          "TikTok",
			PayloadTwitch:          "Twitch",
			PayloadHBO:             "HBO",
			PayloadHulu:            "Hulu",
			PayloadPrimeVideo:      "Prime Video",
			PayloadPandora:         "Pandora",
			PayloadSoundCloud:      "SoundCloud",
			PayloadDAZN:            "DAZN",
			PayloadVimeo:           "Vimeo",
			PayloadBilibili:        "哔哩哔哩",
			PayloadBilibiliHK:      "哔哩哔哩港澳台",
			PayloadIQIYI:           "爱奇艺",
			PayloadIQIYIHK:         "爱奇艺港澳台",
			PayloadTencentVideo:    "腾讯视频",
			PayloadYouku:           "优酷",
			PayloadNeteaseMusic:    "网易云音乐",
			PayloadCCTV:            "CCTV",
			PayloadDouyu:           "斗鱼",
			PayloadHimalaya:        "喜马拉雅",
			PayloadAppStore:        "App Store",
			PayloadICloud:          "iCloud",
			PayloadAppleTV:         "Apple TV",
			PayloadAppleMusic:      "Apple Music",
			PayloadTestFlight:      "TestFlight",
			PayloadApple:           "Apple",
			PayloadOneDrive:        "OneDrive",
			PayloadGoogleDrive:     "Google Drive",
			PayloadDropbox:         "Dropbox",
			PayloadGoogle:          "Google",
			PayloadMicrosoft:       "微软",
			PayloadAmazon:          "亚马逊",
			PayloadFacebook:        "Facebook",
			PayloadAdobe:           "Adobe",
			PayloadGitHub:          "GitHub",
			PayloadGitLab:          "GitLab",
			PayloadDocker:          "Docker",
			PayloadHeroku:          "Heroku",
			PayloadDigitalOcean:    "DigitalOcean",
			PayloadVercel:          "Vercel",
			PayloadCloudflare:      "Cloudflare",
			PayloadBinance:         "币安",
			PayloadOKX:             "OKX",
			PayloadCrypto:          "Crypto.com",
			PayloadCryptocurrency:  "加密货币",
			PayloadPayPal:          "PayPal",
			PayloadTelegram:        "Telegram",
			PayloadTwitter:         "Twitter",
			PayloadInstagram:       "Instagram",
			PayloadWhatsApp:        "WhatsApp",
			PayloadDiscord:         "Discord",
			PayloadLine:            "Line",
			PayloadThreads:         "Threads",
			PayloadReddit:          "Reddit",
			PayloadLinkedIn:        "LinkedIn",
			PayloadWikipedia:       "维基百科",
			PayloadSteam:           "Steam",
			PayloadEpic:            "Epic Games",
			PayloadPlayStation:     "PlayStation",
			PayloadEbay:            "eBay",
			PayloadShopify:         "Shopify",
			PayloadBBC:             "BBC",
			PayloadCNN:             "CNN",
			PayloadBloomberg:       "彭博社",
			PayloadNYTimes:         "纽约时报",
			PayloadScholar:         "学术资源",
		},
		"zh-CN": {
			PayloadUnknown:         "未知",
			PayloadDirect:          "直连",
			PayloadReject:          "拒绝",
			PayloadRejectDrop:      "拒绝丢弃",
			PayloadProxy:           "代理",
			PayloadBilibili:        "哔哩哔哩",
			PayloadBilibiliHK:      "哔哩哔哩港澳台",
			PayloadIQIYI:           "爱奇艺",
			PayloadIQIYIHK:         "爱奇艺港澳台",
			PayloadTencentVideo:    "腾讯视频",
			PayloadYouku:           "优酷",
			PayloadNeteaseMusic:    "网易云音乐",
			PayloadDouyu:           "斗鱼",
			PayloadHimalaya:        "喜马拉雅",
			PayloadMicrosoft:       "微软",
			PayloadAmazon:          "亚马逊",
			PayloadBinance:         "币安",
			PayloadCryptocurrency:  "加密货币",
			PayloadWikipedia:       "维基百科",
			PayloadBloomberg:       "彭博社",
			PayloadNYTimes:         "纽约时报",
			PayloadScholar:         "学术资源",
		},
		"zh-TW": {
			PayloadUnknown:         "未知",
			PayloadDirect:          "直連",
			PayloadReject:          "拒絕",
			PayloadRejectDrop:      "拒絕丟棄",
			PayloadProxy:           "代理",
			PayloadBilibili:        "嗶哩嗶哩",
			PayloadBilibiliHK:      "嗶哩嗶哩港澳台",
			PayloadIQIYI:           "愛奇藝",
			PayloadIQIYIHK:         "愛奇藝港澳台",
			PayloadTencentVideo:    "騰訊視頻",
			PayloadYouku:           "優酷",
			PayloadNeteaseMusic:    "網易雲音樂",
			PayloadDouyu:           "鬥魚",
			PayloadHimalaya:        "喜馬拉雅",
			PayloadMicrosoft:       "微軟",
			PayloadAmazon:          "亞馬遜",
			PayloadBinance:         "幣安",
			PayloadCryptocurrency:  "加密貨幣",
			PayloadWikipedia:       "維基百科",
			PayloadBloomberg:       "彭博社",
			PayloadNYTimes:         "紐約時報",
			PayloadScholar:         "學術資源",
		},
		"ja": {
			PayloadUnknown:         "不明",
			PayloadDirect:          "直接接続",
			PayloadReject:          "拒否",
			PayloadRejectDrop:      "拒否してドロップ",
			PayloadProxy:           "プロキシ",
			PayloadMicrosoft:       "マイクロソフト",
			PayloadAmazon:          "アマゾン",
			PayloadCryptocurrency:  "暗号通貨",
			PayloadWikipedia:       "ウィキペディア",
			PayloadBloomberg:       "ブルームバーグ",
			PayloadScholar:         "学術リソース",
		},
	}
}

// PayloadCategory 服务类别
type PayloadCategory int

const (
	CategoryUnknown    PayloadCategory = 0
	CategoryBase       PayloadCategory = 1 // 基础动作（Direct、Reject、Proxy）
	CategoryAI         PayloadCategory = 2 // AI 服务
	CategoryStreaming  PayloadCategory = 3 // 流媒体服务
	CategorySocial     PayloadCategory = 4 // 社交平台
	CategoryDev        PayloadCategory = 5 // 开发工具
	CategoryGaming     PayloadCategory = 6 // 游戏平台
	CategoryNews       PayloadCategory = 7 // 新闻资讯
	CategoryFinance    PayloadCategory = 8 // 金融服务
	CategoryEducation  PayloadCategory = 9 // 教育学术
	CategoryOther      PayloadCategory = 10 // 其他服务
)

// Category 返回 Payload 的服务类别
func (p Payload) Category() PayloadCategory {
	switch {
	// 基础动作 (0-4)
	case p >= PayloadUnknown && p <= PayloadProxy:
		return CategoryBase

	// AI 服务 (100-119)
	case p >= PayloadOpenAI && p <= PayloadRunwayML:
		return CategoryAI

	// 流媒体服务 (200-309)
	case p >= PayloadYouTube && p <= PayloadYouku:
		return CategoryStreaming

	// 开发工具 (500-509, 700-706)
	case (p >= PayloadAppStore && p <= PayloadDropbox) || (p >= PayloadGitHub && p <= PayloadCloudflare):
		return CategoryDev

	// 金融服务 (800-803, 900)
	case (p >= PayloadBinance && p <= PayloadCryptocurrency) || p == PayloadPayPal:
		return CategoryFinance

	// 社交平台 (1000-1008)
	case p >= PayloadTelegram && p <= PayloadLinkedIn:
		return CategorySocial

	// 游戏平台 (1101-1103)
	case p >= PayloadSteam && p <= PayloadPlayStation:
		return CategoryGaming

	// 新闻资讯 (1108-1109)
	case p >= PayloadBloomberg && p <= PayloadNYTimes:
		return CategoryNews

	// 教育学术 (1110)
	case p == PayloadScholar:
		return CategoryEducation

	default:
		return CategoryOther
	}
}

// IsAIService 判断是否为 AI 服务
func (p Payload) IsAIService() bool {
	return p.Category() == CategoryAI
}

// IsStreamingService 判断是否为流媒体服务
func (p Payload) IsStreamingService() bool {
	return p.Category() == CategoryStreaming
}

// IsSocialService 判断是否为社交平台
func (p Payload) IsSocialService() bool {
	return p.Category() == CategorySocial
}

// IsDevService 判断是否为开发工具
func (p Payload) IsDevService() bool {
	return p.Category() == CategoryDev
}

// IsGamingService 判断是否为游戏平台
func (p Payload) IsGamingService() bool {
	return p.Category() == CategoryGaming
}

// IsBaseAction 判断是否为基础动作
func (p Payload) IsBaseAction() bool {
	return p.Category() == CategoryBase
}
