package authsvc

import (
	"context"
	"fmt"
	"time"

	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/internal/db"
	"github.com/dooleyonline/backend/internal/model"
)

type Service struct {
	cfg *config.Config
	db  *db.DB
}

func New(cfg *config.Config, db *db.DB) *Service {
	return &Service{cfg, db}
}

type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResult struct {
	User         *model.User
	Token        string
	CookieConfig CookieOptionsResult
}

func (s *Service) Login(ctx context.Context, params *LoginParams) (*LoginResult, error) {
	user, err := s.db.User.User.Get(ctx, params.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if verified := VerifyPassword(params.Password, user.Password); !verified {
		return nil, fmt.Errorf("invalid credentials")
	}

	user.Password = ""

	token, err := s.CreateJWT(params.Email, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	res, err := s.CookieOptions()
	if err != nil {
		return nil, fmt.Errorf("failed to get cookie options: %w", err)
	}

	return &LoginResult{
		User:         &user,
		Token:        token,
		CookieConfig: *res,
	}, nil
}

type CookieOptionsResult struct {
	AuthTokenName string        `json:"auth_token_name"`
	AuthTokenExp  time.Duration `json:"auth_token_exp"` // in seconds
	Secure        bool          `json:"secure"`
}

func (s *Service) CookieOptions() (*CookieOptionsResult, error) {
	return &CookieOptionsResult{
		AuthTokenName: s.cfg.AuthTokenName,
		AuthTokenExp:  s.cfg.AuthTokenExp,
		Secure:        s.cfg.IsProd,
	}, nil
}

func (s *Service) GetMe(ctx context.Context, id string) (*model.User, error) {
	user, err := s.db.User.User.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user.Password = ""

	return &user, nil
}
