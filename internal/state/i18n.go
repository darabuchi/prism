package state

import (
	"embed"

	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/i18n"
	"github.com/lazygophers/lrpc/middleware/xerror"
	"github.com/lazygophers/utils/routine"
)

//go:embed ../../resource/localize/*.yaml
var i18nFs embed.FS

// LoadI18n 加载多语言资源
func LoadI18n() (err error) {
	State.I18n = i18n.DefaultI18n

	// 从嵌入的文件系统加载
	err = State.I18n.LoadLocalizesWithFs("resource/localize", i18nFs)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// 设置默认语言为英语
	State.I18n.SetDefaultLang("en")

	// 为 xerror 设置 i18n 支持
	xerror.SetI18n(i18n.NewI18nForXerror(State.I18n))

	// 设置 goroutine 语言传递
	routine.AddBeforeRoutine(func(baseGid, currentGid int64) {
		i18n.SetLanguage(i18n.GetLanguage(baseGid), currentGid)
	})
	routine.AddAfterRoutine(func(currentGid int64) {
		i18n.DelLanguage(currentGid)
	})

	return nil
}

// Localize 本地化消息
func Localize(key string, args ...interface{}) string {
	return State.I18n.Localize(key, args...)
}

// LocalizeWithLang 使用指定语言本地化消息
func LocalizeWithLang(lang string, key string, args ...interface{}) string {
	return State.I18n.LocalizeWithLang(lang, key, args...)
}
