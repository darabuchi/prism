package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/darabuchi/prism"
	"github.com/darabuchi/prism/pkg/rules"
	"github.com/darabuchi/prism/scripts/tools/rule-collector/collector"
	"github.com/pterm/pterm"
)

// ruleWithMeta 带元数据的规则
type ruleWithMeta struct {
	rule         rules.Rule
	source       string // 数据源
	priority     int    // 数据源优先级
	typePriority int    // 规则类型优先级
}

// Collector 规则收集器
type Collector struct {
	cfg *Config

	// 数据源收集器映射
	handlers map[string]collector.Handler

	// 收集的规则列表（带元数据）
	rulesMeta []ruleWithMeta

	// 规则去重映射（Type|Payload -> rulesMeta 数组索引）
	ruleMap map[string]int

	// 数据源优先级映射
	sourcePriority map[string]int

	// 规则类型优先级映射
	typePriority map[rules.RuleType]int
}

// NewCollector 创建新的收集器
func NewCollector(cfg *Config) *Collector {
	c := &Collector{
		cfg:            cfg,
		handlers:       make(map[string]collector.Handler),
		rulesMeta:      make([]ruleWithMeta, 0),
		ruleMap:        make(map[string]int),
		sourcePriority: make(map[string]int),
		typePriority:   make(map[rules.RuleType]int),
	}

	// 初始化规则类型优先级
	// 优先级从低到高：越精确的规则优先级越高
	c.initTypePriority()

	return c
}

// initTypePriority 初始化规则类型优先级
func (c *Collector) initTypePriority() {
	// 基础规则类型优先级（数字越大优先级越高）
	// 1-10: 通用匹配规则（低优先级）
	c.typePriority[rules.TypeMatch] = 1      // MATCH - 最低优先级
	c.typePriority[rules.TypeGEOIP] = 2      // GEOIP
	c.typePriority[rules.TypeGeoSite] = 2    // GEOSITE
	c.typePriority[rules.TypeIPASN] = 3      // IP-ASN

	// 11-30: IP 相关规则
	c.typePriority[rules.TypeIPCIDR] = 15    // IP-CIDR
	c.typePriority[rules.TypeIPCIDR6] = 15   // IP-CIDR6
	c.typePriority[rules.TypeIPSuffix] = 16  // IP-SUFFIX
	c.typePriority[rules.TypeSrcIPCIDR] = 15 // SRC-IP-CIDR
	c.typePriority[rules.TypeSrcIP] = 16     // SRC-IP

	// 31-60: 域名相关规则（从模糊到精确）
	c.typePriority[rules.TypeDomainKeyword] = 35 // DOMAIN-KEYWORD - 关键字匹配，最模糊
	c.typePriority[rules.TypeDomainRegex] = 40   // DOMAIN-REGEX - 正则匹配
	c.typePriority[rules.TypeDomainSuffix] = 45  // DOMAIN-SUFFIX - 后缀匹配
	c.typePriority[rules.TypeDomain] = 50        // DOMAIN - 精确匹配，最精确

	// 61-80: 进程相关规则
	c.typePriority[rules.TypeProcessNameRegex] = 65 // PROCESS-NAME-REGEX
	c.typePriority[rules.TypeProcessPathRegex] = 65 // PROCESS-PATH-REGEX
	c.typePriority[rules.TypeProcess] = 70          // PROCESS-NAME
	c.typePriority[rules.TypeProcessPath] = 70      // PROCESS-PATH

	// 81-100: 其他规则
	c.typePriority[rules.TypeNetwork] = 85   // NETWORK
	c.typePriority[rules.TypeInType] = 85    // IN-TYPE
	c.typePriority[rules.TypeInName] = 85    // IN-NAME
	c.typePriority[rules.TypeInUser] = 85    // IN-USER
	c.typePriority[rules.TypeDstPort] = 90   // DST-PORT
	c.typePriority[rules.TypeSrcPort] = 90   // SRC-PORT
	c.typePriority[rules.TypeUID] = 95       // UID
	c.typePriority[rules.TypeDSCP] = 95      // DSCP
	c.typePriority[rules.TypeRuleSet] = 100  // RULE-SET - 规则集，最高优先级
}

