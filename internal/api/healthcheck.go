package api

import (
	"net/http"
)

func (app *Application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"status": "available", "environment": app.config.Mode, "version": app.version}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.internalErrorResponse(w, r, err)
		return
	}
}
