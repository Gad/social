package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
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
	Activated    bool     `json:"activated"`
	Role         Role     `json:"role"`
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

func (p *password) Compare(passwordS string) error {
	/* log.Printf("hashed %v - non-hashed%v", p.hash, passwordS)
	h, _ := bcrypt.GenerateFromPassword([]byte(passwordS), bcrypt.DefaultCost)
	log.Printf("hashed of passwordS %v", h) */
	return bcrypt.CompareHashAndPassword(p.hash, []byte(passwordS))
}

type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) Create(ctx context.Context, tx *sql.Tx, u *User) error {
	query := `
	INSERT INTO users(email,username, password, role_id)
	VALUES ($1, $2, $3, (SELECT id from roles where name = $4)) RETURNING id, creation_date;
	`

	ctx, cancel := context.WithTimeout(ctx, timeOutDuration)
	defer cancel()

	// make sure the role is filled with a default value
	if u.Role.Name == "" {
		u.Role.Name = "user"
	}

	err := s.db.QueryRowContext(
		ctx,
		query,
		u.Email,
		u.Username,
		// hash password
		u.Password.hash,
		u.Role.Name,
	).Scan(
		&u.ID,
		&u.CreationDate,
	)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "violates unique constraints") && strings.Contains(err.Error(), "users_email_key"):
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
	SELECT users.id, username, email, password, creation_date, roles.* 
	FROM users 
	JOIN roles ON roles.id = users.role_id
	WHERE users.id = $1;
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
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password.hash,
		&u.CreationDate,
		&u.Role.ID,
		&u.Role.Name,
		&u.Role.Level,
		&u.Role.Description,
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

func (s *UsersStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, userID int64, invitationExp time.Duration) error {
	query := `INSERT INTO user_invitations (token, user_id, expiry) VALUES ($1, $2, $3)`
	ctx, cancel := context.WithTimeout(ctx, invitationExp)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(invitationExp))
	if err != nil {
		return err
	}
	return nil
}

func (s *UsersStore) Activate(ctx context.Context, token string) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// find user_id related to token
		ctx, Cancel := context.WithTimeout(ctx, timeOutDuration)
		defer Cancel()

		user_id, err := s.getUserByToken(ctx, tx, token)

		if err != nil {
			return err
		}

		// activate user in users table
		if err = s.toggleUserActivated(ctx, user_id, tx); err != nil {
			return err
		}
		// remove user invitation

		if err = s.deleteInvitation(ctx, user_id, tx); err != nil {
			return err
		}
		return nil
	})
}

func (s *UsersStore) deleteInvitation(ctx context.Context, userID int64, tx *sql.Tx) error {
	query := `
	DELETE FROM user_invitations 
	WHERE user_id = $1;
	`

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	return nil

}

func (s *UsersStore) Delete(ctx context.Context, userID int64) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.delete(ctx, userID, tx); err != nil {
			return err
		}

		if err := s.deleteInvitation(ctx, userID, tx); err != nil {
			return err
		}
		return nil
	})

}

func (s *UsersStore) delete(ctx context.Context, userID int64, tx *sql.Tx) error {
	query := `DELETE FROM users 
	WHERE id = $1;
	`
	ctx, Cancel := context.WithTimeout(ctx, timeOutDuration)
	defer Cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UsersStore) toggleUserActivated(ctx context.Context, user_id int64, tx *sql.Tx) error {
	query := `
	UPDATE users 
	SET activated=true 
	WHERE id=$1;
	`

	_, err := tx.ExecContext(ctx, query, user_id)
	if err != nil {
		return err
	}
	return nil
}

func (s *UsersStore) getUserByToken(ctx context.Context, tx *sql.Tx, token string) (int64, error) {
	query := `
	SELECT user_id 
	FROM user_invitations 
	WHERE token=$1 AND expiry > $2;
	`

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	var user_id int64

	err := tx.QueryRowContext(
		ctx,
		query,
		hashToken,
		time.Now(),
	).Scan(
		&user_id,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return 0, ErrorNotFound
		default:
			return 0, err
		}
	}
	return user_id, nil
}

func (s *UsersStore) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, username, password
	FROM users
	WHERE email = $1 and activated = true`

	user := User{}

	ctx, Cancel := context.WithTimeout(ctx, timeOutDuration)
	defer Cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Password.hash,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}
