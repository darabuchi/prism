package rules

import (
	"net/netip"
	"strings"
	"sync"
)

// Engine 规则引擎
//
// 提供高性能的规则匹配功能，使用多种索引优化
type Engine struct {
	// 原始规则列表（按顺序）
	rules []Rule

	// 域名精确匹配索引
	domainIndex map[string][]Rule

	// 域名后缀索引（使用后缀树）
	suffixTree *SuffixTree

	// IP CIDR 索引（使用前缀树）
	ipv4Trie *IPTrie
	ipv6Trie *IPTrie

	// 其他规则（需要遍历匹配）
	otherRules []Rule

	mu sync.RWMutex
}

// NewEngine 创建规则引擎
func NewEngine() *Engine {
	return &Engine{
		rules:       make([]Rule, 0),
		domainIndex: make(map[string][]Rule),
		suffixTree:  NewSuffixTree(),
		ipv4Trie:    NewIPTrie(),
		ipv6Trie:    NewIPTrie(),
		otherRules:  make([]Rule, 0),
	}
}

// AddRule 添加规则
func (e *Engine) AddRule(rule Rule) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.rules = append(e.rules, rule)

	// 根据规则类型建立索引
	switch rule.Type() {
	case TypeDomain:
		// 域名精确匹配索引
		if domain, ok := rule.(*Domain); ok {
			e.domainIndex[domain.domain] = append(e.domainIndex[domain.domain], rule)
		}

	case TypeDomainSuffix:
		// 域名后缀索引
		if suffix, ok := rule.(*DomainSuffix); ok {
			e.suffixTree.Add(suffix.suffix, rule)
		}

	case TypeIPCIDR, TypeIPCIDR6:
		// IP CIDR 索引
		if cidr, ok := rule.(*IPCIDR); ok {
			if cidr.isIPv6 {
				e.ipv6Trie.Add(cidr.prefix, rule)
			} else {
				e.ipv4Trie.Add(cidr.prefix, rule)
			}
		}

	default:
		// 其他规则需要遍历匹配
		e.otherRules = append(e.otherRules, rule)
	}
}

// AddRules 批量添加规则
func (e *Engine) AddRules(rules []Rule) {
	for _, rule := range rules {
		e.AddRule(rule)
	}
}

// Match 匹配规则
//
// 返回第一个匹配的规则，如果没有匹配则返回 nil
func (e *Engine) Match(metadata *Metadata) (Rule, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// 规范化域名
	if metadata.Domain == "" && metadata.Host != "" {
		metadata.Domain = strings.ToLower(metadata.Host)
	}

	// 1. 尝试域名精确匹配（最快）
	if metadata.Domain != "" {
		if rules, ok := e.domainIndex[metadata.Domain]; ok {
			for _, rule := range rules {
				if rule.Match(metadata) {
					return rule, true
				}
			}
		}
	}

	// 2. 尝试域名后缀匹配（次快）
	if metadata.Domain != "" {
		if rules := e.suffixTree.Match(metadata.Domain); len(rules) > 0 {
			for _, rule := range rules {
				if rule.Match(metadata) {
					return rule, true
				}
			}
		}
	}

	// 3. 尝试 IP CIDR 匹配
	if metadata.DstIP.IsValid() {
		var trie *IPTrie
		if metadata.DstIP.Is4() {
			trie = e.ipv4Trie
		} else {
			trie = e.ipv6Trie
		}

		if rules := trie.Match(metadata.DstIP); len(rules) > 0 {
			for _, rule := range rules {
				if rule.Match(metadata) {
					return rule, true
				}
			}
		}
	}

	// 4. 遍历其他规则
	for _, rule := range e.otherRules {
		if rule.Match(metadata) {
			return rule, true
		}
	}

	return nil, false
}

// Rules 返回所有规则
func (e *Engine) Rules() []Rule {
	e.mu.RLock()
	defer e.mu.RUnlock()

	result := make([]Rule, len(e.rules))
	copy(result, e.rules)
	return result
}

// Clear 清空所有规则
func (e *Engine) Clear() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.rules = make([]Rule, 0)
	e.domainIndex = make(map[string][]Rule)
	e.suffixTree = NewSuffixTree()
	e.ipv4Trie = NewIPTrie()
	e.ipv6Trie = NewIPTrie()
	e.otherRules = make([]Rule, 0)
}

