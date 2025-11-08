package cache

import (
	"context"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/dgraph-io/badger/v4"
	"github.com/gad/social/internal/store"
	"github.com/redis/go-redis/v9"
)

type UserStorage struct {
	Users interface {
		GetUser(context.Context, int64) (*store.User, error)
		SetUser(context.Context, *store.User) error
	}
}

type IPStorage struct {
	IPs interface {
		GetIPCount (context.Context, string) (int, error)
		SetIPCount (context.Context, string, int) error 
	}
}

func NewRedisStorage(rdb *redis.Client, ttl time.Duration) UserStorage {
	return UserStorage{
		Users: &RedisUserStore{rdb: rdb, ttl: ttl},
	}
}


func NewRedisIPStorage(rdb *redis.Client, ttl time.Duration) IPStorage {
	return IPStorage{
		IPs: &RedisIPsStore{rdb: rdb, ttl: ttl},
	}
}
func NewBadgerStorage(bdb *badger.DB, ttl time.Duration) UserStorage {
	return UserStorage{
		Users: &BadgerUserStore{bdb: bdb, ttl: ttl},
	}
}

func NewMemcachedStorage(mdb *memcache.Client, ttl time.Duration) UserStorage {
	return UserStorage{
		Users: &MemcachedUserStore{mdb: mdb, ttl: ttl},
	}
}
