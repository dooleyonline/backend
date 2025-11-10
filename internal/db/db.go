package db

import (
	"context"
	"fmt"

	"github.com/dooleyonline/backend/internal/config"
	chatmessage "github.com/dooleyonline/backend/internal/db/chat/message"
	chatparticipant "github.com/dooleyonline/backend/internal/db/chat/participant"
	chatroom "github.com/dooleyonline/backend/internal/db/chat/room"
	itemcategory "github.com/dooleyonline/backend/internal/db/item/category"
	itemitem "github.com/dooleyonline/backend/internal/db/item/item"
	userliked "github.com/dooleyonline/backend/internal/db/user/liked"
	useruser "github.com/dooleyonline/backend/internal/db/user/user"
	userviewed "github.com/dooleyonline/backend/internal/db/user/viewed"
	userverify "github.com/dooleyonline/backend/internal/db/user/verify"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Item *item
	User *user
	Chat *chat
	Pool *pgxpool.Pool
}

type item struct {
	Item     *itemitem.Queries
	Category *itemcategory.Queries
}

type user struct {
	User   *useruser.Queries
	Liked  *userliked.Queries
	Viewed *userviewed.Queries
	Verify *userverify.Queries
}

type chat struct {
	Message     *chatmessage.Queries
	Participant *chatparticipant.Queries
	Room        *chatroom.Queries
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
		User:   useruser.New(pool),
		Liked:  userliked.New(pool),
		Viewed: userviewed.New(pool),
		Verify: userverify.New(pool),
	}

	chat := &chat{
		Message:     chatmessage.New(pool),
		Participant: chatparticipant.New(pool),
		Room:        chatroom.New(pool),
	}

	return &DB{
		Item: item,
		User: user,
		Chat: chat,
		Pool: pool,
	}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}
