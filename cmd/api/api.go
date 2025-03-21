package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gad/social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr string
	db   dbConfig
	env	string
	version string
	maxByte int64 // max size for incoming http body to mitigate DDOS 
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}



func (app *application) mnt_mux() *chi.Mux {

	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)

	mux.Use(middleware.Timeout(time.Second * 60))

	mux.Route("/v1", func(m chi.Router) {
		m.Get("/health", app.getHealthHandler)
		m.Route("/posts", func(m chi.Router){
			m.Post("/", app.createPostHandler)
		    m.Route("/{postid}", func(m chi.Router){
				m.Use(app.postToContextMiddleware)
				m.Get("/", app.getPostHandler)
				m.Delete("/", app.deletePostHandler)
				m.Patch("/", app.patchPostHandler)
			
				})
		})
		m.Route("/users", func(m chi.Router){
			//m.Post("/", app.createPostHandler)
		    m.Route("/{userid}", func(m chi.Router){
				m.Use(app.userToContextMiddleware)
				m.Get("/", app.getUserHandler)
				
			
				})
		})
		
	})

	return mux

}

func (app *application) run_app(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 60,
	}
	log.Printf("Starting HTTP server at %s \n", app.config.addr)

	return srv.ListenAndServe()

}
