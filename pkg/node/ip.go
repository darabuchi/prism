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

// IPMetadata IP 元数据
type IPMetadata struct {
	Country     string  // 国家
	CountryCode string  // 国家代码
	Region      string  // 地区/省份
	RegionCode  string  // 地区代码
	City        string  // 城市
	Timezone    string  // 时区
	ISP         string  // ISP 供应商
	Org         string  // 组织
	AS          string  // AS 号（如 "AS15169"）
	ASN         int     // ASN 数字
	ASName      string  // AS 名称
	Latitude    float64 // 纬度
	Longitude   float64 // 经度
	Postal      string  // 邮编
	Colo        string  // Cloudflare colo
	Loc         string  // 位置
}

// IPLookupResult IP 查询结果
type IPLookupResult struct {
	IP       string        // IP 地址
	Source   string        // 来源服务名称
	Error    error         // 错误信息（如果查询失败）
	Delay    time.Duration // 查询延迟
	Metadata *IPMetadata   // IP 元数据
}

// IPStatistics IP 统计结果
type IPStatistics struct {
	IP         string   // IP 地址
	Count      int      // 出现次数
	Percentage float64  // 占比（百分比）
	Sources    []string // 来源列表
}

// MetadataStatistics 元数据统计结果
type MetadataStatistics struct {
	Country     map[string]*FieldStatistics // 国家统计
	Region      map[string]*FieldStatistics // 地区统计
	City        map[string]*FieldStatistics // 城市统计
	ISP         map[string]*FieldStatistics // ISP 统计
	ASN         map[int]*FieldStatistics    // ASN 统计
	Timezone    map[string]*FieldStatistics // 时区统计
}

// FieldStatistics 字段统计
type FieldStatistics struct {
	Value      interface{} // 字段值
	Count      int         // 出现次数
	Percentage float64     // 占比（百分比）
	Sources    []string    // 来源列表
}

