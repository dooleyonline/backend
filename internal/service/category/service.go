package categorysvc

import (
	"context"

	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/internal/db"
	categorydb "github.com/dooleyonline/backend/internal/db/category"
)


type Service struct {
	cfg *config.Config
	db *db.DB
}

func New(cfg *config.Config, db *db.DB) *Service {
	return &Service{cfg, db}
}

func (s *Service) GetAll(ctx context.Context) ([]categorydb.Category, error) {
	return s.db.Category.GetAll(ctx)
}

func (s *Service) Get(ctx context.Context, name string) (categorydb.Category, error) {
	return s.db.Category.Get(ctx, name)
}