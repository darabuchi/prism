package main

import (
	"net"
	"sort"
	"strings"

	"github.com/darabuchi/prism/pkg/rules"
)

// RuleMerger 规则合并器
type RuleMerger struct {
	rules []rules.Rule
}

// NewRuleMerger 创建规则合并器
func NewRuleMerger(ruleList []rules.Rule) *RuleMerger {
	return &RuleMerger{
		rules: ruleList,
	}
}

// Merge 执行规则合并
// 1. 相同内容去重（已在 Collector 中实现）
// 2. 相同功效内容去重（DOMAIN 被 DOMAIN-SUFFIX 覆盖）
// 3. IP-CIDR 合并（合并相邻 CIDR）
// 4. 保持规则优先级（靠前规则优先）
func (m *RuleMerger) Merge() []rules.Rule {
	// 按 Action 分组
	groupedRules := make(map[string][]rules.Rule)
	for _, rule := range m.rules {
		action := string(rule.Action())
		groupedRules[action] = append(groupedRules[action], rule)
	}

	// 对每个分组进行合并
	merged := make([]rules.Rule, 0, len(m.rules))
	for _, ruleList := range groupedRules {
		mergedGroup := m.mergeGroup(ruleList)
		merged = append(merged, mergedGroup...)
	}

	return merged
}

// mergeGroup 合并同一 Action 的规则
func (m *RuleMerger) mergeGroup(ruleList []rules.Rule) []rules.Rule {
	// 按规则类型分组
	byType := make(map[string][]rules.Rule)
	for _, rule := range ruleList {
		ruleType := string(rule.Type())
		byType[ruleType] = append(byType[ruleType], rule)
	}

	// 处理域名规则去重（DOMAIN 被 DOMAIN-SUFFIX 覆盖）
	m.deduplicateDomainRules(byType)

	// 合并 IP-CIDR 规则
	m.mergeCIDRRules(byType)

	// 收集所有规则
	result := make([]rules.Rule, 0, len(ruleList))
	for _, rules := range byType {
		result = append(result, rules...)
	}

	return result
}

// deduplicateDomainRules 去重域名规则
// 如果存在 DOMAIN-SUFFIX,example.com，则移除 DOMAIN,example.com 和 DOMAIN,*.example.com
func (m *RuleMerger) deduplicateDomainRules(byType map[string][]rules.Rule) {
	domainSuffixSet := make(map[string]bool)
	domainKeywordSet := make(map[string]bool)

	// 收集所有 DOMAIN-SUFFIX 和 DOMAIN-KEYWORD
	for _, rule := range byType["DOMAIN-SUFFIX"] {
		domainSuffixSet[strings.ToLower(rule.Payload())] = true
	}
	for _, rule := range byType["DOMAIN-KEYWORD"] {
		domainKeywordSet[strings.ToLower(rule.Payload())] = true
	}

	// 过滤 DOMAIN 规则
	if len(domainSuffixSet) > 0 || len(domainKeywordSet) > 0 {
		filtered := make([]rules.Rule, 0, len(byType["DOMAIN"]))
		for _, rule := range byType["DOMAIN"] {
			domain := strings.ToLower(rule.Payload())

			// 检查是否被 DOMAIN-SUFFIX 覆盖
			covered := false
			for suffix := range domainSuffixSet {
				if domain == suffix || strings.HasSuffix(domain, "."+suffix) {
					covered = true
					break
				}
			}

			// 检查是否被 DOMAIN-KEYWORD 覆盖
			if !covered {
				for keyword := range domainKeywordSet {
					if strings.Contains(domain, keyword) {
						covered = true
						break
					}
				}
			}

			if !covered {
				filtered = append(filtered, rule)
			}
		}
		byType["DOMAIN"] = filtered
	}
}

// mergeCIDRRules 合并 IP-CIDR 规则
func (m *RuleMerger) mergeCIDRRules(byType map[string][]rules.Rule) {
	cidrRules := byType["IP-CIDR"]
	if len(cidrRules) == 0 {
		return
	}

	cidr6Rules := byType["IP-CIDR6"]

	// 合并 IPv4 CIDR
	if len(cidrRules) > 0 {
		byType["IP-CIDR"] = m.mergeCIDRList(cidrRules)
	}

	// 合并 IPv6 CIDR
	if len(cidr6Rules) > 0 {
		byType["IP-CIDR6"] = m.mergeCIDRList(cidr6Rules)
	}
}

// mergeCIDRList 合并 CIDR 列表
func (m *RuleMerger) mergeCIDRList(cidrRules []rules.Rule) []rules.Rule {
	// 解析所有 CIDR
	type cidrInfo struct {
		rule    rules.Rule
		network *net.IPNet
		ip      net.IP
		ones    int
		bits    int
	}

	cidrs := make([]cidrInfo, 0, len(cidrRules))
	for _, rule := range cidrRules {
		payload := rule.Payload()
		// 移除可能的 no-resolve 参数
		payload = strings.Split(payload, ",")[0]

		_, network, err := net.ParseCIDR(payload)
		if err != nil {
			// 如果解析失败，保留原规则
			cidrs = append(cidrs, cidrInfo{rule: rule})
			continue
		}

		ones, bits := network.Mask.Size()
		cidrs = append(cidrs, cidrInfo{
			rule:    rule,
			network: network,
			ip:      network.IP,
			ones:    ones,
			bits:    bits,
		})
	}

	// 按 IP 地址和掩码长度排序
	sort.Slice(cidrs, func(i, j int) bool {
		if cidrs[i].network == nil {
			return false
		}
		if cidrs[j].network == nil {
			return true
		}

		// 先按 IP 地址排序
		ipCmp := compareIP(cidrs[i].ip, cidrs[j].ip)
		if ipCmp != 0 {
			return ipCmp < 0
		}

		// 再按掩码长度排序（短掩码优先，因为它覆盖范围更大）
		return cidrs[i].ones < cidrs[j].ones
	})

	// 去除被包含的 CIDR
	result := make([]rules.Rule, 0, len(cidrs))
	for i, cidr := range cidrs {
		if cidr.network == nil {
			result = append(result, cidr.rule)
			continue
		}

		// 检查是否被前面的 CIDR 包含
		covered := false
		for j := 0; j < i; j++ {
			if cidrs[j].network == nil {
				continue
			}

			// 如果前面的 CIDR 掩码更短（覆盖范围更大），且包含当前 IP
			if cidrs[j].ones <= cidr.ones && cidrs[j].network.Contains(cidr.ip) {
				covered = true
				break
			}
		}

		if !covered {
			result = append(result, cidr.rule)
		}
	}

	return result
}

// compareIP 比较两个 IP 地址
func compareIP(a, b net.IP) int {
	// 统一转换为 16 字节格式
	a = a.To16()
	b = b.To16()

	for i := 0; i < len(a); i++ {
		if a[i] < b[i] {
			return -1
		}
		if a[i] > b[i] {
			return 1
		}
	}
	return 0
}
