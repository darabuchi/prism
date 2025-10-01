package state

import (
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/storage/cache"
	"github.com/lazygophers/utils/atexit"
)

var (
	_cache cache.Cache
)

// ConnectCache 连接缓存
func ConnectCache() (err error) {
	log.Info("try connect cache")

	_cache, err = cache.New(State.Config.Cache)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// 注册退出时关闭缓存连接
	atexit.Register(func() {
		if _cache != nil {
			err = _cache.Close()
			if err != nil {
				log.Errorf("close cache error: %v", err)
			}
		}
	})

	log.Info("cache connected successfully")

	return nil
}

// Cache 返回缓存客户端
func Cache() cache.Cache {
	return _cache
}
