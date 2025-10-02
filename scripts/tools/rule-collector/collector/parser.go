package collector

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strings"

	"github.com/darabuchi/prism"
	"github.com/darabuchi/prism/pkg/rules"
	"gopkg.in/yaml.v3"
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
func ParseClashRules(body []byte, payload prism.Payload) ([]rules.Rule, error) {
	// 将 prism.Payload 转换为 string action，用于 pkg/rules
	action := payload.String()
	var ruleList []rules.Rule

	// 尝试使用 YAML 解析（Clash YAML 格式）
	var data struct {
		Payload []string `yaml:"payload"`
	}

	err := yaml.Unmarshal(body, &data)
	if err == nil && len(data.Payload) > 0 {
		// YAML 格式解析成功
		for _, line := range data.Payload {
			// 清理 YAML 格式的残留符号
			// 某些上游源可能在YAML字符串中包含列表语法，如: "- 'domain.com'"
			// 需要按正确顺序清理这些符号

			// 1. 移除 YAML 列表前缀 "- " (如果存在)
			line = strings.TrimPrefix(line, "- ")
			line = strings.TrimSpace(line)

			// 2. 移除前后的引号
			line = strings.Trim(line, "\"'")
			line = strings.TrimSpace(line)

			rule, err := parseClashPayloadLine(line, action)
			if err != nil {
				continue
			}
			if rule != nil {
				ruleList = append(ruleList, rule)
			}
		}
		return ruleList, nil
	}

	// 如果 YAML 解析失败，使用逐行解析（纯文本格式）
	scanner := bufio.NewScanner(bytes.NewReader(body))

	for scanner.Scan() {
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
		line = strings.TrimPrefix(line, "- ")
		line = strings.TrimSpace(line)

		// 移除 YAML 的单引号和双引号
		line = strings.Trim(line, "'\"")

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

// parseClashPayloadLine 解析 Clash payload 中的一行（参考 fire 的实现）
func parseClashPayloadLine(line, action string) (rules.Rule, error) {
	// 移除 ,no-resolve 后缀
	line = strings.ReplaceAll(line, ",no-resolve", "")

	// 处理特殊前缀
	if strings.HasPrefix(line, "+.") {
		// +.example.com -> DOMAIN-SUFFIX,example.com
		return parseRuleLine("DOMAIN-SUFFIX,"+strings.TrimPrefix(line, "+"), action)
	}

	if strings.HasPrefix(line, "*.*.") {
		// *.*.example.com -> DOMAIN-SUFFIX,.example.com
		return parseRuleLine("DOMAIN-SUFFIX,"+strings.TrimPrefix(line, "*.*"), action)
	}

	if strings.HasPrefix(line, ".") {
		// .example.com -> DOMAIN-SUFFIX,.example.com
		return parseRuleLine("DOMAIN-SUFFIX,"+line, action)
	}

	// 如果包含逗号，说明已经是完整的规则格式
	if strings.Contains(line, ",") {
		return parseRuleLine(line, action)
	}

	// 尝试解析为 IP-CIDR
	_, _, err := net.ParseCIDR(line)
	if err == nil {
		return parseRuleLine("IP-CIDR,"+line, action)
	}

	// 包含点号，当作 DOMAIN 处理
	if strings.Contains(line, ".") {
		return parseRuleLine("DOMAIN,"+line, action)
	}

	return nil, fmt.Errorf("unknown format: %s", line)
}

// ParseTextRules 解析纯文本格式的规则（每行一个域名或 IP）
func ParseTextRules(body []byte, payload prism.Payload, ruleTypeStr string) ([]rules.Rule, error) {
	// 将 prism.Payload 转换为 string action，用于 pkg/rules
	action := payload.String()
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
		switch strings.ToUpper(ruleTypeStr) {
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

	ruleTypeStr := strings.ToUpper(strings.TrimSpace(parts[0]))

	// MATCH 规则只需要类型
	if ruleTypeStr == "MATCH" {
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
	// 格式: TYPE,PAYLOAD,ACTION
	ruleStr := fmt.Sprintf("%s,%s,%s", ruleTypeStr, payload, action)

	rule, err := rules.ParseRule(ruleStr)
	if err != nil {
		return nil, fmt.Errorf("parse rule failed: %w", err)
	}

	return rule, nil
}
