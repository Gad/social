package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gad/social/internal/store"
	"github.com/go-chi/chi/v5"
)

type postPayload struct {
	Title   string
	Content string
	Tags    []string
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	// TODO - mock user until auth is implemented

	var payload postPayload

	if err := readJson(app, w, r, &payload); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Println("Payload : ", payload)
	userId := 1
	p := &store.Post{
		Content: payload.Content,
		Title:   payload.Title,
		Tags:    payload.Tags,
		UserID:  int64(userId),
	}

	fmt.Printf("Complete payload : %v", payload)

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, p); err != nil {
		app.internalServerErrorResponse( w, r, err)
		return
	}

	if err := writeJson(w, http.StatusAccepted, &p); err != nil {
		app.internalServerErrorResponse( w, r, err)
		return
	}

}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {

	if postID, err := strconv.Atoi(chi.URLParam(r, "postid")); err != nil {
		
		app.badRequestResponse(w, r, err)

	} else {
		ctx := r.Context()
		if post, err := app.store.Posts.GetPostById(ctx, postID); err != nil {
			switch {
			case errors.Is(err, store.ErrorNotFound):
				app.notFoundResponse(w, r, err)

			default:
				app.internalServerErrorResponse( w, r, err)
				
			}
			return
		} else {
			writeJson(w, http.StatusOK, &post)
			return
		}
	}
}
