package api

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/unknownmemory/img-processing/internal/auth"
)

func (app *Application) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	key := []byte(os.Getenv("JWT_SECRET_KEY"))
	refreshToken, err := r.Cookie("refreshToken")

	if err != nil {
		app.logger.Println(err)
		http.Error(w, "Invalid refresh token", http.StatusBadRequest)
		return
	}
	token, err := jwt.Parse(refreshToken.Value, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		app.logger.Println(err)
		http.Error(w, "Invalid refresh token", http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		app.logger.Println("Invalid refresh token")
		http.Error(w, "Invalid refresh token", http.StatusBadRequest)
		return
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		app.logger.Println("Invalid refresh token")
		http.Error(w, "Invalid refresh token", http.StatusBadRequest)
		return
	}

	userId := int64(sub)

	tokens, err := auth.GenerateTokens(userId)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}

	refreshCookie := http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		SameSite: http.SameSiteStrictMode,
		Path:     "/api/v1/token/refresh",
		Expires:  time.Now().Add(time.Hour * 168),
	}

	http.SetCookie(w, &refreshCookie)
	w.Header().Set("Content-Type", "application/json")
	jsData := map[string]string{"accessToken": tokens.AccessToken}
	err = app.writeJSON(w, http.StatusOK, jsData, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}

}
