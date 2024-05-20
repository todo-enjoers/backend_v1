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
queryMigrateG = `CREATE TABLE IF NOT EXISTS projects
(
    "id" UUID NOT NULL UNIQUE,
    "name" VARCHAR NOT NULL,
    "created_by" UUID NOT NULL ,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("created_by") REFERENCES users(id) ON DELETE CASCADE
);`

queryMigrateGU = `CREATE TABLE IF NOT EXISTS users_in_projects
(
     "user_id" UUID NOT NULL,
    "project_id" UUID NOT NULL,
    FOREIGN KEY ("user_id") REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY ("project_id") REFERENCES projects(id) ON DELETE CASCADE
);`
queryCreate = INSERT INTO projects (id, name, created_by) VALUES ($1, $2, $3); // ???

queryGetByIDG = `SELECT p.id, p.name, p.created_by
FROM projects AS p
WHERE p.id = $1;`

queryGetMyGroups = `SELECT p.id, p.name, p.created_by
FROM projects AS p
WHERE created_by = $1;`
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
_, err = store.pool.Exec(context.Background(), queryMigrateGU)
if err != nil {
return err
}
return err
}

func (store *groupStorage) Create(ctx context.Context, group *model.GroupDTO) error {
_, err := store.pool.Exec(ctx, queryCreate, group.ID, group.Name, group.CreatedBy)
if err != nil {
if errors.As(err, &store.pgErr) && (pgerrcode.UniqueViolation == store.pgErr.Code) {
return storage.ErrAlreadyExists
}
return storage.ErrInserting
}
return nil
}

func (store *groupStorage) GetByID(ctx context.Context, id uuid.UUID) (*model.GroupDTO, error) {
g := new(model.GroupDTO)
err := store.pool.QueryRow(ctx, queryGetByIDG, id).Scan(&g.ID, &g.Name, &g.CreatedBy)
if err != nil {
return nil, storage.ErrGetByID
}
return g, nil
}

func (store *groupStorage) CreateInvite(ctx context.Context) error {
return nil
}

func (store *groupStorage) GetMyGroups(ctx context.Context, createdByID uuid.UUID) ([]model.GroupDTO, error) {
var res []model.GroupDTO
rows, err := store.pool.Query(ctx, queryGetMyGroups, createdByID)
if err != nil {
return nil, fmt.Errorf("error while querying my groups: %w", err)
}

defer rows.Close()

for rows.Next() {
var temp model.GroupDTO
err = rows.Scan(&temp.ID, &temp.Name, &temp.CreatedBy)
if err != nil {
return nil, fmt.Errorf("error while scanning groups: %w", err)
}
res = append(res, temp)
}

if err = rows.Err(); err != nil {
return nil, fmt.Errorf("unwrapped error: %w", err)
}

return res, nil
}