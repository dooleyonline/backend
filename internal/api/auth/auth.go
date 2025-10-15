package authapi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dooleyonline/backend/internal/api/shared"
	"github.com/labstack/echo/v4"
)

const (
	jwtCookieName = "dooleyonline-jwt"
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

	token, err := shared.CreateJWT(params.Email, cfg.HmacSecretKey)
	if err != nil {
		return fmt.Errorf("failed to create token: %w", err)
	}

	cookie := new(http.Cookie)
	cookie.Name = jwtCookieName
	cookie.Value = token
	cookie.Expires = time.Now().Add(240 * time.Hour)
	cookie.HttpOnly = true
	cookie.Secure = cfg.Env == "prod"

	// cookie := &http.Cookie{
	// 	Name:     jwtCookieName,
	// 	Value:    token,
	// 	Expires:  time.Now().Add(240 * time.Hour),
	// 	HttpOnly: true,
	// 	Secure:   true,
	// 	SameSite: http.SameSiteLaxMode,
	// }

	c.SetCookie(cookie)

	return c.String(http.StatusOK, "Logged in")
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
		req = c.Request()
		cfg = c.(shared.Context).Cfg
	)
	defer req.Body.Close()

	cookie := new(http.Cookie)
	cookie.Name = jwtCookieName
	cookie.Value = ""
	cookie.Expires = time.Unix(0, 0)
	cookie.HttpOnly = true
	cookie.Secure = cfg.Env == "prod"

	c.SetCookie(cookie)

	return c.String(http.StatusOK, "Logged out")
}
