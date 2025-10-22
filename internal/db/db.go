package db

import (
	"context"
	"fmt"

	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/internal/db/category"
	"github.com/dooleyonline/backend/internal/db/item"
	"github.com/dooleyonline/backend/internal/db/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Item     *itemdb.Queries
	Category *categorydb.Queries
	User     *userdb.Queries
	Pool     *pgxpool.Pool
}

func New(ctx context.Context, cfg *config.Config) (*DB, error) {
	pool, err := pgxpool.New(ctx, cfg.DatabaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	item := itemdb.New(pool)
	category := categorydb.New(pool)
	user := userdb.New(pool)

	return &DB{
		Item:     item,
		Category: category,
		User:     user,
		Pool:     pool,
	}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}
