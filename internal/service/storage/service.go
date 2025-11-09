package storagesvc

import (
	"context"

	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/internal/db"
	"github.com/dooleyonline/backend/internal/storage"
)

type Service struct {
	cfg *config.Config
	db  *db.DB
}

func New(cfg *config.Config, db *db.DB) *Service {
	return &Service{cfg, db}
}

func (s *Service) GetUploadPresignURL(ctx context.Context, params storage.PresignParams) (*storage.PresignResult, error) {
	storage := storage.New(s.cfg)
	res, err := storage.PresignUpload(ctx, params)
	if err != nil {
		return nil, err
	}
	return res, nil
}