// AddHandle 添加数据源处理器
func (c *Collector) AddHandle(name string, handler collector.Handler) {
	c.handlers[name] = handler
}

// SetSourcePriority 设置数据源优先级
func (c *Collector) SetSourcePriority(source string, priority int) {
	c.sourcePriority[source] = priority
}

// Parse 解析指定数据源的规则
func (c *Collector) Parse(source, path string, payload prism.Payload, priority int) error {
	handler, ok := c.handlers[source]
	if !ok {
		return fmt.Errorf("unknown source: %s", source)
	}

	pterm.Info.Printfln("解析 %s/%s (动作: %s, 优先级: %d)", source, path, payload.String(), priority)

	// 设置数据源优先级
	if priority > 0 {
		c.SetSourcePriority(source, priority)
	}

	// 下载规则数据（带缓存）
	body, err := c.downloadWithCache(source, path, handler)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	pterm.Info.Printfln("已下载 %s/%s (%d 字节)", source, path, len(body))

	// 解析规则
	ruleList, err := handler.Parse(body, payload)
	if err != nil {
		return fmt.Errorf("parse failed: %w", err)
	}

	// 添加规则并去重
	bar, _ := pterm.DefaultProgressbar.
		WithTotal(len(ruleList)).
		WithTitle(fmt.Sprintf("处理 %s/%s", source, path)).
		Start()

	added := 0
	for _, rule := range ruleList {
		if c.addRule(rule, source, priority) {
			added++
		}
		bar.Increment()
	}
	bar.Stop()

	pterm.Success.Printfln("从 %s/%s 添加了 %d/%d 条规则", source, path, added, len(ruleList))

	return nil
}

// downloadWithCache 带缓存的下载
func (c *Collector) downloadWithCache(source, path string, handler collector.Handler) ([]byte, error) {
	// 计算缓存文件路径（使用配置的缓存目录）
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s_%s", source, path)))
	cachePath := filepath.Join(c.cfg.CacheDir, hex.EncodeToString(hash[:]))

	// 检查缓存是否存在且有效
	if info, err := os.Stat(cachePath); err == nil {
		if !handler.NeedUpdate(info, c.cfg.CacheDays) {
			// 使用缓存
			data, err := os.ReadFile(cachePath)
			if err == nil {
				pterm.Info.Printfln("使用缓存: %s/%s", source, path)
				return data, nil
			}
		}
	}

	// 下载新数据
	pterm.Info.Printfln("下载中: %s/%s...", source, path)
	data, err := handler.Download(path)
	if err != nil {
		return nil, err
	}

	// 保存到缓存
	err = os.MkdirAll(filepath.Dir(cachePath), 0755)
	if err != nil {
		return data, nil // 忽略缓存保存失败
	}

	err = os.WriteFile(cachePath, data, 0644)
	if err != nil {
		pterm.Warning.Printfln("保存缓存失败: %v", err)
	}

	return data, nil
}

