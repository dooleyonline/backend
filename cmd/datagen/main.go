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
	"syscall"
	"time"

	sqlcategory "github.com/dooleyonline/backend/sql/category"
	sqlitem "github.com/dooleyonline/backend/sql/item"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"google.golang.org/genai"
)

const (
	ApiBaseUrl      = "http://localhost:8080"
	ReqCookieJwtKey = "dooleyonline_jwt"
	ReqCookieJwtVal = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RAZW1vcnkuZWR1IiwiZXhwIjoxNzYxNjI5NTE0fQ.-rTmjfhLG4GtQLUY9SS7N0n59kHnTIY5Ox0P4HEGopQ"

	Prompt = `Generate a list of items to post on an online secondhand marketplace.
	Make sure the images follow this pattern: sample/<any integer between 0 and 30>.webp.`
)

var (
	NumItems     = int64(10)
	ImagePattern = `sample/([0-9]|[12][0-9]|30)\.webp`
	ConditionMax = 5.0
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	categories, err := getCategories()
	if err != nil {
		slog.Error("failed to get categories", slog.Any("error", err))
		return
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		slog.Error("failed to create genai client", slog.Any("error", err))
		return
	}

	config := &genai.GenerateContentConfig{
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

	slog.Info("generating items...")

	content, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash-lite",
		genai.Text(Prompt),
		config,
	)
	if err != nil {
		slog.Error("failed to generate content", slog.Any("error", err))
		return
	}

	data := content.Text()

	var items []sqlitem.Item
	if err := json.Unmarshal([]byte(data), &items); err != nil {
		slog.Error("failed to unmarshal items", slog.Any("error", err))
		return
	}

	slog.Info("items generated", slog.Int("length", len(items)))
	slog.Info("making POST requests...")

	var wg sync.WaitGroup
	success := 0
	for _, item := range items {
		wg.Go(func() {
			if err := createItem(ctx, item); err == nil {
				success++
			}
		})
	}

	wg.Wait()
	stop()

	slog.Info("requests completed", slog.Int("generated", len(items)), slog.Int("success", success))
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

func getCategories() ([]string, error) {
	url, err := url.JoinPath(ApiBaseUrl, "category")
	if err != nil {
		return nil, fmt.Errorf("failed to join URL path: %w", err)
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request did not respond with 200")
	}

	var categories []sqlcategory.Category
	if err := json.NewDecoder(res.Body).Decode(&categories); err != nil {
		return nil, fmt.Errorf("failed to unmarshal body: %w", err)
	}

	categoriesStr := make([]string, len(categories))
	for i, c := range categories {
		categoriesStr[i] = c.Name
	}

	return categoriesStr, nil
}

func createItem(ctx context.Context, item sqlitem.Item) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	itemBytes, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	url, err := url.JoinPath(ApiBaseUrl, "item")
	if err != nil {
		return fmt.Errorf("failed to join URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(itemBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: ReqCookieJwtKey, Value: ReqCookieJwtVal})

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("request did not respond with 200")
	}

	return nil
}
