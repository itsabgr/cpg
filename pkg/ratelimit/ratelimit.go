package ratelimit

import (
	"context"
	"github.com/itsabgr/ge"
	"github.com/redis/go-redis/v9"
	"time"
)

func NewRateLimit(redisClient IRedisClient) *RateLimit {
	return &RateLimit{redisClient: redisClient}
}

type IRedisClient interface {
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
}

type RateLimit struct {
	redisClient IRedisClient
}

func (rl *RateLimit) Limit(ctx context.Context, invoice string, duration time.Duration) (bool, error) {
	if len(invoice) == 0 {
		return false, ge.New("invalid id")
	}
	return rl.redisClient.SetNX(ctx, invoice, 1, duration).Result()
}
