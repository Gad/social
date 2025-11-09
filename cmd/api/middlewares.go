package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gad/social/internal/store"
	"github.com/golang-jwt/jwt/v5"
)

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read Auth header
			header := r.Header.Get("Authorization")

			if header == "" {
				app.basicauthError(w, r, fmt.Errorf("authentication header not provided"))
				return
			}

			// parse it
			lr := strings.Split(header, " ")
			if len(lr) != 2 || lr[0] != "Basic" {
				app.basicauthError(w, r, fmt.Errorf("authentication header malformed"))
				return
			}
			// decode base64
			credentials, err := base64.StdEncoding.Strict().DecodeString(lr[1])
			if err != err {
				app.basicauthError(w, r, fmt.Errorf("authentication header malformed"))
				return
			}

			user, pass, ok := strings.Cut(string(credentials), ":")
			if !ok {
				app.basicauthError(w, r, fmt.Errorf("authentication header malformed"))
				return
			}
			// check credentials

			if user != app.config.auth.basic.username || pass != app.config.auth.basic.password {
				app.basicauthError(w, r, fmt.Errorf("invalid credentials"))
				return
			}

			next.ServeHTTP(w, r)
		})

	}
}

func (app *application) TokenAuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// read Auth header
		header := r.Header.Get("Authorization")

		if header == "" {
			app.authError(w, r, fmt.Errorf("authentication header not provided"))
			return
		}
		// parse it
		// header is in the form  "Authorization: Bearer token"
		lr := strings.Split(header, " ")
		if len(lr) != 2 || lr[0] != "Bearer" {
			app.authError(w, r, fmt.Errorf("authentication header malformed"))
			return
		}

		token := lr[1]
		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.authError(w, r, err)
			return
		}
		claims := jwtToken.Claims.(jwt.MapClaims)

		userID, err := strconv.ParseInt(fmt.Sprintf("%g", claims["sub"]), 10, 64)
		if err != nil {
			app.authError(w, r, err)
			return
		}

		user, err := app.getUser(r.Context(), userID)
		if err != nil {
			switch err {
			case store.ErrorNotFound:
				app.authError(w, r, err)
			default:
				app.internalServerErrorResponse(w, r, err)
			}
			return
		}

		ctx := context.WithValue(r.Context(), userKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// getUser retrieves the user from the cache or the database
// and sets it in the context
// it returns the user and an error if any
func (app *application) getUser(ctx context.Context, userID int64) (*store.User, error) {

	// if redis is not enabled, get user from db
	if app.config.cacheState == None {
		app.logger.Infow("fetching user from db", "userID", userID)

		return app.store.Users.GetUserById(ctx, userID)
	}

	app.logger.Infow("Trying to get user from cache", "userID", userID)
	user, err := app.cacheStorage.Users.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		app.logger.Infow("fetching user from db", "userID", userID)
		user, err = app.store.Users.GetUserById(ctx, userID)
		if err != nil {
			return user, err
		}
	}

	// set user in cache
	if err := app.cacheStorage.Users.SetUser(ctx, user); err != nil {
		app.logger.Warn("failed to set user in cache", err)
	}

	return user, nil

}

func (app *application) checkOwnerShip(role string, next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// check if user owns the post
		user := getUserFromCtx(r)
		post := getPostFromCtx(r)

		if user.ID == post.UserID {
			next.ServeHTTP(w, r)
			return
		}
		// second chance : check if user has required role

		// extract full role from role name
		role, err := app.store.Roles.GetRoleByName(r.Context(), role)
		if err != nil {
			app.internalServerErrorResponse(w, r, err)
			return
		}

		// compare user role ID with role ID
		if user.Role.Level < role.Level {
			app.forbiddenResponse(w, r, fmt.Errorf("user does not have required role : user.role.level %d < required.role.level %d", user.Role.Level, role.Level))
			return
		}

		next.ServeHTTP(w, r)

	})
}

func (app *application) RateLimiter(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.rateLimitercfg.enabled {
			remoteAddr := strings.Split(r.RemoteAddr, ":")[0]
			allow, retryAfter, err := app.rateLimiter.Allow(remoteAddr)
			if err != nil {
				app.logger.Warnw("error trying to allow IP", "ratelimit type", app.config.rateLimitercfg.rateLimiterType, "remote addr", remoteAddr)
				return
			}
			if !allow {
				app.tooManyCallsResponse(w, r, retryAfter.String(), remoteAddr)
				return
			}

		}
		next.ServeHTTP(w, r)
	})
}
