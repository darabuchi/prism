package node

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/lazygophers/log"
)

// UnlockNetflix 检测 Netflix 解锁状态
//
// 检测原理：
//   1. 访问 Netflix 的原创内容页面（如《纸牌屋》）
//   2. 根据返回的地区代码和内容可用性判断解锁状态
//
// 返回:
//   - *UnlockResult: 解锁检测结果
//   - error: 检测失败时返回错误
func (n *Node) UnlockNetflix(ctx context.Context) (*UnlockResult, error) {
	result := &UnlockResult{
		Platform: "Netflix",
		Status:   UnlockStatusUnknown,
	}

	// 第一步：检测基本的地区信息
	// 使用 Netflix 的 title API
	resp, err := n.doHTTPRequest(ctx, "GET", "https://www.netflix.com/title/70143836", nil, nil)
	if err != nil {
		log.Errorf("Netflix unlock test failed (step 1): %v", err)
		result.Status = UnlockStatusFailed
		result.Message = "Failed to connect to Netflix"
		return result, nil
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read Netflix response failed: %v", err)
		result.Status = UnlockStatusFailed
		result.Message = "Failed to read response"
		return result, nil
	}

	bodyStr := string(body)

	// 检测是否被封禁
	if strings.Contains(bodyStr, "Not Available") || strings.Contains(bodyStr, "Unavailable") {
		result.Status = UnlockStatusNo
		result.Message = "Content not available in your region"
		return result, nil
	}

	// 检测是否使用代理被识别
	if strings.Contains(bodyStr, "proxy") || strings.Contains(bodyStr, "unblocker") || strings.Contains(bodyStr, "VPN") {
		result.Status = UnlockStatusBanned
		result.Message = "Proxy/VPN detected by Netflix"
		return result, nil
	}

	// 第二步：使用 Netflix 的 API 获取详细地区信息
	resp2, err := n.doHTTPRequest(ctx, "GET", "https://www.netflix.com/api/shakti/v23a9c65e/nq/website/nq?q=%7B%22paths%22%3A%5B%5B%22lolomo%22%2C%7B%22from%22%3A0%2C%22to%22%3A20%7D%2C%7B%22from%22%3A0%2C%22to%22%3A20%7D%2C%22summary%22%5D%5D%7D", nil, nil)
	if err != nil {
		// API 请求失败，但前面的检测成功，仍然认为可能解锁
		log.Warnf("Netflix API request failed: %v", err)
	} else {
		defer resp2.Body.Close()
		apiBody, _ := io.ReadAll(resp2.Body)

		// 尝试提取地区信息
		var apiData map[string]interface{}
		if json.Unmarshal(apiBody, &apiData) == nil {
			// 从 API 响应中提取地区代码
			if value, ok := apiData["value"].(map[string]interface{}); ok {
				if paths, ok := value["paths"].([]interface{}); ok && len(paths) > 0 {
					// 成功获取内容，说明已解锁
					result.Status = UnlockStatusYes
				}
			}
		}
	}

	// 尝试从响应头获取地区信息
	if resp.Header.Get("X-Netflix-Country") != "" {
		result.Region = resp.Header.Get("X-Netflix-Country")
	}

	// 如果能访问到内容页面且没有错误提示，认为已解锁
	if result.Status == UnlockStatusUnknown {
		if resp.StatusCode == 200 && !strings.Contains(bodyStr, "Not Available") {
			result.Status = UnlockStatusYes
			result.Message = "Netflix is available"
		} else {
			result.Status = UnlockStatusNo
			result.Message = "Netflix is not available"
		}
	}

	// 尝试从 HTML 中提取地区代码
	if result.Region == "" {
		// 查找类似 "currentCountry":"US" 的模式
		if idx := strings.Index(bodyStr, `"currentCountry":"`); idx != -1 {
			start := idx + len(`"currentCountry":"`)
			if start < len(bodyStr) {
				end := strings.Index(bodyStr[start:], `"`)
				if end != -1 && start+end <= len(bodyStr) {
					result.Region = bodyStr[start : start+end]
				}
			}
		}
	}

	if result.Status == UnlockStatusYes {
		if result.Region != "" {
			result.Message = "Netflix is available in " + result.Region
		} else {
			result.Message = "Netflix is available"
		}
	}

	return result, nil
}
