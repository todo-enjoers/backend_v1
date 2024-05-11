package http

import (
	"context"
	"errors"
	"github.com/todo-enjoers/backend_v1/internal/model"
	"github.com/todo-enjoers/backend_v1/internal/storage/postgres"
	"strings"
)

func ValidateRequest(req model.UserRegisterRequest) (ok bool, err error) {
	err = nil
	if !strings.Contains(req.Login, "@") && len(req.Login) > 7 {
		err = errors.New("email address is required")
		return false, err
	}

	if len(req.Password) < 7 {
		err = errors.New("password is required")
		return false, err
	}

	if postgres.SearchUserByLogin(context.Context, req.Login) {
		return false, errors.New("email address already in use by another user")
	}

	return true, nil
}
