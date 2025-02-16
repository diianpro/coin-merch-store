package config

import (
	"log/slog"
	"time"

	"diianpro/coin-merch-store/pkg/postgres"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTPPort string `env-required:"true" env:"HTTP_PORT"`

	Postgres postgres.Config
	JWT      JWT
	Hasher   Hasher
}

type JWT struct {
	SignKey  string        `env:"JWT_SIGN_KEY"`
	TokenTTL time.Duration `env:"JWT_TOKEN_TTL"`
}
type Hasher struct {
	Salt string `env:"HASHER_SALT"`
}

// New initialize Config structure
func New() (*Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		slog.Error("New config: %v", err)
	}
	return &cfg, nil
}
