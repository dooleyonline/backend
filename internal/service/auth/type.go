package authsvc

import (
	"time"

	usersvc "github.com/dooleyonline/backend/internal/service/user"
)

type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User  usersvc.Me `json:"user"`
	Token string     `json:"token"`
}

type CookieDetailsResponse struct {
	AuthTokenName string        `json:"auth_token_name"`
	AuthTokenExp  time.Duration `json:"auth_token_exp"` // in seconds
	Secure        bool          `json:"secure"`
}
