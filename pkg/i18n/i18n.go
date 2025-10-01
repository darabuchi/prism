package i18n

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// Localizer 实现多语言本地化
type Localizer struct {
	// messages 存储 locale -> code -> message 的映射
	messages map[string]map[int32]string
	mu       sync.RWMutex
	fallback string // 默认语言
}

// NewLocalizer 创建新的本地化器
func NewLocalizer() *Localizer {
	return &Localizer{
		messages: make(map[string]map[int32]string),
		fallback: "en", // 默认使用英语
	}
}

// SetFallback 设置回退语言
func SetFallback(locale string) {
	defaultLocalizer.mu.Lock()
	defer defaultLocalizer.mu.Unlock()
	defaultLocalizer.fallback = locale
}

// LoadFromDir 从目录加载所有语言文件
func (l *Localizer) LoadFromDir(dirPath string) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := filepath.Ext(entry.Name())
		if ext != ".yaml" && ext != ".yml" {
			continue
		}

		// 从文件名提取语言代码（如 zh-CN.yaml -> zh-CN）
		locale := strings.TrimSuffix(entry.Name(), ext)
		filePath := filepath.Join(dirPath, entry.Name())

		if err := l.LoadFromFile(locale, filePath); err != nil {
			return fmt.Errorf("failed to load locale %s: %w", locale, err)
		}
	}

	return nil
}

// LoadFromFile 从文件加载特定语言
func (l *Localizer) LoadFromFile(locale, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	var content map[string]interface{}
	if err := yaml.Unmarshal(data, &content); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	// 扁平化嵌套结构
	messages := make(map[int32]string)
	l.flattenMap("", content, messages)

	l.mu.Lock()
	l.messages[locale] = messages
	l.mu.Unlock()

	return nil
}

// flattenMap 将嵌套的 map 扁平化为 "key.subkey" 格式，并映射到错误码
func (l *Localizer) flattenMap(prefix string, m map[string]interface{}, result map[int32]string) {
	for k, v := range m {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "." + k
		}

		switch val := v.(type) {
		case map[string]interface{}:
			l.flattenMap(fullKey, val, result)
		case string:
			// 将键映射到错误码（需要错误码映射表）
			if code, ok := keyToCodeMap[fullKey]; ok {
				result[code] = val
			}
		}
	}
}

// Localize 实现 xerror.I18n 接口
// 根据错误码和语言获取本地化消息
func (l *Localizer) Localize(code int32, langs ...string) (string, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// 尝试每个提供的语言
	for _, lang := range langs {
		if messages, ok := l.messages[lang]; ok {
			if msg, found := messages[code]; found {
				return msg, true
			}
		}

		// 尝试去掉区域代码（如 zh-CN -> zh）
		if idx := strings.Index(lang, "-"); idx > 0 {
			baseLang := lang[:idx]
			if messages, ok := l.messages[baseLang]; ok {
				if msg, found := messages[code]; found {
					return msg, true
				}
			}
		}
	}

	// 回退到默认语言
	if messages, ok := l.messages[l.fallback]; ok {
		if msg, found := messages[code]; found {
			return msg, true
		}
	}

	return "", false
}

// 默认的全局本地化器
var defaultLocalizer = NewLocalizer()

// LoadFromDir 使用默认本地化器加载
func LoadFromDir(dirPath string) error {
	return defaultLocalizer.LoadFromDir(dirPath)
}

// Localize 使用默认本地化器
func Localize(code int32, langs ...string) (string, bool) {
	return defaultLocalizer.Localize(code, langs...)
}
