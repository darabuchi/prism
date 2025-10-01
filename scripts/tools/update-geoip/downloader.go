package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/lazygophers/log"
)

type Downloader struct {
	client    *http.Client
	cacheDays int
}

func NewDownloader(cfg *Config) *Downloader {
	client := &http.Client{
		Timeout: 5 * time.Minute,
	}

	// 配置代理
	if cfg.Proxy != "" {
		proxyURL, err := url.Parse(cfg.Proxy)
		if err == nil {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			}
			log.Infof("使用代理: %s", cfg.Proxy)
		} else {
			log.Warnf("代理地址无效: %s - %v", cfg.Proxy, err)
		}
	}

	return &Downloader{
		client:    client,
		cacheDays: cfg.CacheDays,
	}
}

func (d *Downloader) Download(urlStr, destPath string, forceUpdate bool, cacheDays int) error {
	// 如果 cacheDays 为 0，使用默认值
	if cacheDays == 0 {
		cacheDays = d.cacheDays
	}

	// 检查是否需要更新
	if !forceUpdate && !d.shouldUpdate(destPath, cacheDays) {
		log.Infof("缓存有效（%d 天），跳过下载", cacheDays)
		return nil
	}

	log.Infof("从 %s 下载...", urlStr)

	// 创建临时文件
	tmpFile := destPath + ".tmp"
	defer os.Remove(tmpFile)

	// 发起 HTTP 请求
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP 状态码: %d", resp.StatusCode)
	}

	// 创建文件
	out, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer out.Close()

	// 写入数据
	written, err := io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	// 重命名临时文件
	if err := os.Rename(tmpFile, destPath); err != nil {
		return fmt.Errorf("重命名文件失败: %w", err)
	}

	sizeMB := float64(written) / (1024 * 1024)
	log.Infof("下载完成: %.2f MB", sizeMB)

	return nil
}

func (d *Downloader) shouldUpdate(filePath string, cacheDays int) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		// 文件不存在，需要下载
		return true
	}

	// 检查文件年龄
	age := time.Since(info.ModTime())
	maxAge := time.Duration(cacheDays) * 24 * time.Hour

	return age > maxAge
}
