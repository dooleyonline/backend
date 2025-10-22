package user

import (
	"fmt"
	"net/http"

	"github.com/dooleyonline/backend/internal/api/shared"
	"github.com/dooleyonline/backend/internal/db/user"
	"github.com/labstack/echo/v4"
)

// GetMany godoc
//
//	@Summary	Get many users
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Success	200	{array}	userdb.User
//	@Router		/user [get]
func GetMany(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
		db  = c.(shared.Context).DB
	)
	defer req.Body.Close()

	users, err := db.User.GetMany(ctx)
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	return c.JSON(http.StatusOK, users)
}

// Create godoc
//
//	@Summary	Create user
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		user	body		userdb.CreateParams	true	"User"
//	@Success	200		{object}	userdb.User
//	@Router		/user [post]
func Create(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
		db  = c.(shared.Context).DB
	)
	defer req.Body.Close()

	var params userdb.CreateParams
	if err := c.Bind(&params); err != nil {
		return fmt.Errorf("failed to bind params: %w", err)
	}

	hash, err := shared.HashPassword(params.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	params.Password = string(hash)

	user, err := db.User.Create(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return c.JSON(http.StatusOK, user)
}

// GetSeller godoc
//
//	@Summary	Get seller by ID
//	@Tags		user
//	@Produce	json
//	@Param		id	path		string	true	"User ID (UUID)"
//	@Success	200	{object}	userdb.User
//	@Router		/user/{id} [get]
func GetSeller(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
		db  = c.(shared.Context).DB
	)
	defer req.Body.Close()

	id := c.Param("id")
	if id == "" {
		return fmt.Errorf("user id is required")
	}

	seller, err := db.User.GetSellerByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	return c.JSON(http.StatusOK, seller)
}

// GetMe godoc
//
//	@Summary	Get current authenticated user
//	@Tags		user
//	@Produce	json
//	@Success	200	{object}	userdb.User
//	@Router		/user/me [get]
//	@Security	ApiKeyAuth
func GetMe(c echo.Context) error {
	var (
		req  = c.Request()
		ctx  = req.Context()
		db   = c.(shared.Context).DB
		user = c.(shared.Context).User
	)
	defer req.Body.Close()

	me, err := db.User.GetFullUserByID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	return c.JSON(http.StatusOK, me)
}
