package api

import "github.com/labstack/echo/v4"

type Controller interface {
	Register(e *echo.Group)
}
