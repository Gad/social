package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gad/social/internal/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T) *application {

	t.Helper()
	return &application{
		logger: zap.Must((zap.NewProduction())).Sugar(),
		store:  store.NewMockStore(),
		
	}
}


func execRequest(mux *chi.Mux, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}