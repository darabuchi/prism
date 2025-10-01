package rules

import "strings"

// DomainKeyword 域名关键字匹配规则
type DomainKeyword struct {
	base
	keyword string
}

// Match 匹配域名关键字
func (d *DomainKeyword) Match(metadata *Metadata) bool {
	if metadata.Domain == "" {
		if metadata.Host != "" {
			metadata.Domain = strings.ToLower(metadata.Host)
		} else {
			return false
		}
	}
	return strings.Contains(metadata.Domain, d.keyword)
}

// NewDomainKeyword 创建域名关键字匹配规则
func NewDomainKeyword(payload string, action ActionType) *DomainKeyword {
	keyword := strings.ToLower(strings.TrimSpace(payload))
	return &DomainKeyword{
		base:    newBase(TypeDomainKeyword, payload, action),
		keyword: keyword,
	}
}
