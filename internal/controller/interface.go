package controller

import "context"

type Controller interface {
	configureRoutes()
	Run() error
	Shutdown(ctx context.Context) error
}
