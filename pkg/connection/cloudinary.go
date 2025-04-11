package connection

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/samber/do/v2"
)

func RegisterCloudinary(di do.Injector) (*cloudinary.Cloudinary, error) {
	return cloudinary.New()
}
