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
		ProjectID   uuid.UUID `json:"project_id"`
		CreatedBy   uuid.UUID `json:"created_by"`
		Column      string    `json:"column"`
	}
	// GroupDTO : Group data transfer object
	GroupDTO struct {
		UserID    uuid.UUID `json:"user_id"`
		ProjectID uuid.UUID `json:"project_id"`
	}
	// ProjectDTO : Projects data transfer object
	ProjectDTO struct {
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		CreatedBy uuid.UUID `json:"created_by"`
	}
	// ColumDTO : Column data transfer object
	ColumDTO struct {
		ProjectId uuid.UUID `json:"project_id"`
		Name      string    `json:"name"`
		Order     int       `json:"order"`
	}
)
