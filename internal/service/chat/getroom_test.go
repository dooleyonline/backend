package chatsvc_test

import (
	"testing"

	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/internal/db"
	chatsvc "github.com/dooleyonline/backend/internal/service/chat"
)

func TestGetroom(t *testing.T) {
	ctx := t.Context()

	cfg, err := config.New()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	db, err := db.New(ctx, cfg)
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}
	svc := chatsvc.New(cfg, db)
	userID := "70458b3b-839b-431f-92a0-b073e16a9f09"
	rooms, err := svc.GetRooms(ctx, userID)
	if err != nil {
		t.Fatalf("GetRooms returned error: %v", err)
	}

	t.Logf("got %d rooms: %+v", len(rooms), rooms)

}
