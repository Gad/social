package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gad/social/internal/auth"
	"github.com/gad/social/internal/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T) *application {

	t.Helper()
	return &application{
		logger:        zap.NewNop().Sugar(), // or logger: zap.Must((zap.NewProduction())).Sugar(),
		store:         store.NewMockStore(),
		authenticator: auth.MockAuthenticator{},
	}
}

func execRequest(mux *chi.Mux, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	t.Helper()
	if expected != actual {
		t.Errorf("Response code : got %v wanted %v", actual, expected)
	}
}
