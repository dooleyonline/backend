package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/dooleyonline/backend/sql"
	"github.com/labstack/echo/v4"
)

// getAllItems godoc
//
//	@Summary	Get all items
//	@Tags		item
//	@Produce	json
//	@Success	200	{array}	sql.Item
//	@Router		/item [get]
func getAllItems(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
		db  = c.(customContext).db
	)
	defer req.Body.Close()

	items, err := db.GetAllItems(ctx)
	if err != nil {
		return fmt.Errorf("failed to get all items: %w", err)
	}

	return c.JSON(http.StatusOK, items)
}

// getItem godoc
//
//	@Summary	Get item by ID
//	@Tags		item
//	@Produce	json
//	@Param		id	path		int	true	"Item ID"
//	@Success	200	{object}	sql.Item
//	@Router		/item/{id} [get]
func getItem(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
		db  = c.(customContext).db
	)
	defer req.Body.Close()

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse id: %w", err)
	}

	item, err := db.GetItem(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get item: %w", err)
	}

	return c.JSON(http.StatusOK, item)
}

// createItem godoc
//
//	@Summary	Create item
//	@Tags		item
//	@Accept		json
//	@Produce	json
//	@Param		item	body		sql.CreateItemParams	true	"Item"
//	@Success	200		{object}	sql.Item
//	@Router		/item [post]
func createItem(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
		db  = c.(customContext).db
	)
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}

	params := sql.CreateItemParams{}

	if err := json.Unmarshal(body, &params); err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}

	item, err := db.CreateItem(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}

	return c.JSON(http.StatusOK, item)
}

// updateItem godoc
//
//	@Summary	Update item by ID
//	@Tags		item
//	@Accept		json
//	@Produce	json
//	@Param		id		path		int						true	"Item ID"
//	@Param		item	body		sql.UpdateItemParams	true	"Item"
//	@Success	200		{object}	sql.Item
//	@Router		/item/{id} [put]
func updateItem(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
		db  = c.(customContext).db
	)
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}

	params := sql.UpdateItemParams{}

	if err := json.Unmarshal(body, &params); err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}

	item, err := db.UpdateItem(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return c.JSON(http.StatusOK, item)
}

// incrementItemViews godoc
//
//	@Summary	Increment item views by ID
//	@Tags		item
//	@Produce	plain
//	@Param		id	path		int		true	"Item ID"
//	@Success	200	{string}	string	"Item ID"
//	@Router		/item/{id}/views [post]
func incrementItemViews(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
		db  = c.(customContext).db
	)
	defer req.Body.Close()

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse id: %w", err)
	}

	if err := db.IncrementItemView(ctx, id); err != nil {
		return fmt.Errorf("failed to increment views: %w", err)
	}

	return c.String(http.StatusOK, idParam)
}

// deleteItem godoc
//
//	@Summary	Delete item by ID
//	@Tags		item
//	@Produce	plain
//	@Param		id	path		int		true	"Item ID"
//	@Success	200	{string}	string	"Item ID"
//	@Router		/item/{id} [delete]
func deleteItem(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
		db  = c.(customContext).db
	)
	defer req.Body.Close()

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse id: %w", err)
	}

	if err := db.DeleteItem(ctx, id); err != nil {
		return fmt.Errorf("failed to get item: %w", err)
	}

	return c.String(http.StatusOK, idParam)
}

// TODO: searchItem
