package node

import (
	"context"
	"strings"
	"time"
)

// SpeedTestResult 速度测试结果
type SpeedTestResult struct {
	TestType       string        // 测试类型：http, ookla
	BytesTransfer  int64         // 传输字节数
	Duration       time.Duration // 持续时间
	Speed          float64       // 速度（字节/秒）
	SpeedMbps      float64       // 速度（Mbps）
	ServerID       string        // 服务器 ID（ookla）
	ServerName     string        // 服务器名称（ookla）
	ServerCountry  string        // 服务器国家（ookla）
	Latency        time.Duration // 延迟（ookla）
	Jitter         time.Duration // 抖动（ookla）
}

// SpeedTestDownload 统一的下载速度测试接口
//
// 通过 URL 判断使用哪种测试方式：
//   - 以 "ookla://" 或 "speedtest://" 开头：使用 Ookla speedtest
//   - 其他 HTTP URL：使用普通 HTTP 下载测试
//
// 参数:
//   - ctx: 上下文，用于控制超时和取消
//   - url: 测试 URL
//     * "ookla://" 或 "speedtest://" - 使用 Ookla speedtest，自动选择最佳服务器
//     * "ookla://server_id" - 使用指定服务器 ID 的 Ookla 测试
//     * "http://..." - 普通 HTTP 下载测试
//     * "" - 使用默认 HTTP 测试文件
//   - duration: 测试持续时间（如果为 0，下载完整个文件或完成完整测试）
//
// 返回:
//   - *SpeedTestResult: 测试结果
//   - error: 测试失败时返回错误
//
// 示例:
//
//	// 使用 Ookla speedtest
//	result, err := node.SpeedTestDownload(ctx, "ookla://", 0)
//
//	// 使用 HTTP 测试
//	result, err := node.SpeedTestDownload(ctx, "http://cachefly.cachefly.net/10mb.test", 10*time.Second)
//
//	// 使用默认 HTTP 测试
//	result, err := node.SpeedTestDownload(ctx, "", 10*time.Second)
func (n *Node) SpeedTestDownload(ctx context.Context, url string, duration time.Duration) (*SpeedTestResult, error) {
	// 判断是否为 Ookla 测试
	if strings.HasPrefix(url, "ookla://") || strings.HasPrefix(url, "speedtest://") {
		// 提取服务器 ID
		serverID := ""
		if strings.HasPrefix(url, "ookla://") {
			serverID = strings.TrimPrefix(url, "ookla://")
		} else {
			serverID = strings.TrimPrefix(url, "speedtest://")
		}

		// 使用 Ookla 测试
		download, latency, err := n.SpeedTestOoklaDownload(ctx, serverID)
		if err != nil {
			return nil, err
		}

		return &SpeedTestResult{
			TestType:  "ookla",
			SpeedMbps: download,
			Latency:   latency,
		}, nil
	}

	// 使用普通 HTTP 测试
	result, err := n.SpeedDownloadHTTP(ctx, url, duration)
	if err != nil {
		return nil, err
	}

	return &SpeedTestResult{
		TestType:      "http",
		BytesTransfer: result.BytesDownloaded,
		Duration:      result.Duration,
		Speed:         result.Speed,
		SpeedMbps:     result.SpeedMbps,
	}, nil
}

