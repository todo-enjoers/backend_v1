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
	_, err := store.pool.Exec(ctx, queryCreateTodo, todo.ID, todo.Name, todo.Description, todo.IsCompleted, todo.CreatedBy, todo.ProjectID, todo.Column)
	if err != nil {
		if errors.As(err, &store.pgErr) && pgerrcode.UniqueViolation == store.pgErr.Code {
			return storage.ErrAlreadyExists
		}
		return storage.ErrInserting
	}
	return nil
}

func (store *todoStorage) GetByID(ctx context.Context, id uuid.UUID) (*model.TodoDTO, error) {
	var todo model.TodoDTO
	err := store.pool.QueryRow(ctx, queryTodoGetByID, id).Scan(&todo.ID, &todo.Name, &todo.Description, &todo.IsCompleted, &todo.CreatedBy, &todo.ProjectID, &todo.Column)
	if err != nil {
		return nil, storage.ErrGetByID
	}
	return &todo, nil
}

func (store *todoStorage) GetAll(ctx context.Context, createdBy uuid.UUID) ([]model.TodoDTO, error) {
	var res []model.TodoDTO

	rows, err := store.pool.Query(ctx, queryGetAllTodos, createdBy)
	if err != nil {
		return nil, fmt.Errorf("error while querying all todos: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var temp model.TodoDTO
		err = rows.Scan(&temp.ID, &temp.Name, &temp.Description, &temp.IsCompleted, &temp.CreatedBy, &temp.ProjectID, &temp.Column)
		if err != nil {
			return nil, fmt.Errorf("error while scanning todos: %w", err)
		}
		res = append(res, temp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return res, nil
}
func (store *todoStorage) Update(ctx context.Context, todo *model.TodoDTO, id uuid.UUID) error {
	_, err := store.pool.Exec(ctx, queryUpdateTodo, todo.Name, todo.Description, todo.IsCompleted, id)
	return err
}

func (store *todoStorage) Delete(ctx context.Context, id uuid.UUID) error {
	commandTag, err := store.pool.Exec(ctx, queryDeleteTodo, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return storage.ErrNotFound
	}

	return nil
}
