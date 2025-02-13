package main

import (
	"SangXanh/cmd/api/controller"
	"SangXanh/pkg/config"
	"SangXanh/pkg/connection"
	"SangXanh/pkg/log"
	"SangXanh/pkg/repository"
	"SangXanh/pkg/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samber/do/v2"
)

func main() {
	di := do.New()
	config.Inject(di)
	connection.Inject(di)
	repository.Inject(di)
	service.Inject(di)

	serverConf := do.MustInvoke[config.Server](di)

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(log.Middleware())

	api := e.Group("/api")
	if err := controller.RegisterAPI(di, api); err != nil {
		panic(err)
	}

	log.Fatal(e.Start(serverConf.Address()))
}
