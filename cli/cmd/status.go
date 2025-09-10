package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/prism/cli/internal/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// statusCmd 状态命令
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "查看 Prism 服务状态",
	Long:  `显示 Prism Core 服务的运行状态、系统信息和代理统计。`,
	RunE:  runStatus,
}

// StatusResponse 状态响应
type StatusResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Version       string  `json:"version"`
		Uptime        int64   `json:"uptime"`
		ProxyStatus   string  `json:"proxy_status"`
		MemoryUsage   float64 `json:"memory_usage"`
		CPUUsage      float64 `json:"cpu_usage"`
		TotalUpload   int64   `json:"total_upload"`
		TotalDownload int64   `json:"total_download"`
	} `json:"data"`
}

func runStatus(cmd *cobra.Command, args []string) error {
	serverURL := viper.GetString("server")
	verbose := viper.GetBool("verbose")

	if verbose {
		fmt.Printf("Connecting to: %s\n", serverURL)
	}

	// 创建 HTTP 客户端
	httpClient := &http.Client{Timeout: 10 * time.Second}
	apiClient := client.NewClient(serverURL, httpClient)

	// 获取系统状态
	resp, err := apiClient.Get("/api/v1/system/status")
	if err != nil {
		return fmt.Errorf("failed to connect to Prism Core: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %s", resp.Status)
	}

	var statusResp StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	if statusResp.Code != 0 {
		return fmt.Errorf("server error: %s", statusResp.Message)
	}

	// 显示状态信息
	printStatus(&statusResp.Data)

	return nil
}

func printStatus(data *struct {
	Version       string  `json:"version"`
	Uptime        int64   `json:"uptime"`
	ProxyStatus   string  `json:"proxy_status"`
	MemoryUsage   float64 `json:"memory_usage"`
	CPUUsage      float64 `json:"cpu_usage"`
	TotalUpload   int64   `json:"total_upload"`
	TotalDownload int64   `json:"total_download"`
}) {
	// 标题
	color.New(color.FgCyan, color.Bold).Println("Prism 服务状态")
	fmt.Println(color.New(color.FgBlue).Sprint("==================="))
	fmt.Println()

	// 基本信息
	table := tablewriter.NewWriter(cmd.OutOrStdout())
	table.SetHeader([]string{"项目", "值"})
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	// 状态颜色
	var statusColor *color.Color
	var statusText string
	switch data.ProxyStatus {
	case "running":
		statusColor = color.New(color.FgGreen)
		statusText = "运行中"
	case "stopped":
		statusColor = color.New(color.FgRed)
		statusText = "已停止"
	default:
		statusColor = color.New(color.FgYellow)
		statusText = "未知"
	}

	// 运行时间格式化
	uptime := time.Duration(data.Uptime) * time.Second
	uptimeStr := formatDuration(uptime)

	// 流量格式化
	uploadStr := formatBytes(data.TotalUpload)
	downloadStr := formatBytes(data.TotalDownload)

	// 添加行
	table.Append([]string{"版本", data.Version})
	table.Append([]string{"状态", statusColor.Sprint(statusText)})
	table.Append([]string{"运行时间", uptimeStr})
	table.Append([]string{"内存使用", fmt.Sprintf("%.1f MB", data.MemoryUsage)})
	table.Append([]string{"CPU 使用率", fmt.Sprintf("%.1f%%", data.CPUUsage)})
	table.Append([]string{"总上传", uploadStr})
	table.Append([]string{"总下载", downloadStr})

	table.Render()
}

// formatDuration 格式化时间间隔
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0f秒", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%.0f分钟", d.Minutes())
	} else if d < 24*time.Hour {
		return fmt.Sprintf("%.1f小时", d.Hours())
	} else {
		days := int(d.Hours() / 24)
		hours := int(d.Hours()) % 24
		return fmt.Sprintf("%d天%d小时", days, hours)
	}
}

// formatBytes 格式化字节数
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func init() {
	rootCmd.AddCommand(statusCmd)
}