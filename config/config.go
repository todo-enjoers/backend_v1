package config

import (
	"context"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/flags"
)

type Config struct {
	BindAddr string `config:"bind_addr, short = a"`
}

func New() (cfg *Config, err error) {
	cfg = &Config{
		BindAddr: ":8080",
	}
	loader := confita.NewLoader(
		env.NewBackend(),
		flags.NewBackend(),
	)
	err = loader.Load(context.Background(), cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
