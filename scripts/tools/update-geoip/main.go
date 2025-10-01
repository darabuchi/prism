package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/darabuchi/prism"
	"github.com/darabuchi/prism/pkg/geoip"
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/xerror"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	forceUpdate bool
	verbose     bool
)

// DownloadedSource 已下载的数据源信息
type DownloadedSource struct {
	Name   string
	Source Source
	Path   string
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "update-geoip",
		Short: "更新 GeoIP 数据库",
		Long:  "从可信的 GitHub 数据源下载并更新 GeoIP 数据库，支持缓存和代理配置",
		Run:   runUpdate,
	}

	rootCmd.Flags().BoolVarP(&forceUpdate, "force", "f", false, "强制更新，忽略缓存")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "详细输出")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runUpdate(cmd *cobra.Command, args []string) {
	// 加载配置
	cfg, err := LoadConfig()
	if err != nil {
		log.Errorf("加载配置失败: %v", err)
		os.Exit(1)
	}

	// 显示环境信息
	showEnvironmentInfo(cfg)

	// 创建必要的目录
	if err := os.MkdirAll(cfg.CacheDir, 0755); err != nil {
		log.Errorf("创建缓存目录失败: %v", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
		log.Errorf("创建数据目录失败: %v", err)
		os.Exit(1)
	}

	log.Infof("开始更新 GeoIP 数据库...")

	// 创建下载器
	downloader := NewDownloader(cfg)

	// 下载所有数据源
	sources := GetSources(cfg)
	var successCount, failCount int
	downloadedSources := make(map[string]DownloadedSource) // name -> downloaded source info

	for name, source := range sources {
		cacheDays := source.CacheDays
		if cacheDays == 0 {
			cacheDays = cfg.CacheDays
		}
		log.Infof("处理数据源: %s (优先级: %d, 缓存: %d 天)", name, source.Priority, cacheDays)

		// 从 URL 中提取正确的文件扩展名（忽略 query 参数）
		cacheFile := filepath.Join(cfg.CacheDir, name+getFileExtension(source.URL))
		if err := downloader.Download(source.URL, cacheFile, forceUpdate, cacheDays); err != nil {
			log.Errorf("下载失败: %s - %v", name, err)
			failCount++
			continue
		}

		downloadedSources[name] = DownloadedSource{
			Source: source,
			Path:   cacheFile,
		}
		successCount++
	}

	log.Infof("下载统计: 成功 %d, 失败 %d", successCount, failCount)

	// 按优先级排序数据源
	sortedSources := sortSourcesByPriority(downloadedSources)
	log.Infof("将按优先级顺序处理 %d 个数据源", len(sortedSources))

	// 转换并合并为 Prism GeoIP 数据库
	log.Infof("开始转换数据库格式...")
	converter, err := NewConverter()
	if err != nil {
		log.Errorf("创建转换器失败: %v", err)
		os.Exit(1)
	}

	// 按优先级处理所有数据源
	for _, ds := range sortedSources {
		log.Infof("转换: %s (类型: %s, 格式: %s)", ds.Name, ds.Source.Type, ds.Source.Format)

		switch ds.Source.Type {
		case SourceTypeMMDB:
			// MaxMind MMDB 格式
			if err := converter.ConvertMaxMindDB(ds.Path); err != nil {
				log.Errorf("转换 MMDB 失败: %s - %v", ds.Name, err)
				failCount++
			}

		case SourceTypeCSV:
			// CSV 格式
			if err := convertCSVSource(converter, ds); err != nil {
				log.Errorf("转换 CSV 失败: %s - %v", ds.Name, err)
				failCount++
			}

		case SourceTypeText:
			// 文本格式（CIDR 列表）
			if err := convertTextSource(converter, ds); err != nil {
				log.Errorf("转换 TXT 失败: %s - %v", ds.Name, err)
				failCount++
			}

		default:
			log.Warnf("未知的数据源类型: %s", ds.Source.Type)
			failCount++
		}
	}

	// 写入最终数据库（先写入临时文件，再替换）
	log.Infof("")
	log.Infof("=================================================")
	log.Infof("写入数据库")
	log.Infof("=================================================")
	log.Infof("")

	outputPath := filepath.Join(cfg.DataDir, "geoip.mmdb")
	tempPath := outputPath + ".tmp"
	backupPath := outputPath + ".bak"

	// 备份旧数据库
	if err := backupOldDatabase(outputPath, backupPath); err != nil {
		log.Warnf("备份旧数据库失败: %v", err)
	}

	// 写入到临时文件
	log.Infof("开始写入数据库到临时文件: %s", filepath.Base(tempPath))
	tempFile, err := os.Create(tempPath)
	if err != nil {
		log.Errorf("创建临时文件失败: %v", err)
		os.Exit(1)
	}

	written, err := converter.WriteTo(tempFile)
	tempFile.Close()

	if err != nil {
		log.Errorf("写入数据库失败: %v", err)
		os.Remove(tempPath)
		// 尝试恢复备份
		if err := restoreBackup(backupPath, outputPath); err != nil {
			log.Errorf("恢复备份失败: %v", err)
		}
		os.Exit(1)
	}

	// 替换旧文件
	log.Infof("替换数据库文件: %s -> %s", filepath.Base(tempPath), filepath.Base(outputPath))
	if err := os.Rename(tempPath, outputPath); err != nil {
		log.Errorf("替换数据库文件失败: %v", err)
		os.Remove(tempPath)
		// 尝试恢复备份
		if err := restoreBackup(backupPath, outputPath); err != nil {
			log.Errorf("恢复备份失败: %v", err)
		}
		os.Exit(1)
	}

	sizeMB := float64(written) / (1024 * 1024)
	log.Infof("")
	log.Infof("✓ 数据库已生成: %s (%.2f MB)", filepath.Base(outputPath), sizeMB)

	// 显示备份信息
	if backupInfo, err := os.Stat(backupPath); err == nil {
		backupSizeMB := float64(backupInfo.Size()) / (1024 * 1024)
		log.Infof("✓ 备份已保存: %s (%.2f MB)", filepath.Base(backupPath), backupSizeMB)
	}

	// 生成 ISP/Org 映射文件
	if err := SaveISPMappings(cfg.DataDir); err != nil {
		log.Errorf("保存 ISP 映射失败: %v", err)
	} else {
		log.Infof("ISP/Org 映射文件已生成")
	}

	// 生成元数据文件
	if err := SaveMetadata(cfg, sources, successCount, failCount); err != nil {
		log.Errorf("保存元数据失败: %v", err)
	} else {
		log.Infof("元数据文件已生成")
	}

	// 显示已安装的数据库
	log.Infof("已安装的数据库:")
	entries, err := os.ReadDir(cfg.DataDir)
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && filepath.Ext(entry.Name()) == ".mmdb" {
				info, err := entry.Info()
				if err == nil {
					size := float64(info.Size()) / (1024 * 1024)
					dbType := "MaxMind"
					if entry.Name() == "prism-geoip.mmdb" {
						dbType = "Prism GeoIP"
					}
					log.Infof("  %s - %.2f MB (%s)", entry.Name(), size, dbType)
				}
			}
		}
	}

	log.Infof("GeoIP 数据库更新完成!")
	log.Infof("")

	// 测试查询
	testGeoIPQuery(outputPath)
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		log.Errorf("读取文件失败，源文件: %s, 错误: %v", src, err)
		return xerror.WrapError(err, ErrFileOperationFailed, "read file "+filepath.Base(src)+" failed")
	}

	if err := os.WriteFile(dst, data, 0644); err != nil {
		log.Errorf("写入文件失败，目标文件: %s, 错误: %v", dst, err)
		return xerror.WrapError(err, ErrFileOperationFailed, "write file "+filepath.Base(dst)+" failed")
	}

	return nil
}

