package itemhandler

import (
	"fmt"
	"net/http"

	"github.com/dooleyonline/backend/internal/api/shared"
	itemsvc "github.com/dooleyonline/backend/internal/service/item"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc *itemsvc.Service
}

func New(svc *itemsvc.Service) *Handler {
	return &Handler{svc}
}

// GetMany godoc
//
//	@Summary	Get many items
//	@Tags		item
//	@Produce	json
//	@Param		seller		query	string	false	"Seller filter"
//	@Param		q			query	string	false	"Search query"
//	@param		category	query	string	false	"Category filter"
//	@Success	200			{array}	itemsvc.Item
//	@Router		/item [get]
func (h *Handler) GetMany(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)
	defer req.Body.Close()

	seller := c.QueryParam("seller")
	query := c.QueryParam("q")
	category := c.QueryParam("category")

	items, err := h.svc.GetMany(ctx, itemsvc.GetManyFilters{
		Seller:   seller,
		Query:    query,
		Category: category,
	})
	if err != nil {
		return fmt.Errorf("failed to get items: %w", err)
	}

	return c.JSON(http.StatusOK, items)
}

// Get godoc
//
//	@Summary	Get item by ID
//	@Tags		item
//	@Produce	json
//	@Param		id	path		int	true	"Item ID"
//	@Success	200	{object}	itemsvc.Item
//	@Router		/item/{id} [get]
func (h *Handler) Get(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)
	defer req.Body.Close()

	var itemId int64
	if err := echo.PathParamsBinder(c).Int64("id", &itemId).BindError(); err != nil {
		return fmt.Errorf("failed to bind id: %w", err)
	}

	item, err := h.svc.Get(ctx, itemId)
	if err != nil {
		return fmt.Errorf("failed to get item: %w", err)
	}

	return c.JSON(http.StatusOK, item)
}

// Create godoc
//
//	@Summary	Create item
//	@Tags		item
//	@Accept		json
//	@Produce	json
//	@Param		item	body		itemsvc.CreateUpdateInput	true	"Item"
//	@Success	200		{object}	itemsvc.Item
//	@Router		/item [post]
func (h *Handler) Create(c echo.Context) error {
	var (
		req  = c.Request()
		ctx  = req.Context()
		user = c.(shared.Context).User
	)
	defer req.Body.Close()

	var params itemsvc.CreateUpdateInput
	if err := c.Bind(&params); err != nil {
		return fmt.Errorf("failed to bind params: %w", err)
	}

	item, err := h.svc.Create(ctx, user.ID, params)
	if err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}

	return c.JSON(http.StatusOK, item)
}

// Update godoc
//
//	@Summary	Update item by ID
//	@Tags		item
//	@Accept		json
//	@Produce	json
//	@Param		id		path		int						true	"Item ID"
//	@Param		item	body		itemsvc.CreateUpdateInput	true	"Item"
//	@Success	200		{object}	itemsvc.Item
//	@Router		/item/{id} [put]
func (h *Handler) Update(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)
	defer req.Body.Close()

	var params itemsvc.CreateUpdateInput
	if err := c.Bind(&params); err != nil {
		return fmt.Errorf("failed to bind params: %w", err)
	}
	var itemId int64
	if err := echo.PathParamsBinder(c).Int64("id", &itemId).BindError(); err != nil {
		return fmt.Errorf("failed to bind id: %w", err)
	}

	item, err := h.svc.Update(ctx, itemId, params)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return c.JSON(http.StatusOK, item)
}

// Delete godoc
//
//	@Summary	Delete item by ID
//	@Tags		item
//	@Param		id	path	int	true	"Item ID"
//	@Success	200
//	@Router		/item/{id} [delete]
func (h *Handler) Delete(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)
	defer req.Body.Close()

	var itemId int64
	if err := echo.PathParamsBinder(c).Int64("id", &itemId).BindError(); err != nil {
		return fmt.Errorf("failed to bind id: %w", err)
	}

	if err := h.svc.Delete(ctx, itemId); err != nil {
		return fmt.Errorf("failed to get item: %w", err)
	}

	return c.NoContent(http.StatusOK)
}

