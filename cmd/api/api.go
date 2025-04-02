package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gad/social/docs"
	"github.com/gad/social/internal/mailer"
	"github.com/gad/social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type application struct {
	config config
	store  store.Storage
	logger *zap.SugaredLogger
	mailer mailer.Client
}

type config struct {
	addr    string
	apiURL  string
	db      dbConfig
	env     string
	version string
	maxByte int64 // max size for incoming http body to mitigate DDOS
	mail mailConfig
}

type mailConfig struct{
	exp time.Duration
	sendGrid sendGridConfig
	mailTrap mailTrapConfig
	fromEmail string
	
}

type sendGridConfig struct{
	apiKey string
}


type mailTrapConfig struct{
	apiKey string
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

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		m.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))
		m.Route("/posts", func(m chi.Router) {
			m.Post("/", app.createPostHandler)
			m.Route("/{postid}", func(m chi.Router) {
				m.Use(app.postToContextMiddleware)
				m.Get("/", app.getPostHandler)
				m.Delete("/", app.deletePostHandler)
				m.Patch("/", app.patchPostHandler)

			})
		})
		m.Route("/users", func(m chi.Router) {
			
			//m.Post("/", app.createPostHandler)
			m.Route("/{userid}", func(m chi.Router) {
				m.Use(app.userToContextMiddleware)
				m.Get("/", app.getUserHandler)
				m.Put("/follow", app.followUserHandler)
				m.Put("/unfollow", app.unfollowUserHandler)

			})
			m.Put("/activate/{token}", app.activateUserHandler)
			m.Group(func(m chi.Router) {
				m.Get("/feed", app.getUserFeedHandler)

			})
		})
		m.Route("/authentication", func(m chi.Router){
			m.Post("/user", app.registerUserHandler)
		})

	})

	return mux

}

func (app *application) run_app(mux http.Handler) error {

	docs.SwaggerInfo.Version = app.config.version
	docs.SwaggerInfo.Host = app.config.apiURL

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 60,
	}
	app.logger.Infof("Starting HTTP server at %s", app.config.addr)

	return srv.ListenAndServe()

}