// showEnvironmentInfo 显示环境配置信息
func showEnvironmentInfo(cfg *Config) {
	log.Infof("=================================================")
	log.Infof("  Prism GeoIP 数据库更新工具")
	log.Infof("=================================================")
	log.Infof("")
	log.Infof("环境配置:")
	log.Infof("  缓存目录: %s", cfg.CacheDir)
	log.Infof("  数据目录: %s", cfg.DataDir)
	log.Infof("  缓存时长: %d 天", cfg.CacheDays)

	if cfg.Proxy != "" {
		log.Infof("  HTTP 代理: %s", cfg.Proxy)
	} else {
		log.Infof("  HTTP 代理: 未设置")
	}

	// 显示数据源统计
	sources := GetSources(cfg)
	log.Infof("")
	log.Infof("数据源统计:")
	log.Infof("  总数据源: %d 个", len(sources))

	// 统计不同类型的数据源
	mmdbCount := 0
	csvCount := 0
	txtCount := 0
	for _, source := range sources {
		switch source.Type {
		case SourceTypeMMDB:
			mmdbCount++
		case SourceTypeCSV:
			csvCount++
		case SourceTypeText:
			txtCount++
		}
	}

	if mmdbCount > 0 {
		log.Infof("  MMDB 格式: %d 个", mmdbCount)
	}
	if csvCount > 0 {
		log.Infof("  CSV 格式: %d 个", csvCount)
	}
	if txtCount > 0 {
		log.Infof("  TXT 格式: %d 个", txtCount)
	}

	log.Infof("")
	log.Infof("运行参数:")
	log.Infof("  强制更新: %v", forceUpdate)
	log.Infof("  详细输出: %v", verbose)

	log.Infof("")
	log.Infof("=================================================")
	log.Infof("")
}

