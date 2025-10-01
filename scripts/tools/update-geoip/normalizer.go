package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ISPMappings 定义 ISP/Org 名称标准化映射
var ISPMappings = map[string]string{
	// Google
	"Google LLC":  "Google",
	"Google Inc.": "Google",
	"GOOGLE":      "Google",

	// Amazon
	"Amazon.com, Inc.":           "Amazon",
	"Amazon Technologies Inc.":   "Amazon",
	"AMAZON-02":                  "Amazon",
	"Amazon Data Services":       "Amazon",
	"Amazon.com":                 "Amazon",
	"AMAZON-AES":                 "Amazon",
	"Amazon Data Services NoVa":  "Amazon",
	"Amazon Data Services Ireland Limited": "Amazon",

	// Cloudflare
	"Cloudflare, Inc.": "Cloudflare",
	"CLOUDFLARENET":    "Cloudflare",

	// Alibaba
	"Alibaba (US) Technology Co., Ltd.": "Alibaba",
	"Alibaba Cloud LLC":                 "Alibaba Cloud",
	"Alibaba Cloud Computing":           "Alibaba Cloud",
	"Hangzhou Alibaba Advertising":      "Alibaba",
	"Alibaba":                           "Alibaba",
	"ALICLOUD":                          "Alibaba Cloud",
	"Aliyun Computing":                  "Alibaba Cloud",

	// Tencent
	"Tencent Cloud Computing":              "Tencent Cloud",
	"TENCENT-NET-AP-CN":                    "Tencent",
	"Shenzhen Tencent Computer Systems":    "Tencent",
	"Tencent":                              "Tencent",
	"Tencent Holdings Limited":             "Tencent",
	"Shenzhen Tencent Computer System":     "Tencent",
	"Tencent Building, Kejizhongyi Avenue": "Tencent",

	// China Telecom
	"China Telecom":    "中国电信",
	"CHINANET":         "中国电信",
	"CHINA TELECOM":    "中国电信",
	"CHINANET-BACKBONE": "中国电信",

	// China Unicom
	"China Unicom":   "中国联通",
	"CHINA169":       "中国联通",
	"CHINA UNICOM":   "中国联通",
	"CNC Group CHINA169 Network": "中国联通",

	// China Mobile
	"China Mobile":   "中国移动",
	"CMNET":          "中国移动",
	"CHINA MOBILE":   "中国移动",

	// Microsoft
	"Microsoft Corporation":       "Microsoft",
	"MICROSOFT-CORP-MSN-AS-BLOCK": "Microsoft",
	"Microsoft":                   "Microsoft",

	// DigitalOcean
	"DigitalOcean, LLC":  "DigitalOcean",
	"DIGITALOCEAN-ASN":   "DigitalOcean",
	"DigitalOcean":       "DigitalOcean",

	// Vultr
	"Vultr Holdings, LLC": "Vultr",
	"AS-CHOOPA":           "Vultr",
	"Vultr":               "Vultr",
	"The Constant Company": "Vultr",

	// Linode
	"Linode, LLC": "Linode",
	"LINODE-AP":   "Linode",
	"Linode":      "Linode",

	// Oracle
	"Oracle Corporation": "Oracle",
	"Oracle":             "Oracle",

	// IBM
	"IBM Cloud": "IBM",
	"IBM":       "IBM",

	// OVH
	"OVH SAS":       "OVH",
	"OVH":           "OVH",

	// Hetzner
	"Hetzner Online GmbH": "Hetzner",
	"Hetzner":             "Hetzner",

	// Huawei
	"Huawei Cloud Service":       "Huawei Cloud",
	"Huawei":                     "Huawei",
	"HUAWEI":                     "Huawei",
}

// NormalizeISPName 标准化 ISP/Org 名称
func NormalizeISPName(name string) string {
	if normalized, ok := ISPMappings[name]; ok {
		return normalized
	}
	return name
}

// SaveISPMappings 保存 ISP 映射到文件
func SaveISPMappings(dataDir string) error {
	filePath := filepath.Join(dataDir, "isp_mappings.txt")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "# ISP/Org Name Standardization Mappings\n")
	fmt.Fprintf(file, "# Format: original_name=standard_name\n")
	fmt.Fprintf(file, "# Generated at: %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(file, "\n")

	for original, standard := range ISPMappings {
		fmt.Fprintf(file, "%s=%s\n", original, standard)
	}

	return nil
}

// Metadata 元数据结构
type Metadata struct {
	UpdatedAt  string            `json:"updated_at"`
	CacheDir   string            `json:"cache_dir"`
	DataDir    string            `json:"data_dir"`
	CacheDays  int               `json:"cache_duration_days"`
	Sources    map[string]Source `json:"sources"`
	Statistics struct {
		TotalSources int `json:"total_sources"`
		Downloaded   int `json:"downloaded"`
		Failed       int `json:"failed"`
	} `json:"statistics"`
}

// SaveMetadata 保存元数据到文件
func SaveMetadata(cfg *Config, sources map[string]Source, successCount, failCount int) error {
	metadata := Metadata{
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		CacheDir:  cfg.CacheDir,
		DataDir:   cfg.DataDir,
		CacheDays: cfg.CacheDays,
		Sources:   sources,
	}
	metadata.Statistics.TotalSources = len(sources)
	metadata.Statistics.Downloaded = successCount
	metadata.Statistics.Failed = failCount

	filePath := filepath.Join(cfg.DataDir, "metadata.json")
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}
