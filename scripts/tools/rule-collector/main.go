package main

import (
	"fmt"
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
		pterm.Success.Printfln("Total time: %v", time.Since(start))
	}()

	pterm.Info.Printfln("PID: %d", os.Getpid())
	pterm.Info.Printfln("Rule Collector for Prism")

	c := NewCollector()

	// 注册数据源收集器
	c.AddHandle(BLACKMATRIX7, collector.NewBlackMatrix7())
	c.AddHandle(LOYALSOLDIER, collector.NewLoyalSoldier())
	c.AddHandle(ACL4SSR, collector.NewACL4SSR())

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

		// 代理规则 - 流媒体
		{BLACKMATRIX7, "YouTube/YouTube.yaml", "PROXY"},           // YouTube
		{BLACKMATRIX7, "YouTubeMusic/YouTubeMusic.yaml", "PROXY"}, // YouTube Music
		{BLACKMATRIX7, "Netflix/Netflix.yaml", "PROXY"},           // Netflix
		{BLACKMATRIX7, "Disney/Disney.yaml", "PROXY"},             // Disney+

		// 代理规则 - 国内流媒体
		{BLACKMATRIX7, "BiliBili/BiliBili.yaml", "DIRECT"},             // 哔哩哔哩
		{BLACKMATRIX7, "BiliBiliIntl/BiliBiliIntl.yaml", "PROXY"},      // 哔哩哔哩国际版
		{BLACKMATRIX7, "iQIYI/iQIYI.yaml", "DIRECT"},                   // 爱奇艺
		{BLACKMATRIX7, "iQIYIIntl/iQIYIIntl.yaml", "PROXY"},            // 爱奇艺国际版
		{BLACKMATRIX7, "TencentVideo/TencentVideo.yaml", "DIRECT"},     // 腾讯视频
		{BLACKMATRIX7, "YoukuYouku/YoukuYouku.yaml", "DIRECT"},         // 优酷
		{BLACKMATRIX7, "NetEaseMusic/NetEaseMusic.yaml", "DIRECT"},     // 网易云音乐

		// 代理规则 - AI 服务
		{BLACKMATRIX7, "OpenAI/OpenAI.yaml", "PROXY"},   // OpenAI
		{BLACKMATRIX7, "Claude/Claude.yaml", "PROXY"},   // Claude
		{BLACKMATRIX7, "Copilot/Copilot.yaml", "PROXY"}, // GitHub Copilot
		{BLACKMATRIX7, "Gemini/Gemini.yaml", "PROXY"},   // Google Gemini
		{BLACKMATRIX7, "BardAI/BardAI.yaml", "PROXY"},   // Bard AI

		// 代理规则 - 开发工具
		{BLACKMATRIX7, "GitHub/GitHub.yaml", "PROXY"},     // GitHub
		{BLACKMATRIX7, "GitLab/GitLab.yaml", "PROXY"},     // GitLab
		{BLACKMATRIX7, "Docker/Docker.yaml", "PROXY"},     // Docker Hub
		{BLACKMATRIX7, "NPM/NPM.yaml", "PROXY"},           // NPM
		{BLACKMATRIX7, "PyPI/PyPI.yaml", "PROXY"},         // PyPI

		// 代理规则 - 社交平台
		{BLACKMATRIX7, "Telegram/Telegram.yaml", "PROXY"}, // Telegram
		{BLACKMATRIX7, "Twitter/Twitter.yaml", "PROXY"},   // Twitter
		{BLACKMATRIX7, "Facebook/Facebook.yaml", "PROXY"}, // Facebook
		{BLACKMATRIX7, "Instagram/Instagram.yaml", "PROXY"}, // Instagram
		{BLACKMATRIX7, "Reddit/Reddit.yaml", "PROXY"},     // Reddit

		// 代理规则 - 科技公司
		{BLACKMATRIX7, "Google/Google.yaml", "PROXY"},       // Google
		{BLACKMATRIX7, "Apple/Apple.yaml", "PROXY"},         // Apple
		{BLACKMATRIX7, "Microsoft/Microsoft.yaml", "PROXY"}, // Microsoft
		{BLACKMATRIX7, "Amazon/Amazon.yaml", "PROXY"},       // Amazon

		// 代理规则 - 其他
		{BLACKMATRIX7, "Spotify/Spotify.yaml", "PROXY"},   // Spotify
		{BLACKMATRIX7, "PayPal/PayPal.yaml", "PROXY"},     // PayPal
		{BLACKMATRIX7, "Steam/Steam.yaml", "PROXY"},       // Steam
		{BLACKMATRIX7, "Wikipedia/Wikipedia.yaml", "PROXY"}, // Wikipedia

		// 国内常用网站和服务
		{BLACKMATRIX7, "ChinaMax/ChinaMax.yaml", "DIRECT"},   // 国内网站合集
		{BLACKMATRIX7, "ByteDance/ByteDance.yaml", "DIRECT"}, // 字节跳动
		{BLACKMATRIX7, "Alibaba/Alibaba.yaml", "DIRECT"},     // 阿里巴巴
		{BLACKMATRIX7, "Tencent/Tencent.yaml", "DIRECT"},     // 腾讯
		{BLACKMATRIX7, "Baidu/Baidu.yaml", "DIRECT"},         // 百度

		// Loyalsoldier 规则
		{LOYALSOLDIER, "proxy.txt", "PROXY"},         // 代理域名列表
		{LOYALSOLDIER, "direct.txt", "DIRECT"},       // 直连域名列表
		{LOYALSOLDIER, "reject.txt", "REJECT"},       // 广告域名列表
		{LOYALSOLDIER, "gfw.txt", "PROXY"},           // GFW 域名列表
		{LOYALSOLDIER, "greatfire.txt", "PROXY"},     // GreatFire 域名列表
		{LOYALSOLDIER, "tld-not-cn.txt", "PROXY"},    // 非中国顶级域名
		{LOYALSOLDIER, "apple-cn.txt", "DIRECT"},     // Apple 中国服务
		{LOYALSOLDIER, "google-cn.txt", "DIRECT"},    // Google 中国服务

		// 广告拦截
		{BLACKMATRIX7, "Advertising/Advertising.yaml", "REJECT"}, // 广告拦截
		{BLACKMATRIX7, "Privacy/Privacy.yaml", "REJECT"},         // 隐私保护
		{BLACKMATRIX7, "Hijacking/Hijacking.yaml", "REJECT"},     // 反劫持

		// ACL4SSR 规则
		{ACL4SSR, "Clash/ProxyGFWlist.list", "PROXY"},  // GFW 列表
		{ACL4SSR, "Clash/ChinaDomain.list", "DIRECT"},  // 国内域名
		{ACL4SSR, "Clash/ChinaCompanyIp.list", "DIRECT"}, // 国内公司 IP
	}

	// 执行规则解析
	for _, cfg := range parseList {
		err := c.Parse(cfg.Source, cfg.Path, cfg.Action)
		if err != nil {
			pterm.Error.Printfln("Failed to parse %s/%s: %v", cfg.Source, cfg.Path, err)
			// 继续处理其他规则，不中断
			continue
		}
	}

	// 导出规则
	err := c.Export()
	if err != nil {
		pterm.Error.Printfln("Failed to export rules: %v", err)
		os.Exit(1)
	}

	pterm.Success.Printfln("Rule collection completed successfully")
	pterm.Info.Printfln("Total rules collected: %d", c.RuleCount())
}
