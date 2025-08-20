package redis

import (
	"YAccount/global"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/google/uuid"
)

type DistributedLock struct {
	client *redis.Client
	key    string
	value  string
	expire time.Duration
}

func NewDistributedLock(key string, expire time.Duration) *DistributedLock {
	return &DistributedLock{
		client: global.Redis,
		key:    key,
		value:  uuid.New().String(), // UUID
		expire: expire,
	}
}

func (l *DistributedLock) Acquire() error {
	ctx := context.Background()
	result := l.client.SetNX(ctx, l.key, l.value, l.expire)
	if result.Err() != nil {
		return result.Err()
	}
	if !result.Val() {
		return fmt.Errorf("获取锁失败")
	}
	return nil
}

func (l *DistributedLock) Release() error {
	// 使用Lua脚本保证原子性
	luaScript := `
        if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("del", KEYS[1])
        else
            return 0
        end
    `
	ctx := context.Background()
	result := l.client.Eval(ctx, luaScript, []string{l.key}, l.value)
	return result.Err()
}
