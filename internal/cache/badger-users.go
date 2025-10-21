package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/gad/social/internal/store"
)

type BadgerUserStore struct {
	bdb *badger.DB
	ttl time.Duration
}

func (s *BadgerUserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	cacheKey := fmt.Sprintf("user-%v", userID)
	var data string
	err := s.bdb.View(func(txn *badger.Txn) error {

		item, err := txn.Get([]byte(cacheKey))
		if err != nil {
			return err
		}
		valCopy, err := item.ValueCopy(nil)

		if err != nil {
			return err
		}
		data = string(valCopy)
		return nil
	})

	if err != nil {
		switch err {
		case badger.ErrKeyNotFound:
			return nil, nil
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

func (s *BadgerUserStore) Set(ctx context.Context, user *store.User) error {
	if user.ID == 0 {
		return fmt.Errorf("user ID is required")
	}
	cacheKey := fmt.Sprintf("user-%v", user.ID)
	json_data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return s.bdb.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(cacheKey), []byte(json_data)).WithTTL(s.ttl)
		err := txn.SetEntry(e)
		return err
	})
}
