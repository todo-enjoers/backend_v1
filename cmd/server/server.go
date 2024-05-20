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
	log, _ = zap.NewProduction() //no error because func recently don't return an error

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

<div class="message spoilers-container" dir="auto"><span class="translatable-message">package main

import (
"context"
"<a class="anchor-url" href="https://github.com/jackc/pgx/v5/pgconn" target="_blank" rel="noopener noreferrer">github.com/jackc/pgx/v5/pgconn</a>"
"<a class="anchor-url" href="https://github.com/jackc/pgx/v5/pgxpool" target="_blank" rel="noopener noreferrer">github.com/jackc/pgx/v5/pgxpool</a>"
"<a class="anchor-url" href="https://github.com/todo-enjoers/backend_v1/internal/config" target="_blank" rel="noopener noreferrer">github.com/todo-enjoers/backend_v1/internal/config</a>"
"<a class="anchor-url" href="https://github.com/todo-enjoers/backend_v1/internal/controller" target="_blank" rel="noopener noreferrer">github.com/todo-enjoers/backend_v1/internal/controller</a>"
"<a class="anchor-url" href="https://github.com/todo-enjoers/backend_v1/internal/controller/http" target="_blank" rel="noopener noreferrer">github.com/todo-enjoers/backend_v1/internal/controller/http</a>"
"<a class="anchor-url" href="https://github.com/todo-enjoers/backend_v1/internal/pkg/token/jwt" target="_blank" rel="noopener noreferrer">github.com/todo-enjoers/backend_v1/internal/pkg/token/jwt</a>"
"<a class="anchor-url" href="https://github.com/todo-enjoers/backend_v1/internal/storage" target="_blank" rel="noopener noreferrer">github.com/todo-enjoers/backend_v1/internal/storage</a>"
"<a class="anchor-url" href="https://github.com/todo-enjoers/backend_v1/internal/storage/pgx" target="_blank" rel="noopener noreferrer">github.com/todo-enjoers/backend_v1/internal/storage/pgx</a>"
"<a class="anchor-url" href="https://github.com/todo-enjoers/backend_v1/pkg/postgres" target="_blank" rel="noopener noreferrer">github.com/todo-enjoers/backend_v1/pkg/postgres</a>"
"<a class="anchor-url" href="https://go.uber.org/zap" target="_blank" rel="noopener noreferrer">go.uber.org/zap</a>"
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
		pgErr    *pgconn.PgError
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
	log, _ = zap.NewProduction() //no error because func recently don't return an error

	//init config
	cfg, err = <a class="anchor-url" href="https://config.New" target="_blank" rel="noopener noreferrer">config.New</a>(ctx)
	if err != nil {
		log.Fatal("Failed to initialize config", zap.Error(err))
	}
	<a class="anchor-url" href="https://log.Info" target="_blank" rel="noopener noreferrer">log.Info</a>("Initialized config", zap.Any("config", cfg))

	//init jwt provider
	provider, err = jwt.NewProvider(cfg, log)
	if err != nil {
		log.Fatal("Failed to initialize jwt provider", zap.Error(err))
	}

	//init pool
	pool, err = <a class="anchor-url" href="https://postgres.New" target="_blank" rel="noopener noreferrer">postgres.New</a>(cfg, log)
	if err != nil {
		log.Fatal("Failed to initialize pool", zap.Error(err))
	}

	//init storage
	store, err = <a class="anchor-url" href="https://pgx.New" target="_blank" rel="noopener noreferrer">pgx.New</a>(pool, log, pgErr)
	if err != nil {
		log.Fatal("Failed to create pgx storage", zap.Error(err))
	}

	//init server
	server, err = <a class="anchor-url" href="https://http.New" target="_blank" rel="noopener noreferrer">http.New</a>(store, log, cfg, provider)
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

	&lt;-ctx.Done()
	<a class="anchor-url" href="https://log.Info" target="_blank" rel="noopener noreferrer">log.Info</a>("[Graceful Shutdown]")
}</span><span class="time"><span class="i18n" dir="auto">12:44</span><div class="time-inner" title="20 May 2024, 12:44:53"><span class="i18n" dir="auto">12:44</span></div></span></div>