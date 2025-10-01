package node

import (
	"context"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/lazygophers/lrpc/middleware/xerror"
	"github.com/lazygophers/utils/candy"
	"github.com/metacubex/mihomo/constant"
)

// DelayHTTP 测试 HTTP 延迟
//
// 通过节点代理访问指定的 HTTP URL，测量从建立连接到接收响应的总时间
// 常用的测试 URL：
//   - http://www.gstatic.com/generate_204 (Google)
//   - http://cp.cloudflare.com/generate_204 (Cloudflare)
//
// 参数:
//   - ctx: 上下文，用于控制超时和取消
//   - url: 测试 URL
//
// 返回:
//   - time.Duration: HTTP 延迟时间
//   - error: 测试失败时返回错误
//
// 示例:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	delay, err := node.DelayHTTP(ctx, "http://www.gstatic.com/generate_204")
//	if err != nil {
//	    log.Errorf("HTTP 延迟测试失败: %v", err)
//	    return
//	}
//	log.Infof("HTTP 延迟: %v", delay)
func (n *Node) DelayHTTP(ctx context.Context, url string) (time.Duration, error) {
	if url == "" {
		url = "http://www.gstatic.com/generate_204"
	}

	// 创建 HTTP 客户端，使用节点作为传输层
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, portStr, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}

			// 解析端口
			var port uint16
			if portStr == "" {
				port = 80
			} else {
				port = uint16(candy.ToInt(portStr))
			}

			// 构建 Mihomo 连接元数据
			metadata := &constant.Metadata{
				NetWork: constant.TCP,
				Host:    host,
				DstPort: port,
			}

			// 使用节点的 DialContext 建立代理连接
			return n.adapter.DialContext(ctx, metadata)
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	// 记录开始时间
	start := time.Now()

	// 发送 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, xerror.WrapError(err, xerror.ErrSystemError, "create HTTP request failed")
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, xerror.WrapError(err, xerror.ErrSystemError, "HTTP request failed")
	}
	defer resp.Body.Close()

	// 读取响应体（确保完整接收）
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return 0, xerror.WrapError(err, xerror.ErrSystemError, "read HTTP response failed")
	}

	// 计算延迟
	delay := time.Since(start)

	return delay, nil
}
