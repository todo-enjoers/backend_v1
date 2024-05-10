package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/todo-enjoers/backend_v1/config"
	"github.com/todo-enjoers/backend_v1/internal/controller/http"
	controller "github.com/todo-enjoers/backend_v1/internal/controller/http"
	"github.com/todo-enjoers/backend_v1/internal/storage"
	"github.com/todo-enjoers/backend_v1/internal/storage/postgres"
	"github.com/todo-enjoers/backend_v1/pkg/client"
	"go.uber.org/zap"
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
	)

	//init logger
	log, err = zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to initialize logger", zap.Error(err))
	}

	//init config
	cfg, err = config.New()
	if err != nil {
		log.Fatal("error initializing config", zap.Error(err))
	}

	//init pool
	pool, err = client.New(context.Background(), cfg)
	if err != nil {
		log.Fatal("Failed to initialize logger", zap.Error(err))
	}

	//init storage
	store, err = postgres.New(pool.P())
	if err != nil {
		log.Fatal("Failed to create pgx storage", zap.Error(err))
	}

	//init server
	server, err = http.New(log, cfg, provider, store)
	defer func() {
		log.Error(
			"Shutting down server",
			zap.Error(server.Shutdown(ctx)),
		)
	}()
	err = server.Run(ctx)
	if err != nil {
		log.Fatal("Failed to initialize server", zap.Error(err))
	}
}
