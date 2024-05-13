package model

import "github.com/google/uuid"

// UserDataInToken : ID and IsAccess are
type UserDataInToken struct {
	ID       uuid.UUID
	IsAccess bool
}
