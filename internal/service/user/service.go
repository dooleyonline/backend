package usersvc

import (
	"context"
	"fmt"

	"github.com/dooleyonline/backend/internal/db"
	userdb "github.com/dooleyonline/backend/internal/db/user"
	authsvc "github.com/dooleyonline/backend/internal/service/auth"
)

type Service struct {
	db *db.DB
}

func New(db *db.DB) *Service {
	return &Service{db}
}

func (s *Service) GetMany(ctx context.Context) (*[]userdb.User, error) {
	users, err := s.db.User.GetMany(ctx)
	if err != nil {
		return nil, err
	}

	var result []userdb.User
	for _, user := range users {
		user.Password = ""
		user.LikedItems = []int64{}
		result = append(result, user)
	}
	return &result, nil
}

type CreateParams struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

func (s *Service) Create(ctx context.Context, p *CreateParams) (*userdb.User, error) {
	dbparams := userdb.CreateParams{
		Email:     p.Email,
		FirstName: p.FirstName,
		LastName:  p.LastName,
	}

	hash, err := authsvc.HashPassword(p.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	dbparams.Password = string(hash)

	user, err := s.db.User.Create(ctx, dbparams)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user.Password = ""
	return &user, nil
}

func (s *Service) Get(ctx context.Context, id string) (*userdb.User, error) {
	user, err := s.db.User.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user.Password = ""
	user.LikedItems = []int64{}

	return &user, nil
}
