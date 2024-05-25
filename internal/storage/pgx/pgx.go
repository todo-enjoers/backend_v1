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
	pool    *pgxpool.Pool
	log     *zap.Logger
	user    *userStorage
	project *projectsStorage
	todo    *todoStorage
	column  *columnStorage
	pgErr   *pgconn.PgError
}

func New(pool *pgxpool.Pool, log *zap.Logger, pgErr *pgconn.PgError) (*Storage, error) {
	users, err := newUserStorage(pool, log, pgErr)
	if err != nil {
		return nil, err
	}

	projects, err := newProjectsStorage(pool, log, pgErr)
	if err != nil {
		return nil, err
	}

	todos, err := newTodoStorage(pool, log, pgErr)
	if err != nil {
		return nil, err
	}

	columns, err := newColumnStorage(pool, log, pgErr)
	if err != nil {
		return nil, err
	}

	store := &Storage{
		pool:    pool,
		log:     log,
		user:    users,
		project: projects,
		todo:    todos,
		column:  columns,
	}

	return store, nil
}

func (s *Storage) User() storage.UserStorage {
	return s.user
}

func (s *Storage) Todo() storage.TodoStorage {
	return s.todo
}

func (s *Storage) Project() storage.ProjectStorage {
	return s.project
}

func (s *Storage) Column() storage.ColumnStorage {
	return s.column
}
