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

// Checking whether the interface "ColumnStorage" implements the structure "columnStorage"
var _ storage.ColumnStorage = (*columnStorage)(nil)

type columnStorage struct {
	pool  *pgxpool.Pool
	log   *zap.Logger
	pgErr *pgconn.PgError
}

func newColumnStorage(pool *pgxpool.Pool, log *zap.Logger, pgErr *pgconn.PgError) (*columnStorage, error) {
	store := &columnStorage{
		pool:  pool,
		log:   log,
		pgErr: pgErr,
	}
	if err := store.migrate(); err != nil {
		return nil, err
	}
	return store, nil
}

func (store *columnStorage) migrate() (err error) {
	_, err = store.pool.Exec(context.Background(), queryMigrateColumnsTable)
	if err != nil {
		return storage.ErrTableMigrations
	}
	return nil
}

func (store *columnStorage) CreateColumn(ctx context.Context, column *model.ColumDTO) error {
	_, err := store.pool.Exec(ctx, queryInsertColumns, column.ProjectId, column.Name, column.Order)
	return err
}
func (store *columnStorage) DeleteColumn(ctx context.Context, name string, projectId uuid.UUID) error {
	_, err := store.pool.Exec(ctx, queryDeleteColumns, name, projectId)
	return err
}

func (store *columnStorage) GetColumnByName(ctx context.Context, name string, projectId uuid.UUID) (*model.ColumDTO, error) {
	var column model.ColumDTO
	err := store.pool.QueryRow(ctx, queryGetColumnByName, name, projectId).Scan(&column.ProjectId, &column.Name, &column.Order)
	if err != nil {
		return nil, err
	}
	return &column, nil
}
func (store *columnStorage) UpdateColumn(ctx context.Context, column *model.ColumDTO, name string, projectId uuid.UUID) error {
	_, err := store.pool.Exec(ctx, queryUpdateColumns, column.Name, name, projectId)
	return err
}
func (store *columnStorage) GetAllColumns(ctx context.Context, projectId uuid.UUID) ([]model.ColumDTO, error) {
	var res []model.ColumDTO

	rows, err := store.pool.Query(ctx, queryGetAllColumns, projectId)
	if err != nil {
		return nil, fmt.Errorf("error while querying all todos: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var temp model.ColumDTO
		err = rows.Scan(&temp.ProjectId, &temp.Name, &temp.Order)
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
