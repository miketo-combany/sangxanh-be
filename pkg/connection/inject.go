package connection

import "github.com/samber/do/v2"

func Inject(di do.Injector) {
	do.Provide(di, NewSupabaseDatabase)
	do.Provide(di, RegisterCloudinary)
}
