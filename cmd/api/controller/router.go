package controller

import (
	"SangXanh/pkg/common/api"

	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func RegisterAPI(di do.Injector, e *echo.Group, auth echo.MiddlewareFunc, enableSwagger bool) error {
	type controllerWithMiddleware func(di do.Injector, auth echo.MiddlewareFunc) (api.Controller, error)

	controllers := []controllerWithMiddleware{
		NewUserController,
		NewProductController,
		NewCategoryController,
		NewProductVariantController,
		NewProductOptionController,
		NewImageController,
		NewAuthController,
		NewCartController,
		NewOrderController,
		NewVideoController,
	}

	for _, c := range controllers {
		ctrl, err := c(di, auth)
		if err != nil {
			return err
		}
		ctrl.Register(e)
	}

	// Register Swagger endpoint only in non-production environments
	if enableSwagger {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	return nil
}
