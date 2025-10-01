package node

import (
	"context"
	"encoding/json"
	"io"

	"github.com/lazygophers/log"
)

// UnlockAbemaTV 检测 AbemaTV 解锁状态
//
// 检测原理：
//   1. 访问 AbemaTV 的 API
//   2. 检查地区限制
//   3. 验证日本地区内容的可用性
//
// 注意：AbemaTV 主要服务日本地区
//
// 返回:
//   - *UnlockResult: 解锁检测结果
//   - error: 检测失败时返回错误
func (n *Node) UnlockAbemaTV(ctx context.Context) (*UnlockResult, error) {
	result := &UnlockResult{
		Platform: "AbemaTV",
		Status:   UnlockStatusUnknown,
	}

	// AbemaTV 的地区检测
	// 使用 AbemaTV 的 license API
	headers := map[string]string{
		"Accept":     "application/json",
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
	}

	// 访问 AbemaTV 的许可证 API 来检测地区
	resp, err := n.doHTTPRequest(ctx, "GET", "https://api.abema.io/v1/ip/check?device=pc", headers, nil)
	if err != nil {
		log.Errorf("AbemaTV unlock test failed: %v", err)
		result.Status = UnlockStatusFailed
		result.Message = "Failed to connect to AbemaTV"
		return result, nil
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode == 403 {
		result.Status = UnlockStatusBanned
		result.Message = "Access blocked by AbemaTV"
		return result, nil
	}

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read AbemaTV response failed: %v", err)
		result.Status = UnlockStatusFailed
		result.Message = "Failed to read response"
		return result, nil
	}

	// 解析 JSON 响应
	var ipCheckData map[string]interface{}
	if err := json.Unmarshal(body, &ipCheckData); err != nil {
		log.Errorf("parse AbemaTV response failed: %v", err)
		result.Status = UnlockStatusFailed
		result.Message = "Failed to parse response"
		return result, nil
	}

	// 检查 isoCountryCode 字段
	// AbemaTV 只在日本可用
	if countryCode, ok := ipCheckData["isoCountryCode"].(string); ok {
		if countryCode == "JP" {
			result.Status = UnlockStatusYes
			result.Region = "JP"
			result.Message = "AbemaTV is available in JP"
		} else {
			result.Status = UnlockStatusNo
			result.Region = countryCode
			result.Message = "AbemaTV is only available in Japan"
		}
		return result, nil
	}

	// 检查 cdnRegion 字段
	if cdnRegion, ok := ipCheckData["cdnRegion"].(string); ok {
		if cdnRegion == "jp" || cdnRegion == "JP" {
			result.Status = UnlockStatusYes
			result.Region = "JP"
			result.Message = "AbemaTV is available in JP"
		} else {
			result.Status = UnlockStatusNo
			result.Message = "AbemaTV is only available in Japan"
		}
		return result, nil
	}

	// 第二步：尝试访问实际内容
	resp2, err := n.doHTTPRequest(ctx, "GET", "https://abema.tv/", headers, nil)
	if err != nil {
		log.Warnf("AbemaTV homepage check failed: %v", err)
		// 如果 IP 检查通过，但主页访问失败，仍返回 IP 检查结果
		if resp.StatusCode == 200 {
			result.Status = UnlockStatusYes
			result.Message = "AbemaTV appears to be available"
			return result, nil
		}
	} else {
		defer resp2.Body.Close()

		if resp2.StatusCode == 200 {
			result.Status = UnlockStatusYes
			result.Region = "JP"
			result.Message = "AbemaTV is available"
		} else if resp2.StatusCode == 403 {
			result.Status = UnlockStatusNo
			result.Message = "AbemaTV is not available in your region"
		}
	}

	// 如果状态仍然未知
	if result.Status == UnlockStatusUnknown {
		if resp.StatusCode == 200 {
			result.Status = UnlockStatusYes
			result.Message = "AbemaTV appears to be available"
		} else {
			result.Status = UnlockStatusNo
			result.Message = "AbemaTV is not available"
		}
	}

	return result, nil
}
