package main

import (
	"net/http"

)


func (app *application) getUserFeedHandler (w http.ResponseWriter, r *http.Request)  {

	ctx := r.Context()
	// TODO : revert to UserID auth
	UserID := int64(27)

	postsWithMetadata, err:=app.store.Feeds.GetUserDefaultFeed(ctx, UserID) 
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	app.jsonResponse(w,http.StatusOK, postsWithMetadata)
	 
}