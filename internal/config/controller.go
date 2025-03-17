package config

import "fmt"

type Controller struct {
	Host string `config:"host" toml:"host"`
	Port int    `config:"port" toml:"port"`
}

func (c Controller) GetBindAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
