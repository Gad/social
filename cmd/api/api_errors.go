package main

import (
	
	"log"
	"net/http"
)



func (app *application) internalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {

	log.Printf("Internal Server Error: method: %s - path: %s - internal error: %s", r.Method, r.URL.Path, err.Error())
	writeJsonError(w, http.StatusInternalServerError, "internal server error")

}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {

	log.Printf("Bad Request Error: method: %s - path: %s - internal error: %s", r.Method, r.URL.Path, err.Error())
	writeJsonError(w, http.StatusBadRequest, "request parsing failed - request probably malformed")

}

func (app *application) notFoundResponse (w http.ResponseWriter, r *http.Request, err error) {

	log.Printf("Not Found Error: method: %s - path: %s - internal error: %s", r.Method, r.URL.Path, err.Error())
	writeJsonError(w, http.StatusNotFound, "ressource not found")
}