// addRule 添加规则（带去重和优先级处理）
func (c *Collector) addRule(rule rules.Rule, source string, sourcePriority int) bool {
	// 生成规则唯一键（不含 Action）
	key := c.ruleKey(rule)

	// 获取规则类型优先级
	typePriority := c.getTypePriority(rule.Type())

	// 检查是否已存在
	if existingIdx, exists := c.ruleMap[key]; exists {
		// 已存在，比较优先级
		existingMeta := c.rulesMeta[existingIdx]

		// 优先级比较顺序：
		// 1. 规则类型优先级（Type Priority）
		// 2. Action 优先级
		// 3. 数据源优先级（Source Priority）

		// 比较规则类型优先级
		if typePriority != existingMeta.typePriority {
			if typePriority > existingMeta.typePriority {
				// 新规则类型优先级更高，替换
				c.rulesMeta[existingIdx] = ruleWithMeta{
					rule:         rule,
					source:       source,
					priority:     sourcePriority,
					typePriority: typePriority,
				}
				return true
			}
			return false
		}

		// 比较 Action 优先级
		actionPriNew := c.actionPriority(rule.Action())
		actionPriOld := c.actionPriority(existingMeta.rule.Action())
		if actionPriNew != actionPriOld {
			if actionPriNew > actionPriOld {
				// 新规则 Action 优先级更高，替换
				c.rulesMeta[existingIdx] = ruleWithMeta{
					rule:         rule,
					source:       source,
					priority:     sourcePriority,
					typePriority: typePriority,
				}
				return true
			}
			return false
		}

		// 比较数据源优先级
		if sourcePriority > existingMeta.priority {
			// 新规则数据源优先级更高，替换
			c.rulesMeta[existingIdx] = ruleWithMeta{
				rule:         rule,
				source:       source,
				priority:     sourcePriority,
				typePriority: typePriority,
			}
			return true
		}

		return false
	}

	// 添加新规则
	idx := len(c.rulesMeta)
	c.rulesMeta = append(c.rulesMeta, ruleWithMeta{
		rule:         rule,
		source:       source,
		priority:     sourcePriority,
		typePriority: typePriority,
	})
	c.ruleMap[key] = idx

	return true
}

// getTypePriority 获取规则类型优先级
func (c *Collector) getTypePriority(ruleType rules.RuleType) int {
	if priority, ok := c.typePriority[ruleType]; ok {
		return priority
	}
	// 未知类型返回默认优先级
	return 50
}

// ruleKey 生成规则唯一键（Type + Payload）
func (c *Collector) ruleKey(rule rules.Rule) string {
	return fmt.Sprintf("%s|%s", rule.Type(), rule.Payload())
}

// actionPriority 返回 Action 的优先级
// 优先级从低到高：reject(1) < direct(2) < proxy(3) < 其他服务(4)
func (c *Collector) actionPriority(action rules.ActionType) int {
	actionStr := strings.ToUpper(string(action))

	switch actionStr {
	case "REJECT":
		return 1
	case "DIRECT":
		return 2
	case "PROXY":
		return 3
	default:
		// 其他服务（OPENAI, NETFLIX, YOUTUBE 等）都是高优先级
		return 4
	}
}

// Export 导出规则
func (c *Collector) Export() error {
	outputDir := c.cfg.OutputDir
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 提取纯规则列表用于合并
	ruleList := make([]rules.Rule, len(c.rulesMeta))
	for i, meta := range c.rulesMeta {
		ruleList[i] = meta.rule
	}

	// 执行规则合并优化
	pterm.Info.Println("执行规则合并优化...")
	originalCount := len(ruleList)
	merger := NewRuleMerger(ruleList)
	ruleList = merger.Merge()
	mergedCount := len(ruleList)

	if originalCount > mergedCount {
		pterm.Success.Printfln("规则合并完成：%d → %d（减少 %d 条）",
			originalCount, mergedCount, originalCount-mergedCount)
	} else {
		pterm.Info.Printfln("规则合并完成：保持 %d 条规则", mergedCount)
	}

	// 按 action 分组规则
	rulesByAction := make(map[string][]rules.Rule)
	for _, rule := range ruleList {
		action := string(rule.Action())
		rulesByAction[action] = append(rulesByAction[action], rule)
	}

	// 导出每个分组
	for action, rules := range rulesByAction {
		err := c.exportRuleFile(outputDir, action, rules)
		if err != nil {
			return fmt.Errorf("failed to export %s rules: %w", action, err)
		}
	}

	// 导出汇总文件
	err = c.exportAllRules(outputDir, ruleList)
	if err != nil {
		return fmt.Errorf("failed to export all rules: %w", err)
	}

	// 导出 subconverter 格式
	err = c.exportSubconverter(outputDir, rulesByAction)
	if err != nil {
		return fmt.Errorf("failed to export subconverter files: %w", err)
	}

	return nil
}

