package state

import (
	"path/filepath"

	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/storage/cache"
	"github.com/lazygophers/lrpc/middleware/storage/db"
	"github.com/lazygophers/utils/config"
	"github.com/lazygophers/utils/validator"
)

type Config struct {
	Db    *db.Config    `json:"db,omitempty" yaml:"db,omitempty" toml:"db,omitempty" validate:"required"`
	Cache *cache.Config `json:"cache,omitempty" yaml:"cache,omitempty" toml:"cache,omitempty"`
}

// apply 应用默认配置
func (p *Config) apply() {
	if p.Db == nil {
		p.Db = &db.Config{
			Type:     db.Sqlite,
			Debug:    false,
			Address:  filepath.Join("resource"),
			Port:     0,
			Name:     "prism",
			Username: "prism",
			Password: "prism2025",
			Extras:   nil,
			Logger:   nil,
		}
	}

	if p.Cache == nil {
		p.Cache = &cache.Config{
			Type:    cache.Bbolt,
			DataDir: "resource/cache/",
		}
	}
}

// LoadConfig 加载配置
func LoadConfig() (err error) {
	State.Config = new(Config)

	// 跳过验证加载配置（允许部分字段为空）
	_ = config.LoadConfigSkipValidate(State.Config)

	// 应用默认值
	State.Config.apply()

	// 验证配置
	err = validator.Struct(State.Config)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
