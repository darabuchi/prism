package node

import (
	"context"
	"io"
	"strings"

	"github.com/lazygophers/log"
)

// UnlockTikTok 检测 TikTok 解锁状态
//
// 检测原理：
//   1. 访问 TikTok 的主页
//   2. 检查地区限制
//   3. 通过响应判断可用地区
//
// 注意：TikTok 在某些地区（如印度）不可用
//
// 返回:
//   - *UnlockResult: 解锁检测结果
//   - error: 检测失败时返回错误
func (n *Node) UnlockTikTok(ctx context.Context) (*UnlockResult, error) {
	result := &UnlockResult{
		Platform: "TikTok",
		Status:   UnlockStatusUnknown,
	}

	// TikTok 的地区检测
	headers := map[string]string{
		"Accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"User-Agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
	}

	// 访问 TikTok 主页
	resp, err := n.doHTTPRequest(ctx, "GET", "https://www.tiktok.com/", headers, nil)
	if err != nil {
		log.Errorf("TikTok unlock test failed: %v", err)
		result.Status = UnlockStatusFailed
		result.Message = "Failed to connect to TikTok"
		return result, nil
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode == 403 {
		result.Status = UnlockStatusBanned
		result.Message = "Access blocked by TikTok"
		return result, nil
	}

	if resp.StatusCode == 451 {
		// HTTP 451 表示因法律原因不可用
		result.Status = UnlockStatusNo
		result.Message = "TikTok is not available in your region (legal restrictions)"
		return result, nil
	}

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read TikTok response failed: %v", err)
		result.Status = UnlockStatusFailed
		result.Message = "Failed to read response"
		return result, nil
	}

	bodyStr := string(body)

	// 检测地区限制
	if strings.Contains(bodyStr, "not available") ||
	   strings.Contains(bodyStr, "region") ||
	   strings.Contains(bodyStr, "blocked") ||
	   strings.Contains(bodyStr, "banned") {
		result.Status = UnlockStatusNo
		result.Message = "TikTok is not available in your region"
		return result, nil
	}

	// 检测是否被识别为使用代理
	if strings.Contains(bodyStr, "proxy") || strings.Contains(bodyStr, "vpn") {
		result.Status = UnlockStatusBanned
		result.Message = "Proxy/VPN detected by TikTok"
		return result, nil
	}

	// 尝试提取地区信息
	// TikTok 可能在 HTML 中包含地区代码
	if idx := strings.Index(bodyStr, `"region":`); idx != -1 {
		start := idx + len(`"region":`)
		// 跳过空格和引号
		for start < len(bodyStr) && (bodyStr[start] == ' ' || bodyStr[start] == '"') {
			start++
		}
		end := start
		for end < len(bodyStr) && bodyStr[end] != '"' && bodyStr[end] != ',' && bodyStr[end] != '}' {
			end++
		}
		if end > start && end-start <= 3 {
			result.Region = strings.ToUpper(bodyStr[start:end])
		}
	}

	// 尝试从 locale 信息中提取地区
	if result.Region == "" {
		if idx := strings.Index(bodyStr, `"locale":`); idx != -1 {
			start := idx + len(`"locale":`)
			for start < len(bodyStr) && (bodyStr[start] == ' ' || bodyStr[start] == '"') {
				start++
			}
			end := start
			for end < len(bodyStr) && bodyStr[end] != '"' && bodyStr[end] != ',' {
				end++
			}
			if end > start {
				locale := bodyStr[start:end]
				// locale 格式通常是 en-US, zh-CN 等
				parts := strings.Split(locale, "-")
				if len(parts) == 2 {
					result.Region = strings.ToUpper(parts[1])
				}
			}
		}
	}

	// 第二步：尝试访问 API 端点
	resp2, err := n.doHTTPRequest(ctx, "GET", "https://www.tiktok.com/api/challenge/list/", headers, nil)
	if err != nil {
		log.Warnf("TikTok API check failed: %v", err)
		// 如果主页可访问，认为服务可用
		if resp.StatusCode == 200 {
			result.Status = UnlockStatusYes
			result.Message = "TikTok is available"
			return result, nil
		}
	} else {
		defer resp2.Body.Close()

		if resp2.StatusCode == 200 {
			result.Status = UnlockStatusYes
		} else if resp2.StatusCode == 403 || resp2.StatusCode == 451 {
			result.Status = UnlockStatusNo
			result.Message = "TikTok is not available in your region"
			return result, nil
		}
	}

	// 如果状态仍然未知，根据主页请求结果判断
	if result.Status == UnlockStatusUnknown {
		if resp.StatusCode == 200 {
			result.Status = UnlockStatusYes
			result.Message = "TikTok is available"
		} else {
			result.Status = UnlockStatusNo
			result.Message = "TikTok is not available"
		}
	}

	// 更新消息包含地区信息
	if result.Status == UnlockStatusYes && result.Region != "" {
		result.Message = "TikTok is available in " + result.Region
	}

	return result, nil
}
