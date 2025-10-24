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
	categorydb "github.com/dooleyonline/backend/internal/db/category"
	itemdb "github.com/dooleyonline/backend/internal/db/item"
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
		fmt.Println(string(body))
		return nil, fmt.Errorf("non-200 response")
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
