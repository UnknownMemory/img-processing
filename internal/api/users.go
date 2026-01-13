package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func (app *Application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}

	q := db.New(app.db)
	user, err := q.GetUser(context.Background(), data.Username)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "Could not find user", http.StatusInternalServerError)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}

	accessT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user.ID,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
	})
	accessToken, err := accessT.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}
	refreshT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user.ID,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 168).Unix(),
	})

	refreshToken, err := refreshT.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}

	refreshCookie := http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		SameSite: http.SameSiteStrictMode,
		Path:     "/api/v1/token/refresh",
		Expires:  time.Now().Add(time.Hour * 168),
	}

	http.SetCookie(w, &refreshCookie)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"accessToken": accessToken})
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}
}
