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
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid refresh token")
		return
	}
	token, err := jwt.Parse(refreshToken.Value, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid refresh token")
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid refresh token")
		return
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid refresh token")
		return
	}

	userId := int64(sub)

	tokens, err := auth.GenerateTokens(userId)
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
	w.Header().Set("Content-Type", "application/json")
	jsData := map[string]string{"accessToken": tokens.AccessToken}
	err = app.writeJSON(w, http.StatusOK, jsData, nil)
	if err != nil {
		app.internalErrorResponse(w, r, err)
	}

}
