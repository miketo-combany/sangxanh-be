package config

import "github.com/samber/do/v2"

func Inject(di do.Injector) {
	do.Provide(di, Parse[Mongo])
	do.Provide(di, Parse[Server])
}
