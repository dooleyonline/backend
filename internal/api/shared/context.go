package shared

import (
	"github.com/dooleyonline/backend/internal/config"
	"github.com/dooleyonline/backend/internal/db"
	"github.com/labstack/echo/v4"
)

type Context struct {
	echo.Context
	Cfg *config.Config
	DB  *db.DB
}
