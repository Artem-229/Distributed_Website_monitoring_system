package app

import (
	"github.com/caarlos0/env"
)

type Config struct {
	PORT        string `env:"PORT"`
	DB_PORT     string `env:"DB_PORT"`
	DB_HOST     string `env:"DB_HOST"`
	DB_USERNAME string `env:"DB_USERNAME"`
	DB_PASSWORD string `env:"DB_PASSWORD"`
	DB_NAME     string `env:"DB_NAME"`
}

func MustGetFromEnv() Config {
	var config Config
	if err := env.Parse(&config); err != nil {
		panic("can not parse files from env")
	}

	return config
}
