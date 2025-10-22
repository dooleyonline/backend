package itemsvc

import (
	"context"
	"fmt"
	"time"

	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/internal/db"
	itemdb "github.com/dooleyonline/backend/internal/db/item"
	userdb "github.com/dooleyonline/backend/internal/db/user"
)


type Service struct {
	cfg *config.Config
	db  *db.DB
}


func New(cfg *config.Config, db *db.DB) *Service {
	return &Service{cfg, db}
}


func (s *Service) GetMany(ctx context.Context, params GetManyFilters) ([]itemdb.Item, error) {
	switch {
	case params.Seller != "":
		return s.db.Item.GetBySeller(ctx, params.Seller)
	case params.Query != "" && params.Category != "":
		searchParams := itemdb.SearchByCategoryParams{
			Category:  params.Category,
			ToTsquery: params.Query,
		}
		return s.db.Item.SearchByCategory(ctx, searchParams)
	case params.Query != "":
		return s.db.Item.Search(ctx, params.Query)
	case params.Category != "":
		return s.db.Item.GetByCategory(ctx, params.Category)
	default:
		return s.db.Item.GetAll(ctx)
	}
}


func (s *Service) Get(ctx context.Context, id int64) (itemdb.Item, error) {
	return s.db.Item.Get(ctx, id)
}


func (s *Service) Create(ctx context.Context, id string, p CreateUpdateInput) (itemdb.Item, error) {
	placeholder, err := generatePlaceholder(s.cfg, p.Images[0])
	if err != nil {
		return itemdb.Item{}, fmt.Errorf("failed to generate placeholders: %w", err)
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
		Seller:       id,
		Placeholder:  placeholder,
	}

	return s.db.Item.Create(ctx, dbparams)
}


func (s *Service) Update(ctx context.Context, id int64, p CreateUpdateInput) (itemdb.Item, error) {
	placeholder, err := generatePlaceholder(s.cfg, p.Images[0])
	if err != nil {
		return itemdb.Item{}, fmt.Errorf("failed to generate placeholders: %w", err)
	}

	dbparams := itemdb.UpdateParams{
		ID:           id,
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

	return s.db.Item.Update(ctx, dbparams)
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


func (s *Service) Like(ctx context.Context, itemId int64, userId string) ([]int64, error) {
	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return []int64{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	userTx := s.db.User.WithTx(tx)
	itemTx := s.db.Item.WithTx(tx)

	likedItems, err := userTx.AddLikedItem(ctx, userdb.AddLikedItemParams{
		ItemID: itemId,
		ID:     userId,
	})
	if err != nil {
		return likedItems, fmt.Errorf("failed to like item: %w", err)
	}

	if err := itemTx.IncrementLike(ctx, itemId); err != nil {
		return likedItems, fmt.Errorf("failed to increment item like: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return likedItems, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return likedItems, nil
}

func (s *Service) Unlike(ctx context.Context, itemId int64, userId string) ([]int64, error) {
	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return []int64{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	userTx := s.db.User.WithTx(tx)
	itemTx := s.db.Item.WithTx(tx)

	likedItems, err := userTx.DeleteLikedItem(ctx, userdb.DeleteLikedItemParams{
		ItemID: itemId,
		ID:     userId,
	})
	if err != nil {
		return likedItems, fmt.Errorf("failed to unlike item: %w", err)
	}

	if err := itemTx.DecrementLike(ctx, itemId); err != nil {
		return likedItems, fmt.Errorf("failed to decrement item like: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return likedItems, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return likedItems, nil
}


func (s *Service) GetBulk(ctx context.Context, ids []int64) ([]itemdb.Item, error) {
	return s.db.Item.GetByIDs(ctx, ids)
}