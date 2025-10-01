package node

import (
	"context"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/xerror"
	"github.com/lazygophers/utils/candy"
	"github.com/metacubex/mihomo/constant"
)

// UnlockStatus 解锁状态
type UnlockStatus int

const (
	UnlockStatusUnknown   UnlockStatus = 0 // 未知
	UnlockStatusYes       UnlockStatus = 1 // 解锁
	UnlockStatusNo        UnlockStatus = 2 // 未解锁（有地区限制）
	UnlockStatusBanned    UnlockStatus = 3 // 被封禁
	UnlockStatusFailed    UnlockStatus = 4 // 检测失败
	UnlockStatusNotAvailable UnlockStatus = 5 // 服务不可用
)

func (s UnlockStatus) String() string {
	switch s {
	case UnlockStatusYes:
		return "Yes"
	case UnlockStatusNo:
		return "No"
	case UnlockStatusBanned:
		return "Banned"
	case UnlockStatusFailed:
		return "Failed"
	case UnlockStatusNotAvailable:
		return "Not Available"
	default:
		return "Unknown"
	}
}

// UnlockResult 解锁检测结果
type UnlockResult struct {
	Platform string       // 平台名称
	Status   UnlockStatus // 解锁状态
	Region   string       // 地区代码（如 US, JP, HK）
	Message  string       // 详细信息
}

// createHTTPClient 创建使用节点代理的 HTTP 客户端
func (n *Node) createHTTPClient(timeout time.Duration) *http.Client {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, portStr, err := net.SplitHostPort(addr)
			if err != nil {
				log.Errorf("split host port failed: %v, addr: %s", err, addr)
				return nil, err
			}

			var port uint16
			if portStr == "" {
				if strings.Contains(addr, "https") {
					port = 443
				} else {
					port = 80
				}
			} else {
				port = uint16(candy.ToInt(portStr))
			}

			metadata := &constant.Metadata{
				NetWork: constant.TCP,
				Host:    host,
				DstPort: port,
			}

			conn, err := n.adapter.DialContext(ctx, metadata)
			if err != nil {
				log.Errorf("dial context failed: %v, host: %s, port: %d", err, host, port)
				return nil, err
			}
			return conn, nil
		},
		DisableKeepAlives: false,
		IdleConnTimeout:   90 * time.Second,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 允许重定向，但最多 10 次
			if len(via) >= 10 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
}

// doHTTPRequest 执行 HTTP 请求的辅助函数
func (n *Node) doHTTPRequest(ctx context.Context, method, url string, headers map[string]string, body io.Reader) (*http.Response, error) {
	client := n.createHTTPClient(15 * time.Second)

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		log.Errorf("create HTTP request failed: %v, url: %s", err, url)
		return nil, xerror.WrapError(err, xerror.ErrSystemError, "create HTTP request failed")
	}

	// 设置默认 headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	// 设置自定义 headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("HTTP request failed: %v, url: %s", err, url)
		return nil, xerror.WrapError(err, xerror.ErrSystemError, "HTTP request failed")
	}

	return resp, nil
}

// UnlockTest 统一的解锁测试接口
//
// 支持的平台：
//   - netflix
//   - disney / disneyplus
//   - youtube / youtube-premium
//   - hulu
//   - hbo / hbomax
//   - primevideo / amazon
//   - bilibili
//   - bahamut
//   - abematv
//   - dazn
//   - tiktok
//
// 参数:
//   - ctx: 上下文，用于控制超时和取消
//   - platform: 平台名称（不区分大小写）
//
// 返回:
//   - *UnlockResult: 解锁检测结果
//   - error: 检测失败时返回错误
//
// 示例:
//
//	result, err := node.UnlockTest(ctx, "netflix")
//	if err != nil {
//	    log.Errorf("unlock test failed: %v", err)
//	    return
//	}
//	log.Infof("Netflix unlock status: %s, region: %s", result.Status, result.Region)
func (n *Node) UnlockTest(ctx context.Context, platform string) (*UnlockResult, error) {
	platform = strings.ToLower(platform)

	switch platform {
	case "netflix":
		return n.UnlockNetflix(ctx)
	case "disney", "disneyplus", "disney+":
		return n.UnlockDisneyPlus(ctx)
	case "youtube", "youtube-premium", "ytb":
		return n.UnlockYouTubePremium(ctx)
	case "hulu":
		return n.UnlockHulu(ctx)
	case "hbo", "hbomax":
		return n.UnlockHBOMax(ctx)
	case "primevideo", "amazon", "prime":
		return n.UnlockPrimeVideo(ctx)
	case "bilibili":
		return n.UnlockBilibili(ctx)
	case "bahamut":
		return n.UnlockBahamut(ctx)
	case "abematv", "abema":
		return n.UnlockAbemaTV(ctx)
	case "dazn":
		return n.UnlockDAZN(ctx)
	case "tiktok":
		return n.UnlockTikTok(ctx)
	default:
		return nil, xerror.New(xerror.ErrSystemError, "unsupported platform: "+platform)
	}
}

// UnlockTestAll 测试所有支持的平台
//
// 参数:
//   - ctx: 上下文，用于控制超时和取消
//
// 返回:
//   - map[string]*UnlockResult: 平台名称到解锁结果的映射
//   - error: 检测失败时返回错误
func (n *Node) UnlockTestAll(ctx context.Context) (map[string]*UnlockResult, error) {
	platforms := []string{
		"netflix",
		"disney",
		"youtube",
		"hulu",
		"hbo",
		"primevideo",
		"bilibili",
		"bahamut",
		"abematv",
		"dazn",
		"tiktok",
	}

	results := make(map[string]*UnlockResult)

	for _, platform := range platforms {
		result, err := n.UnlockTest(ctx, platform)
		if err != nil {
			log.Warnf("unlock test failed for %s: %v", platform, err)
			results[platform] = &UnlockResult{
				Platform: platform,
				Status:   UnlockStatusFailed,
				Message:  err.Error(),
			}
			continue
		}

		results[platform] = result
	}

	return results, nil
}
