package main

import (
	"errors"
	"net/http"
)

var ErrDateFormat = errors.New("incorrect date format")

func (app *application) internalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Errorw("Internal Server Error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJsonError(w, http.StatusInternalServerError, "internal server error")

}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error, showError bool) {

	app.logger.Warnw("Bad Request Error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	var msg string
	if showError {
		msg = err.Error()
	} else {
		msg = "request parsing failed - request probably malformed"
	}
	writeJsonError(w, http.StatusBadRequest, msg)

}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Not Found Error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJsonError(w, http.StatusNotFound, "ressource not found")
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Errorw("Database conflict error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJsonError(w, http.StatusConflict, "database conflict")
}

func (app *application) basicauthError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("basic authorization error", r.Method, "path", r.URL.Path, "error", err.Error())

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	writeJsonError(w, http.StatusUnauthorized, "unauthorized")
}
func (app *application) authError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Token authorization error", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJsonError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Role Precedence failed", r.Method, "path", r.URL.Path,  "with error", err)

	writeJsonError(w, http.StatusForbidden, "unauthorized")
}
