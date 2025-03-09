package main

import (
	"log"

	"github.com/gad/social/internal/env"
)

func main() {

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	app := application{
		config: cfg,
	}

	log.Fatal((app.run_app(app.mnt_mux())))

}
