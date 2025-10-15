package user

import (
	"fmt"
	"net/http"

	"github.com/dooleyonline/backend/internal/api/shared"
	sqluser "github.com/dooleyonline/backend/sql/user"
	"github.com/labstack/echo/v4"
)

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
//	@Param		user	body		sqluser.CreateParams	true	"User"
//	@Success	200		{object}	sqluser.User
//	@Router		/user/create [post]
func Create(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
		db  = c.(shared.Context).DB
	)
	defer req.Body.Close()

	var params sqluser.CreateParams
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
