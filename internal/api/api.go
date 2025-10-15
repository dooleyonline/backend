package api

import (
	"context"
	"log/slog"
	"net/http"

	authapi "github.com/dooleyonline/backend/internal/api/auth"
	categoryapi "github.com/dooleyonline/backend/internal/api/category"
	itemapi "github.com/dooleyonline/backend/internal/api/item"
	storageapi "github.com/dooleyonline/backend/internal/api/storage"
	userapi "github.com/dooleyonline/backend/internal/api/user"
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
	e.Use(corsMiddleware())
	e.Use(authMiddleware(cfg, protectedRoutes))
	e.Use(contextMiddleware(cfg, db))

	e.GET("/", hello)

	// swagger docs
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// item routes
	e.GET("/item", itemapi.GetMany)
	e.POST("/item", itemapi.Create)
	e.GET("/item/:id", itemapi.Get)
	e.PUT("/item/:id", itemapi.Update)
	e.DELETE("/item/:id", itemapi.Delete)
	e.POST("/item/:id/view", itemapi.IncrementView)
	e.POST("/item/:id/sell", itemapi.Sell)

	// category
	e.GET("/category", categoryapi.GetAll)
	e.GET("/category/:name", categoryapi.Get)

	// user routes
	e.GET("/user", userapi.GetMany)
	e.POST("/user", userapi.Create)

	// auth routes
	e.POST("/auth/login", authapi.Login)
	e.POST("/auth/logout", authapi.Logout)

	// storage routes
	e.POST("/storage/presign", storageapi.Presign)

	return e, nil
}

type routesConfig map[string][]string

// define routes to protect with auth middleware
var protectedRoutes = routesConfig{
	"/item":          {http.MethodPost},
	"/item/:id":      {http.MethodPut, http.MethodDelete},
	"/item/:id/sell": {http.MethodPost},
	"/auth/logout":   {http.MethodPost},
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