// sortSourcesByPriority 按优先级排序数据源（优先级数值越小越优先）
func sortSourcesByPriority(sources map[string]DownloadedSource) []DownloadedSource {
	var sorted []DownloadedSource
	for name, ds := range sources {
		ds.Name = name
		sorted = append(sorted, ds)
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Source.Priority < sorted[j].Source.Priority
	})

	return sorted
}

// convertCSVSource 转换 CSV 数据源
func convertCSVSource(converter *Converter, ds DownloadedSource) error {
	switch ds.Source.Format {
	case "ipinfo":
		return converter.ConvertCSV(ds.Path, ConvertIPInfoCSV)
	case "dbip-asn":
		return converter.ConvertCSV(ds.Path, ConvertDBIPASNCSV)
	case "dbip-city":
		return converter.ConvertCSV(ds.Path, ConvertDBIPCityCSV)
	case "dbip-country":
		return converter.ConvertCSV(ds.Path, ConvertDBIPCountryCSV)
	default:
		log.Warnf("未知的 CSV 格式: %s", ds.Source.Format)
		return nil
	}
}

// convertTextSource 转换文本数据源（CIDR 列表）
func convertTextSource(converter *Converter, ds DownloadedSource) error {
	file, err := os.Open(ds.Path)
	if err != nil {
		log.Errorf("打开文本文件失败: %s - %v", ds.Name, err)
		return xerror.WrapError(err, ErrFileOperationFailed, "open text file "+ds.Name+" failed")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	skipCount := 0
	errorCount := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			skipCount++
			continue
		}

		geo, startIP, endIP, err := ConvertChinaIPListTXT(line)
		if err != nil {
			errorCount++
			log.Warnf("解析 CIDR 失败: %s - %v", line, err)
			continue
		}

		if geo == nil {
			skipCount++
			continue
		}

		// 根据不同的格式设置 ISP/Org 信息
		switch ds.Source.Format {
		case "cidr-cernet":
			geo.ISP = "CERNET"
			geo.Org = "China Education and Research Network"
		case "cidr-chinanet":
			geo.ISP = "ChinaNet"
			geo.Org = "China Telecom"
		case "cidr-cmcc":
			geo.ISP = "CMCC"
			geo.Org = "China Mobile"
		case "cidr-unicom":
			geo.ISP = "China Unicom"
			geo.Org = "China Unicom"
		case "cidr-drpeng":
			geo.ISP = "Dr.Peng"
			geo.Org = "Dr.Peng Network"
		}

		if err := converter.writer.InsertGeoIPRange(startIP, endIP, geo); err != nil {
			errorCount++
			log.Errorf("插入 CIDR 失败: %s - %v", line, err)
			continue
		}

		count++
		if count%10000 == 0 {
			log.Infof("已处理 %d 条记录 (跳过 %d, 错误 %d)", count, skipCount, errorCount)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Errorf("读取文本文件失败: %s - %v", ds.Name, err)
		return xerror.WrapError(err, ErrFileOperationFailed, "scan text file "+ds.Name+" failed")
	}

	log.Infof("完成转换 %s: 成功 %d 条，跳过 %d 条，错误 %d 条", ds.Name, count, skipCount, errorCount)
	return nil
}

