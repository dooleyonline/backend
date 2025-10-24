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

func (s *Service) GetAll(ctx context.Context) (*[]categorydb.Category, error) {
	categories, err := s.db.Category.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return &categories, nil
}

func (s *Service) Get(ctx context.Context, name string) (*categorydb.Category, error) {
	category, err := s.db.Category.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return &category, nil
}
