package authsvc

import (
	"context"
	"fmt"

	"github.com/dooleyonline/backend/internal/auth"
	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/internal/db"
	usersvc "github.com/dooleyonline/backend/internal/service/user"
)

type Service struct {
	cfg *config.Config
	db  *db.DB
}

func New(cfg *config.Config, db *db.DB) *Service {
	return &Service{cfg, db}
}

func (s *Service) Login(ctx context.Context, params LoginParams) (LoginResponse, error) {
	user, err := s.db.User.Get(ctx, params.Email)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("failed to get user: %w", err)
	}

	if verified := auth.VerifyPassword(params.Password, user.Password); !verified {
		return LoginResponse{}, fmt.Errorf("invalid credentials")
	}

	token, err := auth.CreateJWT(s.cfg, params.Email, user.ID)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("failed to create token: %w", err)
	}

	return LoginResponse{
		User: usersvc.Me{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
		Token: token,
	}, nil
}

func (s *Service) CookieDetails() (CookieDetailsResponse, error) {
	return CookieDetailsResponse{
		AuthTokenName: s.cfg.AuthTokenName,
		AuthTokenExp:  s.cfg.AuthTokenExp,
		Secure:        s.cfg.IsProd,
	}, nil
}
