package rules

import (
	"strconv"
	"strings"

	"github.com/darabuchi/prism"
)

// Port 端口匹配规则
type Port struct {
	base
	ports      []uint16 // 单个端口列表
	portRanges [][2]uint16 // 端口范围列表
	isDst      bool // true: 目标端口, false: 源端口
}

// Match 匹配端口
func (p *Port) Match(metadata *Metadata) bool {
	var port uint16
	if p.isDst {
		port = metadata.DstPort
	} else {
		port = metadata.SrcPort
	}

	// 检查单个端口
	for _, p := range p.ports {
		if port == p {
			return true
		}
	}

	// 检查端口范围
	for _, r := range p.portRanges {
		if port >= r[0] && port <= r[1] {
			return true
		}
	}

	return false
}

// NewPort 创建端口匹配规则
//
// payload 格式:
// - 单个端口: "80"
// - 多个端口: "80,443,8080"
// - 端口范围: "8000-9000"
// - 混合: "80,443,8000-9000"
func NewPort(payload string, action prism.Payload, isDst bool) (*Port, error) {
	var ports []uint16
	var portRanges [][2]uint16

	parts := strings.Split(payload, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// 检查是否是范围
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, ErrInvalidPortRange
			}

			start, err := strconv.ParseUint(strings.TrimSpace(rangeParts[0]), 10, 16)
			if err != nil {
				return nil, err
			}

			end, err := strconv.ParseUint(strings.TrimSpace(rangeParts[1]), 10, 16)
			if err != nil {
				return nil, err
			}

			if start > end {
				return nil, ErrInvalidPortRange
			}

			portRanges = append(portRanges, [2]uint16{uint16(start), uint16(end)})
		} else {
			// 单个端口
			port, err := strconv.ParseUint(part, 10, 16)
			if err != nil {
				return nil, err
			}
			ports = append(ports, uint16(port))
		}
	}

	ruleType := TypeSrcPort
	if isDst {
		ruleType = TypeDstPort
	}

	return &Port{
		base:       newBase(ruleType, payload, action),
		ports:      ports,
		portRanges: portRanges,
		isDst:      isDst,
	}, nil
}
