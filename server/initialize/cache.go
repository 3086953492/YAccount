package initialize

import (
	"YAccount/global"
	"time"

	"github.com/go-redis/cache/v9"
)

func InitCache() error {
	global.Cache = cache.New(&cache.Options{
		Redis:      global.Redis,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	var isInitCache bool

	err := global.Cache.Once(&cache.Item{
		Key:   "init",
		Value: &isInitCache,
		Do: func(*cache.Item) (any, error) {
			isInitCache = true
			return isInitCache, nil
		},
	})
	if err != nil || !isInitCache {
		return err
	}
	return nil
}
