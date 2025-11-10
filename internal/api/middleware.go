package api

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"

	"github.com/dooleyonline/backend/internal/api/shared"
	"github.com/dooleyonline/backend/internal/config"
	authsvc "github.com/dooleyonline/backend/internal/service/auth"
	"github.com/golang-jwt/jwt/v5"
	echoJWT "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func loggerMiddleware() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			ctx := c.Request().Context()
			if v.Error == nil {
				slog.LogAttrs(ctx, slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				slog.LogAttrs(ctx, slog.LevelError, "REQUEST_ERROR",
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

func contextMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := shared.Context{
				Context: c,
				UserID:  "",
			}

			defer fmt.Println("user", cc.UserID, "request", cc.Request().RequestURI)

			cookie, err := c.Cookie(cfg.AuthTokenName)
			if err != nil {
				return next(cc)
			}

			var claims authsvc.JWTClaims
			token, err := jwt.ParseWithClaims(
				cookie.Value,
				&claims,
				func(t *jwt.Token) (any, error) {
					return []byte(cfg.AuthTokenSecret), nil
				},
			)
			if err != nil || !token.Valid {
				return next(cc)
			}

			cc.UserID = claims.ID

			return next(cc)
		}
	}
}

func corsMiddleware() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins:     []string{"https://dooleyonline.vercel.app", "http://localhost:3000", "https://dooleyonline.net"},
			AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderSetCookie},
			AllowCredentials: true,
		},
	)
}

func authMiddleware(cfg *config.Config, protectedRoutes routesConfig) echo.MiddlewareFunc {
	isPublic := func(c echo.Context) bool {
		r, ok := protectedRoutes[c.Path()]
		if !ok {
			return true
		}
		return !slices.Contains(r, c.Request().Method)
	}

	tokenLookup := func(c echo.Context) ([]string, error) {
		tokenCookie, err := c.Cookie(cfg.AuthTokenName)
		if err != nil {
			return nil, fmt.Errorf("failed to get token cookie: %w", err)
		}

		return []string{tokenCookie.Value}, nil
	}

	config := echoJWT.Config{
		Skipper:    isPublic,
		SigningKey: []byte(cfg.AuthTokenSecret),
		ContextKey: "user",
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(authsvc.JWTClaims)
		},
		TokenLookupFuncs: []middleware.ValuesExtractor{tokenLookup},
	}

	return echoJWT.WithConfig(config)
}
