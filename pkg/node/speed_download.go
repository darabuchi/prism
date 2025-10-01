package node

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/xerror"
	"github.com/lazygophers/utils/candy"
	"github.com/metacubex/mihomo/constant"
)

// SpeedDownloadResult 下载速度测试结果
type SpeedDownloadResult struct {
	BytesDownloaded int64         // 下载的字节数
	Duration        time.Duration // 持续时间
	Speed           float64       // 速度（字节/秒）
	SpeedMbps       float64       // 速度（Mbps）
}

// SpeedDownloadHTTP 测试 HTTP 下载速度
//
// 通过节点代理下载指定的测试文件，测量下载速度
// 常用的测试文件：
//   - 小文件（1MB）：适合快速测试
//   - 中等文件（10MB）：一般测试
//   - 大文件（100MB+）：准确测试
//
// 参数:
//   - ctx: 上下文，用于控制超时和取消
//   - url: 测试文件 URL（如果为空，使用默认测试文件）
//   - duration: 测试持续时间（如果为 0，下载完整个文件）
//
// 返回:
//   - *SpeedDownloadResult: 下载速度测试结果
//   - error: 测试失败时返回错误
//
// 示例:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//	result, err := node.SpeedDownloadHTTP(ctx, "", 10*time.Second)
//	if err != nil {
//	    log.Errorf("download speed test failed: %v", err)
//	    return
//	}
//	log.Infof("download speed: %.2f Mbps", result.SpeedMbps)
func (n *Node) SpeedDownloadHTTP(ctx context.Context, url string, duration time.Duration) (*SpeedDownloadResult, error) {
	if url == "" {
		// 默认使用 10MB 测试文件
		url = "http://cachefly.cachefly.net/10mb.test"
	}

	// 创建 HTTP 客户端，使用节点作为传输层
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, portStr, err := net.SplitHostPort(addr)
			if err != nil {
				log.Errorf("split host port failed: %v, addr: %s", err, addr)
				return nil, err
			}

			var port uint16
			if portStr == "" {
				port = 80
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
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   60 * time.Second,
	}

	// 发送 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("create HTTP request failed: %v, url: %s", err, url)
		return nil, xerror.WrapError(err, xerror.ErrSystemError, "create HTTP request failed")
	}

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("HTTP request failed: %v, url: %s", err, url)
		return nil, xerror.WrapError(err, xerror.ErrSystemError, "HTTP request failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		err := fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		log.Errorf("HTTP request failed: %v, url: %s", err, url)
		return nil, xerror.WrapError(err, xerror.ErrSystemError, "unexpected status code")
	}

	// 读取响应体并计算下载速度
	var totalBytes int64
	buf := make([]byte, 32*1024) // 32KB buffer

	// 如果指定了持续时间，创建一个定时器
	var timer *time.Timer
	var timerCh <-chan time.Time
	if duration > 0 {
		timer = time.NewTimer(duration)
		timerCh = timer.C
		defer timer.Stop()
	}

	for {
		select {
		case <-ctx.Done():
			// 上下文取消
			elapsed := time.Since(start)
			speed := float64(totalBytes) / elapsed.Seconds()
			speedMbps := speed * 8 / 1000000

			return &SpeedDownloadResult{
				BytesDownloaded: totalBytes,
				Duration:        elapsed,
				Speed:           speed,
				SpeedMbps:       speedMbps,
			}, nil

		case <-timerCh:
			// 达到指定持续时间
			elapsed := time.Since(start)
			speed := float64(totalBytes) / elapsed.Seconds()
			speedMbps := speed * 8 / 1000000

			return &SpeedDownloadResult{
				BytesDownloaded: totalBytes,
				Duration:        elapsed,
				Speed:           speed,
				SpeedMbps:       speedMbps,
			}, nil

		default:
			n, err := resp.Body.Read(buf)
			if n > 0 {
				totalBytes += int64(n)
			}

			if err != nil {
				if err == io.EOF {
					// 下载完成
					elapsed := time.Since(start)
					speed := float64(totalBytes) / elapsed.Seconds()
					speedMbps := speed * 8 / 1000000

					return &SpeedDownloadResult{
						BytesDownloaded: totalBytes,
						Duration:        elapsed,
						Speed:           speed,
						SpeedMbps:       speedMbps,
					}, nil
				}

				log.Errorf("read response body failed: %v, url: %s", err, url)
				return nil, xerror.WrapError(err, xerror.ErrSystemError, "read response body failed")
			}
		}
	}
}
