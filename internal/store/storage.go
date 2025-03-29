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
	timeOutDuration = time.Second * 5
	ErrorDuplicateKey = errors.New("duplicate key violates sql unique constraint")
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetPostById(context.Context, int64) (*Post, error)
		DeletePostById(context.Context, int64) error
		UpdatePostById(context.Context, *Post) error
		
	}

	Users interface {
		Create(context.Context, *User) error
		GetUserById(context.Context, int64) (*User, error)
		Follow(context.Context,  int64,  int64) error
		Unfollow( context.Context,  int64, int64) error
		RegisterNew(context.Context, *User, string) error
	}

	

	Comments interface{
		GetCommentsByPostId (context.Context, int64) (*[]Comment, error)
		Create (context.Context, *Comment) error 
	}

	Feeds interface{
		GetUserDefaultFeed(context.Context, int64, FeedPaginationQuery) ([]PostWtMetadata, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostsStore{db},
		Users: &UsersStore{db},
		Comments: &CommentsStore{db},
		Feeds : &FeedsStore{db},
	}
}
