package node_test

import (
	"testing"

	"github.com/darabuchi/prism/pkg/node"
	"github.com/lazygophers/log"
)

func TestGenerateId(t *testing.T) {
	config1 := map[string]any{
		"type":   "vmess",
		"server": "example.com",
		"port":   443,
		"uuid":   "test-uuid",
	}

	config2 := map[string]any{
		"type":   "vmess",
		"server": "example.com",
		"port":   443,
		"uuid":   "test-uuid",
		"name":   "不同的名字", // name 不应该影响 ID
	}

	id1 := node.GenerateId(config1)
	id2 := node.GenerateId(config2)

	if id1 != id2 {
		t.Error("相同配置（除 name 外）应该生成相同的 ID")
	}

	if len(id1) != 64 {
		t.Errorf("expected 64 character ID, got %d", len(id1))
	}

	log.Infof("Generated ID: %s", id1)
}

func TestGenerateIdWithUniqueId(t *testing.T) {
	config := map[string]any{
		"type":      "vmess",
		"server":    "example.com",
		"port":      443,
		"uuid":      "test-uuid",
		"unique_id": "custom-unique-id",
	}

	id := node.GenerateId(config)

	// 如果配置中有 unique_id，应该直接返回该值
	if id != "custom-unique-id" {
		t.Errorf("expected unique_id to be returned, got %s", id)
	}

	log.Infof("Using custom unique_id: %s", id)
}

func TestGenerateIdDifferentConfigs(t *testing.T) {
	config1 := map[string]any{
		"type":   "vmess",
		"server": "example1.com",
		"port":   443,
		"uuid":   "test-uuid",
	}

	config2 := map[string]any{
		"type":   "vmess",
		"server": "example2.com",
		"port":   443,
		"uuid":   "test-uuid",
	}

	id1 := node.GenerateId(config1)
	id2 := node.GenerateId(config2)

	// 不同的 server 应该生成不同的 ID
	if id1 == id2 {
		t.Error("不同配置应该生成不同的 ID")
	}

	log.Infof("ID1: %s", id1)
	log.Infof("ID2: %s", id2)
}

func TestGenerateIdWithDifferentPorts(t *testing.T) {
	config1 := map[string]any{
		"type":   "vmess",
		"server": "example.com",
		"port":   443,
		"uuid":   "test-uuid",
	}

	config2 := map[string]any{
		"type":   "vmess",
		"server": "example.com",
		"port":   8443,
		"uuid":   "test-uuid",
	}

	id1 := node.GenerateId(config1)
	id2 := node.GenerateId(config2)

	// 不同的 port 应该生成不同的 ID
	if id1 == id2 {
		t.Error("不同端口应该生成不同的 ID")
	}

	log.Infof("ID1 (port 443): %s", id1)
	log.Infof("ID2 (port 8443): %s", id2)
}

func TestGenerateIdExcludesMetadataFields(t *testing.T) {
	config1 := map[string]any{
		"type":   "vmess",
		"server": "example.com",
		"port":   443,
		"uuid":   "test-uuid",
	}

	config2 := map[string]any{
		"type":       "vmess",
		"server":     "example.com",
		"port":       443,
		"uuid":       "test-uuid",
		"name":       "节点名",
		"md5":        "some-md5",
		"sha256":     "some-sha256",
		"extra_info": "some-extra-info",
	}

	id1 := node.GenerateId(config1)
	id2 := node.GenerateId(config2)

	// 元数据字段不应该影响 ID
	if id1 != id2 {
		t.Error("元数据字段不应该影响 ID 生成")
	}

	log.Infof("ID (without metadata): %s", id1)
	log.Infof("ID (with metadata): %s", id2)
}

func TestGenerateIdWithoutPort(t *testing.T) {
	config := map[string]any{
		"type":   "vmess",
		"server": "example.com",
		// 没有 port 字段，应使用默认端口 8388
		"uuid": "test-uuid",
	}

	id := node.GenerateId(config)

	if len(id) != 64 {
		t.Errorf("expected 64 character ID, got %d", len(id))
	}

	log.Infof("ID (without port): %s", id)
}

func TestGenerateIdWithComplexTypes(t *testing.T) {
	// 测试 deepField 函数的所有类型分支
	config := map[string]any{
		"type":   "vmess",
		"server": "example.com",
		"port":   443,
		// string
		"string_field": "test-string",
		// []byte
		"bytes_field": []byte("test-bytes"),
		// int
		"int_field": 123,
		// int64
		"int64_field": int64(456),
		// []string
		"string_array": []string{"a", "b", "c"},
		// float64
		"float_field": 3.14,
		// bool
		"bool_field": true,
		// map[string]string
		"string_map": map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
		// []interface{}
		"interface_array": []interface{}{"x", 1, true},
		// map[string]any
		"nested_map": map[string]any{
			"nested_key": "nested_value",
		},
	}

	id := node.GenerateId(config)

	if len(id) != 64 {
		t.Errorf("expected 64 character ID, got %d", len(id))
	}

	log.Infof("ID (with complex types): %s", id)
}

func TestGenerateIdWithEmptyStringArray(t *testing.T) {
	config := map[string]any{
		"type":   "vmess",
		"server": "example.com",
		"port":   443,
		// 包含空字符串的数组
		"string_array": []string{"a", "", "b", ""},
	}

	id := node.GenerateId(config)

	if len(id) != 64 {
		t.Errorf("expected 64 character ID, got %d", len(id))
	}

	log.Infof("ID (with empty strings in array): %s", id)
}

func TestGenerateIdWithEmptyMapValues(t *testing.T) {
	config := map[string]any{
		"type":   "vmess",
		"server": "example.com",
		"port":   443,
		// 包含空值的 map
		"string_map": map[string]string{
			"key1": "value1",
			"key2": "", // 空值应该被跳过
		},
	}

	id := node.GenerateId(config)

	if len(id) != 64 {
		t.Errorf("expected 64 character ID, got %d", len(id))
	}

	log.Infof("ID (with empty map values): %s", id)
}

func TestGenerateIdWithBoolFalse(t *testing.T) {
	config := map[string]any{
		"type":       "vmess",
		"server":     "example.com",
		"port":       443,
		"bool_false": false,
	}

	id := node.GenerateId(config)

	if len(id) != 64 {
		t.Errorf("expected 64 character ID, got %d", len(id))
	}

	log.Infof("ID (with bool false): %s", id)
}

func TestGenerateIdWithNilValue(t *testing.T) {
	config := map[string]any{
		"type":      "vmess",
		"server":    "example.com",
		"port":      443,
		"nil_field": nil,
	}

	id := node.GenerateId(config)

	if len(id) != 64 {
		t.Errorf("expected 64 character ID, got %d", len(id))
	}

	log.Infof("ID (with nil value): %s", id)
}

func TestGenerateIdWithUnsupportedType(t *testing.T) {
	// 定义一个不支持的类型
	type unsupportedType struct {
		field string
	}

	config := map[string]any{
		"type":   "vmess",
		"server": "example.com",
		"port":   443,
		// 不支持的类型
		"unsupported_field": unsupportedType{field: "test"},
	}

	id := node.GenerateId(config)

	if len(id) != 64 {
		t.Errorf("expected 64 character ID, got %d", len(id))
	}

	log.Infof("ID (with unsupported type): %s", id)
}
