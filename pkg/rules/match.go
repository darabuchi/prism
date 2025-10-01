package rules

// Match 匹配所有流量的规则
//
// 通常作为规则列表的最后一条规则，确保所有流量都有匹配的动作
type Match struct {
	base
}

// Match 总是返回 true
func (m *Match) Match(metadata *Metadata) bool {
	return true
}

// NewMatch 创建匹配所有流量的规则
func NewMatch(action ActionType) *Match {
	return &Match{
		base: newBase(TypeMatch, "all", action),
	}
}
