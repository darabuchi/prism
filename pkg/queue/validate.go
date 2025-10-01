package queue

import (
	"fmt"

	"github.com/lazygophers/utils/validator"
)

var defaultValidator *validator.Validator

func init() {
	// 初始化默认验证器
	v, err := validator.New(
		validator.WithLocale("zh"),
		validator.WithUseJSON(true),
	)
	if err != nil {
		// 降级到默认验证器
		defaultValidator = validator.Default()
	} else {
		defaultValidator = v
	}
}

// ValidateStruct 验证结构体
func ValidateStruct(s interface{}) error {
	if err := defaultValidator.Struct(s); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}
