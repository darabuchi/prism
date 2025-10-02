package rules

import (
	"fmt"

	"github.com/darabuchi/prism"
)

// base 规则基础实现
//
// 包含所有规则的通用字段和方法
type base struct {
	ruleType RuleType
	action   prism.Payload
	payload  string
}

// Type 返回规则类型
func (b *base) Type() RuleType {
	return b.ruleType
}

// Action 返回规则动作
func (b *base) Action() prism.Payload {
	return b.action
}

// Content 返回规则内容
func (b *base) Content() string {
	return b.payload
}

// Payload 返回规则载荷（已废弃，使用 Content 代替）
// Deprecated: 使用 Content() 代替
func (b *base) Payload() string {
	return b.payload
}

// String 返回规则的字符串表示
func (b *base) String() string {
	return fmt.Sprintf("%s,%s,%s", b.ruleType, b.payload, b.action)
}

// newBase 创建基础规则
func newBase(ruleType RuleType, payload string, action prism.Payload) base {
	return base{
		ruleType: ruleType,
		payload:  payload,
		action:   action,
	}
}
