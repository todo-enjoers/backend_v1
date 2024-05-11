package token

import (
	"github.com/google/uuid"
	"github.com/todo-enjoers/backend_v1/internal/model"
)

type ProviderI interface {
	GetDataFromToken(token string) (*model.UserDataInToken, error)
	CreateTokenForUser(userID uuid.UUID, isAccess bool) (string, error)
}
