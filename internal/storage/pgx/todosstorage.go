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

const (
	queryMigrate = `CREATE TABLE IF NOT EXISTS todos (
    "id" UUID PRIMARY KEY NOT NULL UNIQUE,
    "name" VARCHAR NOT NULL UNIQUE,
    "description" VARCHAR NOT NULL,
    "is_completed" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_by" UUID NOT NULL,
    "project_id" UUID NOT NULL,
    "column" VARCHAR NOT NULL,
    FOREIGN KEY (project_id, "column") REFERENCES project_columns(project_id, name),
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS todos_created_by_index ON todos(created_by);
`
	queryCreate      = `INSERT INTO todos (id, name, is_completed, description, created_by, project_id, "column" )VALUES ($1, $2, $3, $4, $5, $6, $7)`
	queryTodoGetByID = `SELECT created_by, name, id, description FROM todos WHERE id = $1`
	queryGetAllTodos = `SELECT name, id, description, is_completed 
		FROM todos 
		WHERE created_by = $1;`
	queryUpdate = `UPDATE todos
		SET name = $1, description = $2, is_completed = $3
		WHERE id = $4`
	queryDelete = `DELETE FROM todos WHERE id = $1`
)

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
	_, err := store.pool.Exec(context.Background(), queryMigrate)
	return err
}

func (store *todoStorage) Create(ctx context.Context, todo *model.TodoDTO) error {
	_, err := store.pool.Exec(ctx, queryCreate, todo.ID, todo.Name, todo.Description, todo.IsCompleted, todo.ProjectID, todo.CreatedBy, todo.Column)
	return err
}

func (store *todoStorage) GetByID(ctx context.Context, id uuid.UUID) (*model.TodoDTO, error) {
	var todo model.TodoDTO
	err := store.pool.QueryRow(ctx, queryTodoGetByID, id).Scan(&todo.ID, &todo.Name, &todo.Description, &todo.IsCompleted, &todo.CreatedBy, &todo.ProjectID, &todo.Column)
	if err != nil {
		return nil, err
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
		err = rows.Scan(&temp.Name, &temp.ID, &temp.Description, &temp.IsCompleted)
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
	_, err := store.pool.Exec(ctx, queryUpdate, todo.Name, todo.Description, todo.IsCompleted, id)
	return err
}

func (store *todoStorage) DeleteTodos(ctx context.Context, id uuid.UUID) error {

	commandTag, err := store.pool.Exec(ctx, queryDelete, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return storage.ErrNotFound
	}

	return nil
}
