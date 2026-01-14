package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

func GenerateTokens(userId int64) (*Tokens, error) {
	key := []byte(os.Getenv("JWT_SECRET_KEY"))

	accessT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})
	accessToken, err := accessT.SignedString(key)
	if err != nil {
		return nil, err
	}

	refreshT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 168).Unix(),
	})

	refreshToken, err := refreshT.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return nil, err
	}

	return &Tokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
