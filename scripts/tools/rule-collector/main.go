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

		// 代理规则 - 流媒体（细分类型）
		{BLACKMATRIX7, "YouTube/YouTube.yaml", "YOUTUBE"},           // YouTube
		{BLACKMATRIX7, "YouTubeMusic/YouTubeMusic.yaml", "YOUTUBE"}, // YouTube Music
		{BLACKMATRIX7, "Netflix/Netflix.yaml", "NETFLIX"},           // Netflix
		{BLACKMATRIX7, "Disney/Disney.yaml", "DISNEY"},              // Disney+

		// 代理规则 - 国内流媒体
		{BLACKMATRIX7, "BiliBili/BiliBili.yaml", "BILIBILI"},            // 哔哩哔哩
		{BLACKMATRIX7, "BiliBiliIntl/BiliBiliIntl.yaml", "BILIBILI-HK"}, // 哔哩哔哩国际版
		{BLACKMATRIX7, "iQIYI/iQIYI.yaml", "IQIYI"},                     // 爱奇艺
		{BLACKMATRIX7, "iQIYIIntl/iQIYIIntl.yaml", "IQIYI-HK"},          // 爱奇艺国际版

		// 代理规则 - AI 服务（细分类型）
		{BLACKMATRIX7, "OpenAI/OpenAI.yaml", "OPENAI"},   // OpenAI
		{BLACKMATRIX7, "Claude/Claude.yaml", "CLAUDE"},   // Claude
		{BLACKMATRIX7, "Copilot/Copilot.yaml", "GITHUB"}, // GitHub Copilot
		{BLACKMATRIX7, "Gemini/Gemini.yaml", "GEMINI"},   // Google Gemini

		// 代理规则 - 开发工具
		{BLACKMATRIX7, "GitHub/GitHub.yaml", "GITHUB"}, // GitHub
		{BLACKMATRIX7, "Docker/Docker.yaml", "DOCKER"}, // Docker Hub

		// 代理规则 - 社交平台
		{BLACKMATRIX7, "Telegram/Telegram.yaml", "TELEGRAM"}, // Telegram
		{BLACKMATRIX7, "Twitter/Twitter.yaml", "PROXY"},      // Twitter

		// 代理规则 - 科技公司
		{BLACKMATRIX7, "Google/Google.yaml", "GOOGLE"},         // Google
		{BLACKMATRIX7, "Microsoft/Microsoft.yaml", "MICROSOFT"}, // Microsoft

		// 国内常用网站和服务
		{BLACKMATRIX7, "ChinaMax/ChinaMax.yaml", "DIRECT"}, // 国内网站合集

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
