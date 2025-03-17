package config

import (
	"context"
	"os"
	"path"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	"go.uber.org/zap"
)

var configPath string

func init() {
	configPath = os.Getenv("CONFIG_PATH")
}

type Config struct {
	JWT        *JWT            `config:"JWT" toml:"JWT"`
	Controller *Controller     `config:"Controller" toml:"Controller"`
	Postgres   *PostgresConfig `config:"Postgres" toml:"Postgres"`
}

func New(log *zap.Logger) (*Config, error) {
	ctx := context.Background()

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	log.Info("loading config", zap.String("path", configPath))

	cfg := &Config{
		Controller: &Controller{
			Host: "localhost",
			Port: 8080,
		},
		Postgres: &PostgresConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			Database: "postgres",
		},
		JWT: &JWT{
			AccessTokenLifeTime:  20,
			RefreshTokenLifeTime: 10000,
			PublicKeyPath:        path.Join(wd, "certs", "public.pem"),
			PrivateKeyPath:       path.Join(wd, "certs", "private.pem"),
		},
	}

	loader := confita.NewLoader(
		env.NewBackend(),
		file.NewBackend(configPath),
	)
	if err = loader.Load(ctx, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
