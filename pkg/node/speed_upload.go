package node

import (
	"bytes"
	"context"
	"crypto/rand"
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

// SpeedUploadResult 上传速度测试结果
type SpeedUploadResult struct {
	BytesUploaded int64         // 上传的字节数
	Duration      time.Duration // 持续时间
	Speed         float64       // 速度（字节/秒）
	SpeedMbps     float64       // 速度（Mbps）
}

// generateRandomData 生成随机数据用于上传测试
func generateRandomData(size int64) (io.Reader, error) {
	buf := make([]byte, size)
	_, err := rand.Read(buf)
	if err != nil {
		log.Errorf("generate random data failed: %v, size: %d", err, size)
		return nil, err
	}
	return bytes.NewReader(buf), nil
}

// SpeedUploadHTTP 测试 HTTP 上传速度
//
// 通过节点代理上传数据到指定的测试服务器，测量上传速度
//
// 参数:
//   - ctx: 上下文，用于控制超时和取消
//   - url: 测试服务器 URL（如果为空，使用默认测试服务器）
//   - dataSize: 上传数据大小（字节），如果为 0，使用默认大小（1MB）
//   - duration: 测试持续时间（如果为 0，上传完整个数据）
//
// 返回:
//   - *SpeedUploadResult: 上传速度测试结果
//   - error: 测试失败时返回错误
//
// 示例:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//	result, err := node.SpeedUploadHTTP(ctx, "", 1024*1024*10, 10*time.Second)
//	if err != nil {
//	    log.Errorf("upload speed test failed: %v", err)
//	    return
//	}
//	log.Infof("upload speed: %.2f Mbps", result.SpeedMbps)
func (n *Node) SpeedUploadHTTP(ctx context.Context, url string, dataSize int64, duration time.Duration) (*SpeedUploadResult, error) {
	if url == "" {
		// 默认使用 httpbin.org 的 POST 端点
		url = "http://httpbin.org/post"
	}

	if dataSize <= 0 {
		// 默认 1MB
		dataSize = 1024 * 1024
	}

	// 生成随机数据
	data, err := generateRandomData(dataSize)
	if err != nil {
		return nil, xerror.WrapError(err, xerror.ErrSystemError, "generate random data failed")
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

	// 创建一个包装的 reader 来跟踪上传进度
	pr, pw := io.Pipe()
	var totalBytes int64

	// 如果指定了持续时间，创建一个定时器
	var timer *time.Timer
	var timerCh <-chan time.Time
	if duration > 0 {
		timer = time.NewTimer(duration)
		timerCh = timer.C
		defer timer.Stop()
	}

	// 在 goroutine 中写入数据
	uploadDone := make(chan error, 1)
	start := time.Now()

	go func() {
		defer pw.Close()

		buf := make([]byte, 32*1024) // 32KB buffer

		for {
			select {
			case <-ctx.Done():
				uploadDone <- ctx.Err()
				return
			case <-timerCh:
				// 达到指定持续时间
				uploadDone <- nil
				return
			default:
				n, err := data.Read(buf)
				if n > 0 {
					written, writeErr := pw.Write(buf[:n])
					if writeErr != nil {
						uploadDone <- writeErr
						return
					}
					totalBytes += int64(written)
				}

				if err != nil {
					if err == io.EOF {
						uploadDone <- nil
						return
					}
					uploadDone <- err
					return
				}
			}
		}
	}()

	// 发送 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, pr)
	if err != nil {
		log.Errorf("create HTTP request failed: %v, url: %s", err, url)
		return nil, xerror.WrapError(err, xerror.ErrSystemError, "create HTTP request failed")
	}

	req.ContentLength = dataSize

	resp, err := client.Do(req)
	if err != nil {
		// 检查是否是正常的上传完成
		uploadErr := <-uploadDone
		elapsed := time.Since(start)

		if uploadErr == nil || uploadErr == io.EOF {
			// 上传成功完成
			var speed, speedMbps float64
			if elapsed.Seconds() > 0 {
				speed = float64(totalBytes) / elapsed.Seconds()
				speedMbps = speed * 8 / 1000000
			}

			return &SpeedUploadResult{
				BytesUploaded: totalBytes,
				Duration:      elapsed,
				Speed:         speed,
				SpeedMbps:     speedMbps,
			}, nil
		}

		log.Errorf("HTTP request failed: %v, url: %s", err, url)
		return nil, xerror.WrapError(err, xerror.ErrSystemError, "HTTP request failed")
	}
	defer resp.Body.Close()

	// 等待上传完成
	uploadErr := <-uploadDone
	elapsed := time.Since(start)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		err := fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		log.Errorf("HTTP request failed: %v, url: %s", err, url)
		return nil, xerror.WrapError(err, xerror.ErrSystemError, "unexpected status code")
	}

	if uploadErr != nil && uploadErr != io.EOF {
		log.Errorf("upload failed: %v, url: %s", uploadErr, url)
		return nil, xerror.WrapError(uploadErr, xerror.ErrSystemError, "upload failed")
	}

	// 读取并丢弃响应体
	_, _ = io.Copy(io.Discard, resp.Body)

	// 计算上传速度
	var speed, speedMbps float64
	if elapsed.Seconds() > 0 {
		speed = float64(totalBytes) / elapsed.Seconds()
		speedMbps = speed * 8 / 1000000
	}

	return &SpeedUploadResult{
		BytesUploaded: totalBytes,
		Duration:      elapsed,
		Speed:         speed,
		SpeedMbps:     speedMbps,
	}, nil
}
