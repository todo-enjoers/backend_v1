package pgx

import (
	"context"
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
	queryCreate = `INSERT INTO todos (id, name, description, created_by, project_id, "column" )VALUES ($1, $2, $3, $4, $5, $6)`
)

type todoStorage struct {
	pool *pgxpool.Pool
	log  *zap.Logger
}

func newTodoStorage(pool *pgxpool.Pool, log *zap.Logger) (*todoStorage, error) {
	store := &todoStorage{
		pool: pool,
		log:  log,
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
	_, err := store.pool.Exec(ctx, queryCreate, todo.CreatedBy, todo.Name, todo.ID, todo.Description)
	return err
}

func (store *todoStorage) Get(ctx context.Context, id string) (*model.TodoDTO, error) {
	return store, nil
}
