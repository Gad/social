package cache

import (
	"context"
	"time"

	"github.com/gad/social/internal/store"
	"github.com/redis/go-redis/v9"
)


type Storage struct {
	Users interface {
		Get(context.Context, int64) (*store.User, error)
		Set(context.Context, *store.User) error
	}
	
}

func NewRedisStorage(rdb *redis.Client, ttl time.Duration) Storage {
	return Storage{
		Users: &UserStore{rdb : rdb, ttl: ttl},
	}
}