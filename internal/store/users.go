package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"_"`
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
