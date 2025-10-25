package categoryhandler

import (
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
//	@Success	200	{array}	model.Category
//	@Router		/category [get]
func (h *Handler) GetAll(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	res, err := h.svc.GetAll(ctx)
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(http.StatusOK, res)
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

	var name string
	if err := echo.PathParamsBinder(c).String("name", &name).BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	res, err := h.svc.Get(ctx, name)
	if err != nil {
		return echo.ErrNotFound.WithInternal(err)
	}

	return c.JSON(http.StatusOK, res)
}