// IPLookupService IP 查询服务接口
type IPLookupService struct {
	Name string                                                              // 服务名称
	Func func(ctx context.Context, n *Node) (string, *IPMetadata, error)    // 查询函数
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
func lookupIPFromIPify(ctx context.Context, n *Node) (string, *IPMetadata, error) {
	ip, err := n.doHTTPRequestForIP(ctx, "https://api.ipify.org")
	return ip, nil, err
}

// lookupIPFromAPIIPify 从 api.ipify.org (JSON) 查询 IP
func lookupIPFromAPIIPify(ctx context.Context, n *Node) (string, *IPMetadata, error) {
	client := n.createHTTPClientForIP(10 * time.Second)
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.ipify.org?format=json", nil)
	if err != nil {
		return "", nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", nil, err
	}

	if ip, ok := result["ip"].(string); ok {
		return ip, nil, nil
	}

	return "", nil, nil
}

// lookupIPFromIcanhazip 从 icanhazip.com 查询 IP
func lookupIPFromIcanhazip(ctx context.Context, n *Node) (string, *IPMetadata, error) {
	ip, err := n.doHTTPRequestForIP(ctx, "https://icanhazip.com")
	return ip, nil, err
}

// lookupIPFromIfconfigMe 从 ifconfig.me 查询 IP
func lookupIPFromIfconfigMe(ctx context.Context, n *Node) (string, *IPMetadata, error) {
	ip, err := n.doHTTPRequestForIP(ctx, "https://ifconfig.me")
	return ip, nil, err
}

// lookupIPFromIdentMe 从 ident.me 查询 IP
func lookupIPFromIdentMe(ctx context.Context, n *Node) (string, *IPMetadata, error) {
	ip, err := n.doHTTPRequestForIP(ctx, "https://ident.me")
	return ip, nil, err
}

// lookupIPFromAWS 从 AWS checkip 查询 IP
func lookupIPFromAWS(ctx context.Context, n *Node) (string, *IPMetadata, error) {
	ip, err := n.doHTTPRequestForIP(ctx, "https://checkip.amazonaws.com")
	return ip, nil, err
}

// lookupIPFromIPEcho 从 ipecho.net 查询 IP
func lookupIPFromIPEcho(ctx context.Context, n *Node) (string, *IPMetadata, error) {
	ip, err := n.doHTTPRequestForIP(ctx, "https://ipecho.net/plain")
	return ip, nil, err
}

// lookupIPFromMyIP 从 myip.com 查询 IP
func lookupIPFromMyIP(ctx context.Context, n *Node) (string, *IPMetadata, error) {
	client := n.createHTTPClientForIP(10 * time.Second)
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.myip.com", nil)
	if err != nil {
		return "", nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", nil, err
	}

	metadata := &IPMetadata{}

	if ip, ok := result["ip"].(string); ok {
		if country, ok := result["country"].(string); ok {
			metadata.Country = country
		}
		if cc, ok := result["cc"].(string); ok {
			metadata.CountryCode = cc
		}
		return ip, metadata, nil
	}

	return "", nil, nil
}

// lookupIPFromIPAPI 从 ip-api.com 查询 IP
func lookupIPFromIPAPI(ctx context.Context, n *Node) (string, *IPMetadata, error) {
	client := n.createHTTPClientForIP(10 * time.Second)
	req, err := http.NewRequestWithContext(ctx, "GET", "http://ip-api.com/json/", nil)
	if err != nil {
		return "", nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", nil, err
	}

	metadata := &IPMetadata{}

	if ip, ok := result["query"].(string); ok {
		// 解析国家信息
		if country, ok := result["country"].(string); ok {
			metadata.Country = country
		}
		if countryCode, ok := result["countryCode"].(string); ok {
			metadata.CountryCode = countryCode
		}

		// 解析地区信息
		if region, ok := result["regionName"].(string); ok {
			metadata.Region = region
		}
		if regionCode, ok := result["region"].(string); ok {
			metadata.RegionCode = regionCode
		}

		// 解析城市
		if city, ok := result["city"].(string); ok {
			metadata.City = city
		}

		// 解析 ISP 和组织
		if isp, ok := result["isp"].(string); ok {
			metadata.ISP = isp
		}
		if org, ok := result["org"].(string); ok {
			metadata.Org = org
		}
		if as, ok := result["as"].(string); ok {
			metadata.AS = as
			// 尝试解析 ASN 数字
			if strings.HasPrefix(as, "AS") {
				if asn := candy.ToInt(strings.TrimPrefix(as, "AS")); asn > 0 {
					metadata.ASN = asn
				}
			}
		}

		// 解析坐标
		if lat, ok := result["lat"].(float64); ok {
			metadata.Latitude = lat
		}
		if lon, ok := result["lon"].(float64); ok {
			metadata.Longitude = lon
		}

		// 解析时区
		if timezone, ok := result["timezone"].(string); ok {
			metadata.Timezone = timezone
		}

		// 解析邮编
		if zip, ok := result["zip"].(string); ok {
			metadata.Postal = zip
		}

		return ip, metadata, nil
	}

	return "", nil, nil
}

// lookupIPFromIPInfo 从 ipinfo.io 查询 IP
func lookupIPFromIPInfo(ctx context.Context, n *Node) (string, *IPMetadata, error) {
	client := n.createHTTPClientForIP(10 * time.Second)
	req, err := http.NewRequestWithContext(ctx, "GET", "https://ipinfo.io/json", nil)
	if err != nil {
		return "", nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", nil, err
	}

	metadata := &IPMetadata{}

	if ip, ok := result["ip"].(string); ok {
		// 解析城市
		if city, ok := result["city"].(string); ok {
			metadata.City = city
		}

		// 解析地区
		if region, ok := result["region"].(string); ok {
			metadata.Region = region
		}

		// 解析国家
		if country, ok := result["country"].(string); ok {
			metadata.CountryCode = country // ipinfo.io 返回的是代码
		}

		// 解析位置
		if loc, ok := result["loc"].(string); ok {
			metadata.Loc = loc
			// loc 格式为 "latitude,longitude"
			parts := strings.Split(loc, ",")
			if len(parts) == 2 {
				metadata.Latitude = candy.ToFloat64(parts[0])
				metadata.Longitude = candy.ToFloat64(parts[1])
			}
		}

		// 解析组织/ASN
		if org, ok := result["org"].(string); ok {
			metadata.Org = org
			// org 格式通常为 "AS15169 Google LLC"
			parts := strings.Fields(org)
			if len(parts) > 0 && strings.HasPrefix(parts[0], "AS") {
				metadata.AS = parts[0]
				if asn := candy.ToInt(strings.TrimPrefix(parts[0], "AS")); asn > 0 {
					metadata.ASN = asn
				}
				if len(parts) > 1 {
					metadata.ASName = strings.Join(parts[1:], " ")
				}
			}
		}

		// 解析邮编
		if postal, ok := result["postal"].(string); ok {
			metadata.Postal = postal
		}

		// 解析时区
		if timezone, ok := result["timezone"].(string); ok {
			metadata.Timezone = timezone
		}

		return ip, metadata, nil
	}

	return "", nil, nil
}

// lookupIPFromCloudflare 从 Cloudflare 查询 IP
func lookupIPFromCloudflare(ctx context.Context, n *Node) (string, *IPMetadata, error) {
	client := n.createHTTPClientForIP(10 * time.Second)
	req, err := http.NewRequestWithContext(ctx, "GET", "https://1.1.1.1/cdn-cgi/trace", nil)
	if err != nil {
		return "", nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	metadata := &IPMetadata{}
	var ip string

	// 解析 Cloudflare trace 格式: ip=x.x.x.x, loc=XX, colo=XXX
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "ip=") {
			ip = strings.TrimPrefix(line, "ip=")
		} else if strings.HasPrefix(line, "loc=") {
			metadata.CountryCode = strings.TrimPrefix(line, "loc=")
		} else if strings.HasPrefix(line, "colo=") {
			metadata.Colo = strings.TrimPrefix(line, "colo=")
		}
	}

	if ip != "" {
		return ip, metadata, nil
	}

	return "", nil, nil
}

