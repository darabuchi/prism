package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

// ErrorCodeRegistry 错误码注册表
type ErrorCodeRegistry struct {
	Version     string     `json:"version"`
	Description string     `json:"description"`
	Categories  []Category `json:"categories"`
	Errors      []Error    `json:"errors"`
}

// Category 错误码分类
type Category struct {
	Name        string `json:"name"`
	Range       string `json:"range"`
	Description string `json:"description"`
}

// Error 错误定义
type Error struct {
	Code        int32  `json:"code"`
	Key         string `json:"key"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

var (
	inputFile  string
	outputFile string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gen-error-code",
		Short: "根据 error_code.json 生成 error_code.gen.go",
		Long:  "从 error_code.json 读取错误码定义，自动生成 Go 常量文件 error_code.gen.go",
		Run:   runGenerate,
	}

	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "error_code.json", "输入的错误码 JSON 文件")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "error_code.gen.go", "输出的 Go 文件")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runGenerate(cmd *cobra.Command, args []string) {
	log.Infof("开始生成错误码文件...")
	log.Infof("  输入文件: %s", inputFile)
	log.Infof("  输出文件: %s", outputFile)

	// 读取 JSON 文件
	data, err := os.ReadFile(inputFile)
	if err != nil {
		log.Errorf("读取输入文件失败: %v", err)
		os.Exit(1)
	}

	// 解析 JSON
	var registry ErrorCodeRegistry
	if err := json.Unmarshal(data, &registry); err != nil {
		log.Errorf("解析 JSON 失败: %v", err)
		os.Exit(1)
	}

	log.Infof("成功解析错误码注册表:")
	log.Infof("  版本: %s", registry.Version)
	log.Infof("  分类数: %d", len(registry.Categories))
	log.Infof("  错误数: %d", len(registry.Errors))

	// 验证错误码
	if err := validateRegistry(&registry); err != nil {
		log.Errorf("错误码验证失败: %v", err)
		os.Exit(1)
	}

	// 生成 Go 代码
	output, err := generateGoCode(&registry)
	if err != nil {
		log.Errorf("生成 Go 代码失败: %v", err)
		os.Exit(1)
	}

	// 写入文件
	if err := os.WriteFile(outputFile, []byte(output), 0644); err != nil {
		log.Errorf("写入输出文件失败: %v", err)
		os.Exit(1)
	}

	log.Infof("成功生成错误码文件: %s", outputFile)
}

// validateRegistry 验证错误码注册表
func validateRegistry(registry *ErrorCodeRegistry) error {
	codeMap := make(map[int32]bool)
	keyMap := make(map[string]bool)

	for _, err := range registry.Errors {
		// 检查重复的错误码
		if codeMap[err.Code] {
			return fmt.Errorf("重复的错误码: %d", err.Code)
		}
		codeMap[err.Code] = true

		// 检查重复的 key
		if keyMap[err.Key] {
			return fmt.Errorf("重复的错误 key: %s", err.Key)
		}
		keyMap[err.Key] = true

		// 验证 key 格式 (category.name)
		parts := strings.Split(err.Key, ".")
		if len(parts) != 2 {
			return fmt.Errorf("错误 key 格式无效: %s (应为 category.name)", err.Key)
		}

		// 验证 category 存在
		found := false
		for _, cat := range registry.Categories {
			if cat.Name == err.Category {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("错误 %s 的分类不存在: %s", err.Key, err.Category)
		}
	}

	log.Infof("错误码验证通过")
	return nil
}

// generateGoCode 生成 Go 代码
func generateGoCode(registry *ErrorCodeRegistry) (string, error) {
	tmpl := template.Must(template.New("error_code").Funcs(template.FuncMap{
		"toConstName": toConstName,
		"toLower":     strings.ToLower,
		"toTitle":     strings.Title,
	}).Parse(codeTemplate))

	var buf strings.Builder
	data := map[string]interface{}{
		"Version":     registry.Version,
		"Description": registry.Description,
		"GeneratedAt": time.Now().Format("2006-01-02 15:04:05"),
		"Categories":  registry.Categories,
		"Errors":      registry.Errors,
		"PackageName": getPackageName(outputFile),
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("执行模板失败: %w", err)
	}

	return buf.String(), nil
}

// toConstName 将 key 转换为 Go 常量名
// 例如: "common.success" -> "ErrCommonSuccess"
//      "auth.unauthorized" -> "ErrAuthUnauthorized"
//      "geoip.invalid_ip" -> "ErrGeoipInvalidIp"
func toConstName(key string) string {
	parts := strings.Split(key, ".")
	if len(parts) != 2 {
		return "ErrUnknown"
	}

	category := parts[0]
	name := parts[1]

	// 构建常量名: Err + Category + Name
	var result strings.Builder
	result.WriteString("Err")

	// 处理 category（首字母大写）
	result.WriteString(strings.ToUpper(category[:1]))
	if len(category) > 1 {
		result.WriteString(category[1:])
	}

	// 处理 name（每个单词首字母大写）
	words := strings.Split(name, "_")
	for _, word := range words {
		if word == "" {
			continue
		}
		result.WriteString(strings.ToUpper(word[:1]))
		if len(word) > 1 {
			result.WriteString(word[1:])
		}
	}

	return result.String()
}

// getPackageName 从文件路径获取包名
func getPackageName(filePath string) string {
	// 检查是否是根目录的文件
	dir := filepath.Dir(filePath)

	// 如果是相对路径 ../../../error_code.gen.go，计算实际目录
	if strings.Contains(filePath, "../") {
		absPath, err := filepath.Abs(filePath)
		if err == nil {
			dir = filepath.Dir(absPath)
		}
	}

	// 获取目录名作为包名
	packageName := filepath.Base(dir)

	// 如果包名为 "." 或空，默认使用 "prism"
	if packageName == "." || packageName == "" || packageName == "/" {
		return "prism"
	}

	return packageName
}

const codeTemplate = `// Code generated by gen-error-code. DO NOT EDIT.
// Version: {{.Version}}
// Description: {{.Description}}
// Generated at: {{.GeneratedAt}}

package {{.PackageName}}

// 错误码定义
// 所有错误码都在 error_code.json 中注册
//
// 分类说明:
{{- range .Categories}}
//   {{.Name}}: {{.Description}} ({{.Range}})
{{- end}}
const (
{{- range .Errors}}
	// {{toConstName .Key}} {{.Description}} ({{.Code}})
	{{toConstName .Key}} int32 = {{.Code}}
{{- end}}
)

// ErrorCodeInfo 错误码信息
type ErrorCodeInfo struct {
	Code        int32
	Key         string
	Category    string
	Description string
}

// errorCodeMap 错误码映射表
var errorCodeMap = map[int32]ErrorCodeInfo{
{{- range .Errors}}
	{{.Code}}: {
		Code:        {{.Code}},
		Key:         "{{.Key}}",
		Category:    "{{.Category}}",
		Description: "{{.Description}}",
	},
{{- end}}
}

// GetErrorCodeInfo 获取错误码信息
func GetErrorCodeInfo(code int32) (ErrorCodeInfo, bool) {
	info, ok := errorCodeMap[code]
	return info, ok
}

// IsValidErrorCode 检查错误码是否有效
func IsValidErrorCode(code int32) bool {
	_, ok := errorCodeMap[code]
	return ok
}
`
