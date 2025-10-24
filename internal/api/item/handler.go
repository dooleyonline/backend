package itemhandler

import (
	"net/http"

	"github.com/dooleyonline/backend/internal/api/shared"
	itemsvc "github.com/dooleyonline/backend/internal/service/item"
	"github.com/dooleyonline/backend/internal/storage"
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
//	@Success	200			{array}	itemdb.Item
//	@Router		/item [get]
func (h *Handler) GetMany(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	var (
		seller   string
		query    string
		category string
	)
	if err := echo.QueryParamsBinder(c).String("seller", &seller).String("q", &query).String("category", &category).BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	res, err := h.svc.GetMany(ctx, &itemsvc.GetManyParams{
		Seller:   seller,
		Query:    query,
		Category: category,
	})
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(http.StatusOK, res)
}

// Get godoc
//
//	@Summary	Get item by ID
//	@Tags		item
//	@Produce	json
//	@Param		id	path		int	true	"Item ID"
//	@Success	200	{object}	itemdb.Item
//	@Router		/item/{id} [get]
func (h *Handler) Get(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	var itemId int64
	if err := echo.PathParamsBinder(c).Int64("id", &itemId).BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	res, err := h.svc.Get(ctx, itemId)
	if err != nil {
		return echo.ErrNotFound.WithInternal(err)
	}

	return c.JSON(http.StatusOK, res)
}

// Create godoc
//
//	@Summary	Create item
//	@Tags		item
//	@Accept		json
//	@Produce	json
//	@Param		item	body		itemsvc.MutationParams	true	"Item"
//	@Success	201		{object}	itemdb.Item
//	@Router		/item [post]
func (h *Handler) Create(c echo.Context) error {
	var (
		req    = c.Request()
		ctx    = req.Context()
		userId = c.(shared.Context).UserID
	)

	var params itemsvc.MutationParams
	if err := c.Bind(&params); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	res, err := h.svc.Create(ctx, userId, &params)
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(http.StatusCreated, res)
}

// Update godoc
//
//	@Summary	Update item by ID
//	@Tags		item
//	@Accept		json
//	@Produce	json
//	@Param		id		path		int						true	"Item ID"
//	@Param		item	body		itemsvc.MutationParams	true	"Item"
//	@Success	200		{object}	itemdb.Item
//	@Router		/item/{id} [put]
func (h *Handler) Update(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	var params itemsvc.MutationParams
	if err := c.Bind(&params); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	var itemId int64
	if err := echo.PathParamsBinder(c).Int64("id", &itemId).BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	res, err := h.svc.Update(ctx, itemId, &params)
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(http.StatusOK, *res)
}

// Delete godoc
//
//	@Summary	Delete item by ID
//	@Tags		item
//	@Param		id	path	int	true	"Item ID"
//	@Success	204
//	@Router		/item/{id} [delete]
func (h *Handler) Delete(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	var itemId int64
	if err := echo.PathParamsBinder(c).Int64("id", &itemId).BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := h.svc.Delete(ctx, itemId); err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// Sell godoc
//
//	@Summary	Update sold_at property by ID
//	@Tags		item
//	@Param		id	path	int	true	"Item ID"
//	@Success	204
//	@Router		/item/{id}/sell [post]
func (h *Handler) Sell(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	var itemId int64
	if err := echo.PathParamsBinder(c).Int64("id", &itemId).BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := h.svc.Sell(ctx, itemId); err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// View godoc
//
//	@Summary	Increment item views by ID
//	@Tags		item
//	@Param		id	path	int	true	"Item ID"
//	@Success	204
//	@Router		/item/{id}/view [post]
func (h *Handler) View(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	var itemId int64
	if err := echo.PathParamsBinder(c).Int64("id", &itemId).BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := h.svc.IncrementView(ctx, itemId); err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// Like godoc
//
//	@Summary	Like an item
//	@Tags		item
//	@Param		id	path	int	true	"Item ID"
//	@Success	204
//	@Router		/item/{id}/like [post]
func (h *Handler) Like(c echo.Context) error {
	var (
		req    = c.Request()
		ctx    = req.Context()
		userId = c.(shared.Context).UserID
	)

	var itemId int64
	if err := echo.PathParamsBinder(c).Int64("id", &itemId).BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := h.svc.Like(ctx, itemId, userId); err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// Unlike godoc
//
//	@Summary	Unlike an item
//	@Tags		item
//	@Param		id	path	int	true	"Item ID"
//	@Success	204
//	@Router		/item/{id}/unlike [post]
func (h *Handler) Unlike(c echo.Context) error {
	var (
		req    = c.Request()
		ctx    = req.Context()
		userId = c.(shared.Context).UserID
	)

	var itemId int64
	if err := echo.PathParamsBinder(c).Int64("id", &itemId).BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := h.svc.Unlike(ctx, itemId, userId); err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// GetBulk godoc
//
//	@Summary	Get items in bulk by list of IDs
//	@Tags		item
//	@Accept		json
//	@Produce	json
//	@Param		item_IDs	body	[]int64	true	"Item IDs"
//	@Success	200			{array}	itemdb.Item
//	@Router		/item/batch [post]
func (h *Handler) GetBatch(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	var itemIDs []int64
	if err := c.Bind(&itemIDs); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	res, err := h.svc.GetBatch(ctx, &itemIDs)
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(http.StatusOK, res)
}

// GetUploadPresignURL godoc
//
//	@Summary	Generate presigned URL for item upload
//	@Tags		item
//	@Accept		json
//	@Produce	json
//	@Param		item	body		storage.PresignParams	true	"Presign params"
//	@Success	200		{object}	storage.PresignResult
//	@Router		/item/upload-url [post]
func (h *Handler) GetUploadURL(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	var input storage.PresignParams
	if err := c.Bind(&input); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if input.Method != http.MethodPut {
		return echo.ErrBadRequest
	}

	if input.Bucket == "" || input.Key == "" || input.ContentType == "" {
		return echo.ErrBadRequest
	}

	res, err := h.svc.GetUploadPresignURL(ctx, &input)
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(http.StatusOK, res)
}
