package node

import (
	"context"
	"io"
	"strings"

	"github.com/lazygophers/log"
)

// UnlockBahamut 检测 Bahamut 动画疯解锁状态
//
// 检测原理：
//   1. 访问 Bahamut 的 API
//   2. 检查是否有地区限制
//   3. 验证台湾地区内容的可用性
//
// 注意：Bahamut 动画疯主要服务台湾和港澳地区
//
// 返回:
//   - *UnlockResult: 解锁检测结果
//   - error: 检测失败时返回错误
func (n *Node) UnlockBahamut(ctx context.Context) (*UnlockResult, error) {
	result := &UnlockResult{
		Platform: "Bahamut",
		Status:   UnlockStatusUnknown,
	}

	// Bahamut 的地区检测
	// 使用 Bahamut 动画疯的 API
	headers := map[string]string{
		"Accept":     "application/json, text/javascript, */*; q=0.01",
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
	}

	// 访问 Bahamut 的主页
	resp, err := n.doHTTPRequest(ctx, "GET", "https://ani.gamer.com.tw/", headers, nil)
	if err != nil {
		log.Errorf("Bahamut unlock test failed: %v", err)
		result.Status = UnlockStatusFailed
		result.Message = "Failed to connect to Bahamut"
		return result, nil
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode == 403 {
		result.Status = UnlockStatusBanned
		result.Message = "Access blocked by Bahamut"
		return result, nil
	}

	if resp.StatusCode != 200 {
		result.Status = UnlockStatusNo
		result.Message = "Bahamut is not available"
		return result, nil
	}

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read Bahamut response failed: %v", err)
		result.Status = UnlockStatusFailed
		result.Message = "Failed to read response"
		return result, nil
	}

	bodyStr := string(body)

	// 检测地区限制
	// Bahamut 会显示地区限制信息
	if strings.Contains(bodyStr, "地區限制") ||
	   strings.Contains(bodyStr, "僅限台灣地區") ||
	   strings.Contains(bodyStr, "region") {
		result.Status = UnlockStatusNo
		result.Message = "Bahamut is not available in your region"
		return result, nil
	}

	// 第二步：测试具体动画内容
	// 使用一个已知的动画 SN 进行测试
	resp2, err := n.doHTTPRequest(ctx, "GET", "https://ani.gamer.com.tw/ajax/token.php?adID=89422&sn=14667", headers, nil)
	if err != nil {
		log.Warnf("Bahamut content test failed: %v", err)
		// 如果内容测试失败，但主页可访问，仍认为可能可用
		if resp.StatusCode == 200 {
			result.Status = UnlockStatusYes
			result.Region = "TW"
			result.Message = "Bahamut appears to be available"
			return result, nil
		}
	} else {
		defer resp2.Body.Close()

		contentBody, _ := io.ReadAll(resp2.Body)
		contentStr := string(contentBody)

		// 检查内容响应
		if strings.Contains(contentStr, "\"error\":true") ||
		   strings.Contains(contentStr, "地區限制") ||
		   strings.Contains(contentStr, "地区限制") {
			result.Status = UnlockStatusNo
			result.Message = "Bahamut is not available in your region"
			return result, nil
		}

		// 如果没有错误信息，且能获取到 token，说明可用
		if strings.Contains(contentStr, "token") || resp2.StatusCode == 200 {
			result.Status = UnlockStatusYes
			result.Region = "TW"
			result.Message = "Bahamut is available in TW"
			return result, nil
		}
	}

	// 如果状态仍然未知，根据主页请求结果判断
	if result.Status == UnlockStatusUnknown {
		if resp.StatusCode == 200 && !strings.Contains(bodyStr, "地區限制") {
			result.Status = UnlockStatusYes
			result.Region = "TW"
			result.Message = "Bahamut is available"
		} else {
			result.Status = UnlockStatusNo
			result.Message = "Bahamut is not available"
		}
	}

	return result, nil
}
