package api

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func errorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := next(c); err != nil {
				req := c.Request()

				reqAttr := slog.Group(
					"request",
					slog.String("method", req.Method),
					slog.String("url", req.RequestURI),
				)

				slog.Error("API error", reqAttr, slog.Any("error", err))
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}
			return nil
		}
	}

}
