package db

import (
	"context"
	"fmt"

	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/sql/category"
	"github.com/dooleyonline/backend/sql/item"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Item     *item.Queries
	Category *category.Queries
	Pool     *pgxpool.Pool
}

func New(ctx context.Context, cfg *config.Config) (*DB, error) {
	pool, err := pgxpool.New(ctx, cfg.DatabaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	item := item.New(pool)
	category := category.New(pool)

	return &DB{
		Item:     item,
		Category: category,
		Pool:     pool,
	}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}
