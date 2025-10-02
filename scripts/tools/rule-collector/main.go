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
		Source string
		Path   string
		Action string
	}

	parseList := []ParseConfig{
		// 直连规则
		{BLACKMATRIX7, "Lan/Lan.yaml", "DIRECT"},       // 本地局域网
		{BLACKMATRIX7, "Direct/Direct.yaml", "DIRECT"}, // 直连规则

		// ========== AI 服务（每个独立分类）==========
		{BLACKMATRIX7, "OpenAI/OpenAI.yaml", "OPENAI"},       // OpenAI
		{BLACKMATRIX7, "Claude/Claude.yaml", "CLAUDE"},       // Claude
		{BLACKMATRIX7, "Gemini/Gemini.yaml", "GEMINI"},       // Google Gemini
		{BLACKMATRIX7, "Anthropic/Anthropic.yaml", "CLAUDE"}, // Anthropic (归类到 Claude)

		// ========== 流媒体 - 国际（每个独立分类）==========
		{BLACKMATRIX7, "YouTube/YouTube.yaml", "YOUTUBE"},           // YouTube
		{BLACKMATRIX7, "YouTubeMusic/YouTubeMusic.yaml", "YOUTUBE"}, // YouTube Music
		{BLACKMATRIX7, "Netflix/Netflix.yaml", "NETFLIX"},           // Netflix
		{BLACKMATRIX7, "Disney/Disney.yaml", "DISNEY"},              // Disney+
		{BLACKMATRIX7, "Spotify/Spotify.yaml", "SPOTIFY"},           // Spotify
		{BLACKMATRIX7, "TikTok/TikTok.yaml", "TIKTOK"},              // TikTok
		{BLACKMATRIX7, "Twitch/Twitch.yaml", "TWITCH"},              // Twitch
		{BLACKMATRIX7, "HBO/HBO.yaml", "HBO"},                       // HBO
		{BLACKMATRIX7, "Hulu/Hulu.yaml", "HULU"},                    // Hulu
		{BLACKMATRIX7, "AmazonPrimeVideo/AmazonPrimeVideo.yaml", "PRIME-VIDEO"}, // Amazon Prime Video
		{BLACKMATRIX7, "Pandora/Pandora.yaml", "PANDORA"},                       // Pandora
		{BLACKMATRIX7, "SoundCloud/SoundCloud.yaml", "SOUNDCLOUD"},              // SoundCloud
		{BLACKMATRIX7, "DAZN/DAZN.yaml", "DAZN"},                                // DAZN
		{BLACKMATRIX7, "Vimeo/Vimeo.yaml", "VIMEO"},                             // Vimeo

		// ========== 流媒体 - 国内（每个独立分类）==========
		{BLACKMATRIX7, "BiliBili/BiliBili.yaml", "BILIBILI"},            // 哔哩哔哩
		{BLACKMATRIX7, "BiliBiliIntl/BiliBiliIntl.yaml", "BILIBILI-HK"}, // 哔哩哔哩国际版
		{BLACKMATRIX7, "iQIYI/iQIYI.yaml", "IQIYI"},                     // 爱奇艺
		{BLACKMATRIX7, "iQIYIIntl/iQIYIIntl.yaml", "IQIYI-HK"},          // 爱奇艺国际版
		{BLACKMATRIX7, "TencentVideo/TencentVideo.yaml", "TENCENT-VIDEO"},    // 腾讯视频
		{BLACKMATRIX7, "Youku/Youku.yaml", "YOUKU"},                          // 优酷
		{BLACKMATRIX7, "NetEaseMusic/NetEaseMusic.yaml", "NETEASE-MUSIC"},    // 网易云音乐
		{BLACKMATRIX7, "CCTV/CCTV.yaml", "CCTV"},                             // CCTV
		{BLACKMATRIX7, "Douyu/Douyu.yaml", "DOUYU"},                          // 斗鱼
		{BLACKMATRIX7, "Himalaya/Himalaya.yaml", "HIMALAYA"},                 // 喜马拉雅

		// ========== 苹果服务（细分）==========
		{BLACKMATRIX7, "AppStore/AppStore.yaml", "APP-STORE"}, // App Store
		{BLACKMATRIX7, "iCloud/iCloud.yaml", "ICLOUD"},        // iCloud
		{BLACKMATRIX7, "AppleTV/AppleTV.yaml", "APPLE-TV"},    // Apple TV
		{BLACKMATRIX7, "AppleMusic/AppleMusic.yaml", "APPLE-MUSIC"}, // Apple Music
		{BLACKMATRIX7, "TestFlight/TestFlight.yaml", "TESTFLIGHT"},  // TestFlight
		{BLACKMATRIX7, "Apple/Apple.yaml", "APPLE"},                 // Apple 其他服务

		// ========== 云存储（每个独立分类）==========
		{BLACKMATRIX7, "OneDrive/OneDrive.yaml", "ONEDRIVE"},       // OneDrive
		{BLACKMATRIX7, "GoogleDrive/GoogleDrive.yaml", "GDRIVE"},   // Google Drive
		{BLACKMATRIX7, "Dropbox/Dropbox.yaml", "DROPBOX"},          // Dropbox

		// ========== 科技公司 ==========
		{BLACKMATRIX7, "Google/Google.yaml", "GOOGLE"},         // Google
		{BLACKMATRIX7, "Microsoft/Microsoft.yaml", "MICROSOFT"}, // Microsoft
		{BLACKMATRIX7, "Amazon/Amazon.yaml", "AMAZON"},         // Amazon
		{BLACKMATRIX7, "Facebook/Facebook.yaml", "FACEBOOK"},   // Facebook
		{BLACKMATRIX7, "Adobe/Adobe.yaml", "ADOBE"},            // Adobe

		// ========== 开发工具/VPS（每个独立分类）==========
		{BLACKMATRIX7, "GitHub/GitHub.yaml", "GITHUB"},     // GitHub
		{BLACKMATRIX7, "GitLab/GitLab.yaml", "GITLAB"},     // GitLab
		{BLACKMATRIX7, "Docker/Docker.yaml", "DOCKER"},     // Docker Hub
		{BLACKMATRIX7, "Heroku/Heroku.yaml", "HEROKU"},     // Heroku
		{BLACKMATRIX7, "DigitalOcean/DigitalOcean.yaml", "DIGITALOCEAN"}, // DigitalOcean
		{BLACKMATRIX7, "Vercel/Vercel.yaml", "VERCEL"},     // Vercel
		{BLACKMATRIX7, "Cloudflare/Cloudflare.yaml", "CLOUDFLARE"}, // Cloudflare

		// ========== 交易所（每个独立分类）==========
		{BLACKMATRIX7, "Binance/Binance.yaml", "BINANCE"},         // 币安
		{BLACKMATRIX7, "OKX/OKX.yaml", "OKX"},                     // OKX
		{BLACKMATRIX7, "Crypto/Crypto.yaml", "CRYPTO"},            // Crypto.com
		{BLACKMATRIX7, "Cryptocurrency/Cryptocurrency.yaml", "CRYPTOCURRENCY"}, // 加密货币综合

		// ========== 支付（每个独立分类）==========
		{BLACKMATRIX7, "PayPal/PayPal.yaml", "PAYPAL"}, // PayPal

		// ========== 社交平台（每个独立分类）==========
		{BLACKMATRIX7, "Telegram/Telegram.yaml", "TELEGRAM"},   // Telegram
		{BLACKMATRIX7, "Twitter/Twitter.yaml", "TWITTER"},      // Twitter
		{BLACKMATRIX7, "Instagram/Instagram.yaml", "INSTAGRAM"}, // Instagram
		{BLACKMATRIX7, "WhatsApp/Whatsapp.yaml", "WHATSAPP"},   // WhatsApp
		{BLACKMATRIX7, "Discord/Discord.yaml", "DISCORD"},      // Discord
		{BLACKMATRIX7, "Line/Line.yaml", "LINE"},               // Line
		{BLACKMATRIX7, "Threads/Threads.yaml", "THREADS"},      // Threads
		{BLACKMATRIX7, "Reddit/Reddit.yaml", "REDDIT"},         // Reddit
		{BLACKMATRIX7, "LinkedIn/LinkedIn.yaml", "LINKEDIN"},   // LinkedIn

		// ========== 维基百科 ==========
		{BLACKMATRIX7, "Wikipedia/Wikipedia.yaml", "WIKIPEDIA"}, // Wikipedia

		// ========== 游戏 ==========
		{BLACKMATRIX7, "Steam/Steam.yaml", "STEAM"},   // Steam
		{BLACKMATRIX7, "Epic/Epic.yaml", "EPIC"},      // Epic Games
		{BLACKMATRIX7, "Sony/Sony.yaml", "PLAYSTATION"}, // PlayStation

		// ========== 购物 ==========
		{BLACKMATRIX7, "Amazon/Amazon.yaml", "AMAZON"},   // Amazon
		{BLACKMATRIX7, "eBay/eBay.yaml", "EBAY"},         // eBay
		{BLACKMATRIX7, "Shopify/Shopify.yaml", "SHOPIFY"}, // Shopify

		// ========== 咨询/新闻 ==========
		{BLACKMATRIX7, "BBC/BBC.yaml", "BBC"},             // BBC
		{BLACKMATRIX7, "CNN/CNN.yaml", "CNN"},             // CNN
		{BLACKMATRIX7, "Bloomberg/Bloomberg.yaml", "BLOOMBERG"}, // Bloomberg
		{BLACKMATRIX7, "NYTimes/NYTimes.yaml", "NYTIMES"}, // New York Times

		// ========== 政府相关 ==========
		{BLACKMATRIX7, "GlobalScholar/GlobalScholar.yaml", "SCHOLAR"}, // 学术/政府资源

		// 国内常用网站和服务
		{BLACKMATRIX7, "ChinaMax/ChinaMax.yaml", "DIRECT"}, // 国内网站合集
		{BLACKMATRIX7, "WeChat/WeChat.yaml", "DIRECT"},     // 微信
		{BLACKMATRIX7, "Weibo/Weibo.yaml", "DIRECT"},       // 微博

		// Loyalsoldier 规则
		{LOYALSOLDIER, "proxy.txt", "PROXY"},   // 代理域名列表
		{LOYALSOLDIER, "direct.txt", "DIRECT"}, // 直连域名列表
		{LOYALSOLDIER, "reject.txt", "REJECT"}, // 广告域名列表

		// 广告拦截
		{BLACKMATRIX7, "Advertising/Advertising.yaml", "REJECT"}, // 广告拦截

		// ACL4SSR 规则
		{ACL4SSR, "Clash/ProxyGFWlist.list", "PROXY"}, // GFW 列表
		{ACL4SSR, "Clash/ChinaDomain.list", "DIRECT"}, // 国内域名
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
		err := c.Parse(cfg.Source, cfg.Path, cfg.Action)
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
