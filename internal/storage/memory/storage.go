package memory

import (
	"context"
	"github.com/todo-enjoers/backend_v1/internal/model"
	"sync"
)

/*
var db *gorm.DB
*/

type Storage struct {
	data       *sync.Map
	lastIndex  int64
	lastIndexU int64
}

func NewStorage() *Storage {
	s := &Storage{
		data:       new(sync.Map),
		lastIndex:  0,
		lastIndexU: 0,
	}
	return s
}

func (s *Storage) Store(_ context.Context, item model.TodoCreateRequest) (*model.Todo, error) {
	todo := &model.Todo{
		ID:          s.lastIndex + 1,
		Name:        item.Name,
		Description: item.Description,
		IsDone:      item.IsDone,
	}
	s.lastIndex += 1
	s.data.Store(todo.ID, todo)
	return todo, nil
}

func (s *Storage) GetAll(_ context.Context) (res []*model.Todo, err error) {
	res = make([]*model.Todo, 0, s.lastIndex)
	s.data.Range(func(_, value any) bool {
		res = append(res, value.(*model.Todo))
		return true
	})
	return res, nil
}

func (s *Storage) RegisterUser(_ context.Context, item model.UserCreateRequest) (model.UserDTO, error) {
	user := model.UserDTO{
		ID:       s.lastIndexU + 1,
		Login:    item.Login,
		Password: item.Password,
	}
	s.lastIndexU += 1
	return user, nil
}

func (s *Storage) LoginUser(_ context.Context, item model.TodoCreateRequest) error {
	// logic with logger
	return nil
}

func (s *Storage) GetMe(_ context.Context) (res []*model.UserDTO, err error) {
	res = make([]*model.UserDTO, 0, s.lastIndexU)
	s.data.Range(func(_, value any) bool {
		res = append(res, value.(*model.UserDTO))
		return true
	})
	return res, nil
}
