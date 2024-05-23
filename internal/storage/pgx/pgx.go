package pgx

import (
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/todo-enjoers/backend_v1/internal/storage"
	"go.uber.org/zap"
)

// Checking whether the interface "Controller" implements the structure "Controller"
var _ storage.Interface = (*Storage)(nil)

type Storage struct {
	pool  *pgxpool.Pool
	log   *zap.Logger
	user  *userStorage
	group *groupStorage
	pgErr *pgconn.PgError
}

func New(pool *pgxpool.Pool, log *zap.Logger, pgErr *pgconn.PgError) (*Storage, error) {
	users, err := newUserStorage(pool, log, pgErr)
	if err != nil {
		return nil, err
	}

	groups, err := newGroupStorage(pool, log, pgErr)
	if err != nil {
		return nil, err
	}

	store := &Storage{
		pool:  pool,
		log:   log,
		user:  users,
		group: groups,
	}

	return store, nil
}

func (s *Storage) User() storage.UserStorage {
	return s.user
}

func (s *Storage) Group() storage.GroupStorage {
	return s.group
}
