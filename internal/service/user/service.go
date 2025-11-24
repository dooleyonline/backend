package usersvc

import (
	"context"
	"fmt"

	"github.com/dooleyonline/backend/internal/db"
	useruser "github.com/dooleyonline/backend/internal/db/user/user"
	"github.com/dooleyonline/backend/internal/model"
	authsvc "github.com/dooleyonline/backend/internal/service/auth"
)

type Service struct {
	db *db.DB
}

func New(db *db.DB) *Service {
	return &Service{db}
}

func (s *Service) GetMany(ctx context.Context) ([]model.User, error) {
	users, err := s.db.User.User.GetMany(ctx)
	if err != nil {
		return nil, err
	}

	var result []model.User
	for _, user := range users {
		user.Password = ""
		result = append(result, user)
	}
	return result, nil
}

type CreateParams struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (s *Service) Create(ctx context.Context, p *CreateParams) (*model.User, error) {
	dbparams := useruser.CreateParams{
		Email:     p.Email,
		FirstName: p.FirstName,
		LastName:  p.LastName,
	}

	hash, err := authsvc.HashPassword(p.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	dbparams.Password = string(hash)

	user, err := s.db.User.User.Create(ctx, dbparams)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user.Password = ""
	return &user, nil
}

func (s *Service) Get(ctx context.Context, id string) (*model.User, error) {
	user, err := s.db.User.User.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user.Password = ""

	return &user, nil
}

func (s *Service) GetLikes(ctx context.Context) ([]model.Liked, error) {
	liked, err := s.db.User.Liked.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get liked table: %w", err)
	}

	return liked, nil
}

func (s *Service) GetViews(ctx context.Context) ([]model.Viewed, error) {
	viewed, err := s.db.User.Viewed.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get viewed table: %w", err)
	}

	return viewed, nil
}

type UpdateParams struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Avatar    string `json:"avatar"`
}

func (s *Service) Update(ctx context.Context, params UpdateParams) error {
	err := s.db.User.User.Update(ctx, useruser.UpdateParams{
		ID:        params.ID,
		Email:     params.Email,
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Avatar:    params.Avatar,
	})
	if err != nil {
		return fmt.Errorf("failed to update user avatar: %w", err)
	}
	return nil
}
