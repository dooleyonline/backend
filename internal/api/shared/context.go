package shared

import (
	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/internal/db"
	userdb "github.com/dooleyonline/backend/internal/db/user"
	"github.com/labstack/echo/v4"
)

type Context struct {
	echo.Context
	Cfg  *config.Config
	DB   *db.DB
	User *userdb.User
}
