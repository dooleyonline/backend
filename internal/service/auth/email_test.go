package authsvc_test

import (
	"testing"

	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/internal/db"
	authsvc "github.com/dooleyonline/backend/internal/service/auth"
)

func TestEmailSend(t *testing.T) {
	ctx := t.Context()

	cfg, err := config.New()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	db, err := db.New(ctx, cfg)
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}
	svc := authsvc.New(cfg, db)
	id, err := svc.CreateVerification(ctx, authsvc.SendParams{
		Email:  "taeeunk1208@gmail.com",
		UserId: "0f10d091-43c6-4667-a5f7-bafd9e4795df",
	})
	if err != nil {
		t.Fatalf("failed to send verification email: %v", err)
	}

	err = svc.VerifyUserEmail(ctx, id)
	if err != nil {
		t.Fatalf("failed to verify user email: %v", err)
	}

}
