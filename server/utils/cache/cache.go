package cache

import (
	"context"

	"github.com/3086953492/YaBase/global"
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

// redisInstance 优雅地获取Redis实例
func redisInstance() *redis.Client {
	return global.GetGlobalRedis()
}

// cacheInstance 优雅地获取缓存实例
func cacheInstance() *cache.Cache {
	return global.GetGlobalCache()
}

func GetCacheKeysByPrefix(prefix string) ([]string, error) {
	var keys []string
	iter := redisInstance().Scan(context.Background(), 0, prefix+"*", 0).Iterator()
	for iter.Next(context.Background()) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	return keys, nil
}

func DeleteCacheKeysByPrefix(prefix string) error {
	keys, err := GetCacheKeysByPrefix(prefix)
	if err != nil {
		return err
	}
	for _, key := range keys {
		err := cacheInstance().Delete(context.Background(), key)
		if err != nil {
			return err
		}
	}
	return nil
}
