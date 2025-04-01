package main

import (
	"context"
	"errors"
	"log"

	"net/http"
	"strconv"

	"github.com/gad/social/internal/store"
	"github.com/go-chi/chi/v5"
)

const userKey Key = "user"

type FollowingUser struct {
	UserID int64 `json:"user_id"`
}

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id}	[get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {

	user := getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, &user); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

}

// FollowUser godoc
//
//	@Summary		authenticated user follows another user
//	@Description	follow a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		204	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id}/follow  [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {

	followedUser := getUserFromCtx(r)

	var follower FollowingUser

	if err := readJson(app, w, r, &follower); err != nil {
		app.badRequestResponse(w, r, err, true)
		return
	}

	if err := app.store.Users.Follow(r.Context(), followedUser.ID, follower.UserID); err != nil {

		app.badRequestResponse(w, r, err, true)

	}

	w.WriteHeader(http.StatusNoContent)
	/*if err := app.jsonResponse(w, http.StatusNoContent, ""); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}*/

}

// UnfollowUser godoc
//
//	@Summary		authenticated user unfollows another user
//	@Description	follow a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		204	{object}	store.User
//	@Failure		400	{object}	error	"malformed request"
//	@Failure		404	{object}	error	"User not found"
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id}/unfollow  [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {

	followedUser := getUserFromCtx(r)

	var follower FollowingUser

	if err := readJson(app, w, r, &follower); err != nil {
		app.badRequestResponse(w, r, err, true)
		return
	}

	if err := app.store.Users.Unfollow(r.Context(), followedUser.ID, follower.UserID); err != nil {
		app.internalServerErrorResponse(w, r, err)

	}

	if err := app.jsonResponse(w, http.StatusNoContent, ""); err != nil {
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

// ActivateUser godoc
//
//	@Summary		Activate a user
//	@Description	Activate a user given a token
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		202		{string}	string	"user activated"
//	@Failure		400		{object}	error	"malformed request"
//	@Failure		404		{object}	error	"token not found"
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token}  [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	err := app.store.Users.Activate(r.Context(), token)
	if err != nil {
		switch err {
		case store.ErrorNotFound:
			app.notFoundResponse(w, r, err)

		default:
			app.internalServerErrorResponse(w, r, err)
		}
		return
	}
	if err := app.jsonResponse(w, http.StatusAccepted, nil); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

}
