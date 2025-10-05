package storage

import (
	"fmt"
	"net/http"

	"github.com/dooleyonline/backend/internal/api/shared"
	"github.com/dooleyonline/backend/internal/storage"
	"github.com/labstack/echo/v4"
)

// Presign godoc
//
//	@Summary	Generate presigned URL for s3 operation
//	@Tags		other
//	@Accept		json
//	@Produce	json
//	@Param		item	body		storage.PresignParams	true	"Presign params"
//	@Success	200		{object}	item.Item
//	@Router		/storage/presign [post]
func Presign(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
		cfg = c.(shared.Context).Cfg
	)
	defer req.Body.Close()

	var params storage.PresignParams
	if err := c.Bind(&params); err != nil {
		return fmt.Errorf("failed to bind params: %w", err)
	}

	url, header, err := storage.Presign(ctx, cfg, &params)
	if err != nil {
		return fmt.Errorf("failed to presign request: %w", err)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"url":    url,
		"header": header,
	})
}
