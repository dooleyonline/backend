package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/dooleyonline/backend/internal/config"
	categorydb "github.com/dooleyonline/backend/internal/db/category"
	itemdb "github.com/dooleyonline/backend/internal/db/item"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"google.golang.org/genai"
)

const (
	EnvEmail    = "DATAGEN_EMAIL"
	EnvPassword = "DATAGEN_PASSWORD"

	Prompt = `Generate a list of items to post on an online secondhand marketplace.
	Make sure the images follow this pattern: sample/<any integer between 0 and 10>.webp.`
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
	slog.Info("logged in")

	categories, err := getCategories(cfg, client)
	if err != nil {
		slog.Error("failed to get categories", slog.Any("error", err))
		return
	}

	gemini, err := genai.NewClient(ctx, &genai.ClientConfig{
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		slog.Error("failed to create genai client", slog.Any("error", err))
		return
	}

	slog.Info("generating items...")

	content, err := gemini.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash-lite",
		genai.Text(Prompt),
		geminiConfig(categories),
	)
	if err != nil {
		slog.Error("failed to generate content", slog.Any("error", err))
		return
	}

	data := content.Text()

	var items []itemdb.Item
	if err := json.Unmarshal([]byte(data), &items); err != nil {
		slog.Error("failed to unmarshal items", slog.Any("error", err))
		return
	}

	slog.Info("items generated", slog.Int("length", len(items)))
	slog.Info("making POST requests...")

	var wg sync.WaitGroup
	var success atomic.Int64
	for _, item := range items {
		wg.Go(func() {
			if err := createItem(ctx, cfg, client, cred, item); err == nil {
				success.Add(1)
			}
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

type Credential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func login(ctx context.Context, cfg *config.Config, client *http.Client, email, password string) (*http.Cookie, error) {
	url, err := url.JoinPath(cfg.Url, "auth", "login")
	if err != nil {
		return nil, fmt.Errorf("failed to join URL path: %w", err)
	}

	cred := Credential{email, password}

	credBytes, err := json.Marshal(cred)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal credential: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(credBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer res.Body.Close()

	for _, c := range res.Cookies() {
		if c.Name == cfg.AuthTokenName {
			return c, nil
		}
	}

	return nil, fmt.Errorf("failed to find cookie")
}

func getCategories(cfg *config.Config, client *http.Client) ([]string, error) {
	url, err := url.JoinPath(cfg.Url, "category")
	if err != nil {
		return nil, fmt.Errorf("failed to join URL path: %w", err)
	}

	res, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", res.Status)
	}

	var categories []categorydb.Category
	if err := json.NewDecoder(res.Body).Decode(&categories); err != nil {
		return nil, fmt.Errorf("failed to unmarshal body: %w", err)
	}

	categoriesStr := make([]string, len(categories))
	for i, c := range categories {
		categoriesStr[i] = c.Name
	}

	return categoriesStr, nil
}

func createItem(ctx context.Context, cfg *config.Config, client *http.Client, cred *http.Cookie, item itemdb.Item) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	itemBytes, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	url, err := url.JoinPath(cfg.Url, "item")
	if err != nil {
		return fmt.Errorf("failed to join URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(itemBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.AddCookie(cred)

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("request did not respond with 200")
	}

	return nil
}

var (
	ImagePattern = `sample/([0-9]|10)\.webp`
	ConditionMax = 5.0
)

func geminiConfig(categories []string) *genai.GenerateContentConfig {
	return &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type:     genai.TypeArray,
			MinItems: &NumItems,
			MaxItems: &NumItems,
			Items: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"name":        {Type: genai.TypeString},
					"description": {Type: genai.TypeString},
					"images": {
						Type: genai.TypeArray,
						Items: &genai.Schema{
							Type:    genai.TypeString,
							Pattern: ImagePattern,
						},
					},
					"price": {
						Type:   genai.TypeNumber,
						Format: "float",
					},
					"condition": {
						Type:    genai.TypeInteger,
						Default: 0,
						Maximum: &ConditionMax,
					},
					"is_negotiable": {Type: genai.TypeBoolean},
					"category": {
						Type:   genai.TypeString,
						Format: "enum",
						Enum:   categories,
					},
					"subcategory": {
						Type: genai.TypeString,
						Enum: []string{"Other"},
					},
				},
				Required: []string{"name", "description", "images", "price", "condition", "is_negotiable", "category", "subcategory"},
			},
		},
	}
}
