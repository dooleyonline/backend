package categorysvc

import (
	"context"

	"github.com/dooleyonline/backend/internal/db"
	"github.com/dooleyonline/backend/internal/model"
)

type Service struct {
	db *db.DB
}

func New(db *db.DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetAll(ctx context.Context) ([]model.Category, error) {
	categories, err := s.db.Item.Category.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *Service) Get(ctx context.Context, name string) (*model.Category, error) {
	category, err := s.db.Item.Category.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return &category, nil
}
