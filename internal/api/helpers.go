package api

import (
	"encoding/json"
	"net/http"
)

func (app *Application) writeJSON(w http.ResponseWriter, status int, payload interface{}, headers http.Header) error {

	js, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		app.logger.Println(err)
		return nil
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		app.logger.Println(err)
		return nil
	}

	return nil
}
