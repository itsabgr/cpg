package ratelimit

import (
	"context"
	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"time"
)

var _ IRedisClient = &inMemory{}

type inMemory struct {
	done   *redis.BoolCmd
	failed *redis.BoolCmd
	cache  *cache.Cache
}

func NewInMemory(cleanupInterval time.Duration) IRedisClient {

	failed := &redis.BoolCmd{}
	done := &redis.BoolCmd{}
	done.SetVal(true)

	return &inMemory{done, failed, cache.New(cache.NoExpiration, cleanupInterval)}

}

func (m inMemory) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	if m.cache.Add(key, value, expiration) == nil {
		return m.done
	} else {
		return m.failed
	}
}
