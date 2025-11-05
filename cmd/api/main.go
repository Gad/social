package main

import (
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/gad/social/internal/auth"
	"github.com/gad/social/internal/cache"
	"github.com/gad/social/internal/db"
	"github.com/gad/social/internal/env"
	"github.com/gad/social/internal/mailer"
	"github.com/gad/social/internal/ratelimiter"
	"github.com/gad/social/internal/store"
	"github.com/redis/go-redis/v9"
)

//	@title			Swagger Example API
//	@description	Gophersocial API Server
//	@termsOfService	http://myterms.notyours

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {

	cfg := config{
		addr:        env.GetString("ADDR", ":8080"),
		apiURL:      env.GetString("DEPLOYMENT_ADDR", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_ADDR", "http://localhost:5173"),
		db: dbConfig{
			addr: env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/social?"+
				"sslmode=disable"),

			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env:     env.GetString("ENV", "DEVELOPMENT"),
		version: env.GetString("VERSION", "0.0.2"),
		maxByte: int64(env.GetInt("MAX_BYTES", 1_048_578)),
		mail: mailConfig{
			exp:        time.Hour * 24 * 3,
			fromEmail:  env.GetString("FROM_EMAIL", ""),
			maxRetries: env.GetInt("MAIL_MAX_RETRIES", 3),
			mailTrap: mailTrapConfig{
				apiKey:       env.GetString("MAILTRAP_API_KEY", ""),
				smtpAddr:     env.GetString("MAILTRAP_SMTP_ADDR", "live.smtp.mailtrap.io"),
				smtpPort:     env.GetInt("MAILTRAP_SMTP_PORT", 587),
				smtpUsername: env.GetString("MAILTRAP_USERNAME", "api"),
			},
		},
		auth: authconfig{
			basic: basicConfig{
				username: env.GetString("BASIC_AUTH_USER_NAME", "admin"),
				password: env.GetString("BASIC_AUTH_PASSWORD", "admin"),
			},
			token: tokenConfig{
				secret: env.GetString("TOKEN_AUTH_SECRET", "secret"),
				issuer: env.GetString("TOKEN_AUTH_ISSUER", "Gophersocial"),
				exp:    time.Hour * 24 * 3,
			},
		},
		cacheState: cacheState(env.GetInt("CACHE_STATE", int(None))),
		redisCfg: redisConfig{
			addr:     env.GetString("REDIS_ADDR", "localhost:6379"),
			password: env.GetString("REDIS_PASSWORD", ""),
			db:       env.GetInt("REDIS_DB", 0),
			ttl:      env.GetDuration("REDIS_TTL", time.Minute),
		},
		badgerCfg: badgerConfig{
			path: env.GetString("BADGER_PATH", "/tmp/badger"),
			ttl:  env.GetDuration("BADGER_TTL", time.Minute),
		},
		memcacheCfg: memcachedConfig{
			host:         env.GetString("MEMCACHED_ADDR", "localhost"),
			startingPort: env.GetInt("MEMCACHED_STARTING_PORT", 11211),
			endingPort:   env.GetInt("MEMCACHED_ENDING_PORT", 11211),
			ttl:          env.GetDuration("MEMCACHED_TTL", time.Minute),
		},
		rateLimitercfg: rateLimiterConfig{
			rateLimiterType:      env.GetString("RATELIMITER_TYPE", "FIXED_WINDOW"),
			requestsPerTimeFrame: env.GetInt("RATELIMITER_REQUESTS_PER_TIMEFRAME", 20),
			timeFrame:            env.GetDuration("RATELIMITER_WINDOW_DURATION", time.Second*5),
			enabled:              env.GetBool("RATELIMITER_ENABLED", true),
		},
	}

	// logger setup
	lg := createLogger(cfg.env)

	defer lg.Sync()
	logger := lg.Sugar()
	logger.Infow("Starting GopherSocial", "version", cfg.version, "Env", cfg.env)

	// Database setup
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		logger.Fatal("Could not connect to database")
	}
	logger.Info("Connection with database established")

	defer db.Close()
	store := store.NewStorage(db)

	//Cache setup
	var redisdb *redis.Client
	var badgerdb *badger.DB
	var cacheStorage cache.Storage

	switch cfg.cacheState {
	case Redis:
		redisdb = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.password, cfg.redisCfg.db)
		logger.Info("Connection with Redis cache established")
		cacheStorage = cache.NewRedisStorage(redisdb, cfg.redisCfg.ttl)
		defer redisdb.Close()

	case Badger:
		badgerdb, err = cache.NewBadgerDB()
		if err != nil {
			logger.Errorw("Could not connect to badger database", "error", err)
			badgerdb = nil
		} else {
			logger.Info("Connection with Badger cache established")
			defer badgerdb.Close()
		}
		cacheStorage = cache.NewBadgerStorage(badgerdb, cfg.badgerCfg.ttl)

	case Memcached:
		memDb := cache.NewMemcachedClient(cfg.memcacheCfg.host, cfg.memcacheCfg.startingPort, cfg.memcacheCfg.endingPort)
		if memDb == nil {
			logger.Errorw("Could not connect to memcached database")
		} else {
			logger.Info("Connection with Memcached caches established")
			defer memDb.Close()
		}
		cacheStorage = cache.NewMemcachedStorage(memDb, cfg.memcacheCfg.ttl)

	case None:
		redisdb = nil
		badgerdb = nil
		cacheStorage = cache.Storage{}
		logger.Info("all cache disabled")

	}

	var rateLimiter ratelimiter.Limiter
	switch cfg.rateLimitercfg.rateLimiterType {
	case "FIXED_WINDOW":
		rateLimiter = ratelimiter.NewFixedWindowLimiter(
			cfg.rateLimitercfg.requestsPerTimeFrame,
			cfg.rateLimitercfg.timeFrame,
		)
	}
	// Mailer setup

	mailer := mailer.NewMailtrap(
		cfg.mail.mailTrap.apiKey,
		cfg.mail.fromEmail,
		cfg.mail.mailTrap.smtpAddr,
		cfg.mail.mailTrap.smtpUsername,
		cfg.mail.mailTrap.smtpPort,
		cfg.mail.maxRetries,
	)

	jwtAuth := auth.NewSimpleJWTAuthenticator(
		cfg.auth.token.secret,
		cfg.auth.token.issuer,
		cfg.auth.token.issuer)

	app := application{
		config:        cfg,
		store:         store,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuth,
		cacheStorage:  cacheStorage,
		rateLimiter:   rateLimiter,
	}

	logger.Fatal((app.run_app(app.mnt_mux())))

}
