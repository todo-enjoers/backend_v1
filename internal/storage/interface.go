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
	GetAll(ctx context.Context) ([]model.TodoDTO, error)
	Update(ctx context.Context, todo *model.TodoDTO) error
	Delete(ctx context.Context, id uuid.UUID, columnName string, createdBy uuid.UUID) error
}

type ProjectStorage interface {
	CreateProject(ctx context.Context, project *model.ProjectsDTO) error
	GetMyProjects(ctx context.Context, createdByID uuid.UUID) ([]model.ProjectsDTO, error)
	UpdateProjectName(ctx context.Context, name string, id uuid.UUID) error
	DeleteProject(ctx context.Context, id uuid.UUID) error
}

type Interface interface {
	User() UserStorage
	Todo() TodoStorage
	Project() ProjectStorage
}
