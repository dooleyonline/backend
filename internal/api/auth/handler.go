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

//	CreateVerification godoc
//
// @Summary		Send verification email
// @Description	Sends (or resends) a verification link to the provided email. Returns 204 on success.
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			request	body	authsvc.SendParams	true	"Email payload"
// @Success		204
// @Failure		400	{object}	map[string]string	"invalid request body"
// @Failure		500	{object}	map[string]string	"internal error"
// @Router			/auth/verification [post]
func (h *Handler) CreateVerification(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	var params authsvc.SendParams
	if err := c.Bind(&params); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	_, err := h.svc.CreateVerification(ctx, params)

	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	return nil
}

// VerifyUser godoc
//
//	@Summary		Verify user by token
//	@Description	Consumes the verification token from the path and marks the user verified.
//	@Tags			auth
//	@Produce		json
//	@Param			id	path		string				true	"Verification token (UUID)"
//	@Success		200	{object}	map[string]string	"user verified"
//	@Failure		400	{object}	map[string]string	"invalid token format"
//	@Failure		401	{object}	map[string]string	"invalid or expired token"
//	@Failure		500	{object}	map[string]string	"internal error"
//	@Router			/auth/verification/{id} [post]
func (h *Handler) VerifyUser(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)

	var id string
	if err := echo.PathParamsBinder(c).String("id", &id).BindError(); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	err := h.svc.VerifyUserEmail(ctx, id)
	if err != nil {
		return echo.ErrUnauthorized.WithInternal(err)
	}

	return nil
}
