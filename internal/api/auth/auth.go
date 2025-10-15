package authapi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dooleyonline/backend/internal/api/shared"
	"github.com/labstack/echo/v4"
)

type loginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login godoc
//
//	@Summary	Log in
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		user	body		loginParams	true	"Login Params"
//	@Success	200		{string}	string		"Result"
//	@Router		/auth/login [post]
func Login(c echo.Context) error {
	var (
		req = c.Request()
		ctx = req.Context()
		db  = c.(shared.Context).DB
		cfg = c.(shared.Context).Cfg
	)
	defer req.Body.Close()

	var params loginParams
	if err := c.Bind(&params); err != nil {
		return fmt.Errorf("failed to bind params: %w", err)
	}

	user, err := db.User.Get(ctx, params.Email)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if verified := shared.VerifyPassword(params.Password, user.Password); !verified {
		return c.NoContent(http.StatusUnauthorized)
	}

	token, err := shared.CreateJWT(cfg, params.Email)
	if err != nil {
		return fmt.Errorf("failed to create token: %w", err)
	}

	cookie := &http.Cookie{
		Name:     cfg.AuthTokenName,
		Value:    token,
		Expires:  time.Now().Add(cfg.AuthTokenExp),
		Path:     "/",
		HttpOnly: true,
		Secure:   cfg.IsProd,
		SameSite: http.SameSiteLaxMode,
	}

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, user)
}

// Logout godoc
//
//	@Summary	Log out
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Success	200	{string}	string	"Result"
//	@Router		/auth/logout [post]
func Logout(c echo.Context) error {
	var (
		req  = c.Request()
		cfg  = c.(shared.Context).Cfg
		user = c.(shared.Context).User
	)
	defer req.Body.Close()

	cookie := &http.Cookie{
		Name:     cfg.AuthTokenName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   cfg.IsProd,
		SameSite: http.SameSiteLaxMode,
	}

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, user)
}
