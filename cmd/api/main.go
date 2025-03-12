package main

import (
	"log"

	"github.com/gad/social/internal/db"
	"github.com/gad/social/internal/env"
	"github.com/gad/social/internal/store"
)

func main() {

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr: env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/social?"+
				"sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		
		},
		env: env.GetString("ENV", "DEVELOPMENT"),
		version : env.GetString("VERSION", "0.0.1"),
		maxByte: int64(env.GetInt("MAX_BYTES", 1_048_578)),
	}

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		log.Panic("Could not connect to database")
	}

	defer db.Close()
	store := store.NewStorage(db)

	app := application{
		config: cfg,
		store:  store,
	}

	log.Fatal((app.run_app(app.mnt_mux())))

}
