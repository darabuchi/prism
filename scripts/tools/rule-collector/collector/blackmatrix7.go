package collector

import (
	"fmt"
	"io"
	"net/http"

	"github.com/darabuchi/prism"
	"github.com/darabuchi/prism/pkg/rules"
)

// BlackMatrix7 BlackMatrix7 规则源处理器
// 来源：https://github.com/blackmatrix7/ios_rule_script
type BlackMatrix7 struct {
	*BaseHandler
	baseURL string
}

// NewBlackMatrix7 创建 BlackMatrix7 处理器
func NewBlackMatrix7(proxy string, cacheDays int) *BlackMatrix7 {
	return &BlackMatrix7{
		BaseHandler: NewBaseHandler(proxy, cacheDays),
		baseURL:     "https://raw.githubusercontent.com/blackmatrix7/ios_rule_script/master/rule/Clash/",
	}
}

// Download 下载规则数据
func (h *BlackMatrix7) Download(path string) ([]byte, error) {
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
func (h *BlackMatrix7) Parse(body []byte, payload prism.Payload) ([]rules.Rule, error) {
	return ParseClashRules(body, payload)
}

// NeedUpdate 判断缓存是否需要更新
// 使用基类实现
