package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

)

var (
	ErrorNotFound      = errors.New("record not found")
	ErrorDeleteTooMany = errors.New("expected to affect 1 row, affected more")
	timeOutDuration    = time.Second * 5
	ErrorDuplicateKey  = errors.New("duplicate key violates sql unique constraint")
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetPostById(context.Context, int64) (*Post, error)
		DeletePostById(context.Context, int64) error
		UpdatePostById(context.Context, *Post) error
	}

	Users interface {
		Create(context.Context, *sql.Tx, *User) error
		GetUserById(context.Context, int64) (*User, error)
		Follow(context.Context, int64, int64) error
		Unfollow(context.Context, int64, int64) error
		RegisterNew(context.Context, *User, string, time.Duration) error
		Activate(context.Context, string) error
		DeleteInvitation(context.Context, int64, *sql.Tx) error
	}

	Comments interface {
		GetCommentsByPostId(context.Context, int64) (*[]Comment, error)
		Create(context.Context, *Comment) error
	}

	Feeds interface {
		GetUserDefaultFeed(context.Context, int64, FeedPaginationQuery) ([]PostWtMetadata, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostsStore{db},
		Users:    &UsersStore{db},
		Comments: &CommentsStore{db},
		Feeds:    &FeedsStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return err2
		}
		return err

	}

	return tx.Commit()
}
