package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var ErrorDuplicateEmail = errors.New("duplicate email address")
var ErrorDuplicateµUsername = errors.New("duplicate username")


type User struct {
	ID           int64    `json:"id"`
	Username     string   `json:"username"`
	Email        string   `json:"email"`
	Password     password `json:"-"`
	CreationDate string   `json:"creation_date"`
}

type password struct {
	plainText *string
	hash      []byte
}

func (p *password) Set(plaintext string) error {

	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.hash = hash
	p.plainText = &plaintext
	return nil
}

type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) Create(ctx context.Context, tx *sql.Tx, u *User) error {
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
		switch{
		case strings.Contains(err.Error(),"violates unique constraints") &&  strings.Contains(err.Error(),"users_email_key"):
			return ErrorDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraints "users_username_key"`:
			return ErrorDuplicateµUsername
		default:
			return err
		}
		
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

	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
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

func (s *UsersStore) RegisterNew(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	// Create and invite transaction
	// 1- transaction wrapper
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// 2- create the user
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		} 
		// 3- create the user invite
		if err := s.createUserInvitation(ctx, tx, token, user.ID, invitationExp); err != nil {
			return err
		}
		return nil
	})

}

func (s *UsersStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, userID int64, invitationExp time.Duration) error{
	query := `ÌNSERT INTO user_invitations (token, user_id, expiry) VALUES ($1, $2, $3)`
	ctx, cancel := context.WithTimeout(ctx, invitationExp)
	defer cancel()


	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(invitationExp) )
	if err != nil {
		return err
	}
	return nil
}