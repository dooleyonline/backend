package authhandler

import (
	"net/http"
	"time"

	"github.com/dooleyonline/backend/internal/api/shared"
	authsvc "github.com/dooleyonline/backend/internal/service/auth"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc *authsvc.Service
}

func New(svc *authsvc.Service) *Handler {
	return &Handler{svc}
}

// Login godoc
//
//	@Summary	Log in
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		user	body		authsvc.LoginParams	true	"Login Params"
//	@Success	200		{object}	model.User			"User"
//	@Router		/auth/login [post]
func (h *Handler) Login(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	var params authsvc.LoginParams
	if err := c.Bind(&params); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	res, err := h.svc.Login(ctx, &params)
	if err != nil {
		return echo.ErrUnauthorized.WithInternal(err)
	}

	cookie := &http.Cookie{
		Name:     res.CookieConfig.AuthTokenName,
		Value:    res.Token,
		Expires:  time.Now().Add(res.CookieConfig.AuthTokenExp),
		Path:     "/",
		HttpOnly: true,
		Secure:   res.CookieConfig.Secure,
		SameSite: http.SameSiteNoneMode, // TODO: change to LAX in production
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, res.User)
}

// Logout godoc
//
//	@Summary	Log out
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Success	200	"OK"
//	@Router		/auth/logout [post]
func (h *Handler) Logout(c echo.Context) error {
	res, err := h.svc.CookieOptions()
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	cookie := &http.Cookie{
		Name:     res.AuthTokenName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   res.Secure,
		SameSite: http.SameSiteNoneMode, // TODO: change to LAX in production
	}
	c.SetCookie(cookie)

	return c.NoContent(http.StatusOK)
}

// GetMe godoc
//
//	@Summary	Get current authenticated user
//	@Tags		auth
//	@Produce	json
//	@Success	200	{object}	model.User
//	@Failure	401	"Unauthorized"
//	@Router		/auth/me [get]
func (h *Handler) GetMe(c echo.Context) error {
	var (
		req    = c.Request()
		ctx    = req.Context()
		userId = c.(shared.Context).UserID
	)

	if userId == "" {
		return c.JSON(http.StatusOK, nil)
	}

	res, err := h.svc.GetMe(ctx, userId)
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(http.StatusOK, res)
}
