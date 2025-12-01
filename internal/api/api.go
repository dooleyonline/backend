package api

import (
	"context"
	"net/http"

	authhandler "github.com/dooleyonline/backend/internal/api/auth"
	categoryhandler "github.com/dooleyonline/backend/internal/api/category"
	chathandler "github.com/dooleyonline/backend/internal/api/chat"
	itemhandler "github.com/dooleyonline/backend/internal/api/item"
	storagehandler "github.com/dooleyonline/backend/internal/api/storage"
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
	storage := storagehandler.New(services.Storage)
	chat := chathandler.New(services.Chat)

	// public API
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(loggerMiddleware())
	e.Use(errorMiddleware())
	e.Use(corsMiddleware())
	e.Use(contextMiddleware(cfg))

	e.GET("/", hello)

	// swagger docs
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// item routes
	e.GET("/item", item.GetMany)
	e.GET("/item/:id", item.Get)
	e.POST("/item/:id/view", item.View)
	e.POST("/item/batch", item.GetBatch)

	// category
	e.GET("/category", category.GetAll)
	e.GET("/category/:name", category.Get)

	// user.user routes
	e.GET("/user", user.GetMany)
	e.POST("/user", user.Create)
	e.GET("/user/:id", user.Get)
	e.PUT("/user", user.Update)
	e.GET("/user/interactions", user.GetLikedViewed)

	// auth routes
	e.GET("/auth/me", auth.GetMe)
	e.POST("/auth/login", auth.Login)
	e.POST("/auth/verify", auth.CreateVerification)
	e.GET("/auth/verify/:id", auth.GetVerification)
	e.POST("/auth/verify/:id", auth.VerifyUser)

	// storage routes
	e.POST("/storage/presign", storage.PresignUpload)

	// protected API routes
	protected := e.Group("")
	protected.Use(authMiddleware(cfg))

	protected.POST("/item", item.Create)
	protected.PUT("/item/:id", item.Update)
	protected.DELETE("/item/:id", item.Delete)
	protected.POST("/item/:id/sell", item.Sell)
	protected.POST("/item/:id/like", item.Like)
	protected.POST("/item/:id/unlike", item.Unlike)

	protected.PUT("/user", user.Update)
	protected.GET("/user/liked", user.GetLiked)
	protected.GET("/user/viewed", user.GetViewed)

	protected.POST("/auth/logout", auth.Logout)

	protected.POST("/chat/rooms", chat.CreateRoom)
	protected.GET("/chat/rooms", chat.GetRooms)
	protected.DELETE("/chat/:roomID", chat.DeleteRoom)
	protected.GET("/chat/:roomID/messages", chat.GetMessages)
	protected.PATCH("/chat/messages/:messageID", chat.EditMessage)
	protected.DELETE("/chat/messages/:messageID", chat.DeleteMessage)
	protected.GET("/chat/:roomID/participants", chat.GetParticipants)
	protected.POST("/chat/:roomID/participants/:userID", chat.AddParticipant)
	protected.DELETE("/chat/:roomID/participants/:userID", chat.RemoveParticipant)
	protected.GET("/chat/:roomID/ws", chat.HandleConnections)

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
