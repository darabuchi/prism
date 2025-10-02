package collector

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/darabuchi/prism"
	"github.com/darabuchi/prism/pkg/rules"
)

// LoyalSoldier Loyalsoldier 规则源处理器
// 来源：https://github.com/Loyalsoldier/clash-rules
type LoyalSoldier struct {
	*BaseHandler
	baseURL string
}

// NewLoyalSoldier 创建 LoyalSoldier 处理器
func NewLoyalSoldier(proxy string, cacheDays int) *LoyalSoldier {
	return &LoyalSoldier{
		BaseHandler: NewBaseHandler(proxy, cacheDays),
		baseURL:     "https://raw.githubusercontent.com/Loyalsoldier/clash-rules/release/",
	}
}

// Download 下载规则数据
func (h *LoyalSoldier) Download(path string) ([]byte, error) {
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
// Loyalsoldier 规则格式：
// - YAML 列表格式：使用 "- '域名'" 格式
// - 根据文件名判断规则类型
func (h *LoyalSoldier) Parse(body []byte, payload prism.Payload) ([]rules.Rule, error) {
	// Loyalsoldier 使用 YAML 列表格式，需要使用 Clash 解析器
	return ParseClashRules(body, payload)
}

// NeedUpdate 判断缓存是否需要更新
// 使用基类实现

// ParseWithType 根据文件类型解析规则
func (h *LoyalSoldier) ParseWithType(body []byte, payload prism.Payload, filename string) ([]rules.Rule, error) {
	var ruleTypeStr string

	// 根据文件名判断规则类型
	switch {
	case strings.Contains(filename, "cncidr") || strings.Contains(filename, "lancidr"):
		ruleTypeStr = "IP-CIDR"
	case strings.Contains(filename, "applications"):
		ruleTypeStr = "DOMAIN"
	default:
		ruleTypeStr = "DOMAIN-SUFFIX"
	}

	return ParseTextRules(body, payload, ruleTypeStr)
}
