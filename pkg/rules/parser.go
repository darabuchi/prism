package rules

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/darabuchi/prism"
)

// ParseRule 从字符串解析规则
//
// 规则格式: TYPE,PAYLOAD,ACTION
// 示例:
//   - DOMAIN,google.com,PROXY
//   - DOMAIN-SUFFIX,google.com,PROXY
//   - IP-CIDR,192.168.0.0/16,DIRECT
//   - GEOIP,CN,DIRECT
//   - MATCH,PROXY
func ParseRule(ruleStr string) (Rule, error) {
	ruleStr = strings.TrimSpace(ruleStr)

	// 跳过空行和注释
	if ruleStr == "" || strings.HasPrefix(ruleStr, "#") || strings.HasPrefix(ruleStr, "//") {
		return nil, nil
	}

	// 分割规则字符串
	parts := strings.Split(ruleStr, ",")
	if len(parts) < 2 {
		return nil, fmt.Errorf("%w: %s", ErrInvalidRule, ruleStr)
	}

	ruleType := strings.TrimSpace(strings.ToUpper(parts[0]))

	var payload string
	var action prism.Payload
	var noResolve bool

	// MATCH 规则只需要类型和动作
	if ruleType == "MATCH" {
		if len(parts) < 2 {
			return nil, fmt.Errorf("%w: MATCH rule requires action", ErrInvalidRule)
		}
		actionStr := strings.TrimSpace(strings.ToUpper(parts[1]))
		action = prism.ParsePayload(actionStr)
		return NewMatch(action), nil
	}

	// 其他规则需要至少 3 个部分
	if len(parts) < 3 {
		return nil, fmt.Errorf("%w: %s", ErrInvalidRule, ruleStr)
	}

	// IP-SUFFIX 特殊处理：格式为 IP-SUFFIX,CIDR,SUFFIX,ACTION
	if ruleType == string(TypeIPSuffix) {
		if len(parts) < 4 {
			return nil, fmt.Errorf("%w: IP-SUFFIX requires CIDR, SUFFIX and ACTION", ErrInvalidRule)
		}
		// 组合 CIDR 和 SUFFIX 作为 payload
		payload = strings.TrimSpace(parts[1]) + "," + strings.TrimSpace(parts[2])
		actionStr := strings.TrimSpace(strings.ToUpper(parts[3]))
		action = prism.ParsePayload(actionStr)

		// 检查是否有 no-resolve 选项
		if len(parts) > 4 {
			for i := 4; i < len(parts); i++ {
				option := strings.TrimSpace(strings.ToLower(parts[i]))
				if option == "no-resolve" {
					noResolve = true
				}
			}
		}
	} else {
		payload = strings.TrimSpace(parts[1])
		actionStr := strings.TrimSpace(strings.ToUpper(parts[2]))
		action = prism.ParsePayload(actionStr)

		// 检查是否有 no-resolve 选项
		if len(parts) > 3 {
			for i := 3; i < len(parts); i++ {
				option := strings.TrimSpace(strings.ToLower(parts[i]))
				if option == "no-resolve" {
					noResolve = true
				}
			}
		}
	}

	// 验证 payload
	if payload == "" {
		return nil, fmt.Errorf("%w: empty payload", ErrInvalidPayload)
	}

	// 根据规则类型创建对应的规则
	parsedType := prism.ParseRuleType(ruleType)
	switch parsedType {
	case TypeDomain:
		return NewDomain(payload, action), nil

	case TypeDomainSuffix:
		return NewDomainSuffix(payload, action), nil

	case TypeDomainKeyword:
		return NewDomainKeyword(payload, action), nil

	case TypeDomainRegex:
		rule, err := NewDomainRegex(payload, action)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidRegex, err)
		}
		return rule, nil

	case TypeGEOIP:
		rule := NewGEOIP(payload, action)
		if noResolve {
			rule.NoResolve(true)
		}
		return rule, nil

	case TypeIPCIDR, TypeIPCIDR6:
		isIPv6 := parsedType == TypeIPCIDR6
		rule, err := NewIPCIDR(payload, action, isIPv6)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidIPCIDR, err)
		}
		if noResolve {
			rule.NoResolve(true)
		}
		return rule, nil

	case TypeDstPort:
		return NewPort(payload, action, true)

	case TypeSrcPort:
		return NewPort(payload, action, false)

	case TypeProcess:
		return NewProcess(payload, action, false), nil

	case TypeProcessPath:
		return NewProcess(payload, action, true), nil

	case TypeIPASN:
		rule, err := NewIPASN(payload, action)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidPayload, err)
		}
		if noResolve {
			rule.NoResolve(true)
		}
		return rule, nil

	case TypeInType:
		return NewInType(payload, action), nil

	case TypeInName:
		return NewInName(payload, action), nil

	case TypeInUser:
		return NewInUser(payload, action), nil

	case TypeNetwork:
		return NewNetwork(payload, action), nil

	case TypeUID:
		rule, err := NewUID(payload, action)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidPayload, err)
		}
		return rule, nil

	case TypeDSCP:
		rule, err := NewDSCP(payload, action)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidPayload, err)
		}
		return rule, nil

	case TypeProcessNameRegex:
		rule, err := NewProcessRegex(payload, action, false)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidRegex, err)
		}
		return rule, nil

	case TypeProcessPathRegex:
		rule, err := NewProcessRegex(payload, action, true)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidRegex, err)
		}
		return rule, nil

	case TypeIPSuffix:
		rule, err := NewIPSuffix(payload, action)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidIPCIDR, err)
		}
		if noResolve {
			rule.NoResolve(true)
		}
		return rule, nil

	case TypeGeoSite:
		return NewGeoSite(payload, action), nil

	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidRuleType, ruleType)
	}
}

// ParseRules 从字符串切片解析多条规则
func ParseRules(ruleStrs []string) ([]Rule, error) {
	var rules []Rule
	for i, ruleStr := range ruleStrs {
		rule, err := ParseRule(ruleStr)
		if err != nil {
			return nil, fmt.Errorf("parse rule %d failed: %w", i+1, err)
		}
		if rule != nil {
			rules = append(rules, rule)
		}
	}
	return rules, nil
}

// LoadRulesFromFile 从文件加载规则
//
// 文件格式: 每行一条规则
// 支持注释（# 或 //）和空行
func LoadRulesFromFile(filename string) ([]Rule, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open file failed: %w", err)
	}
	defer file.Close()

	var rules []Rule
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		rule, err := ParseRule(line)
		if err != nil {
			return nil, fmt.Errorf("parse line %d failed: %w", lineNum, err)
		}

		if rule != nil {
			rules = append(rules, rule)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan file failed: %w", err)
	}

	return rules, nil
}

// MatchFirst 使用规则列表匹配元数据，返回第一个匹配的规则
func MatchFirst(rules []Rule, metadata *Metadata) (Rule, bool) {
	for _, rule := range rules {
		if rule.Match(metadata) {
			return rule, true
		}
	}
	return nil, false
}
