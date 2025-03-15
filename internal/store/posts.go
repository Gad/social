package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var (
	ErrorNotFound = errors.New("record not found")
)

type Post struct {
	ID           int64     `json:"id"`
	Content      string    `json:"content"`
	Title        string    `json:"title"`
	UserID       int64     `json:"UserID"`
	Tags         []string  `json:"tags"`
	CreationDate string    `json:"creation_date"`
	UpdateDate   string    `json:"update_date"`
	Comments     []Comment `json:"comments"`
}

type PostsStore struct {
	db *sql.DB
}

func (s *PostsStore) Create(ctx context.Context, p *Post) error {
	query := `
	INSERT INTO posts(content, title, user_id, tags)
	VALUES ($1, $2, $3, $4) RETURNING id, creation_date, update_date
	`

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

func (s *PostsStore) GetPostById(ctx context.Context, postID int) (*Post, error) {
	query := `
	SELECT content, title, user_id, tags, creation_date, update_date 
	FROM posts 
	WHERE id = $1;
	`
	p := &Post{
		ID: int64(postID),
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
