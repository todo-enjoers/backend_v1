package storage

import (
	"context"
	"github.com/todo-enjoers/backend_v1/internal/model"
)

type Interface interface {
	InsertUser(ctx context.Context, item model.UserDTO) error
	GetUserByID(ctx context.Context, id int64) (user *model.UserDTO, err error)
	UpdateUserPassword(ctx context.Context, password string, id int64) error
	SearchUserByLogin(ctx context.Context, login string) error
}
