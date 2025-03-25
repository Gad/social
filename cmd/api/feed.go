package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gad/social/internal/store"
)

func (app *application) setFeedPagination(w http.ResponseWriter, r *http.Request) (store.FeedPaginationQuery, error) {

	var fpq = store.FeedPaginationQuery{
		Limit:  10,
		Offset: 0,
		Sort:   "desc",
	}
	qParams := r.URL.Query()
	limit := qParams.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fpq, err
		}
		fpq.Limit = l
	}
	offset := qParams.Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return fpq, err
		}
		fpq.Offset = o
	}

	sort := qParams.Get("sort")
	if sort != "" {
		fpq.Sort = sort
	}

	search := qParams.Get("search")
	if search != "" {
		fpq.Search = search
	}

	tags := qParams.Get("tags")
	log.Println(tags)
	if tags != "" {
		fpq.Tags = strings.Split(tags, ",")
	} else {
		fpq.Tags = []string{}
	}

	//since does not require initialization as it will be 0001-01-01
	since := qParams.Get("since")
	if since != "" {
		var err error
		fpq.Since, err = time.Parse("2006-01-02", since)
		if err != nil {
			return fpq, ErrDateFormat
		}
	}

	//until does require initialization to something "far into the future"
	until := qParams.Get("until")
	if until != "" {
		var err error
		fpq.Until, err = time.Parse("2006-01-02", until)
		if err != nil {
			return fpq, ErrDateFormat
		}
	} else {
		fpq.Until, _ = time.Parse("2006-01-02", "3000-12-31")
	}

	return fpq, nil

}

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {

	// Query + validate the URL parameters to allocate feed pagination and sorting. fallback to default values otherwise

	fpq, err := app.setFeedPagination(w, r)
	log.Printf("%+v", fpq)

	if err != nil {
		app.badRequestResponse(w, r, err, true)
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
