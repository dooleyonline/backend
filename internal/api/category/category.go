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

// Get godoc
//
//	@Summary	Get category by name
//	@Tags		category
//	@Produce	json
//	@Param		name	path		string	true	"Category name"
//	@Success	200		{object}	category.Category
//	@Router		/category/{name} [get]
func Get(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
		db  = c.(shared.Context).DB
	)
	defer req.Body.Close()

	var name string
	if err := echo.PathParamsBinder(c).String("name", &name).BindError(); err != nil {
		return fmt.Errorf("failed to bind name: %w", err)
	}

	category, err := db.Category.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to get category: %w", err)
	}

	return c.JSON(http.StatusOK, category)
}
