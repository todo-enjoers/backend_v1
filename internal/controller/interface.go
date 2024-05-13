package controller

import (
	"context"
)

type Controller interface {
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
