package client

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/todo-enjoers/backend_v1/config"
)

type Client struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg *config.Config) (*Client, error) {
	pool, err := pgxpool.New(ctx, cfg.DataBaseDNS)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %v\n", err)
	}
	pgInstance := &Client{pool}

	return pgInstance, nil
}

func (c *Client) Ping(ctx context.Context) error {
	return c.pool.Ping(ctx)
}

func (c *Client) Close() {
	c.pool.Close()
}

func (c *Client) P() *pgxpool.Pool {
	return c.pool
}
