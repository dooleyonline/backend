package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/internal/model"
	itemsvc "github.com/dooleyonline/backend/internal/service/item"
)

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

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("non-200 response: %s", string(body))
	}

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

	var categories []model.Category
	if err := json.NewDecoder(res.Body).Decode(&categories); err != nil {
		return nil, fmt.Errorf("failed to unmarshal body: %w", err)
	}

	categoriesStr := make([]string, len(categories))
	for i, c := range categories {
		categoriesStr[i] = c.Name
	}

	return categoriesStr, nil
}

func createItem(ctx context.Context, cfg *config.Config, client *http.Client, cred *http.Cookie, item itemsvc.MutationParams) error {
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

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to read request body: %w", err)
		}
		return fmt.Errorf("request failed: %s", string(body))
	}

	return nil
}
