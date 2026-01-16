package api

import (
	"net/http"
)

func (app *Application) uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
