package rules

import (
	"strings"

	"github.com/darabuchi/prism"
)

// Domain 域名精确匹配规则
type Domain struct {
	base
	domain string
}

// Match 匹配域名
func (d *Domain) Match(metadata *Metadata) bool {
	if metadata.Domain == "" {
		if metadata.Host != "" {
			metadata.Domain = strings.ToLower(metadata.Host)
		} else {
			return false
		}
	}
	return metadata.Domain == d.domain
}

// NewDomain 创建域名精确匹配规则
func NewDomain(payload string, action prism.Payload) *Domain {
	domain := strings.ToLower(strings.TrimSpace(payload))
	return &Domain{
		base:   newBase(TypeDomain, payload, action),
		domain: domain,
	}
}
