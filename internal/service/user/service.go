package usersvc

import (
	"context"
	"fmt"

	"github.com/dooleyonline/backend/internal/db"
	userblock "github.com/dooleyonline/backend/internal/db/user/block"
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
	UserID    string
	Email     string
	FirstName string
	LastName  string
	AvatarID  string
}

func (s *Service) Update(ctx context.Context, params UpdateParams) error {
	err := s.db.User.User.Update(ctx, useruser.UpdateParams{
		ID:        params.UserID,
		Email:     params.Email,
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Avatar:    params.AvatarID,
	})
	if err != nil {
		return fmt.Errorf("failed to update user avatar: %w", err)
	}
	return nil
}

type blockParams struct {
	BlockerID string
	BlockedID string
}

func (s *Service) Block(ctx context.Context, params blockParams) error {
	err := s.db.User.Block.Block(ctx, userblock.BlockParams{
		BlockerID: params.BlockerID,
		BlockedID: params.BlockedID,
	})
	if err != nil {
		return fmt.Errorf("failed to block user: %w", err)
	}
	return nil
}

type UnblockParams struct {
	BlockerID string
	BlockedID string
}

func (s *Service) Unblock(ctx context.Context, params blockParams) error {
	err := s.db.User.Block.Unblock(ctx, userblock.UnblockParams{
		BlockerID: params.BlockerID,
		BlockedID: params.BlockedID,
	})
	if err != nil {
		return fmt.Errorf("failed to block user: %w", err)
	}
	return nil
}

func (s *Service) GetBlocked(ctx context.Context, id string) ([]model.User, error) {
	blocked, err := s.db.User.Block.GetBlocksByBlockerID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user_ids blocked by current user: %w", err)
	}

	blockedUsers := make([]model.User, 5)

	for _, userID := range blocked {
		user, err := s.db.User.User.GetByID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get blocked user: %w", err)
		}
		blockedUsers = append(blockedUsers, user)
	}
	return blockedUsers, nil
}

// GetAllBlocks should only be used in validating chat, search logics
func (s *Service) GetAllBlocks(ctx context.Context, id string) ([]string, error) {
	res := make([]string, 5)
	blocked, err := s.db.User.Block.GetBlocksByBlockerID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get blocked users: %w", err)
	}
	res = append(res, blocked...)

	blockedBy, err := s.db.User.Block.GetBlocksByBlockedID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to users that blocked the current user %w", err)
	}
	res = append(res, blockedBy...)

	return res, nil
}
