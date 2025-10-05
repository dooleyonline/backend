package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dooleyonline/backend/internal/api"
	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/internal/db"
	"github.com/lmittmann/tint"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	lg := slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		}),
	)

	cfg, err := config.New()
	if err != nil {
		lg.Error("failed to initialize config", slog.Any("error", err))
		return
	}

	db, err := db.New(ctx, cfg)
	if err != nil {
		lg.Error("failed to initialize DB", slog.Any("error", err))
		return
	}
	defer db.Close()

	server, err := api.New(ctx, cfg, db, lg)
	if err != nil {
		lg.Error("failed to initialize API", slog.Any("error", err))
		return
	}

	go func() {
		if err := server.Start(cfg.ServerAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			lg.Error("server error", slog.Any("error", err))
			stop()
		}
	}()

	lg.Info("server started",
		slog.String("addr", cfg.ServerAddr),
		slog.String("docs", fmt.Sprintf("%s/swagger/index.html", cfg.ServerAddr)),
	)

	<-ctx.Done()

	stop()

	// force shutdown after 5 seconds
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		lg.Error("failed to stop server", slog.Any("error", err))
	}
	lg.Info("server stopped")
}
