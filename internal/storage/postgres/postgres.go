package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/todo-enjoers/backend_v1/internal/model"
	"github.com/todo-enjoers/backend_v1/internal/storage"
)

type Client struct {
	pool *pgxpool.Pool
}

var (
	_ storage.Interface = (*Client)(nil)
)

const (
	QueryInsertIntoUsers    = `INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id;`
	QueryGetUserByID        = `SELECT login, password FROM users WHERE id = $1;`
	QueryUpdateUserPassword = `UPDATE users SET password = $1 WHERE id = $2;`
	QueryUserByLogin        = `SELECT login, password FROM users WHERE login = $1;`
	QueryMigrateUp          = `CREATE TABLE IF NOT EXISTS users
(
    id       bigserial primary key not null unique,
    login    varchar unique        not null,
    password varchar               not null
);`
)

// constructor
func New(pool *pgxpool.Pool) (*Client, error) {
	cli := &Client{pool}

	return cli, cli.Migrate(context.Background())
}

func (c *Client) InsertUser(ctx context.Context, item *model.UserDTO) error {
	_, err := c.pool.Exec(ctx, QueryInsertIntoUsers, item.Login, item.Password)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}
	return nil
}

func (c *Client) GetUserByID(ctx context.Context, id int64) (user *model.UserDTO, err error) {
	user = new(model.UserDTO)

	rows, err := c.pool.Query(ctx, QueryGetUserByID, id)
	if err != nil {
		return nil, fmt.Errorf("unable to query user: %w", err)
	}
	err = rows.Scan(&user.Login, &user.Password)
	return
}

func (c *Client) SearchUsersByLogin(ctx context.Context, login string) error {
	rows, err := c.pool.Query(ctx, QueryUserByLogin, login)
	if err != nil {
		return fmt.Errorf("unable to query users: %w", err)
	}
	err = rows.Scan(&login)
	return err
}

func (c *Client) UpdateUserPassword(ctx context.Context, password string, id int64) error {
	_, err := c.pool.Exec(ctx, QueryUpdateUserPassword, password, id)
	if err != nil {
		return fmt.Errorf("unable to update user or user not expected: %w", err)
	}
	return err
}

func (c *Client) Migrate(ctx context.Context) error {
	_, err := c.pool.Exec(ctx, QueryMigrateUp)
	return err
}
