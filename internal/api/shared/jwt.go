package shared

import (
	"fmt"
	"time"

	"github.com/dooleyonline/backend/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func CreateJWT(cfg *config.Config, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&JWTClaims{
			Email: email,
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

func ValidateJWT(cfg *config.Config, tokenString string) (*jwt.Token, error) {
	hmacSecretKey := []byte(cfg.AuthTokenSecret)

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (any, error) {
			return hmacSecretKey, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %w", err)
	}

	exp, err := token.Claims.GetExpirationTime()
	if err != nil {
		return nil, fmt.Errorf("failed to get expiration time: %w", err)
	}

	if exp.Before(time.Now()) {
		return nil, fmt.Errorf("token has expired")
	}

	return token, nil
}
