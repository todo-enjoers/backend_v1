package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/todo-enjoers/backend_v1/internal/config"
	"github.com/todo-enjoers/backend_v1/internal/controller"
	"github.com/todo-enjoers/backend_v1/internal/controller/http"
	"github.com/todo-enjoers/backend_v1/internal/pkg/token/jwt"
	"github.com/todo-enjoers/backend_v1/internal/storage"
	"github.com/todo-enjoers/backend_v1/internal/storage/pgx"
	"github.com/todo-enjoers/backend_v1/pkg/postgres"
	"go.uber.org/zap"
	"os/signal"
	"syscall"
)

func main() {
	var (
		ctx      context.Context
		log      *zap.Logger
		err      error
		server   controller.Controller
		provider *jwt.Provider
		cfg      *config.Config
		pool     *pgxpool.Pool
		store    storage.Interface
		cancel   context.CancelFunc
	)

	ctx, cancel = signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	//init logger
	log, _ = zap.NewProduction() //no error because func recently don't  return a error

	//init config
	cfg, err = config.New(ctx)
	if err != nil {
		log.Fatal("Failed to initialize config", zap.Error(err))
	}
	log.Info("Initialized config", zap.Any("config", cfg))

	//init jwt provider
	provider, err = jwt.NewProvider(cfg, log)
	if err != nil {
		log.Fatal("Failed to initialize jwt provider", zap.Error(err))
	}

	//init pool
	pool, err = postgres.New(cfg, log)
	if err != nil {
		log.Fatal("Failed to initialize pool", zap.Error(err))
	}

	//init storage
	store, err = pgx.New(pool, log)
	if err != nil {
		log.Fatal("Failed to create pgx storage", zap.Error(err))
	}

	//init server
	server, err = http.New(store, log, cfg, provider)
	if err != nil {
		log.Fatal("Failed to initialize server", zap.Error(err))
	}

	//close server without fx
	defer func() {
		log.Error(
			"Shutting down server",
			zap.Error(server.Shutdown(ctx)),
		)
	}()
	err = server.Run(ctx)
	if err != nil {
		log.Fatal("Failed to run server", zap.Error(err))
	}

	<-ctx.Done()
	log.Info("[Graceful Shutdown]")
}
