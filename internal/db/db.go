package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/gad/social/internal/env"
)

func New(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {

	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxIdleConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(env.GetDuration(maxIdleTime, 1*time.Minute))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	
	return db, nil
}
