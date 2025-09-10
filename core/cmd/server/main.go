package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/prism/core/internal/api"
	"github.com/prism/core/internal/service"
	"github.com/prism/core/internal/storage"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig           `mapstructure:"server"`
	Database storage.DatabaseConfig `mapstructure:"database"`
	Log      LogConfig              `mapstructure:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

var (
	configFile = flag.String("config", "config.yaml", "配置文件路径")
	version    = flag.Bool("version", false, "显示版本信息")
)

const (
	AppName    = "prism-core"
	AppVersion = "1.0.0"
)

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("%s version %s\n", AppName, AppVersion)
		os.Exit(0)
	}

	// 初始化配置
	config, err := loadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	logger := initLogger(config.Log)
	logger.Infof("Starting %s v%s", AppName, AppVersion)

	// 初始化数据库
	db, err := storage.NewDatabase(&config.Database)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 数据库迁移
	if err := db.AutoMigrate(); err != nil {
		logger.Fatalf("Failed to migrate database: %v", err)
	}
	logger.Info("Database migration completed")

	// 初始化服务
	subscriptionSvc := service.NewSubscriptionService(db)
	nodePoolSvc := service.NewNodePoolService(db)
	nodeSvc := service.NewNodeService(db)
	statsSvc := service.NewStatsService(db)
	schedulerSvc := service.NewSchedulerService(db, subscriptionSvc, nodeSvc)

	// 启动调度器
	if err := schedulerSvc.Start(); err != nil {
		logger.Fatalf("Failed to start scheduler: %v", err)
	}
	logger.Info("Scheduler started")

	// 初始化路由
	router := api.NewRouter(
		logger,
		subscriptionSvc,
		nodePoolSvc,
		nodeSvc,
		statsSvc,
		schedulerSvc,
	)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	logger.Infof("Starting server on %s", addr)

	// 优雅关闭
	go func() {
		if err := router.Run(addr); err != nil {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// 停止调度器
	if err := schedulerSvc.Stop(); err != nil {
		logger.Errorf("Failed to stop scheduler: %v", err)
	}

	logger.Info("Server stopped")
}

// loadConfig 加载配置
func loadConfig(configFile string) (*Config, error) {
	viper.SetConfigFile(configFile)
	viper.AutomaticEnv()

	// 设置默认值
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 9090)
	viper.SetDefault("database.driver", "sqlite")
	viper.SetDefault("database.dsn", "data/prism.db")
	viper.SetDefault("database.debug", false)
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "text")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
		log.Printf("Config file not found, using defaults")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// initLogger 初始化日志
func initLogger(config LogConfig) *logrus.Logger {
	logger := logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// 设置日志格式
	if config.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.RFC3339,
			FullTimestamp:   true,
		})
	}

	return logger
}
