package storagehandler

import (
	"fmt"
	"net/http"

	storagec "github.com/dooleyonline/backend/internal/storage"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	sc *storagec.Client
}

func New(sc *storagec.Client) *Handler {
	return &Handler{sc}
}

// Presign godoc
//
//	@Summary	Generate presigned URL for s3 operation
//	@Tags		other
//	@Accept		json
//	@Produce	json
//	@Param		item	body		storage.PresignRequest	true	"Presign params"
//	@Success	200		{object}	itemdb.Item
//	@Router		/storage/presign [post]
func (h *Handler) Presign(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)
	defer req.Body.Close()

	var input storagec.PresignRequest
	if err := c.Bind(&input); err != nil {
		return fmt.Errorf("failed to bind params: %w", err)
	}

	var (
		url    string
		header http.Header
		err    error
	)

	switch input.Method {
	case http.MethodPut:
		if input.Bucket == "" || input.Key == "" || input.ContentType == "" {
			return fmt.Errorf("invalid presign parameters")
		}
		url, header, err = h.sc.Presign(ctx, &input)
	default:
		return fmt.Errorf("invalid storage request")
	}

	if err != nil {
		return fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return c.JSON(http.StatusOK, storagec.PresignResponse{
		URL:    url,
		Header: header,
	})
}
