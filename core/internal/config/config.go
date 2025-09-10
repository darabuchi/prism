package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config 主配置结构
type Config struct {
	API   APIConfig   `mapstructure:"api" yaml:"api"`
	Proxy ProxyConfig `mapstructure:"proxy" yaml:"proxy"`
	Log   LogConfig   `mapstructure:"log" yaml:"log"`
	DB    DBConfig    `mapstructure:"db" yaml:"db"`
}

// APIConfig API 配置
type APIConfig struct {
	Port         int    `mapstructure:"port" yaml:"port"`
	Host         string `mapstructure:"host" yaml:"host"`
	AllowOrigins []string `mapstructure:"allow_origins" yaml:"allow_origins"`
	Secret       string `mapstructure:"secret" yaml:"secret"`
}

// ProxyConfig 代理配置
type ProxyConfig struct {
	Port      int    `mapstructure:"port" yaml:"port"`
	Host      string `mapstructure:"host" yaml:"host"`
	AllowLAN  bool   `mapstructure:"allow_lan" yaml:"allow_lan"`
	Mode      string `mapstructure:"mode" yaml:"mode"`
	LogLevel  string `mapstructure:"log_level" yaml:"log_level"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `mapstructure:"level" yaml:"level"`
	Format string `mapstructure:"format" yaml:"format"`
	Output string `mapstructure:"output" yaml:"output"`
}

// DBConfig 数据库配置
type DBConfig struct {
	Driver string `mapstructure:"driver" yaml:"driver"`
	DSN    string `mapstructure:"dsn" yaml:"dsn"`
}

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
	// 设置默认值
	setDefaults()

	// 设置配置文件路径
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		// 查找配置文件
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
		viper.AddConfigPath("$HOME/.prism")
		viper.AddConfigPath("/etc/prism")
	}

	// 读取环境变量
	viper.AutomaticEnv()
	viper.SetEnvPrefix("PRISM")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// 创建必要的目录
	if err := createDirectories(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// setDefaults 设置默认配置
func setDefaults() {
	// API 默认配置
	viper.SetDefault("api.port", 9090)
	viper.SetDefault("api.host", "127.0.0.1")
	viper.SetDefault("api.allow_origins", []string{"*"})
	viper.SetDefault("api.secret", "prism-secret-key")

	// Proxy 默认配置
	viper.SetDefault("proxy.port", 7890)
	viper.SetDefault("proxy.host", "127.0.0.1")
	viper.SetDefault("proxy.allow_lan", false)
	viper.SetDefault("proxy.mode", "rule")
	viper.SetDefault("proxy.log_level", "info")

	// Log 默认配置
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "text")
	viper.SetDefault("log.output", "stdout")

	// DB 默认配置
	viper.SetDefault("db.driver", "sqlite")
	viper.SetDefault("db.dsn", "data/prism.db")
}

// createDirectories 创建必要的目录
func createDirectories(config *Config) error {
	// 创建数据库目录
	if config.DB.Driver == "sqlite" {
		dbDir := filepath.Dir(config.DB.DSN)
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return err
		}
	}

	return nil
}