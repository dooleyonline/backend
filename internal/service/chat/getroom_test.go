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
func TestSyncAllMessageCounts(t *testing.T) {
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

	t.Log("Syncing message counts for all rooms...")
	
	if err := svc.SyncAllMessageCounts(ctx); err != nil {
		t.Fatalf("SyncAllMessageCounts returned error: %v", err)
	}

	t.Log("Successfully synced all room message counts")
}
