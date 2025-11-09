package ratelimiter

import (
	"context"
	"sync"
	"time"

	"github.com/gad/social/internal/cache"
)

type RedisFixedWindowRateLimiter struct {
	sync.RWMutex
	store  cache.IPStorage
	limit  int
	window time.Duration
}

func NewRedisFixedWindowLimiter(limit int, window time.Duration, store cache.IPStorage) *RedisFixedWindowRateLimiter {
	return &RedisFixedWindowRateLimiter{
		store:  store,
		limit:  limit,
		window: window,
	}
}

func (rl *RedisFixedWindowRateLimiter) Allow(ip string) (bool, time.Duration, error) {

	ctx := context.Background()

	count, err := rl.store.IPs.IncrIPCount(ctx, ip)

	//log.Printf("IP Key : %s - count : %d - err: %v \n", ip, count, err)

	if err != nil {
		return false, rl.window, err
	}

	if count <= int64(rl.limit) {

		return true, 0, nil
	}

	return false, rl.window, nil

}
