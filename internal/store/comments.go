package store

import (
	"context"
	"database/sql"
)

type CommentsStore struct {
	db *sql.DB
}

type Comment struct {
	ID           int64  `json:"id"`
	PostID       int64  `json:"post_id"`
	UserID       int64  `json:"user_id"`
	Content      string `json:"content"`
	CreationDate string `json:"creation_date"`
	User         User   `json:"user"`
}

func (s *CommentsStore) GetCommentsByPostId(ctx context.Context, postID int64) (*[]Comment, error) {
	query := `
	SELECT c.id, c.post_id, c.user_id, c.content, c.creation_date, u.username, u.id FROM "comments" c 
	LEFT JOIN users u 
	ON  c.user_id = u.id 
	WHERE c.post_id=$1
	ORDER BY c.creation_date DESC;
	`

	Rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer Rows.Close()

	comments := []Comment{}
 
	for Rows.Next() {
		var c Comment

		err = Rows.Scan(
			&c.ID,
			&c.PostID,
			&c.UserID,
			&c.Content,
			&c.CreationDate,
			&c.User.Username,
			&c.User.ID,
		)


		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return &comments, nil

}

func (s *CommentsStore) Create (ctx context.Context,c *Comment) error {
	query := `
	INSERT INTO comments(post_id, user_id, content)
	VALUES ($1, $2, $3) RETURNING id, creation_date
	`

	ctx,Cancel := context.WithTimeout(ctx, timeOutDuration)
	defer Cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		c.PostID,
		c.UserID,
		c.Content,
	).Scan(
		&c.ID,
		&c.CreationDate,
		
	)

	if err != nil {
		return err
	}

	return nil
}
