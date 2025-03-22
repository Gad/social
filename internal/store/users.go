package store

import (
	"context"
	"database/sql"
	"errors"
	

	"github.com/lib/pq"
)

type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"-"`
	CreationDate string `json:"creation_date"`
}

type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) Create(ctx context.Context, u *User) error {
	query := `
	INSERT INTO users(username, email, password)
	VALUES ($1, $2, $3) RETURNING id, creation_date
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		u.Username,
		u.Email,
		u.Password,
	).Scan(
		&u.ID,
		&u.CreationDate,
	)

	if err != nil {
		return err
	}

	return nil

}




func (s *UsersStore) Follow(ctx context.Context, followedID int64, followerID int64) error {
	query := `
	INSERT INTO followers(user_id, follower_id)
	VALUES ($1, $2);
	`
	ctx, Cancel := context.WithTimeout(ctx, timeOutDuration)
	defer Cancel()

	_, err := s.db.ExecContext(
		ctx,
		query,
		followedID,
		followerID,
	)

	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name()=="unique_violation"{
		return ErrorDuplicateKey
	}

	return err
}

func (s *UsersStore) Unfollow(ctx context.Context, followedID int64, followerID int64) error {
	query := `
	DELETE FROM followers
	WHERE user_id = $1 AND follower_id = $2;
	`
	ctx, Cancel := context.WithTimeout(ctx, timeOutDuration)
	defer Cancel()

	_, err := s.db.ExecContext(
		ctx,
		query,
		followedID,
		followerID,
	)

	return err
}

func (s *UsersStore) GetUserById(ctx context.Context, userID int64) (*User, error) {
	query := `
	SELECT username, email, password, creation_date 
	FROM users 
	WHERE id = $1;
	`

	ctx, Cancel := context.WithTimeout(ctx, timeOutDuration)
	defer Cancel()

	u := &User{
		ID: userID,
	}

	err := s.db.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&u.Username,
		&u.Email,
		&u.Password,
		&u.CreationDate,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}
	return u, nil
}
