package model

import (
	"github.com/google/uuid"
)

type (
	// UserDTO : User data transfer object
	UserDTO struct {
		ID       uuid.UUID `json:"id"`
		Login    string    `json:"login"`
		Password string    `json:"password"`
	}
	// TodoDTO : Todos data transfer object
	TodoDTO struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		IsCompleted bool      `json:"is_completed"`
		CreatedBy   uuid.UUID `json:"created_by"`
	}
	// GroupDTO : Group data transfer object
	GroupDTO struct {
		UserID    uuid.UUID `json:"user_id"`
		ProjectID uuid.UUID `json:"project_id"`
	}
	// ProjectsDTO : Projects data transfer object
	ProjectsDTO struct {
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		CreatedBy uuid.UUID `json:"created_by"`
	}
)
