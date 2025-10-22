package categoryhandler

import (
	"fmt"
	"net/http"

	categorysvc "github.com/dooleyonline/backend/internal/service/category"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc *categorysvc.Service
}

func New(svc *categorysvc.Service) *Handler {
	return &Handler{svc}
}

// GetAll godoc
//
//	@Summary	Get all categories
//	@Tags		category
//	@Produce	json
//	@Success	200	{array}	categorydb.Category
//	@Router		/category [get]
func (h *Handler) GetAll(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)
	defer req.Body.Close()

	categories, err := h.svc.GetAll(ctx)
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
//	@Success	200		{object}	categorydb.Category
//	@Router		/category/{name} [get]
func (h *Handler) Get(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)
	defer req.Body.Close()

	var name string
	if err := echo.PathParamsBinder(c).String("name", &name).BindError(); err != nil {
		return fmt.Errorf("failed to bind name: %w", err)
	}

	category, err := h.svc.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to get category: %w", err)
	}

	return c.JSON(http.StatusOK, category)
}
