package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisIPsStore struct {
	rdb *redis.Client
	ttl time.Duration
}

//		GetIPCount (context.Context, string) (int, error)

func (s *RedisIPsStore) IncrIPCount(ctx context.Context, IP string) (int64, error) {

	ipKey := fmt.Sprintf("ip-%v", IP)
	count, err := s.rdb.Incr(ctx, ipKey).Result()
	//log.Printf("INCR IP Key : %s - count : %d - err: %v \n", ipKey, count, err)
	if count == 1 {
		s.rdb.Expire(ctx, ipKey, s.ttl)
	}
	return count, err
}
