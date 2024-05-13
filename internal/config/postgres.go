package config

import "fmt"

type PostgresConfig struct {
	Host     string `config:"postgres-host"`
	Port     int    `config:"postgres-port"`
	User     string `config:"postgres-user"`
	Password string `config:"postgres-password"`
	Database string `config:"postgres-database"`
}

func (cfg *PostgresConfig) DataBaseDNS() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	)
}
