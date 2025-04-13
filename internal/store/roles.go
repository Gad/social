package store

import (
	"context"
	"database/sql"
)

type RolesStore struct {
	db *sql.DB
}

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Level	   int    `json:"level"`
	Description string `json:"description"`
}


func (s *RolesStore) GetRoleByName(ctx context.Context, name string) (*Role, error) {
	
	query := `SELECT id, name, level, description FROM roles WHERE name = $1`
	
	ctx, Cancel := context.WithTimeout(ctx, timeOutDuration)
	defer Cancel()
	// Execute the query and scan the result into a Role struct
	row := s.db.QueryRowContext(ctx, query, name)
	var role Role
	err := row.Scan(&role.ID, &role.Name, &role.Level, &role.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrorNotFound
		}
		return nil, err // Other error
	}

	return &role, nil
}
