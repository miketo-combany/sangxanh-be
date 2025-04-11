package config

type Cloudinary struct {
	URL string `envconfig:"CLOUDINARY_URL" required:"true"`
}
