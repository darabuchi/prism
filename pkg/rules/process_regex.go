package rules

import (
	"regexp"
)

// ProcessRegex 进程正则匹配规则
type ProcessRegex struct {
	base
	pattern *regexp.Regexp
	isPath  bool // true: 进程路径, false: 进程名称
}

// Match 匹配进程
func (p *ProcessRegex) Match(metadata *Metadata) bool {
	if p.isPath {
		// 匹配进程路径
		if metadata.ProcessPath == "" {
			return false
		}
		return p.pattern.MatchString(metadata.ProcessPath)
	} else {
		// 匹配进程名称
		if metadata.ProcessName == "" {
			return false
		}
		return p.pattern.MatchString(metadata.ProcessName)
	}
}

// NewProcessRegex 创建进程正则匹配规则
func NewProcessRegex(payload string, action ActionType, isPath bool) (*ProcessRegex, error) {
	pattern, err := regexp.Compile(payload)
	if err != nil {
		return nil, err
	}

	ruleType := TypeProcessNameRegex
	if isPath {
		ruleType = TypeProcessPathRegex
	}

	return &ProcessRegex{
		base:    newBase(ruleType, payload, action),
		pattern: pattern,
		isPath:  isPath,
	}, nil
}
