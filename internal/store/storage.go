package store

import (
	"context"
	"database/sql"
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
	}

	Comments interface{
		GetCommentsByPostId (context.Context, int64) (*[]Comment, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostsStore{db},
		Users: &UsersStore{db},
		Comments: &CommentsStore{db},
	}
}
