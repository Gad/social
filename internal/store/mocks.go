package store

import (
	"context"
	"database/sql"
	"time"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUsersStore{},
	}
}

type MockUsersStore struct{}

func (m *MockUsersStore) Create(ctx context.Context, tx *sql.Tx, u *User) error {
	return nil
}

func (m *MockUsersStore) GetUserById(ctx context.Context, id int64) (*User, error) {
	return &User{}, nil
}

func (m *MockUsersStore) Follow(ctx context.Context, userID int64, followID int64) error {
	return nil
}

func (m *MockUsersStore) Unfollow(ctx context.Context, userID int64, unfollowID int64) error {
	return nil
}

func (m *MockUsersStore) RegisterNew(ctx context.Context, u *User, token string, ttl time.Duration) error {
	return nil
}

func (m *MockUsersStore) Activate(ctx context.Context, token string) error {
	return nil
}

func (m *MockUsersStore) Delete(ctx context.Context, userID int64) error {
	return nil
}

func (m *MockUsersStore) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return &User{
		ID:    1,
		Email: "user@tests",
	}, nil
}
