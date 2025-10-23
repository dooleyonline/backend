package authsvc

import (
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
