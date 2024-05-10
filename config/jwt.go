package config

import (
	"context"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/flags"
)

type JWT struct {
	AccessTokenLifeTime  int `config:"access_token_lifetime"`
	RefreshTokenLifeTime int `config:"refresh_token_lifetime"`
}

func NewJWT(ctx context.Context) (*Config, error) {
	cfg := new(Config)
	loader := confita.NewLoader(env.NewBackend(), flags.NewBackend())
	if err := loader.Load(ctx, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
