package node

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/lazygophers/log"
)

// UnlockHBOMax 检测 HBO Max 解锁状态
//
// 检测原理：
//   1. 访问 HBO Max 的 API 端点
//   2. 检查地区可用性
//   3. 验证内容访问权限
//
// 返回:
//   - *UnlockResult: 解锁检测结果
//   - error: 检测失败时返回错误
func (n *Node) UnlockHBOMax(ctx context.Context) (*UnlockResult, error) {
	result := &UnlockResult{
		Platform: "HBO Max",
		Status:   UnlockStatusUnknown,
	}

	// HBO Max 的地区检测
	// 使用 HBO Max 的 API 来检测地区
	headers := map[string]string{
		"Accept": "application/json",
	}

	// 访问 HBO Max 的地区检测 API
	resp, err := n.doHTTPRequest(ctx, "GET", "https://www.hbomax.com/", headers, nil)
	if err != nil {
		log.Errorf("HBO Max unlock test failed: %v", err)
		result.Status = UnlockStatusFailed
		result.Message = "Failed to connect to HBO Max"
		return result, nil
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode == 403 {
		result.Status = UnlockStatusBanned
		result.Message = "Access blocked by HBO Max"
		return result, nil
	}

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read HBO Max response failed: %v", err)
		result.Status = UnlockStatusFailed
		result.Message = "Failed to read response"
		return result, nil
	}

	bodyStr := string(body)

	// 检测地区限制
	if strings.Contains(bodyStr, "not available") ||
	   strings.Contains(bodyStr, "isn't available") ||
	   strings.Contains(bodyStr, "outside") {
		result.Status = UnlockStatusNo
		result.Message = "HBO Max is not available in your region"
		return result, nil
	}

	// 检测是否被识别为使用代理
	if strings.Contains(bodyStr, "proxy") || strings.Contains(bodyStr, "vpn") {
		result.Status = UnlockStatusBanned
		result.Message = "Proxy/VPN detected by HBO Max"
		return result, nil
	}

	// 第二步：使用 HBO Max API 获取详细地区信息
	resp2, err := n.doHTTPRequest(ctx, "GET", "https://comet.api.hbo.com/express-content/regions", headers, nil)
	if err != nil {
		log.Warnf("HBO Max API request failed: %v", err)
		// 如果 API 失败，但主页可访问，仍认为可能可用
		if resp.StatusCode == 200 {
			result.Status = UnlockStatusYes
			result.Message = "HBO Max appears to be available"
			return result, nil
		}
	} else {
		defer resp2.Body.Close()

		if resp2.StatusCode == 200 {
			apiBody, _ := io.ReadAll(resp2.Body)

			// 尝试提取地区信息
			var regionData map[string]interface{}
			if json.Unmarshal(apiBody, &regionData) == nil {
				// 尝试从响应中提取地区代码
				if region, ok := regionData["code"].(string); ok {
					result.Region = region
				} else if regions, ok := regionData["regions"].([]interface{}); ok && len(regions) > 0 {
					if firstRegion, ok := regions[0].(map[string]interface{}); ok {
						if code, ok := firstRegion["code"].(string); ok {
							result.Region = code
						}
					}
				}
			}

			result.Status = UnlockStatusYes
			if result.Region != "" {
				result.Message = "HBO Max is available in " + result.Region
			} else {
				result.Message = "HBO Max is available"
			}
		} else if resp2.StatusCode == 403 {
			result.Status = UnlockStatusNo
			result.Message = "HBO Max is not available in your region"
		}
	}

	// 如果状态仍然未知，根据主页请求结果判断
	if result.Status == UnlockStatusUnknown {
		if resp.StatusCode == 200 {
			result.Status = UnlockStatusYes
			result.Message = "HBO Max is available"
		} else {
			result.Status = UnlockStatusNo
			result.Message = "HBO Max is not available"
		}
	}

	return result, nil
}
