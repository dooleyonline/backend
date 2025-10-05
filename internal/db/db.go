package db

import (
	"context"
	"fmt"

	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/sql/category"
	"github.com/dooleyonline/backend/sql/item"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	Item *item.Queries
	Category *category.Queries
	Conn *pgx.Conn
}

func New(ctx context.Context, cfg *config.Config) (*DB, error) {
	conn, err := pgx.Connect(ctx, cfg.DatabaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	item := item.New(conn)
	category := category.New(conn)

	return &DB{
		Item:    item,
		Category: category,
		Conn:    conn,
	}, nil
}

func (db *DB) Close() {
	_ = db.Conn.Close(context.Background())
}
