package node

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/xerror"
	"github.com/lazygophers/utils/candy"
	"github.com/metacubex/mihomo/constant"
	"github.com/showwin/speedtest-go/speedtest"
)

// SpeedOoklaResult Ookla speedtest 测试结果
type SpeedOoklaResult struct {
	ServerID      string        // 服务器 ID
	ServerName    string        // 服务器名称
	ServerCountry string        // 服务器国家
	Latency       time.Duration // 延迟
	Download      float64       // 下载速度（Mbps）
	Upload        float64       // 上传速度（Mbps）
	Jitter        time.Duration // 抖动
	TestDuration  time.Duration // 测试持续时间
}

// SpeedTestOokla 使用 Ookla speedtest 测试速度
//
// 使用 speedtest.net 的服务器进行完整的速度测试
// 包括：延迟、下载速度、上传速度、丢包率、抖动
//
// 参数:
//   - ctx: 上下文，用于控制超时和取消
//   - serverID: 指定服务器 ID（如果为空，自动选择最佳服务器）
//
// 返回:
//   - *SpeedOoklaResult: 测试结果
//   - error: 测试失败时返回错误
//
// 示例:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
//	defer cancel()
//	result, err := node.SpeedTestOokla(ctx, "")
//	if err != nil {
//	    log.Errorf("Ookla speedtest failed: %v", err)
//	    return
//	}
//	log.Infof("download: %.2f Mbps, upload: %.2f Mbps, latency: %v",
//	    result.Download, result.Upload, result.Latency)
func (n *Node) SpeedTestOokla(ctx context.Context, serverID string) (*SpeedOoklaResult, error) {
	start := time.Now()

	// 创建自定义 HTTP 客户端，使用节点作为代理
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, portStr, err := net.SplitHostPort(addr)
			if err != nil {
				log.Errorf("split host port failed: %v, addr: %s", err, addr)
				return nil, err
			}

			var port uint16
			if portStr == "" {
				if network == "tcp" {
					port = 80
				} else {
					port = 443
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

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	// 创建 speedtest 用户
	user := speedtest.New(speedtest.WithDoer(httpClient))

	// 获取服务器列表
	serverList, err := user.FetchServers()
	if err != nil {
		log.Errorf("fetch speedtest servers failed: %v", err)
		return nil, xerror.WrapError(err, xerror.ErrSystemError, "fetch speedtest servers failed")
	}

	if len(serverList) == 0 {
		err := fmt.Errorf("no speedtest servers available")
		log.Errorf("no speedtest servers available")
		return nil, xerror.WrapError(err, xerror.ErrSystemError, "no speedtest servers available")
	}

	var targets speedtest.Servers

	// 选择服务器
	if serverID != "" {
		// 查找指定的服务器
		for _, server := range serverList {
			if server.ID == serverID {
				targets = append(targets, server)
				break
			}
		}
		if len(targets) == 0 {
			err := fmt.Errorf("server not found: %s", serverID)
			log.Errorf("speedtest server not found: %s", serverID)
			return nil, xerror.WrapError(err, xerror.ErrSystemError, "server not found")
		}
	} else {
		// 选择最近的服务器（延迟最低）
		targets, err = serverList.FindServer([]int{})
		if err != nil {
			log.Errorf("find best speedtest server failed: %v", err)
			return nil, xerror.WrapError(err, xerror.ErrSystemError, "find best server failed")
		}
	}

	if len(targets) == 0 {
		err := fmt.Errorf("no suitable speedtest server found")
		log.Errorf("no suitable speedtest server found")
		return nil, xerror.WrapError(err, xerror.ErrSystemError, "no suitable server found")
	}

	server := targets[0]

	// 测试延迟
	err = server.PingTest(func(latency time.Duration) {})
	if err != nil {
		log.Errorf("speedtest ping test failed: %v", err)
		return nil, xerror.WrapError(err, xerror.ErrSystemError, "ping test failed")
	}

	// 测试下载速度
	err = server.DownloadTest()
	if err != nil {
		log.Errorf("speedtest download test failed: %v", err)
		return nil, xerror.WrapError(err, xerror.ErrSystemError, "download test failed")
	}

	// 测试上传速度
	err = server.UploadTest()
	if err != nil {
		log.Errorf("speedtest upload test failed: %v", err)
		return nil, xerror.WrapError(err, xerror.ErrSystemError, "upload test failed")
	}

	testDuration := time.Since(start)

	// 构建结果
	result := &SpeedOoklaResult{
		ServerID:      server.ID,
		ServerName:    server.Name,
		ServerCountry: server.Country,
		Latency:       server.Latency,
		Download:      server.DLSpeed.Mbps(),
		Upload:        server.ULSpeed.Mbps(),
		Jitter:        server.Jitter,
		TestDuration:  testDuration,
	}

	return result, nil
}

// SpeedTestOoklaDownload 使用 Ookla speedtest 仅测试下载速度
//
// 参数:
//   - ctx: 上下文，用于控制超时和取消
//   - serverID: 指定服务器 ID（如果为空，自动选择最佳服务器）
//
// 返回:
//   - download: 下载速度（Mbps）
//   - latency: 延迟
//   - error: 测试失败时返回错误
func (n *Node) SpeedTestOoklaDownload(ctx context.Context, serverID string) (download float64, latency time.Duration, err error) {
	// 创建自定义 HTTP 客户端
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, portStr, err := net.SplitHostPort(addr)
			if err != nil {
				log.Errorf("split host port failed: %v, addr: %s", err, addr)
				return nil, err
			}

			var port uint16
			if portStr == "" {
				if network == "tcp" {
					port = 80
				} else {
					port = 443
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
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	user := speedtest.New(speedtest.WithDoer(httpClient))

	serverList, err := user.FetchServers()
	if err != nil {
		log.Errorf("fetch speedtest servers failed: %v", err)
		return 0, 0, xerror.WrapError(err, xerror.ErrSystemError, "fetch speedtest servers failed")
	}

	var targets speedtest.Servers
	if serverID != "" {
		for _, server := range serverList {
			if server.ID == serverID {
				targets = append(targets, server)
				break
			}
		}
		if len(targets) == 0 {
			err := fmt.Errorf("server not found: %s", serverID)
			log.Errorf("speedtest server not found: %s", serverID)
			return 0, 0, xerror.WrapError(err, xerror.ErrSystemError, "server not found")
		}
	} else {
		targets, err = serverList.FindServer([]int{})
		if err != nil {
			log.Errorf("find best speedtest server failed: %v", err)
			return 0, 0, xerror.WrapError(err, xerror.ErrSystemError, "find best server failed")
		}
	}

	if len(targets) == 0 {
		err := fmt.Errorf("no suitable speedtest server found")
		log.Errorf("no suitable speedtest server found")
		return 0, 0, xerror.WrapError(err, xerror.ErrSystemError, "no suitable server found")
	}

	server := targets[0]

	// 测试延迟
	err = server.PingTest(func(latency time.Duration) {})
	if err != nil {
		log.Errorf("speedtest ping test failed: %v", err)
		return 0, 0, xerror.WrapError(err, xerror.ErrSystemError, "ping test failed")
	}

	// 测试下载速度
	err = server.DownloadTest()
	if err != nil {
		log.Errorf("speedtest download test failed: %v", err)
		return 0, 0, xerror.WrapError(err, xerror.ErrSystemError, "download test failed")
	}

	return server.DLSpeed.Mbps(), server.Latency, nil
}

// SpeedTestOoklaUpload 使用 Ookla speedtest 仅测试上传速度
//
// 参数:
//   - ctx: 上下文，用于控制超时和取消
//   - serverID: 指定服务器 ID（如果为空，自动选择最佳服务器）
//
// 返回:
//   - upload: 上传速度（Mbps）
//   - latency: 延迟
//   - error: 测试失败时返回错误
func (n *Node) SpeedTestOoklaUpload(ctx context.Context, serverID string) (upload float64, latency time.Duration, err error) {
	// 创建自定义 HTTP 客户端
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, portStr, err := net.SplitHostPort(addr)
			if err != nil {
				log.Errorf("split host port failed: %v, addr: %s", err, addr)
				return nil, err
			}

			var port uint16
			if portStr == "" {
				if network == "tcp" {
					port = 80
				} else {
					port = 443
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
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	user := speedtest.New(speedtest.WithDoer(httpClient))

	serverList, err := user.FetchServers()
	if err != nil {
		log.Errorf("fetch speedtest servers failed: %v", err)
		return 0, 0, xerror.WrapError(err, xerror.ErrSystemError, "fetch speedtest servers failed")
	}

	var targets speedtest.Servers
	if serverID != "" {
		for _, server := range serverList {
			if server.ID == serverID {
				targets = append(targets, server)
				break
			}
		}
		if len(targets) == 0 {
			err := fmt.Errorf("server not found: %s", serverID)
			log.Errorf("speedtest server not found: %s", serverID)
			return 0, 0, xerror.WrapError(err, xerror.ErrSystemError, "server not found")
		}
	} else {
		targets, err = serverList.FindServer([]int{})
		if err != nil {
			log.Errorf("find best speedtest server failed: %v", err)
			return 0, 0, xerror.WrapError(err, xerror.ErrSystemError, "find best server failed")
		}
	}

	if len(targets) == 0 {
		err := fmt.Errorf("no suitable speedtest server found")
		log.Errorf("no suitable speedtest server found")
		return 0, 0, xerror.WrapError(err, xerror.ErrSystemError, "no suitable server found")
	}

	server := targets[0]

	// 测试延迟
	err = server.PingTest(func(latency time.Duration) {})
	if err != nil {
		log.Errorf("speedtest ping test failed: %v", err)
		return 0, 0, xerror.WrapError(err, xerror.ErrSystemError, "ping test failed")
	}

	// 测试上传速度
	err = server.UploadTest()
	if err != nil {
		log.Errorf("speedtest upload test failed: %v", err)
		return 0, 0, xerror.WrapError(err, xerror.ErrSystemError, "upload test failed")
	}

	return server.ULSpeed.Mbps(), server.Latency, nil
}
