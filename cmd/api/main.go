package main

import (
	"log"

	"github.com/gad/social/internal/db"
	"github.com/gad/social/internal/env"
	"github.com/gad/social/internal/store"
)
//	@title			Swagger Example API
//	@description	This is a sample server Petstore server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/v1

func main() {

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		apiURL : env.GetString("DEPLOYMENT_ADDR", "localhost:8080"),
		db: dbConfig{
			addr: env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/social?"+
				"sslmode=disable"),
			
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		
		},
		env: env.GetString("ENV", "DEVELOPMENT"),
		version : env.GetString("VERSION", "0.0.2"),
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
