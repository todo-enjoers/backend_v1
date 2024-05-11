package config

import (
	"context"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/flags"
	"os"
	"path"
)

//type Config struct {
//	BindAddress string `config:"bind_address,short=a"`
//	BindPort    string `config:"bind_port,short=p"`
//	DataBaseDNS string `config:"data_base_dns,short=d"`
//	JWT         string `config:"jwt"`
//}

type Config struct {
	JWT        *JWT
	Controller *Controller
	Postgres   *PostgresConfig
}

func New(ctx context.Context) (*Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Controller: &Controller{
			BindAddress: "localhost",
			BindPort:    8080,
		},
		Postgres: &PostgresConfig{
			DataBaseDNS: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
		},
		JWT: &JWT{
			AccessTokenLifeTime:  20,
			RefreshTokenLifeTime: 10000,
			PublicKeyPath:        path.Join(wd, "certs", "public.pem"),
			PrivateKeyPath:       path.Join(wd, "certs", "private.pem"),
		},
	}
	loader := confita.NewLoader(env.NewBackend(), flags.NewBackend())
	if err = loader.Load(context.Background(), cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
