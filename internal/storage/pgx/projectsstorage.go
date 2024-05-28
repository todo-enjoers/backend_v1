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

// Checking whether the interface "ProjectStorage" implements the structure "projectsStorage"
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
		return storage.ErrTableMigrations
	}
	return nil
}

func (store *projectsStorage) Create(ctx context.Context, project *model.ProjectDTO) error {
	_, err := store.pool.Exec(ctx, queryCreateProjects, project.ID, project.Name, project.CreatedBy)
	if err != nil {
		if errors.As(err, &store.pgErr) && (pgerrcode.UniqueViolation == store.pgErr.Code) {
			return storage.ErrAlreadyExists
		}
		return storage.ErrInserting
	}
	return nil
}

func (store *projectsStorage) GetMyByName(ctx context.Context, name string, createdBy uuid.UUID) error {
	project := new(model.ProjectDTO)
	err := store.pool.QueryRow(ctx, queryGetMyProjectsByName, name, createdBy).Scan(&project.ID, &project.Name, &project.CreatedBy)
	if err != nil {
		return err
	}
	return nil
}

func (store *projectsStorage) GetByID(ctx context.Context, id uuid.UUID) (*model.ProjectDTO, error) {
	project := new(model.ProjectDTO)
	err := store.pool.QueryRow(ctx, queryGetProjectsByID, id).Scan(&project.ID, &project.Name, &project.CreatedBy)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (store *projectsStorage) GetMyProjects(ctx context.Context, createdByID uuid.UUID) ([]model.ProjectDTO, error) {
	var projectsList []model.ProjectDTO
	rows, err := store.pool.Query(ctx, queryGetMyProjects, createdByID)
	if err != nil {
		return nil, fmt.Errorf("error while querying my projects: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var temp model.ProjectDTO
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
	commandTag, err := store.pool.Exec(ctx, queryUpdateProjectName, name, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return storage.ErrNotFound
	}
	return nil
}

func (store *projectsStorage) Delete(ctx context.Context, id uuid.UUID) error {
	commandTag, err := store.pool.Exec(ctx, queryDeleteProject, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return storage.ErrNotFound
	}
	return err
}
