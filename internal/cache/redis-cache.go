package cache

import "github.com/redis/go-redis/v9"

func NewRedisClient(add, pw string, db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     add,
		Password: pw,
		DB:       db,
	})

	return client
}
