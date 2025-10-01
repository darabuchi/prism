package node

import (
	"fmt"

	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/xerror"
	"github.com/lazygophers/utils/cryptox"
	"github.com/lazygophers/utils/json"
)

// GenerateID 生成节点唯一ID
// 基于完整配置的 SHA256 哈希
func GenerateID(config map[string]any) (string, error) {
	// 序列化配置
	data, err := json.Marshal(config)
	if err != nil {
		log.Errorf("err:%v", err)
		return "", xerror.Wrap(err, "marshal config failed")
	}

	// 计算 SHA256
	hash := cryptox.SHA256(data)
	return hash, nil
}

// GenerateUniqueKey 生成节点唯一键
// 格式：type://server:port
func GenerateUniqueKey(config map[string]any) (string, error) {
	proxyType, ok := config["type"].(string)
	if !ok || proxyType == "" {
		log.Errorf("missing type in config")
		return "", xerror.New(0, "missing type in config")
	}

	server, ok := config["server"].(string)
	if !ok || server == "" {
		log.Errorf("missing server in config")
		return "", xerror.New(0, "missing server in config")
	}

	// port 可能是 int 或 float64
	var port int
	switch v := config["port"].(type) {
	case int:
		port = v
	case float64:
		port = int(v)
	case int64:
		port = int(v)
	default:
		log.Errorf("invalid port type: %T", config["port"])
		return "", xerror.New(0, "invalid port type")
	}

	if port <= 0 {
		log.Errorf("invalid port: %d", port)
		return "", xerror.New(0, "invalid port")
	}

	return fmt.Sprintf("%s://%s:%d", proxyType, server, port), nil
}

// EncodeConfig 将配置编码为 Base64 JSON 字符串
func EncodeConfig(config map[string]any) (string, error) {
	data, err := json.Marshal(config)
	if err != nil {
		log.Errorf("err:%v", err)
		return "", xerror.Wrap(err, "marshal config failed")
	}

	return cryptox.Base64Encode(data), nil
}

// DecodeConfig 从 Base64 JSON 字符串解码配置
func DecodeConfig(encoded string) (map[string]any, error) {
	data, err := cryptox.Base64Decode(encoded)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, xerror.Wrap(err, "decode base64 failed")
	}

	var config map[string]any
	err = json.Unmarshal([]byte(data), &config)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, xerror.Wrap(err, "unmarshal config failed")
	}

	return config, nil
}
