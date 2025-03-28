package main

import (
	"net/http"
)

// healthCheck godoc
//
//	@Summary		Healthcheck
//	@Description	Healthcheck endpoint
//	@Tags			operations
//	@Produce		json
//	@Success		200	{object}	string	"ok"
//	@Failure		500 {object}	error
//	@Router			/health [get]
func (app *application) getHealthHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": app.config.version,
	}

	if err := writeJson(w, http.StatusOK, data); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}

}