// lookupIPFromGoogleDNS 从 Google DNS 查询 IP
func lookupIPFromGoogleDNS(ctx context.Context, n *Node) (string, *IPMetadata, error) {
	client := n.createHTTPClientForIP(10 * time.Second)
	req, err := http.NewRequestWithContext(ctx, "GET", "https://dns.google/resolve?name=o-o.myaddr.l.google.com&type=TXT", nil)
	if err != nil {
		return "", nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", nil, err
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
							return ipWithCIDR[:idx], nil, nil
						}
						return ipWithCIDR, nil, nil
					}
				}

				return data, nil, nil
			}
		}
	}

	return "", nil, nil
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
//   - *IPMetadata: IP 元数据
//   - error: 查询失败时返回错误
func (n *Node) GetIP(ctx context.Context, serviceName string) (string, *IPMetadata, error) {
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
//   - *MetadataStatistics: 元数据统计结果
//   - error: 查询失败时返回错误
func (n *Node) GetIPConcurrent(ctx context.Context) ([]IPLookupResult, []IPStatistics, *MetadataStatistics, error) {
	services := getAllIPLookupServices()
	results := make([]IPLookupResult, len(services))

	var wg sync.WaitGroup
	wg.Add(len(services))

	// 并发查询所有服务
	for i, service := range services {
		go func(index int, svc IPLookupService) {
			defer wg.Done()

			start := time.Now()
			ip, metadata, err := svc.Func(ctx, n)
			delay := time.Since(start)

			results[index] = IPLookupResult{
				IP:       ip,
				Source:   svc.Name,
				Error:    err,
				Delay:    delay,
				Metadata: metadata,
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
	ipStats := calculateIPStatistics(results)

	// 统计元数据
	metadataStats := calculateMetadataStatistics(results)

	return results, ipStats, metadataStats, nil
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

// calculateMetadataStatistics 计算元数据统计信息
func calculateMetadataStatistics(results []IPLookupResult) *MetadataStatistics {
	stats := &MetadataStatistics{
		Country:  make(map[string]*FieldStatistics),
		Region:   make(map[string]*FieldStatistics),
		City:     make(map[string]*FieldStatistics),
		ISP:      make(map[string]*FieldStatistics),
		ASN:      make(map[int]*FieldStatistics),
		Timezone: make(map[string]*FieldStatistics),
	}

	totalSuccess := 0

	// 统计每个元数据字段
	for _, result := range results {
		if result.Error != nil || result.Metadata == nil {
			continue
		}

		totalSuccess++
		metadata := result.Metadata

		// 统计国家
		if metadata.Country != "" {
			if _, exists := stats.Country[metadata.Country]; !exists {
				stats.Country[metadata.Country] = &FieldStatistics{
					Value:   metadata.Country,
					Sources: make([]string, 0),
				}
			}
			stats.Country[metadata.Country].Count++
			stats.Country[metadata.Country].Sources = append(stats.Country[metadata.Country].Sources, result.Source)
		} else if metadata.CountryCode != "" {
			// 使用国家代码作为备选
			if _, exists := stats.Country[metadata.CountryCode]; !exists {
				stats.Country[metadata.CountryCode] = &FieldStatistics{
					Value:   metadata.CountryCode,
					Sources: make([]string, 0),
				}
			}
			stats.Country[metadata.CountryCode].Count++
			stats.Country[metadata.CountryCode].Sources = append(stats.Country[metadata.CountryCode].Sources, result.Source)
		}

		// 统计地区
		if metadata.Region != "" {
			if _, exists := stats.Region[metadata.Region]; !exists {
				stats.Region[metadata.Region] = &FieldStatistics{
					Value:   metadata.Region,
					Sources: make([]string, 0),
				}
			}
			stats.Region[metadata.Region].Count++
			stats.Region[metadata.Region].Sources = append(stats.Region[metadata.Region].Sources, result.Source)
		}

		// 统计城市
		if metadata.City != "" {
			if _, exists := stats.City[metadata.City]; !exists {
				stats.City[metadata.City] = &FieldStatistics{
					Value:   metadata.City,
					Sources: make([]string, 0),
				}
			}
			stats.City[metadata.City].Count++
			stats.City[metadata.City].Sources = append(stats.City[metadata.City].Sources, result.Source)
		}

		// 统计 ISP
		if metadata.ISP != "" {
			if _, exists := stats.ISP[metadata.ISP]; !exists {
				stats.ISP[metadata.ISP] = &FieldStatistics{
					Value:   metadata.ISP,
					Sources: make([]string, 0),
				}
			}
			stats.ISP[metadata.ISP].Count++
			stats.ISP[metadata.ISP].Sources = append(stats.ISP[metadata.ISP].Sources, result.Source)
		}

		// 统计 ASN
		if metadata.ASN > 0 {
			if _, exists := stats.ASN[metadata.ASN]; !exists {
				stats.ASN[metadata.ASN] = &FieldStatistics{
					Value:   metadata.ASN,
					Sources: make([]string, 0),
				}
			}
			stats.ASN[metadata.ASN].Count++
			stats.ASN[metadata.ASN].Sources = append(stats.ASN[metadata.ASN].Sources, result.Source)
		}

		// 统计时区
		if metadata.Timezone != "" {
			if _, exists := stats.Timezone[metadata.Timezone]; !exists {
				stats.Timezone[metadata.Timezone] = &FieldStatistics{
					Value:   metadata.Timezone,
					Sources: make([]string, 0),
				}
			}
			stats.Timezone[metadata.Timezone].Count++
			stats.Timezone[metadata.Timezone].Sources = append(stats.Timezone[metadata.Timezone].Sources, result.Source)
		}
	}

	// 计算每个字段的百分比
	if totalSuccess > 0 {
		for _, stat := range stats.Country {
			stat.Percentage = float64(stat.Count) / float64(totalSuccess) * 100
		}
		for _, stat := range stats.Region {
			stat.Percentage = float64(stat.Count) / float64(totalSuccess) * 100
		}
		for _, stat := range stats.City {
			stat.Percentage = float64(stat.Count) / float64(totalSuccess) * 100
		}
		for _, stat := range stats.ISP {
			stat.Percentage = float64(stat.Count) / float64(totalSuccess) * 100
		}
		for _, stat := range stats.ASN {
			stat.Percentage = float64(stat.Count) / float64(totalSuccess) * 100
		}
		for _, stat := range stats.Timezone {
			stat.Percentage = float64(stat.Count) / float64(totalSuccess) * 100
		}
	}

	return stats
}

// GetMostCommonCountry 获取占比最大的国家
func (m *MetadataStatistics) GetMostCommonCountry() (string, *FieldStatistics) {
	if len(m.Country) == 0 {
		return "", nil
	}

	var maxKey string
	var maxStat *FieldStatistics
	maxPercentage := -1.0

	for key, stat := range m.Country {
		if stat.Percentage > maxPercentage {
			maxPercentage = stat.Percentage
			maxKey = key
			maxStat = stat
		}
	}

	return maxKey, maxStat
}

// GetLeastCommonCountry 获取占比最小的国家
func (m *MetadataStatistics) GetLeastCommonCountry() (string, *FieldStatistics) {
	if len(m.Country) == 0 {
		return "", nil
	}

	var minKey string
	var minStat *FieldStatistics
	minPercentage := 101.0

	for key, stat := range m.Country {
		if stat.Percentage < minPercentage {
			minPercentage = stat.Percentage
			minKey = key
			minStat = stat
		}
	}

	return minKey, minStat
}

// GetMostCommonRegion 获取占比最大的地区
func (m *MetadataStatistics) GetMostCommonRegion() (string, *FieldStatistics) {
	if len(m.Region) == 0 {
		return "", nil
	}

	var maxKey string
	var maxStat *FieldStatistics
	maxPercentage := -1.0

	for key, stat := range m.Region {
		if stat.Percentage > maxPercentage {
			maxPercentage = stat.Percentage
			maxKey = key
			maxStat = stat
		}
	}

	return maxKey, maxStat
}

// GetLeastCommonRegion 获取占比最小的地区
func (m *MetadataStatistics) GetLeastCommonRegion() (string, *FieldStatistics) {
	if len(m.Region) == 0 {
		return "", nil
	}

	var minKey string
	var minStat *FieldStatistics
	minPercentage := 101.0

	for key, stat := range m.Region {
		if stat.Percentage < minPercentage {
			minPercentage = stat.Percentage
			minKey = key
			minStat = stat
		}
	}

	return minKey, minStat
}

// GetMostCommonCity 获取占比最大的城市
func (m *MetadataStatistics) GetMostCommonCity() (string, *FieldStatistics) {
	if len(m.City) == 0 {
		return "", nil
	}

	var maxKey string
	var maxStat *FieldStatistics
	maxPercentage := -1.0

	for key, stat := range m.City {
		if stat.Percentage > maxPercentage {
			maxPercentage = stat.Percentage
			maxKey = key
			maxStat = stat
		}
	}

	return maxKey, maxStat
}

// GetLeastCommonCity 获取占比最小的城市
func (m *MetadataStatistics) GetLeastCommonCity() (string, *FieldStatistics) {
	if len(m.City) == 0 {
		return "", nil
	}

	var minKey string
	var minStat *FieldStatistics
	minPercentage := 101.0

	for key, stat := range m.City {
		if stat.Percentage < minPercentage {
			minPercentage = stat.Percentage
			minKey = key
			minStat = stat
		}
	}

	return minKey, minStat
}

// GetMostCommonISP 获取占比最大的 ISP
func (m *MetadataStatistics) GetMostCommonISP() (string, *FieldStatistics) {
	if len(m.ISP) == 0 {
		return "", nil
	}

	var maxKey string
	var maxStat *FieldStatistics
	maxPercentage := -1.0

	for key, stat := range m.ISP {
		if stat.Percentage > maxPercentage {
			maxPercentage = stat.Percentage
			maxKey = key
			maxStat = stat
		}
	}

	return maxKey, maxStat
}

// GetLeastCommonISP 获取占比最小的 ISP
func (m *MetadataStatistics) GetLeastCommonISP() (string, *FieldStatistics) {
	if len(m.ISP) == 0 {
		return "", nil
	}

	var minKey string
	var minStat *FieldStatistics
	minPercentage := 101.0

	for key, stat := range m.ISP {
		if stat.Percentage < minPercentage {
			minPercentage = stat.Percentage
			minKey = key
			minStat = stat
		}
	}

	return minKey, minStat
}

// GetMostCommonASN 获取占比最大的 ASN
func (m *MetadataStatistics) GetMostCommonASN() (int, *FieldStatistics) {
	if len(m.ASN) == 0 {
		return 0, nil
	}

	var maxKey int
	var maxStat *FieldStatistics
	maxPercentage := -1.0

	for key, stat := range m.ASN {
		if stat.Percentage > maxPercentage {
			maxPercentage = stat.Percentage
			maxKey = key
			maxStat = stat
		}
	}

	return maxKey, maxStat
}

// GetLeastCommonASN 获取占比最小的 ASN
func (m *MetadataStatistics) GetLeastCommonASN() (int, *FieldStatistics) {
	if len(m.ASN) == 0 {
		return 0, nil
	}

	var minKey int
	var minStat *FieldStatistics
	minPercentage := 101.0

	for key, stat := range m.ASN {
		if stat.Percentage < minPercentage {
			minPercentage = stat.Percentage
			minKey = key
			minStat = stat
		}
	}

	return minKey, minStat
}

// GetMostCommonTimezone 获取占比最大的时区
func (m *MetadataStatistics) GetMostCommonTimezone() (string, *FieldStatistics) {
	if len(m.Timezone) == 0 {
		return "", nil
	}

	var maxKey string
	var maxStat *FieldStatistics
	maxPercentage := -1.0

	for key, stat := range m.Timezone {
		if stat.Percentage > maxPercentage {
			maxPercentage = stat.Percentage
			maxKey = key
			maxStat = stat
		}
	}

	return maxKey, maxStat
}

// GetLeastCommonTimezone 获取占比最小的时区
func (m *MetadataStatistics) GetLeastCommonTimezone() (string, *FieldStatistics) {
	if len(m.Timezone) == 0 {
		return "", nil
	}

	var minKey string
	var minStat *FieldStatistics
	minPercentage := 101.0

	for key, stat := range m.Timezone {
		if stat.Percentage < minPercentage {
			minPercentage = stat.Percentage
			minKey = key
			minStat = stat
		}
	}

	return minKey, minStat
}
