package node

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/lazygophers/utils/json"
)

// GenerateID 生成节点唯一ID
// 基于完整配置的 SHA256 哈希
func GenerateID(config map[string]any) (string, error) {
	// 序列化配置
	data, err := json.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("marshal config failed: %w", err)
	}

	// 计算 SHA256
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

// GenerateUniqueKey 生成节点唯一键
// 格式：type://server:port 或自定义格式
func GenerateUniqueKey(config map[string]any) string {
	proxyType, _ := config["type"].(string)
	server, _ := config["server"].(string)
	port, _ := config["port"].(float64)

	if proxyType == "" || server == "" {
		return ""
	}

	return fmt.Sprintf("%s://%s:%d", proxyType, server, int(port))
}

// EncodeConfig 将配置编码为 Base64 JSON 字符串
func EncodeConfig(config map[string]any) (string, error) {
	data, err := json.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("marshal config failed: %w", err)
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

// DecodeConfig 从 Base64 JSON 字符串解码配置
func DecodeConfig(encoded string) (map[string]any, error) {
	// 尝试多种 Base64 编码方式
	encodings := []*base64.Encoding{
		base64.StdEncoding,
		base64.RawStdEncoding,
		base64.URLEncoding,
		base64.RawURLEncoding,
	}

	var data []byte
	var err error
	for _, encoding := range encodings {
		data, err = encoding.DecodeString(encoded)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, fmt.Errorf("decode base64 failed: %w", err)
	}

	var config map[string]any
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config failed: %w", err)
	}

	return config, nil
}
