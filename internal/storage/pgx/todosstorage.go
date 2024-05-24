package pgx

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/todo-enjoers/backend_v1/internal/model"
	"github.com/todo-enjoers/backend_v1/internal/storage"
	"go.uber.org/zap"
)

var _ storage.TodoStorage = (*todoStorage)(nil)

type todoStorage struct {
	pool  *pgxpool.Pool
	log   *zap.Logger
	pgErr *pgconn.PgError
}

func newTodoStorage(pool *pgxpool.Pool, log *zap.Logger, pgErr *pgconn.PgError) (*todoStorage, error) {
	store := &todoStorage{
		pool:  pool,
		log:   log,
		pgErr: pgErr,
	}
	if err := store.migrateT(); err != nil {
		return nil, err
	}
	return store, nil
}

func (store *todoStorage) migrateT() error {
	_, err := store.pool.Exec(context.Background(), queryMigrateT)
	return err
}

func (store *todoStorage) Create(ctx context.Context, todo *model.TodoDTO) error {
	_, err := store.pool.Exec(ctx, queryCreateTodo, todo.CreatedBy, todo.Name, todo.ID, todo.Description)
	return err
}

func (store *todoStorage) GetByID(ctx context.Context, id uuid.UUID) (*model.TodoDTO, error) {
	var todo model.TodoDTO
	err := store.pool.QueryRow(ctx, queryTodoGetByID, id).Scan(&todo.CreatedBy, &todo.Name, &todo.ID, &todo.Description)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (store *todoStorage) GetAll(ctx context.Context) ([]model.TodoDTO, error) {
	var res []model.TodoDTO
	rows, err := store.pool.Query(ctx, queryGetAllTodos)
	if err != nil {
		return nil, fmt.Errorf("error while querying all todos: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var temp model.TodoDTO
		err = rows.Scan(&temp.IsCompleted, &temp.Name, &temp.ID, &temp.Description)
		if err != nil {
			return nil, fmt.Errorf("error while scanning todos: %w", err)
		}
		res = append(res, temp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("unwrapped error: %w", err)
	}

	return res, err
}
func (store *todoStorage) Update(ctx context.Context, todo *model.TodoDTO) error {
	_, err := store.pool.Exec(ctx, queryUpdateTodo, todo.Name, todo.Description, todo.IsCompleted, todo.ID, todo.CreatedBy)
	return err
}
func (store *todoStorage) Delete(ctx context.Context, id uuid.UUID, columnName string, createdBy uuid.UUID) error {

	commandTag, err := store.pool.Exec(ctx, queryDeleteTodo, id, columnName, createdBy)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return storage.ErrNotFound
	}

	return nil
}
