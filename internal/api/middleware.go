package api

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (app *Application) secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	})
}

func (app *Application) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader != "" {
			accessToken := strings.Split(tokenHeader, "Bearer ")[1]

			token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			})
			if err != nil {
				app.errorResponse(w, r, http.StatusUnauthorized, "Unauthorized")
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)

			if !ok || !token.Valid {
				app.logger.Println("Token isn't valid")
				app.errorResponse(w, r, http.StatusUnauthorized, "Unauthorized")
				return
			}

			sub, ok := claims["sub"].(float64)
			if !ok {
				app.logger.Println("Claim sub doesn't exist.")
				app.errorResponse(w, r, http.StatusUnauthorized, "Unauthorized")
				return
			}
			userId := int64(sub)

			ctx := context.WithValue(r.Context(), "user_id", userId)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		app.errorResponse(w, r, http.StatusUnauthorized, "Unauthorized")
	})
}
