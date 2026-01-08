package api

import (
	"context"
	"encoding/json"
	"net/http"

	db "github.com/unknownmemory/img-processing/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (app *Application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	data := &db.CreateUserParams{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {

	}

	data.Password = string(hashedPassword)
	q := db.New(app.db)
	err = q.CreateUser(context.Background(), *data)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "Could not register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
