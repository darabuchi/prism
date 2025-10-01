package errcode

import (
	"fmt"

	"golang.org/x/xerrors"
)

// ErrCode 错误码结构
type ErrCode struct {
	Code    int    // 错误码（唯一标识）
	Key     string // 错误码键名（用于 i18n 查找）
	HTTPCode int    // HTTP 状态码
}

// Error 实现 error 接口
func (e *ErrCode) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Key)
}

// CodedError 包含错误码和上下文的错误
type CodedError struct {
	ErrCode *ErrCode               // 错误码
	Message string                 // 错误消息（已本地化）
	Cause   error                  // 原始错误
	Details map[string]interface{} // 错误详情
}

// Error 实现 error 接口
func (e *CodedError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%d] %s: %v", e.ErrCode.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%d] %s", e.ErrCode.Code, e.Message)
}

// Unwrap 返回原始错误，支持 errors.Is 和 errors.As
func (e *CodedError) Unwrap() error {
	return e.Cause
}

// Format 实现 fmt.Formatter 接口，支持 %+v 格式化输出错误链
func (e *CodedError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			// %+v: 详细格式，包含错误链
			fmt.Fprintf(s, "[%d] %s", e.ErrCode.Code, e.Message)
			if len(e.Details) > 0 {
				fmt.Fprintf(s, " (details: %+v)", e.Details)
			}
			if e.Cause != nil {
				fmt.Fprintf(s, "\nCaused by: %+v", e.Cause)
			}
			return
		}
		fallthrough
	case 's':
		fmt.Fprint(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

// Code 获取错误码
func (e *CodedError) Code() int {
	return e.ErrCode.Code
}

// HTTPStatus 获取 HTTP 状态码
func (e *CodedError) HTTPStatus() int {
	return e.ErrCode.HTTPCode
}

// WithCause 添加原始错误
func (e *CodedError) WithCause(cause error) *CodedError {
	e.Cause = cause
	return e
}

// WithDetail 添加错误详情
func (e *CodedError) WithDetail(key string, value interface{}) *CodedError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// WithDetails 批量添加错误详情
func (e *CodedError) WithDetails(details map[string]interface{}) *CodedError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	for k, v := range details {
		e.Details[k] = v
	}
	return e
}

// New 创建新的错误码错误
func (ec *ErrCode) New(message string) *CodedError {
	return &CodedError{
		ErrCode: ec,
		Message: message,
	}
}

// Newf 创建新的错误码错误（格式化消息）
func (ec *ErrCode) Newf(format string, args ...interface{}) *CodedError {
	return &CodedError{
		ErrCode: ec,
		Message: fmt.Sprintf(format, args...),
	}
}

// Wrap 包装现有错误
func (ec *ErrCode) Wrap(err error, message string) *CodedError {
	return &CodedError{
		ErrCode: ec,
		Message: message,
		Cause:   err,
	}
}

// Wrapf 包装现有错误（格式化消息）
func (ec *ErrCode) Wrapf(err error, format string, args ...interface{}) *CodedError {
	return &CodedError{
		ErrCode: ec,
		Message: fmt.Sprintf(format, args...),
		Cause:   err,
	}
}

// Is 判断错误是否为指定错误码（支持错误链）
func Is(err error, code *ErrCode) bool {
	if err == nil || code == nil {
		return false
	}

	// 遍历错误链
	for err != nil {
		if codedErr, ok := err.(*CodedError); ok {
			if codedErr.ErrCode.Code == code.Code {
				return true
			}
		}

		// 继续检查错误链
		err = xerrors.Unwrap(err)
	}

	return false
}

// GetCode 从错误中提取错误码（支持错误链）
func GetCode(err error) int {
	if err == nil {
		return 0
	}

	// 遍历错误链，找到第一个 CodedError
	for err != nil {
		if codedErr, ok := err.(*CodedError); ok {
			return codedErr.Code()
		}
		err = xerrors.Unwrap(err)
	}

	return 0
}

// GetHTTPStatus 从错误中提取 HTTP 状态码（支持错误链）
func GetHTTPStatus(err error) int {
	if err == nil {
		return 200
	}

	// 遍历错误链，找到第一个 CodedError
	for err != nil {
		if codedErr, ok := err.(*CodedError); ok {
			return codedErr.HTTPStatus()
		}
		err = xerrors.Unwrap(err)
	}

	return 500
}
