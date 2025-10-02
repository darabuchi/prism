package rules

import (
	"strings"

	"github.com/darabuchi/prism"
)

// DomainSuffix 域名后缀匹配规则
type DomainSuffix struct {
	base
	suffix string
}

// Match 匹配域名后缀
func (d *DomainSuffix) Match(metadata *Metadata) bool {
	if metadata.Domain == "" {
		if metadata.Host != "" {
			metadata.Domain = strings.ToLower(metadata.Host)
		} else {
			return false
		}
	}

	domain := metadata.Domain
	suffix := d.suffix

	// 精确匹配
	if domain == suffix {
		return true
	}

	// 后缀匹配（必须是完整的子域名）
	if strings.HasSuffix(domain, "."+suffix) {
		return true
	}

	return false
}

// NewDomainSuffix 创建域名后缀匹配规则
func NewDomainSuffix(payload string, action prism.Payload) *DomainSuffix {
	suffix := strings.ToLower(strings.TrimSpace(payload))
	// 移除前导点
	suffix = strings.TrimPrefix(suffix, ".")

	return &DomainSuffix{
		base:   newBase(TypeDomainSuffix, payload, action),
		suffix: suffix,
	}
}
