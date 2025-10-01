package rules

import "strings"

// Process 进程匹配规则
type Process struct {
	base
	pattern string
	isPath  bool // true: 进程路径, false: 进程名称
}

// Match 匹配进程
func (p *Process) Match(metadata *Metadata) bool {
	if p.isPath {
		// 匹配进程路径
		if metadata.ProcessPath == "" {
			return false
		}
		return strings.Contains(metadata.ProcessPath, p.pattern)
	} else {
		// 匹配进程名称
		if metadata.ProcessName == "" {
			return false
		}
		return strings.EqualFold(metadata.ProcessName, p.pattern)
	}
}

// NewProcess 创建进程匹配规则
func NewProcess(payload string, action ActionType, isPath bool) *Process {
	ruleType := TypeProcess
	if isPath {
		ruleType = TypeProcessPath
	}

	return &Process{
		base:    newBase(ruleType, payload, action),
		pattern: payload,
		isPath:  isPath,
	}
}
