package api

import (
	"context"
	"log/slog"
	"net/http"

	itemapi "github.com/dooleyonline/backend/internal/api/item"
	storageapi "github.com/dooleyonline/backend/internal/api/storage"
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
func New(ctx context.Context, cfg *config.Config, db *db.DB, lg *slog.Logger) (*echo.Echo, error) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(loggerMiddleware(lg))
	e.Use(errorMiddleware())
	e.Use(contextMiddleware(cfg, db))
	e.Use(corsMiddleware())

	// TODO: auth middleware

	e.GET("/", hello)

	// swagger docs
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// item routes
	e.GET("/item", itemapi.GetAll)
	e.POST("/item", itemapi.Create)
	e.GET("/item/:id", itemapi.Get)
	e.PUT("/item/:id", itemapi.Update)
	e.DELETE("/item/:id", itemapi.Delete)
	e.POST("/item/:id/view", itemapi.IncrementView)
	e.POST("/item/:id/sell", itemapi.Sell)
	e.GET("/item/search", itemapi.Search)

	// TODO: category routes

	// TODO: user routes

	e.POST("/storage/presign", storageapi.Presign)

	return e, nil
}

// hello godoc
//
//	@Summary	Greeting
//	@Produce	plain
//	@Success	200	{string}	string
//	@Router		/ [get]
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome! This is DooleyOnline.")
}
