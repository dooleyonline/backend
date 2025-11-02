package itemsvc

import (
	"context"
	"fmt"
	"time"

	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/internal/db"
	itemitem "github.com/dooleyonline/backend/internal/db/item/item"
	useruser "github.com/dooleyonline/backend/internal/db/user/user"
	"github.com/dooleyonline/backend/internal/model"
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
	Page     int32
}

func (s *Service) GetMany(ctx context.Context, params *GetManyParams) ([]model.Item, error) {
	if params.Seller != "" && params.Query != "" && params.Category != "" {
		return nil, fmt.Errorf("cannot filter by seller, query, and category simultaneously")
	}

	if params.Page < 1 {
		params.Page = 1
	}

	var items []model.Item
	var err error

	switch {
	case params.Seller != "":
		getbysellerParams := itemitem.GetBySellerParams{
			SellerID: params.Seller,
			Page:     params.Page,
			Size:     s.cfg.ItemPageSize,
		}
		items, err = s.db.Item.Item.GetBySeller(ctx, getbysellerParams)
	case params.Query != "" && params.Category != "":
		searchbycategoryParams := itemitem.SearchByCategoryParams{
			Category:  params.Category,
			ToTsquery: params.Query,
			Page:      params.Page,
			Size:      s.cfg.ItemPageSize,
		}
		items, err = s.db.Item.Item.SearchByCategory(ctx, searchbycategoryParams)
	case params.Query != "":
		searchparams := itemitem.SearchParams{
			WebsearchToTsquery: params.Query,
			Page:               params.Page,
			Size:               s.cfg.ItemPageSize,
		}
		items, err = s.db.Item.Item.Search(ctx, searchparams)
	case params.Category != "":
		getbycategoryParams := itemitem.GetByCategoryParams{
			Category: params.Category,
			Page:     params.Page,
			Size:     s.cfg.ItemPageSize,
		}
		items, err = s.db.Item.Item.GetByCategory(ctx, getbycategoryParams)
	default:
		getAllParams := itemitem.GetAllParams{
			Page: params.Page,
			Size: s.cfg.ItemPageSize,
		}
		items, err = s.db.Item.Item.GetAll(ctx, getAllParams)
	}

	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*model.Item, error) {
	item, err := s.db.Item.Item.Get(ctx, id)
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

func (s *Service) Create(ctx context.Context, sellerId string, p *MutationParams) (*model.Item, error) {
	placeholder, err := generatePlaceholder(s.cfg, p.Images[0])
	if err != nil {
		return nil, fmt.Errorf("failed to generate placeholder: %w", err)
	}

	dbparams := itemitem.CreateParams{
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

	item, err := s.db.Item.Item.Create(ctx, dbparams)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *Service) Update(ctx context.Context, itemId int64, p *MutationParams) (*model.Item, error) {
	placeholder, err := generatePlaceholder(s.cfg, p.Images[0])
	if err != nil {
		return nil, fmt.Errorf("failed to generate placeholder: %w", err)
	}

	dbparams := itemitem.UpdateParams{
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

	item, err := s.db.Item.Item.Update(ctx, dbparams)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.db.Item.Item.Delete(ctx, id)
}

func (s *Service) Sell(ctx context.Context, id int64) error {
	now := time.Now()
	return s.db.Item.Item.Sell(ctx, itemitem.SellParams{
		ID:     id,
		SoldAt: &now,
	})
}

func (s *Service) IncrementView(ctx context.Context, id int64) error {
	return s.db.Item.Item.IncrementView(ctx, id)
}

func (s *Service) Like(ctx context.Context, itemId int64, userId string) error {
	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	userTx := s.db.User.User.WithTx(tx)
	itemTx := s.db.Item.Item.WithTx(tx)

	if err := userTx.AddLikedItem(ctx, useruser.AddLikedItemParams{
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

	userTx := s.db.User.User.WithTx(tx)
	itemTx := s.db.Item.Item.WithTx(tx)

	if err := userTx.DeleteLikedItem(ctx, useruser.DeleteLikedItemParams{
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

func (s *Service) GetBatch(ctx context.Context, ids *[]int64, page int32) ([]model.Item, error) {
	getbatchparams := itemitem.GetBatchParams{
		ItemIds: *ids,
		Page:    page,
		Size:    s.cfg.ItemPageSize,
	}
	items, err := s.db.Item.Item.GetBatch(ctx, getbatchparams)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *Service) GetUploadPresignURL(ctx context.Context, contentType string) (*storage.PresignResult, error) {
	storage := storage.New(s.cfg)
	res, err := storage.PresignUpload(ctx, contentType)
	if err != nil {
		return nil, err
	}
	return res, nil
}
