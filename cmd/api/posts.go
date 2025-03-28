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

type Key string

const postKey Key = "post"

type postPayload struct {
	Title   string   `json:"title" validate:"required,min=3,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

// CreatePost godoc
//
//	@Summary		Creates a post
//	@Description	Creates a post
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		postPayload	true	"Post payload"
//	@Success		202		{object}	store.Post
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts [post]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	// TODO - mock user until auth is implemented

	var payload postPayload

	if err := readJson(app, w, r, &payload); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err, true)
		return
	}

	log.Println("Payload : ", payload)
	userId := 27
	p := &store.Post{
		Content: payload.Content,
		Title:   payload.Title,
		Tags:    payload.Tags,
		UserID:  int64(userId),
	}

	log.Printf("Complete payload : %v", payload)

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, p); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusAccepted, &p); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

}
// GetPost godoc
//
//	@Summary		Fetches a post
//	@Description	Fetches a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	store.Post
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [get]
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {

	post := getPostFromCtx(r)

	// fetch potential comments

	comments, err := app.store.Comments.GetCommentsByPostId(r.Context(), post.ID)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	post.Comments = *comments

	if err := app.jsonResponse(w, http.StatusOK, &post); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

}

// DeletePost godoc
//
//	@Summary		Deletes a post
//	@Description	Delete a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		204	{object}	string
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.Atoi(chi.URLParam(r, "postid"))
	if err != nil {

		app.badRequestResponse(w, r, err, false)
		return
	}

	postID := int64(ID)
	ctx := r.Context()
	if err := app.store.Posts.DeletePostById(ctx, postID); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			app.notFoundResponse(w, r, err)

		default:
			app.internalServerErrorResponse(w, r, err)

		}
		return
	}

	app.jsonResponse(w, http.StatusNoContent, nil)
}


type updatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=100"`
}
// UpdatePost godoc
//
//	@Summary		Updates a post
//	@Description	Updates a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Post ID"
//	@Param			payload	body		updatePostPayload	true	"Post payload"
//	@Success		200		{object}	store.Post
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		404		{object}	error
//	@Failure		409		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [patch]
func (app *application) patchPostHandler(w http.ResponseWriter, r *http.Request) {

	post := getPostFromCtx(r)

	var payload updatePostPayload

	if err := readJson(app, w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err, true)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err, true)
		return
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}
	if err := app.store.Posts.UpdatePostById(r.Context(), post); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			app.conflictResponse(w, r, err)

		default:
			app.internalServerErrorResponse(w, r, err)

		}
		return

	}

	if err := app.jsonResponse(w, http.StatusOK, &post); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

}

func getPostFromCtx(r *http.Request) *store.Post {

	post, _ := r.Context().Value(postKey).(*store.Post)
	return post
}

func (app *application) postToContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ID, err := strconv.Atoi(chi.URLParam(r, "postid"))
		if err != nil {

			app.badRequestResponse(w, r, err, false)
			return

		}
		postID := int64(ID)
		ctx := r.Context()
		post, err := app.store.Posts.GetPostById(ctx, postID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrorNotFound):
				app.notFoundResponse(w, r, err)

			default:
				app.internalServerErrorResponse(w, r, err)

			}
			return
		}

		ctx = context.WithValue(ctx, postKey, post)
		next.ServeHTTP(w, r.WithContext(ctx))

	})

}
