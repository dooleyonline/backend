package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/dooleyonline/backend/internal/db"
	"github.com/dooleyonline/backend/sql"
	"github.com/labstack/echo/v4"
)

// getItem godoc
//
//	@Summary		Get item
//	@Description	Get all items or by ID
//	@Tags			item
//	@Produce		json
//	@Param			id	path		int	false	"Item ID"
//	@Success		200	{object}	sql.Item
//	@Router			/item/{id} [get]
func getItem(ctx context.Context, db *db.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		idParam := c.Param("id")

		// return all items if no id provided
		if idParam == "" {
			items, err := db.GetAllItems(ctx)
			if err != nil {
				return fmt.Errorf("failed to get all items: %w", err)
			}

			return c.JSON(http.StatusOK, items)
		}

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
}

// createItem godoc
//
//	@Summary		Create item
//	@Description	Create item
//	@Tags			item
//	@Accept			json
//	@Produce		json
//	@Param			item	body		sql.CreateItemParams	true	"Item"
//	@Success		200		{object}	sql.Item
//	@Router			/item [post]
func createItem(ctx context.Context, db *db.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		defer req.Body.Close()

		body, err := io.ReadAll(req.Body)
		if err != nil {
			return fmt.Errorf("failed to read request body: %w", err)
		}

		params := sql.CreateItemParams{
			// Images: []string{},
		}

		if err := json.Unmarshal(body, &params); err != nil {
			return fmt.Errorf("failed to unmarshal body: %w", err)
		}

		item, err := db.CreateItem(ctx, params)
		if err != nil {
			return fmt.Errorf("failed to create item: %w", err)
		}

		return c.JSON(http.StatusOK, item)
	}
}

// updateItem godoc
//
//	@Summary		Update item
//	@Description	Update item by ID
//	@Tags			item
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						false	"Item ID"
//	@Param			item	body		sql.UpdateItemParams	true	"Item"
//	@Success		200		{object}	sql.Item
//	@Router			/item/{id} [put]
func updateItem(ctx context.Context, db *db.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
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
}

// deleteItem godoc
//
//	@Summary		Delete item
//	@Description	Delete item by ID
//	@Tags			item
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						false	"Item ID"
//	@Param			item	body		sql.UpdateItemParams	true	"Item"
//	@Success		200		{object}	sql.Item
//	@Router			/item/{id} [delete]
func deleteItem(ctx context.Context, db *db.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		idParam := c.Param("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse id: %w", err)
		}

		if err := db.DeleteItem(ctx, id); err != nil {
			return fmt.Errorf("failed to get item: %w", err)
		}

		return c.JSON(http.StatusOK, id)
	}
}

// TODO: incrementView, searchItem
