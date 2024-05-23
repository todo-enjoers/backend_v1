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

const (
	queryMigrateG = `CREATE TABLE IF NOT EXISTS users_in_projects
(
   	"user_id" UUID NOT NULL,
    "project_id" UUID NOT NULL,
    FOREIGN KEY ("user_id") REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY ("project_id") REFERENCES projects(id) ON DELETE CASCADE
);`
	queryGetUsersInGroup = `SELECT uinp.user_id, uinp.project_id
FROM users_in_projects AS uinp
WHERE project_id = $1`

	queryInsertIntoGroup = `INSERT INTO users_in_projects (user_id, project_id) values ($1, $2)`

	queryDeleteFromGroup = `DELETE FROM users_in_projects
WHERE (user_id, project_id) = ($1, $2)`
)

// Checking whether the interface "GroupStorage" implements the structure "groupStorage"
var _ storage.GroupStorage = (*groupStorage)(nil)

type groupStorage struct {
	pool  *pgxpool.Pool
	log   *zap.Logger
	pgErr *pgconn.PgError
}

func newGroupStorage(pool *pgxpool.Pool, log *zap.Logger, pgErr *pgconn.PgError) (*groupStorage, error) {
	store := &groupStorage{
		pool:  pool,
		log:   log,
		pgErr: pgErr,
	}
	if err := store.migrate(); err != nil {
		return nil, err
	}
	return store, nil
}

func (store *groupStorage) migrate() (err error) {
	_, err = store.pool.Exec(context.Background(), queryMigrateG)
	if err != nil {
		return err
	}
	_, err = store.pool.Exec(context.Background(), queryMigrateG)
	if err != nil {
		return err
	}
	return err
}

func (store *groupStorage) CreateGroup(ctx context.Context, group *model.GroupDTO) error {
	_, err := store.pool.Exec(ctx, queryInsertIntoGroup, group.UserID, group.ProjectID)
	if err != nil {
		if errors.As(err, &store.pgErr) && (pgerrcode.UniqueViolation == store.pgErr.Code) {
			return storage.ErrAlreadyExists
		}
		return storage.ErrInserting
	}
	return nil
}

func (store *groupStorage) GetUsersInProjectByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.GroupDTO, error) {
	var usersList []model.GroupDTO
	rows, err := store.pool.Query(ctx, queryGetUsersInGroup, projectID)
	if err != nil {
		return nil, fmt.Errorf("error while querying my groups: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var temp model.GroupDTO
		err = rows.Scan(&temp.UserID, &temp.ProjectID)
		if err != nil {
			return nil, fmt.Errorf("error while scanning groups: %w", err)
		}
		usersList = append(usersList, temp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("unwrapped error: %w", err)
	}

	return usersList, nil
}

func (store *groupStorage) DeleteFromGroup(ctx context.Context, group *model.GroupDTO) error {
	_, err := store.pool.Exec(ctx, queryDeleteFromGroup, group.UserID, group.ProjectID)
	if err != nil {
		if errors.As(err, &store.pgErr) && (pgerrcode.UniqueViolation == store.pgErr.Code) {
			return storage.ErrAlreadyExists
		}
		return storage.ErrInserting
	}
	return nil
}
