package node

import (
	"context"
	"encoding/json"
	"io"

	"github.com/lazygophers/log"
)

// UnlockBilibili 检测 Bilibili 解锁状态
//
// 检测原理：
//   1. 访问 Bilibili 的播放 API
//   2. 检查是否有地区限制
//   3. 测试港澳台内容的可用性
//
// 返回:
//   - *UnlockResult: 解锁检测结果
//   - error: 检测失败时返回错误
func (n *Node) UnlockBilibili(ctx context.Context) (*UnlockResult, error) {
	result := &UnlockResult{
		Platform: "Bilibili",
		Status:   UnlockStatusUnknown,
	}

	// Bilibili 地区检测
	// 使用已知的仅限港澳台的番剧进行测试
	headers := map[string]string{
		"Accept":     "application/json",
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
	}

	// 第一步：测试港澳台番剧（使用 Season ID 测试）
	// 使用一个已知的港澳台限定内容
	resp, err := n.doHTTPRequest(ctx, "GET", "https://api.bilibili.com/pgc/player/web/playurl?avid=50762638&cid=100279344&qn=0&type=&otype=json&ep_id=268176", headers, nil)
	if err != nil {
		log.Errorf("Bilibili unlock test failed: %v", err)
		result.Status = UnlockStatusFailed
		result.Message = "Failed to connect to Bilibili"
		return result, nil
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode == 403 {
		result.Status = UnlockStatusBanned
		result.Message = "Access blocked by Bilibili"
		return result, nil
	}

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read Bilibili response failed: %v", err)
		result.Status = UnlockStatusFailed
		result.Message = "Failed to read response"
		return result, nil
	}

	// 解析 JSON 响应
	var apiData map[string]interface{}
	if err := json.Unmarshal(body, &apiData); err != nil {
		log.Errorf("parse Bilibili response failed: %v", err)
		result.Status = UnlockStatusFailed
		result.Message = "Failed to parse response"
		return result, nil
	}

	// 检查响应代码
	code, ok := apiData["code"].(float64)
	if !ok {
		result.Status = UnlockStatusFailed
		result.Message = "Invalid response format"
		return result, nil
	}

	// code 含义：
	// 0: 成功，说明可以观看港澳台内容
	// -10403: 地区限制
	if code == 0 {
		result.Status = UnlockStatusYes
		result.Region = "HK/TW/MO"
		result.Message = "Bilibili is available (HK/TW/MO content accessible)"
		return result, nil
	} else if code == -10403 {
		// 地区限制，但这不代表完全不可用
		// 可能只是无法观看港澳台内容
		result.Status = UnlockStatusNo
		result.Message = "Bilibili HK/TW/MO content is not available"
	}

	// 第二步：测试大陆番剧
	resp2, err := n.doHTTPRequest(ctx, "GET", "https://api.bilibili.com/pgc/player/web/playurl?avid=18281381&cid=29892777&qn=0&type=&otype=json&ep_id=183799", headers, nil)
	if err != nil {
		log.Warnf("Bilibili mainland content test failed: %v", err)
		return result, nil
	}
	defer resp2.Body.Close()

	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		log.Warnf("read Bilibili mainland response failed: %v", err)
		return result, nil
	}

	var apiData2 map[string]interface{}
	if err := json.Unmarshal(body2, &apiData2); err != nil {
		log.Warnf("parse Bilibili mainland response failed: %v", err)
		return result, nil
	}

	code2, ok := apiData2["code"].(float64)
	if ok && code2 == 0 {
		// 可以观看大陆内容
		if result.Status == UnlockStatusNo {
			// 大陆内容可用，但港澳台内容不可用
			result.Status = UnlockStatusYes
			result.Region = "CN"
			result.Message = "Bilibili is available (Mainland China content only)"
		} else {
			result.Status = UnlockStatusYes
			result.Region = "CN"
			result.Message = "Bilibili is available"
		}
	} else if code2 == -10403 {
		// 大陆内容也受限
		if result.Status == UnlockStatusNo {
			result.Message = "Bilibili is not available"
		}
	}

	// 如果两个都无法访问，检查服务是否完全不可用
	if result.Status == UnlockStatusUnknown {
		result.Status = UnlockStatusNo
		result.Message = "Bilibili is not available"
	}

	return result, nil
}
