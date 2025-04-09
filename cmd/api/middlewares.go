package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read Auth header
			header := r.Header.Get("Authorization")

			if header == "" {
				app.basicAuthError(w, r, fmt.Errorf("authentication header not provided"))
				return
			}

			// parse it
			lr := strings.Split(header, " ")
			if len(lr) != 2 || lr[0] != "Basic" {
				app.basicAuthError(w, r, fmt.Errorf("authentication header malformed"))
				return
			}
			// decode base64
			credentials, err := base64.StdEncoding.Strict().DecodeString(lr[1])
			if err != err {
				app.basicAuthError(w, r, fmt.Errorf("authentication header malformed"))
				return
			}

			user, pass, ok := strings.Cut(string(credentials), ":")
			if !ok {
				app.basicAuthError(w, r, fmt.Errorf("authentication header malformed"))
				return
			}
			// check credentials

			if user != app.config.auth.basic.username || pass != app.config.auth.basic.password {
				app.basicAuthError(w, r, fmt.Errorf("invalid credentials"))
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
			app.AuthError(w, r, fmt.Errorf("authentication header not provided"))
			return
		}
		// parse it
		// header is in the form  "Authorization: Bearer token"
		lr := strings.Split(header, " ")
		if len(lr) != 2 || lr[0] != "Bearer" {
			app.AuthError(w, r, fmt.Errorf("authentication header malformed"))
			return
		}

		token := lr[1]
		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.AuthError(w, r, err)
			return
		}
		claims := jwtToken.Claims.(jwt.MapClaims)
		
		userID, err := strconv.ParseInt(fmt.Sprintf("%g", claims["sub"]), 10, 64)
		if err != nil {
			app.AuthError(w, r, err)
			return
		}
		
		user, err := app.store.Users.GetUserById(r.Context(), userID)
		
		if err != nil {
			app.AuthError(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
