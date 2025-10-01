package node

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/lazygophers/log"
)

// UnlockPrimeVideo 检测 Amazon Prime Video 解锁状态
//
// 检测原理：
//   1. 访问 Prime Video 的内容页面
//   2. 检查地区限制
//   3. 通过 API 验证地区和内容可用性
//
// 返回:
//   - *UnlockResult: 解锁检测结果
//   - error: 检测失败时返回错误
func (n *Node) UnlockPrimeVideo(ctx context.Context) (*UnlockResult, error) {
	result := &UnlockResult{
		Platform: "Prime Video",
		Status:   UnlockStatusUnknown,
	}

	// Prime Video 的地区检测
	// 使用 Prime Video 的已知内容 ID 进行测试
	headers := map[string]string{
		"Accept": "application/json",
	}

	// 访问 Prime Video 主页
	resp, err := n.doHTTPRequest(ctx, "GET", "https://www.primevideo.com/", headers, nil)
	if err != nil {
		log.Errorf("Prime Video unlock test failed: %v", err)
		result.Status = UnlockStatusFailed
		result.Message = "Failed to connect to Prime Video"
		return result, nil
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode == 403 {
		result.Status = UnlockStatusBanned
		result.Message = "Access blocked by Prime Video"
		return result, nil
	}

	if resp.StatusCode != 200 {
		result.Status = UnlockStatusNo
		result.Message = "Prime Video is not available"
		return result, nil
	}

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read Prime Video response failed: %v", err)
		result.Status = UnlockStatusFailed
		result.Message = "Failed to read response"
		return result, nil
	}

	bodyStr := string(body)

	// 检测地区限制
	if strings.Contains(bodyStr, "not available") ||
	   strings.Contains(bodyStr, "isn't available") ||
	   strings.Contains(bodyStr, "location") {
		result.Status = UnlockStatusNo
		result.Message = "Prime Video is not available in your region"
		return result, nil
	}

	// 检测是否被识别为使用代理/VPN
	if strings.Contains(bodyStr, "VPN") || strings.Contains(bodyStr, "proxy") {
		result.Status = UnlockStatusBanned
		result.Message = "Proxy/VPN detected by Prime Video"
		return result, nil
	}

	// 第二步：使用 Prime Video API 获取地区信息
	// 访问一个已知的 Prime Original 内容来确认可用性
	resp2, err := n.doHTTPRequest(ctx, "GET", "https://www.primevideo.com/region/na/detail/0LCFV0C39F8I0A8YRRDM7UJ9CZ", headers, nil)
	if err != nil {
		log.Warnf("Prime Video content check failed: %v", err)
		// 如果内容检查失败，但主页可访问，仍认为可能可用
		if resp.StatusCode == 200 {
			result.Status = UnlockStatusYes
			result.Message = "Prime Video appears to be available"
			return result, nil
		}
	} else {
		defer resp2.Body.Close()

		if resp2.StatusCode == 200 {
			result.Status = UnlockStatusYes
			result.Message = "Prime Video is available"
		} else if resp2.StatusCode == 403 {
			result.Status = UnlockStatusNo
			result.Message = "Prime Video is not available in your region"
			return result, nil
		}
	}

	// 尝试提取地区信息
	// Prime Video 可能在页面中包含地区信息
	if idx := strings.Index(bodyStr, `"currentTerritory":`); idx != -1 {
		start := idx + len(`"currentTerritory":`)
		// 跳过可能的空格和引号
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

	// 从 URL 重定向中提取地区信息
	if result.Region == "" && resp.Request != nil && resp.Request.URL != nil {
		// 检查最终 URL 中的地区代码
		finalURL := resp.Request.URL.String()
		if strings.Contains(finalURL, "/region/") {
			parts := strings.Split(finalURL, "/region/")
			if len(parts) > 1 {
				regionParts := strings.Split(parts[1], "/")
				if len(regionParts) > 0 {
					result.Region = strings.ToUpper(regionParts[0])
				}
			}
		}
	}

	// 尝试使用 API 获取地区信息
	resp3, err := n.doHTTPRequest(ctx, "GET", "https://www.primevideo.com/api/getPaymentMethods", headers, nil)
	if err == nil {
		defer resp3.Body.Close()
		apiBody, _ := io.ReadAll(resp3.Body)

		var apiData map[string]interface{}
		if json.Unmarshal(apiBody, &apiData) == nil {
			if territory, ok := apiData["territory"].(string); ok {
				result.Region = territory
			} else if marketPlace, ok := apiData["marketplace"].(string); ok {
				result.Region = marketPlace
			}
		}
	}

	// 如果状态仍然未知，根据主页请求结果判断
	if result.Status == UnlockStatusUnknown {
		if resp.StatusCode == 200 {
			result.Status = UnlockStatusYes
			result.Message = "Prime Video is available"
		} else {
			result.Status = UnlockStatusNo
			result.Message = "Prime Video is not available"
		}
	}

	// 更新消息包含地区信息
	if result.Status == UnlockStatusYes && result.Region != "" {
		result.Message = "Prime Video is available in " + result.Region
	}

	return result, nil
}
