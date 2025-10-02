package collector

import (
	"os"

	"github.com/darabuchi/prism/pkg/rules"
)

// Handler 数据源处理器接口
type Handler interface {
	// Download 下载指定路径的规则数据
	Download(path string) ([]byte, error)

	// Parse 解析规则数据
	Parse(body []byte, action string) ([]rules.Rule, error)

	// NeedUpdate 判断缓存是否需要更新
	NeedUpdate(info os.FileInfo) bool
}
