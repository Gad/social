package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gad/social/internal/store"
	"github.com/redis/go-redis/v9"
)

type UserStore struct {
	rdb *redis.Client
	ttl time.Duration
}

func (s *UserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	cacheKey := fmt.Sprintf("user-%v", userID)
	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err != nil {
		switch err {
		case redis.Nil:
			return nil,nil
		default:
			return nil, err
		}
	}
	user := &store.User{}
	if err := json.Unmarshal([]byte(data), user); err != nil {
		return nil, err
	}
	return user, nil


}

func (s *UserStore) Set(ctx context.Context, user *store.User) error {
	if user.ID == 0 {
		return fmt.Errorf("user ID is required")
	}
	cacheKey := fmt.Sprintf("user-%v", user.ID)
	json_data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return s.rdb.SetEx(ctx, cacheKey, json_data, s.ttl).Err()
}