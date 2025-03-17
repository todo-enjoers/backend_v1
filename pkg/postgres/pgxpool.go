package postgres

import (
	"context"
	"fmt"

	pgxlib "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/todo-enjoers/backend_v1/internal/config"
)

const (
	enumTypeRole = "role"
)

// New opens new postgres connection, configures it and return prepared client.
func New(cfg *config.Config, log *zap.Logger) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	log.Info("initializing postgres client")

	c, err := pgxpool.ParseConfig(cfg.Postgres.GetURI())
	if err != nil {
		return nil, fmt.Errorf("error while parsing db uri: %w", err)
	}

	// Add UUID support
	c.AfterConnect = func(ctx context.Context, conn *pgxlib.Conn) error {
		var dt *pgtype.Type
		dt, err = conn.LoadType(ctx, enumTypeRole)
		if err != nil {
			return err
		}
		conn.TypeMap().RegisterType(dt)

		return nil
	}

	pool, err = pgxpool.NewWithConfig(context.Background(), c)
	if err != nil {
		return nil, fmt.Errorf("postgres: init pgxpool: %w", err)
	}

	log.Info("created postgres client")
	return pool, nil
}
