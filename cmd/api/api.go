package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gad/social/docs"
	"github.com/gad/social/internal/auth"
	"github.com/gad/social/internal/cache"
	"github.com/gad/social/internal/mailer"
	"github.com/gad/social/internal/ratelimiter"
	"github.com/gad/social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type application struct {
	config        config
	store         store.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
	cacheStorage  cache.Storage
	rateLimiter   ratelimiter.Limiter
}

type config struct {
	addr           string
	apiURL         string
	frontendURL    string
	db             dbConfig
	env            string
	version        string
	maxByte        int64 // max size for incoming http body to mitigate DDOS
	mail           mailConfig
	auth           authconfig
	cacheState     cacheState
	redisCfg       redisConfig
	badgerCfg      badgerConfig
	memcacheCfg    memcachedConfig
	rateLimitercfg rateLimiterConfig
}

type cacheState int

const (
	None cacheState = iota
	Redis
	Badger
	Memcached
)

type memcachedConfig struct {
	host         string
	startingPort int
	endingPort   int
	ttl          time.Duration
}

type badgerConfig struct {
	path string
	ttl  time.Duration
}

type redisConfig struct {
	addr     string
	password string
	db       int
	ttl      time.Duration
}

type authconfig struct {
	basic basicConfig
	token tokenConfig
}

type basicConfig struct {
	username string
	password string
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	issuer string
}

type mailConfig struct {
	exp        time.Duration
	sendGrid   sendGridConfig
	mailTrap   mailTrapConfig
	fromEmail  string
	maxRetries int
}

type sendGridConfig struct {
	apiKey string
}

type mailTrapConfig struct {
	apiKey       string
	smtpAddr     string
	smtpPort     int
	smtpUsername string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type rateLimiterConfig struct {
	rateLimiterType      string
	enabled              bool
	requestsPerTimeFrame int
	timeFrame            time.Duration
}

func (app *application) mnt_mux() *chi.Mux {

	m := chi.NewRouter()

	m.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	m.Use(middleware.Logger)
	m.Use(middleware.Recoverer)
	m.Use(middleware.RequestID)
	m.Use(middleware.RealIP)
	if app.config.rateLimitercfg.enabled {
		m.Use(app.RateLimiter)
	}
	m.Use(middleware.Timeout(time.Second * 60))

	m.Route("/v1", func(m chi.Router) {

		//m.With(app.BasicAuthMiddleware()).Get("/health", app.getHealthHandler)
		//TODO : reactivate Basic Authentication for /health after testing graceful shutdown
		m.Get("/health", app.getHealthHandler)
		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		m.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))
		m.Route("/posts", func(m chi.Router) {
			m.Use(app.TokenAuthMiddleware)
			m.Post("/", app.createPostHandler)
			m.Route("/{postid}", func(m chi.Router) {
				m.Use(app.postToContextMiddleware)
				m.Get("/", app.getPostHandler)
				m.Delete("/", app.checkOwnerShip("admin", app.deletePostHandler))
				m.Patch("/", app.checkOwnerShip("moderator", app.patchPostHandler))

			})
		})
		m.Route("/users", func(m chi.Router) {

			m.Route("/{userid}", func(m chi.Router) {
				m.Use(app.TokenAuthMiddleware)
				m.Get("/", app.getUserHandler)
				m.Put("/follow", app.followUserHandler)
				m.Put("/unfollow", app.unfollowUserHandler)

			})

			m.Group(func(m chi.Router) {
				m.Use(app.TokenAuthMiddleware)
				m.Get("/feed", app.getUserFeedHandler)

			})

			m.Put("/activate/{token}", app.activateUserHandler)
		})
		m.Route("/authentication", func(m chi.Router) {
			m.Post("/user", app.registerUserHandler)
			m.Post("/token", app.setTokenHandler)
		})

	})

	return m

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

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Infow("Signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)

	}()

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	app.logger.Infow("listenAndServe error caught", "error", err)

	err = <-shutdown
	if err != nil {
		return err
	}

	app.logger.Infow("Server stopped", "addr", app.config.addr, "env", app.config.env)

	return nil

}
