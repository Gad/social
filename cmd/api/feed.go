package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gad/social/internal/store"
)



func setFeedPagination(r *http.Request) (store.FeedPaginationQuery, error) {

	var fpq = store.FeedPaginationQuery{
		Limit: 2,
		Offset: 0,
		Sort: "desc",
	}

	limit := r.URL.Query().Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fpq, err
		}
		fpq.Limit = l
	}
	offset := r.URL.Query().Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return fpq, err
		}
		fpq.Offset = o
	}

	sort := r.URL.Query().Get("sort")
	if sort != ""{
		fpq.Sort = sort
	}

	return fpq, nil

}

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {

	// Query + validate the URL parameters to allocate feed pagination and sorting. fallback to default values otherwise

	
	
	fpq, err := setFeedPagination(r) 
	log.Printf("%+v", fpq)

	if err!= nil{
		app.badRequestResponse(w,r,err,true)
		return
	}

	err = validate.Struct(fpq)
	if err != nil {
		app.badRequestResponse(w, r, err, true)
		return
	}

	ctx := r.Context()
	// TODO : revert to UserID auth
	UserID := int64(27)

	postsWithMetadata, err := app.store.Feeds.GetUserDefaultFeed(ctx, UserID, fpq)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, postsWithMetadata)

}
