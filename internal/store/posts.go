package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	ID           int64     `json:"id"`
	Content      string    `json:"content"`
	Title        string    `json:"title"`
	UserID       int64     `json:"UserID"`
	Tags         []string  `json:"tags"`
	Version      int       `json:"version"`
	CreationDate string    `json:"creation_date"`
	UpdateDate   string    `json:"update_date"`
	Comments     []Comment `json:"comments"`
	User         User      `json:"user"`
}


type PostsStore struct {
	db *sql.DB
}

func (s *PostsStore) Create(ctx context.Context, p *Post) error {
	query := `
	INSERT INTO posts(content, title, user_id, tags)
	VALUES ($1, $2, $3, $4) RETURNING id, creation_date, update_date
	`

	ctx, Cancel := context.WithTimeout(ctx, timeOutDuration)
	defer Cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		p.Content,
		p.Title,
		p.UserID,
		pq.Array(p.Tags),
	).Scan(
		&p.ID,
		&p.CreationDate,
		&p.UpdateDate,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *PostsStore) GetPostById(ctx context.Context, postID int64) (*Post, error) {
	query := `
	SELECT content, title, user_id, tags, version, creation_date, update_date 
	FROM posts 
	WHERE id = $1;
	`

	ctx, Cancel := context.WithTimeout(ctx, timeOutDuration)
	defer Cancel()

	p := &Post{
		ID: postID,
	}

	err := s.db.QueryRowContext(
		ctx,
		query,
		int64(postID),
	).Scan(
		&p.Content,
		&p.Title,
		&p.UserID,
		pq.Array(&p.Tags),
		&p.Version,
		&p.CreationDate,
		&p.UpdateDate,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}
	return p, nil
}

func (s *PostsStore) DeletePostById(ctx context.Context, postID int64) error {
	query := `
	DELETE FROM posts 
	WHERE id = $1;
	`
	ctx, Cancel := context.WithTimeout(ctx, timeOutDuration)
	defer Cancel()

	result, err := s.db.ExecContext(
		ctx,
		query,
		postID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	switch rows {
	case 0:
		return ErrorNotFound
	case 1:
		return nil
	default:
		return ErrorDeleteTooMany
	}

}

func (s *PostsStore) UpdatePostById(ctx context.Context, post *Post) error {
	query := `
	UPDATE posts
	SET title = $1, content =$2, version = version + 1		
	WHERE id = $3 AND version = $4
	RETURNING version;`

	ctx, Cancel := context.WithTimeout(ctx, timeOutDuration)
	defer Cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		int64(post.ID),
		post.Version,
	).Scan(&post.Version)
	if err != nil {

		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrorNotFound

		default:
			return err
		}

	}
	return nil

}
