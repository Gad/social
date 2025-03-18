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

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	// TODO - mock user until auth is implemented

	var payload postPayload

	if err := readJson(app, w, r, &payload); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err, true)
		return
	}

	log.Println("Payload : ", payload)
	userId := 1
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

	app.jsonResponse(w, http.StatusOK, nil)
	

}

type updatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=100"`
}

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
		app.internalServerErrorResponse(w, r, err)
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
