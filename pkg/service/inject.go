package service

import "github.com/samber/do/v2"

func Inject(di do.Injector) {
	do.Provide(di, NewUserService)
	do.Provide(di, NewCategoryService)
	do.Provide(di, NewProductService)
	do.Provide(di, NewProductVariantService)
	do.Provide(di, NewProductOptionService)
	do.Provide(di, NewImageService)
	do.Provide(di, NewAuthService)
	do.Provide(di, NewCartService)
	do.Provide(di, NewOrderService)
}
