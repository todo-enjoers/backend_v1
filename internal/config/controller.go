package config

import "fmt"

type Controller struct {
	BindAddress string `config:"bind_address,short=a"`
	BindPort    int    `config:"bind_port,short=p"`
}

func (c Controller) GetBindAddress() string {
	return fmt.Sprintf("%s:%d", c.BindAddress, c.BindPort)
}
