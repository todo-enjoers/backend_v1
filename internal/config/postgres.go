package config

import "fmt"

type PostgresConfig struct {
	Host     string `config:"host" toml:"host"`
	Port     int    `config:"port" toml:"port"`
	User     string `config:"user" toml:"user"`
	Password string `config:"password" toml:"password"`
	Database string `config:"database" toml:"database"`
}

func (cfg *PostgresConfig) GetURI() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	)
}
