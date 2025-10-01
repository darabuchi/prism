package node

import (
	"context"
	"io"
	"strings"

	"github.com/lazygophers/log"
)

// UnlockHulu 检测 Hulu 解锁状态
//
// 检测原理：
//   1. 访问 Hulu 的播放页面
//   2. 检查地区限制信息
//   3. 验证内容可用性
//
// 注意：Hulu 主要在美国和日本提供服务
//
// 返回:
//   - *UnlockResult: 解锁检测结果
//   - error: 检测失败时返回错误
func (n *Node) UnlockHulu(ctx context.Context) (*UnlockResult, error) {
	result := &UnlockResult{
		Platform: "Hulu",
		Status:   UnlockStatusUnknown,
	}

	// Hulu 的地区检测接口
	// 访问播放页面来检测可用性
	resp, err := n.doHTTPRequest(ctx, "GET", "https://auth.hulu.com/v2/livingroom/password/authenticate", nil, nil)
	if err != nil {
		log.Errorf("Hulu unlock test failed: %v", err)
		result.Status = UnlockStatusFailed
		result.Message = "Failed to connect to Hulu"
		return result, nil
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode == 403 {
		result.Status = UnlockStatusBanned
		result.Message = "Access blocked by Hulu"
		return result, nil
	}

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read Hulu response failed: %v", err)
		result.Status = UnlockStatusFailed
		result.Message = "Failed to read response"
		return result, nil
	}

	bodyStr := string(body)

	// 检测地区限制
	// Hulu 会返回特定的错误信息表示地区不可用
	if strings.Contains(bodyStr, "not available") || strings.Contains(bodyStr, "region") {
		result.Status = UnlockStatusNo
		result.Message = "Hulu is not available in your region"
		return result, nil
	}

	// 检测是否被识别为使用代理/VPN
	if strings.Contains(bodyStr, "proxy") || strings.Contains(bodyStr, "vpn") {
		result.Status = UnlockStatusBanned
		result.Message = "Proxy/VPN detected by Hulu"
		return result, nil
	}

	// 第二步：检查播放 API
	resp2, err := n.doHTTPRequest(ctx, "GET", "https://play.hulu.com/v6/playlist/home", nil, nil)
	if err != nil {
		// 如果 playlist API 失败，但前面的检测通过，仍然可能可用
		log.Warnf("Hulu playlist API request failed: %v", err)
		if resp.StatusCode == 200 || resp.StatusCode == 400 {
			// 400 是因为没有提供认证信息，但说明 API 可访问
			result.Status = UnlockStatusYes
			result.Message = "Hulu appears to be available"
			return result, nil
		}
	} else {
		defer resp2.Body.Close()

		// 检查响应
		if resp2.StatusCode == 200 || resp2.StatusCode == 400 {
			// 400 说明需要认证，但服务可用
			result.Status = UnlockStatusYes
		} else if resp2.StatusCode == 403 {
			result.Status = UnlockStatusNo
			result.Message = "Hulu is not available in your region"
			return result, nil
		}
	}

	// 尝试检测具体地区
	// Hulu 主要在美国服务
	if result.Status == UnlockStatusYes {
		// 通过第三方服务或响应头判断地区
		// Hulu 主要服务区域：US (美国)、JP (日本)
		result.Region = "US" // 默认假设美国，实际应该通过其他方式判断
		result.Message = "Hulu is available"
	}

	// 如果状态仍然未知
	if result.Status == UnlockStatusUnknown {
		if resp.StatusCode == 200 || resp.StatusCode == 400 {
			result.Status = UnlockStatusYes
			result.Message = "Hulu is available"
		} else {
			result.Status = UnlockStatusNo
			result.Message = "Hulu is not available"
		}
	}

	return result, nil
}
