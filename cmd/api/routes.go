package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/healthcheck/", app.healthcheckHandler)
	mux.HandleFunc("POST /api/v1/users/", app.registerUserHandler)
	return app.secureHeaders(mux)
}
