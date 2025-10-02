package collector

import (
	"net/http"
	"os"
	"time"

	"github.com/darabuchi/prism/pkg/rules"
)

// Handler 数据源处理器接口
type Handler interface {
	// Download 下载指定路径的规则数据
	Download(path string) ([]byte, error)

	// Parse 解析规则数据
	Parse(body []byte, action string) ([]rules.Rule, error)

	// NeedUpdate 判断缓存是否需要更新
	NeedUpdate(info os.FileInfo, cacheDays int) bool
}

// BaseHandler 基础处理器，提供通用功能
type BaseHandler struct {
	client    *http.Client
	cacheDays int
}

// NewBaseHandler 创建基础处理器
func NewBaseHandler(proxy string, cacheDays int) *BaseHandler {
	return &BaseHandler{
		client:    NewHTTPClient(proxy),
		cacheDays: cacheDays,
	}
}

// NeedUpdate 判断缓存是否需要更新
func (h *BaseHandler) NeedUpdate(info os.FileInfo, cacheDays int) bool {
	if cacheDays == 0 {
		cacheDays = h.cacheDays
	}

	age := time.Since(info.ModTime())
	maxAge := time.Duration(cacheDays) * 24 * time.Hour

	return age > maxAge
}
