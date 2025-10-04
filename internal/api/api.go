package api

import (
	"context"
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
func New(ctx context.Context, cfg *config.Config, db *db.DB) (*http.Server, error) {
	e := echo.New()

	e.Use(errorMiddleware())

	// TODO: auth middleware

	e.GET("/", hello())

	// swagger docs
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// item routes
	e.GET("/item", getItem(ctx, db))
	e.GET("/item/:id", getItem(ctx, db))
	e.PUT("/item/:id", updateItem(ctx, db))
	e.DELETE("/item/:id", deleteItem(ctx, db))
	e.POST("/item", createItem(ctx, db))

	// TODO: category routes

	// TODO: user routes

	e.POST("/storage/presign-upload", presignUpload(ctx, cfg))

	s := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: e,
	}

	return s, nil
}

func hello() echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome! This is DooleyOnline.")
	})
}
