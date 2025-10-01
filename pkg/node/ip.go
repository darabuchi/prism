package node

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/candy"
	"github.com/metacubex/mihomo/constant"
)

// IPLookupResult IP 查询结果
type IPLookupResult struct {
	IP     string        // IP 地址
	Source string        // 来源服务名称
	Error  error         // 错误信息（如果查询失败）
	Delay  time.Duration // 查询延迟
}

// IPStatistics IP 统计结果
type IPStatistics struct {
	IP         string  // IP 地址
	Count      int     // 出现次数
	Percentage float64 // 占比（百分比）
	Sources    []string // 来源列表
}

// IPLookupService IP 查询服务接口
type IPLookupService struct {
	Name string                                                     // 服务名称
	Func func(ctx context.Context, n *Node) (string, error)        // 查询函数
}

// getAllIPLookupServices 获取所有 IP 查询服务
func getAllIPLookupServices() []IPLookupService {
	return []IPLookupService{
		{"ipify", lookupIPFromIPify},
		{"icanhazip", lookupIPFromIcanhazip},
		{"ifconfig.me", lookupIPFromIfconfigMe},
		{"ident.me", lookupIPFromIdentMe},
		{"api.ipify", lookupIPFromAPIIPify},
		{"checkip.amazonaws", lookupIPFromAWS},
		{"ipecho.net", lookupIPFromIPEcho},
		{"myip.com", lookupIPFromMyIP},
		{"ip-api.com", lookupIPFromIPAPI},
		{"ipinfo.io", lookupIPFromIPInfo},
		{"cloudflare", lookupIPFromCloudflare},
		{"google-dns", lookupIPFromGoogleDNS},
	}
}

// createHTTPClientForIP 创建用于 IP 查询的 HTTP 客户端
func (n *Node) createHTTPClientForIP(timeout time.Duration) *http.Client {
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
	}
}

// doHTTPRequestForIP 执行 HTTP 请求获取 IP
func (n *Node) doHTTPRequestForIP(ctx context.Context, url string) (string, error) {
	client := n.createHTTPClientForIP(10 * time.Second)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "curl/7.68.0")
	req.Header.Set("Accept", "*/*")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(body)), nil
}

// lookupIPFromIPify 从 ipify.org 查询 IP
func lookupIPFromIPify(ctx context.Context, n *Node) (string, error) {
	return n.doHTTPRequestForIP(ctx, "https://api.ipify.org")
}

// lookupIPFromAPIIPify 从 api.ipify.org (JSON) 查询 IP
func lookupIPFromAPIIPify(ctx context.Context, n *Node) (string, error) {
	client := n.createHTTPClientForIP(10 * time.Second)
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.ipify.org?format=json", nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if ip, ok := result["ip"].(string); ok {
		return ip, nil
	}

	return "", nil
}

// lookupIPFromIcanhazip 从 icanhazip.com 查询 IP
func lookupIPFromIcanhazip(ctx context.Context, n *Node) (string, error) {
	return n.doHTTPRequestForIP(ctx, "https://icanhazip.com")
}

// lookupIPFromIfconfigMe 从 ifconfig.me 查询 IP
func lookupIPFromIfconfigMe(ctx context.Context, n *Node) (string, error) {
	return n.doHTTPRequestForIP(ctx, "https://ifconfig.me")
}

// lookupIPFromIdentMe 从 ident.me 查询 IP
func lookupIPFromIdentMe(ctx context.Context, n *Node) (string, error) {
	return n.doHTTPRequestForIP(ctx, "https://ident.me")
}

// lookupIPFromAWS 从 AWS checkip 查询 IP
func lookupIPFromAWS(ctx context.Context, n *Node) (string, error) {
	return n.doHTTPRequestForIP(ctx, "https://checkip.amazonaws.com")
}

// lookupIPFromIPEcho 从 ipecho.net 查询 IP
func lookupIPFromIPEcho(ctx context.Context, n *Node) (string, error) {
	return n.doHTTPRequestForIP(ctx, "https://ipecho.net/plain")
}

// lookupIPFromMyIP 从 myip.com 查询 IP
func lookupIPFromMyIP(ctx context.Context, n *Node) (string, error) {
	client := n.createHTTPClientForIP(10 * time.Second)
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.myip.com", nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if ip, ok := result["ip"].(string); ok {
		return ip, nil
	}

	return "", nil
}

