package main

import (
	"SangXanh/cmd/api/controller"
	middleware1 "SangXanh/cmd/api/middleware"
	_ "SangXanh/docs"
	"SangXanh/pkg/config"
	"SangXanh/pkg/connection"
	"SangXanh/pkg/log"
	"SangXanh/pkg/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samber/do/v2"
)

func main() {
	di := do.New()
	config.Inject(di)
	connection.Inject(di)
	service.Inject(di)

	serverConf := do.MustInvoke[config.Server](di)
	appConf := do.MustInvoke[config.App](di)

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(log.Middleware())

	jwtConf := do.MustInvoke[config.JWTKey](di)
	authMiddleware := middleware1.AuthenticationMiddleware(jwtConf.Key)

	// Enable Swagger only in development
	enableSwagger := appConf.IsDevelopment()

	api := e.Group("/api")
	if err := controller.RegisterAPI(di, api, authMiddleware, enableSwagger); err != nil {
		panic(err)
	}
	e.GET("/init-meilisearch", func(c echo.Context) error {
		connection.SyncData(di)
		return c.JSON(200, map[string]string{
			"status": "ok",
		})
	})

	log.Fatal(e.Start(serverConf.Address()))
}
