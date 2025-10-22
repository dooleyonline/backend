package user

import (
	"fmt"
	"net/http"

	"github.com/dooleyonline/backend/internal/api/shared"
	usersvc "github.com/dooleyonline/backend/internal/service/user"
	"github.com/labstack/echo/v4"
)

type Handler struct{
	svc *usersvc.Service
}

func New(svc *usersvc.Service) *Handler {
	return &Handler{svc}
}

// GetMany godoc
//
//	@Summary	Get many users
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Success	200	{array}	usersvc.UserSummary
//	@Router		/user [get]
func (h *Handler) GetMany(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)
	defer req.Body.Close()

	users, err := h.svc.GetMany(ctx)
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
//	@Param		user	body		usersvc.CreateInput	true	"User"
//	@Success	201		{object}	usersvc.User
//	@Router		/user [post]
func (h *Handler) Create(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)
	defer req.Body.Close()

	var params usersvc.CreateInput
	if err := c.Bind(&params); err != nil {
		return fmt.Errorf("failed to bind params: %w", err)
	}

	user, err := h.svc.Create(ctx, params)
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
//	@Success	200	{object}	usersvc.Seller
//	@Router		/user/{id} [get]
func (h *Handler) GetSeller(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)
	defer req.Body.Close()

	id := c.Param("id")
	if id == "" {
		return fmt.Errorf("user id is required")
	}

	seller, err := h.svc.GetSellerByID(ctx, id)
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
//	@Success	200	{object}	usersvc.Me
//	@Router		/user/me [get]
//	@Security	ApiKeyAuth
func (h *Handler) GetMe(c echo.Context) error {
	var (
		req  = c.Request()
		ctx  = req.Context()
		user = c.(shared.Context).User
	)
	defer req.Body.Close()

	me, err := h.svc.GetMe(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	return c.JSON(http.StatusOK, me)
}
