package memo

import (
	"context"
	"github.com/todo-enjoers/valik/internal/model"
)

type Repo struct{}

func New() *Repo {
	return &Repo{}
}

func (repo *Repo) Store(ctx context.Context, todo model.TodoCreateRequest) (*model.TodoDTO, error) {
	return nil, nil
}

func (repo *Repo) Get(ctx context.Context, id int) (*model.TodoDTO, error) {
	return nil, nil
}
