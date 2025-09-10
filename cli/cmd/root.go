package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	serverURL string
	verbose bool
	
	// 版本信息
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "prism",
	Short: "Prism 代理客户端命令行工具",
	Long: `Prism CLI 是用于管理 Prism 代理服务的命令行工具。

支持的功能:
  • 启动/停止代理服务
  • 查看系统状态和流量统计
  • 管理节点池和订阅
  • 测试节点延迟
  • 配置管理

示例:
  prism status              # 查看服务状态
  prism start               # 启动代理服务
  prism node list           # 列出可用节点
  prism config show         # 显示当前配置`,
}

// Execute 执行根命令
func Execute() error {
	return rootCmd.Execute()
}

// SetVersionInfo 设置版本信息
func SetVersionInfo(ver, build, commit string) {
	version = ver
	buildTime = build
	gitCommit = commit
}

func init() {
	cobra.OnInitialize(initConfig)

	// 全局标志
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件路径 (默认: $HOME/.prism.yaml)")
	rootCmd.PersistentFlags().StringVarP(&serverURL, "server", "s", "http://localhost:9090", "Prism Core 服务器地址")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "详细输出")

	// 绑定到 viper
	viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	// 版本命令
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "显示版本信息",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Prism CLI %s\n", version)
			fmt.Printf("Build Time: %s\n", buildTime)
			fmt.Printf("Git Commit: %s\n", gitCommit)
		},
	})
}

// initConfig 初始化配置
func initConfig() {
	if cfgFile != "" {
		// 使用指定的配置文件
		viper.SetConfigFile(cfgFile)
	} else {
		// 查找配置文件
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// 在主目录中查找配置文件
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".prism")
	}

	// 读取环境变量
	viper.AutomaticEnv()
	viper.SetEnvPrefix("PRISM")

	// 如果找到配置文件，则读取它
	if err := viper.ReadInConfig(); err == nil && verbose {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}