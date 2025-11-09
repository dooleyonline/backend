package verifysvc_test

import (
	"testing"

	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/internal/db"
	verifysvc "github.com/dooleyonline/backend/internal/service/verify"
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
	svc := verifysvc.New(cfg, db)
	// err = svc.Send(ctx, verifysvc.SendParams{
	// 	Email:  "changminlee517@gmail.com",
	// 	UserId: "0f10d091-43c6-4667-a5f7-bafd9e4795df",
	// })
	// if err != nil {
	// 	t.Fatalf("failed to send verification email: %v", err)
	// }

	err = svc.VerifyUserEmail(ctx, "9a20f65e-7501-419f-bd94-1516410d63a3")
	if err != nil {
		t.Fatalf("failed to verify user email: %v", err)
	}


}
