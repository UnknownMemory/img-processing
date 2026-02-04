package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/unknownmemory/img-processing/internal/auth"
	db "github.com/unknownmemory/img-processing/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (app *Application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	data := &db.CreateUserParams{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		app.internalErrorResponse(w, r, err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {

	}

	data.Password = string(hashedPassword)
	q := db.New(app.db)
	err = q.CreateUser(context.Background(), *data)
	if err != nil {
		app.logger.Println(err)
		app.errorResponse(w, r, http.StatusInternalServerError, "Could not register user")
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *Application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		app.internalErrorResponse(w, r, err)
	}

	q := db.New(app.db)
	user, err := q.GetUser(context.Background(), data.Username)
	if err != nil {
		app.logger.Println(err)
		app.errorResponse(w, r, http.StatusNotFound, "Could not find user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		app.errorResponse(w, r, http.StatusUnauthorized, "Invalid credentials")
	}

	tokens, err := auth.GenerateTokens(user.ID)
	if err != nil {
		app.internalErrorResponse(w, r, err)
	}

	refreshCookie := http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		SameSite: http.SameSiteStrictMode,
		Path:     "/api/v1/token/refresh",
		Expires:  time.Now().Add(time.Hour * 168),
	}
	http.SetCookie(w, &refreshCookie)

	jsData := map[string]string{"accessToken": tokens.AccessToken}
	err = app.writeJSON(w, http.StatusOK, jsData, nil)
	if err != nil {
		app.internalErrorResponse(w, r, err)
	}
}
