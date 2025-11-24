package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/dooleyonline/backend/internal/config"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

const (
	EnvEmail    = "DATAGEN_EMAIL"
	EnvPassword = "DATAGEN_PASSWORD"

	Prompt = `Generate a list of items to post on an online secondhand marketplace.
	Randomize the images array.`
)

var (
	NumItems = int64(10)
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.New()
	if err != nil {
		slog.Error("failed to initialize config", slog.Any("error", err))
		return
	}

	email, ok := os.LookupEnv(EnvEmail)
	if !ok {
		slog.Error("environment variable not provided", slog.String("env", EnvEmail))
		return
	}

	password, ok := os.LookupEnv(EnvPassword)
	if !ok {
		slog.Error("environment variable not provided", slog.String("env", EnvPassword))
		return
	}

	client := http.DefaultClient

	slog.Info("logging in...")

	cred, err := login(ctx, cfg, client, email, password)
	if err != nil {
		slog.Error("failed to log in", slog.Any("error", err))
		return
	}

	slog.Info("fetching categories ...")

	categories, err := getCategories(cfg, client)
	if err != nil {
		slog.Error("failed to get categories", slog.Any("error", err))
		return
	}

	slog.Info("generating items...")

	items, err := generate(ctx, categories)
	if err != nil {
		slog.Error("failed to generate items", slog.Any("error", err))
	}

	slog.Info("items generated", slog.Int("length", len(items)))
	slog.Info("making POST requests...")

	var wg sync.WaitGroup
	var success atomic.Int64
	for _, item := range items {
		wg.Go(func() {
			if err := createItem(ctx, cfg, client, cred, item); err != nil {
				slog.Error("failed to create item", slog.Any("error", err))
				return
			}
			success.Add(1)
		})
	}

	wg.Wait()
	time.Sleep(time.Second)
	stop()

	slog.Info("requests completed", slog.Int("generated", len(items)), slog.Int("success", int(success.Load())))
}

func init() {
	lg := slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			NoColor:    !isatty.IsTerminal(os.Stderr.Fd()),
			TimeFormat: time.Kitchen,
		}),
	)

	slog.SetDefault(lg)
}
