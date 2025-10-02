package main

import (
	"os"
	"time"

	"github.com/darabuchi/prism/scripts/tools/rule-collector/collector"
	"github.com/pterm/pterm"
)

const (
	BLACKMATRIX7 = "blackmatrix7"
	LOYALSOLDIER = "loyalsoldier"
	ACL4SSR      = "acl4ssr"
)

func main() {
	start := time.Now()
	defer func() {
		pterm.Success.Printfln("总耗时: %v", time.Since(start))
	}()

	pterm.Info.Printfln("PID: %d", os.Getpid())
	pterm.Info.Printfln("Prism 规则收集器")

	// 加载配置
	cfg, err := LoadConfig()
	if err != nil {
		pterm.Error.Printfln("加载配置失败: %v", err)
		os.Exit(1)
	}

	// 显示配置信息
	pterm.Info.Printfln("缓存目录: %s", cfg.CacheDir)
	pterm.Info.Printfln("输出目录: %s", cfg.OutputDir)
	pterm.Info.Printfln("缓存天数: %d", cfg.CacheDays)
	if cfg.Proxy != "" {
		pterm.Info.Printfln("代理地址: %s", cfg.Proxy)
	}

	// 创建收集器
	c := NewCollector(cfg)

	// 注册数据源收集器
	c.AddHandle(BLACKMATRIX7, collector.NewBlackMatrix7(cfg.Proxy, cfg.CacheDays))
	c.AddHandle(LOYALSOLDIER, collector.NewLoyalSoldier(cfg.Proxy, cfg.CacheDays))
	c.AddHandle(ACL4SSR, collector.NewACL4SSR(cfg.Proxy, cfg.CacheDays))

	// 规则解析配置
	type ParseConfig struct {
		Source   string // 数据源名称
		Path     string // 规则文件路径
		Action   string // 规则动作
		Priority int    // 优先级（可选，默认为 0）
	}

	parseList := []ParseConfig{
		// 直连规则
		{BLACKMATRIX7, "Lan/Lan.yaml", "DIRECT", 0},       // 本地局域网
		{BLACKMATRIX7, "Direct/Direct.yaml", "DIRECT", 0}, // 直连规则

		// ========== AI 服务（每个独立分类）==========
		{BLACKMATRIX7, "OpenAI/OpenAI.yaml", "OPENAI", 0},       // OpenAI
		{BLACKMATRIX7, "Claude/Claude.yaml", "CLAUDE", 0},       // Claude
		{BLACKMATRIX7, "Gemini/Gemini.yaml", "GEMINI", 0},       // Google Gemini
		{BLACKMATRIX7, "Anthropic/Anthropic.yaml", "CLAUDE", 0}, // Anthropic (归类到 Claude)

		// ========== 流媒体 - 国际（每个独立分类）==========
		{BLACKMATRIX7, "YouTube/YouTube.yaml", "YOUTUBE", 0},           // YouTube
		{BLACKMATRIX7, "YouTubeMusic/YouTubeMusic.yaml", "YOUTUBE", 0}, // YouTube Music
		{BLACKMATRIX7, "Netflix/Netflix.yaml", "NETFLIX", 0},           // Netflix
		{BLACKMATRIX7, "Disney/Disney.yaml", "DISNEY", 0},              // Disney+
		{BLACKMATRIX7, "Spotify/Spotify.yaml", "SPOTIFY", 0},           // Spotify
		{BLACKMATRIX7, "TikTok/TikTok.yaml", "TIKTOK", 0},              // TikTok
		{BLACKMATRIX7, "Twitch/Twitch.yaml", "TWITCH", 0},              // Twitch
		{BLACKMATRIX7, "HBO/HBO.yaml", "HBO", 0},                       // HBO
		{BLACKMATRIX7, "Hulu/Hulu.yaml", "HULU", 0},                    // Hulu
		{BLACKMATRIX7, "AmazonPrimeVideo/AmazonPrimeVideo.yaml", "PRIME-VIDEO", 0}, // Amazon Prime Video
		{BLACKMATRIX7, "Pandora/Pandora.yaml", "PANDORA", 0},                       // Pandora
		{BLACKMATRIX7, "SoundCloud/SoundCloud.yaml", "SOUNDCLOUD", 0},              // SoundCloud
		{BLACKMATRIX7, "DAZN/DAZN.yaml", "DAZN", 0},                                // DAZN
		{BLACKMATRIX7, "Vimeo/Vimeo.yaml", "VIMEO", 0},                             // Vimeo

		// ========== 流媒体 - 国内（每个独立分类）==========
		{BLACKMATRIX7, "BiliBili/BiliBili.yaml", "BILIBILI", 0},            // 哔哩哔哩
		{BLACKMATRIX7, "BiliBiliIntl/BiliBiliIntl.yaml", "BILIBILI-HK", 0}, // 哔哩哔哩国际版
		{BLACKMATRIX7, "iQIYI/iQIYI.yaml", "IQIYI", 0},                     // 爱奇艺
		{BLACKMATRIX7, "iQIYIIntl/iQIYIIntl.yaml", "IQIYI-HK", 0},          // 爱奇艺国际版
		{BLACKMATRIX7, "TencentVideo/TencentVideo.yaml", "TENCENT-VIDEO", 0},    // 腾讯视频
		{BLACKMATRIX7, "Youku/Youku.yaml", "YOUKU", 0},                          // 优酷
		{BLACKMATRIX7, "NetEaseMusic/NetEaseMusic.yaml", "NETEASE-MUSIC", 0},    // 网易云音乐
		{BLACKMATRIX7, "CCTV/CCTV.yaml", "CCTV", 0},                             // CCTV
		{BLACKMATRIX7, "Douyu/Douyu.yaml", "DOUYU", 0},                          // 斗鱼
		{BLACKMATRIX7, "Himalaya/Himalaya.yaml", "HIMALAYA", 0},                 // 喜马拉雅

		// ========== 苹果服务（细分）==========
		{BLACKMATRIX7, "AppStore/AppStore.yaml", "APP-STORE", 0}, // App Store
		{BLACKMATRIX7, "iCloud/iCloud.yaml", "ICLOUD", 0},        // iCloud
		{BLACKMATRIX7, "AppleTV/AppleTV.yaml", "APPLE-TV", 0},    // Apple TV
		{BLACKMATRIX7, "AppleMusic/AppleMusic.yaml", "APPLE-MUSIC", 0}, // Apple Music
		{BLACKMATRIX7, "TestFlight/TestFlight.yaml", "TESTFLIGHT", 0},  // TestFlight
		{BLACKMATRIX7, "Apple/Apple.yaml", "APPLE", 0},                 // Apple 其他服务

		// ========== 云存储（每个独立分类）==========
		{BLACKMATRIX7, "OneDrive/OneDrive.yaml", "ONEDRIVE", 0},       // OneDrive
		{BLACKMATRIX7, "GoogleDrive/GoogleDrive.yaml", "GDRIVE", 0},   // Google Drive
		{BLACKMATRIX7, "Dropbox/Dropbox.yaml", "DROPBOX", 0},          // Dropbox

		// ========== 科技公司 ==========
		{BLACKMATRIX7, "Google/Google.yaml", "GOOGLE", 0},         // Google
		{BLACKMATRIX7, "Microsoft/Microsoft.yaml", "MICROSOFT", 0}, // Microsoft
		{BLACKMATRIX7, "Amazon/Amazon.yaml", "AMAZON", 0},         // Amazon
		{BLACKMATRIX7, "Facebook/Facebook.yaml", "FACEBOOK", 0},   // Facebook
		{BLACKMATRIX7, "Adobe/Adobe.yaml", "ADOBE", 0},            // Adobe

		// ========== 开发工具/VPS（每个独立分类）==========
		{BLACKMATRIX7, "GitHub/GitHub.yaml", "GITHUB", 0},     // GitHub
		{BLACKMATRIX7, "GitLab/GitLab.yaml", "GITLAB", 0},     // GitLab
		{BLACKMATRIX7, "Docker/Docker.yaml", "DOCKER", 0},     // Docker Hub
		{BLACKMATRIX7, "Heroku/Heroku.yaml", "HEROKU", 0},     // Heroku
		{BLACKMATRIX7, "DigitalOcean/DigitalOcean.yaml", "DIGITALOCEAN", 0}, // DigitalOcean
		{BLACKMATRIX7, "Vercel/Vercel.yaml", "VERCEL", 0},     // Vercel
		{BLACKMATRIX7, "Cloudflare/Cloudflare.yaml", "CLOUDFLARE", 0}, // Cloudflare

		// ========== 交易所（每个独立分类）==========
		{BLACKMATRIX7, "Binance/Binance.yaml", "BINANCE", 0},         // 币安
		{BLACKMATRIX7, "OKX/OKX.yaml", "OKX", 0},                     // OKX
		{BLACKMATRIX7, "Crypto/Crypto.yaml", "CRYPTO", 0},            // Crypto.com
		{BLACKMATRIX7, "Cryptocurrency/Cryptocurrency.yaml", "CRYPTOCURRENCY", 0}, // 加密货币综合

		// ========== 支付（每个独立分类）==========
		{BLACKMATRIX7, "PayPal/PayPal.yaml", "PAYPAL", 0}, // PayPal

		// ========== 社交平台（每个独立分类）==========
		{BLACKMATRIX7, "Telegram/Telegram.yaml", "TELEGRAM", 0},   // Telegram
		{BLACKMATRIX7, "Twitter/Twitter.yaml", "TWITTER", 0},      // Twitter
		{BLACKMATRIX7, "Instagram/Instagram.yaml", "INSTAGRAM", 0}, // Instagram
		{BLACKMATRIX7, "WhatsApp/Whatsapp.yaml", "WHATSAPP", 0},   // WhatsApp
		{BLACKMATRIX7, "Discord/Discord.yaml", "DISCORD", 0},      // Discord
		{BLACKMATRIX7, "Line/Line.yaml", "LINE", 0},               // Line
		{BLACKMATRIX7, "Threads/Threads.yaml", "THREADS", 0},      // Threads
		{BLACKMATRIX7, "Reddit/Reddit.yaml", "REDDIT", 0},         // Reddit
		{BLACKMATRIX7, "LinkedIn/LinkedIn.yaml", "LINKEDIN", 0},   // LinkedIn

		// ========== 维基百科 ==========
		{BLACKMATRIX7, "Wikipedia/Wikipedia.yaml", "WIKIPEDIA", 0}, // Wikipedia

		// ========== 游戏 ==========
		{BLACKMATRIX7, "Steam/Steam.yaml", "STEAM", 0},   // Steam
		{BLACKMATRIX7, "Epic/Epic.yaml", "EPIC", 0},      // Epic Games
		{BLACKMATRIX7, "Sony/Sony.yaml", "PLAYSTATION", 0}, // PlayStation

		// ========== 购物 ==========
		{BLACKMATRIX7, "Amazon/Amazon.yaml", "AMAZON", 0},   // Amazon
		{BLACKMATRIX7, "eBay/eBay.yaml", "EBAY", 0},         // eBay
		{BLACKMATRIX7, "Shopify/Shopify.yaml", "SHOPIFY", 0}, // Shopify

		// ========== 咨询/新闻 ==========
		{BLACKMATRIX7, "BBC/BBC.yaml", "BBC", 0},             // BBC
		{BLACKMATRIX7, "CNN/CNN.yaml", "CNN", 0},             // CNN
		{BLACKMATRIX7, "Bloomberg/Bloomberg.yaml", "BLOOMBERG", 0}, // Bloomberg
		{BLACKMATRIX7, "NYTimes/NYTimes.yaml", "NYTIMES", 0}, // New York Times

		// ========== 政府相关 ==========
		{BLACKMATRIX7, "GlobalScholar/GlobalScholar.yaml", "SCHOLAR", 0}, // 学术/政府资源

		// 国内常用网站和服务
		{BLACKMATRIX7, "ChinaMax/ChinaMax.yaml", "DIRECT", 0}, // 国内网站合集
		{BLACKMATRIX7, "WeChat/WeChat.yaml", "DIRECT", 0},     // 微信
		{BLACKMATRIX7, "Weibo/Weibo.yaml", "DIRECT", 0},       // 微博

		// Loyalsoldier 规则
		{LOYALSOLDIER, "proxy.txt", "PROXY", 0},   // 代理域名列表
		{LOYALSOLDIER, "direct.txt", "DIRECT", 0}, // 直连域名列表
		{LOYALSOLDIER, "reject.txt", "REJECT", 0}, // 广告域名列表

		// 广告拦截
		{BLACKMATRIX7, "Advertising/Advertising.yaml", "REJECT", 0}, // 广告拦截

		// ACL4SSR 规则
		{ACL4SSR, "Clash/ProxyGFWlist.list", "PROXY", 0}, // GFW 列表
		{ACL4SSR, "Clash/ChinaDomain.list", "DIRECT", 0}, // 国内域名
	}

	// 加载自定义规则（在上游规则之前）
	pterm.Info.Printfln("加载自定义规则...")
	err = c.LoadRulesFromDirectory("rules")
	if err != nil {
		pterm.Warning.Printfln("加载自定义规则失败: %v", err)
		// 继续处理，不中断
	}

	// 执行规则解析
	pterm.Info.Printfln("开始收集规则...")
	for _, cfg := range parseList {
		err := c.Parse(cfg.Source, cfg.Path, cfg.Action, cfg.Priority)
		if err != nil {
			pterm.Warning.Printfln("解析 %s/%s 失败: %v", cfg.Source, cfg.Path, err)
			// 继续处理其他规则，不中断
			continue
		}
	}

	// 导出规则
	pterm.Info.Printfln("导出规则...")
	err = c.Export()
	if err != nil {
		pterm.Error.Printfln("导出规则失败: %v", err)
		os.Exit(1)
	}

	pterm.Success.Printfln("规则收集完成！")
	pterm.Info.Printfln("收集规则总数: %d", c.RuleCount())
}
