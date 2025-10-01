package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	CacheDir    string
	DataDir     string
	CacheDays   int
	Proxy       string
	IPInfoToken string
}

func LoadConfig() (*Config, error) {
	// 设置默认值
	cfg := &Config{
		CacheDir:    "/tmp/prism_geoip",
		DataDir:     "./resource/geoip",
		CacheDays:   7,
		Proxy:       "",
		IPInfoToken: "", 
	}

	// 优先从 .env 文件加载配置
	if err := loadDotEnv(cfg); err != nil {
		// .env 文件不存在或读取失败不是致命错误，使用默认配置
		if !os.IsNotExist(err) {
			// 如果不是文件不存在错误，记录警告
			// log.Warnf(".env 文件读取失败，使用默认配置: %v", err)
		}
	}

	// 确保目录路径是绝对路径
	if !filepath.IsAbs(cfg.DataDir) {
		absPath, err := filepath.Abs(cfg.DataDir)
		if err != nil {
			return nil, err
		}
		cfg.DataDir = absPath
	}

	return cfg, nil
}

func loadDotEnv(cfg *Config) error {
	file, err := os.Open(".env")
	if err != nil {
		return err
	}
	defer file.Close()

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
		case "GEOIP_CACHE_DIR":
			if value != "" {
				cfg.CacheDir = value
			}
		case "GEOIP_DATA_DIR":
			if value != "" {
				cfg.DataDir = value
			}
		case "GEOIP_CACHE_DAYS":
			if days, err := strconv.Atoi(value); err == nil && days > 0 {
				cfg.CacheDays = days
			}
		case "GEOIP_PROXY", "HTTP_PROXY", "ALL_PROXY":
			if value != "" {
				cfg.Proxy = value
			}
		case "IPINFO_TOKEN":
			if value != "" {
				cfg.IPInfoToken = value
			}
		}
	}

	return scanner.Err()
}

