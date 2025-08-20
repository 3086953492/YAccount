package cache

import (
	"YAccount/global"
	"context"
)

func GetCacheKeysByPrefix(prefix string) ([]string, error) {
	var keys []string
	iter := global.Redis.Scan(context.Background(), 0, prefix+"*", 0).Iterator()
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
		err := global.Cache.Delete(context.Background(), key)
		if err != nil {
			return err
		}
	}
	return nil
}