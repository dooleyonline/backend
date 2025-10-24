package itemsvc

import (
	"context"
	"fmt"
	"time"

	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/internal/db"
	itemdb "github.com/dooleyonline/backend/internal/db/item"
	userdb "github.com/dooleyonline/backend/internal/db/user"
	"github.com/dooleyonline/backend/internal/storage"
)

type Service struct {
	cfg *config.Config
	db  *db.DB
}

func New(cfg *config.Config, db *db.DB) *Service {
	return &Service{cfg, db}
}

type GetManyParams struct {
	Seller   string
	Query    string
	Category string
}

func (s *Service) GetMany(ctx context.Context, params *GetManyParams) (*[]itemdb.Item, error) {
	if params.Seller != "" && params.Query != "" && params.Category != "" {
		return nil, fmt.Errorf("cannot filter by seller, query, and category simultaneously")
	}

	var items []itemdb.Item
	var err error

	switch {
	case params.Seller != "":
		items, err = s.db.Item.GetBySeller(ctx, params.Seller)
	case params.Query != "" && params.Category != "":
		searchParams := itemdb.SearchByCategoryParams{
			Category:  params.Category,
			ToTsquery: params.Query,
		}
		items, err = s.db.Item.SearchByCategory(ctx, searchParams)
	case params.Query != "":
		items, err = s.db.Item.Search(ctx, params.Query)
	case params.Category != "":
		items, err = s.db.Item.GetByCategory(ctx, params.Category)
	default:
		items, err = s.db.Item.GetAll(ctx)
	}

	if err != nil {
		return nil, err
	}
	return &items, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*itemdb.Item, error) {
	item, err := s.db.Item.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

type MutationParams struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Images       []string `json:"images"`
	Price        float64  `json:"price"`
	Condition    int16    `json:"condition"`
	IsNegotiable bool     `json:"is_negotiable"`
	Category     string   `json:"category"`
	Subcategory  string   `json:"subcategory"`
}

func (s *Service) Create(ctx context.Context, sellerId string, p *MutationParams) (*itemdb.Item, error) {
	placeholder, err := generatePlaceholder(s.cfg, p.Images[0])
	if err != nil {
		return nil, fmt.Errorf("failed to generate placeholder: %w", err)
	}

	dbparams := itemdb.CreateParams{
		Name:         p.Name,
		Description:  p.Description,
		Images:       p.Images,
		Price:        p.Price,
		Condition:    p.Condition,
		IsNegotiable: p.IsNegotiable,
		Category:     p.Category,
		Subcategory:  p.Subcategory,
		Seller:       sellerId,
		Placeholder:  placeholder,
	}

	item, err := s.db.Item.Create(ctx, dbparams)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *Service) Update(ctx context.Context, itemId int64, p *MutationParams) (*itemdb.Item, error) {
	placeholder, err := generatePlaceholder(s.cfg, p.Images[0])
	if err != nil {
		return nil, fmt.Errorf("failed to generate placeholder: %w", err)
	}

	dbparams := itemdb.UpdateParams{
		ID:           itemId,
		Name:         p.Name,
		Description:  p.Description,
		Images:       p.Images,
		Price:        p.Price,
		Condition:    p.Condition,
		IsNegotiable: p.IsNegotiable,
		Category:     p.Category,
		Subcategory:  p.Subcategory,
		Placeholder:  placeholder,
	}

	item, err := s.db.Item.Update(ctx, dbparams)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.db.Item.Delete(ctx, id)
}

func (s *Service) Sell(ctx context.Context, id int64) error {
	now := time.Now()
	return s.db.Item.Sell(ctx, itemdb.SellParams{
		ID:     id,
		SoldAt: &now,
	})
}

func (s *Service) IncrementView(ctx context.Context, id int64) error {
	return s.db.Item.IncrementView(ctx, id)
}

func (s *Service) Like(ctx context.Context, itemId int64, userId string) error {
	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	userTx := s.db.User.WithTx(tx)
	itemTx := s.db.Item.WithTx(tx)

	if err := userTx.AddLikedItem(ctx, userdb.AddLikedItemParams{
		ItemID: itemId,
		ID:     userId,
	}); err != nil {
		return fmt.Errorf("failed to like item: %w", err)
	}

	if err := itemTx.IncrementLike(ctx, itemId); err != nil {
		return fmt.Errorf("failed to increment item like: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Service) Unlike(ctx context.Context, itemId int64, userId string) error {
	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	userTx := s.db.User.WithTx(tx)
	itemTx := s.db.Item.WithTx(tx)

	if err := userTx.DeleteLikedItem(ctx, userdb.DeleteLikedItemParams{
		ItemID: itemId,
		ID:     userId,
	}); err != nil {
		return fmt.Errorf("failed to unlike item: %w", err)
	}

	if err := itemTx.DecrementLike(ctx, itemId); err != nil {
		return fmt.Errorf("failed to decrement item like: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (s *Service) GetBatch(ctx context.Context, ids *[]int64) (*[]itemdb.Item, error) {
	items, err := s.db.Item.GetByIDs(ctx, *ids)
	if err != nil {
		return nil, err
	}
	return &items, nil
}

func (s *Service) GetUploadPresignURL(ctx context.Context, p *storage.PresignParams) (*storage.PresignResult, error) {
	storage := storage.New(s.cfg)
	res, err := storage.PresignUpload(ctx, p)
	if err != nil {
		return nil, err
	}
	return res, nil
}