// SpeedTestUpload 统一的上传速度测试接口
//
// 通过 URL 判断使用哪种测试方式：
//   - 以 "ookla://" 或 "speedtest://" 开头：使用 Ookla speedtest
//   - 其他 HTTP URL：使用普通 HTTP 上传测试
//
// 参数:
//   - ctx: 上下文，用于控制超时和取消
//   - url: 测试 URL
//     * "ookla://" 或 "speedtest://" - 使用 Ookla speedtest，自动选择最佳服务器
//     * "ookla://server_id" - 使用指定服务器 ID 的 Ookla 测试
//     * "http://..." - 普通 HTTP 上传测试
//     * "" - 使用默认 HTTP 测试服务器
//   - dataSize: 上传数据大小（字节），仅用于 HTTP 测试，ookla 测试忽略此参数
//   - duration: 测试持续时间（如果为 0，上传完整个数据或完成完整测试）
//
// 返回:
//   - *SpeedTestResult: 测试结果
//   - error: 测试失败时返回错误
//
// 示例:
//
//	// 使用 Ookla speedtest
//	result, err := node.SpeedTestUpload(ctx, "ookla://", 0, 0)
//
//	// 使用 HTTP 测试
//	result, err := node.SpeedTestUpload(ctx, "http://httpbin.org/post", 1024*1024, 10*time.Second)
//
//	// 使用默认 HTTP 测试
//	result, err := node.SpeedTestUpload(ctx, "", 1024*1024, 10*time.Second)
func (n *Node) SpeedTestUpload(ctx context.Context, url string, dataSize int64, duration time.Duration) (*SpeedTestResult, error) {
	// 判断是否为 Ookla 测试
	if strings.HasPrefix(url, "ookla://") || strings.HasPrefix(url, "speedtest://") {
		// 提取服务器 ID
		serverID := ""
		if strings.HasPrefix(url, "ookla://") {
			serverID = strings.TrimPrefix(url, "ookla://")
		} else {
			serverID = strings.TrimPrefix(url, "speedtest://")
		}

		// 使用 Ookla 测试
		upload, latency, err := n.SpeedTestOoklaUpload(ctx, serverID)
		if err != nil {
			return nil, err
		}

		return &SpeedTestResult{
			TestType:  "ookla",
			SpeedMbps: upload,
			Latency:   latency,
		}, nil
	}

	// 使用普通 HTTP 测试
	result, err := n.SpeedUploadHTTP(ctx, url, dataSize, duration)
	if err != nil {
		return nil, err
	}

	return &SpeedTestResult{
		TestType:      "http",
		BytesTransfer: result.BytesUploaded,
		Duration:      result.Duration,
		Speed:         result.Speed,
		SpeedMbps:     result.SpeedMbps,
	}, nil
}

// SpeedTest 完整的速度测试（包括下载和上传）
//
// 通过 URL 判断使用哪种测试方式：
//   - 以 "ookla://" 或 "speedtest://" 开头：使用 Ookla 完整测试
//   - 其他：不支持（需要分别调用 SpeedTestDownload 和 SpeedTestUpload）
//
// 参数:
//   - ctx: 上下文，用于控制超时和取消
//   - url: 测试 URL（必须是 ookla:// 或 speedtest://）
//
// 返回:
//   - download: 下载速度测试结果
//   - upload: 上传速度测试结果
//   - error: 测试失败时返回错误
//
// 示例:
//
//	// 使用 Ookla 完整测试
//	download, upload, err := node.SpeedTest(ctx, "ookla://")
func (n *Node) SpeedTest(ctx context.Context, url string) (download *SpeedTestResult, upload *SpeedTestResult, err error) {
	// 判断是否为 Ookla 测试
	if !strings.HasPrefix(url, "ookla://") && !strings.HasPrefix(url, "speedtest://") {
		return nil, nil, ErrUnsupportedSpeedTestURL
	}

	// 提取服务器 ID
	serverID := ""
	if strings.HasPrefix(url, "ookla://") {
		serverID = strings.TrimPrefix(url, "ookla://")
	} else {
		serverID = strings.TrimPrefix(url, "speedtest://")
	}

	// 使用 Ookla 完整测试
	result, err := n.SpeedTestOokla(ctx, serverID)
	if err != nil {
		return nil, nil, err
	}

	downloadResult := &SpeedTestResult{
		TestType:      "ookla",
		SpeedMbps:     result.Download,
		Latency:       result.Latency,
		Jitter:        result.Jitter,
		ServerID:      result.ServerID,
		ServerName:    result.ServerName,
		ServerCountry: result.ServerCountry,
	}

	uploadResult := &SpeedTestResult{
		TestType:      "ookla",
		SpeedMbps:     result.Upload,
		Latency:       result.Latency,
		Jitter:        result.Jitter,
		ServerID:      result.ServerID,
		ServerName:    result.ServerName,
		ServerCountry: result.ServerCountry,
	}

	return downloadResult, uploadResult, nil
}
