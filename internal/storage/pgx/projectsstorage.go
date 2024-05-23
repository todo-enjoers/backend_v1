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
	queryMigrateP = `CREATE TABLE IF NOT EXISTS projects
(
    "id" UUID NOT NULL UNIQUE,
    "name" VARCHAR NOT NULL,
    "created_by" UUID NOT NULL ,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("created_by") REFERENCES users(id) ON DELETE CASCADE
);`

	queryGetByIDP = `SELECT p.id, p.name, p.created_by
FROM projects AS p
WHERE p.id = $1;`

	queryGetMyProjects = `SELECT p.id, p.name, p.created_by 
FROM projects AS p
WHERE created_by = $1;`

	queryCreateProjects = `INSERT INTO projects (id, name, created_by) VALUES ($1, $2, $3);`

	queryUpdateName = `UPDATE projects SET name = $1 WHERE id = $2;`
)

// Checking whether the interface "GroupStorage" implements the structure "groupStorage"
var _ storage.ProjectStorage = (*projectsStorage)(nil)

type projectsStorage struct {
	pool  *pgxpool.Pool
	log   *zap.Logger
	pgErr *pgconn.PgError
}

func newProjectsStorage(pool *pgxpool.Pool, log *zap.Logger, pgErr *pgconn.PgError) (*projectsStorage, error) {
	store := &projectsStorage{
		pool:  pool,
		log:   log,
		pgErr: pgErr,
	}
	if err := store.migrate(); err != nil {
		return nil, err
	}
	return store, nil
}

func (store *projectsStorage) migrate() (err error) {
	_, err = store.pool.Exec(context.Background(), queryMigrateP)
	if err != nil {
		return err
	}
	return err
}

func (store *projectsStorage) CreateProjects(ctx context.Context, project *model.ProjectsDTO) error {
	_, err := store.pool.Exec(ctx, queryCreateProjects, project.ID, project.Name, project.CreatedBy)
	if err != nil {
		if errors.As(err, &store.pgErr) && (pgerrcode.UniqueViolation == store.pgErr.Code) {
			return storage.ErrAlreadyExists
		}
		return storage.ErrInserting
	}
	return nil
}

func (store *projectsStorage) GetByID(ctx context.Context, id uuid.UUID) (*model.ProjectsDTO, error) {
	g := new(model.ProjectsDTO)
	err := store.pool.QueryRow(ctx, queryGetByIDP, id).Scan(&g.ID, &g.Name, &g.CreatedBy)
	if err != nil {
		return nil, storage.ErrGetByID
	}
	return g, nil
}

func (store *projectsStorage) GetMyProjects(ctx context.Context, createdByID uuid.UUID) ([]model.ProjectsDTO, error) {
	var projectsList []model.ProjectsDTO
	rows, err := store.pool.Query(ctx, queryGetMyProjects, createdByID)
	if err != nil {
		return nil, fmt.Errorf("error while querying my projects: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var temp model.ProjectsDTO
		err = rows.Scan(&temp.ID, &temp.Name, &temp.CreatedBy)
		if err != nil {
			return nil, fmt.Errorf("error while scanning groups: %w", err)
		}
		projectsList = append(projectsList, temp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("unwrapped error: %w", err)
	}

	return projectsList, nil
}

func (store *projectsStorage) UpdateName(ctx context.Context, name string, id uuid.UUID) error {
	_, err := store.pool.Exec(ctx, queryUpdateName, name, id)
	return err
}
