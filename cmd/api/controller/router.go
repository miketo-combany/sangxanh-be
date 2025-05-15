package controller

import (
	"SangXanh/pkg/common/api"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

func RegisterAPI(di do.Injector, e *echo.Group, auth echo.MiddlewareFunc) error {
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
	}

	for _, c := range controllers {
		ctrl, err := c(di, auth)
		if err != nil {
			return err
		}
		ctrl.Register(e)
	}
	return nil
}
