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

	log.Infof("生成的 ID: %s", id1)
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

	log.Infof("使用自定义 unique_id: %s", id)
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
