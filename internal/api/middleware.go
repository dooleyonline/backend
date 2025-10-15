package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dooleyonline/backend/internal/api/shared"
	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/internal/db"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func loggerMiddleware(lg *slog.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				lg.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				lg.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	})
}

func errorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := next(c); err != nil {
				// 404 if no corresponding row in db
				if errors.Is(err, sql.ErrNoRows) {
					return echo.NewHTTPError(http.StatusNotFound, err.Error())
				}
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
			return nil
		}
	}
}

func contextMiddleware(cfg *config.Config, db *db.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := shared.Context{
				Context: c,
				Cfg:     cfg,
				DB:      db,
			}
			return next(cc)
		}
	}
}

func corsMiddleware() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		},
	)
}

func authMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// token_string := c.Request().Header.Get("jwt-token")

			// _, err := validateJWT(token_string, cfg)
			// if err != nil {
			// 	return echo.NewHTTPError(http.StatusForbidden, err)
			// }
			token, _ := c.Cookie("Token")

			fmt.Println("Printing token:", token)

			return next(c)
		}
	}
}

func validateJWT(tokenString string, cfg *config.Config) (*jwt.Token, error) {
	hmacSecretKey := []byte(cfg.HmacSecretKey)

	// Source: https://pkg.go.dev/github.com/golang-jwt/jwt/v5#example-Parse-Hmac
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return hmacSecretKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		slog.Error("Invalid JWT: " + err.Error())
		return nil, err
	}
	return token, nil
}
