package config

type JWTKey struct {
	Key string `envconfig:"JWT_KEY" required:"true"`
}
