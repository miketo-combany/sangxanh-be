package config

import "github.com/samber/do/v2"

func Inject(di do.Injector) {
	do.Provide(di, Parse[Supabase])
	do.Provide(di, Parse[Server])
	do.Provide(di, Parse[JWTKey])
	do.Provide(di, Parse[Cloudinary])
}
