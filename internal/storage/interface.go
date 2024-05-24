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
	DeleteTodos(ctx context.Context, id uuid.UUID, createdBy uuid.UUID) error
}

type GroupStorage interface {
	CreateGroup(ctx context.Context, group *model.GroupDTO) error
	GetUsersInProjectByProjectID(ctx context.Context, group *model.GroupDTO) ([]model.GroupDTO, error)
	DeleteFromGroup(ctx context.Context, group *model.GroupDTO) error
}

type ProjectStorage interface {
	CreateProjects(ctx context.Context, project *model.ProjectsDTO) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.ProjectsDTO, error)
	GetMyProjects(ctx context.Context, createdByID uuid.UUID) ([]model.ProjectsDTO, error)
	UpdateName(ctx context.Context, name string, id uuid.UUID) error
}

type Interface interface {
	User() UserStorage
	Todo() TodoStorage
	Group() GroupStorage
	Project() ProjectStorage
}