// lookupIPFromIPAPI 从 ip-api.com 查询 IP
func lookupIPFromIPAPI(ctx context.Context, n *Node) (string, error) {
	client := n.createHTTPClientForIP(10 * time.Second)
	req, err := http.NewRequestWithContext(ctx, "GET", "http://ip-api.com/json/", nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if ip, ok := result["query"].(string); ok {
		return ip, nil
	}

	return "", nil
}

// lookupIPFromIPInfo 从 ipinfo.io 查询 IP
func lookupIPFromIPInfo(ctx context.Context, n *Node) (string, error) {
	client := n.createHTTPClientForIP(10 * time.Second)
	req, err := http.NewRequestWithContext(ctx, "GET", "https://ipinfo.io/json", nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if ip, ok := result["ip"].(string); ok {
		return ip, nil
	}

	return "", nil
}

// lookupIPFromCloudflare 从 Cloudflare 查询 IP
func lookupIPFromCloudflare(ctx context.Context, n *Node) (string, error) {
	client := n.createHTTPClientForIP(10 * time.Second)
	req, err := http.NewRequestWithContext(ctx, "GET", "https://1.1.1.1/cdn-cgi/trace", nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析 Cloudflare trace 格式: ip=x.x.x.x
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ip=") {
			return strings.TrimPrefix(line, "ip="), nil
		}
	}

	return "", nil
}

// lookupIPFromGoogleDNS 从 Google DNS 查询 IP
func lookupIPFromGoogleDNS(ctx context.Context, n *Node) (string, error) {
	client := n.createHTTPClientForIP(10 * time.Second)
	req, err := http.NewRequestWithContext(ctx, "GET", "https://dns.google/resolve?name=o-o.myaddr.l.google.com&type=TXT", nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	// 从 DNS TXT 记录中提取 IP
	if answer, ok := result["Answer"].([]interface{}); ok && len(answer) > 0 {
		if firstAnswer, ok := answer[0].(map[string]interface{}); ok {
			if data, ok := firstAnswer["data"].(string); ok {
				// TXT 记录格式: "x.x.x.x" 或 "edns0-client-subnet x.x.x.x/xx"
				data = strings.Trim(data, "\"")

				// 处理 edns0-client-subnet 格式
				if strings.Contains(data, "edns0-client-subnet") {
					// 提取 IP 部分
					parts := strings.Fields(data)
					if len(parts) >= 2 {
						// 移除 CIDR 后缀 /24
						ipWithCIDR := parts[1]
						if idx := strings.Index(ipWithCIDR, "/"); idx != -1 {
							return ipWithCIDR[:idx], nil
						}
						return ipWithCIDR, nil
					}
				}

				return data, nil
			}
		}
	}

	return "", nil
}

// GetIP 获取节点的出口 IP 地址（单次查询）
//
// 使用指定的服务查询 IP 地址
//
// 参数:
//   - ctx: 上下文，用于控制超时和取消
//   - serviceName: 服务名称（如果为空，使用默认服务 ipify）
//
// 返回:
//   - string: IP 地址
//   - error: 查询失败时返回错误
func (n *Node) GetIP(ctx context.Context, serviceName string) (string, error) {
	services := getAllIPLookupServices()

	if serviceName == "" {
		serviceName = "ipify"
	}

	for _, service := range services {
		if service.Name == serviceName {
			return service.Func(ctx, n)
		}
	}

	// 如果未找到指定服务，使用默认服务
	return lookupIPFromIPify(ctx, n)
}

// GetIPConcurrent 并发获取节点的出口 IP 地址
//
// 使用多个服务并发查询 IP 地址，并统计结果
//
// 参数:
//   - ctx: 上下文，用于控制超时和取消
//
// 返回:
//   - []IPLookupResult: 所有查询结果
//   - []IPStatistics: IP 统计结果（按出现次数降序排列）
//   - error: 查询失败时返回错误
func (n *Node) GetIPConcurrent(ctx context.Context) ([]IPLookupResult, []IPStatistics, error) {
	services := getAllIPLookupServices()
	results := make([]IPLookupResult, len(services))

	var wg sync.WaitGroup
	wg.Add(len(services))

	// 并发查询所有服务
	for i, service := range services {
		go func(index int, svc IPLookupService) {
			defer wg.Done()

			start := time.Now()
			ip, err := svc.Func(ctx, n)
			delay := time.Since(start)

			results[index] = IPLookupResult{
				IP:     ip,
				Source: svc.Name,
				Error:  err,
				Delay:  delay,
			}

			if err != nil {
				log.Warnf("IP lookup from %s failed: %v", svc.Name, err)
			} else {
				log.Infof("IP lookup from %s succeeded: %s (delay: %v)", svc.Name, ip, delay)
			}
		}(i, service)
	}

	// 等待所有查询完成
	wg.Wait()

	// 统计 IP 出现次数
	statistics := calculateIPStatistics(results)

	return results, statistics, nil
}

// calculateIPStatistics 计算 IP 统计信息
func calculateIPStatistics(results []IPLookupResult) []IPStatistics {
	ipCount := make(map[string]*IPStatistics)
	totalSuccess := 0

	// 统计每个 IP 的出现次数和来源
	for _, result := range results {
		if result.Error != nil || result.IP == "" {
			continue
		}

		totalSuccess++

		if _, exists := ipCount[result.IP]; !exists {
			ipCount[result.IP] = &IPStatistics{
				IP:      result.IP,
				Count:   0,
				Sources: make([]string, 0),
			}
		}

		ipCount[result.IP].Count++
		ipCount[result.IP].Sources = append(ipCount[result.IP].Sources, result.Source)
	}

	// 计算百分比
	statistics := make([]IPStatistics, 0, len(ipCount))
	for _, stat := range ipCount {
		if totalSuccess > 0 {
			stat.Percentage = float64(stat.Count) / float64(totalSuccess) * 100
		}
		statistics = append(statistics, *stat)
	}

	// 按出现次数降序排序
	for i := 0; i < len(statistics); i++ {
		for j := i + 1; j < len(statistics); j++ {
			if statistics[j].Count > statistics[i].Count {
				statistics[i], statistics[j] = statistics[j], statistics[i]
			}
		}
	}

	return statistics
}
