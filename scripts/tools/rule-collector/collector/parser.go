package collector

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/darabuchi/prism/pkg/rules"
)

// ParseClashRules 解析 Clash 格式的规则
// 支持的格式：
// - DOMAIN,example.com
// - DOMAIN-SUFFIX,example.com
// - DOMAIN-KEYWORD,google
// - IP-CIDR,192.168.0.0/16
// - IP-CIDR6,2001:db8::/32
// - GEOIP,CN
// - MATCH
func ParseClashRules(body []byte, action string) ([]rules.Rule, error) {
	var ruleList []rules.Rule

	scanner := bufio.NewScanner(bytes.NewReader(body))
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 跳过 YAML 格式的 payload: 行
		if strings.HasPrefix(line, "payload:") {
			continue
		}

		// 移除行首的 - 符号（YAML 列表格式）
		line = strings.TrimPrefix(line, "-")
		line = strings.TrimSpace(line)

		// 移除 ,no-resolve 后缀
		line = strings.ReplaceAll(line, ",no-resolve", "")

		// 解析规则
		rule, err := parseRuleLine(line, action)
		if err != nil {
			// 跳过无法解析的规则，不中断整个流程
			continue
		}

		if rule != nil {
			ruleList = append(ruleList, rule)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	return ruleList, nil
}

// ParseTextRules 解析纯文本格式的规则（每行一个域名或 IP）
func ParseTextRules(body []byte, action string, ruleType string) ([]rules.Rule, error) {
	var ruleList []rules.Rule

	scanner := bufio.NewScanner(bytes.NewReader(body))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 根据类型创建规则
		var rule rules.Rule
		switch strings.ToUpper(ruleType) {
		case "DOMAIN":
			rule = rules.NewDomain(line, rules.ActionType(action))
		case "DOMAIN-SUFFIX":
			rule = rules.NewDomainSuffix(line, rules.ActionType(action))
		case "IP-CIDR":
			r, err := rules.NewIPCIDR(line, rules.ActionType(action), false)
			if err != nil {
				continue
			}
			rule = r
		default:
			continue
		}

		if rule != nil {
			ruleList = append(ruleList, rule)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	return ruleList, nil
}

// parseRuleLine 解析单行规则
func parseRuleLine(line, action string) (rules.Rule, error) {
	// 分割规则字符串
	parts := strings.Split(line, ",")
	if len(parts) < 1 {
		return nil, fmt.Errorf("invalid rule format: %s", line)
	}

	ruleType := strings.ToUpper(strings.TrimSpace(parts[0]))

	// MATCH 规则只需要类型
	if ruleType == "MATCH" {
		return rules.NewMatch(rules.ActionType(action)), nil
	}

	// 其他规则需要 payload
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid rule format: %s", line)
	}

	payload := strings.TrimSpace(parts[1])
	if payload == "" {
		return nil, fmt.Errorf("empty payload: %s", line)
	}

	// 构造完整的规则字符串并使用 pkg/rules 的 ParseRule
	// 格式：TYPE,PAYLOAD,ACTION
	ruleStr := fmt.Sprintf("%s,%s,%s", ruleType, payload, action)

	rule, err := rules.ParseRule(ruleStr)
	if err != nil {
		return nil, fmt.Errorf("parse rule failed: %w", err)
	}

	return rule, nil
}
