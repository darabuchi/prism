package node

import (
	"context"
	"encoding/json"
	"io"

	"github.com/lazygophers/log"
)

// UnlockDisneyPlus 检测 Disney+ 解锁状态
//
// 检测原理：
//   1. 访问 Disney+ 的内容页面
//   2. 检查是否有地区限制提示
//   3. 通过 API 验证可用地区
//
// 返回:
//   - *UnlockResult: 解锁检测结果
//   - error: 检测失败时返回错误
func (n *Node) UnlockDisneyPlus(ctx context.Context) (*UnlockResult, error) {
	result := &UnlockResult{
		Platform: "Disney+",
		Status:   UnlockStatusUnknown,
	}

	// 第一步：访问 Disney+ 主页
	headers := map[string]string{
		"Accept": "application/json",
	}

	resp, err := n.doHTTPRequest(ctx, "GET", "https://global.edge.bamgrid.com/token", headers, nil)
	if err != nil {
		log.Errorf("Disney+ unlock test failed: %v", err)
		result.Status = UnlockStatusFailed
		result.Message = "Failed to connect to Disney+"
		return result, nil
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode == 403 {
		result.Status = UnlockStatusBanned
		result.Message = "Access blocked by Disney+"
		return result, nil
	}

	if resp.StatusCode != 200 {
		result.Status = UnlockStatusNo
		result.Message = "Disney+ is not available"
		return result, nil
	}

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read Disney+ response failed: %v", err)
		result.Status = UnlockStatusFailed
		result.Message = "Failed to read response"
		return result, nil
	}

	// 解析 JSON 响应获取地区信息
	var tokenData map[string]interface{}
	if err := json.Unmarshal(body, &tokenData); err == nil {
		// 检查是否包含地区信息
		if region, ok := tokenData["region"].(string); ok {
			result.Region = region
		}
	}

	// 第二步：检查具体内容的可用性
	// 使用已知的 Disney+ 内容 ID
	resp2, err := n.doHTTPRequest(ctx, "GET", "https://disney.api.edge.bamgrid.com/explore/v1.2/page", headers, nil)
	if err != nil {
		log.Warnf("Disney+ API request failed: %v", err)
		// 即使 API 失败，如果 token 请求成功，仍然认为可能可用
		if resp.StatusCode == 200 {
			result.Status = UnlockStatusYes
			result.Message = "Disney+ appears to be available"
			return result, nil
		}
	} else {
		defer resp2.Body.Close()

		if resp2.StatusCode == 200 {
			result.Status = UnlockStatusYes
			if result.Region != "" {
				result.Message = "Disney+ is available in " + result.Region
			} else {
				result.Message = "Disney+ is available"
			}
		} else if resp2.StatusCode == 403 {
			result.Status = UnlockStatusNo
			result.Message = "Disney+ is not available in your region"
		}
	}

	// 如果状态仍然未知，根据 token 请求结果判断
	if result.Status == UnlockStatusUnknown {
		if resp.StatusCode == 200 {
			result.Status = UnlockStatusYes
			result.Message = "Disney+ is available"
		} else {
			result.Status = UnlockStatusNo
			result.Message = "Disney+ is not available"
		}
	}

	return result, nil
}
