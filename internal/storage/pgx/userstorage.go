package pgx

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/todo-enjoers/backend_v1/internal/model"
	"github.com/todo-enjoers/backend_v1/internal/storage"
	"go.uber.org/zap"
)

// Checking whether the interface "TodoStorage" implements the structure "todoStorage"
var _ storage.UserStorage = (*userStorage)(nil)
var pgErr *pgconn.PgError

// Users Query const
const (
	queryInsertInto = `INSERT INTO users (id, login, encrypted_password) VALUES ($1, $2, $3);`

	queryGetByID = `SELECT u.id, u.login, u.encrypted_password
FROM users AS u
WHERE u.id = $1;`

	queryUpdatePassword = `UPDATE users SET encrypted_password = $1 WHERE id = $2;`

	queryGetByLogin = `SELECT u.id, u.login, u.encrypted_password
FROM users AS u
WHERE u.login = $1;`

	queryMigrateUp = `CREATE TABLE IF NOT EXISTS users
(
    id UUID PRIMARY KEY NOT NULL UNIQUE ,
    login VARCHAR NOT NULL UNIQUE ,
    encrypted_password VARCHAR NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS users_login_idx ON users (login);`

	queryGetAll = `SELECT u.id, u.login
FROM users AS u;`
)

type userStorage struct {
	pool *pgxpool.Pool
	log  *zap.Logger
}

func newUserStorage(pool *pgxpool.Pool, log *zap.Logger) (*userStorage, error) {
	store := &userStorage{
		pool: pool,
		log:  log,
	}
	if err := store.migrate(); err != nil {
		return nil, err
	}
	return store, nil
}

func (store *userStorage) migrate() error {
	_, err := store.pool.Exec(context.Background(), queryMigrateUp)
	return err
}

func (store *userStorage) Create(ctx context.Context, user *model.UserDTO) error {
	_, err := store.pool.Exec(ctx, queryInsertInto, user.ID, user.Login, user.Password)
	if err != nil {
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			return storage.ErrAlreadyClosed
		}
		return fmt.Errorf("err while creating user: %w", err) //Need to fix error here
	}
	return nil
}

func (store *userStorage) GetByID(ctx context.Context, id uuid.UUID) (*model.UserDTO, error) {
	u := new(model.UserDTO)
	err := store.pool.QueryRow(ctx, queryGetByID, id).Scan(&u.ID, &u.Login, &u.Password)
	if err != nil {
		return nil, storage.ErrGetByID
	}
	return u, nil
}

func (store *userStorage) GetByLogin(ctx context.Context, login string) (*model.UserDTO, error) {
	u := new(model.UserDTO)
	err := store.pool.QueryRow(ctx, queryGetByLogin, login).Scan(&u.ID, &u.Login, &u.Password)
	if err != nil {
		return nil, storage.ErrGetByLogin
	}
	return u, nil
}

func (store *userStorage) ChangePassword(ctx context.Context, password string, id uuid.UUID) error {
	_, err := store.pool.Exec(ctx, queryUpdatePassword, password, id)
	return err
}

func (store *userStorage) GetAll(ctx context.Context) ([]model.UserDTO, error) {
	var res []model.UserDTO
	rows, err := store.pool.Query(ctx, queryGetAll)
	if err != nil {
		return nil, fmt.Errorf("error while querying all users: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var temp model.UserDTO
		err = rows.Scan(&temp.ID, &temp.Login)
		if err != nil {
			return nil, fmt.Errorf("error while scanning users: %w", err)
		}
		res = append(res, temp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("unwrapped error: %w", err)
	}

	return res, nil
}
