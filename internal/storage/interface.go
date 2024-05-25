package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/todo-enjoers/backend_v1/internal/model"
)

type UserStorage interface {
	Create(ctx context.Context, user *model.UserDTO) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.UserDTO, error)
	GetByLogin(ctx context.Context, login string) (*model.UserDTO, error)
	ChangePassword(ctx context.Context, password string, id uuid.UUID) error
	GetAll(ctx context.Context) ([]model.UserDTO, error)
}

type TodoStorage interface {
	Create(ctx context.Context, todo *model.TodoDTO) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.TodoDTO, error)
	GetAll(ctx context.Context, createdBy uuid.UUID) ([]model.TodoDTO, error)
	Update(ctx context.Context, todo *model.TodoDTO, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProjectStorage interface {
	GetMyByName(ctx context.Context, name string, createdBy uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.ProjectDTO, error)
	GetMyProjects(ctx context.Context, createdByID uuid.UUID) ([]model.ProjectDTO, error)
	UpdateName(ctx context.Context, name string, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	Create(ctx context.Context, project *model.ProjectDTO) error
}
type ColumnStorage interface {
	CreateColumn(ctx context.Context, column *model.ColumDTO) error
	DeleteColumn(ctx context.Context, name string, projectId uuid.UUID) error
	GetColumnByName(ctx context.Context, name string, projectId uuid.UUID) (*model.ColumDTO, error)
	UpdateColumn(ctx context.Context, column *model.ColumDTO, name string, projectId uuid.UUID) error
	GetAllColumns(ctx context.Context, projectId uuid.UUID) ([]model.ColumDTO, error)
}

type Interface interface {
	User() UserStorage
	Todo() TodoStorage
	Project() ProjectStorage
	Column() ColumnStorage
}
