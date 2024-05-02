package memo

import (
	"context"
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
