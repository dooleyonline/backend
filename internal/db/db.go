package db

import (
	"context"
	"fmt"

	"github.com/dooleyonline/backend/internal/config"
	sqlcategory "github.com/dooleyonline/backend/sql/category"
	sqlitem "github.com/dooleyonline/backend/sql/item"
	sqluser "github.com/dooleyonline/backend/sql/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Item     *sqlitem.Queries
	Category *sqlcategory.Queries
	User     *sqluser.Queries
	Pool     *pgxpool.Pool
}

func New(ctx context.Context, cfg *config.Config) (*DB, error) {
	pool, err := pgxpool.New(ctx, cfg.DatabaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	item := sqlitem.New(pool)
	category := sqlcategory.New(pool)
	user := sqluser.New(pool)

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
