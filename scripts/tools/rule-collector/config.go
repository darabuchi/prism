package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Config 规则收集器配置
type Config struct {
	// 缓存目录（默认使用系统临时目录）
	CacheDir string

	// 输出目录
	OutputDir string

	// 缓存天数（超过此天数的缓存将被更新）
	CacheDays int

	// HTTP 代理地址（支持 http:// 和 socks5://）
	Proxy string

	// 是否强制更新缓存
	ForceUpdate bool

	// 是否详细输出
	Verbose bool
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	// 设置默认值
	cfg := &Config{
		CacheDir:    "/tmp/prism_rule_collector",
		OutputDir:   "./resource/rules",  // 默认输出到 resource/rules
		CacheDays:   1, // 规则更新频繁，默认 1 天
		Proxy:       "",
		ForceUpdate: false,
		Verbose:     false,
	}

	// 从 .env 文件加载配置（如果存在）
	if err := loadDotEnv(cfg); err != nil {
		if !os.IsNotExist(err) {
			// 如果不是文件不存在错误，记录警告
			// 这里不是致命错误，使用默认配置继续
		}
	}

	// 确保输出目录是绝对路径
	if !filepath.IsAbs(cfg.OutputDir) {
		absPath, err := filepath.Abs(cfg.OutputDir)
		if err != nil {
			return nil, err
		}
		cfg.OutputDir = absPath
	}

	return cfg, nil
}

// loadDotEnv 从 .env 文件加载配置
func loadDotEnv(cfg *Config) error {
	// 尝试从多个位置加载 .env 文件
	envPaths := []string{
		".env",                     // 当前目录
		"../../../.env",            // 项目根目录（从 scripts/tools/rule-collector）
	}

	var file *os.File
	var err error
	for _, path := range envPaths {
		file, err = os.Open(path)
		if err == nil {
			defer file.Close()
			break
		}
	}

	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过注释和空行
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析 KEY=VALUE
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)

		switch key {
		case "RULE_CACHE_DIR":
			if value != "" {
				cfg.CacheDir = value
			}
		case "RULE_OUTPUT_DIR":
			if value != "" {
				cfg.OutputDir = value
			}
		case "RULE_CACHE_DAYS":
			if days, err := strconv.Atoi(value); err == nil && days > 0 {
				cfg.CacheDays = days
			}
		case "RULE_PROXY", "HTTP_PROXY", "ALL_PROXY":
			if value != "" {
				cfg.Proxy = value
			}
		}
	}

	return scanner.Err()
}
