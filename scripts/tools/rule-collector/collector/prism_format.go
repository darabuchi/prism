package collector

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/darabuchi/prism/pkg/rules"
)

// PrismRuleFile Prism 规则文件格式
// 文件结构：
// - 文件头（16 字节）
//   - Magic Number: "PRISM" (5 字节)
//   - Version: uint8 (1 字节)
//   - Reserved: [10]byte (10 字节)
// - Action 数量: uint16 (2 字节)
// - Action 列表：
//   - Action 名称长度: uint8 (1 字节)
//   - Action 名称: string (变长)
//   - 规则数量: uint32 (4 字节)
//   - 规则列表：
//     - 规则类型: uint8 (1 字节)
//       - 1: DOMAIN
//       - 2: DOMAIN-SUFFIX
//       - 3: IP-CIDR
//       - 4: IP-CIDR6
//     - Payload 长度: uint16 (2 字节)
//     - Payload: string (变长)

const (
	MagicNumber = "PRISM"
	Version     = uint8(1)
)

const (
	RuleTypeDomain       uint8 = 1
	RuleTypeDomainSuffix uint8 = 2
	RuleTypeIPCIDR       uint8 = 3
	RuleTypeIPCIDR6      uint8 = 4
)

// ExportPrismBinary 导出 Prism 二进制规则文件
func ExportPrismBinary(w io.Writer, rulesByAction map[string][]rules.Rule) error {
	// 写入文件头
	if _, err := w.Write([]byte(MagicNumber)); err != nil {
		return fmt.Errorf("write magic number: %w", err)
	}
	if err := binary.Write(w, binary.LittleEndian, Version); err != nil {
		return fmt.Errorf("write version: %w", err)
	}
	// 保留字节
	reserved := make([]byte, 10)
	if _, err := w.Write(reserved); err != nil {
		return fmt.Errorf("write reserved: %w", err)
	}

	// 写入 Action 数量
	actionCount := uint16(len(rulesByAction))
	if err := binary.Write(w, binary.LittleEndian, actionCount); err != nil {
		return fmt.Errorf("write action count: %w", err)
	}

	// 写入每个 Action 及其规则
	for action, ruleList := range rulesByAction {
		// 写入 Action 名称长度
		actionNameLen := uint8(len(action))
		if err := binary.Write(w, binary.LittleEndian, actionNameLen); err != nil {
			return fmt.Errorf("write action name length: %w", err)
		}

		// 写入 Action 名称
		if _, err := w.Write([]byte(action)); err != nil {
			return fmt.Errorf("write action name: %w", err)
		}

		// 统计有效规则数量（只支持的类型）
		validRules := make([]struct {
			ruleType uint8
			payload  string
		}, 0)

		for _, rule := range ruleList {
			ruleType := string(rule.Type())
			payload := rule.Payload()

			var rt uint8
			switch ruleType {
			case "DOMAIN":
				rt = RuleTypeDomain
			case "DOMAIN-SUFFIX":
				rt = RuleTypeDomainSuffix
			case "IP-CIDR":
				rt = RuleTypeIPCIDR
			case "IP-CIDR6":
				rt = RuleTypeIPCIDR6
			default:
				// 跳过不支持的规则类型
				continue
			}

			validRules = append(validRules, struct {
				ruleType uint8
				payload  string
			}{rt, payload})
		}

		// 写入规则数量
		ruleCount := uint32(len(validRules))
		if err := binary.Write(w, binary.LittleEndian, ruleCount); err != nil {
			return fmt.Errorf("write rule count: %w", err)
		}

		// 写入每条规则
		for _, validRule := range validRules {
			// 写入规则类型
			if err := binary.Write(w, binary.LittleEndian, validRule.ruleType); err != nil {
				return fmt.Errorf("write rule type: %w", err)
			}

			// 写入 Payload 长度
			payloadLen := uint16(len(validRule.payload))
			if err := binary.Write(w, binary.LittleEndian, payloadLen); err != nil {
				return fmt.Errorf("write payload length: %w", err)
			}

			// 写入 Payload
			if _, err := w.Write([]byte(validRule.payload)); err != nil {
				return fmt.Errorf("write payload: %w", err)
			}
		}
	}

	return nil
}