// exportRuleFile 导出单个规则文件
func (c *Collector) exportRuleFile(dir, action string, ruleList []rules.Rule) error {
	filename := filepath.Join(dir, fmt.Sprintf("%s.txt", strings.ToLower(action)))

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入文件头
	fmt.Fprintf(file, "# Prism Rules - %s\n", action)
	fmt.Fprintf(file, "# Generated at: %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(file, "# Total rules: %d\n\n", len(ruleList))

	// 写入规则
	for _, rule := range ruleList {
		fmt.Fprintf(file, "%s,%s,%s\n", rule.Type(), rule.Payload(), rule.Action())
	}

	pterm.Success.Printfln("导出 %d 条规则到 %s", len(ruleList), filename)
	return nil
}

// exportAllRules 导出所有规则到单个文件
func (c *Collector) exportAllRules(dir string, ruleList []rules.Rule) error {
	filename := filepath.Join(dir, "all_rules.txt")

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入文件头
	fmt.Fprintf(file, "# Prism Rules - All\n")
	fmt.Fprintf(file, "# Generated at: %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(file, "# Total rules: %d\n\n", len(ruleList))

	// 写入规则
	for _, rule := range ruleList {
		fmt.Fprintf(file, "%s,%s,%s\n", rule.Type(), rule.Payload(), rule.Action())
	}

	pterm.Success.Printfln("导出 %d 条规则到 %s", len(ruleList), filename)
	return nil
}

// exportSubconverter 导出 subconverter 兼容格式
func (c *Collector) exportSubconverter(baseDir string, rulesByAction map[string][]rules.Rule) error {
	// 创建 subconverter 子目录
	subconverterDir := filepath.Join(baseDir, "subconverter")
	err := os.MkdirAll(subconverterDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create subconverter directory: %w", err)
	}

	pterm.Info.Println("导出 subconverter 格式...")

	// 为每个 action 导出 .list 文件
	for action, ruleList := range rulesByAction {
		err := c.exportSubconverterList(subconverterDir, action, ruleList)
		if err != nil {
			return fmt.Errorf("failed to export subconverter list for %s: %w", action, err)
		}
	}

	// 生成 subconverter 配置示例和 README
	err = c.exportSubconverterConfig(subconverterDir, rulesByAction)
	if err != nil {
		return fmt.Errorf("failed to export subconverter config: %w", err)
	}

	// 生成实际的 INI 配置文件
	err = c.exportSubconverterINI(subconverterDir, rulesByAction)
	if err != nil {
		return fmt.Errorf("failed to export subconverter INI: %w", err)
	}

	pterm.Success.Printfln("subconverter 格式导出完成: %s", subconverterDir)
	return nil
}

// exportSubconverterList 导出单个 subconverter .list 文件
func (c *Collector) exportSubconverterList(dir, action string, ruleList []rules.Rule) error {
	filename := filepath.Join(dir, fmt.Sprintf("%s.list", strings.ToLower(action)))

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入文件头（注释）
	fmt.Fprintf(file, "# Prism Rules - %s (Subconverter Format)\n", action)
	fmt.Fprintf(file, "# Generated at: %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(file, "# Total rules: %d\n", len(ruleList))
	fmt.Fprintf(file, "# Repository: https://github.com/darabuchi/prism\n")
	fmt.Fprintf(file, "# Usage: Add this URL to your subconverter ruleset configuration\n")
	fmt.Fprintf(file, "#\n")
	fmt.Fprintf(file, "# Example:\n")
	fmt.Fprintf(file, "# ruleset=%s,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/%s.list\n",
		action, strings.ToLower(action))
	fmt.Fprintf(file, "\n")

	// 写入规则（subconverter 格式：TYPE,PAYLOAD[,no-resolve]）
	for _, rule := range ruleList {
		ruleType := string(rule.Type())
		payload := rule.Payload()

		// 对于 IP-CIDR 规则，添加 no-resolve 标志
		if ruleType == "IP-CIDR" || ruleType == "IP-CIDR6" || ruleType == "SRC-IP-CIDR" {
			fmt.Fprintf(file, "%s,%s,no-resolve\n", ruleType, payload)
		} else {
			fmt.Fprintf(file, "%s,%s\n", ruleType, payload)
		}
	}

	return nil
}

// exportSubconverterConfig 导出 subconverter 配置示例
func (c *Collector) exportSubconverterConfig(dir string, rulesByAction map[string][]rules.Rule) error {
	filename := filepath.Join(dir, "README.md")

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入 README
	fmt.Fprintf(file, "# Prism Rules for Subconverter\n\n")
	fmt.Fprintf(file, "本目录包含适用于 [subconverter](https://github.com/tindy2013/subconverter) 的规则集文件。\n\n")
	fmt.Fprintf(file, "## 使用方法\n\n")
	fmt.Fprintf(file, "在 subconverter 的配置文件中添加以下规则集：\n\n")
	fmt.Fprintf(file, "```ini\n")
	fmt.Fprintf(file, "[custom]\n")
	fmt.Fprintf(file, "; 启用规则生成\n")
	fmt.Fprintf(file, "enable_rule_generator=true\n")
	fmt.Fprintf(file, "overwrite_original_rules=true\n\n")

	// 按 action 排序以保持一致性
	actions := make([]string, 0, len(rulesByAction))
	for action := range rulesByAction {
		actions = append(actions, action)
	}
	// 简单排序：REJECT, DIRECT, PROXY, 其他服务
	sortActions := func(actions []string) {
		order := map[string]int{
			"Reject": 1,
			"Direct": 2,
			"Proxy":  3,
		}
		for i := 0; i < len(actions); i++ {
			for j := i + 1; j < len(actions); j++ {
				orderI := order[actions[i]]
				orderJ := order[actions[j]]
				if orderI == 0 {
					orderI = 100
				}
				if orderJ == 0 {
					orderJ = 100
				}
				if orderI > orderJ || (orderI == orderJ && actions[i] > actions[j]) {
					actions[i], actions[j] = actions[j], actions[i]
				}
			}
		}
	}
	sortActions(actions)

	// 写入规则集配置
	for _, action := range actions {
		ruleCount := len(rulesByAction[action])

		// 生成友好的分组名称
		groupName := c.getGroupName(action)

		fmt.Fprintf(file, "; %s (%d 条规则)\n", action, ruleCount)

		// 将文件名中的空格替换为 %20
		filenameEncoded := strings.ReplaceAll(strings.ToLower(action), " ", "%20")
		fmt.Fprintf(file, "ruleset=%s,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/%s.list\n",
			groupName, filenameEncoded)
	}

	fmt.Fprintf(file, "\n; GEOIP 规则\n")
	fmt.Fprintf(file, "ruleset=🎯 全球直连,[]GEOIP,CN\n")
	fmt.Fprintf(file, "ruleset=🐟 漏网之鱼,[]FINAL\n")
	fmt.Fprintf(file, "```\n\n")

	// 写入规则列表
	fmt.Fprintf(file, "## 可用规则集\n\n")
	fmt.Fprintf(file, "| 规则集 | 规则数量 | 下载链接 |\n")
	fmt.Fprintf(file, "|--------|----------|----------|\n")

	for _, action := range actions {
		ruleCount := len(rulesByAction[action])
		// 将文件名中的空格替换为 %20
		filenameEncoded := strings.ReplaceAll(strings.ToLower(action), " ", "%20")
		downloadURL := fmt.Sprintf("https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/%s.list",
			filenameEncoded)
		fmt.Fprintf(file, "| %s | %d | [下载](%s) |\n", action, ruleCount, downloadURL)
	}

	fmt.Fprintf(file, "\n## 规则格式\n\n")
	fmt.Fprintf(file, "所有规则文件遵循 subconverter 的 `.list` 格式：\n\n")
	fmt.Fprintf(file, "```\n")
	fmt.Fprintf(file, "DOMAIN,example.com\n")
	fmt.Fprintf(file, "DOMAIN-SUFFIX,example.com\n")
	fmt.Fprintf(file, "DOMAIN-KEYWORD,keyword\n")
	fmt.Fprintf(file, "IP-CIDR,192.168.0.0/16,no-resolve\n")
	fmt.Fprintf(file, "IP-CIDR6,2001:db8::/32,no-resolve\n")
	fmt.Fprintf(file, "```\n\n")

	fmt.Fprintf(file, "## 更新频率\n\n")
	fmt.Fprintf(file, "规则集每日自动更新，来源于：\n")
	fmt.Fprintf(file, "- [blackmatrix7/ios_rule_script](https://github.com/blackmatrix7/ios_rule_script)\n")
	fmt.Fprintf(file, "- [Loyalsoldier/clash-rules](https://github.com/Loyalsoldier/clash-rules)\n")
	fmt.Fprintf(file, "- [ACL4SSR/ACL4SSR](https://github.com/ACL4SSR/ACL4SSR)\n\n")

	fmt.Fprintf(file, "## 许可证\n\n")
	fmt.Fprintf(file, "本项目采用 GPL-3.0 许可证。详见 [LICENSE](../../LICENSE) 文件。\n")

	return nil
}

// exportSubconverterINI 导出实际的 INI 配置文件
func (c *Collector) exportSubconverterINI(dir string, rulesByAction map[string][]rules.Rule) error {
	filename := filepath.Join(dir, "prism.ini")

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入 INI 文件头部
	fmt.Fprintf(file, "[custom]\n")
	fmt.Fprintf(file, "; Prism Rules Configuration for Subconverter\n")
	fmt.Fprintf(file, "; Generated at: %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(file, "; Repository: https://github.com/darabuchi/prism\n")
	fmt.Fprintf(file, ";\n")
	fmt.Fprintf(file, "; 使用方法：\n")
	fmt.Fprintf(file, "; 1. 将此文件放到 subconverter 的 base 目录下\n")
	fmt.Fprintf(file, "; 2. 在订阅转换时使用 &config=prism 参数引用此配置\n")
	fmt.Fprintf(file, ";\n\n")

	// 启用规则生成
	fmt.Fprintf(file, "; 启用规则生成\n")
	fmt.Fprintf(file, "enable_rule_generator=true\n")
	fmt.Fprintf(file, "overwrite_original_rules=true\n\n")

	// 按 action 排序
	actions := make([]string, 0, len(rulesByAction))
	for action := range rulesByAction {
		actions = append(actions, action)
	}
	sortActions := func(actions []string) {
		order := map[string]int{
			"Reject": 1,
			"Direct": 2,
			"Proxy":  3,
		}
		for i := 0; i < len(actions); i++ {
			for j := i + 1; j < len(actions); j++ {
				orderI := order[actions[i]]
				orderJ := order[actions[j]]
				if orderI == 0 {
					orderI = 100
				}
				if orderJ == 0 {
					orderJ = 100
				}
				if orderI > orderJ || (orderI == orderJ && actions[i] > actions[j]) {
					actions[i], actions[j] = actions[j], actions[i]
				}
			}
		}
	}
	sortActions(actions)

	// 写入规则集配置
	fmt.Fprintf(file, "; ============ 规则集配置 ============\n\n")

	for _, action := range actions {
		ruleCount := len(rulesByAction[action])

		// 生成友好的分组名称
		groupName := c.getGroupName(action)

		fmt.Fprintf(file, "; %s (%d 条规则)\n", action, ruleCount)

		// 将文件名中的空格替换为 %20
		filenameEncoded := strings.ReplaceAll(strings.ToLower(action), " ", "%20")
		fmt.Fprintf(file, "ruleset=%s,https://raw.githubusercontent.com/darabuchi/prism/main/rules/subconverter/%s.list\n\n",
			groupName, filenameEncoded)
	}

	// 添加 GEOIP 和 FINAL 规则
	fmt.Fprintf(file, "; ============ 最终规则 ============\n\n")
	fmt.Fprintf(file, "; GEOIP 规则\n")
	fmt.Fprintf(file, "ruleset=🎯 全球直连,[]GEOIP,CN\n\n")
	fmt.Fprintf(file, "; 兜底规则\n")
	fmt.Fprintf(file, "ruleset=🐟 漏网之鱼,[]FINAL\n")

	pterm.Success.Printfln("已生成 INI 配置文件: %s", filename)
	return nil
}

// getGroupName 获取友好的分组名称
func (c *Collector) getGroupName(action string) string {
	// 特殊处理基础动作
	switch strings.ToUpper(action) {
	case "REJECT":
		return "🛡️ 广告拦截"
	case "DIRECT":
		return "🎯 全球直连"
	case "PROXY":
		return "🚀 节点选择"
	}

	// AI 服务
	aiServices := map[string]string{
		"OPENAI":           "🤖 OpenAI",
		"CLAUDE":           "🤖 Claude",
		"GEMINI":           "🤖 Gemini",
		"COPILOT":          "🤖 Copilot",
		"BING":             "🤖 Bing AI",
		"PERPLEXITY":       "🤖 Perplexity",
		"CHARACTER.AI":     "🤖 Character.AI",
		"MIDJOURNEY":       "🎨 Midjourney",
		"STABLE DIFFUSION": "🎨 Stable Diffusion",
		"CHATGPT":          "🤖 ChatGPT",
	}
	if groupName, ok := aiServices[strings.ToUpper(action)]; ok {
		return groupName
	}

	// 流媒体服务
	streamingServices := map[string]string{
		"YOUTUBE":        "📹 YouTube",
		"NETFLIX":        "🎬 Netflix",
		"DISNEY":         "🎬 Disney+",
		"SPOTIFY":        "🎵 Spotify",
		"TIKTOK":         "📹 TikTok",
		"BILIBILI":       "📺 哔哩哔哩",
		"BILIBILIINTL":   "📺 哔哩哔哩",
		"BILIBILI HK":    "📺 哔哩哔哩港澳台",
		"IQIYI":          "📺 爱奇艺",
		"TENCENT VIDEO":  "📺 腾讯视频",
		"APPLE TV":       "📺 Apple TV",
		"APPLE MUSIC":    "🎵 Apple Music",
	}
	if groupName, ok := streamingServices[strings.ToUpper(action)]; ok {
		return groupName
	}

	// 社交平台
	socialServices := map[string]string{
		"TELEGRAM":  "💬 Telegram",
		"TWITTER":   "🐦 Twitter",
		"FACEBOOK":  "📘 Facebook",
		"INSTAGRAM": "📷 Instagram",
		"DISCORD":   "💬 Discord",
	}
	if groupName, ok := socialServices[strings.ToUpper(action)]; ok {
		return groupName
	}

	// 开发工具
	devServices := map[string]string{
		"GITHUB":  "💻 GitHub",
		"GITLAB":  "💻 GitLab",
		"DOCKER":  "🐳 Docker",
		"VERCEL":  "⚡ Vercel",
	}
	if groupName, ok := devServices[strings.ToUpper(action)]; ok {
		return groupName
	}

	// 其他服务使用默认格式
	return fmt.Sprintf("📦 %s", action)
}

// RuleCount 返回收集的规则数量
func (c *Collector) RuleCount() int {
	return len(c.rulesMeta)
}
