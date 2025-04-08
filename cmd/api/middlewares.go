package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)



func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read Auth header 
			header := r.Header.Get("Authorization")
			
			if header == "" {
				app.basicAuthError(w,r,fmt.Errorf("authentication header not provided"))
				return
			}
			
			
			lr := strings.Split(header, " ")
			if len(lr) != 2 || lr[0] != "Basic" {
				app.basicAuthError(w,r,fmt.Errorf("authentication header malformed"))
				return
			}

			credentials, err := base64.StdEncoding.Strict().DecodeString(lr[1])
			if err != err {
				app.basicAuthError(w,r,fmt.Errorf("authentication header malformed"))
				return
			}

			user, pass, ok := strings.Cut(string(credentials), ":")
			if !ok {
				app.basicAuthError(w,r,fmt.Errorf("authentication header malformed"))
				return
			}

			if user != app.config.auth.basic.username || pass != app.config.auth.basic.password {
				app.basicAuthError(w,r,fmt.Errorf("invalid credentials"))
				return
			}


			// parse it
			// decode
			// check credentials
			next.ServeHTTP(w, r)
		})

		

	}
}