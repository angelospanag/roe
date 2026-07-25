package config

import "github.com/caarlos0/env/v11"

type Config struct {
	Port        int    `env:"PORT"         envDefault:"8080"`
	DatabaseURL string `env:"DATABASE_URL" envDefault:"postgres://postgres:postgres@localhost:5432/roe_backend?sslmode=disable"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
