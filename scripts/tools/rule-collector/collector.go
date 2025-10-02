package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/darabuchi/prism/pkg/rules"
	"github.com/darabuchi/prism/scripts/tools/rule-collector/collector"
	"github.com/pterm/pterm"
)

// Collector 规则收集器
type Collector struct {
	// 数据源收集器映射
	handlers map[string]collector.Handler

	// 收集的规则列表
	rules []rules.Rule

	// 规则去重映射（用于快速查找）
	ruleMap map[string]bool
}

// NewCollector 创建新的收集器
func NewCollector() *Collector {
	return &Collector{
		handlers: make(map[string]collector.Handler),
		rules:    make([]rules.Rule, 0),
		ruleMap:  make(map[string]bool),
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

	pterm.Info.Printfln("Parsing %s/%s (action: %s)", source, path, action)

	// 下载规则数据（带缓存）
	body, err := c.downloadWithCache(source, path, handler)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	pterm.Info.Printfln("Downloaded %s/%s (%d bytes)", source, path, len(body))

	// 解析规则
	ruleList, err := handler.Parse(body, action)
	if err != nil {
		return fmt.Errorf("parse failed: %w", err)
	}

	// 添加规则并去重
	bar, _ := pterm.DefaultProgressbar.
		WithTotal(len(ruleList)).
		WithTitle(fmt.Sprintf("Processing %s/%s", source, path)).
		Start()

	added := 0
	for _, rule := range ruleList {
		if c.addRule(rule) {
			added++
		}
		bar.Increment()
	}
	bar.Stop()

	pterm.Success.Printfln("Added %d/%d rules from %s/%s", added, len(ruleList), source, path)

	return nil
}

// downloadWithCache 带缓存的下载
func (c *Collector) downloadWithCache(source, path string, handler collector.Handler) ([]byte, error) {
	// 计算缓存文件路径
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s_%s", source, path)))
	cachePath := filepath.Join("scripts", "tools", "rule-collector", "tmp", "cache", hex.EncodeToString(hash[:]))

	// 检查缓存是否存在且有效
	if info, err := os.Stat(cachePath); err == nil {
		if !handler.NeedUpdate(info) {
			// 使用缓存
			data, err := os.ReadFile(cachePath)
			if err == nil {
				pterm.Info.Printfln("Using cached data for %s/%s", source, path)
				return data, nil
			}
		}
	}

	// 下载新数据
	pterm.Info.Printfln("Downloading %s/%s...", source, path)
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
		pterm.Warning.Printfln("Failed to save cache: %v", err)
	}

	return data, nil
}

// addRule 添加规则（带去重）
func (c *Collector) addRule(rule rules.Rule) bool {
	// 生成规则唯一键
	key := c.ruleKey(rule)

	// 检查是否已存在
	if c.ruleMap[key] {
		return false
	}

	// 检查是否被其他规则覆盖
	if c.isRuleCovered(rule) {
		return false
	}

	// 添加规则
	c.rules = append(c.rules, rule)
	c.ruleMap[key] = true

	return true
}

// ruleKey 生成规则唯一键
func (c *Collector) ruleKey(rule rules.Rule) string {
	return fmt.Sprintf("%s|%s|%s", rule.Type(), rule.Payload(), rule.Action())
}

// isRuleCovered 检查规则是否被已有规则覆盖
func (c *Collector) isRuleCovered(newRule rules.Rule) bool {
	// 简单去重逻辑，可以根据需要扩展
	// 例如：
	// - DOMAIN-SUFFIX 可以覆盖 DOMAIN
	// - 更大的 IP-CIDR 可以覆盖更小的

	for _, existingRule := range c.rules {
		// 只比较相同 action 的规则
		if existingRule.Action() != newRule.Action() {
			continue
		}

		// 根据规则类型进行覆盖检查
		if c.ruleCovers(existingRule, newRule) {
			return true
		}
	}

	return false
}

// ruleCovers 检查 existing 规则是否覆盖 new 规则
func (c *Collector) ruleCovers(existing, new rules.Rule) bool {
	// TODO: 实现更智能的规则覆盖检查
	// 例如：
	// - DOMAIN-SUFFIX .example.com 覆盖 DOMAIN www.example.com
	// - IP-CIDR 192.168.0.0/16 覆盖 IP-CIDR 192.168.1.0/24
	return false
}

// Export 导出规则
func (c *Collector) Export() error {
	outputDir := "scripts/tools/rule-collector/output"
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
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
	filename := filepath.Join(dir, fmt.Sprintf("%s.txt", action))

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

	pterm.Success.Printfln("Exported %d rules to %s", len(ruleList), filename)
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

	pterm.Success.Printfln("Exported %d rules to %s", len(c.rules), filename)
	return nil
}

// RuleCount 返回收集的规则数量
func (c *Collector) RuleCount() int {
	return len(c.rules)
}
