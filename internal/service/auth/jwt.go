package authsvc

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	jwt.RegisteredClaims
}

func (s *Service) CreateJWT(email string, id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&JWTClaims{
			Email: email,
			ID:    id,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.AuthTokenExp)),
			},
		},
	)

	tokenString, err := token.SignedString([]byte(s.cfg.AuthTokenSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) ParseJWT(tokenRaw string) (string, error) {
	var claims JWTClaims
	token, err := jwt.ParseWithClaims(
		tokenRaw,
		&claims,
		func(t *jwt.Token) (any, error) {
			return []byte(s.cfg.AuthTokenSecret), nil
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to parse: %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return claims.ID, nil
}
