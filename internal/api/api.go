package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/internal/db"
	"github.com/labstack/echo/v4"

	_ "github.com/dooleyonline/backend/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

//	@title			DooleyOnline API
//	@version		1.0
//	@description	This is the REST API for DooleyOnline.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8080
// @BasePath	/
func New(ctx context.Context, cfg *config.Config, db *db.DB, logger *slog.Logger) (*echo.Echo, error) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(loggerMiddleware(logger))
	e.Use(errorMiddleware())
	e.Use(contextMiddleware(cfg, db))

	// TODO: auth middleware

	e.GET("/", hello())

	// swagger docs
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// item routes
	e.GET("/item", getAllItems)
	e.POST("/item", createItem)
	e.GET("/item/:id", getItem)
	e.PUT("/item/:id", updateItem)
	e.DELETE("/item/:id", deleteItem)
	e.POST("/item/:id/views", incrementItemViews)

	// TODO: category routes

	// TODO: user routes

	e.POST("/storage/presign-upload", presignUpload)

	return e, nil
}

func hello() echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome! This is DooleyOnline.")
	})
}
