package api

import (
	"context"
	"net/http"

	authhandler "github.com/dooleyonline/backend/internal/api/auth"
	categoryhandler "github.com/dooleyonline/backend/internal/api/category"
	chathandler "github.com/dooleyonline/backend/internal/api/chat"
	itemhandler "github.com/dooleyonline/backend/internal/api/item"
	userhandler "github.com/dooleyonline/backend/internal/api/user"
	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/internal/db"
	"github.com/dooleyonline/backend/internal/service"
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
func New(ctx context.Context, cfg *config.Config, db *db.DB) (*echo.Echo, error) {
	services := service.New(cfg, db)

	item := itemhandler.New(services.Item)
	auth := authhandler.New(services.Auth)
	category := categoryhandler.New(services.Category)
	user := userhandler.New(services.User)

	chat := chathandler.New(services.Chat)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(loggerMiddleware())
	e.Use(errorMiddleware())
	e.Use(corsMiddleware())
	e.Use(authMiddleware(cfg, protectedRoutes))
	e.Use(contextMiddleware(cfg))

	e.GET("/", hello)

	// swagger docs
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// item routes
	e.GET("/item", item.GetMany)
	e.POST("/item", item.Create)
	e.GET("/item/:id", item.Get)
	e.PUT("/item/:id", item.Update)
	e.DELETE("/item/:id", item.Delete)
	e.POST("/item/:id/view", item.View)
	e.POST("/item/:id/sell", item.Sell)
	e.POST("/item/:id/like", item.Like)
	e.POST("/item/:id/unlike", item.Unlike)
	e.POST("/item/batch", item.GetBatch)

	// presign upload URL
	e.POST("/item/upload-url", item.GetUploadURL)

	// category
	e.GET("/category", category.GetAll)
	e.GET("/category/:name", category.Get)

	// user routes
	e.GET("/user", user.GetMany)
	e.POST("/user", user.Create)
	e.GET("/user/:id", user.Get)

	// auth routes
	e.POST("/auth/login", auth.Login)
	e.POST("/auth/logout", auth.Logout)
	e.GET("/auth/me", auth.GetMe)

	// chat routes
	e.POST("/chat/rooms", chat.CreateRoom)
	e.GET("/chat/rooms", chat.GetRooms)
	e.DELETE("/chat/rooms/:roomID", chat.DeleteRoom)
	e.GET("/chat/rooms/:roomID/messages", chat.GetLatest)
	e.POST("/chat/rooms/:roomID/messages", chat.CreateMessage)
	e.PATCH("/chat/messages/:messageID", chat.EditMessage)
	e.DELETE("/chat/messages/:messageID", chat.DeleteMessage)
	e.GET("/chat/rooms/:roomID/participants", chat.GetParticipants)
	e.POST("/chat/rooms/:roomID/participants", chat.AddParticipant)
	e.DELETE("/chat/rooms/:roomID/participants/:userID", chat.RemoveParticipant)
	e.GET("/chat/rooms/:roomID/ws", chat.HandleConnections)

	return e, nil
}

type routesConfig map[string][]string

// define routes to protect with auth middleware
var protectedRoutes = routesConfig{
	"/item":            {http.MethodPost},
	"/item/:id":        {http.MethodPut, http.MethodDelete},
	"/item/:id/sell":   {http.MethodPost},
	"/item/:id/like":   {http.MethodPost},
	"/item/:id/unlike": {http.MethodPost},
	"/auth/logout":     {http.MethodPost},
	"/chat":            {http.MethodPost},
	"/chat/:roomID/ws": {http.MethodGet},
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
