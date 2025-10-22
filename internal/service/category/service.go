package categorysvc

import (
	"context"

	"github.com/dooleyonline/backend/internal/db"
	categorydb "github.com/dooleyonline/backend/internal/db/category"
)

type Service struct {
	db *db.DB
}

func New(db *db.DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetAll(ctx context.Context) ([]categorydb.Category, error) {
	return s.db.Category.GetAll(ctx)
}

func (s *Service) Get(ctx context.Context, name string) (categorydb.Category, error) {
	return s.db.Category.Get(ctx, name)
}
