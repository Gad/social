package main

import (
	"context"
	"errors"
	
	"net/http"
	"strconv"

	"github.com/gad/social/internal/store"
	"github.com/go-chi/chi/v5"
)


const userKey Key = "user"


func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {

	user := getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, &user); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
	

}

func getUserFromCtx(r *http.Request) *store.User {

	user, _ := r.Context().Value(userKey).(*store.User)
	return user
}
func (app *application) userToContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		ID, err := strconv.Atoi(chi.URLParam(r, "userid"))
		if err != nil {

			app.badRequestResponse(w, r, err, false)
			return

		}
		userID := int64(ID)
		ctx := r.Context()
		user, err := app.store.Users.GetUserById(ctx, userID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrorNotFound):
				app.notFoundResponse(w, r, err)

			default:
				app.internalServerErrorResponse(w, r, err)

			}
			return
		}

		ctx = context.WithValue(ctx, userKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))

	})

	
}

