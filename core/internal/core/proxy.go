package core

import (
	"context"
	"fmt"
	"sync"

	"github.com/prism/core/internal/config"
	"github.com/prism/core/pkg/logger"
	"github.com/sirupsen/logrus"
)

// ProxyCore 代理核心
type ProxyCore struct {
	config *config.Config
	log    *logrus.Logger
	
	// 运行状态
	running bool
	mutex   sync.RWMutex
	
	// 上下文控制
	ctx    context.Context
	cancel context.CancelFunc
}

// NewProxyCore 创建代理核心实例
func NewProxyCore(cfg *config.Config) (*ProxyCore, error) {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &ProxyCore{
		config: cfg,
		log:    logger.GetLogger(),
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

// Start 启动代理核心
func (p *ProxyCore) Start() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	
	if p.running {
		return fmt.Errorf("proxy core is already running")
	}
	
	p.log.Info("Starting proxy core...")
	
	// TODO: 初始化 mihomo 核心
	// 这里将集成 mihomo/clash 的核心功能
	
	p.running = true
	p.log.WithFields(logrus.Fields{
		"proxy_port": p.config.Proxy.Port,
		"mode":       p.config.Proxy.Mode,
	}).Info("Proxy core started successfully")
	
	return nil
}

// Stop 停止代理核心
func (p *ProxyCore) Stop() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	
	if !p.running {
		return nil
	}
	
	p.log.Info("Stopping proxy core...")
	
	// 取消上下文
	p.cancel()
	
	// TODO: 停止 mihomo 核心
	
	p.running = false
	p.log.Info("Proxy core stopped")
	
	return nil
}

// IsRunning 检查是否运行中
func (p *ProxyCore) IsRunning() bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.running
}

// GetStatus 获取代理状态
func (p *ProxyCore) GetStatus() map[string]interface{} {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	
	return map[string]interface{}{
		"running":     p.running,
		"proxy_port":  p.config.Proxy.Port,
		"mode":        p.config.Proxy.Mode,
		"allow_lan":   p.config.Proxy.AllowLAN,
		"connections": 0, // TODO: 实现连接统计
	}
}

// SetMode 设置代理模式
func (p *ProxyCore) SetMode(mode string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	
	validModes := map[string]bool{
		"rule":   true,
		"global": true,
		"direct": true,
	}
	
	if !validModes[mode] {
		return fmt.Errorf("invalid proxy mode: %s", mode)
	}
	
	p.config.Proxy.Mode = mode
	p.log.WithField("mode", mode).Info("Proxy mode changed")
	
	// TODO: 更新 mihomo 配置
	
	return nil
}

// GetTrafficStats 获取流量统计
func (p *ProxyCore) GetTrafficStats() map[string]interface{} {
	// TODO: 实现流量统计
	return map[string]interface{}{
		"upload":   0,
		"download": 0,
		"total":    0,
	}
}