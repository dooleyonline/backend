package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dooleyonline/backend/db"
	"github.com/jackc/pgx/v5"
)

func main() {
	slog.Info("server")

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "user=pqgotest dbname=pqgotest sslmode=verify-full")
	if err != nil {
		slog.Error("failed to connect to db", slog.Any("error", err))
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	authors, err := queries.ListAuthors(ctx)
	if err != nil {
		slog.Error("failed to execute query", slog.Any("error", err))
	}

	fmt.Println(authors)
}
