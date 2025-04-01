package main

import (
	"encoding/json"

	"net/http"
)

func writeJson(w http.ResponseWriter, status int, data any) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)

}

func readJson(app *application, w http.ResponseWriter, r *http.Request, data any) error {

	r.Body = http.MaxBytesReader(w, r.Body, app.config.maxByte) // to mitigate ddos
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func writeJsonError(w http.ResponseWriter, status int, message string) error {
	type enveloppe struct {
		Error string `json:"error"`
	}

	return writeJson(w, status, &enveloppe{Error: message})
}

func (app *application) jsonResponse(w http.ResponseWriter, status int, data any) error {
	type enveloppe struct {
		Data any `json:"data"`
	}

	return writeJson(w, status, &enveloppe{Data: data})

}
