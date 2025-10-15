package authapi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dooleyonline/backend/internal/api/shared"
	"github.com/labstack/echo/v4"

	sqluser "github.com/dooleyonline/backend/sql/user"
)

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

	hash, err := hashPassword(params.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	params.Password = string(hash)

	user, err := db.User.Create(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}

	return c.JSON(http.StatusOK, user)
}

type loginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login1(c echo.Context) error {
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

	if verified := verifyPassword(params.Password, user.Password); !verified {
		return c.NoContent(http.StatusUnauthorized)
	}

	token, err := createToken(params.Email, cfg.HmacSecretKey)
	if err != nil {
		return fmt.Errorf("failed to create token: %w", err)
	}

	cookie := &http.Cookie{
		Name:     "Token",
		Value:    token,
		Expires:  time.Now().Add(240 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	c.SetCookie(cookie)

	return c.String(http.StatusOK, "Logged in")
}

func Logout1(c echo.Context) error {
	cookie := &http.Cookie{
		Name:     "Token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)

	return c.String(http.StatusOK, "Logged out")
}
