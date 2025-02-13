package log

import (
	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
)

func Middleware() echo.MiddlewareFunc {
	return echozap.ZapLogger(l)
}
