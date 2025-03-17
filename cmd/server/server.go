package main

import (
	"context"
	"errors"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/todo-enjoers/backend_v1/internal/config"
	"github.com/todo-enjoers/backend_v1/internal/controller"
	"github.com/todo-enjoers/backend_v1/internal/controller/http"
	"github.com/todo-enjoers/backend_v1/internal/pkg/tern/migrator"
	"github.com/todo-enjoers/backend_v1/internal/pkg/token"
	"github.com/todo-enjoers/backend_v1/internal/pkg/token/jwt"
	"github.com/todo-enjoers/backend_v1/internal/storage"
	"github.com/todo-enjoers/backend_v1/internal/storage/pgx"
	"github.com/todo-enjoers/backend_v1/migrations"
	"github.com/todo-enjoers/backend_v1/pkg/postgres"
)

var (
	ErrNilReference = errors.New("nil reference")
)

func main() {
	fx.New(CreateApp()).Run()
}

func createLogger(log *zap.Logger) fxevent.Logger {
	return &fxevent.ZapLogger{
		Logger: log.Named("fx"),
	}
}

func newLogger() *zap.Logger {
	l, _ := zap.NewProduction()
	return l.Named("todoer")
}

func migrate(log *zap.Logger, cfg *config.Config) error {
	if log == nil {
		log = zap.L().Named("migrator")
	}
	if cfg == nil {
		return ErrNilReference
	}
	ctx := context.Background()
	m, err := migrator.New(ctx, cfg, log, migrations.Migrations)
	if err != nil {
		return err
	}
	defer m.Close(ctx)
	return m.MigrateUp(ctx)
}

func CreateApp() fx.Option {
	return fx.Options(
		fx.WithLogger(createLogger),
		fx.Provide(
			newLogger,
			config.New,
			postgres.New,

			fx.Annotate(http.New, fx.As(new(controller.Controller))),
			fx.Annotate(pgx.New, fx.As(new(storage.Interface))),
			fx.Annotate(jwt.NewProvider, fx.As(new(token.Provider))),
		),
		fx.Invoke(
			migrate,
			controller.RunControllerFx,
		),
	)
}
