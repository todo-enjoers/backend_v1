package storage

import (
	"context"
	"github.com/todo-enjoers/backend_v1/internal/model"
)

type Interface interface {
	Store(_ context.Context, item model.TodoCreateRequest) (*model.Todo, error)
	GetAll(_ context.Context) (res []*model.Todo, err error)
	RegisterUser(_ context.Context, item model.UserCreateRequest) (*model.UserDTO, error)
}
