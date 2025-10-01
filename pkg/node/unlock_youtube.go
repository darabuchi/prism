package node

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/lazygophers/log"
)

// UnlockYouTubePremium 检测 YouTube Premium 解锁状态
//
// 检测原理：
//   1. 访问 YouTube 的内容页面检测区域
//   2. 检查 Premium 内容的可用性
//   3. 通过 API 验证区域信息
//
// 返回:
//   - *UnlockResult: 解锁检测结果
//   - error: 检测失败时返回错误
func (n *Node) UnlockYouTubePremium(ctx context.Context) (*UnlockResult, error) {
	result := &UnlockResult{
		Platform: "YouTube Premium",
		Status:   UnlockStatusUnknown,
	}

	// 第一步：通过 YouTube 的地区检测接口获取当前位置
	// 使用 YouTube 的 API 获取位置信息
	headers := map[string]string{
		"Accept": "application/json",
	}

	// 访问 YouTube 的 player API 来检测地区
	resp, err := n.doHTTPRequest(ctx, "GET", "https://www.youtube.com/premium", headers, nil)
	if err != nil {
		log.Errorf("YouTube Premium unlock test failed (step 1): %v", err)
		result.Status = UnlockStatusFailed
		result.Message = "Failed to connect to YouTube"
		return result, nil
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode == 403 {
		result.Status = UnlockStatusBanned
		result.Message = "Access blocked by YouTube"
		return result, nil
	}

	if resp.StatusCode != 200 {
		result.Status = UnlockStatusNo
		result.Message = "YouTube Premium is not available"
		return result, nil
	}

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read YouTube Premium response failed: %v", err)
		result.Status = UnlockStatusFailed
		result.Message = "Failed to read response"
		return result, nil
	}

	bodyStr := string(body)

	// 检测是否被识别为使用代理
	if strings.Contains(bodyStr, "vpn") || strings.Contains(bodyStr, "proxy") {
		result.Status = UnlockStatusBanned
		result.Message = "Proxy/VPN detected by YouTube"
		return result, nil
	}

	// 第二步：使用 YouTube 的地区 API 获取详细信息
	resp2, err := n.doHTTPRequest(ctx, "GET", "https://www.youtube.com/red", headers, nil)
	if err != nil {
		// 如果 red API 失败，但 premium 页面可访问，仍认为可能可用
		log.Warnf("YouTube red API request failed: %v", err)
		if resp.StatusCode == 200 {
			result.Status = UnlockStatusYes
			result.Message = "YouTube Premium appears to be available"
			return result, nil
		}
	} else {
		defer resp2.Body.Close()

		if resp2.StatusCode == 200 {
			result.Status = UnlockStatusYes
			result.Message = "YouTube Premium is available"
		} else if resp2.StatusCode == 403 {
			result.Status = UnlockStatusNo
			result.Message = "YouTube Premium is not available in your region"
		}
	}

	// 第三步：尝试获取地区信息
	// 使用 YouTube 的 API 端点来获取地区代码
	resp3, err := n.doHTTPRequest(ctx, "GET", "https://www.youtube.com/getcountry", headers, nil)
	if err != nil {
		log.Warnf("YouTube country API request failed: %v", err)
	} else {
		defer resp3.Body.Close()
		countryBody, _ := io.ReadAll(resp3.Body)

		// 尝试解析地区信息
		// YouTube 的 getcountry API 返回格式：{"country_code":"US"}
		var countryData map[string]interface{}
		if json.Unmarshal(countryBody, &countryData) == nil {
			if code, ok := countryData["country_code"].(string); ok {
				result.Region = code
			}
		}

		// 如果 JSON 解析失败，尝试从纯文本中提取
		if result.Region == "" {
			countryStr := string(countryBody)
			// 有些情况下直接返回国家代码
			countryStr = strings.TrimSpace(countryStr)
			if len(countryStr) == 2 {
				result.Region = strings.ToUpper(countryStr)
			}
		}
	}

	// 如果状态仍然未知，根据 premium 页面请求结果判断
	if result.Status == UnlockStatusUnknown {
		if resp.StatusCode == 200 {
			result.Status = UnlockStatusYes
			result.Message = "YouTube Premium is available"
		} else {
			result.Status = UnlockStatusNo
			result.Message = "YouTube Premium is not available"
		}
	}

	// 设置详细消息
	if result.Status == UnlockStatusYes {
		if result.Region != "" {
			result.Message = "YouTube Premium is available in " + result.Region
		} else {
			result.Message = "YouTube Premium is available"
		}
	}

	return result, nil
}
