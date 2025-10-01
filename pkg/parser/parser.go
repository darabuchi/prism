package parser

import (
	"bufio"
	"bytes"

	"github.com/darabuchi/prism/pkg/node"
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/xerror"
	"github.com/lazygophers/utils/candy"
	"github.com/lazygophers/utils/json"
	clashConvert "github.com/metacubex/mihomo/common/convert"
	"gopkg.in/yaml.v3"
)

// Parse 解析订阅内容
//
// 支持多种格式：
//   - Clash YAML 格式
//   - V2Ray Base64 链接格式
//   - 逐行 JSON 格式
//
// 参数:
//   - data: 订阅内容字节数组
//
// 返回:
//   - []*node.Node: 解析后的节点列表
//   - error: 解析失败时返回错误
//
// 示例:
//
//	data := []byte(`
//	proxies:
//	  - name: "香港节点"
//	    type: ss
//	    server: hk.example.com
//	    port: 443
//	`)
//	nodes, err := Parse(data)
func Parse(data []byte) ([]*node.Node, error) {
	if len(data) == 0 {
		return nil, nil
	}

	reader := bytes.NewBuffer(data)

	// 优先尝试 Clash YAML 格式
	configs, err := parseClash(reader)
	if err == nil && len(configs) > 0 {
		return createNodes(configs)
	}

	// 尝试 V2Ray 格式（Base64 编码的代理链接）
	configs, err = parseV2Ray(reader)
	if err == nil && len(configs) > 0 {
		return createNodes(configs)
	}

	// 尝试逐行 JSON 格式
	configs, err = parseLineByLine(reader)
	if err != nil {
		return nil, err
	}

	return createNodes(configs)
}

// createNodes 从配置列表创建节点列表
//
// 将解析出的配置转换为 node.Node 实例
// 跳过创建失败的节点，记录错误但不中断整个过程
// 使用唯一 ID 进行去重
func createNodes(configs []map[string]any) ([]*node.Node, error) {
	if len(configs) == 0 {
		return nil, nil
	}

	nodes := make([]*node.Node, 0, len(configs))
	for i, config := range configs {
		n, err := node.NewNode(config)
		if err != nil {
			// 记录错误但继续处理其他节点
			log.Warnf("create node %d failed: %v, config: %+v", i, err, config)
			continue
		}
		nodes = append(nodes, n)
	}

	// 如果所有节点都创建失败，返回错误
	if len(nodes) == 0 {
		return nil, xerror.New(xerror.ErrSystemError, "all nodes creation failed")
	}

	// 按照唯一 ID 去重
	nodes = candy.UniqueUsing(nodes, func(n *node.Node) any {
		return n.UniqueId()
	})

	return nodes, nil
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
