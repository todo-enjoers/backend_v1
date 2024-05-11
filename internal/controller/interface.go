package controller

import (
	"context"
	"go.uber.org/zap"
)

type Controller interface {
	Run(_ context.Context, log *zap.Logger) error
	Shutdown(ctx context.Context) error
}
