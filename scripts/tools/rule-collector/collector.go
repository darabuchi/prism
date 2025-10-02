package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/darabuchi/prism/pkg/rules"
	"github.com/darabuchi/prism/scripts/tools/rule-collector/collector"
	"github.com/pterm/pterm"
)

// Collector 规则收集器
type Collector struct {
	cfg *Config

	// 数据源收集器映射
	handlers map[string]collector.Handler

	// 收集的规则列表
	rules []rules.Rule

	// 规则去重映射（Type|Payload -> rules 数组索引）
	ruleMap map[string]int
}

// NewCollector 创建新的收集器
func NewCollector(cfg *Config) *Collector {
	return &Collector{
		cfg:      cfg,
		handlers: make(map[string]collector.Handler),
		rules:    make([]rules.Rule, 0),
		ruleMap:  make(map[string]int),
	}
}

// AddHandle 添加数据源处理器
func (c *Collector) AddHandle(name string, handler collector.Handler) {
	c.handlers[name] = handler
}

// Parse 解析指定数据源的规则
func (c *Collector) Parse(source, path, action string) error {
	handler, ok := c.handlers[source]
	if !ok {
		return fmt.Errorf("unknown source: %s", source)
	}

	pterm.Info.Printfln("解析 %s/%s (动作: %s)", source, path, action)

	// 下载规则数据（带缓存）
	body, err := c.downloadWithCache(source, path, handler)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	pterm.Info.Printfln("已下载 %s/%s (%d 字节)", source, path, len(body))

	// 解析规则
	ruleList, err := handler.Parse(body, action)
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
		if c.addRule(rule) {
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
func (c *Collector) addRule(rule rules.Rule) bool {
	// 生成规则唯一键（不含 Action）
	key := c.ruleKey(rule)

	// 检查是否已存在
	if existingIdx, exists := c.ruleMap[key]; exists {
		// 已存在，比较优先级
		existingRule := c.rules[existingIdx]
		if c.actionPriority(rule.Action()) > c.actionPriority(existingRule.Action()) {
			// 新规则优先级更高，替换旧规则
			c.rules[existingIdx] = rule
			return true
		}
		return false
	}

	// 添加新规则
	idx := len(c.rules)
	c.rules = append(c.rules, rule)
	c.ruleMap[key] = idx

	return true
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

	// 执行规则合并优化
	pterm.Info.Println("执行规则合并优化...")
	originalCount := len(c.rules)
	merger := NewRuleMerger(c.rules)
	c.rules = merger.Merge()
	mergedCount := len(c.rules)

	if originalCount > mergedCount {
		pterm.Success.Printfln("规则合并完成：%d → %d（减少 %d 条）",
			originalCount, mergedCount, originalCount-mergedCount)
	} else {
		pterm.Info.Printfln("规则合并完成：保持 %d 条规则", mergedCount)
	}

	// 按 action 分组规则
	rulesByAction := make(map[string][]rules.Rule)
	for _, rule := range c.rules {
		action := string(rule.Action())
		rulesByAction[action] = append(rulesByAction[action], rule)
	}

	// 导出每个分组
	for action, ruleList := range rulesByAction {
		err := c.exportRuleFile(outputDir, action, ruleList)
		if err != nil {
			return fmt.Errorf("failed to export %s rules: %w", action, err)
		}
	}

	// 导出汇总文件
	err = c.exportAllRules(outputDir)
	if err != nil {
		return fmt.Errorf("failed to export all rules: %w", err)
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
func (c *Collector) exportAllRules(dir string) error {
	filename := filepath.Join(dir, "all_rules.txt")

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入文件头
	fmt.Fprintf(file, "# Prism Rules - All\n")
	fmt.Fprintf(file, "# Generated at: %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(file, "# Total rules: %d\n\n", len(c.rules))

	// 写入规则
	for _, rule := range c.rules {
		fmt.Fprintf(file, "%s,%s,%s\n", rule.Type(), rule.Payload(), rule.Action())
	}

	pterm.Success.Printfln("导出 %d 条规则到 %s", len(c.rules), filename)
	return nil
}

// RuleCount 返回收集的规则数量
func (c *Collector) RuleCount() int {
	return len(c.rules)
}
