package main

import (
	"time"

	"github.com/gad/social/internal/auth"
	"github.com/gad/social/internal/db"
	"github.com/gad/social/internal/env"
	"github.com/gad/social/internal/mailer"
	"github.com/gad/social/internal/store"
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
	}

	logger.Fatal((app.run_app(app.mnt_mux())))

}
