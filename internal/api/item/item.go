package itemapi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dooleyonline/backend/internal/api/shared"
	"github.com/labstack/echo/v4"

	sqlitem "github.com/dooleyonline/backend/sql/item"
	sqluser "github.com/dooleyonline/backend/sql/user"
)

// GetMany godoc
//
//	@Summary	Get many items
//	@Tags		item
//	@Produce	json
//	@Param		q			query	string	false	"Search query"
//	@param		category	query	string	false	"Category filter"
//	@Success	200			{array}	sqlitem.Item
//	@Router		/item [get]
func GetMany(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
		db  = c.(shared.Context).DB
	)
	defer req.Body.Close()

	var (
		items []sqlitem.Item
		err   error
	)

	query := c.QueryParam("q")
	category := c.QueryParam("category")

	switch {
	case query != "" && category != "":
		params := sqlitem.SearchByCategoryParams{
			Category:  category,
			ToTsquery: query,
		}
		items, err = db.Item.SearchByCategory(ctx, params)
	case query != "":
		items, err = db.Item.Search(ctx, query)
	case category != "":
		items, err = db.Item.GetByCategory(ctx, category)
	default:
		items, err = db.Item.GetAll(ctx)
	}

	if err != nil {
		return fmt.Errorf("failed to get items: %w", err)
	}
	if items == nil {
		items = []sqlitem.Item{}
	}
	return c.JSON(http.StatusOK, items)
}

// Get godoc
//
//	@Summary	Get item by ID
//	@Tags		item
//	@Produce	json
//	@Param		id	path		int	true	"Item ID"
//	@Success	200	{object}	sqlitem.Item
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

	item, err := db.Item.Get(ctx, id)
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
//	@Param		item	body		sqlitem.CreateParams	true	"Item"
//	@Success	200		{object}	sqlitem.Item
//	@Router		/item [post]
func Create(c echo.Context) error {
	var (
		req  = c.Request()
		ctx  = req.Context()
		db   = c.(shared.Context).DB
		cfg  = c.(shared.Context).Cfg
		user = c.(shared.Context).User
	)
	defer req.Body.Close()

	var params sqlitem.CreateParams
	if err := c.Bind(&params); err != nil {
		return fmt.Errorf("failed to bind params: %w", err)
	}

	params.Seller = user.ID

	placeholder, err := GeneratePlaceholder(cfg, params.Images[0])
	if err != nil {
		return fmt.Errorf("failed to generate placeholders: %w", err)
	}
	params.Placeholder = placeholder

	item, err := db.Item.Create(ctx, params)
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
//	@Param		item	body		sqlitem.UpdateParams	true	"Item"
//	@Success	200		{object}	sqlitem.Item
//	@Router		/item/{id} [put]
func Update(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
		db  = c.(shared.Context).DB
		cfg = c.(shared.Context).Cfg
	)
	defer req.Body.Close()

	var params sqlitem.UpdateParams
	if err := c.Bind(&params); err != nil {
		return fmt.Errorf("failed to bind params: %w", err)
	}
	if err := echo.PathParamsBinder(c).Int64("id", &params.ID).BindError(); err != nil {
		return fmt.Errorf("failed to bind id: %w", err)
	}

	placeholder, err := GeneratePlaceholder(cfg, params.Images[0])
	if err != nil {
		return fmt.Errorf("failed to generate placeholders: %w", err)
	}
	params.Placeholder = placeholder

	item, err := db.Item.Update(ctx, params)
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

	var params sqlitem.SellParams
	if err := echo.PathParamsBinder(c).Int64("id", &params.ID).BindError(); err != nil {
		return fmt.Errorf("failed to bind id: %w", err)
	}
	now := time.Now()
	params.SoldAt = &now

	if err := db.Item.Sell(ctx, params); err != nil {
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

	if err := db.Item.IncrementView(ctx, id); err != nil {
		return fmt.Errorf("failed to increment views: %w", err)
	}

	return c.NoContent(http.StatusOK)
}

// Delete godoc
//
//	@Summary	Delete item by ID
//	@Tags		item
//	@Param		id	path	int	true	"Item ID"
//	@Success	200
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

	if err := db.Item.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to get item: %w", err)
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
func Like(c echo.Context) error {
	var (
		req  = c.Request()
		ctx  = req.Context()
		db   = c.(shared.Context).DB
		user = c.(shared.Context).User
	)
	defer req.Body.Close()

	var itemId int64
	if err := echo.PathParamsBinder(c).Int64("id", &itemId).BindError(); err != nil {
		return fmt.Errorf("failed to bind item id: %w", err)
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	userTx := db.User.WithTx(tx)
	itemTx := db.Item.WithTx(tx)

	likedItems, err := userTx.AddLikedItem(ctx, sqluser.AddLikedItemParams{
		ItemID: itemId,
		Email:  user.Email,
	})
	if err != nil {
		return fmt.Errorf("failed to like item: %w", err)
	}

	if err := itemTx.IncrementLike(ctx, itemId); err != nil {
		return fmt.Errorf("failed to increment item like: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
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
func Unlike(c echo.Context) error {
	var (
		req  = c.Request()
		ctx  = req.Context()
		db   = c.(shared.Context).DB
		user = c.(shared.Context).User
	)
	defer req.Body.Close()

	var itemId int64
	if err := echo.PathParamsBinder(c).Int64("id", &itemId).BindError(); err != nil {
		return fmt.Errorf("failed to bind item id: %w", err)
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	userTx := db.User.WithTx(tx)
	itemTx := db.Item.WithTx(tx)

	likedItems, err := userTx.DeleteLikedItem(ctx, sqluser.DeleteLikedItemParams{
		ItemID: itemId,
		Email:  user.Email,
	})
	if err != nil {
		return fmt.Errorf("failed to unlike item: %w", err)
	}

	if err := itemTx.DecrementLike(ctx, itemId); err != nil {
		return fmt.Errorf("failed to decrement item like: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return c.JSON(http.StatusOK, likedItems)
}
