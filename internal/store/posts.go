package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Post struct {
	ID           int64    `json:"id"`
	Content      string   `json:"content"`
	Title        string   `json:"title"`
	UserID       int64    `json:"UserID"`
	Tags         []string `json:"tags"`
	CreationDate string   `json:"creation_date"`
	UpdateDate   string   `json:"update_date"`
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
		return nil
	}

	return nil

}
