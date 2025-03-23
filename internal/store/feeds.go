package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type FeedsStore struct {
	db *sql.DB
}

type PostWtMetadata struct{
	Post
	CommentsCount int `json:"comments_count"`

}


func (s *FeedsStore) GetUserDefaultFeed(ctx context.Context, userId int64) ([]PostWtMetadata, error){

	query := `select 
		p.id, p.user_id, p.title, p.content, p.creation_date, p.version, p.tags,
		u.username, 
		count(c.id) as comments_count
	from (
		SELECT p.id, p.user_id, p.title, p.content, p.creation_date, p.version, p.tags
		FROM posts p
		where p.user_id = $1
		UNION ALL
		SELECT p.id, p.user_id, p.title, p.content, p.creation_date, p.version, p.tags
		FROM posts p
		JOIN followers f ON p.user_id = f.user_id
		WHERE f.follower_id = $1
		ORDER BY creation_date  
	) as p
	left join comments c on c.post_id = p.id
	left join users u on p.user_id = u.id
	group by p.id, p.user_id, p.title, p.content, p.creation_date, p.version, p.tags, u.username `

	ctx, Cancel := context.WithTimeout(ctx, timeOutDuration)
	defer Cancel()

	rows, err := s.db.QueryContext(
		ctx,
		query,
		userId,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}

	var posts = []PostWtMetadata{}

	for rows.Next() {
		var p PostWtMetadata
		err = rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.CreationDate,
			&p.Version,
			pq.Array(&p.Tags),
			&p.User.Username,
			&p.CommentsCount,
		)
		if err != nil {
			return []PostWtMetadata{}, err
		}
		posts = append(posts, p)

	}
	


	return posts,nil
}