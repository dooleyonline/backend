package categoryapi

import (
	"fmt"
	"net/http"

	"github.com/dooleyonline/backend/internal/api/shared"
	"github.com/labstack/echo/v4"
)

// GetAll godoc
//
//	@Summary	Get all categories
//	@Tags		category
//	@Produce	json
//	@Success	200	{array}	category.Category
//	@Router		/category [get]
func GetAll(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
		db  = c.(shared.Context).DB
	)
	defer req.Body.Close()

	categories, err := db.Category.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get all categories: %w", err)
	}

	return c.JSON(http.StatusOK, categories)
}
