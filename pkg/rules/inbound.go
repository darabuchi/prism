package rules

import "strings"

// InType 入站类型匹配规则
type InType struct {
	base
	inType string
}

// Match 匹配入站类型
func (i *InType) Match(metadata *Metadata) bool {
	if metadata.InboundType == "" {
		return false
	}
	return strings.EqualFold(metadata.InboundType, i.inType)
}

// NewInType 创建入站类型匹配规则
func NewInType(payload string, action ActionType) *InType {
	inType := strings.ToUpper(strings.TrimSpace(payload))
	return &InType{
		base:   newBase(TypeInType, payload, action),
		inType: inType,
	}
}

// InName 入站名称匹配规则
type InName struct {
	base
	inName string
}

// Match 匹配入站名称
func (i *InName) Match(metadata *Metadata) bool {
	if metadata.InboundName == "" {
		return false
	}
	return strings.EqualFold(metadata.InboundName, i.inName)
}

// NewInName 创建入站名称匹配规则
func NewInName(payload string, action ActionType) *InName {
	inName := strings.TrimSpace(payload)
	return &InName{
		base:   newBase(TypeInName, payload, action),
		inName: inName,
	}
}

// InUser 入站用户匹配规则
type InUser struct {
	base
	user string
}

// Match 匹配入站用户
func (i *InUser) Match(metadata *Metadata) bool {
	if metadata.InboundUser == "" {
		return false
	}
	return metadata.InboundUser == i.user
}

// NewInUser 创建入站用户匹配规则
func NewInUser(payload string, action ActionType) *InUser {
	user := strings.TrimSpace(payload)
	return &InUser{
		base: newBase(TypeInUser, payload, action),
		user: user,
	}
}