// testGeoIPQuery 测试 GeoIP 查询并美化展示
func testGeoIPQuery(dbPath string) {
	testIP := "223.5.5.5" // 阿里云公共 DNS

	log.Infof("")
	log.Infof("=================================================")
	log.Infof("测试 GeoIP 查询")
	log.Infof("=================================================")
	log.Infof("")

	// 打开数据库
	reader, err := geoip.NewMaxMindReader(dbPath)
	if err != nil {
		log.Errorf("打开 GeoIP 数据库失败: %v", err)
		return
	}
	defer reader.Close()

	// 查询 IP
	result, err := reader.LookupString(testIP)
	if err != nil {
		log.Errorf("查询 IP 失败: %v", err)
		return
	}

	// 使用 pterm 美化展示
	displayGeoIPResult(testIP, result)

	log.Infof("")
	log.Infof("数据库类型: %s", reader.Metadata().DatabaseType)
	log.Infof("使用方法:")
	log.Infof("  reader, _ := geoip.NewMaxMindReader(\"%s\")", dbPath)
	log.Infof("  result, _ := reader.LookupString(\"8.8.8.8\")")
	log.Infof("")
}

// displayGeoIPResult 使用 pterm 美化展示 GeoIP 查询结果
func displayGeoIPResult(ip string, geo *prism.GeoIP) {
	if geo == nil {
		pterm.Error.Println("未找到 IP 信息")
		return
	}

	// 标题
	pterm.DefaultHeader.WithFullWidth().
		WithBackgroundStyle(pterm.NewStyle(pterm.BgCyan)).
		WithTextStyle(pterm.NewStyle(pterm.FgBlack)).
		Printfln("GeoIP 查询结果: %s", ip)

	pterm.Println()

	// IP 基本信息
	pterm.DefaultSection.Println("IP 信息")
	ipTableData := pterm.TableData{
		{"字段", "值"},
		{"IP 地址", geo.IP},
	}
	if geo.IPVersion != 0 {
		ipTableData = append(ipTableData, []string{"IP 版本", fmt.Sprintf("IPv%d", geo.IPVersion)})
	}
	pterm.DefaultTable.WithHasHeader().WithData(ipTableData).Render()
	pterm.Println()

	// 地理位置信息
	if geo.CountryCode != "" || geo.Country != "" {
		pterm.DefaultSection.Println("地理位置")
		tableData := pterm.TableData{
			{"字段", "值"},
		}

		if geo.Country != "" {
			tableData = append(tableData, []string{"国家", fmt.Sprintf("%s (%s)", geo.Country, geo.CountryCode)})
		} else if geo.CountryCode != "" {
			tableData = append(tableData, []string{"国家代码", geo.CountryCode})
		}

		if geo.Continent != "" {
			tableData = append(tableData, []string{"大洲", fmt.Sprintf("%s (%s)", geo.Continent, geo.ContinentCode)})
		} else if geo.ContinentCode != "" {
			tableData = append(tableData, []string{"大洲代码", geo.ContinentCode})
		}

		if geo.Region != "" {
			if geo.RegionCode != "" {
				tableData = append(tableData, []string{"地区", fmt.Sprintf("%s (%s)", geo.Region, geo.RegionCode)})
			} else {
				tableData = append(tableData, []string{"地区", geo.Region})
			}
		} else if geo.RegionCode != "" {
			tableData = append(tableData, []string{"地区代码", geo.RegionCode})
		}

		if geo.City != "" {
			tableData = append(tableData, []string{"城市", geo.City})
		}

		if geo.Latitude != 0 || geo.Longitude != 0 {
			tableData = append(tableData, []string{"经纬度", fmt.Sprintf("%.4f, %.4f", geo.Latitude, geo.Longitude)})
		}

		if geo.Timezone != "" {
			tableData = append(tableData, []string{"时区", geo.Timezone})
		}

		if geo.Postal != "" {
			tableData = append(tableData, []string{"邮编", geo.Postal})
		}

		pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		pterm.Println()
	}

	// ASN 信息
	if geo.ASN != 0 || geo.ASName != "" {
		pterm.DefaultSection.Println("自治系统 (ASN)")
		tableData := pterm.TableData{
			{"字段", "值"},
		}

		if geo.ASN != 0 {
			tableData = append(tableData, []string{"ASN", fmt.Sprintf("%d", geo.ASN)})
		}

		if geo.AS != "" {
			tableData = append(tableData, []string{"AS", geo.AS})
		}

		if geo.ASName != "" {
			tableData = append(tableData, []string{"AS 名称", geo.ASName})
		}

		pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		pterm.Println()
	}

	// ISP 信息
	if geo.ISP != "" || geo.Org != "" {
		pterm.DefaultSection.Println("ISP / 组织")
		tableData := pterm.TableData{
			{"字段", "值"},
		}

		if geo.ISP != "" {
			tableData = append(tableData, []string{"ISP", geo.ISP})
		}

		if geo.Org != "" {
			tableData = append(tableData, []string{"组织", geo.Org})
		}

		pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		pterm.Println()
	}

	// 其他标记
	if geo.Proxy || geo.Hosting {
		pterm.DefaultSection.Println("特殊标记")
		tableData := pterm.TableData{
			{"类型", "状态"},
		}

		if geo.Proxy {
			tableData = append(tableData, []string{"代理/VPN", pterm.Green("是")})
		}

		if geo.Hosting {
			tableData = append(tableData, []string{"托管服务器", pterm.Green("是")})
		}

		pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		pterm.Println()
	}

	// 成功提示
	pterm.Success.Printfln("成功查询到 %s 的 GeoIP 信息", ip)
}