// ParsePrismBinary 解析 Prism 二进制规则文件
func ParsePrismBinary(r io.Reader) (map[string][]rules.Rule, error) {
	// 读取文件头
	magic := make([]byte, 5)
	if _, err := io.ReadFull(r, magic); err != nil {
		return nil, fmt.Errorf("read magic number: %w", err)
	}
	if !bytes.Equal(magic, []byte(MagicNumber)) {
		return nil, fmt.Errorf("invalid magic number: %s", magic)
	}

	var version uint8
	if err := binary.Read(r, binary.LittleEndian, &version); err != nil {
		return nil, fmt.Errorf("read version: %w", err)
	}
	if version != Version {
		return nil, fmt.Errorf("unsupported version: %d", version)
	}

	// 跳过保留字节
	reserved := make([]byte, 10)
	if _, err := io.ReadFull(r, reserved); err != nil {
		return nil, fmt.Errorf("read reserved: %w", err)
	}

	// 读取 Action 数量
	var actionCount uint16
	if err := binary.Read(r, binary.LittleEndian, &actionCount); err != nil {
		return nil, fmt.Errorf("read action count: %w", err)
	}

	rulesByAction := make(map[string][]rules.Rule)

	// 读取每个 Action 及其规则
	for i := 0; i < int(actionCount); i++ {
		// 读取 Action 名称长度
		var actionNameLen uint8
		if err := binary.Read(r, binary.LittleEndian, &actionNameLen); err != nil {
			return nil, fmt.Errorf("read action name length: %w", err)
		}

		// 读取 Action 名称
		actionNameBytes := make([]byte, actionNameLen)
		if _, err := io.ReadFull(r, actionNameBytes); err != nil {
			return nil, fmt.Errorf("read action name: %w", err)
		}
		action := string(actionNameBytes)

		// 读取规则数量
		var ruleCount uint32
		if err := binary.Read(r, binary.LittleEndian, &ruleCount); err != nil {
			return nil, fmt.Errorf("read rule count: %w", err)
		}

		ruleList := make([]rules.Rule, 0, ruleCount)

		// 读取每条规则
		for j := 0; j < int(ruleCount); j++ {
			// 读取规则类型
			var ruleType uint8
			if err := binary.Read(r, binary.LittleEndian, &ruleType); err != nil {
				return nil, fmt.Errorf("read rule type: %w", err)
			}

			// 读取 Payload 长度
			var payloadLen uint16
			if err := binary.Read(r, binary.LittleEndian, &payloadLen); err != nil {
				return nil, fmt.Errorf("read payload length: %w", err)
			}

			// 读取 Payload
			payloadBytes := make([]byte, payloadLen)
			if _, err := io.ReadFull(r, payloadBytes); err != nil {
				return nil, fmt.Errorf("read payload: %w", err)
			}
			payload := string(payloadBytes)

			// 构造规则
			var rule rules.Rule
			var err error
			actionType := rules.ActionType(action)

			switch ruleType {
			case RuleTypeDomain:
				rule = rules.NewDomain(payload, actionType)
			case RuleTypeDomainSuffix:
				rule = rules.NewDomainSuffix(payload, actionType)
			case RuleTypeIPCIDR, RuleTypeIPCIDR6:
				rule, err = rules.NewIPCIDR(payload, actionType, false)
				if err != nil {
					return nil, fmt.Errorf("parse IP-CIDR: %w", err)
				}
			default:
				return nil, fmt.Errorf("unknown rule type: %d", ruleType)
			}

			ruleList = append(ruleList, rule)
		}

		rulesByAction[action] = ruleList
	}

	return rulesByAction, nil
}
