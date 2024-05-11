package model

import "github.com/google/uuid"

type (
	// UserDTO : User Data Transfer Object
	UserDTO struct {
		ID       int64  `json:"id"`
		Login    string `json:"email"`
		Password string `json:"password"`
	}
	// TodoDTO : User Data Transfer Object
	TodoDTO struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		IsCompleted bool      `json:"is_completed"`
		CreatedBy   uuid.UUID `json:"created_by"`
	}
)