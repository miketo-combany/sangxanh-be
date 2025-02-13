package config

type Mongo struct {
	URI      string `envconfig:"MONGO_URI" required:"true"`
	Database string `envconfig:"MONGO_DATABASE" required:"true"`
}
