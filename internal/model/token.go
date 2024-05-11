package model

import "github.com/google/uuid"

// UserDataInToken : ID and IsAccess are checking ...
type UserDataInToken struct {
	ID       uuid.UUID
	IsAccess bool
}
