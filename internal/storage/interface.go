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
	//GetByUserID(ctx context.Context, userID uuid.UUID) (*model.TodoDTO, error)
	//ChangeTodos(ctx context.Context, todo *model.TodoDTO) error
	//DeleteTodos(ctx context.Context, id uuid.UUID) error
	//ChangeStatus(ctx context.Context, todo *model.TodoDTO) error
}
type Interface interface {
	User() UserStorage
	Todo() TodoStorage
}
