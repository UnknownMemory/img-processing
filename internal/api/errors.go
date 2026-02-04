package api

import "net/http"

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	errInt := map[string]interface{}{"error": message}

	err := app.writeJSON(w, status, errInt, nil)
	if err != nil {
		app.logger.Println(err)
		w.WriteHeader(500)
	}
}

func (app *Application) internalErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Println(err)
	app.errorResponse(w, r, http.StatusInternalServerError, "The server encountered a problem and could not process your request")
}
