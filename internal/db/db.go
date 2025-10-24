package db

import (
	"context"
	"fmt"

	"github.com/dooleyonline/backend/internal/config"
	itemcategory "github.com/dooleyonline/backend/internal/db/item/category"
	itemitem "github.com/dooleyonline/backend/internal/db/item/item"
	useruser "github.com/dooleyonline/backend/internal/db/user/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Item *item
	User *user
	Pool *pgxpool.Pool
}

type item struct {
	Item     *itemitem.Queries
	Category *itemcategory.Queries
}

type user struct {
	User *useruser.Queries
}

func New(ctx context.Context, cfg *config.Config) (*DB, error) {
	pool, err := pgxpool.New(ctx, cfg.DatabaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	item := &item{
		Item:     itemitem.New(pool),
		Category: itemcategory.New(pool),
	}

	user := &user{
		User: useruser.New(pool),
	}

	return &DB{
		Item: item,
		User: user,
		Pool: pool,
	}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}
