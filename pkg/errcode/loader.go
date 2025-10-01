package errcode

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/xerrors"
	"gopkg.in/yaml.v3"
)

// Localizer 本地化器
type Localizer struct {
	messages map[string]map[string]string // locale -> key -> message
	mu       sync.RWMutex
	fallback string // 默认语言
}

var (
	defaultLocalizer *Localizer
	once             sync.Once
)

// GetLocalizer 获取默认本地化器
func GetLocalizer() *Localizer {
	once.Do(func() {
		defaultLocalizer = NewLocalizer("zh-CN")
	})
	return defaultLocalizer
}

// NewLocalizer 创建新的本地化器
func NewLocalizer(fallbackLocale string) *Localizer {
	return &Localizer{
		messages: make(map[string]map[string]string),
		fallback: fallbackLocale,
	}
}

// LoadFromFile 从 YAML 文件加载消息
func (l *Localizer) LoadFromFile(locale, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return xerrors.Errorf("failed to read file %s: %w", filePath, err)
	}

	return l.LoadFromBytes(locale, data)
}

// LoadFromBytes 从字节数据加载消息
func (l *Localizer) LoadFromBytes(locale string, data []byte) error {
	var rawData map[string]interface{}
	if err := yaml.Unmarshal(data, &rawData); err != nil {
		return xerrors.Errorf("failed to unmarshal yaml: %w", err)
	}

	messages := flattenMap(rawData, "")

	l.mu.Lock()
	defer l.mu.Unlock()

	l.messages[locale] = messages
	return nil
}

// LoadFromDir 从目录加载所有语言文件
func (l *Localizer) LoadFromDir(dirPath string) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return xerrors.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// 检查文件扩展名
		ext := filepath.Ext(entry.Name())
		if ext != ".yaml" && ext != ".yml" {
			continue
		}

		// 从文件名提取 locale（例如：zh-CN.yaml -> zh-CN）
		locale := strings.TrimSuffix(entry.Name(), ext)

		// 加载文件
		filePath := filepath.Join(dirPath, entry.Name())
		if err := l.LoadFromFile(locale, filePath); err != nil {
			return xerrors.Errorf("failed to load locale %s: %w", locale, err)
		}
	}

	return nil
}

// Get 获取本地化消息
func (l *Localizer) Get(locale, key string) string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// 尝试从指定 locale 获取
	if messages, ok := l.messages[locale]; ok {
		if msg, ok := messages[key]; ok {
			return msg
		}
	}

	// 尝试从 fallback locale 获取
	if locale != l.fallback {
		if messages, ok := l.messages[l.fallback]; ok {
			if msg, ok := messages[key]; ok {
				return msg
			}
		}
	}

	// 返回 key 本身作为后备
	return key
}

// Getf 获取本地化消息并格式化
func (l *Localizer) Getf(locale, key string, args ...interface{}) string {
	msg := l.Get(locale, key)
	return fmt.Sprintf(msg, args...)
}

// NewError 创建本地化错误
func (l *Localizer) NewError(locale string, code *ErrCode) *CodedError {
	message := l.Get(locale, code.Key)
	return code.New(message)
}

// NewErrorf 创建本地化错误（格式化）
func (l *Localizer) NewErrorf(locale string, code *ErrCode, args ...interface{}) *CodedError {
	message := l.Getf(locale, code.Key, args...)
	return code.New(message)
}

// WrapError 包装错误并本地化
func (l *Localizer) WrapError(locale string, code *ErrCode, err error) *CodedError {
	message := l.Get(locale, code.Key)
	return code.Wrap(err, message)
}

// WrapErrorf 包装错误并本地化（格式化）
func (l *Localizer) WrapErrorf(locale string, code *ErrCode, err error, args ...interface{}) *CodedError {
	message := l.Getf(locale, code.Key, args...)
	return code.Wrap(err, message)
}

// flattenMap 将嵌套的 map 展平为点分隔的键
func flattenMap(data map[string]interface{}, prefix string) map[string]string {
	result := make(map[string]string)

	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch v := value.(type) {
		case string:
			result[fullKey] = v
		case map[string]interface{}:
			// 递归处理嵌套 map
			nested := flattenMap(v, fullKey)
			for k, val := range nested {
				result[k] = val
			}
		default:
			// 其他类型转换为字符串
			result[fullKey] = fmt.Sprint(v)
		}
	}

	return result
}

// 便捷函数：使用默认本地化器

// LoadLocalizations 加载本地化文件
func LoadLocalizations(dirPath string) error {
	return GetLocalizer().LoadFromDir(dirPath)
}

// T 翻译（Translate 的简写）
func T(locale, key string) string {
	return GetLocalizer().Get(locale, key)
}

// Tf 翻译并格式化
func Tf(locale, key string, args ...interface{}) string {
	return GetLocalizer().Getf(locale, key, args...)
}

// NewLocalizedError 创建本地化错误
func NewLocalizedError(locale string, code *ErrCode) *CodedError {
	return GetLocalizer().NewError(locale, code)
}

// NewLocalizedErrorf 创建本地化错误（格式化）
func NewLocalizedErrorf(locale string, code *ErrCode, args ...interface{}) *CodedError {
	return GetLocalizer().NewErrorf(locale, code, args...)
}

// WrapLocalizedError 包装错误并本地化
func WrapLocalizedError(locale string, code *ErrCode, err error) *CodedError {
	return GetLocalizer().WrapError(locale, code, err)
}

// WrapLocalizedErrorf 包装错误并本地化（格式化）
func WrapLocalizedErrorf(locale string, code *ErrCode, err error, args ...interface{}) *CodedError {
	return GetLocalizer().WrapErrorf(locale, code, err, args...)
}
