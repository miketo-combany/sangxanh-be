package config

import (
	"SangXanh/pkg/log"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/samber/do/v2"
	"os"
)

func init() {
	if _, err := os.Stat(".env"); err == nil {
		log.Info(".env file detected, load env variables")
		if err := godotenv.Load(".env"); err != nil {
			panic(err)
		}
	}
}

func Parse[T any](_ do.Injector) (T, error) {
	var conf T
	return conf, envconfig.Process("", &conf)
}
