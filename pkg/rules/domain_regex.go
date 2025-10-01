package rules

import (
	"regexp"
	"strings"
)

// DomainRegex 域名正则表达式匹配规则
type DomainRegex struct {
	base
	pattern *regexp.Regexp
}

// Match 匹配域名正则表达式
func (d *DomainRegex) Match(metadata *Metadata) bool {
	if metadata.Domain == "" {
		if metadata.Host != "" {
			metadata.Domain = strings.ToLower(metadata.Host)
		} else {
			return false
		}
	}
	return d.pattern.MatchString(metadata.Domain)
}

// NewDomainRegex 创建域名正则表达式匹配规则
func NewDomainRegex(payload string, action ActionType) (*DomainRegex, error) {
	pattern, err := regexp.Compile(payload)
	if err != nil {
		return nil, err
	}

	return &DomainRegex{
		base:    newBase(TypeDomainRegex, payload, action),
		pattern: pattern,
	}, nil
}
