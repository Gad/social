package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
}

type config struct {
	addr string
}

func (app *application) getHealthHandler(w http.ResponseWriter, r *http.Request) {

	//TODO : implement
	w.Write([]byte("Everything ok boss !"))

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
