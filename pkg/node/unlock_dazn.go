package node

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/lazygophers/log"
)

// UnlockDAZN 检测 DAZN 解锁状态
//
// 检测原理：
//   1. 访问 DAZN 的地区检测 API
//   2. 检查服务可用地区
//   3. 验证内容访问权限
//
// 注意：DAZN 在不同地区提供不同的内容
//
// 返回:
//   - *UnlockResult: 解锁检测结果
//   - error: 检测失败时返回错误
func (n *Node) UnlockDAZN(ctx context.Context) (*UnlockResult, error) {
	result := &UnlockResult{
		Platform: "DAZN",
		Status:   UnlockStatusUnknown,
	}

	// DAZN 的地区检测
	headers := map[string]string{
		"Accept":     "application/json",
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
	}

	// 访问 DAZN 的地区检测接口
	resp, err := n.doHTTPRequest(ctx, "GET", "https://startup.api.indazn.com/misl/v5/Startup", headers, nil)
	if err != nil {
		log.Errorf("DAZN unlock test failed: %v", err)
		result.Status = UnlockStatusFailed
		result.Message = "Failed to connect to DAZN"
		return result, nil
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode == 403 {
		result.Status = UnlockStatusBanned
		result.Message = "Access blocked by DAZN"
		return result, nil
	}

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read DAZN response failed: %v", err)
		result.Status = UnlockStatusFailed
		result.Message = "Failed to read response"
		return result, nil
	}

	bodyStr := string(body)

	// 检查是否包含地区限制信息
	if strings.Contains(bodyStr, "NOT_AVAILABLE") ||
	   strings.Contains(bodyStr, "GEO_BLOCKED") ||
	   strings.Contains(bodyStr, "not available") {
		result.Status = UnlockStatusNo
		result.Message = "DAZN is not available in your region"
		return result, nil
	}

	// 解析 JSON 响应获取地区信息
	var startupData map[string]interface{}
	if err := json.Unmarshal(body, &startupData); err == nil {
		// 尝试提取地区信息
		if region, ok := startupData["Region"].(string); ok {
			result.Region = region
		} else if region, ok := startupData["region"].(string); ok {
			result.Region = region
		}

		// 检查 GeolocatedCountry
		if country, ok := startupData["GeolocatedCountry"].(string); ok {
			result.Region = country
		}

		// 检查是否允许访问
		if isAllowed, ok := startupData["isAllowed"].(bool); ok {
			if isAllowed {
				result.Status = UnlockStatusYes
			} else {
				result.Status = UnlockStatusNo
				result.Message = "DAZN is not available in your region"
				return result, nil
			}
		}
	}

	// 第二步：检查 DAZN 主页
	resp2, err := n.doHTTPRequest(ctx, "GET", "https://www.dazn.com/", headers, nil)
	if err != nil {
		log.Warnf("DAZN homepage check failed: %v", err)
		// 如果 API 请求成功，但主页失败，仍基于 API 结果判断
		if resp.StatusCode == 200 {
			if result.Status == UnlockStatusUnknown {
				result.Status = UnlockStatusYes
				result.Message = "DAZN appears to be available"
			}
			return result, nil
		}
	} else {
		defer resp2.Body.Close()

		if resp2.StatusCode == 200 {
			if result.Status == UnlockStatusUnknown {
				result.Status = UnlockStatusYes
			}
		} else if resp2.StatusCode == 403 {
			result.Status = UnlockStatusNo
			result.Message = "DAZN is not available in your region"
			return result, nil
		}
	}

	// 如果状态仍然未知，根据 API 请求结果判断
	if result.Status == UnlockStatusUnknown {
		if resp.StatusCode == 200 {
			result.Status = UnlockStatusYes
			result.Message = "DAZN is available"
		} else {
			result.Status = UnlockStatusNo
			result.Message = "DAZN is not available"
		}
	}

	// 更新消息包含地区信息
	if result.Status == UnlockStatusYes {
		if result.Region != "" {
			result.Message = "DAZN is available in " + result.Region
		} else {
			result.Message = "DAZN is available"
		}
	}

	return result, nil
}
