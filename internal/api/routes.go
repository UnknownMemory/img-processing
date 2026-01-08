package api

import "net/http"

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/healthcheck/", app.healthcheckHandler)
	mux.HandleFunc("POST /api/v1/users/", app.registerUserHandler)
	return app.secureHeaders(mux)
}
