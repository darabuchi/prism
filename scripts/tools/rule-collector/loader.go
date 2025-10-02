package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/darabuchi/prism/pkg/rules"
	"github.com/pterm/pterm"
)

// LoadRulesFromDirectory 从指定目录递归加载所有 .rule 文件
func (c *Collector) LoadRulesFromDirectory(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		pterm.Info.Printfln("规则目录不存在，跳过: %s", dir)
		return nil
	}

	pterm.Info.Printfln("正在加载自定义规则目录: %s", dir)

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if d.IsDir() {
			return nil
		}

		// 只处理 .rule 文件
		if filepath.Ext(path) != ".rule" {
			return nil
		}

		return c.loadRuleFile(path)
	})

	if err != nil {
		return fmt.Errorf("遍历规则目录失败: %w", err)
	}

	return nil
}

// loadRuleFile 加载单个 .rule 文件
func (c *Collector) loadRuleFile(path string) error {
	pterm.Debug.Printfln("加载规则文件: %s", path)

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("打开文件失败 %s: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析规则
		rule, err := rules.ParseRule(line)
		if err != nil {
			pterm.Warning.Printfln("解析规则失败 %s:%d: %s (错误: %v)",
				filepath.Base(path), lineNum, line, err)
			continue
		}

		// 添加到规则集合
		c.addRule(rule)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取文件失败 %s: %w", path, err)
	}

	return nil
}

// LoadRuleFile 加载单个规则文件（支持嵌入式文件系统）
func (c *Collector) LoadRuleFile(path string) error {
	pterm.Debug.Printfln("尝试加载规则文件: %s", path)

	// 尝试从文件系统加载
	if _, err := os.Stat(path); err == nil {
		return c.loadRuleFile(path)
	}

	pterm.Warning.Printfln("规则文件不存在，跳过: %s", path)
	return nil
}
