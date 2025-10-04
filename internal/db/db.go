package db

import (
	"context"
	"fmt"

	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/sql"
	"github.com/jackc/pgx/v5"
)

type DB struct {
	*sql.Queries
	Conn *pgx.Conn
}

func New(ctx context.Context, cfg *config.Config) (*DB, error) {
	conn, err := pgx.Connect(ctx, cfg.DatabaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	queries := sql.New(conn)

	return &DB{
		Queries: queries,
		Conn:    conn,
	}, nil
}

func (db *DB) Close() {
	_ = db.Conn.Close(context.Background())
}
