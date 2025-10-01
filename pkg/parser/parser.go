package parser

import (
	"bufio"
	"bytes"

	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/json"
	clashConvert "github.com/metacubex/mihomo/common/convert"
	"gopkg.in/yaml.v3"
)

// Parse 解析订阅内容
// 支持多种格式：Clash YAML、V2Ray Base64 链接、逐行 JSON
// 返回 []map[string]any 可直接用于 mihomo adapter.ParseProxy
func Parse(data []byte) ([]map[string]any, error) {
	if len(data) == 0 {
		return nil, nil
	}

	reader := bytes.NewBuffer(data)

	// 优先尝试 Clash YAML 格式
	proxies, err := parseClash(reader)
	if err == nil && len(proxies) > 0 {
		return proxies, nil
	}

	// 尝试 V2Ray 格式（Base64 编码的代理链接）
	proxies, err = parseV2Ray(reader)
	if err == nil && len(proxies) > 0 {
		return proxies, nil
	}

	// 尝试逐行 JSON 格式
	return parseLineByLine(reader)
}

// parseClash 解析 Clash YAML 格式
func parseClash(reader *bytes.Buffer) ([]map[string]any, error) {
	var config struct {
		Proxies []map[string]any `yaml:"proxies"`
	}

	err := yaml.NewDecoder(reader).Decode(&config)
	if err != nil {
		log.Debugf("parseClash failed: %v", err)
		return nil, err
	}

	return config.Proxies, nil
}

// parseV2Ray 解析 V2Ray 格式（Base64 编码的代理链接）
// 使用 mihomo 的内置转换器
func parseV2Ray(reader *bytes.Buffer) ([]map[string]any, error) {
	proxies, err := clashConvert.ConvertsV2Ray(reader.Bytes())
	if err != nil {
		log.Debugf("parseV2Ray failed: %v", err)
		return nil, err
	}

	return proxies, nil
}

// parseLineByLine 逐行解析 JSON 格式
// 每行可以是一个 JSON 对象（Clash 节点格式）
func parseLineByLine(reader *bytes.Buffer) ([]map[string]any, error) {
	var proxies []map[string]any
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Bytes()

		// 清理前后缀空格和 Clash YAML 的 "-" 前缀
		for len(line) > 0 && (line[0] == '-' || line[0] == ' ') {
			line = line[1:]
		}
		for len(line) > 0 && line[len(line)-1] == ' ' {
			line = line[:len(line)-1]
		}

		if len(line) == 0 {
			continue
		}

		var proxy map[string]any
		err := json.Unmarshal(line, &proxy)
		if err != nil {
			log.Debugf("parseLineByLine failed on line: %v", err)
			return nil, err
		}

		proxies = append(proxies, proxy)
	}

	return proxies, scanner.Err()
}
