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
		user.LikedItems = []int64{}
		result = append(result, user)
	}
	return result, nil
}

type CreateParams struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
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
	user.LikedItems = []int64{}

	return &user, nil
}

type UpdateAvatarParams struct {
	UserID   string `json:"-"`
	AvatarID string `json:"avatar_id"`
}

func (s *Service) UpdateAvatar(ctx context.Context, params UpdateAvatarParams) error {
	err := s.db.User.User.UpdateAvatar(ctx, useruser.UpdateAvatarParams{
		ID:     params.UserID,
		Avatar: params.AvatarID,
	})
	if err != nil {
		return fmt.Errorf("failed to update user avatar: %w", err)
	}
	return nil
}