// getFileExtension 从 URL 中提取文件扩展名，忽略 query 参数
// 支持双扩展名如 .csv.gz, .tar.gz 等
func getFileExtension(urlStr string) string {
	// 移除 query 参数
	if idx := strings.Index(urlStr, "?"); idx != -1 {
		urlStr = urlStr[:idx]
	}

	// 移除 fragment
	if idx := strings.Index(urlStr, "#"); idx != -1 {
		urlStr = urlStr[:idx]
	}

	// 提取基础文件名
	base := filepath.Base(urlStr)

	// 检查是否是双扩展名（如 .csv.gz, .tar.gz）
	if strings.HasSuffix(base, ".csv.gz") {
		return ".csv.gz"
	} else if strings.HasSuffix(base, ".tar.gz") {
		return ".tar.gz"
	}

	// 返回普通扩展名
	return filepath.Ext(base)
}

// backupOldDatabase 备份旧数据库
func backupOldDatabase(oldPath, backupPath string) error {
	// 检查旧文件是否存在
	oldInfo, err := os.Stat(oldPath)
	if os.IsNotExist(err) {
		log.Infof("首次创建数据库，无需备份")
		return nil
	}
	if err != nil {
		return xerror.WrapError(err, ErrFileOperationFailed, "stat old database failed")
	}

	// 删除旧的备份文件
	if _, err := os.Stat(backupPath); err == nil {
		log.Infof("删除旧备份: %s", filepath.Base(backupPath))
		if err := os.Remove(backupPath); err != nil {
			return xerror.WrapError(err, ErrFileOperationFailed, "remove old backup failed")
		}
	}

	// 复制旧文件到备份位置
	oldSizeMB := float64(oldInfo.Size()) / (1024 * 1024)
	log.Infof("备份当前数据库: %s (%.2f MB) -> %s",
		filepath.Base(oldPath), oldSizeMB, filepath.Base(backupPath))

	return copyFile(oldPath, backupPath)
}

// restoreBackup 恢复备份
func restoreBackup(backupPath, targetPath string) error {
	// 检查备份文件是否存在
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return xerror.WrapError(err, ErrFileOperationFailed, "backup file not found")
	}

	log.Infof("恢复备份: %s -> %s", filepath.Base(backupPath), filepath.Base(targetPath))

	// 删除损坏的文件
	if _, err := os.Stat(targetPath); err == nil {
		if err := os.Remove(targetPath); err != nil {
			return xerror.WrapError(err, ErrFileOperationFailed, "remove corrupted file failed")
		}
	}

	// 恢复备份
	return os.Rename(backupPath, targetPath)
}
