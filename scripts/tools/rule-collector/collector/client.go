package collector

import (
	"net/http"
	"net/url"
	"time"

	"github.com/pterm/pterm"
)

// NewHTTPClient 创建 HTTP 客户端（支持代理）
func NewHTTPClient(proxy string) *http.Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 配置代理
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err == nil {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			}
			pterm.Info.Printfln("Using proxy: %s", proxy)
		} else {
			pterm.Warning.Printfln("Invalid proxy URL: %s - %v", proxy, err)
		}
	}

	return client
}