// SuffixTree 域名后缀树
//
// 用于高效匹配域名后缀
type SuffixTree struct {
	root *SuffixNode
	mu   sync.RWMutex
}

// SuffixNode 后缀树节点
type SuffixNode struct {
	children map[string]*SuffixNode
	rules    []Rule
}

// NewSuffixTree 创建后缀树
func NewSuffixTree() *SuffixTree {
	return &SuffixTree{
		root: &SuffixNode{
			children: make(map[string]*SuffixNode),
		},
	}
}

// Add 添加后缀和规则
func (t *SuffixTree) Add(suffix string, rule Rule) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// 分割域名部分（反向）
	parts := strings.Split(suffix, ".")
	// 反转，从 TLD 开始
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}

	node := t.root
	for _, part := range parts {
		if node.children[part] == nil {
			node.children[part] = &SuffixNode{
				children: make(map[string]*SuffixNode),
			}
		}
		node = node.children[part]
	}

	node.rules = append(node.rules, rule)
}

// Match 匹配域名，返回所有匹配的规则
func (t *SuffixTree) Match(domain string) []Rule {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// 分割域名部分（反向）
	parts := strings.Split(domain, ".")
	// 反转，从 TLD 开始
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}

	var result []Rule
	node := t.root

	// 从 TLD 向下匹配，收集所有匹配的规则
	for _, part := range parts {
		// 继续向下
		child := node.children[part]
		if child == nil {
			break
		}
		node = child
		// 收集当前节点的规则
		result = append(result, node.rules...)
	}

	return result
}

// IPTrie IP 前缀树
//
// 用于高效匹配 IP CIDR
type IPTrie struct {
	root *IPTrieNode
	mu   sync.RWMutex
}

// IPTrieNode IP 前缀树节点
type IPTrieNode struct {
	left  *IPTrieNode // 0
	right *IPTrieNode // 1
	rules []Rule
}

// NewIPTrie 创建 IP 前缀树
func NewIPTrie() *IPTrie {
	return &IPTrie{
		root: &IPTrieNode{},
	}
}

// Add 添加 IP 前缀和规则
func (t *IPTrie) Add(prefix netip.Prefix, rule Rule) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// 获取前缀的二进制表示
	addr := prefix.Addr()
	bits := prefix.Bits()

	node := t.root
	var ipBytes []byte

	if addr.Is4() {
		bytes := addr.As4()
		ipBytes = bytes[:]
	} else {
		bytes := addr.As16()
		ipBytes = bytes[:]
	}

	// 遍历前缀位
	bitPos := 0
	for bitPos < bits {
		byteIndex := bitPos / 8
		bitIndex := 7 - (bitPos % 8)

		bit := (ipBytes[byteIndex] >> bitIndex) & 1

		if bit == 0 {
			if node.left == nil {
				node.left = &IPTrieNode{}
			}
			node = node.left
		} else {
			if node.right == nil {
				node.right = &IPTrieNode{}
			}
			node = node.right
		}

		bitPos++
	}

	node.rules = append(node.rules, rule)
}

// Match 匹配 IP，返回所有匹配的规则
func (t *IPTrie) Match(ip netip.Addr) []Rule {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var result []Rule
	node := t.root

	var ipBytes []byte
	if ip.Is4() {
		bytes := ip.As4()
		ipBytes = bytes[:]
	} else {
		bytes := ip.As16()
		ipBytes = bytes[:]
	}

	// 遍历 IP 的每一位
	maxBits := len(ipBytes) * 8
	for bitPos := 0; bitPos < maxBits && node != nil; bitPos++ {
		// 收集当前节点的规则
		result = append(result, node.rules...)

		byteIndex := bitPos / 8
		bitIndex := 7 - (bitPos % 8)

		bit := (ipBytes[byteIndex] >> bitIndex) & 1

		if bit == 0 {
			node = node.left
		} else {
			node = node.right
		}
	}

	// 收集最后一个节点的规则
	if node != nil {
		result = append(result, node.rules...)
	}

	return result
}
