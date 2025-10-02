package collector

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/darabuchi/prism/pkg/rules"
)

// ACL4SSR ACL4SSR 规则源处理器
// 来源：https://github.com/ACL4SSR/ACL4SSR
type ACL4SSR struct {
	baseURL string
	client  *http.Client
}

// NewACL4SSR 创建 ACL4SSR 处理器
func NewACL4SSR() *ACL4SSR {
	return &ACL4SSR{
		baseURL: "https://raw.githubusercontent.com/ACL4SSR/ACL4SSR/master/",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Download 下载规则数据
func (h *ACL4SSR) Download(path string) ([]byte, error) {
	url := h.baseURL + path

	resp, err := h.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body failed: %w", err)
	}

	return body, nil
}

// Parse 解析规则数据
// ACL4SSR 规则使用 Clash 格式
func (h *ACL4SSR) Parse(body []byte, action string) ([]rules.Rule, error) {
	return ParseClashRules(body, action)
}

// NeedUpdate 判断缓存是否需要更新
// ACL4SSR 规则每天更新一次
func (h *ACL4SSR) NeedUpdate(info os.FileInfo) bool {
	return time.Since(info.ModTime()) > 24*time.Hour
}
