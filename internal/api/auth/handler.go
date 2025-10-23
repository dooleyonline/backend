package authhandler

import (
	"fmt"
	"net/http"
	"time"

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
//	@Success	200		{object}	usersvc.Me			"User"
//	@Router		/auth/login [post]
func (h *Handler) Login(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
	)
	defer req.Body.Close()

	var params authsvc.LoginParams
	if err := c.Bind(&params); err != nil {
		return fmt.Errorf("failed to bind params: %w", err)
	}

	resp, err := h.svc.Login(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	cookieDetails, err := h.svc.CookieDetails()
	if err != nil {
		return fmt.Errorf("failed to get cookie details: %w", err)
	}

	cookie := &http.Cookie{
		Name:     cookieDetails.AuthTokenName,
		Value:    resp.Token,
		Expires:  time.Now().Add(cookieDetails.AuthTokenExp),
		Path:     "/",
		HttpOnly: true,
		Secure:   cookieDetails.Secure,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, resp.User)
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
	cookieDetails, err := h.svc.CookieDetails()
	if err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}

	cookie := &http.Cookie{
		Name:     cookieDetails.AuthTokenName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   cookieDetails.Secure,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)

	return c.NoContent(http.StatusOK)
}
