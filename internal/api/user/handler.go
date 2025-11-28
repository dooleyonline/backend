package user

import (
	"net/http"

	"github.com/dooleyonline/backend/internal/api/shared"
	usersvc "github.com/dooleyonline/backend/internal/service/user"
	"github.com/labstack/echo/v4"
)

type Handler struct {
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
//	@Success	200	{array}	model.User
//	@Router		/user [get]
func (h *Handler) GetMany(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	res, err := h.svc.GetMany(ctx)
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(http.StatusOK, res)
}

// Create godoc
//
//	@Summary	Create user
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		user	body		usersvc.CreateParams	true	"User"
//	@Success	201		{object}	model.User
//	@Router		/user [post]
func (h *Handler) Create(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	var params usersvc.CreateParams
	if err := c.Bind(&params); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	res, err := h.svc.Create(ctx, &params)
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(http.StatusOK, res)
}

// Get godoc
//
//	@Summary	Get user by ID
//	@Tags		user
//	@Produce	json
//	@Param		id	path		string	true	"User ID (UUID)"
//	@Success	200	{object}	model.User
//	@Router		/user/{id} [get]
func (h *Handler) Get(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	id := c.Param("id")
	if id == "" {
		return echo.ErrBadRequest
	}

	res, err := h.svc.Get(ctx, id)
	if err != nil {
		return echo.ErrNotFound.WithInternal(err)
	}

	return c.JSON(http.StatusOK, res)
}

// TODO: Update email verification if email is changed

// Update godoc
//
//	@Summary	Update user
//	@Tags		user
//	@Accept		json
//	@Produce	json
//	@Param		user	body	usersvc.UpdateParams	true	"User"
//	@Success	200		"OK"
//	@Router		/user [put]
func (h *Handler) Update(c echo.Context) error {
	var (
		req    = c.Request()
		ctx    = req.Context()
		userID = c.(shared.Context).UserID
	)

	var params usersvc.UpdateParams
	if err := c.Bind(&params); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}
	params.UserID = userID

	err := h.svc.Update(ctx, params)
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.NoContent(http.StatusOK)
}

// GetLikedViewed godoc
//
//	@Summary	Get the liked, viewed items group by userID
//	@Tags		user
//	@Produce	json
//	@Success	200	{object}	[]useruser.GetAllLikedViewedRow
//	@Router		/user/interactions [get]
func (h *Handler) GetLikedViewed(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	res, err := h.svc.GetLikedViewed(ctx)
	if err != nil {
		return echo.ErrNotFound.WithInternal(err)
	}

	return c.JSON(http.StatusOK, res)
}
