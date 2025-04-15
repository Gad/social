package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gad/social/internal/store"
)

type MemcachedUserStore struct {
	mdb *memcache.Client
	ttl time.Duration
}

func (s *MemcachedUserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	cacheKey := fmt.Sprintf("user-%v", userID)
	data, err := s.mdb.Get(cacheKey)
	if err != nil {
		switch err {
		case memcache.ErrCacheMiss:
			return nil, nil
		default:
			return nil, err
		}
	}
	user := &store.User{}
	if err := json.Unmarshal([]byte(data.Value), user); err != nil {
		return nil, err
	}
	return user, nil

}

func (s *MemcachedUserStore) Set(ctx context.Context, user *store.User) error {
	if user.ID == 0 {
		return fmt.Errorf("user ID is required")
	}
	
	cacheKey := fmt.Sprintf("user-%v", user.ID)
	ttlSecondsInt := int32(s.ttl.Seconds())
	
	
	json_data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return s.mdb.Set(&memcache.Item{
		Key:        cacheKey,
		Value:      json_data,
		Expiration: ttlSecondsInt,
	})
}
