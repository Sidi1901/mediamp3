package config

import(
	"github.com/caarlos0/env/v11"
)

type Config struct {
	ServerPort string `env:"SERVER_PORT" envDefault:"8080"`
}


func LoadEnv() (Config,error) {

	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}