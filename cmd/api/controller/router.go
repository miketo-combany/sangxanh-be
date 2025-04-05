package controller

import (
	"SangXanh/pkg/common/api"
	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"
)

func RegisterAPI(di do.Injector, e *echo.Group) error {
	type controller func(di do.Injector) (api.Controller, error)
	controllers := []controller{
		NewUserController,
		NewProductController,
		NewCategoryController,
		NewProductVariantController,
		NewProductOptionController,
	}

	for _, c := range controllers {
		ctrl, err := c(di)
		if err != nil {
			return err
		}
		ctrl.Register(e)
	}
	return nil
}