// Sell godoc
//
//	@Summary	Update sold_at property by ID
//	@Tags		item
//	@Param		id	path		int		true	"Item ID"
//	@Success	200	{string}	string	"Item ID"
//	@Router		/item/{id}/sell [post]
func (h *Handler) Sell(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)
	defer req.Body.Close()

	var itemId int64
	if err := echo.PathParamsBinder(c).Int64("id", &itemId).BindError(); err != nil {
		return fmt.Errorf("failed to bind id: %w", err)
	}

	if err := h.svc.Sell(ctx, itemId); err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return c.NoContent(http.StatusOK)
}

// IncrementView godoc
//
//	@Summary	Increment item views by ID
//	@Tags		item
//	@Param		id	path	int	true	"Item ID"
//	@Success	200
//	@Router		/item/{id}/view [post]
func (h *Handler) IncrementView(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)
	defer req.Body.Close()

	var itemId int64
	if err := echo.PathParamsBinder(c).Int64("id", &itemId).BindError(); err != nil {
		return fmt.Errorf("failed to bind id: %w", err)
	}

	if err := h.svc.IncrementView(ctx, itemId); err != nil {
		return fmt.Errorf("failed to increment views: %w", err)
	}

	return c.NoContent(http.StatusOK)
}

// Like godoc
//
//	@Summary	Like an item
//	@Tags		item
//	@Param		id	path	int		true	"Item ID"
//	@Success	200	{array}	int64	"Updated liked items"
//	@Router		/item/{id}/like [post]
func (h *Handler) Like(c echo.Context) error {
	var (
		req  = c.Request()
		ctx  = req.Context()
		user = c.(shared.Context).User
	)
	defer req.Body.Close()

	var itemId int64
	if err := echo.PathParamsBinder(c).Int64("id", &itemId).BindError(); err != nil {
		return fmt.Errorf("failed to bind item id: %w", err)
	}

	likedItems, err := h.svc.Like(ctx, itemId, user.ID)
	if err != nil {
		return fmt.Errorf("failed to like item: %w", err)
	}

	return c.JSON(http.StatusOK, likedItems)
}

// Unlike godoc
//
//	@Summary	Unlike an item
//	@Tags		item
//	@Param		id	path	int		true	"Item ID"
//	@Success	200	{array}	int64	"Updated liked items"
//	@Router		/item/{id}/unlike [post]
func (h *Handler) Unlike(c echo.Context) error {
	var (
		req  = c.Request()
		ctx  = req.Context()
		user = c.(shared.Context).User
	)
	defer req.Body.Close()

	var itemId int64
	if err := echo.PathParamsBinder(c).Int64("id", &itemId).BindError(); err != nil {
		return fmt.Errorf("failed to bind item id: %w", err)
	}

	likedItems, err := h.svc.Unlike(ctx, itemId, user.ID)
	if err != nil {
		return fmt.Errorf("failed to unlike item: %w", err)
	}

	return c.JSON(http.StatusOK, likedItems)
}

// GetBulk godoc
//
//	@Summary	Get items in bulk by list of IDs
//	@Tags		item
//	@Accept		json
//	@Produce	json
//	@Param		item_IDs	body	[]int64	true	"Item IDs"
//	@Success	200			{array}	itemsvc.Item
//	@Router		/item/bulk [post]
func (h *Handler) GetBulk(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)
	defer req.Body.Close()

	var itemIDs []int64
	if err := c.Bind(&itemIDs); err != nil {
		return fmt.Errorf("failed to bind item ids: %w", err)
	}

	items, err := h.svc.GetBulk(ctx, itemIDs)
	if err != nil {
		return fmt.Errorf("failed to get items: %w", err)
	}

	return c.JSON(http.StatusOK, items)
}
