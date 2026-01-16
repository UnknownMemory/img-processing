package api

import "net/http"

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/healthcheck/", app.healthcheckHandler)
	mux.HandleFunc("POST /api/v1/users/register", app.registerUserHandler)
	mux.HandleFunc("POST /api/v1/users/login", app.loginUserHandler)
	mux.HandleFunc("GET /api/v1/token/refresh", app.refreshTokenHandler)
	mux.Handle("POST /api/v1/images/upload", app.authMiddleware(http.HandlerFunc(app.uploadImageHandler)))
	return app.secureHeaders(mux)
}
