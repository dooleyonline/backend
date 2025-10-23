package auth

import (
	"time"

	"github.com/dooleyonline/backend/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	jwt.RegisteredClaims
}

func CreateJWT(cfg *config.Config, email string, id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&JWTClaims{
			Email: email,
			ID:    id,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.AuthTokenExp)),
			},
		},
	)

	tokenString, err := token.SignedString([]byte(cfg.AuthTokenSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
