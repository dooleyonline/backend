package itemapi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dooleyonline/backend/internal/api/shared"
	"github.com/dooleyonline/backend/sql"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

// GetAll godoc
//
//	@Summary	Get all items
//	@Tags		item
//	@Produce	json
//	@Success	200	{array}	sql.Item
//	@Router		/item [get]
func GetAll(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
		db  = c.(shared.Context).DB
	)
	defer req.Body.Close()

	items, err := db.GetAllItems(ctx)
	if err != nil {
		return fmt.Errorf("failed to get all items: %w", err)
	}

	return c.JSON(http.StatusOK, items)
}

// Get godoc
//
//	@Summary	Get item by ID
//	@Tags		item
//	@Produce	json
//	@Param		id	path		int	true	"Item ID"
//	@Success	200	{object}	sql.Item
//	@Router		/item/{id} [get]
func Get(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
		db  = c.(shared.Context).DB
	)
	defer req.Body.Close()

	var id int64
	if err := echo.PathParamsBinder(c).Int64("id", &id).BindError(); err != nil {
		return fmt.Errorf("failed to bind id: %w", err)
	}

	item, err := db.GetItem(ctx, id)
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
//	@Param		item	body		sql.CreateItemParams	true	"Item"
//	@Success	200		{object}	sql.Item
//	@Router		/item [post]
func Create(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
		db  = c.(shared.Context).DB
	)
	defer req.Body.Close()

	var params sql.CreateItemParams
	if err := c.Bind(&params); err != nil {
		return fmt.Errorf("failed to bind params: %w", err)
	}

	item, err := db.CreateItem(ctx, params)
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
//	@Param		item	body		sql.UpdateItemParams	true	"Item"
//	@Success	200		{object}	sql.Item
//	@Router		/item/{id} [put]
func Update(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
		db  = c.(shared.Context).DB
	)
	defer req.Body.Close()

	var params sql.UpdateItemParams
	if err := c.Bind(&params); err != nil {
		return fmt.Errorf("failed to bind params: %w", err)
	}
	if err := echo.PathParamsBinder(c).Int64("id", &params.ID).BindError(); err != nil {
		return fmt.Errorf("failed to bind id: %w", err)
	}

	item, err := db.UpdateItem(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return c.JSON(http.StatusOK, item)
}

// Sell godoc
//
//	@Summary	Update sold_at property by ID
//	@Tags		item
//	@Param		id	path		int		true	"Item ID"
//	@Success	200	{string}	string	"Item ID"
//	@Router		/item/{id}/sell [post]
func Sell(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
		db  = c.(shared.Context).DB
	)
	defer req.Body.Close()

	var params sql.SellItemParams
	if err := echo.PathParamsBinder(c).Int64("id", &params.ID).BindError(); err != nil {
		return fmt.Errorf("failed to bind id: %w", err)
	}
	params.SoldAt = pgtype.Timestamptz{
		Time:  time.Now(),
		Valid: true,
	}

	if err := db.SellItem(ctx, params); err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return c.NoContent(http.StatusOK)
}

// IncrementView godoc
//
//	@Summary	Increment item views by ID
//	@Tags		item
//	@Param		id	path		int		true	"Item ID"
//	@Success	200	{string}	string	"Item ID"
//	@Router		/item/{id}/view [post]
func IncrementView(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
		db  = c.(shared.Context).DB
	)
	defer req.Body.Close()

	var id int64
	if err := echo.PathParamsBinder(c).Int64("id", &id).BindError(); err != nil {
		return fmt.Errorf("failed to bind id: %w", err)
	}

	if err := db.IncrementItemView(ctx, id); err != nil {
		return fmt.Errorf("failed to increment views: %w", err)
	}

	return c.NoContent(http.StatusOK)
}

// Delete godoc
//
//	@Summary	Delete item by ID
//	@Tags		item
//	@Param		id	path		int		true	"Item ID"
//	@Success	200	{string}	string	"Item ID"
//	@Router		/item/{id} [delete]
func Delete(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
		db  = c.(shared.Context).DB
	)
	defer req.Body.Close()

	var id int64
	if err := echo.PathParamsBinder(c).Int64("id", &id).BindError(); err != nil {
		return fmt.Errorf("failed to bind id: %w", err)
	}

	if err := db.DeleteItem(ctx, id); err != nil {
		return fmt.Errorf("failed to get item: %w", err)
	}

	return c.NoContent(http.StatusOK)
}

// TODO: searchItem
