package postgres

import (
	"context"
	"github.com/todo-enjoers/backend_v1/internal/model"
	"github.com/todo-enjoers/backend_v1/internal/storage"
)

var _ storage.Interface = (*Storage)(nil)

type Storage struct {
}

func (s Storage) Store(_ context.Context, item model.TodoCreateRequest) (*model.Todo, error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) GetAll(_ context.Context) (res []*model.Todo, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) RegisterUser(_ context.Context, item model.UserCreateRequest) (*model.UserDTO, error) {
	//TODO implement me
	panic("implement me")
}
