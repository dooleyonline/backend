package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/internal/storage"
	"github.com/labstack/echo/v4"
)

func presignUpload(ctx context.Context, cfg *config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		defer req.Body.Close()

		body, err := io.ReadAll(req.Body)
		if err != nil {
			return fmt.Errorf("failed to read request body: %w", err)
		}

		upload := &storage.StorageRequest{
			Method: http.MethodPut,
		}

		if err := json.Unmarshal(body, upload); err != nil {
			return fmt.Errorf("failed to unmarshal body: %w", err)
		}

		// presign put request
		url, header, err := storage.Presign(ctx, cfg, upload)
		if err != nil {
			return fmt.Errorf("failed to presign request: %w", err)
		}

		return c.JSON(http.StatusOK, map[string]any{
			"url":    url,
			"header": header,
		})
	}
}
