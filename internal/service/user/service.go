package usersvc

import (
	"context"
	"fmt"

	"github.com/dooleyonline/backend/internal/auth"
	"github.com/dooleyonline/backend/internal/db"
	userdb "github.com/dooleyonline/backend/internal/db/user"
)

type Service struct {
	db *db.DB
}

func New(db *db.DB) *Service {
	return &Service{db}
}

func (s *Service) GetMany(ctx context.Context) ([]UserSummary, error) {
	users, err := s.db.User.GetMany(ctx)
	if err != nil {
		return nil, err
	}

	var result []UserSummary
	for _, user := range users {
		result = append(result, UserSummary{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		})
	}
	return result, nil
}

func (s *Service) Create(ctx context.Context, p CreateInput) (UserSummary, error) {
	dbparams := userdb.CreateParams{
		Email:     p.Email,
		FirstName: p.FirstName,
		LastName:  p.LastName,
	}

	hash, err := auth.HashPassword(p.Password)
	if err != nil {
		return UserSummary{}, fmt.Errorf("failed to hash password: %w", err)
	}
	dbparams.Password = string(hash)

	user, err := s.db.User.Create(ctx, dbparams)
	if err != nil {
		return UserSummary{}, fmt.Errorf("failed to create user: %w", err)
	}

	return UserSummary{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}

func (s *Service) GetSellerByID(ctx context.Context, id string) (Seller, error) {
	seller, err := s.db.User.GetSellerByID(ctx, id)
	if err != nil {
		return Seller{}, fmt.Errorf("failed to get seller: %w", err)
	}
	return Seller{
		ID:        seller.ID,
		FirstName: seller.FirstName,
		LastName:  seller.LastName,
	}, nil
}

func (s *Service) GetMe(ctx context.Context, id string) (Me, error) {
	full, err := s.db.User.GetFullUserByID(ctx, id)
	if err != nil {
		return Me{}, fmt.Errorf("failed to get user: %w", err)
	}

	return Me{
		ID:        full.ID,
		Email:     full.Email,
		FirstName: full.FirstName,
		LastName:  full.LastName,
		LikedItems: full.LikedItems,
	}, nil
}