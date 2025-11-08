package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisIPsStore struct {
	rdb *redis.Client
	ttl time.Duration
}
//		GetIPCount (context.Context, string) (int, error)

func (s *RedisIPsStore) GetIPCount(ctx context.Context, IP string) (int, error) {
	cacheKey := fmt.Sprintf("ip-%v", IP)
	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err != nil {
		switch err {
		case redis.Nil:
			return 0, nil
		default:
			return 0, err
		}
	}
	count, err := strconv.ParseInt(data,10, 0)
	if err != nil {
		return 0, err
	}
	return int(count), nil

}



func (s *RedisIPsStore) SetIPCount(ctx context.Context, IP string, count int) error {
	if IP == "" {
		return fmt.Errorf("IP is required")
	}
	cacheKey := fmt.Sprintf("ip-%v", IP)
	data := strconv.Itoa(count)
	
	return s.rdb.SetEx(ctx, cacheKey, data, s.ttl).Err()
}
