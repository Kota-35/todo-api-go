package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DatabaseUrl string `env:"DATABASE_URL"`
	JwtSecret   string `env:"JWT_SECRET"`
	Env         string `env:"ENV"`
}

func LoadEnv() Config {
	cfg := Config{}
	env.Parse(&cfg)
	fmt.Println(cfg)

	return cfg
}
