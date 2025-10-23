package shared

import (
	"github.com/labstack/echo/v4"
)

type Context struct {
	echo.Context
	UserId string
}
