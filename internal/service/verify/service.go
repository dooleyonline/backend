package verifysvc

import (
	"context"

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

func (s *Service) Get(ctx context.Context, id string) (*model.Verify, error) {
	verify, err := s.db.User.Verify.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &verify, nil
}
