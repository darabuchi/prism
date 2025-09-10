package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prism/core/internal/api"
	"github.com/prism/core/internal/config"
	"github.com/prism/core/internal/core"
	"github.com/prism/core/pkg/logger"
	"github.com/sirupsen/logrus"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	var (
		configPath = flag.String("config", "config/config.yaml", "配置文件路径")
		showVer    = flag.Bool("version", false, "显示版本信息")
	)
	flag.Parse()

	if *showVer {
		fmt.Printf("Prism Core %s\n", version)
		fmt.Printf("Build Time: %s\n", buildTime)
		fmt.Printf("Git Commit: %s\n", gitCommit)
		return
	}

	// 初始化配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		logrus.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	logger.Init(cfg.Log.Level, cfg.Log.Format)
	log := logger.GetLogger()

	log.WithFields(logrus.Fields{
		"version":   version,
		"buildTime": buildTime,
		"gitCommit": gitCommit,
	}).Info("Starting Prism Core")

	// 初始化代理核心
	proxyCore, err := core.NewProxyCore(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize proxy core: %v", err)
	}

	// 启动代理核心
	if err := proxyCore.Start(); err != nil {
		log.Fatalf("Failed to start proxy core: %v", err)
	}
	defer proxyCore.Stop()

	// 初始化 API 服务器
	apiServer := api.NewServer(cfg, proxyCore)

	// 启动 API 服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.API.Port),
		Handler: apiServer.Engine(),
	}

	go func() {
		log.Infof("API server starting on port %d", cfg.API.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start API server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("Server forced to shutdown: %v", err)
	}

	log.Info("Server exited")
}
