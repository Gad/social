package main

import (
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {

	app := newTestApplication(t)
	mux := app.mnt_mux()

	
	t.Run("should not allow unauthenticated user", func(t *testing.T) {
		
		// Test getting a user that exists
		req, err := http.NewRequest("GET", "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := execRequest(mux, req)
		checkResponseCode(t, http.StatusUnauthorized, rr.Code)
	})


	testToken, err := app.authenticator.GenerateToken(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("should allow authenticated user", func(t *testing.T) {
		
		// Test getting a user that exists
		req, err := http.NewRequest("GET", "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)
		rr := execRequest(mux, req)
		checkResponseCode(t, http.StatusOK, rr.Code)
	})
}
