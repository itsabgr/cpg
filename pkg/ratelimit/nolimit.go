package ratelimit

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

func NoLimit() IRedisClient {
	result := &redis.BoolCmd{}
	result.SetVal(true)
	return (*noLimit)(result)
}

type noLimit redis.BoolCmd

func (n *noLimit) SetNX(_ context.Context, _ string, _ interface{}, _ time.Duration) *redis.BoolCmd {
	return (*redis.BoolCmd)(n)
}
