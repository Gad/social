package main

import (
	
	"log"
	"net/http"
)



func (app *application) internalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {

	log.Printf("Internal Server Error: method: %s - path: %s - internal error: %s", r.Method, r.URL.Path, err.Error())
	writeJsonError(w, http.StatusInternalServerError, "internal server error")

}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error, showError bool) {

	log.Printf("Bad Request Error: method: %s - path: %s - internal error: %s", r.Method, r.URL.Path, err.Error())
	var msg string
	if showError { 
		msg = err.Error()
	} else {
		msg = "request parsing failed - request probably malformed"
	} 
	writeJsonError(w, http.StatusBadRequest, msg)

}

func (app *application) notFoundResponse (w http.ResponseWriter, r *http.Request, err error) {

	log.Printf("Not Found Error: method: %s - path: %s - internal error: %s", r.Method, r.URL.Path, err.Error())
	writeJsonError(w, http.StatusNotFound, "ressource not found")
}