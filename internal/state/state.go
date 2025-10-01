package state

import (
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/i18n"
	"github.com/lazygophers/lrpc/middleware/storage/db"
	"github.com/lazygophers/utils/app"
	"github.com/pterm/pterm"
)

type state struct {
	Config *Config

	// I18n 多语言支持
	I18n *i18n.I18n
}

var State = new(state)

// Load 初始化全局状态
// 按顺序加载：配置 -> i18n -> 数据库 -> 缓存 -> 队列
func Load() (err error) {
	log.SetPrefixMsg(app.Name)

	// 1. 加载配置
	spinner, _ := pterm.DefaultSpinner.Start("正在加载配置......")
	err = LoadConfig()
	if err != nil {
		log.Errorf("err:%v", err)
		spinner.Fail("配置加载失败")
		pterm.Debug.Println(err)
		return err
	}
	spinner.Success("配置加载成功")

	// 2. 加载 i18n
	spinner, _ = pterm.DefaultSpinner.Start("正在加载多语言资源......")
	err = LoadI18n()
	if err != nil {
		log.Errorf("err:%v", err)
		spinner.Fail("多语言资源加载失败")
		pterm.Debug.Println(err)
		return err
	}
	spinner.Success("多语言资源加载成功")

	// 3. 设置数据库默认驱动
	db.DefaultDriver = State.Config.Db.Type

	// 4. 连接数据库
	spinner, _ = pterm.DefaultSpinner.Start("正在连接数据库......")
	err = ConnectDatabase()
	if err != nil {
		log.Errorf("err:%v", err)
		spinner.Fail("数据库连接失败")
		pterm.Debug.Println(err)
		return err
	}
	spinner.Success("数据库连接成功")

	// 5. 连接缓存
	spinner, _ = pterm.DefaultSpinner.Start("正在连接缓存......")
	err = ConnectCache()
	if err != nil {
		log.Errorf("err:%v", err)
		spinner.Fail("缓存连接失败")
		pterm.Debug.Println(err)
		return err
	}
	spinner.Success("缓存连接成功")

	// 6. 初始化队列
	spinner, _ = pterm.DefaultSpinner.Start("正在启动队列......")
	err = InitQueue()
	if err != nil {
		log.Errorf("err:%v", err)
		spinner.Fail("队列启动失败")
		pterm.Debug.Println(err)
		return err
	}
	spinner.Success("队列启动成功")

	return nil
}
