package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/healthcheck", app.healthcheckHandler)
	return app.secureHeaders(mux)
}
